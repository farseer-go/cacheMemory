package memoryCache

import "time"

// 检查缓存失效
func checkTtl() {
	for {
		for key, value := range localCache {
			if value.ttl > 0 && time.Now().UnixMicro() > value.ttl {
				delete(localCache, key)
			}
		}
		time.Sleep(100 * time.Millisecond)
	}
}
