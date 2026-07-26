package sqlite

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/raiyin/artserver/internal/domain"
)

// ChatThreadRepository implements port.ChatThreadRepository.
type ChatThreadRepository struct {
	db *sql.DB
}

// NewChatThreadRepository creates a new ChatThreadRepository.
func NewChatThreadRepository(db *sql.DB) *ChatThreadRepository {
	return &ChatThreadRepository{db: db}
}

func (r *ChatThreadRepository) Create(ctx context.Context, thread *domain.ChatThread) error {
	query := `INSERT INTO chat_threads (user_id, subject, status, created_at, updated_at) VALUES (?, ?, ?, ?, ?)`
	now := time.Now()
	if thread.CreatedAt.IsZero() {
		thread.CreatedAt = now
	}
	if thread.UpdatedAt.IsZero() {
		thread.UpdatedAt = now
	}
	if thread.Status == "" {
		thread.Status = "open"
	}

	result, err := r.db.ExecContext(ctx, query,
		thread.UserID, nullString(thread.Subject), thread.Status, thread.CreatedAt, thread.UpdatedAt,
	)
	if err != nil {
		return err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return err
	}
	thread.ID = id
	return nil
}

func (r *ChatThreadRepository) GetByID(ctx context.Context, id int64) (*domain.ChatThread, error) {
	t := &domain.ChatThread{}
	var subject sql.NullString
	err := r.db.QueryRowContext(ctx,
		"SELECT id, user_id, subject, status, created_at, updated_at FROM chat_threads WHERE id = ?", id,
	).Scan(&t.ID, &t.UserID, &subject, &t.Status, &t.CreatedAt, &t.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, domain.ErrNotFound
	}
	if subject.Valid {
		t.Subject = subject.String
	}
	return t, err
}

