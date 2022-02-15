package project

import (
	"context"
	"encoding/json"
	log "github.com/sirupsen/logrus"
	"shopware-cli/shop"
)

type SystemConfigSync struct{}

func (s SystemConfigSync) Push(ctx context.Context, client *shop.Client, config *shop.Config, operation *ConfigSyncOperation) error {
	if config.Sync == nil {
		return nil
	}

	c := shop.Criteria{}
	c.Includes = map[string][]string{"sales_channel": {"id", "name"}}
	salesChannelResponse, err := client.SearchAll(ctx, "sales_channel", c)

	if err != nil {
		return err
	}

	for _, config := range config.Sync.Config {
		if config.SalesChannel != nil && len(*config.SalesChannel) != 32 {
			foundId := false

			for _, scRowRaw := range salesChannelResponse.Data {
				scRow := scRowRaw.(map[string]interface{})

				if *config.SalesChannel == scRow["name"] {
					val, _ := scRow["id"].(string)
					config.SalesChannel = &val

					foundId = true
				}
			}

			if !foundId {
				log.Errorf("Cannot find sales channel id for %s", *config.SalesChannel)
				continue
			}
		}

		currentConfig, err := readSystemConfig(ctx, client, config.SalesChannel)

		if err != nil {
			return err
		}

		for newK, newV := range config.Settings {
			foundKey := false

			for _, existingConfigRaw := range currentConfig.Data {
				existingConfig := existingConfigRaw.(map[string]interface{})

				if existingConfig["configurationKey"] == newK {
					foundKey = true

					encodedSource, _ := json.Marshal(existingConfig["configurationValue"])
					encodedTarget, _ := json.Marshal(newV)

					if string(encodedSource) != string(encodedTarget) {
						operation.Upsert["system_config"] = append(operation.Upsert["system_config"], map[string]interface{}{
							"id":                 existingConfig["id"],
							"configurationKey":   newK,
							"configurationValue": newV,
						})
					}

					break
				}
			}

			if !foundKey {
				operation.Upsert["system_config"] = append(operation.Upsert["system_config"], map[string]interface{}{
					"id":                 shop.NewUuid(),
					"configurationKey":   newK,
					"configurationValue": newV,
					"salesChannelId":     config.SalesChannel,
				})
			}
		}
	}

	return nil
}

func (s SystemConfigSync) Pull(ctx context.Context, client *shop.Client, config *shop.Config) error {
	config.Sync.Config = make([]shop.ConfigSyncConfig, 0)

	c := shop.Criteria{}
	c.Includes = map[string][]string{"sales_channel": {"id", "name"}}
	salesChannelResponse, err := client.SearchAll(ctx, "sales_channel", c)

	if err != nil {
		return err
	}

	salesChannelList := make([]map[string]interface{}, 0)
	salesChannelList = append(salesChannelList, nil)

	for _, row := range salesChannelResponse.Data {
		r := row.(map[string]interface{})
		salesChannelList = append(salesChannelList, r)
	}

	for _, sc := range salesChannelList {
		var sysConfigs *shop.SearchResponse
		var err error

		cfg := shop.ConfigSyncConfig{
			Settings: map[string]interface{}{},
		}

		if sc == nil {
			sysConfigs, err = readSystemConfig(ctx, client, nil)
		} else {
			scId, _ := sc["id"].(string)
			scName, _ := sc["name"].(string)

			cfg.SalesChannel = &scName

			sysConfigs, err = readSystemConfig(ctx, client, &scId)
		}

		if err != nil {
			return err
		}

		for _, recordRaw := range sysConfigs.Data {
			record := recordRaw.(map[string]interface{})

			key, _ := record["configurationKey"].(string)
			val := record["configurationValue"]

			// app system shopId
			if key == "core.app.shopId" {
				continue
			}

			cfg.Settings[key] = val
		}

		config.Sync.Config = append(config.Sync.Config, cfg)
	}

	return nil
}
