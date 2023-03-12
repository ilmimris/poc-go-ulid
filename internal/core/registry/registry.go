package registry

import (
	"github.com/jmoiron/sqlx"
	"go.mongodb.org/mongo-driver/mongo"

	restRegistry "github.com/ilmimris/poc-go-ulid/internal/adapter/inbound/rest/registry"
)

type registry struct {
	pg    *sqlx.DB
	maria *sqlx.DB
	mongo *mongo.Database
}

type Registry interface {
	NewRestRegistry() *restRegistry.ServiceRegistry
}

type Option func(r *registry)

func NewRegistry(opts ...Option) Registry {
	r := &registry{}

	for _, opt := range opts {
		opt(r)
	}

	return r
}

func (r *registry) NewRestRegistry() *restRegistry.ServiceRegistry {
	return restRegistry.NewServiceRegistry()
}

func WithPg(pg *sqlx.DB) Option {
	return func(r *registry) {
		r.pg = pg
	}
}

func WithMaria(maria *sqlx.DB) Option {
	return func(r *registry) {
		r.maria = maria
	}
}

func WithMongo(mongo *mongo.Database) Option {
	return func(r *registry) {
		r.mongo = mongo
	}
}
