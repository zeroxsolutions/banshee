// Package redis provides a Redis-based implementation of the cache interface.
// This package offers a production-ready cache solution using Redis as the backend,
// supporting all standard cache operations including key-value storage, pattern matching,
// expiration handling, and connection management.
package redis

import (
	"context"
	"errors"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/zeroxsolutions/alex"
	"github.com/zeroxsolutions/barbatos/cache"
)

// NewRedisCache creates a new Redis-based cache implementation with the provided configuration.
// This constructor establishes a connection to Redis, verifies connectivity, and returns
// a fully initialized cache instance ready for use.
//
// The function performs the following initialization steps:
//   - Creates a Redis client with the provided configuration
//   - Tests the connection using a PING command
//   - Returns an error if connection fails
//   - Wraps the client in a RedisCache struct implementing the Cache interface
//
// Configuration options include:
//   - Addr: Redis server address (host:port format)
//   - Password: Redis authentication password (if required)
//   - DB: Redis database number to use (0-15 typically)
//
// Parameters:
//   - config: Pointer to alex.RedisConfig containing Redis connection settings
//
// Returns:
//   - cache.Cache: A Redis cache implementation ready for use
//   - error: Connection error if Redis is unreachable or authentication fails
//
// Example:
//
//	config := &alex.RedisConfig{
//	    Addr:     "localhost:6379",
//	    Password: "",
//	    DB:       0,
//	}
//	cache, err := redis.NewRedisCache(config)
//	if err != nil {
//	    log.Fatal("Failed to connect to Redis:", err)
//	}
//	defer cache.Close()
func NewRedisCache(config *alex.RedisConfig) (cache.Cache, error) {
	client := redis.NewClient(
		&redis.Options{
			Addr:     config.Addr,
			Password: config.Password,
			DB:       config.DB,
		},
	)
	_, err := client.Ping(context.Background()).Result()
	if err != nil {
		return nil, err
	}
	return &RedisCache{client: client}, nil
}

// RedisCache implements the Cache interface using Redis as the backend storage.
// This struct wraps a Redis client and provides thread-safe cache operations
// with full Redis feature support including persistence, clustering, and advanced data types.
//
// The Redis implementation offers several advantages:
//   - High performance with sub-millisecond operations
//   - Persistence options for data durability
//   - Built-in expiration and memory management
//   - Support for clustering and high availability
//   - Rich pattern matching capabilities
//   - Atomic operations for consistency
//
// Thread safety: All operations are thread-safe as they delegate to the
// underlying Redis client which handles concurrent access properly.
type RedisCache struct {
	client *redis.Client
}

// IsConnected verifies the current connection status to the Redis server.
// This method uses Redis PING command to test connectivity and is useful for
// health checks, monitoring, and determining cache availability.
//
// The method performs a lightweight operation that:
//   - Sends a PING command to Redis
//   - Returns true if Redis responds with PONG
//   - Returns false for any network, authentication, or server errors
//   - Respects the provided context for timeout and cancellation
//
// This is particularly useful for:
//   - Application health checks
//   - Circuit breaker implementations
//   - Monitoring and alerting systems
//   - Graceful degradation scenarios
//
// Parameters:
//   - ctx: Context for timeout control and request cancellation
//
// Returns:
//   - bool: true if Redis is reachable and responsive, false otherwise
//
// Example:
//
//	if !cache.IsConnected(ctx) {
//	    log.Warn("Redis cache is unavailable, falling back to local cache")
//	    // Implement fallback logic
//	}
func (r *RedisCache) IsConnected(ctx context.Context) bool {
	_, err := r.client.Ping(ctx).Result()
	return err == nil
}

// Keys retrieves all Redis keys matching the specified pattern.
// This method uses Redis KEYS command to find keys that match the given pattern
// using Redis glob-style pattern matching.
//
// Pattern syntax supports:
//   - '*' matches zero or more characters
//   - '?' matches exactly one character
//   - '[abc]' matches any character in the set
//   - '[^abc]' matches any character not in the set
//   - '[a-z]' matches any character in the range
//   - '\' escapes special characters
//
// Performance considerations:
//   - KEYS command can be slow on large databases (O(N) complexity)
//   - May block Redis server during execution on large datasets
//   - Consider using SCAN for production environments with large key sets
//   - Use specific patterns to limit the result set
//
// Parameters:
//   - ctx: Context for timeout control and request cancellation
//   - pattern: Glob-style pattern to match keys against
//
// Returns:
//   - []string: Slice of keys matching the pattern (empty if no matches)
//   - error: Redis connection error or command execution error
//
// Examples:
//
//	keys, err := cache.Keys(ctx, "user:*")           // All user keys
//	keys, err := cache.Keys(ctx, "session:2023*")    // Sessions from 2023
//	keys, err := cache.Keys(ctx, "temp:???")         // 3-character temp keys
//	keys, err := cache.Keys(ctx, "cache:[0-9]*")     // Numbered cache keys
func (r *RedisCache) Keys(ctx context.Context, pattern string) ([]string, error) {
	keys, err := r.client.Keys(ctx, pattern).Result()
	if err != nil {
		return nil, err
	}
	return keys, nil
}

