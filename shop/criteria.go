package shop

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"io/ioutil"
	"net/http"
	"strings"
)

const (
	TotalCountModeDefault  = 0
	TotalCountModeExact    = 1
	TotalCountModeNextPage = 2

	SearchFilterTypeEquals    = "equals"
	SearchFilterTypeEqualsAny = "equalsAny"

	SearchSortDirectionAscending  = "ASC"
	SearchSortDirectionDescending = "DESC"
)

type Criteria struct {
	Includes       map[string][]string `json:"includes,omitempty"`
	Page           int64               `json:"page,omitempty"`
	Limit          int64               `json:"limit,omitempty"`
	IDs            []string            `json:"ids,omitempty"`
	Filter         []CriteriaFilter    `json:"filter,omitempty"`
	PostFilter     []CriteriaFilter    `json:"postFilter,omitempty"`
	Sort           []CriteriaSort      `json:"sort,omitempty"`
	Term           string              `json:"term,omitempty"`
	TotalCountMode int                 `json:"totalCountMode,omitempty"`
}

type CriteriaFilter struct {
	Type  string      `json:"type"`
	Field string      `json:"field"`
	Value interface{} `json:"value"`
}

type CriteriaSort struct {
	Direction      string `json:"order"`
	Field          string `json:"field"`
	NaturalSorting bool   `json:"naturalSorting"`
}

type SearchResponse struct {
	Total        int64         `json:"total"`
	Data         []interface{} `json:"data"`
	Aggregations interface{}   `json:"aggregations"`
}

func (c Client) Search(ctx context.Context, entity string, criteria Criteria) (*SearchResponse, error) {
	content, err := json.Marshal(criteria)

	if err != nil {
		return nil, err
	}

	req, err := c.newRequest(ctx, http.MethodPost, fmt.Sprintf("/api/search/%s", entity), bytes.NewReader(content))

	if err != nil {
		return nil, err
	}

	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)

	if err != nil {
		return nil, errors.Wrap(err, "Search")
	}

	defer resp.Body.Close()

	content, err = ioutil.ReadAll(resp.Body)

	if err != nil {
		return nil, errors.Wrap(err, "Search")
	}

	var result *SearchResponse
	if err := json.Unmarshal(content, &result); err != nil {
		return nil, err
	}

	return result, nil
}

func (c Client) SearchAll(ctx context.Context, entity string, criteria Criteria) (*SearchResponse, error) {
	entity = strings.ReplaceAll(entity, "_", "-")

	criteria.Page = 1
	criteria.Limit = 200

	result, err := c.Search(ctx, entity, criteria)

	if err != nil {
		return nil, err
	}

	for {
		criteria.Page++

		pagedResult, err := c.Search(ctx, entity, criteria)

		if err != nil {
			return nil, err
		}

		if len(pagedResult.Data) == 0 {
			break
		}

		result.Data = append(result.Data, pagedResult.Data...)
	}

	return result, nil
}

func (c Client) Sync(ctx context.Context, payload map[string]SyncOperation) error {
	content, err := json.Marshal(payload)

	if err != nil {
		return err
	}

	req, err := c.newRequest(ctx, http.MethodPost, "/api/_action/sync", bytes.NewReader(content))

	if err != nil {
		return err
	}

	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)

	if err != nil {
		return errors.Wrap(err, "Sync")
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent && resp.StatusCode != http.StatusOK {
		content, err := ioutil.ReadAll(resp.Body)

		if err != nil {
			return err
		}

		return fmt.Errorf("Sync: got http code %d from api: %s", resp.StatusCode, string(content))
	}

	return nil
}

func (c Client) UpdateSystemConfig(ctx context.Context, payload string) error {
	req, err := c.newRequest(ctx, http.MethodPost, "/api/_action/system-config/batch", strings.NewReader(payload))

	if err != nil {
		return err
	}

	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)

	if err != nil {
		return errors.Wrap(err, "UpdateSystemConfig")
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent && resp.StatusCode != http.StatusOK {
		content, err := ioutil.ReadAll(resp.Body)

		if err != nil {
			return err
		}

		return fmt.Errorf("UpdateSystemConfig: got http code %d from api: %s", resp.StatusCode, string(content))
	}

	return nil
}

type SyncOperation struct {
	Entity  string                   `json:"entity"`
	Action  string                   `json:"action"`
	Payload []map[string]interface{} `json:"payload"`
}
