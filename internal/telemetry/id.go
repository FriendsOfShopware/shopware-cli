package telemetry

import (
	"crypto/sha256"
	"encoding/hex"
	"os"
	"path"

	"github.com/FriendsOfShopware/shopware-cli/internal/system"
	"github.com/google/uuid"
)

// gatherDistinctId returns a distinct ID for the current user, so we anoynmously track the usage of the CLI
func gatherDistinctId() string {
	if os.Getenv("CI") == "true" {
		return gatherByCI()
	}

	return gatherByMachine()
}

func gatherByCI() string {
	if os.Getenv("GITHUB_REPOSITORY") != "" {
		return hash(os.Getenv("GITHUB_REPOSITORY"))
	}

	// GitLab
	if os.Getenv("CI_PROJECT_NAME") != "" {
		return hash(os.Getenv("CI_PROJECT_NAME"))
	}

	// Bitbucket
	if os.Getenv("BITBUCKET_REPO_FULL_NAME") != "" {
		return hash(os.Getenv("BITBUCKET_REPO_FULL_NAME"))
	}

	// We cannot determine the CI system, so we generate a random UUID
	return hash(uuid.New().String())
}

func gatherByMachine() string {
	configDir := system.GetShopwareCliConfigDir()

	telemetryFile := path.Join(configDir, "telemetry-id")

	if contents, err := os.ReadFile(telemetryFile); err == nil {
		return string(contents)
	}

	if _, err := os.Stat(configDir); os.IsNotExist(err) {
		_ = os.MkdirAll(configDir, 0o700)
	}

	id := hash(uuid.New().String())
	_ = os.WriteFile(telemetryFile, []byte(id), 0o600)

	return id
}

func hash(s string) string {
	return hex.EncodeToString(sha256.New().Sum([]byte(s)))
}
