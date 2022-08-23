package extension

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

type AssetBuildConfig struct {
	EnableESBuildForAdmin      bool
	EnableESBuildForStorefront bool
}

func BuildAssetsForExtensions(shopwareRoot string, extensions []Extension, assetConfig AssetBuildConfig) error {
	cfgs := buildAssetConfigFromExtensions(extensions, shopwareRoot)

	if len(cfgs) == 1 {
		return nil
	}

	if !cfgs.RequiresAdminBuild() && !cfgs.RequiresStorefrontBuild() {
		log.Infof("Building assets has been skipped as not required")
		return nil
	}

	buildWithoutShopwareSource := assetConfig.EnableESBuildForStorefront && assetConfig.EnableESBuildForAdmin

	var err error
	if shopwareRoot == "" && !buildWithoutShopwareSource {
		shopwareRoot, err = setupShopwareInTemp()

		if err != nil {
			return err
		}

		defer func(path string) {
			err := os.RemoveAll(path)
			if err != nil {
				log.Println(err)
			}
		}(shopwareRoot)
	}

	if !buildWithoutShopwareSource {
		if err := prepareShopwareForAsset(shopwareRoot, cfgs); err != nil {
			return err
		}
	}

	if cfgs.RequiresAdminBuild() {
		for _, entry := range cfgs {
			// If extension has package.json install it
			if _, err := os.Stat(filepath.Join(entry.BasePath, "Resources/app/administration/package.json")); err == nil {
				err := npmInstall(filepath.Join(entry.BasePath, "Resources/app/administration/"))

				if err != nil {
					return err
				}
			}
		}

		if assetConfig.EnableESBuildForAdmin {
			for _, extension := range extensions {
				name, _ := extension.GetName()
				if !cfgs.Has(name) {
					continue
				}

				options := NewAssetCompileOptionsAdmin()
				options.ProductionMode = true
				if _, err := CompileExtensionAsset(extension, options); err != nil {
					return err
				}
			}
		} else {
			administrationRoot := PlatformPath(shopwareRoot, "Administration", "Resources/app/administration")
			err := npmInstallAndBuild(
				administrationRoot,
				"build",
				[]string{fmt.Sprintf("PROJECT_ROOT=%s", shopwareRoot), fmt.Sprintf("PATH=%s", os.Getenv("PATH")), "SHOPWARE_ADMIN_BUILD_ONLY_EXTENSIONS=1"},
			)

			if err != nil {
				return err
			}
		}
	}

	if cfgs.RequiresStorefrontBuild() {
		if assetConfig.EnableESBuildForStorefront {
			for _, extension := range extensions {
				name, _ := extension.GetName()
				if !cfgs.Has(name) {
					continue
				}

				options := NewAssetCompileOptionsStorefront()
				options.ProductionMode = true
				if _, err := CompileExtensionAsset(extension, options); err != nil {
					return err
				}
			}
		} else {
			for _, entry := range cfgs {
				// If extension has package.json install it
				if _, err := os.Stat(filepath.Join(entry.BasePath, "Resources/app/storefront/package.json")); err == nil {
					err := npmInstall(filepath.Join(entry.BasePath, "Resources/app/storefront/"))

					if err != nil {
						return err
					}
				}
			}

			storefrontRoot := PlatformPath(shopwareRoot, "Storefront", "Resources/app/storefront")
			err := npmInstallAndBuild(
				storefrontRoot,
				"production",
				[]string{fmt.Sprintf("PROJECT_ROOT=%s", shopwareRoot), fmt.Sprintf("PATH=%s", os.Getenv("PATH")), fmt.Sprintf("STOREFRONT_ROOT=%s", storefrontRoot)},
			)

			if err != nil {
				return err
			}
		}
	}

	return nil
}

func npmInstallAndBuild(path string, buildCmd string, buildEnvVariables []string) error {
	if err := npmInstall(path); err != nil {
		return err
	}

	npmBuildCmd := exec.Command("npm", "--prefix", path, "run", buildCmd) //nolint:gosec
	npmBuildCmd.Env = buildEnvVariables
	npmBuildCmd.Stdout = os.Stdout
	npmBuildCmd.Stderr = os.Stderr

	if err := npmBuildCmd.Run(); err != nil {
		return err
	}

	return nil
}

