package cacheMemory

type List[T any] struct {
	items map[string]*T
}

// SetIfNotExists 写入缓存，如果Key不存在
func (receiver *List[T]) SetIfNotExists(key string, getItem func() *T) {
	if _, isExists := receiver.items[key]; !isExists {
		receiver.items[key] = getItem()
	}
}

// Get 获取数据
func (receiver *List[T]) Get(key string) *T {
	return receiver.items[key]
}
