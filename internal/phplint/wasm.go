package phplint

import (
	"context"
	"path"

	"github.com/tetratelabs/wazero"
	"github.com/tetratelabs/wazero/imports/wasi_snapshot_preview1"
)

func getWazeroRuntime(ctx context.Context) (wazero.Runtime, error) {
	cache, err := wazero.NewCompilationCacheWithDir(path.Join(getShopwareCliCacheDir(), "wasm", "cache"))
	if err != nil {
		return nil, err
	}

	r := wazero.NewRuntimeWithConfig(ctx, wazero.NewRuntimeConfig().WithCompilationCache(cache))

	wasi_snapshot_preview1.MustInstantiate(ctx, r)

	return r, nil
}
