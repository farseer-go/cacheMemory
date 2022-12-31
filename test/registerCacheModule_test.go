package test

import (
	"github.com/farseer-go/cache"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestRegisterCacheModule(t *testing.T) {
	assert.Panics(t, func() {
		cache.RegisterCacheModule[po]("", "", "", nil)
	})

	assert.Panics(t, func() {
		cache.RegisterCacheModule[po]("", "", "a", nil)
	})
}
