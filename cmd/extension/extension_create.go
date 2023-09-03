package extension

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"regexp"

	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"

	"github.com/FriendsOfShopware/shopware-cli/internal/config"
)

var extensionCreateCmd = &cobra.Command{
	Use:   "create [name]",
	Short: "Create an extension boilerplate",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		rootPath, err := filepath.Abs(".")
		if err != nil {
			return err
		}

		var extensionConfig config.ExtensionConfig = config.ExtensionConfig{
			Name:            args[0],
			Namespace:       args[0],
			ShopwareVersion: "^6.4",
			License:         "MIT",
		}

		namespace, _ := cmd.PersistentFlags().GetString("namespace")
		createInCustomPlugins, _ := cmd.PersistentFlags().GetBool("create-in-custom-plugins")

		if namespace != "" {
			extensionConfig.Namespace = namespace
		}

		fmt.Printf("Using namespace: %s\n", extensionConfig.Namespace)

		var pluginPath string
		if createInCustomPlugins {
			pluginPath = fmt.Sprintf("%s/custom/plugins/%s", rootPath, extensionConfig.Name)
		} else {
			pluginPath = fmt.Sprintf("%s/%s", rootPath, extensionConfig.Name)
		}

		if _, err := os.Stat(pluginPath); err == nil {
			return fmt.Errorf("the directory '%s' already exists", pluginPath)
		}

		extensionConfig.ComposerPackage = askExtension(promptui.Prompt{
			Label:    "Composer package",
			Validate: validComposerPackage,
		})

		extensionConfig.ShopwareVersion = askExtension(promptui.Prompt{
			Label:    "Required shopware/core version",
			Validate: emptyValidator,
			Default:  extensionConfig.ShopwareVersion,
		})

		extensionConfig.Label = askExtension(promptui.Prompt{
			Label:    "Plugin label",
			Validate: emptyValidator,
		})

		extensionConfig.Description = askExtension(promptui.Prompt{
			Label:    "Plugin description",
			Validate: emptyValidator,
		})

		extensionConfig.License = askExtension(promptui.Prompt{
			Label:    "License",
			Validate: emptyValidator,
			Default:  extensionConfig.License,
		})

		extensionConfig.ManufacturerLink = askExtension(promptui.Prompt{
			Label:    "Manufacturer link",
			Validate: emptyValidator,
		})

		extensionConfig.SupportLink = askExtension(promptui.Prompt{
			Label:    "Support link",
			Validate: emptyValidator,
			Default:  extensionConfig.ManufacturerLink,
		})

		err = os.MkdirAll(fmt.Sprintf("%s/src/Resources/config", pluginPath), 0o755)
		if err != nil {
			return err
		}

		err = createComposerJson(fmt.Sprintf("%s/composer.json", pluginPath), extensionConfig)
		if err != nil {
			return err
		}

		err = makeBootstrap(pluginPath, extensionConfig)
		if err != nil {
			return err
		}

		err = makeChangelog(pluginPath)
		if err != nil {
			return err
		}

		err = makeDefaultServices(pluginPath)
		if err != nil {
			return err
		}

		return nil
	},
}

func init() {
	extensionRootCmd.AddCommand(extensionCreateCmd)
	extensionCreateCmd.PersistentFlags().String("namespace", "", "Namespace for the plugin")
	extensionCreateCmd.PersistentFlags().Bool("create-in-custom-plugins", true, "Create the plugin in custom/plugins directory")
}

func askExtension(inputPrompt promptui.Prompt) string {
	if id, err := inputPrompt.Run(); err != nil {
		panic(err)
	} else {
		return id
	}
}

func emptyValidator(s string) error {
	if len(s) == 0 {
		return errors.New("this cannot be empty")
	}
	return nil
}

func validComposerPackage(s string) error {
	validComposerPackageRegExp := regexp.MustCompile("^[a-z0-9]([_.-]?[a-z0-9]+)*/[a-z0-9](([_.]?|-{0,2})[a-z0-9]+)*$")

	if !validComposerPackageRegExp.MatchString(s) {
		return fmt.Errorf("'%s' is not a valid composer package", s)
	}
	return nil
}

func createComposerJson(composerFile string, extensionConfig config.ExtensionConfig) error {
	composerData := composerStruct{
		Name:        extensionConfig.ComposerPackage,
		Version:     "1.0.0",
		Description: extensionConfig.Description,
		Type:        "shopware-platform-plugin",
		License:     extensionConfig.License,
		Autoload: map[string]any{
			"psr-4": map[string]string{
				fmt.Sprintf("%s\\", extensionConfig.Namespace): "src/",
			},
		},
		Require: map[string]string{
			"shopware/core": extensionConfig.ShopwareVersion,
		},
		Extra: composerExtra{
			ShopwarePluginClass: fmt.Sprintf("%s\\%s", extensionConfig.Namespace, extensionConfig.Name),
			Label: map[string]string{
				"en-GB": extensionConfig.Label,
			},
			Description: map[string]string{
				"en-GB": extensionConfig.Description,
			},
			ManufacturerLink: map[string]string{
				"en-GB": extensionConfig.ManufacturerLink,
			},
			SupportLink: map[string]string{
				"en-GB": extensionConfig.SupportLink,
			},
		},
	}

	jsonContent, err := json.MarshalIndent(composerData, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(composerFile, jsonContent, 0o600)
}

func makeDefaultServices(pluginPath string) error {
	xml := `<?xml version="1.0" ?>
<container xmlns="http://symfony.com/schema/dic/services"
           xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance"
           xsi:schemaLocation="http://symfony.com/schema/dic/services http://symfony.com/schema/dic/services/services-1.0.xsd">
    <services>

    </services>
</container>
`

	servicesFilename := fmt.Sprintf("%s/src/Resources/config/services.xml", pluginPath)
	return os.WriteFile(servicesFilename, []byte(xml), 0o600)
}

func makeChangelog(pluginPath string) error {
	changelogContent := `# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

## [1.0.0]

### Added

- Initial version
`

	changelogFilename := fmt.Sprintf("%s/CHANGELOG_en-GB.md", pluginPath)
	return os.WriteFile(changelogFilename, []byte(changelogContent), 0o600)
}

func makeBootstrap(pluginPath string, extensionConfig config.ExtensionConfig) error {
	fileContentTemplate := `<?php declare(strict_types=1);

namespace %s;

use Shopware\Core\Framework\Plugin;

class %s extends Plugin
{
}
`
	bootStrapFilename := fmt.Sprintf("%s/src/%s.php", pluginPath, extensionConfig.Name)
	fileContent := fmt.Sprintf(fileContentTemplate, extensionConfig.Namespace, extensionConfig.Name)
	return os.WriteFile(bootStrapFilename, []byte(fileContent), 0o600)
}

type composerExtra struct {
	ShopwarePluginClass string            `json:"shopware-plugin-class"`
	Label               map[string]string `json:"label"`
	Description         map[string]string `json:"description"`
	ManufacturerLink    map[string]string `json:"manufacturerLink"`
	SupportLink         map[string]string `json:"supportLink"`
}

type composerStruct struct {
	Name        string            `json:"name"`
	Version     string            `json:"version"`
	Description string            `json:"description"`
	Type        string            `json:"type"`
	License     string            `json:"license"`
	Autoload    map[string]any    `json:"autoload"`
	Require     map[string]string `json:"require"`
	Extra       composerExtra     `json:"extra"`
}
