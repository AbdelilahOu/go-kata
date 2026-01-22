package main

import (
	"encoding/binary"
	"fmt"
	"hash/fnv"
	"sync"
)

type ShardedMap2[K comparable, V any] struct {
	shards    []*Bucket[K, V]
	numShards uint64
}

type Bucket[K comparable, V any] struct {
	sync.RWMutex
	items map[K]V
}

func NewShardedMap2[K comparable, V any](numShards uint64) *ShardedMap2[K, V] {
	if numShards == 0 {
		numShards = 1
	}

	sm := &ShardedMap2[K, V]{
		shards:    make([]*Bucket[K, V], numShards),
		numShards: numShards,
	}

	for i := range sm.shards {
		sm.shards[i] = &Bucket[K, V]{items: make(map[K]V)}
	}

	return sm
}

func (sm *ShardedMap2[K, V]) getShardIndex(key K) uint64 {
	h := fnv.New64a()

	switch k := any(key).(type) {
	case string:
		h.Write([]byte(k))
	case int:
		var buf [8]byte
		binary.LittleEndian.PutUint64(buf[:], uint64(k))
		h.Write(buf[:])
	case uint64:
		var buf [8]byte
		binary.LittleEndian.PutUint64(buf[:], k)
		h.Write(buf[:])
	default:
		s := fmt.Sprint(key)
		h.Write([]byte(s))
	}
	hash := h.Sum64()
	return hash % sm.numShards
}

func (sm *ShardedMap2[K, V]) Get(key K) (V, bool) {
	idx := sm.getShardIndex(key)
	shard := sm.shards[idx]
	shard.RLock()
	val, ok := shard.items[key]
	shard.RUnlock()
	return val, ok
}

func (sm *ShardedMap2[K, V]) Set(key K, value V) {
	idx := sm.getShardIndex(key)
	shard := sm.shards[idx]
	shard.Lock()
	shard.items[key] = value
	shard.Unlock()
}

func (sm *ShardedMap2[K, V]) Delete(key K) {
	idx := sm.getShardIndex(key)
	shard := sm.shards[idx]
	shard.Lock()
	delete(shard.items, key)
	shard.Unlock()
}

func (sm *ShardedMap2[K, V]) Keys() []K {
	keys := make([]K, 0, 100)
	for i := uint64(0); i < sm.numShards; i++ {
		shard := sm.shards[i]
		shard.RLock()
		for k := range shard.items {
			keys = append(keys, k)
		}
		shard.RUnlock()
	}
	return keys
}

func main() {
	// Choose which one to test:
	// sm := NewShardedMap[string, int](16)     // parallel slices version
	sm := NewShardedMap2[string, int](16) // embedded Bucket version

	// Write some values
	sm.Set("alice", 42)
	sm.Set("bob", 19)
	sm.Set("charlie", 7)
	sm.Set("dave", 100)
	sm.Set("eve", 33)

	// Read them back
	if v, ok := sm.Get("alice"); ok {
		fmt.Printf("alice → %d\n", v)
	}
	if v, ok := sm.Get("bob"); ok {
		fmt.Printf("bob → %d\n", v)
	}
	if v, ok := sm.Get("nonexistent"); ok {
		fmt.Printf("should not see this: %d\n", v)
	} else {
		fmt.Println("nonexistent → not found (correct)")
	}

	// Delete one
	sm.Delete("charlie")
	if _, ok := sm.Get("charlie"); !ok {
		fmt.Println("charlie was deleted successfully")
	}

	// List all keys
	fmt.Println("\nAll keys:")
	for _, k := range sm.Keys() {
		fmt.Println(" -", k)
	}
}
