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

// AdminListNewsResponse represents a paginated list of news for the admin panel
type AdminListNewsResponse struct {
	Items      []AdminNewsItem `json:"items"`
	Total      int64           `json:"total"`
	Page       int             `json:"page"`
	PerPage    int             `json:"per_page"`
	TotalPages int             `json:"total_pages"`
}

// AdminNewsItem is a lightweight news representation for the admin table
type AdminNewsItem struct {
	ID          string   `json:"id"`
	TitleRu     string   `json:"title_ru"`
	TitleEn     string   `json:"title_en"`
	Datetime    string   `json:"datetime"`
	Dir         string   `json:"dir"`
	ImgBack     string   `json:"img_back"`
	ImgBackfull string   `json:"img_backfull"`
	TextRu      string   `json:"text_ru"`
	TextEn      string   `json:"text_en"`
	Images      []string `json:"images"`
	Videos      []string `json:"videos"`
}

// AdminGetNews returns a paginated, searchable list of news for the admin panel
func AdminGetNews(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	perPage, _ := strconv.Atoi(c.DefaultQuery("per_page", "20"))
	search := c.Query("search")
	sortBy := c.DefaultQuery("sort_by", "datetime")
	sortDir := c.DefaultQuery("sort_dir", "desc")

	if page < 1 {
		page = 1
	}
	if perPage < 1 || perPage > 100 {
		perPage = 20
	}

	// Validate sort fields
	allowedSortFields := map[string]bool{
		"id": true, "datetime": true, "title_ru": true, "title_en": true,
	}
	if !allowedSortFields[sortBy] {
		sortBy = "datetime"
	}
	if sortDir != "asc" && sortDir != "desc" {
		sortDir = "desc"
	}

	// Build query
	baseQuery := "FROM news n"
	var whereClauses []string
	var args []interface{}

	if search != "" {
		whereClauses = append(whereClauses, "(n.title_ru LIKE ? OR n.title_en LIKE ? OR n.text_ru LIKE ? OR n.text_en LIKE ?)")
		searchPattern := "%" + search + "%"
		args = append(args, searchPattern, searchPattern, searchPattern, searchPattern)
	}

	if len(whereClauses) > 0 {
		baseQuery += " WHERE " + strings.Join(whereClauses, " AND ")
	}

	// Count total
	var total int64
	countQuery := "SELECT COUNT(*) " + baseQuery
	err := db.QueryRow(countQuery, args...).Scan(&total)
	if err != nil {
		log.Printf("Error counting news: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "database error"})
		return
	}

	// Fetch page
	offset := (page - 1) * perPage
	dataQuery := fmt.Sprintf(
		"SELECT n.id, n.datetime, n.title_ru, n.title_en, n.dir, n.img_back, n.img_backfull, n.text_ru, n.text_en, n.images, n.videos %s ORDER BY n.%s %s LIMIT ? OFFSET ?",
		baseQuery, sortBy, sortDir,
	)
	dataArgs := append(args, perPage, offset)

	rows, err := db.Query(dataQuery, dataArgs...)
	if err != nil {
		log.Printf("Error querying news: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "database error"})
		return
	}
	defer rows.Close()

	items := []AdminNewsItem{}
	for rows.Next() {
		var item AdminNewsItem
		var imagesStr, videosStr sql.NullString
		err := rows.Scan(
			&item.ID, &item.Datetime, &item.TitleRu, &item.TitleEn,
			&item.Dir, &item.ImgBack, &item.ImgBackfull,
			&item.TextRu, &item.TextEn, &imagesStr, &videosStr,
		)
		if err != nil {
			log.Printf("Error scanning news row: %v", err)
			continue
		}

		if imagesStr.Valid && imagesStr.String != "" {
			item.Images = strings.Split(imagesStr.String, ";")
		} else {
			item.Images = []string{}
		}

		if videosStr.Valid && videosStr.String != "" {
			item.Videos = strings.Split(videosStr.String, ";")
		} else {
			item.Videos = []string{}
		}

		items = append(items, item)
	}

	totalPages := int(total) / perPage
	if int(total)%perPage > 0 {
		totalPages++
	}

	c.JSON(http.StatusOK, AdminListNewsResponse{
		Items:      items,
		Total:      total,
		Page:       page,
		PerPage:    perPage,
		TotalPages: totalPages,
	})
}

// AdminDeleteNews handles bulk deletion of news
func AdminDeleteNews(c *gin.Context) {
	var req struct {
		IDs []string `json:"ids"`
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
		var dir string
		err := tx.QueryRow("SELECT dir FROM news WHERE id = ?", id).Scan(&dir)
		if err != nil {
			continue
		}

		_, err = tx.Exec("DELETE FROM news WHERE id = ?", id)
		if err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": "database error"})
			return
		}

		// Clean up directory
		dirPath := config.AppConfigInstance.Directories.AbsNewsDir + dir
		if _, err := os.Stat(dirPath); err == nil {
			os.RemoveAll(dirPath)
		}
	}

	err = tx.Commit()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "database error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "news deleted successfully"})
}
