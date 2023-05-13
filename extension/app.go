package extension

import (
	"context"
	"encoding/xml"
	"fmt"
	"os"
	"path"
	"path/filepath"

	"github.com/FriendsOfShopware/shopware-cli/version"
)

type translatedXmlNode []struct {
	Text string `xml:",chardata"`
	Lang string `xml:"lang,attr"`
}

type appManifest struct {
	XMLName                   xml.Name        `xml:"manifest"`
	Text                      string          `xml:",chardata"`
	Xsi                       string          `xml:"xsi,attr"`
	NoNamespaceSchemaLocation string          `xml:"noNamespaceSchemaLocation,attr"`
	Meta                      appManifestMeta `xml:"meta"`
	Setup                     struct {
		Text            string `xml:",chardata"`
		RegistrationUrl string `xml:"registrationUrl"`
		Secret          string `xml:"secret"`
	} `xml:"setup"`
	Permissions struct {
		Text   string `xml:",chardata"`
		Read   string `xml:"read"`
		Create string `xml:"create"`
		Update string `xml:"update"`
		Delete string `xml:"delete"`
	} `xml:"permissions"`
	Webhooks struct {
		Text    string `xml:",chardata"`
		Webhook struct {
			Text  string `xml:",chardata"`
			Name  string `xml:"name,attr"`
			URL   string `xml:"url,attr"`
			Event string `xml:"event,attr"`
		} `xml:"webhook"`
	} `xml:"webhooks"`
	Admin struct {
		Text   string `xml:",chardata"`
		Module []struct {
			Text     string `xml:",chardata"`
			Name     string `xml:"name,attr"`
			Parent   string `xml:"parent,attr"`
			Position string `xml:"position,attr"`
			Source   string `xml:"source,attr"`
			Label    []struct {
				Text string `xml:",chardata"`
				Lang string `xml:"lang,attr"`
			} `xml:"label"`
		} `xml:"module"`
		MainModule struct {
			Text   string `xml:",chardata"`
			Source string `xml:"source,attr"`
		} `xml:"main-module"`
		ActionButton []struct {
			Text   string `xml:",chardata"`
			Action string `xml:"action,attr"`
			Entity string `xml:"entity,attr"`
			View   string `xml:"view,attr"`
			URL    string `xml:"url,attr"`
			Label  string `xml:"label"`
		} `xml:"action-button"`
	} `xml:"admin"`
	CustomFields struct {
		Text           string `xml:",chardata"`
		CustomFieldSet struct {
			Text  string `xml:",chardata"`
			Name  string `xml:"name"`
			Label []struct {
				Text string `xml:",chardata"`
				Lang string `xml:"lang,attr"`
			} `xml:"label"`
			RelatedEntities struct {
				Text  string `xml:",chardata"`
				Order string `xml:"order"`
			} `xml:"related-entities"`
			Fields struct {
				Chardata string `xml:",chardata"`
				Text     struct {
					Text     string `xml:",chardata"`
					Name     string `xml:"name,attr"`
					Label    string `xml:"label"`
					Position string `xml:"position"`
					Required string `xml:"required"`
					HelpText string `xml:"help-text"`
				} `xml:"text"`
				Float struct {
					Text  string `xml:",chardata"`
					Name  string `xml:"name,attr"`
					Label []struct {
						Text string `xml:",chardata"`
						Lang string `xml:"lang,attr"`
					} `xml:"label"`
					HelpText    string `xml:"help-text"`
					Position    string `xml:"position"`
					Placeholder string `xml:"placeholder"`
					Min         string `xml:"min"`
					Max         string `xml:"max"`
					Steps       string `xml:"steps"`
				} `xml:"float"`
			} `xml:"fields"`
		} `xml:"custom-field-set"`
	} `xml:"custom-fields"`
	Cookies struct {
		Text   string `xml:",chardata"`
		Cookie struct {
			Text               string `xml:",chardata"`
			Cookie             string `xml:"cookie"`
			SnippetName        string `xml:"snippet-name"`
			SnippetDescription string `xml:"snippet-description"`
			Value              string `xml:"value"`
			Expiration         string `xml:"expiration"`
		} `xml:"cookie"`
		Group struct {
			Text               string `xml:",chardata"`
			SnippetName        string `xml:"snippet-name"`
			SnippetDescription string `xml:"snippet-description"`
			Entries            struct {
				Text   string `xml:",chardata"`
				Cookie struct {
					Text               string `xml:",chardata"`
					Cookie             string `xml:"cookie"`
					SnippetName        string `xml:"snippet-name"`
					SnippetDescription string `xml:"snippet-description"`
					Value              string `xml:"value"`
					Expiration         string `xml:"expiration"`
				} `xml:"cookie"`
			} `xml:"entries"`
		} `xml:"group"`
	} `xml:"cookies"`
	Payments struct {
		Text          string `xml:",chardata"`
		PaymentMethod struct {
			Text       string `xml:",chardata"`
			Identifier string `xml:"identifier"`
			Name       []struct {
				Text string `xml:",chardata"`
				Lang string `xml:"lang,attr"`
			} `xml:"name"`
			Description []struct {
				Text string `xml:",chardata"`
				Lang string `xml:"lang,attr"`
			} `xml:"description"`
			PayURL      string `xml:"pay-url"`
			FinalizeURL string `xml:"finalize-url"`
			Icon        string `xml:"icon"`
		} `xml:"payment-method"`
	} `xml:"payments"`
}

