package service

import (
	"bytes"
	"context"
	"errors"
	"io"
	"strings"
	"testing"

	"github.com/raiyin/artserver/internal/domain"
)

// trackedMockFileRepo extends mockFileRepo to track Save/Delete calls.
type trackedMockFileRepo struct {
	mockFileRepo
	savedPaths   []string
	deletedPaths []string
}

func (m *trackedMockFileRepo) Save(_ context.Context, path string, _ io.Reader) error {
	if m.err != nil {
		return m.err
	}
	m.savedPaths = append(m.savedPaths, path)
	return nil
}

func (m *trackedMockFileRepo) Delete(_ context.Context, path string) error {
	if m.err != nil {
		return m.err
	}
	m.deletedPaths = append(m.deletedPaths, path)
	return nil
}

func newUserService(userRepo *mockUserRepo, fileRepo *trackedMockFileRepo) *UserService {
	return NewUserService(userRepo, fileRepo, "/tmp/avatars")
}

func TestUploadAvatar(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		userRepo := &mockUserRepo{
			users: []domain.User{
				{ID: 1, Username: "testuser", Email: "test@example.com", AvatarPath: ""},
			},
		}
		fileRepo := &trackedMockFileRepo{}
		svc := newUserService(userRepo, fileRepo)

		reader := strings.NewReader("fake-image-data")
		path, err := svc.UploadAvatar(context.Background(), 1, "photo.jpg", reader)
		if err != nil {
			t.Fatalf("UploadAvatar failed: %v", err)
		}

		if !strings.HasSuffix(path, "/tmp/avatars/avatar_1.jpg") {
			t.Errorf("expected path ending with /tmp/avatars/avatar_1.jpg, got %s", path)
		}

		if len(fileRepo.savedPaths) != 1 {
			t.Fatalf("expected 1 Save call, got %d", len(fileRepo.savedPaths))
		}
		if fileRepo.savedPaths[0] != path {
			t.Errorf("saved path mismatch: %s vs %s", fileRepo.savedPaths[0], path)
		}

		// Verify user's avatar_path was updated
		if len(userRepo.users) != 1 {
			t.Fatal("expected 1 user")
		}
		if userRepo.users[0].AvatarPath != path {
			t.Errorf("expected user.AvatarPath = %s, got %s", path, userRepo.users[0].AvatarPath)
		}
	})

	t.Run("replaces old avatar", func(t *testing.T) {
		userRepo := &mockUserRepo{
			users: []domain.User{
				{ID: 1, Username: "testuser", Email: "test@example.com", AvatarPath: "/tmp/avatars/avatar_1_old.jpg"},
			},
		}
		fileRepo := &trackedMockFileRepo{}
		svc := newUserService(userRepo, fileRepo)

		reader := strings.NewReader("new-image-data")
		path, err := svc.UploadAvatar(context.Background(), 1, "newphoto.png", reader)
		if err != nil {
			t.Fatalf("UploadAvatar failed: %v", err)
		}

		// Old file was deleted
		if len(fileRepo.deletedPaths) != 1 {
			t.Fatalf("expected 1 Delete call, got %d", len(fileRepo.deletedPaths))
		}
		if fileRepo.deletedPaths[0] != "/tmp/avatars/avatar_1_old.jpg" {
			t.Errorf("expected delete of old path, got %s", fileRepo.deletedPaths[0])
		}

		// New file was saved
		if len(fileRepo.savedPaths) != 1 {
			t.Fatalf("expected 1 Save call, got %d", len(fileRepo.savedPaths))
		}
		if !strings.HasSuffix(path, ".png") {
			t.Errorf("expected .png extension, got %s", path)
		}

		// User path was updated
		if userRepo.users[0].AvatarPath != path {
			t.Errorf("expected user.AvatarPath = %s, got %s", path, userRepo.users[0].AvatarPath)
		}
	})

	t.Run("invalid file extension", func(t *testing.T) {
		userRepo := &mockUserRepo{
			users: []domain.User{
				{ID: 1, Username: "testuser", Email: "test@example.com"},
			},
		}
		fileRepo := &trackedMockFileRepo{}
		svc := newUserService(userRepo, fileRepo)

		_, err := svc.UploadAvatar(context.Background(), 1, "photo.exe", bytes.NewReader([]byte("data")))
		if !errors.Is(err, domain.ErrInvalidFileType) {
			t.Errorf("expected ErrInvalidFileType, got %v", err)
		}
		if len(fileRepo.savedPaths) != 0 {
			t.Errorf("expected 0 Save calls for invalid extension, got %d", len(fileRepo.savedPaths))
		}
	})

	t.Run("user not found", func(t *testing.T) {
		userRepo := &mockUserRepo{}
		fileRepo := &trackedMockFileRepo{}
		svc := newUserService(userRepo, fileRepo)

		_, err := svc.UploadAvatar(context.Background(), 1, "photo.jpg", bytes.NewReader([]byte("data")))
		if !errors.Is(err, domain.ErrNotFound) {
			t.Errorf("expected ErrNotFound, got %v", err)
		}
		if len(fileRepo.savedPaths) != 0 {
			t.Errorf("expected 0 Save calls when user not found, got %d", len(fileRepo.savedPaths))
		}
	})

	t.Run("file save fails", func(t *testing.T) {
		userRepo := &mockUserRepo{
			users: []domain.User{
				{ID: 1, Username: "testuser", Email: "test@example.com"},
			},
		}
		fileRepo := &trackedMockFileRepo{}
		fileRepo.err = errors.New("disk full")
		svc := newUserService(userRepo, fileRepo)

		_, err := svc.UploadAvatar(context.Background(), 1, "photo.jpg", bytes.NewReader([]byte("data")))
		if err == nil {
			t.Fatal("expected error, got nil")
		}

		// user avatar_path should NOT be updated on save failure
		if userRepo.users[0].AvatarPath != "" {
			t.Errorf("expected empty AvatarPath on save failure, got %s", userRepo.users[0].AvatarPath)
		}
	})
}

