package shop

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"io/ioutil"
	"net/http"
)

func (c Client) GetThemeConfiguration(ctx context.Context, themeId string) (*ThemeConfiguration, error) {
	req, err := c.NewRequest(ctx, http.MethodGet, fmt.Sprintf("/api/_action/theme/%s/configuration", themeId), nil)

	if err != nil {
		return nil, err
	}

	resp, err := c.httpClient.Do(req)

	if err != nil {
		return nil, errors.Wrap(err, "GetThemeConfiguration")
	}

	defer resp.Body.Close()

	content, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.Wrap(err, "GetThemeConfiguration")
	}

	var result *ThemeConfiguration
	if err := json.Unmarshal(content, &result); err != nil {
		return nil, err
	}

	// Old shopware version, use fields instead
	if result.CurrentFields == nil {
		result.CurrentFields = &result.Fields
	}

	return result, nil
}

type ThemeConfiguration struct {
	CurrentFields *map[string]ThemeConfigValue `json:"currentFields"`
	Fields        map[string]ThemeConfigValue  `json:"fields"`
}

type ThemeUpdateRequest struct {
	Config map[string]ThemeConfigValue `json:"config"`
}

func (c Client) SaveThemeConfiguration(ctx context.Context, themeId string, update ThemeUpdateRequest) error {
	content, err := json.Marshal(update)

	if err != nil {
		return err
	}

	req, err := c.NewRequest(ctx, http.MethodPatch, fmt.Sprintf("/api/_action/theme/%s", themeId), bytes.NewReader(content))

	if err != nil {
		return err
	}

	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)

	if err != nil {
		return errors.Wrap(err, "SaveThemeConfiguration")
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent && resp.StatusCode != http.StatusOK {
		content, err := ioutil.ReadAll(resp.Body)

		if err != nil {
			return err
		}

		return fmt.Errorf("SaveThemeConfiguration: got http code %d from api: %s", resp.StatusCode, string(content))
	}

	return nil
}
