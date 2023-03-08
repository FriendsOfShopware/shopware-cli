package project

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/manifoldco/promptui"
	"github.com/mholt/archiver/v3"
	"github.com/schollz/progressbar/v3"
	"github.com/spf13/cobra"

	"github.com/FriendsOfShopware/shopware-cli/logging"
	update_api "github.com/FriendsOfShopware/shopware-cli/update-api"
)

var projectCreateCmd = &cobra.Command{
	Use:   "create [name] [version]",
	Short: "Create a new Shopware 6 project",
	Args:  cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		projectFolder := args[0]

		if _, err := os.Stat(projectFolder); err == nil {
			return fmt.Errorf("the folder %s exists already", projectFolder)
		}

		if err := os.Mkdir(projectFolder, os.ModePerm); err != nil {
			return err
		}

		releases, err := update_api.GetLatestReleases(cmd.Context())

		if err != nil {
			return err
		}

		var result string

		if len(args) == 2 {
			result = args[1]
		} else {
			var chooseVersions []string

			for _, release := range releases {
				chooseVersions = append(chooseVersions, release.Version)
			}

			prompt := promptui.Select{
				Label: "Select Version",
				Items: chooseVersions,
			}

			if _, result, err = prompt.Run(); err != nil {
				return err
			}
		}

		var chooseVersion *update_api.ShopwareInstallRelease

		for _, release := range releases {
			if release.Version == result {
				chooseVersion = release
				break
			}
		}

		if chooseVersion == nil {
			_ = os.RemoveAll(projectFolder)
			return fmt.Errorf("cannot find version %s", result)
		}

		fileName := filepath.Base(chooseVersion.Uri)

		req, _ := http.NewRequest("GET", chooseVersion.Uri, nil)
		resp, _ := http.DefaultClient.Do(req)
		defer func(Body io.ReadCloser) {
			_ = Body.Close()
		}(resp.Body)

		f, _ := os.OpenFile(fileName, os.O_CREATE|os.O_WRONLY, 0644)
		defer func(f *os.File) {
			_ = f.Close()
		}(f)

		bar := progressbar.DefaultBytes(
			resp.ContentLength,
			"downloading",
		)

		if _, err := io.Copy(io.MultiWriter(f, bar), resp.Body); err != nil {
			return err
		}

		defer func(name string) {
			_ = os.Remove(name)
		}(fileName)

		logging.FromContext(cmd.Context()).Infof("Unpacking now the zip")

		if err := archiver.Unarchive(fileName, projectFolder); err != nil {
			return err
		}

		logging.FromContext(cmd.Context()).Infof("Shopware %s is created in folder %s", chooseVersion.Version, projectFolder)

		return nil
	},
}

func init() {
	projectRootCmd.AddCommand(projectCreateCmd)
}
