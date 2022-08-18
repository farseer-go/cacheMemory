package memoryCache

import (
	"github.com/farseer-go/collections"
	"time"
)

// 二级缓存-本地缓存操作
type cacheInMemory struct {
}

// localCache 缓存的存储
var localCache map[string]cacheValue

type cacheValue struct {
	// 缓存的数据
	data any
	// 缓存失效时间
	ttl int64
}

func newCacheInMemory() ICache {
	return cacheInMemory{}
}

func (r cacheInMemory) Get(cacheKey CacheKey) collections.ListAny {
	var defValue collections.ListAny
	value, ok := localCache[cacheKey.Key]
	if !ok {
		return defValue
	}
	return value.data.(collections.ListAny)
}

func (r cacheInMemory) GetItem(cacheKey CacheKey, cacheId string) any {
	lst := r.Get(cacheKey)
	for _, item := range lst.ToArray() {
		id := cacheKey.GetUniqueId(item)
		if cacheId == id {
			return item
		}
	}
	return nil
}

func (r cacheInMemory) Set(cacheKey CacheKey, val collections.ListAny) {
	localCache[cacheKey.Key] = cacheValue{
		data: val,
	}

	if cacheKey.MemoryExpiry > 0 {
		value, ok := localCache[cacheKey.Key]
		if !ok {
			return
		}
		value.ttl = time.Now().Add(cacheKey.MemoryExpiry).UnixMicro()
		localCache[cacheKey.Key] = value
	}
}

func (r cacheInMemory) SaveItem(cacheKey CacheKey, newVal any) {
	var list = r.Get(cacheKey)
	if list.Count() == 0 {
		return
	}

	// CacheKey.DataKey=null，说明实际缓存的是单个对象。所以此处直接替换新的对象即可，而不用查找。
	if cacheKey.UniqueField == "" {
		list.Clear()
	} else {
		// 从新对象中，获取唯一标识
		newValDataKey := cacheKey.GetUniqueId(newVal)
		for index := 0; index < list.Count(); index++ {
			// 从当前缓存item中，获取唯一标识
			itemDataKey := cacheKey.GetUniqueId(list.Index(index))
			// 找到了
			if itemDataKey == newValDataKey {
				list.Set(index, newVal)
				return
			}
		}
	}
	list.Add(newVal)
	// 保存
	r.Set(cacheKey, list)
}

func (r cacheInMemory) Remove(cacheKey CacheKey, cacheId string) {
	var list = r.Get(cacheKey)
	if list.Count() == 0 {
		return
	}
	list.RemoveAll(func(item any) bool { return cacheKey.GetUniqueId(item) == cacheId })
	// 保存
	r.Set(cacheKey, list)
}

func (r cacheInMemory) Clear(cacheKey CacheKey) {
	var list = r.Get(cacheKey)
	if list.Count() > 0 {
		list.Clear()
		r.Set(cacheKey, list)
	}
	delete(localCache, cacheKey.Key)
}

func (r cacheInMemory) Count(cacheKey CacheKey) int {
	return r.Get(cacheKey).Count()
}

func (r cacheInMemory) ExistsItem(cacheKey CacheKey, cacheId string) bool {
	var list = r.Get(cacheKey)
	if list.Count() == 0 {
		return false
	}
	for index := 0; index < list.Count(); index++ {
		// 从当前缓存item中，获取唯一标识
		itemDataKey := cacheKey.GetUniqueId(list.Index(index))
		// 找到了
		if itemDataKey == cacheId {
			return true
		}
	}
	return false
}

func (r cacheInMemory) ExistsKey(cacheKey CacheKey) bool {
	_, ok := localCache[cacheKey.Key]
	return ok
}
