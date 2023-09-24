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
	logger := logging.NewLogger()

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
	_, err := CompileExtensionAsset(getTestContext(), options)

	assert.NoError(t, err)

	compiledFilePath := path.Join(dir, "Resources", "public", "administration", "js", "bla.js")
	_, err = os.Stat(compiledFilePath)
	assert.NoError(t, err)
}

func TestESBuildAdminTypeScript(t *testing.T) {
	dir := t.TempDir()

	adminDir := path.Join(dir, "Resources", "app", "administration", "src")
	_ = os.MkdirAll(adminDir, os.ModePerm)

	_ = os.WriteFile(path.Join(adminDir, "main.ts"), []byte("console.log('bla')"), os.ModePerm)

	options := NewAssetCompileOptionsAdmin("Bla", dir)
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

	options := NewAssetCompileOptionsStorefront("Bla", dir)
	_, err := CompileExtensionAsset(getTestContext(), options)

	assert.NoError(t, err)

	compiledFilePath := path.Join(dir, "Resources", "app", "storefront", "dist", "storefront", "js", "bla.js")
	_, err = os.Stat(compiledFilePath)
	assert.NoError(t, err)
}
