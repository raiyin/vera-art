package service

import (
	"context"
	"errors"
	"io"
	"reflect"
	"testing"

	"github.com/raiyin/artserver/internal/domain"
	"github.com/raiyin/artserver/internal/port"
)

// mockWorkRepo implements port.WorkRepository.
type mockWorkRepo struct {
	works []domain.Work
	err   error
}

func (m *mockWorkRepo) Create(_ context.Context, work *domain.Work) error {
	if m.err != nil {
		return m.err
	}
	work.ID = int64(len(m.works) + 1)
	m.works = append(m.works, *work)
	return nil
}

func (m *mockWorkRepo) GetByID(_ context.Context, id int64) (*domain.Work, error) {
	if m.err != nil {
		return nil, m.err
	}
	for _, w := range m.works {
		if w.ID == id {
			return &w, nil
		}
	}
	return nil, domain.ErrNotFound
}

func (m *mockWorkRepo) List(_ context.Context, filter domain.WorkFilter) ([]domain.Work, int, error) {
	if m.err != nil {
		return nil, 0, m.err
	}
	result := m.works
	if filter.Query != "" {
		var filtered []domain.Work
		for _, w := range result {
			if w.NameRu == filter.Query || w.NameEn == filter.Query {
				filtered = append(filtered, w)
			}
		}
		result = filtered
	}
	return result, len(result), nil
}

func (m *mockWorkRepo) Update(_ context.Context, work *domain.Work) error {
	if m.err != nil {
		return m.err
	}
	for i, w := range m.works {
		if w.ID == work.ID {
			m.works[i] = *work
			return nil
		}
	}
	return domain.ErrNotFound
}

func (m *mockWorkRepo) Delete(_ context.Context, id int64) error {
	if m.err != nil {
		return m.err
	}
	for i, w := range m.works {
		if w.ID == id {
			m.works = append(m.works[:i], m.works[i+1:]...)
			return nil
		}
	}
	return domain.ErrNotFound
}

// mockSaleRepo implements port.SaleRepository.
type mockSaleRepo struct {
	sales []domain.Sale
	err   error
}

func (m *mockSaleRepo) Create(_ context.Context, sale *domain.Sale) error {
	if m.err != nil {
		return m.err
	}
	sale.ID = int64(len(m.sales) + 1)
	m.sales = append(m.sales, *sale)
	return nil
}

func (m *mockSaleRepo) GetByID(_ context.Context, id int64) (*domain.Sale, error) {
	if m.err != nil {
		return nil, m.err
	}
	for _, s := range m.sales {
		if s.ID == id {
			return &s, nil
		}
	}
	return nil, domain.ErrNotFound
}

func (m *mockSaleRepo) List(_ context.Context, filter domain.SaleFilter) ([]domain.Sale, int, error) {
	if m.err != nil {
		return nil, 0, m.err
	}
	return m.sales, len(m.sales), nil
}

func (m *mockSaleRepo) Update(_ context.Context, sale *domain.Sale) error {
	if m.err != nil {
		return m.err
	}
	for i, s := range m.sales {
		if s.ID == sale.ID {
			m.sales[i] = *sale
			return nil
		}
	}
	return domain.ErrNotFound
}

func (m *mockSaleRepo) Delete(_ context.Context, id int64) error {
	if m.err != nil {
		return m.err
	}
	for i, s := range m.sales {
		if s.ID == id {
			m.sales = append(m.sales[:i], m.sales[i+1:]...)
			return nil
		}
	}
	return domain.ErrNotFound
}

func (m *mockSaleRepo) SetMaterials(_ context.Context, _ int64, _ []int64) error {
	return m.err
}
func (m *mockSaleRepo) SetBases(_ context.Context, _ int64, _ []int64) error {
	return m.err
}
func (m *mockSaleRepo) GetMaterialIDs(_ context.Context, _ int64) ([]int64, error) {
	return nil, m.err
}
func (m *mockSaleRepo) GetBaseIDs(_ context.Context, _ int64) ([]int64, error) {
	return nil, m.err
}

// mockFileRepo implements port.FileRepository.
type mockFileRepo struct {
	err error
}

