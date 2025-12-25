package namedroute

import "github.com/gin-gonic/gin"

// Wrap 接收 *gin.Engine 或 *gin.RouterGroup，返回命名路由组
func Wrap(group *gin.RouterGroup) *NamedRouteGroup {
	return &NamedRouteGroup{group: group}
}

// New 从 *gin.Engine 创建根命名路由组
func New(r *gin.Engine) *NamedRouteGroup {
	return &NamedRouteGroup{group: &r.RouterGroup}
}
