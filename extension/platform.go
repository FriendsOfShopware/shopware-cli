package extension

import (
	"encoding/json"
	"fmt"
	"github.com/hashicorp/go-version"
	"io/ioutil"
	"os"
	"strings"
)

type PlatformPlugin struct {
	path     string
	composer platformComposerJson
}

func newPlatformPlugin(path string) (*PlatformPlugin, error) {
	composerJsonFile := fmt.Sprintf("%s/composer.json", path)
	if _, err := os.Stat(composerJsonFile); err != nil {
		return nil, err
	}

	jsonFile, err := ioutil.ReadFile(composerJsonFile)

	if err != nil {
		return nil, fmt.Errorf("newPlatformPlugin: %v", err)
	}

	var composerJson platformComposerJson
	err = json.Unmarshal(jsonFile, &composerJson)

	if err != nil {
		return nil, fmt.Errorf("newPlatformPlugin: %v", err)
	}

	parts := strings.Split(composerJson.Extra.ShopwarePluginClass, "\\")
	shopwareConstraintString, ok := composerJson.Require["shopware/core"]

	if !ok {
		return nil, fmt.Errorf("newPlatformPlugin: require.shopware/core is required")
	}

	shopwareConstraint, err := version.NewConstraint(shopwareConstraintString)

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

func (p PlatformPlugin) GetType() string {
	return "platform"
}

func (p PlatformPlugin) GetVersion() (*version.Version, error) {
	v, err := version.NewVersion(p.composer.Version)

	if err != nil {
		return nil, err
	}

	return v, nil
}

func (p PlatformPlugin) GetChangelog() (*extensionTranslated, error) {
	v, err := p.GetVersion()

	if err != nil {
		return nil, err
	}

	changelogs, err := parseMarkdownChangelogInPath(p.path)

	if err != nil {
		return nil, err
	}

	changelogDe, ok := changelogs["de-DE"]

	if !ok {
		return nil, fmt.Errorf("german changelog is missing")
	}

	changelogDeVersion, ok := changelogDe[v.String()]

	if !ok {
		return nil, fmt.Errorf("german changelog in version %s is missing", v.String())
	}

	changelogEn, ok := changelogs["en-GB"]

	changelogEnVersion, ok := changelogEn[v.String()]

	if !ok {
		return nil, fmt.Errorf("english changelog in version %s is missing", v.String())
	}

	if !ok {
		return nil, fmt.Errorf("english changelog is missing")
	}

	return &extensionTranslated{German: changelogDeVersion, English: changelogEnVersion}, nil
}

func (p PlatformPlugin) GetLicense() (string, error) {
	return p.composer.License, nil
}

func (p PlatformPlugin) GetPath() string {
	return p.path
}

func (p PlatformPlugin) GetMetaData() (*extensionMetadata, error) {
	return nil, fmt.Errorf("not implemented")
}
