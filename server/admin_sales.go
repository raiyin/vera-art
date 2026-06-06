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

// AdminListSalesResponse represents a paginated list of sales for the admin panel
type AdminListSalesResponse struct {
	Items      []AdminSaleItem `json:"items"`
	Total      int64           `json:"total"`
	Page       int             `json:"page"`
	PerPage    int             `json:"per_page"`
	TotalPages int             `json:"total_pages"`
}

// AdminSaleItem is a lightweight sale representation for the admin table
type AdminSaleItem struct {
	ID          int      `json:"id"`
	StrID       string   `json:"str_id"`
	Dir         string   `json:"dir"`
	NameRu      string   `json:"name_ru"`
	NameEn      string   `json:"name_en"`
	Year        int      `json:"year"`
	Width       int      `json:"width"`
	Height      int      `json:"height"`
	Price       int      `json:"price"`
	BaseRu      string   `json:"base_ru"`
	BaseEn      string   `json:"base_en"`
	Images      []string `json:"images"`
	MaterialsRu []string `json:"materials_ru"`
	MaterialsEn []string `json:"materials_en"`
	CreatedAt   string   `json:"created_at"`
}

// AdminGetSales returns a paginated, filterable list of sales for the admin panel
func AdminGetSales(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	perPage, _ := strconv.Atoi(c.DefaultQuery("per_page", "20"))
	search := c.Query("search")
	baseFilter := c.Query("base_id")
	sortBy := c.DefaultQuery("sort_by", "id")
	sortDir := c.DefaultQuery("sort_dir", "desc")

	if page < 1 {
		page = 1
	}
	if perPage < 1 || perPage > 100 {
		perPage = 20
	}

	// Validate sort fields
	allowedSortFields := map[string]bool{
		"id": true, "year": true, "name_ru": true, "name_en": true,
		"width": true, "height": true, "price": true, "created_at": true,
	}
	if !allowedSortFields[sortBy] {
		sortBy = "id"
	}
	if sortDir != "asc" && sortDir != "desc" {
		sortDir = "desc"
	}

	// Build query
	baseQuery := "FROM sales s LEFT JOIN bases b ON s.base_id = b.id"
	var whereClauses []string
	var args []interface{}

	if search != "" {
		whereClauses = append(whereClauses, "(s.name_ru LIKE ? OR s.name_en LIKE ?)")
		searchPattern := "%" + search + "%"
		args = append(args, searchPattern, searchPattern)
	}

	if baseFilter != "" {
		if baseInt, err := strconv.Atoi(baseFilter); err == nil {
			whereClauses = append(whereClauses, "s.base_id = ?")
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
		log.Printf("Error counting sales: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "database error"})
		return
	}

	// Fetch page
	offset := (page - 1) * perPage
	dataQuery := fmt.Sprintf(
		"SELECT s.id, s.str_id, s.name_ru, s.name_en, s.year, s.width, s.height, s.price, s.images, s.created_at, COALESCE(b.base_ru, ''), COALESCE(b.base_en, '') %s ORDER BY s.%s %s LIMIT ? OFFSET ?",
		baseQuery, sortBy, sortDir,
	)
	dataArgs := append(args, perPage, offset)

	rows, err := db.Query(dataQuery, dataArgs...)
	if err != nil {
		log.Printf("Error querying sales: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "database error"})
		return
	}
	defer rows.Close()

	items := []AdminSaleItem{}
	for rows.Next() {
		var item AdminSaleItem
		var imagesStr sql.NullString
		var createdAt sql.NullString
		err := rows.Scan(
			&item.ID, &item.StrID, &item.NameRu, &item.NameEn,
			&item.Year, &item.Width, &item.Height, &item.Price,
			&imagesStr, &createdAt, &item.BaseRu, &item.BaseEn,
		)
		if err != nil {
			log.Printf("Error scanning sale row: %v", err)
			continue
		}

		dir := config.AppConfigInstance.Directories.RelSalesDir + item.StrID + "/"
		item.Dir = dir

		if imagesStr.Valid && imagesStr.String != "" {
			item.Images = strings.Split(imagesStr.String, ";")
		} else {
			item.Images = []string{}
		}

		if createdAt.Valid {
			item.CreatedAt = createdAt.String
		}

		// Fetch materials for this sale
		matRows, err := db.Query(
			"SELECT m.material_ru, m.material_en FROM materials m JOIN sales_materials sm ON m.id = sm.material_id WHERE sm.sale_id = ?",
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

	c.JSON(http.StatusOK, AdminListSalesResponse{
		Items:      items,
		Total:      total,
		Page:       page,
		PerPage:    perPage,
		TotalPages: totalPages,
	})
}

// AdminDeleteSales handles bulk deletion of sales
func AdminDeleteSales(c *gin.Context) {
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
		var strId string
		err := tx.QueryRow("SELECT str_id FROM sales WHERE id = ?", id).Scan(&strId)
		if err != nil {
			continue
		}

		_, err = tx.Exec("DELETE FROM sales_materials WHERE sale_id = ?", id)
		if err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": "database error"})
			return
		}

		_, err = tx.Exec("DELETE FROM sales WHERE id = ?", id)
		if err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": "database error"})
			return
		}

		dirPath := config.AppConfigInstance.Directories.AbsSalesDir +
			config.AppConfigInstance.Directories.RelSalesDir + strId + "/"
		if _, err := os.Stat(dirPath); err == nil {
			os.RemoveAll(dirPath)
		}
	}

	err = tx.Commit()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "database error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "sales deleted successfully"})
}
