package handler

import (
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/raiyin/artserver/internal/domain"
	"github.com/raiyin/artserver/internal/dto"
	"github.com/raiyin/artserver/internal/port"
	"github.com/raiyin/artserver/pkg/apperror"
)

// LearningHandler handles learning (lessons, progress) HTTP requests.
type LearningHandler struct {
	learningService port.LearningService
}

// NewLearningHandler creates a new LearningHandler.
func NewLearningHandler(learningService port.LearningService) *LearningHandler {
	return &LearningHandler{learningService: learningService}
}

// GetLessonsByProduct returns lessons for a product.
func (h *LearningHandler) GetLessonsByProduct(c *gin.Context) {
	productID, ok := ParseInt64Param(c, "product_id")
	if !ok {
		return
	}

	userID := c.GetInt64("user_id")

	lessons, err := h.learningService.GetLessonsByProduct(c.Request.Context(), productID, userID)
	if err != nil {
		slog.Error("GetLessonsByProduct: failed to get lessons",
			"product_id", productID,
			"user_id", userID,
			"error", err,
		)
		apiErr := apperror.FromError(err)
		c.JSON(apiErr.Status, apiErr)
		return
	}

	responses := make([]dto.LessonResponse, len(lessons))
	for i, l := range lessons {
		responses[i] = dto.LessonResponse{
			ID:              l.ID,
			ProductID:       l.ProductID,
			Title:           l.Title,
			Description:     l.Description,
			Content:         l.Content,
			VideoURL:        l.VideoURL,
			Resources:       l.Resources,
			DurationMinutes: l.DurationMinutes,
			SortOrder:       l.SortOrder,
			Status:          l.Status,
			CreatedAt:       l.CreatedAt,
			UpdatedAt:       l.UpdatedAt,
		}
	}

	c.JSON(http.StatusOK, gin.H{"lessons": responses})
}

// GetLessonByID returns a lesson by ID.
func (h *LearningHandler) GetLessonByID(c *gin.Context) {
	id, ok := ParseInt64Param(c, "id")
	if !ok {
		return
	}

	userID := c.GetInt64("user_id")

	lesson, err := h.learningService.GetLessonByID(c.Request.Context(), id, userID)
	if err != nil {
		slog.Error("GetLessonByID: failed to get lesson",
			"lesson_id", id,
			"user_id", userID,
			"error", err,
		)
		apiErr := apperror.FromError(err)
		c.JSON(apiErr.Status, apiErr)
		return
	}

	c.JSON(http.StatusOK, dto.LessonResponse{
		ID:              lesson.ID,
		ProductID:       lesson.ProductID,
		Title:           lesson.Title,
		Description:     lesson.Description,
		Content:         lesson.Content,
		VideoURL:        lesson.VideoURL,
		Resources:       lesson.Resources,
		DurationMinutes: lesson.DurationMinutes,
		SortOrder:       lesson.SortOrder,
		Status:          lesson.Status,
		CreatedAt:       lesson.CreatedAt,
		UpdatedAt:       lesson.UpdatedAt,
	})
}

// AdminGetLessons returns all lessons (admin).
func (h *LearningHandler) AdminGetLessons(c *gin.Context) {
	filter := domain.LessonFilter{
		ProductID: ParseIntQueryAsInt64(c, "product_id"),
		Status:    c.Query("status"),
		Query:     c.Query("q"),
	}

	lessons, total, err := h.learningService.AdminGetLessons(c.Request.Context(), filter)
	if err != nil {
		slog.Error("AdminGetLessons: failed to list lessons",
			"error", err,
			"filter", filter,
		)
		apiErr := apperror.FromError(err)
		c.JSON(apiErr.Status, apiErr)
		return
	}

	responses := make([]dto.LessonResponse, len(lessons))
	for i, l := range lessons {
		responses[i] = dto.LessonResponse{
			ID:              l.ID,
			ProductID:       l.ProductID,
			Title:           l.Title,
			Description:     l.Description,
			Content:         l.Content,
			VideoURL:        l.VideoURL,
			Resources:       l.Resources,
			DurationMinutes: l.DurationMinutes,
			SortOrder:       l.SortOrder,
			Status:          l.Status,
			CreatedAt:       l.CreatedAt,
			UpdatedAt:       l.UpdatedAt,
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"lessons": responses,
		"total":   total,
	})
}

// CreateLesson creates a new lesson.
func (h *LearningHandler) CreateLesson(c *gin.Context) {
	var req dto.CreateLessonRequest
	if !BindJSON(c, &req) {
		return
	}

	lesson := &domain.Lesson{
		ProductID:       req.ProductID,
		Title:           req.Title,
		Description:     req.Description,
		Content:         req.Content,
		VideoURL:        req.VideoURL,
		Resources:       req.Resources,
		DurationMinutes: req.DurationMinutes,
		SortOrder:       req.SortOrder,
		Status:          req.Status,
	}

	if err := h.learningService.CreateLesson(c.Request.Context(), lesson); err != nil {
		slog.Error("CreateLesson: failed to create lesson",
			"title", req.Title,
			"product_id", req.ProductID,
			"error", err,
		)
		apiErr := apperror.FromError(err)
		c.JSON(apiErr.Status, apiErr)
		return
	}

	slog.Info("Lesson created successfully",
		"lesson_id", lesson.ID,
		"title", lesson.Title,
	)
	c.JSON(http.StatusCreated, dto.LessonResponse{
		ID:              lesson.ID,
		ProductID:       lesson.ProductID,
		Title:           lesson.Title,
		Description:     lesson.Description,
		Content:         lesson.Content,
		VideoURL:        lesson.VideoURL,
		Resources:       lesson.Resources,
		DurationMinutes: lesson.DurationMinutes,
		SortOrder:       lesson.SortOrder,
		Status:          lesson.Status,
		CreatedAt:       lesson.CreatedAt,
		UpdatedAt:       lesson.UpdatedAt,
	})
}

