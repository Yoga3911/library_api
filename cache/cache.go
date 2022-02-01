package cache

import (
	r "github.com/dgraph-io/ristretto"
	"log"
)

type Cache interface {
	Set(key string, value interface{})
	Get(key string) interface{}
	Del(key ...string)
}

type cache struct {
	cache *r.Cache
}

func NewCache() Cache {
	res, err := r.NewCache(&r.Config{
		BufferItems: 64,
		MaxCost:     1 << 30,
		NumCounters: 1e7,
	})
	if err != nil {
		log.Println(err)
	}

	return &cache{
		cache: res,
	}
}

func (c *cache) Set(key string, value interface{}) {
	c.cache.Set(key, value, 1)
}

func (c *cache) Get(key string) interface{} {
	val, _ := c.cache.Get(key)

	return val
}

func (c *cache) Del(key ...string) {
	for _, i := range key {
		c.cache.Del(i)
	}
}