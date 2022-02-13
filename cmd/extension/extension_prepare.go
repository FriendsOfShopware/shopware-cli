package extension

import (
	"path/filepath"
	"shopware-cli/extension"

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

		err = extension.PrepareFolderForZipping(cmd.Context(), path+"/", ext)
		if err != nil {
			return errors.Wrap(err, "prepare zip")
		}

		return nil
	},
}

func init() {
	extensionRootCmd.AddCommand(extensionPrepareCmd)
}
