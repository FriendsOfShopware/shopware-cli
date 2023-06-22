package extension

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"

	"github.com/FriendsOfShopware/shopware-cli/esbuild"
	"github.com/FriendsOfShopware/shopware-cli/logging"
	"github.com/FriendsOfShopware/shopware-cli/version"
)

const (
	StorefrontWebpackConfig     = "Resources/app/storefront/build/webpack.config.js"
	StorefrontEntrypointJS      = "Resources/app/storefront/src/main.js"
	StorefrontEntrypointTS      = "Resources/app/storefront/src/main.ts"
	StorefrontBaseCSS           = "Resources/app/storefront/src/scss/base.scss"
	AdministrationWebpackConfig = "Resources/app/administration/build/webpack.config.js"
	AdministrationEntrypointJS  = "Resources/app/administration/src/main.js"
	AdministrationEntrypointTS  = "Resources/app/administration/src/main.ts"
)

type AssetBuildConfig struct {
	EnableESBuildForAdmin      bool
	EnableESBuildForStorefront bool
	CleanupNodeModules         bool
	DisableAdminBuild          bool
	DisableStorefrontBuild     bool
}

func BuildAssetsForExtensions(ctx context.Context, shopwareRoot string, extensions []Extension, assetConfig AssetBuildConfig) error { // nolint:gocyclo
	cfgs := buildAssetConfigFromExtensions(ctx, extensions, shopwareRoot)

	if len(cfgs) == 1 {
		return nil
	}

	if !cfgs.RequiresAdminBuild() && !cfgs.RequiresStorefrontBuild() {
		logging.FromContext(ctx).Infof("Building assets has been skipped as not required")
		return nil
	}

	buildWithoutShopwareSource := assetConfig.EnableESBuildForStorefront && assetConfig.EnableESBuildForAdmin

	var err error
	if shopwareRoot == "" && !buildWithoutShopwareSource {
		shopwareRoot, err = setupShopwareInTemp(ctx, extensions[0])

		if err != nil {
			return err
		}

		defer deletePath(ctx, shopwareRoot)
	}

	if !buildWithoutShopwareSource {
		if err := prepareShopwareForAsset(shopwareRoot, cfgs); err != nil {
			return err
		}
	}

	// Install shared node_modules between admin and storefront
	for _, entry := range cfgs {
		// Install also shared node_modules
		if _, err := os.Stat(filepath.Join(entry.BasePath, "Resources", "app", "package.json")); err == nil {
			npmPath := filepath.Join(entry.BasePath, "Resources", "app")
			if err := npmInstall(npmPath); err != nil {
				return err
			}

			if assetConfig.CleanupNodeModules {
				defer deletePath(ctx, path.Join(npmPath, "node_modules"))
			}
		}

		if _, err := os.Stat(filepath.Join(entry.BasePath, "Resources", "app", "administration", "package.json")); err == nil {
			npmPath := filepath.Join(entry.BasePath, "Resources", "app", "administration")
			if err := npmInstall(npmPath); err != nil {
				return err
			}

			if assetConfig.CleanupNodeModules {
				defer deletePath(ctx, path.Join(npmPath, "node_modules"))
			}
		}

		if _, err := os.Stat(filepath.Join(entry.BasePath, "Resources", "app", "storefront", "package.json")); err == nil {
			npmPath := filepath.Join(entry.BasePath, "Resources", "app", "storefront")
			err := npmInstall(npmPath)
			if err != nil {
				return err
			}

			if assetConfig.CleanupNodeModules {
				defer deletePath(ctx, path.Join(npmPath, "node_modules"))
			}
		}
	}

	if !assetConfig.DisableAdminBuild && cfgs.RequiresAdminBuild() {
		if assetConfig.EnableESBuildForAdmin {
			for _, extension := range extensions {
				name, _ := extension.GetName()
				if !cfgs.Has(name) {
					continue
				}

				options := esbuild.NewAssetCompileOptionsAdmin(name, extension.GetPath(), extension.GetType())

				if _, err := esbuild.CompileExtensionAsset(ctx, options); err != nil {
					return err
				}
			}
		} else {
			administrationRoot := PlatformPath(shopwareRoot, "Administration", "Resources/app/administration")
			err := npmInstallAndBuild(
				administrationRoot,
				"build",
				[]string{fmt.Sprintf("PROJECT_ROOT=%s", shopwareRoot), "SHOPWARE_ADMIN_BUILD_ONLY_EXTENSIONS=1"},
			)

			if assetConfig.CleanupNodeModules {
				defer deletePath(ctx, path.Join(administrationRoot, "node_modules"))
				defer deletePath(ctx, path.Join(administrationRoot, "twigVuePlugin"))
			}

			if err != nil {
				return err
			}
		}
	}

	if !assetConfig.DisableStorefrontBuild && cfgs.RequiresStorefrontBuild() {
		if assetConfig.EnableESBuildForStorefront {
			for _, extension := range extensions {
				name, _ := extension.GetName()
				if !cfgs.Has(name) {
					continue
				}

				options := esbuild.NewAssetCompileOptionsStorefront(name, extension.GetPath(), extension.GetType())
				if _, err := esbuild.CompileExtensionAsset(ctx, options); err != nil {
					return err
				}
			}
		} else {
			storefrontRoot := PlatformPath(shopwareRoot, "Storefront", "Resources/app/storefront")
			err := npmInstallAndBuild(
				storefrontRoot,
				"production",
				[]string{fmt.Sprintf("PROJECT_ROOT=%s", shopwareRoot), fmt.Sprintf("STOREFRONT_ROOT=%s", storefrontRoot)},
			)

			if assetConfig.CleanupNodeModules {
				defer deletePath(ctx, path.Join(storefrontRoot, "node_modules"))
			}

			if err != nil {
				return err
			}
		}
	}

	return nil
}

