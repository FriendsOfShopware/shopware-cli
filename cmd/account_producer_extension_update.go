package cmd

import (
	"fmt"
	termColor "github.com/fatih/color"
	"github.com/spf13/cobra"
	"log"
	"os"
	"path/filepath"
	"shopware-cli/extension"
)

var accountCompanyProducerExtensionUpdateCmd = &cobra.Command{
	Use:   "update [zip or path]",
	Short: "Update store information of extension",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		client := getAccountApiByConfig()

		path, err := filepath.Abs(args[0])

		if err != nil {
			log.Fatalln(fmt.Errorf("update: %v", err))
		}

		stat, err := os.Stat(path)

		if err != nil {
			log.Fatalln(fmt.Errorf("update: %v", err))
		}

		var zipExt extension.Extension

		if stat.IsDir() {
			zipExt, err = extension.GetExtensionByFolder(path)
		} else {
			zipExt, err = extension.GetExtensionByZip(path)
		}

		if err != nil {
			log.Fatalln(fmt.Errorf("update: %v", err))
		}

		zipName, err := zipExt.GetName()

		if err != nil {
			log.Fatalln(fmt.Errorf("update: %v", err))
		}

		p, err := client.Producer()

		if err != nil {
			log.Fatalln(fmt.Errorf("update: %v", err))
		}

		storeExt, err := p.GetExtensionByName(zipName)

		if err != nil {
			log.Fatalln(fmt.Errorf("update: %v", err))
		}

		metadata := zipExt.GetMetaData()

		for _, info := range storeExt.Infos {
			language := info.Locale.Name[0:2]

			if language == "de" {
				info.Name = metadata.Label.German
				info.ShortDescription = metadata.Description.German
			} else {
				info.Name = metadata.Label.English
				info.ShortDescription = metadata.Description.English
			}
		}

		err = p.UpdateExtension(storeExt)

		if err != nil {
			log.Fatalln(fmt.Errorf("update: %v", err))
		}

		termColor.Green("Store information has been updated")
	},
}

func init() {
	accountCompanyProducerExtensionCmd.AddCommand(accountCompanyProducerExtensionUpdateCmd)
}
