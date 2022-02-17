package project

import (
	"context"
	"shopware-cli/shop"
)

type ThemeSync struct{}

func (s ThemeSync) Push(ctx context.Context, client *shop.Client, config *shop.Config, operation *ConfigSyncOperation) error {
	if len(config.Sync.Theme) == 0 {
		return nil
	}

	criteria := shop.Criteria{}
	criteria.Includes = map[string][]string{"theme": {"id", "name"}}
	themes, err := client.SearchAll(ctx, "theme", criteria)

	if err != nil {
		return err
	}

	for _, theme := range themes.Data {
		t := theme.(map[string]interface{})

		remoteConfigs, err := client.GetThemeConfiguration(ctx, t["id"].(string))

		if err != nil {
			return err
		}

		for _, localThemeConfig := range config.Sync.Theme {
			if localThemeConfig.Name == t["name"] {
				op := ThemeSyncOperation{
					Id:       t["id"].(string),
					Name:     t["name"].(string),
					Settings: map[string]shop.ThemeConfigValue{},
				}

				for remoteFieldName, remoteFieldValue := range *remoteConfigs.CurrentFields {
					for localFieldName, localFieldValue := range localThemeConfig.Settings {
						if remoteFieldName == localFieldName {
							if remoteFieldValue != localFieldValue {
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

func (s ThemeSync) Pull(ctx context.Context, client *shop.Client, config *shop.Config) error {
	config.Sync.Theme = make([]shop.ThemeConfig, 0)

	criteria := shop.Criteria{}
	criteria.Includes = map[string][]string{"theme": {"id", "name"}}
	themes, err := client.SearchAll(ctx, "theme", criteria)

	if err != nil {
		return err
	}

	for _, theme := range themes.Data {
		t := theme.(map[string]interface{})

		cfg := shop.ThemeConfig{
			Name:     t["name"].(string),
			Settings: map[string]shop.ThemeConfigValue{},
		}

		themeConfig, err := client.GetThemeConfiguration(ctx, t["id"].(string))

		if err != nil {
			return err
		}

		cfg.Settings = *themeConfig.CurrentFields
		config.Sync.Theme = append(config.Sync.Theme, cfg)
	}

	return nil
}
