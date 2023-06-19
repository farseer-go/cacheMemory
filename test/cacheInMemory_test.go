package test

import (
	"github.com/farseer-go/cache"
	"github.com/farseer-go/cache/eumExpiryType"
	"github.com/farseer-go/cacheMemory"
	"github.com/farseer-go/collections"
	"github.com/farseer-go/fs"
	"github.com/farseer-go/fs/container"
	"github.com/farseer-go/fs/flog"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

type po struct {
	Name string
	Age  int
}

func init() {
	fs.Initialize[cacheMemory.Module]("unit test")
}

func TestCacheInMemory_GetItem(t *testing.T) {

	cacheMemory.SetProfiles[po]("test1", "Name")

	cacheManage := container.Resolve[cache.ICacheManage[po]]("test1")
	cacheManage.SetItemSource(func(cacheId any) (po, bool) {
		if cacheId == "laoLi" {
			return po{Name: "laoLi"}, true
		}
		return po{}, false
	})
	cacheManage.SetListSource(func() collections.List[po] {
		return collections.NewList(po{Name: "xiaoLi"})
	})
	cacheManage.EnableItemNullToLoadAll()

	item, b := cacheManage.GetItem("laoLi")
	assert.Equal(t, "laoLi", item.Name)
	assert.True(t, b)

	item, b = cacheManage.GetItem("xiaoLi")
	assert.Equal(t, "xiaoLi", item.Name)
	assert.True(t, b)

	cacheManage.Set(po{Name: "steden", Age: 18}, po{Name: "steden2", Age: 19})
	item1, _ := cacheManage.GetItem("steden")

	assert.Equal(t, item1.Name, "steden")
	assert.Equal(t, item1.Age, 18)

	item2, _ := cacheManage.GetItem("steden2")

	assert.Equal(t, item2.Name, "steden2")
	assert.Equal(t, item2.Age, 19)

	cacheManage.Single()
}

func TestCacheInMemory_Set(t *testing.T) {
	cacheMemory.SetProfiles[po]("test2", "Name")
	assert.Panics(t, func() {
		cacheMemory.SetProfiles[po]("test2", "")
	})

	assert.Panics(t, func() {
		cacheMemory.SetProfiles[po]("test2", "ClientName")
	})
	cacheManage := container.Resolve[cache.ICacheManage[po]]("test2")
	cacheManage.SetListSource(func() collections.List[po] {
		return collections.NewList[po]()
	})
	cacheManage.SetItemSource(func(cacheId any) (po, bool) {
		return po{}, false
	})
	lstDb := collections.NewDictionary[string, po]()
	cacheManage.SetSyncSource(10*time.Millisecond, func(val po) {
		if !lstDb.ContainsKey(val.Name) {
			lstDb.Add(val.Name, val)
			flog.Info(val.Name)
		}
	})
	cacheManage.EnableItemNullToLoadAll()

	lstEmpty := cacheManage.Get()
	assert.Equal(t, 0, lstEmpty.Count())

	assert.False(t, cacheManage.ExistsItem("steden"))
	lst := collections.NewList(po{Name: "steden", Age: 18}, po{Name: "steden2", Age: 19})
	cacheManage.Set(lst.ToArray()...)
	assert.True(t, cacheManage.ExistsItem("steden"))
	assert.False(t, cacheManage.ExistsItem("steden3"))
	lst2 := cacheManage.Get()

	assert.Equal(t, lst.Count(), lst2.Count())

	for i := 0; i < lst.Count(); i++ {
		assert.Equal(t, lst.Index(i).Name, lst2.Index(i).Name)
		assert.Equal(t, lst.Index(i).Age, lst2.Index(i).Age)
	}

	time.Sleep(20 * time.Millisecond)
	assert.True(t, lstDb.ContainsKey("steden"))
	assert.True(t, lstDb.ContainsKey("steden2"))
}

func TestCacheInMemory_SaveItem(t *testing.T) {
	cacheMemory.SetProfiles[po]("test3", "Name")
	cacheManage := container.Resolve[cache.ICacheManage[po]]("test3")

	cacheManage.SaveItem(po{Name: "steden3", Age: 121})
	item0, _ := cacheManage.GetItem("steden3")
	assert.Equal(t, item0.Name, "steden3")
	assert.Equal(t, item0.Age, 121)

	cacheManage.Set(po{Name: "steden", Age: 18}, po{Name: "steden2", Age: 19})
	cacheManage.SaveItem(po{Name: "steden", Age: 99})
	item1, _ := cacheManage.GetItem("steden")

	assert.Equal(t, item1.Name, "steden")
	assert.Equal(t, item1.Age, 99)

	item2, _ := cacheManage.GetItem("steden2")

	assert.Equal(t, item2.Name, "steden2")
	assert.Equal(t, item2.Age, 19)
}

func TestCacheInMemory_Remove(t *testing.T) {
	cacheMemory.SetProfiles[po]("test4", "Name")
	cacheManage := container.Resolve[cache.ICacheManage[po]]("test4")
	cacheManage.Set(po{Name: "steden", Age: 18}, po{Name: "steden2", Age: 19})
	cacheManage.Remove("steden")

	_, exists := cacheManage.GetItem("steden")
	assert.False(t, exists)

	item2, _ := cacheManage.GetItem("steden2")
	assert.Equal(t, item2.Name, "steden2")
	assert.Equal(t, item2.Age, 19)
}

func TestCacheInMemory_Clear(t *testing.T) {
	cacheMemory.SetProfiles[po]("test5", "Name")
	cacheManage := container.Resolve[cache.ICacheManage[po]]("test5")
	cacheManage.Set(po{Name: "steden", Age: 18}, po{Name: "steden2", Age: 19})
	assert.Equal(t, cacheManage.Count(), 2)
	cacheManage.Clear()
	assert.Equal(t, cacheManage.Count(), 0)
}

func TestCacheInMemory_Exists(t *testing.T) {
	cacheMemory.SetProfiles[po]("test6", "Name")
	cacheManage := container.Resolve[cache.ICacheManage[po]]("test6")
	assert.False(t, cacheManage.ExistsKey())
	cacheManage.Set(po{Name: "steden", Age: 18}, po{Name: "steden2", Age: 19})
	assert.True(t, cacheManage.ExistsKey())
}

func TestCacheInMemory_Ttl(t *testing.T) {
	cacheMemory.SetProfiles[po]("test7", "Name", func(op *cache.Op) {
		op.ExpiryType = eumExpiryType.SlidingExpiration
		op.Expiry = 200 * time.Millisecond
	})
	
	cacheManage := container.Resolve[cache.ICacheManage[po]]("test7")
	lst := collections.NewList(po{Name: "steden", Age: 18}, po{Name: "steden2", Age: 19})
	cacheManage.Set(lst.ToArray()...)
	lst2 := cacheManage.Get()
	assert.Equal(t, lst.Count(), lst2.Count())
	for i := 0; i < lst.Count(); i++ {
		assert.Equal(t, lst.Index(i).Name, lst2.Index(i).Name)
		assert.Equal(t, lst.Index(i).Age, lst2.Index(i).Age)
	}
	time.Sleep(100 * time.Millisecond)
	cacheManage.Get()
	time.Sleep(150 * time.Millisecond)
	assert.True(t, cacheManage.ExistsKey())
	time.Sleep(300 * time.Millisecond)
	assert.False(t, cacheManage.ExistsKey())
}
