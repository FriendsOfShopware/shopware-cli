package extension

import (
	"encoding/json"
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/FriendsOfShopware/shopware-cli/version"
)

type PlatformPlugin struct {
	path     string
	composer platformComposerJson
}

func (p PlatformPlugin) GetRootDir() string {
	return path.Join(p.path, "src")
}

func (p PlatformPlugin) GetResourcesDir() string {
	return path.Join(p.GetRootDir(), "Resources")
}

func newPlatformPlugin(path string) (*PlatformPlugin, error) {
	composerJsonFile := fmt.Sprintf("%s/composer.json", path)
	if _, err := os.Stat(composerJsonFile); err != nil {
		return nil, err
	}

	jsonFile, err := os.ReadFile(composerJsonFile)

	if err != nil {
		return nil, fmt.Errorf("newPlatformPlugin: %v", err)
	}

	var composerJson platformComposerJson
	err = json.Unmarshal(jsonFile, &composerJson)

	if err != nil {
		return nil, fmt.Errorf("newPlatformPlugin: %v", err)
	}

	extension := PlatformPlugin{
		composer: composerJson,
		path:     path,
	}

	return &extension, nil
}

type platformComposerJson struct {
	Name        string   `json:"name"`
	Keywords    []string `json:"keywords"`
	Description string   `json:"description"`
	Version     string   `json:"version"`
	Type        string   `json:"type"`
	License     string   `json:"license"`
	Authors     []struct {
		Name     string `json:"name"`
		Homepage string `json:"homepage"`
	} `json:"authors"`
	Require map[string]string `json:"require"`
	Extra   struct {
		ShopwarePluginClass string            `json:"shopware-plugin-class"`
		Label               map[string]string `json:"label"`
		Description         map[string]string `json:"description"`
		ManufacturerLink    map[string]string `json:"manufacturerLink"`
		SupportLink         map[string]string `json:"supportLink"`
	} `json:"extra"`
	Autoload struct {
		Psr0 map[string]string `json:"psr-0"`
		Psr4 map[string]string `json:"psr-4"`
	} `json:"autoload"`
}

func (p PlatformPlugin) GetName() (string, error) {
	if p.composer.Extra.ShopwarePluginClass == "" {
		return "", fmt.Errorf("extension name is empty")
	}

	parts := strings.Split(p.composer.Extra.ShopwarePluginClass, "\\")

	return parts[len(parts)-1], nil
}

func (p PlatformPlugin) GetShopwareVersionConstraint() (*version.Constraints, error) {
	shopwareConstraintString, ok := p.composer.Require["shopware/core"]

	if !ok {
		return nil, fmt.Errorf("require.shopware/core is required")
	}

	shopwareConstraint, err := version.NewConstraint(shopwareConstraintString)

	if err != nil {
		return nil, err
	}

	return &shopwareConstraint, err
}

func (PlatformPlugin) GetType() string {
	return TypePlatformPlugin
}

func (p PlatformPlugin) GetVersion() (*version.Version, error) {
	return version.NewVersion(p.composer.Version)
}

func (p PlatformPlugin) GetChangelog() (*extensionTranslated, error) {
	return parseExtensionMarkdownChangelog(p)
}

func (p PlatformPlugin) GetLicense() (string, error) {
	return p.composer.License, nil
}

func (p PlatformPlugin) GetPath() string {
	return p.path
}

func (p PlatformPlugin) GetMetaData() *extensionMetadata {
	return &extensionMetadata{
		Label: extensionTranslated{
			German:  p.composer.Extra.Label["de-DE"],
			English: p.composer.Extra.Label["en-GB"],
		},
		Description: extensionTranslated{
			German:  p.composer.Extra.Description["de-DE"],
			English: p.composer.Extra.Description["en-GB"],
		},
	}
}

func (p PlatformPlugin) Validate(ctx *validationContext) {
	if len(p.composer.Name) == 0 {
		ctx.AddError("Key `name` is required")
	}

	if len(p.composer.Type) == 0 {
		ctx.AddError("Key `type` is required")
	} else if p.composer.Type != "shopware-platform-plugin" {
		ctx.AddError("The composer type must be shopware-platform-plugin")
	}

	if len(p.composer.Description) == 0 {
		ctx.AddError("Key `description` is required")
	}

	if len(p.composer.License) == 0 {
		ctx.AddError("Key `license` is required")
	}

	if len(p.composer.Version) == 0 {
		ctx.AddError("Key `version` is required")
	}

	if len(p.composer.Authors) == 0 {
		ctx.AddError("Key `authors` is required")
	}

	if len(p.composer.Require) == 0 {
		ctx.AddError("Key `require` is required")
	} else {
		_, exists := p.composer.Require["shopware/core"]

		if !exists {
			ctx.AddError("You need to require \"shopware/core\" package")
		}
	}

	requiredKeys := []string{"de-DE", "en-GB"}

	for _, key := range requiredKeys {
		_, hasLabel := p.composer.Extra.Label[key]
		_, hasDescription := p.composer.Extra.Description[key]
		_, hasManufacturer := p.composer.Extra.ManufacturerLink[key]
		_, hasSupportLink := p.composer.Extra.SupportLink[key]

		if !hasLabel {
			ctx.AddError(fmt.Sprintf("extra.label for language %s is required", key))
		}

		if !hasDescription {
			ctx.AddError(fmt.Sprintf("extra.description for language %s is required", key))
		}

		if !hasManufacturer {
			ctx.AddError(fmt.Sprintf("extra.manufacturerLink for language %s is required", key))
		}

		if !hasSupportLink {
			ctx.AddError(fmt.Sprintf("extra.supportLink for language %s is required", key))
		}
	}

	if len(p.composer.Autoload.Psr0) == 0 && len(p.composer.Autoload.Psr4) == 0 {
		ctx.AddError("At least one of the properties psr-0 or psr-4 are required in the composer.json")
	}

	validateTheme(ctx)
}
