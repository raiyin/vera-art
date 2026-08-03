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

// TestSaleCRUD reproduces the full user flow for a store sale (art-store):
// creating it with all fields and an image, verifying the database state,
// editing it (adding a second image), verifying the changes, then deleting.
func TestSaleCRUD(t *testing.T) {
	r, db, baseDir := setupServer(t)
	seedAdminUser(t, db)
	token := login(t, r)

	// abs dir already points at the public root, so the rel dir prefix from the
	// service is repeated inside it (mirrors the production layout).
	salesRootDir := filepath.Join(baseDir, "content", "sales", "content", "sales")

	// ------------------------------------------------------------------
	// Create
	// ------------------------------------------------------------------
	createBody, createCT := buildMultipart(t,
		map[string][]string{
			"name_ru":      {"Тестовая картина"},
			"name_en":      {"Test painting"},
			"description":  {"Картина маслом на холсте"},
			"price":        {"12500.5"},
			"year":         {"2023"},
			"technique":    {"масло"},
			"width":        {"30"},
			"height":       {"45"},
			"status":       {"published"},
			"sort_order":   {"5"},
			"sold":         {"false"},
			"material_ids": {"1", "2"},
			"base_ids":     {"1"},
		},
		map[string][]formFile{
			"image": {
				{filename: "1.webp", content: testFile(t, "1.webp")},
			},
		},
	)

	w := doRequest(t, r, http.MethodPost, "/sales", token, createBody, createCT)
	mustStatus(t, w, http.StatusCreated)

	var created dto.SaleResponse
	if err := json.Unmarshal(w.Body.Bytes(), &created); err != nil {
		t.Fatalf("decode create sale response: %v", err)
	}
	mustEqual(t, "created id", created.ID > 0, true)
	mustEqual(t, "created name_ru", created.NameRu, "Тестовая картина")
	mustEqual(t, "created status", created.Status, "published")

	// The sale must exist in the database with the submitted fields.
	oldImages, oldSalePath, price := scanSaleRow(t, db, created.ID)
	mustEqual(t, "db name_ru", created.NameRu, "Тестовая картина")
	mustEqual(t, "db price", price, 12500.5)
	mustEqual(t, "db image_path set", oldImages != "", true)
	mustEqual(t, "db sale_path set", oldSalePath != "", true)
	mustEqual(t, "db material count", countRows(t, db, "SELECT COUNT(*) FROM sales_materials WHERE sale_id = ?", created.ID), 2)
	mustEqual(t, "db base count", countRows(t, db, "SELECT COUNT(*) FROM sales_bases WHERE sale_id = ?", created.ID), 1)

	// Uploaded image must be physically present on disk in the sale directory.
	oldImageDisk := filepath.Join(salesRootDir, oldSalePath, "1.webp")
	if _, err := os.Stat(oldImageDisk); err != nil {
		t.Fatalf("sale image missing on disk: %v", err)
	}

	// ------------------------------------------------------------------
	// Update: change every field and add a second image.
	// ------------------------------------------------------------------
	updateBody, updateCT := buildMultipart(t,
		map[string][]string{
			"name_ru":      {"Обновлённая картина"},
			"name_en":      {"Updated painting"},
			"description":  {"Акрил на холсте"},
			"price":        {"20000"},
			"year":         {"2024"},
			"technique":    {"акрил"},
			"width":        {"50"},
			"height":       {"80"},
			"status":       {"published"},
			"sort_order":   {"2"},
			"sold":         {"true"},
			"material_ids": {"3"},
			"base_ids":     {"2"},
			"images":       {"1.webp"},
		},
		map[string][]formFile{
			"image": {
				{filename: "2.webp", content: testFile(t, "2.webp")},
			},
		},
	)

	w = doRequest(t, r, http.MethodPut, "/sales/"+int64Str(created.ID), token, updateBody, updateCT)
	mustStatus(t, w, http.StatusOK)

	var updated dto.UpdateSaleResponse
	if err := json.Unmarshal(w.Body.Bytes(), &updated); err != nil {
		t.Fatalf("decode update sale response: %v", err)
	}
	mustEqual(t, "updated name_ru", updated.NameRu, "Обновлённая картина")
	mustEqual(t, "updated sold", updated.Sold, true)
	mustEqual(t, "updated material_ids", updated.MaterialIDs, []int64{3})

	newImages, newSalePath, newPrice := scanSaleRow(t, db, created.ID)
	mustEqual(t, "db updated price", newPrice, 20000.0)
	mustEqual(t, "db updated image changed", newImages != oldImages, true)
	mustEqual(t, "db updated sale_path unchanged", newSalePath, oldSalePath)
	mustEqual(t, "db updated material count", countRows(t, db, "SELECT COUNT(*) FROM sales_materials WHERE sale_id = ?", created.ID), 1)
	mustEqual(t, "db updated base count", countRows(t, db, "SELECT COUNT(*) FROM sales_bases WHERE sale_id = ?", created.ID), 1)

	// Both the kept and the newly added images must exist on disk.
	for _, name := range []string{"1.webp", "2.webp"} {
		if _, err := os.Stat(filepath.Join(salesRootDir, newSalePath, name)); err != nil {
			t.Fatalf("sale image %s missing on disk: %v", name, err)
		}
	}

	// The edited sale must be returned by the public read endpoint.
	w = doRequest(t, r, http.MethodGet, "/sales/"+int64Str(created.ID), "", nil, "")
	mustStatus(t, w, http.StatusOK)
	var byID dto.SaleResponse
	if err := json.Unmarshal(w.Body.Bytes(), &byID); err != nil {
		t.Fatalf("decode get sale response: %v", err)
	}
	mustEqual(t, "get sale name_ru", byID.NameRu, "Обновлённая картина")

	// ------------------------------------------------------------------
	// Delete
	// ------------------------------------------------------------------
	w = doRequest(t, r, http.MethodDelete, "/sales/"+int64Str(created.ID), token, nil, "")
	mustStatus(t, w, http.StatusOK)

	mustEqual(t, "sales after delete", countRows(t, db, "SELECT COUNT(*) FROM sales WHERE id = ?", created.ID), 0)
	mustEqual(t, "sales_materials after delete", countRows(t, db, "SELECT COUNT(*) FROM sales_materials WHERE sale_id = ?", created.ID), 0)
	mustEqual(t, "sales_bases after delete", countRows(t, db, "SELECT COUNT(*) FROM sales_bases WHERE sale_id = ?", created.ID), 0)
	for _, name := range []string{"1.webp", "2.webp"} {
		if _, err := os.Stat(filepath.Join(salesRootDir, newSalePath, name)); !os.IsNotExist(err) {
			t.Fatalf("sale image %s should be deleted after sale removal, stat err: %v", name, err)
		}
	}

	// The deleted sale must now be reported as not found.
	w = doRequest(t, r, http.MethodGet, "/sales/"+int64Str(created.ID), "", nil, "")
	mustStatus(t, w, http.StatusNotFound)
}

func scanSaleRow(t *testing.T, db *sql.DB, id int64) (imagePath, salePath string, price float64) {
	t.Helper()
	err := db.QueryRow(`
		SELECT image_path, sale_path, price, name_ru, name_en, descr_ru, descr_en, year, technique, width, height, status, sort_order, sold
		FROM sales WHERE id = ?`, id).
		Scan(&imagePath, &salePath, &price,
			new(string), new(string), new(sql.NullString), new(sql.NullString), new(sql.NullInt64),
			new(string), new(sql.NullInt64), new(sql.NullInt64),
			new(string), new(int), new(bool))
	if err != nil {
		t.Fatalf("scan sale row: %v", err)
	}
	return imagePath, salePath, price
}
