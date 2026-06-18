package handler

import (
	"io"
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
	// imagePath is a directory name like "ajax/" or "gnome/"
	dir := strings.TrimSuffix(imagePath, "/")
	if dir == "" {
		return nil
	}
	fullDir := filepath.Join(h.absWorksDir, h.relWorksDir, dir)
	entries, err := os.ReadDir(fullDir)
	if err != nil {
		return nil
	}

	// Filter and sort image files
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

	// Sort numerically by extracting number before extension
	sort.Slice(images, func(i, j int) bool {
		ni, _ := strconv.Atoi(strings.TrimSuffix(images[i], filepath.Ext(images[i])))
		nj, _ := strconv.Atoi(strings.TrimSuffix(images[j], filepath.Ext(images[j])))
		return ni < nj
	})

	return images
}

// GetWorks returns a list of works.
func (h *GalleryHandler) GetWorks(c *gin.Context) {
	filter := domain.WorkFilter{
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

	works, total, err := h.galleryService.GetWorks(c.Request.Context(), filter)
	if err != nil {
		apiErr := apperror.FromError(err)
		c.JSON(apiErr.Status, apiErr)
		return
	}

	responses := make([]dto.WorkResponse, len(works))
	for i, w := range works {
		images := h.listImageFiles(w.ImagePath)
		responses[i] = dto.WorkResponse{
			ID:          w.ID,
			Title:       w.Title,
			Description: w.Description,
			ImagePath:   w.ImagePath,
			Images:      images,
			Year:        w.Year,
			Technique:   w.Technique,
			Size:        w.Size,
			Status:      w.Status,
			SortOrder:   w.SortOrder,
			MaterialIDs: w.MaterialIDs,
			BaseIDs:     w.BaseIDs,
			CreatedAt:   w.CreatedAt,
			UpdatedAt:   w.UpdatedAt,
		}
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
		apiErr := apperror.FromError(err)
		c.JSON(apiErr.Status, apiErr)
		return
	}

	images := h.listImageFiles(work.ImagePath)
	c.JSON(http.StatusOK, dto.WorkResponse{
		ID:          work.ID,
		Title:       work.Title,
		Description: work.Description,
		ImagePath:   work.ImagePath,
		Images:      images,
		Year:        work.Year,
		Technique:   work.Technique,
		Size:        work.Size,
		Status:      work.Status,
		SortOrder:   work.SortOrder,
		MaterialIDs: work.MaterialIDs,
		BaseIDs:     work.BaseIDs,
		CreatedAt:   work.CreatedAt,
		UpdatedAt:   work.UpdatedAt,
	})
}

// CreateWork creates a new work.
func (h *GalleryHandler) CreateWork(c *gin.Context) {
	var req dto.CreateWorkRequest
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, apperror.APIError{
			Status:  http.StatusBadRequest,
			Code:    "INVALID_REQUEST",
			Message: "Invalid request",
		})
		return
	}

	work := &domain.Work{
		Title:       req.Title,
		Description: req.Description,
		Year:        req.Year,
		Technique:   req.Technique,
		Size:        req.Size,
		Status:      req.Status,
		SortOrder:   req.SortOrder,
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

	if err := h.galleryService.CreateWork(c.Request.Context(), work, filename, reader); err != nil {
		apiErr := apperror.FromError(err)
		c.JSON(apiErr.Status, apiErr)
		return
	}

	c.JSON(http.StatusCreated, dto.WorkResponse{
		ID:          work.ID,
		Title:       work.Title,
		Description: work.Description,
		ImagePath:   work.ImagePath,
		Year:        work.Year,
		Technique:   work.Technique,
		Size:        work.Size,
		Status:      work.Status,
		SortOrder:   work.SortOrder,
		MaterialIDs: work.MaterialIDs,
		BaseIDs:     work.BaseIDs,
		CreatedAt:   work.CreatedAt,
		UpdatedAt:   work.UpdatedAt,
	})
}

// UpdateWork updates a work.
func (h *GalleryHandler) UpdateWork(c *gin.Context) {
	id, ok := ParseInt64Param(c, "id")
	if !ok {
		return
	}

	var req dto.UpdateWorkRequest
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, apperror.APIError{
			Status:  http.StatusBadRequest,
			Code:    "INVALID_REQUEST",
			Message: "Invalid request",
		})
		return
	}

	work := &domain.Work{
		ID:          id,
		Title:       req.Title,
		Description: req.Description,
		Year:        req.Year,
		Technique:   req.Technique,
		Size:        req.Size,
		Status:      req.Status,
		SortOrder:   req.SortOrder,
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

	if err := h.galleryService.UpdateWork(c.Request.Context(), work, filename, reader); err != nil {
		apiErr := apperror.FromError(err)
		c.JSON(apiErr.Status, apiErr)
		return
	}

	c.JSON(http.StatusOK, dto.UpdateWorkResponse{
		ID:          work.ID,
		Title:       work.Title,
		Description: work.Description,
		ImagePath:   work.ImagePath,
		Year:        work.Year,
		Technique:   work.Technique,
		Size:        work.Size,
		Status:      work.Status,
		SortOrder:   work.SortOrder,
		MaterialIDs: work.MaterialIDs,
		BaseIDs:     work.BaseIDs,
		CreatedAt:   work.CreatedAt,
		UpdatedAt:   work.UpdatedAt,
	})
}

