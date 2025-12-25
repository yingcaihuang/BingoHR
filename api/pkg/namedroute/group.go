package namedroute

import (
	"strings"
	"sync"

	"github.com/gin-gonic/gin"
)

type RouteKey string

var (
	routeNameMap = make(map[RouteKey]string)
	routeMu      sync.RWMutex
)

//	GetRouteNameMap 返回路由名称及其uri定义
//
// "GET::/api/v2/user/list": "rest.user.list",
// "POST::/api/v2/user/create": "rest.user.create",
// "PUT::/api/v2/user/update": "rest.user.update"
// "DELETE::/api/v2/user/delete/:id": "rest.user.delete",
// "GET::/api/v2/user/perms": "rest.user.perms",
// "DELETE::/api/v2/job/delete/:id": "rest.job.delete",
func GetRouteNameMap() map[string]string {
	routeMu.RLock()
	defer routeMu.RUnlock()
	m := make(map[string]string, len(routeNameMap))
	for k, v := range routeNameMap {
		m[string(k)] = v
	}
	return m
}

// GetRouteName 供 Handler 使用
func GetRouteName(c *gin.Context) string {
	if v, ok := c.Get("routeName"); ok {
		if name, ok := v.(string); ok {
			return name
		}
	}
	return ""
}

// NamedRouteGroup 包装 gin.RouterGroup
type NamedRouteGroup struct {
	group        *gin.RouterGroup
	method, path string
}

// Name 将当前路由注册名称
func (nrg *NamedRouteGroup) Name(name string) *NamedRouteGroup {
	if nrg.method != "" && nrg.path != "" {
		key := RouteKey(nrg.method + "::" + nrg.path)
		routeMu.Lock()
		routeNameMap[key] = name
		routeMu.Unlock()
	}
	return nrg
}

// --- 封装 HTTP 方法 ---
func (nrg *NamedRouteGroup) GET(relativePath string, handlers ...gin.HandlerFunc) *NamedRouteGroup {
	fullPath := joinPath(nrg.group.BasePath(), relativePath)
	nrg.group.GET(relativePath, handlers...)
	return &NamedRouteGroup{
		group:  nrg.group,
		method: "GET",
		path:   fullPath,
	}
}

func (nrg *NamedRouteGroup) POST(relativePath string, handlers ...gin.HandlerFunc) *NamedRouteGroup {
	fullPath := joinPath(nrg.group.BasePath(), relativePath)
	nrg.group.POST(relativePath, handlers...)
	return &NamedRouteGroup{
		group:  nrg.group,
		method: "POST",
		path:   fullPath,
	}
}

func (nrg *NamedRouteGroup) PUT(relativePath string, handlers ...gin.HandlerFunc) *NamedRouteGroup {
	fullPath := joinPath(nrg.group.BasePath(), relativePath)
	nrg.group.PUT(relativePath, handlers...)
	return &NamedRouteGroup{
		group:  nrg.group,
		method: "PUT",
		path:   fullPath,
	}
}

func (nrg *NamedRouteGroup) DELETE(relativePath string, handlers ...gin.HandlerFunc) *NamedRouteGroup {
	fullPath := joinPath(nrg.group.BasePath(), relativePath)
	nrg.group.DELETE(relativePath, handlers...)
	return &NamedRouteGroup{
		group:  nrg.group,
		method: "DELETE",
		path:   fullPath,
	}
}

// Group 支持嵌套分组
func (nrg *NamedRouteGroup) Group(prefix string, handlers ...gin.HandlerFunc) *NamedRouteGroup {
	newGroup := nrg.group.Group(prefix, handlers...)
	return &NamedRouteGroup{group: newGroup}
}

// Use 添加中间件
func (nrg *NamedRouteGroup) Use(middleware ...gin.HandlerFunc) *NamedRouteGroup {
	nrg.group.Use(middleware...)
	return nrg
}

func joinPath(base, rel string) string {
	if !strings.HasPrefix(rel, "/") {
		rel = "/" + rel
	}
	if base == "/" {
		return rel
	}
	return base + rel
}
