package main

import (
	"fmt"
	"sync"
	"time"
)

const defaultCalue = 100

type Cache struct {
	mu    sync.Mutex
	items map[int]int
}

func NewCache() *Cache {
	m := make(map[int]int)
	return &Cache{
		items: m,
	}
}

func (c *Cache) Set(key int, value int) {
	c.mu.Lock()
	c.items[key] = value
	c.mu.Unlock()
}

func (c *Cache) Get(key int) int {
	c.mu.Lock()
	v, ok := c.items[key]
	c.mu.Unlock()
	if ok {
		fmt.Println("cache hit")
		return v
	}

	// 非同期にキャッシュを更新する
	go func() {
		v = HeavyGet(key)
		c.Set(key, v)
	}()

	// とりあえずデフォルト値を返す
	return defaultCalue
}

func HeavyGet(key int) int {
	fmt.Println("HeavyGet")
	time.Sleep(2 * time.Second)
	value := key * 2
	return value
}

type user struct {
	userID int
}

func main() {
	mCache := NewCache()
	fmt.Println(mCache.Get(1))
	time.Sleep(3 * time.Second)
	fmt.Println(mCache.Get(1))

}
