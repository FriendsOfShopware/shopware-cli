package git

import (
	"context"
	"os"
	"os/exec"
	"path"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInvalidGitRepository(t *testing.T) {
	repo := "invalid"
	ctx := context.Background()

	_, err := getPreviousTag(ctx, repo)
	if err == nil {
		t.Errorf("expected error, got nil")
	}
}

func TestNoTags(t *testing.T) {
	tmpDir := t.TempDir()
	prepareRepository(t, tmpDir)
	_ = os.WriteFile(path.Join(tmpDir, "a"), []byte(""), os.ModePerm)
	runCommand(t, tmpDir, "git", "add", "a")
	runCommand(t, tmpDir, "git", "commit", "-m", "initial commit", "--no-verify", "--no-gpg-sign")

	tag, err := getPreviousTag(context.Background(), tmpDir)
	assert.NoError(t, err)
	assert.NotEmpty(t, tag)

	commits, err := GetCommits(context.Background(), tmpDir)
	assert.NoError(t, err)
	assert.Len(t, commits, 0)
}

func TestWithOneTagAndCommit(t *testing.T) {
	tmpDir := t.TempDir()
	prepareRepository(t, tmpDir)
	_ = os.WriteFile(path.Join(tmpDir, "a"), []byte(""), os.ModePerm)
	runCommand(t, tmpDir, "git", "add", "a")
	runCommand(t, tmpDir, "git", "commit", "-m", "initial commit", "--no-verify", "--no-gpg-sign")
	runCommand(t, tmpDir, "git", "tag", "v1.0.0", "-m", "initial release")
	_ = os.WriteFile(path.Join(tmpDir, "b"), []byte(""), os.ModePerm)
	runCommand(t, tmpDir, "git", "add", "b")
	runCommand(t, tmpDir, "git", "commit", "-m", "second commit", "--no-verify", "--no-gpg-sign")

	tag, err := getPreviousTag(context.Background(), tmpDir)
	assert.NoError(t, err)
	assert.Equal(t, tag, "v1.0.0")

	commits, err := GetCommits(context.Background(), tmpDir)
	assert.NoError(t, err)
	assert.Len(t, commits, 1)
	assert.Equal(t, commits[0].Message, "second commit")
}

func TestGetPublicVCSURL(t *testing.T) {
	tmpDir := t.TempDir()
	prepareRepository(t, tmpDir)

	url, err := GetPublicVCSURL(context.Background(), tmpDir)
	assert.Equal(t, "", url)
	assert.Error(t, err)

	runCommand(t, tmpDir, "git", "remote", "add", "origin", "https://github.com/FriendsOfShopware/FroshTools.git")

	url, err = GetPublicVCSURL(context.Background(), tmpDir)
	assert.Equal(t, "https://github.com/FriendsOfShopware/FroshTools/commit", url)
	assert.NoError(t, err)

	runCommand(t, tmpDir, "git", "remote", "set-url", "origin", "git@github.com:FriendsOfShopware/FroshTools.git")

	url, err = GetPublicVCSURL(context.Background(), tmpDir)
	assert.Equal(t, "https://github.com/FriendsOfShopware/FroshTools/commit", url)
	assert.NoError(t, err)

	runCommand(t, tmpDir, "git", "remote", "set-url", "origin", "https://gitlab.com/xxx")
	t.Setenv("CI_PROJECT_URL", "https://example.com/gitlab-org/gitlab-foss")

	url, err = GetPublicVCSURL(context.Background(), tmpDir)
	assert.Equal(t, "https://example.com/gitlab-org/gitlab-foss/-/commit", url)
	assert.NoError(t, err)
}

func runCommand(t *testing.T, tmpDir, cmd string, args ...string) {
	t.Helper()

	c := exec.Command(cmd, args...)
	c.Dir = tmpDir

	out, err := c.CombinedOutput()
	if err != nil {
		t.Fatalf("command failed: %s", string(out))
	}
}

func prepareRepository(t *testing.T, tmpDir string) {
	t.Helper()

	runCommand(t, tmpDir, "git", "init")
	runCommand(t, tmpDir, "git", "config", "commit.gpgsign", "false")
	runCommand(t, tmpDir, "git", "config", "tag.gpgsign", "false")
	runCommand(t, tmpDir, "git", "config", "user.name", "test")
	runCommand(t, tmpDir, "git", "config", "user.email", "test@test.de")
}
