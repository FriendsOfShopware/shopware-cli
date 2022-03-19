package project

import (
	"encoding/json"
	adminSdk "github.com/friendsofshopware/go-shopware-admin-api-sdk"
	log "github.com/sirupsen/logrus"
	"shopware-cli/shop"
)

type SystemConfigSync struct{}

func (s SystemConfigSync) Push(ctx adminSdk.ApiContext, client *adminSdk.Client, config *shop.Config, operation *ConfigSyncOperation) error {
	if config.Sync == nil {
		return nil
	}

	c := adminSdk.Criteria{}
	c.Includes = map[string][]string{"sales_channel": {"id", "name"}}
	salesChannelResponse, _, err := client.Repository.SalesChannel.SearchAll(ctx, c)

	if err != nil {
		return err
	}

	for _, config := range config.Sync.Config {
		if config.SalesChannel != nil && len(*config.SalesChannel) != 32 {
			foundId := false

			for _, scRow := range salesChannelResponse.Data {
				if *config.SalesChannel == scRow.Name {
					config.SalesChannel = &scRow.Id //nolint:exportloopref

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
				if existingConfig.ConfigurationKey == newK {
					foundKey = true

					encodedSource, _ := json.Marshal(existingConfig.ConfigurationValue)
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

func (s SystemConfigSync) Pull(ctx adminSdk.ApiContext, client *adminSdk.Client, config *shop.Config) error {
	config.Sync.Config = make([]shop.ConfigSyncConfig, 0)

	c := adminSdk.Criteria{}
	c.Includes = map[string][]string{"sales_channel": {"id", "name"}}
	salesChannelResponse, _, err := client.Repository.SalesChannel.SearchAll(ctx, c)

	if err != nil {
		return err
	}

	salesChannelList := make([]adminSdk.SalesChannel, 0)
	salesChannelList = append(salesChannelList, adminSdk.SalesChannel{Id: ""})
	salesChannelList = append(salesChannelList, salesChannelResponse.Data...)

	for _, sc := range salesChannelList {
		var sysConfigs *adminSdk.SystemConfigCollection
		var err error

		cfg := shop.ConfigSyncConfig{
			Settings: map[string]interface{}{},
		}

		if sc.Id == "" {
			sysConfigs, err = readSystemConfig(ctx, client, nil)
		} else {
			cfg.SalesChannel = &sc.Name

			sysConfigs, err = readSystemConfig(ctx, client, &sc.Id)
		}

		if err != nil {
			return err
		}

		for _, record := range sysConfigs.Data {
			// app system shopId
			if record.ConfigurationKey == "core.app.shopId" {
				continue
			}

			cfg.Settings[record.ConfigurationKey] = record.ConfigurationValue
		}

		config.Sync.Config = append(config.Sync.Config, cfg)
	}

	return nil
}
