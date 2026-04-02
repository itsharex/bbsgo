package cache

import (
	"bbsgo/config"
	"log"
	"time"

	"github.com/dgraph-io/ristretto"
)

var Cache *ristretto.Cache

func Init() {
	var err error
	Cache, err = ristretto.NewCache(&ristretto.Config{
		NumCounters: int64(config.GetConfigInt("cache_num_counters", 10000)),
		MaxCost:     int64(config.GetConfigInt("cache_max_cost", 10000000)),
		BufferItems: 64,
	})
	if err != nil {
		log.Fatalf("Failed to initialize cache: %v", err)
	}
	log.Println("Cache initialized successfully")
}

func Set(key string, value interface{}, ttl time.Duration) {
	Cache.SetWithTTL(key, value, 1, ttl)
}

func Get(key string) (interface{}, bool) {
	return Cache.Get(key)
}

func Delete(key string) {
	Cache.Del(key)
}

func DeletePattern(pattern string) {
}
