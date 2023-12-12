package esbuild

import (
	"context"
	"fmt"
	"github.com/evanw/esbuild/pkg/api"
	"io"
	"net/http"
	"path/filepath"
	"regexp"
)

const RemoteLoaderName = "sw-remote-loader"

type RemoteLoaderOptions struct {
	BaseUrl  string
	Matchers map[string]RemoteLoaderReplacer
}

type RemoteLoaderReplacer struct {
	Matching *regexp.Regexp
	Replace  string
}

func newRemoteLoaderPlugin(ctx context.Context, options RemoteLoaderOptions) api.Plugin {
	return api.Plugin{
		Name: RemoteLoaderName,
		Setup: func(build api.PluginBuild) {
			for matcher, matcherOptions := range options.Matchers {
				build.OnResolve(api.OnResolveOptions{Filter: matcher}, func(args api.OnResolveArgs) (api.OnResolveResult, error) {
					path := matcherOptions.Matching.ReplaceAllString(args.Path, matcherOptions.Replace)

					return api.OnResolveResult{
						Path:      options.BaseUrl + path,
						Namespace: RemoteLoaderName,
					}, nil
				})
			}

			build.OnResolve(api.OnResolveOptions{Filter: "deepmerge", Namespace: RemoteLoaderName}, func(args api.OnResolveArgs) (api.OnResolveResult, error) {
				return api.OnResolveResult{
					Path:      "https://unpkg.com/deepmerge@4.3.1",
					Namespace: RemoteLoaderName,
				}, nil
			})

			// When our namespace is used, we load the remote file
			build.OnLoad(api.OnLoadOptions{Filter: "/.*/", Namespace: RemoteLoaderName}, func(args api.OnLoadArgs) (api.OnLoadResult, error) {
				ext := filepath.Ext(args.Path)

				// When we have a file extension, try direct load. But maybe the file has two file extensions, therefore, we need to fallback to .js/.ts
				if ext != "" {
					if content, err := fetchRemoteAsset(ctx, args.Path); err == nil {
						return api.OnLoadResult{
							Contents: &content,
							Loader:   api.LoaderTS,
						}, nil
					}
				}

				// Try to load the file with .ts and .js extension
				if content, err := fetchRemoteAsset(ctx, args.Path+".ts"); err == nil {
					return api.OnLoadResult{
						Contents: &content,
						Loader:   api.LoaderTS,
					}, nil
				}

				if content, err := fetchRemoteAsset(ctx, args.Path+".js"); err == nil {
					return api.OnLoadResult{
						Contents: &content,
						Loader:   api.LoaderTS,
					}, nil
				}

				return api.OnLoadResult{}, fmt.Errorf("file does not exists")
			})
		},
	}
}

func fetchRemoteAsset(ctx context.Context, url string) (string, error) {
	r, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return "", err
	}

	r.Header.Add("User-Agent", "Shopware CLI")

	resp, err := http.DefaultClient.Do(r)
	if err != nil {
		return "", err
	}

	if resp.StatusCode != 200 {
		return "", fmt.Errorf("url does not exists")
	}

	content, err := io.ReadAll(resp.Body)

	if err != nil {
		return "", err
	}

	if err := resp.Body.Close(); err != nil {
		return "", err
	}

	return string(content), nil
}
