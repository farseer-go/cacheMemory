package memoryCache

import (
	"github.com/farseer-go/collections"
	"github.com/farseer-go/fs/modules"
	"testing"
)

func TestCacheKey_Set(t *testing.T) {
	modules.StartModules(Module{})
	type po struct {
		Name string
		Age  int
	}
	lst := collections.NewList(po{Name: "steden", Age: 18}, po{Name: "steden2", Age: 19})
	SetProfilesInMemory[po]("test", "Name", 0)
	cacheKey := GetCacheManage[po]("test")
	cacheKey.Set(lst)
	lst2 := cacheKey.Get()
	if lst.Count() != lst2.Count() {
		t.Error()
	}

	for i := 0; i < lst.Count(); i++ {
		if lst.Index(i).Name != lst2.Index(i).Name || lst.Index(i).Age != lst2.Index(i).Age {
			t.Error()
		}
	}
}

func TestCacheKey_GetItem(t *testing.T) {
	modules.StartModules(Module{})
	type po struct {
		Name string
		Age  int
	}
	SetProfilesInMemory[po]("test", "Name", 0)
	cacheKey := GetCacheManage[po]("test")
	cacheKey.Set(collections.NewList(po{Name: "steden", Age: 18}, po{Name: "steden2", Age: 19}))
	item1, _ := cacheKey.GetItem("steden")
	if item1.Name != "steden" || item1.Age != 18 {
		t.Error()
	}

	item2, _ := cacheKey.GetItem("steden2")
	if item2.Name != "steden2" || item2.Age != 19 {
		t.Error()
	}
}

func TestCacheKey_SaveItem(t *testing.T) {
	modules.StartModules(Module{})
	type po struct {
		Name string
		Age  int
	}
	SetProfilesInMemory[po]("test", "Name", 0)
	cacheKey := GetCacheManage[po]("test")
	cacheKey.Set(collections.NewList(po{Name: "steden", Age: 18}, po{Name: "steden2", Age: 19}))
	cacheKey.SaveItem(po{Name: "steden", Age: 99})
	item1, _ := cacheKey.GetItem("steden")
	if item1.Name != "steden" || item1.Age != 99 {
		t.Error()
	}

	item2, _ := cacheKey.GetItem("steden2")
	if item2.Name != "steden2" || item2.Age != 19 {
		t.Error()
	}
}

func TestCacheKey_Remove(t *testing.T) {
	modules.StartModules(Module{})
	type po struct {
		Name string
		Age  int
	}
	SetProfilesInMemory[po]("test", "Name", 0)
	cacheKey := GetCacheManage[po]("test")
	cacheKey.Set(collections.NewList(po{Name: "steden", Age: 18}, po{Name: "steden2", Age: 19}))
	cacheKey.Remove("steden")

	_, exists := cacheKey.GetItem("steden")
	if exists {
		t.Error()
	}

	item2, _ := cacheKey.GetItem("steden2")
	if item2.Name != "steden2" || item2.Age != 19 {
		t.Error()
	}
}

func TestCacheKey_Clear(t *testing.T) {
	modules.StartModules(Module{})
	type po struct {
		Name string
		Age  int
	}
	SetProfilesInMemory[po]("test", "Name", 0)
	cacheKey := GetCacheManage[po]("test")
	cacheKey.Set(collections.NewList(po{Name: "steden", Age: 18}, po{Name: "steden2", Age: 19}))
	if cacheKey.Count() != 2 {
		t.Error()
	}
	cacheKey.Clear()
	if cacheKey.Count() != 0 {
		t.Error()
	}
}

func TestCacheKey_Exists(t *testing.T) {
	modules.StartModules(Module{})
	type po struct {
		Name string
		Age  int
	}
	SetProfilesInMemory[po]("test", "Name", 0)
	cacheKey := GetCacheManage[po]("test")
	if cacheKey.Exists() {
		t.Error()
	}
	cacheKey.Set(collections.NewList(po{Name: "steden", Age: 18}, po{Name: "steden2", Age: 19}))
	if !cacheKey.Exists() {
		t.Error()
	}
}
