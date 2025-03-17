// Package tenant 提供了通用的多租户数据库和Redis管理器
package tenant

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
)

// RedisHelper Redis辅助工具接口
// 提供常用的Redis操作封装，自动处理租户隔离
type RedisHelper interface {
	// Set 设置键值对
	Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error

	// Get 获取值
	Get(ctx context.Context, key string) (string, error)

	// GetObject 获取并反序列化对象
	GetObject(ctx context.Context, key string, dest interface{}) error

	// Delete 删除键
	Delete(ctx context.Context, key string) error

	// Exists 检查键是否存在
	Exists(ctx context.Context, key string) (bool, error)

	// Expire 设置过期时间
	Expire(ctx context.Context, key string, expiration time.Duration) error

	// Incr 自增
	Incr(ctx context.Context, key string) (int64, error)

	// HSet 设置哈希表字段
	HSet(ctx context.Context, key, field string, value interface{}) error

	// HGet 获取哈希表字段
	HGet(ctx context.Context, key, field string) (string, error)

	// HGetAll 获取哈希表所有字段
	HGetAll(ctx context.Context, key string) (map[string]string, error)

	// HDel 删除哈希表字段
	HDel(ctx context.Context, key string, fields ...string) error

	// LPush 将值推入列表左端
	LPush(ctx context.Context, key string, values ...interface{}) error

	// RPush 将值推入列表右端
	RPush(ctx context.Context, key string, values ...interface{}) error

	// LRange 获取列表范围
	LRange(ctx context.Context, key string, start, stop int64) ([]string, error)

	// SAdd 添加集合成员
	SAdd(ctx context.Context, key string, members ...interface{}) error

	// SMembers 获取集合所有成员
	SMembers(ctx context.Context, key string) ([]string, error)

	// SRem 移除集合成员
	SRem(ctx context.Context, key string, members ...interface{}) error

	// ZAdd 添加有序集合成员
	ZAdd(ctx context.Context, key string, members ...*redis.Z) error

	// ZRange 获取有序集合范围
	ZRange(ctx context.Context, key string, start, stop int64) ([]string, error)

	// ZRangeWithScores 获取有序集合范围及分数
	ZRangeWithScores(ctx context.Context, key string, start, stop int64) ([]redis.Z, error)

	// ZRem 移除有序集合成员
	ZRem(ctx context.Context, key string, members ...interface{}) error

	// Lock 获取分布式锁
	Lock(ctx context.Context, key string, value string, expiration time.Duration) (bool, error)

	// Unlock 释放分布式锁
	Unlock(ctx context.Context, key string, value string) (bool, error)
}

// tenantRedisHelper Redis辅助工具实现
type tenantRedisHelper struct {
	manager RedisManager
}

// NewRedisHelper 创建Redis辅助工具
func NewRedisHelper(manager RedisManager) RedisHelper {
	return &tenantRedisHelper{
		manager: manager,
	}
}

// Set 设置键值对
func (h *tenantRedisHelper) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	client := h.manager.GetClientFromContext(ctx)
	key = h.manager.WithTenantPrefix(ctx, key)

	var strValue string
	switch v := value.(type) {
	case string:
		strValue = v
	case []byte:
		strValue = string(v)
	default:
		// 序列化对象
		bytes, err := json.Marshal(value)
		if err != nil {
			return fmt.Errorf("failed to marshal value: %w", err)
		}
		strValue = string(bytes)
	}

	return client.Set(ctx, key, strValue, expiration).Err()
}

// Get 获取值
func (h *tenantRedisHelper) Get(ctx context.Context, key string) (string, error) {
	client := h.manager.GetClientFromContext(ctx)
	key = h.manager.WithTenantPrefix(ctx, key)

	result, err := client.Get(ctx, key).Result()
	if err == redis.Nil {
		return "", nil // 键不存在
	}
	return result, err
}

// GetObject 获取并反序列化对象
func (h *tenantRedisHelper) GetObject(ctx context.Context, key string, dest interface{}) error {
	value, err := h.Get(ctx, key)
	if err != nil {
		return err
	}

	if value == "" {
		return redis.Nil
	}

	return json.Unmarshal([]byte(value), dest)
}

// Delete 删除键
func (h *tenantRedisHelper) Delete(ctx context.Context, key string) error {
	client := h.manager.GetClientFromContext(ctx)
	key = h.manager.WithTenantPrefix(ctx, key)

	return client.Del(ctx, key).Err()
}

// Exists 检查键是否存在
func (h *tenantRedisHelper) Exists(ctx context.Context, key string) (bool, error) {
	client := h.manager.GetClientFromContext(ctx)
	key = h.manager.WithTenantPrefix(ctx, key)

	result, err := client.Exists(ctx, key).Result()
	return result > 0, err
}

// Expire 设置过期时间
func (h *tenantRedisHelper) Expire(ctx context.Context, key string, expiration time.Duration) error {
	client := h.manager.GetClientFromContext(ctx)
	key = h.manager.WithTenantPrefix(ctx, key)

	return client.Expire(ctx, key, expiration).Err()
}

// Incr 自增
func (h *tenantRedisHelper) Incr(ctx context.Context, key string) (int64, error) {
	client := h.manager.GetClientFromContext(ctx)
	key = h.manager.WithTenantPrefix(ctx, key)

	return client.Incr(ctx, key).Result()
}

