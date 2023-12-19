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
	CleanupNodeModules           bool
	DisableAdminBuild            bool
	DisableStorefrontBuild       bool
	ShopwareRoot                 string
	ShopwareVersion              *version.Constraints
	Browserslist                 string
	SkipExtensionsWithBuildFiles bool
}

func BuildAssetsForExtensions(ctx context.Context, sources []asset.Source, assetConfig AssetBuildConfig) error { // nolint:gocyclo
	cfgs := BuildAssetConfigFromExtensions(ctx, sources, assetConfig)

	if len(cfgs) == 0 {
		return nil
	}

	if !cfgs.RequiresAdminBuild() && !cfgs.RequiresStorefrontBuild() {
		logging.FromContext(ctx).Infof("Building assets has been skipped as not required")
		return nil
	}

	requiresShopwareSources := cfgs.RequiresShopwareRepository()

	shopwareRoot := assetConfig.ShopwareRoot
	var err error
	if shopwareRoot == "" && requiresShopwareSources {
		shopwareRoot, err = setupShopwareInTemp(ctx, assetConfig.ShopwareVersion)

		if err != nil {
			return err
		}

		defer deletePaths(ctx, shopwareRoot)
	}

	paths, err := InstallNodeModulesOfConfigs(cfgs, true)
	if err != nil {
		return err
	}

	defer deletePaths(ctx, paths...)

	if !assetConfig.DisableAdminBuild && cfgs.RequiresAdminBuild() {
		// Build all extensions compatible with esbuild first
		for name, entry := range cfgs.FilterByAdminAndEsBuild(true) {
			options := esbuild.NewAssetCompileOptionsAdmin(name, entry.BasePath)

			if _, err := esbuild.CompileExtensionAsset(ctx, options); err != nil {
				return err
			}
		}

		nonCompatibleExtensions := cfgs.FilterByAdminAndEsBuild(false)

		if len(nonCompatibleExtensions) != 0 {
			if err := prepareShopwareForAsset(shopwareRoot, nonCompatibleExtensions); err != nil {
				return err
			}

			administrationRoot := PlatformPath(shopwareRoot, "Administration", "Resources/app/administration")
			err := npmRunBuild(
				administrationRoot,
				"build",
				[]string{fmt.Sprintf("PROJECT_ROOT=%s", shopwareRoot), "SHOPWARE_ADMIN_BUILD_ONLY_EXTENSIONS=1", "SHOPWARE_ADMIN_SKIP_SOURCEMAP_GENERATION=1"},
			)

			if assetConfig.CleanupNodeModules {
				defer deletePaths(ctx, path.Join(administrationRoot, "node_modules"), path.Join(administrationRoot, "twigVuePlugin"))
			}

			if err != nil {
				return err
			}
		}
	}

	if !assetConfig.DisableStorefrontBuild && cfgs.RequiresStorefrontBuild() {
		// Build all extensions compatible with esbuild first
		for name, entry := range cfgs.FilterByStorefrontAndEsBuild(true) {
			options := esbuild.NewAssetCompileOptionsStorefront(name, entry.BasePath)

			if _, err := esbuild.CompileExtensionAsset(ctx, options); err != nil {
				return err
			}
		}

		nonCompatibleExtensions := cfgs.FilterByStorefrontAndEsBuild(false)

		if len(nonCompatibleExtensions) != 0 {
			// add the storefront itself as plugin into json
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
			nonCompatibleExtensions["Storefront"] = ExtensionAssetConfigEntry{
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

			if err := prepareShopwareForAsset(shopwareRoot, nonCompatibleExtensions); err != nil {
				return err
			}

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
				defer deletePaths(ctx, path.Join(storefrontRoot, "node_modules"))
			}

			if err != nil {
				return err
			}
		}
	}

	return nil
}

func InstallNodeModulesOfConfigs(cfgs ExtensionAssetConfig, force bool) ([]string, error) {
	paths := make([]string, 0)

	// Install shared node_modules between admin and storefront
	for _, entry := range cfgs {
		possibleNodePaths := []string{
			// shared between admin and storefront
			filepath.Join(entry.BasePath, "Resources", "app", "package.json"),
			// only admin
			filepath.Join(entry.BasePath, "Resources", "app", "administration", "package.json"),
			filepath.Join(entry.BasePath, "Resources", "app", "administration", "src", "package.json"),

			// only storefront
			filepath.Join(entry.BasePath, "Resources", "app", "storefront", "package.json"),
			filepath.Join(entry.BasePath, "Resources", "app", "storefront", "src", "package.json"),
		}

		for _, possibleNodePath := range possibleNodePaths {
			if _, err := os.Stat(possibleNodePath); err == nil {
				npmPath := filepath.Dir(possibleNodePath)

				if _, err := os.Stat(filepath.Join(npmPath, "node_modules")); err == nil && !force {
					continue
				}

				if err := installDependencies(npmPath); err != nil {
					return nil, err
				}

				paths = append(paths, path.Join(npmPath, "node_modules"))
			}
		}
	}

	return paths, nil
}