// Get retrieves the string value associated with the specified key from Redis.
// This method uses Redis GET command to fetch the value and handles the special
// case of non-existent keys by returning a standardized cache error.
//
// The method handles different scenarios:
//   - Returns the string value if the key exists
//   - Returns cache.ErrCacheNil if the key doesn't exist (Redis Nil reply)
//   - Returns other errors for connection issues or Redis errors
//   - Respects context timeout and cancellation
//
// Redis data types supported:
//   - String values (direct retrieval)
//   - Numbers stored as strings (returned as string representation)
//   - Other data types may return errors or unexpected results
//
// Parameters:
//   - ctx: Context for timeout control and request cancellation
//   - key: Redis key to retrieve the value for
//
// Returns:
//   - string: The value stored under the key
//   - error: cache.ErrCacheNil if key doesn't exist, other errors for failures
//
// Example:
//
//	value, err := cache.Get(ctx, "user:123")
//	if errors.Is(err, cache.ErrCacheNil) {
//	    // Key not found, handle accordingly
//	    return handleMissingUser()
//	} else if err != nil {
//	    // Redis error, handle connection issues
//	    return err
//	}
//	// Use the retrieved value
//	return processUser(value)
func (r *RedisCache) Get(ctx context.Context, key string) (string, error) {
	value, err := r.client.Get(ctx, key).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return "", cache.ErrCacheNil
		}
		return "", err
	}
	return value, nil
}

// Set stores a value in Redis under the specified key without expiration.
// This method provides a simple interface for persistent key-value storage
// by delegating to SetWithExpiration with zero expiration time.
//
// The stored value will persist until:
//   - Explicitly deleted using Del or DelWithPattern
//   - Overwritten by another Set operation
//   - Redis server is restarted (unless persistence is configured)
//   - Redis memory limits are reached and key is evicted
//
// Data type handling:
//   - Strings are stored directly
//   - Numbers are converted to string representation
//   - Complex types should be serialized by the caller
//   - Binary data should be base64 encoded or use Redis binary-safe commands
//
// Parameters:
//   - ctx: Context for timeout control and request cancellation
//   - key: Redis key to store the value under
//   - value: Value to store (will be converted to string by Redis client)
//
// Returns:
//   - error: Redis connection error or command execution error
//
// Example:
//
//	err := cache.Set(ctx, "user:123", "john_doe")
//	err := cache.Set(ctx, "counter", 42)
//	err := cache.Set(ctx, "config", jsonString)
func (r *RedisCache) Set(ctx context.Context, key string, value interface{}) error {
	return r.SetWithExpiration(ctx, key, value, 0)
}

// SetWithExpiration stores a value in Redis with automatic expiration after the specified duration.
// This method uses Redis SET command with EX/PX options to set both value and TTL atomically.
// It's essential for implementing session management, caching strategies, and temporary data storage.
//
// Expiration behavior:
//   - expiration = 0: Key persists indefinitely (same as Set)
//   - expiration > 0: Key automatically expires after the duration
//   - Sub-second precision supported using Redis PSETEX for milliseconds
//   - Expiration is absolute from the time of setting, not from last access
//
// TTL management:
//   - Redis handles expiration automatically
//   - Expired keys are removed by Redis background processes
//   - TTL can be updated by calling SetWithExpiration again
//   - TTL can be removed by calling Set (makes key persistent)
//
// Use cases:
//   - Session storage with automatic cleanup
//   - Rate limiting with time windows
//   - Temporary data that should auto-expire
//   - Caching with refresh intervals
//
// Parameters:
//   - ctx: Context for timeout control and request cancellation
//   - key: Redis key to store the value under
//   - value: Value to store (will be converted to string by Redis client)
//   - expiration: Duration after which the key should automatically expire
//
// Returns:
//   - error: Redis connection error or command execution error
//
// Examples:
//
//	// Session expires in 1 hour
//	err := cache.SetWithExpiration(ctx, "session:abc123", sessionData, time.Hour)
//
//	// Rate limit counter expires in 1 minute
//	err := cache.SetWithExpiration(ctx, "rate:user:123", "1", time.Minute)
//
//	// Cache entry expires in 5 minutes
//	err := cache.SetWithExpiration(ctx, "cache:expensive_query", result, 5*time.Minute)
//
//	// No expiration (equivalent to Set)
//	err := cache.SetWithExpiration(ctx, "permanent_config", config, 0)
func (r *RedisCache) SetWithExpiration(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	err := r.client.Set(ctx, key, value, expiration).Err()
	if err != nil {
		return err
	}
	return nil
}

