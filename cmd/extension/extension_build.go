package extension

import (
	"path/filepath"
	"shopware-cli/extension"

	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"

	"github.com/spf13/cobra"
)

var extensionAssetBundleCmd = &cobra.Command{
	Use:   "build [path]",
	Short: "Builds assets for extensions",
	Args:  cobra.MinimumNArgs(1),
	RunE: func(_ *cobra.Command, args []string) error {
		validatedExtensions := make([]extension.Extension, 0)

		for _, arg := range args {
			path, err := filepath.Abs(arg)

			if err != nil {
				return errors.Wrap(err, "cannot open file")
			}

			ext, err := extension.GetExtensionByFolder(path)

			if err != nil {
				return errors.Wrap(err, "cannot open extension")
			}

			validatedExtensions = append(validatedExtensions, ext)
		}

		err := extension.BuildAssetsForExtensions("", validatedExtensions)

		if err != nil {
			return errors.Wrap(err, "cannot build assets")
		}

		log.Infof("Assets has been built")

		return nil
	},
}

func init() {
	extensionRootCmd.AddCommand(extensionAssetBundleCmd)
}
