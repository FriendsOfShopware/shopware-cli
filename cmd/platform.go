package cmd

import (
	"fmt"
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

func getPlatformPath(component, path string) string {
	if _, err := os.Stat("src/Core/composer.json"); err == nil {
		return fmt.Sprintf("src/%s/%s", component, path)
	} else if _, err := os.Stat("vendor/shopware/platform/"); err == nil {
		return fmt.Sprintf("vendor/shopware/platform/src/%s/%s", component, path)
	}

	return fmt.Sprintf("vendor/shopware/%s/%s", strings.ToLower(component), path)
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
	adminRoot := getPlatformPath("Storefront", "Resources/app/storefront")

	if err := runSimpleCommand(projectRoot, "php", "bin/console", "bundle:dump"); err != nil {
		return err
	}

	// Optional command, allowed to failure
	_ = runSimpleCommand(projectRoot, "php", "bin/console", "feature:dump")

	// Optional npm install
	_, err := os.Stat(getPlatformPath("Storefront", "Resources/app/storefront/node_modules"))

	if forceNpmInstall || os.IsNotExist(err) {
		if installErr := runSimpleCommand(projectRoot, "npm", "install", "--prefix", adminRoot); err != nil {
			return installErr
		}
	}

	if err := runSimpleCommand(projectRoot, "node", getPlatformPath("Storefront", "Resources/app/storefront/copy-to-vendor.js")); err != nil {
		return err
	}

	envs := []string{
		fmt.Sprintf("PATH=%s", os.Getenv("PATH")),
		fmt.Sprintf("PROJECT_ROOT=%s", projectRoot),
	}

	npmRun := exec.Command("npm", "--prefix", adminRoot, "run", "production")
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
	adminRoot := getPlatformPath("Administration", "Resources/app/administration")

	if err := runSimpleCommand(projectRoot, "php", "bin/console", "bundle:dump"); err != nil {
		return err
	}

	// Optional command, allowed to failure
	_ = runSimpleCommand(projectRoot, "php", "bin/console", "feature:dump")

	// Optional npm install

	_, err := os.Stat(getPlatformPath("Administration", "Resources/app/administration/node_modules"))

	if forceNpmInstall || os.IsNotExist(err) {
		if installErr := runSimpleCommand(projectRoot, "npm", "install", "--prefix", adminRoot); err != nil {
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
