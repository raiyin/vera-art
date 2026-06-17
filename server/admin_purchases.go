package main

import (
	"database/sql"
	"math"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

// AdminPurchaseItem represents a purchase in the admin list
type AdminPurchaseItem struct {
	ID             int    `json:"id"`
	UserID         int    `json:"user_id"`
	ProductID      int    `json:"product_id"`
	PurchaseDate   string `json:"purchase_date"`
	AccessStart    string `json:"access_start"`
	AccessEnd      string `json:"access_end,omitempty"`
	Status         string `json:"status"`
	PaymentID      *int   `json:"payment_id,omitempty"`
	PromoCodeID    *int   `json:"promo_code_id,omitempty"`
	PricePaid      int    `json:"price_paid"`
	CreatedAt      string `json:"created_at"`
	Username       string `json:"username"`
	UserFullName   string `json:"user_full_name"`
	UserEmail      string `json:"user_email"`
	ProductTitleRu string `json:"product_title_ru"`
	ProductTitleEn string `json:"product_title_en"`
	ProductType    string `json:"product_type"`
	DaysRemaining  int    `json:"days_remaining"`
}

// AdminListPurchasesResponse is the response for the purchases list
type AdminListPurchasesResponse struct {
	Items      []AdminPurchaseItem `json:"items"`
	Total      int                 `json:"total"`
	Page       int                 `json:"page"`
	PerPage    int                 `json:"per_page"`
	TotalPages int                 `json:"total_pages"`
}

// AdminGetPurchasesList returns a paginated list of purchases
func AdminGetPurchasesList(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	perPage, _ := strconv.Atoi(c.DefaultQuery("per_page", "20"))
	search := c.Query("search")
	status := c.Query("status")
	productType := c.Query("product_type")
	sortBy := c.DefaultQuery("sort_by", "created_at")
	sortDir := c.DefaultQuery("sort_dir", "desc")

	if page < 1 {
		page = 1
	}
	if perPage < 1 || perPage > 100 {
		perPage = 20
	}

	allowedSortFields := map[string]bool{
		"created_at":    true,
		"purchase_date": true,
		"price_paid":    true,
		"status":        true,
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
		where += " AND (u.username LIKE ? OR u.full_name LIKE ? OR u.email LIKE ? OR p.title_ru LIKE ? OR p.title_en LIKE ?)"
		s := "%" + search + "%"
		args = append(args, s, s, s, s, s)
	}
	if status != "" {
		where += " AND pur.status = ?"
		args = append(args, status)
	}
	if productType != "" {
		where += " AND p.type = ?"
		args = append(args, productType)
	}

	// Count total
	countQuery := `SELECT COUNT(*) FROM purchases pur
		JOIN users u ON pur.user_id = u.id
		JOIN products p ON pur.product_id = p.id
		WHERE ` + where

	var total int
	err := db.QueryRow(countQuery, args...).Scan(&total)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка подсчёта покупок"})
		return
	}

	totalPages := int(math.Ceil(float64(total) / float64(perPage)))
	offset := (page - 1) * perPage

	// Fetch items
	query := `SELECT pur.id, pur.user_id, pur.product_id, pur.purchase_date, pur.access_start,
		pur.access_end, pur.status, pur.payment_id, pur.promo_code_id, pur.price_paid, pur.created_at,
		u.username, u.full_name, u.email,
		p.title_ru, p.title_en, p.type
		FROM purchases pur
		JOIN users u ON pur.user_id = u.id
		JOIN products p ON pur.product_id = p.id
		WHERE ` + where + ` ORDER BY pur.` + sortBy + ` ` + sortDir + ` LIMIT ? OFFSET ?`

	args = append(args, perPage, offset)
	rows, err := db.Query(query, args...)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка загрузки покупок"})
		return
	}
	defer rows.Close()

	items := []AdminPurchaseItem{}
	for rows.Next() {
		var item AdminPurchaseItem
		var accessEnd sql.NullTime
		var purchaseDate time.Time
		var accessStart time.Time
		var createdAt time.Time
		var paymentID sql.NullInt64
		var promoCodeID sql.NullInt64

		err := rows.Scan(
			&item.ID, &item.UserID, &item.ProductID, &purchaseDate, &accessStart,
			&accessEnd, &item.Status, &paymentID, &promoCodeID, &item.PricePaid, &createdAt,
			&item.Username, &item.UserFullName, &item.UserEmail,
			&item.ProductTitleRu, &item.ProductTitleEn, &item.ProductType,
		)
		if err != nil {
			continue
		}

		item.PurchaseDate = purchaseDate.Format("2006-01-02 15:04:05")
		item.AccessStart = accessStart.Format("2006-01-02 15:04:05")
		item.CreatedAt = createdAt.Format("2006-01-02 15:04:05")

		if accessEnd.Valid {
			item.AccessEnd = accessEnd.Time.Format("2006-01-02 15:04:05")
			remaining := time.Until(accessEnd.Time)
			days := int(remaining.Hours() / 24)
			if days < 0 {
				days = 0
			}
			item.DaysRemaining = days
		} else {
			item.AccessEnd = ""
			item.DaysRemaining = -1 // unlimited
		}

		if paymentID.Valid {
			v := int(paymentID.Int64)
			item.PaymentID = &v
		}
		if promoCodeID.Valid {
			v := int(promoCodeID.Int64)
			item.PromoCodeID = &v
		}

		items = append(items, item)
	}

	c.JSON(http.StatusOK, AdminListPurchasesResponse{
		Items:      items,
		Total:      total,
		Page:       page,
		PerPage:    perPage,
		TotalPages: totalPages,
	})
}

