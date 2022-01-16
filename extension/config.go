package extension

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"os"
)

type ConfigStore struct {
	Availabilities                      *[]string                  `yaml:"availabilities"`
	DefaultLocale                       *string                    `yaml:"default_locale"`
	Localizations                       *[]string                  `yaml:"localizations"`
	Categories                          *[]string                  `yaml:"categories"`
	Type                                *string                    `yaml:"type"`
	Icon                                *string                    `yaml:"icon"`
	AutomaticBugfixVersionCompatibility *bool                      `yaml:"automatic_bugfix_version_compatibility"`
	Description                         ConfigTranslatedString     `yaml:"description"`
	InstallationManual                  ConfigTranslatedString     `yaml:"installation_manual"`
	Tags                                ConfigTranslatedStringList `yaml:"tags"`
	Videos                              ConfigTranslatedStringList `yaml:"videos"`
	Highlights                          ConfigTranslatedStringList `yaml:"highlights"`
	Features                            ConfigTranslatedStringList `yaml:"features"`
	Faq                                 ConfigStoreTranslatedFaq   `yaml:"faq"`
	Images                              *[]ConfigStoreImage        `yaml:"images"`
}

type ConfigTranslatedString struct {
	German  *string `yaml:"de"`
	English *string `yaml:"en"`
}

type ConfigTranslatedStringList struct {
	German  *[]string `yaml:"de"`
	English *[]string `yaml:"en"`
}

type ConfigStoreTranslatedFaq struct {
	German  *[]ConfigStoreFaq `yaml:"de"`
	English *[]ConfigStoreFaq `yaml:"en"`
}

type ConfigStoreFaq struct {
	Question string `yaml:"question"`
	Answer   string `yaml:"answer"`
}

type ConfigStoreImage struct {
	File     string                   `yaml:"file"`
	Activate ConfigStoreImageActivate `yaml:"activate"`
	Preview  ConfigStoreImagePreview  `yaml:"preview"`
	Priority int                      `yaml:"priority"`
}

type ConfigStoreImageActivate struct {
	German  bool `yaml:"de"`
	English bool `yaml:"en"`
}

type ConfigStoreImagePreview struct {
	German  bool `yaml:"de"`
	English bool `yaml:"en"`
}

type Config struct {
	Store ConfigStore `yaml:"store"`
}

func ReadExtensionConfig(dir string) (*Config, error) {
	var config *Config

	fileName := fmt.Sprintf("%s/.shopware-extension.yml", dir)
	_, err := os.Stat(fileName)

	if err != nil {
		return nil, nil
	}

	fileHandle, err := ioutil.ReadFile(fileName)

	if err != nil {
		return nil, fmt.Errorf("NewExtensionConfig: %v", err)
	}

	err = yaml.Unmarshal(fileHandle, &config)

	if err != nil {
		return nil, fmt.Errorf("NewExtensionConfig: %v", err)
	}

	err = validateExtensionConfig(config)
	if err != nil {
		return nil, fmt.Errorf("NewExtensionConfig: %v", err)
	}

	return config, nil
}

func validateExtensionConfig(config *Config) error {
	if config.Store.Tags.English != nil && len(*config.Store.Tags.English) > 5 {
		return fmt.Errorf("store.info.tags.en can contain maximal 5 items")
	}

	if config.Store.Tags.German != nil && len(*config.Store.Tags.German) > 5 {
		return fmt.Errorf("store.info.tags.en can contain maximal 5 items")
	}

	if config.Store.Videos.English != nil && len(*config.Store.Videos.English) > 2 {
		return fmt.Errorf("store.info.videos.en can contain maximal 2 items")
	}

	if config.Store.Videos.German != nil && len(*config.Store.Videos.German) > 2 {
		return fmt.Errorf("store.info.videos.de can contain maximal 2 items")
	}

	return nil
}
