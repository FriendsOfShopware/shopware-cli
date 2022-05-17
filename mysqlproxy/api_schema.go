package mysqlproxy

import (
	"context"
	"github.com/dolthub/go-mysql-server/sql"
	adminSdk "github.com/friendsofshopware/go-shopware-admin-api-sdk"
)

type ApiSchema struct {
	entities        map[string]Entity
	mappingEntities map[string]bool
}

func NewApiSchema(ctx context.Context, client *adminSdk.Client) (*ApiSchema, error) {

	req, err := client.NewRequest(adminSdk.NewApiContext(ctx), "GET", "/api/_info/entity-schema.json", nil)

	if err != nil {
		return nil, err
	}

	var schema map[string]Entity
	_, err = client.Do(ctx, req, &schema)
	if err != nil {
		return nil, err
	}

	mappingEntities := map[string]bool{}

	for _, entity := range schema {
		for _, table := range entity.MappingTables() {
			mappingEntities[table] = true
		}
	}
	return &ApiSchema{
		entities:        schema,
		mappingEntities: mappingEntities,
	}, nil
}

func (apiSchema *ApiSchema) BuildTables(client *adminSdk.Client) map[string]sql.Table {
	tables := map[string]sql.Table{}
	for name, entity := range apiSchema.entities {
		_, ok := apiSchema.mappingEntities[name]
		tables[name] = &AdminTable{
			Client:      client,
			adminEntity: entity,
			columns:     entityColumns(entity),
			isMapping:   ok,
		}
	}
	return tables
}

func entityColumns(entity Entity) []*sql.Column {
	var columns []*sql.Column

	for name, property := range entity.Properties {
		if property.Type == "association" {
			continue
		}
		columns = append(columns, &sql.Column{
			Name:          name,
			Type:          property.GetType(),
			Default:       nil,
			AutoIncrement: false,
			Nullable:      !property.IsPrimary(),
			Source:        entity.Name,
			PrimaryKey:    property.IsPrimary(),
			Comment:       property.Comment(),
			Extra:         "",
		})
	}
	return columns
}