// AdminExtendPurchaseAccess extends the access end date for a purchase
func AdminExtendPurchaseAccess(c *gin.Context) {
	id := c.Param("id")

	var req struct {
		Days int `json:"days" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный запрос"})
		return
	}
	if req.Days <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Количество дней должно быть положительным"})
		return
	}

	// Get current access_end
	var currentEnd sql.NullTime
	err := db.QueryRow("SELECT access_end FROM purchases WHERE id = ?", id).Scan(&currentEnd)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Покупка не найдена"})
		return
	}

	var newEnd time.Time
	now := time.Now()
	if currentEnd.Valid && currentEnd.Time.After(now) {
		newEnd = currentEnd.Time.AddDate(0, 0, req.Days)
	} else {
		newEnd = now.AddDate(0, 0, req.Days)
	}

	_, err = db.Exec("UPDATE purchases SET access_end = ?, status = 'active' WHERE id = ?", newEnd, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка продления доступа"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":     "success",
		"access_end": newEnd.Format("2006-01-02 15:04:05"),
	})
}

// AdminCancelPurchase cancels a purchase (sets status to "cancelled")
func AdminCancelPurchase(c *gin.Context) {
	id := c.Param("id")

	_, err := db.Exec("UPDATE purchases SET status = 'cancelled' WHERE id = ?", id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка отмены покупки"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success"})
}

// AdminGetPurchaseDetail returns detailed info about a purchase
func AdminGetPurchaseDetail(c *gin.Context) {
	id := c.Param("id")

	var item AdminPurchaseItem
	var accessEnd sql.NullTime
	var purchaseDate time.Time
	var accessStart time.Time
	var createdAt time.Time
	var paymentID sql.NullInt64
	var promoCodeID sql.NullInt64

	err := db.QueryRow(`
		SELECT pur.id, pur.user_id, pur.product_id, pur.purchase_date, pur.access_start,
			pur.access_end, pur.status, pur.payment_id, pur.promo_code_id, pur.price_paid, pur.created_at,
			u.username, u.full_name, u.email,
			p.title_ru, p.title_en, p.type
		FROM purchases pur
		JOIN users u ON pur.user_id = u.id
		JOIN products p ON pur.product_id = p.id
		WHERE pur.id = ?
	`, id).Scan(
		&item.ID, &item.UserID, &item.ProductID, &purchaseDate, &accessStart,
		&accessEnd, &item.Status, &paymentID, &promoCodeID, &item.PricePaid, &createdAt,
		&item.Username, &item.UserFullName, &item.UserEmail,
		&item.ProductTitleRu, &item.ProductTitleEn, &item.ProductType,
	)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Покупка не найдена"})
		return
	}

	item.PurchaseDate = purchaseDate.Format("2006-01-02 15:04:05")
	item.AccessStart = accessStart.Format("2006-01-02 15:04:05")
	item.CreatedAt = createdAt.Format("2006-01-02 15:04:05")

	if accessEnd.Valid {
		item.AccessEnd = accessEnd.Time.Format("2006-01-02 15:04:05")
		remaining := time.Until(accessEnd.Time)
		days := int(remaining.Hours() / 24)
		if days < 0 {
			days = 0
		}
		item.DaysRemaining = days
	} else {
		item.AccessEnd = ""
		item.DaysRemaining = -1
	}

	if paymentID.Valid {
		v := int(paymentID.Int64)
		item.PaymentID = &v
	}
	if promoCodeID.Valid {
		v := int(promoCodeID.Int64)
		item.PromoCodeID = &v
	}

	c.JSON(http.StatusOK, item)
}
