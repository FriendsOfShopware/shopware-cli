package config

import (
	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v3"
	"os"
	"path"
	"strconv"
	"sync"
	"testing"
)

func TestParseEnvConfig(t *testing.T) {
	testEnv := newTestEnv(t)
	defer testEnv.restore()
	defer resetState()

	testData := struct {
		email, password string
		companyId       int
	}{
		email:     "test@test.com",
		password:  "test123",
		companyId: 456,
	}

	testEnv.set("SHOPWARE_CLI_ACCOUNT_EMAIL", testData.email)
	testEnv.set("SHOPWARE_CLI_ACCOUNT_PASSWORD", testData.password)
	testEnv.set("SHOPWARE_CLI_ACCOUNT_COMPANY", strconv.Itoa(testData.companyId))

	assert.NoError(t, InitConfig(""))
	assert.True(t, state.loadedFromEnv)

	confService := Config{}
	assert.Equal(t, testData.email, confService.GetAccountEmail())
	assert.Equal(t, testData.password, confService.GetAccountPassword())
	assert.Equal(t, testData.companyId, confService.GetAccountCompanyId())
}

func TestParseFileConfig(t *testing.T) {
	defer resetState()
	testData := struct {
		email, password string
		companyId       int
	}{
		email:     "test@test.com",
		password:  "test123",
		companyId: 456,
	}

	cwd, err := os.Getwd()
	assert.NoError(t, err)
	testConfig := path.Join(cwd, "testdata/.shopware-cli.yml")

	assert.NoError(t, InitConfig(testConfig))
	assert.False(t, state.loadedFromEnv)

	confService := Config{}
	assert.Equal(t, testData.email, confService.GetAccountEmail())
	assert.Equal(t, testData.password, confService.GetAccountPassword())
	assert.Equal(t, testData.companyId, confService.GetAccountCompanyId())
	assert.Equal(t, testConfig, state.cfgPath)
}

func TestSaveConfig(t *testing.T) {
	defer resetState()
	testData := struct {
		email, password string
		companyId       int
	}{
		email:     "test@new.com",
		password:  "test",
		companyId: 111,
	}

	cwd, err := os.Getwd()
	assert.NoError(t, err)
	testConfig := path.Join(cwd, "testdata/write-test.yml")
	configBackup, err := os.ReadFile(testConfig)
	assert.NoError(t, err)
	defer func() {
		assert.NoError(t, os.WriteFile(testConfig, configBackup, os.ModePerm))
	}()

	assert.NoError(t, InitConfig(testConfig))

	configService := Config{}

	assert.NoError(t, configService.SetAccountEmail(testData.email))

	assert.NoError(t, configService.SetAccountPassword(testData.password))

	assert.NoError(t, configService.SetAccountCompanyId(testData.companyId))

	assert.True(t, state.modified)

	assert.NoError(t, SaveConfig())

	assert.False(t, state.modified)

	newConfData, err := os.ReadFile(testConfig)
	assert.NoError(t, err)

	var newConf configData
	assert.NoError(t, yaml.Unmarshal(newConfData, &newConf))

	assert.Equal(t, testData.email, newConf.Account.Email)
	assert.Equal(t, testData.password, newConf.Account.Password)
	assert.Equal(t, testData.companyId, newConf.Account.Company)
}

func TestDontWriteEnvConfig(t *testing.T) {
	testEnv := newTestEnv(t)
	defer testEnv.restore()
	defer resetState()

	testData := struct {
		email, password string
		companyId       int
	}{
		email:     "test@test.com",
		password:  "test123",
		companyId: 456,
	}

	testEnv.set("SHOPWARE_CLI_ACCOUNT_EMAIL", testData.email)
	testEnv.set("SHOPWARE_CLI_ACCOUNT_PASSWORD", testData.password)
	testEnv.set("SHOPWARE_CLI_ACCOUNT_COMPANY", strconv.Itoa(testData.companyId))

	assert.NoError(t, InitConfig(""))
	assert.True(t, state.loadedFromEnv)

	confService := Config{}
	assert.Error(t, confService.SetAccountEmail("test@foo.com"))
	assert.Error(t, confService.SetAccountPassword("S3CR3TF4RT3St"))
	assert.Error(t, confService.SetAccountCompanyId(111))
}

type testEnv struct {
	t       *testing.T
	oldVars map[string]string
}

func newTestEnv(t *testing.T) *testEnv {
	t.Helper()
	return &testEnv{
		t,
		map[string]string{},
	}
}

func (e *testEnv) set(key, value string) {
	val := os.Getenv(key)
	e.oldVars[key] = val
	if err := os.Setenv(key, value); err != nil {
		e.t.Fatal(err)
	}
}

func (e *testEnv) restore() {
	for key, value := range e.oldVars {
		var err error
		if len(value) > 0 {
			err = os.Setenv(key, value)
		} else {
			err = os.Unsetenv(key)
		}
		if err != nil {
			e.t.Error(err)
		}
	}
}

func resetState() {
	state = &configState{
		mu:      sync.RWMutex{},
		cfgPath: "",
		inner:   defaultConfig(),
	}
}
