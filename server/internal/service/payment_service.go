package service

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/raiyin/artserver/internal/domain"
	"github.com/raiyin/artserver/internal/port"
)

// PaymentService implements port.PaymentService.
type PaymentService struct {
	paymentRepo  port.PaymentRepository
	purchaseRepo port.PurchaseRepository
	productRepo  port.ProductRepository
	promoService port.ShopService
}

// NewPaymentService creates a new PaymentService.
func NewPaymentService(
	paymentRepo port.PaymentRepository,
	purchaseRepo port.PurchaseRepository,
	productRepo port.ProductRepository,
	promoService port.ShopService,
) *PaymentService {
	return &PaymentService{
		paymentRepo:  paymentRepo,
		purchaseRepo: purchaseRepo,
		productRepo:  productRepo,
		promoService: promoService,
	}
}

// CreatePayment creates a new payment (YooKassa stub).
func (s *PaymentService) CreatePayment(ctx context.Context, userID, productID int64, promoCode string) (*domain.Payment, string, error) {
	product, err := s.productRepo.GetByID(ctx, productID)
	if err != nil {
		return nil, "", err
	}

	amount := product.Price
	var discount float64

	if promoCode != "" {
		promo, err := s.promoService.ValidatePromoCode(ctx, promoCode)
		if err != nil {
			return nil, "", err
		}
		discount = amount * promo.DiscountPercent / 100
		amount -= discount
	}

	payment := &domain.Payment{
		UserID:    userID,
		ProductID: productID,
		Amount:    amount,
		Currency:  "RUB",
		Status:    "pending",
		PromoCode: promoCode,
		Discount:  discount,
	}

	if err := s.paymentRepo.Create(ctx, payment); err != nil {
		return nil, "", err
	}

	// Stub: return a mock confirmation URL
	confirmationURL := fmt.Sprintf("https://yookassa.ru/payment?payment_id=%d", payment.ID)

	return payment, confirmationURL, nil
}

// GetPaymentStatus retrieves the status of a payment.
func (s *PaymentService) GetPaymentStatus(ctx context.Context, paymentID int64) (*domain.Payment, error) {
	return s.paymentRepo.GetByID(ctx, paymentID)
}

// HandleWebhook handles a YooKassa webhook (stub).
func (s *PaymentService) HandleWebhook(ctx context.Context, payload []byte) error {
	var webhook struct {
		Event  string `json:"event"`
		Object struct {
			ID       string                 `json:"id"`
			Status   string                 `json:"status"`
			Metadata map[string]interface{} `json:"metadata"`
		} `json:"object"`
	}

	if err := json.Unmarshal(payload, &webhook); err != nil {
		return domain.ErrInvalidInput
	}

	// Find payment by metadata or yookassa_id
	payment, err := s.paymentRepo.GetByYooKassaID(ctx, webhook.Object.ID)
	if err != nil {
		return err
	}

	switch webhook.Event {
	case "payment.waiting_for_capture":
		payment.Status = "waiting_for_capture"
	case "payment.succeeded":
		payment.Status = "succeeded"
		payment.YooKassaID = webhook.Object.ID

		// Create purchase record
		purchase := &domain.Purchase{
			UserID:    payment.UserID,
			ProductID: payment.ProductID,
			PaymentID: payment.ID,
			Amount:    payment.Amount,
			Status:    "completed",
		}
		if err := s.purchaseRepo.Create(ctx, purchase); err != nil {
			return err
		}
	case "payment.canceled":
		payment.Status = "canceled"
	}

	return s.paymentRepo.Update(ctx, payment)
}

// GetPayments retrieves all payments.
func (s *PaymentService) GetPayments(ctx context.Context) ([]domain.Payment, int, error) {
	return s.paymentRepo.List(ctx)
}

// GetPurchases retrieves all purchases.
func (s *PaymentService) GetPurchases(ctx context.Context) ([]domain.Purchase, int, error) {
	return s.purchaseRepo.List(ctx)
}

// GetUserPurchases retrieves purchases for a user.
func (s *PaymentService) GetUserPurchases(ctx context.Context, userID int64) ([]domain.Purchase, error) {
	return s.purchaseRepo.ListByUser(ctx, userID)
}

// HasUserPurchasedProduct checks if a user has purchased a product.
func (s *PaymentService) HasUserPurchasedProduct(ctx context.Context, userID, productID int64) (bool, error) {
	purchase, err := s.purchaseRepo.GetByUserAndProduct(ctx, userID, productID)
	if err == domain.ErrNotFound {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return purchase != nil, nil
}

// Ensure interface compliance.
var _ port.PaymentService = (*PaymentService)(nil)