// DeleteWork deletes a work.
func (h *GalleryHandler) DeleteWork(c *gin.Context) {
	id, ok := ParseInt64Param(c, "id")
	if !ok {
		return
	}

	if err := h.galleryService.DeleteWork(c.Request.Context(), id); err != nil {
		apiErr := apperror.FromError(err)
		c.JSON(apiErr.Status, apiErr)
		return
	}

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
		apiErr := apperror.FromError(err)
		c.JSON(apiErr.Status, apiErr)
		return
	}

	responses := make([]dto.SaleResponse, len(sales))
	for i, s := range sales {
		responses[i] = dto.SaleResponse{
			ID:          s.ID,
			Title:       s.Title,
			Description: s.Description,
			ImagePath:   s.ImagePath,
			Price:       s.Price,
			OldPrice:    s.OldPrice,
			Year:        s.Year,
			Technique:   s.Technique,
			Size:        s.Size,
			Status:      s.Status,
			SortOrder:   s.SortOrder,
			Sold:        s.Sold,
			MaterialIDs: s.MaterialIDs,
			BaseIDs:     s.BaseIDs,
			CreatedAt:   s.CreatedAt,
			UpdatedAt:   s.UpdatedAt,
		}
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
		apiErr := apperror.FromError(err)
		c.JSON(apiErr.Status, apiErr)
		return
	}

	c.JSON(http.StatusOK, dto.SaleResponse{
		ID:          sale.ID,
		Title:       sale.Title,
		Description: sale.Description,
		ImagePath:   sale.ImagePath,
		Price:       sale.Price,
		OldPrice:    sale.OldPrice,
		Year:        sale.Year,
		Technique:   sale.Technique,
		Size:        sale.Size,
		Status:      sale.Status,
		SortOrder:   sale.SortOrder,
		Sold:        sale.Sold,
		MaterialIDs: sale.MaterialIDs,
		BaseIDs:     sale.BaseIDs,
		CreatedAt:   sale.CreatedAt,
		UpdatedAt:   sale.UpdatedAt,
	})
}

// CreateSale creates a new sale.
func (h *GalleryHandler) CreateSale(c *gin.Context) {
	var req dto.CreateSaleRequest
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, apperror.APIError{
			Status:  http.StatusBadRequest,
			Code:    "INVALID_REQUEST",
			Message: "Invalid request",
		})
		return
	}

	sale := &domain.Sale{
		Title:       req.Title,
		Description: req.Description,
		Price:       req.Price,
		OldPrice:    req.OldPrice,
		Year:        req.Year,
		Technique:   req.Technique,
		Size:        req.Size,
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

	if err := h.galleryService.CreateSale(c.Request.Context(), sale, filename, reader); err != nil {
		apiErr := apperror.FromError(err)
		c.JSON(apiErr.Status, apiErr)
		return
	}

	c.JSON(http.StatusCreated, dto.SaleResponse{
		ID:          sale.ID,
		Title:       sale.Title,
		Description: sale.Description,
		ImagePath:   sale.ImagePath,
		Price:       sale.Price,
		OldPrice:    sale.OldPrice,
		Year:        sale.Year,
		Technique:   sale.Technique,
		Size:        sale.Size,
		Status:      sale.Status,
		SortOrder:   sale.SortOrder,
		Sold:        sale.Sold,
		MaterialIDs: sale.MaterialIDs,
		BaseIDs:     sale.BaseIDs,
		CreatedAt:   sale.CreatedAt,
		UpdatedAt:   sale.UpdatedAt,
	})
}

// UpdateSale updates a sale.
func (h *GalleryHandler) UpdateSale(c *gin.Context) {
	id, ok := ParseInt64Param(c, "id")
	if !ok {
		return
	}

	var req dto.UpdateSaleRequest
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, apperror.APIError{
			Status:  http.StatusBadRequest,
			Code:    "INVALID_REQUEST",
			Message: "Invalid request",
		})
		return
	}

	sale := &domain.Sale{
		ID:          id,
		Title:       req.Title,
		Description: req.Description,
		Price:       req.Price,
		OldPrice:    req.OldPrice,
		Year:        req.Year,
		Technique:   req.Technique,
		Size:        req.Size,
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
		apiErr := apperror.FromError(err)
		c.JSON(apiErr.Status, apiErr)
		return
	}

	c.JSON(http.StatusOK, dto.UpdateSaleResponse{
		ID:          sale.ID,
		Title:       sale.Title,
		Description: sale.Description,
		ImagePath:   sale.ImagePath,
		Price:       sale.Price,
		OldPrice:    sale.OldPrice,
		Year:        sale.Year,
		Technique:   sale.Technique,
		Size:        sale.Size,
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
		apiErr := apperror.FromError(err)
		c.JSON(apiErr.Status, apiErr)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Sale deleted successfully"})
}