func deletePaths(ctx context.Context, nodeModulesPaths ...string) {
	for _, nodeModulesPath := range nodeModulesPaths {
		if err := os.RemoveAll(nodeModulesPath); err != nil {
			logging.FromContext(ctx).Errorf("Failed to remove path %s: %s", nodeModulesPath, err.Error())
			return
		}
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

	// Bun can migrate on the fly the package-lock.json to a bun.lockdb and is much faster than NPM
	if _, err := exec.LookPath("bun"); err == nil && canRunBunOnPackage(path) {
		return exec.Command("bun", "install", "--no-save")
	}

	return exec.Command("npm", "install", "--no-audit", "--no-fund", "--prefer-offline")
}

func installDependencies(path string) error {
	installCmd := getInstallCommand(path)
	installCmd.Dir = path
	installCmd.Stdout = os.Stdout
	installCmd.Stderr = os.Stderr
	installCmd.Env = os.Environ()
	installCmd.Env = append(installCmd.Env, "PUPPETEER_SKIP_DOWNLOAD=1", "npm_config_engine_strict=false", "npm_config_fund=false", "npm_config_audit=false", "npm_config_update_notifier=false")

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

func BuildAssetConfigFromExtensions(ctx context.Context, sources []asset.Source, assetCfg AssetBuildConfig) ExtensionAssetConfig {
	list := make(ExtensionAssetConfig)

	for _, source := range sources {
		if source.Name == "" {
			continue
		}

		resourcesDir := path.Join(source.Path, "Resources", "app")

		if _, err := os.Stat(resourcesDir); os.IsNotExist(err) {
			continue
		}

		sourceConfig := createConfigFromPath(source.Name, source.Path)
		sourceConfig.EnableESBuildForAdmin = source.AdminEsbuildCompatible
		sourceConfig.EnableESBuildForStorefront = source.StorefrontEsbuildCompatible

		if assetCfg.SkipExtensionsWithBuildFiles {
			expectedAdminCompiledFile := path.Join(source.Path, "Resources", "public", "administration", "js", esbuild.ToKebabCase(source.Name)+".js")

			if _, err := os.Stat(expectedAdminCompiledFile); err == nil {
				// clear out the entrypoint, so the admin does not build it
				sourceConfig.Administration.EntryFilePath = nil
				sourceConfig.Administration.Webpack = nil

				logging.FromContext(ctx).Infof("Skipping building administration assets for \"%s\" as compiled files are present", source.Name)
			}
		}

		list[source.Name] = sourceConfig
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

	branch := "v" + minVersion

	if minVersion == DevVersionNumber || minVersion == "6.6.0.0" {
		branch = "trunk"
	}

	logging.FromContext(ctx).Infof("Cloning shopware with branch: %s into %s", branch, dir)

	gitCheckoutCmd := exec.Command("git", "clone", "https://github.com/shopware/shopware.git", "--depth=1", "-b", branch, dir)
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

func (c ExtensionAssetConfig) RequiresShopwareRepository() bool {
	for _, entry := range c {
		if entry.Administration.EntryFilePath != nil && !entry.EnableESBuildForAdmin {
			return true
		}

		if entry.Storefront.EntryFilePath != nil && !entry.EnableESBuildForStorefront {
			return true
		}
	}

	return false
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
		if entry.Storefront.EntryFilePath != nil {
			return true
		}
	}

	return false
}

func (c ExtensionAssetConfig) FilterByAdmin() ExtensionAssetConfig {
	filtered := make(ExtensionAssetConfig)

	for name, entry := range c {
		if entry.Administration.EntryFilePath != nil {
			filtered[name] = entry
		}
	}

	return filtered
}

func (c ExtensionAssetConfig) FilterByAdminAndEsBuild(esbuildEnabled bool) ExtensionAssetConfig {
	filtered := make(ExtensionAssetConfig)

	for name, entry := range c {
		if entry.Administration.EntryFilePath != nil && entry.EnableESBuildForAdmin == esbuildEnabled {
			filtered[name] = entry
		}
	}

	return filtered
}

func (c ExtensionAssetConfig) FilterByStorefrontAndEsBuild(esbuildEnabled bool) ExtensionAssetConfig {
	filtered := make(ExtensionAssetConfig)

	for name, entry := range c {
		if entry.Storefront.EntryFilePath != nil && entry.EnableESBuildForStorefront == esbuildEnabled {
			filtered[name] = entry
		}
	}

	return filtered
}

type ExtensionAssetConfigEntry struct {
	BasePath                   string                         `json:"basePath"`
	Views                      []string                       `json:"views"`
	TechnicalName              string                         `json:"technicalName"`
	Administration             ExtensionAssetConfigAdmin      `json:"administration"`
	Storefront                 ExtensionAssetConfigStorefront `json:"storefront"`
	EnableESBuildForAdmin      bool
	EnableESBuildForStorefront bool
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
