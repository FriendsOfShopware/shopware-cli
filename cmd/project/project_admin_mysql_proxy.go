package project

import (
	sqle "github.com/dolthub/go-mysql-server"
	"github.com/dolthub/go-mysql-server/memory"
	"github.com/dolthub/go-mysql-server/server"
	"github.com/dolthub/go-mysql-server/sql"
	"github.com/dolthub/go-mysql-server/sql/information_schema"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"shopware-cli/mysqlproxy"
	"shopware-cli/shop"
	"time"
)

var projectAdminMysqlCmd = &cobra.Command{
	Use:   "admin-mysql-proxy",
	Short: "Spawns a MySQL server to access the entities",
	RunE: func(cobraCmd *cobra.Command, args []string) error {
		var cfg *shop.Config
		var err error

		if cfg, err = shop.ReadConfig(projectConfigPath); err != nil {
			return err
		}

		client, err := shop.NewShopClient(cobraCmd.Context(), cfg, nil)
		if err != nil {
			return err
		}

		shopwareDb := mysqlproxy.NewAdminDatabase(client)

		engine := sqle.NewDefault(
			sql.NewDatabaseProvider(
				shopwareDb,
				information_schema.NewInformationSchemaDatabase(),
			))

		config := server.Config{
			Protocol: "tcp",
			Address:  "localhost:3307",
		}

		s, err := server.NewDefaultServer(config, engine)
		if err != nil {
			panic(err)
		}

		log.Infof("Started MySQL Proxy at localhost:3307")

		s.Start()

		return nil
	},
}

func init() {
	projectRootCmd.AddCommand(projectAdminMysqlCmd)
}

func createTestDatabase() *memory.Database {
	const (
		dbName    = "mydb"
		tableName = "mytable"
	)

	db := memory.NewDatabase(dbName)
	table := memory.NewTable(tableName, sql.NewPrimaryKeySchema(sql.Schema{
		{Name: "name", Type: sql.Text, Nullable: false, Source: tableName},
		{Name: "email", Type: sql.Text, Nullable: false, Source: tableName},
		{Name: "phone_numbers", Type: sql.JSON, Nullable: false, Source: tableName},
		{Name: "created_at", Type: sql.Timestamp, Nullable: false, Source: tableName},
	}), db.GetForeignKeyCollection())

	creationTime := time.Unix(1524044473, 0).UTC()
	db.AddTable(tableName, table)
	ctx := sql.NewEmptyContext()
	table.Insert(ctx, sql.NewRow("John Doe", "john@doe.com", []string{"555-555-555"}, creationTime))
	table.Insert(ctx, sql.NewRow("John Doe", "johnalt@doe.com", []string{}, creationTime))
	table.Insert(ctx, sql.NewRow("Jane Doe", "jane@doe.com", []string{}, creationTime))
	table.Insert(ctx, sql.NewRow("Evil Bob", "evilbob@gmail.com", []string{"555-666-555", "666-666-666"}, creationTime))
	return db
}
