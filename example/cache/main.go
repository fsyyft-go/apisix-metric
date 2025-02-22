// Copyright 2024 fsyyft-go
//
// Licensed under the MIT License. See LICENSE file in the project root for full license information.

// Package main 提供了缓存包的使用示例，展示了缓存包的各种功能和最佳实践。
// 包含以下示例：
//   - 基本的缓存操作（设置、获取、删除）
//   - 带过期时间的缓存操作（TTL 支持）
//   - 类型安全的缓存操作（泛型支持）
//   - 全局缓存的使用（单例模式）
//   - 结构体类型的缓存（复杂类型支持）
package main

import (
	"fmt"
	"time"

	"github.com/fsyyft-go/apisix-metric/pkg/cache"
)

// User 表示用户信息结构体，用于演示复杂类型的缓存操作。
// 包含用户的基本信息和最后在线时间，可用于用户会话管理。
type User struct {
	// ID 表示用户的唯一标识符。
	ID int
	// Name 表示用户的名称。
	Name string
	// Email 表示用户的电子邮件地址。
	Email string
	// LastSeen 表示用户最后一次在线的时间。
	LastSeen time.Time
}

func main() {
	// 示例 1：展示基本的缓存操作，包括创建、设置、获取和删除缓存。
	fmt.Println("=== 基本的缓存操作 ===")
	basicCacheExample()

	// 示例 2：展示带过期时间的缓存操作，演示 TTL 功能。
	fmt.Println("\n=== 带过期时间的缓存操作 ===")
	ttlCacheExample()

	// 示例 3：展示类型安全的缓存操作，演示泛型功能。
	fmt.Println("\n=== 类型安全的缓存操作 ===")
	typedCacheExample()

	// 示例 4：展示全局缓存的使用，演示单例模式。
	fmt.Println("\n=== 全局缓存的使用 ===")
	globalCacheExample()

	// 示例 5：展示结构体类型的缓存，演示复杂类型支持。
	fmt.Println("\n=== 结构体类型的缓存 ===")
	structCacheExample()
}

// basicCacheExample 展示基本的缓存操作。
// 演示了如何创建缓存实例、设置和获取缓存值，以及删除缓存项。
// 这些是最常用的缓存操作，适用于简单的缓存场景。
func basicCacheExample() {
	// 使用默认配置创建一个新的缓存实例。
	c, err := cache.NewCache(cache.DefaultConfig())
	if err != nil {
		fmt.Printf("创建缓存失败：%v。\n", err)
		return
	}
	defer c.Close()

	// 设置一个简单的字符串值到缓存中。
	key := "greeting"
	value := "Hello, World!"
	if !c.Set(key, value) {
		fmt.Println("设置缓存失败。")
		return
	}

	// 从缓存中获取之前设置的值。
	if val, exists := c.Get(key); exists {
		fmt.Printf("获取到缓存值：%v。\n", val)
	} else {
		fmt.Println("缓存值不存在。")
	}

	// 从缓存中删除该值。
	c.Delete(key)
	if _, exists := c.Get(key); !exists {
		fmt.Println("缓存已被删除。")
	}
}

// ttlCacheExample 展示带过期时间的缓存操作。
// 演示了如何设置带 TTL 的缓存值，以及如何检查剩余过期时间。
// TTL（生存时间）功能常用于临时数据存储，如会话管理。
func ttlCacheExample() {
	// 创建缓存实例。
	c, err := cache.NewCache(cache.DefaultConfig())
	if err != nil {
		fmt.Printf("创建缓存失败：%v。\n", err)
		return
	}
	defer c.Close()

	// 设置一个 2 秒后过期的缓存项。
	key := "temp"
	value := "临时数据"
	ttl := 2 * time.Second

	if !c.SetWithTTL(key, value, ttl) {
		fmt.Println("设置缓存失败。")
		return
	}

	// 立即获取值和剩余过期时间。
	if val, exists, remainingTTL := c.GetWithTTL(key); exists {
		fmt.Printf("值：%v，剩余时间：%v。\n", val, remainingTTL)
	}

	// 等待 1 秒后再次获取，演示剩余时间的变化。
	time.Sleep(time.Second)
	if val, exists, remainingTTL := c.GetWithTTL(key); exists {
		fmt.Printf("1 秒后 - 值：%v，剩余时间：%v。\n", val, remainingTTL)
	}

	// 等待直到过期（额外等待 100ms 确保过期）。
	time.Sleep(1100 * time.Millisecond)
	if _, exists := c.Get(key); !exists {
		fmt.Println("缓存已过期。")
	}
}

