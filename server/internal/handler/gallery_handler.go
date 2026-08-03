package handler

import (
	"io"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/raiyin/artserver/internal/domain"
	"github.com/raiyin/artserver/internal/dto"
	"github.com/raiyin/artserver/internal/port"
	"github.com/raiyin/artserver/pkg/apperror"
)

// GalleryHandler handles gallery (works + sales) HTTP requests.
type GalleryHandler struct {
	galleryService port.GalleryService
	absWorksDir    string
	relWorksDir    string
	absSalesDir    string
	relSalesDir    string
}

// NewGalleryHandler creates a new GalleryHandler.
func NewGalleryHandler(
	galleryService port.GalleryService,
	absWorksDir, relWorksDir, absSalesDir, relSalesDir string,
) *GalleryHandler {
	return &GalleryHandler{
		galleryService: galleryService,
		absWorksDir:    absWorksDir,
		relWorksDir:    relWorksDir,
		absSalesDir:    absSalesDir,
		relSalesDir:    relSalesDir,
	}
}

// listImageFiles reads the image directory for a work and returns sorted filenames.
func (h *GalleryHandler) listImageFiles(imagePath string) []string {
	dir := strings.TrimSuffix(imagePath, "/")
	if dir == "" {
		return nil
	}
	fullDir := filepath.Join(h.absWorksDir, h.relWorksDir, dir)
	entries, err := os.ReadDir(fullDir)
	if err != nil {
		slog.Debug("listImageFiles: could not read directory",
			"dir", fullDir,
			"error", err,
		)
		return nil
	}

	var images []string
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		name := strings.ToLower(entry.Name())
		if strings.HasSuffix(name, ".jpg") || strings.HasSuffix(name, ".jpeg") ||
			strings.HasSuffix(name, ".png") || strings.HasSuffix(name, ".webp") ||
			strings.HasSuffix(name, ".gif") {
			images = append(images, entry.Name())
		}
	}

	sort.Slice(images, func(i, j int) bool {
		ni, _ := strconv.Atoi(strings.TrimSuffix(images[i], filepath.Ext(images[i])))
		nj, _ := strconv.Atoi(strings.TrimSuffix(images[j], filepath.Ext(images[j])))
		return ni < nj
	})

	return images
}

func (h *GalleryHandler) workToResponse(w *domain.Work) dto.WorkResponse {
	return dto.WorkResponse{
		ID:          w.ID,
		StrID:       w.StrID,
		Width:       w.Width,
		Height:      w.Height,
		Year:        w.Year,
		NameRu:      w.NameRu,
		NameEn:      w.NameEn,
		BaseID:      w.BaseID,
		DescrRu:     w.DescrRu,
		DescrEn:     w.DescrEn,
		WorkPath:    w.WorkPath,
		Dir:         h.relWorksDir + w.WorkPath,
		Images:      dto.SplitImages(w.Images),
		MaterialIDs: w.MaterialIDs,
	}
}

func (h *GalleryHandler) saleToResponse(s *domain.Sale) dto.SaleResponse {
	images := dto.SplitImages(s.ImagePath)
	dir := h.relSalesDir
	if s.SalePath != "" {
		dir = h.relSalesDir + s.SalePath + "/"
	}
	return dto.SaleResponse{
		ID:          s.ID,
		NameRu:      s.NameRu,
		NameEn:      s.NameEn,
		DescrRu:     s.DescrRu,
		DescrEn:     s.DescrEn,
		ImagePath:   s.ImagePath,
		SalePath:    s.SalePath,
		Dir:         dir,
		Images:      images,
		Price:       s.Price,
		Year:        s.Year,
		Technique:   s.Technique,
		Width:       s.Width,
		Height:      s.Height,
		Status:      s.Status,
		SortOrder:   s.SortOrder,
		Sold:        s.Sold,
		MaterialIDs: s.MaterialIDs,
		BaseIDs:     s.BaseIDs,
		CreatedAt:   s.CreatedAt,
		UpdatedAt:   s.UpdatedAt,
	}
}

// GetWorks returns a list of works.
func (h *GalleryHandler) GetWorks(c *gin.Context) {
	filter := domain.WorkFilter{
		Query: c.Query("q"),
		Page:  ParseIntQuery(c, "page", 1),
		Limit: ParseIntQuery(c, "limit", 20),
	}

	works, total, err := h.galleryService.GetWorks(c.Request.Context(), filter)
	if err != nil {
		slog.Error("GetWorks: failed to list works",
			"error", err,
			"filter", filter,
		)
		apiErr := apperror.FromError(err)
		c.JSON(apiErr.Status, apiErr)
		return
	}

	responses := make([]dto.WorkResponse, len(works))
	for i, w := range works {
		responses[i] = h.workToResponse(&w)
	}

	c.JSON(http.StatusOK, gin.H{
		"works": responses,
		"total": total,
	})
}