func deletePath(ctx context.Context, path string) {
	if err := os.RemoveAll(path); err != nil {
		logging.FromContext(ctx).Errorf("Failed to remove path %s: %s", path, err.Error())
		return
	}
}

func npmInstallAndBuild(path string, buildCmd string, buildEnvVariables []string) error {
	if err := npmInstall(path); err != nil {
		return err
	}

	npmBuildCmd := exec.Command("npm", "--prefix", path, "run", buildCmd) //nolint:gosec
	npmBuildCmd.Env = append(os.Environ(), buildEnvVariables...)
	npmBuildCmd.Stdout = os.Stdout
	npmBuildCmd.Stderr = os.Stderr

	if err := npmBuildCmd.Run(); err != nil {
		return err
	}

	return nil
}

func npmInstall(path string) error {
	npmInstallCmd := exec.Command("npm", "--prefix", path, "install", "--no-audit", "--no-fund", "--prefer-offline") //nolint:gosec
	npmInstallCmd.Stdout = os.Stdout
	npmInstallCmd.Stderr = os.Stderr
	npmInstallCmd.Env = append(os.Environ(), "PUPPETEER_SKIP_DOWNLOAD=1")

	if err := npmInstallCmd.Run(); err != nil {
		return err
	}

	return nil
}

func prepareShopwareForAsset(shopwareRoot string, cfgs map[string]ExtensionAssetConfigEntry) error {
	varFolder := fmt.Sprintf("%s/var", shopwareRoot)
	if _, err := os.Stat(varFolder); os.IsNotExist(err) {
		err := os.Mkdir(varFolder, os.ModePerm)
		if err != nil {
			return fmt.Errorf("prepareShopwareForAsset: %w", err)
		}
	}

	pluginJson, err := json.Marshal(cfgs)
	if err != nil {
		return fmt.Errorf("prepareShopwareForAsset: %w", err)
	}

	err = os.WriteFile(fmt.Sprintf("%s/var/plugins.json", shopwareRoot), pluginJson, os.ModePerm)

	if err != nil {
		return fmt.Errorf("prepareShopwareForAsset: %w", err)
	}

	err = os.WriteFile(fmt.Sprintf("%s/var/features.json", shopwareRoot), []byte("{}"), os.ModePerm)

	if err != nil {
		return fmt.Errorf("prepareShopwareForAsset: %w", err)
	}

	return nil
}

func buildAssetConfigFromExtensions(ctx context.Context, extensions []Extension, shopwareRoot string) ExtensionAssetConfig {
	list := make(ExtensionAssetConfig)

	for _, extension := range extensions {
		extName, err := extension.GetName()
		if err != nil {
			logging.FromContext(ctx).Errorf("Skipping extension %s as it has a invalid name", extension.GetPath())
			continue
		}

		extPath := extension.GetPath()

		if _, err := os.Stat(path.Join(extension.GetRootDir(), "Resources")); os.IsNotExist(err) {
			logging.FromContext(ctx).Infof("Skipping building of assets for extension %s as it doesnt contain assets", extName)
			continue
		}

		list[extName] = createConfigFromPath(extName, extension.GetRootDir())

		extCfg, err := readExtensionConfig(extPath)
		if err != nil {
			logging.FromContext(ctx).Errorf("Skipping extension additional bundles %s as it has a invalid config", extPath)
			continue
		}

		for _, bundle := range extCfg.Build.ExtraBundles {
			bundleName := bundle.Name

			if bundleName == "" {
				bundleName = filepath.Base(bundle.Path)
			}

			list[bundleName] = createConfigFromPath(bundleName, path.Join(extension.GetRootDir(), bundle.Path))
		}
	}

	var basePath string
	if shopwareRoot == "" {
		basePath = "src/Storefront/"
	} else {
		basePath = strings.TrimLeft(
			strings.Replace(PlatformPath(shopwareRoot, "Storefront", ""), shopwareRoot, "", 1),
			"/",
		) + "/"
	}

	entryPath := "Resources/app/storefront/src/main.js"
	list["Storefront"] = ExtensionAssetConfigEntry{
		BasePath:      basePath,
		Views:         []string{"Resources/views"},
		TechnicalName: "storefront",
		Storefront: ExtensionAssetConfigStorefront{
			Path:          "Resources/app/storefront/src",
			EntryFilePath: &entryPath,
			StyleFiles:    []string{},
		},
		Administration: ExtensionAssetConfigAdmin{
			Path: "Resources/app/administration/src",
		},
	}

	return list
}

