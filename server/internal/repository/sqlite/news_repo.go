package sqlite

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/raiyin/artserver/internal/domain"
)

// NewsRepository implements port.NewsRepository.
type NewsRepository struct {
	db *sql.DB
}

// NewNewsRepository creates a new NewsRepository.
func NewNewsRepository(db *sql.DB) *NewsRepository {
	return &NewsRepository{db: db}
}

const newsColumns = `id, title, description, content, image_path, video_path, video_paths, image_paths, status, created_at, updated_at`

func (r *NewsRepository) scanNews(scanner interface {
	Scan(dest ...interface{}) error
}) (*domain.News, error) {
	n := &domain.News{}
	var description, content, videoPath, videoPaths, imagePaths sql.NullString

	err := scanner.Scan(
		&n.ID, &n.Title, &description, &content,
		&n.ImagePath, &videoPath, &videoPaths, &imagePaths, &n.Status,
		&n.CreatedAt, &n.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	if description.Valid {
		n.Description = description.String
	}
	if content.Valid {
		n.Content = content.String
	}
	if videoPath.Valid {
		n.VideoPath = videoPath.String
	}
	if videoPaths.Valid && videoPaths.String != "" {
		n.VideoPaths = strings.Split(videoPaths.String, ";")
	}
	if imagePaths.Valid && imagePaths.String != "" {
		n.ImagePaths = strings.Split(imagePaths.String, ";")
	}

	return n, nil
}

// Create inserts a new news entry.
func (r *NewsRepository) Create(ctx context.Context, news *domain.News) error {
	query := `INSERT INTO news (title, description, content, image_path, video_path, video_paths, image_paths, status, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`

	now := time.Now()
	if news.CreatedAt.IsZero() {
		news.CreatedAt = now
	}
	if news.UpdatedAt.IsZero() {
		news.UpdatedAt = now
	}
	if news.Status == "" {
		news.Status = "draft"
	}

	videoPathsStr := strings.Join(news.VideoPaths, ";")
	imagePathsStr := strings.Join(news.ImagePaths, ";")

	result, err := r.db.ExecContext(ctx, query,
		news.Title, nullString(news.Description), nullString(news.Content),
		news.ImagePath, nullString(news.VideoPath), nullString(videoPathsStr), nullString(imagePathsStr), news.Status,
		news.CreatedAt, news.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("create news: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("get last insert id: %w", err)
	}
	news.ID = id
	return nil
}

// GetByID retrieves a news entry by ID.
func (r *NewsRepository) GetByID(ctx context.Context, id int64) (*domain.News, error) {
	query := fmt.Sprintf("SELECT %s FROM news WHERE id = ?", newsColumns)

	news, err := r.scanNews(r.db.QueryRowContext(ctx, query, id))
	if err == sql.ErrNoRows {
		return nil, domain.ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("get news by id: %w", err)
	}
	return news, nil
}

// List retrieves news entries with optional filtering.
func (r *NewsRepository) List(ctx context.Context, filter domain.NewsFilter) ([]domain.News, int, error) {
	var conditions []string
	var args []interface{}

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

	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM news %s", whereClause)
	var total int
	if err := r.db.QueryRowContext(ctx, countQuery, args...).Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("count news: %w", err)
	}

	page := filter.Page
	if page < 1 {
		page = 1
	}
	limit := filter.Limit
	if limit < 1 {
		limit = 50
	}
	offset := (page - 1) * limit

	listQuery := fmt.Sprintf(`SELECT %s FROM news %s ORDER BY created_at DESC LIMIT ? OFFSET ?`,
		newsColumns, whereClause)
	args = append(args, limit, offset)

	rows, err := r.db.QueryContext(ctx, listQuery, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("list news: %w", err)
	}
	defer rows.Close()

	var newsList []domain.News
	for rows.Next() {
		n, err := r.scanNews(rows)
		if err != nil {
			return nil, 0, fmt.Errorf("scan news: %w", err)
		}
		newsList = append(newsList, *n)
	}

	return newsList, total, nil
}

// Update updates a news entry.
func (r *NewsRepository) Update(ctx context.Context, news *domain.News) error {
	query := `UPDATE news SET title = ?, description = ?, content = ?, image_path = ?,
		video_path = ?, video_paths = ?, image_paths = ?, status = ?, updated_at = ? WHERE id = ?`

	news.UpdatedAt = time.Now()

	videoPathsStr := strings.Join(news.VideoPaths, ";")
	imagePathsStr := strings.Join(news.ImagePaths, ";")

	_, err := r.db.ExecContext(ctx, query,
		news.Title, nullString(news.Description), nullString(news.Content),
		news.ImagePath, nullString(news.VideoPath), nullString(videoPathsStr), nullString(imagePathsStr), news.Status,
		news.UpdatedAt, news.ID,
	)
	if err != nil {
		return fmt.Errorf("update news: %w", err)
	}
	return nil
}

// Delete deletes a news entry by ID.
func (r *NewsRepository) Delete(ctx context.Context, id int64) error {
	result, err := r.db.ExecContext(ctx, "DELETE FROM news WHERE id = ?", id)
	if err != nil {
		return fmt.Errorf("delete news: %w", err)
	}
	affected, _ := result.RowsAffected()
	if affected == 0 {
		return domain.ErrNotFound
	}
	return nil
}
