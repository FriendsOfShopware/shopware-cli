package extension

import (
	"fmt"
	"path/filepath"

	"github.com/FriendsOfShopware/shopware-cli/extension"
	"github.com/FriendsOfShopware/shopware-cli/logging"

	"github.com/pkg/errors"

	"github.com/spf13/cobra"

	"github.com/FriendsOfShopware/shopware-cli/extension"
)

var extensionPrepareCmd = &cobra.Command{
	Use:   "prepare [path]",
	Short: "Prepare a extension for zipping",
	Args:  cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		path, err := filepath.Abs(args[0])
		if err != nil {
			return fmt.Errorf("path not found: %w", err)
		}

		ext, err := extension.GetExtensionByFolder(path)
		if err != nil {
			return fmt.Errorf("detect extension type: %w", err)
		}

		extCfg, err := extension.ReadExtensionConfig(ext.GetPath())
		if err != nil {
			logging.FromContext(cmd.Context()).Warnf("error reading config: %v", err)
		}

		err = extension.PrepareFolderForZipping(cmd.Context(), path+"/", ext, extCfg)
		if err != nil {
			return fmt.Errorf("prepare zip: %w", err)
		}

		return nil
	},
}

func init() {
	extensionRootCmd.AddCommand(extensionPrepareCmd)
}
