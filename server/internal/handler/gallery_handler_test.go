package handler

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"

	"github.com/raiyin/artserver/internal/domain"
	"github.com/raiyin/artserver/internal/dto"
	"github.com/raiyin/artserver/pkg/apperror"
)

// mockGalleryService implements port.GalleryService.
type mockGalleryService struct {
	works        []domain.Work
	sales        []domain.Sale
	getWorkByID  func(ctx context.Context, id int64) (*domain.Work, error)
	createWork   func(ctx context.Context, work *domain.Work, filename string, reader io.Reader) error
	updateWork   func(ctx context.Context, work *domain.Work, filename string, reader io.Reader) error
	deleteWork   func(ctx context.Context, id int64) error
	createSale   func(ctx context.Context, sale *domain.Sale, filename string, reader io.Reader) error
	updateSale   func(ctx context.Context, sale *domain.Sale, filename string, reader io.Reader) error
	deleteSale   func(ctx context.Context, id int64) error
}

func (m *mockGalleryService) GetWorks(_ context.Context, filter domain.WorkFilter) ([]domain.Work, int, error) {
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

func (m *mockGalleryService) GetWorkByID(_ context.Context, id int64) (*domain.Work, error) {
	if m.getWorkByID != nil {
		return m.getWorkByID(context.Background(), id)
	}
	for _, w := range m.works {
		if w.ID == id {
			return &w, nil
		}
	}
	return nil, domain.ErrNotFound
}

func (m *mockGalleryService) CreateWork(_ context.Context, work *domain.Work, _ string, _ io.Reader) error {
	if m.createWork != nil {
		return m.createWork(context.Background(), work, "", nil)
	}
	work.ID = int64(len(m.works) + 1)
	m.works = append(m.works, *work)
	return nil
}

func (m *mockGalleryService) UpdateWork(_ context.Context, work *domain.Work, _ string, _ io.Reader) error {
	if m.updateWork != nil {
		return m.updateWork(context.Background(), work, "", nil)
	}
	for i, w := range m.works {
		if w.ID == work.ID {
			m.works[i] = *work
			return nil
		}
	}
	return domain.ErrNotFound
}

func (m *mockGalleryService) DeleteWork(_ context.Context, id int64) error {
	if m.deleteWork != nil {
		return m.deleteWork(context.Background(), id)
	}
	for i, w := range m.works {
		if w.ID == id {
			m.works = append(m.works[:i], m.works[i+1:]...)
			return nil
		}
	}
	return domain.ErrNotFound
}

func (m *mockGalleryService) GetSales(_ context.Context, _ domain.SaleFilter) ([]domain.Sale, int, error) {
	return m.sales, len(m.sales), nil
}

func (m *mockGalleryService) GetSaleByID(_ context.Context, id int64) (*domain.Sale, error) {
	for _, s := range m.sales {
		if s.ID == id {
			return &s, nil
		}
	}
	return nil, domain.ErrNotFound
}

func (m *mockGalleryService) CreateSale(_ context.Context, sale *domain.Sale, _ string, _ io.Reader) error {
	if m.createSale != nil {
		return m.createSale(context.Background(), sale, "", nil)
	}
	sale.ID = int64(len(m.sales) + 1)
	m.sales = append(m.sales, *sale)
	return nil
}

func (m *mockGalleryService) UpdateSale(_ context.Context, sale *domain.Sale, _ string, _ io.Reader) error {
	if m.updateSale != nil {
		return m.updateSale(context.Background(), sale, "", nil)
	}
	for i, s := range m.sales {
		if s.ID == sale.ID {
			m.sales[i] = *sale
			return nil
		}
	}
	return domain.ErrNotFound
}

func (m *mockGalleryService) DeleteSale(_ context.Context, id int64) error {
	if m.deleteSale != nil {
		return m.deleteSale(context.Background(), id)
	}
	for i, s := range m.sales {
		if s.ID == id {
			m.sales = append(m.sales[:i], m.sales[i+1:]...)
			return nil
		}
	}
	return domain.ErrNotFound
}

func (m *mockGalleryService) BulkDeleteWorks(_ context.Context, ids []int64) error {
	for _, id := range ids {
		if err := m.DeleteWork(context.Background(), id); err != nil {
			return err
		}
	}
	return nil
}

func (m *mockGalleryService) BulkDeleteSales(_ context.Context, ids []int64) error {
	for _, id := range ids {
		if err := m.DeleteSale(context.Background(), id); err != nil {
			return err
		}
	}
	return nil
}

func setupGalleryRouter(h *GalleryHandler) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.GET("/api/works", h.GetWorks)
	r.GET("/api/works/:id", h.GetWorkByID)
	r.POST("/api/works", h.CreateWork)
	r.PUT("/api/works/:id", h.UpdateWork)
	r.DELETE("/api/works/:id", h.DeleteWork)
	r.GET("/api/sales", h.GetSales)
	r.GET("/api/sales/:id", h.GetSaleByID)
	r.POST("/api/sales", h.CreateSale)
	r.PUT("/api/sales/:id", h.UpdateSale)
	r.DELETE("/api/sales/:id", h.DeleteSale)
	return r
}