func (m *mockFileRepo) Save(_ context.Context, _ string, _ io.Reader) error {
	return m.err
}
func (m *mockFileRepo) Delete(_ context.Context, _ string) error {
	return m.err
}
func (m *mockFileRepo) GetPath(_, _ string) string {
	return ""
}
func (m *mockFileRepo) Exists(_ context.Context, _ string) (bool, error) {
	return false, m.err
}
func (m *mockFileRepo) Copy(_ context.Context, _, _ string) error {
	return m.err
}
func (m *mockFileRepo) MkdirAll(_ context.Context, _ string) error {
	return m.err
}
func (m *mockFileRepo) RemoveDir(_ context.Context, _ string) error {
	return m.err
}
func (m *mockFileRepo) RenameDir(_ context.Context, _, _ string) error {
	return m.err
}

func newGalleryService(workRepo *mockWorkRepo, saleRepo *mockSaleRepo, fileRepo *mockFileRepo) *GalleryService {
	return NewGalleryService(workRepo, saleRepo, fileRepo, "/tmp/images")
}

var testWork = domain.Work{
	NameRu:   "Test Work RU",
	NameEn:   "Test Work EN",
	StrID:    "test-work",
	Width:    100,
	Height:   200,
	Year:     2024,
	BaseID:   1,
	DescrRu:  "Описание",
	DescrEn:  "Description",
	WorkPath: "test-work/",
	Images:   "img1.jpg",
}

var testSale = domain.Sale{
	Title:       "Test Sale",
	Description: "Sale description",
	Price:       1000,
	OldPrice:    1500,
	Status:      "available",
	Sold:        false,
}

func TestCreateWork(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		svc := newGalleryService(&mockWorkRepo{}, &mockSaleRepo{}, &mockFileRepo{})
		w := testWork
		err := svc.CreateWork(context.Background(), &w, "", nil)
		if err != nil {
			t.Fatalf("CreateWork failed: %v", err)
		}
		if w.ID != 1 {
			t.Errorf("expected work ID 1, got %d", w.ID)
		}
	})

	t.Run("empty name_ru returns ErrInvalidInput", func(t *testing.T) {
		svc := newGalleryService(&mockWorkRepo{}, &mockSaleRepo{}, &mockFileRepo{})
		w := domain.Work{NameRu: ""}
		err := svc.CreateWork(context.Background(), &w, "", nil)
		if !errors.Is(err, domain.ErrInvalidInput) {
			t.Errorf("expected ErrInvalidInput, got %v", err)
		}
	})

	t.Run("repo error", func(t *testing.T) {
		workRepo := &mockWorkRepo{err: errors.New("db error")}
		svc := newGalleryService(workRepo, &mockSaleRepo{}, &mockFileRepo{})
		err := svc.CreateWork(context.Background(), &domain.Work{NameRu: "Test"}, "", nil)
		if err == nil {
			t.Fatal("expected error, got nil")
		}
	})
}

func TestGetWorkByID(t *testing.T) {
	t.Run("found", func(t *testing.T) {
		workRepo := &mockWorkRepo{}
		svc := newGalleryService(workRepo, &mockSaleRepo{}, &mockFileRepo{})
		w := testWork
		if err := svc.CreateWork(context.Background(), &w, "", nil); err != nil {
			t.Fatalf("setup failed: %v", err)
		}
		got, err := svc.GetWorkByID(context.Background(), w.ID)
		if err != nil {
			t.Fatalf("GetWorkByID failed: %v", err)
		}
		if got.NameRu != testWork.NameRu {
			t.Errorf("NameRu = %q, want %q", got.NameRu, testWork.NameRu)
		}
	})

	t.Run("not found", func(t *testing.T) {
		svc := newGalleryService(&mockWorkRepo{}, &mockSaleRepo{}, &mockFileRepo{})
		_, err := svc.GetWorkByID(context.Background(), 999)
		if !errors.Is(err, domain.ErrNotFound) {
			t.Errorf("expected ErrNotFound, got %v", err)
		}
	})
}

