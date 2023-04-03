package main

import (
	"fmt"
	"sync"
	"time"

	"golang.org/x/sync/singleflight"
)

var group singleflight.Group

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

	// singleflight で同時に同じキーで呼ばれた場合は一つだけ処理を実行する
	vv, err, _ := group.Do(fmt.Sprintf("cacheGet_%d", key), func() (interface{}, error) {
		value := HeavyGet(key)
		v = HeavyGet(key)
		c.Set(key, v)
		return value, nil
	})

	if err != nil {
		panic(err)
	}

	// interface{} 型なのでint にキャスト
	return vv.(int)
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
	fmt.Println(mCache.Get(1))
	fmt.Println(mCache.Get(1))
	fmt.Println(mCache.Get(1))

}
