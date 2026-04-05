package storage

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

// Local stores files under a single root directory on disk.
type Local struct {
	root string
}

// NewLocal creates local storage, ensuring root exists (0700).
func NewLocal(root string) (*Local, error) {
	abs, err := filepath.Abs(root)
	if err != nil {
		return nil, fmt.Errorf("storage root: %w", err)
	}
	if err := os.MkdirAll(abs, 0700); err != nil {
		return nil, fmt.Errorf("create storage root: %w", err)
	}
	return &Local{root: abs}, nil
}

func (l *Local) resolve(key string) (string, error) {
	if key == "" || strings.Contains(key, "..") || strings.Contains(key, string(os.PathSeparator)) {
		return "", fmt.Errorf("invalid storage key")
	}
	// key is a single path segment (e.g. UUID string)
	candidate := filepath.Join(l.root, key)
	rootClean := filepath.Clean(l.root) + string(os.PathSeparator)
	candClean := filepath.Clean(candidate)
	if candClean != l.root && !strings.HasPrefix(candClean+string(os.PathSeparator), rootClean) {
		return "", fmt.Errorf("path escapes storage root")
	}
	return candidate, nil
}

func (l *Local) Store(ctx context.Context, key string, r io.Reader, size int64) error {
	path, err := l.resolve(key)
	if err != nil {
		return err
	}
	tmp := path + ".tmp"
	f, err := os.OpenFile(tmp, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		return fmt.Errorf("open temp: %w", err)
	}
	n, err := io.Copy(f, r)
	if err != nil {
		f.Close()
		os.Remove(tmp)
		return fmt.Errorf("write: %w", err)
	}
	if size > 0 && n != size {
		f.Close()
		os.Remove(tmp)
		return fmt.Errorf("size mismatch: wrote %d expected %d", n, size)
	}
	if err := f.Close(); err != nil {
		os.Remove(tmp)
		return err
	}
	if err := os.Rename(tmp, path); err != nil {
		os.Remove(tmp)
		return fmt.Errorf("rename: %w", err)
	}
	_ = ctx
	return nil
}

func (l *Local) Open(ctx context.Context, key string) (io.ReadCloser, error) {
	path, err := l.resolve(key)
	if err != nil {
		return nil, err
	}
	f, err := os.Open(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("object not found: %w", err)
		}
		return nil, err
	}
	_ = ctx
	return f, nil
}

func (l *Local) Remove(ctx context.Context, key string) error {
	path, err := l.resolve(key)
	if err != nil {
		return err
	}
	if err := os.Remove(path); err != nil && !os.IsNotExist(err) {
		return err
	}
	_ = ctx
	return nil
}
