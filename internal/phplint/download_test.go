package phplint

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDownloadPHPFile(t *testing.T) {
	_, err := findPHPWasmFile(context.Background(), "7.4")
	assert.NoError(t, err)
}