func TestGetWorks(t *testing.T) {
	t.Run("empty list", func(t *testing.T) {
		svc := newGalleryService(&mockWorkRepo{}, &mockSaleRepo{}, &mockFileRepo{})
		works, total, err := svc.GetWorks(context.Background(), domain.WorkFilter{})
		if err != nil {
			t.Fatalf("GetWorks failed: %v", err)
		}
		if len(works) != 0 {
			t.Errorf("expected 0 works, got %d", len(works))
		}
		if total != 0 {
			t.Errorf("expected total 0, got %d", total)
		}
	})

	t.Run("multiple works", func(t *testing.T) {
		workRepo := &mockWorkRepo{}
		svc := newGalleryService(workRepo, &mockSaleRepo{}, &mockFileRepo{})
		for i := 0; i < 3; i++ {
			w := testWork
			w.NameRu = "Work"
			if err := svc.CreateWork(context.Background(), &w, "", nil); err != nil {
				t.Fatalf("setup failed: %v", err)
			}
		}
		works, total, err := svc.GetWorks(context.Background(), domain.WorkFilter{})
		if err != nil {
			t.Fatalf("GetWorks failed: %v", err)
		}
		if len(works) != 3 {
			t.Errorf("expected 3 works, got %d", len(works))
		}
		if total != 3 {
			t.Errorf("expected total 3, got %d", total)
		}
	})

	t.Run("repo error", func(t *testing.T) {
		workRepo := &mockWorkRepo{err: errors.New("db error")}
		svc := newGalleryService(workRepo, &mockSaleRepo{}, &mockFileRepo{})
		_, _, err := svc.GetWorks(context.Background(), domain.WorkFilter{})
		if err == nil {
			t.Fatal("expected error, got nil")
		}
	})
}

func TestUpdateWork(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		workRepo := &mockWorkRepo{}
		svc := newGalleryService(workRepo, &mockSaleRepo{}, &mockFileRepo{})
		w := testWork
		if err := svc.CreateWork(context.Background(), &w, "", nil); err != nil {
			t.Fatalf("setup failed: %v", err)
		}
		w.NameRu = "Updated"
		if err := svc.UpdateWork(context.Background(), &w, "", nil); err != nil {
			t.Fatalf("UpdateWork failed: %v", err)
		}
		got, _ := svc.GetWorkByID(context.Background(), w.ID)
		if got.NameRu != "Updated" {
			t.Errorf("NameRu = %q, want %q", got.NameRu, "Updated")
		}
	})

	t.Run("repo error", func(t *testing.T) {
		workRepo := &mockWorkRepo{err: errors.New("db error")}
		svc := newGalleryService(workRepo, &mockSaleRepo{}, &mockFileRepo{})
		w := domain.Work{ID: 1, NameRu: "Test"}
		err := svc.UpdateWork(context.Background(), &w, "", nil)
		if err == nil {
			t.Fatal("expected error, got nil")
		}
	})
}

func TestDeleteWork(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		workRepo := &mockWorkRepo{}
		svc := newGalleryService(workRepo, &mockSaleRepo{}, &mockFileRepo{})
		w := testWork
		if err := svc.CreateWork(context.Background(), &w, "", nil); err != nil {
			t.Fatalf("setup failed: %v", err)
		}
		if err := svc.DeleteWork(context.Background(), w.ID); err != nil {
			t.Fatalf("DeleteWork failed: %v", err)
		}
		_, err := svc.GetWorkByID(context.Background(), w.ID)
		if !errors.Is(err, domain.ErrNotFound) {
			t.Errorf("expected ErrNotFound after delete, got %v", err)
		}
	})

	t.Run("not found", func(t *testing.T) {
		svc := newGalleryService(&mockWorkRepo{}, &mockSaleRepo{}, &mockFileRepo{})
		err := svc.DeleteWork(context.Background(), 999)
		if !errors.Is(err, domain.ErrNotFound) {
			t.Errorf("expected ErrNotFound, got %v", err)
		}
	})

	t.Run("repo error", func(t *testing.T) {
		workRepo := &mockWorkRepo{err: errors.New("db error")}
		svc := newGalleryService(workRepo, &mockSaleRepo{}, &mockFileRepo{})
		err := svc.DeleteWork(context.Background(), 1)
		if err == nil {
			t.Fatal("expected error, got nil")
		}
	})
}

func TestCreateSale(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		svc := newGalleryService(&mockWorkRepo{}, &mockSaleRepo{}, &mockFileRepo{})
		s := testSale
		err := svc.CreateSale(context.Background(), &s, "", nil)
		if err != nil {
			t.Fatalf("CreateSale failed: %v", err)
		}
		if s.ID != 1 {
			t.Errorf("expected sale ID 1, got %d", s.ID)
		}
	})

	t.Run("empty title returns ErrInvalidInput", func(t *testing.T) {
		svc := newGalleryService(&mockWorkRepo{}, &mockSaleRepo{}, &mockFileRepo{})
		s := domain.Sale{Title: ""}
		err := svc.CreateSale(context.Background(), &s, "", nil)
		if !errors.Is(err, domain.ErrInvalidInput) {
			t.Errorf("expected ErrInvalidInput, got %v", err)
		}
	})
}

