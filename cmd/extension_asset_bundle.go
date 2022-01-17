package cmd

import (
	termColor "github.com/fatih/color"
	"github.com/spf13/cobra"
	"log"
	"path/filepath"
	"shopware-cli/extension"
)

var extensionAssetBundleCmd = &cobra.Command{
	Use:   "build [path]",
	Short: "Builds assets for extensions",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		validatedExtensions := make([]extension.Extension, 0)

		for _, arg := range args {
			path, err := filepath.Abs(arg)

			if err != nil {
				log.Fatalln(err)
			}

			ext, err := extension.GetExtensionByFolder(path)

			if err != nil {
				log.Fatalln(err)
			}

			validatedExtensions = append(validatedExtensions, ext)
		}

		err := extension.BuildAssetsForExtensions("", validatedExtensions)

		if err != nil {
			log.Fatalln(err)
		}

		termColor.Green("Assets has been built")
	},
}

func init() {
	extensionRootCmd.AddCommand(extensionAssetBundleCmd)
}
