package extension

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"reflect"
	"strings"

	"github.com/wI2L/jsondiff"
)

func validateStorefrontSnippets(context *ValidationContext) {
	rootDir := context.Extension.GetRootDir()
	if err := validateStorefrontSnippetsByPath(rootDir, rootDir, context); err != nil {
		return
	}

	for _, extraBundle := range context.Extension.GetExtensionConfig().Build.ExtraBundles {
		bundlePath := rootDir

		if extraBundle.Path != "" {
			bundlePath = path.Join(bundlePath, extraBundle.Path)
		} else {
			bundlePath = path.Join(bundlePath, extraBundle.Name)
		}

		if err := validateStorefrontSnippetsByPath(bundlePath, rootDir, context); err != nil {
			return
		}
	}
}

func validateStorefrontSnippetsByPath(extensionRoot, rootDir string, context *ValidationContext) error {
	snippetFolder := path.Join(extensionRoot, "Resources", "snippet")

	if _, err := os.Stat(snippetFolder); err != nil {
		return nil //nolint:nilerr
	}

	snippetFiles := make(map[string][]string)

	err := filepath.WalkDir(snippetFolder, func(path string, d os.DirEntry, err error) error {
		if d.IsDir() {
			return nil
		}

		if filepath.Ext(path) != ".json" {
			return nil
		}

		containingFolder := filepath.Dir(path)

		if _, ok := snippetFiles[containingFolder]; !ok {
			snippetFiles[containingFolder] = []string{}
		}

		snippetFiles[containingFolder] = append(snippetFiles[containingFolder], path)

		return nil
	})

	if err != nil {
		return err
	}

	for _, files := range snippetFiles {
		if len(files) == 1 {
			// We have no other file to compare against
			continue
		}

		var mainFile string

		for _, file := range files {
			if strings.HasSuffix(filepath.Base(file), "en-GB.json") {
				mainFile = file
			}
		}

		if len(mainFile) == 0 {
			context.AddWarning(fmt.Sprintf("No en-GB.json file found in %s, using %s", snippetFolder, files[0]))
			mainFile = files[0]
		}

		mainFileContent, err := os.ReadFile(mainFile)
		if err != nil {
			return err
		}

		for _, file := range files {
			// makes no sense to compare to ourself
			if file == mainFile {
				continue
			}

			if err := compareSnippets(mainFileContent, file, context, rootDir); err != nil {
				return err
			}
		}
	}

	return nil
}

func validateAdministrationSnippets(context *ValidationContext) {
	rootDir := context.Extension.GetRootDir()
	if err := validateAdministrationByPath(rootDir, rootDir, context); err != nil {
		return
	}

	for _, extraBundle := range context.Extension.GetExtensionConfig().Build.ExtraBundles {
		bundlePath := rootDir

		if extraBundle.Path != "" {
			bundlePath = path.Join(bundlePath, extraBundle.Path)
		} else {
			bundlePath = path.Join(bundlePath, extraBundle.Name)
		}

		if err := validateAdministrationByPath(bundlePath, rootDir, context); err != nil {
			return
		}
	}
}

func validateAdministrationByPath(extensionRoot, rootDir string, context *ValidationContext) error {
	adminFolder := path.Join(extensionRoot, "Resources", "app", "administration")

	if _, err := os.Stat(adminFolder); err != nil {
		return nil //nolint:nilerr
	}

	snippetFiles := make(map[string][]string)

	err := filepath.WalkDir(adminFolder, func(path string, d os.DirEntry, err error) error {
		if d.IsDir() {
			return nil
		}

		if filepath.Ext(path) != ".json" {
			return nil
		}

		containingFolder := filepath.Dir(path)

		if filepath.Base(containingFolder) != "snippet" {
			return nil
		}

		if _, ok := snippetFiles[containingFolder]; !ok {
			snippetFiles[containingFolder] = []string{}
		}

		snippetFiles[containingFolder] = append(snippetFiles[containingFolder], path)

		return nil
	})
	if err != nil {
		return err
	}

	for folder, files := range snippetFiles {
		if len(files) == 1 {
			// We have no other file to compare against
			continue
		}

		var mainFile string

		for _, file := range files {
			if strings.HasSuffix(filepath.Base(file), "en-GB.json") {
				mainFile = file
			}
		}

		if len(mainFile) == 0 {
			context.AddWarning(fmt.Sprintf("No en-GB.json file found in %s, using %s", folder, files[0]))
			mainFile = files[0]
		}

		mainFileContent, err := os.ReadFile(mainFile)
		if err != nil {
			return err
		}

		for _, file := range files {
			// makes no sense to compare to ourself
			if file == mainFile {
				continue
			}

			if err := compareSnippets(mainFileContent, file, context, rootDir); err != nil {
				return err
			}
		}
	}

	return nil
}

func compareSnippets(mainFile []byte, file string, context *ValidationContext, extensionRoot string) error {
	checkFile, err := os.ReadFile(file)
	if err != nil {
		return err
	}

	compare, err := jsondiff.CompareJSON(mainFile, checkFile)
	if err != nil {
		return err
	}

	for _, diff := range compare {
		normalizedPath := strings.ReplaceAll(file, extensionRoot+"/", "")

		if diff.Type == jsondiff.OperationReplace && reflect.TypeOf(diff.OldValue) != reflect.TypeOf(diff.Value) {
			context.AddError(fmt.Sprintf("Snippet file: %s, key: %s, has the type %s, but in the main language it is %s", normalizedPath, diff.Path, reflect.TypeOf(diff.OldValue), reflect.TypeOf(diff.Value)))
			continue
		}

		if diff.Type == jsondiff.OperationAdd {
			context.AddError(fmt.Sprintf("Snippet file: %s, missing key \"%s\" in this snippet file, but defined in the main language", normalizedPath, diff.Path))
			continue
		}

		if diff.Type == jsondiff.OperationRemove {
			context.AddError(fmt.Sprintf("Snippet file: %s, key: %s, is not defined in the main language", normalizedPath, diff.Path))
			continue
		}
	}

	return nil
}
