package main

import (
	"database/sql"
	"fmt"
	"log"
	"math"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

// AdminUserItem — пользователь для списка в админке
type AdminUserItem struct {
	ID             int    `json:"id"`
	Username       string `json:"username"`
	FullName       string `json:"full_name"`
	Email          string `json:"email"`
	Role           string `json:"role"`
	Blocked        bool   `json:"blocked"`
	PurchasesCount int    `json:"purchases_count"`
	ReviewsCount   int    `json:"reviews_count"`
	CreatedAt      string `json:"created_at"`
	UpdatedAt      string `json:"updated_at"`
}

// AdminListUsersResponse — ответ со списком пользователей (с пагинацией)
type AdminListUsersResponse struct {
	Items      []AdminUserItem `json:"items"`
	Total      int64           `json:"total"`
	Page       int             `json:"page"`
	PerPage    int             `json:"per_page"`
	TotalPages int             `json:"total_pages"`
}

// AdminUserDetail — детальная информация о пользователе
type AdminUserDetail struct {
	ID        int                 `json:"id"`
	Username  string              `json:"username"`
	FullName  string              `json:"full_name"`
	Email     string              `json:"email"`
	Role      string              `json:"role"`
	Blocked   bool                `json:"blocked"`
	CreatedAt string              `json:"created_at"`
	UpdatedAt string              `json:"updated_at"`
	Purchases []AdminUserPurchase `json:"purchases"`
	Reviews   []AdminUserReview   `json:"reviews"`
}

// AdminUserPurchase — покупка пользователя
type AdminUserPurchase struct {
	ID             int    `json:"id"`
	ProductID      int    `json:"product_id"`
	ProductTitleRu string `json:"product_title_ru"`
	ProductTitleEn string `json:"product_title_en"`
	ProductType    string `json:"product_type"`
	PricePaid      int    `json:"price_paid"`
	Status         string `json:"status"`
	PurchaseDate   string `json:"purchase_date"`
}

// AdminUserReview — отзыв пользователя
type AdminUserReview struct {
	ID             int    `json:"id"`
	ProductID      int    `json:"product_id"`
	ProductTitleRu string `json:"product_title_ru"`
	ProductTitleEn string `json:"product_title_en"`
	Rating         int    `json:"rating"`
	CommentRu      string `json:"comment_ru"`
	CommentEn      string `json:"comment_en"`
	IsApproved     bool   `json:"is_approved"`
	CreatedAt      string `json:"created_at"`
}

// AdminGetUsersList возвращает список пользователей с пагинацией, поиском и фильтрацией
func AdminGetUsersList(c *gin.Context) {
	pageStr := c.DefaultQuery("page", "1")
	perPageStr := c.DefaultQuery("per_page", "20")
	search := c.Query("search")
	roleFilter := c.Query("role")
	blockedFilter := c.Query("blocked")
	sortBy := c.DefaultQuery("sort_by", "u.id")
	sortDir := c.DefaultQuery("sort_dir", "desc")

	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1
	}

	perPage, err := strconv.Atoi(perPageStr)
	if err != nil || perPage < 1 || perPage > 100 {
		perPage = 20
	}

	// Whitelist sort columns
	allowedSortColumns := map[string]bool{
		"u.id":         true,
		"u.username":   true,
		"u.full_name":  true,
		"u.email":      true,
		"u.role":       true,
		"u.blocked":    true,
		"u.created_at": true,
		"u.updated_at": true,
	}

	if !allowedSortColumns[sortBy] {
		sortBy = "u.id"
	}

	if sortDir != "asc" && sortDir != "desc" {
		sortDir = "desc"
	}

	// Build WHERE clause
	var whereClauses []string
	var args []interface{}

	if search != "" {
		whereClauses = append(whereClauses, "(u.username LIKE ? OR u.full_name LIKE ? OR u.email LIKE ?)")
		searchPattern := "%" + search + "%"
		args = append(args, searchPattern, searchPattern, searchPattern)
	}

	if roleFilter != "" {
		whereClauses = append(whereClauses, "u.role = ?")
		args = append(args, roleFilter)
	}

	if blockedFilter != "" {
		if blockedFilter == "1" || blockedFilter == "true" {
			whereClauses = append(whereClauses, "u.blocked = 1")
		} else if blockedFilter == "0" || blockedFilter == "false" {
			whereClauses = append(whereClauses, "u.blocked = 0")
		}
	}

	whereSQL := ""
	if len(whereClauses) > 0 {
		whereSQL = " WHERE " + strings.Join(whereClauses, " AND ")
	}

	// Count total
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM users u%s", whereSQL)
	var total int64
	err = db.QueryRow(countQuery, args...).Scan(&total)
	if err != nil {
		log.Printf("Error counting users: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "database error"})
		return
	}

	totalPages := int(math.Ceil(float64(total) / float64(perPage)))
	if page > totalPages && totalPages > 0 {
		page = totalPages
	}

	offset := (page - 1) * perPage

	// Fetch users with purchase and review counts
	query := fmt.Sprintf(`
		SELECT
			u.id, u.username, COALESCE(u.full_name, ''), COALESCE(u.email, ''),
			u.role, u.blocked, u.created_at, u.updated_at,
			(SELECT COUNT(*) FROM purchases p WHERE p.user_id = u.id) as purchases_count,
			(SELECT COUNT(*) FROM reviews r WHERE r.user_id = u.id) as reviews_count
		FROM users u
		%s
		ORDER BY %s %s
		LIMIT ? OFFSET ?
	`, whereSQL, sortBy, sortDir)

	queryArgs := append(args, perPage, offset)
	rows, err := db.Query(query, queryArgs...)
	if err != nil {
		log.Printf("Error querying users: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "database error"})
		return
	}
	defer rows.Close()

	items := []AdminUserItem{}
	for rows.Next() {
		var item AdminUserItem
		var createdAt, updatedAt time.Time
		err := rows.Scan(
			&item.ID, &item.Username, &item.FullName, &item.Email,
			&item.Role, &item.Blocked, &createdAt, &updatedAt,
			&item.PurchasesCount, &item.ReviewsCount,
		)
		if err != nil {
			log.Printf("Error scanning user row: %v", err)
			continue
		}
		item.CreatedAt = createdAt.Format("2006-01-02 15:04:05")
		item.UpdatedAt = updatedAt.Format("2006-01-02 15:04:05")
		items = append(items, item)
	}

	if items == nil {
		items = []AdminUserItem{}
	}

	c.JSON(http.StatusOK, AdminListUsersResponse{
		Items:      items,
		Total:      total,
		Page:       page,
		PerPage:    perPage,
		TotalPages: totalPages,
	})
}

// AdminGetUserDetail возвращает детальную информацию о пользователе
func AdminGetUserDetail(c *gin.Context) {
	userIDStr := c.Param("id")
	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user id"})
		return
	}

	// Get user basic info
	var detail AdminUserDetail
	var createdAt, updatedAt time.Time
	err = db.QueryRow(`
		SELECT id, username, COALESCE(full_name, ''), COALESCE(email, ''),
			role, blocked, created_at, updated_at
		FROM users WHERE id = ?
	`, userID).Scan(
		&detail.ID, &detail.Username, &detail.FullName, &detail.Email,
		&detail.Role, &detail.Blocked, &createdAt, &updatedAt,
	)
	if err == sql.ErrNoRows {
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}
	if err != nil {
		log.Printf("Error fetching user detail: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "database error"})
		return
	}
	detail.CreatedAt = createdAt.Format("2006-01-02 15:04:05")
	detail.UpdatedAt = updatedAt.Format("2006-01-02 15:04:05")

	// Get user purchases
	purchaseRows, err := db.Query(`
		SELECT p.id, p.product_id, COALESCE(pr.title_ru, ''), COALESCE(pr.title_en, ''),
			COALESCE(pr.type, ''), p.price_paid, p.status, p.purchase_date
		FROM purchases p
		LEFT JOIN products pr ON p.product_id = pr.id
		WHERE p.user_id = ?
		ORDER BY p.purchase_date DESC
	`, userID)
	if err == nil {
		defer purchaseRows.Close()
		for purchaseRows.Next() {
			var purchase AdminUserPurchase
			var purchaseDate time.Time
			err := purchaseRows.Scan(
				&purchase.ID, &purchase.ProductID,
				&purchase.ProductTitleRu, &purchase.ProductTitleEn,
				&purchase.ProductType, &purchase.PricePaid,
				&purchase.Status, &purchaseDate,
			)
			if err != nil {
				log.Printf("Error scanning purchase row: %v", err)
				continue
			}
			purchase.PurchaseDate = purchaseDate.Format("2006-01-02 15:04:05")
			detail.Purchases = append(detail.Purchases, purchase)
		}
	}
	if detail.Purchases == nil {
		detail.Purchases = []AdminUserPurchase{}
	}

	// Get user reviews
	reviewRows, err := db.Query(`
		SELECT r.id, r.product_id, COALESCE(pr.title_ru, ''), COALESCE(pr.title_en, ''),
			r.rating, COALESCE(r.comment_ru, ''), COALESCE(r.comment_en, ''),
			r.is_approved, r.created_at
		FROM reviews r
		LEFT JOIN products pr ON r.product_id = pr.id
		WHERE r.user_id = ?
		ORDER BY r.created_at DESC
	`, userID)
	if err == nil {
		defer reviewRows.Close()
		for reviewRows.Next() {
			var review AdminUserReview
			var reviewCreatedAt time.Time
			err := reviewRows.Scan(
				&review.ID, &review.ProductID,
				&review.ProductTitleRu, &review.ProductTitleEn,
				&review.Rating, &review.CommentRu, &review.CommentEn,
				&review.IsApproved, &reviewCreatedAt,
			)
			if err != nil {
				log.Printf("Error scanning review row: %v", err)
				continue
			}
			review.CreatedAt = reviewCreatedAt.Format("2006-01-02 15:04:05")
			detail.Reviews = append(detail.Reviews, review)
		}
	}
	if detail.Reviews == nil {
		detail.Reviews = []AdminUserReview{}
	}

	c.JSON(http.StatusOK, detail)
}

// AdminUpdateUserRole изменяет роль пользователя
func AdminUpdateUserRole(c *gin.Context) {
	userIDStr := c.Param("id")
	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user id"})
		return
	}

	var req struct {
		Role string `json:"role" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	if req.Role != "user" && req.Role != "admin" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "role must be 'user' or 'admin'"})
		return
	}

	result, err := db.Exec("UPDATE users SET role = ?, updated_at = CURRENT_TIMESTAMP WHERE id = ?", req.Role, userID)
	if err != nil {
		log.Printf("Error updating user role: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "database error"})
		return
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "role updated"})
}

// AdminToggleUserBlock блокирует или разблокирует пользователя
func AdminToggleUserBlock(c *gin.Context) {
	userIDStr := c.Param("id")
	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user id"})
		return
	}

	var req struct {
		Blocked bool `json:"blocked"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	blockedInt := 0
	if req.Blocked {
		blockedInt = 1
	}

	result, err := db.Exec("UPDATE users SET blocked = ?, updated_at = CURRENT_TIMESTAMP WHERE id = ?", blockedInt, userID)
	if err != nil {
		log.Printf("Error toggling user block: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "database error"})
		return
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}

	action := "unblocked"
	if req.Blocked {
		action = "blocked"
	}
	c.JSON(http.StatusOK, gin.H{"message": fmt.Sprintf("user %s", action)})
}
