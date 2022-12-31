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
var sConfig *config.Config

func Run(ctx context.Context) {
	// Load environment variables
	sConfig = config.LoadConfig()

	// Connect to MongoDB
	dbClient = database.Connect(ctx, sConfig)

	// Set up router
	router := gin.Default()

	// Endpoints
	p := router.Group("/products")
	{
		p.GET("/live", Liveness)
		p.GET("/ready", Readiness)

		p.GET("/openapi", func(c *gin.Context) {
			c.Redirect(http.StatusMovedPermanently, "/products/openapi/index.html")
		})

		p.GET("/", func(c *gin.Context) {
			c.Redirect(http.StatusMovedPermanently, "/products/openapi/index.html")
		})

		p.GET("/openapi/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

		v1 := p.Group("/v1")
		{
			v1.GET("/all", Products)
			v1.GET("/:id", Product)
		}
	}

	// Start HTTP server
	log.Fatal(router.Run(fmt.Sprintf(":%s", sConfig.Port)))
}

// Liveness is a simple endpoint to check if the server is alive
// @Summary Get liveness status of the microservice
// @Description Get liveness status of the microservice
// @Tags Kubernetes
// @Success 200 {string} string
// @Router /live [get]
func Liveness(c *gin.Context) {
	// Check if the config value 'broken' is set to 1
	if val, err := GetConfig("broken"); err == nil && val == "1" {
		c.String(http.StatusServiceUnavailable, "dead")
		return
	}

	c.String(http.StatusOK, "alive")
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
		c.String(http.StatusServiceUnavailable, "not ready")
	} else {
		c.String(http.StatusOK, "ready")
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
	if params.Error != "" {
		c.String(http.StatusBadRequest, params.Error)
		return
	}

	products, err := database.Products(params)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
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
	if params.Error != "" {
		c.String(http.StatusBadRequest, params.Error)
		return
	}

	product, err := database.Product(params)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
	} else {
		c.IndentedJSON(http.StatusOK, product)
	}
}
