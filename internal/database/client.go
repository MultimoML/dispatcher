package database

import (
	"context"
	"fmt"
	"log"
	"os"
	"sync"

	"github.com/multimoml/dispatcher/internal/model"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var dbClient *mongo.Client

func Connect(ctx context.Context) *mongo.Client {
	var once sync.Once

	once.Do(func() {
		user := os.Getenv("M_USERNAME")
		pass := os.Getenv("M_PASSWORD")
		server := os.Getenv("M_SERVER")
		mongoUrl := fmt.Sprintf("mongodb://%s:%s@%s/", user, pass, server)

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

func Collection(_ context.Context, client *mongo.Client, name string) *mongo.Collection {
	return client.Database(os.Getenv("DATABASE")).Collection(name)
}

func Products(ctx context.Context) []model.Product {
	productCollection := Collection(ctx, dbClient, "extractor-timer")

	var products []model.Product
	cursor, err := productCollection.Find(ctx, nil)
	if err != nil {
		log.Fatal(err)
	}

	if err = cursor.All(ctx, &products); err != nil {
		log.Fatal(err)
	}

	return products
}
