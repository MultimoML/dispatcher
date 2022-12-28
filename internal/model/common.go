package model

import "github.com/gin-gonic/gin"

type SortBy string
type SortDirection int
type Depth int

type QueryParameters struct {
	ProductId     string
	Limit         int
	Offset        int
	Depth         Depth
	SortDirection SortDirection
	SortBy        SortBy
	Filter        string
	Error         gin.H
}

const (
	Default  SortBy = "" // ID
	Name     SortBy = "name"
	Price    SortBy = "price"    // Latest price
	Category SortBy = "category" // Category name

	Ascending  SortDirection = 1
	Descending SortDirection = -1

	Full Depth = 0
	Last Depth = 1
	None Depth = 2
)

func Sorts() map[string]SortBy {
	return map[string]SortBy{
		"":         Default,
		"name":     Name,
		"price":    Price,
		"category": Category,
	}
}

func Directions() map[string]SortDirection {
	return map[string]SortDirection{
		"asc":  Ascending,
		"desc": Descending,
	}
}

func Depths() map[string]Depth {
	return map[string]Depth{
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
