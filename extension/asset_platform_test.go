package extension

import (
	"io"
	"os"
	"path"
	"testing"

	log "github.com/sirupsen/logrus"
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

	if err := os.MkdirAll(path.Join(dir, "src", "Resources", "app", "administration", "src"), os.ModePerm); err != nil {
		t.Fatal(err)
	}

	if err := os.MkdirAll(path.Join(dir, "src", "Resources", "app", "storefront", "src"), os.ModePerm); err != nil {
		t.Fatal(err)
	}

	if err := os.WriteFile(path.Join(dir, "src", "Resources", "app", "administration", "src", "main.js"), []byte("test"), os.ModePerm); err != nil {
		t.Fatal(err)
	}

	if err := os.WriteFile(path.Join(dir, "src", "Resources", "app", "storefront", "src", "main.js"), []byte("test"), os.ModePerm); err != nil {
		t.Fatal(err)
	}

	config := buildAssetConfigFromExtensions([]Extension{plugin}, "")

	if config.Has("FroshTools") == false {
		t.Error("Expected to have FroshTools")
	}

	if config.RequiresAdminBuild() == false {
		t.Error("Expected to require admin build")
	}

	if config.RequiresStorefrontBuild() == false {
		t.Error("Expected to require storefront build")
	}

	if config["FroshTools"].TechnicalName != "frosh-tools" {
		t.Error("Expected to have frosh-tools")
	}

	if *config["FroshTools"].Administration.EntryFilePath != "app/administration/src/main.js" {
		t.Error("Expected to have Administration JS")
	}

	if *config["FroshTools"].Storefront.EntryFilePath != "app/storefront/src/main.js" {
		t.Error("Expected to have Storefront JS")
	}

	if config["FroshTools"].Administration.Webpack != nil {
		t.Error("Webpack is not overriden for admin")
	}

	if config["FroshTools"].Storefront.Webpack != nil {
		t.Error("Webpack is not overriden for storefront")
	}

	if config["FroshTools"].Administration.Path != "app/administration/src" {
		t.Error("Expected to have Administration Path")
	}

	if config["FroshTools"].Storefront.Path != "app/storefront/src" {
		t.Error("Expected to have Storefront Path")
	}
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

	if err := os.MkdirAll(path.Join(dir, "src", "Resources", "app", "administration", "src"), os.ModePerm); err != nil {
		t.Fatal(err)
	}

	if err := os.MkdirAll(path.Join(dir, "src", "Resources", "app", "administration", "build"), os.ModePerm); err != nil {
		t.Fatal(err)
	}

	if err := os.MkdirAll(path.Join(dir, "src", "Resources", "app", "storefront", "src"), os.ModePerm); err != nil {
		t.Fatal(err)
	}

	if err := os.MkdirAll(path.Join(dir, "src", "Resources", "app", "storefront", "build"), os.ModePerm); err != nil {
		t.Fatal(err)
	}

	if err := os.WriteFile(path.Join(dir, "src", "Resources", "app", "administration", "src", "main.ts"), []byte("test"), os.ModePerm); err != nil {
		t.Fatal(err)
	}

	if err := os.WriteFile(path.Join(dir, "src", "Resources", "app", "administration", "build", "webpack.config.js"), []byte("test"), os.ModePerm); err != nil {
		t.Fatal(err)
	}

	if err := os.WriteFile(path.Join(dir, "src", "Resources", "app", "storefront", "src", "main.ts"), []byte("test"), os.ModePerm); err != nil {
		t.Fatal(err)
	}

	if err := os.WriteFile(path.Join(dir, "src", "Resources", "app", "storefront", "build", "webpack.config.js"), []byte("test"), os.ModePerm); err != nil {
		t.Fatal(err)
	}

	config := buildAssetConfigFromExtensions([]Extension{plugin}, "")

	if config.Has("FroshTools") == false {
		t.Error("Expected to have FroshTools")
	}

	if config.RequiresAdminBuild() == false {
		t.Error("Expected to require admin build")
	}

	if config.RequiresStorefrontBuild() == false {
		t.Error("Expected to require storefront build")
	}

	if config["FroshTools"].TechnicalName != "frosh-tools" {
		t.Error("Expected to have frosh-tools")
	}

	if *config["FroshTools"].Administration.EntryFilePath != "app/administration/src/main.ts" {
		t.Error("Expected to have Administration TS")
	}

	if *config["FroshTools"].Administration.Webpack != "app/administration/build/webpack.config.js" {
		t.Error("Expected to have Administration Webpack")
	}

	if *config["FroshTools"].Storefront.EntryFilePath != "app/storefront/src/main.ts" {
		t.Error("Expected to have Storefront TS")
	}

	if *config["FroshTools"].Storefront.Webpack != "app/storefront/build/webpack.config.js" {
		t.Error("Expected to have Storefront Webpack")
	}

	if config["FroshTools"].Administration.Path != "app/administration/src" {
		t.Error("Expected to have Administration Path")
	}

	if config["FroshTools"].Storefront.Path != "app/storefront/src" {
		t.Error("Expected to have Storefront Path")
	}
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

	if err := os.MkdirAll(path.Join(dir, "Resources", "app", "storefront", "src"), os.ModePerm); err != nil {
		t.Fatal(err)
	}

	if err := os.MkdirAll(path.Join(dir, "Resources", "app", "storefront", "src"), os.ModePerm); err != nil {
		t.Fatal(err)
	}

	if err := os.MkdirAll(path.Join(dir, "Resources", "app", "storefront", "build"), os.ModePerm); err != nil {
		t.Fatal(err)
	}

	if err := os.WriteFile(path.Join(dir, "Resources", "app", "storefront", "src", "main.ts"), []byte("test"), os.ModePerm); err != nil {
		t.Fatal(err)
	}

	if err := os.WriteFile(path.Join(dir, "Resources", "app", "storefront", "build", "webpack.config.js"), []byte("test"), os.ModePerm); err != nil {
		t.Fatal(err)
	}

	config := buildAssetConfigFromExtensions([]Extension{app}, "")

	if config.Has("FroshApp") == false {
		t.Error("Expected to have FroshApp")
	}

	if config.RequiresAdminBuild() == true {
		t.Error("Expected to not require admin build")
	}

	if config.RequiresStorefrontBuild() == false {
		t.Error("Expected to require storefront build")
	}

	if config["FroshApp"].TechnicalName != "frosh-app" {
		t.Error("Expected to have frosh-app")
	}

	if *config["FroshApp"].Storefront.EntryFilePath != "app/storefront/src/main.ts" {
		t.Error("Expected to have Storefront TS")
	}

	if *config["FroshApp"].Storefront.Webpack != "app/storefront/build/webpack.config.js" {
		t.Error("Expected to have Storefront Webpack")
	}
}

