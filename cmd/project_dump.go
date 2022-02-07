package cmd

import (
	"database/sql"
	"github.com/doutorfinancas/go-mad/core"
	"github.com/doutorfinancas/go-mad/database"
	"github.com/doutorfinancas/go-mad/generator"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
	"io"
	"os"
)

var projectDatabaseDumpCmd = &cobra.Command{
	Use:   "dump [database]",
	Short: "Dumps the Shopware database",
	Args:  cobra.ExactArgs(1),
	RunE: func(cobraCmd *cobra.Command, args []string) error {
		host, _ := cobraCmd.Flags().GetString("host")
		port, _ := cobraCmd.Flags().GetString("port")
		username, _ := cobraCmd.Flags().GetString("username")
		password, _ := cobraCmd.Flags().GetString("password")
		output, _ := cobraCmd.Flags().GetString("output")
		clean, _ := cobraCmd.Flags().GetBool("clean")

		cfg := database.NewConfig(username, password, host, port, args[0])

		db, err := sql.Open("mysql", cfg.ConnectionString())

		if err != nil {
			return err
		}

		service := generator.NewService()
		var opt []database.Option
		opt = append(opt, database.OptionValue("hex-encode", "1"))
		opt = append(opt, database.OptionValue("set-charset", "utf8mb4"))

		logger, _ := zap.NewProduction()
		dumper, err := database.NewMySQLDumper(db, logger, service, opt...)

		if err != nil {
			return err
		}

		pConf := core.Rules{Ignore: []string{}, NoData: []string{}, Where: map[string]string{}, Rewrite: map[string]core.Rewrite{}}

		if clean {
			pConf.NoData = append(pConf.NoData, "cart", "customer_recovery", "dead_message", "enqueue", "elasticsearch_index_task", "log_entry", "message_queue_stats", "notification", "payment_token", "refresh_token", "version", "version_commit", "version_commit_data", "webhook_event_log")
		}

		dumper.SetSelectMap(pConf.RewriteToMap())
		dumper.SetWhereMap(pConf.Where)
		if dErr := dumper.SetFilterMap(pConf.NoData, pConf.Ignore); dErr != nil {
			return dErr
		}

		var w io.Writer
		if w, err = os.Create(output); err != nil {
			return err
		}

		if err = dumper.Dump(w); err != nil {
			return err
		}

		return nil
	},
}

func init() {
	projectRootCmd.AddCommand(projectDatabaseDumpCmd)
	projectDatabaseDumpCmd.Flags().String("host", "127.0.0.1", "hostname")
	projectDatabaseDumpCmd.Flags().String("username", "root", "mysql user")
	projectDatabaseDumpCmd.Flags().String("password", "root", "mysql password")
	projectDatabaseDumpCmd.Flags().String("port", "3306", "mysql port")
	projectDatabaseDumpCmd.Flags().String("output", "dump.sql", "file")
	projectDatabaseDumpCmd.Flags().Bool("clean", false, "Ignores cart, enqueue, message_queue_stats")
}
