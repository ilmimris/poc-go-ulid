package cmd

import (
	"fmt"
	"math/rand"
	"strings"
	"time"

	faker "github.com/go-faker/faker/v4"
	"github.com/ilmimris/poc-go-ulid/config"
	pgDataStore "github.com/ilmimris/poc-go-ulid/internal/adapter/outbound/datastore/pgsql"
	"github.com/ilmimris/poc-go-ulid/internal/core/domain"
	"github.com/ilmimris/poc-go-ulid/pkg/log"
	"github.com/oklog/ulid"
	"github.com/spf13/cobra"
)

const nEmployees = 500000 // 1 Million

var seedpg = &cobra.Command{
	Use:   "seedpg",
	Short: "golang seed data into postgresql",
	Long:  `golang seed data into postgresql`,
	Run: func(cmd *cobra.Command, args []string) {
		log.Info("seedpg called")
		SeedPg()
	},
}

func SeedPg() {
	// Read config file
	log.Info("Reading config file")
	config.LoadConfig("config.json")

	// Init Entropy
	log.Info("Init Entropy")
	t := time.Unix(10, 0)
	entropy := ulid.Monotonic(rand.New(rand.NewSource(t.UnixNano())), 0)

	// Connect to postgresql
	log.Info("Connect to postgresql")
	db := pgDataStore.New(pgDataStore.OptPgsql{
		Host:     config.GetConfig().Database.Postgres.Host,
		Port:     config.GetConfig().Database.Postgres.Port,
		User:     config.GetConfig().Database.Postgres.Username,
		Password: config.GetConfig().Database.Postgres.Password,
		Database: config.GetConfig().Database.Postgres.Database,
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
			DeletedAt:  time.Time{},
			IsDeleted:  false,
		}

		employees[i] = employee
	}

	// Migrate
	var schema = `
DROP TABLE IF EXISTS employees;

CREATE TABLE employees (
    id text PRIMARY KEY,
    name text,
    age integer,
	email text,
	phone text,
	department text,
	created_at timestamp,
	updated_at timestamp,
	deleted_at timestamp,
	is_deleted boolean
);

CREATE INDEX employees_id_idx ON employees (id);
`

	log.Info("Migrate Schema")
	db.MustExec(schema)

	// Insert data
	// set elapsed time
	tInsert := time.Now()
	log.Info("Insert data")
	tx := db.MustBegin()
	sqlStr := `INSERT INTO employees (id, name, age, email, phone, department, created_at, updated_at, deleted_at, is_deleted) VALUES `
	qInsert := sqlStr
	nbatch := 1000
	vals := make([]interface{}, 10*nbatch)
	idx := 0
	c := 0
	for _, employee := range employees {
		idx = c * 10
		if len(vals) < idx+10 {
			vals = append(vals, make([]interface{}, 10*nbatch)...)
		}

		vals[idx] = employee.ID
		vals[idx+1] = employee.Name
		vals[idx+2] = employee.Age
		vals[idx+3] = employee.Email
		vals[idx+4] = employee.Phone
		vals[idx+5] = employee.Department
		vals[idx+6] = employee.CreatedAt
		vals[idx+7] = employee.UpdatedAt
		vals[idx+8] = employee.DeletedAt
		vals[idx+9] = employee.IsDeleted
		c++

		sqlStr += fmt.Sprintf("($%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d),", idx+1, idx+2, idx+3, idx+4, idx+5, idx+6, idx+7, idx+8, idx+9, idx+10)

		// commit every 1000 rows
		if c >= nbatch {
			sqlStr = strings.TrimSuffix(sqlStr, ",")
			_, err := tx.Exec(sqlStr, vals...)
			if err != nil {
				log.Error(err)
			}

			sqlStr = qInsert
			vals = make([]interface{}, 10*nbatch)
			idx = 0
			c = 0
		}
	}

	log.Info("Commit data")
	tx.Commit()
	duration := time.Since(tInsert).Milliseconds()
	log.Infof("Insert %d data in %d ms", nEmployees, duration)

	// Close connection
	log.Info("Close connection")
	db.Close()

	log.Info("Done")
}
