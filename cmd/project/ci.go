package project

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"

	"dario.cat/mergo"
	"github.com/FriendsOfShopware/shopware-cli/extension"
	"github.com/FriendsOfShopware/shopware-cli/logging"
	"github.com/FriendsOfShopware/shopware-cli/shop"
	"github.com/spf13/cobra"
)

// cleanupPaths are paths that are not nesscarry for the production build.
var cleanupPaths = []string{
	"vendor/shopware/storefront/Resources/app/storefront/vendor/bootstrap/dist",
	"vendor/shopware/storefront/Resources/app/storefront/test",
	"vendor/shopware/storefront/Test",
	"vendor/shopware/core/Framework/Test",
	"vendor/shopware/core/Content/Test",
	"vendor/shopware/core/Checkout/Test",
	"vendor/shopware/core/System/Test",
	"vendor/tecnickcom/tcpdf/examples",
}

var projectCI = &cobra.Command{
	Use:   "ci",
	Short: "Build Shopware in the CI",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		var err error
		args[0], err = filepath.Abs(args[0])
		if err != nil {
			return err
		}

		if os.Getenv("APP_ENV") == "" {
			if err := os.Setenv("APP_ENV", "prod"); err != nil {
				return err
			}
		}

		shopCfg, err := shop.ReadConfig(filepath.Join(args[0], ".shopware-project.yml"), true)
		if err != nil {
			return err
		}

		cleanupPaths = append(cleanupPaths, shopCfg.Build.CleanupPaths...)

		logging.FromContext(cmd.Context()).Infof("Installing dependencies using Composer")

		composer := exec.CommandContext(cmd.Context(), "composer", "install", "--no-dev", "--no-interaction", "--no-progress", "--optimize-autoloader", "--classmap-authoritative")
		composer.Dir = args[0]
		composer.Stdin = os.Stdin
		composer.Stdout = os.Stdout
		composer.Stderr = os.Stderr

		if err := composer.Run(); err != nil {
			return err
		}

		logging.FromContext(cmd.Context()).Infof("Looking for extensions to build assets in project")

		sources := extension.FindAssetSourcesOfProject(cmd.Context(), args[0])
		constraint, err := extension.GetShopwareProjectConstraint(args[0])
		if err != nil {
			return err
		}

		assetCfg := extension.AssetBuildConfig{
			EnableESBuildForAdmin:      false,
			EnableESBuildForStorefront: false,
			CleanupNodeModules:         true,
			ShopwareRoot:               args[0],
			ShopwareVersion:            constraint,
		}

		if err := extension.BuildAssetsForExtensions(cmd.Context(), sources, assetCfg); err != nil {
			return err
		}

		logging.FromContext(cmd.Context()).Infof("Optimizing Administration sources")
		if err := cleanupAdministrationFiles(cmd.Context(), path.Join(args[0], "vendor", "shopware", "administration")); err != nil {
			return err
		}

		if !shopCfg.Build.KeepExtensionSource {
			for _, source := range sources {
				if err := cleanupAdministrationFiles(cmd.Context(), source.Path); err != nil {
					return err
				}
			}
		}

		for _, removePath := range cleanupPaths {
			logging.FromContext(cmd.Context()).Infof("Removing %s", removePath)

			if err := os.RemoveAll(path.Join(args[0], removePath)); err != nil {
				return err
			}
		}

		logging.FromContext(cmd.Context()).Infof("Remove unnecessary fonts from tcpdf")

		if err := cleanupTcpdf(args[0]); err != nil {
			return err
		}

		logging.FromContext(cmd.Context()).Infof("Warmup container cache")

		if err := runTransparentCommand(exec.CommandContext(cmd.Context(), "php", path.Join(args[0], "bin", "ci"), "--version")); err != nil { //nolint: gosec
			return fmt.Errorf("failed to warmup container cache (php bin/ci --version): %w", err)
		}

		if !shopCfg.Build.DisableAssetCopy {
			logging.FromContext(cmd.Context()).Infof("Copying extension assets to final public/bundles folder")

			// Delete asset manifest to force a new build
			manifestPath := path.Join(args[0], "public", "asset-manifest.json")
			if _, err := os.Stat(manifestPath); err == nil {
				if err := os.Remove(manifestPath); err != nil {
					return err
				}
			}

			if err := runTransparentCommand(exec.CommandContext(cmd.Context(), "php", path.Join(args[0], "bin", "ci"), "asset:install")); err != nil { //nolint: gosec
				return fmt.Errorf("failed to install assets (php bin/ci asset:install): %w", err)
			}
		}

		if shopCfg.Build.RemoveExtensionAssets {
			logging.FromContext(cmd.Context()).Infof("Deleting assets of extensions")

			for _, source := range sources {
				if _, err := os.Stat(path.Join(source.Path, "Resources", "public", "administration", "css")); err == nil {
					if err := os.WriteFile(path.Join(source.Path, "Resources", ".administration-css"), []byte{}, os.ModePerm); err != nil {
						return err
					}
				}

				if _, err := os.Stat(path.Join(source.Path, "Resources", "public", "administration", "js")); err == nil {
					if err := os.WriteFile(path.Join(source.Path, "Resources", ".administration-js"), []byte{}, os.ModePerm); err != nil {
						return err
					}
				}

				if err := os.RemoveAll(path.Join(source.Path, "Resources", "public")); err != nil {
					return err
				}
			}

			if err := os.RemoveAll(path.Join(args[0], "vendor", "shopware", "administration", "Resources", "public")); err != nil {
				return err
			}

			if err := os.WriteFile(path.Join(args[0], "vendor", "shopware", "administration", "Resources", ".administration-js"), []byte{}, os.ModePerm); err != nil {
				return err
			}

			if err := os.WriteFile(path.Join(args[0], "vendor", "shopware", "administration", "Resources", ".administration-css"), []byte{}, os.ModePerm); err != nil {
				return err
			}
		}

		return nil
	},
}

