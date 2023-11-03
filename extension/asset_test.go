package extension

import (
	"path"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConvertPlugin(t *testing.T) {
	plugin := PlatformPlugin{
		path:   t.TempDir(),
		config: &Config{},
		composer: platformComposerJson{
			Extra: platformComposerJsonExtra{
				ShopwarePluginClass: "FroshTools\\FroshTools",
			},
		},
	}

	assetSource := ConvertExtensionsToSources(getTestContext(), []Extension{plugin})

	assert.Len(t, assetSource, 1)
	froshTools := assetSource[0]

	assert.Equal(t, "FroshTools", froshTools.Name)
	assert.Equal(t, path.Join(plugin.path, "src"), froshTools.Path)
}

func TestConvertApp(t *testing.T) {
	app := App{
		path:   t.TempDir(),
		config: &Config{},
		manifest: appManifest{
			Meta: appManifestMeta{
				Name: "TestApp",
			},
		},
	}

	assetSource := ConvertExtensionsToSources(getTestContext(), []Extension{app})

	assert.Len(t, assetSource, 1)
	froshTools := assetSource[0]

	assert.Equal(t, "TestApp", froshTools.Name)
	assert.Equal(t, app.path, froshTools.Path)
}

func TestConvertExtraBundlesOfConfig(t *testing.T) {
	app := App{
		path: t.TempDir(),
		manifest: appManifest{
			Meta: appManifestMeta{
				Name: "TestApp",
			},
		},
		config: &Config{
			Build: ConfigBuild{
				ExtraBundles: []ConfigExtraBundle{
					{
						Path: "src/Fooo",
					},
				},
			},
		},
	}

	assetSource := ConvertExtensionsToSources(getTestContext(), []Extension{app})

	assert.Len(t, assetSource, 2)
	sourceOne := assetSource[0]

	assert.Equal(t, "TestApp", sourceOne.Name)
	assert.Equal(t, app.path, sourceOne.Path)

	sourceExtra := assetSource[1]

	assert.Equal(t, "Fooo", sourceExtra.Name)
	assert.Equal(t, path.Join(app.path, "src", "Fooo"), sourceExtra.Path)
}

func TestConvertExtraBundlesOfConfigWithOverride(t *testing.T) {
	app := App{
		path: t.TempDir(),
		manifest: appManifest{
			Meta: appManifestMeta{
				Name: "TestApp",
			},
		},
		config: &Config{
			Build: ConfigBuild{
				ExtraBundles: []ConfigExtraBundle{
					{
						Name: "Bla",
						Path: "src/Fooo",
					},
				},
			},
		},
	}

	assetSource := ConvertExtensionsToSources(getTestContext(), []Extension{app})

	assert.Len(t, assetSource, 2)
	sourceOne := assetSource[0]

	assert.Equal(t, "TestApp", sourceOne.Name)
	assert.Equal(t, app.path, sourceOne.Path)

	sourceExtra := assetSource[1]

	assert.Equal(t, "Bla", sourceExtra.Name)
	assert.Equal(t, path.Join(app.path, "src", "Fooo"), sourceExtra.Path)
}
