package esbuild

import (
	"regexp"
	"strings"
)

var matchLetter = regexp.MustCompile(`[A-Z]`)

// @see https://github.com/symfony/symfony/blob/6.3/src/Symfony/Component/Serializer/NameConverter/CamelCaseToSnakeCaseNameConverter.php#L31
func toSnakeCase(str string) string {
	converted := matchLetter.ReplaceAllStringFunc(str, func(match string) string {
		return "_" + strings.ToLower(match)
	})
	return converted[1:]
}
