package mysqlproxy

import (
	json2 "encoding/json"
	"github.com/dolthub/go-mysql-server/memory"
	"github.com/dolthub/go-mysql-server/sql"
	adminSdk "github.com/friendsofshopware/go-shopware-admin-api-sdk"
	"strings"
)

func NewAdminTable(client *adminSdk.Client, entity Entity) *AdminTable {
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

	t := &AdminTable{
		Client:      client,
		adminEntity: entity,
		columns:     columns,
	}

	return t
}

type AdminTable struct {
	Client      *adminSdk.Client
	adminEntity Entity
	memory      *memory.Table
	columns     []*sql.Column
}

func (at *AdminTable) Name() string {
	return at.adminEntity.Name
}

func (at *AdminTable) String() string {
	return at.adminEntity.Name
}

func (at *AdminTable) Schema() sql.Schema {
	return at.columns
}

func (at *AdminTable) Partitions(context *sql.Context) (sql.PartitionIter, error) {
	at.memory = memory.NewTable(at.Name(), sql.NewPrimaryKeySchema(at.Schema()), nil)

	ctx := adminSdk.NewApiContext(context)
	req, err := at.Client.NewRequest(ctx, "POST", "/api/search/"+strings.ReplaceAll(at.Name(), "_", "-"), nil)

	if err != nil {
		return nil, err
	}

	var resp searchResp
	if _, err := at.Client.Do(context, req, &resp); err != nil {
		return nil, err
	}

	for _, row := range resp.Data {
		row := row.(map[string]interface{})
		insertRow := make([]interface{}, len(at.Schema()))

		for k, v := range row {
			i := at.Schema().IndexOf(k, at.Name())
			if i == -1 {
				continue
			}

			if at.Schema()[i].Comment == "json" {
				json, _ := json2.Marshal(v)
				insertRow[i] = json
			} else {
				insertRow[i] = v
			}
		}

		if err := at.memory.Insert(context, sql.NewRow(insertRow...)); err != nil {
			return nil, err
		}
	}

	return at.memory.Partitions(context)
}

func (at *AdminTable) PartitionRows(context *sql.Context, partition sql.Partition) (sql.RowIter, error) {
	return at.memory.PartitionRows(context, partition)
}

type searchResp struct {
	Total int           `json:"total"`
	Data  []interface{} `json:"data"`
}
