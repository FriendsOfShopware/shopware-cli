package project

import (
	"archive/zip"
	"bytes"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"shopware-cli/extension"
	"shopware-cli/shop"
	"shopware-cli/version"
	"strings"

	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var projectExtensionUploadCmd = &cobra.Command{
	Use:   "upload",
	Short: "Upload local extension to external shop",
	RunE: func(cmd *cobra.Command, args []string) error {
		var cfg *shop.Config
		var err error

		doLifecycleEvents, _ := cmd.PersistentFlags().GetBool("activate")
		increaseVersionBeforeUpload, _ := cmd.PersistentFlags().GetBool("increase-version")

		path, err := filepath.Abs(args[0])

		if err != nil {
			return errors.Wrap(err, "cannot find path")
		}

		stat, err := os.Stat(path)

		if err != nil {
			return errors.Wrap(err, "cannot find path")
		}

		var ext extension.Extension

		if stat.IsDir() {
			ext, err = extension.GetExtensionByFolder(path)
		} else {
			ext, err = extension.GetExtensionByZip(path)
		}

		if err != nil {
			return err
		}

		if increaseVersionBeforeUpload {
			if err := increaseExtensionVersion(ext); err != nil {
				return err
			}

			ext, err = extension.GetExtensionByFolder(ext.GetPath())

			if err != nil {
				return err
			}
		}

		if cfg, err = shop.ReadConfig(projectConfigPath); err != nil {
			return err
		}

		client, err := shop.NewShopClient(cmd.Context(), cfg, nil)
		if err != nil {
			return err
		}

		name, err := ext.GetName()

		if err != nil {
			return err
		}

		version, err := ext.GetVersion()

		if err != nil {
			return err
		}

		var buf bytes.Buffer
		w := zip.NewWriter(&buf)
		extension.AddZipFiles(w, ext.GetPath()+"/", name+"/")

		if err := w.Close(); err != nil {
			return err
		}

		shopInfo, err := client.Info(cmd.Context())

		if err != nil {
			return errors.Wrap(err, "cannot get shop info")
		}

		extensions, err := client.GetAvailableExtensions(cmd.Context())

		if err != nil {
			return err
		}

		if !shopInfo.IsCloudShop() || extensions.GetByName(name) == nil {
			if err := client.UploadExtension(cmd.Context(), &buf); err != nil {
				return errors.Wrap(err, "cannot upload extension")
			}

			extensions, err = client.GetAvailableExtensions(cmd.Context())

			if err != nil {
				return err
			}
		} else {
			if err := client.UploadExtensionUpdateCloud(cmd.Context(), name, &buf); err != nil {
				return errors.Wrap(err, "cannot upload extension update")
			}
		}

		log.Infof("Uploaded extension %s with version %s", name, version)

		if err := client.RefreshExtensions(cmd.Context()); err != nil {
			return errors.Wrap(err, "cannot refresh extension list")
		}

		log.Infof("Refreshed extension list")

		if doLifecycleEvents {
			remoteExtension := extensions.GetByName(name)

			if remoteExtension.InstalledAt == nil {
				if err := client.InstallExtension(cmd.Context(), remoteExtension.Type, remoteExtension.Name); err != nil {
					return errors.Wrap(err, "cannot install extension")
				}

				log.Infof("Installed %s", name)
			}

			if !remoteExtension.Active {
				if err := client.ActivateExtension(cmd.Context(), remoteExtension.Type, remoteExtension.Name); err != nil {
					return errors.Wrap(err, "cannot activate extension")
				}

				log.Infof("Activated %s", name)
			}

			if remoteExtension.IsUpdateAble() {
				if err := client.UpdateExtension(cmd.Context(), remoteExtension.Type, remoteExtension.Name); err != nil {
					return errors.Wrap(err, "cannot update extension")
				}

				log.Infof("Updated %s from %s to %s", name, remoteExtension.Version, remoteExtension.LatestVersion)
			}
		}

		if ext.GetType() == "plugin" {
			if err := client.ClearCache(cmd.Context()); err != nil {
				return err
			}

			log.Infof("Cleared cache")
		}

		return nil
	},
}

func increaseExtensionVersion(ext extension.Extension) error {
	if ext.GetType() == "app" {
		manifestPath := fmt.Sprintf("%s/manifest.xml", ext.GetPath())
		file, err := os.Open(manifestPath)

		if err != nil {
			return errors.Wrap(err, "cannot read manifest file")
		}

		defer file.Close()

		var buf bytes.Buffer
		decoder := xml.NewDecoder(file)
		encoder := xml.NewEncoder(&buf)

		for {
			token, err := decoder.Token()
			if err == io.EOF {
				break
			}
			if err != nil {
				log.Printf("error getting token: %v\n", err)
				break
			}

			if v, ok := token.(xml.StartElement); ok {
				if v.Name.Local == "version" {
					var versionStr string
					if err = decoder.DecodeElement(&versionStr, &v); err != nil {
						return err
					}

					ver, err := version.NewVersion(versionStr)

					if err != nil {
						return err
					}

					ver.Increase()

					if err = encoder.EncodeElement(ver.String(), v); err != nil {
						return err
					}
					continue
				}
			}

			if err := encoder.EncodeToken(token); err != nil {
				return err
			}
		}

		// must call flush, otherwise some elements will be missing
		if err := encoder.Flush(); err != nil {
			return err
		}

		newManifest := buf.String()
		newManifest = strings.ReplaceAll(newManifest, "xmlns:_xmlns=\"xmlns\" _xmlns:xsi=", "xmlns:xsi=")
		newManifest = strings.ReplaceAll(newManifest, "xmlns:_XMLSchema-instance=\"http://www.w3.org/2001/XMLSchema-instance\" _XMLSchema-instance:noNamespaceSchemaLocation=", "xsi:noNamespaceSchemaLocation=")

		if err := ioutil.WriteFile(manifestPath, []byte(newManifest), os.ModePerm); err != nil {
			return err
		}

		return nil
	}

	composerJsonPath := fmt.Sprintf("%s/composer.json", ext.GetPath())

	composerJsonContent, err := ioutil.ReadFile(composerJsonPath)

	if err != nil {
		return err
	}

	var composerJson map[string]interface{}

	if err := json.Unmarshal(composerJsonContent, &composerJson); err != nil {
		return err
	}

	versionStr, ok := composerJson["version"].(string)

	if !ok {
		return nil
	}

	ver, err := version.NewVersion(versionStr)

	if err != nil {
		return err
	}

	ver.Increase()

	composerJson["version"] = ver.String()

	composerJsonContent, err = json.Marshal(composerJson)

	if err != nil {
		return err
	}

	if err := ioutil.WriteFile(composerJsonPath, composerJsonContent, os.ModePerm); err != nil {
		return err
	}

	return nil
}

func init() {
	projectExtensionCmd.AddCommand(projectExtensionUploadCmd)
	projectExtensionUploadCmd.PersistentFlags().Bool("activate", false, "Installs, Activates, Updates the extension")
	projectExtensionUploadCmd.PersistentFlags().Bool("increase-version", false, "Increases extension version before uploading")
}
