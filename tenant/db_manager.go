package tenant

import (
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/plugin/opentelemetry/tracing"
)

// MigrateFunc 定义数据库迁移函数类型
type MigrateFunc func(db *gorm.DB) error

// DBConfig 数据库管理器配置
type DBConfig struct {
	// DSN连接字符串或模板
	// 对于多租户情况，可以是带有%s占位符的模板，如：
	// "user:pass@tcp(host:port)/%s?charset=utf8mb4&parseTime=True&loc=Local"
	// 默认值: 空字符串，必须由用户提供
	DSNTemplate string

	// 预配置的特定实例DSN映射
	// 键为租户ID，值为完整DSN
	// 默认值: 空map
	TenantDSNs map[string]string

	// 数据库默认配置
	// 默认值: 见NewDefaultDBConfig函数
	DBConfig *gorm.Config

	// 是否启用追踪
	// 默认值: false
	EnableTracing bool

	// 日志级别
	// 默认值: logger.Error
	LogLevel logger.LogLevel

	// 慢查询阈值
	// 默认值: 200 * time.Millisecond
	SlowThreshold time.Duration

	// 迁移函数
	// 默认值: nil (不执行迁移)
	MigrateFunc MigrateFunc
}

// NewDefaultDBConfig 创建带有默认值的配置
func NewDefaultDBConfig() *DBConfig {
	return &DBConfig{
		TenantDSNs:    make(map[string]string),
		EnableTracing: false,
		LogLevel:      logger.Error,
		SlowThreshold: 200 * time.Millisecond,
		DBConfig: &gorm.Config{
			PrepareStmt:            true,
			SkipDefaultTransaction: true,
			Logger: logger.New(
				log.New(os.Stdout, "\r\n", log.LstdFlags),
				logger.Config{
					LogLevel:      logger.Error,
					SlowThreshold: 200 * time.Millisecond,
					Colorful:      true,
				},
			),
		},
	}
}

// TenantDBManager 租户数据库管理器
type TenantDBManager struct {
	dbs           map[string]*gorm.DB
	mutex         sync.RWMutex
	config        *DBConfig
	defaultConfig *gorm.Config
}

// NewTenantDBManager 创建租户数据库管理器
func NewTenantDBManager(config *DBConfig) *TenantDBManager {
	// 如果传入nil配置，则使用默认配置
	if config == nil {
		config = NewDefaultDBConfig()
	} else {
		// 填充默认值
		defaultConfig := NewDefaultDBConfig()

		// 租户DSNs为nil时初始化
		if config.TenantDSNs == nil {
			config.TenantDSNs = make(map[string]string)
		}

		// 数据库配置为nil时使用默认值
		if config.DBConfig == nil {
			config.DBConfig = defaultConfig.DBConfig
		} else {
			// 配置日志
			if config.DBConfig.Logger == nil {
				logLevel := config.LogLevel
				if logLevel == 0 {
					logLevel = defaultConfig.LogLevel
				}

				slowThreshold := config.SlowThreshold
				if slowThreshold == 0 {
					slowThreshold = defaultConfig.SlowThreshold
				}

				config.DBConfig.Logger = logger.New(
					log.New(os.Stdout, "\r\n", log.LstdFlags),
					logger.Config{
						LogLevel:      logLevel,
						SlowThreshold: slowThreshold,
						Colorful:      true,
					},
				)
			}
		}
	}

	return &TenantDBManager{
		dbs:           make(map[string]*gorm.DB),
		config:        config,
		defaultConfig: config.DBConfig,
	}
}

// GetDB 获取指定租户的数据库连接
func (m *TenantDBManager) GetDB(tenantID string) (*gorm.DB, error) {
	if tenantID == "" {
		return nil, fmt.Errorf("tenant ID cannot be empty")
	}

	m.mutex.RLock()
	db, exists := m.dbs[tenantID]
	m.mutex.RUnlock()

	if exists {
		return db, nil
	}

	// 如果不存在，创建新连接
	return m.createDB(tenantID)
}

// GetDefaultDB 获取默认数据库连接
func (m *TenantDBManager) GetDefaultDB() (*gorm.DB, error) {
	// 默认库使用空字符串作为key
	return m.GetDB("")
}

// createDB 创建租户数据库连接
func (m *TenantDBManager) createDB(tenantID string) (*gorm.DB, error) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	// 再次检查，防止并发创建
	if db, exists := m.dbs[tenantID]; exists {
		return db, nil
	}

	// 确定DSN
	var dsn string

	if specificDSN, ok := m.config.TenantDSNs[tenantID]; ok {
		// 使用预配置的租户特定DSN
		dsn = specificDSN
	} else if m.config.DSNTemplate != "" {
		// 使用模板构建租户特定DSN
		dsn = fmt.Sprintf(m.config.DSNTemplate, tenantID)
	} else {
		return nil, fmt.Errorf("no DSN configuration found for tenant: %s", tenantID)
	}

	// 创建新的DB连接
	db, err := gorm.Open(mysql.Open(dsn), m.defaultConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database for tenant %s: %w", tenantID, err)
	}

	// 添加OpenTelemetry
	if m.config.EnableTracing {
		if err := db.Use(tracing.NewPlugin(tracing.WithoutMetrics())); err != nil {
			return nil, fmt.Errorf("failed to setup tracing for tenant %s: %w", tenantID, err)
		}
	}

	// 存储连接
	m.dbs[tenantID] = db

	// 如果是租户库，执行迁移
	if tenantID != "" && m.config.MigrateFunc != nil {
		if err := m.config.MigrateFunc(db); err != nil {
			return nil, fmt.Errorf("failed to migrate database for tenant %s: %w", tenantID, err)
		}
	}

	return db, nil
}

// RegisterTenant 注册特定租户的DSN
func (m *TenantDBManager) RegisterTenant(tenantID string, dsn string) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	m.config.TenantDSNs[tenantID] = dsn
}

// SetMigrateFunc 设置迁移函数
func (m *TenantDBManager) SetMigrateFunc(migrateFunc MigrateFunc) {
	m.config.MigrateFunc = migrateFunc
}

// CloseAll 关闭所有数据库连接
func (m *TenantDBManager) CloseAll() {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	for tenantID, db := range m.dbs {
		sqlDB, err := db.DB()
		if err != nil {
			continue
		}
		_ = sqlDB.Close()
		delete(m.dbs, tenantID)
	}
}
