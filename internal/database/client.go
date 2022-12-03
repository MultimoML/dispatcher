package database

import (
	"context"
	"fmt"
	"log"
	"os"
	"sync"

	"go.mongodb.org/mongo-driver/bson"

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

func Products(ctx context.Context) []model.Product {
	productCollection := dbClient.Database(os.Getenv("DATABASE")).Collection("spar")

	cursor, err := productCollection.Find(ctx, bson.D{})
	if err != nil {
		log.Fatal("Failed to find products in database: ", err)
	}

	var products []model.Product
	if err = cursor.All(ctx, &products); err != nil {
		log.Fatal("Failed to save products into struct: ", err)
	}

	return products
}
