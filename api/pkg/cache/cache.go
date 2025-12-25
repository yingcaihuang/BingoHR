package cache

import (
	"time"

	"github.com/gin-gonic/gin"
)

const (
	CACHE_ARTICLE   = "ARTICLE"
	CACHE_TAG       = "TAG"
	CACHE_USER      = "USER"
	CACHE_ROLE      = "ROLE"
	CACHE_ROLE_PERM = "ROLE_PERM"
	CACHE_JOB       = "JOB"
	CACHE_RESUME    = "RESUME"
)

// Cache 定义缓存接口，便于后续扩展其他缓存实现
type Cache interface {
	// 基础操作
	Set(ctx *gin.Context, key string, value interface{}, expiration time.Duration) error
	Get(ctx *gin.Context, key string, dest interface{}) error
	Delete(ctx *gin.Context, keys ...string) error
	Exists(ctx *gin.Context, key string) (bool, error)
	Expire(ctx *gin.Context, key string, expiration time.Duration) error

	// 哈希表操作
	HSet(ctx *gin.Context, key string, values ...interface{}) error
	HGet(ctx *gin.Context, key, field string, dest interface{}) error
	HGetAll(ctx *gin.Context, key string) (map[string]string, error)
	HDel(ctx *gin.Context, key string, fields ...string) error

	// 集合操作
	SAdd(ctx *gin.Context, key string, members ...interface{}) error
	SRem(ctx *gin.Context, key string, members ...interface{}) error
	SMembers(ctx *gin.Context, key string) ([]string, error)
	SIsMember(ctx *gin.Context, key string, member interface{}) (bool, error)

	// 列表操作
	LPush(ctx *gin.Context, key string, values ...interface{}) error
	RPop(ctx *gin.Context, key string, dest interface{}) error
	LLen(ctx *gin.Context, key string) (int64, error)

	// 批量操作
	MGet(ctx *gin.Context, keys []string) ([]interface{}, error)

	// 高级功能
	Incr(ctx *gin.Context, key string) (int64, error)
	Decr(ctx *gin.Context, key string) (int64, error)
	Lock(ctx *gin.Context, key string, expiration time.Duration) (bool, error)
	Unlock(ctx *gin.Context, key string) error

	// 连接管理
	Ping(ctx *gin.Context) error
	Close() error
}
