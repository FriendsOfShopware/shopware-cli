package extension

import (
	"fmt"
	"github.com/FriendsOfShopware/shopware-cli/extension"
	"os"
	"path/filepath"

	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"

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

		if !context.HasErrors() {
			log.Infof("Validation passed without errors")
			return nil
		}

		table := tablewriter.NewWriter(os.Stdout)
		table.SetColWidth(100)
		table.SetHeader([]string{"Message"})

		for _, msg := range context.Errors() {
			table.Append([]string{msg})
		}

		table.Render()

		return fmt.Errorf("validation failed")
	},
}

func init() {
	extensionRootCmd.AddCommand(extensionValidateCmd)
}
