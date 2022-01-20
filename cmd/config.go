package cmd

import (
	"fmt"
	"github.com/caarlos0/env/v6"
	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"os"
)

type Config struct {
	loadedWithEnv bool
	Account       struct {
		Email    string `env:"SHOPWARE_CLI_ACCOUNT_EMAIL" yaml:"email"`
		Password string `env:"SHOPWARE_CLI_ACCOUNT_PASSWORD" yaml:"password"`
		Company  int    `env:"SHOPWARE_CLI_ACCOUNT_COMPANY" yaml:"company"`
	} `yaml:"account"`
}

var appConfig *Config

func getApplicationConfigPath() string {
	if cfgFile != "" {
		return cfgFile
	}

	configDir, err := os.UserConfigDir()

	if err != nil {
		return ".shopware-cli.yml"
	}

	cfgFile = fmt.Sprintf("%s/.shopware-cli.yml", configDir)

	return cfgFile
}

func initApplicationConfig() error {
	appConfig = &Config{}

	err := env.Parse(appConfig)
	if err != nil {
		return err
	}

	if len(appConfig.Account.Email) > 0 {
		appConfig.loadedWithEnv = true

		log.Tracef("Loaded config with environment variables")

		return nil
	}

	cfg := getApplicationConfigPath()
	if _, err := os.Stat(cfg); os.IsNotExist(err) {
		return nil
	}

	content, err := ioutil.ReadFile(cfgFile)

	if err != nil {
		return err
	}

	err = yaml.Unmarshal(content, &appConfig)

	if err != nil {
		return err
	}

	log.Tracef("Using config file from %s", cfgFile)

	return nil
}

func saveApplicationConfig() error {
	if appConfig.loadedWithEnv {
		return nil
	}

	content, err := yaml.Marshal(appConfig)

	if err != nil {
		return err
	}

	return ioutil.WriteFile(getApplicationConfigPath(), content, os.ModePerm)
}
