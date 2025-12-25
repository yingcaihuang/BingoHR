package cache

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"hr-api/pkg/setting"
	"sync"
	"time"

	"github.com/go-redis/redis/v8"
)

var (
	instance *RedisCache
	once     sync.Once
	mu       sync.RWMutex
)

// RedisCache Redis缓存实现
type RedisCache struct {
	client *redis.Client
}

// Get 返回 RedisCache 单例
func GetInstance() (*RedisCache, error) {
	if instance == nil {
		return nil, errors.New("RedisCache 尚未初始化，请先调用 cache.Init()")
	}
	return instance, nil
}

// Init 初始化Redis连接
func Init() error {
	mu.Lock()
	defer mu.Unlock()

	if instance != nil {
		return nil
	}

	client := redis.NewClient(&redis.Options{
		Addr:         setting.RedisSetting.Host,
		DB:           setting.RedisSetting.DB,
		Password:     setting.RedisSetting.Password,
		MaxRetries:   setting.RedisSetting.MaxRetries,
		PoolSize:     setting.RedisSetting.PoolSize,
		PoolTimeout:  setting.RedisSetting.PoolTimeout,
		DialTimeout:  setting.RedisSetting.DialTimeout,
		ReadTimeout:  setting.RedisSetting.ReadTimeout,
		WriteTimeout: setting.RedisSetting.WriteTimeout,
	})

	// 测试连接
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := client.Ping(ctx).Err(); err != nil {
		return fmt.Errorf("redis连接测试失败: %w", err)
	}

	instance = &RedisCache{
		client: client,
	}

	return nil
}

// Set 设置缓存
func (r *RedisCache) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	if key == "" {
		return errors.New("key不能为空")
	}

	var err error
	switch v := value.(type) {
	case string:
		err = r.client.Set(ctx, key, v, expiration).Err()
	case []byte:
		err = r.client.Set(ctx, key, v, expiration).Err()
	case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64, float32, float64, bool:
		err = r.client.Set(ctx, key, v, expiration).Err()
	default:
		// 尝试JSON序列化
		data, marshalErr := json.Marshal(v)
		if marshalErr != nil {
			return fmt.Errorf("JSON序列化失败: %w", marshalErr)
		}
		err = r.client.Set(ctx, key, data, expiration).Err()
	}

	if err != nil {
		return fmt.Errorf("设置缓存失败[key=%s]: %w", key, err)
	}

	return nil
}

// Get 获取缓存
func (r *RedisCache) Get(ctx context.Context, key string, dest interface{}) error {
	val, err := r.client.Get(ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			return ErrCacheMiss
		}
		return fmt.Errorf("获取缓存失败[key=%s]: %w", key, err)
	}

	// 根据目标类型进行解析
	switch d := dest.(type) {
	case *string:
		*d = val
	case *[]byte:
		*d = []byte(val)
	case *int:
		_, err = fmt.Sscanf(val, "%d", d)
	case *int64:
		_, err = fmt.Sscanf(val, "%d", d)
	case *float64:
		_, err = fmt.Sscanf(val, "%f", d)
	case *bool:
		*d = val == "true"
	default:
		// 尝试JSON反序列化
		err = json.Unmarshal([]byte(val), dest)
	}

	if err != nil {
		return fmt.Errorf("解析缓存值失败[key=%s]: %w", key, err)
	}

	return nil
}

// GetString 获取字符串值（简化方法）
func (r *RedisCache) GetString(ctx context.Context, key string) (string, error) {
	var result string
	err := r.Get(ctx, key, &result)
	return result, err
}

// GetInt 获取整数值
func (r *RedisCache) GetInt(ctx context.Context, key string) (int64, error) {
	var result int64
	err := r.Get(ctx, key, &result)
	return result, err
}

// GetJSON 获取并反序列化JSON
func (r *RedisCache) GetJSON(ctx context.Context, key string, dest interface{}) error {
	var data string
	if err := r.Get(ctx, key, &data); err != nil {
		return err
	}
	return json.Unmarshal([]byte(data), dest)
}

// Delete 删除缓存
func (r *RedisCache) Delete(ctx context.Context, keys ...string) error {
	if len(keys) == 0 {
		return nil
	}

	err := r.client.Del(ctx, keys...).Err()
	if err != nil {
		return fmt.Errorf("删除缓存失败[keys=%v]: %w", keys, err)
	}

	return nil
}

// Exists 检查key是否存在
func (r *RedisCache) Exists(ctx context.Context, key string) (bool, error) {
	result, err := r.client.Exists(ctx, key).Result()
	if err != nil {
		return false, fmt.Errorf("检查缓存是否存在失败[key=%s]: %w", key, err)
	}

	return result > 0, nil
}

