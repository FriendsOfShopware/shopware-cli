package extension

import (
	"encoding/json"
	"fmt"
	"os"
)

func validateTheme(ctx *ValidationContext) {
	themeJSONPath := fmt.Sprintf("%s/src/Resources/theme.json", ctx.Extension.GetPath())

	if _, err := os.Stat(themeJSONPath); !os.IsNotExist(err) {
		content, err := os.ReadFile(themeJSONPath)
		if err != nil {
			ctx.AddError("Invalid theme.json")
			return
		}

		var theme themeJSON
		err = json.Unmarshal(content, &theme)

		if err != nil {
			ctx.AddError("Cannot decode theme.json")
			return
		}

		if len(theme.PreviewMedia) == 0 {
			ctx.AddError("Required field \"previewMedia\" in theme.json is not in")
			return
		}

		expectedMediaPath := fmt.Sprintf("%s/src/Resources/%s", ctx.Extension.GetPath(), theme.PreviewMedia)

		if _, err := os.Stat(expectedMediaPath); os.IsNotExist(err) {
			ctx.AddError(fmt.Sprintf("Theme preview image file is expected to be placed at %s, but not found there.", expectedMediaPath))
		}
	}
}

type themeJSON struct {
	PreviewMedia string `json:"previewMedia"`
}
