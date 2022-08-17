package memoryCache

import (
	"github.com/farseer-go/fs/exception"
	"github.com/farseer-go/memoryCache/eumCacheStoreType"
	"reflect"
	"time"
)

// SetProfilesInMemory 设置内存缓存（集合）
func SetProfilesInMemory[TEntity any](key string, uniqueField string, memoryExpiry time.Duration) {
	if uniqueField == "" {
		exception.ThrowRefuseException("缓存集合数据时，需要设置UniqueField字段")
	}
	var entity TEntity
	_, isExists := reflect.TypeOf(entity).FieldByName(uniqueField)
	if !isExists {
		exception.ThrowRefuseException(uniqueField + "字段，在缓存集合中不存在")
	}

	cacheConfigure[key] = CacheManage[TEntity]{
		CacheKey: CacheKey{
			Key:            key,
			CacheStoreType: eumCacheStoreType.Memory,
			MemoryExpiry:   memoryExpiry,
			UniqueField:    uniqueField,
			Cache:          newCacheInMemory(),
		},
	}
}

// SetProfilesInRedis 设置Redis缓存（集合）
func SetProfilesInRedis[TEntity any](key string, redisConfigName string, uniqueField string, redisExpiry time.Duration) {
	if uniqueField == "" {
		exception.ThrowRefuseException("缓存集合数据时，需要设置UniqueField字段")
	}
	var entity TEntity
	_, isExists := reflect.TypeOf(entity).FieldByName(uniqueField)
	if !isExists {
		exception.ThrowRefuseException(uniqueField + "字段，在缓存集合中不存在")
	}

	cacheConfigure[key] = CacheManage[TEntity]{
		CacheKey: CacheKey{
			Key:             key,
			CacheStoreType:  eumCacheStoreType.Redis,
			RedisExpiry:     redisExpiry,
			UniqueField:     uniqueField,
			RedisConfigName: redisConfigName,
			Cache:           newCacheInMemory(),
		},
	}
}

// SetProfilesInMemoryAndRedis 设置内存-Redis缓存（集合）
func SetProfilesInMemoryAndRedis[TEntity any](key string, redisConfigName string, uniqueField string, redisExpiry time.Duration, memoryExpiry time.Duration) {
	if uniqueField == "" {
		exception.ThrowRefuseException("缓存集合数据时，需要设置UniqueField字段")
	}
	var entity TEntity
	_, isExists := reflect.TypeOf(entity).FieldByName(uniqueField)
	if !isExists {
		exception.ThrowRefuseException(uniqueField + "字段，在缓存集合中不存在")
	}

	cacheConfigure[key] = CacheManage[TEntity]{
		CacheKey: CacheKey{
			Key:             key,
			CacheStoreType:  eumCacheStoreType.MemoryAndRedis,
			RedisExpiry:     redisExpiry,
			MemoryExpiry:    memoryExpiry,
			UniqueField:     uniqueField,
			RedisConfigName: redisConfigName,
			Cache:           newCacheInMemory(),
		},
	}
}

// SetSingleProfilesInMemory 设置内存缓存（缓存单个对象）
func SetSingleProfilesInMemory[TEntity any](key string, memoryExpiry time.Duration) {
	cacheConfigure[key] = CacheManage[TEntity]{
		CacheKey: CacheKey{
			Key:            key,
			CacheStoreType: eumCacheStoreType.Memory,
			MemoryExpiry:   memoryExpiry,
			Cache:          newCacheInMemory(),
		},
	}
}

// SetSingleProfilesInRedis 设置Redis缓存（缓存单个对象）
func SetSingleProfilesInRedis[TEntity any](key string, redisConfigName string, redisExpiry time.Duration) {
	cacheConfigure[key] = CacheManage[TEntity]{
		CacheKey: CacheKey{
			Key:             key,
			CacheStoreType:  eumCacheStoreType.Redis,
			RedisExpiry:     redisExpiry,
			RedisConfigName: redisConfigName,
			Cache:           newCacheInMemory(),
		},
	}
}

// SetSingleProfilesInMemoryAndRedis 设置内存-Redis缓存（缓存单个对象）
func SetSingleProfilesInMemoryAndRedis[TEntity any](key string, redisConfigName string, redisExpiry time.Duration, memoryExpiry time.Duration) {
	cacheConfigure[key] = CacheManage[TEntity]{
		CacheKey: CacheKey{
			Key:             key,
			CacheStoreType:  eumCacheStoreType.MemoryAndRedis,
			RedisExpiry:     redisExpiry,
			MemoryExpiry:    memoryExpiry,
			RedisConfigName: redisConfigName,
			Cache:           newCacheInMemory(),
		},
	}
}
