package database

import (
	"context"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/multimoml/dispatcher/internal/model"
)

func Products(params *model.QueryParameters) ([]model.Product, error) {
	ctx := context.TODO()
	productCollection := dbClient.Database("products").Collection("spar")

	// Set up query options
	findOptions := options.Find()
	if params.Limit != -1 {
		findOptions.SetLimit(int64(params.Limit))
	}

	if params.Offset != 0 {
		findOptions.SetSkip(int64(params.Offset))
	}

	// Set depth options (default is Full, needs no projection)
	switch params.History {
	case model.None:
		findOptions.SetProjection(bson.D{{"price-in-time", 0}})
	case model.Last:
		findOptions.SetProjection(bson.D{{"price-in-time", bson.D{{"$slice", -1}}}})
	}

	// Set sort options
	if params.SortBy != model.Default {
		sortBy := model.MapQueryParamToDbField(params.SortBy)
		findOptions.SetSort(bson.D{{sortBy, params.SortOrder}})
	}

	// Set filter options
	var filter bson.M
	if params.Category != "" {
		category := model.MapQueryParamToDbField(model.Category)
		filter = bson.M{category: params.Category}
	}

	// Execute query
	cursor, err := productCollection.Find(ctx, filter, findOptions)
	if err != nil {
		log.Println("Failed to query elements in database: ", err)
		return nil, err
	}

	// Unmarshal results
	var products []model.Product
	if err = cursor.All(ctx, &products); err != nil {
		log.Println("Failed to unmarshal query results: ", err)
		return nil, err
	}

	return products, nil
}

func Product(params *model.QueryParameters) (model.Product, error) {
	ctx := context.TODO()
	productCollection := dbClient.Database("products").Collection("spar")

	findOptions := options.FindOne()

	// Set depth options (default is Full, needs no projection)
	switch params.History {
	case model.None:
		findOptions.SetProjection(bson.D{{"price-in-time", 0}})
	case model.Last:
		findOptions.SetProjection(bson.D{{"price-in-time", bson.D{{"$slice", -1}}}})
	}

	// Execute query
	var product model.Product
	id, err := primitive.ObjectIDFromHex(params.ProductId)
	if err != nil {
		log.Println("Invalid product id: ", err)
		return product, err
	}

	res := productCollection.FindOne(ctx, bson.D{{"_id", id}}, findOptions)

	// If error is ErrNoDocuments, return empty product
	if res.Err() == mongo.ErrNoDocuments {
		return product, nil
	}

	if err := res.Decode(&product); err != nil {
		log.Println("Failed to unmarshal query results: ", err)
		return product, err
	}

	return product, nil
}
