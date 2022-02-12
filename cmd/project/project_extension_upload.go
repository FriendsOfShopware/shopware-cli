package project

import (
	"archive/zip"
	"bytes"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"os"
	"path/filepath"
	"shopware-cli/extension"
	"shopware-cli/shop"
)

var projectExtensionUploadCmd = &cobra.Command{
	Use:   "upload",
	Short: "Upload local extension to external shop",
	RunE: func(cmd *cobra.Command, args []string) error {
		var cfg *shop.Config
		var err error

		doLifecycleEvents, _ := cmd.PersistentFlags().GetBool("activate")

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

		if err := client.UploadExtension(cmd.Context(), &buf); err != nil {
			return err
		}

		log.Infof("Uploaded extension %s with version %s", name, version)

		if err := client.RefreshExtensions(cmd.Context()); err != nil {
			return err
		}

		log.Infof("Refreshed extension list")

		if doLifecycleEvents {
			extensions, err := client.GetAvailableExtensions(cmd.Context())

			if err != nil {
				return err
			}

			remoteExtension := extensions.GetByName(name)

			if remoteExtension.InstalledAt == nil {
				if err := client.InstallExtension(cmd.Context(), remoteExtension.Type, remoteExtension.Name); err != nil {
					return err
				}

				log.Infof("Installed %s", name)
			}

			if !remoteExtension.Active {
				if err := client.ActivateExtension(cmd.Context(), remoteExtension.Type, remoteExtension.Name); err != nil {
					return err
				}

				log.Infof("Activated %s", name)
			}

			if remoteExtension.IsUpdateAble() {
				if err := client.UpdateExtension(cmd.Context(), remoteExtension.Type, remoteExtension.Name); err != nil {
					return err
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

func init() {
	projectExtensionCmd.AddCommand(projectExtensionUploadCmd)
	projectExtensionUploadCmd.PersistentFlags().Bool("activate", false, "Installs, Activates, Updates the extension")
}
