package config

import (
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

	if err := InitConfig(""); err != nil {
		t.Fatalf("unexpectd err: %q", err)
	}
	if !state.loadedFromEnv {
		t.Fatal("expected loadedWithEnv to be true")
	}

	confService := Config{}
	if email := confService.GetAccountEmail(); email != testData.email {
		t.Errorf("expected Email to be %q got %q", testData.email, email)
	}

	if passw := confService.GetAccountPassword(); passw != testData.password {
		t.Errorf("expected password to be %q got %q", testData.password, passw)
	}

	if cID := confService.GetAccountCompanyId(); cID != testData.companyId {
		t.Errorf("expected Email to be %d got %d", testData.companyId, cID)
	}
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
	if err != nil {
		t.Fatal(err)
	}
	testConfig := path.Join(cwd, "testdata/.shopware-cli.yml")

	if err := InitConfig(testConfig); err != nil {
		t.Fatalf("unexpectd err: %q", err)
	}
	if state.loadedFromEnv {
		t.Fatal("expected loadedWithEnv to be false")
	}

	confService := Config{}
	if email := confService.GetAccountEmail(); email != testData.email {
		t.Errorf("expected Email to be %q got %q", testData.email, email)
	}

	if pass := confService.GetAccountPassword(); pass != testData.password {
		t.Errorf("expected password to be %q got %q", testData.password, pass)
	}

	if cID := confService.GetAccountCompanyId(); cID != testData.companyId {
		t.Errorf("expected Email to be %d got %d", testData.companyId, cID)
	}

	if state.cfgPath != testConfig {
		t.Errorf("unexpected change to cfgFile. expected %q got %q", testConfig, state.cfgPath)
	}
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
	if err != nil {
		t.Fatal(err)
	}
	testConfig := path.Join(cwd, "testdata/write-test.yml")
	configBackup, err := os.ReadFile(testConfig)
	if err != nil {
		t.Fatalf("could not open fixture %q %s", testConfig, err)
	}
	defer func() {
		if err := os.WriteFile(testConfig, configBackup, os.ModePerm); err != nil {
			t.Error(err)
		}
	}()

	if err := InitConfig(testConfig); err != nil {
		t.Fatalf("unexpectd err: %q", err)
	}

	configService := Config{}

	if err := configService.SetAccountEmail(testData.email); err != nil {
		t.Errorf("unexpected error %s", err)
	}
	if err := configService.SetAccountPassword(testData.password); err != nil {
		t.Errorf("unexpected error %s", err)
	}
	if err := configService.SetAccountCompanyId(testData.companyId); err != nil {
		t.Errorf("unexpected error %s", err)
	}

	if err := SaveConfig(); err != nil {
		t.Fatalf("unexpected error %s", err)
	}

	if state.modified {
		t.Errorf("config state must be marked unmodified after a save")
	}

	newConfData, err := os.ReadFile(testConfig)
	if err != nil {
		t.Fatalf("unexpectd err: %q", err)
	}

	var newConf configData
	if err := yaml.Unmarshal(newConfData, &newConf); err != nil {
		t.Fatalf("encountered an error reading new config err: %q config: %q", err, string(newConfData))
	}

	if newConf.Account.Email != testData.email {
		t.Errorf("expected Email to be %q got %q", testData.email, newConf.Account.Email)
	}

	if newConf.Account.Password != testData.password {
		t.Errorf("expected password to be %q got %q", testData.password, newConf.Account.Password)
	}

	if newConf.Account.Company != testData.companyId {
		t.Errorf("expected Email to be %d got %d", testData.companyId, newConf.Account.Company)
	}
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

	if err := InitConfig(""); err != nil {
		t.Fatalf("unexpectd err: %q", err)
	}
	if !state.loadedFromEnv {
		t.Fatal("expected loadedWithEnv to be true")
	}

	confService := Config{}
	if confService.SetAccountEmail("test@foo.com") == nil {
		t.Error("expected an error when trying to write env config")
	}
	if confService.SetAccountPassword("S3CR3TF4RT3St") == nil {
		t.Error("expected an error when trying to write env config")
	}
	if confService.SetAccountCompanyId(111) == nil {
		t.Error("expected an error when trying to write env config")
	}
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
