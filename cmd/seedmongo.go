package cmd

import (
	"context"
	"math/rand"
	"time"

	faker "github.com/go-faker/faker/v4"
	"github.com/ilmimris/poc-go-ulid/config"
	mongoDataStore "github.com/ilmimris/poc-go-ulid/internal/adapter/outbound/datastore/mongo"
	"github.com/ilmimris/poc-go-ulid/internal/core/domain"
	"github.com/ilmimris/poc-go-ulid/pkg/log"
	"github.com/oklog/ulid"
	"github.com/spf13/cobra"
	"go.mongodb.org/mongo-driver/mongo"
)

var seedmongo = &cobra.Command{
	Use:   "seedmongo",
	Short: "golang seed data into mongodb",
	Long:  `golang seed data into mongodb`,
	Run: func(cmd *cobra.Command, args []string) {
		log.Info("seedmongo called")
		SeedMongo()
	},
}

func SeedMongo() {
	// Read config file
	log.Info("Reading config file")
	config.LoadConfig("config.json")

	// Init Entropy
	log.Info("Init Entropy")
	t := time.Unix(1000000, 0)
	entropy := ulid.Monotonic(rand.New(rand.NewSource(t.UnixNano())), 0)

	// Connect to mongodb
	log.Info("Connect to mongodb")
	db := mongoDataStore.New(mongoDataStore.OptMongo{
		URI:               config.GetConfig().Database.Mongo.URI,
		DB:                config.GetConfig().Database.Mongo.Db,
		AppName:           "poc",
		ConnectionTimeOut: int(config.GetConfig().Database.Mongo.ConnectionTimeOut),
		PingTimeOut:       int(config.GetConfig().Database.Mongo.PingTimeOut),
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

	log.Info("Inserting data")
	ctx := context.Background()

	// bulk write insert mongodb
	writeModels := make([]mongo.WriteModel, len(employees))
	for i, employee := range employees {
		writeModels[i] = mongo.NewInsertOneModel().SetDocument(employee)
	}

	tInsert := time.Now()
	// using transaction
	session, err := db.Client().StartSession()
	if err != nil {
		log.Fatal(err)
	}

	err = session.StartTransaction()
	if err != nil {
		log.Fatal(err)
	}

	// drop collection
	err = db.Collection("employees").Drop(ctx)
	if err != nil {
		log.Fatal(err)
	}

	res, err := db.Collection("employees").BulkWrite(ctx, writeModels)
	if err != nil {
		log.Fatal(err)
	}

	err = session.CommitTransaction(ctx)
	if err != nil {
		log.Fatal(err)
	}
	duration := time.Since(tInsert).Milliseconds()
	log.Infof("Insert %d data in %d ms", res.InsertedCount, duration)

	// close connection
	log.Info("Close connection")
	db.Client().Disconnect(ctx)

	log.Info("Done")
}
