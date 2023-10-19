package system

import (
	"os"
	"path"
)

func GetShopwareCliCacheDir() string {
	cacheDir, _ := os.UserCacheDir()

	return path.Join(cacheDir, "shopware-cli")
}
