package mysqlproxy

import (
	"context"
	"github.com/dolthub/go-mysql-server/sql"
	adminSdk "github.com/friendsofshopware/go-shopware-admin-api-sdk"
	log "github.com/sirupsen/logrus"
)

func NewAdminDatabase(client *adminSdk.Client) *AdminDatabase {
	db := &AdminDatabase{
		Client: client,
	}

	db.tables = map[string]sql.Table{}

	return db
}

type AdminDatabase struct {
	Client *adminSdk.Client
	tables map[string]sql.Table
}

func (db *AdminDatabase) Name() string {
	return "shopware"
}

func (db *AdminDatabase) GetTableInsensitive(ctx *sql.Context, tblName string) (sql.Table, bool, error) {
	if len(db.tables) == 0 {
		if err := db.loadTables(ctx); err != nil {
			return nil, false, err
		}
	}

	log.Infof("searching for table %s", tblName)

	table, ok := db.tables[tblName]

	if !ok {
		return nil, false, sql.ErrTableNotFound.New(tblName)
	}

	return table, true, nil
}

func (db *AdminDatabase) GetTableNames(ctx *sql.Context) ([]string, error) {
	if len(db.tables) == 0 {
		if err := db.loadTables(ctx); err != nil {
			return nil, err
		}
	}

	var tables []string

	for name, _ := range db.tables {
		tables = append(tables, name)
	}

	return tables, nil
}

func (db AdminDatabase) loadTables(ctx context.Context) error {
	req, err := db.Client.NewRequest(adminSdk.NewApiContext(ctx), "GET", "/api/_info/entity-schema.json", nil)

	if err != nil {
		return err
	}

	var schema map[string]Entity

	_, err = db.Client.Do(ctx, req, &schema)
	if err != nil {
		return err
	}

	for s, entity := range schema {
		db.tables[s] = NewAdminTable(db.Client, entity)
	}

	return nil
}