// GetWorkByID returns a work by ID.
func (h *GalleryHandler) GetWorkByID(c *gin.Context) {
	id, ok := ParseInt64Param(c, "id")
	if !ok {
		return
	}

	work, err := h.galleryService.GetWorkByID(c.Request.Context(), id)
	if err != nil {
		slog.Error("GetWorkByID: failed to get work",
			"work_id", id,
			"error", err,
		)
		apiErr := apperror.FromError(err)
		c.JSON(apiErr.Status, apiErr)
		return
	}

	c.JSON(http.StatusOK, h.workToResponse(work))
}

// CreateWork creates a new work.
func (h *GalleryHandler) CreateWork(c *gin.Context) {
	var req dto.CreateWorkRequest
	if err := c.ShouldBind(&req); err != nil {
		slog.Warn("CreateWork: invalid request",
			"error", err,
		)
		c.JSON(http.StatusBadRequest, apperror.APIError{
			Status:  http.StatusBadRequest,
			Code:    "INVALID_REQUEST",
			Message: "Invalid request",
		})
		return
	}

	work := &domain.Work{
		StrID:       req.StrID,
		Width:       req.Width,
		Height:      req.Height,
		Year:        req.Year,
		NameRu:      req.NameRu,
		NameEn:      req.NameEn,
		BaseID:      req.BaseID,
		DescrRu:     req.DescrRu,
		DescrEn:     req.DescrEn,
		MaterialIDs: req.MaterialIDs,
	}

	var files []domain.UploadedFile
	if c.Request.MultipartForm != nil {
		for _, headers := range c.Request.MultipartForm.File["image"] {
			file, err := headers.Open()
			if err != nil {
				continue
			}
			defer file.Close()
			files = append(files, domain.UploadedFile{
				Filename: headers.Filename,
				Size:     headers.Size,
				Reader:   file,
			})
		}
	}

	if err := h.galleryService.CreateWork(c.Request.Context(), work, files); err != nil {
		slog.Error("CreateWork: failed to create work",
			"name_ru", req.NameRu,
			"error", err,
		)
		apiErr := apperror.FromError(err)
		c.JSON(apiErr.Status, apiErr)
		return
	}

	slog.Info("Work created successfully",
		"work_id", work.ID,
		"name_ru", work.NameRu,
	)
	c.JSON(http.StatusCreated, h.workToResponse(work))
}

// UpdateWork updates a work.
func (h *GalleryHandler) UpdateWork(c *gin.Context) {
	id, ok := ParseInt64Param(c, "id")
	if !ok {
		return
	}

	var req dto.UpdateWorkRequest
	if err := c.ShouldBind(&req); err != nil {
		slog.Warn("UpdateWork: invalid request",
			"work_id", id,
			"error", err,
		)
		c.JSON(http.StatusBadRequest, apperror.APIError{
			Status:  http.StatusBadRequest,
			Code:    "INVALID_REQUEST",
			Message: "Invalid request",
		})
		return
	}

	work := &domain.Work{
		ID:          id,
		StrID:       req.StrID,
		Width:       req.Width,
		Height:      req.Height,
		Year:        req.Year,
		NameRu:      req.NameRu,
		NameEn:      req.NameEn,
		BaseID:      req.BaseID,
		DescrRu:     req.DescrRu,
		DescrEn:     req.DescrEn,
		MaterialIDs: req.MaterialIDs,
		Images:      dto.JoinImages(req.Images),
	}

	var files []domain.UploadedFile
	if c.Request.MultipartForm != nil {
		for _, headers := range c.Request.MultipartForm.File["image"] {
			file, err := headers.Open()
			if err != nil {
				continue
			}
			defer file.Close()
			files = append(files, domain.UploadedFile{
				Filename: headers.Filename,
				Size:     headers.Size,
				Reader:   file,
			})
		}
	}

	if err := h.galleryService.UpdateWork(c.Request.Context(), work, files); err != nil {
		slog.Error("UpdateWork: failed to update work",
			"work_id", id,
			"error", err,
		)
		apiErr := apperror.FromError(err)
		c.JSON(apiErr.Status, apiErr)
		return
	}

	slog.Info("Work updated successfully",
		"work_id", id,
	)
		c.JSON(http.StatusOK, dto.UpdateWorkResponse{
			ID:          work.ID,
			StrID:       work.StrID,
			Width:       work.Width,
			Height:      work.Height,
			Year:        work.Year,
			NameRu:      work.NameRu,
			NameEn:      work.NameEn,
			BaseID:      work.BaseID,
			DescrRu:     work.DescrRu,
			DescrEn:     work.DescrEn,
			WorkPath:    work.WorkPath,
			Dir:         h.relWorksDir + work.WorkPath,
			Images:      dto.SplitImages(work.Images),
			MaterialIDs: work.MaterialIDs,
		})
}

