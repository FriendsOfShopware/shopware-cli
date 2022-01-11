package extension

import (
	"archive/zip"
	"bytes"
	"fmt"
	"os/exec"
	"strings"
)

func gitTagOrBranchOfFolder(source string) (string, error) {
	tagCmd := exec.Command("git", "-C", source, "tag", "--sort=-creatordate")

	stdout, err := tagCmd.Output()

	if err != nil {
		return "", err
	}

	versions := strings.Split(string(stdout), "\n")

	if len(versions) > 0 {
		return versions[0], fmt.Errorf("gitTagOrBranchOfFolder: %v", nil)
	}

	branchCmd := exec.Command("git", "-C", source, "branch")

	stdout, err = branchCmd.Output()

	if err != nil {
		return "", fmt.Errorf("gitTagOrBranchOfFolder: %v", err)
	}

	return strings.TrimLeft(string(stdout), "* "), nil
}

func GitCopyFolder(source, target string) (string, error) {
	tag, err := gitTagOrBranchOfFolder(source)

	if err != nil {
		return "", fmt.Errorf("GitCopyFolder: %v", err)
	}

	archiveCmd := exec.Command("git", "-C", source, "archive", tag, "--format=zip")

	stdout, err := archiveCmd.Output()
	if err != nil {
		return "", fmt.Errorf("GitCopyFolder: %v", err)
	}

	zipReader, err := zip.NewReader(bytes.NewReader(stdout), int64(len(stdout)))
	if err != nil {
		return "", fmt.Errorf("GitCopyFolder: %v", err)
	}

	err = Unzip(zipReader, target)
	if err != nil {
		return "", fmt.Errorf("GitCopyFolder: %v", err)
	}

	return tag, err
}
