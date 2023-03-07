package project

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"github.com/FriendsOfShopware/shopware-cli/shop"
)

var projectProxyCmd = &cobra.Command{
	Use:   "proxy",
	Short: "Proxy the local Shop to Cloudflare",
	RunE: func(cobraCmd *cobra.Command, _ []string) error {
		var cfg *shop.Config
		var err error

		if cfg, err = shop.ReadConfig(projectConfigPath); err != nil {
			return err
		}

		cloudflareInstalled, err := exec.LookPath("cloudflared")

		if err != nil {
			message := "Cloudflare Tunnel is not installed. Please use your system package manager to install it\n"

			if runtime.GOOS == "darwin" {
				message += "You can install it with Homebrew using brew install cloudflared"
			} else if runtime.GOOS == "windows" {
				message += "You can install it with Winget using winget install -e --id Cloudflare.cloudflared"
			} else {
				message += "See cloudflare for more information: https://developers.cloudflare.com/cloudflare-one/connections/connect-apps/install-and-setup/installation"
			}

			return fmt.Errorf(message)
		}

		log.Infof("Make sure you have set TRUSTED_PROXIES=127.0.0.1,::1 inside your .env file")

		command := exec.Command(cloudflareInstalled, "tunnel", "--url", cfg.URL)
		command.Stdin = os.Stdin
		command.Stdout = os.Stdout
		command.Stderr = os.Stderr

		return command.Run()
	},
}

func init() {
	projectRootCmd.AddCommand(projectProxyCmd)
}
