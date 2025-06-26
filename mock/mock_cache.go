// Package mock provides a mock implementation of the cache interface for testing purposes.
// It uses the testify/mock framework to enable controlled testing scenarios with
// configurable return values and behavior verification.
package mock

import (
	"context"
	"time"

	"github.com/stretchr/testify/mock"
	aliasCache "github.com/zeroxsolutions/barbatos/cache"
)

// MockCache is a mock implementation of the Cache interface designed for testing.
// It leverages the testify/mock package to simulate cache interactions and enables
// setting expectations on method calls, return values, and call verification.
//
// This mock is particularly useful for:
//   - Unit testing cache-dependent code without actual cache infrastructure
//   - Simulating cache failures and edge cases
//   - Verifying cache interaction patterns in application logic
//   - Performance testing without I/O overhead
//
// Example usage:
//
//	mockCache := mock.NewMockCache(t)
//	mockCache.On("Get", mock.Anything, "key").Return("value", nil)
//	// Use mockCache in your tests
type MockCache struct {
	mock.Mock
}

// IsConnected mocks the cache connectivity check method.
// This method simulates checking the connection status to the cache system
// and allows tests to control whether the cache appears connected or not.
//
// The mock can be configured to return different values for different test scenarios:
//   - Return true to simulate a healthy cache connection
//   - Return false to simulate connection failures or unavailable cache
//   - Use function-based returns for dynamic behavior based on context
//
// Parameters:
//   - ctx: Context for request lifecycle management and cancellation
//
// Returns:
//   - bool: Mocked connection status (true if connected, false otherwise)
//
// Example:
//
//	mockCache.On("IsConnected", mock.Anything).Return(true)
//	connected := mockCache.IsConnected(ctx) // returns true
func (m *MockCache) IsConnected(ctx context.Context) bool {
	ret := m.Called(ctx)
	var r0 bool
	if rf, ok := ret.Get(0).(func(context.Context) bool); ok {
		r0 = rf(ctx)
	} else {
		r0 = ret.Bool(0)
	}
	return r0
}

// Keys mocks the pattern-based key retrieval method.
// This method simulates retrieving keys that match a given pattern from the cache.
// It allows tests to control which keys are returned for specific patterns.
//
// The mock supports various return scenarios:
//   - Return a predefined list of keys matching the pattern
//   - Return an empty slice to simulate no matches
//   - Return an error to simulate retrieval failures
//   - Use function-based returns for dynamic key generation
//
// Parameters:
//   - ctx: Context for request lifecycle management and cancellation
//   - pattern: Pattern string to match keys (format depends on cache implementation)
//
// Returns:
//   - []string: Slice of mocked keys matching the pattern
//   - error: Mocked error if the operation should fail
//
// Example:
//
//	mockCache.On("Keys", mock.Anything, "user:*").Return([]string{"user:1", "user:2"}, nil)
//	keys, err := mockCache.Keys(ctx, "user:*") // returns ["user:1", "user:2"], nil
func (m *MockCache) Keys(ctx context.Context, pattern string) ([]string, error) {
	ret := m.Called(ctx, pattern)
	var r0 []string
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string) ([]string, error)); ok {
		return rf(ctx, pattern)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string) []string); ok {
		r0 = rf(ctx, pattern)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]string)
		}
	}
	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, pattern)
	} else {
		r1 = ret.Error(1)
	}
	return r0, r1
}

// Get mocks the value retrieval method for cache keys.
// This method simulates retrieving a value from the cache based on a provided key.
// It enables tests to control what values are returned for specific keys.
//
// The mock supports various return scenarios:
//   - Return a specific value for a given key
//   - Return an empty string and error to simulate key not found
//   - Return different values for different keys
//   - Use function-based returns for dynamic value generation
//
// Parameters:
//   - ctx: Context for request lifecycle management and cancellation
//   - key: Cache key to retrieve the value for
//
// Returns:
//   - string: Mocked value associated with the key
//   - error: Mocked error if the operation should fail (e.g., key not found)
//
// Example:
//
//	mockCache.On("Get", mock.Anything, "user:123").Return("john_doe", nil)
//	mockCache.On("Get", mock.Anything, "missing").Return("", cache.ErrCacheNil)
//	value, err := mockCache.Get(ctx, "user:123") // returns "john_doe", nil
func (m *MockCache) Get(ctx context.Context, key string) (string, error) {
	ret := m.Called(ctx, key)
	var r0 string
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string) (string, error)); ok {
		return rf(ctx, key)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string) string); ok {
		r0 = rf(ctx, key)
	} else {
		r0 = ret.Get(0).(string)
	}
	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, key)
	} else {
		r1 = ret.Error(1)
	}
	return r0, r1
}

