package extension

import (
	"fmt"
	"os"
	"os/exec"
	"path"
	"path/filepath"

	cp "github.com/otiai10/copy"
	"github.com/spf13/cobra"

	"github.com/FriendsOfShopware/shopware-cli/extension"
	"github.com/FriendsOfShopware/shopware-cli/logging"
)

var (
	disableGit           = false
	extensionReleaseMode = false
)

var extensionZipCmd = &cobra.Command{
	Use:   "zip [path] [branch]",
	Short: "Zip a Extension",
	Args:  cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		extPath, err := filepath.Abs(args[0])
		if err != nil {
			return err
		}

		var branch string
		if len(args) == 2 {
			branch = args[1]
		}

		ext, err := extension.GetExtensionByFolder(extPath)
		if err != nil {
			return fmt.Errorf("detect extension type: %w", err)
		}

		extCfg := ext.GetExtensionConfig()

		name, err := ext.GetName()
		if err != nil {
			return fmt.Errorf("get name: %w", err)
		}

		// Clear previous zips
		existingFiles, err := filepath.Glob(fmt.Sprintf("%s-*.zip", name))
		if err != nil {
			return err
		}

		for _, file := range existingFiles {
			err = os.Remove(file)
			if err != nil {
				return fmt.Errorf("remove existing file: %w", err)
			}
		}

		// Create temp dir
		tempDir, err := os.MkdirTemp("", "extension")
		if err != nil {
			return fmt.Errorf("create temp directory: %w", err)
		}

		extName, err := ext.GetName()
		if err != nil {
			return fmt.Errorf("get extension name: %w", err)
		}

		extDir := fmt.Sprintf("%s/%s/", tempDir, extName)

		err = os.Mkdir(extDir, os.ModePerm)
		if err != nil {
			return fmt.Errorf("create temp directory: %w", err)
		}

		tempDir += "/"

		defer func(path string) {
			_ = os.RemoveAll(path)
		}(tempDir)

		var tag string

		// Extract files using strategy
		if disableGit {
			err = cp.Copy(extPath, extDir, copyOptions())
			if err != nil {
				return fmt.Errorf("copy files: %w", err)
			}
		} else {
			gitCommit, _ := cmd.Flags().GetString("git-commit")

			tag, err = extension.GitCopyFolder(extPath, extDir, gitCommit)
			if err != nil {
				return fmt.Errorf("copy via git: %w", err)
			}

			logging.FromContext(cmd.Context()).Infof("Checking out %s using Git", tag)
		}

		// User input wins
		if len(branch) > 0 {
			tag = branch
		}

		if extCfg.Build.Zip.Composer.Enabled {
			if err := executeHooks(ext, extCfg.Build.Zip.Composer.BeforeHooks, extDir); err != nil {
				return fmt.Errorf("before hooks composer: %w", err)
			}

			if err := extension.PrepareFolderForZipping(cmd.Context(), extDir, ext, extCfg); err != nil {
				return fmt.Errorf("prepare package: %w", err)
			}

			if err := executeHooks(ext, extCfg.Build.Zip.Composer.AfterHooks, extDir); err != nil {
				return fmt.Errorf("after hooks composer: %w", err)
			}
		}

		if extCfg.Build.Zip.Assets.Enabled {
			if err := executeHooks(ext, extCfg.Build.Zip.Assets.BeforeHooks, extDir); err != nil {
				return fmt.Errorf("before hooks assets: %w", err)
			}

			var tempExt extension.Extension
			if tempExt, err = extension.GetExtensionByFolder(extDir); err != nil {
				return err
			}

			shopwareConstraint, err := tempExt.GetShopwareVersionConstraint()
			if err != nil {
				return fmt.Errorf("get shopware version constraint: %w", err)
			}

			assetBuildConfig := extension.AssetBuildConfig{
				CleanupNodeModules: true,
				ShopwareRoot:       os.Getenv("SHOPWARE_PROJECT_ROOT"),
				ShopwareVersion:    shopwareConstraint,
			}

			if err := extension.BuildAssetsForExtensions(cmd.Context(), extension.ConvertExtensionsToSources(cmd.Context(), []extension.Extension{tempExt}), assetBuildConfig); err != nil {
				return fmt.Errorf("building assets: %w", err)
			}

			if err := executeHooks(ext, extCfg.Build.Zip.Assets.AfterHooks, extDir); err != nil {
				return fmt.Errorf("after hooks assets: %w", err)
			}
		}

		// Cleanup not wanted files
		if err := extension.CleanupExtensionFolder(extDir, extCfg.Build.Zip.Pack.Excludes.Paths); err != nil {
			return fmt.Errorf("cleanup package: %w", err)
		}

		if extensionReleaseMode {
			if err := extension.PrepareExtensionForRelease(cmd.Context(), extPath, extDir, ext); err != nil {
				return fmt.Errorf("prepare for release: %w", err)
			}
		}

		fileName := fmt.Sprintf("%s-%s.zip", name, tag)
		if len(tag) == 0 {
			fileName = fmt.Sprintf("%s.zip", name)
		}

		outputDir, _ := cmd.Flags().GetString("output-directory")

		if len(outputDir) > 0 {
			if _, err := os.Stat(outputDir); os.IsNotExist(err) {
				if err := os.MkdirAll(outputDir, os.ModePerm); err != nil {
					return fmt.Errorf("create output directory: %w", err)
				}
			}

			fileName = path.Join(outputDir, fileName)
		}

		if err := executeHooks(ext, extCfg.Build.Zip.Pack.BeforeHooks, extDir); err != nil {
			return fmt.Errorf("before hooks pack: %w", err)
		}

		if err := extension.CreateZip(tempDir, fileName); err != nil {
			return fmt.Errorf("create zip file: %w", err)
		}

		logging.FromContext(cmd.Context()).Infof("Created file %s", fileName)

		return nil
	},
}

func init() {
	extensionRootCmd.AddCommand(extensionZipCmd)
	extensionZipCmd.Flags().BoolVar(&disableGit, "disable-git", false, "Use the source folder as it is")
	extensionZipCmd.Flags().BoolVar(&extensionReleaseMode, "release", false, "Release mode (remove app secrets)")
	extensionZipCmd.Flags().String("output-directory", "", "Output directory for the zip file")
	extensionZipCmd.Flags().String("git-commit", "", "Commit Hash / Tag to use")
}

func executeHooks(ext extension.Extension, hooks []string, extDir string) error {
	env := []string{
		fmt.Sprintf("EXTENSION_DIR=%s", extDir),
		fmt.Sprintf("ORIGINAL_EXTENSION_DIR=%s", ext.GetPath()),
	}

	for _, hook := range hooks {
		hookCmd := exec.Command("sh", "-c", hook)
		hookCmd.Stdout = os.Stdout
		hookCmd.Stderr = os.Stderr
		hookCmd.Dir = extDir
		hookCmd.Env = append(os.Environ(), env...)
		err := hookCmd.Run()
		if err != nil {
			return err
		}
	}

	return nil
}

func copyOptions() cp.Options {
	return cp.Options{
		OnSymlink: func(string) cp.SymlinkAction {
			return cp.Skip
		},
	}
}
