package sqlite

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/raiyin/artserver/internal/domain"
)

// PaymentRepository implements port.PaymentRepository.
type PaymentRepository struct {
	db *sql.DB
}

// NewPaymentRepository creates a new PaymentRepository.
func NewPaymentRepository(db *sql.DB) *PaymentRepository {
	return &PaymentRepository{db: db}
}

const paymentColumns = `id, user_id, product_id, amount, currency, status, payment_method, yookassa_id, promo_code, discount, metadata, created_at, updated_at`

func (r *PaymentRepository) scanPayment(scanner interface {
	Scan(dest ...interface{}) error
}) (*domain.Payment, error) {
	p := &domain.Payment{}
	var paymentMethod, yookassaID, promoCode, metadata sql.NullString
	var discount sql.NullFloat64

	err := scanner.Scan(
		&p.ID, &p.UserID, &p.ProductID, &p.Amount, &p.Currency, &p.Status,
		&paymentMethod, &yookassaID, &promoCode, &discount, &metadata,
		&p.CreatedAt, &p.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	if paymentMethod.Valid {
		p.PaymentMethod = paymentMethod.String
	}
	if yookassaID.Valid {
		p.YooKassaID = yookassaID.String
	}
	if promoCode.Valid {
		p.PromoCode = promoCode.String
	}
	if discount.Valid {
		p.Discount = discount.Float64
	}
	if metadata.Valid {
		p.Metadata = metadata.String
	}

	return p, nil
}

// Create inserts a new payment.
func (r *PaymentRepository) Create(ctx context.Context, payment *domain.Payment) error {
	query := `INSERT INTO payments (user_id, product_id, amount, currency, status, payment_method, yookassa_id, promo_code, discount, metadata, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`

	now := time.Now()
	if payment.CreatedAt.IsZero() {
		payment.CreatedAt = now
	}
	if payment.UpdatedAt.IsZero() {
		payment.UpdatedAt = now
	}
	if payment.Currency == "" {
		payment.Currency = "RUB"
	}
	if payment.Status == "" {
		payment.Status = "pending"
	}

	result, err := r.db.ExecContext(ctx, query,
		payment.UserID, payment.ProductID, payment.Amount, payment.Currency, payment.Status,
		nullString(payment.PaymentMethod), nullString(payment.YooKassaID),
		nullString(payment.PromoCode), nullFloat(payment.Discount), nullString(payment.Metadata),
		payment.CreatedAt, payment.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("create payment: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}
	payment.ID = id
	return nil
}

// GetByID retrieves a payment by ID.
func (r *PaymentRepository) GetByID(ctx context.Context, id int64) (*domain.Payment, error) {
	query := fmt.Sprintf("SELECT %s FROM payments WHERE id = ?", paymentColumns)
	payment, err := r.scanPayment(r.db.QueryRowContext(ctx, query, id))
	if err == sql.ErrNoRows {
		return nil, domain.ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	return payment, nil
}

// GetByYooKassaID retrieves a payment by YooKassa ID.
func (r *PaymentRepository) GetByYooKassaID(ctx context.Context, yookassaID string) (*domain.Payment, error) {
	query := fmt.Sprintf("SELECT %s FROM payments WHERE yookassa_id = ?", paymentColumns)
	payment, err := r.scanPayment(r.db.QueryRowContext(ctx, query, yookassaID))
	if err == sql.ErrNoRows {
		return nil, domain.ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	return payment, nil
}

// List retrieves all payments.
func (r *PaymentRepository) List(ctx context.Context) ([]domain.Payment, int, error) {
	countQuery := "SELECT COUNT(*) FROM payments"
	var total int
	if err := r.db.QueryRowContext(ctx, countQuery).Scan(&total); err != nil {
		return nil, 0, err
	}

	rows, err := r.db.QueryContext(ctx, fmt.Sprintf("SELECT %s FROM payments ORDER BY created_at DESC", paymentColumns))
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var payments []domain.Payment
	for rows.Next() {
		p, err := r.scanPayment(rows)
		if err != nil {
			return nil, 0, err
		}
		payments = append(payments, *p)
	}
	return payments, total, nil
}

// Update updates a payment.
func (r *PaymentRepository) Update(ctx context.Context, payment *domain.Payment) error {
	query := `UPDATE payments SET user_id = ?, product_id = ?, amount = ?, currency = ?, status = ?,
		payment_method = ?, yookassa_id = ?, promo_code = ?, discount = ?, metadata = ?, updated_at = ? WHERE id = ?`

	payment.UpdatedAt = time.Now()

	_, err := r.db.ExecContext(ctx, query,
		payment.UserID, payment.ProductID, payment.Amount, payment.Currency, payment.Status,
		nullString(payment.PaymentMethod), nullString(payment.YooKassaID),
		nullString(payment.PromoCode), nullFloat(payment.Discount), nullString(payment.Metadata),
		payment.UpdatedAt, payment.ID,
	)
	return err
}

// PurchaseRepository implements port.PurchaseRepository.
type PurchaseRepository struct {
	db *sql.DB
}

// NewPurchaseRepository creates a new PurchaseRepository.
func NewPurchaseRepository(db *sql.DB) *PurchaseRepository {
	return &PurchaseRepository{db: db}
}

func (r *PurchaseRepository) Create(ctx context.Context, purchase *domain.Purchase) error {
	query := `INSERT INTO purchases (user_id, product_id, payment_id, price_paid, status, created_at) VALUES (?, ?, ?, ?, ?, ?)`
	now := time.Now()
	if purchase.CreatedAt.IsZero() {
		purchase.CreatedAt = now
	}
	if purchase.Status == "" {
		purchase.Status = "active"
	}

	result, err := r.db.ExecContext(ctx, query,
		purchase.UserID, purchase.ProductID, purchase.PaymentID, purchase.PricePaid, purchase.Status, purchase.CreatedAt,
	)
	if err != nil {
		return err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return err
	}
	purchase.ID = id
	return nil
}

func (r *PurchaseRepository) GetByID(ctx context.Context, id int64) (*domain.Purchase, error) {
	p := &domain.Purchase{}
	err := r.db.QueryRowContext(ctx,
		"SELECT id, user_id, product_id, payment_id, price_paid, status, created_at FROM purchases WHERE id = ?", id,
	).Scan(&p.ID, &p.UserID, &p.ProductID, &p.PaymentID, &p.PricePaid, &p.Status, &p.CreatedAt)
	if err == sql.ErrNoRows {
		return nil, domain.ErrNotFound
	}
	return p, err
}

func (r *PurchaseRepository) List(ctx context.Context) ([]domain.Purchase, int, error) {
	var total int
	r.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM purchases").Scan(&total)

	rows, err := r.db.QueryContext(ctx, "SELECT id, user_id, product_id, payment_id, price_paid, status, created_at FROM purchases ORDER BY created_at DESC")
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var purchases []domain.Purchase
	for rows.Next() {
		var p domain.Purchase
		if err := rows.Scan(&p.ID, &p.UserID, &p.ProductID, &p.PaymentID, &p.PricePaid, &p.Status, &p.CreatedAt); err != nil {
			return nil, 0, err
		}
		purchases = append(purchases, p)
	}
	return purchases, total, nil
}

func (r *PurchaseRepository) ListByUser(ctx context.Context, userID int64) ([]domain.Purchase, error) {
	rows, err := r.db.QueryContext(ctx,
		"SELECT id, user_id, product_id, payment_id, price_paid, status, created_at FROM purchases WHERE user_id = ? ORDER BY created_at DESC", userID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var purchases []domain.Purchase
	for rows.Next() {
		var p domain.Purchase
		if err := rows.Scan(&p.ID, &p.UserID, &p.ProductID, &p.PaymentID, &p.PricePaid, &p.Status, &p.CreatedAt); err != nil {
			return nil, err
		}
		purchases = append(purchases, p)
	}
	return purchases, nil
}

func (r *PurchaseRepository) GetByUserAndProduct(ctx context.Context, userID, productID int64) (*domain.Purchase, error) {
	p := &domain.Purchase{}
	err := r.db.QueryRowContext(ctx,
		"SELECT id, user_id, product_id, payment_id, price_paid, status, created_at FROM purchases WHERE user_id = ? AND product_id = ?",
		userID, productID,
	).Scan(&p.ID, &p.UserID, &p.ProductID, &p.PaymentID, &p.PricePaid, &p.Status, &p.CreatedAt)
	if err == sql.ErrNoRows {
		return nil, domain.ErrNotFound
	}
	return p, err
}

// PromoCodeRepository implements port.PromoCodeRepository.
type PromoCodeRepository struct {
	db *sql.DB
}

// NewPromoCodeRepository creates a new PromoCodeRepository.
func NewPromoCodeRepository(db *sql.DB) *PromoCodeRepository {
	return &PromoCodeRepository{db: db}
}

func (r *PromoCodeRepository) Create(ctx context.Context, code *domain.PromoCode) error {
	query := `INSERT INTO promo_codes (code, discount_percent, max_uses, current_uses, expires_at, is_active, created_at) VALUES (?, ?, ?, ?, ?, ?, ?)`
	now := time.Now()
	if code.CreatedAt.IsZero() {
		code.CreatedAt = now
	}

	result, err := r.db.ExecContext(ctx, query,
		code.Code, code.DiscountPercent, code.MaxUses, code.CurrentUses, code.ExpiresAt, code.IsActive, code.CreatedAt,
	)
	if err != nil {
		return err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return err
	}
	code.ID = id
	return nil
}

func (r *PromoCodeRepository) GetByCode(ctx context.Context, code string) (*domain.PromoCode, error) {
	p := &domain.PromoCode{}
	err := r.db.QueryRowContext(ctx,
		"SELECT id, code, discount_percent, max_uses, current_uses, expires_at, is_active, created_at FROM promo_codes WHERE code = ?", code,
	).Scan(&p.ID, &p.Code, &p.DiscountPercent, &p.MaxUses, &p.CurrentUses, &p.ExpiresAt, &p.IsActive, &p.CreatedAt)
	if err == sql.ErrNoRows {
		return nil, domain.ErrNotFound
	}
	return p, err
}

func (r *PromoCodeRepository) List(ctx context.Context) ([]domain.PromoCode, int, error) {
	var total int
	r.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM promo_codes").Scan(&total)

	rows, err := r.db.QueryContext(ctx, "SELECT id, code, discount_percent, max_uses, current_uses, expires_at, is_active, created_at FROM promo_codes ORDER BY created_at DESC")
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var codes []domain.PromoCode
	for rows.Next() {
		var p domain.PromoCode
		if err := rows.Scan(&p.ID, &p.Code, &p.DiscountPercent, &p.MaxUses, &p.CurrentUses, &p.ExpiresAt, &p.IsActive, &p.CreatedAt); err != nil {
			return nil, 0, err
		}
		codes = append(codes, p)
	}
	return codes, total, nil
}

func (r *PromoCodeRepository) Update(ctx context.Context, code *domain.PromoCode) error {
	_, err := r.db.ExecContext(ctx,
		"UPDATE promo_codes SET code = ?, discount_percent = ?, max_uses = ?, current_uses = ?, expires_at = ?, is_active = ? WHERE id = ?",
		code.Code, code.DiscountPercent, code.MaxUses, code.CurrentUses, code.ExpiresAt, code.IsActive, code.ID,
	)
	return err
}

func (r *PromoCodeRepository) Delete(ctx context.Context, id int64) error {
	result, err := r.db.ExecContext(ctx, "DELETE FROM promo_codes WHERE id = ?", id)
	if err != nil {
		return err
	}
	affected, _ := result.RowsAffected()
	if affected == 0 {
		return domain.ErrNotFound
	}
	return nil
}

// ReviewRepository implements port.ReviewRepository.
type ReviewRepository struct {
	db *sql.DB
}

// NewReviewRepository creates a new ReviewRepository.
func NewReviewRepository(db *sql.DB) *ReviewRepository {
	return &ReviewRepository{db: db}
}

func (r *ReviewRepository) Create(ctx context.Context, review *domain.Review) error {
	query := `INSERT INTO reviews (user_id, product_id, rating, text, status, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, ?)`
	now := time.Now()
	if review.CreatedAt.IsZero() {
		review.CreatedAt = now
	}
	if review.UpdatedAt.IsZero() {
		review.UpdatedAt = now
	}
	if review.Status == "" {
		review.Status = "pending"
	}

	result, err := r.db.ExecContext(ctx, query,
		review.UserID, review.ProductID, review.Rating, nullString(review.Text), review.Status, review.CreatedAt, review.UpdatedAt,
	)
	if err != nil {
		return err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return err
	}
	review.ID = id
	return nil
}

func (r *ReviewRepository) GetByID(ctx context.Context, id int64) (*domain.Review, error) {
	rev := &domain.Review{}
	var text sql.NullString
	err := r.db.QueryRowContext(ctx,
		"SELECT id, user_id, product_id, rating, text, status, created_at, updated_at FROM reviews WHERE id = ?", id,
	).Scan(&rev.ID, &rev.UserID, &rev.ProductID, &rev.Rating, &text, &rev.Status, &rev.CreatedAt, &rev.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, domain.ErrNotFound
	}
	if text.Valid {
		rev.Text = text.String
	}
	return rev, err
}

func (r *ReviewRepository) ListByProduct(ctx context.Context, productID int64) ([]domain.Review, error) {
	rows, err := r.db.QueryContext(ctx,
		"SELECT id, user_id, product_id, rating, text, status, created_at, updated_at FROM reviews WHERE product_id = ? ORDER BY created_at DESC", productID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var reviews []domain.Review
	for rows.Next() {
		var rev domain.Review
		var text sql.NullString
		if err := rows.Scan(&rev.ID, &rev.UserID, &rev.ProductID, &rev.Rating, &text, &rev.Status, &rev.CreatedAt, &rev.UpdatedAt); err != nil {
			return nil, err
		}
		if text.Valid {
			rev.Text = text.String
		}
		reviews = append(reviews, rev)
	}
	return reviews, nil
}

func (r *ReviewRepository) List(ctx context.Context) ([]domain.Review, int, error) {
	var total int
	r.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM reviews").Scan(&total)

	rows, err := r.db.QueryContext(ctx, "SELECT id, user_id, product_id, rating, text, status, created_at, updated_at FROM reviews ORDER BY created_at DESC")
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var reviews []domain.Review
	for rows.Next() {
		var rev domain.Review
		var text sql.NullString
		if err := rows.Scan(&rev.ID, &rev.UserID, &rev.ProductID, &rev.Rating, &text, &rev.Status, &rev.CreatedAt, &rev.UpdatedAt); err != nil {
			return nil, 0, err
		}
		if text.Valid {
			rev.Text = text.String
		}
		reviews = append(reviews, rev)
	}
	return reviews, total, nil
}

func (r *ReviewRepository) Update(ctx context.Context, review *domain.Review) error {
	review.UpdatedAt = time.Now()
	_, err := r.db.ExecContext(ctx,
		"UPDATE reviews SET rating = ?, text = ?, status = ?, updated_at = ? WHERE id = ?",
		review.Rating, nullString(review.Text), review.Status, review.UpdatedAt, review.ID,
	)
	return err
}

func (r *ReviewRepository) Delete(ctx context.Context, id int64) error {
	result, err := r.db.ExecContext(ctx, "DELETE FROM reviews WHERE id = ?", id)
	if err != nil {
		return err
	}
	affected, _ := result.RowsAffected()
	if affected == 0 {
		return domain.ErrNotFound
	}
	return nil
}
