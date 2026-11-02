package integration

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"os"
	"path/filepath"
	"testing"

	"github.com/raiyin/artserver/internal/dto"
)

// TestWorkCRUD reproduces the full user flow for a gallery work (works):
// creating it with all fields and two images, verifying the data landed in the
// database, editing it (including a new image and replaced materials), verifying
// the changes again, and finally deleting it.
func TestWorkCRUD(t *testing.T) {
	r, db, baseDir := setupServer(t)
	seedAdminUser(t, db)
	token := login(t, r)

	// ------------------------------------------------------------------
	// Create
	// ------------------------------------------------------------------
	createBody, createCT := buildMultipart(t,
		map[string][]string{
			"str_id":       {"test-work-01"},
			"width":        {"40"},
			"height":       {"60"},
			"year":         {"2024"},
			"name_ru":      {"Тестовая работа"},
			"name_en":      {"Test work"},
			"base_id":      {"1"},
			"descr_ru":     {"Описание работы"},
			"descr_en":     {"Work description"},
			"material_ids": {"1", "2"},
		},
		map[string][]formFile{
			"image": {
				{filename: "1.webp", content: testFile(t, "1.webp")},
				{filename: "2.webp", content: testFile(t, "2.webp")},
			},
		},
	)

	w := doRequest(t, r, http.MethodPost, "/works", token, createBody, createCT)
	mustStatus(t, w, http.StatusCreated)

	var created dto.WorkResponse
	if err := json.Unmarshal(w.Body.Bytes(), &created); err != nil {
		t.Fatalf("decode create work response: %v", err)
	}
	mustEqual(t, "created id", created.ID > 0, true)
	mustEqual(t, "created name_ru", created.NameRu, "Тестовая работа")
	mustEqual(t, "created name_en", created.NameEn, "Test work")
	mustEqual(t, "created width", created.Width, 40)
	mustEqual(t, "created height", created.Height, 60)
	mustEqual(t, "created year", created.Year, 2024)
	mustEqual(t, "created base_id", created.BaseID, int64(1))
	mustEqual(t, "created work_path", created.WorkPath, "test-work-01/")
	mustEqual(t, "created images", len(created.Images), 2)

	// The work must exist in the database with the submitted fields.
	row := db.QueryRow(`
		SELECT str_id, width, height, year, name_ru, name_en, base_id, descr_ru, descr_en, work_path, images
		FROM works WHERE id = ?`, created.ID)
	var (
		strID string
		width, height, year int
		nameRu, nameEn      string
		baseID              int64
		descrRu, descrEn    string
		workPath, images    string
	)
	if err := row.Scan(&strID, &width, &height, &year, &nameRu, &nameEn, &baseID, &descrRu, &descrEn, &workPath, &images); err != nil {
		t.Fatalf("scan created work row: %v", err)
	}
	mustEqual(t, "db str_id", strID, "test-work-01")
	mustEqual(t, "db name_ru", nameRu, "Тестовая работа")
	mustEqual(t, "db name_en", nameEn, "Test work")
	mustEqual(t, "db width", width, 40)
	mustEqual(t, "db height", height, 60)
	mustEqual(t, "db year", year, 2024)
	mustEqual(t, "db base_id", baseID, int64(1))
	mustEqual(t, "db descr_ru", descrRu, "Описание работы")
	mustEqual(t, "db descr_en", descrEn, "Work description")
	mustEqual(t, "db work_path", workPath, "test-work-01/")
	mustEqual(t, "db images", images, "1.webp;2.webp")

	// Both materials must be linked.
	mustEqual(t, "works_materials count", countRows(t, db, "SELECT COUNT(*) FROM works_materials WHERE work_id = ?", created.ID), 2)
	mustEqual(t, "works_materials ids", workMaterialIDs(t, db, created.ID), []int64{1, 2})

	// Uploaded images must be physically present on disk.
	// abs dir already points at the public root, so the rel dir prefix from the
	// service is repeated inside it (mirrors the production layout).
	workDir := filepath.Join(baseDir, "content", "works", "content", "works", "test-work-01")
	for _, img := range []string{"1.webp", "2.webp"} {
		if _, err := os.Stat(filepath.Join(workDir, img)); err != nil {
			t.Fatalf("uploaded image %s missing on disk: %v", img, err)
		}
	}

	// ------------------------------------------------------------------
	// Update: change every editable field, keep the old preview images,
	// upload one new image and replace the materials.
	// ------------------------------------------------------------------
	updateBody, updateCT := buildMultipart(t,
		map[string][]string{
			"str_id":       {"test-work-01"},
			"width":        {"50"},
			"height":       {"70"},
			"year":         {"2025"},
			"name_ru":      {"Обновлённая работа"},
			"name_en":      {"Updated work"},
			"base_id":      {"2"},
			"descr_ru":     {"Новое описание"},
			"descr_en":     {"New description"},
			"material_ids": {"3"},
		},
		map[string][]formFile{
			"image": {
				{filename: "3.webp", content: testFile(t, "3.webp")},
			},
		},
	)

	w = doRequest(t, r, http.MethodPut, "/works/"+int64Str(created.ID), token, updateBody, updateCT)
	mustStatus(t, w, http.StatusOK)

	var updated dto.UpdateWorkResponse
	if err := json.Unmarshal(w.Body.Bytes(), &updated); err != nil {
		t.Fatalf("decode update work response: %v", err)
	}
	mustEqual(t, "updated name_ru", updated.NameRu, "Обновлённая работа")
	mustEqual(t, "updated images", updated.Images, []string{"1.webp", "2.webp", "3.webp"})
	mustEqual(t, "updated material_ids", updated.MaterialIDs, []int64{3})

	// Changes must be visible in the database.
	if err := db.QueryRow(`
		SELECT name_ru, name_en, width, height, year, base_id, descr_ru, descr_en, work_path, images
		FROM works WHERE id = ?`, created.ID).
		Scan(&nameRu, &nameEn, &width, &height, &year, &baseID, &descrRu, &descrEn, &workPath, &images); err != nil {
		t.Fatalf("scan updated work row: %v", err)
	}
	mustEqual(t, "db updated name_ru", nameRu, "Обновлённая работа")
	mustEqual(t, "db updated name_en", nameEn, "Updated work")
	mustEqual(t, "db updated width", width, 50)
	mustEqual(t, "db updated height", height, 70)
	mustEqual(t, "db updated year", year, 2025)
	mustEqual(t, "db updated base_id", baseID, int64(2))
	mustEqual(t, "db updated descr_ru", descrRu, "Новое описание")
	mustEqual(t, "db updated work_path", workPath, "test-work-01/")
	mustEqual(t, "db updated images", images, "1.webp;2.webp;3.webp")
	mustEqual(t, "db updated materials", workMaterialIDs(t, db, created.ID), []int64{3})

	// The new image must exist on disk next to the originals.
	if _, err := os.Stat(filepath.Join(workDir, "3.webp")); err != nil {
		t.Fatalf("updated image 3.webp missing on disk: %v", err)
	}

	// The edited work must be returned by the public read endpoints.
	w = doRequest(t, r, http.MethodGet, "/works/"+int64Str(created.ID), "", nil, "")
	mustStatus(t, w, http.StatusOK)
	var byID dto.WorkResponse
	if err := json.Unmarshal(w.Body.Bytes(), &byID); err != nil {
		t.Fatalf("decode get work response: %v", err)
	}
	mustEqual(t, "get work name_ru", byID.NameRu, "Обновлённая работа")

	// ------------------------------------------------------------------
	// Delete
	// ------------------------------------------------------------------
	w = doRequest(t, r, http.MethodDelete, "/works/"+int64Str(created.ID), token, nil, "")
	mustStatus(t, w, http.StatusOK)

	mustEqual(t, "works after delete", countRows(t, db, "SELECT COUNT(*) FROM works WHERE id = ?", created.ID), 0)
	mustEqual(t, "works_materials after delete", countRows(t, db, "SELECT COUNT(*) FROM works_materials WHERE work_id = ?", created.ID), 0)

	// The deleted work must now be reported as not found.
	w = doRequest(t, r, http.MethodGet, "/works/"+int64Str(created.ID), "", nil, "")
	mustStatus(t, w, http.StatusNotFound)
}

func workMaterialIDs(t *testing.T, db *sql.DB, workID int64) []int64 {
	t.Helper()
	rows, err := db.Query("SELECT material_id FROM works_materials WHERE work_id = ? ORDER BY material_id", workID)
	if err != nil {
		t.Fatalf("query work materials: %v", err)
	}
	defer rows.Close()

	var ids []int64
	for rows.Next() {
		var id int64
		if err := rows.Scan(&id); err != nil {
			t.Fatalf("scan work material: %v", err)
		}
		ids = append(ids, id)
	}
	return ids
}