func TestGenerateConfigAddsStorefrontAlwaysAsEntrypoint(t *testing.T) {
	config := buildAssetConfigFromExtensions([]Extension{}, "")

	if config.RequiresStorefrontBuild() == true {
		t.Error("Storefront entrypoint is ignored for build")
	}

	if config.RequiresAdminBuild() == true {
		t.Error("Storefront bundle does not offer a admin entrypoint")
	}
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

	if config.Has("FroshApp") == true {
		t.Error("Expected to not have FroshApp")
	}
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

	if len(config) != 1 {
		t.Error("Expected no to add plugin")
	}
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

	if err := os.MkdirAll(path.Join(dir, "src", "Resources", "app", "administration", "src"), os.ModePerm); err != nil {
		t.Fatal(err)
	}

	if err := os.WriteFile(path.Join(dir, "src", "Resources", "app", "administration", "src", "main.ts"), []byte("test"), os.ModePerm); err != nil {
		t.Fatal(err)
	}

	if err := os.MkdirAll(path.Join(dir, "src", "Foo", "Resources", "app", "administration", "src"), os.ModePerm); err != nil {
		t.Fatal(err)
	}

	if err := os.WriteFile(path.Join(dir, "src", "Foo", "Resources", "app", "administration", "src", "main.ts"), []byte("test"), os.ModePerm); err != nil {
		t.Fatal(err)
	}

	cfg := "{\"build\": {\"extraBundles\": [{\"path\": \"Foo\"}]}}"

	if err := os.WriteFile(path.Join(dir, ".shopware-extension.yml"), []byte(cfg), os.ModePerm); err != nil {
		t.Fatal(err)
	}

	config := buildAssetConfigFromExtensions([]Extension{plugin}, "")

	if config.RequiresAdminBuild() == false {
		t.Error("Expected to require admin build")
	}

	if config.Has("FroshTools") == false {
		t.Error("Expected to have FroshTools")
	}

	if config.Has("Foo") == false {
		t.Error("Expected to have Foo")
	}

	if config["Foo"].TechnicalName != "foo" {
		t.Error("Expected to have foo")
	}

	if *config["FroshTools"].Administration.EntryFilePath != "app/administration/src/main.ts" {
		t.Error("Expected to have Admin JS")
	}

	if *config["Foo"].Administration.EntryFilePath != "app/administration/src/main.ts" {
		t.Error("Expected to have Admin JS")
	}
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

	if err := os.MkdirAll(path.Join(dir, "src", "Resources", "app", "administration", "src"), os.ModePerm); err != nil {
		t.Fatal(err)
	}

	if err := os.WriteFile(path.Join(dir, "src", "Resources", "app", "administration", "src", "main.ts"), []byte("test"), os.ModePerm); err != nil {
		t.Fatal(err)
	}

	if err := os.MkdirAll(path.Join(dir, "src", "Foo", "Resources", "app", "administration", "src"), os.ModePerm); err != nil {
		t.Fatal(err)
	}

	if err := os.WriteFile(path.Join(dir, "src", "Foo", "Resources", "app", "administration", "src", "main.ts"), []byte("test"), os.ModePerm); err != nil {
		t.Fatal(err)
	}

	cfg := "{\"build\": {\"extraBundles\": [{\"path\": \"Foo\", \"name\": \"Bla\"}]}}"

	if err := os.WriteFile(path.Join(dir, ".shopware-extension.yml"), []byte(cfg), os.ModePerm); err != nil {
		t.Fatal(err)
	}

	config := buildAssetConfigFromExtensions([]Extension{plugin}, "")

	if config.RequiresAdminBuild() == false {
		t.Error("Expected to require admin build")
	}

	if config.Has("FroshTools") == false {
		t.Error("Expected to have FroshTools")
	}

	if config.Has("Bla") == false {
		t.Error("Expected to have Bla")
	}

	if config["Bla"].TechnicalName != "bla" {
		t.Error("Expected to have bla")
	}

	if *config["FroshTools"].Administration.EntryFilePath != "app/administration/src/main.ts" {
		t.Error("Expected to have Admin JS")
	}

	if *config["Bla"].Administration.EntryFilePath != "app/administration/src/main.ts" {
		t.Error("Expected to have Admin JS")
	}
}
