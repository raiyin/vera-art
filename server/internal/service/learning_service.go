package service

import (
	"context"

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
		return nil, err
	}

	status := ""
	if purchase == nil {
		// Only show published lessons for non-purchasers
		status = "published"
	}

	return s.lessonRepo.ListByProduct(ctx, productID, status)
}

// GetLessonByID retrieves a lesson by ID, checking access.
func (s *LearningService) GetLessonByID(ctx context.Context, id int64, userID int64) (*domain.Lesson, error) {
	lesson, err := s.lessonRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Check access
	purchase, _ := s.purchaseRepo.GetByUserAndProduct(ctx, userID, lesson.ProductID)
	if purchase == nil && lesson.Status != "published" {
		return nil, domain.ErrForbidden
	}

	return lesson, nil
}

// AdminGetLessons retrieves all lessons for admin.
func (s *LearningService) AdminGetLessons(ctx context.Context, filter domain.LessonFilter) ([]domain.Lesson, int, error) {
	return s.lessonRepo.List(ctx, filter)
}

// CreateLesson creates a new lesson.
func (s *LearningService) CreateLesson(ctx context.Context, lesson *domain.Lesson) error {
	if lesson.Title == "" {
		return domain.ErrInvalidInput
	}
	return s.lessonRepo.Create(ctx, lesson)
}

// UpdateLesson updates a lesson.
func (s *LearningService) UpdateLesson(ctx context.Context, lesson *domain.Lesson) error {
	return s.lessonRepo.Update(ctx, lesson)
}

// DeleteLesson deletes a lesson.
func (s *LearningService) DeleteLesson(ctx context.Context, id int64) error {
	return s.lessonRepo.Delete(ctx, id)
}

// GetMyCourses retrieves course IDs for a user.
func (s *LearningService) GetMyCourses(ctx context.Context, userID int64) ([]int64, error) {
	return s.progressRepo.ListCoursesByUser(ctx, userID)
}

// GetLearningProgress retrieves learning progress for a user and product.
func (s *LearningService) GetLearningProgress(ctx context.Context, userID, productID int64) ([]domain.LearningProgress, error) {
	return s.progressRepo.ListByUserAndProduct(ctx, userID, productID)
}

// UpdateLessonProgress updates the progress for a lesson.
func (s *LearningService) UpdateLessonProgress(ctx context.Context, userID, lessonID int64, completed bool) error {
	lesson, err := s.lessonRepo.GetByID(ctx, lessonID)
	if err != nil {
		return err
	}

	progress := &domain.LearningProgress{
		UserID:    userID,
		LessonID:  lessonID,
		ProductID: lesson.ProductID,
		Completed: completed,
	}

	return s.progressRepo.Upsert(ctx, progress)
}

// Ensure interface compliance.
var _ port.LearningService = (*LearningService)(nil)
