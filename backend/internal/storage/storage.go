package storage

import (
	"context"
	"io"
)

// FileStorage persists binary blobs addressed by a relative key (no ".." segments).
type FileStorage interface {
	Store(ctx context.Context, key string, r io.Reader, size int64) error
	Open(ctx context.Context, key string) (io.ReadCloser, error)
	Remove(ctx context.Context, key string) error
}