// UpdateLesson updates a lesson.
func (h *LearningHandler) UpdateLesson(c *gin.Context) {
	id, ok := ParseInt64Param(c, "id")
	if !ok {
		return
	}

	var req dto.UpdateLessonRequest
	if !BindJSON(c, &req) {
		return
	}

	lesson := &domain.Lesson{
		ID:              id,
		Title:           req.Title,
		Description:     req.Description,
		Content:         req.Content,
		VideoURL:        req.VideoURL,
		Resources:       req.Resources,
		DurationMinutes: req.DurationMinutes,
		SortOrder:       req.SortOrder,
		Status:          req.Status,
	}

	if err := h.learningService.UpdateLesson(c.Request.Context(), lesson); err != nil {
		slog.Error("UpdateLesson: failed to update lesson",
			"lesson_id", id,
			"error", err,
		)
		apiErr := apperror.FromError(err)
		c.JSON(apiErr.Status, apiErr)
		return
	}

	slog.Info("Lesson updated successfully",
		"lesson_id", id,
	)
	c.JSON(http.StatusOK, dto.LessonResponse{
		ID:              lesson.ID,
		ProductID:       lesson.ProductID,
		Title:           lesson.Title,
		Description:     lesson.Description,
		Content:         lesson.Content,
		VideoURL:        lesson.VideoURL,
		Resources:       lesson.Resources,
		DurationMinutes: lesson.DurationMinutes,
		SortOrder:       lesson.SortOrder,
		Status:          lesson.Status,
		CreatedAt:       lesson.CreatedAt,
		UpdatedAt:       lesson.UpdatedAt,
	})
}

// DeleteLesson deletes a lesson.
func (h *LearningHandler) DeleteLesson(c *gin.Context) {
	id, ok := ParseInt64Param(c, "id")
	if !ok {
		return
	}

	if err := h.learningService.DeleteLesson(c.Request.Context(), id); err != nil {
		slog.Error("DeleteLesson: failed to delete lesson",
			"lesson_id", id,
			"error", err,
		)
		apiErr := apperror.FromError(err)
		c.JSON(apiErr.Status, apiErr)
		return
	}

	slog.Info("Lesson deleted successfully",
		"lesson_id", id,
	)
	c.JSON(http.StatusOK, gin.H{"message": "Lesson deleted successfully"})
}

// GetMyCourses returns course IDs for the authenticated user.
func (h *LearningHandler) GetMyCourses(c *gin.Context) {
	userID := c.GetInt64("user_id")

	courseIDs, err := h.learningService.GetMyCourses(c.Request.Context(), userID)
	if err != nil {
		slog.Error("GetMyCourses: failed to get courses",
			"user_id", userID,
			"error", err,
		)
		apiErr := apperror.FromError(err)
		c.JSON(apiErr.Status, apiErr)
		return
	}

	c.JSON(http.StatusOK, gin.H{"course_ids": courseIDs})
}

// GetLearningProgress returns learning progress for a user and product.
func (h *LearningHandler) GetLearningProgress(c *gin.Context) {
	userID := c.GetInt64("user_id")
	productID, ok := ParseInt64Param(c, "product_id")
	if !ok {
		return
	}

	progress, err := h.learningService.GetLearningProgress(c.Request.Context(), userID, productID)
	if err != nil {
		slog.Error("GetLearningProgress: failed to get progress",
			"user_id", userID,
			"product_id", productID,
			"error", err,
		)
		apiErr := apperror.FromError(err)
		c.JSON(apiErr.Status, apiErr)
		return
	}

	responses := make([]dto.LearningProgressResponse, len(progress))
	for i, p := range progress {
		responses[i] = dto.LearningProgressResponse{
			ID:        p.ID,
			UserID:    p.UserID,
			LessonID:  p.LessonID,
			ProductID: p.ProductID,
			Completed: p.Completed,
			CreatedAt: p.CreatedAt,
			UpdatedAt: p.UpdatedAt,
		}
	}

	c.JSON(http.StatusOK, gin.H{"progress": responses})
}

// UpdateLessonProgress updates the progress for a lesson.
func (h *LearningHandler) UpdateLessonProgress(c *gin.Context) {
	userID := c.GetInt64("user_id")
	lessonID, ok := ParseInt64Param(c, "id")
	if !ok {
		return
	}

	var req dto.UpdateLessonProgressRequest
	if !BindJSON(c, &req) {
		return
	}

	if err := h.learningService.UpdateLessonProgress(c.Request.Context(), userID, lessonID, req.Completed); err != nil {
		slog.Error("UpdateLessonProgress: failed to update progress",
			"user_id", userID,
			"lesson_id", lessonID,
			"error", err,
		)
		apiErr := apperror.FromError(err)
		c.JSON(apiErr.Status, apiErr)
		return
	}

	slog.Info("Lesson progress updated",
		"user_id", userID,
		"lesson_id", lessonID,
		"completed", req.Completed,
	)
	c.JSON(http.StatusOK, gin.H{"message": "Progress updated successfully"})
}
