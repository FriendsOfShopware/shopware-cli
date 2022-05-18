package extension

import (
	"archive/zip"
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/FriendsOfShopware/shopware-cli/version"
)

const (
	TypePlatformApp    = "app"
	TypePlatformPlugin = "plugin"
)

func GetExtensionByFolder(path string) (Extension, error) {
	if _, err := os.Stat(fmt.Sprintf("%s/plugin.xml", path)); err == nil {
		return nil, fmt.Errorf("Shopware 5 is not supported. Please use https://github.com/FriendsOfShopware/FroshPluginUploader instead")
	}

	if _, err := os.Stat(fmt.Sprintf("%s/manifest.xml", path)); err == nil {
		return newApp(path)
	}

	if _, err := os.Stat(fmt.Sprintf("%s/composer.json", path)); err != nil {
		return nil, fmt.Errorf("unknown extension type")
	}

	return newPlatformPlugin(path)
}

func GetExtensionByZip(filePath string) (Extension, error) {
	dir, err := ioutil.TempDir("", "extension")
	if err != nil {
		return nil, err
	}

	content, err := ioutil.ReadFile(filePath)

	if err != nil {
		return nil, err
	}

	file, err := zip.NewReader(bytes.NewReader(content), int64(len(content)))

	if err != nil {
		return nil, err
	}

	err = Unzip(file, dir)

	if err != nil {
		return nil, err
	}

	extName := strings.Split(file.File[0].Name, "/")[0]
	return GetExtensionByFolder(fmt.Sprintf("%s/%s", dir, extName))
}

type extensionTranslated struct {
	German  string `json:"german"`
	English string `json:"english"`
}

type extensionMetadata struct {
	Label       extensionTranslated
	Description extensionTranslated
}

type Extension interface {
	GetName() (string, error)
	GetVersion() (*version.Version, error)
	GetLicense() (string, error)
	GetShopwareVersionConstraint() (*version.Constraints, error)
	GetType() string
	GetPath() string
	GetChangelog() (*extensionTranslated, error)
	GetMetaData() *extensionMetadata
	Validate(context *validationContext)
}
