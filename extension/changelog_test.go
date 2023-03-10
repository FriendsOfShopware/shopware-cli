package extension

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestChangelogParsing(t *testing.T) {
	content := parseMarkdownChangelog("# 1.0.0\n- Test\n- Test2\n# 2.0.0\n- Test3\n- Test4\n")

	assert.Equal(t, "Test<br> Test2<br>", content["1.0.0"])
}
