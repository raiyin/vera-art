package file

import (
	"context"
	"io"
	"os"
	"path/filepath"
	"strings"
)

// Repository implements port.FileRepository for local filesystem storage.
type Repository struct {
	baseDir string
}

// NewRepository creates a new file repository.
func NewRepository(baseDir string) *Repository {
	return &Repository{baseDir: baseDir}
}

// Save saves a file to the given path.
func (r *Repository) Save(_ context.Context, path string, reader io.Reader) error {
	fullPath := filepath.Join(r.baseDir, path)

	if err := os.MkdirAll(filepath.Dir(fullPath), 0755); err != nil {
		return err
	}

	f, err := os.Create(fullPath)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = io.Copy(f, reader)
	return err
}

// Delete deletes a file at the given path.
func (r *Repository) Delete(_ context.Context, path string) error {
	fullPath := filepath.Join(r.baseDir, strings.ReplaceAll(path, "\\", "/"))
	return os.Remove(fullPath)
}

// GetPath returns the full filesystem path for a given relative path.
func (r *Repository) GetPath(dir, filename string) string {
	return filepath.Join(r.baseDir, dir, filename)
}

// Exists checks if a file exists at the given path.
func (r *Repository) Exists(_ context.Context, path string) (bool, error) {
	fullPath := filepath.Join(r.baseDir, path)
	_, err := os.Stat(fullPath)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

// Copy copies a file from src to dst.
func (r *Repository) Copy(_ context.Context, src, dst string) error {
	srcPath := filepath.Join(r.baseDir, src)
	dstPath := filepath.Join(r.baseDir, dst)

	if err := os.MkdirAll(filepath.Dir(dstPath), 0755); err != nil {
		return err
	}

	srcFile, err := os.Open(srcPath)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	dstFile, err := os.Create(dstPath)
	if err != nil {
		return err
	}
	defer dstFile.Close()

	_, err = io.Copy(dstFile, srcFile)
	return err
}

// MkdirAll creates a directory and all parent directories.
func (r *Repository) MkdirAll(_ context.Context, path string) error {
	fullPath := filepath.Join(r.baseDir, path)
	return os.MkdirAll(fullPath, 0755)
}

// RemoveDir removes a directory and all its contents.
func (r *Repository) RemoveDir(_ context.Context, path string) error {
	fullPath := filepath.Join(r.baseDir, path)
	return os.RemoveAll(fullPath)
}

// RenameDir renames a directory from oldPath to newPath.
func (r *Repository) RenameDir(_ context.Context, oldPath, newPath string) error {
	oldFullPath := filepath.Join(r.baseDir, oldPath)
	newFullPath := filepath.Join(r.baseDir, newPath)

	if err := os.MkdirAll(filepath.Dir(newFullPath), 0755); err != nil {
		return err
	}

	return os.Rename(oldFullPath, newFullPath)
}
