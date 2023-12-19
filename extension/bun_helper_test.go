package extension

import (
	"github.com/stretchr/testify/assert"
	"os"
	"path"
	"testing"
)

func TestValidPackageJsonBun(t *testing.T) {
	tmpDir := t.TempDir()

	packageJson := `{
		"dependencies": {
			"foo": "1.0.0"
		}
	}`

	if err := os.WriteFile(path.Join(tmpDir, "package.json"), []byte(packageJson), 0644); err != nil {
		t.Fatal(err)
	}

	assert.True(t, canRunBunOnPackage(tmpDir))
}

func TestValidPackageJsonWithDevBun(t *testing.T) {
	tmpDir := t.TempDir()

	packageJson := `{
		"dependencies": {
			"foo": "1.0.0"
		},
		"devDependencies": {
			"bar": "1.0.0"
		}
	}`

	if err := os.WriteFile(path.Join(tmpDir, "package.json"), []byte(packageJson), 0644); err != nil {
		t.Fatal(err)
	}

	assert.True(t, canRunBunOnPackage(tmpDir))
}

func TestInvalidPackageJsonBun(t *testing.T) {
	tmpDir := t.TempDir()

	packageJson := `{
		"dependencies": {
			"foo": "1.0.0"
		},
		"devDependencies": {
			"foo": "1.0.0"
		}
	}`

	if err := os.WriteFile(path.Join(tmpDir, "package.json"), []byte(packageJson), 0644); err != nil {
		t.Fatal(err)
	}

	assert.False(t, canRunBunOnPackage(tmpDir))
}
