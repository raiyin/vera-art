package integration

import (
	"encoding/json"
	"net/http"
	"os"
	"path/filepath"
	"testing"

	"github.com/raiyin/artserver/internal/dto"
)

// TestNewsCRUD reproduces the full user flow for a news item: creating it with
// a main image, additional images and a video, verifying the database and disk
// state, editing the texts, verifying the change, then deleting it.
func TestNewsCRUD(t *testing.T) {
	r, db, baseDir := setupServer(t)
	seedAdminUser(t, db)
	token := login(t, r)

	// abs dir already points at the public root, so the /content/news prefix in
	// news.Dir is repeated inside it (mirrors the production layout).
	newsDir := filepath.Join(baseDir, "content", "news", "content", "news", "2025", "01", "15")

	// ------------------------------------------------------------------
	// Create
	// ------------------------------------------------------------------
	dataJSON, err := json.Marshal(map[string]interface{}{
		"datetime":  "2025-01-15",
		"title_ru":  "Тестовая новость",
		"title_en":  "Test news",
		"dir":       "",
		"main_image": "",
		"text_ru":   "Текст новости",
		"text_en":   "News text",
		"images":    []string{},
		"videos":    []string{},
	})
	if err != nil {
		t.Fatalf("marshal news data: %v", err)
	}

	createBody, createCT := buildMultipart(t,
		map[string][]string{
			"data": {string(dataJSON)},
		},
		map[string][]formFile{
			"main_image": {
				{filename: "1.webp", content: testFile(t, "1.webp")},
			},
			"images": {
				{filename: "2.webp", content: testFile(t, "2.webp")},
			},
			"videos": {
				{filename: "sample-5s.mp4", content: testFile(t, "sample-5s.mp4")},
			},
		},
	)

	w := doRequest(t, r, http.MethodPost, "/news", token, createBody, createCT)
	mustStatus(t, w, http.StatusCreated)

	var created dto.NewsResponse
	if err := json.Unmarshal(w.Body.Bytes(), &created); err != nil {
		t.Fatalf("decode create news response: %v", err)
	}
	mustEqual(t, "created id", created.ID, "20250115")
	mustEqual(t, "created title_ru", created.TitleRu, "Тестовая новость")
	mustEqual(t, "created dir", created.Dir, "/content/news/2025/01/15/")
	mustEqual(t, "created main_image", created.MainImage, "1.webp")
	mustEqual(t, "created images", created.Images, []string{"2.webp"})
	mustEqual(t, "created videos", created.Videos, []string{"sample-5s"})

	// The news must exist in the database with the generated id and dir.
	row := db.QueryRow(`
		SELECT id, datetime, title_ru, title_en, dir, main_image, text_ru, text_en, images, videos
		FROM news WHERE id = ?`, created.ID)
	var (
		id, datetime, titleRu, titleEn, dir, mainImage, textRu, textEn, images, videos string
	)
	if err := row.Scan(&id, &datetime, &titleRu, &titleEn, &dir, &mainImage, &textRu, &textEn, &images, &videos); err != nil {
		t.Fatalf("scan created news row: %v", err)
	}
	mustEqual(t, "db id", id, "20250115")
	mustEqual(t, "db datetime", datetime, "2025-01-15")
	mustEqual(t, "db title_ru", titleRu, "Тестовая новость")
	mustEqual(t, "db title_en", titleEn, "Test news")
	mustEqual(t, "db dir", dir, "/content/news/2025/01/15/")
	mustEqual(t, "db main_image", mainImage, "1.webp")
	mustEqual(t, "db text_ru", textRu, "Текст новости")
	mustEqual(t, "db images", images, "2.webp")
	mustEqual(t, "db videos", videos, "sample-5s")

	// Uploaded files must be present on disk.
	for _, p := range []string{
		filepath.Join(newsDir, "1.webp"),
		filepath.Join(newsDir, "2.webp"),
		filepath.Join(newsDir, "videos", "sample-5s", "sample-5s.mp4"),
	} {
		if _, err := os.Stat(p); err != nil {
			t.Fatalf("uploaded news file missing on disk %s: %v", p, err)
		}
	}

	// The news must be returned by the public read endpoint.
	w = doRequest(t, r, http.MethodGet, "/news/"+created.ID, "", nil, "")
	mustStatus(t, w, http.StatusOK)
	var byID dto.NewsResponse
	if err := json.Unmarshal(w.Body.Bytes(), &byID); err != nil {
		t.Fatalf("decode get news response: %v", err)
	}
	mustEqual(t, "get news title_ru", byID.TitleRu, "Тестовая новость")

	// ------------------------------------------------------------------
	// Update: change texts, keep the existing images and video on the preview.
	// ------------------------------------------------------------------
	updateData, err := json.Marshal(map[string]interface{}{
		"datetime":   "2025-01-15",
		"title_ru":   "Обновлённая новость",
		"title_en":   "Updated news",
		"dir":        "/content/news/2025/01/15/",
		"main_image": "1.webp",
		"text_ru":    "Новый текст новости",
		"text_en":    "New news text",
		"images":     []string{"2.webp"},
		"videos":     []string{"sample-5s"},
	})
	if err != nil {
		t.Fatalf("marshal updated news data: %v", err)
	}

	updateBody, updateCT := buildMultipart(t,
		map[string][]string{
			"data": {string(updateData)},
		},
		nil,
	)

	w = doRequest(t, r, http.MethodPut, "/news/"+created.ID, token, updateBody, updateCT)
	mustStatus(t, w, http.StatusOK)

	var updated dto.NewsResponse
	if err := json.Unmarshal(w.Body.Bytes(), &updated); err != nil {
		t.Fatalf("decode update news response: %v", err)
	}
	mustEqual(t, "updated title_ru", updated.TitleRu, "Обновлённая новость")
	mustEqual(t, "updated images", updated.Images, []string{"2.webp"})
	mustEqual(t, "updated videos", updated.Videos, []string{"sample-5s"})

	// Changes must be visible in the database; media lists must be preserved.
	if err := db.QueryRow(`SELECT title_ru, text_ru, images, videos FROM news WHERE id = ?`, created.ID).
		Scan(&titleRu, &textRu, &images, &videos); err != nil {
		t.Fatalf("scan updated news row: %v", err)
	}
	mustEqual(t, "db updated title_ru", titleRu, "Обновлённая новость")
	mustEqual(t, "db updated text_ru", textRu, "Новый текст новости")
	mustEqual(t, "db updated images", images, "2.webp")
	mustEqual(t, "db updated videos", videos, "sample-5s")

	// ------------------------------------------------------------------
	// Delete
	// ------------------------------------------------------------------
	w = doRequest(t, r, http.MethodDelete, "/news/"+created.ID, token, nil, "")
	mustStatus(t, w, http.StatusOK)

	mustEqual(t, "news after delete", countRows(t, db, "SELECT COUNT(*) FROM news WHERE id = ?", created.ID), 0)

	// The news directory (with all its media) must be removed from disk.
	if _, err := os.Stat(filepath.Join(newsDir, "..")); err != nil {
		t.Fatalf("news parent dir missing: %v", err)
	}
	if _, err := os.Stat(newsDir); !os.IsNotExist(err) {
		t.Fatalf("news dir should be deleted after removal, stat err: %v", err)
	}

	// The deleted news must now be reported as not found.
	w = doRequest(t, r, http.MethodGet, "/news/"+created.ID, "", nil, "")
	mustStatus(t, w, http.StatusNotFound)
}
