package cache_service

import (
	"strconv"
	"strings"
)

type Cache struct {
	Name    string
	Id      int
	Keyword string

	Page  int
	Limit int
}

func (c *Cache) GetTagsKey() string {
	keys := []string{
		"LIST",
		c.Name,
	}

	if c.Id > 0 {
		keys = append(keys, strconv.Itoa(c.Id))
	}
	if c.Keyword != "" {
		keys = append(keys, c.Keyword)
	}
	if c.Page > 0 {
		keys = append(keys, strconv.Itoa(c.Page))
	}
	if c.Limit > 0 {
		keys = append(keys, strconv.Itoa(c.Limit))
	}

	return strings.Join(keys, "_")
}