func createConfigFromPath(entryPointName string, extensionRoot string) ExtensionAssetConfigEntry {
	var entryFilePathAdmin, entryFilePathStorefront, webpackFileAdmin, webpackFileStorefront *string
	storefrontStyles := make([]string, 0)

	if _, err := os.Stat(path.Join(extensionRoot, AdministrationEntrypointJS)); err == nil {
		val := AdministrationEntrypointJS
		entryFilePathAdmin = &val
	}

	if _, err := os.Stat(path.Join(extensionRoot, AdministrationEntrypointTS)); err == nil {
		val := AdministrationEntrypointTS
		entryFilePathAdmin = &val
	}

	if _, err := os.Stat(path.Join(extensionRoot, AdministrationWebpackConfig)); err == nil {
		val := AdministrationWebpackConfig
		webpackFileAdmin = &val
	}

	if _, err := os.Stat(path.Join(extensionRoot, StorefrontEntrypointJS)); err == nil {
		val := StorefrontEntrypointJS
		entryFilePathStorefront = &val
	}

	if _, err := os.Stat(path.Join(extensionRoot, StorefrontEntrypointTS)); err == nil {
		val := StorefrontEntrypointTS
		entryFilePathStorefront = &val
	}

	if _, err := os.Stat(path.Join(extensionRoot, StorefrontWebpackConfig)); err == nil {
		val := StorefrontWebpackConfig
		webpackFileStorefront = &val
	}

	if _, err := os.Stat(path.Join(extensionRoot, StorefrontBaseCSS)); err == nil {
		storefrontStyles = append(storefrontStyles, StorefrontBaseCSS)
	}

	extensionRoot = strings.TrimRight(extensionRoot, "/") + "/"

	cfg := ExtensionAssetConfigEntry{
		BasePath: extensionRoot,
		Views: []string{
			"Resources/views",
		},
		TechnicalName: strings.ReplaceAll(ToSnakeCase(entryPointName), "_", "-"),
		Administration: ExtensionAssetConfigAdmin{
			Path:          "Resources/app/administration/src",
			EntryFilePath: entryFilePathAdmin,
			Webpack:       webpackFileAdmin,
		},
		Storefront: ExtensionAssetConfigStorefront{
			Path:          "Resources/app/storefront/src",
			EntryFilePath: entryFilePathStorefront,
			Webpack:       webpackFileStorefront,
			StyleFiles:    storefrontStyles,
		},
	}
	return cfg
}

func setupShopwareInTemp(ctx context.Context, ext Extension) (string, error) {
	minVersion, err := lookupForMinMatchingVersion(ctx, ext)
	if err != nil {
		return "", err
	}

	dir, err := os.MkdirTemp("", "extension")
	if err != nil {
		return "", err
	}

	cloneBranch := "6.4"

	shopware65Constraint, _ := version.NewConstraint("~6.5.0")

	if shopware65Constraint.Check(version.Must(version.NewVersion(minVersion))) {
		cloneBranch = "trunk"
	}

	logging.FromContext(ctx).Infof("Cloning shopware with branch: %s into %s", cloneBranch, dir)

	gitCheckoutCmd := exec.Command("git", "clone", "https://github.com/shopware/platform.git", "--depth=1", "-b", cloneBranch, dir)
	gitCheckoutCmd.Stdout = os.Stdout
	gitCheckoutCmd.Stderr = os.Stderr
	err = gitCheckoutCmd.Run()

	if err != nil {
		return "", err
	}

	return dir, nil
}

type ExtensionAssetConfig map[string]ExtensionAssetConfigEntry

func (c ExtensionAssetConfig) Has(name string) bool {
	_, ok := c[name]

	return ok
}

func (c ExtensionAssetConfig) RequiresAdminBuild() bool {
	for _, entry := range c {
		if entry.Administration.EntryFilePath != nil {
			return true
		}
	}

	return false
}

func (c ExtensionAssetConfig) RequiresStorefrontBuild() bool {
	for _, entry := range c {
		if entry.TechnicalName == "storefront" {
			continue
		}

		if entry.Storefront.EntryFilePath != nil {
			return true
		}
	}

	return false
}

type ExtensionAssetConfigEntry struct {
	BasePath       string                         `json:"basePath"`
	Views          []string                       `json:"views"`
	TechnicalName  string                         `json:"technicalName"`
	Administration ExtensionAssetConfigAdmin      `json:"administration"`
	Storefront     ExtensionAssetConfigStorefront `json:"storefront"`
}

type ExtensionAssetConfigAdmin struct {
	Path          string  `json:"path"`
	EntryFilePath *string `json:"entryFilePath"`
	Webpack       *string `json:"webpack"`
}

type ExtensionAssetConfigStorefront struct {
	Path          string   `json:"path"`
	EntryFilePath *string  `json:"entryFilePath"`
	Webpack       *string  `json:"webpack"`
	StyleFiles    []string `json:"styleFiles"`
}
