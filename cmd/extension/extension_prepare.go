package extension

import (
	"github.com/FriendsOfShopware/shopware-cli/extension"
	log "github.com/sirupsen/logrus"
	"path/filepath"

	"github.com/pkg/errors"

	"github.com/spf13/cobra"
)

var extensionPrepareCmd = &cobra.Command{
	Use:   "prepare [path]",
	Short: "Prepare a extension for zipping",
	Args:  cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		path, err := filepath.Abs(args[0])
		if err != nil {
			return errors.Wrap(err, "path not found")
		}

		ext, err := extension.GetExtensionByFolder(path)
		if err != nil {
			return errors.Wrap(err, "detect extension type")
		}

		extCfg, err := extension.ReadExtensionConfig(ext.GetPath())
		if err != nil {
			log.Warningf("error reading config: %v", err)
		}

		err = extension.PrepareFolderForZipping(cmd.Context(), path+"/", ext, extCfg)
		if err != nil {
			return errors.Wrap(err, "prepare zip")
		}

		return nil
	},
}

func init() {
	extensionRootCmd.AddCommand(extensionPrepareCmd)
}
