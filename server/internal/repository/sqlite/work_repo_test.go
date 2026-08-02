package sqlite

import (
	"context"
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"os"
	"path/filepath"
	"testing"

	"github.com/raiyin/artserver/internal/domain"
)

func setupWorkDB(t *testing.T) *sql.DB {
	t.Helper()
	dir, err := os.MkdirTemp("", "worktest-*")
	if err != nil {
		t.Fatalf("mkdir temp: %v", err)
	}
	t.Cleanup(func() { os.RemoveAll(dir) })
	dbPath := filepath.Join(dir, "test.db")

	db, err := sql.Open("sqlite3", dbPath+"?_journal_mode=WAL")
	if err != nil {
		t.Fatalf("open db: %v", err)
	}
	t.Cleanup(func() { db.Close() })

	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS works (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			str_id TEXT NOT NULL DEFAULT '',
			width INTEGER NOT NULL DEFAULT 0,
			height INTEGER NOT NULL DEFAULT 0,
			year INTEGER NOT NULL DEFAULT 0,
			name_ru TEXT NOT NULL DEFAULT '',
			name_en TEXT NOT NULL DEFAULT '',
			base_id INTEGER NOT NULL DEFAULT 0,
			descr_ru TEXT DEFAULT '',
			descr_en TEXT DEFAULT '',
			work_path TEXT NOT NULL DEFAULT '',
			images TEXT NOT NULL DEFAULT ''
		);
		CREATE TABLE IF NOT EXISTS works_materials (
			work_id INTEGER NOT NULL,
			material_id INTEGER NOT NULL,
			PRIMARY KEY (work_id, material_id)
		)
	`)
	if err != nil {
		t.Fatalf("create table: %v", err)
	}
	return db
}

func newWork(t *testing.T, repo *WorkRepository, work *domain.Work) *domain.Work {
	t.Helper()
	ctx := context.Background()
	if err := repo.Create(ctx, work); err != nil {
		t.Fatalf("Create: %v", err)
	}
	return work
}

func TestWorkRepoCreate(t *testing.T) {
	db := setupWorkDB(t)
	repo := NewWorkRepository(db)
	ctx := context.Background()

	w := &domain.Work{
		StrID:    "test",
		Width:    100,
		Height:   200,
		Year:     2024,
		NameRu:   "Test RU",
		NameEn:   "Test EN",
		BaseID:   1,
		DescrRu:  "Desc RU",
		DescrEn:  "Desc EN",
		WorkPath: "test/",
		Images:   "img1.jpg",
	}
	if err := repo.Create(ctx, w); err != nil {
		t.Fatalf("Create: %v", err)
	}
	if w.ID == 0 {
		t.Fatal("expected non-zero ID after create")
	}
}

func TestWorkRepoGetByID(t *testing.T) {
	db := setupWorkDB(t)
	repo := NewWorkRepository(db)
	ctx := context.Background()

	w := newWork(t, repo, &domain.Work{NameRu: "Test", NameEn: "Test"})

	got, err := repo.GetByID(ctx, w.ID)
	if err != nil {
		t.Fatalf("GetByID: %v", err)
	}
	if got.NameRu != "Test" {
		t.Errorf("NameRu = %q, want %q", got.NameRu, "Test")
	}
}

func TestWorkRepoGetByIDNotFound(t *testing.T) {
	db := setupWorkDB(t)
	repo := NewWorkRepository(db)
	_, err := repo.GetByID(context.Background(), 999)
	if err != domain.ErrNotFound {
		t.Errorf("expected ErrNotFound, got %v", err)
	}
}

func TestWorkRepoListEmpty(t *testing.T) {
	db := setupWorkDB(t)
	repo := NewWorkRepository(db)
	works, total, err := repo.List(context.Background(), domain.WorkFilter{})
	if err != nil {
		t.Fatalf("List: %v", err)
	}
	if len(works) != 0 {
		t.Errorf("expected 0 works, got %d", len(works))
	}
	if total != 0 {
		t.Errorf("expected total 0, got %d", total)
	}
}

func TestWorkRepoListMultiple(t *testing.T) {
	db := setupWorkDB(t)
	repo := NewWorkRepository(db)
	ctx := context.Background()

	for i := 0; i < 5; i++ {
		newWork(t, repo, &domain.Work{NameRu: "Test", NameEn: "Test"})
	}

	works, total, err := repo.List(ctx, domain.WorkFilter{})
	if err != nil {
		t.Fatalf("List: %v", err)
	}
	if len(works) != 5 {
		t.Errorf("expected 5 works, got %d", len(works))
	}
	if total != 5 {
		t.Errorf("expected total 5, got %d", total)
	}
}

func TestWorkRepoListPagination(t *testing.T) {
	db := setupWorkDB(t)
	repo := NewWorkRepository(db)
	ctx := context.Background()

	for i := 0; i < 10; i++ {
		newWork(t, repo, &domain.Work{NameRu: "Test", NameEn: "Test"})
	}

	t.Run("first page", func(t *testing.T) {
		works, total, err := repo.List(ctx, domain.WorkFilter{Page: 1, Limit: 3})
		if err != nil {
			t.Fatalf("List: %v", err)
		}
		if len(works) != 3 {
			t.Errorf("expected 3 works, got %d", len(works))
		}
		if total != 10 {
			t.Errorf("expected total 10, got %d", total)
		}
	})

	t.Run("second page", func(t *testing.T) {
		works, _, err := repo.List(ctx, domain.WorkFilter{Page: 2, Limit: 3})
		if err != nil {
			t.Fatalf("List: %v", err)
		}
		if len(works) != 3 {
			t.Errorf("expected 3 works, got %d", len(works))
		}
	})

	t.Run("last page partial", func(t *testing.T) {
		works, _, err := repo.List(ctx, domain.WorkFilter{Page: 4, Limit: 3})
		if err != nil {
			t.Fatalf("List: %v", err)
		}
		if len(works) != 1 {
			t.Errorf("expected 1 work, got %d", len(works))
		}
	})

	t.Run("page beyond end", func(t *testing.T) {
		works, _, err := repo.List(ctx, domain.WorkFilter{Page: 10, Limit: 3})
		if err != nil {
			t.Fatalf("List: %v", err)
		}
		if len(works) != 0 {
			t.Errorf("expected 0 works, got %d", len(works))
		}
	})
}

func TestWorkRepoListFilterQuery(t *testing.T) {
	db := setupWorkDB(t)
	repo := NewWorkRepository(db)
	ctx := context.Background()

	newWork(t, repo, &domain.Work{NameRu: "Абстракция", NameEn: "Abstract"})
	newWork(t, repo, &domain.Work{NameRu: "Пейзаж", NameEn: "Landscape"})
	newWork(t, repo, &domain.Work{NameRu: "Портрет", NameEn: "Portrait"})

	t.Run("filter by name_ru", func(t *testing.T) {
		works, total, err := repo.List(ctx, domain.WorkFilter{Query: "Абстракция"})
		if err != nil {
			t.Fatalf("List: %v", err)
		}
		if total != 1 {
			t.Errorf("expected total 1, got %d", total)
		}
		if len(works) != 1 || works[0].NameRu != "Абстракция" {
			t.Errorf("expected [Абстракция], got %v", works)
		}
	})

	t.Run("filter by name_en", func(t *testing.T) {
		works, total, err := repo.List(ctx, domain.WorkFilter{Query: "Landscape"})
		if err != nil {
			t.Fatalf("List: %v", err)
		}
		if total != 1 {
			t.Errorf("expected total 1, got %d", total)
		}
		if len(works) != 1 || works[0].NameEn != "Landscape" {
			t.Errorf("expected [Landscape], got %v", works)
		}
	})

	t.Run("filter no match", func(t *testing.T) {
		works, total, err := repo.List(ctx, domain.WorkFilter{Query: "Nonexistent"})
		if err != nil {
			t.Fatalf("List: %v", err)
		}
		if total != 0 {
			t.Errorf("expected total 0, got %d", total)
		}
		if len(works) != 0 {
			t.Errorf("expected 0 works, got %d", len(works))
		}
	})

	t.Run("filter empty query returns all", func(t *testing.T) {
		works, total, err := repo.List(ctx, domain.WorkFilter{Query: ""})
		if err != nil {
			t.Fatalf("List: %v", err)
		}
		if total != 3 {
			t.Errorf("expected total 3, got %d", total)
		}
		if len(works) != 3 {
			t.Errorf("expected 3 works, got %d", len(works))
		}
	})
}

func TestWorkRepoUpdate(t *testing.T) {
	db := setupWorkDB(t)
	repo := NewWorkRepository(db)
	ctx := context.Background()

	w := newWork(t, repo, &domain.Work{
		NameRu:   "Original",
		NameEn:   "Original",
		StrID:    "orig",
		Width:    100,
		Height:   200,
		Year:     2024,
		BaseID:   1,
		DescrRu:  "Old desc",
		DescrEn:  "Old desc",
		WorkPath: "orig/",
		Images:   "old.jpg",
	})

	w.NameRu = "Updated"
	w.Width = 300
	w.Height = 400
	w.Year = 2025
	w.DescrRu = "New desc"
	w.Images = "new.jpg"

	if err := repo.Update(ctx, w); err != nil {
		t.Fatalf("Update: %v", err)
	}

	got, err := repo.GetByID(ctx, w.ID)
	if err != nil {
		t.Fatalf("GetByID after update: %v", err)
	}
	if got.NameRu != "Updated" {
		t.Errorf("NameRu = %q, want %q", got.NameRu, "Updated")
	}
	if got.Width != 300 {
		t.Errorf("Width = %d, want 300", got.Width)
	}
	if got.Height != 400 {
		t.Errorf("Height = %d, want 400", got.Height)
	}
	if got.Year != 2025 {
		t.Errorf("Year = %d, want 2025", got.Year)
	}
	if got.DescrRu != "New desc" {
		t.Errorf("DescrRu = %q, want %q", got.DescrRu, "New desc")
	}
	if got.Images != "new.jpg" {
		t.Errorf("Images = %q, want %q", got.Images, "new.jpg")
	}
}

func TestWorkRepoUpdatePreservesUnchangedFields(t *testing.T) {
	db := setupWorkDB(t)
	repo := NewWorkRepository(db)
	ctx := context.Background()

	w := newWork(t, repo, &domain.Work{
		NameRu:   "Keep",
		NameEn:   "Keep",
		StrID:    "keep",
		Width:    100,
		Height:   200,
		Year:     2024,
		BaseID:   5,
		DescrRu:  "Keep desc",
		DescrEn:  "Keep desc",
		WorkPath: "keep/",
		Images:   "keep.jpg",
	})

	w.Width = 999
	if err := repo.Update(ctx, w); err != nil {
		t.Fatalf("Update: %v", err)
	}

	got, err := repo.GetByID(ctx, w.ID)
	if err != nil {
		t.Fatalf("GetByID: %v", err)
	}
	if got.NameRu != "Keep" {
		t.Errorf("NameRu changed to %q", got.NameRu)
	}
	if got.BaseID != 5 {
		t.Errorf("BaseID changed to %d", got.BaseID)
	}
	if got.WorkPath != "keep/" {
		t.Errorf("WorkPath changed to %q", got.WorkPath)
	}
	if got.Width != 999 {
		t.Errorf("Width = %d, want 999", got.Width)
	}
}

func TestWorkRepoDelete(t *testing.T) {
	db := setupWorkDB(t)
	repo := NewWorkRepository(db)
	ctx := context.Background()

	w := newWork(t, repo, &domain.Work{NameRu: "Delete Me", NameEn: "Delete Me"})

	if err := repo.Delete(ctx, w.ID); err != nil {
		t.Fatalf("Delete: %v", err)
	}

	_, err := repo.GetByID(ctx, w.ID)
	if err != domain.ErrNotFound {
		t.Errorf("expected ErrNotFound after delete, got %v", err)
	}
}

func TestWorkRepoDeleteNotFound(t *testing.T) {
	db := setupWorkDB(t)
	repo := NewWorkRepository(db)
	err := repo.Delete(context.Background(), 999)
	if err != domain.ErrNotFound {
		t.Errorf("expected ErrNotFound, got %v", err)
	}
}

func TestWorkRepoOrderByIDAscending(t *testing.T) {
	db := setupWorkDB(t)
	repo := NewWorkRepository(db)
	ctx := context.Background()

	var ids []int64
	for i := 0; i < 5; i++ {
		w := newWork(t, repo, &domain.Work{NameRu: "Test", NameEn: "Test"})
		ids = append(ids, w.ID)
	}

	works, _, err := repo.List(ctx, domain.WorkFilter{})
	if err != nil {
		t.Fatalf("List: %v", err)
	}
	for i, w := range works {
		if w.ID != ids[i] {
			t.Errorf("position %d: expected ID %d, got %d", i, ids[i], w.ID)
		}
	}
}

func TestWorkRepoCreateAllFields(t *testing.T) {
	db := setupWorkDB(t)
	repo := NewWorkRepository(db)
	ctx := context.Background()

	w := &domain.Work{
		StrID:    "all-fields",
		Width:    800,
		Height:   600,
		Year:     2023,
		NameRu:   "Все поля",
		NameEn:   "All Fields",
		BaseID:   3,
		DescrRu:  "Описание",
		DescrEn:  "Description",
		WorkPath: "all-fields/",
		Images:   "img1.jpg;img2.jpg",
	}
	if err := repo.Create(ctx, w); err != nil {
		t.Fatalf("Create: %v", err)
	}

	got, err := repo.GetByID(ctx, w.ID)
	if err != nil {
		t.Fatalf("GetByID: %v", err)
	}

	if got.StrID != "all-fields" {
		t.Errorf("StrID = %q, want %q", got.StrID, "all-fields")
	}
	if got.Width != 800 {
		t.Errorf("Width = %d, want 800", got.Width)
	}
	if got.Height != 600 {
		t.Errorf("Height = %d, want 600", got.Height)
	}
	if got.Year != 2023 {
		t.Errorf("Year = %d, want 2023", got.Year)
	}
	if got.NameRu != "Все поля" {
		t.Errorf("NameRu = %q, want %q", got.NameRu, "Все поля")
	}
	if got.NameEn != "All Fields" {
		t.Errorf("NameEn = %q, want %q", got.NameEn, "All Fields")
	}
	if got.BaseID != 3 {
		t.Errorf("BaseID = %d, want 3", got.BaseID)
	}
	if got.DescrRu != "Описание" {
		t.Errorf("DescrRu = %q, want %q", got.DescrRu, "Описание")
	}
	if got.DescrEn != "Description" {
		t.Errorf("DescrEn = %q, want %q", got.DescrEn, "Description")
	}
	if got.WorkPath != "all-fields/" {
		t.Errorf("WorkPath = %q, want %q", got.WorkPath, "all-fields/")
	}
	if got.Images != "img1.jpg;img2.jpg" {
		t.Errorf("Images = %q, want %q", got.Images, "img1.jpg;img2.jpg")
	}
}

func TestWorkRepoSetMaterials(t *testing.T) {
	db := setupWorkDB(t)
	repo := NewWorkRepository(db)
	ctx := context.Background()

	w := newWork(t, repo, &domain.Work{NameRu: "Materials", NameEn: "Materials"})

	if err := repo.SetMaterials(ctx, w.ID, []int64{1, 2, 3}); err != nil {
		t.Fatalf("SetMaterials: %v", err)
	}

	got, err := repo.GetByID(ctx, w.ID)
	if err != nil {
		t.Fatalf("GetByID: %v", err)
	}
	if len(got.MaterialIDs) != 3 {
		t.Fatalf("expected 3 material ids, got %v", got.MaterialIDs)
	}

	if err := repo.SetMaterials(ctx, w.ID, []int64{2}); err != nil {
		t.Fatalf("SetMaterials replace: %v", err)
	}

	got, err = repo.GetByID(ctx, w.ID)
	if err != nil {
		t.Fatalf("GetByID after replace: %v", err)
	}
	if len(got.MaterialIDs) != 1 || got.MaterialIDs[0] != 2 {
		t.Errorf("expected material ids [2], got %v", got.MaterialIDs)
	}
}

func TestWorkRepoDeleteRemovesFromList(t *testing.T) {	db := setupWorkDB(t)
	repo := NewWorkRepository(db)
	ctx := context.Background()

	w1 := newWork(t, repo, &domain.Work{NameRu: "A", NameEn: "A"})
	w2 := newWork(t, repo, &domain.Work{NameRu: "B", NameEn: "B"})
	w3 := newWork(t, repo, &domain.Work{NameRu: "C", NameEn: "C"})

	if err := repo.Delete(ctx, w2.ID); err != nil {
		t.Fatalf("Delete: %v", err)
	}

	works, total, err := repo.List(ctx, domain.WorkFilter{})
	if err != nil {
		t.Fatalf("List: %v", err)
	}
	if total != 2 {
		t.Errorf("expected total 2, got %d", total)
	}
	if len(works) != 2 {
		t.Errorf("expected 2 works, got %d", len(works))
	}
	if works[0].ID != w1.ID || works[1].ID != w3.ID {
		t.Error("expected remaining works to be A and C in order")
	}
}

func TestWorkRepoFilterWithPagination(t *testing.T) {
	db := setupWorkDB(t)
	repo := NewWorkRepository(db)
	ctx := context.Background()

	names := []string{"Apple", "Banana", "Cherry", "Date", "Elderberry"}
	for _, n := range names {
		newWork(t, repo, &domain.Work{NameRu: n, NameEn: n})
	}

	t.Run("filter + page 1", func(t *testing.T) {
		works, total, err := repo.List(ctx, domain.WorkFilter{Query: "Apple", Page: 1, Limit: 2})
		if err != nil {
			t.Fatalf("List: %v", err)
		}
		if total != 1 {
			t.Errorf("expected total 1, got %d", total)
		}
		if len(works) != 1 {
			t.Errorf("expected 1 work, got %d", len(works))
		}
	})

	t.Run("filter + pagination with no matches on page", func(t *testing.T) {
		works, total, err := repo.List(ctx, domain.WorkFilter{Query: "xyz", Page: 1, Limit: 2})
		if err != nil {
			t.Fatalf("List: %v", err)
		}
		if total != 0 {
			t.Errorf("expected total 0, got %d", total)
		}
		if len(works) != 0 {
			t.Errorf("expected 0 works, got %d", len(works))
		}
	})
}

func TestWorkRepoDefaultPaginationValues(t *testing.T) {
	db := setupWorkDB(t)
	repo := NewWorkRepository(db)
	ctx := context.Background()

	for i := 0; i < 60; i++ {
		newWork(t, repo, &domain.Work{NameRu: "Test", NameEn: "Test"})
	}

	t.Run("default page 1 with no filter", func(t *testing.T) {
		works, _, err := repo.List(ctx, domain.WorkFilter{Page: 0, Limit: 0})
		if err != nil {
			t.Fatalf("List: %v", err)
		}
		if len(works) != 50 {
			t.Errorf("expected 50 works (default limit), got %d", len(works))
		}
	})
}
