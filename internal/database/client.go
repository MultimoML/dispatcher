package database

import (
	"context"
	"fmt"
	"log"
	"sync"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/multimoml/dispatcher/internal/config"
)

var dbClient *mongo.Client

func Connect(ctx context.Context, config *config.Config) *mongo.Client {
	var once sync.Once

	log.Println("Connecting to MongoDB...")
	once.Do(func() {
		mongoUrl := fmt.Sprintf("mongodb+srv://%s:%s@%s/%s?tls=true&authSource=admin&replicaSet=prod",
			config.DBUsername, config.DBPassword, config.DBHost, config.DBName)

		client, err := mongo.NewClient(options.Client().ApplyURI(mongoUrl))
		if err != nil {
			log.Fatal(err)
		}

		if err = client.Connect(ctx); err != nil {
			log.Fatal(err)
		}

		if err = client.Ping(ctx, nil); err != nil {
			log.Fatal(err)
		}

		dbClient = client
		log.Println("Connected to MongoDB")
	})

	return dbClient
}
