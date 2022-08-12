package extension

import (
	"fmt"
	"io/ioutil"
	"os"

	"gopkg.in/yaml.v3"
)

type ConfigBuild struct {
	Zip struct {
		Composer struct {
			Enabled          bool     `yaml:"enabled"`
			BeforeHooks      []string `yaml:"before_hooks"`
			AfterHooks       []string `yaml:"after_hooks"`
			ExcludedPackages []string `yaml:"excluded_packages"`
		} `yaml:"composer"`
		Assets struct {
			Enabled               bool     `yaml:"enabled"`
			BeforeHooks           []string `yaml:"before_hooks"`
			AfterHooks            []string `yaml:"after_hooks"`
			EnableESBuildForAdmin bool     `yaml:"enable_es_build_for_admin"`
		} `yaml:"assets"`
		Pack struct {
			Excludes struct {
				Paths []string `yaml:"paths"`
			} `yaml:"excludes"`
			BeforeHooks []string `yaml:"before_hooks"`
		} `yaml:"pack"`
	} `yaml:"zip"`
}

type ConfigStore struct {
	Availabilities                      *[]string                          `yaml:"availabilities"`
	DefaultLocale                       *string                            `yaml:"default_locale"`
	Localizations                       *[]string                          `yaml:"localizations"`
	Categories                          *[]string                          `yaml:"categories"`
	Type                                *string                            `yaml:"type"`
	Icon                                *string                            `yaml:"icon"`
	AutomaticBugfixVersionCompatibility *bool                              `yaml:"automatic_bugfix_version_compatibility"`
	Description                         ConfigTranslated[string]           `yaml:"description"`
	InstallationManual                  ConfigTranslated[string]           `yaml:"installation_manual"`
	Tags                                ConfigTranslated[[]string]         `yaml:"tags"`
	Videos                              ConfigTranslated[[]string]         `yaml:"videos"`
	Highlights                          ConfigTranslated[[]string]         `yaml:"highlights"`
	Features                            ConfigTranslated[[]string]         `yaml:"features"`
	Faq                                 ConfigTranslated[[]ConfigStoreFaq] `yaml:"faq"`
	Images                              *[]ConfigStoreImage                `yaml:"images"`
}

type Translatable interface {
	string | []string | []ConfigStoreFaq
}

type ConfigTranslated[T Translatable] struct {
	German  *T `yaml:"de"`
	English *T `yaml:"en"`
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
	Build ConfigBuild `yaml:"build"`
}

func ReadExtensionConfig(dir string) (*Config, error) {
	errorFormat := "ReadExtensionConfig: %v"
	config := &Config{}
	config.Build.Zip.Assets.Enabled = true
	config.Build.Zip.Composer.Enabled = true

	fileName := fmt.Sprintf("%s/.shopware-extension.yml", dir)
	_, err := os.Stat(fileName)

	if os.IsNotExist(err) {
		return config, nil
	}

	if err != nil {
		return nil, err
	}

	fileHandle, err := ioutil.ReadFile(fileName)

	if err != nil {
		return nil, fmt.Errorf(errorFormat, err)
	}

	err = yaml.Unmarshal(fileHandle, &config)

	if err != nil {
		return nil, fmt.Errorf(errorFormat, err)
	}

	err = validateExtensionConfig(config)
	if err != nil {
		return nil, fmt.Errorf(errorFormat, err)
	}

	return config, nil
}

func validateExtensionConfig(config *Config) error {
	if config.Store.Tags.English != nil && len(*config.Store.Tags.English) > 5 {
		return fmt.Errorf("store.info.tags.en can contain maximal 5 items")
	}

	if config.Store.Tags.German != nil && len(*config.Store.Tags.German) > 5 {
		return fmt.Errorf("store.info.tags.de can contain maximal 5 items")
	}

	if config.Store.Videos.English != nil && len(*config.Store.Videos.English) > 2 {
		return fmt.Errorf("store.info.videos.en can contain maximal 2 items")
	}

	if config.Store.Videos.German != nil && len(*config.Store.Videos.German) > 2 {
		return fmt.Errorf("store.info.videos.de can contain maximal 2 items")
	}

	return nil
}
