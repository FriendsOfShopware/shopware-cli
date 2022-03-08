package project

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"shopware-cli/shop"
)

type EntitySync struct{}

func (s EntitySync) Push(ctx context.Context, client *shop.Client, config *shop.Config, operation *ConfigSyncOperation) error {
	for _, entity := range config.Sync.Entity {

		if entity.Exists != nil && len(*entity.Exists) > 0 {
			criteria := make(map[string]interface{})
			criteria["filter"] = entity.Exists

			searchPayload, err := json.Marshal(criteria)

			if err != nil {
				return err
			}

			r, err := client.NewRequest(ctx, "POST", fmt.Sprintf("/api/search-ids/%s", entity.Entity), bytes.NewReader(searchPayload))

			if err != nil {
				return err
			}

			r.Header.Set("Accept", "application/json")
			r.Header.Set("Content-Type", "application/json")

			resp, err := client.Do(r)

			if err != nil {
				return err
			}

			defer resp.Body.Close()

			content, err := ioutil.ReadAll(resp.Body)

			if err != nil {
				return err
			}

			if resp.StatusCode != 200 {
				return fmt.Errorf("request failed with error: %s", string(content))
			}

			var res criteriaApiResponse
			if err := json.Unmarshal(content, &res); err != nil {
				return err
			}

			if res.Total > 0 {
				continue
			}
		}

		operation.Operations[shop.NewUuid()] = shop.SyncOperation{
			Action:  "upsert",
			Entity:  entity.Entity,
			Payload: []map[string]interface{}{entity.Payload},
		}
	}

	return nil
}

func (s EntitySync) Pull(ctx context.Context, client *shop.Client, config *shop.Config) error {
	return nil
}

type criteriaApiResponse struct {
	Total int      `json:"total"`
	Data  []string `json:"data"`
}
