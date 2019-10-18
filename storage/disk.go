package storage

import (
	"io"
	"os"
)

// diskStorage disk storage
type diskStorage struct {
}

func newDiskStorage() *diskStorage {
	return &diskStorage{}
}

// Open open path
func (l *diskStorage) Open(path string) (io.ReadWriteCloser, error) {
	return os.OpenFile(path, os.O_RDWR|os.O_CREATE, 0644)
}
