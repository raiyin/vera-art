package sqlite

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

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

const workColumns = `id, title, description, image_path, year, technique, width, height, status, sort_order, created_at, updated_at`

func (r *WorkRepository) scanWork(scanner interface {
	Scan(dest ...interface{}) error
}) (*domain.Work, error) {
	w := &domain.Work{}
	var description, technique sql.NullString
	var year, width, height sql.NullInt64

	err := scanner.Scan(
		&w.ID, &w.Title, &description, &w.ImagePath,
		&year, &technique, &width, &height,
		&w.Status, &w.SortOrder,
		&w.CreatedAt, &w.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	if description.Valid {
		w.Description = description.String
	}
	if technique.Valid {
		w.Technique = technique.String
	}
	if width.Valid {
		w.Width = int(width.Int64)
	}
	if height.Valid {
		w.Height = int(height.Int64)
	}
	if year.Valid {
		w.Year = int(year.Int64)
	}

	return w, nil
}

// Create inserts a new work.
func (r *WorkRepository) Create(ctx context.Context, work *domain.Work) error {
	query := `INSERT INTO works (title, description, image_path, year, technique, width, height, status, sort_order, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`

	now := time.Now()
	if work.CreatedAt.IsZero() {
		work.CreatedAt = now
	}
	if work.UpdatedAt.IsZero() {
		work.UpdatedAt = now
	}
	if work.Status == "" {
		work.Status = "draft"
	}

	result, err := r.db.ExecContext(ctx, query,
		work.Title, nullString(work.Description), work.ImagePath,
		nullInt(int64(work.Year)), nullString(work.Technique),
		nullInt(int64(work.Width)), nullInt(int64(work.Height)),
		work.Status, work.SortOrder, work.CreatedAt, work.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("create work: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("get last insert id: %w", err)
	}
	work.ID = id

	// Set material and base associations
	if len(work.MaterialIDs) > 0 {
		if err := r.SetMaterials(ctx, work.ID, work.MaterialIDs); err != nil {
			return err
		}
	}
	if len(work.BaseIDs) > 0 {
		if err := r.SetBases(ctx, work.ID, work.BaseIDs); err != nil {
			return err
		}
	}

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

	// Load associations
	materialIDs, _ := r.GetMaterialIDs(ctx, id)
	baseIDs, _ := r.GetBaseIDs(ctx, id)
	work.MaterialIDs = materialIDs
	work.BaseIDs = baseIDs

	return work, nil
}

// List retrieves works with optional filtering.
func (r *WorkRepository) List(ctx context.Context, filter domain.WorkFilter) ([]domain.Work, int, error) {
	var conditions []string
	var args []interface{}

	if filter.Status != "" {
		conditions = append(conditions, "w.status = ?")
		args = append(args, filter.Status)
	}
	if filter.Query != "" {
		conditions = append(conditions, "w.title LIKE ?")
		args = append(args, "%"+filter.Query+"%")
	}
	if filter.MaterialID > 0 {
		conditions = append(conditions, "wm.material_id = ?")
		args = append(args, filter.MaterialID)
	}
	if filter.BaseID > 0 {
		conditions = append(conditions, "wb.base_id = ?")
		args = append(args, filter.BaseID)
	}

	whereClause := ""
	if len(conditions) > 0 {
		whereClause = "WHERE " + strings.Join(conditions, " AND ")
	}

	joinClause := ""
	if filter.MaterialID > 0 {
		joinClause += " JOIN works_materials wm ON w.id = wm.work_id"
	}
	if filter.BaseID > 0 {
		joinClause += " JOIN works_bases wb ON w.id = wb.work_id"
	}

	// Count
	countQuery := fmt.Sprintf("SELECT COUNT(DISTINCT w.id) FROM works w%s %s", joinClause, whereClause)
	var total int
	if err := r.db.QueryRowContext(ctx, countQuery, args...).Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("count works: %w", err)
	}

	// Sort
	sortBy := "w.sort_order"
	sortOrder := "ASC"
	if filter.SortBy != "" {
		sortBy = "w." + filter.SortBy
	}
	if filter.SortOrder == "desc" {
		sortOrder = "DESC"
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

	listQuery := fmt.Sprintf(`SELECT DISTINCT w.%s FROM works w%s %s ORDER BY %s %s LIMIT ? OFFSET ?`,
		workColumns, joinClause, whereClause, sortBy, sortOrder)
	args = append(args, limit, offset)

	rows, err := r.db.QueryContext(ctx, listQuery, args...)
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

	// Load material and base associations for all works
	materialIDs, _ := r.GetBulkMaterialIDs(ctx, works)
	baseIDs, _ := r.GetBulkBaseIDs(ctx, works)
	for i := range works {
		works[i].MaterialIDs = materialIDs[works[i].ID]
		works[i].BaseIDs = baseIDs[works[i].ID]
	}

	return works, total, nil
}

// Update updates a work.
func (r *WorkRepository) Update(ctx context.Context, work *domain.Work) error {
	query := `UPDATE works SET title = ?, description = ?, image_path = ?, year = ?, technique = ?,
		width = ?, height = ?, status = ?, sort_order = ?, updated_at = ? WHERE id = ?`

	work.UpdatedAt = time.Now()

	_, err := r.db.ExecContext(ctx, query,
		work.Title, nullString(work.Description), work.ImagePath,
		nullInt(int64(work.Year)), nullString(work.Technique),
		nullInt(int64(work.Width)), nullInt(int64(work.Height)),
		work.Status, work.SortOrder, work.UpdatedAt, work.ID,
	)
	if err != nil {
		return fmt.Errorf("update work: %w", err)
	}

	// Update associations
	if work.MaterialIDs != nil {
		if err := r.SetMaterials(ctx, work.ID, work.MaterialIDs); err != nil {
			return err
		}
	}
	if work.BaseIDs != nil {
		if err := r.SetBases(ctx, work.ID, work.BaseIDs); err != nil {
			return err
		}
	}

	return nil
}

// Delete deletes a work by ID.
func (r *WorkRepository) Delete(ctx context.Context, id int64) error {
	// Delete associations first
	if _, err := r.db.ExecContext(ctx, "DELETE FROM works_materials WHERE work_id = ?", id); err != nil {
		return fmt.Errorf("delete work materials: %w", err)
	}
	if _, err := r.db.ExecContext(ctx, "DELETE FROM works_bases WHERE work_id = ?", id); err != nil {
		return fmt.Errorf("delete work bases: %w", err)
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

// SetMaterials sets material associations for a work.
func (r *WorkRepository) SetMaterials(ctx context.Context, workID int64, materialIDs []int64) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin tx: %w", err)
	}
	defer tx.Rollback()

	if _, err := tx.ExecContext(ctx, "DELETE FROM works_materials WHERE work_id = ?", workID); err != nil {
		return fmt.Errorf("delete work materials: %w", err)
	}

	for _, materialID := range materialIDs {
		if _, err := tx.ExecContext(ctx, "INSERT INTO works_materials (work_id, material_id) VALUES (?, ?)", workID, materialID); err != nil {
			return fmt.Errorf("insert work material: %w", err)
		}
	}

	return tx.Commit()
}

// SetBases sets base associations for a work.
func (r *WorkRepository) SetBases(ctx context.Context, workID int64, baseIDs []int64) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin tx: %w", err)
	}
	defer tx.Rollback()

	if _, err := tx.ExecContext(ctx, "DELETE FROM works_bases WHERE work_id = ?", workID); err != nil {
		return fmt.Errorf("delete work bases: %w", err)
	}

	for _, baseID := range baseIDs {
		if _, err := tx.ExecContext(ctx, "INSERT INTO works_bases (work_id, base_id) VALUES (?, ?)", workID, baseID); err != nil {
			return fmt.Errorf("insert work base: %w", err)
		}
	}

	return tx.Commit()
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

// GetBaseIDs retrieves base IDs for a work.
func (r *WorkRepository) GetBaseIDs(ctx context.Context, workID int64) ([]int64, error) {
	rows, err := r.db.QueryContext(ctx, "SELECT base_id FROM works_bases WHERE work_id = ? ORDER BY base_id", workID)
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

// GetBulkBaseIDs retrieves base IDs for multiple works at once.
func (r *WorkRepository) GetBulkBaseIDs(ctx context.Context, works []domain.Work) (map[int64][]int64, error) {
	if len(works) == 0 {
		return nil, nil
	}
	ids := make([]int64, len(works))
	for i, w := range works {
		ids[i] = w.ID
	}
	result := make(map[int64][]int64, len(works))

	query := "SELECT work_id, base_id FROM works_bases WHERE work_id IN (?" + strings.Repeat(",?", len(ids)-1) + ") ORDER BY work_id, base_id"
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
		var workID, baseID int64
		if err := rows.Scan(&workID, &baseID); err != nil {
			return nil, err
		}
		result[workID] = append(result[workID], baseID)
	}
	return result, nil
}

func nullInt(n int64) interface{} {
	if n == 0 {
		return nil
	}
	return n
}