func (r *ChatThreadRepository) ListByUser(ctx context.Context, userID int64) ([]domain.ChatThread, error) {
	rows, err := r.db.QueryContext(ctx,
		"SELECT id, user_id, subject, status, created_at, updated_at FROM chat_threads WHERE user_id = ? ORDER BY updated_at DESC", userID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var threads []domain.ChatThread
	for rows.Next() {
		var t domain.ChatThread
		var subject sql.NullString
		if err := rows.Scan(&t.ID, &t.UserID, &subject, &t.Status, &t.CreatedAt, &t.UpdatedAt); err != nil {
			return nil, err
		}
		if subject.Valid {
			t.Subject = subject.String
		}
		threads = append(threads, t)
	}
	return threads, nil
}

func (r *ChatThreadRepository) List(ctx context.Context) ([]domain.ChatThread, error) {
	rows, err := r.db.QueryContext(ctx,
		"SELECT id, user_id, subject, status, created_at, updated_at FROM chat_threads ORDER BY updated_at DESC",
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var threads []domain.ChatThread
	for rows.Next() {
		var t domain.ChatThread
		var subject sql.NullString
		if err := rows.Scan(&t.ID, &t.UserID, &subject, &t.Status, &t.CreatedAt, &t.UpdatedAt); err != nil {
			return nil, err
		}
		if subject.Valid {
			t.Subject = subject.String
		}
		threads = append(threads, t)
	}
	return threads, nil
}

func (r *ChatThreadRepository) Update(ctx context.Context, thread *domain.ChatThread) error {
	thread.UpdatedAt = time.Now()
	_, err := r.db.ExecContext(ctx,
		"UPDATE chat_threads SET subject = ?, status = ?, updated_at = ? WHERE id = ?",
		nullString(thread.Subject), thread.Status, thread.UpdatedAt, thread.ID,
	)
	return err
}

// ChatMessageRepository implements port.ChatMessageRepository.
type ChatMessageRepository struct {
	db *sql.DB
}

// NewChatMessageRepository creates a new ChatMessageRepository.
func NewChatMessageRepository(db *sql.DB) *ChatMessageRepository {
	return &ChatMessageRepository{db: db}
}

func (r *ChatMessageRepository) Create(ctx context.Context, message *domain.ChatMessage) error {
	query := `INSERT INTO chat_messages (thread_id, user_id, content, is_admin, is_read, created_at) VALUES (?, ?, ?, ?, ?, ?)`
	now := time.Now()
	if message.CreatedAt.IsZero() {
		message.CreatedAt = now
	}

	result, err := r.db.ExecContext(ctx, query,
		message.ThreadID, message.UserID, message.Content, message.IsAdmin, message.IsRead, message.CreatedAt,
	)
	if err != nil {
		return err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return err
	}
	message.ID = id
	return nil
}

func (r *ChatMessageRepository) ListByThread(ctx context.Context, threadID int64) ([]domain.ChatMessage, error) {
	rows, err := r.db.QueryContext(ctx,
		"SELECT id, thread_id, user_id, content, is_admin, is_read, read_at, created_at FROM chat_messages WHERE thread_id = ? ORDER BY created_at", threadID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var messages []domain.ChatMessage
	for rows.Next() {
		var m domain.ChatMessage
		var readAt sql.NullTime
		if err := rows.Scan(&m.ID, &m.ThreadID, &m.UserID, &m.Content, &m.IsAdmin, &m.IsRead, &readAt, &m.CreatedAt); err != nil {
			return nil, err
		}
		if readAt.Valid {
			m.ReadAt = &readAt.Time
		}
		messages = append(messages, m)
	}
	return messages, nil
}

func (r *ChatMessageRepository) MarkAsRead(ctx context.Context, messageID int64) error {
	now := time.Now()
	_, err := r.db.ExecContext(ctx,
		"UPDATE chat_messages SET is_read = ?, read_at = ? WHERE id = ?", true, now, messageID,
	)
	return err
}

func (r *ChatMessageRepository) GetUnreadCount(ctx context.Context, threadID int64) (int, error) {
	var count int
	err := r.db.QueryRowContext(ctx,
		"SELECT COUNT(*) FROM chat_messages WHERE thread_id = ? AND is_read = 0 AND is_admin = 0", threadID,
	).Scan(&count)
	return count, err
}

func (r *ChatMessageRepository) PollNewMessages(ctx context.Context, threadID int64, lastMessageID int64) ([]domain.ChatMessage, bool, error) {
	rows, err := r.db.QueryContext(ctx,
		`SELECT id, thread_id, user_id, content, is_admin, is_read, read_at, created_at
		FROM chat_messages WHERE thread_id = ? AND id > ? ORDER BY created_at`, threadID, lastMessageID,
	)
	if err != nil {
		return nil, false, err
	}
	defer rows.Close()

	var messages []domain.ChatMessage
	for rows.Next() {
		var m domain.ChatMessage
		var readAt sql.NullTime
		if err := rows.Scan(&m.ID, &m.ThreadID, &m.UserID, &m.Content, &m.IsAdmin, &m.IsRead, &readAt, &m.CreatedAt); err != nil {
			return nil, false, err
		}
		if readAt.Valid {
			m.ReadAt = &readAt.Time
		}
		messages = append(messages, m)
	}

	hasMore := len(messages) > 0
	return messages, hasMore, nil
}

// MasterClassRepository implements port.MasterClassRepository.
type MasterClassRepository struct {
	db *sql.DB
}

// NewMasterClassRepository creates a new MasterClassRepository.
func NewMasterClassRepository(db *sql.DB) *MasterClassRepository {
	return &MasterClassRepository{db: db}
}

func (r *MasterClassRepository) Create(ctx context.Context, mc *domain.MasterClass) error {
	query := `INSERT INTO master_classes (title, description, price, image_path, status, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, ?)`
	now := time.Now()
	if mc.CreatedAt.IsZero() {
		mc.CreatedAt = now
	}
	if mc.UpdatedAt.IsZero() {
		mc.UpdatedAt = now
	}
	if mc.Status == "" {
		mc.Status = "draft"
	}

	result, err := r.db.ExecContext(ctx, query,
		mc.Title, nullString(mc.Description), mc.Price, mc.ImagePath, mc.Status, mc.CreatedAt, mc.UpdatedAt,
	)
	if err != nil {
		return err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return err
	}
	mc.ID = id
	return nil
}

func (r *MasterClassRepository) GetByID(ctx context.Context, id int64) (*domain.MasterClass, error) {
	mc := &domain.MasterClass{}
	var description sql.NullString
	err := r.db.QueryRowContext(ctx,
		"SELECT id, title, description, price, image_path, status, created_at, updated_at FROM master_classes WHERE id = ?", id,
	).Scan(&mc.ID, &mc.Title, &description, &mc.Price, &mc.ImagePath, &mc.Status, &mc.CreatedAt, &mc.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, domain.ErrNotFound
	}
	if description.Valid {
		mc.Description = description.String
	}
	return mc, err
}

func (r *MasterClassRepository) List(ctx context.Context) ([]domain.MasterClass, int, error) {
	var total int
	r.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM master_classes").Scan(&total)

	rows, err := r.db.QueryContext(ctx, "SELECT id, title, description, price, image_path, status, created_at, updated_at FROM master_classes ORDER BY created_at DESC")
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var classes []domain.MasterClass
	for rows.Next() {
		var mc domain.MasterClass
		var description sql.NullString
		if err := rows.Scan(&mc.ID, &mc.Title, &description, &mc.Price, &mc.ImagePath, &mc.Status, &mc.CreatedAt, &mc.UpdatedAt); err != nil {
			return nil, 0, err
		}
		if description.Valid {
			mc.Description = description.String
		}
		classes = append(classes, mc)
	}
	return classes, total, nil
}

func (r *MasterClassRepository) Update(ctx context.Context, mc *domain.MasterClass) error {
	mc.UpdatedAt = time.Now()
	_, err := r.db.ExecContext(ctx,
		"UPDATE master_classes SET title = ?, description = ?, price = ?, image_path = ?, status = ?, updated_at = ? WHERE id = ?",
		mc.Title, nullString(mc.Description), mc.Price, mc.ImagePath, mc.Status, mc.UpdatedAt, mc.ID,
	)
	return err
}

func (r *MasterClassRepository) Delete(ctx context.Context, id int64) error {
	result, err := r.db.ExecContext(ctx, "DELETE FROM master_classes WHERE id = ?", id)
	if err != nil {
		return err
	}
	affected, _ := result.RowsAffected()
	if affected == 0 {
		return domain.ErrNotFound
	}
	return nil
}

// TagRepository implements port.TagRepository.
type TagRepository struct {
	db *sql.DB
}

// NewTagRepository creates a new TagRepository.
func NewTagRepository(db *sql.DB) *TagRepository {
	return &TagRepository{db: db}
}

func (r *TagRepository) Create(ctx context.Context, tag *domain.Tag) error {
	query := `INSERT INTO tags (name_ru, name_en, slug, created_at) VALUES (?, ?, ?, ?)`
	now := time.Now()
	if tag.CreatedAt.IsZero() {
		tag.CreatedAt = now
	}
	result, err := r.db.ExecContext(ctx, query, tag.NameRu, tag.NameEn, tag.Slug, tag.CreatedAt)
	if err != nil {
		return fmt.Errorf("create tag: %w", err)
	}
	id, err := result.LastInsertId()
	if err != nil {
		return err
	}
	tag.ID = id
	return nil
}

func (r *TagRepository) Update(ctx context.Context, tag *domain.Tag) error {
	query := `UPDATE tags SET name_ru = ?, name_en = ?, slug = ? WHERE id = ?`
	result, err := r.db.ExecContext(ctx, query, tag.NameRu, tag.NameEn, tag.Slug, tag.ID)
	if err != nil {
		return fmt.Errorf("update tag: %w", err)
	}
	affected, _ := result.RowsAffected()
	if affected == 0 {
		return domain.ErrNotFound
	}
	return nil
}

func (r *TagRepository) List(ctx context.Context) ([]domain.Tag, error) {
	rows, err := r.db.QueryContext(ctx, "SELECT id, name_ru, name_en, slug, created_at FROM tags ORDER BY name_ru")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tags []domain.Tag
	for rows.Next() {
		var t domain.Tag
		if err := rows.Scan(&t.ID, &t.NameRu, &t.NameEn, &t.Slug, &t.CreatedAt); err != nil {
			return nil, err
		}
		tags = append(tags, t)
	}
	return tags, nil
}

func (r *TagRepository) Delete(ctx context.Context, id int64) error {
	result, err := r.db.ExecContext(ctx, "DELETE FROM tags WHERE id = ?", id)
	if err != nil {
		return err
	}
	affected, _ := result.RowsAffected()
	if affected == 0 {
		return domain.ErrNotFound
	}
	return nil
}

// MaterialRepository implements port.MaterialRepository.
type MaterialRepository struct {
	db *sql.DB
}

// NewMaterialRepository creates a new MaterialRepository.
func NewMaterialRepository(db *sql.DB) *MaterialRepository {
	return &MaterialRepository{db: db}
}

func (r *MaterialRepository) Create(ctx context.Context, material *domain.Material) error {
	query := `INSERT INTO materials (name_ru, name_en, created_at) VALUES (?, ?, ?)`
	now := time.Now()
	if material.CreatedAt.IsZero() {
		material.CreatedAt = now
	}
	result, err := r.db.ExecContext(ctx, query, material.NameRu, material.NameEn, material.CreatedAt)
	if err != nil {
		return err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return err
	}
	material.ID = id
	return nil
}

func (r *MaterialRepository) List(ctx context.Context) ([]domain.Material, error) {
	rows, err := r.db.QueryContext(ctx, "SELECT id, name_ru, name_en, created_at FROM materials ORDER BY name_ru")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var materials []domain.Material
	for rows.Next() {
		var m domain.Material
		if err := rows.Scan(&m.ID, &m.NameRu, &m.NameEn, &m.CreatedAt); err != nil {
			return nil, err
		}
		materials = append(materials, m)
	}
	return materials, nil
}

func (r *MaterialRepository) Delete(ctx context.Context, id int64) error {
	result, err := r.db.ExecContext(ctx, "DELETE FROM materials WHERE id = ?", id)
	if err != nil {
		return err
	}
	affected, _ := result.RowsAffected()
	if affected == 0 {
		return domain.ErrNotFound
	}
	return nil
}

// BaseRepository implements port.BaseRepository.
type BaseRepository struct {
	db *sql.DB
}

// NewBaseRepository creates a new BaseRepository.
func NewBaseRepository(db *sql.DB) *BaseRepository {
	return &BaseRepository{db: db}
}

func (r *BaseRepository) Create(ctx context.Context, base *domain.Base) error {
	query := `INSERT INTO bases (name_ru, name_en, created_at) VALUES (?, ?, ?)`
	now := time.Now()
	if base.CreatedAt.IsZero() {
		base.CreatedAt = now
	}
	result, err := r.db.ExecContext(ctx, query, base.NameRu, base.NameEn, base.CreatedAt)
	if err != nil {
		return err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return err
	}
	base.ID = id
	return nil
}

func (r *BaseRepository) List(ctx context.Context) ([]domain.Base, error) {
	rows, err := r.db.QueryContext(ctx, "SELECT id, name_ru, name_en, created_at FROM bases ORDER BY name_ru")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var bases []domain.Base
	for rows.Next() {
		var b domain.Base
		if err := rows.Scan(&b.ID, &b.NameRu, &b.NameEn, &b.CreatedAt); err != nil {
			return nil, err
		}
		bases = append(bases, b)
	}
	return bases, nil
}

func (r *BaseRepository) Delete(ctx context.Context, id int64) error {
	result, err := r.db.ExecContext(ctx, "DELETE FROM bases WHERE id = ?", id)
	if err != nil {
		return err
	}
	affected, _ := result.RowsAffected()
	if affected == 0 {
		return domain.ErrNotFound
	}
	return nil
}

// ConsentRepository implements port.ConsentRepository.
type ConsentRepository struct {
	db *sql.DB
}

// NewConsentRepository creates a new ConsentRepository.
func NewConsentRepository(db *sql.DB) *ConsentRepository {
	return &ConsentRepository{db: db}
}

func (r *ConsentRepository) Create(ctx context.Context, consent *domain.UserConsent) error {
	query := `INSERT INTO user_consents (user_id, consent_type, granted, ip_address, created_at) VALUES (?, ?, ?, ?, ?)`
	now := time.Now()
	if consent.CreatedAt.IsZero() {
		consent.CreatedAt = now
	}
	_, err := r.db.ExecContext(ctx, query,
		consent.UserID, consent.ConsentType, consent.Granted, nullString(consent.IPAddress), consent.CreatedAt,
	)
	return err
}

func (r *ConsentRepository) ListByUser(ctx context.Context, userID int64) ([]domain.UserConsent, error) {
	rows, err := r.db.QueryContext(ctx,
		"SELECT id, user_id, consent_type, granted, ip_address, created_at FROM user_consents WHERE user_id = ? ORDER BY created_at DESC", userID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var consents []domain.UserConsent
	for rows.Next() {
		var c domain.UserConsent
		var ipAddress sql.NullString
		if err := rows.Scan(&c.ID, &c.UserID, &c.ConsentType, &c.Granted, &ipAddress, &c.CreatedAt); err != nil {
			return nil, err
		}
		if ipAddress.Valid {
			c.IPAddress = ipAddress.String
		}
		consents = append(consents, c)
	}
	return consents, nil
}
