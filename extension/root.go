package extension

import (
	"fmt"
	"github.com/hashicorp/go-version"
	"os"
)

func GetExtensionByFolder(path string) (Extension, error) {
	if _, err := os.Stat(fmt.Sprintf("%s/plugin.xml", path)); err == nil {
		return nil, fmt.Errorf("shopware 5 is currently not supported")
	}

	if _, err := os.Stat(fmt.Sprintf("%s/manifest.xml", path)); err == nil {
		return nil, fmt.Errorf("apps are currently not supported")
	}

	if _, err := os.Stat(fmt.Sprintf("%s/composer.json", path)); err == nil {
		return newPlatformPlugin(path)
	}

	return nil, fmt.Errorf("cannot detect extension type")
}

type Changelog struct {
	German  string `json:"german"`
	English string `json:"english"`
}

type Extension interface {
	GetName() string
	GetVersion() string
	GetShopwareVersionConstraint() version.Constraints
	GetType() string
	GetChangelog() (*Changelog, error)
}