// Expire 设置过期时间
func (r *RedisCache) Expire(ctx context.Context, key string, expiration time.Duration) error {
	ok, err := r.client.Expire(ctx, key, expiration).Result()
	if err != nil {
		return fmt.Errorf("设置过期时间失败[key=%s]: %w", key, err)
	}

	if !ok {
		return errors.New("key不存在或设置过期时间失败")
	}

	return nil
}

// 哈希表操作
func (r *RedisCache) HSet(ctx context.Context, key string, values ...interface{}) error {
	err := r.client.HSet(ctx, key, values...).Err()
	if err != nil {
		return fmt.Errorf("设置哈希表失败[key=%s]: %w", key, err)
	}

	return nil
}

func (r *RedisCache) HGet(ctx context.Context, key, field string, dest interface{}) error {
	val, err := r.client.HGet(ctx, key, field).Result()
	if err != nil {
		if err == redis.Nil {
			return ErrCacheMiss
		}
		return fmt.Errorf("获取哈希字段失败[key=%s, field=%s]: %w", key, field, err)
	}

	// 解析值
	switch d := dest.(type) {
	case *string:
		*d = val
	default:
		err = json.Unmarshal([]byte(val), dest)
	}

	return err
}

func (r *RedisCache) HGetAll(ctx context.Context, key string) (map[string]string, error) {
	result, err := r.client.HGetAll(ctx, key).Result()
	if err != nil {
		return nil, fmt.Errorf("获取哈希表失败[key=%s]: %w", key, err)
	}

	return result, nil
}

// 集合操作
func (r *RedisCache) SAdd(ctx context.Context, key string, members ...interface{}) error {
	err := r.client.SAdd(ctx, key, members...).Err()
	if err != nil {
		return fmt.Errorf("添加集合成员失败[key=%s]: %w", key, err)
	}

	return nil
}

func (r *RedisCache) SMembers(ctx context.Context, key string) ([]string, error) {
	members, err := r.client.SMembers(ctx, key).Result()
	if err != nil {
		return nil, fmt.Errorf("获取集合成员失败[key=%s]: %w", key, err)
	}

	return members, nil
}

// 列表操作
func (r *RedisCache) LPush(ctx context.Context, key string, values ...interface{}) error {
	err := r.client.LPush(ctx, key, values...).Err()
	if err != nil {
		return fmt.Errorf("列表左推失败[key=%s]: %w", key, err)
	}

	return nil
}

// 批量获取
func (r *RedisCache) MGet(ctx context.Context, keys []string) ([]interface{}, error) {
	if len(keys) == 0 {
		return []interface{}{}, nil
	}

	vals, err := r.client.MGet(ctx, keys...).Result()
	if err != nil {
		return nil, fmt.Errorf("批量获取失败[keys=%v]: %w", keys, err)
	}

	return vals, nil
}

// 原子操作
func (r *RedisCache) Incr(ctx context.Context, key string) (int64, error) {
	val, err := r.client.Incr(ctx, key).Result()
	if err != nil {
		return 0, fmt.Errorf("自增失败[key=%s]: %w", key, err)
	}

	return val, nil
}

func (r *RedisCache) Decr(ctx context.Context, key string) (int64, error) {
	val, err := r.client.Decr(ctx, key).Result()
	if err != nil {
		return 0, fmt.Errorf("自减失败[key=%s]: %w", key, err)
	}

	return val, nil
}

// 分布式锁
func (r *RedisCache) Lock(ctx context.Context, key string, expiration time.Duration) (bool, error) {
	// 使用SET NX EX实现分布式锁
	result, err := r.client.SetNX(ctx, key, "1", expiration).Result()
	if err != nil {
		return false, fmt.Errorf("获取锁失败[key=%s]: %w", key, err)
	}

	return result, nil
}

func (r *RedisCache) Unlock(ctx context.Context, key string) error {
	err := r.client.Del(ctx, key).Err()
	if err != nil {
		return fmt.Errorf("释放锁失败[key=%s]: %w", key, err)
	}

	return nil
}

// 连接管理
func (r *RedisCache) Ping(ctx context.Context) error {
	return r.client.Ping(ctx).Err()
}

func (r *RedisCache) Close() error {
	return r.client.Close()
}

// 获取客户端统计信息
func (r *RedisCache) Stats() *redis.PoolStats {
	return r.client.PoolStats()
}

// 错误定义
var (
	ErrCacheMiss = errors.New("缓存未命中")
	ErrCacheType = errors.New("缓存类型错误")
)
