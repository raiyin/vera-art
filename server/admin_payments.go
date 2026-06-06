package main

import (
	"math"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

// AdminPaymentItem represents a payment in the admin list
type AdminPaymentItem struct {
	ID            int    `json:"id"`
	UserID        int    `json:"user_id"`
	ExternalID    string `json:"external_id,omitempty"`
	Status        string `json:"status"`
	Amount        int    `json:"amount"`
	Currency      string `json:"currency"`
	Description   string `json:"description,omitempty"`
	PaymentMethod string `json:"payment_method,omitempty"`
	CreatedAt     string `json:"created_at"`
	UpdatedAt     string `json:"updated_at"`
	Username      string `json:"username"`
	UserFullName  string `json:"user_full_name"`
	UserEmail     string `json:"user_email"`
}

// AdminListPaymentsResponse is the response for the payments list
type AdminListPaymentsResponse struct {
	Items      []AdminPaymentItem `json:"items"`
	Total      int                `json:"total"`
	Page       int                `json:"page"`
	PerPage    int                `json:"per_page"`
	TotalPages int                `json:"total_pages"`
}

// AdminGetPaymentsList returns a paginated list of payments
func AdminGetPaymentsList(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	perPage, _ := strconv.Atoi(c.DefaultQuery("per_page", "20"))
	search := c.Query("search")
	status := c.Query("status")
	sortBy := c.DefaultQuery("sort_by", "created_at")
	sortDir := c.DefaultQuery("sort_dir", "desc")

	if page < 1 {
		page = 1
	}
	if perPage < 1 || perPage > 100 {
		perPage = 20
	}

	allowedSortFields := map[string]bool{
		"created_at": true,
		"amount":     true,
		"status":     true,
	}
	if !allowedSortFields[sortBy] {
		sortBy = "created_at"
	}
	if sortDir != "asc" && sortDir != "desc" {
		sortDir = "desc"
	}

	where := "1=1"
	args := []any{}

	if search != "" {
		where += " AND (u.username LIKE ? OR u.full_name LIKE ? OR u.email LIKE ? OR pay.description LIKE ?)"
		s := "%" + search + "%"
		args = append(args, s, s, s, s)
	}
	if status != "" {
		where += " AND pay.status = ?"
		args = append(args, status)
	}

	// Count total
	countQuery := `SELECT COUNT(*) FROM payments pay
		JOIN users u ON pay.user_id = u.id
		WHERE ` + where

	var total int
	err := db.QueryRow(countQuery, args...).Scan(&total)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка подсчёта платежей"})
		return
	}

	totalPages := int(math.Ceil(float64(total) / float64(perPage)))
	offset := (page - 1) * perPage

	// Fetch items
	query := `SELECT pay.id, pay.user_id, pay.external_id, pay.status, pay.amount, pay.currency,
		pay.description, pay.payment_method, pay.created_at, pay.updated_at,
		u.username, u.full_name, u.email
		FROM payments pay
		JOIN users u ON pay.user_id = u.id
		WHERE ` + where + ` ORDER BY pay.` + sortBy + ` ` + sortDir + ` LIMIT ? OFFSET ?`

	args = append(args, perPage, offset)
	rows, err := db.Query(query, args...)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка загрузки платежей"})
		return
	}
	defer rows.Close()

	items := []AdminPaymentItem{}
	for rows.Next() {
		var item AdminPaymentItem
		var externalID, description, paymentMethod interface{}
		var createdAt, updatedAt time.Time

		err := rows.Scan(
			&item.ID, &item.UserID, &externalID, &item.Status, &item.Amount, &item.Currency,
			&description, &paymentMethod, &createdAt, &updatedAt,
			&item.Username, &item.UserFullName, &item.UserEmail,
		)
		if err != nil {
			continue
		}

		if externalID != nil {
			item.ExternalID = string(externalID.([]byte))
		}
		if description != nil {
			item.Description = string(description.([]byte))
		}
		if paymentMethod != nil {
			item.PaymentMethod = string(paymentMethod.([]byte))
		}

		item.CreatedAt = createdAt.Format("2006-01-02 15:04:05")
		item.UpdatedAt = updatedAt.Format("2006-01-02 15:04:05")

		items = append(items, item)
	}

	c.JSON(http.StatusOK, AdminListPaymentsResponse{
		Items:      items,
		Total:      total,
		Page:       page,
		PerPage:    perPage,
		TotalPages: totalPages,
	})
}

// AdminGetPaymentDetail returns detailed info about a payment
func AdminGetPaymentDetail(c *gin.Context) {
	id := c.Param("id")

	var item AdminPaymentItem
	var externalID, description, paymentMethod interface{}
	var createdAt, updatedAt time.Time

	err := db.QueryRow(`
		SELECT pay.id, pay.user_id, pay.external_id, pay.status, pay.amount, pay.currency,
			pay.description, pay.payment_method, pay.created_at, pay.updated_at,
			u.username, u.full_name, u.email
		FROM payments pay
		JOIN users u ON pay.user_id = u.id
		WHERE pay.id = ?
	`, id).Scan(
		&item.ID, &item.UserID, &externalID, &item.Status, &item.Amount, &item.Currency,
		&description, &paymentMethod, &createdAt, &updatedAt,
		&item.Username, &item.UserFullName, &item.UserEmail,
	)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Платеж не найден"})
		return
	}

	if externalID != nil {
		item.ExternalID = string(externalID.([]byte))
	}
	if description != nil {
		item.Description = string(description.([]byte))
	}
	if paymentMethod != nil {
		item.PaymentMethod = string(paymentMethod.([]byte))
	}

	item.CreatedAt = createdAt.Format("2006-01-02 15:04:05")
	item.UpdatedAt = updatedAt.Format("2006-01-02 15:04:05")

	c.JSON(http.StatusOK, item)
}

// AdminRefundPayment initiates a refund for a payment
func AdminRefundPayment(c *gin.Context) {
	id := c.Param("id")

	// Check payment exists and is in refundable status
	var status string
	err := db.QueryRow("SELECT status FROM payments WHERE id = ?", id).Scan(&status)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Платеж не найден"})
		return
	}

	if status != "succeeded" && status != "waiting_for_capture" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Платеж не может быть возвращен (текущий статус: " + status + ")"})
		return
	}

	_, err = db.Exec("UPDATE payments SET status = 'refunded', updated_at = ? WHERE id = ?", time.Now(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка возврата платежа"})
		return
	}

	// Also cancel associated purchases
	_, _ = db.Exec("UPDATE purchases SET status = 'cancelled' WHERE payment_id = ?", id)

	c.JSON(http.StatusOK, gin.H{"status": "success", "message": "Платеж возвращен"})
}
