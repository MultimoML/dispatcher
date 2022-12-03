package server

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"

	"github.com/multimoml/dispatcher/internal/config"
	"github.com/multimoml/dispatcher/internal/database"
)

func Run(ctx context.Context) {
	// Load environment variables
	config.Environment()

	// Connect to MongoDB
	database.Connect(ctx)

	// Start HTTP server
	router := httprouter.New()

	// Endpoints
	router.GET("/", Index)
	router.GET("/products", Products)

	log.Fatal(http.ListenAndServe(":6001", router))
}

func Index(w http.ResponseWriter, _ *http.Request, _ httprouter.Params) {
	if _, err := fmt.Fprint(w, "You should call the /products endpoint\n"); err != nil {
		log.Println(err)
	}
}

func Products(w http.ResponseWriter, _ *http.Request, _ httprouter.Params) {
	products := database.Products(context.TODO())

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	err := json.NewEncoder(w).Encode(products)
	if err != nil {
		log.Fatal("Failed to encode products into JSON")
	}
}
