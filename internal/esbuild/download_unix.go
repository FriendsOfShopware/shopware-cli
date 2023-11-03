//go:build !windows

package esbuild

import (
	"archive/tar"
	"compress/gzip"
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

func downloadDartSass(ctx context.Context, cacheDir string) error {
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

	request, _ := http.NewRequestWithContext(ctx, "GET", fmt.Sprintf("https://github.com/sass/dart-sass/releases/download/%s/dart-sass-%s-%s-%s.tar.gz", dartSassVersion, dartSassVersion, osType, arch), nil)

	tarFile, err := http.DefaultClient.Do(request)
	if err != nil {
		return fmt.Errorf("cannot download dart-sass: %w", err)
	}

	defer tarFile.Body.Close()

	if tarFile.StatusCode != 200 {
		return fmt.Errorf("cannot download dart-sass: %s with http code %s", tarFile.Request.URL, tarFile.Status)
	}

	uncompressedStream, err := gzip.NewReader(tarFile.Body)
	if err != nil {
		return fmt.Errorf("cannot open gzip tar file: %w", err)
	}

	tarReader := tar.NewReader(uncompressedStream)

	for {
		header, err := tarReader.Next()

		if err == io.EOF {
			break
		}

		name := strings.TrimPrefix(header.Name, "dart-sass/")

		if strings.Contains(name, "..") {
			continue
		}

		folder := filepath.Join(cacheDir, filepath.Dir(name))
		file := filepath.Join(cacheDir, name)

		if !strings.HasSuffix(folder, ".") {
			if _, err := os.Stat(folder); os.IsNotExist(err) {
				if err := os.MkdirAll(folder, os.ModePerm); err != nil {
					return err
				}
			}
		}

		outFile, err := os.Create(file)
		if err != nil {
			return fmt.Errorf("cannot create dart-sass in temp: %w", err)
		}
		if _, err := io.CopyN(outFile, tarReader, header.Size); err != nil {
			return fmt.Errorf("cannot copy dart-sass in temp: %w", err)
		}
		if err := outFile.Close(); err != nil {
			return fmt.Errorf("cannot close dart-sass in temp: %w", err)
		}

		if err := os.Chmod(file, os.FileMode(header.Mode)); err != nil {
			return fmt.Errorf("cannot chmod dart-sass in temp: %w", err)
		}
	}

	return nil
}
