package extension

import (
	"encoding/json"
	"os"
	"path"
)

type npmPackage struct {
	Dependencies    map[string]string `json:"dependencies"`
	DevDependencies map[string]string `json:"devDependencies"`
}

// When a package is defined in both dependencies and devDependencies, bun will crash.
func canRunBunOnPackage(packagePath string) bool {
	packageJson, err := os.ReadFile(path.Join(packagePath, "package.json"))

	if err != nil {
		return false
	}

	var npmPackage npmPackage

	if json.Unmarshal(packageJson, &npmPackage) != nil {
		return false
	}

	for name := range npmPackage.Dependencies {
		if _, ok := npmPackage.DevDependencies[name]; ok {
			return false
		}
	}

	return true
}
