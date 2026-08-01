package service

import (
	"context"
	"io"
	"reflect"
	"strings"
	"testing"

	"github.com/raiyin/artserver/internal/domain"
)

type mockNewsRepo struct {
	news   *domain.News
	err    error
	update *domain.News
}

func (m *mockNewsRepo) Create(_ context.Context, _ *domain.News) error { return m.err }

func (m *mockNewsRepo) GetByID(_ context.Context, _ string) (*domain.News, error) {
	if m.err != nil {
		return nil, m.err
	}
	return m.news, nil
}

func (m *mockNewsRepo) List(_ context.Context, _ domain.NewsFilter) ([]domain.News, int, error) {
	return nil, 0, m.err
}

func (m *mockNewsRepo) Update(_ context.Context, news *domain.News) error {
	m.update = news
	return m.err
}

func (m *mockNewsRepo) Delete(_ context.Context, _ string) error { return m.err }

type trackingFileRepo struct {
	deleted   []string
	removed   []string
	saved     []string
}

func (m *trackingFileRepo) Save(_ context.Context, path string, _ io.Reader) error {
	m.saved = append(m.saved, path)
	return nil
}
func (m *trackingFileRepo) Delete(_ context.Context, path string) error {
	m.deleted = append(m.deleted, path)
	return nil
}
func (m *trackingFileRepo) GetPath(_, _ string) string { return "" }
func (m *trackingFileRepo) Exists(_ context.Context, _ string) (bool, error) {
	return false, nil
}
func (m *trackingFileRepo) Copy(_ context.Context, _, _ string) error { return nil }
func (m *trackingFileRepo) MkdirAll(_ context.Context, _ string) error {
	return nil
}
func (m *trackingFileRepo) RemoveDir(_ context.Context, path string) error {
	m.removed = append(m.removed, path)
	return nil
}
func (m *trackingFileRepo) RenameDir(_ context.Context, _, _ string) error {
	return nil
}

func newTestNewsService(newsRepo *mockNewsRepo, fileRepo *trackingFileRepo) *NewsService {
	return NewNewsService(newsRepo, fileRepo, "/data/news")
}

func existingNews() *domain.News {
	return &domain.News{
		ID:          "20240101",
		Dir:         "/content/news/2024/01/01/",
		ImgBack:     "back.jpg",
		ImgBackfull: "back_full.jpg",
		Images:      []string{"1.jpg", "2.jpg", "3.jpg"},
		Videos:      []string{"clip1", "clip2"},
	}
}

func TestUpdateNews_RemovesFilesNotLeftOnPreview(t *testing.T) {
	newsRepo := &mockNewsRepo{news: existingNews()}
	fileRepo := &trackingFileRepo{}
	svc := newTestNewsService(newsRepo, fileRepo)

	// User removed "2.jpg" and "clip1" from the preview.
	news := &domain.News{
		ID:          "20240101",
		ImgBack:     "back.jpg",
		ImgBackfull: "back_full.jpg",
		Images:      []string{"1.jpg", "3.jpg"},
		Videos:      []string{"clip2"},
	}

	if err := svc.UpdateNews(context.Background(), news, nil, nil, nil, nil); err != nil {
		t.Fatalf("UpdateNews failed: %v", err)
	}

	if newsRepo.update == nil {
		t.Fatal("UpdateNews did not persist the news")
	}

	if want := []string{"1.jpg", "3.jpg"}; !reflect.DeepEqual(newsRepo.update.Images, want) {
		t.Errorf("Images = %v, want %v", newsRepo.update.Images, want)
	}
	if want := []string{"clip2"}; !reflect.DeepEqual(newsRepo.update.Videos, want) {
		t.Errorf("Videos = %v, want %v", newsRepo.update.Videos, want)
	}

	// Only the removed file should have been deleted from disk.
	if !reflect.DeepEqual(fileRepo.deleted, []string{"/data/news/content/news/2024/01/01/2.jpg"}) {
		t.Errorf("deleted files = %v, want only 2.jpg", fileRepo.deleted)
	}
	if !reflect.DeepEqual(fileRepo.removed, []string{"/data/news/content/news/2024/01/01/videos/clip1"}) {
		t.Errorf("removed dirs = %v, want only clip1", fileRepo.removed)
	}
}

