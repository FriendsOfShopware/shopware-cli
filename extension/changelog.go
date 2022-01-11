package extension

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"
)

func parseMarkdownChangelogInPath(path string) (map[string]map[string]string, error) {
	files, err := filepath.Glob(fmt.Sprintf("%s/CHANGELOG*.md", path))

	if err != nil {
		return nil, err
	}

	changelogs := make(map[string]map[string]string, 0)

	for _, file := range files {
		language := strings.Trim(strings.ReplaceAll(strings.ReplaceAll(filepath.Base(file), "CHANGELOG", ""), ".md", ""), "_")

		content, err := ioutil.ReadFile(file)

		if err != nil {
			return nil, fmt.Errorf("parseMarkdownChangelogInPath: %v", err)
		}

		changelogs[language] = parseMarkdownChangelog(string(content))
	}

	return changelogs, nil
}

func parseMarkdownChangelog(content string) map[string]string {
	versions := make(map[string]string, 0)
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
