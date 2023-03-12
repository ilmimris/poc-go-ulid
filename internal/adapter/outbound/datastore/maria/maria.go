package maria

import (
	"sync"
	"time"

	"github.com/ilmimris/poc-go-ulid/pkg/database"
	log "github.com/sirupsen/logrus"

	"github.com/jmoiron/sqlx"
)

var once sync.Once
var maria *sqlx.DB

type OptMaria struct {
	Dsn string
}

func New(configs ...OptMaria) *sqlx.DB {
	var cfg OptMaria
	if len(configs) > 0 {
		cfg = configs[0]
	}

	if maria != nil {
		return maria
	}

	once.Do(func() {
		db := database.DB{
			DBString:        cfg.Dsn,
			MaxConn:         10,
			MaxIdleConn:     5,
			RetryInterval:   5,
			ConnMaxLifetime: time.Duration(60) * time.Minute,
		}
		err := db.ConnectAndMonitor(database.DriverMySQL)
		if err != nil {
			panic(err)
		}
		log.Info("Maria connected")
		maria = db.DBConnection
	})

	return maria
}

func Close() {
	err := maria.Close()
	if err != nil {
		log.Error(err)
	}
	log.Info("Maria closed")
}
