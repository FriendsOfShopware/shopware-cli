package phplint

import (
	"context"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLintTestData(t *testing.T) {
	if os.Getenv("NIX_CC") != "" {
		t.Skip("Downloading does not work in Nix build")
	}

	errors, err := LintFolder(context.Background(), "7.4", "testdata")

	assert.NoError(t, err)

	assert.Len(t, errors, 1)

	assert.Equal(t, "invalid.php", errors[0].File)
	assert.Contains(t, errors[0].Message, "syntax error, unexpected end of file")
}
