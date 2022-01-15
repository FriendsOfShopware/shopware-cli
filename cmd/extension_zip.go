package cmd

import (
	"fmt"
	termColor "github.com/fatih/color"
	cp "github.com/otiai10/copy"
	"github.com/spf13/cobra"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"shopware-cli/extension"
)

var disableGit = false

var extensionZipCmd = &cobra.Command{
	Use:   "zip [path] [branch]",
	Short: "Zip a Extension",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		path, err := filepath.Abs(args[0])
		branch := ""

		if len(args) == 2 {
			branch = args[1]
		}

		if err != nil {
			log.Fatalln(err)
		}

		ext, err := extension.GetExtensionByFolder(path)

		if err != nil {
			log.Fatalln(err)
		}

		name, err := ext.GetName()

		if err != nil {
			log.Fatalln(fmt.Errorf("zip: %v", err))
		}

		// Clear previous zips
		existingFiles, err := filepath.Glob(fmt.Sprintf("%s-*.zip", name))
		if err != nil {
			log.Fatalln(err)
		}

		for _, file := range existingFiles {
			_ = os.Remove(file)
		}

		// Create temp dir
		tempDir, err := ioutil.TempDir("", "extension")
		extName, err := ext.GetName()
		if err != nil {
			log.Fatalln(err)
		}

		extDir := fmt.Sprintf("%s/%s/", tempDir, extName)

		err = os.Mkdir(extDir, os.ModePerm)
		tempDir = tempDir + "/"

		if err != nil {
			log.Fatalln(err)
		}

		defer func(path string) {
			_ = os.RemoveAll(path)
		}(tempDir)

		tag := ""

		// Extract files using strategy
		if disableGit {
			err = cp.Copy(path, extDir)
		} else {
			tag, err = extension.GitCopyFolder(path, extDir)
		}

		// User input wins
		if len(branch) > 0 {
			tag = branch
		}

		if err != nil {
			log.Fatalln(err)
		}

		err = extension.PrepareFolderForZipping(extDir, ext)

		if err != nil {
			log.Fatalln(err)
		}

		// Cleanup not wanted files
		err = extension.CleanupExtensionFolder(extDir)
		if err != nil {
			log.Fatalln(err)
		}

		fileName := fmt.Sprintf("%s-%s.zip", name, tag)

		if len(tag) == 0 {
			fileName = fmt.Sprintf("%s.zip", name)
		}

		extension.CreateZip(tempDir, fileName)

		termColor.Green("Created file %s", fileName)
	},
}

func init() {
	extensionRootCmd.AddCommand(extensionZipCmd)
	extensionZipCmd.Flags().BoolVar(&disableGit, "disable-git", false, "Use the source folder as it is")
}
