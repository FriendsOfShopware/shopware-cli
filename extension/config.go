package extension

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"os"
)

type storeInfo struct {
	Tags   []string `yaml:"tags"`
	Videos []string `yaml:"videos"`
}

type Config struct {
	Store struct {
		Availabilities []string `yaml:"availabilities"`
		DefaultLocale  string   `yaml:"default_locale"`
		Localizations  []string `yaml:"localizations"`
		Categories     []string `yaml:"categories"`
		Info           struct {
			German  storeInfo `yaml:"de"`
			English storeInfo `yaml:"en"`
		} `yaml:"info"`
	} `yaml:"store"`
}

func ReadExtensionConfig(dir string) (*Config, error) {
	config := Config{}

	fileName := fmt.Sprintf("%s/.shopware-extension.yml", dir)
	_, err := os.Stat(fileName)

	if err != nil {
		return nil, err
	}

	fileHandle, err := ioutil.ReadFile(fileName)

	if err != nil {
		return nil, fmt.Errorf("NewExtensionConfig: %v", err)
	}

	err = yaml.Unmarshal(fileHandle, &config)

	if err != nil {
		return nil, fmt.Errorf("NewExtensionConfig: %v", err)
	}

	return &config, nil
}
