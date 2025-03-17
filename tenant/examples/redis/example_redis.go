package main

import (
	"context"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/onebids/onecommon/consts"
	"github.com/onebids/onecommon/tenant"
)

type User struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
}

func main() {
	// 创建Redis管理器配置
	config := &tenant.RedisConfig{
		DefaultOptions: &redis.Options{
			Addr:     "localhost:6379",
			Password: "", // 无密码
			DB:       0,  // 默认数据库
		},
		TenantOptions: map[string]*redis.Options{
			"tenant1": {
				Addr:     "localhost:6379",
				Password: "",
				DB:       1, // 使用不同的数据库隔离租户
			},
			"tenant2": {
				Addr:     "localhost:6379",
				Password: "",
				DB:       2,
			},
		},
		EnableTenantIsolation: true,
		TenantSeparator:       ":",
	}

	// 创建Redis管理器
	manager, err := tenant.NewRedisManager(config)
	if err != nil {
		fmt.Printf("创建Redis管理器失败: %v\n", err)
		return
	}
	defer manager.Close()

	// 创建Redis辅助工具
	helper := tenant.NewRedisHelper(manager)

	// 创建默认上下文
	defaultCtx := context.Background()

	// 创建带有租户信息的上下文
	tenant1Ctx := context.WithValue(defaultCtx, consts.TenantID, "tenant1")
	tenant2Ctx := context.WithValue(defaultCtx, consts.TenantID, "tenant2")

	// 示例1: 使用默认上下文（无租户）
	fmt.Println("示例1: 使用默认上下文")
	if err := helper.Set(defaultCtx, "test_key", "default_value", 10*time.Minute); err != nil {
		fmt.Printf("设置键值对失败: %v\n", err)
	}

	value, err := helper.Get(defaultCtx, "test_key")
	if err != nil {
		fmt.Printf("获取值失败: %v\n", err)
	} else {
		fmt.Printf("获取到的值: %s\n", value)
	}
	fmt.Println()

	// 示例2: 使用租户1上下文
	fmt.Println("示例2: 使用租户1上下文")
	if err := helper.Set(tenant1Ctx, "test_key", "tenant1_value", 10*time.Minute); err != nil {
		fmt.Printf("设置键值对失败: %v\n", err)
	}

	value, err = helper.Get(tenant1Ctx, "test_key")
	if err != nil {
		fmt.Printf("获取值失败: %v\n", err)
	} else {
		fmt.Printf("获取到的值: %s\n", value)
	}
	fmt.Println()

	// 示例3: 使用租户2上下文
	fmt.Println("示例3: 使用租户2上下文")
	if err := helper.Set(tenant2Ctx, "test_key", "tenant2_value", 10*time.Minute); err != nil {
		fmt.Printf("设置键值对失败: %v\n", err)
	}

	value, err = helper.Get(tenant2Ctx, "test_key")
	if err != nil {
		fmt.Printf("获取值失败: %v\n", err)
	} else {
		fmt.Printf("获取到的值: %s\n", value)
	}
	fmt.Println()

	// 示例4: 存储和检索对象
	fmt.Println("示例4: 存储和检索对象")
	user := User{
		ID:       1,
		Username: "johndoe",
		Email:    "john.doe@example.com",
	}

	if err := helper.Set(tenant1Ctx, "user:1", user, 10*time.Minute); err != nil {
		fmt.Printf("存储对象失败: %v\n", err)
	}

	var retrievedUser User
	if err := helper.GetObject(tenant1Ctx, "user:1", &retrievedUser); err != nil {
		fmt.Printf("检索对象失败: %v\n", err)
	} else {
		fmt.Printf("检索到的对象: %+v\n", retrievedUser)
	}
	fmt.Println()

	// 示例5: 哈希表操作
	fmt.Println("示例5: 哈希表操作")
	if err := helper.HSet(tenant1Ctx, "user:hash:1", "username", "johndoe"); err != nil {
		fmt.Printf("设置哈希表字段失败: %v\n", err)
	}
	if err := helper.HSet(tenant1Ctx, "user:hash:1", "email", "john.doe@example.com"); err != nil {
		fmt.Printf("设置哈希表字段失败: %v\n", err)
	}

	fields, err := helper.HGetAll(tenant1Ctx, "user:hash:1")
	if err != nil {
		fmt.Printf("获取哈希表所有字段失败: %v\n", err)
	} else {
		fmt.Printf("哈希表字段: %+v\n", fields)
	}
	fmt.Println()

	// 示例6: 列表操作
	fmt.Println("示例6: 列表操作")
	if err := helper.RPush(tenant1Ctx, "user:list", "user1", "user2", "user3"); err != nil {
		fmt.Printf("推入列表失败: %v\n", err)
	}

	items, err := helper.LRange(tenant1Ctx, "user:list", 0, -1)
	if err != nil {
		fmt.Printf("获取列表范围失败: %v\n", err)
	} else {
		fmt.Printf("列表项: %v\n", items)
	}
	fmt.Println()

	// 示例7: 分布式锁
	fmt.Println("示例7: 分布式锁")
	lockKey := "my_lock"
	lockValue := "process_1"

	acquired, err := helper.Lock(tenant1Ctx, lockKey, lockValue, 10*time.Second)
	if err != nil {
		fmt.Printf("获取锁失败: %v\n", err)
	} else {
		fmt.Printf("获取锁: %v\n", acquired)

		// 执行需要锁保护的操作
		fmt.Println("执行需要锁保护的操作...")

		// 释放锁
		released, err := helper.Unlock(tenant1Ctx, lockKey, lockValue)
		if err != nil {
			fmt.Printf("释放锁失败: %v\n", err)
		} else {
			fmt.Printf("释放锁: %v\n", released)
		}
	}
}
