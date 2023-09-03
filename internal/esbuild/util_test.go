package esbuild

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSnakeCase(t *testing.T) {
	assert.Equal(t, "foo_bar", toSnakeCase("FooBar"))
	assert.Equal(t, "frosh_tools", toSnakeCase("FroshTools"))
	assert.Equal(t, "my_module_name_s_w6", toSnakeCase("MyModuleNameSW6"))
}
