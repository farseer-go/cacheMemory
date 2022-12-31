package cacheMemory

import (
	"github.com/farseer-go/cache"
	"github.com/farseer-go/collections"
	"github.com/farseer-go/fs/parse"
	"reflect"
	"sync"
	"time"
)

// localCache 缓存的存储
var localCache map[string]*cacheValue

// 缓存锁
var lock sync.RWMutex

type cacheValue struct {
	// 缓存的数据
	data any
	// 缓存失效时间
	ttl int64
	// 失效后，拿到chan通知
	ttlAfter <-chan time.Time
}

// 二级缓存-本地缓存操作
type cacheInMemory struct {
	expiry      time.Duration // 设置Memory缓存过期时间
	uniqueField string        // hash中的主键（唯一ID的字段名称）
	itemType    reflect.Type  // itemType
	key         string        // 缓存KEY
}

// 创建实例
func newCache(key string, uniqueField string, itemType reflect.Type, expiry time.Duration) cache.ICache {
	return &cacheInMemory{
		expiry:      expiry,
		uniqueField: uniqueField,
		itemType:    itemType,
		key:         key,
	}
}

func (r *cacheInMemory) Get() collections.ListAny {
	var defValue collections.ListAny
	value, ok := localCache[r.key]
	if !ok {
		return defValue
	}
	return value.data.(collections.ListAny)
}

func (r *cacheInMemory) GetItem(cacheId string) any {
	lst := r.Get()
	for _, item := range lst.ToArray() {
		id := r.GetUniqueId(item)
		if cacheId == id {
			return item
		}
	}
	return nil
}

func (r *cacheInMemory) Set(val collections.ListAny) {
	lock.Lock()
	defer lock.Unlock()
	localCache[r.key] = &cacheValue{
		data: val,
	}

	if r.expiry > 0 {
		value := localCache[r.key]
		value.ttl = time.Now().Add(r.expiry).UnixMilli()
		value.ttlAfter = time.After(r.expiry)

		// ttl到期后，自动删除缓存
		go r.ttl()()
	}
}

// ttl到期后，自动删除缓存
func (r *cacheInMemory) ttl() func() {
	return func() {
		<-localCache[r.key].ttlAfter
		lock.Lock()
		defer lock.Unlock()
		delete(localCache, r.key)
	}
}

func (r *cacheInMemory) SaveItem(newVal any) {
	var list = r.Get()
	// if list.Count() == 0 {
	//	return
	// }

	// cacheKey.DataKey=null，说明实际缓存的是单个对象。所以此处直接替换新的对象即可，而不用查找。
	// if cacheKey.uniqueField == "" {
	//	list.Clear()
	// } else {
	// 从新对象中，获取唯一标识
	newValDataKey := r.GetUniqueId(newVal)
	for index := 0; index < list.Count(); index++ {
		// 从当前缓存item中，获取唯一标识
		itemDataKey := r.GetUniqueId(list.Index(index))
		// 找到了
		if itemDataKey == newValDataKey {
			list.Set(index, newVal)
			return
		}
	}
	// }
	if list.Count() == 0 {
		list = collections.NewListAny()
	}
	list.Add(newVal)
	// 保存
	r.Set(list)
}

func (r *cacheInMemory) Remove(cacheId string) {
	var list = r.Get()
	if list.Count() > 0 {
		list.RemoveAll(func(item any) bool { return r.GetUniqueId(item) == cacheId })
		// 保存
		r.Set(list)
	}
}

func (r *cacheInMemory) Clear() {
	var list = r.Get()
	if list.Count() > 0 {
		list.Clear()
		r.Set(list)
	}
	delete(localCache, r.key)
}

func (r *cacheInMemory) Count() int {
	return r.Get().Count()
}

func (r *cacheInMemory) ExistsItem(cacheId string) bool {
	var list = r.Get()
	if list.Count() == 0 {
		return false
	}
	for index := 0; index < list.Count(); index++ {
		// 从当前缓存item中，获取唯一标识
		itemDataKey := r.GetUniqueId(list.Index(index))
		// 找到了
		if itemDataKey == cacheId {
			return true
		}
	}
	return false
}

func (r *cacheInMemory) ExistsKey() bool {
	lock.RLock()
	defer lock.RUnlock()
	_, ok := localCache[r.key]
	return ok
}

// GetUniqueId 获取唯一字段数据
func (r *cacheInMemory) GetUniqueId(item any) string {
	val := reflect.ValueOf(item).FieldByName(r.uniqueField).Interface()
	return parse.Convert(val, "")
}