func init() {
	projectRootCmd.AddCommand(projectCI)
}

func commandWithRoot(cmd *exec.Cmd, root string) *exec.Cmd {
	cmd.Dir = root

	return cmd
}

func runTransparentCommand(cmd *exec.Cmd) error {
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Env = append(os.Environ(), "APP_SECRET=test", "LOCK_DSN=flock")

	return cmd.Run()
}

func cleanupTcpdf(folder string) error {
	return filepath.WalkDir(path.Join(folder, "vendor", "tecnickcom/tcpdf/fonts"), func(path string, d os.DirEntry, err error) error {
		if d.IsDir() {
			return nil
		}

		if filepath.Base(path) == ".z" {
			return os.Remove(path)
		}

		baseName := filepath.Base(path)

		if strings.Contains(baseName, "courier") || strings.Contains(baseName, "helvetica") {
			return nil
		}

		return os.Remove(path)
	})
}

func cleanupAdministrationFiles(ctx context.Context, folder string) error {
	adminFolder := path.Join(folder, "Resources", "app", "administration")

	if _, err := os.Stat(adminFolder); err == nil {
		logging.FromContext(ctx).Infof("Merging Administration snippet for %s", folder)

		snippetFiles := make(map[string][]string)

		err = filepath.WalkDir(adminFolder, func(path string, d os.DirEntry, err error) error {
			if d.IsDir() {
				return nil
			}

			if filepath.Ext(path) != ".json" {
				return nil
			}

			if filepath.Base(filepath.Dir(path)) != "snippet" {
				return nil
			}

			name := filepath.Base(path)
			extension := filepath.Ext(name)
			language := name[0 : len(name)-len(extension)]

			if _, ok := snippetFiles[language]; !ok {
				snippetFiles[language] = []string{}
			}

			snippetFiles[language] = append(snippetFiles[language], path)

			return nil
		})

		if err != nil {
			return err
		}

		for language, files := range snippetFiles {
			if len(files) == 1 {
				data, err := os.ReadFile(files[0])
				if err != nil {
					return err
				}

				if err := os.WriteFile(path.Join(folder, language), data, os.ModePerm); err != nil {
					return err
				}

				continue
			}

			merged := make(map[string]interface{})

			for _, file := range files {
				snippetFile := make(map[string]interface{})

				data, err := os.ReadFile(file)
				if err != nil {
					return err
				}

				if err := json.Unmarshal(data, &snippetFile); err != nil {
					return err
				}

				if err := mergo.Merge(&merged, snippetFile, mergo.WithOverride); err != nil {
					return err
				}

				if err != nil {
					return err
				}
			}

			mergedData, err := json.Marshal(merged)
			if err != nil {
				return err
			}

			if err := os.WriteFile(path.Join(folder, language), mergedData, os.ModePerm); err != nil {
				return err
			}
		}

		logging.FromContext(ctx).Infof("Deleting Administration source files for %s", folder)

		if err := os.RemoveAll(adminFolder); err != nil {
			return err
		}

		logging.FromContext(ctx).Infof("Migrating generated snippet file for %s", folder)

		snippetFolder := path.Join(adminFolder, "src", "app", "snippet")
		if err := os.MkdirAll(snippetFolder, os.ModePerm); err != nil {
			return err
		}

		for language := range snippetFiles {
			if err := os.Rename(path.Join(folder, language), path.Join(snippetFolder, language+".json")); err != nil {
				return err
			}
		}

		logging.FromContext(ctx).Infof("Creating empty main.js for %s", folder)
		return os.WriteFile(path.Join(adminFolder, "src", "main.js"), []byte(""), os.ModePerm)
	}

	return nil
}
