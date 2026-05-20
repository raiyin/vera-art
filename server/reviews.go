package main

import (
	"database/sql"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/raiyin/artserver/models"
)

// GetProductReviews возвращает отзывы на продукт (публичный доступ)
func GetProductReviews(c *gin.Context) {
	productID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid product id"})
		return
	}

	// Проверяем, существует ли продукт
	var exists bool
	err = db.QueryRow("SELECT EXISTS(SELECT 1 FROM products WHERE id = ? AND status = 'published')", productID).Scan(&exists)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "database error"})
		return
	}
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "product not found"})
		return
	}

	// Получаем отзывы с информацией о пользователе
	query := `
		SELECT
			r.id, r.user_id, r.product_id, r.purchase_id,
			r.rating, r.title_ru, r.title_en, r.comment_ru, r.comment_en,
			r.is_approved, r.is_visible, r.created_at, r.updated_at,
			u.id, u.username, u.full_name, u.email
		FROM reviews r
		JOIN users u ON r.user_id = u.id
		WHERE r.product_id = ? AND r.is_approved = TRUE AND r.is_visible = TRUE
		ORDER BY r.created_at DESC
	`
	rows, err := db.Query(query, productID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "database error"})
		return
	}
	defer rows.Close()

	reviews := []models.Review{}
	for rows.Next() {
		var r models.Review
		var u models.User
		err := rows.Scan(
			&r.Id, &r.UserId, &r.ProductId, &r.PurchaseId,
			&r.Rating, &r.TitleRu, &r.TitleEn, &r.CommentRu, &r.CommentEn,
			&r.IsApproved, &r.IsVisible, &r.CreatedAt, &r.UpdatedAt,
			&u.Id, &u.Username, &u.FullName, &u.Email,
		)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "scan error"})
			return
		}
		r.User = &u
		reviews = append(reviews, r)
	}

	c.JSON(http.StatusOK, reviews)
}

// CreateReview создает отзыв на продукт (только пользователь, купивший продукт)
func CreateReview(c *gin.Context) {
	productID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid product id"})
		return
	}

	// Получаем ID пользователя из контекста (установлено middleware аутентификации)
	claimsInterface, _ := c.Get("claims")
	claims, ok := claimsInterface.(*models.Claims)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "authentication required"})
		return
	}
	userID := claims.UserID

	// Проверяем, есть ли у пользователя активная покупка этого продукта
	var purchaseID int
	err = db.QueryRow(`
		SELECT id FROM purchases
		WHERE user_id = ? AND product_id = ? AND status = 'active'
		AND (access_end IS NULL OR access_end > CURRENT_TIMESTAMP)
		LIMIT 1
	`, userID, productID).Scan(&purchaseID)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusForbidden, gin.H{"error": "you must purchase the product before leaving a review"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "database error"})
		return
	}

	// Проверяем, не оставлял ли пользователь уже отзыв на этот продукт
	var existingReviewID int
	err = db.QueryRow("SELECT id FROM reviews WHERE user_id = ? AND product_id = ?", userID, productID).Scan(&existingReviewID)
	if err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "you have already reviewed this product"})
		return
	} else if err != sql.ErrNoRows {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "database error"})
		return
	}

	// Парсим тело запроса
	var req struct {
		Rating    int    `json:"rating" binding:"required,min=1,max=5"`
		TitleRu   string `json:"title_ru,omitempty"`
		TitleEn   string `json:"title_en,omitempty"`
		CommentRu string `json:"comment_ru,omitempty"`
		CommentEn string `json:"comment_en,omitempty"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	// Вставляем отзыв (по умолчанию is_approved = false, требуется модерация)
	result, err := db.Exec(`
		INSERT INTO reviews (
			user_id, product_id, purchase_id, rating,
			title_ru, title_en, comment_ru, comment_en,
			is_approved, is_visible, created_at, updated_at
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, FALSE, TRUE, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
	`, userID, productID, purchaseID, req.Rating,
		req.TitleRu, req.TitleEn, req.CommentRu, req.CommentEn)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create review"})
		return
	}

	reviewID, _ := result.LastInsertId()
	c.JSON(http.StatusCreated, gin.H{
		"id":          reviewID,
		"message":     "review submitted, awaiting moderation",
		"is_approved": false,
	})
}

// UpdateReview обновляет отзыв (только автор)
func UpdateReview(c *gin.Context) {
	reviewID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid review id"})
		return
	}

	claimsInterface, _ := c.Get("claims")
	claims, ok := claimsInterface.(*models.Claims)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "authentication required"})
		return
	}
	userID := claims.UserID

	// Проверяем, принадлежит ли отзыв пользователю
	var authorID int
	err = db.QueryRow("SELECT user_id FROM reviews WHERE id = ?", reviewID).Scan(&authorID)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "review not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "database error"})
		return
	}
	if authorID != userID && claims.Role != "admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "you can only edit your own reviews"})
		return
	}

	// Парсим обновления
	var req struct {
		Rating    *int   `json:"rating"`
		TitleRu   string `json:"title_ru,omitempty"`
		TitleEn   string `json:"title_en,omitempty"`
		CommentRu string `json:"comment_ru,omitempty"`
		CommentEn string `json:"comment_en,omitempty"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	// Строим динамический запрос
	updates := []string{}
	args := []interface{}{}
	if req.Rating != nil {
		updates = append(updates, "rating = ?")
		args = append(args, *req.Rating)
	}
	if req.TitleRu != "" {
		updates = append(updates, "title_ru = ?")
		args = append(args, req.TitleRu)
	}
	if req.TitleEn != "" {
		updates = append(updates, "title_en = ?")
		args = append(args, req.TitleEn)
	}
	if req.CommentRu != "" {
		updates = append(updates, "comment_ru = ?")
		args = append(args, req.CommentRu)
	}
	if req.CommentEn != "" {
		updates = append(updates, "comment_en = ?")
		args = append(args, req.CommentEn)
	}
	if len(updates) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "no fields to update"})
		return
	}
	updates = append(updates, "updated_at = CURRENT_TIMESTAMP")
	args = append(args, reviewID)

	query := "UPDATE reviews SET " + strings.Join(updates, ", ") + " WHERE id = ?"
	_, err = db.Exec(query, args...)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update review"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "review updated"})
}

