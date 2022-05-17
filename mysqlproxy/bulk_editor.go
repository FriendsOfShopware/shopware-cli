package mysqlproxy

import (
	"encoding/json"
	"fmt"
	"github.com/dolthub/go-mysql-server/sql"
	"github.com/dolthub/vitess/go/vt/log"
	adminSdk "github.com/friendsofshopware/go-shopware-admin-api-sdk"
	"strings"
)

type editor interface {
	sql.RowDeleter
	sql.RowUpdater
}

type bulkEditor struct {
	table      *AdminTable
	deletes    []map[string]interface{}
	upserts    []map[string]interface{}
	memDeleter sql.RowDeleter
	memUpdater sql.RowUpdater
}

type bulkOp struct {
	Entity  string                   `json:"entity"`
	Action  string                   `json:"action"`
	Payload []map[string]interface{} `json:"payload"`
}

func (b *bulkEditor) StatementBegin(ctx *sql.Context) {
	b.memUpdater = b.table.memory.Updater(ctx)
	b.memDeleter = b.table.memory.Deleter(ctx)
	b.memUpdater.StatementBegin(ctx)
	b.memDeleter.StatementBegin(ctx)
}

func (b *bulkEditor) DiscardChanges(ctx *sql.Context, errorEncountered error) error {
	updateERR := b.memUpdater.DiscardChanges(ctx, errorEncountered)
	deleteERR := b.memDeleter.DiscardChanges(ctx, errorEncountered)
	if updateERR != nil || deleteERR != nil {
		return fmt.Errorf("err while syncing mem tables: %q, %q", updateERR, deleteERR)
	}
	return nil
}

func (b *bulkEditor) StatementComplete(ctx *sql.Context) error {
	updateERR := b.memUpdater.StatementComplete(ctx)
	deleteERR := b.memDeleter.StatementComplete(ctx)
	if updateERR != nil || deleteERR != nil {
		return fmt.Errorf("err while syncing mem tables: %q, %q", updateERR, deleteERR)
	}
	return nil
}

func (b *bulkEditor) Delete(context *sql.Context, row sql.Row) error {
	partialEntity := map[string]interface{}{}
	for _, idx := range b.table.pkSchema.PkOrdinals {
		property := b.table.pkSchema.Schema[idx].Name
		value := fmt.Sprintf("%v", row[idx])
		partialEntity[property] = value
	}
	b.deletes = append(b.deletes, partialEntity)
	return b.memDeleter.Delete(context, row)
}

func (b *bulkEditor) Update(ctx *sql.Context, old sql.Row, new sql.Row) error {
	partialEntity := map[string]interface{}{}
	for _, idx := range b.table.pkSchema.PkOrdinals {
		property := b.table.pkSchema.Schema[idx].Name
		value := fmt.Sprintf("%v", old[idx])
		partialEntity[property] = value
	}

	for idx, column := range b.table.columns {
		if new[idx] == nil {
			continue
		}
		if column.Comment == "json" {
			var bytes []byte
			switch v := new[idx].(type) {
			case []byte:
				bytes = v
			case string:
				bytes = []byte(v)
			default:
				continue
			}
			var data json.RawMessage
			err := json.Unmarshal(bytes, &data)
			if err == nil {
				partialEntity[column.Name] = data
			}
			continue
		}
		partialEntity[column.Name] = new[idx]
	}
	b.upserts = append(b.upserts, partialEntity)
	return b.memUpdater.Update(ctx, old, new)
}

func (b *bulkEditor) Close(context *sql.Context) error {
	client := b.table.Client
	ctx := adminSdk.NewApiContext(context)

	bulk := map[string]bulkOp{}

	if len(b.deletes) > 0 {
		operation := bulkOp{
			Entity:  b.table.Name(),
			Action:  "delete",
			Payload: b.deletes,
		}
		bulk[fmt.Sprintf("delete-%s", strings.ReplaceAll(operation.Entity, "_", "-"))] = operation
	}
	if len(b.upserts) > 0 {
		operation := bulkOp{
			Entity:  b.table.Name(),
			Action:  "upsert",
			Payload: b.upserts,
		}
		bulk[fmt.Sprintf("write-%s", strings.ReplaceAll(operation.Entity, "_", "-"))] = operation
		log.Infof("%+v", bulk)
	}

	req, err := client.NewRequest(ctx, "POST", "/api/_action/sync", bulk)
	if err != nil {
		return err
	}

	if _, err := client.Do(context, req, nil); err != nil {
		return err
	}

	updateERR := b.memUpdater.Close(context)
	deleteERR := b.memDeleter.Close(context)
	if updateERR != nil || deleteERR != nil {
		return fmt.Errorf("err while syncing mem tables: %q, %q", updateERR, deleteERR)
	}
	return nil
}
