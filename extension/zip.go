package extension

import (
	"archive/zip"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/fs"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"

	"github.com/pkg/errors"

	"github.com/FriendsOfShopware/shopware-cli/version"
)

var (
	defaultNotAllowedPaths = []string{
		".travis.yml",
		".gitlab-ci.yml",
		"build.sh",
		".editorconfig",
		".php_cs.dist",
		".php_cs.cache",
		"ISSUE_TEMPLATE.md",
		".sw-zip-blacklist",
		"tests",
		"Resources/store",
		"src/Resources/store",
		".github",
		".git",
		".shopware-extension.yml",
		"src/Resources/app/storefront/node_modules",
		"src/Resources/app/administration/node_modules",
		"src/Resources/app/node_modules",
		"var",
	}

	defaultNotAllowedFiles = []string{
		".DS_Store",
		"Thumbs.db",
		"__MACOSX",
	}

	defaultNotAllowedExtensions = []string{
		".zip",
		".tar",
		".gz",
		".phar",
		".rar",
	}
)

func Unzip(r *zip.Reader, dest string) error {
	errorFormat := "unzip: %w"

	for _, f := range r.File {
		// Store filename/path for returning and using later on
		fpath := filepath.Join(dest, f.Name) //nolint:gosec

		// Check for ZipSlip. More Info: http://bit.ly/2MsjAWE
		if !strings.HasPrefix(fpath, filepath.Clean(dest)+string(os.PathSeparator)) {
			return fmt.Errorf("Unzip: %s: illegal file path", fpath)
		}

		if f.FileInfo().IsDir() {
			// Make Folder
			_ = os.MkdirAll(fpath, os.ModePerm)
			continue
		}

		// Make File
		if err := os.MkdirAll(filepath.Dir(fpath), os.ModePerm); err != nil {
			return err
		}

		outFile, err := os.OpenFile(fpath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
		if err != nil {
			return fmt.Errorf(errorFormat, err)
		}

		rc, err := f.Open()
		if err != nil {
			return fmt.Errorf(errorFormat, err)
		}

		_, err = io.Copy(outFile, rc) //nolint:gosec

		// Close the file without defer to close before next iteration of loop
		_ = outFile.Close()
		_ = rc.Close()

		if err != nil {
			return fmt.Errorf(errorFormat, err)
		}
	}

	return nil
}

func CreateZip(baseFolder, zipFile string) error {
	// Get a Buffer to Write To
	outFile, err := os.Create(zipFile)
	if err != nil {
		return errors.Wrap(err, "create zipfile")
	}
	defer outFile.Close()

	// Create a new zip archive.
	w := zip.NewWriter(outFile)
	defer w.Close()

	return AddZipFiles(w, baseFolder, "")
}

func AddZipFiles(w *zip.Writer, basePath, baseInZip string) error {
	files, err := ioutil.ReadDir(basePath)
	if err != nil {
		return fmt.Errorf("could not zip dir, basePath: %q, baseInZip: %q, %w", basePath, baseInZip, err)
	}

	for _, file := range files {
		if file.IsDir() {
			// Add files of directory recursively
			if err = AddZipFiles(w, filepath.Join(basePath, file.Name()), filepath.Join(baseInZip, file.Name())); err != nil {
				return err
			}
		} else {
			if err = addFileToZip(w, filepath.Join(basePath, file.Name()), filepath.Join(baseInZip, file.Name())); err != nil {
				return err
			}
		}
	}

	return nil
}

func CleanupExtensionFolder(path string, additionalPaths []string) error {
	defaultNotAllowedPaths = append(defaultNotAllowedPaths, additionalPaths...)

	for _, folder := range defaultNotAllowedPaths {
		if _, err := os.Stat(path + folder); !os.IsNotExist(err) {
			err := os.RemoveAll(path + folder)
			if err != nil {
				return err
			}
		}
	}

	return filepath.Walk(path, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// When we delete a folder, this function will be called also the files in it
		if _, err := os.Stat(path); os.IsNotExist(err) {
			return nil
		}

		base := filepath.Base(path)

		for _, file := range defaultNotAllowedFiles {
			if file == base {
				return os.RemoveAll(path)
			}
		}

		for _, ext := range defaultNotAllowedExtensions {
			if strings.HasSuffix(base, ext) {
				return os.RemoveAll(path)
			}
		}

		return nil
	})
}

