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

const productColumns = `p.id, p.title, p.slug, p.description, p.full_description, p.price, p.old_price,
	p.image_path, p.category_id, p.status, p.is_digital, p.is_master_class, p.sort_order, p.tags, p.created_at, p.updated_at`

func (r *ProductRepository) scanProduct(scanner interface {
	Scan(dest ...interface{}) error
}) (*domain.Product, error) {
	p := &domain.Product{}
	var description, fullDescription, tagsStr sql.NullString
	var oldPrice sql.NullFloat64
	var categoryID sql.NullInt64

	err := scanner.Scan(
		&p.ID, &p.Title, &p.Slug, &description, &fullDescription,
		&p.Price, &oldPrice, &p.ImagePath, &categoryID, &p.Status,
		&p.IsDigital, &p.IsMasterClass, &p.SortOrder, &tagsStr,
		&p.CreatedAt, &p.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	if description.Valid {
		p.Description = description.String
	}
	if fullDescription.Valid {
		p.FullDescription = fullDescription.String
	}
	if oldPrice.Valid {
		p.OldPrice = oldPrice.Float64
	}
	if categoryID.Valid {
		p.CategoryID = categoryID.Int64
	}
	if tagsStr.Valid {
		json.Unmarshal([]byte(tagsStr.String), &p.Tags)
	}

	return p, nil
}

// Create inserts a new product.
func (r *ProductRepository) Create(ctx context.Context, product *domain.Product) error {
	query := `INSERT INTO products (title, slug, description, full_description, price, old_price,
		image_path, category_id, status, is_digital, is_master_class, sort_order, tags, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`

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

	tagsJSON, _ := json.Marshal(product.Tags)

	result, err := r.db.ExecContext(ctx, query,
		product.Title, product.Slug, nullString(product.Description), nullString(product.FullDescription),
		product.Price, nullFloat(product.OldPrice), product.ImagePath,
		nullInt(product.CategoryID), product.Status, product.IsDigital, product.IsMasterClass,
		product.SortOrder, string(tagsJSON), product.CreatedAt, product.UpdatedAt,
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
	if filter.CategoryID > 0 {
		conditions = append(conditions, "p.category_id = ?")
		args = append(args, filter.CategoryID)
	}
	if filter.IsDigital != nil {
		conditions = append(conditions, "p.is_digital = ?")
		args = append(args, *filter.IsDigital)
	}
	if filter.IsMasterClass != nil {
		conditions = append(conditions, "p.is_master_class = ?")
		args = append(args, *filter.IsMasterClass)
	}
	if filter.Query != "" {
		conditions = append(conditions, "(p.title LIKE ? OR p.description LIKE ?)")
		q := "%" + filter.Query + "%"
		args = append(args, q, q)
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

	sortBy := "p.sort_order"
	sortOrder := "ASC"
	if filter.SortBy != "" {
		sortBy = "p." + filter.SortBy
	}
	if filter.SortOrder == "desc" {
		sortOrder = "DESC"
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
	query := `UPDATE products SET title = ?, slug = ?, description = ?, full_description = ?,
		price = ?, old_price = ?, image_path = ?, category_id = ?, status = ?,
		is_digital = ?, is_master_class = ?, sort_order = ?, tags = ?, updated_at = ? WHERE id = ?`

	product.UpdatedAt = time.Now()
	tagsJSON, _ := json.Marshal(product.Tags)

	_, err := r.db.ExecContext(ctx, query,
		product.Title, product.Slug, nullString(product.Description), nullString(product.FullDescription),
		product.Price, nullFloat(product.OldPrice), product.ImagePath,
		nullInt(product.CategoryID), product.Status,
		product.IsDigital, product.IsMasterClass, product.SortOrder, string(tagsJSON),
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
	query := `INSERT INTO product_categories (name, slug, sort_order, created_at) VALUES (?, ?, ?, ?)`
	now := time.Now()
	if category.CreatedAt.IsZero() {
		category.CreatedAt = now
	}

	result, err := r.db.ExecContext(ctx, query, category.Name, category.Slug, category.SortOrder, category.CreatedAt)
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
	err := r.db.QueryRowContext(ctx, "SELECT id, name, slug, sort_order, created_at FROM product_categories WHERE id = ?", id).
		Scan(&c.ID, &c.Name, &c.Slug, &c.SortOrder, &c.CreatedAt)
	if err == sql.ErrNoRows {
		return nil, domain.ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	return c, nil
}

func (r *ProductCategoryRepository) List(ctx context.Context) ([]domain.ProductCategory, error) {
	rows, err := r.db.QueryContext(ctx, "SELECT id, name, slug, sort_order, created_at FROM product_categories ORDER BY sort_order")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var categories []domain.ProductCategory
	for rows.Next() {
		var c domain.ProductCategory
		if err := rows.Scan(&c.ID, &c.Name, &c.Slug, &c.SortOrder, &c.CreatedAt); err != nil {
			return nil, err
		}
		categories = append(categories, c)
	}
	return categories, nil
}

func (r *ProductCategoryRepository) Update(ctx context.Context, category *domain.ProductCategory) error {
	_, err := r.db.ExecContext(ctx, "UPDATE product_categories SET name = ?, slug = ?, sort_order = ? WHERE id = ?",
		category.Name, category.Slug, category.SortOrder, category.ID)
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
