package extension

type npmPackage struct {
	Dependencies    map[string]string `json:"dependencies"`
	DevDependencies map[string]string `json:"devDependencies"`
}

// When a package is defined in both dependencies and devDependencies, bun will crash.
func canRunBunOnPackage(npmPackage npmPackage) bool {
	for name := range npmPackage.Dependencies {
		if _, ok := npmPackage.DevDependencies[name]; ok {
			return false
		}
	}

	return true
}
