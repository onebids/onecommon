package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/onebids/onecommon/tenant"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// 示例模型 - 包含租户ID字段
type TESUser struct {
	ID       uint   `gorm:"primarykey"`
	Name     string `gorm:"size:100"`
	TenantID string `gorm:"size:50;index"` // 租户ID字段，会被自动设置
}

func main() {
	// 示例1：使用默认配置
	config := tenant.NewDefaultFieldDBConfig()
	// 设置DSN
	config.DSN = "root:123456@tcp(10.3.5.103:3306)/sharedb?charset=utf8mb4&parseTime=True&loc=Local"
	// 设置迁移函数
	config.MigrateFunc = func(db *gorm.DB) error {
		return db.AutoMigrate(&TESUser{})
	}

	// 创建字段级租户隔离数据库管理器
	dbManager := tenant.NewFieldDBManager(config)
	fmt.Println("成功创建字段级租户隔离数据库管理器")

	// 连接到数据库
	if err := dbManager.Connect(); err != nil {
		log.Fatalf("连接数据库失败: %v", err)
	}

	// 关闭连接
	defer dbManager.Close()

	// 示例2：为不同租户创建数据
	// 创建租户A的上下文
	ctxA := tenant.WithTenant(context.Background(), "tenant-a")
	// 获取租户A的数据库连接
	dbA, err := dbManager.GetDB(ctxA)
	if err != nil {
		log.Fatalf("获取租户A的数据库连接失败: %v", err)
	}

	// 为租户A创建用户
	userA := TESUser{Name: "Alice"}
	if err := dbA.Create(&userA).Error; err != nil {
		log.Fatalf("创建租户A的用户失败: %v", err)
	}
	fmt.Printf("为租户A创建用户: %+v\n", userA)

	// 创建租户B的上下文
	ctxB := tenant.WithTenant(context.Background(), "tenant-b")
	// 获取租户B的数据库连接
	dbB, err := dbManager.GetDB(ctxB)
	if err != nil {
		log.Fatalf("获取租户B的数据库连接失败: %v", err)
	}

	// 为租户B创建用户
	userB := TESUser{Name: "Bob"}
	if err := dbB.Create(&userB).Error; err != nil {
		log.Fatalf("创建租户B的用户失败: %v", err)
	}
	fmt.Printf("为租户B创建用户: %+v\n", userB)

	// 示例3：查询租户数据
	// 查询租户A的所有用户
	var usersA []TESUser
	if err := dbA.Find(&usersA).Error; err != nil {
		log.Fatalf("查询租户A的用户失败: %v", err)
	}
	fmt.Printf("租户A的用户: %+v\n", usersA)

	// 查询租户B的所有用户
	var usersB []TESUser
	if err := dbB.Find(&usersB).Error; err != nil {
		log.Fatalf("查询租户B的用户失败: %v", err)
	}
	fmt.Printf("租户B的用户: %+v\n", usersB)

	// 示例4：使用自定义配置
	customConfig := &tenant.FieldDBConfig{
		DSN:           "user:password@tcp(localhost:3306)/shared_db?charset=utf8mb4&parseTime=True&loc=Local",
		EnableTracing: true,
		LogLevel:      logger.Info,
		SlowThreshold: 100 * time.Millisecond,
		TenantIDField: "organization_id", // 使用自定义的租户ID字段名
	}

	// 创建自定义字段级租户隔离数据库管理器
	_ = tenant.NewFieldDBManager(customConfig)
	fmt.Println("成功创建自定义字段级租户隔离数据库管理器")
}
