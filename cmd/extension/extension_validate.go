package extension

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/FriendsOfShopware/shopware-cli/extension"
	log "github.com/sirupsen/logrus"

	"github.com/pkg/errors"

	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
)

var extensionValidateCmd = &cobra.Command{
	Use:   "validate [path]",
	Short: "Validate a Extension",
	Args:  cobra.MinimumNArgs(1),
	RunE: func(_ *cobra.Command, args []string) error {
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
			return errors.Wrap(err, "cannot open extension")
		}

		context := extension.RunValidation(ext)

		if context.HasErrors() || context.HasWarnings() {
			table := tablewriter.NewWriter(os.Stdout)
			table.SetHeader([]string{"Type", "Message"})
			table.SetAutoWrapText(false)

			for _, msg := range context.Errors() {
				table.Append([]string{"Error", msg})
			}

			for _, msg := range context.Warnings() {
				table.Append([]string{"Warning", msg})
			}

			table.Render()
		}

		if context.HasErrors() {
			return fmt.Errorf("validation failed")
		}

		log.Infof("Validation has been successful")

		return nil
	},
}

func init() {
	extensionRootCmd.AddCommand(extensionValidateCmd)
}
