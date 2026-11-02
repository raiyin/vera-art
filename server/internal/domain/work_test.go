package domain

import (
	"testing"
)

func TestWorkDefaults(t *testing.T) {
	w := Work{}
	if w.ID != 0 {
		t.Errorf("expected ID 0, got %d", w.ID)
	}
	if w.StrID != "" {
		t.Errorf("expected empty StrID, got %q", w.StrID)
	}
	if w.BaseID != 0 {
		t.Errorf("expected BaseID 0, got %d", w.BaseID)
	}
	if w.WorkPath != "" {
		t.Errorf("expected empty WorkPath, got %q", w.WorkPath)
	}
	if w.Images != "" {
		t.Errorf("expected empty Images, got %q", w.Images)
	}
}

func TestWorkAssignment(t *testing.T) {
	w := Work{
		ID:       1,
		StrID:    "test-work",
		Width:    100,
		Height:   200,
		Year:     2024,
		NameRu:   "Тестовая работа",
		NameEn:   "Test Work",
		BaseID:   5,
		DescrRu:  "Описание",
		DescrEn:  "Description",
		WorkPath: "test-work/",
		Images:   "img1.jpg;img2.jpg",
	}

	if w.ID != 1 {
		t.Errorf("ID = %d, want 1", w.ID)
	}
	if w.StrID != "test-work" {
		t.Errorf("StrID = %q, want %q", w.StrID, "test-work")
	}
	if w.Width != 100 {
		t.Errorf("Width = %d, want 100", w.Width)
	}
	if w.Height != 200 {
		t.Errorf("Height = %d, want 200", w.Height)
	}
	if w.Year != 2024 {
		t.Errorf("Year = %d, want 2024", w.Year)
	}
	if w.NameRu != "Тестовая работа" {
		t.Errorf("NameRu = %q, want %q", w.NameRu, "Тестовая работа")
	}
	if w.NameEn != "Test Work" {
		t.Errorf("NameEn = %q, want %q", w.NameEn, "Test Work")
	}
	if w.BaseID != 5 {
		t.Errorf("BaseID = %d, want 5", w.BaseID)
	}
	if w.DescrRu != "Описание" {
		t.Errorf("DescrRu = %q, want %q", w.DescrRu, "Описание")
	}
	if w.DescrEn != "Description" {
		t.Errorf("DescrEn = %q, want %q", w.DescrEn, "Description")
	}
	if w.WorkPath != "test-work/" {
		t.Errorf("WorkPath = %q, want %q", w.WorkPath, "test-work/")
	}
	if w.Images != "img1.jpg;img2.jpg" {
		t.Errorf("Images = %q, want %q", w.Images, "img1.jpg;img2.jpg")
	}
}

func TestWorkFilterDefaults(t *testing.T) {
	f := WorkFilter{}
	if f.Query != "" {
		t.Errorf("expected empty Query, got %q", f.Query)
	}
	if f.Page != 0 {
		t.Errorf("expected Page 0, got %d", f.Page)
	}
	if f.Limit != 0 {
		t.Errorf("expected Limit 0, got %d", f.Limit)
	}
}

func TestWorkFilterAssignment(t *testing.T) {
	f := WorkFilter{
		Query: "test",
		Page:  2,
		Limit: 10,
	}

	if f.Query != "test" {
		t.Errorf("Query = %q, want %q", f.Query, "test")
	}
	if f.Page != 2 {
		t.Errorf("Page = %d, want 2", f.Page)
	}
	if f.Limit != 10 {
		t.Errorf("Limit = %d, want 10", f.Limit)
	}
}
