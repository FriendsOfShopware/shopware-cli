package extension

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path"

	"github.com/FriendsOfShopware/shopware-cli/version"
)

type ShopwareBundle struct {
	path     string
	composer shopwareBundleComposerJson
	config   *Config
}

func newShopwareBundle(path string) (*ShopwareBundle, error) {
	composerJsonFile := fmt.Sprintf("%s/composer.json", path)
	if _, err := os.Stat(composerJsonFile); err != nil {
		return nil, err
	}

	jsonFile, err := os.ReadFile(composerJsonFile)
	if err != nil {
		return nil, fmt.Errorf("newShopwareBundle: %v", err)
	}

	var composerJson shopwareBundleComposerJson
	err = json.Unmarshal(jsonFile, &composerJson)

	if err != nil {
		return nil, fmt.Errorf("newShopwareBundle: %v", err)
	}

	if composerJson.Type != "shopware-bundle" {
		return nil, fmt.Errorf("newShopwareBundle: composer.json type is not shopware-bundle")
	}

	if composerJson.Extra.BundleName == "" {
		return nil, fmt.Errorf("composer.json does not contain shopware-bundle-name in extra")
	}

	cfg, err := readExtensionConfig(path)
	if err != nil {
		return nil, fmt.Errorf("newShopwareBundle: %v", err)
	}

	extension := ShopwareBundle{
		composer: composerJson,
		path:     path,
		config:   cfg,
	}

	return &extension, nil
}

type shopwareBundleComposerJson struct {
	Name    string                          `json:"name"`
	Type    string                          `json:"type"`
	License string                          `json:"license"`
	Version string                          `json:"version"`
	Require map[string]string               `json:"require"`
	Extra   shopwareBundleComposerJsonExtra `json:"extra"`
}

type shopwareBundleComposerJsonExtra struct {
	BundleName string `json:"shopware-bundle-name"`
}

// GetRootDir returns the src directory of the bundle.
func (p ShopwareBundle) GetRootDir() string {
	return path.Join(p.path, "src")
}

// GetResourcesDir returns the resources directory of the shopware bundle.
func (p ShopwareBundle) GetResourcesDir() string {
	return path.Join(p.GetRootDir(), "Resources")
}

func (p ShopwareBundle) GetName() (string, error) {
	return p.composer.Extra.BundleName, nil
}

func (p ShopwareBundle) GetExtensionConfig() *Config {
	return p.config
}

func (p ShopwareBundle) GetShopwareVersionConstraint() (*version.Constraints, error) {
	if p.config != nil && p.config.Build.ShopwareVersionConstraint != "" {
		constraint, err := version.NewConstraint(p.config.Build.ShopwareVersionConstraint)
		if err != nil {
			return nil, err
		}

		return &constraint, nil
	}

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

func (ShopwareBundle) GetType() string {
	return TypeShopwareBundle
}

func (p ShopwareBundle) GetVersion() (*version.Version, error) {
	return version.NewVersion(p.composer.Version)
}

func (p ShopwareBundle) GetChangelog() (*extensionTranslated, error) {
	return parseExtensionMarkdownChangelog(p)
}

func (p ShopwareBundle) GetLicense() (string, error) {
	return p.composer.License, nil
}

func (p ShopwareBundle) GetPath() string {
	return p.path
}

func (p ShopwareBundle) GetMetaData() *extensionMetadata {
	return &extensionMetadata{
		Label: extensionTranslated{
			German:  "FALLBACK",
			English: "FALLBACK",
		},
		Description: extensionTranslated{
			German:  "FALLBACK",
			English: "FALLBACK",
		},
	}
}

func (p ShopwareBundle) Validate(c context.Context, ctx *ValidationContext) {
}