// typedCacheExample 展示类型安全的缓存操作。
// 演示了如何使用泛型包装器创建类型安全的缓存，以及类型不匹配的处理。
// 类型安全的缓存可以在编译时捕获类型错误，提供更好的类型检查。
func typedCacheExample() {
	// 创建基础缓存实例。
	baseCache, err := cache.NewCache(cache.DefaultConfig())
	if err != nil {
		fmt.Printf("创建缓存失败：%v。\n", err)
		return
	}
	defer baseCache.Close()

	// 创建字符串类型和整数类型的类型安全缓存。
	strCache := cache.NewTypedCache[string](baseCache)
	intCache := cache.NewTypedCache[int](baseCache)

	// 使用字符串类型缓存，无需类型断言。
	strCache.Set("name", "Alice")
	if name, exists := strCache.Get("name"); exists {
		fmt.Printf("名字：%s，长度：%d。\n", name, len(name))
	}

	// 使用整数类型缓存，无需类型断言。
	intCache.Set("age", 25)
	if age, exists := intCache.Get("age"); exists {
		fmt.Printf("年龄：%d。\n", age)
	}

	// 演示类型不匹配的情况。
	baseCache.Set("mixed", 42)
	if _, exists := strCache.Get("mixed"); !exists {
		fmt.Println("无法将整数值作为字符串获取。")
	}
}

// globalCacheExample 展示全局缓存的使用。
// 演示了如何使用全局缓存实例，适用于需要在多个包之间共享缓存的场景。
// 全局缓存使用单例模式，确保只有一个缓存实例。
func globalCacheExample() {
	// 使用默认配置初始化全局缓存。
	if err := cache.InitCache(cache.DefaultConfig()); err != nil {
		fmt.Printf("初始化全局缓存失败：%v。\n", err)
		return
	}
	defer cache.Close()

	// 使用全局缓存存储数据。
	cache.Set("global_key", "全局数据")
	if val, exists := cache.Get("global_key"); exists {
		fmt.Printf("从全局缓存获取：%v。\n", val)
	}

	// 使用全局缓存存储临时数据。
	cache.SetWithTTL("temp_global", "临时全局数据", time.Second)
	if val, exists, ttl := cache.GetWithTTL("temp_global"); exists {
		fmt.Printf("临时数据：%v，剩余时间：%v。\n", val, ttl)
	}
}

// structCacheExample 展示结构体类型的缓存。
// 演示了如何缓存复杂的结构体类型，以及如何使用带过期时间的结构体缓存。
// 这种方式适用于缓存用户会话、配置信息等复杂数据。
func structCacheExample() {
	// 创建基础缓存实例。
	baseCache, err := cache.NewCache(cache.DefaultConfig())
	if err != nil {
		fmt.Printf("创建缓存失败：%v。\n", err)
		return
	}
	defer baseCache.Close()

	// 创建用于存储 User 结构体的类型安全缓存。
	userCache := cache.NewTypedCache[User](baseCache)

	// 创建示例用户数据。
	user := User{
		ID:       1,
		Name:     "Alice",
		Email:    "alice@example.com",
		LastSeen: time.Now(),
	}

	// 将用户数据存储到缓存中。
	userCache.Set("user:1", user)

	// 获取并显示缓存的用户数据。
	if cached, exists := userCache.Get("user:1"); exists {
		fmt.Printf("用户信息：\n")
		fmt.Printf("  ID：%d\n", cached.ID)
		fmt.Printf("  名字：%s\n", cached.Name)
		fmt.Printf("  邮箱：%s\n", cached.Email)
		fmt.Printf("  最后在线：%v\n", cached.LastSeen)
	}

	// 创建用于存储用户会话的缓存，并设置 30 分钟的过期时间。
	sessionCache := cache.NewTypedCache[User](baseCache)
	sessionCache.SetWithTTL("session:1", user, 30*time.Minute)

	// 获取并显示会话信息及其过期时间。
	if cached, exists, ttl := sessionCache.GetWithTTL("session:1"); exists {
		fmt.Printf("用户 %s 的会话将在 %v 后过期。\n", cached.Name, ttl)
	}
}
