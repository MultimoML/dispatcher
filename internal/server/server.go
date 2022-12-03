package server

import (
	"context"
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
	if _, err := fmt.Fprint(w, "You should call the /products endpoint"); err != nil {
		log.Println(err)
	}
}

func Products(w http.ResponseWriter, _ *http.Request, _ httprouter.Params) {
	if _, err := fmt.Fprint(w, database.Products(context.TODO())); err != nil {
		log.Println(err)
	}
}
