package cacheMemory

import (
	"github.com/farseer-go/cache"
	"github.com/farseer-go/fs/flog"
	"reflect"
)

// SetProfiles 设置内存缓存（集合）
// key：这批数据的key，用于区分不同数据集合。
// uniqueField：数据集合中，用于区分Item的唯一值的字段名称（主键）
// ops：选项
func SetProfiles[TEntity any](key string, uniqueField string, ops ...cache.Option) cache.ICacheManage[TEntity] {
	if uniqueField == "" {
		flog.Panicf("cacheMemory.SetProfiles：缓存集合数据时，需要设置UniqueField字段")
	}
	var entity TEntity
	entityType := reflect.TypeOf(entity)
	_, isExists := entityType.FieldByName(uniqueField)
	if !isExists {
		flog.Panicf("cacheMemory.SetProfiles：%s字段，在缓存集合中不存在", uniqueField)
	}

	cacheIns := newCache(key, uniqueField, entityType, ops...)
	return cache.RegisterCacheModule[TEntity](key, "memory", uniqueField, cacheIns)
}
