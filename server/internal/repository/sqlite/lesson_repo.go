package sqlite

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/raiyin/artserver/internal/domain"
)

// LessonRepository implements port.LessonRepository.
type LessonRepository struct {
	db *sql.DB
}

// NewLessonRepository creates a new LessonRepository.
func NewLessonRepository(db *sql.DB) *LessonRepository {
	return &LessonRepository{db: db}
}

const lessonColumns = `id, product_id, title, description, content, video_url, resources, duration_minutes, sort_order, status, created_at, updated_at`

func (r *LessonRepository) scanLesson(scanner interface {
	Scan(dest ...interface{}) error
}) (*domain.Lesson, error) {
	l := &domain.Lesson{}
	var description, content, videoURL, resourcesStr sql.NullString
	var durationMinutes sql.NullInt64

	err := scanner.Scan(
		&l.ID, &l.ProductID, &l.Title, &description, &content,
		&videoURL, &resourcesStr, &durationMinutes,
		&l.SortOrder, &l.Status, &l.CreatedAt, &l.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	if description.Valid {
		l.Description = description.String
	}
	if content.Valid {
		l.Content = content.String
	}
	if videoURL.Valid {
		l.VideoURL = videoURL.String
	}
	if durationMinutes.Valid {
		l.DurationMinutes = int(durationMinutes.Int64)
	}
	if resourcesStr.Valid {
		json.Unmarshal([]byte(resourcesStr.String), &l.Resources)
	}

	return l, nil
}

// Create inserts a new lesson.
func (r *LessonRepository) Create(ctx context.Context, lesson *domain.Lesson) error {
	query := `INSERT INTO lessons (product_id, title, description, content, video_url, resources, duration_minutes, sort_order, status, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`

	now := time.Now()
	if lesson.CreatedAt.IsZero() {
		lesson.CreatedAt = now
	}
	if lesson.UpdatedAt.IsZero() {
		lesson.UpdatedAt = now
	}
	if lesson.Status == "" {
		lesson.Status = "draft"
	}

	resourcesJSON, _ := json.Marshal(lesson.Resources)

	result, err := r.db.ExecContext(ctx, query,
		lesson.ProductID, lesson.Title, nullString(lesson.Description), nullString(lesson.Content),
		nullString(lesson.VideoURL), string(resourcesJSON), nullInt(int64(lesson.DurationMinutes)),
		lesson.SortOrder, lesson.Status, lesson.CreatedAt, lesson.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("create lesson: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("get last insert id: %w", err)
	}
	lesson.ID = id
	return nil
}

// GetByID retrieves a lesson by ID.
func (r *LessonRepository) GetByID(ctx context.Context, id int64) (*domain.Lesson, error) {
	query := fmt.Sprintf("SELECT %s FROM lessons WHERE id = ?", lessonColumns)

	lesson, err := r.scanLesson(r.db.QueryRowContext(ctx, query, id))
	if err == sql.ErrNoRows {
		return nil, domain.ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("get lesson by id: %w", err)
	}
	return lesson, nil
}

// ListByProduct retrieves lessons for a product.
func (r *LessonRepository) ListByProduct(ctx context.Context, productID int64, status string) ([]domain.Lesson, error) {
	var conditions []string
	var args []interface{}

	conditions = append(conditions, "product_id = ?")
	args = append(args, productID)

	if status != "" {
		conditions = append(conditions, "status = ?")
		args = append(args, status)
	}

	whereClause := "WHERE " + strings.Join(conditions, " AND ")
	query := fmt.Sprintf("SELECT %s FROM lessons %s ORDER BY sort_order", lessonColumns, whereClause)

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var lessons []domain.Lesson
	for rows.Next() {
		l, err := r.scanLesson(rows)
		if err != nil {
			return nil, err
		}
		lessons = append(lessons, *l)
	}
	return lessons, nil
}

// List retrieves lessons with optional filtering.
func (r *LessonRepository) List(ctx context.Context, filter domain.LessonFilter) ([]domain.Lesson, int, error) {
	var conditions []string
	var args []interface{}

	if filter.ProductID > 0 {
		conditions = append(conditions, "product_id = ?")
		args = append(args, filter.ProductID)
	}
	if filter.Status != "" {
		conditions = append(conditions, "status = ?")
		args = append(args, filter.Status)
	}
	if filter.Query != "" {
		conditions = append(conditions, "title LIKE ?")
		args = append(args, "%"+filter.Query+"%")
	}

	whereClause := ""
	if len(conditions) > 0 {
		whereClause = "WHERE " + strings.Join(conditions, " AND ")
	}

	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM lessons %s", whereClause)
	var total int
	if err := r.db.QueryRowContext(ctx, countQuery, args...).Scan(&total); err != nil {
		return nil, 0, err
	}

	listQuery := fmt.Sprintf("SELECT %s FROM lessons %s ORDER BY sort_order", lessonColumns, whereClause)
	rows, err := r.db.QueryContext(ctx, listQuery, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var lessons []domain.Lesson
	for rows.Next() {
		l, err := r.scanLesson(rows)
		if err != nil {
			return nil, 0, err
		}
		lessons = append(lessons, *l)
	}

	return lessons, total, nil
}

// Update updates a lesson.
func (r *LessonRepository) Update(ctx context.Context, lesson *domain.Lesson) error {
	query := `UPDATE lessons SET product_id = ?, title = ?, description = ?, content = ?,
		video_url = ?, resources = ?, duration_minutes = ?, sort_order = ?, status = ?, updated_at = ? WHERE id = ?`

	lesson.UpdatedAt = time.Now()
	resourcesJSON, _ := json.Marshal(lesson.Resources)

	_, err := r.db.ExecContext(ctx, query,
		lesson.ProductID, lesson.Title, nullString(lesson.Description), nullString(lesson.Content),
		nullString(lesson.VideoURL), string(resourcesJSON), nullInt(int64(lesson.DurationMinutes)),
		lesson.SortOrder, lesson.Status, lesson.UpdatedAt, lesson.ID,
	)
	if err != nil {
		return fmt.Errorf("update lesson: %w", err)
	}
	return nil
}

// Delete deletes a lesson by ID.
func (r *LessonRepository) Delete(ctx context.Context, id int64) error {
	result, err := r.db.ExecContext(ctx, "DELETE FROM lessons WHERE id = ?", id)
	if err != nil {
		return fmt.Errorf("delete lesson: %w", err)
	}
	affected, _ := result.RowsAffected()
	if affected == 0 {
		return domain.ErrNotFound
	}
	return nil
}

// LearningProgressRepository implements port.LearningProgressRepository.
type LearningProgressRepository struct {
	db *sql.DB
}

// NewLearningProgressRepository creates a new LearningProgressRepository.
func NewLearningProgressRepository(db *sql.DB) *LearningProgressRepository {
	return &LearningProgressRepository{db: db}
}

func (r *LearningProgressRepository) Upsert(ctx context.Context, progress *domain.LearningProgress) error {
	query := `INSERT INTO learning_progress (user_id, lesson_id, product_id, completed, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?)
		ON CONFLICT(user_id, lesson_id) DO UPDATE SET completed = ?, updated_at = ?`

	now := time.Now()
	if progress.CreatedAt.IsZero() {
		progress.CreatedAt = now
	}
	progress.UpdatedAt = now

	_, err := r.db.ExecContext(ctx, query,
		progress.UserID, progress.LessonID, progress.ProductID, progress.Completed,
		progress.CreatedAt, progress.UpdatedAt,
		progress.Completed, progress.UpdatedAt,
	)
	return err
}

func (r *LearningProgressRepository) GetByUserAndLesson(ctx context.Context, userID, lessonID int64) (*domain.LearningProgress, error) {
	p := &domain.LearningProgress{}
	err := r.db.QueryRowContext(ctx,
		"SELECT id, user_id, lesson_id, product_id, completed, created_at, updated_at FROM learning_progress WHERE user_id = ? AND lesson_id = ?",
		userID, lessonID,
	).Scan(&p.ID, &p.UserID, &p.LessonID, &p.ProductID, &p.Completed, &p.CreatedAt, &p.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, domain.ErrNotFound
	}
	return p, err
}

func (r *LearningProgressRepository) ListByUserAndProduct(ctx context.Context, userID, productID int64) ([]domain.LearningProgress, error) {
	rows, err := r.db.QueryContext(ctx,
		"SELECT id, user_id, lesson_id, product_id, completed, created_at, updated_at FROM learning_progress WHERE user_id = ? AND product_id = ?",
		userID, productID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var progress []domain.LearningProgress
	for rows.Next() {
		var p domain.LearningProgress
		if err := rows.Scan(&p.ID, &p.UserID, &p.LessonID, &p.ProductID, &p.Completed, &p.CreatedAt, &p.UpdatedAt); err != nil {
			return nil, err
		}
		progress = append(progress, p)
	}
	return progress, nil
}

func (r *LearningProgressRepository) ListCoursesByUser(ctx context.Context, userID int64) ([]int64, error) {
	rows, err := r.db.QueryContext(ctx,
		"SELECT DISTINCT product_id FROM learning_progress WHERE user_id = ?", userID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var ids []int64
	for rows.Next() {
		var id int64
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		ids = append(ids, id)
	}
	return ids, nil
}
