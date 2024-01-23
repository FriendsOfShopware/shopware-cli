package project

import (
	"bytes"
	"crypto/rand"
	_ "embed"
	"encoding/base64"
	"github.com/FriendsOfShopware/shopware-cli/shop"
	"github.com/spf13/cobra"
	"os"
	"os/exec"
	"text/template"
)

//go:embed templates/compose.yaml
var composeFileTemplate string

var dockerUpCmd = &cobra.Command{
	Use:   "up",
	Short: "Start local setup",
	RunE: func(cmd *cobra.Command, args []string) error {
		shopCfg, err := shop.ReadConfig(projectConfigPath, true)
		if err != nil {
			return err
		}

		if err = dumpDockerfile(cmd.Context(), shopCfg); err != nil {
			return err
		}

		if err = dumpDockerIgnore(shopCfg); err != nil {
			return err
		}

		if err = dumpComposeFile(); err != nil {
			return err
		}

		return runTransparentCommand(exec.CommandContext(cmd.Context(), "docker", "compose", "watch"))
	},
}

func dumpComposeFile() error {
	templateVariables, err := configureComposeTemplate()
	if err != nil {
		return err
	}

	composeFile, err := renderComposeFile(templateVariables)
	if err != nil {
		return err
	}

	if err = os.WriteFile("compose.yaml", composeFile, os.ModePerm); err != nil {
		return err
	}

	return nil
}

func configureComposeTemplate() (map[string]interface{}, error) {
	publicKey, privateKey, err := generatePrivatePublicKey(2048)
	if err != nil {
		return nil, err
	}

	appSecret := make([]byte, 16)
	instanceID := make([]byte, 16)

	if _, err = rand.Read(appSecret); err != nil {
		return nil, err
	}

	if _, err = rand.Read(instanceID); err != nil {
		return nil, err
	}

	config := map[string]interface{}{
		"appSecret":     base64.URLEncoding.EncodeToString(appSecret),
		"instanceID":    base64.URLEncoding.EncodeToString(instanceID),
		"jwtPrivateKey": base64.StdEncoding.EncodeToString(privateKey),
		"jwtPublicKey":  base64.StdEncoding.EncodeToString(publicKey),
	}

	return config, nil
}

func renderComposeFile(templateVars map[string]interface{}) ([]byte, error) {
	var buf bytes.Buffer

	err := template.
		Must(template.New("compose.yaml").Parse(composeFileTemplate)).
		Execute(&buf, templateVars)

	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func init() {
	dockerRootCmd.AddCommand(dockerUpCmd)
}
