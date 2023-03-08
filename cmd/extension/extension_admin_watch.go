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

	"github.com/FriendsOfShopware/shopware-cli/esbuild"
	"github.com/FriendsOfShopware/shopware-cli/extension"
	"github.com/FriendsOfShopware/shopware-cli/logging"
	"github.com/NYTimes/gziphandler"
	"github.com/evanw/esbuild/pkg/api"
	"github.com/spf13/cobra"
	"github.com/vulcand/oxy/v2/forward"
)

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
	RunE: func(cmd *cobra.Command, args []string) error {
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

		name, _ := ext.GetName()

		options := esbuild.NewAssetCompileOptionsAdmin(name, ext.GetPath(), ext.GetType())
		options.ProductionMode = false

		esbuildContext, esBuildError := esbuild.Context(options, cmd.Context())

		if esBuildError != nil && len(esBuildError.Errors) > 0 {
			return err
		}

		if err := esbuildContext.Watch(api.WatchOptions{}); err != nil {
			return err
		}

		esbuildServer, err := esbuildContext.Serve(api.ServeOptions{
			Host: "127.0.0.1",
		})

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

		fwd := forward.New(true)

		redirect := http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			logging.FromContext(cmd.Context()).Debugf("Got request %s %s", req.Method, req.URL.Path)

			// Our custom live reload script
			if req.URL.Path == "/__internal-admin-proxy/live-reload.js" {
				w.Header().Set("content-type", "application/javascript")
				_, _ = w.Write(liveReloadJS)

				return
			}

			// Serve the local static folder to the cdn url
			assetPrefix := fmt.Sprintf(targetShopUrl.Path+"/bundles/%s/static/", strings.ToLower(name))
			if strings.HasPrefix(req.URL.Path, assetPrefix) {
				newFilePath := strings.TrimPrefix(req.URL.Path, assetPrefix)

				expectedLocation := filepath.Join(filepath.Dir(filepath.Dir(filepath.Join(ext.GetPath(), "Resources/app/administration/src"))), "static", newFilePath)

				http.ServeFile(w, req, expectedLocation)
				return
			}

			// Modify admin url index page to load anything from our watcher
			if req.URL.Path == targetShopUrl.Path+"/admin" {
				resp, err := http.Get(fmt.Sprintf("%s/admin", args[1]))

				if err != nil {
					logging.FromContext(cmd.Context()).Errorf("proxy failed %v", err)
					w.WriteHeader(http.StatusInternalServerError)
					return
				}

				body, err := io.ReadAll(resp.Body)

				if err != nil {
					logging.FromContext(cmd.Context()).Errorf("proxy reading failed %v", err)
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
						logging.FromContext(cmd.Context()).Infof("cannot parse url: %s, err: %s", s, err.Error())
						return org
					}

					if parsedUrl.Host != targetShopUrl.Host {
						return org
					}

					parsedUrl.Host = browserUrl.Host
					parsedUrl.Scheme = browserUrl.Scheme

					return firstPart + parsedUrl.String() + "\""
				})

				w.Header().Set("content-type", "text/html")
				if _, err := w.Write([]byte(bodyStr)); err != nil {
					logging.FromContext(cmd.Context()).Error(err)
				}
				logging.FromContext(cmd.Context()).Debugf("Served modified admin")
				return
			}

			// Inject our custom extension JS
			if req.URL.Path == targetShopUrl.Path+"/api/_info/config" {
				logging.FromContext(cmd.Context()).Debugf("intercept plugins call")

				proxyReq, _ := http.NewRequest("GET", targetShopUrl.Scheme+"://"+targetShopUrl.Host+req.URL.Path, nil)

				proxyReq.Header.Set("Authorization", req.Header.Get("Authorization"))

				resp, err := http.DefaultClient.Do(proxyReq)

				if err != nil {
					logging.FromContext(cmd.Context()).Errorf("proxy failed %v", err)
					w.WriteHeader(http.StatusInternalServerError)
					return
				}

				body, err := io.ReadAll(resp.Body)

				if err != nil {
					logging.FromContext(cmd.Context()).Errorf("proxy reading failed %v", err)
					w.WriteHeader(http.StatusInternalServerError)
					return
				}

				var bundleInfo adminBundlesInfo
				if err := json.Unmarshal(body, &bundleInfo); err != nil {
					logging.FromContext(cmd.Context()).Errorf("could not decode bundle info %v", err)
					w.WriteHeader(http.StatusInternalServerError)
					return
				}

				if bundleInfo.Bundles == nil {
					logging.FromContext(cmd.Context()).Errorf("cannot inject bundles. got invalid response %s", body)
					w.WriteHeader(http.StatusInternalServerError)
					return
				}

				for name, bundle := range bundleInfo.Bundles {
					newCss := []string{}

					for _, assetUrl := range bundle.Css {
						parsedUrl, _ := url.Parse(assetUrl)

						if parsedUrl.Host == targetShopUrl.Host {
							parsedUrl.Host = browserUrl.Host
							parsedUrl.Scheme = browserUrl.Scheme
						}

						newCss = append(newCss, parsedUrl.String())
					}

					newJS := []string{}

					for _, assetUrl := range bundle.Js {
						parsedUrl, _ := url.Parse(assetUrl)
						if parsedUrl.Host == targetShopUrl.Host {
							parsedUrl.Host = browserUrl.Host
							parsedUrl.Scheme = browserUrl.Scheme
						}

						newJS = append(newJS, parsedUrl.String())
					}

					bundleInfo.Bundles[name] = adminBundlesInfoAsset{Css: newCss, Js: newJS}
				}

				bundleInfo.Bundles[name] = adminBundlesInfoAsset{Css: []string{browserUrl.String() + "/extension.css"}, Js: []string{browserUrl.String() + "/extension.js"}}
				bundleInfo.Bundles["live-reload"] = adminBundlesInfoAsset{Css: []string{}, Js: []string{browserUrl.String() + "/__internal-admin-proxy/live-reload.js"}}

				newJson, _ := json.Marshal(bundleInfo)

				w.Header().Set("content-type", "application/json")
				if _, err := w.Write(newJson); err != nil {
					logging.FromContext(cmd.Context()).Error(err)
				}

				return
			}

			if req.URL.Path == "/extension.css" || req.URL.Path == "/extension.js" || req.URL.Path == "/esbuild" {
				req.URL = &url.URL{Scheme: "http", Host: fmt.Sprintf("%s:%d", esbuildServer.Host, esbuildServer.Port), Path: req.URL.Path}
				fwd.ServeHTTP(w, req)
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
		logging.FromContext(cmd.Context()).Infof("Admin Watcher started at "+browserUrl.String()+"%s/admin", targetShopUrl.Path)
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
		EnableUrlFeature  bool `json:"enableUrlFeature"`
		AppUrlReachable   bool `json:"appUrlReachable"`
		AppsRequireAppUrl bool `json:"appsRequireAppUrl"`
	} `json:"settings"`
}

type adminBundlesInfoAsset struct {
	Css []string `json:"css"`
	Js  []string `json:"js"`
}
