package eumCacheStoreType

type Enum int

const (
	Memory Enum = iota
	Redis
	MemoryAndRedis
)
