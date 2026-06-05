package storage

import (
	"context"
	"errors"
	"io"
)

// ErrNotFound is returned when a storage key has no object on disk.
var ErrNotFound = errors.New("storage: object not found")

// FileStorage persists binary blobs addressed by a relative key (no ".." segments).
type FileStorage interface {
	Store(ctx context.Context, key string, r io.Reader, size int64) error
	Open(ctx context.Context, key string) (io.ReadCloser, error)
	Exists(ctx context.Context, key string) (bool, error)
	Remove(ctx context.Context, key string) error
}
