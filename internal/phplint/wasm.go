package phplint

import (
	"context"
	"os"
	"path"

	"github.com/FriendsOfShopware/shopware-cli/internal/system"
	"github.com/tetratelabs/wazero"
	"github.com/tetratelabs/wazero/imports/wasi_snapshot_preview1"
)

func getWazeroRuntime(ctx context.Context) (wazero.Runtime, error) {
	wazeroCacheDir := path.Join(system.GetShopwareCliCacheDir(), "wasm", "cache")

	if _, err := os.Stat(wazeroCacheDir); os.IsNotExist(err) {
		if err := os.MkdirAll(wazeroCacheDir, os.ModePerm); err != nil {
			return nil, err
		}
	}

	cache, err := wazero.NewCompilationCacheWithDir(wazeroCacheDir)
	if err != nil {
		return nil, err
	}

	r := wazero.NewRuntimeWithConfig(ctx, wazero.NewRuntimeConfig().WithCompilationCache(cache))

	wasi_snapshot_preview1.MustInstantiate(ctx, r)

	return r, nil
}
