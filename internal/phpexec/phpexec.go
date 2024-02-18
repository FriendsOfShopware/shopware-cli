package phpexec

import (
	"context"
	"os"
	"os/exec"
	"sync"
)

var pathToSymfonyCLI = sync.OnceValue[string](func() string {
	path, err := exec.LookPath("symfony")
	if err != nil {
		return ""
	}
	return path
})

func symfonyCliAllowed() bool {
	return os.Getenv("SHOPWARE_CLI_NO_SYMFONY_CLI") != "1"
}

func ConsoleCommand(ctx context.Context, args ...string) *exec.Cmd {
	if path := pathToSymfonyCLI(); path != "" && symfonyCliAllowed() {
		return exec.CommandContext(ctx, path, append([]string{"console"}, args...)...)
	}
	return exec.CommandContext(ctx, "php", append([]string{"bin/console"}, args...)...)
}

func ComposerCommand(ctx context.Context, args ...string) *exec.Cmd {
	if path := pathToSymfonyCLI(); path != "" && symfonyCliAllowed() {
		return exec.CommandContext(ctx, path, append([]string{"composer"}, args...)...)
	}
	return exec.CommandContext(ctx, "composer", args...)
}

func PHPCommand(ctx context.Context, args ...string) *exec.Cmd {
	if path := pathToSymfonyCLI(); path != "" && symfonyCliAllowed() {
		return exec.CommandContext(ctx, path, append([]string{"php"}, args...)...)
	}
	return exec.CommandContext(ctx, "php", args...)
}
