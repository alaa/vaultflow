package cache

import (
	"io/ioutil"
	"os"
)

type Cache struct {
}

func New() (*Cache, error) {
	cache := &Cache{}
	if _, err := os.Stat(".cache/"); os.IsNotExist(err) {
		os.Mkdir(".cache/", 0777)
		if _, err := os.Stat(".cache/revision"); os.IsNotExist(err) {
			return cache, cache.UpdateRevision("")
		}
	}
	return cache, nil
}

func (c *Cache) Put(key, value []byte) error {
	cacheDir := ".cache/" + string(key)
	return ioutil.WriteFile(cacheDir, value, 0644)
}

func (c *Cache) Get(key []byte) ([]byte, error) {
	cacheDir := ".cache/" + string(key)
	return ioutil.ReadFile(cacheDir)
}

// Get Revision ID
func (c *Cache) GetRevision() (string, error) {
	b, err := c.Get([]byte("revision"))
	return string(b), err
}

// Update Revision ID
func (c *Cache) UpdateRevision(id string) error {
	return c.Put([]byte("revision"), []byte(id))
}
