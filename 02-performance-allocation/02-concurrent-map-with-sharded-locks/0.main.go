package main

import (
	"encoding/binary"
	"fmt"
	"hash/fnv"
	"sync"
	"time"
)

type ShardedMap[K comparable, V any] struct {
	shards    []map[K]V
	locks     []sync.RWMutex
	numShards uint64
}

func NewShardedMap[K comparable, V any](numShards uint64) *ShardedMap[K, V] {
	if numShards == 0 {
		numShards = 1
	}

	sm := &ShardedMap[K, V]{
		shards:    make([]map[K]V, numShards),
		locks:     make([]sync.RWMutex, numShards),
		numShards: numShards,
	}

	for i := range sm.shards {
		sm.shards[i] = make(map[K]V)
	}

	return sm
}

func (sm *ShardedMap[K, V]) getShardIndex(key K) uint64 {
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

func (sm *ShardedMap[K, V]) Get(key K) (V, bool) {
	idx := sm.getShardIndex(key)
	sm.locks[idx].RLock()
	val, ok := sm.shards[idx][key]
	sm.locks[idx].RUnlock()
	return val, ok
}

func (sm *ShardedMap[K, V]) Set(key K, value V) {
	idx := sm.getShardIndex(key)
	sm.locks[idx].Lock()
	sm.shards[idx][key] = value
	sm.locks[idx].Unlock()
}

func (sm *ShardedMap[K, V]) Delete(key K) {
	idx := sm.getShardIndex(key)
	sm.locks[idx].Lock()
	delete(sm.shards[idx], key)
	sm.locks[idx].Unlock()
}

func (sm *ShardedMap[K, V]) Keys() []K {
	keys := make([]K, 0, 100)
	for i := uint64(0); i < sm.numShards; i++ {
		sm.locks[i].RLock()
		for k := range sm.shards[i] {
			keys = append(keys, k)
		}
		sm.locks[i].RUnlock()
	}
	return keys
}

func runTest(sm interface {
	Set(string, int)
	Get(string) (int, bool)
}, name string, workers, ops int) time.Duration {
	var wg sync.WaitGroup
	wg.Add(workers)

	start := time.Now()

	for w := 0; w < workers; w++ {
		go func() {
			defer wg.Done()
			for i := 0; i < ops; i++ {
				key := fmt.Sprintf("k%d", i%5000) // 5000 hot keys → high contention possible
				if i%7 == 0 {
					sm.Set(key, i)
				} else {
					_, _ = sm.Get(key)
				}
			}
		}()
	}

	wg.Wait()
	return time.Since(start)
}

func main() {
	const workers = 16
	const opsPerWorker = 100_000

	fmt.Println("=== 1 shard (high contention expected) ===")
	sm1 := NewShardedMap[string, int](1)
	d1 := runTest(sm1, "1-shard", workers, opsPerWorker)
	fmt.Printf("1 shard → %v\n", d1)

	fmt.Println("\n=== 64 shards (should be much faster) ===")
	sm64 := NewShardedMap[string, int](64)
	d64 := runTest(sm64, "64-shard", workers, opsPerWorker)
	fmt.Printf("64 shards → %v\n", d64)

	speedup := float64(d1) / float64(d64)
	fmt.Printf("\nSpeedup with 64 shards: ≈ %.1fx faster\n", speedup)
}
