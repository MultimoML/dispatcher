package server

import (
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/multimoml/dispatcher/internal/model"
)

func parseQuery(c *gin.Context) (params *model.QueryParameters) {
	params = &model.QueryParameters{}

	// Get query parameters
	reqLimit := c.DefaultQuery("limit", "-1")
	reqOffset := c.DefaultQuery("offset", "0")
	reqDepth := c.DefaultQuery("depth", "last")
	reqSortDirection := c.DefaultQuery("sortDirection", "asc")
	reqSort := c.DefaultQuery("sortBy", "")
	reqFilter := c.DefaultQuery("filter", "")

	// Check parameter validity
	if val, err := strconv.Atoi(reqLimit); err != nil || (val < 0 && val != -1) {
		params.Error = gin.H{"error": "Invalid limit"}
		return
	} else {
		params.Limit = val
	}

	if val, err := strconv.Atoi(reqOffset); err != nil || val < 0 {
		params.Error = gin.H{"error": "Invalid offset"}
		return
	} else {
		params.Offset = val
	}

	if val, ok := model.Depths()[reqDepth]; !ok {
		params.Error = gin.H{"error": "Invalid depth"}
		return
	} else {
		params.Depth = val
	}

	if val, ok := model.Directions()[reqSortDirection]; !ok {
		params.Error = gin.H{"error": "Invalid sort direction"}
		return
	} else {
		params.SortDirection = val
	}

	if _, ok := model.Sorts()[reqSort]; !ok {
		params.Error = gin.H{"error": "Invalid sort parameter"}
		return
	} else {
		params.SortBy = model.SortBy(reqSort)
	}

	params.Filter = reqFilter
	params.ProductId = c.Param("id")

	return
}