// Set mocks the value storage method for cache keys.
// This method simulates storing a value in the cache under a specified key.
// It allows tests to verify that values are stored correctly and handle storage failures.
//
// The mock supports various return scenarios:
//   - Return nil to simulate successful storage
//   - Return an error to simulate storage failures
//   - Use function-based returns for dynamic behavior based on input
//   - Verify that the correct key-value pairs are being stored
//
// Parameters:
//   - ctx: Context for request lifecycle management and cancellation
//   - key: Cache key to store the value under
//   - value: Value to be stored (can be any type supported by the cache)
//
// Returns:
//   - error: Mocked error if the storage operation should fail
//
// Example:
//
//	mockCache.On("Set", mock.Anything, "user:123", "john_doe").Return(nil)
//	mockCache.On("Set", mock.Anything, "readonly", mock.Anything).Return(errors.New("read-only"))
//	err := mockCache.Set(ctx, "user:123", "john_doe") // returns nil
func (m *MockCache) Set(ctx context.Context, key string, value interface{}) error {
	ret := m.Called(ctx, key, value)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string, interface{}) error); ok {
		r0 = rf(ctx, key, value)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// SetWithExpiration mocks the value storage method with automatic expiration.
// This method simulates storing a value in the cache with a time-to-live (TTL) setting.
// It allows tests to verify TTL behavior and handle expiration-related scenarios.
//
// The mock supports various return scenarios:
//   - Return nil to simulate successful storage with expiration
//   - Return an error to simulate storage failures
//   - Use function-based returns for dynamic behavior
//   - Verify that correct expiration times are being set
//
// Parameters:
//   - ctx: Context for request lifecycle management and cancellation
//   - key: Cache key to store the value under
//   - value: Value to be stored (can be any type supported by the cache)
//   - expiration: Duration after which the key should automatically expire
//
// Returns:
//   - error: Mocked error if the storage operation should fail
//
// Example:
//
//	mockCache.On("SetWithExpiration", mock.Anything, "session:abc", sessionData, time.Hour).Return(nil)
//	err := mockCache.SetWithExpiration(ctx, "session:abc", sessionData, time.Hour) // returns nil
func (m *MockCache) SetWithExpiration(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	ret := m.Called(ctx, key, value, expiration)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string, interface{}, time.Duration) error); ok {
		r0 = rf(ctx, key, value, expiration)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Del mocks the key deletion method for one or more cache keys.
// This method simulates removing keys from the cache and allows tests to verify
// deletion behavior and handle deletion failures.
//
// The mock supports various return scenarios:
//   - Return nil to simulate successful deletion
//   - Return an error to simulate deletion failures
//   - Use function-based returns for dynamic behavior
//   - Verify that the correct keys are being deleted
//
// Note: The mock handles variadic arguments by converting them to []interface{}
// for compatibility with the testify/mock framework.
//
// Parameters:
//   - ctx: Context for request lifecycle management and cancellation
//   - keys: Variable number of keys to delete from the cache
//
// Returns:
//   - error: Mocked error if the deletion operation should fail
//
// Example:
//
//	mockCache.On("Del", mock.Anything, "user:123", "user:456").Return(nil)
//	err := mockCache.Del(ctx, "user:123", "user:456") // returns nil
func (m *MockCache) Del(ctx context.Context, keys ...string) error {
	_keys := make([]interface{}, len(keys))
	for _idx := range keys {
		_keys[_idx] = keys[_idx]
	}
	var _args []interface{}
	_args = append(_args, ctx)
	_args = append(_args, _keys...)
	ret := m.Called(_args...)
	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, ...string) error); ok {
		r0 = rf(ctx, keys...)
	} else {
		r0 = ret.Error(0)
	}
	return r0
}

// DelWithPattern mocks the pattern-based key deletion method.
// This method simulates deleting all keys that match a specified pattern.
// It allows tests to verify bulk deletion behavior and handle pattern deletion failures.
//
// The mock supports various return scenarios:
//   - Return nil to simulate successful pattern-based deletion
//   - Return an error to simulate deletion failures
//   - Use function-based returns for dynamic behavior
//   - Verify that the correct patterns are being used for deletion
//
// Parameters:
//   - ctx: Context for request lifecycle management and cancellation
//   - pattern: Pattern string to match keys for deletion
//
// Returns:
//   - error: Mocked error if the deletion operation should fail
//
// Example:
//
//	mockCache.On("DelWithPattern", mock.Anything, "session:*").Return(nil)
//	err := mockCache.DelWithPattern(ctx, "session:*") // returns nil
func (m *MockCache) DelWithPattern(ctx context.Context, pattern string) error {
	ret := m.Called(ctx, pattern)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string) error); ok {
		r0 = rf(ctx, pattern)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Close mocks the cache connection cleanup method.
// This method simulates closing the connection to the cache system and releasing
// associated resources. It allows tests to verify proper cleanup behavior.
//
// The mock supports various return scenarios:
//   - Return nil to simulate successful connection closure
//   - Return an error to simulate close failures
//   - Use function-based returns for dynamic behavior
//   - Verify that Close is called when expected
//
// Returns:
//   - error: Mocked error if the close operation should fail
//
// Example:
//
//	mockCache.On("Close").Return(nil)
//	err := mockCache.Close() // returns nil
func (m *MockCache) Close() error {
	ret := m.Called()

	var r0 error
	if rf, ok := ret.Get(0).(func() error); ok {
		r0 = rf()
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// NewMockCache creates and configures a new MockCache instance for testing.
// This constructor sets up the mock with proper test integration and automatic
// expectation verification when the test completes.
//
// The function configures the mock to:
//   - Integrate with the testing framework for proper error reporting
//   - Automatically verify all expectations when the test finishes
//   - Clean up resources properly after test completion
//
// Parameters:
//   - t: Testing interface that supports both mock.TestingT and cleanup functionality
//
// Returns:
//   - aliasCache.Cache: A new MockCache instance implementing the Cache interface
//
// Example:
//
//	func TestMyFunction(t *testing.T) {
//	    mockCache := mock.NewMockCache(t)
//	    mockCache.On("Get", mock.Anything, "key").Return("value", nil)
//
//	    // Use mockCache in your test
//	    result := myFunction(mockCache)
//
//	    // Expectations are automatically verified when test ends
//	}
func NewMockCache(t interface {
	mock.TestingT
	Cleanup(func())
}) aliasCache.Cache {
	m := &MockCache{}
	m.Mock.Test(t)

	t.Cleanup(func() { m.AssertExpectations(t) })
	return m
}
