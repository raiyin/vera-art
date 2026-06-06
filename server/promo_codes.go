package main

import (
	"net/http"

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
	var promo models.PromoCode
	if err := c.ShouldBindJSON(&promo); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный запрос"})
		return
	}

	// Проверяем обязательные поля
	if promo.Code == "" || promo.DiscountType == "" || promo.DiscountValue <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Не заполнены обязательные поля"})
		return
	}

	_, err := db.Exec(`
		INSERT INTO promo_codes (code, discount_type, discount_value, max_uses, used_count, valid_from, valid_until, is_active, created_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
	`, promo.Code, promo.DiscountType, promo.DiscountValue, promo.MaxUses, 0,
		promo.ValidFrom, promo.ValidUntil, promo.IsActive, promo.CreatedAt)
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
	var promo models.PromoCode
	if err := c.ShouldBindJSON(&promo); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный запрос"})
		return
	}

	_, err := db.Exec(`
		UPDATE promo_codes
		SET code = ?, discount_type = ?, discount_value = ?, max_uses = ?, valid_from = ?, valid_until = ?, is_active = ?
		WHERE id = ?
	`, promo.Code, promo.DiscountType, promo.DiscountValue, promo.MaxUses,
		promo.ValidFrom, promo.ValidUntil, promo.IsActive, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка обновления промокода"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success"})
}
