package main

import (
	"os"
	"os/exec"
	"regexp"
	"strings"
)

func main() {
	cmd := exec.Command("nix", "build")
	nixOutput, _ := cmd.CombinedOutput()

	if cmd.ProcessState.ExitCode() == 0 {
		return
	}

	re := regexp.MustCompile(`(?m)got:\s(.*)\serror`)

	newSha := strings.Trim(re.FindStringSubmatch(string(nixOutput))[1], " ")

	flake, err := os.ReadFile("flake.nix")
	if err != nil {
		panic(err)
	}

	replace := regexp.MustCompile(`(?m)vendorSha256.*`)

	content := replace.ReplaceAllString(string(flake), "vendorSha256 = \""+newSha+"\";")

	err = os.WriteFile("flake.nix", []byte(content), 0o644)

	if err != nil {
		panic(err)
	}
}
