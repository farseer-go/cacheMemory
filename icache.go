package memoryCache

import (
	"github.com/farseer-go/collections"
)

// ICache 缓存操作
type ICache interface {
	// Get 从本地内存中获取
	Get(key CacheKey) collections.ListAny
	// GetItem 获取集合中的item
	GetItem(key CacheKey, cacheId string) any
	// Set 整个数据保存到缓存
	Set(key CacheKey, val collections.ListAny)
	// SaveItem 添加或修改集合的Item
	SaveItem(key CacheKey, newVal any)
	// Remove 移除缓存
	Remove(key CacheKey, cacheId string)
	// Clear 清空缓存
	Clear(key CacheKey)
	// Count 获取集合数量
	Count(key CacheKey) int
	// ExistsItem 指定数据是否存在
	ExistsItem(key CacheKey, cacheId string) bool
	// ExistsKey 缓存是否存在
	ExistsKey(key CacheKey) bool
}
