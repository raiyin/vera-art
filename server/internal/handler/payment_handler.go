package handler

import (
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/raiyin/artserver/internal/dto"
	"github.com/raiyin/artserver/internal/port"
	"github.com/raiyin/artserver/pkg/apperror"
)

// PaymentHandler handles payment HTTP requests.
type PaymentHandler struct {
	paymentService port.PaymentService
}

// NewPaymentHandler creates a new PaymentHandler.
func NewPaymentHandler(paymentService port.PaymentService) *PaymentHandler {
	return &PaymentHandler{paymentService: paymentService}
}

// CreatePayment creates a new payment.
func (h *PaymentHandler) CreatePayment(c *gin.Context) {
	userID := c.GetInt64("user_id")

	var req dto.CreatePaymentRequest
	if !BindJSON(c, &req) {
		return
	}

	payment, confirmationURL, err := h.paymentService.CreatePayment(c.Request.Context(), userID, req.ProductID, req.PromoCode)
	if err != nil {
		slog.Error("CreatePayment: failed to create payment",
			"user_id", userID,
			"product_id", req.ProductID,
			"error", err,
		)
		apiErr := apperror.FromError(err)
		c.JSON(apiErr.Status, apiErr)
		return
	}

	slog.Info("Payment created",
		"payment_id", payment.ID,
		"user_id", userID,
		"amount", payment.Amount,
	)
	c.JSON(http.StatusOK, dto.CreatePaymentResponse{
		PaymentID:       payment.ID,
		Amount:          payment.Amount,
		Currency:        payment.Currency,
		ConfirmationURL: confirmationURL,
		Description:     "Payment for product",
		Status:          payment.Status,
	})
}

// GetPaymentStatus returns the status of a payment.
func (h *PaymentHandler) GetPaymentStatus(c *gin.Context) {
	paymentID, ok := ParseInt64Param(c, "id")
	if !ok {
		return
	}

	payment, err := h.paymentService.GetPaymentStatus(c.Request.Context(), paymentID)
	if err != nil {
		slog.Error("GetPaymentStatus: failed to get payment status",
			"payment_id", paymentID,
			"error", err,
		)
		apiErr := apperror.FromError(err)
		c.JSON(apiErr.Status, apiErr)
		return
	}

	c.JSON(http.StatusOK, dto.PaymentResponse{
		ID:            payment.ID,
		UserID:        payment.UserID,
		Status:        payment.Status,
		Amount:        payment.Amount,
		Currency:      payment.Currency,
		Description:   "Payment",
		PaymentMethod: payment.PaymentMethod,
		CreatedAt:     payment.CreatedAt,
		UpdatedAt:     payment.UpdatedAt,
	})
}

// HandleWebhook handles YooKassa webhook.
func (h *PaymentHandler) HandleWebhook(c *gin.Context) {
	payload, err := c.GetRawData()
	if err != nil {
		slog.Error("HandleWebhook: failed to read webhook payload",
			"error", err,
		)
		c.JSON(http.StatusBadRequest, apperror.APIError{
			Status:  http.StatusBadRequest,
			Code:    "INVALID_REQUEST",
			Message: "Invalid webhook payload",
		})
		return
	}

	if err := h.paymentService.HandleWebhook(c.Request.Context(), payload); err != nil {
		slog.Error("HandleWebhook: webhook processing failed",
			"error", err,
		)
		apiErr := apperror.FromError(err)
		c.JSON(apiErr.Status, apiErr)
		return
	}

	slog.Info("Webhook processed successfully")
	c.Status(http.StatusOK)
}

