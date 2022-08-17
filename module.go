package memoryCache

import (
	"github.com/farseer-go/fs/core/container"
	"github.com/farseer-go/fs/modules"
)

type Module struct {
}

func (module Module) DependsModule() []modules.FarseerModule {
	return []modules.FarseerModule{modules.FarseerKernelModule{}}
}

func (module Module) PreInitialize() {
	localCache = make(map[string]cacheValue)
	cacheConfigure = make(map[string]any)
	_ = container.Register(func() ICache { return newCacheInMemory() })
}

func (module Module) Initialize() {
}

func (module Module) PostInitialize() {
	go checkTtl()
}

func (module Module) Shutdown() {
}
