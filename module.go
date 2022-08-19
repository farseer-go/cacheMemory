package cacheMemory

import (
	"github.com/farseer-go/cache"
	"github.com/farseer-go/fs/container"
	"github.com/farseer-go/fs/modules"
)

type Module struct {
}

func (module Module) DependsModule() []modules.FarseerModule {
	return []modules.FarseerModule{modules.FarseerKernelModule{}, cache.Module{}}
}

func (module Module) PreInitialize() {
	localCache = make(map[string]*cacheValue)
	container.Use[cache.ICache](func() cache.ICache { return newCacheInMemory() }).Name("memory").Register()
}

func (module Module) Initialize() {
}

func (module Module) PostInitialize() {
}

func (module Module) Shutdown() {
}
