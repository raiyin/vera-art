package handler

import (
	"io"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/raiyin/artserver/internal/domain"
	"github.com/raiyin/artserver/internal/dto"
	"github.com/raiyin/artserver/internal/port"
	"github.com/raiyin/artserver/pkg/apperror"
)

// NewsHandler handles news HTTP requests.
type NewsHandler struct {
	newsService port.NewsService
}

// NewNewsHandler creates a new NewsHandler.
func NewNewsHandler(newsService port.NewsService) *NewsHandler {
	return &NewsHandler{newsService: newsService}
}

// GetNews returns a list of news entries.
func (h *NewsHandler) GetNews(c *gin.Context) {
	filter := domain.NewsFilter{
		Status: c.Query("status"),
		Query:  c.Query("q"),
		Page:   ParseIntQuery(c, "page", 1),
		Limit:  ParseIntQuery(c, "limit", 20),
	}

	newsList, total, err := h.newsService.GetNews(c.Request.Context(), filter)
	if err != nil {
		slog.Error("GetNews: failed to list news",
			"error", err,
			"filter", filter,
		)
		apiErr := apperror.FromError(err)
		c.JSON(apiErr.Status, apiErr)
		return
	}

	responses := make([]dto.NewsResponse, len(newsList))
	for i, n := range newsList {
		responses[i] = dto.NewsResponse{
			ID:          n.ID,
			Title:       n.Title,
			Description: n.Description,
			Content:     n.Content,
			ImagePath:   n.ImagePath,
			VideoPath:   n.VideoPath,
			VideoPaths:  n.VideoPaths,
			ImagePaths:  n.ImagePaths,
			Status:      n.Status,
			CreatedAt:   n.CreatedAt,
			UpdatedAt:   n.UpdatedAt,
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"news":  responses,
		"total": total,
	})
}

// GetNewsByID returns a news entry by ID.
func (h *NewsHandler) GetNewsByID(c *gin.Context) {
	id, ok := ParseInt64Param(c, "id")
	if !ok {
		return
	}

	news, err := h.newsService.GetNewsByID(c.Request.Context(), id)
	if err != nil {
		slog.Error("GetNewsByID: failed to get news",
			"news_id", id,
			"error", err,
		)
		apiErr := apperror.FromError(err)
		c.JSON(apiErr.Status, apiErr)
		return
	}

	c.JSON(http.StatusOK, dto.NewsResponse{
		ID:          news.ID,
		Title:       news.Title,
		Description: news.Description,
		Content:     news.Content,
		ImagePath:   news.ImagePath,
		VideoPath:   news.VideoPath,
		VideoPaths:  news.VideoPaths,
		ImagePaths:  news.ImagePaths,
		Status:      news.Status,
		CreatedAt:   news.CreatedAt,
		UpdatedAt:   news.UpdatedAt,
	})
}

// CreateNews creates a new news entry.
func (h *NewsHandler) CreateNews(c *gin.Context) {
	var req dto.CreateNewsRequest
	if err := c.ShouldBind(&req); err != nil {
		slog.Warn("CreateNews: invalid request",
			"error", err,
		)
		c.JSON(http.StatusBadRequest, apperror.APIError{
			Status:  http.StatusBadRequest,
			Code:    "INVALID_REQUEST",
			Message: "Invalid request",
		})
		return
	}

	news := &domain.News{
		Title:       req.Title,
		Description: req.Description,
		Content:     req.Content,
		Status:      req.Status,
	}

	var imageFile *domain.UploadedFile
	var videoFile *domain.UploadedFile

	if imgFile, imgHeader, err := c.Request.FormFile("image"); err == nil {
		defer imgFile.Close()
		imageFile = &domain.UploadedFile{
			Filename: imgHeader.Filename,
			Reader:   imgFile,
		}
	}

	if vidFile, vidHeader, err := c.Request.FormFile("video"); err == nil {
		defer vidFile.Close()
		videoFile = &domain.UploadedFile{
			Filename: vidHeader.Filename,
			Reader:   vidFile,
		}
	}

	if err := h.newsService.CreateNews(c.Request.Context(), news, imageFile, videoFile); err != nil {
		slog.Error("CreateNews: failed to create news",
			"title", req.Title,
			"error", err,
		)
		apiErr := apperror.FromError(err)
		c.JSON(apiErr.Status, apiErr)
		return
	}

	slog.Info("News created successfully",
		"news_id", news.ID,
		"title", news.Title,
	)
	c.JSON(http.StatusCreated, dto.NewsResponse{
		ID:          news.ID,
		Title:       news.Title,
		Description: news.Description,
		Content:     news.Content,
		ImagePath:   news.ImagePath,
		VideoPath:   news.VideoPath,
		VideoPaths:  news.VideoPaths,
		ImagePaths:  news.ImagePaths,
		Status:      news.Status,
		CreatedAt:   news.CreatedAt,
		UpdatedAt:   news.UpdatedAt,
	})
}