type appManifestMeta struct {
	Text                    string            `xml:",chardata"`
	Name                    string            `xml:"name"`
	Label                   translatedXmlNode `xml:"label"`
	Description             translatedXmlNode `xml:"description"`
	Author                  string            `xml:"author"`
	Copyright               string            `xml:"copyright"`
	Version                 string            `xml:"version"`
	License                 string            `xml:"license"`
	Icon                    string            `xml:"icon"`
	Privacy                 string            `xml:"privacy"`
	PrivacyPolicyExtensions []struct {
		Text string `xml:",chardata"`
		Lang string `xml:"lang,attr"`
	} `xml:"privacyPolicyExtensions"`
}

func getTranslatedTextFromXmlNode(node translatedXmlNode, keys []string) string {
	for _, n := range node {
		for _, key := range keys {
			if n.Lang == key {
				return n.Text
			}
		}
	}

	return ""
}

type App struct {
	path     string
	manifest appManifest
	config   *Config
}

func (a App) GetRootDir() string {
	return a.path
}

func (a App) GetResourcesDir() string {
	return path.Join(a.path, "Resources")
}

func newApp(path string) (*App, error) {
	appFileName := fmt.Sprintf("%s/manifest.xml", path)

	if _, err := os.Stat(appFileName); err != nil {
		return nil, err
	}

	appFile, err := os.ReadFile(appFileName)
	if err != nil {
		return nil, fmt.Errorf("newApp: %v", err)
	}

	var manifest appManifest
	err = xml.Unmarshal(appFile, &manifest)

	if err != nil {
		return nil, fmt.Errorf("newApp: %v", err)
	}

	cfg, err := readExtensionConfig(path)
	if err != nil {
		return nil, fmt.Errorf("newApp: %v", err)
	}

	app := App{
		path:     path,
		manifest: manifest,
		config:   cfg,
	}

	return &app, nil
}

func (a App) GetName() (string, error) {
	return a.manifest.Meta.Name, nil
}

func (a App) GetVersion() (*version.Version, error) {
	return version.NewVersion(a.manifest.Meta.Version)
}

func (a App) GetLicense() (string, error) {
	return a.manifest.Meta.License, nil
}

func (a App) GetExtensionConfig() *Config {
	return a.config
}

func (a App) GetShopwareVersionConstraint() (*version.Constraints, error) {
	if a.config.Build.ShopwareVersionConstraint != "" {
		v, err := version.NewConstraint(a.config.Build.ShopwareVersionConstraint)
		if err != nil {
			return nil, err
		}

		return &v, err
	}

	v, err := version.NewConstraint("~6.4")
	if err != nil {
		return nil, err
	}

	return &v, err
}

func (App) GetType() string {
	return TypePlatformApp
}

func (a App) GetPath() string {
	return a.path
}

func (a App) GetChangelog() (*extensionTranslated, error) {
	return parseExtensionMarkdownChangelog(a)
}

func (a App) GetMetaData() *extensionMetadata {
	german := []string{"de-DE", "de"}
	english := []string{"en-GB", "en-US", "en", ""}

	return &extensionMetadata{
		Label: extensionTranslated{
			German:  getTranslatedTextFromXmlNode(a.manifest.Meta.Label, german),
			English: getTranslatedTextFromXmlNode(a.manifest.Meta.Label, english),
		},
		Description: extensionTranslated{
			German:  getTranslatedTextFromXmlNode(a.manifest.Meta.Label, german),
			English: getTranslatedTextFromXmlNode(a.manifest.Meta.Label, english),
		},
	}
}

func (a App) Validate(_ context.Context, ctx *ValidationContext) {
	validateTheme(ctx)

	appIcon := a.manifest.Meta.Icon

	if appIcon == "" {
		appIcon = "Resources/config/plugin.png"
	}

	if _, err := os.Stat(filepath.Join(a.GetPath(), appIcon)); os.IsNotExist(err) {
		ctx.AddError(fmt.Sprintf("Cannot find app icon at %s", appIcon))
	}
}
