package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/raiyin/artserver/internal/config"
)

// AdminListWorksResponse represents a paginated list of works for the admin panel
type AdminListWorksResponse struct {
	Items      []AdminWorkItem `json:"items"`
	Total      int64           `json:"total"`
	Page       int             `json:"page"`
	PerPage    int             `json:"per_page"`
	TotalPages int             `json:"total_pages"`
}

// AdminWorkItem is a lightweight work representation for the admin table
type AdminWorkItem struct {
	ID          int      `json:"id"`
	StrID       string   `json:"str_id"`
	Dir         string   `json:"dir"`
	NameRu      string   `json:"name_ru"`
	NameEn      string   `json:"name_en"`
	Year        int      `json:"year"`
	Width       int      `json:"width"`
	Height      int      `json:"height"`
	Type        int      `json:"type"`
	BaseRu      string   `json:"base_ru"`
	BaseEn      string   `json:"base_en"`
	Images      []string `json:"images"`
	MaterialsRu []string `json:"materials_ru"`
	MaterialsEn []string `json:"materials_en"`
	CreatedAt   string   `json:"created_at"`
}

// AdminGetWorks returns a paginated, filterable list of works for the admin panel
func AdminGetWorks(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	perPage, _ := strconv.Atoi(c.DefaultQuery("per_page", "20"))
	search := c.Query("search")
	typeFilter := c.Query("type")
	baseFilter := c.Query("base_id")
	sortBy := c.DefaultQuery("sort_by", "id")
	sortDir := c.DefaultQuery("sort_dir", "desc")

	if page < 1 {
		page = 1
	}
	if perPage < 1 || perPage > 100 {
		perPage = 20
	}

	// Validate sort fields to prevent SQL injection
	allowedSortFields := map[string]bool{
		"id": true, "year": true, "name_ru": true, "name_en": true,
		"width": true, "height": true, "type": true, "created_at": true,
	}
	if !allowedSortFields[sortBy] {
		sortBy = "id"
	}
	if sortDir != "asc" && sortDir != "desc" {
		sortDir = "desc"
	}

	// Build query
	baseQuery := "FROM works w LEFT JOIN bases b ON w.base_id = b.id"
	var whereClauses []string
	var args []interface{}

	if search != "" {
		whereClauses = append(whereClauses, "(w.name_ru LIKE ? OR w.name_en LIKE ?)")
		searchPattern := "%" + search + "%"
		args = append(args, searchPattern, searchPattern)
	}

	if typeFilter != "" {
		if typeInt, err := strconv.Atoi(typeFilter); err == nil {
			whereClauses = append(whereClauses, "w.type = ?")
			args = append(args, typeInt)
		}
	}

	if baseFilter != "" {
		if baseInt, err := strconv.Atoi(baseFilter); err == nil {
			whereClauses = append(whereClauses, "w.base_id = ?")
			args = append(args, baseInt)
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
		log.Printf("Error counting works: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "database error"})
		return
	}

	// Fetch page
	offset := (page - 1) * perPage
	dataQuery := fmt.Sprintf(
		"SELECT w.id, w.str_id, w.name_ru, w.name_en, w.year, w.width, w.height, w.type, w.images, w.created_at, COALESCE(b.base_ru, ''), COALESCE(b.base_en, '') %s ORDER BY w.%s %s LIMIT ? OFFSET ?",
		baseQuery, sortBy, sortDir,
	)
	dataArgs := append(args, perPage, offset)

	rows, err := db.Query(dataQuery, dataArgs...)
	if err != nil {
		log.Printf("Error querying works: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "database error"})
		return
	}
	defer rows.Close()

	items := []AdminWorkItem{}
	for rows.Next() {
		var item AdminWorkItem
		var imagesStr sql.NullString
		var createdAt sql.NullString
		err := rows.Scan(
			&item.ID, &item.StrID, &item.NameRu, &item.NameEn,
			&item.Year, &item.Width, &item.Height, &item.Type,
			&imagesStr, &createdAt, &item.BaseRu, &item.BaseEn,
		)
		if err != nil {
			log.Printf("Error scanning work row: %v", err)
			continue
		}

		dir := config.AppConfigInstance.Directories.RelWorksDir + item.StrID + "/"
		item.Dir = dir

		if imagesStr.Valid && imagesStr.String != "" {
			item.Images = strings.Split(imagesStr.String, ";")
		} else {
			item.Images = []string{}
		}

		if createdAt.Valid {
			item.CreatedAt = createdAt.String
		}

		// Fetch materials for this work
		matRows, err := db.Query(
			"SELECT m.material_ru, m.material_en FROM materials m JOIN works_materials wm ON m.id = wm.material_id WHERE wm.work_id = ?",
			item.ID,
		)
		if err == nil {
			for matRows.Next() {
				var matRu, matEn string
				if err := matRows.Scan(&matRu, &matEn); err == nil {
					item.MaterialsRu = append(item.MaterialsRu, matRu)
					item.MaterialsEn = append(item.MaterialsEn, matEn)
				}
			}
			matRows.Close()
		}
		if item.MaterialsRu == nil {
			item.MaterialsRu = []string{}
		}
		if item.MaterialsEn == nil {
			item.MaterialsEn = []string{}
		}

		items = append(items, item)
	}

	totalPages := int(total) / perPage
	if int(total)%perPage > 0 {
		totalPages++
	}

	c.JSON(http.StatusOK, AdminListWorksResponse{
		Items:      items,
		Total:      total,
		Page:       page,
		PerPage:    perPage,
		TotalPages: totalPages,
	})
}

// AdminDeleteWorks handles bulk deletion of works
func AdminDeleteWorks(c *gin.Context) {
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
		// Get str_id for directory cleanup
		var strId string
		err := tx.QueryRow("SELECT str_id FROM works WHERE id = ?", id).Scan(&strId)
		if err != nil {
			continue
		}

		// Delete materials associations
		_, err = tx.Exec("DELETE FROM works_materials WHERE work_id = ?", id)
		if err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": "database error"})
			return
		}

		// Delete work
		_, err = tx.Exec("DELETE FROM works WHERE id = ?", id)
		if err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": "database error"})
			return
		}

		// Clean up directory
		dirPath := config.AppConfigInstance.Directories.AbsWorksDir +
			config.AppConfigInstance.Directories.RelWorksDir + strId + "/"
		if _, err := os.Stat(dirPath); err == nil {
			os.RemoveAll(dirPath)
		}
	}

	err = tx.Commit()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "database error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "works deleted successfully"})
}