func TestGetSaleByID(t *testing.T) {
	t.Run("found", func(t *testing.T) {
		saleRepo := &mockSaleRepo{}
		svc := newGalleryService(&mockWorkRepo{}, saleRepo, &mockFileRepo{})
		s := testSale
		if err := svc.CreateSale(context.Background(), &s, "", nil); err != nil {
			t.Fatalf("setup failed: %v", err)
		}
		got, err := svc.GetSaleByID(context.Background(), s.ID)
		if err != nil {
			t.Fatalf("GetSaleByID failed: %v", err)
		}
		if got.Title != testSale.Title {
			t.Errorf("Title = %q, want %q", got.Title, testSale.Title)
		}
	})

	t.Run("not found", func(t *testing.T) {
		svc := newGalleryService(&mockWorkRepo{}, &mockSaleRepo{}, &mockFileRepo{})
		_, err := svc.GetSaleByID(context.Background(), 999)
		if !errors.Is(err, domain.ErrNotFound) {
			t.Errorf("expected ErrNotFound, got %v", err)
		}
	})
}

func TestGetSales(t *testing.T) {
	t.Run("empty list", func(t *testing.T) {
		svc := newGalleryService(&mockWorkRepo{}, &mockSaleRepo{}, &mockFileRepo{})
		sales, total, err := svc.GetSales(context.Background(), domain.SaleFilter{})
		if err != nil {
			t.Fatalf("GetSales failed: %v", err)
		}
		if len(sales) != 0 {
			t.Errorf("expected 0 sales, got %d", len(sales))
		}
		if total != 0 {
			t.Errorf("expected total 0, got %d", total)
		}
	})

	t.Run("multiple sales", func(t *testing.T) {
		saleRepo := &mockSaleRepo{}
		svc := newGalleryService(&mockWorkRepo{}, saleRepo, &mockFileRepo{})
		for i := 0; i < 3; i++ {
			s := testSale
			s.Title = "Sale"
			if err := svc.CreateSale(context.Background(), &s, "", nil); err != nil {
				t.Fatalf("setup failed: %v", err)
			}
		}
		sales, total, err := svc.GetSales(context.Background(), domain.SaleFilter{})
		if err != nil {
			t.Fatalf("GetSales failed: %v", err)
		}
		if len(sales) != 3 {
			t.Errorf("expected 3 sales, got %d", len(sales))
		}
		if total != 3 {
			t.Errorf("expected total 3, got %d", total)
		}
	})
}

func TestUpdateSale(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		saleRepo := &mockSaleRepo{}
		fileRepo := &mockFileRepo{}
		svc := newGalleryService(&mockWorkRepo{}, saleRepo, fileRepo)
		s := testSale
		if err := svc.CreateSale(context.Background(), &s, "test.jpg", nil); err != nil {
			t.Fatalf("setup failed: %v", err)
		}
		s.Title = "Updated Sale"
		s.ImagePath = "sale_image.jpg"
		if err := svc.UpdateSale(context.Background(), &s, "", nil); err != nil {
			t.Fatalf("UpdateSale failed: %v", err)
		}
		got, _ := svc.GetSaleByID(context.Background(), s.ID)
		if got.Title != "Updated Sale" {
			t.Errorf("Title = %q, want %q", got.Title, "Updated Sale")
		}
	})
}

func TestDeleteSale(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		saleRepo := &mockSaleRepo{}
		fileRepo := &mockFileRepo{}
		svc := newGalleryService(&mockWorkRepo{}, saleRepo, fileRepo)
		s := testSale
		if err := svc.CreateSale(context.Background(), &s, "test.jpg", nil); err != nil {
			t.Fatalf("setup failed: %v", err)
		}
		if err := svc.DeleteSale(context.Background(), s.ID); err != nil {
			t.Fatalf("DeleteSale failed: %v", err)
		}
		_, err := svc.GetSaleByID(context.Background(), s.ID)
		if !errors.Is(err, domain.ErrNotFound) {
			t.Errorf("expected ErrNotFound after delete, got %v", err)
		}
	})

	t.Run("not found", func(t *testing.T) {
		svc := newGalleryService(&mockWorkRepo{}, &mockSaleRepo{}, &mockFileRepo{})
		err := svc.DeleteSale(context.Background(), 999)
		if !errors.Is(err, domain.ErrNotFound) {
			t.Errorf("expected ErrNotFound, got %v", err)
		}
	})
}

