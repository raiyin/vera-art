package handler

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/raiyin/artserver/internal/domain"
	"github.com/raiyin/artserver/internal/dto"
	"github.com/raiyin/artserver/internal/port"
	"github.com/raiyin/artserver/pkg/apperror"
)

type NewsHandler struct {
	newsService port.NewsService
}

func NewNewsHandler(newsService port.NewsService) *NewsHandler {
	return &NewsHandler{newsService: newsService}
}

func (h *NewsHandler) GetNews(c *gin.Context) {
	filter := domain.NewsFilter{
		Status: c.Query("status"),
		Query:  c.Query("q"),
		Page:   ParseIntQuery(c, "page", 1),
		Limit:  ParseIntQuery(c, "limit", 20),
	}

	newsList, total, err := h.newsService.GetNews(c.Request.Context(), filter)
	if err != nil {
		slog.Error("GetNews: failed to list news", "error", err, "filter", filter)
		apiErr := apperror.FromError(err)
		c.JSON(apiErr.Status, apiErr)
		return
	}

	responses := make([]dto.NewsResponse, len(newsList))
	for i, n := range newsList {
		responses[i] = newsToResponse(&n)
	}

	c.JSON(http.StatusOK, gin.H{
		"news":  responses,
		"total": total,
	})
}

func (h *NewsHandler) GetNewsByID(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, apperror.APIError{
			Status:  http.StatusBadRequest,
			Code:    "INVALID_PARAM",
			Message: "Invalid id",
		})
		return
	}

	news, err := h.newsService.GetNewsByID(c.Request.Context(), id)
	if err != nil {
		slog.Error("GetNewsByID: failed to get news", "news_id", id, "error", err)
		apiErr := apperror.FromError(err)
		c.JSON(apiErr.Status, apiErr)
		return
	}

	c.JSON(http.StatusOK, newsToResponse(news))
}

func (h *NewsHandler) CreateNews(c *gin.Context) {
	dataJSON := c.PostForm("data")
	if dataJSON == "" {
		c.JSON(http.StatusBadRequest, apperror.APIError{
			Status:  http.StatusBadRequest,
			Code:    "INVALID_REQUEST",
			Message: "Missing data field",
		})
		return
	}

	var payload dto.NewsDataPayload
	if err := json.Unmarshal([]byte(dataJSON), &payload); err != nil {
		slog.Warn("CreateNews: invalid data JSON", "error", err)
		c.JSON(http.StatusBadRequest, apperror.APIError{
			Status:  http.StatusBadRequest,
			Code:    "INVALID_REQUEST",
			Message: "Invalid data JSON",
		})
		return
	}

	form, err := c.MultipartForm()
	if err != nil {
		slog.Warn("CreateNews: invalid multipart form", "error", err)
		c.JSON(http.StatusBadRequest, apperror.APIError{
			Status:  http.StatusBadRequest,
			Code:    "INVALID_REQUEST",
			Message: "Invalid multipart form",
		})
		return
	}

	var imgBackFile, imgBackfullFile *domain.UploadedFile
	if files := form.File["img_back"]; len(files) > 0 {
		f, err := files[0].Open()
		if err == nil {
			defer f.Close()
			imgBackFile = &domain.UploadedFile{Filename: files[0].Filename, Reader: f}
		}
	}
	if files := form.File["img_backfull"]; len(files) > 0 {
		f, err := files[0].Open()
		if err == nil {
			defer f.Close()
			imgBackfullFile = &domain.UploadedFile{Filename: files[0].Filename, Reader: f}
		}
	}

	var imageFiles []domain.UploadedFile
	if files := form.File["images"]; len(files) > 0 {
		for _, fh := range files {
			f, err := fh.Open()
			if err != nil {
				continue
			}
			defer f.Close()
			imageFiles = append(imageFiles, domain.UploadedFile{Filename: fh.Filename, Reader: f})
		}
	}

	var videoFiles []domain.UploadedFile
	if files := form.File["videos"]; len(files) > 0 {
		for _, fh := range files {
			f, err := fh.Open()
			if err != nil {
				continue
			}
			defer f.Close()
			videoFiles = append(videoFiles, domain.UploadedFile{Filename: fh.Filename, Reader: f})
		}
	}

	news := &domain.News{
		DateTime:   payload.Datetime,
		TitleRu:    payload.TitleRu,
		TitleEn:    payload.TitleEn,
		SubTitleRu: payload.SubTitleRu,
		SubTitleEn: payload.SubTitleEn,
		Dir:        payload.Dir,
		ImgBack:    payload.ImgBack,
		ImgBackfull: payload.ImgBackfull,
		TextRu:     payload.TextRu,
		TextEn:     payload.TextEn,
		Images:     payload.Images,
		Videos:     payload.Videos,
	}

	if err := h.newsService.CreateNews(c.Request.Context(), news, imgBackFile, imgBackfullFile, imageFiles, videoFiles); err != nil {
		slog.Error("CreateNews: failed", "error", err)
		apiErr := apperror.FromError(err)
		c.JSON(apiErr.Status, apiErr)
		return
	}

	slog.Info("News created", "news_id", news.ID, "title_ru", news.TitleRu)
	c.JSON(http.StatusCreated, newsToResponse(news))
}

