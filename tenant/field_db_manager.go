package tenant

import (
	"context"
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"gorm.io/gorm/logger"
	"gorm.io/plugin/opentelemetry/tracing"
)

// FieldDBConfig 字段级隔离数据库管理器配置
type FieldDBConfig struct {
	// DSN连接字符串
	// 默认值: 空字符串，必须由用户提供
	DSN string

	// 数据库默认配置
	// 默认值: 见NewDefaultFieldDBConfig函数
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

	// 租户ID字段名
	// 默认值: "tenant_id"
	TenantIDField string
}

// NewDefaultFieldDBConfig 创建带有默认值的字段级隔离配置
func NewDefaultFieldDBConfig() *FieldDBConfig {
	return &FieldDBConfig{
		EnableTracing: false,
		LogLevel:      logger.Error,
		SlowThreshold: 200 * time.Millisecond,
		TenantIDField: "tenant_id",
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

// FieldDBManager 字段级租户隔离数据库管理器
type FieldDBManager struct {
	db            *gorm.DB
	mutex         sync.RWMutex
	config        *FieldDBConfig
	defaultConfig *gorm.Config
}

// NewFieldDBManager 创建字段级租户隔离数据库管理器
func NewFieldDBManager(config *FieldDBConfig) *FieldDBManager {
	// 如果传入nil配置，则使用默认配置
	if config == nil {
		config = NewDefaultFieldDBConfig()
	} else {
		// 填充默认值
		defaultConfig := NewDefaultFieldDBConfig()

		// 租户ID字段为空时使用默认值
		if config.TenantIDField == "" {
			config.TenantIDField = defaultConfig.TenantIDField
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

	return &FieldDBManager{
		config:        config,
		defaultConfig: config.DBConfig,
	}
}

// Connect 连接到数据库
func (m *FieldDBManager) Connect() error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	if m.db != nil {
		return nil // 已经连接
	}

	if m.config.DSN == "" {
		return fmt.Errorf("DSN cannot be empty")
	}

	// 创建新的DB连接
	db, err := gorm.Open(mysql.Open(m.config.DSN), m.defaultConfig)
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}

	// 添加OpenTelemetry
	if m.config.EnableTracing {
		if err := db.Use(tracing.NewPlugin(tracing.WithoutMetrics())); err != nil {
			return fmt.Errorf("failed to setup tracing: %w", err)
		}
	}

	// 存储连接
	m.db = db

	// 执行迁移
	if m.config.MigrateFunc != nil {
		if err := m.config.MigrateFunc(db); err != nil {
			return fmt.Errorf("failed to migrate database: %w", err)
		}
	}

	return nil
}

// tenantContext 用于存储租户ID的上下文键
type tenantContextKey struct{}

// WithTenant 创建包含租户ID的上下文
func WithTenant(ctx context.Context, tenantID string) context.Context {
	return context.WithValue(ctx, tenantContextKey{}, tenantID)
}

// GetTenantFromContext 从上下文中获取租户ID
func GetTenantFromContext(ctx context.Context) (string, bool) {
	tenantID, ok := ctx.Value(tenantContextKey{}).(string)
	return tenantID, ok
}

// GetDB 获取带有租户过滤的数据库连接
func (m *FieldDBManager) GetDB(ctx context.Context) (*gorm.DB, error) {
	m.mutex.RLock()
	db := m.db
	m.mutex.RUnlock()

	if db == nil {
		return nil, fmt.Errorf("database not connected, call Connect() first")
	}

	// 从上下文中获取租户ID
	tenantID, ok := GetTenantFromContext(ctx)
	if !ok || tenantID == "" {
		fmt.Println("tenantID为空")
		return db, nil // 没有租户ID，返回原始DB
	}
	fmt.Println("tenantID:", tenantID)
	// 创建带有租户过滤的新会话
	tenantDB := db.Session(&gorm.Session{NewDB: true})

	// 添加全局租户过滤回调
	tenantDB.Callback().Query().Before("gorm:query").Register("tenant_filter", func(db *gorm.DB) {
		fmt.Println("执行租户过滤回调")
		// 检查模型是否有租户ID字段
		if db.Statement.Schema != nil {
			fmt.Println("租户ID字段存在")
			if _, ok := db.Statement.Schema.FieldsByDBName[m.config.TenantIDField]; ok {
				// 添加租户过滤条件
				db.Statement.AddClause(clause.Where{
					Exprs: []clause.Expression{
						clause.Eq{
							Column: clause.Column{Table: clause.CurrentTable, Name: m.config.TenantIDField},
							Value:  tenantID,
						},
					},
				})
			}
		}
	})

	// 添加创建回调以自动设置租户ID
	tenantDB.Callback().Create().Before("gorm:create").Register("tenant_create", func(db *gorm.DB) {
		// 检查模型是否有租户ID字段
		if db.Statement.Schema != nil {
			if field, ok := db.Statement.Schema.FieldsByDBName[m.config.TenantIDField]; ok {
				field.Set(db.Statement.Context, db.Statement.ReflectValue, tenantID)

			}
		}
	})

	return tenantDB, nil
}

// SetMigrateFunc 设置迁移函数
func (m *FieldDBManager) SetMigrateFunc(migrateFunc MigrateFunc) {
	m.config.MigrateFunc = migrateFunc
}

// Close 关闭数据库连接
func (m *FieldDBManager) Close() error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	if m.db == nil {
		return nil
	}

	sqlDB, err := m.db.DB()
	if err != nil {
		return err
	}

	err = sqlDB.Close()
	if err == nil {
		m.db = nil
	}
	return err
}
