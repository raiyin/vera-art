package main

import (
	"database/sql"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/raiyin/artserver/models"
)

// ValidatePromoCode проверяет промокод и возвращает информацию о скидке
func ValidatePromoCode(c *gin.Context) {
	var req struct {
		Code      string `json:"code" binding:"required"`
		ProductID int    `json:"product_id,omitempty"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный запрос"})
		return
	}

	var promo models.PromoCode
	err := db.QueryRow("SELECT * FROM promo_codes WHERE code = ?", req.Code).Scan(
		&promo.Id, &promo.Code, &promo.DiscountType, &promo.DiscountValue,
		&promo.MaxUses, &promo.UsedCount, &promo.ValidFrom, &promo.ValidUntil,
		&promo.IsActive, &promo.CreatedAt,
	)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Промокод не найден"})
		return
	}

	if !promo.IsValid() {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Промокод недействителен"})
		return
	}

	// Если указан product_id, можно рассчитать скидку для конкретного продукта
	var discountAmount int
	if req.ProductID > 0 {
		var price int
		err := db.QueryRow("SELECT price FROM products WHERE id = ?", req.ProductID).Scan(&price)
		if err == nil {
			discountAmount = promo.CalculateDiscount(price)
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"promo_code":      promo.Code,
		"discount_type":   promo.DiscountType,
		"discount_value":  promo.DiscountValue,
		"discount_amount": discountAmount,
		"is_valid":        true,
	})
}

// AdminGetPromoCodes возвращает все промокоды (админ)
func AdminGetPromoCodes(c *gin.Context) {
	rows, err := db.Query("SELECT * FROM promo_codes ORDER BY created_at DESC")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка базы данных"})
		return
	}
	defer rows.Close()

	var promoCodes []models.PromoCode
	for rows.Next() {
		var promo models.PromoCode
		err := rows.Scan(
			&promo.Id, &promo.Code, &promo.DiscountType, &promo.DiscountValue,
			&promo.MaxUses, &promo.UsedCount, &promo.ValidFrom, &promo.ValidUntil,
			&promo.IsActive, &promo.CreatedAt,
		)
		if err != nil {
			continue
		}
		promoCodes = append(promoCodes, promo)
	}

	c.JSON(http.StatusOK, promoCodes)
}

// AdminCreatePromoCode создает новый промокод (админ)
func AdminCreatePromoCode(c *gin.Context) {
	var req struct {
		Code          string     `json:"code" binding:"required"`
		DiscountType  string     `json:"discount_type" binding:"required"`
		DiscountValue int        `json:"discount_value" binding:"required"`
		MaxUses       *int64     `json:"max_uses"`
		ValidFrom     *time.Time `json:"valid_from"`
		ValidUntil    *time.Time `json:"valid_until"`
		IsActive      bool       `json:"is_active"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный запрос"})
		return
	}

	// Проверяем обязательные поля
	if req.Code == "" || req.DiscountType == "" || req.DiscountValue <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Не заполнены обязательные поля"})
		return
	}

	// Convert request to sql.Null* for DB insertion
	var maxUses sql.NullInt64
	if req.MaxUses != nil {
		maxUses = sql.NullInt64{Int64: *req.MaxUses, Valid: true}
	}
	var validFrom sql.NullTime
	if req.ValidFrom != nil {
		validFrom = sql.NullTime{Time: *req.ValidFrom, Valid: true}
	}
	var validUntil sql.NullTime
	if req.ValidUntil != nil {
		validUntil = sql.NullTime{Time: *req.ValidUntil, Valid: true}
	}

	_, err := db.Exec(`
		INSERT INTO promo_codes (code, discount_type, discount_value, max_uses, used_count, valid_from, valid_until, is_active, created_at)
		VALUES (?, ?, ?, ?, 0, ?, ?, ?, ?)
	`, req.Code, req.DiscountType, req.DiscountValue, maxUses,
		validFrom, validUntil, req.IsActive, time.Now())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка создания промокода"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success"})
}

// AdminDeletePromoCode удаляет промокод (админ)
func AdminDeletePromoCode(c *gin.Context) {
	id := c.Param("id")

	_, err := db.Exec("DELETE FROM promo_codes WHERE id = ?", id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка удаления промокода"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success"})
}

// AdminUpdatePromoCode обновляет промокод (админ)
func AdminUpdatePromoCode(c *gin.Context) {
	id := c.Param("id")
	var req struct {
		Code          string     `json:"code" binding:"required"`
		DiscountType  string     `json:"discount_type" binding:"required"`
		DiscountValue int        `json:"discount_value" binding:"required"`
		MaxUses       *int64     `json:"max_uses"`
		ValidFrom     *time.Time `json:"valid_from"`
		ValidUntil    *time.Time `json:"valid_until"`
		IsActive      bool       `json:"is_active"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный запрос"})
		return
	}

	// Convert request to sql.Null* for DB update
	var maxUses sql.NullInt64
	if req.MaxUses != nil {
		maxUses = sql.NullInt64{Int64: *req.MaxUses, Valid: true}
	}
	var validFrom sql.NullTime
	if req.ValidFrom != nil {
		validFrom = sql.NullTime{Time: *req.ValidFrom, Valid: true}
	}
	var validUntil sql.NullTime
	if req.ValidUntil != nil {
		validUntil = sql.NullTime{Time: *req.ValidUntil, Valid: true}
	}

	_, err := db.Exec(`
		UPDATE promo_codes
		SET code = ?, discount_type = ?, discount_value = ?, max_uses = ?, valid_from = ?, valid_until = ?, is_active = ?
		WHERE id = ?
	`, req.Code, req.DiscountType, req.DiscountValue, maxUses,
		validFrom, validUntil, req.IsActive, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка обновления промокода"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success"})
}
