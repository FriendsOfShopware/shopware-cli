package shop

import (
	"context"
	"encoding/json"
	"github.com/pkg/errors"
	"io/ioutil"
	"net/http"
)

type InfoResponse struct {
	Version         string `json:"version"`
	VersionRevision string `json:"versionRevision"`
	AdminWorker     struct {
		EnableAdminWorker bool     `json:"enableAdminWorker"`
		Transports        []string `json:"transports"`
	} `json:"adminWorker"`
	Bundles  map[string]infoResponseBundle `json:"bundles"`
	Settings struct {
		EnableURLFeature bool `json:"enableUrlFeature"`
	} `json:"settings"`
}

type infoResponseBundle struct {
	CSS []string `json:"css"`
	Js  []string `json:"js"`
}

func (r InfoResponse) IsCloudShop() bool {
	_, ok := r.Bundles["SaasRufus"]

	return ok
}

func (c Client) Info(ctx context.Context) (*InfoResponse, error) {
	req, err := c.newRequest(ctx, http.MethodGet, "/api/_info/config", nil)
	if err != nil {
		return nil, errors.Wrap(err, "Info")
	}

	resp, err := c.httpClient.Do(req)

	if err != nil {
		return nil, errors.Wrap(err, "Info")
	}

	defer resp.Body.Close()

	content, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return nil, errors.Wrap(err, "Info")
	}

	var infoResponse *InfoResponse
	if err := json.Unmarshal(content, &infoResponse); err != nil {
		return nil, err
	}

	return infoResponse, nil
}
