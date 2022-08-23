package update_api

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
)

type ShopwareInstallRelease struct {
	Version string `json:"version"`
	Uri     string `json:"uri"`
	Size    string `json:"size"`
	Sha1    string `json:"sha1"`
	Sha256  string `json:"sha256"`
}

func GetLatestReleases(ctx context.Context) ([]*ShopwareInstallRelease, error) {
	r, err := http.NewRequestWithContext(ctx, "GET", "https://update-api.shopware.com/v1/releases/install?major=6", nil)

	if err != nil {
		return nil, err
	}

	resp, err := http.DefaultClient.Do(r)

	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	content, err := io.ReadAll(resp.Body)

	if err != nil {
		return nil, err
	}

	var releases []*ShopwareInstallRelease

	if err := json.Unmarshal(content, &releases); err != nil {
		return nil, err
	}

	return releases, nil
}