func PrepareFolderForZipping(ctx context.Context, path string, ext Extension, extCfg *Config) error {
	errorFormat := "PrepareFolderForZipping: %v"
	composerJSONPath := filepath.Join(path, "composer.json")
	composerLockPath := filepath.Join(path, "composer.lock")

	if _, err := os.Stat(composerJSONPath); os.IsNotExist(err) {
		return nil
	}

	content, err := ioutil.ReadFile(composerJSONPath)

	if err != nil {
		return fmt.Errorf(errorFormat, err)
	}

	var composer map[string]interface{}
	err = json.Unmarshal(content, &composer)

	if err != nil {
		return fmt.Errorf(errorFormat, err)
	}

	// Add replacements
	composer, err = addComposerReplacements(ctx, composer, ext)
	if err != nil {
		return errors.Wrap(err, "add composer replacements")
	}

	filtered := filterRequires(composer, extCfg)

	if len(filtered["require"].(map[string]interface{})) == 0 {
		return nil
	}

	// Remove the composer.lock
	if _, err := os.Stat(composerLockPath); !os.IsNotExist(err) {
		err := os.Remove(composerLockPath)
		if err != nil {
			return fmt.Errorf(errorFormat, err)
		}
	}

	newContent, err := json.Marshal(&composer)

	if err != nil {
		return fmt.Errorf(errorFormat, err)
	}

	err = ioutil.WriteFile(composerJSONPath, newContent, 0644) //nolint:gosec
	if err != nil {
		// Revert on failure
		_ = ioutil.WriteFile(composerJSONPath, content, 0644) //nolint:gosec
		return fmt.Errorf(errorFormat, err)
	}

	// Execute composer in this directory
	composerInstallCmd := exec.Command("composer", "install", "-d", path, "--no-dev", "-n", "-o")
	composerInstallCmd.Stdout = os.Stdout
	composerInstallCmd.Stderr = os.Stderr
	err = composerInstallCmd.Run()
	if err != nil {
		// Revert on failure
		_ = ioutil.WriteFile(composerJSONPath, content, 0644) //nolint:gosec
		return fmt.Errorf(errorFormat, err)
	}

	_ = ioutil.WriteFile(composerJSONPath, content, 0644) //nolint:gosec

	return nil
}

func addFileToZip(zipWriter *zip.Writer, sourcePath string, zipPath string) error {
	zipErrorFormat := "could not zip file, sourcePath: %q, zipPath: %q, %w"

	dat, err := ioutil.ReadFile(sourcePath)
	if err != nil {
		return fmt.Errorf(zipErrorFormat, sourcePath, zipPath, err)
	}

	f, err := zipWriter.Create(zipPath)
	if err != nil {
		return fmt.Errorf(zipErrorFormat, sourcePath, zipPath, err)
	}

	if _, err := f.Write(dat); err != nil {
		return fmt.Errorf(zipErrorFormat, sourcePath, zipPath, err)
	}

	return nil
}

func filterRequires(composer map[string]interface{}, extCfg *Config) map[string]interface{} {
	if _, ok := composer["provide"]; !ok {
		composer["provide"] = make(map[string]interface{})
	}
	if _, ok := composer["require"]; !ok {
		composer["require"] = make(map[string]interface{})
	}

	provide := composer["provide"]
	require := composer["require"]

	keys := []string{"shopware/platform", "shopware/core", "shopware/shopware", "shopware/storefront", "shopware/administration", "shopware/elasticsearch", "composer/installers"}
	if extCfg != nil {
		keys = append(keys, extCfg.Build.Zip.Composer.ExcludedPackages...)
	}

	for _, key := range keys {
		if _, ok := require.(map[string]interface{})[key]; ok {
			delete(require.(map[string]interface{}), key)
			provide.(map[string]interface{})[key] = "*"
		}
	}

	return composer
}

func addComposerReplacements(ctx context.Context, composer map[string]interface{}, ext Extension) (map[string]interface{}, error) {
	if _, ok := composer["replace"]; !ok {
		composer["replace"] = make(map[string]interface{})
	}

	if _, ok := composer["require"]; !ok {
		composer["require"] = make(map[string]interface{})
	}

	replace := composer["replace"]
	require := composer["require"]

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "https://swagger.docs.fos.gg/composer/versions.json", nil)
	if err != nil {
		return nil, errors.Wrap(err, "create composer version request")
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, errors.Wrap(err, "fetch composer versions")
	}
	defer resp.Body.Close()

	versionString, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.Wrap(err, "read version body")
	}

	var versions []string
	err = json.Unmarshal(versionString, &versions)
	if err != nil {
		return nil, errors.Wrap(err, "unmarshal composer versions")
	}

	versionConstraint, err := ext.GetShopwareVersionConstraint()
	if err != nil {
		return nil, errors.Wrap(err, "get shopware version constraint")
	}

	minVersion := getMinMatchingVersion(versionConstraint, versions)
	components := []string{"core", "administration", "storefront", "administration"}

	for _, component := range components {
		packageName := fmt.Sprintf("shopware/%s", component)

		if _, ok := require.(map[string]interface{})[packageName]; ok {
			req, err := http.NewRequestWithContext(ctx, http.MethodGet, fmt.Sprintf("https://swagger.docs.fos.gg/composer/%s/%s.json", minVersion, component), nil)
			if err != nil {
				return nil, errors.Wrap(err, "create component request")
			}

			resp, err := http.DefaultClient.Do(req)
			if err != nil {
				return nil, fmt.Errorf("get packte version %s: %w", component, err)
			}

			composerPartByte, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				return nil, errors.Wrap(err, "read component version body")
			}

			_ = resp.Body.Close()

			var composerPart map[string]string
			err = json.Unmarshal(composerPartByte, &composerPart)
			if err != nil {
				return nil, errors.Wrap(err, "unmarshal component version")
			}

			for k, v := range composerPart {
				replace.(map[string]interface{})[k] = v

				delete(require.(map[string]interface{}), k)
			}
		}
	}

	return composer, nil
}

func getMinMatchingVersion(constraint *version.Constraints, versions []string) string {
	vs := make([]*version.Version, 0)

	for _, r := range versions {
		v, err := version.NewVersion(r)
		if err != nil {
			continue
		}

		vs = append(vs, v)
	}

	sort.Sort(version.Collection(vs))

	for _, v := range vs {
		if constraint.Check(v) {
			return v.String()
		}
	}

	return vs[0].String()
}
