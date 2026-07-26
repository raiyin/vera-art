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

// ProductRepository implements port.ProductRepository.
type ProductRepository struct {
	db *sql.DB
}

// NewProductRepository creates a new ProductRepository.
func NewProductRepository(db *sql.DB) *ProductRepository {
	return &ProductRepository{db: db}
}

const productColumns = `p.id, p.type, p.title_ru, p.title_en, p.description_ru, p.description_en,
	p.short_description_ru, p.short_description_en, p.price, p.duration_days,
	p.thumbnail_url, p.video_url, p.status, p.difficulty, p.total_lessons,
	p.total_duration_minutes, p.category_id, p.instructor_id, p.tags,
	p.prerequisites_ru, p.prerequisites_en, p.learning_outcomes_ru, p.learning_outcomes_en,
	p.certificate_available, p.max_students, p.start_date, p.language,
	p.is_featured, p.view_count, p.created_at, p.updated_at`

func (r *ProductRepository) scanProduct(scanner interface {
	Scan(dest ...interface{}) error
}) (*domain.Product, error) {
	p := &domain.Product{}
	var descriptionRu, descriptionEn, shortDescriptionRu, shortDescriptionEn sql.NullString
	var thumbnailURL, videoURL, difficulty sql.NullString
	var durationDays, categoryID, instructorID, maxStudents sql.NullInt64
	var prerequisitesRu, prerequisitesEn, learningOutcomesRu, learningOutcomesEn sql.NullString
	var startDate sql.NullString
	var tagsStr sql.NullString
	var certificateAvailable, isFeatured sql.NullBool
	var totalLessons, totalDurationMinutes, viewCount sql.NullInt64

	err := scanner.Scan(
		&p.ID, &p.Type, &p.TitleRu, &p.TitleEn,
		&descriptionRu, &descriptionEn, &shortDescriptionRu, &shortDescriptionEn,
		&p.Price, &durationDays, &thumbnailURL, &videoURL,
		&p.Status, &difficulty, &totalLessons, &totalDurationMinutes,
		&categoryID, &instructorID, &tagsStr,
		&prerequisitesRu, &prerequisitesEn, &learningOutcomesRu, &learningOutcomesEn,
		&certificateAvailable, &maxStudents, &startDate, &p.Language,
		&isFeatured, &viewCount, &p.CreatedAt, &p.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	if descriptionRu.Valid {
		p.DescriptionRu = descriptionRu.String
	}
	if descriptionEn.Valid {
		p.DescriptionEn = descriptionEn.String
	}
	if shortDescriptionRu.Valid {
		p.ShortDescriptionRu = shortDescriptionRu.String
	}
	if shortDescriptionEn.Valid {
		p.ShortDescriptionEn = shortDescriptionEn.String
	}
	if thumbnailURL.Valid {
		p.ThumbnailURL = thumbnailURL.String
	}
	if videoURL.Valid {
		p.VideoURL = videoURL.String
	}
	if difficulty.Valid {
		p.Difficulty = difficulty.String
	}
	if durationDays.Valid {
		d := int(durationDays.Int64)
		p.DurationDays = &d
	}
	if categoryID.Valid {
		id := categoryID.Int64
		p.CategoryID = &id
	}
	if instructorID.Valid {
		id := instructorID.Int64
		p.InstructorID = &id
	}
	if maxStudents.Valid {
		m := int(maxStudents.Int64)
		p.MaxStudents = &m
	}
	if startDate.Valid {
		p.StartDate = &startDate.String
	}
	if certificateAvailable.Valid {
		p.CertificateAvailable = certificateAvailable.Bool
	}
	if isFeatured.Valid {
		p.IsFeatured = isFeatured.Bool
	}
	if totalLessons.Valid {
		p.TotalLessons = int(totalLessons.Int64)
	}
	if totalDurationMinutes.Valid {
		p.TotalDurationMinutes = int(totalDurationMinutes.Int64)
	}
	if viewCount.Valid {
		p.ViewCount = int(viewCount.Int64)
	}
	if prerequisitesRu.Valid {
		p.PrerequisitesRu = prerequisitesRu.String
	}
	if prerequisitesEn.Valid {
		p.PrerequisitesEn = prerequisitesEn.String
	}
	if learningOutcomesRu.Valid {
		p.LearningOutcomesRu = learningOutcomesRu.String
	}
	if learningOutcomesEn.Valid {
		p.LearningOutcomesEn = learningOutcomesEn.String
	}
	if tagsStr.Valid {
		json.Unmarshal([]byte(tagsStr.String), &p.Tags)
	}

	return p, nil
}

// Create inserts a new product.
func (r *ProductRepository) Create(ctx context.Context, product *domain.Product) error {
	query := `INSERT INTO products (type, title_ru, title_en, description_ru, description_en,
		short_description_ru, short_description_en, price, duration_days,
		thumbnail_url, video_url, status, difficulty, total_lessons,
		total_duration_minutes, category_id, instructor_id, tags,
		prerequisites_ru, prerequisites_en, learning_outcomes_ru, learning_outcomes_en,
		certificate_available, max_students, start_date, language,
		is_featured, view_count, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`

	now := time.Now()
	if product.CreatedAt.IsZero() {
		product.CreatedAt = now
	}
	if product.UpdatedAt.IsZero() {
		product.UpdatedAt = now
	}
	if product.Status == "" {
		product.Status = "draft"
	}
	if product.Type == "" {
		product.Type = "course"
	}
	if product.Language == "" {
		product.Language = "ru"
	}

	tagsJSON, _ := json.Marshal(product.Tags)

	result, err := r.db.ExecContext(ctx, query,
		product.Type, product.TitleRu, product.TitleEn,
		nullString(product.DescriptionRu), nullString(product.DescriptionEn),
		nullString(product.ShortDescriptionRu), nullString(product.ShortDescriptionEn),
		product.Price, nullIntPtr(product.DurationDays),
		nullString(product.ThumbnailURL), nullString(product.VideoURL),
		product.Status, nullString(product.Difficulty),
		nullInt(int64(product.TotalLessons)), nullInt(int64(product.TotalDurationMinutes)),
		nullInt64Ptr(product.CategoryID), nullInt64Ptr(product.InstructorID),
		string(tagsJSON),
		nullString(product.PrerequisitesRu), nullString(product.PrerequisitesEn),
		nullString(product.LearningOutcomesRu), nullString(product.LearningOutcomesEn),
		product.CertificateAvailable, nullIntPtr(product.MaxStudents),
		nullStringPtr(product.StartDate), product.Language,
		product.IsFeatured, product.ViewCount,
		product.CreatedAt, product.UpdatedAt,
	)
	if err != nil {
		if strings.Contains(err.Error(), "UNIQUE") {
			return domain.ErrDuplicate
		}
		return fmt.Errorf("create product: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("get last insert id: %w", err)
	}
	product.ID = id
	return nil
}

// GetByID retrieves a product by ID.
func (r *ProductRepository) GetByID(ctx context.Context, id int64) (*domain.Product, error) {
	query := fmt.Sprintf(`SELECT %s FROM products p LEFT JOIN product_categories pc ON p.category_id = pc.id WHERE p.id = ?`,
		productColumns)

	product, err := r.scanProduct(r.db.QueryRowContext(ctx, query, id))
	if err == sql.ErrNoRows {
		return nil, domain.ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("get product by id: %w", err)
	}
	return product, nil
}

// GetBySlug retrieves a product by slug.
func (r *ProductRepository) GetBySlug(ctx context.Context, slug string) (*domain.Product, error) {
	query := fmt.Sprintf(`SELECT %s FROM products p LEFT JOIN product_categories pc ON p.category_id = pc.id WHERE p.slug = ?`,
		productColumns)

	product, err := r.scanProduct(r.db.QueryRowContext(ctx, query, slug))
	if err == sql.ErrNoRows {
		return nil, domain.ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("get product by slug: %w", err)
	}
	return product, nil
}

// List retrieves products with optional filtering.
func (r *ProductRepository) List(ctx context.Context, filter domain.ProductFilter) ([]domain.Product, int, error) {
	var conditions []string
	var args []interface{}

	if filter.Status != "" {
		conditions = append(conditions, "p.status = ?")
		args = append(args, filter.Status)
	}
	if filter.Type != "" {
		conditions = append(conditions, "p.type = ?")
		args = append(args, filter.Type)
	}
	if filter.CategoryID > 0 {
		conditions = append(conditions, "p.category_id = ?")
		args = append(args, filter.CategoryID)
	}
	if filter.Difficulty != "" {
		conditions = append(conditions, "p.difficulty = ?")
		args = append(args, filter.Difficulty)
	}
	if filter.IsFeatured != nil {
		conditions = append(conditions, "p.is_featured = ?")
		args = append(args, *filter.IsFeatured)
	}
	if filter.Query != "" {
		conditions = append(conditions, "(p.title_ru LIKE ? OR p.title_en LIKE ? OR p.description_ru LIKE ? OR p.description_en LIKE ?)")
		q := "%" + filter.Query + "%"
		args = append(args, q, q, q, q)
	}

	whereClause := ""
	if len(conditions) > 0 {
		whereClause = "WHERE " + strings.Join(conditions, " AND ")
	}

	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM products p %s", whereClause)
	var total int
	if err := r.db.QueryRowContext(ctx, countQuery, args...).Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("count products: %w", err)
	}

	sortBy := "p.created_at"
	sortOrder := "DESC"
	if filter.SortBy != "" {
		sortBy = "p." + filter.SortBy
	}
	if filter.SortOrder == "asc" {
		sortOrder = "ASC"
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

	listQuery := fmt.Sprintf(`SELECT %s FROM products p LEFT JOIN product_categories pc ON p.category_id = pc.id %s ORDER BY %s %s LIMIT ? OFFSET ?`,
		productColumns, whereClause, sortBy, sortOrder)
	args = append(args, limit, offset)

	rows, err := r.db.QueryContext(ctx, listQuery, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("list products: %w", err)
	}
	defer rows.Close()

	var products []domain.Product
	for rows.Next() {
		p, err := r.scanProduct(rows)
		if err != nil {
			return nil, 0, fmt.Errorf("scan product: %w", err)
		}
		products = append(products, *p)
	}

	return products, total, nil
}

// Update updates a product.
func (r *ProductRepository) Update(ctx context.Context, product *domain.Product) error {
	query := `UPDATE products SET type = ?, title_ru = ?, title_en = ?,
		description_ru = ?, description_en = ?,
		short_description_ru = ?, short_description_en = ?,
		price = ?, duration_days = ?, thumbnail_url = ?, video_url = ?,
		status = ?, difficulty = ?, total_lessons = ?, total_duration_minutes = ?,
		category_id = ?, instructor_id = ?, tags = ?,
		prerequisites_ru = ?, prerequisites_en = ?,
		learning_outcomes_ru = ?, learning_outcomes_en = ?,
		certificate_available = ?, max_students = ?, start_date = ?, language = ?,
		is_featured = ?, view_count = ?, updated_at = ? WHERE id = ?`

	product.UpdatedAt = time.Now()
	tagsJSON, _ := json.Marshal(product.Tags)

	_, err := r.db.ExecContext(ctx, query,
		product.Type, product.TitleRu, product.TitleEn,
		nullString(product.DescriptionRu), nullString(product.DescriptionEn),
		nullString(product.ShortDescriptionRu), nullString(product.ShortDescriptionEn),
		product.Price, nullIntPtr(product.DurationDays),
		nullString(product.ThumbnailURL), nullString(product.VideoURL),
		product.Status, nullString(product.Difficulty),
		nullInt(int64(product.TotalLessons)), nullInt(int64(product.TotalDurationMinutes)),
		nullInt64Ptr(product.CategoryID), nullInt64Ptr(product.InstructorID),
		string(tagsJSON),
		nullString(product.PrerequisitesRu), nullString(product.PrerequisitesEn),
		nullString(product.LearningOutcomesRu), nullString(product.LearningOutcomesEn),
		product.CertificateAvailable, nullIntPtr(product.MaxStudents),
		nullStringPtr(product.StartDate), product.Language,
		product.IsFeatured, product.ViewCount,
		product.UpdatedAt, product.ID,
	)
	if err != nil {
		return fmt.Errorf("update product: %w", err)
	}
	return nil
}

// Delete deletes a product by ID.
func (r *ProductRepository) Delete(ctx context.Context, id int64) error {
	result, err := r.db.ExecContext(ctx, "DELETE FROM products WHERE id = ?", id)
	if err != nil {
		return fmt.Errorf("delete product: %w", err)
	}
	affected, _ := result.RowsAffected()
	if affected == 0 {
		return domain.ErrNotFound
	}
	return nil
}

// ProductCategoryRepository implements port.ProductCategoryRepository.
type ProductCategoryRepository struct {
	db *sql.DB
}

// NewProductCategoryRepository creates a new ProductCategoryRepository.
func NewProductCategoryRepository(db *sql.DB) *ProductCategoryRepository {
	return &ProductCategoryRepository{db: db}
}

func (r *ProductCategoryRepository) Create(ctx context.Context, category *domain.ProductCategory) error {
	query := `INSERT INTO product_categories (name_ru, name_en, slug, description_ru, description_en, sort_order, is_active, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`
	now := time.Now()
	if category.CreatedAt.IsZero() {
		category.CreatedAt = now
	}
	if category.UpdatedAt.IsZero() {
		category.UpdatedAt = now
	}

	result, err := r.db.ExecContext(ctx, query,
		category.NameRu, category.NameEn, category.Slug,
		nullString(category.DescriptionRu), nullString(category.DescriptionEn),
		category.SortOrder, category.IsActive, category.CreatedAt, category.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("create category: %w", err)
	}
	id, err := result.LastInsertId()
	if err != nil {
		return err
	}
	category.ID = id
	return nil
}

func (r *ProductCategoryRepository) GetByID(ctx context.Context, id int64) (*domain.ProductCategory, error) {
	c := &domain.ProductCategory{}
	var descriptionRu, descriptionEn sql.NullString
	err := r.db.QueryRowContext(ctx,
		"SELECT id, name_ru, name_en, slug, description_ru, description_en, sort_order, is_active, created_at, updated_at FROM product_categories WHERE id = ?", id,
	).Scan(&c.ID, &c.NameRu, &c.NameEn, &c.Slug, &descriptionRu, &descriptionEn, &c.SortOrder, &c.IsActive, &c.CreatedAt, &c.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, domain.ErrNotFound
	}
	if descriptionRu.Valid {
		c.DescriptionRu = descriptionRu.String
	}
	if descriptionEn.Valid {
		c.DescriptionEn = descriptionEn.String
	}
	return c, err
}

func (r *ProductCategoryRepository) List(ctx context.Context) ([]domain.ProductCategory, error) {
	rows, err := r.db.QueryContext(ctx,
		"SELECT id, name_ru, name_en, slug, description_ru, description_en, sort_order, is_active, created_at, updated_at FROM product_categories ORDER BY sort_order")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var categories []domain.ProductCategory
	for rows.Next() {
		var c domain.ProductCategory
		var descriptionRu, descriptionEn sql.NullString
		if err := rows.Scan(&c.ID, &c.NameRu, &c.NameEn, &c.Slug, &descriptionRu, &descriptionEn, &c.SortOrder, &c.IsActive, &c.CreatedAt, &c.UpdatedAt); err != nil {
			return nil, err
		}
		if descriptionRu.Valid {
			c.DescriptionRu = descriptionRu.String
		}
		if descriptionEn.Valid {
			c.DescriptionEn = descriptionEn.String
		}
		categories = append(categories, c)
	}
	return categories, nil
}

func (r *ProductCategoryRepository) Update(ctx context.Context, category *domain.ProductCategory) error {
	category.UpdatedAt = time.Now()
	_, err := r.db.ExecContext(ctx,
		"UPDATE product_categories SET name_ru = ?, name_en = ?, slug = ?, description_ru = ?, description_en = ?, sort_order = ?, is_active = ?, updated_at = ? WHERE id = ?",
		category.NameRu, category.NameEn, category.Slug,
		nullString(category.DescriptionRu), nullString(category.DescriptionEn),
		category.SortOrder, category.IsActive, category.UpdatedAt, category.ID,
	)
	return err
}

func (r *ProductCategoryRepository) Delete(ctx context.Context, id int64) error {
	result, err := r.db.ExecContext(ctx, "DELETE FROM product_categories WHERE id = ?", id)
	if err != nil {
		return err
	}
	affected, _ := result.RowsAffected()
	if affected == 0 {
		return domain.ErrNotFound
	}
	return nil
}

// Helper functions for nullable fields.
func nullIntPtr(p *int) interface{} {
	if p == nil {
		return nil
	}
	return *p
}

func nullInt64Ptr(p *int64) interface{} {
	if p == nil {
		return nil
	}
	return *p
}

func nullStringPtr(p *string) interface{} {
	if p == nil {
		return nil
	}
	return *p
}
