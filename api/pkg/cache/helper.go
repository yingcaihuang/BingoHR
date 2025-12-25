package cache

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"github.com/gin-gonic/gin"
	"time"
)

// GenerateCacheKey 生成缓存key
func GenerateCacheKey(prefix string, params ...interface{}) string {
	hash := md5.New()
	for _, param := range params {
		hash.Write([]byte(fmt.Sprintf("%v", param)))
	}
	return prefix + ":" + hex.EncodeToString(hash.Sum(nil))
}

// CacheWithRetry 带重试的缓存操作
func CacheWithRetry(
	ctx *gin.Context,
	operation func() error,
	maxRetries int,
	delay time.Duration,
) error {
	var lastErr error

	for i := 0; i < maxRetries; i++ {
		err := operation()
		if err == nil {
			return nil
		}

		lastErr = err

		// 指数退避
		sleepTime := delay * (1 << uint(i))
		if sleepTime > 5*time.Second {
			sleepTime = 5 * time.Second
		}

		select {
		case <-time.After(sleepTime):
			continue
		case <-ctx.Done():
			return ctx.Err()
		}
	}

	return fmt.Errorf("操作失败，重试%d次后仍失败: %w", maxRetries, lastErr)
}

// CacheStats 获取缓存统计
func (cm *CacheManager) Stats() *CacheStats {
	return cm.stats
}

// ResetStats 重置统计
func (cm *CacheManager) ResetStats() {
	cm.stats = &CacheStats{}
}
