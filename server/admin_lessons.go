package main

import (
	"database/sql"
	"fmt"
	"log"
	"math"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

// AdminLessonItem — урок для списка в админке (с информацией о продукте)
type AdminLessonItem struct {
	ID              int    `json:"id"`
	ProductID       int    `json:"product_id"`
	ProductTitleRu  string `json:"product_title_ru"`
	ProductTitleEn  string `json:"product_title_en"`
	ProductType     string `json:"product_type"`
	TitleRu         string `json:"title_ru"`
	TitleEn         string `json:"title_en"`
	ContentType     string `json:"content_type"`
	DurationMinutes int    `json:"duration_minutes"`
	SortOrder       int    `json:"sort_order"`
	IsPreview       bool   `json:"is_preview"`
	IsRequired      bool   `json:"is_required"`
	ResourcesCount  int    `json:"resources_count"`
	CreatedAt       string `json:"created_at"`
	UpdatedAt       string `json:"updated_at"`
}

// AdminListLessonsResponse — ответ со списком уроков (с пагинацией)
type AdminListLessonsResponse struct {
	Items      []AdminLessonItem `json:"items"`
	Total      int64             `json:"total"`
	Page       int               `json:"page"`
	PerPage    int               `json:"per_page"`
	TotalPages int               `json:"total_pages"`
}

// AdminGetLessonsList возвращает список уроков с пагинацией, поиском и фильтрацией
func AdminGetLessonsList(c *gin.Context) {
	pageStr := c.DefaultQuery("page", "1")
	perPageStr := c.DefaultQuery("per_page", "20")
	search := c.Query("search")
	productIDStr := c.Query("product_id")
	contentType := c.Query("content_type")
	sortBy := c.DefaultQuery("sort_by", "l.sort_order")
	sortDir := c.DefaultQuery("sort_dir", "asc")

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
		"l.sort_order":       true,
		"l.id":               true,
		"l.title_ru":         true,
		"l.title_en":         true,
		"l.content_type":     true,
		"l.duration_minutes": true,
		"l.created_at":       true,
		"l.updated_at":       true,
	}

	if !allowedSortColumns[sortBy] {
		sortBy = "l.sort_order"
	}

	if sortDir != "asc" && sortDir != "desc" {
		sortDir = "asc"
	}

	// Build WHERE clause
	var whereClauses []string
	var args []interface{}

	if search != "" {
		whereClauses = append(whereClauses, "(l.title_ru LIKE ? OR l.title_en LIKE ?)")
		searchPattern := "%" + search + "%"
		args = append(args, searchPattern, searchPattern)
	}

	if productIDStr != "" {
		productID, err := strconv.Atoi(productIDStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid product_id"})
			return
		}
		whereClauses = append(whereClauses, "l.product_id = ?")
		args = append(args, productID)
	}

	if contentType != "" {
		whereClauses = append(whereClauses, "l.content_type = ?")
		args = append(args, contentType)
	}

	whereSQL := ""
	if len(whereClauses) > 0 {
		whereSQL = "WHERE " + strings.Join(whereClauses, " AND ")
	}

	// Count total
	countQuery := fmt.Sprintf(`
		SELECT COUNT(*)
		FROM lessons l
		%s
	`, whereSQL)

	var total int64
	err = db.QueryRow(countQuery, args...).Scan(&total)
	if err != nil {
		log.Printf("Error counting lessons: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "database error"})
		return
	}

	totalPages := int(math.Ceil(float64(total) / float64(perPage)))
	if totalPages < 1 {
		totalPages = 1
	}
	if page > totalPages {
		page = totalPages
	}

	offset := (page - 1) * perPage

	// Fetch items
	query := fmt.Sprintf(`
		SELECT
			l.id, l.product_id,
			COALESCE(p.title_ru, '') as product_title_ru,
			COALESCE(p.title_en, '') as product_title_en,
			COALESCE(p.type, '') as product_type,
			l.title_ru, l.title_en,
			l.content_type, l.duration_minutes, l.sort_order,
			l.is_preview, l.is_required,
			l.resources,
			l.created_at, l.updated_at
		FROM lessons l
		LEFT JOIN products p ON l.product_id = p.id
		%s
		ORDER BY %s %s
		LIMIT ? OFFSET ?
	`, whereSQL, sortBy, sortDir)

	queryArgs := append(args, perPage, offset)
	rows, err := db.Query(query, queryArgs...)
	if err != nil {
		log.Printf("Error querying lessons: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "database error"})
		return
	}
	defer rows.Close()

	items := []AdminLessonItem{}
	for rows.Next() {
		var item AdminLessonItem
		var resourcesJSON sql.NullString
		var createdAt, updatedAt sql.NullString

		err := rows.Scan(
			&item.ID, &item.ProductID,
			&item.ProductTitleRu, &item.ProductTitleEn, &item.ProductType,
			&item.TitleRu, &item.TitleEn,
			&item.ContentType, &item.DurationMinutes, &item.SortOrder,
			&item.IsPreview, &item.IsRequired,
			&resourcesJSON,
			&createdAt, &updatedAt,
		)
		if err != nil {
			log.Printf("Error scanning lesson row: %v", err)
			continue
		}

		// Count resources from JSON
		if resourcesJSON.Valid && resourcesJSON.String != "" && resourcesJSON.String != "[]" {
			// Simple count by counting commas in JSON array
			trimmed := strings.TrimSpace(resourcesJSON.String)
			if strings.HasPrefix(trimmed, "[") && strings.HasSuffix(trimmed, "]") {
				inner := trimmed[1 : len(trimmed)-1]
				if inner == "" {
					item.ResourcesCount = 0
				} else {
					// Count objects by counting "{" occurrences
					item.ResourcesCount = strings.Count(inner, "{")
				}
			}
		}

		if createdAt.Valid {
			item.CreatedAt = createdAt.String
		}
		if updatedAt.Valid {
			item.UpdatedAt = updatedAt.String
		}

		items = append(items, item)
	}

	if items == nil {
		items = []AdminLessonItem{}
	}

	c.JSON(http.StatusOK, AdminListLessonsResponse{
		Items:      items,
		Total:      total,
		Page:       page,
		PerPage:    perPage,
		TotalPages: totalPages,
	})
}

// AdminDeleteLessons удаляет уроки массово
func AdminDeleteLessons(c *gin.Context) {
	var req struct {
		IDs []int `json:"ids" binding:"required"`
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
	defer tx.Rollback()

	// Build placeholders for IN clause
	placeholders := make([]string, len(req.IDs))
	args := make([]interface{}, len(req.IDs))
	for i, id := range req.IDs {
		placeholders[i] = "?"
		args[i] = id
	}

	query := fmt.Sprintf("DELETE FROM lessons WHERE id IN (%s)", strings.Join(placeholders, ","))
	result, err := tx.Exec(query, args...)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "database error: " + err.Error()})
		return
	}

	affected, _ := result.RowsAffected()

	if err := tx.Commit(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "database error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": fmt.Sprintf("deleted %d lessons", affected),
		"deleted": affected,
	})
}
