package curl

import (
	"fmt"
	"net/url"
	"os"
	"os/exec"
	"strings"
)

type CurlCommand struct {
	options []curlOption
	args    []string
}

type curlOption struct {
	flag, value string
}

type CurlConfig func(*CurlCommand)

func Method(method string) CurlConfig {
	return func(command *CurlCommand) {
		command.addOption(curlOption{
			flag:  "-X",
			value: strings.ToUpper(method),
		})
	}
}

func BearerToken(token string) CurlConfig {
	return func(command *CurlCommand) {
		command.addOption(curlOption{
			flag:  "--header",
			value: fmt.Sprintf("Authorization: %s", token),
		})
	}
}

func Args(args []string) CurlConfig {
	return func(command *CurlCommand) {
		if len(args) == 0 {
			return
		}
		if len(command.args) > 0 {
			command.args = append(command.args, args...)
		} else {
			command.args = args
		}
	}
}

func Url(url *url.URL) CurlConfig {
	return func(command *CurlCommand) {
		if len(command.args) > 0 {
			command.args = append([]string{url.String()}, command.args...)
		} else {
			command.args = []string{url.String()}
		}
	}
}

func InitCurlCommand(options ...CurlConfig) *CurlCommand {
	cmd := &CurlCommand{}
	for _, opt := range options {
		opt(cmd)
	}

	return cmd
}

func (c *CurlCommand) addOption(o curlOption) {
	c.options = append(c.options, o)
}

func (c *CurlCommand) getCmdOptions() []string {
	var cmdOptions []string
	for _, opt := range c.options {
		cmdOptions = append(cmdOptions, opt.flag, opt.value)
	}
	return append(cmdOptions, c.args...)
}

func (c *CurlCommand) Run() error {
	// The user wants to execute code with custom parameters
	/* #nosec G204 */
	cmd := exec.Command("curl", c.getCmdOptions()...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