func newHandler(svc *mockGalleryService) *GalleryHandler {
	return NewGalleryHandler(svc, "", "", "", "")
}

func TestHandlerGetWorks(t *testing.T) {
	t.Run("empty list", func(t *testing.T) {
		svc := &mockGalleryService{}
		h := newHandler(svc)
		r := setupGalleryRouter(h)

		w := httptest.NewRecorder()
		req := httptest.NewRequest("GET", "/api/works", nil)
		r.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("expected 200, got %d", w.Code)
		}

		var resp struct {
			Works []dto.WorkResponse `json:"works"`
			Total int                `json:"total"`
		}
		if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
			t.Fatalf("unmarshal: %v", err)
		}
		if len(resp.Works) != 0 {
			t.Errorf("expected 0 works, got %d", len(resp.Works))
		}
		if resp.Total != 0 {
			t.Errorf("expected total 0, got %d", resp.Total)
		}
	})

	t.Run("with works", func(t *testing.T) {
		svc := &mockGalleryService{
			works: []domain.Work{
				{ID: 1, NameRu: "Work 1", NameEn: "Work 1"},
				{ID: 2, NameRu: "Work 2", NameEn: "Work 2"},
			},
		}
		h := newHandler(svc)
		r := setupGalleryRouter(h)

		w := httptest.NewRecorder()
		req := httptest.NewRequest("GET", "/api/works", nil)
		r.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("expected 200, got %d", w.Code)
		}

		var resp struct {
			Works []dto.WorkResponse `json:"works"`
			Total int                `json:"total"`
		}
		if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
			t.Fatalf("unmarshal: %v", err)
		}
		if len(resp.Works) != 2 {
			t.Errorf("expected 2 works, got %d", len(resp.Works))
		}
		if resp.Total != 2 {
			t.Errorf("expected total 2, got %d", resp.Total)
		}
	})
}

func TestHandlerGetWorkByID(t *testing.T) {
	t.Run("found", func(t *testing.T) {
		svc := &mockGalleryService{
			getWorkByID: func(_ context.Context, id int64) (*domain.Work, error) {
				return &domain.Work{ID: id, NameRu: "Found", NameEn: "Found"}, nil
			},
		}
		h := newHandler(svc)
		r := setupGalleryRouter(h)

		w := httptest.NewRecorder()
		req := httptest.NewRequest("GET", "/api/works/1", nil)
		r.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("expected 200, got %d", w.Code)
		}

		var resp dto.WorkResponse
		if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
			t.Fatalf("unmarshal: %v", err)
		}
		if resp.ID != 1 || resp.NameRu != "Found" {
			t.Errorf("got %+v, want ID=1 NameRu=Found", resp)
		}
	})

	t.Run("not found", func(t *testing.T) {
		svc := &mockGalleryService{
			getWorkByID: func(_ context.Context, id int64) (*domain.Work, error) {
				return nil, apperror.ErrNotFound
			},
		}
		h := newHandler(svc)
		r := setupGalleryRouter(h)

		w := httptest.NewRecorder()
		req := httptest.NewRequest("GET", "/api/works/999", nil)
		r.ServeHTTP(w, req)

		if w.Code != http.StatusNotFound {
			t.Errorf("expected 404, got %d", w.Code)
		}
	})

	t.Run("invalid id", func(t *testing.T) {
		svc := &mockGalleryService{}
		h := newHandler(svc)
		r := setupGalleryRouter(h)

		w := httptest.NewRecorder()
		req := httptest.NewRequest("GET", "/api/works/abc", nil)
		r.ServeHTTP(w, req)

		if w.Code != http.StatusBadRequest {
			t.Errorf("expected 400, got %d", w.Code)
		}
	})
}

