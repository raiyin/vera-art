package service

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"math"

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
		slog.Error("PaymentService.CreatePayment: failed to get product",
			"user_id", userID,
			"product_id", productID,
			"error", err,
		)
		return nil, "", err
	}

	amount := float64(product.Price) / 100
	var discount float64

	if promoCode != "" {
		promo, err := s.promoService.ValidatePromoCode(ctx, promoCode)
		if err != nil {
			slog.Warn("PaymentService.CreatePayment: invalid promo code",
				"user_id", userID,
				"product_id", productID,
				"promo_code", promoCode,
				"error", err,
			)
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
		slog.Error("PaymentService.CreatePayment: failed to create payment",
			"user_id", userID,
			"product_id", productID,
			"amount", amount,
			"error", err,
		)
		return nil, "", err
	}

	// Stub: return a mock confirmation URL
	confirmationURL := fmt.Sprintf("https://yookassa.ru/payment?payment_id=%d", payment.ID)

	slog.Info("PaymentService.CreatePayment: payment created",
		"user_id", userID,
		"product_id", productID,
		"payment_id", payment.ID,
		"amount", amount,
		"discount", discount,
	)

	return payment, confirmationURL, nil
}

// GetPaymentStatus retrieves the status of a payment.
func (s *PaymentService) GetPaymentStatus(ctx context.Context, paymentID int64) (*domain.Payment, error) {
	payment, err := s.paymentRepo.GetByID(ctx, paymentID)
	if err != nil {
		slog.Error("PaymentService.GetPaymentStatus: failed to get payment",
			"payment_id", paymentID,
			"error", err,
		)
		return nil, err
	}
	slog.Debug("PaymentService.GetPaymentStatus: payment retrieved",
		"payment_id", paymentID,
		"status", payment.Status,
	)
	return payment, nil
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
		slog.Error("PaymentService.HandleWebhook: failed to unmarshal payload",
			"error", err,
		)
		return domain.ErrInvalidInput
	}

	// Find payment by metadata or yookassa_id
	payment, err := s.paymentRepo.GetByYooKassaID(ctx, webhook.Object.ID)
	if err != nil {
		slog.Error("PaymentService.HandleWebhook: failed to find payment by yookassa_id",
			"yookassa_id", webhook.Object.ID,
			"event", webhook.Event,
			"error", err,
		)
		return err
	}

	slog.Info("PaymentService.HandleWebhook: processing webhook",
		"payment_id", payment.ID,
		"yookassa_id", webhook.Object.ID,
		"event", webhook.Event,
		"current_status", payment.Status,
	)

	switch webhook.Event {
	case "payment.waiting_for_capture":
		payment.Status = "waiting_for_capture"
		slog.Info("PaymentService.HandleWebhook: payment waiting for capture",
			"payment_id", payment.ID,
		)
	case "payment.succeeded":
		payment.Status = "succeeded"
		payment.YooKassaID = webhook.Object.ID

		// Create purchase record
		purchase := &domain.Purchase{
			UserID:    payment.UserID,
			ProductID: payment.ProductID,
			PaymentID: payment.ID,
			PricePaid: int64(math.Round(payment.Amount * 100)),
			Status:    "active",
		}
		if err := s.purchaseRepo.Create(ctx, purchase); err != nil {
			slog.Error("PaymentService.HandleWebhook: failed to create purchase",
				"payment_id", payment.ID,
				"user_id", payment.UserID,
				"product_id", payment.ProductID,
				"error", err,
			)
			return err
		}
		slog.Info("PaymentService.HandleWebhook: payment succeeded, purchase created",
			"payment_id", payment.ID,
			"user_id", payment.UserID,
			"product_id", payment.ProductID,
			"amount", payment.Amount,
		)
	case "payment.canceled":
		payment.Status = "canceled"
		slog.Info("PaymentService.HandleWebhook: payment canceled",
			"payment_id", payment.ID,
		)
	default:
		slog.Warn("PaymentService.HandleWebhook: unknown event type",
			"payment_id", payment.ID,
			"event", webhook.Event,
		)
	}

	if err := s.paymentRepo.Update(ctx, payment); err != nil {
		slog.Error("PaymentService.HandleWebhook: failed to update payment",
			"payment_id", payment.ID,
			"new_status", payment.Status,
			"error", err,
		)
		return err
	}

	slog.Info("PaymentService.HandleWebhook: payment updated",
		"payment_id", payment.ID,
		"new_status", payment.Status,
	)
	return nil
}

// GetPayments retrieves all payments.
func (s *PaymentService) GetPayments(ctx context.Context) ([]domain.Payment, int, error) {
	payments, total, err := s.paymentRepo.List(ctx)
	if err != nil {
		slog.Error("PaymentService.GetPayments: failed to list payments",
			"error", err,
		)
		return nil, 0, err
	}
	slog.Debug("PaymentService.GetPayments: payments listed",
		"count", len(payments),
		"total", total,
	)
	return payments, total, nil
}

