package model

import "github.com/gin-gonic/gin"

type SortBy string
type SortOrder int
type History int

type QueryParameters struct {
	ProductId string
	Limit     int
	Offset    int
	History   History
	SortOrder SortOrder
	SortBy    SortBy
	Category  string
	Error     gin.H
}

const (
	Default  SortBy = "" // ID
	Name     SortBy = "name"
	Price    SortBy = "price"    // Latest price
	Category SortBy = "category" // Category name

	Ascending  SortOrder = 1
	Descending SortOrder = -1

	Full History = 0
	Last History = 1
	None History = 2
)

func SortableBy() map[string]SortBy {
	return map[string]SortBy{
		"":         Default,
		"name":     Name,
		"price":    Price,
		"category": Category,
	}
}

func SortOrders() map[string]SortOrder {
	return map[string]SortOrder{
		"asc":  Ascending,
		"desc": Descending,
	}
}

func ResultHistory() map[string]History {
	return map[string]History{
		"full": Full,
		"last": Last,
		"none": None,
	}
}

func MapQueryParamToDbField(param SortBy) string {
	switch param {
	case "name":
		return "name"
	case "price":
		return "price-in-time.price"
	case "category":
		return "category-name"
	default:
		return ""
	}
}
