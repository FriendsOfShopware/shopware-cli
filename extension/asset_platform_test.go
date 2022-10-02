package extension

import (
	"os"
	"path"
	"testing"
)

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