// GetPurchases retrieves all purchases.
func (s *PaymentService) GetPurchases(ctx context.Context) ([]domain.Purchase, int, error) {
	purchases, total, err := s.purchaseRepo.List(ctx)
	if err != nil {
		slog.Error("PaymentService.GetPurchases: failed to list purchases",
			"error", err,
		)
		return nil, 0, err
	}
	slog.Debug("PaymentService.GetPurchases: purchases listed",
		"count", len(purchases),
		"total", total,
	)
	return purchases, total, nil
}

// GetUserPurchases retrieves purchases for a user.
func (s *PaymentService) GetUserPurchases(ctx context.Context, userID int64) ([]domain.Purchase, error) {
	purchases, err := s.purchaseRepo.ListByUser(ctx, userID)
	if err != nil {
		slog.Error("PaymentService.GetUserPurchases: failed to list user purchases",
			"user_id", userID,
			"error", err,
		)
		return nil, err
	}
	slog.Debug("PaymentService.GetUserPurchases: user purchases listed",
		"user_id", userID,
		"count", len(purchases),
	)
	return purchases, nil
}

// HasUserPurchasedProduct checks if a user has purchased a product.
func (s *PaymentService) HasUserPurchasedProduct(ctx context.Context, userID, productID int64) (bool, error) {
	purchase, err := s.purchaseRepo.GetByUserAndProduct(ctx, userID, productID)
	if err == domain.ErrNotFound {
		slog.Debug("PaymentService.HasUserPurchasedProduct: no purchase found",
			"user_id", userID,
			"product_id", productID,
		)
		return false, nil
	}
	if err != nil {
		slog.Error("PaymentService.HasUserPurchasedProduct: failed to check purchase",
			"user_id", userID,
			"product_id", productID,
			"error", err,
		)
		return false, err
	}
	purchased := purchase != nil
	slog.Debug("PaymentService.HasUserPurchasedProduct: checked",
		"user_id", userID,
		"product_id", productID,
		"purchased", purchased,
	)
	return purchased, nil
}

func (s *PaymentService) GetPaymentByID(ctx context.Context, id int64) (*domain.Payment, error) {
	payment, err := s.paymentRepo.GetByID(ctx, id)
	if err != nil {
		slog.Error("PaymentService.GetPaymentByID: failed to get payment",
			"payment_id", id,
			"error", err,
		)
		return nil, err
	}
	return payment, nil
}

func (s *PaymentService) GetPurchaseByID(ctx context.Context, id int64) (*domain.Purchase, error) {
	purchase, err := s.purchaseRepo.GetByID(ctx, id)
	if err != nil {
		slog.Error("PaymentService.GetPurchaseByID: failed to get purchase",
			"purchase_id", id,
			"error", err,
		)
		return nil, err
	}
	return purchase, nil
}

func (s *PaymentService) ExtendPurchaseAccess(ctx context.Context, id int64, days int) error {
	purchase, err := s.purchaseRepo.GetByID(ctx, id)
	if err != nil {
		slog.Error("PaymentService.ExtendPurchaseAccess: failed to get purchase",
			"purchase_id", id,
			"error", err,
		)
		return err
	}
	if purchase.AccessEnd != nil {
		newEnd := purchase.AccessEnd.AddDate(0, 0, days)
		purchase.AccessEnd = &newEnd
	}
	slog.Info("PaymentService.ExtendPurchaseAccess: access extended",
		"purchase_id", id,
		"days", days,
	)
	return nil
}

func (s *PaymentService) CancelPurchase(ctx context.Context, id int64) error {
	purchase, err := s.purchaseRepo.GetByID(ctx, id)
	if err != nil {
		slog.Error("PaymentService.CancelPurchase: failed to get purchase",
			"purchase_id", id,
			"error", err,
		)
		return err
	}
	purchase.Status = "cancelled"
	slog.Info("PaymentService.CancelPurchase: purchase cancelled",
		"purchase_id", id,
	)
	return nil
}

func (s *PaymentService) RefundPayment(ctx context.Context, id int64) error {
	payment, err := s.paymentRepo.GetByID(ctx, id)
	if err != nil {
		slog.Error("PaymentService.RefundPayment: failed to get payment",
			"payment_id", id,
			"error", err,
		)
		return err
	}
	payment.Status = "refunded"
	if err := s.paymentRepo.Update(ctx, payment); err != nil {
		slog.Error("PaymentService.RefundPayment: failed to update payment",
			"payment_id", id,
			"error", err,
		)
		return err
	}
	slog.Info("PaymentService.RefundPayment: payment refunded",
		"payment_id", id,
	)
	return nil
}

// Ensure interface compliance.
var _ port.PaymentService = (*PaymentService)(nil)
