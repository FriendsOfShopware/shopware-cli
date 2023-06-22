package extension

import (
	"context"
	"encoding/json"
	"os"
	"path"

	"github.com/FriendsOfShopware/shopware-cli/logging"
)

func FindExtensionsFromProject(ctx context.Context, project string) []Extension {
	extensions := make(map[string]Extension)

	for _, ext := range addExtensionsByComposer(project) {
		name, err := ext.GetName()
		if err != nil {
			continue
		}

		version, _ := ext.GetVersion()

		logging.FromContext(ctx).Infof("Found extension using Composer: %s (%s)", name, version)

		extensions[name] = ext
	}

	for _, ext := range addExtensionsByWildcard(path.Join(project, "custom", "plugins")) {
		name, err := ext.GetName()
		if err != nil {
			continue
		}

		// Skip if extension is already added by composer
		if _, ok := extensions[name]; ok {
			continue
		}

		version, _ := ext.GetVersion()

		logging.FromContext(ctx).Infof("Found extension in custom/plugins: %s (%s)", name, version)
		logging.FromContext(ctx).Errorf("Extension %s should be installed using Composer. Please remove the extension from custom/plugins.", name)

		extensions[name] = ext
	}

	for _, ext := range addExtensionsByWildcard(path.Join(project, "custom", "apps")) {
		name, err := ext.GetName()
		if err != nil {
			continue
		}
		version, _ := ext.GetVersion()

		logging.FromContext(ctx).Infof("Found extension in custom/apps: %s (%s)", name, version)

		extensions[name] = ext
	}

	extensionsSlice := make([]Extension, 0, len(extensions))

	for _, ext := range extensions {
		extensionsSlice = append(extensionsSlice, ext)
	}

	return extensionsSlice
}

func addExtensionsByComposer(project string) []Extension {
	var list []Extension

	lock, err := os.ReadFile(path.Join(project, "composer.lock"))
	if err != nil {
		return list
	}

	var composer composerLock
	err = json.Unmarshal(lock, &composer)

	if err != nil {
		return list
	}

	for _, pkg := range composer.Packages {
		if pkg.PackageType == "shopware-platform-plugin" {
			ext, err := GetExtensionByFolder(path.Join(project, "vendor", pkg.Name))
			if err != nil {
				continue
			}

			list = append(list, ext)
		}
	}

	return list
}

func addExtensionsByWildcard(extensionDir string) []Extension {
	var list []Extension

	extensions, err := os.ReadDir(extensionDir)
	if err != nil {
		return list
	}

	for _, file := range extensions {
		if file.IsDir() {
			ext, err := GetExtensionByFolder(path.Join(extensionDir, file.Name()))
			if err != nil {
				continue
			}

			list = append(list, ext)
		}
	}

	return list
}

type composerLock struct {
	Packages []struct {
		Name        string `json:"name"`
		Version     string `json:"version"`
		PackageType string `json:"type"`
	} `json:"packages"`
}
