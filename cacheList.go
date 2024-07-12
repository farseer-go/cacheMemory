package cacheMemory

import "github.com/farseer-go/collections"

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
		item := getItem()
		if item != nil {
			receiver.items[key] = item
		}
	}
}

// Get 获取数据
func (receiver *CacheList[T]) Get(key string) *T {
	return receiver.items[key]
}

// ToList 得到所有缓存项
func (receiver *CacheList[T]) ToList() collections.List[T] {
	lst := collections.NewList[T]()
	for _, t := range receiver.items {
		lst.Add(*t)
	}
	return lst
}
