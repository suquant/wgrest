package storage

import (
	"io"
)

// Storage storage interface
type Storage interface {
	Open(path string) (io.ReadWriteCloser, error)
}

// Type storage type
type Type int

const (
	// DiskStorage disk storage
	DiskStorage Type = 1 << iota
	// MemoryStorage memory storage
	MemoryStorage
)

// NewStorage create new storage by type
func NewStorage(t Type) Storage {
	switch t {
	case DiskStorage:
		return newDiskStorage()
	case MemoryStorage:
		return newMemoryStorage()
	default:
		return newMemoryStorage()
	}
}
