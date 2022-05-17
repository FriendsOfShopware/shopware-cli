package mysqlproxy

import (
	"fmt"
	"github.com/dolthub/go-mysql-server/sql"
	adminSdk "github.com/friendsofshopware/go-shopware-admin-api-sdk"
	"strings"
)

type bulkEditor struct {
	table      *AdminTable
	deletes    []map[string]string
	memDeleter sql.RowDeleter
}

type bulkOp struct {
	Entity  string              `json:"entity"`
	Action  string              `json:"action"`
	Payload []map[string]string `json:"payload"`
}

func (b *bulkEditor) StatementBegin(ctx *sql.Context) {
	b.memDeleter = b.table.memory.Deleter(ctx)
}

func (b *bulkEditor) DiscardChanges(ctx *sql.Context, errorEncountered error) error {
	return b.memDeleter.DiscardChanges(ctx, errorEncountered)
}

func (b *bulkEditor) StatementComplete(ctx *sql.Context) error {
	return b.memDeleter.StatementComplete(ctx)
}

func (b *bulkEditor) Delete(context *sql.Context, row sql.Row) error {
	partialEntity := map[string]string{}
	for _, idx := range b.table.pkSchema.PkOrdinals {
		property := b.table.pkSchema.Schema[idx].Name
		value := fmt.Sprintf("%v", row[idx])
		partialEntity[property] = value
	}
	b.deletes = append(b.deletes, partialEntity)
	return b.memDeleter.Delete(context, row)
}

func (b *bulkEditor) Close(context *sql.Context) error {
	client := b.table.Client
	ctx := adminSdk.NewApiContext(context)

	operation := bulkOp{
		Entity:  b.table.Name(),
		Action:  "delete",
		Payload: b.deletes,
	}

	bulk := map[string]bulkOp{
		fmt.Sprintf("delete-%s", strings.ReplaceAll(operation.Entity, "_", "-")): operation,
	}

	req, err := client.NewRequest(ctx, "POST", "/api/_action/sync", bulk)
	if err != nil {
		return err
	}

	if _, err := client.Do(context, req, nil); err != nil {
		return err
	}

	return b.memDeleter.Close(context)
}
