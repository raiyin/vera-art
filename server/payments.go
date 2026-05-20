package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/raiyin/artserver/models"
)

// CreatePayment создает платеж в системе (заглушка для ЮKassa)
func CreatePayment(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Не авторизован"})
		return
	}

	var req struct {
		ProductID int    `json:"product_id" binding:"required"`
		PromoCode string `json:"promo_code,omitempty"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный запрос"})
		return
	}

	// Получаем продукт
	var product models.Product
	err := db.QueryRow("SELECT * FROM products WHERE id = ?", req.ProductID).Scan(
		&product.Id, &product.Type, &product.TitleRu, &product.TitleEn,
		&product.DescriptionRu, &product.DescriptionEn, &product.ShortDescriptionRu, &product.ShortDescriptionEn,
		&product.Price, &product.DurationDays, &product.ThumbnailUrl, &product.VideoUrl,
		&product.Status, &product.Difficulty, &product.TotalLessons, &product.TotalDurationMinutes,
		&product.CategoryId, &product.InstructorId, &product.Tags, &product.PrerequisitesRu,
		&product.PrerequisitesEn, &product.LearningOutcomesRu, &product.LearningOutcomesEn,
		&product.CertificateAvailable, &product.MaxStudents, &product.StartDate,
		&product.Language, &product.IsFeatured, &product.ViewCount, &product.CreatedAt, &product.UpdatedAt,
	)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Продукт не найден"})
		return
	}

	// Проверяем промокод
	var discount int
	if req.PromoCode != "" {
		var promo models.PromoCode
		err := db.QueryRow("SELECT * FROM promo_codes WHERE code = ?", req.PromoCode).Scan(
			&promo.Id, &promo.Code, &promo.DiscountType, &promo.DiscountValue,
			&promo.MaxUses, &promo.UsedCount, &promo.ValidFrom, &promo.ValidUntil,
			&promo.IsActive, &promo.CreatedAt,
		)
		if err == nil && promo.IsValid() {
			discount = promo.CalculateDiscount(product.Price)
		}
	}

	finalPrice := product.Price - discount
	if finalPrice < 0 {
		finalPrice = 0
	}

	// Создаем запись платежа в статусе pending
	result, err := db.Exec(
		"INSERT INTO payments (user_id, status, amount, currency, description, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, ?)",
		userID, "pending", finalPrice, "RUB", "Оплата курса "+product.TitleRu, time.Now(), time.Now(),
	)
	if err != nil {
		log.Printf("Ошибка создания платежа: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка создания платежа"})
		return
	}
	paymentID, _ := result.LastInsertId()

	// В реальной интеграции здесь был бы вызов API ЮKassa для создания платежа
	// и получение confirmation_url для перенаправления пользователя
	// Для заглушки возвращаем фейковый URL
	confirmationURL := "https://yookassa.ru/fake-payment/" + strconv.FormatInt(paymentID, 10)

	c.JSON(http.StatusOK, gin.H{
		"payment_id":       paymentID,
		"amount":           finalPrice,
		"currency":         "RUB",
		"confirmation_url": confirmationURL,
		"description":      "Оплата курса " + product.TitleRu,
		"status":           "pending",
	})
}

// GetPaymentStatus возвращает статус платежа
func GetPaymentStatus(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Не авторизован"})
		return
	}

	paymentID := c.Param("id")
	var payment models.Payment
	err := db.QueryRow("SELECT * FROM payments WHERE id = ? AND user_id = ?", paymentID, userID).Scan(
		&payment.Id, &payment.UserId, &payment.ExternalId, &payment.Status,
		&payment.Amount, &payment.Currency, &payment.Description, &payment.PaymentMethod,
		&payment.Metadata, &payment.CreatedAt, &payment.UpdatedAt,
	)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Платеж не найден"})
		return
	}

	c.JSON(http.StatusOK, payment)
}

// WebhookYooKassa обработчик вебхуков от ЮKassa (заглушка)
func WebhookYooKassa(c *gin.Context) {
	// В реальности здесь должна быть проверка подписи
	var notification struct {
		Event  string `json:"event"`
		Object struct {
			ID     string `json:"id"`
			Status string `json:"status"`
			Amount struct {
				Value    string `json:"value"`
				Currency string `json:"currency"`
			} `json:"amount"`
			Metadata map[string]any `json:"metadata"`
		} `json:"object"`
	}
	if err := c.ShouldBindJSON(&notification); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный формат уведомления"})
		return
	}

	// Обновляем статус платежа по external_id
	_, err := db.Exec("UPDATE payments SET status = ?, updated_at = ? WHERE external_id = ?",
		notification.Object.Status, time.Now(), notification.Object.ID)
	if err != nil {
		log.Printf("Ошибка обновления платежа: %v", err)
	}

	// Если платеж успешен, создаем покупку
	if notification.Object.Status == "succeeded" {
		// Получаем payment_id из metadata или иным способом
		// Для простоты предположим, что metadata содержит product_id и user_id
		// В реальности нужно более надежное сопоставление
		metadataJSON, _ := json.Marshal(notification.Object.Metadata)
		log.Printf("Платеж успешен: %s, metadata: %s", notification.Object.ID, string(metadataJSON))
		// Здесь должна быть логика создания purchase
	}

	c.Status(http.StatusOK)
}

// CreatePurchase создает покупку после успешной оплаты
func CreatePurchase(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Не авторизован"})
		return
	}

	var req struct {
		PaymentID int `json:"payment_id" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный запрос"})
		return
	}

	// Проверяем, что платеж успешен
	var payment models.Payment
	err := db.QueryRow("SELECT * FROM payments WHERE id = ? AND user_id = ?", req.PaymentID, userID).Scan(
		&payment.Id, &payment.UserId, &payment.ExternalId, &payment.Status,
		&payment.Amount, &payment.Currency, &payment.Description, &payment.PaymentMethod,
		&payment.Metadata, &payment.CreatedAt, &payment.UpdatedAt,
	)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Платеж не найден"})
		return
	}
	if payment.Status != "succeeded" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Платеж не завершен успешно"})
		return
	}

	// Получаем product_id из metadata (в реальности нужно хранить отдельно)
	// Для простоты предположим, что product_id хранится в metadata
	var productID int
	// Поскольку Metadata является map[string]any, можно использовать напрямую
	if payment.Metadata != nil {
		if pid, ok := payment.Metadata["product_id"].(float64); ok {
			productID = int(pid)
		}
	}
	if productID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Не удалось определить продукт"})
		return
	}

	// Получаем продукт для определения длительности доступа
	var durationDays *int
	err = db.QueryRow("SELECT duration_days FROM products WHERE id = ?", productID).Scan(&durationDays)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Продукт не найден"})
		return
	}

	accessStart := time.Now()
	var accessEnd sql.NullTime
	if durationDays != nil {
		accessEnd = sql.NullTime{Time: accessStart.Add(time.Duration(*durationDays) * 24 * time.Hour), Valid: true}
	}

	// Создаем покупку
	result, err := db.Exec(
		"INSERT INTO purchases (user_id, product_id, purchase_date, access_start, access_end, status, payment_id, price_paid, created_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)",
		userID, productID, accessStart, accessStart, accessEnd, "active", payment.Id, payment.Amount, time.Now(),
	)
	if err != nil {
		log.Printf("Ошибка создания покупки: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка создания покупки"})
		return
	}
	purchaseID, _ := result.LastInsertId()

	c.JSON(http.StatusOK, gin.H{
		"purchase_id":  purchaseID,
		"access_start": accessStart,
		"access_end":   accessEnd,
		"status":       "active",
	})
}
