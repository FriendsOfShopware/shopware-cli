package esbuild

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestKebabCase(t *testing.T) {
	assert.Equal(t, "foo-bar", ToKebabCase("FooBar"))
	assert.Equal(t, "f-o-o-bar-baz", ToKebabCase("FOOBarBaz"))
	assert.Equal(t, "frosh-tools", ToKebabCase("FroshTools"))
	assert.Equal(t, "my-module-name-s-w6", ToKebabCase("MyModuleNameSW6"))
	assert.Equal(t, "a-i-search", ToKebabCase("AISearch"))
}
