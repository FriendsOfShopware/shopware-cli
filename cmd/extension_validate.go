package cmd

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"shopware-cli/extension"

	termColor "github.com/fatih/color"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
)

var extensionValidateCmd = &cobra.Command{
	Use:   "validate [path]",
	Short: "Validate a Extension",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		path, err := filepath.Abs(args[0])

		if err != nil {
			log.Fatalln(fmt.Errorf("validate: %v", err))
		}

		stat, err := os.Stat(path)

		if err != nil {
			log.Fatalln(fmt.Errorf("validate: %v", err))
		}

		var ext extension.Extension

		if stat.IsDir() {
			ext, err = extension.GetExtensionByFolder(path)
		} else {
			ext, err = extension.GetExtensionByZip(path)
		}

		if err != nil {
			log.Fatalln(fmt.Errorf("validate: %v", err))
		}

		context := extension.RunValidation(ext)

		if !context.HasErrors() {
			termColor.Green("Validation passed without errors")
			os.Exit(0)
		}

		termColor.Red("Validation failed")

		table := tablewriter.NewWriter(os.Stdout)
		table.SetColWidth(100)
		table.SetHeader([]string{"Message"})

		for _, msg := range context.Errors() {
			table.Append([]string{msg})
		}

		table.Render()

		os.Exit(1)
	},
}

func init() {
	extensionRootCmd.AddCommand(extensionValidateCmd)
}
