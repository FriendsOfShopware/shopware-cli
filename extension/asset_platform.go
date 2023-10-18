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

	"github.com/FriendsOfShopware/shopware-cli/internal/asset"
	"github.com/FriendsOfShopware/shopware-cli/internal/esbuild"
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
	ShopwareRoot               string
	ShopwareVersion            *version.Constraints
	Browserslist               string
}

func BuildAssetsForExtensions(ctx context.Context, sources []asset.Source, assetConfig AssetBuildConfig) error { // nolint:gocyclo
	cfgs := buildAssetConfigFromExtensions(sources, assetConfig.ShopwareRoot)

	if len(cfgs) == 1 {
		return nil
	}

	if !cfgs.RequiresAdminBuild() && !cfgs.RequiresStorefrontBuild() {
		logging.FromContext(ctx).Infof("Building assets has been skipped as not required")
		return nil
	}

	buildWithoutShopwareSource := assetConfig.EnableESBuildForStorefront && assetConfig.EnableESBuildForAdmin

	shopwareRoot := assetConfig.ShopwareRoot
	var err error
	if shopwareRoot == "" && !buildWithoutShopwareSource {
		shopwareRoot, err = setupShopwareInTemp(ctx, assetConfig.ShopwareVersion)

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
			if err := installDependencies(npmPath); err != nil {
				return err
			}

			if assetConfig.CleanupNodeModules {
				defer deletePath(ctx, path.Join(npmPath, "node_modules"))
			}
		}

		if _, err := os.Stat(filepath.Join(entry.BasePath, "Resources", "app", "administration", "package.json")); err == nil {
			npmPath := filepath.Join(entry.BasePath, "Resources", "app", "administration")
			if err := installDependencies(npmPath); err != nil {
				return err
			}

			if assetConfig.CleanupNodeModules {
				defer deletePath(ctx, path.Join(npmPath, "node_modules"))
			}
		}

		if _, err := os.Stat(filepath.Join(entry.BasePath, "Resources", "app", "storefront", "package.json")); err == nil {
			npmPath := filepath.Join(entry.BasePath, "Resources", "app", "storefront")
			err := installDependencies(npmPath)
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
			for _, source := range sources {
				if !cfgs.Has(source.Name) {
					continue
				}

				options := esbuild.NewAssetCompileOptionsAdmin(source.Name, source.Path)

				if _, err := esbuild.CompileExtensionAsset(ctx, options); err != nil {
					return err
				}
			}
		} else {
			administrationRoot := PlatformPath(shopwareRoot, "Administration", "Resources/app/administration")
			err := npmRunBuild(
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
			for _, source := range sources {
				if !cfgs.Has(source.Name) {
					continue
				}

				options := esbuild.NewAssetCompileOptionsStorefront(source.Name, source.Path)
				if _, err := esbuild.CompileExtensionAsset(ctx, options); err != nil {
					return err
				}
			}
		} else {
			storefrontRoot := PlatformPath(shopwareRoot, "Storefront", "Resources/app/storefront")

			envList := []string{
				fmt.Sprintf("PROJECT_ROOT=%s", shopwareRoot),
				fmt.Sprintf("STOREFRONT_ROOT=%s", storefrontRoot),
			}

			if assetConfig.Browserslist != "" {
				npx := exec.CommandContext(ctx, "npx", "--yes", "update-browserslist-db", "--quiet")
				npx.Stdout = os.Stdout
				npx.Stderr = os.Stderr
				npx.Dir = storefrontRoot

				if err := npx.Run(); err != nil {
					return err
				}

				envList = append(envList, fmt.Sprintf("BROWSERSLIST=%s", assetConfig.Browserslist))
			}

			err := npmRunBuild(
				storefrontRoot,
				"production",
				envList,
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

func npmRunBuild(path string, buildCmd string, buildEnvVariables []string) error {
	if err := installDependencies(path); err != nil {
		return err
	}

	npmBuildCmd := exec.Command("npm", "--prefix", path, "run", buildCmd) //nolint:gosec
	npmBuildCmd.Env = os.Environ()
	npmBuildCmd.Env = append(npmBuildCmd.Env, buildEnvVariables...)
	npmBuildCmd.Stdout = os.Stdout
	npmBuildCmd.Stderr = os.Stderr

	if err := npmBuildCmd.Run(); err != nil {
		return err
	}

	return nil
}

func getInstallCommand(path string) *exec.Cmd {
	if _, err := os.Stat(filepath.Join(path, "pnpm-lock.yaml")); err == nil {
		return exec.Command("pnpm", "install")
	}

	if _, err := os.Stat(filepath.Join(path, "yarn.lock")); err == nil {
		return exec.Command("yarn", "install")
	}

	if _, err := os.Stat(filepath.Join(path, "bun.lockdb")); err == nil {
		return exec.Command("bun", "install")
	}

	return exec.Command("npm", "install", "--no-audit", "--no-fund", "--prefer-offline")
}

func installDependencies(path string) error {
	installCmd := getInstallCommand(path)
	installCmd.Dir = path
	installCmd.Stdout = os.Stdout
	installCmd.Stderr = os.Stderr
	installCmd.Env = os.Environ()
	installCmd.Env = append(installCmd.Env, "PUPPETEER_SKIP_DOWNLOAD=1")

	if err := installCmd.Run(); err != nil {
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

func buildAssetConfigFromExtensions(sources []asset.Source, shopwareRoot string) ExtensionAssetConfig {
	list := make(ExtensionAssetConfig)

	for _, source := range sources {
		if source.Name == "" {
			continue
		}

		resourcesDir := path.Join(source.Path, "Resources", "app")

		if _, err := os.Stat(resourcesDir); os.IsNotExist(err) {
			continue
		}

		list[source.Name] = createConfigFromPath(source.Name, source.Path)
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
		TechnicalName: esbuild.ToKebabCase(entryPointName),
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

func setupShopwareInTemp(ctx context.Context, shopwareVersionConstraint *version.Constraints) (string, error) {
	minVersion, err := lookupForMinMatchingVersion(ctx, shopwareVersionConstraint)
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
