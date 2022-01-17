package extension

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
)

func validateTheme(ctx *validationContext) {
	themeJSONPath := fmt.Sprintf("%s/src/Resources/theme.json", ctx.Extension.GetPath())

	if _, err := os.Stat(themeJSONPath); !os.IsNotExist(err) {
		content, err := ioutil.ReadFile(themeJSONPath)

		if err != nil {
			log.Fatalln(err)
		}

		var theme themeJSON
		err = json.Unmarshal(content, &theme)

		if err != nil {
			log.Fatalln(err)
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
