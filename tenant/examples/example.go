package main

import (
	"fmt"
	"time"

	"github.com/onebids/onecommon/tenant"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// 示例模型
type DBUser struct {
	ID   uint   `gorm:"primarykey"`
	Name string `gorm:"size:100"`
}

func main() {
	// 示例1：使用默认配置
	// 获取默认配置
	config := tenant.NewDefaultConfig()
	// 只需要设置必须的DSN模板
	config.DSNTemplate = "user:password@tcp(localhost:3306)/%s?charset=utf8mb4&parseTime=True&loc=Local"
	// 设置迁移函数
	config.MigrateFunc = func(db *gorm.DB) error {
		return db.AutoMigrate(&DBUser{})
	}

	// 创建数据库管理器
	dbManager1 := tenant.NewTenantDBManager(config)
	fmt.Println("成功创建数据库管理器 (使用默认配置)")

	// 示例2：完全自定义配置
	customConfig := &tenant.Config{
		DSNTemplate:   "user:password@tcp(localhost:3306)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		EnableTracing: true,
		LogLevel:      logger.Info, // 设置为Info级别以查看更多日志
		SlowThreshold: 100 * time.Millisecond,
		MigrateFunc: func(db *gorm.DB) error {
			return db.AutoMigrate(&User{})
		},
	}

	// 创建数据库管理器
	dbManager2 := tenant.NewTenantDBManager(customConfig)
	fmt.Println("成功创建数据库管理器 (使用自定义配置)")

	// 示例3：使用nil配置（将使用所有默认值）
	// 注意：这种情况下您需要在获取数据库连接前设置DSNTemplate
	dbManager3 := tenant.NewTenantDBManager(nil)
	// 设置DSN模板
	dbManager3.RegisterTenant("", "user:password@tcp(localhost:3306)/default_db?charset=utf8mb4&parseTime=True&loc=Local")
	fmt.Println("成功创建数据库管理器 (使用nil配置)")

	// 关闭所有连接
	defer func() {
		dbManager1.CloseAll()
		dbManager2.CloseAll()
		dbManager3.CloseAll()
	}()
}
