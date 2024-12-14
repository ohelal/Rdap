package cmd

import (
    "fmt"
    "time"
)

// CacheStats represents cache statistics
type CacheStats struct {
    TotalItems int
    Items      map[string]interface{}
}

var cache map[string]interface{}

func init() {
    cache = make(map[string]interface{})
}

func getCachedResult(key string) (interface{}, bool) {
    if result, ok := cache[key]; ok {
        fmt.Printf("%s\n", successStyle("Cache hit"))
        return result, true
    }
    return nil, false
}

func cacheResult(key string, value interface{}, duration time.Duration) {
    cache[key] = value
    fmt.Printf("%s\n", successStyle("Result cached"))
}

func clearCache() {
    cache = make(map[string]interface{})
    fmt.Printf("%s\n", successStyle("Cache cleared"))
}

func getCacheStats() CacheStats {
    return CacheStats{
        TotalItems: len(cache),
        Items:      cache,
    }
}
