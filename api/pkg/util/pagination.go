package util

import (
	"github.com/gin-gonic/gin"
	"github.com/unknwon/com"
)

// GetPage get page parameter
func GetPage(c *gin.Context) int {
	return GetIntQuery("page", c)
}

// GetLimit get limit parameter
func GetLimit(c *gin.Context) int {
	return GetIntQuery("limit", c)
}

// GetIntQuery get a query param named by name, it is int
func GetIntQuery(name string, c *gin.Context) int {
	page := c.DefaultQuery(name, "0")
	return com.StrTo(page).MustInt()
}

// GetCacheClear get cache_clear parameter
func GetCacheClear(c *gin.Context) int {
	cache_clear := c.DefaultQuery("cache_clear", "0")
	return com.StrTo(cache_clear).MustInt()
}
