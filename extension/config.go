package extension

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"os"
)

type storeFaq struct {
	Question string `yaml:"question"`
	Answer   string `yaml:"answer"`
}

type storeImage struct {
	File     string `yaml:"file"`
	Activate struct {
		German  bool `yaml:"de"`
		English bool `yaml:"en"`
	}
	Preview struct {
		German  bool `yaml:"de"`
		English bool `yaml:"en"`
	}
	Priority int `yaml:"priority"`
}

type Config struct {
	Store struct {
		Availabilities                      *[]string `yaml:"availabilities"`
		DefaultLocale                       *string   `yaml:"default_locale"`
		Localizations                       *[]string `yaml:"localizations"`
		Categories                          *[]string `yaml:"categories"`
		Type                                *string   `yaml:"type"`
		Icon                                *string   `yaml:"icon"`
		AutomaticBugfixVersionCompatibility *bool     `yaml:"automatic_bugfix_version_compatibility"`
		Description                         struct {
			German  *string `yaml:"de"`
			English *string `yaml:"en"`
		} `yaml:"description"`
		InstallationManual struct {
			German  *string `yaml:"de"`
			English *string `yaml:"en"`
		} `yaml:"installation_manual"`
		Tags struct {
			German  *[]string `yaml:"de"`
			English *[]string `yaml:"en"`
		} `yaml:"tags"`
		Videos struct {
			German  *[]string `yaml:"de"`
			English *[]string `yaml:"en"`
		} `yaml:"videos"`
		Highlights struct {
			German  *[]string `yaml:"de"`
			English *[]string `yaml:"en"`
		} `yaml:"highlights"`
		Features struct {
			German  *[]string `yaml:"de"`
			English *[]string `yaml:"en"`
		} `yaml:"features"`
		Faq struct {
			German  *[]storeFaq `yaml:"de"`
			English *[]storeFaq `yaml:"en"`
		} `yaml:"faq"`
		Images *[]storeImage `yaml:"images"`
	} `yaml:"store"`
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
