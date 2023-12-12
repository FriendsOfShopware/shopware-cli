package esbuild

import (
	"context"
	"github.com/evanw/esbuild/pkg/api"
	"github.com/stretchr/testify/assert"
	"os"
	"regexp"
	"testing"
)

func TestLoadRegularJSFileWithImportPrefixAdmin(t *testing.T) {
	tmpDir := t.TempDir()

	// Create a file
	filePath := tmpDir + "/test.js"
	fileContent := "import {searchRankingPoint} from \"@administration/app/service/search-ranking.service\";"

	assert.NoError(t, os.WriteFile(filePath, []byte(fileContent), 0644))

	options := RemoteLoaderOptions{
		BaseUrl: "https://raw.githubusercontent.com/shopware/shopware/v6.5.7.3/src/Administration/Resources/app/administration/",
		Matchers: map[string]RemoteLoaderReplacer{
			"^@administration/": {Matching: regexp.MustCompile("^@administration/"), Replace: "src/"},
		},
	}

	build := api.BuildOptions{
		EntryPoints: []string{filePath},
		Plugins:     []api.Plugin{newRemoteLoaderPlugin(context.Background(), options)},
		Outfile:     "extension.js",
		Bundle:      true,
		Write:       false,
	}

	result := api.Build(build)
	assert.Len(t, result.Errors, 0)

	assert.Contains(t, string(result.OutputFiles[0].Contents), "HIGH_SEARCH_RANKING")
}

func TestLoadRegularJSFileAdmin(t *testing.T) {
	tmpDir := t.TempDir()

	// Create a file
	filePath := tmpDir + "/test.js"
	fileContent := "import {searchRankingPoint} from \"src/app/service/search-ranking.service\";"

	assert.NoError(t, os.WriteFile(filePath, []byte(fileContent), 0644))

	options := RemoteLoaderOptions{
		BaseUrl: "https://raw.githubusercontent.com/shopware/shopware/v6.5.7.3/src/Administration/Resources/app/administration/",
		Matchers: map[string]RemoteLoaderReplacer{
			"^src\\/": {Matching: regexp.MustCompile("^src/"), Replace: "src/"},
		},
	}

	build := api.BuildOptions{
		EntryPoints: []string{filePath},
		Plugins:     []api.Plugin{newRemoteLoaderPlugin(context.Background(), options)},
		Outfile:     "extension.js",
		Bundle:      true,
		Write:       false,
	}

	result := api.Build(build)
	assert.Len(t, result.Errors, 0)

	assert.Len(t, result.OutputFiles, 1)

	assert.Contains(t, string(result.OutputFiles[0].Contents), "HIGH_SEARCH_RANKING")
}

func TestLoadRegularJSFileStorefront(t *testing.T) {
	tmpDir := t.TempDir()

	// Create a file
	filePath := tmpDir + "/test.js"
	fileContent := "import Plugin from \"src/plugin-system/plugin.class\"; class Foo extends Plugin {}; export {Foo}"

	assert.NoError(t, os.WriteFile(filePath, []byte(fileContent), 0644))

	options := RemoteLoaderOptions{
		BaseUrl: "https://raw.githubusercontent.com/shopware/shopware/v6.5.7.3/src/Storefront/Resources/app/storefront/",
		Matchers: map[string]RemoteLoaderReplacer{
			"^src\\/": {Matching: regexp.MustCompile("^src/"), Replace: "src/"},
		},
	}

	build := api.BuildOptions{
		EntryPoints: []string{filePath},
		Plugins:     []api.Plugin{newRemoteLoaderPlugin(context.Background(), options)},
		Outfile:     "extension.js",
		Bundle:      true,
		Write:       false,
	}

	result := api.Build(build)
	assert.Len(t, result.Errors, 0)

	assert.Len(t, result.OutputFiles, 1)

	assert.Contains(t, string(result.OutputFiles[0].Contents), "There is no valid element given")
}
