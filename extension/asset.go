package extension

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strings"

	termColor "github.com/fatih/color"
)

func BuildAssetsForExtensions(shopwareRoot string, extensions []Extension) error {
	cfgs := buildAssetConfigFromExtensions(extensions)

	if len(cfgs) == 0 {
		termColor.Yellow("Skipping asset building as all extensions can't be processed")

		return nil
	}

	if !cfgs.RequiresAdminBuild() && cfgs.RequiresStorefrontBuild() {
		termColor.Yellow("Building assets has been skipped as not required")
		return nil
	}

	var err error
	if shopwareRoot == "" {
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

	prepareShopwareForAsset(shopwareRoot, cfgs)

	if cfgs.RequiresAdminBuild() {
		for _, entry := range cfgs {
			// If extension has package.json install it
			if _, err := os.Stat(fmt.Sprintf("%s/Resources/app/administration/src/package.json", entry.BasePath)); err == nil {
				npmInstallCmd := exec.Command("npm", "--prefix", fmt.Sprintf("%s/Resources/app/administration/src/", entry.BasePath), "install") //nolint:gosec
				npmInstallCmd.Stdout = os.Stdout
				npmInstallCmd.Stderr = os.Stderr
				err := npmInstallCmd.Run()

				if err != nil {
					return err
				}
			}
		}

		npmInstallCmd := exec.Command("npm", "--prefix", fmt.Sprintf("%s/src/Administration/Resources/app/administration/", shopwareRoot), "install") //nolint:gosec
		npmInstallCmd.Stdout = os.Stdout
		npmInstallCmd.Stderr = os.Stderr
		err := npmInstallCmd.Run()

		if err != nil {
			return err
		}

		npmBuildCmd := exec.Command("npm", "--prefix", fmt.Sprintf("%s/src/Administration/Resources/app/administration/", shopwareRoot), "run", "build") //nolint:gosec
		npmBuildCmd.Env = []string{fmt.Sprintf("PROJECT_ROOT=%s", shopwareRoot), "SHOPWARE_ADMIN_BUILD_ONLY_EXTENSIONS=1", fmt.Sprintf("PATH=%s", os.Getenv("PATH"))}
		npmBuildCmd.Stdout = os.Stdout
		npmBuildCmd.Stderr = os.Stderr
		err = npmBuildCmd.Run()

		if err != nil {
			return err
		}
	}

	if cfgs.RequiresStorefrontBuild() {
		for _, entry := range cfgs {
			// If extension has package.json install it
			if _, err := os.Stat(fmt.Sprintf("%s/Resources/app/storefront/src/package.json", entry.BasePath)); err == nil {
				npmInstallCmd := exec.Command("npm", "--prefix", fmt.Sprintf("%s/Resources/app/storefront/src/", entry.BasePath), "install") //nolint:gosec
				npmInstallCmd.Stdout = os.Stdout
				npmInstallCmd.Stderr = os.Stderr
				err := npmInstallCmd.Run()

				if err != nil {
					return err
				}
			}
		}

		npmInstallCmd := exec.Command("npm", "--prefix", fmt.Sprintf("%s/src/Storefront/Resources/app/storefront/", shopwareRoot), "install") //nolint:gosec
		npmInstallCmd.Stdout = os.Stdout
		npmInstallCmd.Stderr = os.Stderr
		err := npmInstallCmd.Run()

		if err != nil {
			return err
		}

		npmBuildCmd := exec.Command("npm", "--prefix", fmt.Sprintf("%s/src/Storefront/Resources/app/storefront/", shopwareRoot), "run", "production") //nolint:gosec
		npmBuildCmd.Env = []string{fmt.Sprintf("PROJECT_ROOT=%s", shopwareRoot), fmt.Sprintf("PATH=%s", os.Getenv("PATH"))}
		npmBuildCmd.Stdout = os.Stdout
		npmBuildCmd.Stderr = os.Stderr
		err = npmBuildCmd.Run()

		if err != nil {
			return err
		}
	}

	return nil
}

func prepareShopwareForAsset(shopwareRoot string, cfgs map[string]extensionAssetConfigEntry) {
	varFolder := fmt.Sprintf("%s/var", shopwareRoot)
	if _, err := os.Stat(varFolder); os.IsNotExist(err) {
		err := os.Mkdir(varFolder, os.ModePerm)

		if err != nil {
			log.Fatalln(err)
		}
	}

	pluginJson, err := json.Marshal(cfgs)

	if err != nil {
		log.Fatalln(err)
	}

	err = os.WriteFile(fmt.Sprintf("%s/var/plugins.json", shopwareRoot), pluginJson, os.ModePerm)

	if err != nil {
		log.Fatalln(err)
	}

	err = os.WriteFile(fmt.Sprintf("%s/var/features.json", shopwareRoot), []byte("{}"), os.ModePerm)

	if err != nil {
		log.Fatalln(err)
	}
}

func buildAssetConfigFromExtensions(extensions []Extension) extensionAssetConfig {
	list := make(extensionAssetConfig)

	for _, extension := range extensions {
		extName, err := extension.GetName()

		if err != nil {
			termColor.Red("Skipping extension %s as it has a invalid name", extension.GetPath())
			continue
		}

		extPath := extension.GetPath()

		if _, err := os.Stat(fmt.Sprintf("%s/src/Resources/app", extPath)); os.IsNotExist(err) {
			termColor.Yellow("Skipping extension %s as it doesnt contain assets", extName)
			continue
		}

		var entryFilePathAdmin, entryFilePathStorefront, webpackFileAdmin, webpackFileStorefront *string
		storefrontStyles := make([]string, 0)

		if _, err := os.Stat(fmt.Sprintf("%s/src/Resources/app/administration/src/main.js", extPath)); err == nil {
			val := "Resources/app/administration/src/main.js"
			entryFilePathAdmin = &val
		}

		if _, err := os.Stat(fmt.Sprintf("%s/src/Resources/app/administration/src/main.ts", extPath)); err == nil {
			val := "Resources/app/administration/src/main.ts"
			entryFilePathAdmin = &val
		}

		if _, err := os.Stat(fmt.Sprintf("%s/src/Resources/app/administration/src/build/webpack.config.js", extPath)); err == nil {
			val := "Resources/app/administration/src/build/webpack.config.js"
			webpackFileAdmin = &val
		}

		if _, err := os.Stat(fmt.Sprintf("%s/src/Resources/app/storefront/src/main.js", extPath)); err == nil {
			val := "Resources/app/storefront/src/main.js"
			entryFilePathStorefront = &val
		}

		if _, err := os.Stat(fmt.Sprintf("%s/src/Resources/app/storefront/src/main.ts", extPath)); err == nil {
			val := "Resources/app/storefront/src/main.ts"
			entryFilePathStorefront = &val
		}

		if _, err := os.Stat(fmt.Sprintf("%s/src/Resources/app/storefront/src/build/webpack.config.js", extPath)); err == nil {
			val := "Resources/app/storefront/src/build/webpack.config.js"
			webpackFileStorefront = &val
		}

		if _, err := os.Stat(fmt.Sprintf("%s/Resources/app/storefront/src/scss/base.scss", extPath)); err == nil {
			storefrontStyles = append(storefrontStyles, "Resources/app/storefront/src/scss/base.scss")
		}

		cfg := extensionAssetConfigEntry{
			BasePath: fmt.Sprintf("%s/src/", extPath),
			Views: []string{
				"Resources/views",
			},
			TechnicalName: strings.ReplaceAll(toSnakeCase(extName), "_", "-"),
			Administration: extensionAssetConfigAdmin{
				Path:          "Resources/app/administration/src",
				EntryFilePath: entryFilePathAdmin,
				Webpack:       webpackFileAdmin,
			},
			Storefront: extensionAssetConfigStorefront{
				Path:          "Resources/app/storefront/src",
				EntryFilePath: entryFilePathStorefront,
				Webpack:       webpackFileStorefront,
				StyleFiles:    storefrontStyles,
			},
		}

		list[extName] = cfg
	}

	entryPath := "Resources/app/storefront/src/main.js"
	list["Storefront"] = extensionAssetConfigEntry{
		BasePath:      "src/Storefront/",
		Views:         []string{"Resources/views"},
		TechnicalName: "storefront",
		Storefront: extensionAssetConfigStorefront{
			Path:          "Resources/app/storefront/src",
			EntryFilePath: &entryPath,
			StyleFiles:    []string{},
		},
		Administration: extensionAssetConfigAdmin{
			Path: "Resources/app/storefront/src",
		},
	}

	return list
}

func setupShopwareInTemp() (string, error) {
	dir, err := ioutil.TempDir("", "extension")
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

type extensionAssetConfig map[string]extensionAssetConfigEntry

func (c extensionAssetConfig) RequiresAdminBuild() bool {
	for _, entry := range c {
		if entry.Administration.EntryFilePath != nil {
			return true
		}
	}

	return false
}

func (c extensionAssetConfig) RequiresStorefrontBuild() bool {
	for _, entry := range c {
		if entry.Storefront.EntryFilePath != nil {
			return true
		}
	}

	return false
}

type extensionAssetConfigEntry struct {
	BasePath       string                         `json:"basePath"`
	Views          []string                       `json:"views"`
	TechnicalName  string                         `json:"technicalName"`
	Administration extensionAssetConfigAdmin      `json:"administration"`
	Storefront     extensionAssetConfigStorefront `json:"storefront"`
}

type extensionAssetConfigAdmin struct {
	Path          string  `json:"path"`
	EntryFilePath *string `json:"entryFilePath"`
	Webpack       *string `json:"webpack"`
}

type extensionAssetConfigStorefront struct {
	Path          string   `json:"path"`
	EntryFilePath *string  `json:"entryFilePath"`
	Webpack       *string  `json:"webpack"`
	StyleFiles    []string `json:"styleFiles"`
}
