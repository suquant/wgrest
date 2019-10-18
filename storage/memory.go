package storage

import (
	"fmt"
	"io"

	"github.com/patrickmn/go-cache"
)

type memoryStorage struct {
	cache *cache.Cache
}

type memoryDescriptor struct {
	cache  *cache.Cache
	key    string
	rIndex int
	wIndex int
	closed bool
}

func newMemoryStorage() *memoryStorage {
	return &memoryStorage{
		cache: cache.New(cache.NoExpiration, cache.NoExpiration),
	}
}

// Open open path
func (m *memoryStorage) Open(path string) (io.ReadWriteCloser, error) {

	return &memoryDescriptor{
		cache: m.cache,
		key:   path,
	}, nil
}

func (d *memoryDescriptor) Read(p []byte) (int, error) {
	if d.closed {
		return 0, ErrClosed
	}

	res, found := d.cache.Get(d.key)
	if found {
		if resBytes, ok := res.([]byte); ok {
			if d.rIndex >= len(resBytes) {
				return 0, io.EOF
			}

			n := copy(p, resBytes[d.rIndex:])
			d.rIndex += n
			return n, nil
		}

		return 0, fmt.Errorf("%s contain inconsist data", d.key)
	}

	return 0, fmt.Errorf("%s not found", d.key)
}

func (d *memoryDescriptor) Write(p []byte) (int, error) {
	if d.closed {
		return 0, ErrClosed
	}

	d.cache.SetDefault(d.key, p)

	return len(p), nil
}

func (d *memoryDescriptor) Close() error {
	d.closed = true

	return nil
}
