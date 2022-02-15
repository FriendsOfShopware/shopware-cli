package project

import (
	"context"
	"shopware-cli/shop"
)

func readSystemConfig(ctx context.Context, client *shop.Client, salesChannelId *string) (*shop.SearchResponse, error) {
	c := shop.Criteria{}
	c.Includes = map[string][]string{"system_config": {"id", "configurationKey", "configurationValue"}}

	if salesChannelId != nil {
		c.Filter = []shop.CriteriaFilter{
			{Type: shop.SearchFilterTypeEquals, Field: "salesChannelId", Value: salesChannelId},
		}
	}

	configs, err := client.SearchAll(ctx, "system_config", c)

	if err != nil {
		return nil, err
	}

	return configs, nil
}

type ConfigSyncApplyer interface {
	Push(ctx context.Context, client *shop.Client, config *shop.Config, operation *ConfigSyncOperation) error
	Pull(ctx context.Context, client *shop.Client, config *shop.Config) error
}

func NewSyncApplyers() []ConfigSyncApplyer {
	return []ConfigSyncApplyer{SystemConfigSync{}}
}

type ConfigSyncOperation struct {
	Upsert map[string][]map[string]interface{}
	Delete map[string][]map[string]interface{}
}

func (o ConfigSyncOperation) HasChanges() bool {
	for _, i := range o.Upsert {
		if len(i) > 0 {
			return true
		}
	}

	for _, i := range o.Delete {
		if len(i) > 0 {
			return true
		}
	}

	return false
}
