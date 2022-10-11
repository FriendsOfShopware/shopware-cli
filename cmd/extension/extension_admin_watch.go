package extension

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/FriendsOfShopware/shopware-cli/extension"
	"github.com/NYTimes/gziphandler"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/vulcand/oxy/forward"
	"gopkg.in/antage/eventsource.v1"
)

var es eventsource.EventSource
var hostRegExp = regexp.MustCompile(`(?m)host:\s'.*,`)
var portRegExp = regexp.MustCompile(`(?m)port:\s.*,`)
var schemeRegExp = regexp.MustCompile(`(?m)scheme:\s.*,`)
var schemeAndHttpHostRegExp = regexp.MustCompile(`(?m)schemeAndHttpHost:\s.*,`)
var uriRegExp = regexp.MustCompile(`(?m)uri:\s.*,`)
var assetPathRegExp = regexp.MustCompile(`(?m)assetPath:\s.*`)
var assetRegExp = regexp.MustCompile(`(?m)(src|href|content)="(https?.*\/bundles.*)"`)

//go:embed static/live-reload.js
var liveReloadJS []byte

var adminWatchListen = ""
var adminWatchURL = ""

var extensionAdminWatchCmd = &cobra.Command{
	Use:   "admin-watch [path] [host]",
	Short: "Builds assets for extensions",
	Args:  cobra.ExactArgs(2),
	RunE: func(_ *cobra.Command, args []string) error {
		ext, err := extension.GetExtensionByFolder(args[0])

		if err != nil {
			return err
		}

		listenSplit := strings.Split(adminWatchListen, ":")

		if len(listenSplit) != 2 {
			return fmt.Errorf("listen should contain a colon")
		}

		if len(adminWatchURL) == 0 {
			adminWatchURL = "http://localhost:" + listenSplit[1]
		}

		browserUrl, err := url.Parse(adminWatchURL)

		if err != nil {
			return err
		}

		es = eventsource.New(nil, func(request *http.Request) [][]byte {
			return [][]byte{[]byte("Access-Control-Allow-Origin: " + browserUrl.String())}
		})

		options := extension.NewAssetCompileOptionsAdmin()
		options.ProductionMode = false
		options.WatchMode = &extension.WatchMode{
			OnRebuild: func(onlyCssChanges bool) {
				if onlyCssChanges {
					es.SendEventMessage("reloadCss", "message", "1")
				} else {
					es.SendEventMessage("reload", "message", "1")
				}
			},
		}

		compileResult, err := extension.CompileExtensionAsset(ext, options)

		if err != nil {
			return err
		}

		targetShopUrl, err := url.Parse(strings.TrimSuffix(args[1], "/"))

		if err != nil {
			return err
		}

		browserPort := browserUrl.Port()

		if len(browserPort) == 0 {
			if browserUrl.Scheme == "https" {
				browserPort = "443"
			} else {
				browserPort = "80"
			}
		}

		fwd, _ := forward.New()

		redirect := http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			log.Debugf("Got request %s %s", req.Method, req.URL.Path)

			// Real time updates, that the browser should reload
			if strings.HasPrefix(req.URL.Path, "/__internal-admin-proxy/events") {
				es.ServeHTTP(w, req)
				return
			}

			// Our custom live reload script
			if req.URL.Path == "/__internal-admin-proxy/live-reload.js" {
				w.Header().Set("content-type", "application/json")
				_, _ = w.Write(liveReloadJS)

				return
			}

			// Serve the local static folder to the cdn url
			assetPrefix := fmt.Sprintf(targetShopUrl.Path+"/bundles/%s/static/", strings.ToLower(compileResult.Name))
			if strings.HasPrefix(req.URL.Path, assetPrefix) {
				newFilePath := strings.TrimPrefix(req.URL.Path, assetPrefix)

				expectedLocation := filepath.Join(filepath.Dir(filepath.Dir(compileResult.Entrypoint)), "static", newFilePath)

				http.ServeFile(w, req, expectedLocation)
				return
			}

			// Modify admin url index page to load anything from our watcher
			if req.URL.Path == targetShopUrl.Path+"/admin" {
				resp, err := http.Get(fmt.Sprintf("%s/admin", args[1]))

				if err != nil {
					log.Errorf("proxy failed %v", err)
					w.WriteHeader(http.StatusInternalServerError)
					return
				}

				body, err := io.ReadAll(resp.Body)

				if err != nil {
					log.Errorf("proxy reading failed %v", err)
					w.WriteHeader(http.StatusInternalServerError)
					return
				}

				bodyStr := string(body)

				bodyStr = hostRegExp.ReplaceAllString(bodyStr, "host: '"+browserUrl.Host+"',")
				bodyStr = portRegExp.ReplaceAllString(bodyStr, "port: "+browserPort+",")
				bodyStr = schemeRegExp.ReplaceAllString(bodyStr, "scheme: '"+browserUrl.Scheme+"',")
				bodyStr = schemeAndHttpHostRegExp.ReplaceAllString(bodyStr, "schemeAndHttpHost: '"+browserUrl.Scheme+"://"+browserUrl.Host+"',")
				bodyStr = uriRegExp.ReplaceAllString(bodyStr, "uri: '"+browserUrl.Scheme+"://"+browserUrl.Host+targetShopUrl.Path+"/admin',")
				bodyStr = assetPathRegExp.ReplaceAllString(bodyStr, "assetPath: '"+browserUrl.Scheme+"://"+browserUrl.Host+targetShopUrl.Path+"'")

				bodyStr = assetRegExp.ReplaceAllStringFunc(bodyStr, func(s string) string {
					firstPart := ""

					if strings.HasPrefix(s, "href=\"") {
						firstPart = "href=\""
					} else if strings.HasPrefix(s, "content=\"") {
						firstPart = "content=\""
					} else if strings.HasPrefix(s, "src=\"") {
						firstPart = "src=\""
					}

					org := s
					s = strings.TrimPrefix(s, firstPart)
					s = strings.TrimSuffix(s, "\"")

					parsedUrl, err := url.Parse(s)

					if err != nil {
						log.Infof("cannot parse url: %s, err: %s", s, err.Error())
						return org
					}

					parsedUrl.Host = browserUrl.Host
					parsedUrl.Scheme = browserUrl.Scheme

					return firstPart + parsedUrl.String() + "\""
				})

				w.Header().Set("content-type", "text/html")
				if _, err := w.Write([]byte(bodyStr)); err != nil {
					log.Error(err)
				}
				log.Debugf("Served modified admin")
				return
			}

			// Inject our custom extension JS
			if req.URL.Path == targetShopUrl.Path+"/api/_info/config" {
				log.Debugf("intercept plugins call")

				proxyReq, _ := http.NewRequest("GET", targetShopUrl.Scheme+"://"+targetShopUrl.Host+req.URL.Path, nil)

				proxyReq.Header.Set("Authorization", req.Header.Get("Authorization"))

				resp, err := http.DefaultClient.Do(proxyReq)

				if err != nil {
					log.Errorf("proxy failed %v", err)
					w.WriteHeader(http.StatusInternalServerError)
					return
				}

				body, err := io.ReadAll(resp.Body)

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

				for name, bundle := range bundleInfo.Bundles {
					newCss := []string{}

					for _, assetUrl := range bundle.Css {
						parsedUrl, _ := url.Parse(assetUrl)
						parsedUrl.Host = browserUrl.Host
						parsedUrl.Scheme = browserUrl.Scheme

						newCss = append(newCss, parsedUrl.String())
					}

					newJS := []string{}

					for _, assetUrl := range bundle.Js {
						parsedUrl, _ := url.Parse(assetUrl)
						parsedUrl.Host = browserUrl.Host
						parsedUrl.Scheme = browserUrl.Scheme

						newJS = append(newJS, parsedUrl.String())
					}

					bundleInfo.Bundles[name] = adminBundlesInfoAsset{Css: newCss, Js: newJS}
				}

				bundleInfo.Bundles[compileResult.Name] = adminBundlesInfoAsset{Css: []string{browserUrl.String() + "/extension.css"}, Js: []string{browserUrl.String() + "/extension.js"}}
				bundleInfo.Bundles["live-reload"] = adminBundlesInfoAsset{Css: []string{}, Js: []string{browserUrl.String() + "/__internal-admin-proxy/live-reload.js"}}

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

			// let us forward this request to another server
			req.URL = targetShopUrl
			fwd.ServeHTTP(w, req)
		})

		wrapper, _ := gziphandler.GzipHandlerWithOpts(gziphandler.ContentTypes([]string{"application/vnd.api+json", "application/json ", "text/html", "text/javascript", "text/css", "image/png"}))

		s := &http.Server{
			Addr:              adminWatchListen,
			Handler:           wrapper(redirect),
			ReadHeaderTimeout: time.Second,
		}
		log.Infof("Admin Watcher started at "+browserUrl.String()+"%s/admin", targetShopUrl.Path)
		if err := s.ListenAndServe(); err != nil {
			return err
		}

		return nil
	},
}

func init() {
	extensionRootCmd.AddCommand(extensionAdminWatchCmd)
	extensionAdminWatchCmd.PersistentFlags().StringVar(&adminWatchListen, "listen", ":8080", "Listen (default :8080)")
	extensionAdminWatchCmd.PersistentFlags().StringVar(&adminWatchURL, "external-url", "", "External reachable url for admin watcher. Needed for reverse proxy setups")
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