func TestWorkByIDLogsNameRu(t *testing.T) {
	workRepo := &mockWorkRepo{}
	svc := newGalleryService(workRepo, &mockSaleRepo{}, &mockFileRepo{})
	w := testWork
	if err := svc.CreateWork(context.Background(), &w, "", nil); err != nil {
		t.Fatalf("setup failed: %v", err)
	}
	got, err := svc.GetWorkByID(context.Background(), w.ID)
	if err != nil {
		t.Fatalf("GetWorkByID failed: %v", err)
	}
	if got.NameRu != testWork.NameRu {
		t.Errorf("NameRu = %q, want %q", got.NameRu, testWork.NameRu)
	}
}

func TestServiceWorkFilterPassthrough(t *testing.T) {
	workRepo := &mockWorkRepo{}
	svc := newGalleryService(workRepo, &mockSaleRepo{}, &mockFileRepo{})

	works := []domain.Work{
		{NameRu: "Alpha", NameEn: "Alpha"},
		{NameRu: "Beta", NameEn: "Beta"},
		{NameRu: "Gamma", NameEn: "Gamma"},
	}
	for _, w := range works {
		cw := w
		svc.CreateWork(context.Background(), &cw, "", nil)
	}

	result, total, err := svc.GetWorks(context.Background(), domain.WorkFilter{Query: "Beta"})
	if err != nil {
		t.Fatalf("GetWorks failed: %v", err)
	}
	if total != 1 {
		t.Errorf("expected total 1, got %d", total)
	}
	if len(result) != 1 || result[0].NameRu != "Beta" {
		t.Errorf("expected [Beta], got %v", result)
	}
}

func TestServiceGetSalesPassthrough(t *testing.T) {
	_, _, err := newGalleryService(&mockWorkRepo{}, &mockSaleRepo{}, &mockFileRepo{}).GetSales(context.Background(), domain.SaleFilter{Status: "available"})
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestNewGalleryService(t *testing.T) {
	svc := NewGalleryService(&mockWorkRepo{}, &mockSaleRepo{}, &mockFileRepo{}, "/tmp")
	if svc == nil {
		t.Fatal("expected non-nil service")
	}
}

func TestCreateWorkMaxID(t *testing.T) {
	workRepo := &mockWorkRepo{}
	svc := newGalleryService(workRepo, &mockSaleRepo{}, &mockFileRepo{})

	for i := 0; i < 30; i++ {
		w := domain.Work{NameRu: "Test"}
		if err := svc.CreateWork(context.Background(), &w, "", nil); err != nil {
			t.Fatalf("CreateWork %d failed: %v", i, err)
		}
	}

	works, total, err := svc.GetWorks(context.Background(), domain.WorkFilter{})
	if err != nil {
		t.Fatalf("GetWorks failed: %v", err)
	}
	if total != 30 {
		t.Errorf("expected total 30, got %d", total)
	}
	if len(works) != 30 {
		t.Errorf("expected 30 works, got %d", len(works))
	}
}

func TestGalleryServiceImplementsPort(t *testing.T) {
	svc := NewGalleryService(&mockWorkRepo{}, &mockSaleRepo{}, &mockFileRepo{}, "/tmp")
	var _ port.GalleryService = svc
}

func TestWorkFieldsAfterCreate(t *testing.T) {
	workRepo := &mockWorkRepo{}
	svc := newGalleryService(workRepo, &mockSaleRepo{}, &mockFileRepo{})

	w := domain.Work{
		NameRu:   "Full Test",
		NameEn:   "Full Test EN",
		StrID:    "full-test",
		Width:    800,
		Height:   600,
		Year:     2023,
		BaseID:   2,
		DescrRu:  "Полное описание",
		DescrEn:  "Full description",
		WorkPath: "full-test/",
		Images:   "img1.jpg;img2.jpg",
	}

	if err := svc.CreateWork(context.Background(), &w, "", nil); err != nil {
		t.Fatalf("CreateWork failed: %v", err)
	}
	if w.ID == 0 {
		t.Error("expected non-zero ID after create")
	}

	got, err := svc.GetWorkByID(context.Background(), w.ID)
	if err != nil {
		t.Fatalf("GetWorkByID failed: %v", err)
	}
	if !reflect.DeepEqual(w, *got) {
		t.Errorf("work mismatch after round-trip:\nhave %+v\nwant %+v", *got, w)
	}
}
