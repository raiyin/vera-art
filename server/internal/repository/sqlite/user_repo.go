package sqlite

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/raiyin/artserver/internal/domain"
)

// UserRepository implements port.UserRepository.
type UserRepository struct {
	db *sql.DB
}

// NewUserRepository creates a new UserRepository.
func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

// Create inserts a new user.
func (r *UserRepository) Create(ctx context.Context, user *domain.User) error {
	query := `INSERT INTO users (username, email, password_hash, name, role, email_verified, verification_token, verification_sent_at, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`

	now := time.Now()
	if user.CreatedAt.IsZero() {
		user.CreatedAt = now
	}
	if user.UpdatedAt.IsZero() {
		user.UpdatedAt = now
	}
	if user.Role == "" {
		user.Role = "user"
	}

	result, err := r.db.ExecContext(ctx, query,
		user.Username, user.Email, user.PasswordHash, user.Name, user.Role,
		user.EmailVerified, user.VerificationToken, user.VerificationSentAt,
		user.CreatedAt, user.UpdatedAt,
	)
	if err != nil {
		if strings.Contains(err.Error(), "UNIQUE") {
			return domain.ErrDuplicate
		}
		return fmt.Errorf("create user: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("get last insert id: %w", err)
	}
	user.ID = id
	return nil
}

// GetByID retrieves a user by ID.
func (r *UserRepository) GetByID(ctx context.Context, id int64) (*domain.User, error) {
	query := `SELECT id, username, email, password_hash, name, role, avatar_path, email_verified, verification_token, verification_sent_at, created_at, updated_at
		FROM users WHERE id = ?`

	user := &domain.User{}
	var avatarPath, verificationToken sql.NullString
	var verificationSentAt sql.NullTime

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&user.ID, &user.Username, &user.Email, &user.PasswordHash, &user.Name, &user.Role,
		&avatarPath, &user.EmailVerified, &verificationToken, &verificationSentAt,
		&user.CreatedAt, &user.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, domain.ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("get user by id: %w", err)
	}

	if avatarPath.Valid {
		user.AvatarPath = avatarPath.String
	}
	if verificationToken.Valid {
		user.VerificationToken = verificationToken.String
	}
	if verificationSentAt.Valid {
		user.VerificationSentAt = &verificationSentAt.Time
	}

	return user, nil
}

// GetByUsername retrieves a user by username.
func (r *UserRepository) GetByUsername(ctx context.Context, username string) (*domain.User, error) {
	query := `SELECT id, username, email, password_hash, name, role, avatar_path, email_verified, verification_token, verification_sent_at, created_at, updated_at
		FROM users WHERE username = ?`

	user := &domain.User{}
	var avatarPath, verificationToken sql.NullString
	var verificationSentAt sql.NullTime

	err := r.db.QueryRowContext(ctx, query, username).Scan(
		&user.ID, &user.Username, &user.Email, &user.PasswordHash, &user.Name, &user.Role,
		&avatarPath, &user.EmailVerified, &verificationToken, &verificationSentAt,
		&user.CreatedAt, &user.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, domain.ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("get user by username: %w", err)
	}

	if avatarPath.Valid {
		user.AvatarPath = avatarPath.String
	}
	if verificationToken.Valid {
		user.VerificationToken = verificationToken.String
	}
	if verificationSentAt.Valid {
		user.VerificationSentAt = &verificationSentAt.Time
	}

	return user, nil
}

// GetByEmail retrieves a user by email.
func (r *UserRepository) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	query := `SELECT id, username, email, password_hash, name, role, avatar_path, email_verified, verification_token, verification_sent_at, created_at, updated_at
		FROM users WHERE email = ?`

	user := &domain.User{}
	var avatarPath, verificationToken sql.NullString
	var verificationSentAt sql.NullTime

	err := r.db.QueryRowContext(ctx, query, email).Scan(
		&user.ID, &user.Username, &user.Email, &user.PasswordHash, &user.Name, &user.Role,
		&avatarPath, &user.EmailVerified, &verificationToken, &verificationSentAt,
		&user.CreatedAt, &user.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, domain.ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("get user by email: %w", err)
	}

	if avatarPath.Valid {
		user.AvatarPath = avatarPath.String
	}
	if verificationToken.Valid {
		user.VerificationToken = verificationToken.String
	}
	if verificationSentAt.Valid {
		user.VerificationSentAt = &verificationSentAt.Time
	}

	return user, nil
}