// Del removes one or more keys from Redis atomically.
// This method uses Redis DEL command which can delete multiple keys in a single operation.
// The operation is atomic for single keys and uses Redis transaction semantics for multiple keys.
//
// Deletion behavior:
//   - Non-existent keys are silently ignored (no error returned)
//   - Returns success even if some/all keys don't exist
//   - Multiple keys are deleted in a single Redis command for efficiency
//   - Operation is atomic - either all keys are processed or none
//
// Performance considerations:
//   - Deleting multiple keys is more efficient than individual Del calls
//   - Very large key lists may impact Redis performance
//   - Consider using DelWithPattern for pattern-based bulk deletion
//
// Parameters:
//   - ctx: Context for timeout control and request cancellation
//   - keys: Variable number of Redis keys to delete
//
// Returns:
//   - error: Redis connection error or command execution error
//
// Examples:
//
//	// Delete single key
//	err := cache.Del(ctx, "user:123")
//
//	// Delete multiple keys atomically
//	err := cache.Del(ctx, "user:123", "user:456", "user:789")
//
//	// Delete session and related data
//	err := cache.Del(ctx, "session:abc", "session_data:abc", "session_meta:abc")
//
//	// Safe to call with non-existent keys
//	err := cache.Del(ctx, "might_not_exist") // No error if key doesn't exist
func (r *RedisCache) Del(ctx context.Context, keys ...string) error {
	err := r.client.Del(ctx, keys...).Err()
	if err != nil {
		return err
	}
	return nil
}

// DelWithPattern deletes all Redis keys matching the specified pattern.
// This method combines the Keys and Del operations to provide pattern-based bulk deletion.
// It's a convenience method for cleaning up multiple related keys at once.
//
// Operation steps:
//  1. Uses Keys() to find all keys matching the pattern
//  2. If no keys found, returns successfully without Redis calls
//  3. Uses Del() to remove all found keys atomically
//  4. Returns any errors from either operation
//
// Performance and safety considerations:
//   - Can be expensive on large databases due to Keys() operation
//   - May block Redis during execution with large result sets
//   - Use specific patterns to limit scope and improve performance
//   - Be extremely careful with broad patterns like "*"
//   - Consider using Redis SCAN-based approaches for large datasets
//
// Pattern matching uses same rules as Keys():
//   - '*' matches any number of characters
//   - '?' matches exactly one character
//   - Character classes and ranges supported
//
// Parameters:
//   - ctx: Context for timeout control and request cancellation
//   - pattern: Glob-style pattern to match keys for deletion
//
// Returns:
//   - error: Pattern matching error, key deletion error, or Redis connection error
//
// Examples:
//
//	// Delete all user sessions
//	err := cache.DelWithPattern(ctx, "session:user:*")
//
//	// Clean up temporary data
//	err := cache.DelWithPattern(ctx, "temp:*")
//
//	// Remove cache entries for specific date
//	err := cache.DelWithPattern(ctx, "cache:2023-12-*")
//
//	// Clean up test data
//	err := cache.DelWithPattern(ctx, "test:*")
//
// Warning: Use patterns carefully in production:
//
//	cache.DelWithPattern(ctx, "*") // DANGEROUS: Deletes ALL keys!
func (r *RedisCache) DelWithPattern(ctx context.Context, pattern string) error {
	keys, err := r.Keys(ctx, pattern)
	if err != nil {
		return err
	}
	if len(keys) == 0 {
		return nil
	}
	err = r.Del(ctx, keys...)
	if err != nil {
		return err
	}
	return nil
}

// Close gracefully shuts down the Redis connection and releases all associated resources.
// This method should be called when the cache instance is no longer needed, typically
// during application shutdown or when disposing of cache instances.
//
// Cleanup operations performed:
//   - Closes the underlying Redis client connection
//   - Terminates any background goroutines managed by the client
//   - Releases connection pool resources
//   - Frees memory allocated for client state
//
// Behavior after closing:
//   - Subsequent cache operations will return connection errors
//   - Multiple calls to Close() are safe (idempotent)
//   - The cache instance becomes unusable after closing
//
// Connection pool considerations:
//   - If using a shared Redis client, closing may affect other users
//   - For applications with multiple cache instances, consider connection sharing
//   - Connection pools are properly drained before closing
//
// Returns:
//   - error: Connection close error (rare, usually indicates network issues)
//
// Example usage patterns:
//
//	// Basic cleanup
//	defer cache.Close()
//
//	// With error handling
//	defer func() {
//	    if err := cache.Close(); err != nil {
//	        log.Printf("Error closing Redis cache: %v", err)
//	    }
//	}()
//
//	// In application shutdown
//	func (app *App) Shutdown() {
//	    if err := app.cache.Close(); err != nil {
//	        app.logger.Error("Failed to close cache", "error", err)
//	    }
//	}
func (r *RedisCache) Close() error {
	return r.client.Close()
}