func TestDeleteAvatar(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		userRepo := &mockUserRepo{
			users: []domain.User{
				{ID: 1, Username: "testuser", AvatarPath: "/tmp/avatars/avatar_1.jpg"},
			},
		}
		fileRepo := &trackedMockFileRepo{}
		svc := newUserService(userRepo, fileRepo)

		err := svc.DeleteAvatar(context.Background(), 1)
		if err != nil {
			t.Fatalf("DeleteAvatar failed: %v", err)
		}

		if len(fileRepo.deletedPaths) != 1 {
			t.Fatalf("expected 1 Delete call, got %d", len(fileRepo.deletedPaths))
		}
		if fileRepo.deletedPaths[0] != "/tmp/avatars/avatar_1.jpg" {
			t.Errorf("expected delete of avatar file, got %s", fileRepo.deletedPaths[0])
		}

		if userRepo.users[0].AvatarPath != "" {
			t.Errorf("expected empty AvatarPath after delete, got %s", userRepo.users[0].AvatarPath)
		}
	})

	t.Run("no avatar to delete", func(t *testing.T) {
		userRepo := &mockUserRepo{
			users: []domain.User{
				{ID: 1, Username: "testuser", AvatarPath: ""},
			},
		}
		fileRepo := &trackedMockFileRepo{}
		svc := newUserService(userRepo, fileRepo)

		err := svc.DeleteAvatar(context.Background(), 1)
		if err != nil {
			t.Fatalf("DeleteAvatar failed: %v", err)
		}

		if len(fileRepo.deletedPaths) != 0 {
			t.Errorf("expected 0 Delete calls when no avatar, got %d", len(fileRepo.deletedPaths))
		}
	})

	t.Run("user not found", func(t *testing.T) {
		userRepo := &mockUserRepo{}
		fileRepo := &trackedMockFileRepo{}
		svc := newUserService(userRepo, fileRepo)

		err := svc.DeleteAvatar(context.Background(), 999)
		if !errors.Is(err, domain.ErrNotFound) {
			t.Errorf("expected ErrNotFound, got %v", err)
		}
	})
}

func TestServeAvatar(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		userRepo := &mockUserRepo{
			users: []domain.User{
				{ID: 1, Username: "testuser", AvatarPath: "/tmp/avatars/avatar_1.jpg"},
			},
		}
		svc := newUserService(userRepo, &trackedMockFileRepo{})

		path, err := svc.ServeAvatar(context.Background(), 1)
		if err != nil {
			t.Fatalf("ServeAvatar failed: %v", err)
		}
		if path != "/tmp/avatars/avatar_1.jpg" {
			t.Errorf("expected /tmp/avatars/avatar_1.jpg, got %s", path)
		}
	})

	t.Run("no avatar set", func(t *testing.T) {
		userRepo := &mockUserRepo{
			users: []domain.User{
				{ID: 1, Username: "testuser", AvatarPath: ""},
			},
		}
		svc := newUserService(userRepo, &trackedMockFileRepo{})

		_, err := svc.ServeAvatar(context.Background(), 1)
		if !errors.Is(err, domain.ErrNotFound) {
			t.Errorf("expected ErrNotFound, got %v", err)
		}
	})

	t.Run("user not found", func(t *testing.T) {
		userRepo := &mockUserRepo{}
		svc := newUserService(userRepo, &trackedMockFileRepo{})

		_, err := svc.ServeAvatar(context.Background(), 999)
		if !errors.Is(err, domain.ErrNotFound) {
			t.Errorf("expected ErrNotFound, got %v", err)
		}
	})
}
