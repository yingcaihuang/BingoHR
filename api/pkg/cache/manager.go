package cache

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"time"
)

// CacheManager 缓存管理器
type CacheManager struct {
	cache *RedisCache
	stats *CacheStats
}

// CacheStats 缓存统计
type CacheStats struct {
	Hits   int64
	Misses int64
	Sets   int64
	Dels   int64
}

// NewCacheManager 创建缓存管理器
func NewCacheManager(cache *RedisCache) *CacheManager {
	return &CacheManager{
		cache: cache,
		stats: &CacheStats{},
	}
}

// GetWithFallback 带降级的缓存获取
func (cm *CacheManager) GetWithFallback(
	ctx *gin.Context,
	key string,
	dest interface{},
	fallback func() (interface{}, error),
	expiration time.Duration,
) error {
	// 尝试从缓存获取
	err := cm.cache.Get(ctx, key, dest)
	if err == nil {
		cm.stats.Hits++
		return nil
	}

	// 缓存未命中，调用fallback函数
	cm.stats.Misses++

	val, err := fallback()
	if err != nil {
		return fmt.Errorf("fallback失败: %w", err)
	}

	// 设置缓存
	err = cm.cache.Set(ctx, key, val, expiration)
	if err != nil {
		log.Printf("设置缓存失败: %v", err)
		// 继续执行，不返回错误
	}

	// 将结果赋值给dest
	data, err := json.Marshal(val)
	if err != nil {
		return err
	}

	return json.Unmarshal(data, dest)
}

// GetOrSet 获取或设置缓存
func (cm *CacheManager) GetOrSet(
	ctx *gin.Context,
	key string,
	dest interface{},
	loader func() (interface{}, error),
	expiration time.Duration,
) error {
	// 尝试获取
	err := cm.cache.Get(ctx, key, dest)
	if err == nil {
		cm.stats.Hits++
		return nil
	}

	cm.stats.Misses++

	// 获取锁，防止缓存击穿
	lockKey := "lock:" + key
	locked, err := cm.cache.Lock(ctx, lockKey, 5*time.Second)
	if err != nil || !locked {
		// 获取锁失败，可能重试或直接调用loader
		log.Printf("获取锁失败: %v", err)
	}

	defer cm.cache.Unlock(ctx, lockKey)

	// 再次检查缓存（双重检查锁）
	err = cm.cache.Get(ctx, key, dest)
	if err == nil {
		return nil
	}

	// 加载数据
	val, err := loader()
	if err != nil {
		return fmt.Errorf("加载数据失败: %w", err)
	}

	// 设置缓存
	err = cm.cache.Set(ctx, key, val, expiration)
	if err != nil {
		log.Printf("设置缓存失败: %v", err)
	}

	// 返回结果
	data, err := json.Marshal(val)
	if err != nil {
		return err
	}

	return json.Unmarshal(data, dest)
}
