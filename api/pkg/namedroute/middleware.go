package namedroute

import "github.com/gin-gonic/gin"

func Middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		method := c.Request.Method
		path := c.FullPath()
		if path == "" {
			c.Next()
			return
		}

		key := RouteKey(method + "::" + path)

		routeMu.RLock()
		name := routeNameMap[key]
		routeMu.RUnlock()

		if name != "" {
			c.Set("routeName", name)
		}
		c.Next()
	}
}
