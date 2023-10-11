package phplint

import (
	"context"
	"io/fs"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/tetratelabs/wazero"
	"github.com/tetratelabs/wazero/sys"
)

type LintError struct {
	File    string
	Message string
}

type LintErrors []LintError

func LintFolder(ctx context.Context, phpVersion, folder string) (LintErrors, error) {
	wasmFile, err := findPHPWasmFile(ctx, phpVersion)
	if err != nil {
		return nil, err
	}

	wasmRuntime, err := getWazeroRuntime(ctx)
	if err != nil {
		return nil, err
	}

	defer wasmRuntime.Close(ctx)

	wasmCompiled, _ := wasmRuntime.CompileModule(ctx, wasmFile)

	dirFs := os.DirFS(folder)

	paths := make([]string, 0)

	_ = filepath.Walk(folder, func(path string, _ fs.FileInfo, _ error) error {
		name := filepath.Base(path)

		if strings.HasSuffix(name, ".php") {
			paths = append(paths, path)
		}

		return nil
	})

	errorsChain := make(chan *LintError, len(paths))

	runtime.GOMAXPROCS(2)

	for _, file := range paths {
		go func(file string) {
			file, _ = filepath.Rel(folder, file)
			stderr := new(strings.Builder)

			config := wazero.NewModuleConfig().
				WithStderr(stderr).
				WithStdout(stderr).
				WithArgs("php", "-l", file).
				WithFS(dirFs)

			if wasmModule, err := wasmRuntime.InstantiateModule(ctx, wasmCompiled, config); err != nil {
				if exitErr, ok := err.(*sys.ExitError); ok && exitErr.ExitCode() != 0 {
					errorsChain <- &LintError{
						File:    file,
						Message: stderr.String(),
					}
				} else if !ok {
					errorsChain <- &LintError{
						File:    file,
						Message: err.Error(),
					}
				} else {
					errorsChain <- nil
				}

				if wasmModule != nil {
					wasmModule.Close(ctx)
				}
			} else {
				wasmModule.Close(ctx)
				errorsChain <- nil
			}
		}(file)
	}

	listOfErrors := make(LintErrors, 0)

	for i := 0; i < len(paths); i++ {
		err := <-errorsChain
		if err != nil {
			listOfErrors = append(listOfErrors, *err)
		}
	}

	return listOfErrors, nil
}
