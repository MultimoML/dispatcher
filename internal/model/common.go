package model

import "github.com/gin-gonic/gin"

type SortBy string
type SortDirection int

const (
	Default  SortBy = "" // ID
	Name     SortBy = "name"
	Price    SortBy = "price"    // Latest price
	Category SortBy = "category" // Category name

	Ascending  SortDirection = 1
	Descending SortDirection = -1
)

type QueryParameters struct {
	ProductId     string
	Limit         int
	Offset        int
	Full          bool
	SortDirection SortDirection
	SortBy        SortBy
	Filter        string
	Error         gin.H
}