func TestHandlerCreateWork(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		svc := &mockGalleryService{
			createWork: func(_ context.Context, work *domain.Work, _ string, _ io.Reader) error {
				work.ID = 1
				work.WorkPath = "test/"
				return nil
			},
		}
		h := newHandler(svc)
		r := setupGalleryRouter(h)

		body := `{"name_ru":"New Work","name_en":"New Work EN","width":100,"height":200,"year":2024}`
		w := httptest.NewRecorder()
		req := httptest.NewRequest("POST", "/api/works", strings.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		r.ServeHTTP(w, req)

		if w.Code != http.StatusCreated {
			t.Errorf("expected 201, got %d. Body: %s", w.Code, w.Body.String())
		}

		var resp dto.WorkResponse
		if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
			t.Fatalf("unmarshal: %v", err)
		}
		if resp.NameRu != "New Work" {
			t.Errorf("NameRu = %q, want %q", resp.NameRu, "New Work")
		}
	})

	t.Run("missing required fields", func(t *testing.T) {
		svc := &mockGalleryService{}
		h := newHandler(svc)
		r := setupGalleryRouter(h)

		body := `{"name_ru":""}`
		w := httptest.NewRecorder()
		req := httptest.NewRequest("POST", "/api/works", strings.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		r.ServeHTTP(w, req)

		if w.Code != http.StatusBadRequest {
			t.Errorf("expected 400, got %d", w.Code)
		}
	})

	t.Run("invalid JSON", func(t *testing.T) {
		svc := &mockGalleryService{}
		h := newHandler(svc)
		r := setupGalleryRouter(h)

		w := httptest.NewRecorder()
		req := httptest.NewRequest("POST", "/api/works", strings.NewReader("not-json"))
		req.Header.Set("Content-Type", "application/json")
		r.ServeHTTP(w, req)

		if w.Code != http.StatusBadRequest {
			t.Errorf("expected 400, got %d", w.Code)
		}
	})

	t.Run("service error", func(t *testing.T) {
		svc := &mockGalleryService{
			createWork: func(_ context.Context, _ *domain.Work, _ string, _ io.Reader) error {
				return apperror.ErrInvalidInput
			},
		}
		h := newHandler(svc)
		r := setupGalleryRouter(h)

		body := `{"name_ru":"Test","name_en":"Test"}`
		w := httptest.NewRecorder()
		req := httptest.NewRequest("POST", "/api/works", strings.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		r.ServeHTTP(w, req)

		if w.Code != http.StatusBadRequest {
			t.Errorf("expected 400, got %d", w.Code)
		}
	})
}

func TestHandlerUpdateWork(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		svc := &mockGalleryService{
			getWorkByID: func(_ context.Context, id int64) (*domain.Work, error) {
				return &domain.Work{ID: id, NameRu: "Old", NameEn: "Old"}, nil
			},
			updateWork: func(_ context.Context, work *domain.Work, _ string, _ io.Reader) error {
				return nil
			},
		}
		h := newHandler(svc)
		r := setupGalleryRouter(h)

		body := `{"name_ru":"Updated","name_en":"Updated EN","width":200,"height":300}`
		w := httptest.NewRecorder()
		req := httptest.NewRequest("PUT", "/api/works/1", strings.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		r.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("expected 200, got %d. Body: %s", w.Code, w.Body.String())
		}
	})

	t.Run("invalid id", func(t *testing.T) {
		svc := &mockGalleryService{}
		h := newHandler(svc)
		r := setupGalleryRouter(h)

		body := `{"name_ru":"Test","name_en":"Test"}`
		w := httptest.NewRecorder()
		req := httptest.NewRequest("PUT", "/api/works/abc", strings.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		r.ServeHTTP(w, req)

		if w.Code != http.StatusBadRequest {
			t.Errorf("expected 400, got %d", w.Code)
		}
	})
}

