package file

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestSaveWithAbsolutePath(t *testing.T) {
	baseDir := t.TempDir()
	targetDir := filepath.Join(t.TempDir(), "content", "works")
	r := NewRepository(baseDir)

	target := filepath.Join(targetDir, "1.jpg")
	if err := r.Save(context.Background(), target, strings.NewReader("data")); err != nil {
		t.Fatalf("Save failed: %v", err)
	}

	got, err := os.ReadFile(target)
	if err != nil {
		t.Fatalf("file was not written to absolute path: %v", err)
	}
	if string(got) != "data" {
		t.Errorf("content = %q, want %q", string(got), "data")
	}

	if err := r.Delete(context.Background(), target); err != nil {
		t.Fatalf("Delete failed: %v", err)
	}
	if _, err := os.Stat(target); !os.IsNotExist(err) {
		t.Errorf("expected file to be deleted, got err %v", err)
	}
}

func TestSaveWithRelativePathJoinsBaseDir(t *testing.T) {
	baseDir := t.TempDir()
	r := NewRepository(baseDir)

	if err := r.Save(context.Background(), "content/works/1.jpg", strings.NewReader("data")); err != nil {
		t.Fatalf("Save failed: %v", err)
	}

	fullPath := filepath.Join(baseDir, "content", "works", "1.jpg")
	if _, err := os.Stat(fullPath); err != nil {
		t.Fatalf("file was not written to base dir: %v", err)
	}
}
