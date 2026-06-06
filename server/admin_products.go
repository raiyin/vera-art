package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

// AdminListProductsResponse represents a paginated list of products for the admin panel
type AdminListProductsResponse struct {
	Items      []AdminProductItem `json:"items"`
	Total      int64              `json:"total"`
	Page       int                `json:"page"`
	PerPage    int                `json:"per_page"`
	TotalPages int                `json:"total_pages"`
}

// AdminProductItem is a lightweight product representation for the admin table
type AdminProductItem struct {
	ID               int    `json:"id"`
	Type             string `json:"type"` // "course" or "masterclass"
	TitleRu          string `json:"title_ru"`
	TitleEn          string `json:"title_en"`
	Price            int    `json:"price"` // in kopecks
	Status           string `json:"status"`
	Difficulty       string `json:"difficulty"`
	TotalLessons     int    `json:"total_lessons"`
	TotalDurationMin int    `json:"total_duration_minutes"`
	ViewCount        int    `json:"view_count"`
	IsFeatured       bool   `json:"is_featured"`
	Language         string `json:"language"`
	ThumbnailUrl     string `json:"thumbnail_url"`
	CategoryNameRu   string `json:"category_name_ru"`
	CategoryNameEn   string `json:"category_name_en"`
	CertificateAvail bool   `json:"certificate_available"`
	CreatedAt        string `json:"created_at"`
	UpdatedAt        string `json:"updated_at"`
}

// AdminGetProductsList returns a paginated, filterable list of products for the admin panel
func AdminGetProductsList(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	perPage, _ := strconv.Atoi(c.DefaultQuery("per_page", "20"))
	search := c.Query("search")
	typeFilter := c.Query("type")             // "course" or "masterclass"
	statusFilter := c.Query("status")         // "draft", "published", "archived"
	difficultyFilter := c.Query("difficulty") // "beginner", "intermediate", "advanced"
	categoryFilter := c.Query("category_id")
	sortBy := c.DefaultQuery("sort_by", "created_at")
	sortDir := c.DefaultQuery("sort_dir", "desc")

	if page < 1 {
		page = 1
	}
	if perPage < 1 || perPage > 100 {
		perPage = 20
	}

	// Validate sort fields
	allowedSortFields := map[string]bool{
		"id": true, "title_ru": true, "title_en": true,
		"price": true, "status": true, "difficulty": true,
		"total_lessons": true, "view_count": true,
		"created_at": true, "updated_at": true,
	}
	if !allowedSortFields[sortBy] {
		sortBy = "created_at"
	}
	if sortDir != "asc" && sortDir != "desc" {
		sortDir = "desc"
	}

	// Build query
	baseQuery := `FROM products p LEFT JOIN product_categories pc ON p.category_id = pc.id`
	var whereClauses []string
	var args []interface{}

	if search != "" {
		whereClauses = append(whereClauses, "(p.title_ru LIKE ? OR p.title_en LIKE ?)")
		searchPattern := "%" + search + "%"
		args = append(args, searchPattern, searchPattern)
	}

	if typeFilter != "" {
		whereClauses = append(whereClauses, "p.type = ?")
		args = append(args, typeFilter)
	}

	if statusFilter != "" {
		whereClauses = append(whereClauses, "p.status = ?")
		args = append(args, statusFilter)
	}

	if difficultyFilter != "" {
		whereClauses = append(whereClauses, "p.difficulty = ?")
		args = append(args, difficultyFilter)
	}

	if categoryFilter != "" {
		if catID, err := strconv.Atoi(categoryFilter); err == nil {
			whereClauses = append(whereClauses, "p.category_id = ?")
			args = append(args, catID)
		}
	}

	if len(whereClauses) > 0 {
		baseQuery += " WHERE " + strings.Join(whereClauses, " AND ")
	}

	// Count total
	var total int64
	countQuery := "SELECT COUNT(*) " + baseQuery
	err := db.QueryRow(countQuery, args...).Scan(&total)
	if err != nil {
		log.Printf("Error counting products: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "database error"})
		return
	}

	// Fetch page
	offset := (page - 1) * perPage
	dataQuery := fmt.Sprintf(
		`SELECT p.id, p.type, p.title_ru, p.title_en,
		 p.price, p.status, p.difficulty, p.total_lessons,
		 p.total_duration_minutes, p.view_count, p.is_featured,
		 p.language, p.thumbnail_url, p.certificate_available,
		 p.created_at, p.updated_at,
		 COALESCE(pc.name_ru, ''), COALESCE(pc.name_en, '')
		 %s ORDER BY p.%s %s LIMIT ? OFFSET ?`,
		baseQuery, sortBy, sortDir,
	)
	dataArgs := append(args, perPage, offset)

	rows, err := db.Query(dataQuery, dataArgs...)
	if err != nil {
		log.Printf("Error querying products: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "database error"})
		return
	}
	defer rows.Close()

	items := []AdminProductItem{}
	for rows.Next() {
		var item AdminProductItem
		var createdAt, updatedAt sql.NullString
		var thumbnailURL sql.NullString

		err := rows.Scan(
			&item.ID, &item.Type, &item.TitleRu, &item.TitleEn,
			&item.Price, &item.Status, &item.Difficulty, &item.TotalLessons,
			&item.TotalDurationMin, &item.ViewCount, &item.IsFeatured,
			&item.Language, &thumbnailURL, &item.CertificateAvail,
			&createdAt, &updatedAt,
			&item.CategoryNameRu, &item.CategoryNameEn,
		)
		if err != nil {
			log.Printf("Error scanning product row: %v", err)
			continue
		}

		if thumbnailURL.Valid {
			item.ThumbnailUrl = thumbnailURL.String
		}
		if createdAt.Valid {
			item.CreatedAt = createdAt.String
		}
		if updatedAt.Valid {
			item.UpdatedAt = updatedAt.String
		}

		items = append(items, item)
	}

	totalPages := int(total) / perPage
	if int(total)%perPage > 0 {
		totalPages++
	}

	c.JSON(http.StatusOK, AdminListProductsResponse{
		Items:      items,
		Total:      total,
		Page:       page,
		PerPage:    perPage,
		TotalPages: totalPages,
	})
}

// AdminDeleteProducts handles bulk deletion of products
func AdminDeleteProducts(c *gin.Context) {
	var req struct {
		IDs []int `json:"ids"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	if len(req.IDs) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "no ids provided"})
		return
	}

	tx, err := db.Begin()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "database error"})
		return
	}

	for _, id := range req.IDs {
		// Delete related lessons first
		_, err = tx.Exec("DELETE FROM lessons WHERE product_id = ?", id)
		if err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": "database error"})
			return
		}

		// Delete the product
		_, err = tx.Exec("DELETE FROM products WHERE id = ?", id)
		if err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": "database error"})
			return
		}
	}

	err = tx.Commit()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "database error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "products deleted successfully"})
}

// AdminUpdateProductStatus handles status changes for products
func AdminUpdateProductStatus(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid product id"})
		return
	}

	var req struct {
		Status string `json:"status" binding:"required,oneof=draft published archived"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request: " + err.Error()})
		return
	}

	result, err := db.Exec("UPDATE products SET status = ?, updated_at = CURRENT_TIMESTAMP WHERE id = ?", req.Status, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "database error"})
		return
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "product not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "status updated successfully"})
}
