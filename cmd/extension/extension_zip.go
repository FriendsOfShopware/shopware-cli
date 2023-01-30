package extension

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/FriendsOfShopware/shopware-cli/extension"

	"github.com/pkg/errors"

	cp "github.com/otiai10/copy"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var disableGit = false
var extensionReleaseMode = false

var extensionZipCmd = &cobra.Command{
	Use:   "zip [path] [branch]",
	Short: "Zip a Extension",
	Args:  cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		path, err := filepath.Abs(args[0])
		if err != nil {
			return err
		}

		var branch string
		if len(args) == 2 {
			branch = args[1]
		}

		ext, err := extension.GetExtensionByFolder(path)
		if err != nil {
			return errors.Wrap(err, "detect extension type")
		}

		extCfg, err := extension.ReadExtensionConfig(ext.GetPath())
		if err != nil {
			log.Fatalln(fmt.Errorf("update: %v", err))
		}

		name, err := ext.GetName()
		if err != nil {
			return errors.Wrap(err, "get name")
		}

		// Clear previous zips
		existingFiles, err := filepath.Glob(fmt.Sprintf("%s-*.zip", name))
		if err != nil {
			return err
		}

		for _, file := range existingFiles {
			err = os.Remove(file)
			if err != nil {
				return errors.Wrap(err, "remove existing file")
			}
		}

		// Create temp dir
		tempDir, err := os.MkdirTemp("", "extension")
		if err != nil {
			return errors.Wrap(err, "create temp directory")
		}

		extName, err := ext.GetName()
		if err != nil {
			return errors.Wrap(err, "get extension name")
		}

		extDir := fmt.Sprintf("%s/%s/", tempDir, extName)

		err = os.Mkdir(extDir, os.ModePerm)
		if err != nil {
			return errors.Wrap(err, "create temp directory")
		}

		tempDir += "/"

		defer func(path string) {
			_ = os.RemoveAll(path)
		}(tempDir)

		var tag string

		// Extract files using strategy
		if disableGit {
			err = cp.Copy(path, extDir, copyOptions())
			if err != nil {
				return errors.Wrap(err, "copy files")
			}
		} else {
			tag, err = extension.GitCopyFolder(path, extDir)
			if err != nil {
				return errors.Wrap(err, "copy via git")
			}
		}

		// User input wins
		if len(branch) > 0 {
			tag = branch
		}

		if extCfg.Build.Zip.Composer.Enabled {
			if err := executeHooks(ext, extCfg.Build.Zip.Composer.BeforeHooks, extDir); err != nil {
				return errors.Wrap(err, "before hooks composer")
			}

			if err := extension.PrepareFolderForZipping(cmd.Context(), extDir, ext, extCfg); err != nil {
				return errors.Wrap(err, "prepare package")
			}

			if err := executeHooks(ext, extCfg.Build.Zip.Composer.AfterHooks, extDir); err != nil {
				return errors.Wrap(err, "after hooks composer")
			}
		}

		if extCfg.Build.Zip.Assets.Enabled {
			if err := executeHooks(ext, extCfg.Build.Zip.Assets.BeforeHooks, extDir); err != nil {
				return errors.Wrap(err, "before hooks assets")
			}

			var tempExt extension.Extension
			if tempExt, err = extension.GetExtensionByFolder(extDir); err != nil {
				return err
			}

			assetBuildConfig := extension.AssetBuildConfig{
				EnableESBuildForAdmin:      extCfg.Build.Zip.Assets.EnableESBuildForAdmin,
				EnableESBuildForStorefront: extCfg.Build.Zip.Assets.EnableESBuildForStorefront,
			}

			if err := extension.BuildAssetsForExtensions(os.Getenv("SHOPWARE_PROJECT_ROOT"), []extension.Extension{tempExt}, assetBuildConfig); err != nil {
				return errors.Wrap(err, "building assets")
			}

			if err := executeHooks(ext, extCfg.Build.Zip.Assets.AfterHooks, extDir); err != nil {
				return errors.Wrap(err, "after hooks assets")
			}
		}

		// Cleanup not wanted files
		if err := extension.CleanupExtensionFolder(extDir, extCfg.Build.Zip.Pack.Excludes.Paths); err != nil {
			return errors.Wrap(err, "cleanup package")
		}

		if extensionReleaseMode {
			if err := extension.PrepareExtensionForRelease(extDir, ext); err != nil {
				return errors.Wrap(err, "prepare for release")
			}
		}

		fileName := fmt.Sprintf("%s-%s.zip", name, tag)
		if len(tag) == 0 {
			fileName = fmt.Sprintf("%s.zip", name)
		}

		if err := executeHooks(ext, extCfg.Build.Zip.Pack.BeforeHooks, extDir); err != nil {
			return errors.Wrap(err, "before hooks pack")
		}

		if err := extension.CreateZip(tempDir, fileName); err != nil {
			return errors.Wrap(err, "create zip file")
		}

		log.Infof("Created file %s", fileName)

		return nil
	},
}

func init() {
	extensionRootCmd.AddCommand(extensionZipCmd)
	extensionZipCmd.Flags().BoolVar(&disableGit, "disable-git", false, "Use the source folder as it is")
	extensionZipCmd.Flags().BoolVar(&extensionReleaseMode, "release", false, "Release mode (remove app secrets)")
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
