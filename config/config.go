package config

import (
	"fmt"
	"github.com/caarlos0/env/v6"
	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"os"
	"strconv"
	"sync"
)

var state *configState

type configState struct {
	mu            sync.RWMutex
	cfgPath       string
	inner         *configData
	loadedFromEnv bool
	isReady       bool
	modified      bool
}

type configData struct {
	Account struct {
		Email    string `env:"SHOPWARE_CLI_ACCOUNT_EMAIL" yaml:"email"`
		Password string `env:"SHOPWARE_CLI_ACCOUNT_PASSWORD" yaml:"password"`
		Company  int    `env:"SHOPWARE_CLI_ACCOUNT_COMPANY" yaml:"company"`
	} `yaml:"account"`
}

type Config struct{}

func init() {
	state = &configState{
		mu:      sync.RWMutex{},
		cfgPath: "",
		inner:   defaultConfig(),
	}
}

func defaultConfig() *configData {
	config := &configData{}
	config.Account.Email = ""
	config.Account.Password = ""
	config.Account.Company = 1
	return config
}

func InitConfig(configPath string) error {
	state.mu.Lock()
	defer state.mu.Unlock()
	if state.isReady {
		return nil
	}

	if len(configPath) > 0 {
		state.cfgPath = configPath
	} else {
		configDir, err := os.UserConfigDir()

		if err != nil {
			return err
		}

		state.cfgPath = fmt.Sprintf("%s/.shopware-cli.yml", configDir)
	}

	err := env.Parse(state.inner)
	if err != nil {
		return err
	}
	if len(state.inner.Account.Email) > 0 {
		state.loadedFromEnv = true

		state.isReady = true
		log.Tracef("Loaded config with environment variables")

		return nil
	}
	if _, err := os.Stat(state.cfgPath); os.IsNotExist(err) {
		if err := createNewConfig(state.cfgPath); err != nil {
			return err
		}
	}

	content, err := ioutil.ReadFile(state.cfgPath)

	if err != nil {
		return err
	}

	err = yaml.Unmarshal(content, &state.inner)

	if err != nil {
		return err
	}

	log.Tracef("Using config file from %s", state.cfgPath)
	state.isReady = true
	return nil
}

func SaveConfig() error {
	state.mu.Lock()
	defer state.mu.Unlock()
	if !state.modified || state.loadedFromEnv {
		return nil
	}

	configFile, err := os.OpenFile(state.cfgPath, os.O_CREATE|os.O_WRONLY, os.ModePerm)
	if err != nil {
		return err
	}
	configWriter := yaml.NewEncoder(configFile)
	defer func() {
		state.modified = false
		_ = configWriter.Close()
	}()

	return configWriter.Encode(state.inner)
}

func createNewConfig(path string) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()
	encoder := yaml.NewEncoder(f)
	return encoder.Encode(defaultConfig())
}

func (_ Config) GetAccountEmail() string {
	state.mu.RLock()
	defer state.mu.RUnlock()
	return state.inner.Account.Email
}

func (_ Config) GetAccountPassword() string {
	state.mu.RLock()
	defer state.mu.RUnlock()
	return state.inner.Account.Password
}

func (_ Config) GetAccountCompanyId() int {
	state.mu.RLock()
	defer state.mu.RUnlock()
	return state.inner.Account.Company
}

func (_ Config) SetAccountEmail(email string) error {
	state.mu.Lock()
	defer state.mu.Unlock()
	if state.loadedFromEnv {
		return fmt.Errorf("could not set config value %s to %q config was loaded from env",
			"account.email",
			email,
		)
	}
	state.modified = true
	state.inner.Account.Email = email
	return nil
}

func (_ Config) SetAccountPassword(password string) error {
	state.mu.Lock()
	defer state.mu.Unlock()
	if state.loadedFromEnv {
		return fmt.Errorf("could not set config value %s to %q config was loaded from env",
			"account.password",
			"***",
		)
	}
	state.modified = true
	state.inner.Account.Password = password
	return nil
}

func (_ Config) SetAccountCompanyId(id int) error {
	state.mu.Lock()
	defer state.mu.Unlock()
	if state.loadedFromEnv {
		return fmt.Errorf("could not set config value %s to %q config was loaded from env",
			"account.company",
			strconv.Itoa(id),
		)
	}
	state.modified = true
	state.inner.Account.Company = id
	return nil
}

func (_ Config) Save() error {
	return SaveConfig()
}
