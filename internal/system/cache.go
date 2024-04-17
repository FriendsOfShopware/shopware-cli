package system

import (
	"os"
	"path"
)

func GetShopwareCliCacheDir() string {
	cacheDir, _ := os.UserCacheDir()

	return path.Join(cacheDir, "shopware-cli")
}

func GetShopwareCliConfigDir() string {
	configDir, _ := os.UserConfigDir()

	return path.Join(configDir, "shopware-cli")
}
