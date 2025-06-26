// Package redis_test contains tests for Redis cache operations using the cache interface.
// It ensures the proper functioning of setting, retrieving, deleting, and managing keys
// in Redis, as well as connection management and expiration handling.
package redis_test

import (
	"context"
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/zeroxsolutions/alex"
	"github.com/zeroxsolutions/banshee/redis"
	"github.com/zeroxsolutions/barbatos/cache"
	"github.com/zeroxsolutions/strike/ssutil"
)

// initRedisCache initializes a Redis cache instance using environment variables
// and returns a cache.Cache implementation. It will terminate the test if configuration
// fails.
func initRedisCache(t *testing.T) cache.Cache {
	addr := os.Getenv("REDIS_ADDRESS")
	password := os.Getenv("REDIS_PASSWORD")
	dbRaw := os.Getenv("REDIS_DB")
	db := 0

	// Convert Redis DB from string to integer, failing the test if conversion fails.
	if dbRaw != "" {
		dbConverted, err := strconv.Atoi(dbRaw)
		if err != nil {
			t.Fatal(err)
		}
		db = dbConverted
	}

	// Initialize Redis configuration.
	redisCacheConfig := alex.RedisConfig{
		Addr:     addr,
		Password: password,
		DB:       db,
	}

	// Create a Redis cache instance, terminating the test on error.
	redisCache, err := redis.NewRedisCache(&redisCacheConfig)
	if err != nil {
		t.Fatal(err)
	}
	return redisCache
}

// TestRedisCache groups multiple test cases to validate Redis cache behavior.
func TestRedisCache(t *testing.T) {

	// Test the Redis connection state.
	t.Run("IsConnected", func(t *testing.T) {
		redisCache := initRedisCache(t)

		defer func(redisCache cache.Cache) {
			if err := redisCache.Close(); err != nil {
				t.Log("Close Redis cache connection err", err)
			}
		}(redisCache)

		if isConnected := redisCache.IsConnected(context.Background()); !isConnected {
			t.FailNow() // Terminate if Redis is not connected.
		}
	})

	// Test the storage and retrieval of multiple keys.
	t.Run("Keys", func(t *testing.T) {
		redisCache := initRedisCache(t)

		defer func(redisCache cache.Cache) {
			if err := redisCache.Close(); err != nil {
				t.Log("Close Redis cache connection err", err)
			}
		}(redisCache)

		// Delete any keys matching the pattern "key*".
		if err := redisCache.DelWithPattern(context.Background(), "key*"); err != nil {
			t.Error(err)
		}

		// Generate and store multiple key-value pairs.
		kv := map[string]string{}

		for range make([]int, 10) {
			key := "key" + ssutil.MakeString(10)
			value := ssutil.MakeString(12)
			kv[key] = value
			if err := redisCache.Set(context.Background(), key, value); err != nil {
				t.Error(err)
			}
		}

		if len(kv) == 0 {
			t.FailNow()
		}

		// Retrieve and verify keys match stored values.
		keys, err := redisCache.Keys(context.Background(), "key*")

		if err != nil {
			t.Error(err)
		}

		if len(keys) != len(kv) {
			t.FailNow()
		}

		for _, key := range keys {
			value, err := redisCache.Get(context.Background(), key)
			if err != nil {
				t.Error(err)
			}
			if value != kv[key] {
				t.FailNow()
			}
		}
	})

	// Test the Get operation to retrieve a value by key.
	t.Run("Get", func(t *testing.T) {
		redisCache := initRedisCache(t)

		defer func(redisCache cache.Cache) {
			if err := redisCache.Close(); err != nil {
				t.Log("Close Redis cache connection err", err)
			}
		}(redisCache)

		key := ssutil.MakeString(10)
		value := ssutil.MakeString(12)

		if err := redisCache.Set(context.Background(), key, value); err != nil {
			t.Error(err)
		}

		v, err := redisCache.Get(context.Background(), key)

		if err != nil {
			t.Error(err)
		}

		if v != value {
			t.FailNow()
		}

	})

	// Test the Set operation to store a value by key.
	t.Run("Set", func(t *testing.T) {
		redisCache := initRedisCache(t)
		key := ssutil.MakeString(10)
		value := ssutil.MakeString(12)

		if err := redisCache.Set(context.Background(), key, value); err != nil {
			t.Error(err)
		}

	})

	// Test setting a key with expiration and verifying expiration works.
	t.Run("SetWithExpiration", func(t *testing.T) {
		redisCache := initRedisCache(t)

		defer func(redisCache cache.Cache) {
			if err := redisCache.Close(); err != nil {
				t.Log("Close Redis cache connection err", err)
			}
		}(redisCache)

		key := ssutil.MakeString(10)
		value := ssutil.MakeString(12)

		if err := redisCache.SetWithExpiration(context.Background(), key, value, 2*time.Second); err != nil {
			t.Error(err)
		}

		v, err := redisCache.Get(context.Background(), key)

		if err != nil {
			t.Error(err)
		}

		if v != value {
			t.FailNow()
		}

		time.Sleep(3 * time.Second)
		if _, err := redisCache.Get(context.Background(), key); err != cache.ErrCacheNil {
			t.Log(err)
			t.FailNow()
		}
	})

	// Test deleting a specific key.
	t.Run("Del", func(t *testing.T) {
		redisCache := initRedisCache(t)

		defer func(redisCache cache.Cache) {
			if err := redisCache.Close(); err != nil {
				t.Log("Close Redis cache connection err", err)
			}
		}(redisCache)

		key := ssutil.MakeString(10)
		value := ssutil.MakeString(12)

		if err := redisCache.Set(context.Background(), key, value); err != nil {
			t.Error(err)
		}

		v, err := redisCache.Get(context.Background(), key)

		if err != nil {
			t.Error(err)
		}

		if v != value {
			t.FailNow()
		}

		if err := redisCache.Del(context.Background(), key); err != nil {
			t.Error(err)
		}

		if _, err := redisCache.Get(context.Background(), key); err != cache.ErrCacheNil {
			t.Log(err)
			t.FailNow()
		}

	})

	// Test deleting keys matching a pattern.
	t.Run("DelWithPattern", func(t *testing.T) {
		redisCache := initRedisCache(t)

		defer func(redisCache cache.Cache) {
			if err := redisCache.Close(); err != nil {
				t.Log("Close Redis cache connection err", err)
			}
		}(redisCache)

		if err := redisCache.DelWithPattern(context.Background(), "key*"); err != nil {
			t.Error(err)
		}

		for range make([]int, 10) {
			key := "key" + ssutil.MakeString(10)
			value := ssutil.MakeString(12)
			if err := redisCache.Set(context.Background(), key, value); err != nil {
				t.Error(err)
			}
		}

		if err := redisCache.DelWithPattern(context.Background(), "key*"); err != nil {
			t.Error(err)
		}

		keys, err := redisCache.Keys(context.Background(), "key*")

		if err != nil {
			t.Error(err)
		}

		if len(keys) != 0 {
			t.FailNow()
		}

	})

	// Test closing the Redis connection.
	t.Run("Close", func(t *testing.T) {
		redisCache := initRedisCache(t)

		if isConnected := redisCache.IsConnected(context.Background()); !isConnected {
			t.FailNow()
		}

		if err := redisCache.Close(); err != nil {
			t.Error(err)
		}

		if isConnected := redisCache.IsConnected(context.Background()); isConnected {
			t.FailNow()
		}
	})
}
