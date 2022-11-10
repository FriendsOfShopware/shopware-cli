package extension

import (
	"archive/tar"
	"compress/gzip"
	"context"
	"crypto/md5" //nolint:gosec
	_ "embed"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	log "github.com/sirupsen/logrus"

	"github.com/bep/godartsass"
	"github.com/evanw/esbuild/pkg/api"
	"github.com/pkg/errors"
)

//go:embed static/variables.scss
var scssVariables []byte

var scssPlugin = api.Plugin{
	Name: "scss",
	Setup: func(build api.PluginBuild) {
		dartSassBinary, err := downloadDartSass()

		if err != nil {
			log.Fatalln(err)
		}

		log.Infof("Using dart-sass binary %s", dartSassBinary)

		start, err := godartsass.Start(godartsass.Options{
			DartSassEmbeddedFilename: dartSassBinary,
			Timeout:                  0,
			LogEventHandler:          nil,
		})

		if err != nil {
			log.Fatalln(err)
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
					ImportResolver: scssImporter{},
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

func downloadDartSass() (string, error) {
	if path, err := exec.LookPath("dart-sass-embedded"); err == nil {
		return path, nil
	}

	cacheDir, err := os.UserCacheDir()

	if err != nil {
		cacheDir = "/tmp"
	}

	cacheDir += "/dart-sass-embedded"

	expectedPath := fmt.Sprintf("%s/dart-sass-embedded", cacheDir)

	if _, err := os.Stat(expectedPath); err == nil {
		return expectedPath, nil
	}

	if _, err := os.Stat(filepath.Dir(expectedPath)); os.IsNotExist(err) {
		if err := os.MkdirAll(filepath.Dir(expectedPath), os.ModePerm); err != nil {
			return "", err
		}
	}

	log.Infof("Downloading dart-sass")

	osType := runtime.GOOS
	arch := runtime.GOARCH

	switch runtime.GOARCH {
	case "arm64":
		arch = "arm64"
	case "amd64":
		arch = "x64"
	case "386":
		arch = "ia32"
	}

	if osType == "darwin" {
		osType = "macos"
	}

	request, _ := http.NewRequestWithContext(context.Background(), "GET", fmt.Sprintf("https://github.com/sass/dart-sass-embedded/releases/download/1.56.1/sass_embedded-1.56.1-%s-%s.tar.gz", osType, arch), nil)

	tarFile, err := http.DefaultClient.Do(request)
	if err != nil {
		return "", errors.Wrap(err, "cannot download dart-sass")
	}

	defer tarFile.Body.Close()

	uncompressedStream, err := gzip.NewReader(tarFile.Body)
	if err != nil {
		return "", errors.Wrap(err, "cannot open gzip tar file")
	}

	tarReader := tar.NewReader(uncompressedStream)

	for {
		header, err := tarReader.Next()

		if err == io.EOF {
			break
		}

		name := strings.TrimPrefix(header.Name, "sass_embedded/")
		folder := filepath.Join(cacheDir, filepath.Dir(name))
		file := filepath.Join(cacheDir, name)

		if !strings.HasSuffix(folder, ".") {
			if _, err := os.Stat(folder); os.IsNotExist(err) {
				if err := os.MkdirAll(folder, os.ModePerm); err != nil {
					return "", err
				}
			}
		}

		outFile, err := os.Create(file)
		if err != nil {
			return "", errors.Wrap(err, "cannot create dart-sass in temp")
		}
		if _, err := io.CopyN(outFile, tarReader, header.Size); err != nil {
			return "", errors.Wrap(err, "cannot copy dart-sass in temp")
		}
		if err := outFile.Close(); err != nil {
			return "", errors.Wrap(err, "cannot close dart-sass in temp")
		}

		if err := os.Chmod(file, os.FileMode(header.Mode)); err != nil {
			return "", errors.Wrap(err, "cannot chmod dart-sass in temp")
		}
	}

	return expectedPath, nil
}

type WatchMode struct {
	OnRebuild func(bool)
}

type AssetCompileResult struct {
	Name       string
	Entrypoint string
	JsFile     string
	CssFile    string
}

type AssetCompileOptions struct {
	ProductionMode bool
	WatchMode      *WatchMode
	EntrypointDir  string
	OutputDir      string
}

func NewAssetCompileOptionsAdmin() AssetCompileOptions {
	return AssetCompileOptions{
		EntrypointDir: "src/Resources/app/administration/src",
		OutputDir:     "src/Resources/public/administration",
	}
}

func NewAssetCompileOptionsStorefront() AssetCompileOptions {
	return AssetCompileOptions{
		EntrypointDir: "src/Resources/app/storefront/src",
		OutputDir:     "src/Resources/app/storefront/dist/storefront",
	}
}

func CompileExtensionAsset(ext Extension, options AssetCompileOptions) (*AssetCompileResult, error) {
	entryPoint := filepath.Join(ext.GetPath(), options.EntrypointDir, "main.js")

	if _, err := os.Stat(entryPoint); os.IsNotExist(err) {
		entryPointTS := filepath.Join(ext.GetPath(), options.EntrypointDir, "main.ts")

		if _, err := os.Stat(entryPointTS); os.IsNotExist(err) {
			return nil, fmt.Errorf("cannot find entrypoint at %s as main.js or main.ts", options.EntrypointDir)
		}

		entryPoint = entryPointTS
	}

	name, err := ext.GetName()
	if err != nil {
		return nil, err
	}

	technicalName := strings.ReplaceAll(ToSnakeCase(name), "_", "-")
	jsFile := filepath.Join(ext.GetPath(), options.OutputDir, "js", technicalName+".js")
	cssFile := filepath.Join(ext.GetPath(), options.OutputDir, "css", technicalName+".css")

	bundlerOptions := api.BuildOptions{
		MinifySyntax:      options.ProductionMode,
		MinifyWhitespace:  options.ProductionMode,
		MinifyIdentifiers: options.ProductionMode,
		EntryPoints:       []string{entryPoint},
		Outfile:           "extension.js",
		Bundle:            true,
		Write:             false,
		LogLevel:          api.LogLevelInfo,
		Plugins:           []api.Plugin{scssPlugin},
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

	jsMD5 := ""

	if options.WatchMode != nil {
		bundlerOptions.Watch = &api.WatchMode{
			OnRebuild: func(br api.BuildResult) {
				currentMD5 := jsMD5
				if err := writeBundlerResultToDisk(br, jsFile, cssFile, &jsMD5); err != nil {
					log.Error(err)
				}

				options.WatchMode.OnRebuild(currentMD5 == jsMD5)
			},
		}
	}

	result := api.Build(bundlerOptions)

	if len(result.Errors) > 0 {
		return nil, fmt.Errorf("initial compile failed")
	}

	if err := writeBundlerResultToDisk(result, jsFile, cssFile, &jsMD5); err != nil {
		return nil, err
	}

	compileResult := AssetCompileResult{
		Name:       name,
		Entrypoint: entryPoint,
		JsFile:     jsFile,
		CssFile:    cssFile,
	}

	return &compileResult, nil
}

func writeBundlerResultToDisk(result api.BuildResult, jsFile, cssFile string, jsMD5 *string) error {
	for _, file := range result.OutputFiles {
		outFile := jsFile

		if strings.HasSuffix(file.Path, ".css") {
			outFile = cssFile
		} else {
			hash := md5.New() //nolint:gosec
			hash.Write(file.Contents)

			*jsMD5 = fmt.Sprintf("%x", hash.Sum(nil))
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

type scssImporter struct{}

const InternalScssPath = "file://internal//variables.scss"

func (s scssImporter) CanonicalizeURL(url string) (string, error) {
	if url == "~scss/variables" {
		return InternalScssPath, nil
	}

	return "", nil
}

func (s scssImporter) Load(canonicalizedURL string) (string, error) {
	if canonicalizedURL == InternalScssPath {
		return string(scssVariables), nil
	}

	log.Infof("Load: %s", canonicalizedURL)

	return "", nil
}
