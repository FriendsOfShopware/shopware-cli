package extension

import (
	"fmt"
	"io/fs"
	"path/filepath"
	"strings"
)

type validationContext struct {
	Extension Extension
	errors    *[]string
}

func newValidationContext(ext Extension) *validationContext {
	context := validationContext{Extension: ext}
	str := make([]string, 0)

	context.errors = &str

	return &context
}

func (c validationContext) AddError(message string) {
	*c.errors = append(*c.errors, message)
}

func (c validationContext) HasErrors() bool {
	return len(*c.errors) != 0
}

func (c validationContext) Errors() []string {
	return *c.errors
}

type Validator interface {
	Validate(context *validationContext)
}

func RunValidation(ext Extension) *validationContext {
	context := newValidationContext(ext)

	validators := []Validator{generalChecker{}}

	for _, validator := range validators {
		validator.Validate(context)
	}

	return context
}

type generalChecker struct{}

func (c generalChecker) Validate(context *validationContext) {
	_, versionErr := context.Extension.GetVersion()
	_, nameErr := context.Extension.GetName()
	_, shopwareVersionErr := context.Extension.GetShopwareVersionConstraint()

	if versionErr != nil {
		context.AddError(versionErr.Error())
	}

	if nameErr != nil {
		context.AddError(nameErr.Error())
	}

	if shopwareVersionErr != nil {
		context.AddError(shopwareVersionErr.Error())
	}

	_ = filepath.Walk(context.Extension.GetPath(), func(path string, info fs.FileInfo, err error) error {
		name := filepath.Base(path)

		if name == ".." {
			context.AddError("Path travel detected in zip file")
		}

		for _, file := range defaultNotAllowedPaths {
			if strings.HasPrefix(path, file) {
				context.AddError(fmt.Sprintf("file %s is not allowed in the zip file", path))
			}
		}

		for _, file := range defaultNotAllowedFiles {
			if file == name {
				context.AddError(fmt.Sprintf("file %s is not allowed in the zip file", path))
			}
		}

		for _, ext := range defaultNotAllowedExtensions {
			if strings.HasSuffix(name, ext) {
				context.AddError(fmt.Sprintf("file %s is not allowed in the zip file", path))
			}
		}

		return nil
	})
}
