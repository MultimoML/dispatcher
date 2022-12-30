package server

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"go.mongodb.org/mongo-driver/mongo"

	_ "github.com/multimoml/dispatcher/docs"
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

	// Redirect /products and /products/openapi to /products/openapi/index.html
	router.GET("/products", func(c *gin.Context) {
		c.Redirect(http.StatusMovedPermanently, "/products/openapi/index.html")
	})
	router.GET("/products/openapi", func(c *gin.Context) {
		c.Redirect(http.StatusMovedPermanently, "/products/openapi/index.html")
	})

	router.GET("/products/openapi/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	v1 := router.Group("/products/v1")
	{
		v1.GET("/all", Products)
		v1.GET("/:id", Product)
	}

	// Start HTTP server
	log.Fatal(router.Run(fmt.Sprintf(":%s", cfg.Port)))
}

// Liveness is a simple endpoint to check if the server is alive
// @Summary Get liveness status of the microservice
// @Description Get liveness status of the microservice
// @Tags Kubernetes
// @Success 200 {string} string
// @Router /live [get]
func Liveness(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, gin.H{"status": "alive"})
}

// Readiness is a simple endpoint to check if the server is ready
// @Summary Get readiness status of the microservice
// @Description Get readiness status of the microservice
// @Tags Kubernetes
// @Success 200 {string} string
// @Failure 503 {string} string
// @Router /ready [get]
func Readiness(c *gin.Context) {
	if err := dbClient.Ping(context.TODO(), nil); err != nil {
		c.IndentedJSON(http.StatusServiceUnavailable, gin.H{"status": "not ready"})
	} else {
		c.IndentedJSON(http.StatusOK, gin.H{"status": "ready"})
	}
}

// Products returns a list of products from the database
// @Summary Get a list of products
// @Description Get a list of products
// @Tags Products
// @Produce json
// @Success 200 {array} object
// @Failure 400 {string} string
// @Failure 500 {string} string
// @Router /v1/all [get]
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

// Product returns a product from the database
// @Summary Get a product
// @Description Get a product
// @Tags Products
// @Produce json
// @Param id path string true "Product ID"
// @Success 200 {object} object
// @Failure 400 {string} string
// @Failure 500 {string} string
// @Router /v1/{id} [get]
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
