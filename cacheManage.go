package memoryCache

import (
	"github.com/farseer-go/collections"
	"github.com/farseer-go/fs/parse"
	"github.com/farseer-go/mapper"
	"github.com/farseer-go/memoryCache/eumCacheStoreType"
	"reflect"
	"strings"
	"time"
)

// value=CacheManage
var cacheConfigure map[string]any

type CacheKey struct {
	// 缓存KEY
	Key string
	// 缓存策略（默认Memory模式）
	CacheStoreType eumCacheStoreType.Enum
	// 设置Redis缓存过期时间
	RedisExpiry time.Duration
	// 设置Memory缓存过期时间
	MemoryExpiry time.Duration
	// hash中的主键（唯一ID的字段名称）
	UniqueField string
	// Redis配置名称
	RedisConfigName string
	//// 缓存的数据Item Type
	//ListType reflect.Type
	//// 缓存的数据Item Type
	//ItemType reflect.Type
	// 获取缓存实现
	Cache ICache
}

// getUniqueId 获取唯一字段数据
func (receiver CacheKey) getUniqueId(item any) (T string) {
	val := reflect.ValueOf(item).FieldByName(receiver.UniqueField).Interface()
	return parse.Convert(val, "")
}

type CacheManage[TEntity any] struct {
	CacheKey
}

// GetCacheManage 获取CacheKey
func GetCacheManage[TEntity any](key string) CacheManage[TEntity] {
	cacheKey, exists := cacheConfigure[key]
	if !exists {
		panic(key + "不存在，要使用Cache缓存，需提前初始化")
	}
	return cacheKey.(CacheManage[TEntity])
}

// Get 获取缓存数据
func (receiver CacheManage[TEntity]) Get() collections.List[TEntity] {
	lst := receiver.Cache.Get(receiver.CacheKey)
	return mapper.ToList[TEntity](lst)
}

// GetItem 从集合中获取指定cacheId的元素
func (receiver CacheManage[TEntity]) GetItem(cacheId any) (TEntity, bool) {
	item := receiver.Cache.GetItem(receiver.CacheKey, parse.Convert(cacheId, ""))
	if item == nil {
		var entity TEntity
		return entity, false
	}
	return item.(TEntity), true
}

// Set 保存缓存
func (receiver CacheManage[TEntity]) Set(val any) {
	valValue := reflect.ValueOf(val)
	valType := valValue.Type()
	lst := collections.NewListAny()

	// 数据是数组切片时
	if valType.Kind() == reflect.Slice || valType.Kind() == reflect.Array {
		for i := 0; i < valValue.Len(); i++ {
			itemValue := valValue.Index(i)
			lst.Add(itemValue.Interface())
		}
	} else if strings.HasPrefix(valType.String(), "collections.List[") { // 数据是List时
		arrValue := valValue.MethodByName("ToArray").Call(nil)[0]
		for i := 0; i < arrValue.Len(); i++ {
			itemValue := arrValue.Index(i)
			lst.Add(itemValue.Interface())
		}
	} else { // 数据作为单个item保存
		lst.Add(val)
	}
	receiver.Cache.Set(receiver.CacheKey, lst)
}

// SaveItem 更新缓存
func (receiver CacheManage[TEntity]) SaveItem(newVal TEntity) {
	receiver.Cache.SaveItem(receiver.CacheKey, newVal)
}

// Remove 移除缓存
func (receiver CacheManage[TEntity]) Remove(cacheId string) {
	receiver.Cache.Remove(receiver.CacheKey, cacheId)
}

// Clear 清空缓存
func (receiver CacheManage[TEntity]) Clear() {
	receiver.Cache.Clear(receiver.CacheKey)
}

// Exists 缓存是否存在
func (receiver CacheManage[TEntity]) Exists() bool {
	return receiver.Cache.ExistsKey(receiver.CacheKey)
}

// Count 数据集合的数量
func (receiver CacheManage[TEntity]) Count() int {
	if !receiver.Exists() {
		return 0
	}
	return receiver.Cache.Count(receiver.CacheKey)
}
