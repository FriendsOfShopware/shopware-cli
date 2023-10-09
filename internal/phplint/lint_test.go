package phplint

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLintTestData(t *testing.T) {
	errors, err := LintFolder(context.Background(), "7.4", "testdata")

	assert.NoError(t, err)

	assert.Len(t, errors, 1)

	assert.Equal(t, "invalid.php", errors[0].File)
	assert.Contains(t, errors[0].Message, "syntax error, unexpected end of file")
}
