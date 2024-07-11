package cacheMemory

type CacheList[T any] struct {
	items map[string]*T
}

func NewCache[T any]() CacheList[T] {
	return CacheList[T]{
		items: make(map[string]*T),
	}
}

// SetIfNotExists 写入缓存，如果Key不存在
func (receiver *CacheList[T]) SetIfNotExists(key string, getItem func() *T) {
	if _, isExists := receiver.items[key]; !isExists {
		receiver.items[key] = getItem()
	}
}

// Get 获取数据
func (receiver *CacheList[T]) Get(key string) *T {
	return receiver.items[key]
}
