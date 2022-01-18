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

	if len(versions) > 0 && len(versions[0]) > 0 {
		return versions[0], nil
	}

	branchCmd := exec.Command("git", "-C", source, "rev-parse", "--abbrev-ref", "HEAD")

	stdout, err = branchCmd.Output()

	if err != nil {
		return "", fmt.Errorf("gitTagOrBranchOfFolder: %v", err)
	}

	return strings.Trim(strings.TrimLeft(string(stdout), "* "), "\n"), nil
}

func GitCopyFolder(source, target string) (string, error) {
	tag, err := gitTagOrBranchOfFolder(source)

	if err != nil {
		return "", fmt.Errorf("GitCopyFolder: cannot find checkout tag or branch: %v", err)
	}

	archiveCmd := exec.Command("git", "-C", source, "archive", tag, "--format=zip")

	stdout, err := archiveCmd.Output()
	if err != nil {
		return "", fmt.Errorf("GitCopyFolder: cannot archive %s:  %v", tag, err)
	}

	zipReader, err := zip.NewReader(bytes.NewReader(stdout), int64(len(stdout)))
	if err != nil {
		return "", fmt.Errorf("GitCopyFolder: cannot open the zip file produced by git archive: %v", err)
	}

	err = Unzip(zipReader, target)
	if err != nil {
		return "", fmt.Errorf("GitCopyFolder: cannot unzip the zip archive: %v", err)
	}

	return tag, err
}
