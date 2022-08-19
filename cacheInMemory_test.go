package memoryCache

import (
	"github.com/farseer-go/cache"
	"github.com/farseer-go/collections"
	"github.com/farseer-go/fs/modules"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCacheKey_Set(t *testing.T) {
	modules.StartModules(Module{})
	type po struct {
		Name string
		Age  int
	}
	cache.SetProfilesInMemory[po]("test", "Name", 0)
	cacheManage := cache.GetCacheManage[po]("test")
	lst := collections.NewList(po{Name: "steden", Age: 18}, po{Name: "steden2", Age: 19})
	cacheManage.Set(lst)
	lst2 := cacheManage.Get()

	assert.Equal(t, lst.Count(), lst2.Count())

	for i := 0; i < lst.Count(); i++ {
		assert.Equal(t, lst.Index(i).Name, lst2.Index(i).Name)
		assert.Equal(t, lst.Index(i).Age, lst2.Index(i).Age)
	}
}

func TestCacheKey_GetItem(t *testing.T) {
	modules.StartModules(Module{})
	type po struct {
		Name string
		Age  int
	}
	cache.SetProfilesInMemory[po]("test", "Name", 0)
	cacheManage := cache.GetCacheManage[po]("test")
	cacheManage.Set(collections.NewList(po{Name: "steden", Age: 18}, po{Name: "steden2", Age: 19}))
	item1, _ := cacheManage.GetItem("steden")

	assert.Equal(t, item1.Name, "steden")
	assert.Equal(t, item1.Age, 18)

	item2, _ := cacheManage.GetItem("steden2")

	assert.Equal(t, item2.Name, "steden2")
	assert.Equal(t, item2.Age, 19)
}

func TestCacheKey_SaveItem(t *testing.T) {
	modules.StartModules(Module{})
	type po struct {
		Name string
		Age  int
	}
	cache.SetProfilesInMemory[po]("test", "Name", 0)
	cacheManage := cache.GetCacheManage[po]("test")
	cacheManage.Set(collections.NewList(po{Name: "steden", Age: 18}, po{Name: "steden2", Age: 19}))
	cacheManage.SaveItem(po{Name: "steden", Age: 99})
	item1, _ := cacheManage.GetItem("steden")

	assert.Equal(t, item1.Name, "steden")
	assert.Equal(t, item1.Age, 99)

	item2, _ := cacheManage.GetItem("steden2")

	assert.Equal(t, item2.Name, "steden2")
	assert.Equal(t, item2.Age, 19)
}

func TestCacheKey_Remove(t *testing.T) {
	modules.StartModules(Module{})
	type po struct {
		Name string
		Age  int
	}
	cache.SetProfilesInMemory[po]("test", "Name", 0)
	cacheManage := cache.GetCacheManage[po]("test")
	cacheManage.Set(collections.NewList(po{Name: "steden", Age: 18}, po{Name: "steden2", Age: 19}))
	cacheManage.Remove("steden")

	_, exists := cacheManage.GetItem("steden")
	assert.False(t, exists)

	item2, _ := cacheManage.GetItem("steden2")
	assert.Equal(t, item2.Name, "steden2")
	assert.Equal(t, item2.Age, 19)
}

func TestCacheKey_Clear(t *testing.T) {
	modules.StartModules(Module{})
	type po struct {
		Name string
		Age  int
	}
	cache.SetProfilesInMemory[po]("test", "Name", 0)
	cacheManage := cache.GetCacheManage[po]("test")
	cacheManage.Set(collections.NewList(po{Name: "steden", Age: 18}, po{Name: "steden2", Age: 19}))
	assert.Equal(t, cacheManage.Count(), 2)
	cacheManage.Clear()
	assert.Equal(t, cacheManage.Count(), 0)
}

func TestCacheKey_Exists(t *testing.T) {
	modules.StartModules(Module{})
	type po struct {
		Name string
		Age  int
	}
	cache.SetProfilesInMemory[po]("test", "Name", 0)
	cacheManage := cache.GetCacheManage[po]("test")
	assert.False(t, cacheManage.Exists())
	cacheManage.Set(collections.NewList(po{Name: "steden", Age: 18}, po{Name: "steden2", Age: 19}))
	assert.True(t, cacheManage.Exists())
}
