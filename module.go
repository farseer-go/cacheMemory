package memoryCache

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
	localCache = make(map[string]cacheValue)
	_ = container.RegisterSingle(func() cache.ICache { return newCacheInMemory() })
}

func (module Module) Initialize() {
}

func (module Module) PostInitialize() {
	go checkTtl()
}

func (module Module) Shutdown() {
}
