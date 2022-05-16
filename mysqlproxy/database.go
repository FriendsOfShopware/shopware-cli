package mysqlproxy

import (
	"github.com/dolthub/go-mysql-server/sql"
	adminSdk "github.com/friendsofshopware/go-shopware-admin-api-sdk"
)

type AdminDatabase struct {
	client adminSdk.Client
}

func (db *AdminDatabase) Name() string {
	return "shopware"
}

func (db *AdminDatabase) GetTableInsensitive(ctx *sql.Context, tblName string) (sql.Table, bool, error) {
	panic("implement me")
}

func (db *AdminDatabase) GetTableNames(ctx *sql.Context) ([]string, error) {
	panic("implement me")
}
