package esbuild

import (
	"context"
	_ "embed"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/evanw/esbuild/pkg/api"
)

type AssetCompileResult struct {
	Name       string
	Entrypoint string
	JsFile     string
	CssFile    string
}

type AssetCompileOptions struct {
	ProductionMode bool
	EntrypointDir  string
	OutputDir      string
	Name           string
	Path           string
}

func NewAssetCompileOptionsAdmin(name, path, extType string) AssetCompileOptions {
	root := "src/"

	if extType == "app" {
		root = ""
	}

	return AssetCompileOptions{
		Name:           name,
		Path:           path,
		EntrypointDir:  root + "Resources/app/administration/src",
		OutputDir:      root + "Resources/public/administration",
		ProductionMode: true,
	}
}

func NewAssetCompileOptionsStorefront(name, path, extType string) AssetCompileOptions {
	root := "src/"

	if extType == "app" {
		root = ""
	}

	return AssetCompileOptions{
		Name:           name,
		Path:           path,
		EntrypointDir:  root + "Resources/app/storefront/src",
		OutputDir:      root + "Resources/app/storefront/dist/storefront",
		ProductionMode: true,
	}
}

func getEsbuildOptions(options AssetCompileOptions, ctx context.Context) (*api.BuildOptions, error) {
	entryPoint := filepath.Join(options.Path, options.EntrypointDir, "main.js")

	if _, err := os.Stat(entryPoint); os.IsNotExist(err) {
		entryPointTS := filepath.Join(options.Path, options.EntrypointDir, "main.ts")

		if _, err := os.Stat(entryPointTS); os.IsNotExist(err) {
			return nil, fmt.Errorf("cannot find entrypoint at %s as main.js or main.ts", options.EntrypointDir)
		}

		entryPoint = entryPointTS
	}

	bundlerOptions := api.BuildOptions{
		MinifySyntax:      options.ProductionMode,
		MinifyWhitespace:  options.ProductionMode,
		MinifyIdentifiers: options.ProductionMode,
		EntryPoints:       []string{entryPoint},
		Outfile:           "extension.js",
		Bundle:            true,
		Write:             false,
		LogLevel:          api.LogLevelWarning,
		Plugins:           []api.Plugin{newScssPlugin(ctx)},
		Loader: map[string]api.Loader{
			".twig": api.LoaderText,
			".scss": api.LoaderCSS,
			".css":  api.LoaderCSS,
			".png":  api.LoaderFile,
			".jpg":  api.LoaderFile,
			".jpeg": api.LoaderFile,
			".ts":   api.LoaderTS,
		},
	}

	return &bundlerOptions, nil
}

func Context(options AssetCompileOptions, ctx context.Context) (api.BuildContext, *api.ContextError) {
	bundlerOptions, err := getEsbuildOptions(options, ctx)
	if err != nil {
		panic(err)
	}

	return api.Context(*bundlerOptions)
}

func CompileExtensionAsset(options AssetCompileOptions, ctx context.Context) (*AssetCompileResult, error) {
	technicalName := strings.ReplaceAll(ToSnakeCase(options.Name), "_", "-")
	jsFile := filepath.Join(options.Path, options.OutputDir, "js", technicalName+".js")
	cssFile := filepath.Join(options.Path, options.OutputDir, "css", technicalName+".css")

	bundlerOptions, err := getEsbuildOptions(options, ctx)
	if err != nil {
		return nil, err
	}

	result := api.Build(*bundlerOptions)

	if len(result.Errors) > 0 {
		return nil, fmt.Errorf("initial compile failed")
	}

	if err := writeBundlerResultToDisk(result, jsFile, cssFile); err != nil {
		return nil, err
	}

	compileResult := AssetCompileResult{
		Name:       options.Name,
		Entrypoint: bundlerOptions.EntryPoints[0],
		JsFile:     jsFile,
		CssFile:    cssFile,
	}

	return &compileResult, nil
}

func writeBundlerResultToDisk(result api.BuildResult, jsFile, cssFile string) error {
	for _, file := range result.OutputFiles {
		outFile := jsFile

		if strings.HasSuffix(file.Path, ".css") {
			outFile = cssFile
		}

		outFolder := filepath.Dir(outFile)

		if _, err := os.Stat(outFolder); os.IsNotExist(err) {
			if err := os.MkdirAll(outFolder, os.ModePerm); err != nil {
				return err
			}
		}

		if err := os.WriteFile(outFile, file.Contents, os.ModePerm); err != nil {
			return err
		}
	}

	return nil
}
