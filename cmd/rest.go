package cmd

import (
	"fmt"

	"github.com/ilmimris/poc-go-ulid/config"
	"github.com/ilmimris/poc-go-ulid/internal/adapter/inbound/rest"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var restCmd = &cobra.Command{
	Use:   "rest",
	Short: "Run rest server",
	Run: func(cmd *cobra.Command, args []string) {
		log.Info("Starting rest server")
		RestRun()
	},
}

func RestRun() {
	bst := NewBootstrap()
	bst.Initialized(cfgFile)
	bst.AddService(initServices()...)
	bst.RunServices()
}

func initServices() (services OptionsService) {
	// ...
	services = initRootServices()
	services.Add(
		NewServiceRest(rest.Options{
			Port:         fmt.Sprintf("%d", (config.GetConfig().Rest.Port)),
			BodyLimit:    config.GetConfig().Rest.BodyLimit,
			ReadTimeout:  config.GetConfig().Rest.ReadTimeOut,
			WriteTimeout: config.GetConfig().Rest.WriteTimeOut,
		}),
	)

	return
}
