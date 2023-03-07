package extension

import (
	"io"
	"os"
	"path"
	"testing"

	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func init() {
	log.SetOutput(io.Discard)
}

func TestGenerateConfigForPlugin(t *testing.T) {
	dir := t.TempDir()

	plugin := PlatformPlugin{
		path: dir,
		composer: platformComposerJson{
			Extra: platformComposerJsonExtra{
				ShopwarePluginClass: "FroshTools\\FroshTools",
			},
		},
	}

	assert.NoError(t, os.MkdirAll(path.Join(dir, "src", "Resources", "app", "administration", "src"), os.ModePerm))
	assert.NoError(t, os.WriteFile(path.Join(dir, "src", "Resources", "app", "administration", "src", "main.js"), []byte("test"), os.ModePerm))
	assert.NoError(t, os.MkdirAll(path.Join(dir, "src", "Resources", "app", "storefront", "src"), os.ModePerm))
	assert.NoError(t, os.WriteFile(path.Join(dir, "src", "Resources", "app", "storefront", "src", "main.js"), []byte("test"), os.ModePerm))

	config := buildAssetConfigFromExtensions([]Extension{plugin}, "")

	assert.True(t, config.Has("FroshTools"))
	assert.True(t, config.RequiresAdminBuild())
	assert.True(t, config.RequiresStorefrontBuild())
	assert.Equal(t, "frosh-tools", config["FroshTools"].TechnicalName)
	assert.Equal(t, "Resources/app/administration/src/main.js", *config["FroshTools"].Administration.EntryFilePath)
	assert.Equal(t, "Resources/app/storefront/src/main.js", *config["FroshTools"].Storefront.EntryFilePath)
	assert.Nil(t, config["FroshTools"].Administration.Webpack)
	assert.Nil(t, config["FroshTools"].Storefront.Webpack)
	assert.Equal(t, "Resources/app/administration/src", config["FroshTools"].Administration.Path)
	assert.Equal(t, "Resources/app/storefront/src", config["FroshTools"].Storefront.Path)
}

func TestGenerateConfigForPluginInTypeScript(t *testing.T) {
	dir := t.TempDir()

	plugin := PlatformPlugin{
		path: dir,
		composer: platformComposerJson{
			Extra: platformComposerJsonExtra{
				ShopwarePluginClass: "FroshTools\\FroshTools",
			},
		},
	}

	assert.NoError(t, os.MkdirAll(path.Join(dir, "src", "Resources", "app", "administration", "src"), os.ModePerm))
	assert.NoError(t, os.MkdirAll(path.Join(dir, "src", "Resources", "app", "administration", "build"), os.ModePerm))

	assert.NoError(t, os.MkdirAll(path.Join(dir, "src", "Resources", "app", "storefront", "src"), os.ModePerm))
	assert.NoError(t, os.MkdirAll(path.Join(dir, "src", "Resources", "app", "storefront", "build"), os.ModePerm))

	assert.NoError(t, os.WriteFile(path.Join(dir, "src", "Resources", "app", "administration", "src", "main.ts"), []byte("test"), os.ModePerm))

	assert.NoError(t, os.WriteFile(path.Join(dir, "src", "Resources", "app", "administration", "build", "webpack.config.js"), []byte("test"), os.ModePerm))

	assert.NoError(t, os.WriteFile(path.Join(dir, "src", "Resources", "app", "storefront", "src", "main.ts"), []byte("test"), os.ModePerm))
	assert.NoError(t, os.WriteFile(path.Join(dir, "src", "Resources", "app", "storefront", "build", "webpack.config.js"), []byte("test"), os.ModePerm))

	config := buildAssetConfigFromExtensions([]Extension{plugin}, "")

	assert.True(t, config.Has("FroshTools"))
	assert.True(t, config.RequiresAdminBuild())
	assert.True(t, config.RequiresStorefrontBuild())
	assert.Equal(t, "frosh-tools", config["FroshTools"].TechnicalName)
	assert.Equal(t, "Resources/app/administration/src/main.ts", *config["FroshTools"].Administration.EntryFilePath)
	assert.Equal(t, "Resources/app/storefront/src/main.ts", *config["FroshTools"].Storefront.EntryFilePath)
	assert.Equal(t, "Resources/app/administration/build/webpack.config.js", *config["FroshTools"].Administration.Webpack)
	assert.Equal(t, "Resources/app/storefront/build/webpack.config.js", *config["FroshTools"].Storefront.Webpack)
	assert.Equal(t, "Resources/app/administration/src", config["FroshTools"].Administration.Path)
	assert.Equal(t, "Resources/app/storefront/src", config["FroshTools"].Storefront.Path)
}

func TestGenerateConfigForApp(t *testing.T) {
	dir := t.TempDir()

	app := App{
		path: dir,
		manifest: appManifest{
			Meta: appManifestMeta{
				Name: "FroshApp",
			},
		},
	}

	assert.NoError(t, os.MkdirAll(path.Join(dir, "Resources", "app", "storefront", "src"), os.ModePerm))

	assert.NoError(t, os.MkdirAll(path.Join(dir, "Resources", "app", "storefront", "build"), os.ModePerm))

	assert.NoError(t, os.WriteFile(path.Join(dir, "Resources", "app", "storefront", "src", "main.ts"), []byte("test"), os.ModePerm))

	assert.NoError(t, os.WriteFile(path.Join(dir, "Resources", "app", "storefront", "build", "webpack.config.js"), []byte("test"), os.ModePerm))

	config := buildAssetConfigFromExtensions([]Extension{app}, "")

	assert.True(t, config.Has("FroshApp"))
	assert.False(t, config.RequiresAdminBuild())
	assert.True(t, config.RequiresStorefrontBuild())

	assert.Equal(t, "frosh-app", config["FroshApp"].TechnicalName)
	assert.Nil(t, config["FroshApp"].Administration.EntryFilePath)
	assert.Equal(t, "Resources/app/storefront/src/main.ts", *config["FroshApp"].Storefront.EntryFilePath)
	assert.Nil(t, config["FroshApp"].Administration.Webpack)
	assert.Equal(t, "Resources/app/storefront/build/webpack.config.js", *config["FroshApp"].Storefront.Webpack)
	assert.Equal(t, "Resources/app/administration/src", config["FroshApp"].Administration.Path)
	assert.Equal(t, "Resources/app/storefront/src", config["FroshApp"].Storefront.Path)
}

func TestGenerateConfigAddsStorefrontAlwaysAsEntrypoint(t *testing.T) {
	config := buildAssetConfigFromExtensions([]Extension{}, "")

	assert.False(t, config.RequiresStorefrontBuild())
	assert.False(t, config.RequiresAdminBuild())
}

func TestGenerateConfigDoesNotAddExtensionWithoutConfig(t *testing.T) {
	dir := t.TempDir()

	app := App{
		path: dir,
		manifest: appManifest{
			Meta: appManifestMeta{
				Name: "FroshApp",
			},
		},
	}

	config := buildAssetConfigFromExtensions([]Extension{app}, "")

	assert.False(t, config.Has("FroshApp"))
}

func TestGenerateConfigDoesNotAddExtensionWithoutName(t *testing.T) {
	dir := t.TempDir()

	plugin := PlatformPlugin{
		path: dir,
		composer: platformComposerJson{
			Extra: platformComposerJsonExtra{
				ShopwarePluginClass: "",
			},
		},
	}

	config := buildAssetConfigFromExtensions([]Extension{plugin}, "")

	assert.Len(t, config, 1)
}

func TestGenerateConfigWithSubBundles(t *testing.T) {
	dir := t.TempDir()

	plugin := PlatformPlugin{
		path: dir,
		composer: platformComposerJson{
			Extra: platformComposerJsonExtra{
				ShopwarePluginClass: "FroshTools",
			},
		},
	}

	assert.NoError(t, os.MkdirAll(path.Join(dir, "src", "Resources", "app", "administration", "src"), os.ModePerm))

	assert.NoError(t, os.WriteFile(path.Join(dir, "src", "Resources", "app", "administration", "src", "main.ts"), []byte("test"), os.ModePerm))

	assert.NoError(t, os.MkdirAll(path.Join(dir, "src", "Foo", "Resources", "app", "administration", "src"), os.ModePerm))

	assert.NoError(t, os.WriteFile(path.Join(dir, "src", "Foo", "Resources", "app", "administration", "src", "main.ts"), []byte("test"), os.ModePerm))

	cfg := "{\"build\": {\"extraBundles\": [{\"path\": \"Foo\"}]}}"

	assert.NoError(t, os.WriteFile(path.Join(dir, ".shopware-extension.yml"), []byte(cfg), os.ModePerm))

	config := buildAssetConfigFromExtensions([]Extension{plugin}, "")

	assert.True(t, config.RequiresAdminBuild())
	assert.False(t, config.RequiresStorefrontBuild())
	assert.True(t, config.Has("FroshTools"))
	assert.True(t, config.Has("Foo"))

	assert.Equal(t, "frosh-tools", config["FroshTools"].TechnicalName)
	assert.Equal(t, "Resources/app/administration/src/main.ts", *config["FroshTools"].Administration.EntryFilePath)
	assert.Nil(t, config["FroshTools"].Administration.Webpack)
	assert.Nil(t, config["FroshTools"].Storefront.EntryFilePath)
	assert.Nil(t, config["FroshTools"].Storefront.Webpack)
	assert.Equal(t, "Resources/app/administration/src", config["FroshTools"].Administration.Path)
	assert.Equal(t, "Resources/app/storefront/src", config["FroshTools"].Storefront.Path)

	assert.Equal(t, "foo", config["Foo"].TechnicalName)
	assert.Contains(t, config["Foo"].BasePath, "src/Foo")
	assert.Equal(t, "Resources/app/administration/src/main.ts", *config["Foo"].Administration.EntryFilePath)
	assert.Nil(t, config["Foo"].Administration.Webpack)
	assert.Nil(t, config["Foo"].Storefront.EntryFilePath)
	assert.Nil(t, config["Foo"].Storefront.Webpack)
	assert.Equal(t, "Resources/app/administration/src", config["Foo"].Administration.Path)
	assert.Equal(t, "Resources/app/storefront/src", config["Foo"].Storefront.Path)
}

func TestGenerateConfigWithSubBundlesWithNameOverride(t *testing.T) {
	dir := t.TempDir()

	plugin := PlatformPlugin{
		path: dir,
		composer: platformComposerJson{
			Extra: platformComposerJsonExtra{
				ShopwarePluginClass: "FroshTools",
			},
		},
	}

	assert.NoError(t, os.MkdirAll(path.Join(dir, "src", "Resources", "app", "administration", "src"), os.ModePerm))

	assert.NoError(t, os.WriteFile(path.Join(dir, "src", "Resources", "app", "administration", "src", "main.ts"), []byte("test"), os.ModePerm))

	assert.NoError(t, os.MkdirAll(path.Join(dir, "src", "Foo", "Resources", "app", "administration", "src"), os.ModePerm))

	assert.NoError(t, os.WriteFile(path.Join(dir, "src", "Foo", "Resources", "app", "administration", "src", "main.ts"), []byte("test"), os.ModePerm))

	cfg := "{\"build\": {\"extraBundles\": [{\"path\": \"Foo\", \"name\": \"Bla\"}]}}"

	assert.NoError(t, os.WriteFile(path.Join(dir, ".shopware-extension.yml"), []byte(cfg), os.ModePerm))

	config := buildAssetConfigFromExtensions([]Extension{plugin}, "")

	assert.True(t, config.RequiresAdminBuild())
	assert.False(t, config.RequiresStorefrontBuild())

	assert.True(t, config.Has("FroshTools"))
	assert.True(t, config.Has("Bla"))

	assert.Equal(t, "frosh-tools", config["FroshTools"].TechnicalName)
	assert.Equal(t, "Resources/app/administration/src/main.ts", *config["FroshTools"].Administration.EntryFilePath)
	assert.Nil(t, config["FroshTools"].Administration.Webpack)
	assert.Nil(t, config["FroshTools"].Storefront.EntryFilePath)
	assert.Nil(t, config["FroshTools"].Storefront.Webpack)
	assert.Equal(t, "Resources/app/administration/src", config["FroshTools"].Administration.Path)
	assert.Equal(t, "Resources/app/storefront/src", config["FroshTools"].Storefront.Path)

	assert.Equal(t, "bla", config["Bla"].TechnicalName)
	assert.Contains(t, config["Bla"].BasePath, "src/Foo")
	assert.Equal(t, "Resources/app/administration/src/main.ts", *config["Bla"].Administration.EntryFilePath)
	assert.Nil(t, config["Bla"].Administration.Webpack)
	assert.Nil(t, config["Bla"].Storefront.EntryFilePath)
	assert.Nil(t, config["Bla"].Storefront.Webpack)
	assert.Equal(t, "Resources/app/administration/src", config["Bla"].Administration.Path)
	assert.Equal(t, "Resources/app/storefront/src", config["Bla"].Storefront.Path)
}
