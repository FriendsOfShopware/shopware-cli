package project

import (
	"bytes"
	"context"
	_ "embed"
	"fmt"
	"github.com/FriendsOfShopware/shopware-cli/extension"
	"github.com/FriendsOfShopware/shopware-cli/logging"
	"github.com/FriendsOfShopware/shopware-cli/shop"
	"github.com/spf13/cobra"
	"os"
	"os/exec"
	"strings"
	"text/template"
)

//go:embed templates/Dockerfile.tpl
var dockerFileTemplate string

var dockerBuildCmd = &cobra.Command{
	Use:   "build [name]",
	Short: "Build Docker Image",
	Args:  cobra.ExactArgs(1),
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

		return runTransparentCommand(exec.CommandContext(cmd.Context(), "docker", "build", "-t", args[0], "."))
	},
}

func dumpDockerIgnore(shopCfg *shop.Config) error {
	shopCfg.Docker.ExcludePaths = append(shopCfg.Docker.ExcludePaths, "/var", ".git", "node_modules", ".idea")

	if err := os.WriteFile(".dockerignore", []byte(strings.Join(shopCfg.Docker.ExcludePaths, "\n")), os.ModePerm); err != nil {
		return err
	}

	return nil
}

func dumpDockerfile(ctx context.Context, shopCfg *shop.Config) error {
	templateVars, err := configureDockerfileTemplate(ctx, shopCfg)
	if err != nil {
		return err
	}

	dockerfile, err := renderDockerfile(templateVars)
	if err != nil {
		return err
	}

	if err := os.WriteFile("Dockerfile", dockerfile, os.ModePerm); err != nil {
		return err
	}

	return nil
}

func configureDockerfileTemplate(ctx context.Context, shopCfg *shop.Config) (map[string]interface{}, error) {
	if shopCfg.Docker.PHP.PhpVersion == "" {
		projectRoot, err := os.Getwd()
		if err != nil {
			return nil, err
		}

		constraint, err := extension.GetShopwareProjectConstraint(projectRoot)
		if err != nil {
			return nil, err
		}

		phpVersion, err := extension.GetPhpVersion(ctx, constraint)
		if err != nil {
			return nil, err
		}

		shopCfg.Docker.PHP.PhpVersion = phpVersion
		logging.FromContext(ctx).Infof("No PHP version set, using PHP version %s", phpVersion)
	}

	buildEnvironments := make([]string, 0)
	runEnvironments := make([]string, 0)

	for _, value := range shopCfg.Docker.Environment {
		envLine := fmt.Sprintf("%s=%s", value.Name, value.Value)

		switch value.Only {
		case "build":
			buildEnvironments = append(buildEnvironments, envLine)
		case "runtime":
			runEnvironments = append(runEnvironments, envLine)
		default:
			buildEnvironments = append(buildEnvironments, envLine)
			runEnvironments = append(runEnvironments, envLine)
		}
	}

	hooks := make(map[string]string)
	if shopCfg.Docker.Hooks.PreUpdate != "" {
		hooks["pre_update"] = shopCfg.Docker.Hooks.PreUpdate
	}

	if shopCfg.Docker.Hooks.PostUpdate != "" {
		hooks["post_update"] = shopCfg.Docker.Hooks.PostUpdate
	}

	if shopCfg.Docker.Hooks.PreInstall != "" {
		hooks["pre_install"] = shopCfg.Docker.Hooks.PreInstall
	}

	if shopCfg.Docker.Hooks.PostInstall != "" {
		hooks["post_install"] = shopCfg.Docker.Hooks.PostInstall
	}

	templateVars := map[string]interface{}{
		"PHP":          shopCfg.Docker.PHP,
		"Variant":      shopCfg.Docker.Variant,
		"ExcludePaths": shopCfg.Docker.ExcludePaths,
		"BuildEnv":     strings.Join(buildEnvironments, " "),
		"RunEnv":       strings.Join(runEnvironments, " "),
		"Hooks":        hooks,
	}

	return templateVars, nil
}

func renderDockerfile(cfg map[string]interface{}) ([]byte, error) {
	var buf bytes.Buffer

	err := template.
		Must(template.New("Dockerfile").Parse(dockerFileTemplate)).
		Execute(&buf, cfg)

	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func init() {
	dockerRootCmd.AddCommand(dockerBuildCmd)
}
