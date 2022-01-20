package update_api

import (
	"encoding/json"
	"io"
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

func GetLatestReleases() ([]ShopwareInstallRelease, error) {
	resp, err := http.Get("https://update-api.shopware.com/v1/releases/install?major=6")

	if err != nil {
		return nil, err
	}

	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(resp.Body)

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