func TestUpdateNews_KeepsExistingAndAddsNew(t *testing.T) {
	newsRepo := &mockNewsRepo{news: existingNews()}
	fileRepo := &trackingFileRepo{}
	svc := newTestNewsService(newsRepo, fileRepo)

	// User kept "1.jpg", removed "2.jpg"/"3.jpg" and added "4.jpg".
	newImage := domain.UploadedFile{Filename: "4.jpg", Reader: strings.NewReader("img")}
	news := &domain.News{
		ID:          "20240101",
		ImgBack:     "back.jpg",
		ImgBackfull: "back_full.jpg",
		Images:      []string{"1.jpg", "4.jpg"},
		Videos:      []string{"clip1", "clip2"},
	}

	if err := svc.UpdateNews(context.Background(), news, nil, nil, []domain.UploadedFile{newImage}, nil); err != nil {
		t.Fatalf("UpdateNews failed: %v", err)
	}

	if want := []string{"1.jpg", "4.jpg"}; !reflect.DeepEqual(newsRepo.update.Images, want) {
		t.Errorf("Images = %v, want %v", newsRepo.update.Images, want)
	}

	if !reflect.DeepEqual(fileRepo.deleted, []string{"/data/news/content/news/2024/01/01/2.jpg", "/data/news/content/news/2024/01/01/3.jpg"}) {
		t.Errorf("deleted files = %v, want 2.jpg and 3.jpg", fileRepo.deleted)
	}

	found := false
	for _, p := range fileRepo.saved {
		if p == "/data/news/content/news/2024/01/01/4.jpg" {
			found = true
		}
	}
	if !found {
		t.Errorf("saved files = %v, want new image 4.jpg to be saved", fileRepo.saved)
	}
}

func TestUpdateNews_RemovesMainImage(t *testing.T) {
	newsRepo := &mockNewsRepo{news: existingNews()}
	fileRepo := &trackingFileRepo{}
	svc := newTestNewsService(newsRepo, fileRepo)

	// User removed img_back and img_backfull.
	news := &domain.News{
		ID:          "20240101",
		ImgBack:     "",
		ImgBackfull: "",
		Images:      []string{"1.jpg", "2.jpg", "3.jpg"},
		Videos:      []string{"clip1", "clip2"},
	}

	if err := svc.UpdateNews(context.Background(), news, nil, nil, nil, nil); err != nil {
		t.Fatalf("UpdateNews failed: %v", err)
	}

	if newsRepo.update.ImgBack != "" || newsRepo.update.ImgBackfull != "" {
		t.Errorf("ImgBack = %q, ImgBackfull = %q, want both empty", newsRepo.update.ImgBack, newsRepo.update.ImgBackfull)
	}

	for _, f := range []string{"/data/news/content/news/2024/01/01/back.jpg", "/data/news/content/news/2024/01/01/back_full.jpg"} {
		if !contains(fileRepo.deleted, f) {
			t.Errorf("expected %s to be deleted, deleted = %v", f, fileRepo.deleted)
		}
	}
}

func TestUpdateNews_KeepsEverythingWhenPayloadNil(t *testing.T) {
	newsRepo := &mockNewsRepo{news: existingNews()}
	fileRepo := &trackingFileRepo{}
	svc := newTestNewsService(newsRepo, fileRepo)

	// Payload does not include images/videos fields (nil) -> keep existing.
	news := &domain.News{
		ID:          "20240101",
		ImgBack:     "back.jpg",
		ImgBackfull: "back_full.jpg",
	}

	if err := svc.UpdateNews(context.Background(), news, nil, nil, nil, nil); err != nil {
		t.Fatalf("UpdateNews failed: %v", err)
	}

	if want := []string{"1.jpg", "2.jpg", "3.jpg"}; !reflect.DeepEqual(newsRepo.update.Images, want) {
		t.Errorf("Images = %v, want %v", newsRepo.update.Images, want)
	}
	if want := []string{"clip1", "clip2"}; !reflect.DeepEqual(newsRepo.update.Videos, want) {
		t.Errorf("Videos = %v, want %v", newsRepo.update.Videos, want)
	}
	if len(fileRepo.deleted) > 0 || len(fileRepo.removed) > 0 {
		t.Errorf("no files should be deleted, deleted=%v removed=%v", fileRepo.deleted, fileRepo.removed)
	}
}

func contains(list []string, item string) bool {
	for _, v := range list {
		if v == item {
			return true
		}
	}
	return false
}
