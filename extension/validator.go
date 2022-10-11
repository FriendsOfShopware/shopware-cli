package extension

import (
	"fmt"
	"io/fs"
	"path/filepath"
	"strings"
)

type validationContext struct {
	Extension Extension
	errors    []string
	warnings  []string
}

func newValidationContext(ext Extension) *validationContext {
	return &validationContext{Extension: ext}
}

func (c *validationContext) AddError(message string) {
	c.errors = append(c.errors, message)
}

func (c validationContext) HasErrors() bool {
	return len(c.errors) > 0
}

func (c validationContext) Errors() []string {
	return c.errors
}

func (c *validationContext) AddWarning(message string) {
	c.warnings = append(c.warnings, message)
}

func (c validationContext) HasWarnings() bool {
	return len(c.warnings) > 0
}

func (c validationContext) Warnings() []string {
	return c.warnings
}

func RunValidation(ext Extension) *validationContext {
	context := newValidationContext(ext)

	runDefaultValidate(context)
	ext.Validate(context)

	return context
}

func runDefaultValidate(context *validationContext) {
	_, versionErr := context.Extension.GetVersion()
	name, nameErr := context.Extension.GetName()
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

	if len(name) == 0 {
		context.AddError("Extension name cannot be empty")
	}

	notAllowedErrorFormat := "file %s is not allowed in the zip file"
	_ = filepath.Walk(context.Extension.GetPath(), func(path string, info fs.FileInfo, err error) error {
		name := filepath.Base(path)

		if name == ".." {
			context.AddError("Path travel detected in zip file")
		}

		for _, file := range defaultNotAllowedPaths {
			if strings.HasPrefix(path, file) {
				context.AddError(fmt.Sprintf(notAllowedErrorFormat, path))
			}
		}

		for _, file := range defaultNotAllowedFiles {
			if file == name {
				context.AddError(fmt.Sprintf(notAllowedErrorFormat, path))
			}
		}

		for _, ext := range defaultNotAllowedExtensions {
			if strings.HasSuffix(name, ext) {
				context.AddError(fmt.Sprintf(notAllowedErrorFormat, path))
			}
		}

		return nil
	})

	metaData := context.Extension.GetMetaData()

	if len(metaData.Label.German) == 0 {
		context.AddError("label is not translated in german")
	}

	if len(metaData.Label.English) == 0 {
		context.AddError("label is not translated in english")
	}

	if len(metaData.Description.German) == 0 {
		context.AddError("description is not translated in german")
	}

	if len(metaData.Description.English) == 0 {
		context.AddError("description is not translated in english")
	}

	if len(metaData.Description.German) < 150 || len(metaData.Description.German) > 185 {
		context.AddError(fmt.Sprintf("the %s description with length of %d should have a length from 150 up to 185 characters.", "german", len(metaData.Description.German)))
	}

	if len(metaData.Description.English) < 150 || len(metaData.Description.English) > 185 {
		context.AddError(fmt.Sprintf("the %s description with length of %d should have a length from 150 up to 185 characters.", "english", len(metaData.Description.English)))
	}
}
