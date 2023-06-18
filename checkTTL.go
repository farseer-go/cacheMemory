package cacheMemory

import (
	"github.com/farseer-go/collections"
	"time"
)

var checkList = collections.NewList[*cacheInMemory]()

// 检查所有缓存过期项，缓存失效时，移除缓存
func checkTTL() {
	ticker := time.NewTicker(200 * time.Millisecond)
	for {
		<-ticker.C
		for i := 0; i < checkList.Count(); i++ {
			r := checkList.Index(i)

			if r.ttlExpiry > 0 && time.Now().UnixMilli() >= r.ttlExpiry {
				r.lock.Lock()
				r.Clear()
				// 不能使用NewListAny，因为需要data = nil
				r.data = collections.ListAny{}
				// 重新计算下一次的失效时间
				r.ttlExpiry = time.Now().Add(r.expiry).UnixMilli()
				r.lock.Unlock()
			}
		}
	}
}
