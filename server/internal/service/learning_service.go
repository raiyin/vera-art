package service

import (
	"context"
	"log/slog"

	"github.com/raiyin/artserver/internal/domain"
	"github.com/raiyin/artserver/internal/port"
)

// LearningService implements port.LearningService.
type LearningService struct {
	lessonRepo   port.LessonRepository
	progressRepo port.LearningProgressRepository
	purchaseRepo port.PurchaseRepository
}

// NewLearningService creates a new LearningService.
func NewLearningService(
	lessonRepo port.LessonRepository,
	progressRepo port.LearningProgressRepository,
	purchaseRepo port.PurchaseRepository,
) *LearningService {
	return &LearningService{
		lessonRepo:   lessonRepo,
		progressRepo: progressRepo,
		purchaseRepo: purchaseRepo,
	}
}

// GetLessonsByProduct retrieves lessons for a product, checking access.
func (s *LearningService) GetLessonsByProduct(ctx context.Context, productID int64, userID int64) ([]domain.Lesson, error) {
	// Check if user has purchased the product or is admin
	purchase, err := s.purchaseRepo.GetByUserAndProduct(ctx, userID, productID)
	if err != nil && err != domain.ErrNotFound {
		slog.Error("LearningService.GetLessonsByProduct: failed to check purchase",
			"product_id", productID,
			"user_id", userID,
			"error", err,
		)
		return nil, err
	}

	onlyPublic := false
	if purchase == nil {
		onlyPublic = true
	}

	lessons, err := s.lessonRepo.ListByProduct(ctx, productID, onlyPublic)
	if err != nil {
		slog.Error("LearningService.GetLessonsByProduct: failed to list lessons",
			"product_id", productID,
			"user_id", userID,
			"only_public", onlyPublic,
			"error", err,
		)
		return nil, err
	}
	slog.Debug("LearningService.GetLessonsByProduct: lessons listed",
		"product_id", productID,
		"user_id", userID,
		"count", len(lessons),
	)
	return lessons, nil
}

// GetLessonByID retrieves a lesson by ID, checking access.
func (s *LearningService) GetLessonByID(ctx context.Context, id int64, userID int64) (*domain.Lesson, error) {
	lesson, err := s.lessonRepo.GetByID(ctx, id)
	if err != nil {
		slog.Error("LearningService.GetLessonByID: failed to get lesson",
			"lesson_id", id,
			"user_id", userID,
			"error", err,
		)
		return nil, err
	}

	// Check access
	purchase, _ := s.purchaseRepo.GetByUserAndProduct(ctx, userID, lesson.ProductID)
	if purchase == nil && !lesson.IsPreview {
		slog.Warn("LearningService.GetLessonByID: forbidden access",
			"lesson_id", id,
			"user_id", userID,
			"product_id", lesson.ProductID,
			"is_preview", lesson.IsPreview,
		)
		return nil, domain.ErrForbidden
	}

	slog.Debug("LearningService.GetLessonByID: lesson retrieved",
		"lesson_id", id,
		"user_id", userID,
		"product_id", lesson.ProductID,
	)
	return lesson, nil
}

// AdminGetLessons retrieves all lessons for admin.
func (s *LearningService) AdminGetLessons(ctx context.Context, filter domain.LessonFilter) ([]domain.Lesson, int, error) {
	lessons, total, err := s.lessonRepo.List(ctx, filter)
	if err != nil {
		slog.Error("LearningService.AdminGetLessons: failed to list lessons",
			"filter", filter,
			"error", err,
		)
		return nil, 0, err
	}
	slog.Debug("LearningService.AdminGetLessons: lessons listed",
		"count", len(lessons),
		"total", total,
	)
	return lessons, total, nil
}

// CreateLesson creates a new lesson.
func (s *LearningService) CreateLesson(ctx context.Context, lesson *domain.Lesson) error {
	if lesson.TitleRu == "" && lesson.TitleEn == "" {
		slog.Warn("LearningService.CreateLesson: empty title",
			"lesson", lesson,
		)
		return domain.ErrInvalidInput
	}
	if err := s.lessonRepo.Create(ctx, lesson); err != nil {
		slog.Error("LearningService.CreateLesson: failed to create lesson",
			"lesson_title_ru", lesson.TitleRu,
			"product_id", lesson.ProductID,
			"error", err,
		)
		return err
	}
	slog.Info("LearningService.CreateLesson: lesson created",
		"lesson_id", lesson.ID,
		"lesson_title_ru", lesson.TitleRu,
		"product_id", lesson.ProductID,
	)
	return nil
}