// GetPayments returns all payments (admin).
func (h *PaymentHandler) GetPayments(c *gin.Context) {
	payments, total, err := h.paymentService.GetPayments(c.Request.Context())
	if err != nil {
		slog.Error("GetPayments: failed to list payments",
			"error", err,
		)
		apiErr := apperror.FromError(err)
		c.JSON(apiErr.Status, apiErr)
		return
	}

	responses := make([]dto.PaymentResponse, len(payments))
	for i, p := range payments {
		responses[i] = dto.PaymentResponse{
			ID:            p.ID,
			UserID:        p.UserID,
			Status:        p.Status,
			Amount:        p.Amount,
			Currency:      p.Currency,
			Description:   "Payment",
			PaymentMethod: p.PaymentMethod,
			CreatedAt:     p.CreatedAt,
			UpdatedAt:     p.UpdatedAt,
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"payments": responses,
		"total":    total,
	})
}

// GetPurchases returns all purchases (admin).
func (h *PaymentHandler) GetPurchases(c *gin.Context) {
	purchases, total, err := h.paymentService.GetPurchases(c.Request.Context())
	if err != nil {
		slog.Error("GetPurchases: failed to list purchases",
			"error", err,
		)
		apiErr := apperror.FromError(err)
		c.JSON(apiErr.Status, apiErr)
		return
	}

	responses := make([]dto.PurchaseResponse, len(purchases))
	for i, p := range purchases {
		responses[i] = dto.PurchaseResponse{
			ID:          p.ID,
			UserID:      p.UserID,
			ProductID:   p.ProductID,
			PaymentID:   p.PaymentID,
			PricePaid:   p.PricePaid,
			Status:      p.Status,
			AccessStart: p.AccessStart,
			AccessEnd:   p.AccessEnd,
			CreatedAt:   p.CreatedAt,
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"purchases": responses,
		"total":     total,
	})
}

// GetUserPurchases returns purchases for the authenticated user.
func (h *PaymentHandler) GetUserPurchases(c *gin.Context) {
	userID := c.GetInt64("user_id")

	purchases, err := h.paymentService.GetUserPurchases(c.Request.Context(), userID)
	if err != nil {
		slog.Error("GetUserPurchases: failed to list user purchases",
			"user_id", userID,
			"error", err,
		)
		apiErr := apperror.FromError(err)
		c.JSON(apiErr.Status, apiErr)
		return
	}

	responses := make([]dto.PurchaseResponse, len(purchases))
	for i, p := range purchases {
		responses[i] = dto.PurchaseResponse{
			ID:          p.ID,
			UserID:      p.UserID,
			ProductID:   p.ProductID,
			PaymentID:   p.PaymentID,
			PricePaid:   p.PricePaid,
			Status:      p.Status,
			AccessStart: p.AccessStart,
			AccessEnd:   p.AccessEnd,
			CreatedAt:   p.CreatedAt,
		}
	}

	c.JSON(http.StatusOK, gin.H{"purchases": responses})
}

// HasUserPurchasedProduct checks if the user has purchased a product.
func (h *PaymentHandler) HasUserPurchasedProduct(c *gin.Context) {
	userID := c.GetInt64("user_id")
	productID, ok := ParseInt64Param(c, "product_id")
	if !ok {
		return
	}

	hasPurchased, err := h.paymentService.HasUserPurchasedProduct(c.Request.Context(), userID, productID)
	if err != nil {
		slog.Error("HasUserPurchasedProduct: failed to check purchase",
			"user_id", userID,
			"product_id", productID,
			"error", err,
		)
		apiErr := apperror.FromError(err)
		c.JSON(apiErr.Status, apiErr)
		return
	}

	c.JSON(http.StatusOK, gin.H{"has_purchased": hasPurchased})
}

// CreatePurchase creates a purchase after successful payment.
func (h *PaymentHandler) CreatePurchase(c *gin.Context) {
	userID := c.GetInt64("user_id")

	var req dto.CreatePurchaseRequest
	if !BindJSON(c, &req) {
		return
	}

	// This is a simplified version - in the real implementation,
	// the purchase is created via webhook. This endpoint is kept
	// for backward compatibility.
	_ = userID
	_ = req

	c.JSON(http.StatusOK, gin.H{"message": "Purchase creation is handled via webhook"})
}