// GetByVerificationToken retrieves a user by verification token.
func (r *UserRepository) GetByVerificationToken(ctx context.Context, token string) (*domain.User, error) {
	query := `SELECT id, username, email, password_hash, name, role, avatar_path, email_verified, verification_token, verification_sent_at, created_at, updated_at
		FROM users WHERE verification_token = ?`

	user := &domain.User{}
	var avatarPath, verificationToken sql.NullString
	var verificationSentAt sql.NullTime

	err := r.db.QueryRowContext(ctx, query, token).Scan(
		&user.ID, &user.Username, &user.Email, &user.PasswordHash, &user.Name, &user.Role,
		&avatarPath, &user.EmailVerified, &verificationToken, &verificationSentAt,
		&user.CreatedAt, &user.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, domain.ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("get user by verification token: %w", err)
	}

	if avatarPath.Valid {
		user.AvatarPath = avatarPath.String
	}
	if verificationToken.Valid {
		user.VerificationToken = verificationToken.String
	}
	if verificationSentAt.Valid {
		user.VerificationSentAt = &verificationSentAt.Time
	}

	return user, nil
}

// List retrieves users with optional filtering.
func (r *UserRepository) List(ctx context.Context, filter domain.UserFilter) ([]domain.User, int, error) {
	var conditions []string
	var args []interface{}

	if filter.Role != "" {
		conditions = append(conditions, "role = ?")
		args = append(args, filter.Role)
	}
	if filter.Query != "" {
		conditions = append(conditions, "(username LIKE ? OR email LIKE ? OR name LIKE ?)")
		q := "%" + filter.Query + "%"
		args = append(args, q, q, q)
	}

	whereClause := ""
	if len(conditions) > 0 {
		whereClause = "WHERE " + strings.Join(conditions, " AND ")
	}

	// Count
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM users %s", whereClause)
	var total int
	if err := r.db.QueryRowContext(ctx, countQuery, args...).Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("count users: %w", err)
	}

	// List
	listQuery := fmt.Sprintf(`SELECT id, username, email, password_hash, name, role, avatar_path, email_verified, verification_token, verification_sent_at, created_at, updated_at
		FROM users %s ORDER BY created_at DESC`, whereClause)

	rows, err := r.db.QueryContext(ctx, listQuery, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("list users: %w", err)
	}
	defer rows.Close()

	var users []domain.User
	for rows.Next() {
		var u domain.User
		var avatarPath, email, name, verificationToken sql.NullString
		var verificationSentAt sql.NullTime

		if err := rows.Scan(
			&u.ID, &u.Username, &email, &u.PasswordHash, &name, &u.Role,
			&avatarPath, &u.EmailVerified, &verificationToken, &verificationSentAt,
			&u.CreatedAt, &u.UpdatedAt,
		); err != nil {
			return nil, 0, fmt.Errorf("scan user: %w", err)
		}

		if email.Valid {
			u.Email = email.String
		}
		if name.Valid {
			u.Name = name.String
		}
		if avatarPath.Valid {
			u.AvatarPath = avatarPath.String
		}
		if verificationToken.Valid {
			u.VerificationToken = verificationToken.String
		}
		if verificationSentAt.Valid {
			u.VerificationSentAt = &verificationSentAt.Time
		}

		users = append(users, u)
	}

	return users, total, nil
}

// Update updates a user.
func (r *UserRepository) Update(ctx context.Context, user *domain.User) error {
	query := `UPDATE users SET username = ?, email = ?, password_hash = ?, name = ?, role = ?, avatar_path = ?,
		email_verified = ?, verification_token = ?, verification_sent_at = ?, updated_at = ? WHERE id = ?`

	user.UpdatedAt = time.Now()

	_, err := r.db.ExecContext(ctx, query,
		user.Username, user.Email, user.PasswordHash, user.Name, user.Role, nullString(user.AvatarPath),
		user.EmailVerified, nullString(user.VerificationToken), nullTime(user.VerificationSentAt),
		user.UpdatedAt, user.ID,
	)
	if err != nil {
		return fmt.Errorf("update user: %w", err)
	}
	return nil
}

// Delete deletes a user by ID.
func (r *UserRepository) Delete(ctx context.Context, id int64) error {
	result, err := r.db.ExecContext(ctx, "DELETE FROM users WHERE id = ?", id)
	if err != nil {
		return fmt.Errorf("delete user: %w", err)
	}
	affected, _ := result.RowsAffected()
	if affected == 0 {
		return domain.ErrNotFound
	}
	return nil
}

// Helper functions for nullable fields.
func nullString(s string) interface{} {
	if s == "" {
		return nil
	}
	return s
}

func nullTime(t *time.Time) interface{} {
	if t == nil {
		return nil
	}
	return *t
}
