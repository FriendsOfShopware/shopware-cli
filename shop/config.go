package shop

import (
	"dario.cat/mergo"
	"fmt"
	"os"
	"strings"

	"github.com/doutorfinancas/go-mad/core"
	adminSdk "github.com/friendsofshopware/go-shopware-admin-api-sdk"
	"github.com/google/uuid"
	"gopkg.in/yaml.v3"
)

type Config struct {
	AdditionalConfigs []string        `yaml:"include,omitempty"`
	URL               string          `yaml:"url"`
	Build             *ConfigBuild    `yaml:"build,omitempty"`
	AdminApi          *ConfigAdminApi `yaml:"admin_api,omitempty"`
	ConfigDump        *ConfigDump     `yaml:"dump,omitempty"`
	Sync              *ConfigSync     `yaml:"sync,omitempty"`
	Docker            *ConfigDocker   `yaml:"docker,omitempty"`
}

type ConfigBuild struct {
	DisableAssetCopy      bool     `yaml:"disable_asset_copy,omitempty"`
	RemoveExtensionAssets bool     `yaml:"remove_extension_assets,omitempty"`
	KeepExtensionSource   bool     `yaml:"keep_extension_source,omitempty"`
	CleanupPaths          []string `yaml:"cleanup_paths,omitempty"`
	Browserslist          string   `yaml:"browserslist,omitempty"`
	ExcludeExtensions     []string `yaml:"exclude_extensions,omitempty"`
}

type ConfigDockerPHP struct {
	PhpVersion string            `yaml:"version,omitempty"`
	Extensions []string          `yaml:"extensions,omitempty"`
	Settings   map[string]string `yaml:"ini,omitempty"`
}

type ConfigDocker struct {
	PHP          ConfigDockerPHP `yaml:"php"`
	ExcludePaths []string        `yaml:"exclude_paths,omitempty"`
}

type ConfigAdminApi struct {
	ClientId        string `yaml:"client_id,omitempty"`
	ClientSecret    string `yaml:"client_secret,omitempty"`
	Username        string `yaml:"username,omitempty"`
	Password        string `yaml:"password,omitempty"`
	DisableSSLCheck bool   `yaml:"disable_ssl_check,omitempty"`
}

type ConfigDump struct {
	Rewrite map[string]core.Rewrite `yaml:"rewrite,omitempty"`
	NoData  []string                `yaml:"nodata,omitempty"`
	Ignore  []string                `yaml:"ignore,omitempty"`
	Where   map[string]string       `yaml:"where,omitempty"`
}

type ConfigSync struct {
	Config       []ConfigSyncConfig `yaml:"config"`
	Theme        []ThemeConfig      `yaml:"theme"`
	MailTemplate []MailTemplate     `yaml:"mail_template"`
	Entity       []EntitySync       `yaml:"entity"`
}

type ConfigSyncConfig struct {
	SalesChannel *string                `yaml:"sales_channel,omitempty"`
	Settings     map[string]interface{} `yaml:"settings"`
}

type ThemeConfig struct {
	Name     string                               `yaml:"name"`
	Settings map[string]adminSdk.ThemeConfigValue `yaml:"settings"`
}

type MailTemplate struct {
	Id           string                    `yaml:"id"`
	Translations []MailTemplateTranslation `yaml:"translations"`
}

type EntitySync struct {
	Entity  string                 `yaml:"entity"`
	Exists  *[]interface{}         `yaml:"exists"`
	Payload map[string]interface{} `yaml:"payload"`
}

type MailTemplateTranslation struct {
	Language     string      `yaml:"language"`
	SenderName   string      `yaml:"sender_name"`
	Subject      string      `yaml:"subject"`
	HTML         string      `yaml:"html"`
	Plain        string      `yaml:"plain"`
	CustomFields interface{} `yaml:"custom_fields"`
}

func ReadConfig(fileName string, allowFallback bool) (*Config, error) {
	config := &Config{}

	_, err := os.Stat(fileName)

	if os.IsNotExist(err) {
		if allowFallback {
			return fillEmptyConfig(config), nil
		}

		return nil, fmt.Errorf("cannot find project configuration file \"%s\", use shopware-cli project config init to create one", fileName)
	}

	if err != nil {
		return nil, err
	}

	fileHandle, err := os.ReadFile(fileName)
	if err != nil {
		return nil, fmt.Errorf("ReadConfig: %v", err)
	}

	substitutedConfig := os.ExpandEnv(string(fileHandle))
	err = yaml.Unmarshal([]byte(substitutedConfig), &config)

	if len(config.AdditionalConfigs) > 0 {
		for _, additionalConfigFile := range config.AdditionalConfigs {
			additionalConfig, err := ReadConfig(additionalConfigFile, allowFallback)
			if err != nil {
				return nil, fmt.Errorf("error while reading included config: %s", err.Error())
			}

			err = mergo.Merge(additionalConfig, config, mergo.WithOverride, mergo.WithSliceDeepCopy)
			if err != nil {
				return nil, fmt.Errorf("error while merging included config: %s", err.Error())
			}

			config = additionalConfig
		}
	}

	if err != nil {
		return nil, fmt.Errorf("ReadConfig: %v", err)
	}

	return fillEmptyConfig(config), nil
}

func fillEmptyConfig(c *Config) *Config {
	if c.Build == nil {
		c.Build = &ConfigBuild{}
	}

	if c.Docker == nil {
		c.Docker = &ConfigDocker{
			PHP: ConfigDockerPHP{
				Extensions: make([]string, 0),
				Settings:   make(map[string]string),
			},
		}
	}

	return c
}

func NewUuid() string {
	return strings.ReplaceAll(uuid.New().String(), "-", "")
}
