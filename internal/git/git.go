package git

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path"
	"strings"
)

type GitCommit struct {
	Hash    string
	Message string
}

func runGit(ctx context.Context, repo string, args ...string) (string, error) {
	cmd := exec.CommandContext(ctx, "git", args...)
	cmd.Dir = repo

	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("cannot run git: %w, %s", err, output)
	}

	if cmd.ProcessState.ExitCode() != 0 {
		return "", fmt.Errorf("cannot run git: %s", string(output))
	}

	gitOuput := string(output)
	return strings.Trim(gitOuput, " "), nil
}

func getPreviousTag(ctx context.Context, repo string) (string, error) {
	if err := unshallowRepository(ctx, repo); err != nil {
		return "", err
	}

	previousVersion := os.Getenv("SHOPWARE_CLI_PREVIOUS_TAG")
	if previousVersion != "" {
		return previousVersion, nil
	}

	commits, err := runGit(ctx, repo, "log", "--pretty=format:%h", "--no-merges")
	if err != nil {
		return "", fmt.Errorf("cannot get previous tag: %w", err)
	}

	commitsArray := strings.Split(commits, "\n")
	for commit := range commitsArray {
		contains, err := runGit(ctx, repo, "tag", "--contains", commitsArray[commit])
		if err != nil {
			return "", fmt.Errorf("cannot get previous tag: %w", err)
		}

		if contains == "" {
			continue
		}

		matchingTags := strings.Split(contains, "\n")

		if len(matchingTags) == 0 {
			continue
		}

		return matchingTags[0], nil
	}

	// if no tag was found, return the first commit
	return commitsArray[len(commitsArray)-1], nil
}

func GetCommits(ctx context.Context, repo string) ([]GitCommit, error) {
	if err := unshallowRepository(ctx, repo); err != nil {
		return nil, err
	}

	previousTag, err := getPreviousTag(ctx, repo)
	if err != nil {
		return nil, err
	}

	commits, err := runGit(ctx, repo, "log", "--pretty=format:%h|%s", previousTag+"..HEAD", "--no-merges")
	if err != nil {
		return nil, fmt.Errorf("cannot get commits: %w", err)
	}

	if commits == "" {
		return []GitCommit{}, nil
	}

	commitsArray := strings.Split(commits, "\n")
	gitCommits := make([]GitCommit, len(commitsArray))

	for commit := range commitsArray {
		splitCommit := strings.Split(commitsArray[commit], "|")
		gitCommits[commit] = GitCommit{
			Hash:    splitCommit[0],
			Message: strings.Join(splitCommit[1:], "|"),
		}
	}

	return gitCommits, nil
}

func GetPublicVCSURL(ctx context.Context, repo string) (string, error) {
	origin, err := runGit(ctx, repo, "config", "--get", "remote.origin.url")
	if err != nil {
		return "", fmt.Errorf("failed to run git command: %w", err)
	}

	origin = strings.Trim(origin, "\n")

	switch {
	case strings.HasPrefix(origin, "https://github.com/"):
		origin = strings.TrimSuffix(origin, ".git")

		return fmt.Sprintf("%s/commit", origin), nil
	case strings.HasPrefix(origin, "git@github.com:"):
		origin = origin[15:]
		origin = strings.TrimSuffix(origin, ".git")

		return fmt.Sprintf("https://github.com/%s/commit", origin), nil
	case os.Getenv("CI_PROJECT_URL") != "":
		return fmt.Sprintf("%s/-/commit", os.Getenv("CI_PROJECT_URL")), nil
	}

	return "", fmt.Errorf("unsupported vcs provider")
}

func unshallowRepository(ctx context.Context, repo string) error {
	if _, err := os.Stat(path.Join(repo, ".git", "shallow")); os.IsNotExist(err) {
		return nil
	}

	_, err := runGit(ctx, repo, "fetch", "--unshallow")

	return err
}
