package cmd

import (
	"math/rand"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/ilmimris/poc-go-ulid/config"
	"github.com/ilmimris/poc-go-ulid/internal/adapter/inbound/interfaces"
	"github.com/ilmimris/poc-go-ulid/internal/adapter/inbound/rest"
	restRegistry "github.com/ilmimris/poc-go-ulid/internal/adapter/inbound/rest/registry"
	mariaDataStore "github.com/ilmimris/poc-go-ulid/internal/adapter/outbound/datastore/maria"
	mongoDataStore "github.com/ilmimris/poc-go-ulid/internal/adapter/outbound/datastore/mongo"
	pgDataStore "github.com/ilmimris/poc-go-ulid/internal/adapter/outbound/datastore/pgsql"
	"github.com/ilmimris/poc-go-ulid/internal/adapter/outbound/logger"
	"github.com/ilmimris/poc-go-ulid/internal/adapter/outbound/tracing"
	"github.com/ilmimris/poc-go-ulid/internal/core/registry"
	"github.com/jmoiron/sqlx"
	"go.mongodb.org/mongo-driver/mongo"
)

type (
	OptionService func(b *bootstrap)

	OptionsService []OptionService

	bootstrap struct {
		pg          *sqlx.DB
		maria       *sqlx.DB
		mongodb     *mongo.Database
		optServices []OptionService
		rest        interfaces.Rest
	}

	Bootstrap interface {
		Initialized(cfgFile string)
		AddService(service ...OptionService)
		GetRegistryRest() *restRegistry.ServiceRegistry
		RunServices()
		Shutdown()
	}
)

func (o *OptionsService) Add(option ...OptionService) {
	*o = append(*o, option...)
}

func (b *bootstrap) Initialized(cfgFile string) {
	config.LoadConfig(cfgFile)

	logger.InitLogger(config.GetConfig().Logger.Level)

	// initialize random seed
	rand.Seed(time.Now().UnixNano())

	tracing.InitTracing(config.GetConfig().TracerConfig)
}

func (b *bootstrap) AddService(services ...OptionService) {
	b.optServices = append(b.optServices, services...)
}

func (b *bootstrap) GetRegistryRest() *restRegistry.ServiceRegistry {
	return registry.NewRegistry(
		registry.WithPg(b.pg),
		registry.WithMaria(b.maria),
		registry.WithMongo(b.mongodb),
	).NewRestRegistry()
}

func (b *bootstrap) RunServices() {
	for _, value := range b.optServices {
		value(b)
	}

	// add graceful shutdown when interrupt signal detected
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGTERM, syscall.SIGINT)
	go func() {
		<-c
		b.Shutdown()
	}()
}

func (b *bootstrap) Shutdown() {

}

func NewBootstrap() Bootstrap {
	return &bootstrap{}
}

func NewServicePgSQL(param pgDataStore.OptPgsql) OptionService {
	return func(b *bootstrap) {
		b.pg = pgDataStore.New(param)
	}
}

func NewServiceMariaDB(param mariaDataStore.OptMaria) OptionService {
	return func(b *bootstrap) {
		b.maria = mariaDataStore.New(param)
	}
}

func NewServiceMongoDB(param mongoDataStore.OptMongo) OptionService {
	return func(b *bootstrap) {
		b.mongodb = mongoDataStore.New(param)
	}
}

func NewServiceRest(param rest.Options) OptionService {
	return func(b *bootstrap) {
		b.rest = rest.NewRest(&param, b.GetRegistryRest())
		b.rest.Serve()
	}
}
