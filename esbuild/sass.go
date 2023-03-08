package esbuild

import (
	"archive/tar"
	"compress/gzip"
	"context"
	_ "embed"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/FriendsOfShopware/shopware-cli/logging"
)

const dartSassVersion = "1.57.1"

//go:embed static/variables.scss
var scssVariables []byte

//go:embed static/mixins.scss
var scssMixins []byte

func downloadDartSass(ctx context.Context) (string, error) {
	if path, err := exec.LookPath("dart-sass-embedded"); err == nil {
		return path, nil
	}

	cacheDir, err := os.UserCacheDir()
	if err != nil {
		cacheDir = "/tmp"
	}

	cacheDir += "/dart-sass-embedded-" + dartSassVersion

	expectedPath := fmt.Sprintf("%s/dart-sass-embedded", cacheDir)

	if _, err := os.Stat(expectedPath); err == nil {
		return expectedPath, nil
	}

	if _, err := os.Stat(filepath.Dir(expectedPath)); os.IsNotExist(err) {
		if err := os.MkdirAll(filepath.Dir(expectedPath), os.ModePerm); err != nil {
			return "", err
		}
	}

	logging.FromContext(ctx).Infof("Downloading dart-sass")

	osType := runtime.GOOS
	arch := runtime.GOARCH

	switch runtime.GOARCH {
	case "arm64":
		arch = "arm64"
	case "amd64":
		arch = "x64"
	case "386":
		arch = "ia32"
	}

	if osType == "darwin" {
		osType = "macos"
	}

	request, _ := http.NewRequestWithContext(context.Background(), "GET", fmt.Sprintf("https://github.com/sass/dart-sass-embedded/releases/download/%s/sass_embedded-%s-%s-%s.tar.gz", dartSassVersion, dartSassVersion, osType, arch), nil)

	tarFile, err := http.DefaultClient.Do(request)
	if err != nil {
		return "", fmt.Errorf("cannot download dart-sass: %w", err)
	}

	defer tarFile.Body.Close()

	uncompressedStream, err := gzip.NewReader(tarFile.Body)
	if err != nil {
		return "", fmt.Errorf("cannot open gzip tar file: %w", err)
	}

	tarReader := tar.NewReader(uncompressedStream)

	for {
		header, err := tarReader.Next()

		if err == io.EOF {
			break
		}

		name := strings.TrimPrefix(header.Name, "sass_embedded/")
		folder := filepath.Join(cacheDir, filepath.Dir(name))
		file := filepath.Join(cacheDir, name)

		if !strings.HasSuffix(folder, ".") {
			if _, err := os.Stat(folder); os.IsNotExist(err) {
				if err := os.MkdirAll(folder, os.ModePerm); err != nil {
					return "", err
				}
			}
		}

		outFile, err := os.Create(file)
		if err != nil {
			return "", fmt.Errorf("cannot create dart-sass in temp: %w", err)
		}
		if _, err := io.CopyN(outFile, tarReader, header.Size); err != nil {
			return "", fmt.Errorf("cannot copy dart-sass in temp: %w", err)
		}
		if err := outFile.Close(); err != nil {
			return "", fmt.Errorf("cannot close dart-sass in temp: %w", err)
		}

		if err := os.Chmod(file, os.FileMode(header.Mode)); err != nil {
			return "", fmt.Errorf("cannot chmod dart-sass in temp: %w", err)
		}
	}

	return expectedPath, nil
}
