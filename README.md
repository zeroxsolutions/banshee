# Banshee Cache Library

[![Go Version](https://img.shields.io/badge/Go-%3E%3D%201.18-blue.svg)](https://golang.org/)
[![License](https://img.shields.io/badge/License-MIT-green.svg)](LICENSE)
[![Version](https://img.shields.io/badge/Version-v0.0.1-orange.svg)](https://github.com/zeroxsolutions/banshee/releases/tag/v0.0.1)

A high-performance, production-ready Go cache library providing a unified interface for different cache backends. Banshee offers Redis implementation with comprehensive testing support through mock implementations.

## 🚀 Features

- **Unified Cache Interface**: Clean, consistent API for all cache operations
- **Redis Backend**: Production-ready Redis implementation with full feature support
- **Mock Implementation**: Complete mock for unit testing without external dependencies
- **Context Support**: All operations support context for timeout and cancellation
- **Pattern Operations**: Advanced pattern-based key operations (Keys, DelWithPattern)
- **TTL Support**: Flexible expiration handling with automatic cleanup
- **Thread-Safe**: All implementations are safe for concurrent use
- **Comprehensive Testing**: Extensive test coverage for both unit and integration tests

## 📦 Installation

### Prerequisites

- Go 1.18 or higher
- Redis server (for Redis implementation)

### Install the Package

1. Install the library:
```bash
go get github.com/zeroxsolutions/banshee
```

2. Install specific implementations:
```bash
# For Redis implementation
go get github.com/zeroxsolutions/banshee/redis

# For mock implementation (testing)
go get github.com/zeroxsolutions/banshee/mock
```

## 🔧 Quick Start

### Redis Cache Example

```go
package main

import (
    "context"
    "fmt"
    "log"
    "time"

    "github.com/zeroxsolutions/alex"
    "github.com/zeroxsolutions/banshee/redis"
)

func main() {
    // Configure Redis connection
    config := &alex.RedisConfig{
        Addr:     "localhost:6379",
        Password: "", // no password
        DB:       0,  // default database
    }

    // Create cache instance
    cache, err := redis.NewRedisCache(config)
    if err != nil {
        log.Fatal("Failed to connect to Redis:", err)
    }
    defer cache.Close()

    ctx := context.Background()

    // Basic operations
    err = cache.Set(ctx, "user:123", "john_doe")
    if err != nil {
        log.Fatal("Set failed:", err)
    }

    value, err := cache.Get(ctx, "user:123")
    if err != nil {
        log.Fatal("Get failed:", err)
    }
    fmt.Printf("Retrieved: %s\n", value)

    // Set with expiration
    err = cache.SetWithExpiration(ctx, "session:abc", "session_data", 30*time.Minute)
    if err != nil {
        log.Fatal("SetWithExpiration failed:", err)
    }

    // Pattern operations
    keys, err := cache.Keys(ctx, "user:*")
    if err != nil {
        log.Fatal("Keys failed:", err)
    }
    fmt.Printf("Found keys: %v\n", keys)
}
```

### Mock Testing Example

```go
package main

import (
    "context"
    "testing"

    "github.com/stretchr/testify/mock"
    "github.com/zeroxsolutions/banshee/mock"
)

func TestUserService(t *testing.T) {
    // Create mock cache
    mockCache := mock.NewMockCache(t)
    
    // Set expectations
    mockCache.On("Get", mock.Anything, "user:123").Return("john_doe", nil)
    mockCache.On("Set", mock.Anything, "user:123", "jane_doe").Return(nil)
    
    // Test your service with the mock
    service := NewUserService(mockCache)
    
    // Your test logic here
    user, err := service.GetUser(context.Background(), "123")
    assert.NoError(t, err)
    assert.Equal(t, "john_doe", user)
    
    // Expectations are automatically verified when test ends
}
```

## 📚 API Reference

### Cache Interface

The core cache interface provides these methods:

```go
type Cache interface {
    // Connection management
    IsConnected(ctx context.Context) bool
    Close() error
    
    // Basic operations
    Get(ctx context.Context, key string) (string, error)
    Set(ctx context.Context, key string, value interface{}) error
    SetWithExpiration(ctx context.Context, key string, value interface{}, expiration time.Duration) error
    
    // Deletion operations
    Del(ctx context.Context, keys ...string) error
    DelWithPattern(ctx context.Context, pattern string) error
    
    // Pattern operations
    Keys(ctx context.Context, pattern string) ([]string, error)
}
```

### Available Implementations

| Implementation | Package | Use Case |
|----------------|---------|----------|
| **Redis** | `github.com/zeroxsolutions/banshee/redis` | Production caching with Redis backend |
| **Mock** | `github.com/zeroxsolutions/banshee/mock` | Unit testing without external dependencies |

### Pattern Syntax

Redis pattern matching supports:
- `*` - matches zero or more characters
- `?` - matches exactly one character  
- `[abc]` - matches any character in the set
- `[^abc]` - matches any character not in the set
- `[a-z]` - matches any character in the range

Examples:
```go
keys, _ := cache.Keys(ctx, "user:*")        // All user keys
keys, _ := cache.Keys(ctx, "session:???")   // 3-character session IDs
keys, _ := cache.Keys(ctx, "cache:[0-9]*")  // Numbered cache entries
```

## 🧪 Testing

### Environment Variables

Configure the following environment variables for testing:

| Variable | Description | Default |
|----------|-------------|---------|
| `REDIS_ADDRESS` | Redis server address | `localhost:6379` |
| `REDIS_PASSWORD` | Redis authentication password | _(empty)_ |
| `REDIS_DB` | Redis database number | `0` |

### Running Tests

1. **Install dependencies**:
```bash
go mod tidy
```

2. **Unit tests only** (using mocks):
```bash
go test -v ./mock/...
```

3. **Integration tests** (requires Redis):
```bash
# Start Redis server first
redis-server

# Run Redis integration tests
go test -v ./redis/...
```

4. **All tests**:
```bash
go test -v ./...
```

5. **Test with coverage**:
```bash
go test -v -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

### Test Script

Use the provided test script for convenience:
```bash
chmod +x bin/test.sh
./bin/test.sh
```

## 🏗️ Architecture

```
banshee/
├── cache.go              # Core interface (if created)
├── redis/
│   ├── redis_cache.go    # Redis implementation
│   └── redis_cache_test.go
├── mock/
│   ├── mock_cache.go     # Mock implementation
│   └── mock_cache_test.go
└── bin/
    └── test.sh           # Test runner script
```

## 📊 Performance

### Redis Implementation
- **Latency**: Sub-millisecond operations for local Redis
- **Throughput**: Supports thousands of operations per second
- **Memory**: Efficient Redis protocol usage
- **Concurrency**: Thread-safe with connection pooling

### Benchmarks
```bash
# Run benchmarks
go test -bench=. -benchmem ./...
```

## 🔒 Security Considerations

1. **Redis Security**:
   - Use authentication in production (`REDIS_PASSWORD`)
   - Configure Redis with appropriate security settings
   - Use TLS for network connections in production

2. **Pattern Operations**:
   - Be cautious with broad patterns like `*`
   - `DelWithPattern(ctx, "*")` will delete ALL keys
   - Use specific patterns to limit scope

3. **Input Validation**:
   - Validate keys and values before caching
   - Sanitize user input used in patterns

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch: `git checkout -b feature/amazing-feature`
3. Commit your changes: `git commit -m 'Add amazing feature'`
4. Push to the branch: `git push origin feature/amazing-feature`
5. Open a Pull Request

### Development Guidelines

- Follow Go conventions and best practices
- Add comprehensive tests for new features
- Update documentation for API changes
- Run `go fmt` and `go vet` before committing
- Ensure all tests pass before submitting PR

## 📝 Changelog

### v0.0.1 (Initial Release)
- ✅ Core cache interface design
- ✅ Redis implementation with full feature support
- ✅ Mock implementation for testing
- ✅ Comprehensive test coverage
- ✅ Context support for all operations
- ✅ Pattern-based operations (Keys, DelWithPattern)
- ✅ TTL support with automatic expiration
- ✅ Thread-safe implementations
- ✅ Documentation and examples

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🆘 Support

- **Documentation**: Check the GoDoc comments in source files
- **Issues**: Report bugs on [GitHub Issues](https://github.com/zeroxsolutions/banshee/issues)
- **Discussions**: Join discussions on [GitHub Discussions](https://github.com/zeroxsolutions/banshee/discussions)

## 🔗 Related Projects

- [Redis](https://redis.io/) - The Redis data structure server
- [go-redis](https://github.com/redis/go-redis) - Redis client for Go
- [testify](https://github.com/stretchr/testify) - Testing toolkit for Go

---

**Made with ❤️ by ZeroX Solutions**
