package cmd

import (
	"github.com/ilmimris/poc-go-ulid/config"
	"github.com/ilmimris/poc-go-ulid/internal/adapter/outbound/datastore/maria"
	"github.com/ilmimris/poc-go-ulid/internal/adapter/outbound/datastore/mongo"
	"github.com/ilmimris/poc-go-ulid/internal/adapter/outbound/datastore/pgsql"
	"github.com/ilmimris/poc-go-ulid/shared/constant"
	log "github.com/sirupsen/logrus"

	"github.com/spf13/cobra"
)

var cfgFile string

var rootCmd = &cobra.Command{
	Use: "myapp",
	Run: func(cmd *cobra.Command, args []string) {
		log.Println("use -h to show available commands")
	},
}

func Run() {
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "config.json", "config file (default is config.json)")
	rootCmd.AddCommand(restCmd)
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}

func initRootServices() (services OptionsService) {
	services.Add(
		NewServiceMariaDB(maria.OptMaria{
			Dsn: config.GetConfig().Database.Maria.Dsn,
		}),
		NewServicePgSQL(pgsql.OptPgsql{
			Host:     config.GetConfig().Database.Postgres.Host,
			Port:     config.GetConfig().Database.Postgres.Port,
			User:     config.GetConfig().Database.Postgres.Username,
			Password: config.GetConfig().Database.Postgres.Password,
			Database: config.GetConfig().Database.Postgres.Database,
		}),
		NewServiceMongoDB(mongo.OptMongo{
			URI:               config.GetConfig().Database.Mongo.URI,
			DB:                config.GetConfig().Database.Mongo.Db,
			AppName:           constant.AppName,
			ConnectionTimeOut: 10,
			PingTimeOut:       10,
		}),
	)
	return
}