// UpdateLesson updates a lesson.
func (s *LearningService) UpdateLesson(ctx context.Context, lesson *domain.Lesson) error {
	if err := s.lessonRepo.Update(ctx, lesson); err != nil {
		slog.Error("LearningService.UpdateLesson: failed to update lesson",
			"lesson_id", lesson.ID,
			"lesson_title_ru", lesson.TitleRu,
			"error", err,
		)
		return err
	}
	slog.Info("LearningService.UpdateLesson: lesson updated",
		"lesson_id", lesson.ID,
		"lesson_title_ru", lesson.TitleRu,
	)
	return nil
}

// DeleteLesson deletes a lesson.
func (s *LearningService) DeleteLesson(ctx context.Context, id int64) error {
	if err := s.lessonRepo.Delete(ctx, id); err != nil {
		slog.Error("LearningService.DeleteLesson: failed to delete lesson",
			"lesson_id", id,
			"error", err,
		)
		return err
	}
	slog.Info("LearningService.DeleteLesson: lesson deleted",
		"lesson_id", id,
	)
	return nil
}

// GetMyCourses retrieves course IDs for a user.
func (s *LearningService) GetMyCourses(ctx context.Context, userID int64) ([]int64, error) {
	courses, err := s.progressRepo.ListCoursesByUser(ctx, userID)
	if err != nil {
		slog.Error("LearningService.GetMyCourses: failed to list courses",
			"user_id", userID,
			"error", err,
		)
		return nil, err
	}
	slog.Debug("LearningService.GetMyCourses: courses listed",
		"user_id", userID,
		"count", len(courses),
	)
	return courses, nil
}

// GetLearningProgress retrieves learning progress for a user and product.
func (s *LearningService) GetLearningProgress(ctx context.Context, userID, productID int64) ([]domain.LearningProgress, error) {
	progress, err := s.progressRepo.ListByUserAndProduct(ctx, userID, productID)
	if err != nil {
		slog.Error("LearningService.GetLearningProgress: failed to get progress",
			"user_id", userID,
			"product_id", productID,
			"error", err,
		)
		return nil, err
	}
	slog.Debug("LearningService.GetLearningProgress: progress retrieved",
		"user_id", userID,
		"product_id", productID,
		"count", len(progress),
	)
	return progress, nil
}

// UpdateLessonProgress updates the progress for a lesson.
func (s *LearningService) UpdateLessonProgress(ctx context.Context, userID, lessonID int64, completed bool) error {
	lesson, err := s.lessonRepo.GetByID(ctx, lessonID)
	if err != nil {
		slog.Error("LearningService.UpdateLessonProgress: failed to get lesson",
			"lesson_id", lessonID,
			"user_id", userID,
			"error", err,
		)
		return err
	}

	progress := &domain.LearningProgress{
		UserID:    userID,
		LessonID:  lessonID,
		Completed: completed,
	}

	if err := s.progressRepo.Upsert(ctx, progress); err != nil {
		slog.Error("LearningService.UpdateLessonProgress: failed to upsert progress",
			"user_id", userID,
			"lesson_id", lessonID,
			"product_id", lesson.ProductID,
			"completed", completed,
			"error", err,
		)
		return err
	}
	slog.Info("LearningService.UpdateLessonProgress: progress updated",
		"user_id", userID,
		"lesson_id", lessonID,
		"product_id", lesson.ProductID,
		"completed", completed,
	)
	return nil
}

func (s *LearningService) BulkDeleteLessons(ctx context.Context, ids []int64) error {
	for _, id := range ids {
		if err := s.DeleteLesson(ctx, id); err != nil {
			slog.Error("LearningService.BulkDeleteLessons: failed to delete lesson",
				"lesson_id", id,
				"error", err,
			)
			return err
		}
	}
	slog.Info("LearningService.BulkDeleteLessons: lessons deleted",
		"count", len(ids),
	)
	return nil
}

// Ensure interface compliance.
var _ port.LearningService = (*LearningService)(nil)
