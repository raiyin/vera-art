package sqlite

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/raiyin/artserver/internal/domain"
)

// SaleRepository implements port.SaleRepository.
type SaleRepository struct {
	db *sql.DB
}

// NewSaleRepository creates a new SaleRepository.
func NewSaleRepository(db *sql.DB) *SaleRepository {
	return &SaleRepository{db: db}
}

const saleColumns = `id, name_ru, name_en, description, image_path, sale_path, price, year, technique, width, height, status, sort_order, sold, created_at, updated_at`

func (r *SaleRepository) scanSale(scanner interface {
	Scan(dest ...interface{}) error
}) (*domain.Sale, error) {
	s := &domain.Sale{}
	var description, technique sql.NullString
	var year sql.NullFloat64
	var width, height sql.NullInt64

	err := scanner.Scan(
		&s.ID, &s.NameRu, &s.NameEn, &description, &s.ImagePath, &s.SalePath,
		&s.Price, &year, &technique,
		&width, &height,
		&s.Status, &s.SortOrder, &s.Sold,
		&s.CreatedAt, &s.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	if description.Valid {
		s.Description = description.String
	}
	if technique.Valid {
		s.Technique = technique.String
	}
	if width.Valid {
		s.Width = int(width.Int64)
	}
	if height.Valid {
		s.Height = int(height.Int64)
	}
	if year.Valid {
		s.Year = int(year.Float64)
	}

	return s, nil
}

// Create inserts a new sale.
func (r *SaleRepository) Create(ctx context.Context, sale *domain.Sale) error {
	query := `INSERT INTO sales (name_ru, name_en, description, image_path, sale_path, price, year, technique, width, height, status, sort_order, sold, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`

	now := time.Now()
	if sale.CreatedAt.IsZero() {
		sale.CreatedAt = now
	}
	if sale.UpdatedAt.IsZero() {
		sale.UpdatedAt = now
	}
	if sale.Status == "" {
		sale.Status = "draft"
	}

	result, err := r.db.ExecContext(ctx, query,
		sale.NameRu, sale.NameEn, nullString(sale.Description), sale.ImagePath, sale.SalePath,
		sale.Price, nullInt(int64(sale.Year)),
		nullString(sale.Technique),
		nullInt(int64(sale.Width)), nullInt(int64(sale.Height)),
		sale.Status, sale.SortOrder, sale.Sold,
		sale.CreatedAt, sale.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("create sale: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("get last insert id: %w", err)
	}
	sale.ID = id

	if len(sale.MaterialIDs) > 0 {
		if err := r.SetMaterials(ctx, sale.ID, sale.MaterialIDs); err != nil {
			return err
		}
	}
	if len(sale.BaseIDs) > 0 {
		if err := r.SetBases(ctx, sale.ID, sale.BaseIDs); err != nil {
			return err
		}
	}

	return nil
}

// GetByID retrieves a sale by ID.
func (r *SaleRepository) GetByID(ctx context.Context, id int64) (*domain.Sale, error) {
	query := fmt.Sprintf("SELECT %s FROM sales WHERE id = ?", saleColumns)

	sale, err := r.scanSale(r.db.QueryRowContext(ctx, query, id))
	if err == sql.ErrNoRows {
		return nil, domain.ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("get sale by id: %w", err)
	}

	materialIDs, _ := r.GetMaterialIDs(ctx, id)
	baseIDs, _ := r.GetBaseIDs(ctx, id)
	sale.MaterialIDs = materialIDs
	sale.BaseIDs = baseIDs

	return sale, nil
}

// List retrieves sales with optional filtering.
func (r *SaleRepository) List(ctx context.Context, filter domain.SaleFilter) ([]domain.Sale, int, error) {
	var conditions []string
	var args []interface{}

	if filter.Status != "" {
		conditions = append(conditions, "s.status = ?")
		args = append(args, filter.Status)
	}
	if filter.Query != "" {
		conditions = append(conditions, "(s.name_ru LIKE ? OR s.name_en LIKE ?)")
		args = append(args, "%"+filter.Query+"%", "%"+filter.Query+"%")
	}
	if filter.MaterialID > 0 {
		conditions = append(conditions, "sm.material_id = ?")
		args = append(args, filter.MaterialID)
	}
	if filter.BaseID > 0 {
		conditions = append(conditions, "sb.base_id = ?")
		args = append(args, filter.BaseID)
	}

	whereClause := ""
	if len(conditions) > 0 {
		whereClause = "WHERE " + strings.Join(conditions, " AND ")
	}

	joinClause := ""
	if filter.MaterialID > 0 {
		joinClause += " JOIN sale_materials sm ON s.id = sm.sale_id"
	}
	if filter.BaseID > 0 {
		joinClause += " JOIN sales_bases sb ON s.id = sb.sale_id"
	}

	countQuery := fmt.Sprintf("SELECT COUNT(DISTINCT s.id) FROM sales s%s %s", joinClause, whereClause)
	var total int
	if err := r.db.QueryRowContext(ctx, countQuery, args...).Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("count sales: %w", err)
	}

	sortBy := "s.sort_order"
	sortOrder := "ASC"
	if filter.SortBy != "" {
		sortBy = "s." + filter.SortBy
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

	listQuery := fmt.Sprintf(`SELECT DISTINCT s.%s FROM sales s%s %s ORDER BY %s %s LIMIT ? OFFSET ?`,
		saleColumns, joinClause, whereClause, sortBy, sortOrder)
	args = append(args, limit, offset)

	rows, err := r.db.QueryContext(ctx, listQuery, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("list sales: %w", err)
	}
	defer rows.Close()

	var sales []domain.Sale
	for rows.Next() {
		sale, err := r.scanSale(rows)
		if err != nil {
			return nil, 0, fmt.Errorf("scan sale: %w", err)
		}
		sales = append(sales, *sale)
	}

	// Load material and base associations for all sales
	materialIDs, _ := r.GetBulkMaterialIDs(ctx, sales)
	baseIDs, _ := r.GetBulkBaseIDs(ctx, sales)
	for i := range sales {
		sales[i].MaterialIDs = materialIDs[sales[i].ID]
		sales[i].BaseIDs = baseIDs[sales[i].ID]
	}

	return sales, total, nil
}

// Update updates a sale.
func (r *SaleRepository) Update(ctx context.Context, sale *domain.Sale) error {
	query := `UPDATE sales SET name_ru = ?, name_en = ?, description = ?, image_path = ?, sale_path = ?, price = ?,
		year = ?, technique = ?, width = ?, height = ?, status = ?, sort_order = ?, sold = ?, updated_at = ? WHERE id = ?`

	sale.UpdatedAt = time.Now()

	_, err := r.db.ExecContext(ctx, query,
		sale.NameRu, sale.NameEn, nullString(sale.Description), sale.ImagePath, sale.SalePath,
		sale.Price, nullInt(int64(sale.Year)),
		nullString(sale.Technique),
		nullInt(int64(sale.Width)), nullInt(int64(sale.Height)),
		sale.Status, sale.SortOrder, sale.Sold, sale.UpdatedAt, sale.ID,
	)
	if err != nil {
		return fmt.Errorf("update sale: %w", err)
	}

	if sale.MaterialIDs != nil {
		if err := r.SetMaterials(ctx, sale.ID, sale.MaterialIDs); err != nil {
			return err
		}
	}
	if sale.BaseIDs != nil {
		if err := r.SetBases(ctx, sale.ID, sale.BaseIDs); err != nil {
			return err
		}
	}

	return nil
}

// Delete deletes a sale by ID.
func (r *SaleRepository) Delete(ctx context.Context, id int64) error {
	if _, err := r.db.ExecContext(ctx, "DELETE FROM sale_materials WHERE sale_id = ?", id); err != nil {
		return fmt.Errorf("delete sale materials: %w", err)
	}
	if _, err := r.db.ExecContext(ctx, "DELETE FROM sales_bases WHERE sale_id = ?", id); err != nil {
		return fmt.Errorf("delete sale bases: %w", err)
	}

	result, err := r.db.ExecContext(ctx, "DELETE FROM sales WHERE id = ?", id)
	if err != nil {
		return fmt.Errorf("delete sale: %w", err)
	}
	affected, _ := result.RowsAffected()
	if affected == 0 {
		return domain.ErrNotFound
	}
	return nil
}

// SetMaterials sets material associations for a sale.
func (r *SaleRepository) SetMaterials(ctx context.Context, saleID int64, materialIDs []int64) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin tx: %w", err)
	}
	defer tx.Rollback()

	if _, err := tx.ExecContext(ctx, "DELETE FROM sale_materials WHERE sale_id = ?", saleID); err != nil {
		return fmt.Errorf("delete sale materials: %w", err)
	}
	for _, materialID := range materialIDs {
		if _, err := tx.ExecContext(ctx, "INSERT INTO sale_materials (sale_id, material_id) VALUES (?, ?)", saleID, materialID); err != nil {
			return fmt.Errorf("insert sale material: %w", err)
		}
	}
	return tx.Commit()
}

