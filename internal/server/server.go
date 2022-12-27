package server

import (
	"context"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/multimoml/dispatcher/internal/config"
	"github.com/multimoml/dispatcher/internal/database"
)

var dbClient *mongo.Client

func Run(ctx context.Context) {
	// Load environment variables
	dbConfig := config.LoadConfig()

	// Connect to MongoDB
	dbClient = database.Connect(ctx, dbConfig)

	// Start HTTP server
	router := gin.Default()

	// Endpoints
	router.GET("/products/live", Liveness)
	router.GET("/products/ready", Readiness)
	router.GET("/products/v1/all", AllProducts)

	log.Fatal(router.Run("localhost:6001"))
}

func Liveness(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, gin.H{"status": "alive"})
}

func Readiness(c *gin.Context) {
	if err := dbClient.Ping(context.TODO(), nil); err != nil {
		c.IndentedJSON(http.StatusServiceUnavailable, gin.H{"status": "not ready"})
	} else {
		c.IndentedJSON(http.StatusOK, gin.H{"status": "ready"})
	}
}

func AllProducts(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, database.Products(context.TODO()))
}
