package extension

import (
	"os"
	"path"
	"testing"

	"github.com/stretchr/testify/assert"
)

const testAppManifest = `<?xml version="1.0" encoding="UTF-8"?>
<manifest xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance" xsi:noNamespaceSchemaLocation="https://raw.githubusercontent.com/shopware/platform/trunk/src/Core/Framework/App/Manifest/Schema/manifest-2.0.xsd">
	<meta>
		<name>MyExampleApp</name>
		<label>Label</label>
		<label lang="de-DE">Name</label>
		<description>A description</description>
		<description lang="de-DE">Eine Beschreibung</description>
		<author>Your Company Ltd.</author>
		<copyright>(c) by Your Company Ltd.</copyright>
		<version>1.0.0</version>
		<license>MIT</license>
	</meta>
</manifest>`

const testAppManifestIcon = `<?xml version="1.0" encoding="UTF-8"?>
<manifest xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance" xsi:noNamespaceSchemaLocation="https://raw.githubusercontent.com/shopware/platform/trunk/src/Core/Framework/App/Manifest/Schema/manifest-2.0.xsd">
	<meta>
		<name>MyExampleApp</name>
		<label>Label</label>
		<label lang="de-DE">Name</label>
		<description>A description</description>
		<description lang="de-DE">Eine Beschreibung</description>
		<author>Your Company Ltd.</author>
		<copyright>(c) by Your Company Ltd.</copyright>
		<version>1.0.0</version>
		<license>MIT</license>
		<icon>app.png</icon>
	</meta>
</manifest>`

func TestIconNotExists(t *testing.T) {
	appPath := t.TempDir()

	os.WriteFile(path.Join(appPath, "manifest.xml"), []byte(testAppManifest), os.ModePerm)

	app, err := newApp(appPath)

	assert.NoError(t, err)

	assert.Equal(t, "MyExampleApp", app.manifest.Meta.Name)
	assert.Equal(t, "", app.manifest.Meta.Icon)

	ctx := newValidationContext(app)
	app.Validate(ctx)

	assert.Equal(t, 1, len(ctx.errors))
	assert.Equal(t, "Cannot find app icon at Resources/config/plugin.png", ctx.errors[0])
}

func TestIconExistsDefaultsPath(t *testing.T) {
	appPath := t.TempDir()

	os.MkdirAll(path.Join(appPath, "Resources/config"), os.ModePerm)

	os.WriteFile(path.Join(appPath, "manifest.xml"), []byte(testAppManifest), os.ModePerm)
	os.WriteFile(path.Join(appPath, "Resources/config/plugin.png"), []byte("test"), os.ModePerm)

	app, err := newApp(appPath)

	assert.NoError(t, err)

	assert.Equal(t, "MyExampleApp", app.manifest.Meta.Name)
	assert.Equal(t, "", app.manifest.Meta.Icon)

	ctx := newValidationContext(app)
	app.Validate(ctx)

	assert.Equal(t, 0, len(ctx.errors))
}

func TestIconExistsDifferentPath(t *testing.T) {
	appPath := t.TempDir()

	os.WriteFile(path.Join(appPath, "manifest.xml"), []byte(testAppManifestIcon), os.ModePerm)
	os.WriteFile(path.Join(appPath, "app.png"), []byte("test"), os.ModePerm)

	app, err := newApp(appPath)

	assert.NoError(t, err)

	assert.Equal(t, "MyExampleApp", app.manifest.Meta.Name)
	assert.Equal(t, "app.png", app.manifest.Meta.Icon)

	ctx := newValidationContext(app)
	app.Validate(ctx)

	assert.Equal(t, 0, len(ctx.errors))
}
