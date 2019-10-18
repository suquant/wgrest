package storage

import (
	"errors"
)

var (
	// ErrClosed error for read/write in cloased storage
	ErrClosed = errors.New("storage already closed")
)
