package utils

import (
	"strconv"

	"github.com/gin-gonic/gin"
)

// https://go.dev/tour/generics/1
func Paginate[T any](c *gin.Context, items []T) (data []T, pageInt, limitInt int) {
	var err error
	page := c.Query("page")
	limit := c.Query("limit")

	if page == "" {
		page = "1"
	}
	if limit == "" {
		limit = "10"
	}

	pageInt, err = strconv.Atoi(page)
	if err != nil || pageInt < 1 {
		pageInt = 1
	}

	limitInt, err = strconv.Atoi(limit)
	if err != nil || limitInt < 1 {
		limitInt = 10
	}

	start := (pageInt - 1) * limitInt
	end := start + limitInt

	if start > len(items) {
		return []T{}, pageInt, limitInt
	}
	if end > len(items) {
		end = len(items)
	}

	return items[start:end], pageInt, limitInt
}