func TestHandlerDeleteWork(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		svc := &mockGalleryService{
			deleteWork: func(_ context.Context, id int64) error {
				return nil
			},
		}
		h := newHandler(svc)
		r := setupGalleryRouter(h)

		w := httptest.NewRecorder()
		req := httptest.NewRequest("DELETE", "/api/works/1", nil)
		r.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("expected 200, got %d", w.Code)
		}

		var resp map[string]string
		if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
			t.Fatalf("unmarshal: %v", err)
		}
		if resp["message"] != "Work deleted successfully" {
			t.Errorf("message = %q, want %q", resp["message"], "Work deleted successfully")
		}
	})

	t.Run("not found", func(t *testing.T) {
		svc := &mockGalleryService{
			deleteWork: func(_ context.Context, id int64) error {
				return apperror.ErrNotFound
			},
		}
		h := newHandler(svc)
		r := setupGalleryRouter(h)

		w := httptest.NewRecorder()
		req := httptest.NewRequest("DELETE", "/api/works/999", nil)
		r.ServeHTTP(w, req)

		if w.Code != http.StatusNotFound {
			t.Errorf("expected 404, got %d", w.Code)
		}
	})
}

func TestHandlerWorkToResponseConversion(t *testing.T) {
	w := &domain.Work{
		ID:       1,
		StrID:    "test",
		Width:    100,
		Height:   200,
		Year:     2024,
		NameRu:   "Тест",
		NameEn:   "Test",
		BaseID:   5,
		DescrRu:  "Описание",
		DescrEn:  "Description",
		WorkPath: "test/",
		Images:   "img1.jpg;img2.jpg",
	}

	h := NewGalleryHandler(nil, "", "/content/works/", "", "")
	resp := h.workToResponse(w)
	if resp.ID != w.ID {
		t.Errorf("ID = %d, want %d", resp.ID, w.ID)
	}
	if resp.NameRu != w.NameRu {
		t.Errorf("NameRu = %q, want %q", resp.NameRu, w.NameRu)
	}
	if len(resp.Images) != 2 {
		t.Errorf("expected 2 images, got %d", len(resp.Images))
	}
	if resp.Images[0] != "img1.jpg" {
		t.Errorf("Images[0] = %q, want %q", resp.Images[0], "img1.jpg")
	}
	if resp.Dir != "/content/works/test/" {
		t.Errorf("Dir = %q, want %q", resp.Dir, "/content/works/test/")
	}
}

func TestHandlerGetWorksQueryParsing(t *testing.T) {
	svc := &mockGalleryService{
		works: []domain.Work{
			{ID: 1, NameRu: "Apple", NameEn: "Apple"},
			{ID: 2, NameRu: "Banana", NameEn: "Banana"},
			{ID: 3, NameRu: "Cherry", NameEn: "Cherry"},
		},
	}
	h := newHandler(svc)
	r := setupGalleryRouter(h)

	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/api/works?q=Apple", nil)
	r.ServeHTTP(w, req)

	var resp struct {
		Works []dto.WorkResponse `json:"works"`
		Total int                `json:"total"`
	}
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if len(resp.Works) != 1 {
		t.Errorf("expected 1 filtered work, got %d", len(resp.Works))
	}
}

func TestHandlerGetWorksDefaultPagination(t *testing.T) {
	svc := &mockGalleryService{}
	h := newHandler(svc)
	r := setupGalleryRouter(h)

	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/api/works", nil)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}
}

func TestHandlerCreateWorkMultipartForm(t *testing.T) {
	svc := &mockGalleryService{
		createWork: func(_ context.Context, work *domain.Work, _ string, _ io.Reader) error {
			work.ID = 1
			return nil
		},
	}
	h := newHandler(svc)
	r := setupGalleryRouter(h)

	body := `--boundary
Content-Disposition: form-data; name="name_ru"

Test Work
--boundary
Content-Disposition: form-data; name="name_en"

Test Work EN
--boundary--
`
	w := httptest.NewRecorder()
	req := httptest.NewRequest("POST", "/api/works", strings.NewReader(body))
	req.Header.Set("Content-Type", "multipart/form-data; boundary=boundary")
	r.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("expected 201, got %d. Body: %s", w.Code, w.Body.String())
	}
}
