package server

import (
	"context"
	"fmt"
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
	cfg := config.LoadConfig()

	// Connect to MongoDB
	dbClient = database.Connect(ctx, cfg)

	// Set up router
	router := gin.Default()

	// Endpoints
	router.GET("/products/live", Liveness)
	router.GET("/products/ready", Readiness)

	v1 := router.Group("/products/v1")
	{
		v1.GET("/all", Products)
		v1.GET("/:id", Product)
	}

	// Start HTTP server
	log.Fatal(router.Run(fmt.Sprintf(":%s", cfg.Port)))
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

func Products(c *gin.Context) {
	params := parseQuery(c)
	if params.Error != nil {
		c.IndentedJSON(http.StatusBadRequest, params.Error)
		return
	}

	products, err := database.Products(params)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	} else {
		c.IndentedJSON(http.StatusOK, products)
	}
}

func Product(c *gin.Context) {
	params := parseQuery(c)
	if params.Error != nil {
		c.IndentedJSON(http.StatusBadRequest, params.Error)
		return
	}

	product, err := database.Product(params)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	} else {
		c.IndentedJSON(http.StatusOK, product)
	}
}
