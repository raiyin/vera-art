package main

import (
	"database/sql"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

// AdminReviewItem represents a review in the admin list
type AdminReviewItem struct {
	ID             int    `json:"id"`
	UserID         int    `json:"user_id"`
	ProductID      int    `json:"product_id"`
	PurchaseID     int    `json:"purchase_id"`
	Rating         int    `json:"rating"`
	TitleRu        string `json:"title_ru"`
	TitleEn        string `json:"title_en"`
	CommentRu      string `json:"comment_ru"`
	CommentEn      string `json:"comment_en"`
	IsApproved     bool   `json:"is_approved"`
	IsVisible      bool   `json:"is_visible"`
	Status         string `json:"status"`
	CreatedAt      string `json:"created_at"`
	UpdatedAt      string `json:"updated_at"`
	Username       string `json:"username"`
	UserFullName   string `json:"user_full_name"`
	UserEmail      string `json:"user_email"`
	ProductTitleRu string `json:"product_title_ru"`
	ProductTitleEn string `json:"product_title_en"`
	ProductType    string `json:"product_type"`
}

// AdminListReviewsResponse is the paginated response for admin reviews
type AdminListReviewsResponse struct {
	Items      []AdminReviewItem `json:"items"`
	Total      int64             `json:"total"`
	Page       int               `json:"page"`
	PerPage    int               `json:"per_page"`
	TotalPages int               `json:"total_pages"`
	Stats      AdminReviewsStats `json:"stats"`
}

// AdminReviewsStats holds review statistics
type AdminReviewsStats struct {
	TotalReviews  int64   `json:"total_reviews"`
	PendingCount  int64   `json:"pending_count"`
	ApprovedCount int64   `json:"approved_count"`
	RejectedCount int64   `json:"rejected_count"`
	AverageRating float64 `json:"average_rating"`
	FiveStarCount int64   `json:"five_star_count"`
}

// AdminGetReviewsList returns paginated reviews with filters for admin
func AdminGetReviewsList(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	perPage, _ := strconv.Atoi(c.DefaultQuery("per_page", "20"))
	search := c.Query("search")
	status := c.Query("status") // "pending", "approved", "rejected", "all"
	rating := c.Query("rating") // "1", "2", "3", "4", "5", "all"
	sortBy := c.DefaultQuery("sort_by", "created_at")
	sortDir := c.DefaultQuery("sort_dir", "desc")

	if page < 1 {
		page = 1
	}
	if perPage < 1 || perPage > 100 {
		perPage = 20
	}

	// Validate sort_by
	allowedSortFields := map[string]bool{
		"created_at": true,
		"rating":     true,
		"updated_at": true,
	}
	if !allowedSortFields[sortBy] {
		sortBy = "created_at"
	}
	if sortDir != "asc" && sortDir != "desc" {
		sortDir = "desc"
	}

	// Build WHERE conditions
	whereClauses := []string{"1=1"}
	args := []interface{}{}

	// Status filter
	if status == "pending" {
		whereClauses = append(whereClauses, "r.is_approved = FALSE AND r.is_visible = TRUE")
	} else if status == "approved" {
		whereClauses = append(whereClauses, "r.is_approved = TRUE")
	} else if status == "rejected" {
		whereClauses = append(whereClauses, "r.is_approved = FALSE AND r.is_visible = FALSE")
	}

	// Rating filter
	if rating != "" && rating != "all" {
		ratingVal, err := strconv.Atoi(rating)
		if err == nil && ratingVal >= 1 && ratingVal <= 5 {
			whereClauses = append(whereClauses, "r.rating = ?")
			args = append(args, ratingVal)
		}
	}

	// Search filter
	if search != "" {
		whereClauses = append(whereClauses, "(r.comment_ru LIKE ? OR r.comment_en LIKE ? OR r.title_ru LIKE ? OR r.title_en LIKE ? OR u.username LIKE ? OR u.full_name LIKE ?)")
		searchPattern := "%" + search + "%"
		args = append(args, searchPattern, searchPattern, searchPattern, searchPattern, searchPattern, searchPattern)
	}

	whereSQL := strings.Join(whereClauses, " AND ")

	// Count total
	var total int64
	countQuery := `SELECT COUNT(*) FROM reviews r JOIN users u ON r.user_id = u.id JOIN products p ON r.product_id = p.id WHERE ` + whereSQL
	err := db.QueryRow(countQuery, args...).Scan(&total)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "database error"})
		return
	}

	totalPages := int(total) / perPage
	if int(total)%perPage != 0 {
		totalPages++
	}

	offset := (page - 1) * perPage

	// Get stats
	stats := getAdminReviewsStats()

	// Fetch reviews
	query := `
		SELECT
			r.id, r.user_id, r.product_id, r.purchase_id,
			r.rating, COALESCE(r.title_ru, ''), COALESCE(r.title_en, ''),
			COALESCE(r.comment_ru, ''), COALESCE(r.comment_en, ''),
			r.is_approved, r.is_visible, r.created_at, r.updated_at,
			u.username, COALESCE(u.full_name, ''), COALESCE(u.email, ''),
			COALESCE(p.title_ru, ''), COALESCE(p.title_en, ''), COALESCE(p.type, '')
		FROM reviews r
		JOIN users u ON r.user_id = u.id
		JOIN products p ON r.product_id = p.id
		WHERE ` + whereSQL + `
		ORDER BY r.` + sortBy + ` ` + sortDir + `
		LIMIT ? OFFSET ?
	`
	queryArgs := append(args, perPage, offset)
	rows, err := db.Query(query, queryArgs...)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "database error"})
		return
	}
	defer rows.Close()

	items := []AdminReviewItem{}
	for rows.Next() {
		var item AdminReviewItem
		var createdAt, updatedAt string
		err := rows.Scan(
			&item.ID, &item.UserID, &item.ProductID, &item.PurchaseID,
			&item.Rating, &item.TitleRu, &item.TitleEn,
			&item.CommentRu, &item.CommentEn,
			&item.IsApproved, &item.IsVisible, &createdAt, &updatedAt,
			&item.Username, &item.UserFullName, &item.UserEmail,
			&item.ProductTitleRu, &item.ProductTitleEn, &item.ProductType,
		)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "scan error"})
			return
		}
		item.CreatedAt = createdAt
		item.UpdatedAt = updatedAt

		// Compute status
		if item.IsApproved {
			item.Status = "approved"
		} else if !item.IsVisible {
			item.Status = "rejected"
		} else {
			item.Status = "pending"
		}

		items = append(items, item)
	}

	c.JSON(http.StatusOK, AdminListReviewsResponse{
		Items:      items,
		Total:      total,
		Page:       page,
		PerPage:    perPage,
		TotalPages: totalPages,
		Stats:      stats,
	})
}

