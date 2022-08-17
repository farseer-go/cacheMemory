## memoryCache
本地缓存操作

## What are the functions?
* memoryCache
  * func
      * GetCacheManage （获取缓存对象）
      * SetProfilesInMemory （设置内存缓存（集合））
      * SetProfilesInRedis （设置Redis缓存（集合））
      * SetProfilesInMemoryAndRedis （设置内存-Redis缓存（集合））
      * SetSingleProfilesInMemory （设置内存缓存（缓存单个对象））
      * SetSingleProfilesInRedis （设置Redis缓存（缓存单个对象））
      * SetSingleProfilesInMemoryAndRedis （设置内存-Redis缓存（缓存单个对象））
  * struct
    * CacheManage（缓存对象）
      * Get（获取缓存数据）
      * GetItem（从集合中获取指定cacheId的元素）
      * Set（保存缓存）
      * SaveItem（更新缓存）
      * Remove（移除缓存）
      * Clear（清空缓存）
      * Exists（缓存是否存在）
      * Count（数据集合的数量）

## Getting Started
```go
type po struct {
    Name string
    Age  int
}
```

### 单缓存操作
```go
// 配置缓存
// arg1: Key name
// arg2: UniqueField
// arg3：MemoryExpiry
memoryCache.SetProfilesInMemory[po]("test", "Name", 0)
// 获取key=test的缓存对象
cacheManage := memoryCache.GetCacheManage[po]("test")
// 设置缓存
cacheManage.Set(collections.NewList(po{Name: "steden", Age: 18}, po{Name: "steden2", Age: 19}))
// 获取Name=steden 的缓存item
cacheManage.GetItem("steden")
// 更新Name=steden的缓存item
cacheManage.SaveItem(po{Name: "steden", Age: 99})
// 移除Name=steden缓存item
cacheManage.Remove("steden")
// 清除缓存
cacheManage.Clear()
// 检查缓存是否存在
cacheManage.Exists()
// 获取缓存的集合数量
cacheManage.Count()
```