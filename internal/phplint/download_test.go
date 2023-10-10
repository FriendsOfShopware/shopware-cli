package phplint

import (
	"context"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDownloadPHPFile(t *testing.T) {
	if os.Getenv("NIX_CC") != "" {
		t.Skip("Downloading does not work in Nix build")
	}

	_, err := findPHPWasmFile(context.Background(), "7.4")
	assert.NoError(t, err)
}
