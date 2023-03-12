package pgsql

import (
	"fmt"
	"sync"
	"time"

	"github.com/ilmimris/poc-go-ulid/pkg/database"
	"github.com/jmoiron/sqlx"
	log "github.com/sirupsen/logrus"
)

var once sync.Once
var pgsql *sqlx.DB

type OptPgsql struct {
	Host     string
	Port     int
	User     string
	Password string
	Database string
}

func New(configs ...OptPgsql) *sqlx.DB {
	var cfg OptPgsql
	if len(configs) > 0 {
		cfg = configs[0]
	}

	if pgsql != nil {
		return pgsql
	}

	once.Do(func() {
		dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.Database)
		db := database.DB{
			DBString:        dsn,
			MaxConn:         10,
			MaxIdleConn:     5,
			RetryInterval:   5,
			ConnMaxLifetime: time.Duration(60) * time.Minute,
		}
		err := db.ConnectAndMonitor(database.DriverPostgres)
		if err != nil {
			panic(err)
		}
		log.Info("Postgres connected")
		pgsql = db.DBConnection
	})

	return pgsql
}

func Close() {
	err := pgsql.Close()
	if err != nil {
		log.Error(err)
	}
	log.Info("Postgres closed")
}
