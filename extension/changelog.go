package extension

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/sirupsen/logrus"
)

func parseMarkdownChangelogInPath(path string) (map[string]map[string]string, error) {
	files, err := filepath.Glob(fmt.Sprintf("%s/CHANGELOG*.md", path))

	if err != nil {
		return nil, err
	}

	changelogs := make(map[string]map[string]string)

	for _, file := range files {
		language := strings.Trim(strings.ReplaceAll(strings.ReplaceAll(filepath.Base(file), "CHANGELOG", ""), ".md", ""), "_")

		if len(language) == 0 {
			language = "en-GB"
		}

		content, err := os.ReadFile(file)

		if err != nil {
			return nil, fmt.Errorf("parseMarkdownChangelogInPath: %v", err)
		}

		changelogs[language] = parseMarkdownChangelog(string(content))
	}

	return changelogs, nil
}

func parseMarkdownChangelog(content string) map[string]string {
	versions := make(map[string]string)
	currentVersion := ""
	versionText := ""

	for _, line := range strings.Split(content, "\n") {
		if strings.HasPrefix(line, "#") {
			if len(currentVersion) > 0 && len(versionText) > 0 {
				versions[currentVersion] = versionText
			}

			currentVersion = strings.Trim(strings.TrimPrefix(line, "#"), " ")
			versionText = ""
		} else if strings.HasPrefix(line, "-") || strings.HasPrefix(line, "*") {
			versionText = strings.Trim(versionText+line[1:]+"\n", " ")
		}
	}

	versions[currentVersion] = versionText

	return versions
}

func parseExtensionMarkdownChangelog(ext Extension) (*extensionTranslated, error) {
	v, err := ext.GetVersion()
	if err != nil {
		return nil, err
	}

	changelogs, err := parseMarkdownChangelogInPath(ext.GetPath())
	if err != nil {
		return nil, err
	}

	changelogEn, ok := changelogs["en-GB"]
	if !ok {
		return nil, fmt.Errorf("english changelog in version %s is missing", v.String())
	}

	changelogEnVersion, ok := changelogEn[v.String()]
	if !ok {
		return nil, fmt.Errorf("english changelog is missing")
	}

	changelogDe, ok := changelogs["de-DE"]
	if !ok {
		logrus.Debugf("german changelog is missing. using english as fallback")
		changelogDe = changelogEn
	}

	changelogDeVersion, ok := changelogDe[v.String()]
	if !ok {
		return nil, fmt.Errorf("german changelog in version %s is missing", v.String())
	}

	return &extensionTranslated{German: changelogDeVersion, English: changelogEnVersion}, nil
}
