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
		return "", nil
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
