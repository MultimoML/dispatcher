package server

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/multimoml/dispatcher/internal/config"
	"github.com/multimoml/dispatcher/internal/database"
)

var dbClient *mongo.Client

func Run(ctx context.Context) {
	// Load environment variables
	config.Environment()

	// Connect to MongoDB
	dbClient = database.Connect(ctx)

	// Start HTTP server
	router := httprouter.New()

	// Endpoints
	router.GET("/products/live", Liveliness)
	router.GET("/products/ready", Readiness)
	router.GET("/products/v1/all", AllProducts)

	log.Fatal(http.ListenAndServe(":6001", router))
}

func Liveliness(w http.ResponseWriter, _ *http.Request, _ httprouter.Params) {
	w.WriteHeader(http.StatusOK)
	if _, err := fmt.Fprint(w, "I'm alive!\n"); err != nil {
		log.Println(err)
	}
}

func Readiness(w http.ResponseWriter, _ *http.Request, _ httprouter.Params) {
	if err := dbClient.Ping(context.TODO(), nil); err != nil {
		w.WriteHeader(http.StatusServiceUnavailable)
		if _, err = fmt.Fprint(w, "I'm NOT ready!\n"); err != nil {
			log.Println(err)
		}
	} else {
		w.WriteHeader(http.StatusOK)
		if _, err = fmt.Fprint(w, "I'm ready!\n"); err != nil {
			log.Println(err)
		}
	}
}

func AllProducts(w http.ResponseWriter, _ *http.Request, _ httprouter.Params) {
	products := database.Products(context.TODO())

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	err := json.NewEncoder(w).Encode(products)
	if err != nil {
		log.Fatal("Failed to encode products into JSON")
	}
}
