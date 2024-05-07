package esbuild

import (
	"context"
	"os"
	"path"
	"testing"

	"github.com/FriendsOfShopware/shopware-cli/logging"
	"github.com/stretchr/testify/assert"
)

func getTestContext() context.Context {
	logger := logging.NewLogger(false)

	return logging.WithLogger(context.TODO(), logger)
}

func TestESBuildAdminNoEntrypoint(t *testing.T) {
	dir := t.TempDir()

	options := NewAssetCompileOptionsAdmin("Bla", dir)
	_, err := CompileExtensionAsset(getTestContext(), options)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "cannot find entrypoint")
}

func TestESBuildAdmin(t *testing.T) {
	dir := t.TempDir()

	adminDir := path.Join(dir, "Resources", "app", "administration", "src")
	_ = os.MkdirAll(adminDir, os.ModePerm)

	_ = os.WriteFile(path.Join(adminDir, "main.js"), []byte("console.log('bla')"), os.ModePerm)

	options := NewAssetCompileOptionsAdmin("Bla", dir)
	options.DisableSass = true
	_, err := CompileExtensionAsset(getTestContext(), options)

	assert.NoError(t, err)

	compiledFilePath := path.Join(dir, "Resources", "public", "administration", "js", "bla.js")
	_, err = os.Stat(compiledFilePath)
	assert.NoError(t, err)
}

func TestESBuildAdminWithSCSS(t *testing.T) {
	if os.Getenv("NIX_CC") != "" {
		t.Skip("Downloading does not work in Nix build")
	}

	dir := t.TempDir()

	adminDir := path.Join(dir, "Resources", "app", "administration", "src")
	_ = os.MkdirAll(adminDir, os.ModePerm)

	_ = os.WriteFile(path.Join(adminDir, "main.js"), []byte("import './a.scss'"), os.ModePerm)
	_ = os.WriteFile(path.Join(adminDir, "a.scss"), []byte(".a { .b { color: red } }"), os.ModePerm)

	options := NewAssetCompileOptionsAdmin("Bla", dir)
	_, err := CompileExtensionAsset(getTestContext(), options)

	assert.NoError(t, err)

	compiledFilePathJS := path.Join(dir, "Resources", "public", "administration", "js", "bla.js")
	_, err = os.Stat(compiledFilePathJS)
	assert.NoError(t, err)

	compiledFilePathCSS := path.Join(dir, "Resources", "public", "administration", "css", "bla.css")
	_, err = os.Stat(compiledFilePathCSS)
	assert.NoError(t, err)

	bytes, err := os.ReadFile(compiledFilePathCSS)
	assert.NoError(t, err)

	assert.Contains(t, string(bytes), ".a .b")
}

func TestESBuildAdminTypeScript(t *testing.T) {
	dir := t.TempDir()

	adminDir := path.Join(dir, "Resources", "app", "administration", "src")
	_ = os.MkdirAll(adminDir, os.ModePerm)

	_ = os.WriteFile(path.Join(adminDir, "main.ts"), []byte("console.log('bla')"), os.ModePerm)

	options := NewAssetCompileOptionsAdmin("Bla", dir)
	options.DisableSass = true
	result, err := CompileExtensionAsset(getTestContext(), options)

	assert.NoError(t, err)
	assert.Contains(t, result.Entrypoint, "main.ts")

	compiledFilePath := path.Join(dir, "Resources", "public", "administration", "js", "bla.js")
	_, err = os.Stat(compiledFilePath)
	assert.NoError(t, err)
}

func TestESBuildStorefront(t *testing.T) {
	dir := t.TempDir()

	storefrontDir := path.Join(dir, "Resources", "app", "storefront", "src")
	_ = os.MkdirAll(storefrontDir, os.ModePerm)

	_ = os.WriteFile(path.Join(storefrontDir, "main.js"), []byte("console.log('bla')"), os.ModePerm)

	options := NewAssetCompileOptionsStorefront("Bla", dir, false)
	_, err := CompileExtensionAsset(getTestContext(), options)

	assert.NoError(t, err)

	compiledFilePath := path.Join(dir, "Resources", "app", "storefront", "dist", "storefront", "js", "bla.js")
	_, err = os.Stat(compiledFilePath)
	assert.NoError(t, err)
}

func TestESBuildStorefrontNewLayout(t *testing.T) {
	dir := t.TempDir()

	storefrontDir := path.Join(dir, "Resources", "app", "storefront", "src")
	_ = os.MkdirAll(storefrontDir, os.ModePerm)

	_ = os.WriteFile(path.Join(storefrontDir, "main.js"), []byte("console.log('bla')"), os.ModePerm)

	options := NewAssetCompileOptionsStorefront("Bla", dir, true)
	_, err := CompileExtensionAsset(getTestContext(), options)

	assert.NoError(t, err)

	compiledFilePath := path.Join(dir, "Resources", "app", "storefront", "dist", "storefront", "js", "bla", "bla.js")
	_, err = os.Stat(compiledFilePath)
	assert.NoError(t, err)
}
