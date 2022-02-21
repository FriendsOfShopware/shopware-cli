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

			for _, scRow := range salesChannelResponse.Data {
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
			_, ok := operation.SystemSettings[config.SalesChannel]

			if !ok {
				operation.SystemSettings[config.SalesChannel] = map[string]interface{}{}
			}

			foundKey := false

			for _, existingConfig := range currentConfig.Data {
				if existingConfig["configurationKey"] == newK {
					foundKey = true

					encodedSource, _ := json.Marshal(existingConfig["configurationValue"])
					encodedTarget, _ := json.Marshal(newV)

					if string(encodedSource) != string(encodedTarget) {
						operation.SystemSettings[config.SalesChannel][newK] = newV
					}

					break
				}
			}

			if !foundKey {
				operation.SystemSettings[config.SalesChannel][newK] = newV
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
	salesChannelList = append(salesChannelList, salesChannelResponse.Data...)

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

		for _, record := range sysConfigs.Data {
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