// SetBases sets base associations for a sale.
func (r *SaleRepository) SetBases(ctx context.Context, saleID int64, baseIDs []int64) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin tx: %w", err)
	}
	defer tx.Rollback()

	if _, err := tx.ExecContext(ctx, "DELETE FROM sales_bases WHERE sale_id = ?", saleID); err != nil {
		return fmt.Errorf("delete sale bases: %w", err)
	}
	for _, baseID := range baseIDs {
		if _, err := tx.ExecContext(ctx, "INSERT INTO sales_bases (sale_id, base_id) VALUES (?, ?)", saleID, baseID); err != nil {
			return fmt.Errorf("insert sale base: %w", err)
		}
	}
	return tx.Commit()
}

// GetMaterialIDs retrieves material IDs for a sale.
func (r *SaleRepository) GetMaterialIDs(ctx context.Context, saleID int64) ([]int64, error) {
	rows, err := r.db.QueryContext(ctx, "SELECT material_id FROM sales_materials WHERE sale_id = ? ORDER BY material_id", saleID)
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

// GetBaseIDs retrieves base IDs for a sale.
func (r *SaleRepository) GetBaseIDs(ctx context.Context, saleID int64) ([]int64, error) {
	rows, err := r.db.QueryContext(ctx, "SELECT base_id FROM sales_bases WHERE sale_id = ? ORDER BY base_id", saleID)
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

// GetBulkMaterialIDs retrieves material IDs for multiple sales at once.
func (r *SaleRepository) GetBulkMaterialIDs(ctx context.Context, sales []domain.Sale) (map[int64][]int64, error) {
	if len(sales) == 0 {
		return nil, nil
	}
	ids := make([]int64, len(sales))
	for i, s := range sales {
		ids[i] = s.ID
	}
	result := make(map[int64][]int64, len(sales))

	query := "SELECT sale_id, material_id FROM sales_materials WHERE sale_id IN (?" + strings.Repeat(",?", len(ids)-1) + ") ORDER BY sale_id, material_id"
	args := make([]interface{}, len(ids))
	for i, id := range ids {
		args[i] = id
	}

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var saleID, materialID int64
		if err := rows.Scan(&saleID, &materialID); err != nil {
			return nil, err
		}
		result[saleID] = append(result[saleID], materialID)
	}
	return result, nil
}

// GetBulkBaseIDs retrieves base IDs for multiple sales at once.
func (r *SaleRepository) GetBulkBaseIDs(ctx context.Context, sales []domain.Sale) (map[int64][]int64, error) {
	if len(sales) == 0 {
		return nil, nil
	}
	ids := make([]int64, len(sales))
	for i, s := range sales {
		ids[i] = s.ID
	}
	result := make(map[int64][]int64, len(sales))

	query := "SELECT sale_id, base_id FROM sales_bases WHERE sale_id IN (?" + strings.Repeat(",?", len(ids)-1) + ") ORDER BY sale_id, base_id"
	args := make([]interface{}, len(ids))
	for i, id := range ids {
		args[i] = id
	}

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var saleID, baseID int64
		if err := rows.Scan(&saleID, &baseID); err != nil {
			return nil, err
		}
		result[saleID] = append(result[saleID], baseID)
	}
	return result, nil
}

func nullFloat(f float64) interface{} {
	if f == 0 {
		return nil
	}
	return f
}
