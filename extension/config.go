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

type storeInfo struct {
	Tags       *[]string   `yaml:"tags"`
	Videos     *[]string   `yaml:"videos"`
	Hightlight *[]string   `yaml:"hightlights"`
	Features   *[]string   `yaml:"features"`
	Faq        *[]storeFaq `yaml:"faq"`
}

type Config struct {
	Store struct {
		Availabilities                      *[]string `yaml:"availabilities"`
		DefaultLocale                       *string   `yaml:"default_locale"`
		Localizations                       *[]string `yaml:"localizations"`
		Categories                          *[]string `yaml:"categories"`
		Type                                *string   `yaml:"type"`
		AutomaticBugfixVersionCompatibility *bool     `yaml:"automatic_bugfix_version_compatibility"`
		Info                                struct {
			German  storeInfo `yaml:"de"`
			English storeInfo `yaml:"en"`
		} `yaml:"info"`
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
	if config.Store.Info.English.Tags != nil && len(*config.Store.Info.English.Tags) > 5 {
		return fmt.Errorf("store.info.en.tags can contain maximal 5 items")
	}

	if config.Store.Info.German.Tags != nil && len(*config.Store.Info.German.Tags) > 5 {
		return fmt.Errorf("store.info.de.tags can contain maximal 5 items")
	}

	if config.Store.Info.English.Videos != nil && len(*config.Store.Info.English.Videos) > 2 {
		return fmt.Errorf("store.info.en.videos can contain maximal 2 items")
	}

	if config.Store.Info.German.Videos != nil && len(*config.Store.Info.German.Videos) > 2 {
		return fmt.Errorf("store.info.de.videos can contain maximal 2 items")
	}

	return nil
}
