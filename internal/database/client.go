package database

import (
	"context"
	"fmt"
	"log"
	"sync"

	"github.com/multimoml/dispatcher/internal/config"
	"go.mongodb.org/mongo-driver/bson"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/multimoml/dispatcher/internal/model"
)

var dbClient *mongo.Client

func Connect(ctx context.Context, config *config.Config) *mongo.Client {
	var once sync.Once

	once.Do(func() {
		username := config.DBUsername
		password := config.DBPassword
		host := config.DBHost
		mongoUrl := fmt.Sprintf("mongodb+srv://%s:%s@%s/", username, password, host)

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
	productCollection := dbClient.Database("products").Collection("spar")

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
