package mysqlproxy

import (
	json2 "encoding/json"
	"fmt"
	"github.com/dolthub/go-mysql-server/memory"
	"github.com/dolthub/go-mysql-server/sql"
	adminSdk "github.com/friendsofshopware/go-shopware-admin-api-sdk"
	"strings"
)

type AdminTable struct {
	Client      *adminSdk.Client
	adminEntity Entity
	memory      *memory.Table
	columns     []*sql.Column
	pkSchema    sql.PrimaryKeySchema
	isMapping   bool
}

func (at *AdminTable) Deleter(context *sql.Context) sql.RowDeleter {
	return &bulkEditor{table: at}
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
	at.pkSchema = sql.NewPrimaryKeySchema(at.Schema())
	at.memory = memory.NewTable(at.Name(), at.pkSchema, nil)

	ctx := adminSdk.NewApiContext(context)
	searchUrl := fmt.Sprintf("/api/search/%s", strings.ReplaceAll(at.Name(), "_", "-"))

	if at.isMapping {
		searchUrl = strings.ReplaceAll(searchUrl, "/search/", "/search-ids/")
	}

	req, err := at.Client.NewRequest(ctx, "POST", searchUrl, nil)

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
