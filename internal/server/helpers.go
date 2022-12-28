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
	reqType := c.DefaultQuery("full", "false")
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

	if val, err := strconv.ParseBool(reqType); err != nil {
		params.Error = gin.H{"error": "Invalid request type (should be true or false)"}
		return
	} else {
		params.Full = val
	}

	if reqSortDirection != "asc" && reqSortDirection != "desc" {
		params.Error = gin.H{"error": "Invalid sort direction (should be asc or desc)"}
		return
	} else {
		if reqSortDirection == "asc" {
			params.SortDirection = model.Ascending
		} else {
			params.SortDirection = model.Descending
		}
	}

	if _, ok := model.SortableBy()[reqSort]; !ok {
		params.Error = gin.H{"error": "Invalid sort parameter"}
		return
	} else {
		params.SortBy = model.SortBy(reqSort)
	}

	params.Filter = reqFilter
	params.ProductId = c.Param("id")

	return
}