// HSet 设置哈希表字段
func (h *tenantRedisHelper) HSet(ctx context.Context, key, field string, value interface{}) error {
	client := h.manager.GetClientFromContext(ctx)
	key = h.manager.WithTenantPrefix(ctx, key)

	var strValue string
	switch v := value.(type) {
	case string:
		strValue = v
	case []byte:
		strValue = string(v)
	default:
		// 序列化对象
		bytes, err := json.Marshal(value)
		if err != nil {
			return fmt.Errorf("failed to marshal value: %w", err)
		}
		strValue = string(bytes)
	}

	return client.HSet(ctx, key, field, strValue).Err()
}

// HGet 获取哈希表字段
func (h *tenantRedisHelper) HGet(ctx context.Context, key, field string) (string, error) {
	client := h.manager.GetClientFromContext(ctx)
	key = h.manager.WithTenantPrefix(ctx, key)

	result, err := client.HGet(ctx, key, field).Result()
	if err == redis.Nil {
		return "", nil // 字段不存在
	}
	return result, err
}

// HGetAll 获取哈希表所有字段
func (h *tenantRedisHelper) HGetAll(ctx context.Context, key string) (map[string]string, error) {
	client := h.manager.GetClientFromContext(ctx)
	key = h.manager.WithTenantPrefix(ctx, key)

	return client.HGetAll(ctx, key).Result()
}

// HDel 删除哈希表字段
func (h *tenantRedisHelper) HDel(ctx context.Context, key string, fields ...string) error {
	client := h.manager.GetClientFromContext(ctx)
	key = h.manager.WithTenantPrefix(ctx, key)

	return client.HDel(ctx, key, fields...).Err()
}

// LPush 将值推入列表左端
func (h *tenantRedisHelper) LPush(ctx context.Context, key string, values ...interface{}) error {
	client := h.manager.GetClientFromContext(ctx)
	key = h.manager.WithTenantPrefix(ctx, key)

	return client.LPush(ctx, key, values...).Err()
}

// RPush 将值推入列表右端
func (h *tenantRedisHelper) RPush(ctx context.Context, key string, values ...interface{}) error {
	client := h.manager.GetClientFromContext(ctx)
	key = h.manager.WithTenantPrefix(ctx, key)

	return client.RPush(ctx, key, values...).Err()
}

// LRange 获取列表范围
func (h *tenantRedisHelper) LRange(ctx context.Context, key string, start, stop int64) ([]string, error) {
	client := h.manager.GetClientFromContext(ctx)
	key = h.manager.WithTenantPrefix(ctx, key)

	return client.LRange(ctx, key, start, stop).Result()
}

// SAdd 添加集合成员
func (h *tenantRedisHelper) SAdd(ctx context.Context, key string, members ...interface{}) error {
	client := h.manager.GetClientFromContext(ctx)
	key = h.manager.WithTenantPrefix(ctx, key)

	return client.SAdd(ctx, key, members...).Err()
}

// SMembers 获取集合所有成员
func (h *tenantRedisHelper) SMembers(ctx context.Context, key string) ([]string, error) {
	client := h.manager.GetClientFromContext(ctx)
	key = h.manager.WithTenantPrefix(ctx, key)

	return client.SMembers(ctx, key).Result()
}

// SRem 移除集合成员
func (h *tenantRedisHelper) SRem(ctx context.Context, key string, members ...interface{}) error {
	client := h.manager.GetClientFromContext(ctx)
	key = h.manager.WithTenantPrefix(ctx, key)

	return client.SRem(ctx, key, members...).Err()
}

// ZAdd 添加有序集合成员
func (h *tenantRedisHelper) ZAdd(ctx context.Context, key string, members ...*redis.Z) error {
	client := h.manager.GetClientFromContext(ctx)
	key = h.manager.WithTenantPrefix(ctx, key)

	return client.ZAdd(ctx, key, members...).Err()
}

// ZRange 获取有序集合范围
func (h *tenantRedisHelper) ZRange(ctx context.Context, key string, start, stop int64) ([]string, error) {
	client := h.manager.GetClientFromContext(ctx)
	key = h.manager.WithTenantPrefix(ctx, key)

	return client.ZRange(ctx, key, start, stop).Result()
}

// ZRangeWithScores 获取有序集合范围及分数
func (h *tenantRedisHelper) ZRangeWithScores(ctx context.Context, key string, start, stop int64) ([]redis.Z, error) {
	client := h.manager.GetClientFromContext(ctx)
	key = h.manager.WithTenantPrefix(ctx, key)

	return client.ZRangeWithScores(ctx, key, start, stop).Result()
}

// ZRem 移除有序集合成员
func (h *tenantRedisHelper) ZRem(ctx context.Context, key string, members ...interface{}) error {
	client := h.manager.GetClientFromContext(ctx)
	key = h.manager.WithTenantPrefix(ctx, key)

	return client.ZRem(ctx, key, members...).Err()
}

// Lock 获取分布式锁
func (h *tenantRedisHelper) Lock(ctx context.Context, key string, value string, expiration time.Duration) (bool, error) {
	client := h.manager.GetClientFromContext(ctx)
	key = h.manager.WithTenantPrefix(ctx, "lock:"+key)

	return client.SetNX(ctx, key, value, expiration).Result()
}

// Unlock 释放分布式锁
func (h *tenantRedisHelper) Unlock(ctx context.Context, key string, value string) (bool, error) {
	client := h.manager.GetClientFromContext(ctx)
	key = h.manager.WithTenantPrefix(ctx, "lock:"+key)

	// 使用Lua脚本确保只删除自己的锁
	script := `
		if redis.call("get", KEYS[1]) == ARGV[1] then
			return redis.call("del", KEYS[1])
		else
			return 0
		end
	`

	result, err := client.Eval(ctx, script, []string{key}, value).Result()
	if err != nil {
		return false, err
	}

	return result.(int64) == 1, nil
}
