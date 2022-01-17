package cmd

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"shopware-cli/extension"

	"github.com/pkg/errors"

	termColor "github.com/fatih/color"
	cp "github.com/otiai10/copy"
	"github.com/spf13/cobra"
)

var disableGit = false

var extensionZipCmd = &cobra.Command{
	Use:   "zip [path] [branch]",
	Short: "Zip a Extension",
	Args:  cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		path, err := filepath.Abs(args[0])
		if err != nil {
			return err
		}

		var branch string
		if len(args) == 2 {
			branch = args[1]
		}

		ext, err := extension.GetExtensionByFolder(path)
		if err != nil {
			return errors.Wrap(err, "detect extension type")
		}

		name, err := ext.GetName()
		if err != nil {
			return errors.Wrap(err, "get name")
		}

		// Clear previous zips
		existingFiles, err := filepath.Glob(fmt.Sprintf("%s-*.zip", name))
		if err != nil {
			return err
		}

		for _, file := range existingFiles {
			err = os.Remove(file)
			if err != nil {
				return errors.Wrap(err, "remove existing file")
			}
		}

		// Create temp dir
		tempDir, err := ioutil.TempDir("", "extension")
		if err != nil {
			return errors.Wrap(err, "create temp directory")
		}

		extName, err := ext.GetName()
		if err != nil {
			return errors.Wrap(err, "get extension name")
		}

		extDir := fmt.Sprintf("%s/%s/", tempDir, extName)

		err = os.Mkdir(extDir, os.ModePerm)
		if err != nil {
			return errors.Wrap(err, "create temp directory")
		}

		tempDir += "/"

		defer func(path string) {
			_ = os.RemoveAll(path)
		}(tempDir)

		var tag string

		// Extract files using strategy
		if disableGit {
			err = cp.Copy(path, extDir)
			if err != nil {
				return errors.Wrap(err, "copy files")
			}
		} else {
			tag, err = extension.GitCopyFolder(path, extDir)
			if err != nil {
				return errors.Wrap(err, "copy via git")
			}
		}

		// User input wins
		if len(branch) > 0 {
			tag = branch
		}

		if err := extension.PrepareFolderForZipping(cmd.Context(), extDir, ext); err != nil {
			return errors.Wrap(err, "prepare package")
		}

		if err := extension.BuildAssetsForExtensions("", []extension.Extension{ext}); err != nil {
			return errors.Wrap(err, "building assets")
		}

		// Cleanup not wanted files
		if err := extension.CleanupExtensionFolder(extDir); err != nil {
			return errors.Wrap(err, "cleanup package")
		}

		fileName := fmt.Sprintf("%s-%s.zip", name, tag)
		if len(tag) == 0 {
			fileName = fmt.Sprintf("%s.zip", name)
		}

		if err := extension.CreateZip(tempDir, fileName); err != nil {
			return errors.Wrap(err, "create zip file")
		}

		termColor.Green("Created file %s", fileName)

		return nil
	},
}

func init() {
	extensionRootCmd.AddCommand(extensionZipCmd)
	extensionZipCmd.Flags().BoolVar(&disableGit, "disable-git", false, "Use the source folder as it is")
}
