package extension

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/FriendsOfShopware/shopware-cli/extension"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/vulcand/oxy/forward"
	"github.com/vulcand/oxy/testutils"
	"gopkg.in/antage/eventsource.v1"
)

var es eventsource.EventSource
var hostRegExp = regexp.MustCompile(`(?m)host:\s'.*,`)
var portRegExp = regexp.MustCompile(`(?m)port:\s.*,`)
var schemeRegExp = regexp.MustCompile(`(?m)scheme:\s.*,`)
var schemeAndHttpHostRegExp = regexp.MustCompile(`(?m)schemeAndHttpHost:\s.*,`)
var uriRegExp = regexp.MustCompile(`(?m)uri:\s.*,`)
var assetPathRegExp = regexp.MustCompile(`(?m)assetPath:\s.*`)

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

		es = eventsource.New(nil, func(request *http.Request) [][]byte {
			return [][]byte{[]byte("Access-Control-Allow-Origin: http://localhost:8080")}
		})

		options := extension.NewAssetCompileOptionsAdmin()
		options.ProductionMode = false
		options.WatchMode = &extension.WatchMode{
			OnRebuild: func() {
				es.SendEventMessage("reload", "message", "1")
			},
		}

		compileResult, err := extension.CompileExtensionAsset(ext, options)

		if err != nil {
			return err
		}

		fwd, _ := forward.New()

		redirect := http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			log.Debugf("Got request %s %s", req.Method, req.URL.Path)

			if strings.HasPrefix(req.URL.Path, "/events") {
				es.ServeHTTP(w, req)
				return
			}

			assetPrefix := fmt.Sprintf("/bundles/%s/static/", strings.ToLower(compileResult.Name))
			if strings.HasPrefix(req.URL.Path, assetPrefix) {
				newFilePath := strings.TrimPrefix(req.URL.Path, assetPrefix)

				expectedLocation := filepath.Join(filepath.Dir(filepath.Dir(compileResult.Entrypoint)), "static", newFilePath)

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

				if bundleInfo.Bundles == nil {
					log.Errorf("cannot inject bundles. got invalid response %s", body)
					w.WriteHeader(http.StatusInternalServerError)
					return
				}

				bundleInfo.Bundles[compileResult.Name] = adminBundlesInfoAsset{Css: []string{"http://localhost:8080/extension.css"}, Js: []string{"http://localhost:8080/extension.js"}}
				bundleInfo.Bundles["live-reload"] = adminBundlesInfoAsset{Css: []string{}, Js: []string{"http://localhost:8080/live-reload.js"}}

				newJson, _ := json.Marshal(bundleInfo)

				w.Header().Set("content-type", "application/json")
				if _, err := w.Write(newJson); err != nil {
					log.Error(err)
				}

				return
			}

			if req.URL.Path == "/extension.css" {
				http.ServeFile(w, req, compileResult.CssFile)
				return
			}

			if req.URL.Path == "/extension.js" {
				http.ServeFile(w, req, compileResult.JsFile)
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