// DeleteWork deletes a work.
func (h *GalleryHandler) DeleteWork(c *gin.Context) {
	id, ok := ParseInt64Param(c, "id")
	if !ok {
		return
	}

	if err := h.galleryService.DeleteWork(c.Request.Context(), id); err != nil {
		slog.Error("DeleteWork: failed to delete work",
			"work_id", id,
			"error", err,
		)
		apiErr := apperror.FromError(err)
		c.JSON(apiErr.Status, apiErr)
		return
	}

	slog.Info("Work deleted successfully",
		"work_id", id,
	)
	c.JSON(http.StatusOK, gin.H{"message": "Work deleted successfully"})
}

// GetSales returns a list of sales.
func (h *GalleryHandler) GetSales(c *gin.Context) {
	filter := domain.SaleFilter{
		Status:    c.Query("status"),
		Query:     c.Query("q"),
		SortBy:    c.Query("sort_by"),
		SortOrder: c.Query("sort_order"),
		Page:      ParseIntQuery(c, "page", 1),
		Limit:     ParseIntQuery(c, "limit", 20),
	}

	if materialID, ok := ParseInt64Query(c, "material_id"); ok && materialID > 0 {
		filter.MaterialID = materialID
	}
	if baseID, ok := ParseInt64Query(c, "base_id"); ok && baseID > 0 {
		filter.BaseID = baseID
	}

	sales, total, err := h.galleryService.GetSales(c.Request.Context(), filter)
	if err != nil {
		slog.Error("GetSales: failed to list sales",
			"error", err,
			"filter", filter,
		)
		apiErr := apperror.FromError(err)
		c.JSON(apiErr.Status, apiErr)
		return
	}

	responses := make([]dto.SaleResponse, len(sales))
	for i, s := range sales {
		responses[i] = h.saleToResponse(&s)
	}

	c.JSON(http.StatusOK, gin.H{
		"sales": responses,
		"total": total,
	})
}

// GetSaleByID returns a sale by ID.
func (h *GalleryHandler) GetSaleByID(c *gin.Context) {
	id, ok := ParseInt64Param(c, "id")
	if !ok {
		return
	}

	sale, err := h.galleryService.GetSaleByID(c.Request.Context(), id)
	if err != nil {
		slog.Error("GetSaleByID: failed to get sale",
			"sale_id", id,
			"error", err,
		)
		apiErr := apperror.FromError(err)
		c.JSON(apiErr.Status, apiErr)
		return
	}

	c.JSON(http.StatusOK, h.saleToResponse(sale))
}

// CreateSale creates a new sale.
func (h *GalleryHandler) CreateSale(c *gin.Context) {
	var req dto.CreateSaleRequest
	if err := c.ShouldBind(&req); err != nil {
		slog.Warn("CreateSale: invalid request",
			"error", err,
		)
		c.JSON(http.StatusBadRequest, apperror.APIError{
			Status:  http.StatusBadRequest,
			Code:    "INVALID_REQUEST",
			Message: "Invalid request",
		})
		return
	}

	sale := &domain.Sale{
		NameRu:      req.NameRu,
		NameEn:      req.NameEn,
		DescrRu:     req.DescrRu,
		DescrEn:     req.DescrEn,
		Price:       req.Price,
		Year:        req.Year,
		Technique:   req.Technique,
		Width:       req.Width,
		Height:      req.Height,
		Status:      req.Status,
		SortOrder:   req.SortOrder,
		Sold:        req.Sold,
		MaterialIDs: req.MaterialIDs,
		BaseIDs:     req.BaseIDs,
	}

	file, header, err := c.Request.FormFile("image")
	if err != nil {
		slog.Warn("CreateSale: missing image file")
		c.JSON(http.StatusBadRequest, apperror.APIError{
			Status:  http.StatusBadRequest,
			Code:    "MISSING_IMAGE",
			Message: "Image file is required",
		})
		return
	}
	defer file.Close()
	filename := header.Filename
	reader := file

	if err := h.galleryService.CreateSale(c.Request.Context(), sale, filename, reader); err != nil {
		slog.Error("CreateSale: failed to create sale",
			"name_ru", req.NameRu,
			"error", err,
		)
		apiErr := apperror.FromError(err)
		c.JSON(apiErr.Status, apiErr)
		return
	}

	slog.Info("Sale created successfully",
		"sale_id", sale.ID,
		"name_ru", sale.NameRu,
	)
	c.JSON(http.StatusCreated, h.saleToResponse(sale))
}