// UpdateNews updates a news entry.
func (h *NewsHandler) UpdateNews(c *gin.Context) {
	id, ok := ParseInt64Param(c, "id")
	if !ok {
		return
	}

	var req dto.UpdateNewsRequest
	if err := c.ShouldBind(&req); err != nil {
		slog.Warn("UpdateNews: invalid request",
			"news_id", id,
			"error", err,
		)
		c.JSON(http.StatusBadRequest, apperror.APIError{
			Status:  http.StatusBadRequest,
			Code:    "INVALID_REQUEST",
			Message: "Invalid request",
		})
		return
	}

	news := &domain.News{
		ID:          id,
		Title:       req.Title,
		Description: req.Description,
		Content:     req.Content,
		Status:      req.Status,
	}

	var imageFile *domain.UploadedFile
	var videoFile *domain.UploadedFile

	if imgFile, imgHeader, err := c.Request.FormFile("image"); err == nil {
		defer imgFile.Close()
		imageFile = &domain.UploadedFile{
			Filename: imgHeader.Filename,
			Reader:   imgFile,
		}
	}

	if vidFile, vidHeader, err := c.Request.FormFile("video"); err == nil {
		defer vidFile.Close()
		videoFile = &domain.UploadedFile{
			Filename: vidHeader.Filename,
			Reader:   vidFile,
		}
	}

	if err := h.newsService.UpdateNews(c.Request.Context(), news, imageFile, videoFile); err != nil {
		slog.Error("UpdateNews: failed to update news",
			"news_id", id,
			"error", err,
		)
		apiErr := apperror.FromError(err)
		c.JSON(apiErr.Status, apiErr)
		return
	}

	slog.Info("News updated successfully",
		"news_id", id,
	)
	c.JSON(http.StatusOK, dto.NewsResponse{
		ID:          news.ID,
		Title:       news.Title,
		Description: news.Description,
		Content:     news.Content,
		ImagePath:   news.ImagePath,
		VideoPath:   news.VideoPath,
		VideoPaths:  news.VideoPaths,
		ImagePaths:  news.ImagePaths,
		Status:      news.Status,
		CreatedAt:   news.CreatedAt,
		UpdatedAt:   news.UpdatedAt,
	})
}

// DeleteNews deletes a news entry.
func (h *NewsHandler) DeleteNews(c *gin.Context) {
	id, ok := ParseInt64Param(c, "id")
	if !ok {
		return
	}

	if err := h.newsService.DeleteNews(c.Request.Context(), id); err != nil {
		slog.Error("DeleteNews: failed to delete news",
			"news_id", id,
			"error", err,
		)
		apiErr := apperror.FromError(err)
		c.JSON(apiErr.Status, apiErr)
		return
	}

	slog.Info("News deleted successfully",
		"news_id", id,
	)
	c.JSON(http.StatusOK, gin.H{"message": "News deleted successfully"})
}

// BulkDeleteNews deletes multiple news entries.
func (h *NewsHandler) BulkDeleteNews(c *gin.Context) {
	var req struct {
		IDs []int64 `json:"ids"`
	}
	if !BindJSON(c, &req) {
		return
	}

	if err := h.newsService.BulkDeleteNews(c.Request.Context(), req.IDs); err != nil {
		slog.Error("BulkDeleteNews: failed to bulk delete news",
			"error", err,
		)
		apiErr := apperror.FromError(err)
		c.JSON(apiErr.Status, apiErr)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "News deleted successfully"})
}

// Ensure io is used
var _ io.Reader
