package cmd

import (
	"encoding/json"
	"fmt"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"shopware-cli/extension"
	"strings"
)

func findClosestShopwareProject() (string, error) {
	currentDir, err := os.Getwd()
	if err != nil {
		return "", err
	}

	for {
		files := []string{
			fmt.Sprintf("%s/composer.json", currentDir),
			fmt.Sprintf("%s/composer.lock", currentDir),
		}

		for _, file := range files {
			if _, err := os.Stat(file); err == nil {
				content, err := ioutil.ReadFile(file)

				if err != nil {
					return "", err
				}
				contentString := string(content)

				if strings.Contains(contentString, "shopware/core") {
					return currentDir, nil
				}
			}
		}

		currentDir = filepath.Dir(currentDir)

		if currentDir == filepath.Dir(currentDir) {
			break
		}
	}

	return "", fmt.Errorf("cannot find Shopware project in current directory")
}

func getPlatformPath(projectRoot, component, path string) string {
	if _, err := os.Stat(projectRoot + "/src/Core/composer.json"); err == nil {
		return fmt.Sprintf(projectRoot+"/src/%s/%s", component, path)
	} else if _, err := os.Stat(projectRoot + "/vendor/shopware/platform/"); err == nil {
		return fmt.Sprintf(projectRoot+"/vendor/shopware/platform/src/%s/%s", component, path)
	}

	return fmt.Sprintf(projectRoot+"/vendor/shopware/%s/%s", strings.ToLower(component), path)
}

func runSimpleCommand(root string, app string, args ...string) error {
	cmd := exec.Command(app, args...)
	cmd.Dir = root
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}

func buildStorefront(projectRoot string, forceNpmInstall bool) error {
	storefrontRoot := getPlatformPath(projectRoot, "Storefront", "Resources/app/storefront")

	if err := runSimpleCommand(projectRoot, "php", "bin/console", "bundle:dump"); err != nil {
		return err
	}

	if err := setupExtensionNodeModules(projectRoot, forceNpmInstall); err != nil {
		return err
	}

	// Optional command, allowed to failure
	_ = runSimpleCommand(projectRoot, "php", "bin/console", "feature:dump")

	// Optional npm install
	_, err := os.Stat(getPlatformPath(projectRoot, "Storefront", "Resources/app/storefront/node_modules"))

	if forceNpmInstall || os.IsNotExist(err) {
		if installErr := runSimpleCommand(projectRoot, "npm", "install", "--prefix", storefrontRoot, "--no-save"); err != nil {
			return installErr
		}
	}

	if err := runSimpleCommand(projectRoot, "node", getPlatformPath(projectRoot, "Storefront", "Resources/app/storefront/copy-to-vendor.js")); err != nil {
		return err
	}

	envs := []string{
		fmt.Sprintf("PATH=%s", os.Getenv("PATH")),
		fmt.Sprintf("PROJECT_ROOT=%s", projectRoot),
	}

	npmRun := exec.Command("npm", "--prefix", storefrontRoot, "run", "production")
	npmRun.Env = envs
	npmRun.Stdin = os.Stdin
	npmRun.Stdout = os.Stdout
	npmRun.Stderr = os.Stderr

	if err := npmRun.Run(); err != nil {
		return err
	}

	if err := runSimpleCommand(projectRoot, "php", "bin/console", "assets:install"); err != nil {
		return err
	}

	return runSimpleCommand(projectRoot, "php", "bin/console", "theme:compile")
}

func buildAdministration(projectRoot string, forceNpmInstall bool) error {
	adminRoot := getPlatformPath(projectRoot, "Administration", "Resources/app/administration")

	if err := runSimpleCommand(projectRoot, "php", "bin/console", "bundle:dump"); err != nil {
		return err
	}

	if err := setupExtensionNodeModules(projectRoot, forceNpmInstall); err != nil {
		return err
	}

	// Optional command, allowed to failure
	_ = runSimpleCommand(projectRoot, "php", "bin/console", "feature:dump")

	// Optional npm install

	_, err := os.Stat(getPlatformPath(projectRoot, "Administration", "Resources/app/administration/node_modules"))

	if forceNpmInstall || os.IsNotExist(err) {
		if installErr := runSimpleCommand(projectRoot, "npm", "install", "--prefix", adminRoot, "--no-save"); err != nil {
			return installErr
		}
	}

	envs := []string{
		fmt.Sprintf("PATH=%s", os.Getenv("PATH")),
		fmt.Sprintf("PROJECT_ROOT=%s", projectRoot),
	}

	npmRun := exec.Command("npm", "--prefix", adminRoot, "run", "build")
	npmRun.Env = envs
	npmRun.Stdin = os.Stdin
	npmRun.Stdout = os.Stdout
	npmRun.Stderr = os.Stderr

	if err := npmRun.Run(); err != nil {
		return err
	}

	return runSimpleCommand(projectRoot, "php", "bin/console", "theme:compile")
}

func setupExtensionNodeModules(projectRoot string, forceNpmInstall bool) error {
	// Skip if plugins.json is missing
	if _, err := os.Stat(projectRoot + "/var/plugins.json"); os.IsNotExist(err) {
		log.Infof("Cannot find a var/plugins.json")
		return nil
	}

	var assetCfg extension.ExtensionAssetConfig
	var content []byte
	var err error

	if content, err = ioutil.ReadFile(projectRoot + "/var/plugins.json"); err != nil {
		return err
	}

	if err := json.Unmarshal(content, &assetCfg); err != nil {
		return err
	}

	for _, ext := range assetCfg {
		_, adminPathPackage := os.Stat(fmt.Sprintf("%s/%s/%s/package.json", projectRoot, ext.BasePath, filepath.Dir(ext.Administration.Path)))
		_, adminPathNodeModules := os.Stat(fmt.Sprintf("%s/%s/%s/node_modules", projectRoot, ext.BasePath, filepath.Dir(ext.Administration.Path)))

		_, storefrontPathPackage := os.Stat(fmt.Sprintf("%s/%s/%s/package.json", projectRoot, ext.BasePath, filepath.Dir(ext.Storefront.Path)))
		_, storefrontPathNodeModules := os.Stat(fmt.Sprintf("%s/%s/%s/node_modules", projectRoot, ext.BasePath, filepath.Dir(ext.Storefront.Path)))

		if ext.Administration.EntryFilePath != nil && adminPathPackage == nil && os.IsNotExist(adminPathNodeModules) || forceNpmInstall {
			if err := runSimpleCommand(projectRoot, "npm", "install", "--prefix", fmt.Sprintf("%s/%s/%s", projectRoot, ext.BasePath, filepath.Dir(ext.Administration.Path))); err != nil {
				return err
			}
		}

		if ext.Storefront.EntryFilePath != nil && storefrontPathPackage == nil && os.IsNotExist(storefrontPathNodeModules) || forceNpmInstall {
			if err := runSimpleCommand(projectRoot, "npm", "install", "--prefix", fmt.Sprintf("%s/%s/%s", projectRoot, ext.BasePath, filepath.Dir(ext.Storefront.Path))); err != nil {
				fmt.Println(ext.TechnicalName)
				fmt.Println(ext.Storefront.EntryFilePath)

				return err
			}
		}
	}

	return nil
}