// UpdateSale updates a sale.
func (h *GalleryHandler) UpdateSale(c *gin.Context) {
	id, ok := ParseInt64Param(c, "id")
	if !ok {
		return
	}

	var req dto.UpdateSaleRequest
	if err := c.ShouldBind(&req); err != nil {
		slog.Warn("UpdateSale: invalid request",
			"sale_id", id,
			"error", err,
		)
		c.JSON(http.StatusBadRequest, apperror.APIError{
			Status:  http.StatusBadRequest,
			Code:    "INVALID_REQUEST",
			Message: "Invalid request",
		})
		return
	}

	sale := &domain.Sale{
		ID:          id,
		NameRu:      req.NameRu,
		NameEn:      req.NameEn,
		DescrRu:     req.DescrRu,
		DescrEn:     req.DescrEn,
		Price:       req.Price,
		Year:        req.Year,
		Technique:   req.Technique,
		Width:       req.Width,
		Height:      req.Height,
		Status:      req.Status,
		SortOrder:   req.SortOrder,
		Sold:        req.Sold,
		MaterialIDs: req.MaterialIDs,
		BaseIDs:     req.BaseIDs,
	}

	var filename string
	var reader io.ReadCloser

	file, header, err := c.Request.FormFile("image")
	if err == nil {
		defer file.Close()
		filename = header.Filename
		reader = file
	}

	if err := h.galleryService.UpdateSale(c.Request.Context(), sale, filename, reader); err != nil {
		slog.Error("UpdateSale: failed to update sale",
			"sale_id", id,
			"error", err,
		)
		apiErr := apperror.FromError(err)
		c.JSON(apiErr.Status, apiErr)
		return
	}

	slog.Info("Sale updated successfully",
		"sale_id", id,
	)
	images := dto.SplitImages(sale.ImagePath)
	dir := h.relSalesDir
	if sale.SalePath != "" {
		dir = h.relSalesDir + sale.SalePath + "/"
	}
	c.JSON(http.StatusOK, dto.UpdateSaleResponse{
		ID:          sale.ID,
		NameRu:      sale.NameRu,
		NameEn:      sale.NameEn,
		DescrRu:     sale.DescrRu,
		DescrEn:     sale.DescrEn,
		ImagePath:   sale.ImagePath,
		SalePath:    sale.SalePath,
		Dir:         dir,
		Images:      images,
		Price:       sale.Price,
		Year:        sale.Year,
		Technique:   sale.Technique,
		Width:       sale.Width,
		Height:      sale.Height,
		Status:      sale.Status,
		SortOrder:   sale.SortOrder,
		Sold:        sale.Sold,
		MaterialIDs: sale.MaterialIDs,
		BaseIDs:     sale.BaseIDs,
		CreatedAt:   sale.CreatedAt,
		UpdatedAt:   sale.UpdatedAt,
	})
}

// DeleteSale deletes a sale.
func (h *GalleryHandler) DeleteSale(c *gin.Context) {
	id, ok := ParseInt64Param(c, "id")
	if !ok {
		return
	}

	if err := h.galleryService.DeleteSale(c.Request.Context(), id); err != nil {
		slog.Error("DeleteSale: failed to delete sale",
			"sale_id", id,
			"error", err,
		)
		apiErr := apperror.FromError(err)
		c.JSON(apiErr.Status, apiErr)
		return
	}

	slog.Info("Sale deleted successfully",
		"sale_id", id,
	)
	c.JSON(http.StatusOK, gin.H{"message": "Sale deleted successfully"})
}

// BulkDeleteWorks deletes multiple works.
func (h *GalleryHandler) BulkDeleteWorks(c *gin.Context) {
	var req struct {
		IDs []int64 `json:"ids"`
	}
	if !BindJSON(c, &req) {
		return
	}

	if err := h.galleryService.BulkDeleteWorks(c.Request.Context(), req.IDs); err != nil {
		slog.Error("BulkDeleteWorks: failed to bulk delete works",
			"error", err,
		)
		apiErr := apperror.FromError(err)
		c.JSON(apiErr.Status, apiErr)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Works deleted successfully"})
}

// BulkDeleteSales deletes multiple sales.
func (h *GalleryHandler) BulkDeleteSales(c *gin.Context) {
	var req struct {
		IDs []int64 `json:"ids"`
	}
	if !BindJSON(c, &req) {
		return
	}

	if err := h.galleryService.BulkDeleteSales(c.Request.Context(), req.IDs); err != nil {
		slog.Error("BulkDeleteSales: failed to bulk delete sales",
			"error", err,
		)
		apiErr := apperror.FromError(err)
		c.JSON(apiErr.Status, apiErr)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Sales deleted successfully"})
}
