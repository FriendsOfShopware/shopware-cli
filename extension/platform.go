package extension

import (
	"archive/zip"
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"mime/multipart"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"

	log "github.com/sirupsen/logrus"

	"github.com/FriendsOfShopware/shopware-cli/version"
)

type PlatformPlugin struct {
	path     string
	composer platformComposerJson
}

// GetRootDir returns the src directory of the plugin.
func (p PlatformPlugin) GetRootDir() string {
	return path.Join(p.path, "src")
}

// GetResourcesDir returns the resources directory of the plugin.
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
	Require  map[string]string         `json:"require"`
	Extra    platformComposerJsonExtra `json:"extra"`
	Autoload struct {
		Psr0 map[string]string `json:"psr-0"`
		Psr4 map[string]string `json:"psr-4"`
	} `json:"autoload"`
}

type platformComposerJsonExtra struct {
	ShopwarePluginClass string            `json:"shopware-plugin-class"`
	Label               map[string]string `json:"label"`
	Description         map[string]string `json:"description"`
	ManufacturerLink    map[string]string `json:"manufacturerLink"`
	SupportLink         map[string]string `json:"supportLink"`
	PluginIcon          string            `json:"plugin-icon"`
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

	pluginIcon := p.composer.Extra.PluginIcon

	if len(pluginIcon) == 0 {
		pluginIcon = "src/Resources/config/plugin.png"
	}

	// check if the plugin icon exists
	if _, err := os.Stat(filepath.Join(p.GetPath(), pluginIcon)); os.IsNotExist(err) {
		ctx.AddError(fmt.Sprintf("The plugin icon %s does not exist", pluginIcon))
	}

	validateTheme(ctx)
	validatePHPFiles(ctx)
}

type phpSyntaxCheckerResult struct {
	Errors []string `json:"errors"`
}

func validatePHPFiles(ctx *validationContext) {
	var b bytes.Buffer
	bufferW := bufio.NewWriter(&b)

	phpZip := zip.NewWriter(bufferW)

	_ = filepath.Walk(ctx.Extension.GetPath(), func(path string, info fs.FileInfo, err error) error {
		name := filepath.Base(path)

		if strings.HasSuffix(name, ".php") {
			zipFile, err := phpZip.Create(strings.TrimPrefix(path, ctx.Extension.GetPath()))

			if err != nil {
				return err
			}

			file, err := os.Open(path)

			if err != nil {
				return err
			}

			_, err = io.Copy(zipFile, file)

			if err != nil {
				return err
			}

			_ = file.Close()
		}

		return nil
	})

	_ = phpZip.Close()

	_ = bufferW.Flush()

	body := new(bytes.Buffer)
	multipartWriter := multipart.NewWriter(body)

	part, err := multipartWriter.CreateFormFile("file", "file.zip")

	if err != nil {
		ctx.AddError(fmt.Sprintf("Could not create form file: %s", err.Error()))
		return
	}

	_, err = part.Write(b.Bytes())

	if err != nil {
		ctx.AddError(fmt.Sprintf("Could not write zip file to multipart form: %s", err.Error()))
		return
	}

	_ = multipartWriter.Close()

	constraint, err := ctx.Extension.GetShopwareVersionConstraint()

	if err != nil {
		ctx.AddError(fmt.Sprintf("Could not parse shopware version constraint: %s", err.Error()))
		return
	}

	phpVersion, err := getPhpVersion(constraint)
	if err != nil {
		ctx.AddWarning(fmt.Sprintf("Could not find min php version for plugin: %s", err.Error()))
		return
	}

	log.Infof("Using php version %s for syntax check with https://github.com/FriendsOfShopware/aws-php-syntax-checker-lambda", phpVersion)

	req, err := http.NewRequestWithContext(context.Background(), http.MethodPost, fmt.Sprintf("https://php-syntax-checker.fos.gg/?version=%s", phpVersion), body)

	if err != nil {
		ctx.AddWarning(fmt.Sprintf("Could not create request to validate php files: %s", err.Error()))
		return
	}

	req.Header.Set("Content-Type", multipartWriter.FormDataContentType())

	resp, err := http.DefaultClient.Do(req)

	if err != nil {
		ctx.AddWarning(fmt.Sprintf("Could not validate php files: %s", err.Error()))
		return
	}

	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		return
	}

	var result phpSyntaxCheckerResult

	err = json.NewDecoder(resp.Body).Decode(&result)

	if err != nil {
		ctx.AddWarning(fmt.Sprintf("cannot decode php syntax checker response: %s", err.Error()))
		return
	}

	for _, error := range result.Errors {
		ctx.AddError(error)
	}
}

func getPhpVersion(constraint *version.Constraints) (string, error) {
	r, _ := http.NewRequestWithContext(context.Background(), http.MethodGet, "https://raw.githubusercontent.com/FriendsOfShopware/shopware-static-data/main/data/php-version.json", nil)

	resp, err := http.DefaultClient.Do(r)

	if err != nil {
		return "", err
	}

	defer resp.Body.Close()

	var shopwareToPHPVersion map[string]string

	err = json.NewDecoder(resp.Body).Decode(&shopwareToPHPVersion)

	if err != nil {
		return "", err
	}

	for shopwareVersion, phpVersion := range shopwareToPHPVersion {
		shopwareVersionConstraint, err := version.NewVersion(shopwareVersion)

		if err != nil {
			continue
		}

		if constraint.Check(shopwareVersionConstraint) {
			return phpVersion, nil
		}
	}

	return "", errors.New("could not find php version for shopware version")
}
