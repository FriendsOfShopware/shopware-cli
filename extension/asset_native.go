package extension

import (
	"archive/tar"
	"compress/gzip"
	_ "embed"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
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
				content, err := ioutil.ReadFile(args.Path)
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

	tarFile, err := http.Get(fmt.Sprintf("https://github.com/sass/dart-sass-embedded/releases/download/1.51.0/sass_embedded-1.51.0-%s-%s.tar.gz", osType, arch))
	if err != nil {
		return "", errors.Wrap(err, "cannot download dart-sass")
	}

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
	OnRebuild func()
}

type CompileAdminExtensionResult struct {
	Name       string
	Entrypoint string
	JsFile     string
	CssFile    string
}

type CompileAdminExtensionOptions struct {
	ProductionMode bool
	WatchMode      *WatchMode
}

func CompileAdminExtension(ext Extension, options CompileAdminExtensionOptions) (*CompileAdminExtensionResult, error) {
	npmFile := filepath.Join(ext.GetPath(), "src/Resources/app/administration/package.json")

	if _, err := os.Stat(npmFile); err == nil {
		if err := npmInstall(filepath.Dir(npmFile)); err != nil {
			return nil, err
		}
	}

	entryPoint := filepath.Join(ext.GetPath(), "src/Resources/app/administration/src/main.js")

	if _, err := os.Stat(entryPoint); os.IsNotExist(err) {
		return nil, fmt.Errorf("cannot find entrypoint at %s", entryPoint)
	}

	name, err := ext.GetName()
	if err != nil {
		return nil, err
	}

	technicalName := strings.ReplaceAll(ToSnakeCase(name), "_", "-")
	jsFile := filepath.Join(ext.GetPath(), "src/Resources/public/administration/js", technicalName+".js")
	cssFile := filepath.Join(ext.GetPath(), "src/Resources/public/administration/css", technicalName+".css")

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

	if options.WatchMode != nil {
		bundlerOptions.Watch = &api.WatchMode{
			OnRebuild: func(br api.BuildResult) {
				if err := writeBundlerResultToDisk(br, jsFile, cssFile); err != nil {
					log.Error(err)
				}

				options.WatchMode.OnRebuild()
			},
		}
	}

	result := api.Build(bundlerOptions)

	if len(result.Errors) > 0 {
		return nil, fmt.Errorf("initial compile failed")
	}

	if err := writeBundlerResultToDisk(result, jsFile, cssFile); err != nil {
		return nil, err
	}

	compileResult := CompileAdminExtensionResult{
		Name:       name,
		Entrypoint: entryPoint,
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
