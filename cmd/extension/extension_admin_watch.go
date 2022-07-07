package extension

import (
	"archive/tar"
	"compress/gzip"
	_ "embed"
	"encoding/json"
	"fmt"
	"github.com/FriendsOfShopware/shopware-cli/extension"
	"github.com/bep/godartsass"
	"github.com/evanw/esbuild/pkg/api"
	"github.com/fsnotify/fsnotify"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/vulcand/oxy/forward"
	"github.com/vulcand/oxy/testutils"
	"gopkg.in/antage/eventsource.v1"
	"io"
	"io/fs"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"
)

var es eventsource.EventSource
var hostRegExp = regexp.MustCompile(`(?m)host:\s'.*,`)
var portRegExp = regexp.MustCompile(`(?m)port:\s.*,`)
var schemeRegExp = regexp.MustCompile(`(?m)scheme:\s.*,`)
var schemeAndHttpHostRegExp = regexp.MustCompile(`(?m)schemeAndHttpHost:\s.*,`)
var uriRegExp = regexp.MustCompile(`(?m)uri:\s.*,`)
var assetPathRegExp = regexp.MustCompile(`(?m)assetPath:\s.*`)

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
		folder := fmt.Sprintf("%s/%s", cacheDir, filepath.Dir(name))
		file := fmt.Sprintf("%s/%s", cacheDir, name)

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

var extensionAdminWatchCmd = &cobra.Command{
	Use:   "admin-watch [path] [host]",
	Short: "Builds assets for extensions",
	Args:  cobra.ExactArgs(2),
	RunE: func(_ *cobra.Command, args []string) error {
		log.Infof("!!!This command is ALPHA and does not support any features of the actual Shopware watcher!!!")

		ext, err := extension.GetExtensionByFolder(args[0])

		if err != nil {
			return err
		}

		entryPoint := fmt.Sprintf("%s/src/Resources/app/administration/src/main.js", ext.GetPath())

		if _, err := os.Stat(entryPoint); os.IsNotExist(err) {
			return fmt.Errorf("cannot find entrypoint at %s", entryPoint)
		}

		name, err := ext.GetName()
		if err != nil {
			return err
		}

		technicalName := strings.ReplaceAll(extension.ToSnakeCase(name), "_", "-")
		jsFile := fmt.Sprintf("%s/src/Resources/public/administration/js/%s.js", ext.GetPath(), technicalName)
		cssFile := fmt.Sprintf("%s/src/Resources/public/administration/css/%s.css", ext.GetPath(), technicalName)

		if err := compileExtension(entryPoint, jsFile, cssFile); err != nil {
			return err
		}

		go setupWatcher(filepath.Dir(entryPoint), entryPoint, jsFile, cssFile)

		fwd, _ := forward.New()
		es = eventsource.New(nil, func(request *http.Request) [][]byte {
			return [][]byte{[]byte("Access-Control-Allow-Origin: http://localhost:8080")}
		})

		redirect := http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			log.Debugf("Got request %s %s", req.Method, req.URL.Path)

			if strings.HasPrefix(req.URL.Path, "/events") {
				es.ServeHTTP(w, req)
				return
			}

			assetPrefix := fmt.Sprintf("/bundles/%s/static/", strings.ToLower(name))
			if strings.HasPrefix(req.URL.Path, assetPrefix) {
				newFilePath := strings.TrimPrefix(req.URL.Path, assetPrefix)

				expectedLocation := filepath.Join(filepath.Dir(filepath.Dir(entryPoint)), "static", newFilePath)

				http.ServeFile(w, req, expectedLocation)
				return
			}

			if req.URL.Path == "/admin" {
				resp, err := http.Get(fmt.Sprintf("%s/admin", args[1]))

				if err != nil {
					log.Errorf("proxy failed %v", err)
					w.WriteHeader(http.StatusInternalServerError)
					return
				}

				body, err := ioutil.ReadAll(resp.Body)

				if err != nil {
					log.Errorf("proxy reading failed %v", err)
					w.WriteHeader(http.StatusInternalServerError)
					return
				}

				bodyStr := string(body)

				bodyStr = hostRegExp.ReplaceAllString(bodyStr, "host: 'localhost',")
				bodyStr = portRegExp.ReplaceAllString(bodyStr, "port: 8080,")
				bodyStr = schemeRegExp.ReplaceAllString(bodyStr, "scheme: 'http',")
				bodyStr = schemeAndHttpHostRegExp.ReplaceAllString(bodyStr, "schemeAndHttpHost: 'http://localhost:8080',")
				bodyStr = uriRegExp.ReplaceAllString(bodyStr, "uri: 'http://localhost:8080/admin',")
				bodyStr = assetPathRegExp.ReplaceAllString(bodyStr, "assetPath: 'http://localhost:8080'")

				w.Header().Set("content-type", "text/html")
				if _, err := w.Write([]byte(bodyStr)); err != nil {
					log.Error(err)
				}
				log.Debugf("Served modified admin")
				return
			}
			if req.URL.Path == "/api/_info/config" {
				log.Debugf("intercept plugins call")

				proxyReq, _ := http.NewRequest("GET", fmt.Sprintf("%s%s", args[1], req.URL.Path), nil)

				proxyReq.Header.Set("Authorization", req.Header.Get("Authorization"))

				resp, err := http.DefaultClient.Do(proxyReq)

				if err != nil {
					log.Errorf("proxy failed %v", err)
					w.WriteHeader(http.StatusInternalServerError)
					return
				}

				body, err := ioutil.ReadAll(resp.Body)

				if err != nil {
					log.Errorf("proxy reading failed %v", err)
					w.WriteHeader(http.StatusInternalServerError)
					return
				}

				var bundleInfo adminBundlesInfo
				if err := json.Unmarshal(body, &bundleInfo); err != nil {
					log.Errorf("could not decode bundle info %v", err)
					w.WriteHeader(http.StatusInternalServerError)
					return
				}

				bundleInfo.Bundles[name] = adminBundlesInfoAsset{Css: []string{"http://localhost:8080/extension.css"}, Js: []string{"http://localhost:8080/extension.js"}}
				bundleInfo.Bundles["live-reload"] = adminBundlesInfoAsset{Css: []string{}, Js: []string{"http://localhost:8080/live-reload.js"}}

				newJson, _ := json.Marshal(bundleInfo)

				w.Header().Set("content-type", "application/json")
				if _, err := w.Write(newJson); err != nil {
					log.Error(err)
				}

				return
			}

			if req.URL.Path == "/extension.css" {
				http.ServeFile(w, req, cssFile)
				return
			}

			if req.URL.Path == "/extension.js" {
				http.ServeFile(w, req, jsFile)
				return
			}

			if req.URL.Path == "/live-reload.js" {
				w.Header().Set("content-type", "application/json")
				_, _ = w.Write([]byte(("let eventSource = new EventSource('/events');\n\neventSource.onmessage = function (message) {\n    window.location.reload();\n}")))

				return
			}

			// let us forward this request to another server
			req.URL = testutils.ParseURI(args[1])
			fwd.ServeHTTP(w, req)
		})

		s := &http.Server{
			Addr:    ":8080",
			Handler: redirect,
		}
		log.Infof("Admin Watcher started at http://localhost:8080/admin")
		if err := s.ListenAndServe(); err != nil {
			return err
		}

		return nil
	},
}