// getAdminReviewsStats returns review statistics
func getAdminReviewsStats() AdminReviewsStats {
	var stats AdminReviewsStats

	// Total reviews
	db.QueryRow("SELECT COUNT(*) FROM reviews").Scan(&stats.TotalReviews)

	// Pending (not approved, visible)
	db.QueryRow("SELECT COUNT(*) FROM reviews WHERE is_approved = FALSE AND is_visible = TRUE").Scan(&stats.PendingCount)

	// Approved
	db.QueryRow("SELECT COUNT(*) FROM reviews WHERE is_approved = TRUE").Scan(&stats.ApprovedCount)

	// Rejected (not approved, not visible)
	db.QueryRow("SELECT COUNT(*) FROM reviews WHERE is_approved = FALSE AND is_visible = FALSE").Scan(&stats.RejectedCount)

	// Average rating
	var avgRating sql.NullFloat64
	db.QueryRow("SELECT AVG(CAST(rating AS REAL)) FROM reviews WHERE is_approved = TRUE").Scan(&avgRating)
	if avgRating.Valid {
		stats.AverageRating = avgRating.Float64
	}

	// Five star count
	db.QueryRow("SELECT COUNT(*) FROM reviews WHERE rating = 5 AND is_approved = TRUE").Scan(&stats.FiveStarCount)

	return stats
}

// AdminApproveReview approves a review
func AdminApproveReview(c *gin.Context) {
	reviewID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid review id"})
		return
	}

	result, err := db.Exec("UPDATE reviews SET is_approved = TRUE, is_visible = TRUE, updated_at = CURRENT_TIMESTAMP WHERE id = ?", reviewID)
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

// AdminRejectReview rejects (hides) a review
func AdminRejectReview(c *gin.Context) {
	reviewID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid review id"})
		return
	}

	result, err := db.Exec("UPDATE reviews SET is_approved = FALSE, is_visible = FALSE, updated_at = CURRENT_TIMESTAMP WHERE id = ?", reviewID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "database error"})
		return
	}
	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "review not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "review rejected"})
}

// AdminBulkApproveReviews approves multiple reviews at once
func AdminBulkApproveReviews(c *gin.Context) {
	var req struct {
		IDs []int `json:"ids"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	if len(req.IDs) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "no ids provided"})
		return
	}

	// Build placeholders
	placeholders := make([]string, len(req.IDs))
	args := make([]interface{}, len(req.IDs))
	for i, id := range req.IDs {
		placeholders[i] = "?"
		args[i] = id
	}

	query := "UPDATE reviews SET is_approved = TRUE, is_visible = TRUE, updated_at = CURRENT_TIMESTAMP WHERE id IN (" + strings.Join(placeholders, ",") + ")"
	result, err := db.Exec(query, args...)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "database error"})
		return
	}

	affected, _ := result.RowsAffected()
	c.JSON(http.StatusOK, gin.H{"message": "reviews approved", "affected": affected})
}

// AdminBulkRejectReviews rejects multiple reviews at once
func AdminBulkRejectReviews(c *gin.Context) {
	var req struct {
		IDs []int `json:"ids"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	if len(req.IDs) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "no ids provided"})
		return
	}

	placeholders := make([]string, len(req.IDs))
	args := make([]interface{}, len(req.IDs))
	for i, id := range req.IDs {
		placeholders[i] = "?"
		args[i] = id
	}

	query := "UPDATE reviews SET is_approved = FALSE, is_visible = FALSE, updated_at = CURRENT_TIMESTAMP WHERE id IN (" + strings.Join(placeholders, ",") + ")"
	result, err := db.Exec(query, args...)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "database error"})
		return
	}

	affected, _ := result.RowsAffected()
	c.JSON(http.StatusOK, gin.H{"message": "reviews rejected", "affected": affected})
}

// AdminBulkDeleteReviews deletes multiple reviews at once
func AdminBulkDeleteReviews(c *gin.Context) {
	var req struct {
		IDs []int `json:"ids"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	if len(req.IDs) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "no ids provided"})
		return
	}

	placeholders := make([]string, len(req.IDs))
	args := make([]interface{}, len(req.IDs))
	for i, id := range req.IDs {
		placeholders[i] = "?"
		args[i] = id
	}

	query := "DELETE FROM reviews WHERE id IN (" + strings.Join(placeholders, ",") + ")"
	result, err := db.Exec(query, args...)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "database error"})
		return
	}

	affected, _ := result.RowsAffected()
	c.JSON(http.StatusOK, gin.H{"message": "reviews deleted", "affected": affected})
}
