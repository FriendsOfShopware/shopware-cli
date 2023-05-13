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

	"github.com/FriendsOfShopware/shopware-cli/extension"
	"github.com/FriendsOfShopware/shopware-cli/logging"
	"github.com/imdario/mergo"
	"github.com/spf13/cobra"
)

// cleanupPaths are paths that are not nesscarry for the production build.
var cleanupPaths = []string{
	"vendor/shopware/storefront/Resources/app/storefront/vendor/bootstrap/dist",
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
	RunE: func(cmd *cobra.Command, args []string) error {
		logging.FromContext(cmd.Context()).Infof("Installing dependencies using Composer")

		composer := exec.CommandContext(cmd.Context(), "composer", "install", "--no-dev", "--no-interaction", "--no-progress", "--optimize-autoloader", "--classmap-authoritative")
		composer.Dir = args[0]
		composer.Stdout = os.Stdout
		composer.Stderr = os.Stderr

		if err := composer.Run(); err != nil {
			return err
		}

		logging.FromContext(cmd.Context()).Infof("Looking for extensions to build assets in project")

		extensions := findExtensionsFromProject(cmd.Context(), args[0])

		assetCfg := extension.AssetBuildConfig{EnableESBuildForAdmin: false, EnableESBuildForStorefront: false, CleanupNodeModules: true}

		err := extension.BuildAssetsForExtensions(cmd.Context(), args[0], extensions, assetCfg)
		if err != nil {
			return err
		}

		logging.FromContext(cmd.Context()).Infof("Optimizing Administration sources")
		if err := cleanupAdministrationFiles(cmd.Context(), path.Join(args[0], "vendor", "shopware", "administration")); err != nil {
			return err
		}

		for _, ext := range extensions {
			if err := cleanupAdministrationFiles(cmd.Context(), ext.GetRootDir()); err != nil {
				return err
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

		if err := exec.CommandContext(cmd.Context(), "php", path.Join(args[0], "bin", "ci"), "--version").Run(); err != nil { //nolint: gosec
			return fmt.Errorf("failed to warmup container cache (php bin/ci --version): %w", err)
		}

		return nil
	},
}

func init() {
	projectRootCmd.AddCommand(projectCI)
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

		snippetFolder := path.Join(adminFolder, "src", "snippet")
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

func findExtensionsFromProject(ctx context.Context, project string) []extension.Extension {
	extensions := make(map[string]extension.Extension)

	for _, ext := range addExtensionsByComposer(project) {
		name, err := ext.GetName()
		if err != nil {
			continue
		}

		version, _ := ext.GetVersion()

		logging.FromContext(ctx).Infof("Found extension using Composer: %s (%s)", name, version)

		extensions[name] = ext
	}

	for _, ext := range addExtensionsByWildcard(path.Join(project, "custom", "plugins")) {
		name, err := ext.GetName()
		if err != nil {
			continue
		}

		// Skip if extension is already added by composer
		if _, ok := extensions[name]; ok {
			continue
		}

		version, _ := ext.GetVersion()

		logging.FromContext(ctx).Infof("Found extension in custom/plugins: %s (%s)", name, version)
		logging.FromContext(ctx).Errorf("Extension %s should be installed using Composer. Please remove the extension from custom/plugins.", name)

		extensions[name] = ext
	}

	for _, ext := range addExtensionsByWildcard(path.Join(project, "custom", "apps")) {
		name, err := ext.GetName()
		if err != nil {
			continue
		}
		version, _ := ext.GetVersion()

		logging.FromContext(ctx).Infof("Found extension in custom/apps: %s (%s)", name, version)

		extensions[name] = ext
	}

	extensionsSlice := make([]extension.Extension, 0, len(extensions))

	for _, ext := range extensions {
		extensionsSlice = append(extensionsSlice, ext)
	}

	return extensionsSlice
}

func addExtensionsByComposer(project string) []extension.Extension {
	var list []extension.Extension

	lock, err := os.ReadFile(path.Join(project, "composer.lock"))
	if err != nil {
		return list
	}

	var composer composerLock
	err = json.Unmarshal(lock, &composer)

	if err != nil {
		return list
	}

	for _, pkg := range composer.Packages {
		if pkg.PackageType == "shopware-platform-plugin" {
			ext, err := extension.GetExtensionByFolder(path.Join(project, "vendor", pkg.Name))
			if err != nil {
				continue
			}

			list = append(list, ext)
		}
	}

	return list
}

func addExtensionsByWildcard(extensionDir string) []extension.Extension {
	var list []extension.Extension

	extensions, err := os.ReadDir(extensionDir)
	if err != nil {
		return list
	}

	for _, file := range extensions {
		if file.IsDir() {
			ext, err := extension.GetExtensionByFolder(path.Join(extensionDir, file.Name()))
			if err != nil {
				continue
			}

			list = append(list, ext)
		}
	}

	return list
}

type composerLock struct {
	Packages []struct {
		Name        string `json:"name"`
		Version     string `json:"version"`
		PackageType string `json:"type"`
	} `json:"packages"`
}
