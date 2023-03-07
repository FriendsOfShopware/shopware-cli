package extension

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/FriendsOfShopware/shopware-cli/version"
)

func TestDetermineMinVersion(t *testing.T) {
	constraint, _ := version.NewConstraint("~6.5.0")

	matchingVersion, _ := getMinMatchingVersion(&constraint, []string{"6.4.0.0", "6.5.0.0-rc1", "6.5.0.0"})
	assert.Equal(t, "6.5.0.0", matchingVersion)
	matchingVersion, _ = getMinMatchingVersion(&constraint, []string{"6.4.0.0", "6.5.0.0-rc1"})
	assert.Equal(t, "6.5.0.0-rc1", matchingVersion)
	matchingVersion, _ = getMinMatchingVersion(&constraint, []string{"6.5.0.0-rc1", "6.4.0.0"})
	assert.Equal(t, "6.5.0.0-rc1", matchingVersion)

	_, err := getMinMatchingVersion(&constraint, []string{"1.0.0", "2.0.0"})
	assert.Error(t, err)

	matchingVersion, _ = getMinMatchingVersion(&constraint, []string{"6.5.0.0-rc1", "abc", "6.4.0.0"})
	assert.Equal(t, "6.5.0.0-rc1", matchingVersion)
}
