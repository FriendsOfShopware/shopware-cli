package app

import (
	"archive/zip"
	"bytes"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"shopware-cli/extension"
	"shopware-cli/shop"
	"shopware-cli/tui"
)

var credentialFile string
var url string

var pushCommand = &cobra.Command{
	Use:   "push --url url [--credentials file.json] [--dir path]",
	Short: "Install the app in --dir to a shopware instance",
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()
		tui := &tui.TUI{}
		var credentials shop.Credentials
		if len(credentialFile) > 0 {
			c, err := shop.ReadCredentialsFromFile(credentialFile)
			if err != nil {
				return err
			}
			credentials = c
		}
		if credentials == nil {
			user, err := tui.AskForUsername()
			if err != nil {
				return err
			}
			pw, err := tui.AskForPassword(user)
			if err != nil {
				return err
			}
			credentials = shop.PasswordCredentials{
				Username: user,
				Password: pw,
			}

		}

		shopClient, err := shop.NewShopClient(ctx, url, credentials, nil)
		if err != nil {
			return err
		}
		shopClient.TUI = tui

		app, err := extension.GetExtensionByFolder(appDir)
		if err != nil {
			return err
		}
		app, ok := app.(*extension.App)
		if !ok {
			return errors.New("given directory contains a plugin")
		}
		appName, err := app.GetName()
		if err != nil {
			return err
		}
		var buf bytes.Buffer
		w := zip.NewWriter(&buf)
		extension.AddZipFiles(w, appDir, appName+"/")
		w.Close()

		if err := shopClient.UploadExtension(ctx, &buf); err != nil {
			return err
		}

		tl := tui.ShowTaskList("uploading", "installing", "activating")
		tl.Done()

		if err := shopClient.InstallApp(ctx, appName); err != nil {
			return err
		}
		tl.Done()

		if err := shopClient.ActivateApp(ctx, appName); err != nil {
			return err
		}
		tl.Done()
		return nil
	},
}

func init() {
	pushCommand.Flags().StringVar(&credentialFile, "credentials", "", "a json file containing a oauth grant")
	pushCommand.Flags().StringVar(&url, "url", "", "url of shopware instance")
	_ = pushCommand.MarkFlagRequired("url")
	appRootCommand.AddCommand(pushCommand)
}
