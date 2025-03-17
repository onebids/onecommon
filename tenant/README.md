# 通用租户数据库管理器 (TenantDBManager)

这是一个基于GORM的通用租户数据库管理器，支持多租户场景，也可用于单数据库场景。

## 特点

- 不依赖于特定的模型和配置结构
- 支持多租户数据库管理
- 提供灵活的配置选项
- 支持自定义迁移功能
- 支持OpenTelemetry追踪（可选）
- 线程安全

## 使用方法

### 基本使用

```go
package main

import (
    "fmt"
    "github.com/onebids/onecommon/tenant"
    "gorm.io/gorm"
)

// 定义模型
type User struct {
    ID   uint
    Name string
}

func main() {
    // 创建配置
    config := &tenant.Config{
        DSNTemplate: "user:password@tcp(localhost:3306)/%s?charset=utf8mb4&parseTime=True&loc=Local",
        EnableTracing: true,
    }
    
    // 创建数据库管理器
    dbManager := tenant.NewTenantDBManager(config)
    
    // 设置迁移函数
    dbManager.SetMigrateFunc(func(db *gorm.DB) error {
        return db.AutoMigrate(&User{})
    })
    
    // 获取默认数据库连接
    defaultDB, err := dbManager.GetDefaultDB()
    if err != nil {
        panic(err)
    }
    
    // 获取特定租户的数据库连接
    tenantDB, err := dbManager.GetDB("tenant1")
    if err != nil {
        panic(err)
    }
    
    // 使用数据库连接
    defaultDB.Create(&User{Name: "DefaultUser"})
    tenantDB.Create(&User{Name: "TenantUser"})
    
    // 应用结束时关闭所有连接
    defer dbManager.CloseAll()
}
```

### 高级配置

```go
package main

import (
    "github.com/onebids/onecommon/tenant"
    "gorm.io/gorm"
    "gorm.io/gorm/logger"
    "time"
)

func main() {
    // 创建详细配置
    config := &tenant.Config{
        DSNTemplate: "user:password@tcp(localhost:3306)/%s?charset=utf8mb4&parseTime=True&loc=Local",
        TenantDSNs: map[string]string{
            "tenant1": "admin:secretpass@tcp(special-host:3306)/tenant1_db?charset=utf8mb4&parseTime=True&loc=Local",
        },
        EnableTracing: true,
        LogLevel: logger.Info,  // 设置为Info级别以查看更多日志
        SlowThreshold: 100 * time.Millisecond,  // 将慢查询阈值设置为100ms
    }
    
    dbManager := tenant.NewTenantDBManager(config)
    
    // 使用...
    
    // 动态注册新的租户
    dbManager.RegisterTenant("tenant2", "user:pass@tcp(localhost:3306)/tenant2_db?charset=utf8mb4&parseTime=True&loc=Local")
    
    // 获取新注册租户的连接
    tenant2DB, _ := dbManager.GetDB("tenant2")
    
    // 使用...
}
```

## 与原有代码的区别

相比原来的`TenantDBManager`，新版本有以下改进：

1. 不依赖于特定项目的配置结构（如`conf.GetConf().MySQL.DSN`）
2. 支持更灵活的配置，如日志级别、慢查询阈值等
3. 迁移功能通过回调函数实现，不依赖于特定的模型
4. 增加了`RegisterTenant`方法，支持运行时动态注册新租户
5. 支持可选的OpenTelemetry追踪
6. 提供了更通用的接口，可以用于任何基于GORM的项目
