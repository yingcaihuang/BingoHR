package util

import (
	"github.com/gin-gonic/gin"
	"github.com/unknwon/com"
)

// GetPage get page parameter
func GetPage(c *gin.Context) int {
	page := c.DefaultQuery("page", "0")
	return com.StrTo(page).MustInt()
}

// GetLimit get limit parameter
func GetLimit(c *gin.Context) int {
	limit := c.DefaultQuery("limit", "0")
	return com.StrTo(limit).MustInt()
}

// GetCacheClear get cache_clear parameter
func GetCacheClear(c *gin.Context) int {
	cache_clear := c.DefaultQuery("cache_clear", "0")
	return com.StrTo(cache_clear).MustInt()
}
