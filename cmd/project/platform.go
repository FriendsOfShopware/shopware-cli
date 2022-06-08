package project

import (
	"encoding/json"
	"fmt"
	"github.com/FriendsOfShopware/shopware-cli/extension"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
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
					if _, err := os.Stat(fmt.Sprintf("%s/bin/console", currentDir)); err == nil {
						return currentDir, nil
					}
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

func runConsoleCommand(projectRoot string, command string) error {
	return runSimpleCommand(projectRoot, "php", "bin/console", command)
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
	storefrontRoot := extension.PlatformPath(projectRoot, "Storefront", "Resources/app/storefront")

	if err := runConsoleCommand(projectRoot, "bundle:dump"); err != nil {
		return err
	}

	if err := setupExtensionNodeModules(projectRoot, forceNpmInstall); err != nil {
		return err
	}

	// Optional command, allowed to fail
	_ = runConsoleCommand(projectRoot, "feature:dump")

	// Optional npm install
	_, err := os.Stat(extension.PlatformPath(projectRoot, "Storefront", "Resources/app/storefront/node_modules"))

	if forceNpmInstall || os.IsNotExist(err) {
		if installErr := runSimpleCommand(projectRoot, "npm", "install", "--prefix", storefrontRoot, "--no-save"); err != nil {
			return installErr
		}
	}

	if err := runSimpleCommand(projectRoot, "node", extension.PlatformPath(projectRoot, "Storefront", "Resources/app/storefront/copy-to-vendor.js")); err != nil {
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

	if err := runConsoleCommand(projectRoot, "assets:install"); err != nil {
		return err
	}

	return runConsoleCommand(projectRoot, "theme:compile")
}

func buildAdministration(projectRoot string, forceNpmInstall bool) error {
	adminRoot := extension.PlatformPath(projectRoot, "Administration", "Resources/app/administration")

	if err := runConsoleCommand(projectRoot, "bundle:dump"); err != nil {
		return err
	}

	if err := setupExtensionNodeModules(projectRoot, forceNpmInstall); err != nil {
		return err
	}

	// Optional command, allowed to fail
	_ = runConsoleCommand(projectRoot, "feature:dump")

	// Optional npm install

	_, err := os.Stat(extension.PlatformPath(projectRoot, "Administration", "Resources/app/administration/node_modules"))

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

	return runConsoleCommand(projectRoot, "assets:install")
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

		if ext.Administration.EntryFilePath != nil && adminPathPackage == nil && (os.IsNotExist(adminPathNodeModules) || forceNpmInstall) {
			if err := runSimpleCommand(projectRoot, "npm", "install", "--prefix", fmt.Sprintf("%s/%s/%s", projectRoot, ext.BasePath, filepath.Dir(ext.Administration.Path)), "--no-save"); err != nil {
				return err
			}
		}

		if ext.Storefront.EntryFilePath != nil && storefrontPathPackage == nil && (os.IsNotExist(storefrontPathNodeModules) || forceNpmInstall) {
			if err := runSimpleCommand(projectRoot, "npm", "install", "--prefix", fmt.Sprintf("%s/%s/%s", projectRoot, ext.BasePath, filepath.Dir(ext.Storefront.Path)), "--no-save"); err != nil {
				return err
			}
		}
	}

	return nil
}
