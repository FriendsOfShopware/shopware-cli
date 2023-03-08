package esbuild

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/bep/godartsass"
	"github.com/evanw/esbuild/pkg/api"

	"github.com/FriendsOfShopware/shopware-cli/logging"
)

func newScssPlugin(ctx context.Context) api.Plugin {
	return api.Plugin{
		Name: "scss",
		Setup: func(build api.PluginBuild) {
			dartSassBinary, err := downloadDartSass(ctx)

			if err != nil {
				logging.FromContext(ctx).Fatalln(err)
			}

			logging.FromContext(ctx).Infof("Using dart-sass binary %s", dartSassBinary)

			start, err := godartsass.Start(godartsass.Options{
				DartSassEmbeddedFilename: dartSassBinary,
				Timeout:                  0,
				LogEventHandler:          nil,
			})

			if err != nil {
				logging.FromContext(ctx).Fatalln(err)
			}

			build.OnLoad(api.OnLoadOptions{Filter: `\.scss`},
				func(args api.OnLoadArgs) (api.OnLoadResult, error) {
					content, err := os.ReadFile(args.Path)
					if err != nil {
						return api.OnLoadResult{}, err
					}

					execute, err := start.Execute(godartsass.Args{
						Source:          string(content),
						URL:             fmt.Sprintf("file://%s", args.Path),
						EnableSourceMap: true,
						IncludePaths: []string{
							filepath.Dir(args.Path),
						},
						ImportResolver: scssImporter{ctx: ctx},
					})

					if err != nil {
						return api.OnLoadResult{}, err
					}

					return api.OnLoadResult{
						Contents: &execute.CSS,
						Loader:   api.LoaderCSS,
					}, nil
				})
		},
	}
}

type scssImporter struct {
	ctx context.Context
}

const InternalVariablesScssPath = "file://internal//variables.scss"
const InternalMixinsScssPath = "file://internal//mixins.scss"

func (s scssImporter) CanonicalizeURL(url string) (string, error) {
	if url == "~scss/variables" {
		return InternalVariablesScssPath, nil
	}

	if url == "~scss/mixins" {
		return InternalMixinsScssPath, nil
	}

	return "", nil
}

func (s scssImporter) Load(canonicalizedURL string) (string, error) {
	if canonicalizedURL == InternalVariablesScssPath {
		return string(scssVariables), nil
	}

	if canonicalizedURL == InternalMixinsScssPath {
		return string(scssMixins), nil
	}

	logging.FromContext(s.ctx).Infof("Load: %s", canonicalizedURL)

	return "", nil
}
