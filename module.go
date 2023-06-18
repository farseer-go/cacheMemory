package cacheMemory

import (
	"github.com/farseer-go/cache"
	"github.com/farseer-go/fs/modules"
)

type Module struct {
}

func (module Module) DependsModule() []modules.FarseerModule {
	return []modules.FarseerModule{cache.Module{}}
}

func (module Module) PostInitialize() {
	go checkTTL()
}