// DeleteReview удаляет отзыв (автор или админ)
func DeleteReview(c *gin.Context) {
	reviewID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid review id"})
		return
	}

	claimsInterface, _ := c.Get("claims")
	claims, ok := claimsInterface.(*models.Claims)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "authentication required"})
		return
	}
	userID := claims.UserID

	// Проверяем, принадлежит ли отзыв пользователю или пользователь админ
	var authorID int
	err = db.QueryRow("SELECT user_id FROM reviews WHERE id = ?", reviewID).Scan(&authorID)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "review not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "database error"})
		return
	}
	if authorID != userID && claims.Role != "admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "you can only delete your own reviews"})
		return
	}

	// Удаляем отзыв
	_, err = db.Exec("DELETE FROM reviews WHERE id = ?", reviewID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete review"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "review deleted"})
}

// GetPendingReviews возвращает отзывы, ожидающие модерации (только админ)
func GetPendingReviews(c *gin.Context) {
	query := `
		SELECT
			r.id, r.user_id, r.product_id, r.purchase_id,
			r.rating, r.title_ru, r.title_en, r.comment_ru, r.comment_en,
			r.is_approved, r.is_visible, r.created_at, r.updated_at,
			u.id, u.username, u.full_name, u.email,
			p.title_ru, p.title_en
		FROM reviews r
		JOIN users u ON r.user_id = u.id
		JOIN products p ON r.product_id = p.id
		WHERE r.is_approved = FALSE
		ORDER BY r.created_at DESC
	`
	rows, err := db.Query(query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "database error"})
		return
	}
	defer rows.Close()

	reviews := []models.Review{}
	for rows.Next() {
		var r models.Review
		var u models.User
		var productTitleRu, productTitleEn string
		err := rows.Scan(
			&r.Id, &r.UserId, &r.ProductId, &r.PurchaseId,
			&r.Rating, &r.TitleRu, &r.TitleEn, &r.CommentRu, &r.CommentEn,
			&r.IsApproved, &r.IsVisible, &r.CreatedAt, &r.UpdatedAt,
			&u.Id, &u.Username, &u.FullName, &u.Email,
			&productTitleRu, &productTitleEn,
		)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "scan error"})
			return
		}
		r.User = &u
		r.Product = &models.Product{TitleRu: productTitleRu, TitleEn: productTitleEn}
		reviews = append(reviews, r)
	}

	c.JSON(http.StatusOK, reviews)
}

// ApproveReview одобряет отзыв (админ)
func ApproveReview(c *gin.Context) {
	reviewID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid review id"})
		return
	}

	// Обновляем отзыв
	result, err := db.Exec("UPDATE reviews SET is_approved = TRUE, updated_at = CURRENT_TIMESTAMP WHERE id = ?", reviewID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "database error"})
		return
	}
	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "review not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "review approved"})
}
