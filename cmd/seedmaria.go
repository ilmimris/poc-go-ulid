package cmd

import (
	"math/rand"
	"strings"
	"time"

	faker "github.com/go-faker/faker/v4"
	"github.com/ilmimris/poc-go-ulid/config"
	ds "github.com/ilmimris/poc-go-ulid/internal/adapter/outbound/datastore/maria"
	"github.com/ilmimris/poc-go-ulid/internal/core/domain"
	"github.com/ilmimris/poc-go-ulid/pkg/log"
	"github.com/oklog/ulid"
	"github.com/spf13/cobra"
)

var seedmaria = &cobra.Command{
	Use:   "seedmaria",
	Short: "golang seed into mariadb",
	Long:  `golang seed into mariadb`,
	Run: func(cmd *cobra.Command, args []string) {
		log.Info("seedmaria called")
		SeedMaria()
	},
}

func SeedMaria() {
	// Read config file
	log.Info("Reading config file")
	config.LoadConfig("config.json")

	// Init Entropy
	log.Info("Init Entropy")
	t := time.Unix(1000000, 0)
	entropy := ulid.Monotonic(rand.New(rand.NewSource(t.UnixNano())), 0)

	// Connect to mariadb
	log.Info("Connect to mariadb")
	db := ds.New(ds.OptMaria{
		Dsn: config.GetConfig().Database.Maria.Dsn,
	})

	// Seed data
	log.Info("Seed data")
	employees := make([]*domain.Employee, nEmployees)

	log.Info("Generating data")
	for i := 0; i < nEmployees; i++ {
		id := ulid.MustNew(ulid.Timestamp(time.Now()), entropy).String()

		employee := &domain.Employee{
			ID:         id,
			Name:       faker.FirstName() + " " + faker.LastName(),
			Email:      faker.Email(),
			Age:        rand.Intn(60-18+1) + 18,
			Phone:      faker.Phonenumber(),
			Department: faker.Word(),
			CreatedAt:  time.Now(),
			UpdatedAt:  time.Now(),
			IsDeleted:  false,
		}

		employees[i] = employee
	}

	// Migrate
	var schema = `
	CREATE TABLE IF NOT EXISTS employees (
		id varchar(26) PRIMARY KEY,
		name text,
		age integer,
		email text,
		phone text,
		department text,
		created_at timestamp,
		updated_at timestamp,
		deleted_at timestamp,
		is_deleted boolean,
		INDEX idx_employees_id (id)
	);
`

	log.Info("Migrate Schema")
	db.MustExec("DROP TABLE IF EXISTS employees")
	db.MustExec(schema)

	// Insert data
	// set elapsed time
	tInsert := time.Now()
	log.Info("Insert data")
	tx := db.MustBegin()
	sqlStr := `INSERT INTO employees (id, name, age, email, phone, department, created_at, updated_at, deleted_at, is_deleted) VALUES `
	qInsert := sqlStr
	vals := make([]interface{}, 0, nEmployees)
	for _, employee := range employees {
		sqlStr += "(?, ?, ?, ?, ?, ?, ?, ?, ?, ?),"
		vals = append(vals, employee.ID, employee.Name, employee.Age, employee.Email, employee.Phone, employee.Department, employee.CreatedAt, employee.UpdatedAt, employee.DeletedAt, employee.IsDeleted)

		// commit every 1000 rows
		if len(vals) >= 1000 {
			sqlStr = strings.TrimSuffix(sqlStr, ",")
			stmt, err := tx.Prepare(sqlStr)
			if err != nil {
				log.Error(err)
			}

			_, err = stmt.Exec(vals...)
			if err != nil {
				log.Error(err)
			}

			sqlStr = qInsert
			vals = vals[:0]
		}
	}

	log.Info("Commit data")
	err := tx.Commit()
	if err != nil {
		log.Error(err)
	}

	duration := time.Since(tInsert).Milliseconds()
	log.Infof("Insert %d data in %d ms", nEmployees, duration)

	// release memory mariadb
	log.Info("Release memory mariadb")
	db.MustExec("FLUSH TABLES")

	// Close connection
	log.Info("Close connection")
	db.Close()

	log.Info("Done")
}
