package cmd

import (
	"github.com/spf13/cobra"
	"log"
	"path/filepath"
	"shopware-cli/extension"
)

var extensionPrepareCmd = &cobra.Command{
	Use:   "prepare [path]",
	Short: "Prepare a extension for zipping",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		path, err := filepath.Abs(args[0])

		if err != nil {
			log.Fatalln(err)
		}

		ext, err := extension.GetExtensionByFolder(path)

		if err != nil {
			log.Fatalln(err)
		}

		err = extension.PrepareFolderForZipping(path+"/", ext)

		if err != nil {
			log.Fatalln(err)
		}
	},
}

func init() {
	extensionRootCmd.AddCommand(extensionPrepareCmd)
}
