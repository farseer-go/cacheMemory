package memoryCache

import "time"

func Get() {

}

// Set 设置缓存
func Set(key string, val any, absoluteExpiration time.Duration) {
	localCache[key] = val
}