func (h *NewsHandler) UpdateNews(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, apperror.APIError{
			Status:  http.StatusBadRequest,
			Code:    "INVALID_PARAM",
			Message: "Invalid id",
		})
		return
	}

	dataJSON := c.PostForm("data")
	if dataJSON == "" {
		c.JSON(http.StatusBadRequest, apperror.APIError{
			Status:  http.StatusBadRequest,
			Code:    "INVALID_REQUEST",
			Message: "Missing data field",
		})
		return
	}

	var payload dto.NewsDataPayload
	if err := json.Unmarshal([]byte(dataJSON), &payload); err != nil {
		slog.Warn("UpdateNews: invalid data JSON", "error", err)
		c.JSON(http.StatusBadRequest, apperror.APIError{
			Status:  http.StatusBadRequest,
			Code:    "INVALID_REQUEST",
			Message: "Invalid data JSON",
		})
		return
	}

	payload.ID = id

	form, err := c.MultipartForm()
	if err != nil {
		form = nil
	}

	var imgBackFile, imgBackfullFile *domain.UploadedFile
	if form != nil {
		if files := form.File["img_back"]; len(files) > 0 {
			f, err := files[0].Open()
			if err == nil {
				defer f.Close()
				imgBackFile = &domain.UploadedFile{Filename: files[0].Filename, Reader: f}
			}
		}
		if files := form.File["img_backfull"]; len(files) > 0 {
			f, err := files[0].Open()
			if err == nil {
				defer f.Close()
				imgBackfullFile = &domain.UploadedFile{Filename: files[0].Filename, Reader: f}
			}
		}
	}

	var imageFiles []domain.UploadedFile
	if form != nil {
		if files := form.File["images"]; len(files) > 0 {
			for _, fh := range files {
				f, err := fh.Open()
				if err != nil {
					continue
				}
				defer f.Close()
				imageFiles = append(imageFiles, domain.UploadedFile{Filename: fh.Filename, Reader: f})
			}
		}
	}

	var videoFiles []domain.UploadedFile
	if form != nil {
		if files := form.File["videos"]; len(files) > 0 {
			for _, fh := range files {
				f, err := fh.Open()
				if err != nil {
					continue
				}
				defer f.Close()
				videoFiles = append(videoFiles, domain.UploadedFile{Filename: fh.Filename, Reader: f})
			}
		}
	}

	news := &domain.News{
		ID:          id,
		DateTime:    payload.Datetime,
		TitleRu:     payload.TitleRu,
		TitleEn:     payload.TitleEn,
		SubTitleRu:  payload.SubTitleRu,
		SubTitleEn:  payload.SubTitleEn,
		Dir:         payload.Dir,
		ImgBack:     payload.ImgBack,
		ImgBackfull: payload.ImgBackfull,
		TextRu:      payload.TextRu,
		TextEn:      payload.TextEn,
		Images:      payload.Images,
		Videos:      payload.Videos,
	}

	if err := h.newsService.UpdateNews(c.Request.Context(), news, imgBackFile, imgBackfullFile, imageFiles, videoFiles); err != nil {
		slog.Error("UpdateNews: failed", "news_id", id, "error", err)
		apiErr := apperror.FromError(err)
		c.JSON(apiErr.Status, apiErr)
		return
	}

	slog.Info("News updated", "news_id", id)
	c.JSON(http.StatusOK, newsToResponse(news))
}

func (h *NewsHandler) DeleteNews(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, apperror.APIError{
			Status:  http.StatusBadRequest,
			Code:    "INVALID_PARAM",
			Message: "Invalid id",
		})
		return
	}

	if err := h.newsService.DeleteNews(c.Request.Context(), id); err != nil {
		slog.Error("DeleteNews: failed", "news_id", id, "error", err)
		apiErr := apperror.FromError(err)
		c.JSON(apiErr.Status, apiErr)
		return
	}

	slog.Info("News deleted", "news_id", id)
	c.JSON(http.StatusOK, gin.H{"message": "News deleted successfully"})
}

func (h *NewsHandler) BulkDeleteNews(c *gin.Context) {
	var req struct {
		IDs []string `json:"ids"`
	}
	if !BindJSON(c, &req) {
		return
	}

	if err := h.newsService.BulkDeleteNews(c.Request.Context(), req.IDs); err != nil {
		slog.Error("BulkDeleteNews: failed", "error", err)
		apiErr := apperror.FromError(err)
		c.JSON(apiErr.Status, apiErr)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "News deleted successfully"})
}

func newsToResponse(n *domain.News) dto.NewsResponse {
	return dto.NewsResponse{
		ID:          n.ID,
		DateTime:    n.DateTime,
		TitleRu:     n.TitleRu,
		TitleEn:     n.TitleEn,
		SubTitleRu:  n.SubTitleRu,
		SubTitleEn:  n.SubTitleEn,
		Dir:         n.Dir,
		ImgBack:     n.ImgBack,
		ImgBackfull: n.ImgBackfull,
		TextRu:      n.TextRu,
		TextEn:      n.TextEn,
		Images:      n.Images,
		Videos:      n.Videos,
	}
}