func setupWatcher(watchDir string, entryPoint string, jsFile string, cssFile string) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer watcher.Close()

	done := make(chan bool)
	go func() {
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}

				if strings.HasSuffix(event.Name, "~") {
					continue
				}

				if stat, err := os.Stat(event.Name); err == nil && stat.IsDir() {
					err = watcher.Add(event.Name)
					if err != nil {
						log.Fatal(err)
					}

					log.Debugf("Added watch path: %s", event.Name)
				}

				es.SendEventMessage("reload", "message", "1")
				log.Infof("File %s has been changed", event.Name)
				if err := compileExtension(entryPoint, jsFile, cssFile); err != nil {
					log.Errorf("compile failed: %v", err)
				}
			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				log.Println("error:", err)
			}
		}
	}()

	log.Infof("Watching for changes in %s", watchDir)

	err = filepath.WalkDir(watchDir, func(path string, d fs.DirEntry, err error) error {
		if d.IsDir() {
			watchErr := watcher.Add(path)
			if watchErr != nil {
				log.Fatal(err)
			}
		}

		return nil
	})

	if err != nil {
		log.Error(err)
	}

	<-done
}

func init() {
	extensionRootCmd.AddCommand(extensionAdminWatchCmd)
}

type adminBundlesInfo struct {
	Version         string `json:"version"`
	VersionRevision string `json:"versionRevision"`
	AdminWorker     struct {
		EnableAdminWorker bool     `json:"enableAdminWorker"`
		Transports        []string `json:"transports"`
	} `json:"adminWorker"`
	Bundles  map[string]adminBundlesInfoAsset `json:"bundles"`
	Settings struct {
		EnableUrlFeature bool `json:"enableUrlFeature"`
	} `json:"settings"`
}

type adminBundlesInfoAsset struct {
	Css []string `json:"css"`
	Js  []string `json:"js"`
}

func compileExtension(entryPoint, jsFile, cssFile string) error {
	result := api.Build(api.BuildOptions{
		EntryPoints: []string{entryPoint},
		Outfile:     "extension.js",
		Bundle:      true,
		Write:       false,
		LogLevel:    api.LogLevelInfo,
		Plugins:     []api.Plugin{scssPlugin},
		Loader: map[string]api.Loader{
			".twig": api.LoaderText,
			".scss": api.LoaderCSS,
			".css":  api.LoaderCSS,
			".png":  api.LoaderFile,
			".jpg":  api.LoaderFile,
			".jpeg": api.LoaderFile,
			".ts":   api.LoaderTS,
		},
	})

	if len(result.Errors) > 0 {
		return fmt.Errorf("initial compile failed")
	}

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
