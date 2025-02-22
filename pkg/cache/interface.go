// Copyright 2024 fsyyft-go
//
// Licensed under the MIT License. See LICENSE file in the project root for full license information.

// Package cache 提供了一个统一的缓存接口和多种缓存实现。
// 这个包的主要特性包括：
//   - 支持多种缓存后端（默认使用 ristretto）
//   - 提供统一的缓存接口
//   - 支持泛型和类型安全
//   - 支持 TTL（生存时间）设置
//   - 支持全局缓存实例
//   - 线程安全
//
// 基本使用示例：
//
//	// 创建缓存实例
//	cache, err := cache.NewCache(cache.DefaultConfig())
//	if err != nil {
//	    panic(err)
//	}
//	defer cache.Close()
//
//	// 基本操作
//	cache.Set("key", "value")                    // 设置永不过期的值
//	cache.SetWithTTL("temp", "value", time.Hour) // 设置 1 小时后过期的值
//
//	// 获取值
//	if val, exists := cache.Get("key"); exists {
//	    fmt.Println(val)
//	}
//
//	// 获取带 TTL 的值
//	if val, exists, ttl := cache.GetWithTTL("temp"); exists {
//	    fmt.Printf("值：%v，剩余时间：%v\n", val, ttl)
//	}
//
// 更多示例请参考 example/cache 目录。
package cache

import "time"

// Cache 定义了统一的缓存接口。
// 这个接口提供了基本的缓存操作功能，可以通过不同的实现来支持不同的缓存后端。
// 所有的实现都必须保证线程安全。
type Cache interface {
	// Get 获取缓存中的值。
	// 如果键不存在或已过期，exists 将返回 false。
	// 返回值：
	//   - value：缓存的值，如果不存在则为 nil
	//   - exists：值是否存在且未过期
	Get(key interface{}) (value interface{}, exists bool)

	// GetWithTTL 获取缓存中的值及其剩余过期时间。
	// 如果键不存在或已过期，exists 将返回 false。
	// 返回值：
	//   - value：缓存的值，如果不存在则为 nil
	//   - exists：值是否存在且未过期
	//   - remainingTTL：剩余过期时间：
	//     * 如果值不存在或已过期，返回 0
	//     * 如果值永不过期，返回 -1
	//     * 否则返回实际的剩余时间
	GetWithTTL(key interface{}) (value interface{}, exists bool, remainingTTL time.Duration)

	// Set 设置缓存值，该值永不过期。
	// 参数：
	//   - key：缓存键，可以是任意类型
	//   - value：要缓存的值，可以是任意类型
	// 返回值：
	//   - bool：是否设置成功
	Set(key interface{}, value interface{}) bool

	// SetWithTTL 设置带过期时间的缓存值。
	// 参数：
	//   - key：缓存键，可以是任意类型
	//   - value：要缓存的值，可以是任意类型
	//   - ttl：过期时间，如果 <= 0 则表示永不过期
	// 返回值：
	//   - bool：是否设置成功
	SetWithTTL(key interface{}, value interface{}, ttl time.Duration) bool

	// Delete 从缓存中删除指定的键。
	// 如果键不存在，该操作也会成功返回。
	// 删除后，该键的所有后续访问都将返回不存在。
	Delete(key interface{})

	// Clear 清空缓存中的所有内容。
	// 这个操作会立即使所有缓存项失效。
	// 在清空后，所有之前的键都将返回不存在。
	Clear()

	// Close 关闭缓存，释放相关资源。
	// 关闭后的缓存不应该再被使用。
	// 重复调用 Close 是安全的。
	Close() error
}

// Config 定义了缓存的配置选项。
// 这些配置项会影响缓存的性能和资源使用。
type Config struct {
	// NumCounters 定义了缓存跟踪的最大条目数。
	// 这个数字应该是预期独特条目数的大约 10 倍。
	// 例如，如果预计会有 1 万个不同的键，则应设置为 10 万。
	NumCounters int64

	// MaxCost 定义了缓存的最大成本。
	// 对于简单的缓存使用场景，这可以理解为最大条目数。
	// 当缓存达到这个限制时，最少使用的项目将被驱逐。
	MaxCost int64

	// BufferItems 定义了在写入操作时的缓冲大小。
	// 更大的缓冲区会导致更好的并发性能，但会使用更多的内存。
	// 对于大多数场景，默认值 64 是合适的。
	BufferItems int64
}

