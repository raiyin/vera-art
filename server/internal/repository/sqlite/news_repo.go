package sqlite

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/raiyin/artserver/internal/domain"
)

type NewsRepository struct {
	db *sql.DB
}

func NewNewsRepository(db *sql.DB) *NewsRepository {
	return &NewsRepository{db: db}
}

const newsColumns = `id, datetime, title_ru, title_en, subTitle_ru, subTitle_en, dir, img_back, img_backfull, text_ru, text_en, images, videos`

func (r *NewsRepository) scanNews(scanner interface {
	Scan(dest ...interface{}) error
}) (*domain.News, error) {
	n := &domain.News{}
	var subTitleRu, subTitleEn, imagesStr, videosStr sql.NullString

	err := scanner.Scan(
		&n.ID, &n.DateTime, &n.TitleRu, &n.TitleEn,
		&subTitleRu, &subTitleEn,
		&n.Dir, &n.ImgBack, &n.ImgBackfull,
		&n.TextRu, &n.TextEn,
		&imagesStr, &videosStr,
	)
	if err != nil {
		return nil, err
	}

	if subTitleRu.Valid {
		n.SubTitleRu = subTitleRu.String
	}
	if subTitleEn.Valid {
		n.SubTitleEn = subTitleEn.String
	}
	if imagesStr.Valid && imagesStr.String != "" {
		n.Images = splitAndTrim(imagesStr.String, ";")
	}
	if videosStr.Valid && videosStr.String != "" {
		n.Videos = splitAndTrim(videosStr.String, ";")
	}

	if n.Images == nil {
		n.Images = []string{}
	}
	if n.Videos == nil {
		n.Videos = []string{}
	}

	return n, nil
}

func (r *NewsRepository) Create(ctx context.Context, news *domain.News) error {
	query := `INSERT INTO news_legacy (id, datetime, title_ru, title_en, subTitle_ru, subTitle_en, dir, img_back, img_backfull, text_ru, text_en, images, videos)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`

	_, err := r.db.ExecContext(ctx, query,
		news.ID, news.DateTime, news.TitleRu, news.TitleEn,
		nullString(news.SubTitleRu), nullString(news.SubTitleEn),
		news.Dir, news.ImgBack, news.ImgBackfull,
		news.TextRu, news.TextEn,
		joinOrNull(news.Images, ";"), joinOrNull(news.Videos, ";"),
	)
	if err != nil {
		return fmt.Errorf("create news: %w", err)
	}
	return nil
}

func (r *NewsRepository) GetByID(ctx context.Context, id string) (*domain.News, error) {
	query := fmt.Sprintf("SELECT %s FROM news_legacy WHERE id = ?", newsColumns)
	news, err := r.scanNews(r.db.QueryRowContext(ctx, query, id))
	if err == sql.ErrNoRows {
		return nil, domain.ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("get news by id: %w", err)
	}
	return news, nil
}

func (r *NewsRepository) List(ctx context.Context, filter domain.NewsFilter) ([]domain.News, int, error) {
	var conditions []string
	var args []interface{}

	if filter.Query != "" {
		conditions = append(conditions, "(title_ru LIKE ? OR title_en LIKE ? OR text_ru LIKE ? OR text_en LIKE ?)")
		q := "%" + filter.Query + "%"
		args = append(args, q, q, q, q)
	}

	whereClause := ""
	if len(conditions) > 0 {
		whereClause = "WHERE " + strings.Join(conditions, " AND ")
	}

	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM news_legacy %s", whereClause)
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

	listQuery := fmt.Sprintf(`SELECT %s FROM news_legacy %s ORDER BY datetime DESC LIMIT ? OFFSET ?`,
		newsColumns, whereClause)
	queryArgs := append(args, limit, offset)

	rows, err := r.db.QueryContext(ctx, listQuery, queryArgs...)
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

func (r *NewsRepository) Update(ctx context.Context, news *domain.News) error {
	query := `UPDATE news_legacy SET datetime=?, title_ru=?, title_en=?, subTitle_ru=?, subTitle_en=?,
		dir=?, img_back=?, img_backfull=?, text_ru=?, text_en=?, images=?, videos=? WHERE id=?`

	_, err := r.db.ExecContext(ctx, query,
		news.DateTime, news.TitleRu, news.TitleEn,
		nullString(news.SubTitleRu), nullString(news.SubTitleEn),
		news.Dir, news.ImgBack, news.ImgBackfull,
		news.TextRu, news.TextEn,
		joinOrNull(news.Images, ";"), joinOrNull(news.Videos, ";"),
		news.ID,
	)
	if err != nil {
		return fmt.Errorf("update news: %w", err)
	}
	return nil
}

func (r *NewsRepository) Delete(ctx context.Context, id string) error {
	result, err := r.db.ExecContext(ctx, "DELETE FROM news_legacy WHERE id = ?", id)
	if err != nil {
		return fmt.Errorf("delete news: %w", err)
	}
	affected, _ := result.RowsAffected()
	if affected == 0 {
		return domain.ErrNotFound
	}
	return nil
}

func splitAndTrim(s, sep string) []string {
	parts := strings.Split(s, sep)
	var result []string
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p != "" {
			result = append(result, p)
		}
	}
	return result
}

func joinOrNull(items []string, sep string) interface{} {
	if len(items) == 0 {
		return nil
	}
	return strings.Join(items, sep)
}
