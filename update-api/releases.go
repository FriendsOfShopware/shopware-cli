package update_api

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
)

type ShopwareInstallRelease struct {
	Version string `json:"version"`
	Uri     string `json:"uri"`
	Size    string `json:"size"`
	Sha1    string `json:"sha1"`
	Sha256  string `json:"sha256"`
}

func GetLatestReleases(ctx context.Context) ([]ShopwareInstallRelease, error) {
	resp, err := http.NewRequestWithContext(ctx, "GET", "https://update-api.shopware.com/v1/releases/install?major=6", nil)

	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	content, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return nil, err
	}

	var releases []ShopwareInstallRelease

	if err := json.Unmarshal(content, &releases); err != nil {
		return nil, err
	}

	return releases, nil
}
