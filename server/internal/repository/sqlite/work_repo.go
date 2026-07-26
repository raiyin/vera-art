package sqlite

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/raiyin/artserver/internal/domain"
)

// WorkRepository implements port.WorkRepository.
type WorkRepository struct {
	db *sql.DB
}

// NewWorkRepository creates a new WorkRepository.
func NewWorkRepository(db *sql.DB) *WorkRepository {
	return &WorkRepository{db: db}
}

const workColumns = `id, str_id, width, height, year, name_ru, name_en, base_id, descr_ru, descr_en, work_path, images`

func (r *WorkRepository) scanWork(scanner interface {
	Scan(dest ...interface{}) error
}) (*domain.Work, error) {
	w := &domain.Work{}
	var descrRu, descrEn sql.NullString

	err := scanner.Scan(
		&w.ID, &w.StrID,
		&w.Width, &w.Height, &w.Year,
		&w.NameRu, &w.NameEn, &w.BaseID,
		&descrRu, &descrEn,
		&w.WorkPath, &w.Images,
	)
	if err != nil {
		return nil, err
	}

	if descrRu.Valid {
		w.DescrRu = descrRu.String
	}
	if descrEn.Valid {
		w.DescrEn = descrEn.String
	}

	return w, nil
}

// Create inserts a new work.
func (r *WorkRepository) Create(ctx context.Context, work *domain.Work) error {
	query := `INSERT INTO works (str_id, width, height, year, name_ru, name_en, base_id, descr_ru, descr_en, work_path, images)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`

	result, err := r.db.ExecContext(ctx, query,
		work.StrID, work.Width, work.Height, work.Year,
		work.NameRu, work.NameEn, work.BaseID,
		nullString(work.DescrRu), nullString(work.DescrEn),
		work.WorkPath, work.Images,
	)
	if err != nil {
		return fmt.Errorf("create work: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("get last insert id: %w", err)
	}
	work.ID = id

	return nil
}

// GetByID retrieves a work by ID.
func (r *WorkRepository) GetByID(ctx context.Context, id int64) (*domain.Work, error) {
	query := fmt.Sprintf("SELECT %s FROM works WHERE id = ?", workColumns)

	work, err := r.scanWork(r.db.QueryRowContext(ctx, query, id))
	if err == sql.ErrNoRows {
		return nil, domain.ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("get work by id: %w", err)
	}

	materialIDs, _ := r.GetMaterialIDs(ctx, id)
	work.MaterialIDs = materialIDs

	return work, nil
}

// List retrieves works with optional filtering.
func (r *WorkRepository) List(ctx context.Context, filter domain.WorkFilter) ([]domain.Work, int, error) {
	var conditions []string
	var args []interface{}

	if filter.Query != "" {
		conditions = append(conditions, "name_ru LIKE ? OR name_en LIKE ?")
		args = append(args, "%"+filter.Query+"%", "%"+filter.Query+"%")
	}

	whereClause := ""
	if len(conditions) > 0 {
		whereClause = "WHERE " + strings.Join(conditions, " AND ")
	}

	// Count
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM works %s", whereClause)
	var total int
	if err := r.db.QueryRowContext(ctx, countQuery, args...).Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("count works: %w", err)
	}

	// Pagination
	page := filter.Page
	if page < 1 {
		page = 1
	}
	limit := filter.Limit
	if limit < 1 {
		limit = 50
	}
	offset := (page - 1) * limit

	listQuery := fmt.Sprintf("SELECT %s FROM works %s ORDER BY id ASC LIMIT ? OFFSET ?",
		workColumns, whereClause)
	listArgs := append(args, limit, offset)

	rows, err := r.db.QueryContext(ctx, listQuery, listArgs...)
	if err != nil {
		return nil, 0, fmt.Errorf("list works: %w", err)
	}
	defer rows.Close()

	var works []domain.Work
	for rows.Next() {
		work, err := r.scanWork(rows)
		if err != nil {
			return nil, 0, fmt.Errorf("scan work: %w", err)
		}
		works = append(works, *work)
	}

	// Load material associations for all works
	materialIDs, _ := r.GetBulkMaterialIDs(ctx, works)
	for i := range works {
		works[i].MaterialIDs = materialIDs[works[i].ID]
	}

	return works, total, nil
}

// Update updates a work.
func (r *WorkRepository) Update(ctx context.Context, work *domain.Work) error {
	query := `UPDATE works SET str_id = ?, width = ?, height = ?, year = ?,
		name_ru = ?, name_en = ?, base_id = ?, descr_ru = ?, descr_en = ?, work_path = ?, images = ?
		WHERE id = ?`

	_, err := r.db.ExecContext(ctx, query,
		work.StrID, work.Width, work.Height, work.Year,
		work.NameRu, work.NameEn, work.BaseID,
		nullString(work.DescrRu), nullString(work.DescrEn),
		work.WorkPath, work.Images,
		work.ID,
	)
	if err != nil {
		return fmt.Errorf("update work: %w", err)
	}

	return nil
}

// Delete deletes a work by ID.
func (r *WorkRepository) Delete(ctx context.Context, id int64) error {
	if _, err := r.db.ExecContext(ctx, "DELETE FROM works_materials WHERE work_id = ?", id); err != nil {
		return fmt.Errorf("delete work materials: %w", err)
	}

	result, err := r.db.ExecContext(ctx, "DELETE FROM works WHERE id = ?", id)
	if err != nil {
		return fmt.Errorf("delete work: %w", err)
	}
	affected, _ := result.RowsAffected()
	if affected == 0 {
		return domain.ErrNotFound
	}
	return nil
}

// GetMaterialIDs retrieves material IDs for a work.
func (r *WorkRepository) GetMaterialIDs(ctx context.Context, workID int64) ([]int64, error) {
	rows, err := r.db.QueryContext(ctx, "SELECT material_id FROM works_materials WHERE work_id = ? ORDER BY material_id", workID)
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

// GetBulkMaterialIDs retrieves material IDs for multiple works at once.
func (r *WorkRepository) GetBulkMaterialIDs(ctx context.Context, works []domain.Work) (map[int64][]int64, error) {
	if len(works) == 0 {
		return nil, nil
	}
	ids := make([]int64, len(works))
	for i, w := range works {
		ids[i] = w.ID
	}
	result := make(map[int64][]int64, len(works))

	query := "SELECT work_id, material_id FROM works_materials WHERE work_id IN (?" + strings.Repeat(",?", len(ids)-1) + ") ORDER BY work_id, material_id"
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
		var workID, materialID int64
		if err := rows.Scan(&workID, &materialID); err != nil {
			return nil, err
		}
		result[workID] = append(result[workID], materialID)
	}
	return result, nil
}

func nullInt(n int64) interface{} {
	if n == 0 {
		return nil
	}
	return n
}
