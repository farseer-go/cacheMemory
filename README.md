## memoryCache
`github.com/farseer-go/cache`组件的本地缓存实现

实现了`cache.ICache`接口，并通过`container`注册到`IOC`

如果你需要开启本地缓存，则需要在`startupModule`中依赖`cacheMemory.Module`模块