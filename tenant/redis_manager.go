// Package tenant 提供了通用的多租户数据库和Redis管理器
package tenant

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/onebids/onecommon/consts"
)

// RedisConfig Redis管理器配置
type RedisConfig struct {
	// 默认Redis配置，当租户没有特定配置时使用
	DefaultOptions *redis.Options
	// 租户特定的Redis配置，key为租户ID，value为Redis连接选项
	TenantOptions map[string]*redis.Options
	// 是否启用租户隔离，如果启用，则会使用租户ID作为key前缀
	EnableTenantIsolation bool
	// 租户隔离时的分隔符
	TenantSeparator string
}

// NewDefaultRedisConfig 创建默认Redis配置
func NewDefaultRedisConfig() *RedisConfig {
	return &RedisConfig{
		DefaultOptions: &redis.Options{
			Addr:     "localhost:6379",
			Password: "",
			DB:       0,
		},
		TenantOptions:         make(map[string]*redis.Options),
		EnableTenantIsolation: true,
		TenantSeparator:       ":",
	}
}

// RedisManager Redis管理器接口
type RedisManager interface {
	// GetClient 获取Redis客户端
	// 如果tenantID为空，则返回默认客户端
	GetClient(ctx context.Context, tenantID string) *redis.Client

	// GetClientFromContext 从上下文中获取租户ID，然后获取对应的Redis客户端
	GetClientFromContext(ctx context.Context) *redis.Client

	// Close 关闭所有Redis连接
	Close() error

	// WithTenantPrefix 在key前面添加租户前缀
	WithTenantPrefix(ctx context.Context, key string) string
}

// tenantRedisManager Redis管理器实现
type tenantRedisManager struct {
	clients       map[string]*redis.Client
	mutex         sync.RWMutex
	config        *RedisConfig
	defaultClient *redis.Client
}

// NewRedisManager 创建Redis管理器
func NewRedisManager(config *RedisConfig) (RedisManager, error) {
	if config == nil {
		config = NewDefaultRedisConfig()
	}

	if config.DefaultOptions == nil {
		return nil, fmt.Errorf("default Redis options cannot be nil")
	}

	manager := &tenantRedisManager{
		clients: make(map[string]*redis.Client),
		config:  config,
	}

	// 创建默认客户端
	defaultClient := redis.NewClient(config.DefaultOptions)

	// 测试默认连接
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := defaultClient.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("failed to connect to default Redis: %w", err)
	}

	manager.defaultClient = defaultClient

	// 初始化租户客户端
	for tenantID, options := range config.TenantOptions {
		client := redis.NewClient(options)

		// 测试连接
		if err := client.Ping(ctx).Err(); err != nil {
			// 关闭已创建的连接
			_ = manager.Close()
			return nil, fmt.Errorf("failed to connect to Redis for tenant %s: %w", tenantID, err)
		}

		manager.clients[tenantID] = client
	}

	return manager, nil
}

// GetClient 获取Redis客户端
func (m *tenantRedisManager) GetClient(ctx context.Context, tenantID string) *redis.Client {
	if tenantID == "" {
		return m.defaultClient
	}

	m.mutex.RLock()
	client, exists := m.clients[tenantID]
	m.mutex.RUnlock()

	if exists {
		return client
	}

	// 如果没有找到特定租户的客户端，使用默认客户端
	return m.defaultClient
}

// GetClientFromContext 从上下文中获取租户ID，然后获取对应的Redis客户端
func (m *tenantRedisManager) GetClientFromContext(ctx context.Context) *redis.Client {
	tenantID := getTenantIDFromContext(ctx)
	return m.GetClient(ctx, tenantID)
}

// Close 关闭所有Redis连接
func (m *tenantRedisManager) Close() error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	var lastErr error

	// 关闭默认客户端
	if m.defaultClient != nil {
		if err := m.defaultClient.Close(); err != nil {
			lastErr = err
		}
	}

	// 关闭所有租户客户端
	for tenantID, client := range m.clients {
		if err := client.Close(); err != nil {
			lastErr = fmt.Errorf("failed to close Redis client for tenant %s: %w", tenantID, err)
		}
		delete(m.clients, tenantID)
	}

	return lastErr
}

// WithTenantPrefix 在key前面添加租户前缀
func (m *tenantRedisManager) WithTenantPrefix(ctx context.Context, key string) string {
	if !m.config.EnableTenantIsolation {
		return key
	}

	tenantID := getTenantIDFromContext(ctx)
	if tenantID == "" {
		return key
	}

	return fmt.Sprintf("%s%s%s", tenantID, m.config.TenantSeparator, key)
}

// getTenantIDFromContext 从上下文中获取租户ID
func getTenantIDFromContext(ctx context.Context) string {
	if ctx == nil {
		return ""
	}

	tenantID, ok := ctx.Value(consts.TenantID).(string)
	if !ok || tenantID == "" {
		return ""
	}

	return tenantID
}