func npmInstall(path string) error {
	npmInstallCmd := exec.Command("npm", "--prefix", path, "install") //nolint:gosec
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
			return errors.Wrap(err, "prepareShopwareForAsset")
		}
	}

	pluginJson, err := json.Marshal(cfgs)

	if err != nil {
		return errors.Wrap(err, "prepareShopwareForAsset")
	}

	err = os.WriteFile(fmt.Sprintf("%s/var/plugins.json", shopwareRoot), pluginJson, os.ModePerm)

	if err != nil {
		return errors.Wrap(err, "prepareShopwareForAsset")
	}

	err = os.WriteFile(fmt.Sprintf("%s/var/features.json", shopwareRoot), []byte("{}"), os.ModePerm)

	if err != nil {
		return errors.Wrap(err, "prepareShopwareForAsset")
	}

	return nil
}

func buildAssetConfigFromExtensions(extensions []Extension, shopwareRoot string) ExtensionAssetConfig {
	list := make(ExtensionAssetConfig)

	for _, extension := range extensions {
		extName, err := extension.GetName()

		if err != nil {
			log.Errorf("Skipping extension %s as it has a invalid name", extension.GetPath())
			continue
		}

		pathPrefix := "src/Resources"
		extensionRoot := "src/"
		if extension.GetType() == TypePlatformApp {
			pathPrefix = "Resources"
			extensionRoot = ""
		}

		extPath := extension.GetPath()

		if _, err := os.Stat(fmt.Sprintf("%s/%s/app", extPath, pathPrefix)); os.IsNotExist(err) {
			log.Infof("Skipping building of assets for extension %s as it doesnt contain assets", extName)
			continue
		}

		var entryFilePathAdmin, entryFilePathStorefront, webpackFileAdmin, webpackFileStorefront *string
		storefrontStyles := make([]string, 0)

		if _, err := os.Stat(fmt.Sprintf("%s/%s/app/administration/src/main.js", extPath, pathPrefix)); err == nil {
			val := "Resources/app/administration/src/main.js"
			entryFilePathAdmin = &val
		}

		if _, err := os.Stat(fmt.Sprintf("%s/%s/app/administration/src/main.ts", extPath, pathPrefix)); err == nil {
			val := "Resources/app/administration/src/main.ts"
			entryFilePathAdmin = &val
		}

		if _, err := os.Stat(fmt.Sprintf("%s/%s/app/administration/build/webpack.config.js", extPath, pathPrefix)); err == nil {
			val := "Resources/app/administration/build/webpack.config.js"
			webpackFileAdmin = &val
		}

		if _, err := os.Stat(fmt.Sprintf("%s/%s/app/storefront/src/main.js", extPath, pathPrefix)); err == nil {
			val := "Resources/app/storefront/src/main.js"
			entryFilePathStorefront = &val
		}

		if _, err := os.Stat(fmt.Sprintf("%s/%s/app/storefront/src/main.ts", extPath, pathPrefix)); err == nil {
			val := "Resources/app/storefront/src/main.ts"
			entryFilePathStorefront = &val
		}

		if _, err := os.Stat(fmt.Sprintf("%s/%s/app/storefront/build/webpack.config.js", extPath, pathPrefix)); err == nil {
			val := "Resources/app/storefront/build/webpack.config.js"
			webpackFileStorefront = &val
		}

		if _, err := os.Stat(fmt.Sprintf("%s/%s/app/storefront/src/scss/base.scss", extPath, pathPrefix)); err == nil {
			storefrontStyles = append(storefrontStyles, "Resources/app/storefront/src/scss/base.scss")
		}

		cfg := ExtensionAssetConfigEntry{
			BasePath: fmt.Sprintf("%s/%s", extPath, extensionRoot),
			Views: []string{
				"Resources/views",
			},
			TechnicalName: strings.ReplaceAll(ToSnakeCase(extName), "_", "-"),
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

		list[extName] = cfg
	}

	var basePath string
	if shopwareRoot == "" {
		basePath = "src/Storefront"
	} else {
		basePath = strings.TrimLeft(
			strings.Replace(PlatformPath(shopwareRoot, "Storefront", ""), shopwareRoot, "", 1),
			"/",
		)
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

func setupShopwareInTemp() (string, error) {
	dir, err := os.MkdirTemp("", "extension")
	if err != nil {
		return "", err
	}

	gitCheckoutCmd := exec.Command("git", "clone", "https://github.com/shopware/platform.git", "--depth=1", dir)
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
