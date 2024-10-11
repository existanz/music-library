package query

import (
	"strconv"

	"github.com/gin-gonic/gin"
)

var filterNames = []string{"group", "song", "releaseDate", "text", "link"}

type Paginator struct {
	Limit  string
	Offset string
}

type Filter struct {
	Field string
	Value string
}

type Options struct {
	Paginator Paginator
	Filters   []Filter
}

func GetPaginator(c *gin.Context) Paginator {
	page, _ := c.GetQuery("page")
	limit, _ := c.GetQuery("limit")

	p, err := strconv.Atoi(page)
	if err != nil {
		p = 1
	}

	l, err := strconv.Atoi(limit)
	if err != nil {
		l = 10
	}

	offset := (p - 1) * l

	return Paginator{
		Offset: strconv.Itoa(offset),
		Limit:  strconv.Itoa(l),
	}
}

func GetFilters(c *gin.Context) []Filter {
	filters := make([]Filter, 0, len(filterNames))

	for _, name := range filterNames {
		val, _ := c.GetQuery(name)
		if val != "" {
			filters = append(filters, Filter{
				Field: name,
				Value: val,
			})
		}
	}

	return filters
}

func GetOptions(c *gin.Context) Options {
	return Options{
		Paginator: GetPaginator(c),
		Filters:   GetFilters(c),
	}
}
