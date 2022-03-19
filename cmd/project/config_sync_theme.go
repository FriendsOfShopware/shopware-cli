package project

import (
	"encoding/json"
	adminSdk "github.com/friendsofshopware/go-shopware-admin-api-sdk"
	"shopware-cli/shop"
)

type ThemeSync struct{}

func (s ThemeSync) Push(ctx adminSdk.ApiContext, client *adminSdk.Client, config *shop.Config, operation *ConfigSyncOperation) error {
	if len(config.Sync.Theme) == 0 {
		return nil
	}

	criteria := adminSdk.Criteria{}
	criteria.Includes = map[string][]string{"theme": {"id", "name"}}
	themes, _, err := client.Repository.Theme.SearchAll(ctx, criteria)

	if err != nil {
		return err
	}

	for _, t := range themes.Data {
		remoteConfigs, _, err := client.ThemeManager.GetConfiguration(ctx, t.Id)

		if err != nil {
			return err
		}

		for _, localThemeConfig := range config.Sync.Theme {
			if localThemeConfig.Name == t.Name {
				op := ThemeSyncOperation{
					Id:       t.Id,
					Name:     t.Name,
					Settings: map[string]adminSdk.ThemeConfigValue{},
				}

				for remoteFieldName, remoteFieldValue := range *remoteConfigs.CurrentFields {
					for localFieldName, localFieldValue := range localThemeConfig.Settings {
						if remoteFieldName == localFieldName {
							localJson, _ := json.Marshal(localFieldValue)
							remoteJson, _ := json.Marshal(remoteFieldValue)

							if string(localJson) != string(remoteJson) {
								op.Settings[remoteFieldName] = localFieldValue
							}
						}
					}
				}

				operation.ThemeSettings = append(operation.ThemeSettings, op)
			}
		}
	}

	return nil
}

func (s ThemeSync) Pull(ctx adminSdk.ApiContext, client *adminSdk.Client, config *shop.Config) error {
	config.Sync.Theme = make([]shop.ThemeConfig, 0)

	criteria := adminSdk.Criteria{}
	criteria.Includes = map[string][]string{"theme": {"id", "name"}}
	themes, _, err := client.Repository.Theme.SearchAll(ctx, criteria)

	if err != nil {
		return err
	}

	for _, t := range themes.Data {
		cfg := shop.ThemeConfig{
			Name:     t.Name,
			Settings: map[string]adminSdk.ThemeConfigValue{},
		}

		themeConfig, _, err := client.ThemeManager.GetConfiguration(ctx, t.Id)

		if err != nil {
			return err
		}

		cfg.Settings = *themeConfig.CurrentFields
		config.Sync.Theme = append(config.Sync.Theme, cfg)
	}

	return nil
}
