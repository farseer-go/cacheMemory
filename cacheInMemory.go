package cacheMemory

import (
	"github.com/farseer-go/cache"
	"github.com/farseer-go/cache/eumExpiryType"
	"github.com/farseer-go/collections"
	"github.com/farseer-go/fs/parse"
	"reflect"
	"sync"
	"time"
)

// 二级缓存-本地缓存操作
type cacheInMemory struct {
	expiry      time.Duration       // 设置Memory缓存过期时间
	expiryType  eumExpiryType.Enum  // 过期策略
	uniqueField string              // hash中的主键（唯一ID的字段名称）
	itemType    reflect.Type        // itemType
	key         string              // 缓存KEY
	lock        *sync.RWMutex       // 锁
	data        collections.ListAny // 缓存的数据
	lastVisitAt time.Time           // 最后一次访问时间
}

// 创建实例
func newCache(key string, uniqueField string, itemType reflect.Type, ops ...cache.Option) cache.ICache {
	op := &cache.Op{}
	for _, option := range ops {
		option(op)
	}

	r := &cacheInMemory{
		expiry:      op.Expiry,
		expiryType:  op.ExpiryType,
		uniqueField: uniqueField,
		itemType:    itemType,
		key:         key,
		lock:        &sync.RWMutex{},
		lastVisitAt: time.Now(),
	}

	if r.expiry > 0 && r.expiryType == eumExpiryType.SlidingExpiration {
		go r.updateTtl()
	}

	return r
}

func (r *cacheInMemory) Get() collections.ListAny {
	r.lock.RLock()
	defer r.lock.RUnlock()
	// 更新缓存过期时间
	r.updateExpiry()
	return r.data
}

func (r *cacheInMemory) GetItem(cacheId any) any {
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
	r.lock.Lock()
	defer r.lock.Unlock()

	// 更新缓存过期时间
	r.updateExpiry()
	r.data = val
}

func (r *cacheInMemory) SaveItem(newVal any) {
	var list = r.Get()
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

	if list.Count() == 0 {
		list = collections.NewListAny()
	}
	list.Add(newVal)

	// 保存
	r.Set(list)
}

func (r *cacheInMemory) Remove(cacheId any) {
	var list = r.Get()
	if list.Count() > 0 {
		list.RemoveAll(func(item any) bool { return r.GetUniqueId(item) == cacheId })
	}
}

func (r *cacheInMemory) Clear() {
	if r.data.Count() > 0 {
		r.data.Clear()
	}
}

func (r *cacheInMemory) Count() int {
	return r.Get().Count()
}

func (r *cacheInMemory) ExistsItem(cacheId any) bool {
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
	r.lock.RLock()
	defer r.lock.RUnlock()
	return !r.data.IsNil()
}

func (r *cacheInMemory) GetUniqueId(item any) string {
	val := reflect.ValueOf(item).FieldByName(r.uniqueField).Interface()
	return parse.Convert(val, "")
}

// 更新缓存过期时间
func (r *cacheInMemory) updateExpiry() {
	if r.expiry > 0 && r.expiryType == eumExpiryType.SlidingExpiration {
		r.lastVisitAt = time.Now()
	}
}

// TTL时间到期后没有访问数据，则移除缓存数据
func (r *cacheInMemory) updateTtl() {
	expiry := r.expiry
	if r.expiry >= 2*time.Second {
		expiry = r.expiry - time.Second
	} else if r.expiry >= time.Second {
		expiry = r.expiry - 500*time.Millisecond
	} else if r.expiry >= 500*time.Millisecond {
		expiry = r.expiry - 100*time.Millisecond
	}

	ticker := time.NewTicker(expiry)
	for range ticker.C {
		if time.Now().Sub(r.lastVisitAt) > r.expiry {
			// 重新计算下一次的失效时间
			r.lastVisitAt = time.Time{}
			r.lock.Lock()
			r.Clear()
			// 不能使用NewListAny，因为需要data = nil
			r.data = collections.ListAny{}
			r.lock.Unlock()
		}
	}
}