// DefaultConfig 返回默认的缓存配置。
// 默认配置适用于大多数中等规模的应用场景：
//   - NumCounters：1000 万（适合跟踪 100 万个不同的键）
//   - MaxCost：1GB（适合存储较大的数据集）
//   - BufferItems：64（提供良好的并发性能）
func DefaultConfig() Config {
	return Config{
		NumCounters: 1e7,     // 1000 万
		MaxCost:     1 << 30, // 1GB
		BufferItems: 64,      // 64 个条目的缓冲区
	}
}

// TypedCache 是一个泛型包装器，提供类型安全的缓存操作。
// 通过泛型参数 T 指定缓存值的类型，避免了手动类型断言。
// 这个包装器适用于需要类型安全的场景，例如：
//   - 存储特定类型的数据
//   - 避免运行时类型错误
//   - 提供更好的 IDE 支持
type TypedCache[T any] struct {
	cache Cache
}

// NewTypedCache 创建一个新的类型安全的缓存包装器。
// 参数：
//   - cache：底层的缓存实现
//
// 返回值：
//   - *TypedCache[T]：类型安全的缓存包装器
//
// 示例：
//
//	baseCache := cache.NewCache(cache.DefaultConfig())
//	strCache := cache.NewTypedCache[string](baseCache)
//	intCache := cache.NewTypedCache[int](baseCache)
func NewTypedCache[T any](cache Cache) *TypedCache[T] {
	return &TypedCache[T]{
		cache: cache,
	}
}

// Get 获取缓存中的值，并进行类型转换。
// 如果键不存在、已过期或类型不匹配，exists 将返回 false。
// 返回的值已经是正确的类型，无需进行类型断言。
func (tc *TypedCache[T]) Get(key interface{}) (value T, exists bool) {
	if v, ok := tc.cache.Get(key); ok {
		if typed, ok := v.(T); ok {
			return typed, true
		}
	}
	return value, false
}

// GetWithTTL 获取缓存中的值及其剩余过期时间，并进行类型转换。
// 如果键不存在、已过期或类型不匹配，exists 将返回 false。
// 返回的值已经是正确的类型，无需进行类型断言。
func (tc *TypedCache[T]) GetWithTTL(key interface{}) (value T, exists bool, remainingTTL time.Duration) {
	if v, ok, ttl := tc.cache.GetWithTTL(key); ok {
		if typed, ok := v.(T); ok {
			return typed, true, ttl
		}
	}
	return value, false, 0
}

// Set 设置缓存值，该值永不过期。
// 参数值必须匹配泛型类型 T，这在编译时就能保证。
func (tc *TypedCache[T]) Set(key interface{}, value T) bool {
	return tc.cache.Set(key, value)
}

// SetWithTTL 设置带过期时间的缓存值。
// 参数值必须匹配泛型类型 T，这在编译时就能保证。
func (tc *TypedCache[T]) SetWithTTL(key interface{}, value T, ttl time.Duration) bool {
	return tc.cache.SetWithTTL(key, value, ttl)
}

// Delete 从缓存中删除指定的键。
// 这个操作与底层缓存的 Delete 操作相同。
func (tc *TypedCache[T]) Delete(key interface{}) {
	tc.cache.Delete(key)
}

// Clear 清空缓存中的所有内容。
// 这个操作与底层缓存的 Clear 操作相同。
func (tc *TypedCache[T]) Clear() {
	tc.cache.Clear()
}

// Close 关闭缓存，释放相关资源。
// 这个操作与底层缓存的 Close 操作相同。
func (tc *TypedCache[T]) Close() error {
	return tc.cache.Close()
}
