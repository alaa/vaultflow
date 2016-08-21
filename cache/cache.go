package cache

import (
	"io/ioutil"
	"os"
	"path/filepath"
)

const revision = "revision"

type Cache struct {
	basePath string
}

func New(directory string) (*Cache, error) {
	revision := filepath.Join(directory, revision)
	cache := &Cache{basePath: directory}

	if _, err := os.Stat(directory); os.IsNotExist(err) {
		os.Mkdir(directory, 0777)
		revisionPath := filepath.Join(cache.basePath, revision)
		if _, err := os.Stat(revisionPath); os.IsNotExist(err) {
			return cache, cache.UpdateRevision("")
		}
	}

	return cache, nil
}

func (c *Cache) Put(key string, value []byte) error {
	cacheDir := filepath.Join(c.basePath, key)
	return ioutil.WriteFile(cacheDir, value, 0644)
}

func (c *Cache) Get(key string) ([]byte, error) {
	cacheDir := filepath.Join(c.basePath, key)
	return ioutil.ReadFile(cacheDir)
}

// Get Revision ID
func (c *Cache) GetRevision() (string, error) {
	b, err := c.Get(revision)
	return string(b), err
}

// Update Revision ID
func (c *Cache) UpdateRevision(id string) error {
	return c.Put(revision, []byte(id))
}
