package cache

import (
	"fmt"
	"log"
	"sync"
	"time"
)

type Cache struct {
	lock sync.RWMutex
	data map[string][]byte
}

func New() *Cache {
	return &Cache{
		data: make(map[string][]byte),
	}
}

func (c *Cache) Set(k []byte, v []byte, ttl time.Duration) error {
	c.lock.Lock()
	defer c.lock.Unlock()

	c.data[string(k)] = v
	log.Printf("SET %s to %s\n", string(k), string(v))

	//ticker := time.NewTicker(ttl)
	if ttl > 0 {
		go func() {
			<-time.After(ttl)
			delete(c.data, string(k))
		}()
	}

	return nil
}

func (c *Cache) Get(k []byte) ([]byte, error) {
	c.lock.RLock()
	defer c.lock.RUnlock()

	v, ok := c.data[string(k)]
	if !ok {
		return []byte{}, fmt.Errorf("there is no value for this key: %s", string(k))
	}

	//log.Printf("GET %s -> %s\n", string(k), string(v))

	return v, nil
}

func (c *Cache) Has(k []byte) bool {
	c.lock.RLock()
	defer c.lock.RUnlock()

	_, ok := c.data[string(k)]
	if !ok {
		return false
	}

	return true
}

func (c *Cache) Delete(k []byte) error {
	c.lock.Lock()
	defer c.lock.Unlock()

	if _, ok := c.data[string(k)]; !ok {
		return fmt.Errorf("there is no value for this key: %s", string(k))
	}

	delete(c.data, string(k))

	return nil
}
