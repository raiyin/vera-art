package main

import (
	"database/sql"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/raiyin/artserver/models"
)

// GetMasterClasses возвращает список опубликованных мастер-классов с тегами.
// GET /master-classes
// Query params: difficulty, featured, search, limit, offset
func GetMasterClasses(c *gin.Context) {
	difficulty := c.Query("difficulty")
	featured := c.Query("featured")
	search := c.Query("search")
	limitStr := c.Query("limit")
	offsetStr := c.Query("offset")

	// Базовый запрос: выбираем мастер-классы (type='masterclass', status='published')
	query := `
		SELECT
			p.id, p.title_ru, p.title_en,
			p.description_ru, p.description_en,
			p.short_description_ru, p.short_description_en,
			p.price, p.thumbnail_url, p.video_url,
			p.total_duration_minutes, p.difficulty,
			p.category_id, p.is_featured, p.view_count
		FROM products p
		WHERE p.type = 'masterclass' AND p.status = 'published'
	`
	args := []interface{}{}

	if difficulty != "" {
		query += " AND p.difficulty = ?"
		args = append(args, difficulty)
	}

	if featured == "true" {
		query += " AND p.is_featured = 1"
	}

	if search != "" {
		query += " AND (p.title_ru LIKE ? OR p.title_en LIKE ? OR p.description_ru LIKE ? OR p.description_en LIKE ?)"
		searchPattern := "%" + search + "%"
		args = append(args, searchPattern, searchPattern, searchPattern, searchPattern)
	}

	query += " ORDER BY p.is_featured DESC, p.created_at DESC"

	if limitStr != "" {
		limit, err := strconv.Atoi(limitStr)
		if err == nil && limit > 0 {
			query += " LIMIT ?"
			args = append(args, limit)
		}
	}
	if offsetStr != "" {
		offset, err := strconv.Atoi(offsetStr)
		if err == nil && offset >= 0 {
			query += " OFFSET ?"
			args = append(args, offset)
		}
	}

	rows, err := db.Query(query, args...)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "database error: " + err.Error()})
		return
	}
	defer rows.Close()

	masterClasses := []models.MasterClassResponse{}
	for rows.Next() {
		var mc models.MasterClassResponse
		var descriptionRu, descriptionEn sql.NullString
		var shortDescriptionRu, shortDescriptionEn sql.NullString
		var thumbnailUrl, videoUrl sql.NullString
		var difficulty sql.NullString
		var categoryID sql.NullInt64

		err := rows.Scan(
			&mc.Id, &mc.TitleRu, &mc.TitleEn,
			&descriptionRu, &descriptionEn,
			&shortDescriptionRu, &shortDescriptionEn,
			&mc.Price, &thumbnailUrl, &videoUrl,
			&mc.DurationMinutes, &difficulty,
			&categoryID, &mc.IsFeatured, &mc.ViewCount,
		)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "scan error: " + err.Error()})
			return
		}

		if descriptionRu.Valid {
			mc.DescriptionRu = descriptionRu.String
		}
		if descriptionEn.Valid {
			mc.DescriptionEn = descriptionEn.String
		}
		if shortDescriptionRu.Valid {
			mc.ShortDescriptionRu = shortDescriptionRu.String
		}
		if shortDescriptionEn.Valid {
			mc.ShortDescriptionEn = shortDescriptionEn.String
		}
		if thumbnailUrl.Valid {
			mc.ThumbnailUrl = thumbnailUrl.String
		}
		if videoUrl.Valid {
			mc.VideoUrl = videoUrl.String
		}
		if difficulty.Valid {
			mc.Difficulty = difficulty.String
		}
		if categoryID.Valid {
			mc.CategoryId = int(categoryID.Int64)
		}

		mc.IsFree = mc.Price == 0

		// Загружаем теги для мастер-класса
		tags, err := getTagsForProduct(mc.Id)
		if err != nil {
			// Логируем, но не прерываем — теги не критичны
			mc.Tags = []models.TagDTO{}
		} else {
			mc.Tags = tags
		}

		masterClasses = append(masterClasses, mc)
	}

	if masterClasses == nil {
		masterClasses = []models.MasterClassResponse{}
	}

	c.JSON(http.StatusOK, masterClasses)
}

// GetMasterClassByID возвращает один мастер-класс по ID.
// GET /master-classes/:id
func GetMasterClassByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	query := `
		SELECT
			p.id, p.title_ru, p.title_en,
			p.description_ru, p.description_en,
			p.short_description_ru, p.short_description_en,
			p.price, p.thumbnail_url, p.video_url,
			p.total_duration_minutes, p.difficulty,
			p.category_id, p.is_featured, p.view_count
		FROM products p
		WHERE p.id = ? AND p.type = 'masterclass' AND p.status = 'published'
	`

	var mc models.MasterClassResponse
	var descriptionRu, descriptionEn sql.NullString
	var shortDescriptionRu, shortDescriptionEn sql.NullString
	var thumbnailUrl, videoUrl sql.NullString
	var difficulty sql.NullString
	var categoryID sql.NullInt64

	err = db.QueryRow(query, id).Scan(
		&mc.Id, &mc.TitleRu, &mc.TitleEn,
		&descriptionRu, &descriptionEn,
		&shortDescriptionRu, &shortDescriptionEn,
		&mc.Price, &thumbnailUrl, &videoUrl,
		&mc.DurationMinutes, &difficulty,
		&categoryID, &mc.IsFeatured, &mc.ViewCount,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "master class not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "database error: " + err.Error()})
		return
	}

	if descriptionRu.Valid {
		mc.DescriptionRu = descriptionRu.String
	}
	if descriptionEn.Valid {
		mc.DescriptionEn = descriptionEn.String
	}
	if shortDescriptionRu.Valid {
		mc.ShortDescriptionRu = shortDescriptionRu.String
	}
	if shortDescriptionEn.Valid {
		mc.ShortDescriptionEn = shortDescriptionEn.String
	}
	if thumbnailUrl.Valid {
		mc.ThumbnailUrl = thumbnailUrl.String
	}
	if videoUrl.Valid {
		mc.VideoUrl = videoUrl.String
	}
	if difficulty.Valid {
		mc.Difficulty = difficulty.String
	}
	if categoryID.Valid {
		mc.CategoryId = int(categoryID.Int64)
	}

	mc.IsFree = mc.Price == 0

	// Загружаем теги
	tags, err := getTagsForProduct(mc.Id)
	if err != nil {
		mc.Tags = []models.TagDTO{}
	} else {
		mc.Tags = tags
	}

	// Увеличиваем счётчик просмотров
	_, _ = db.Exec("UPDATE products SET view_count = view_count + 1 WHERE id = ?", id)

	c.JSON(http.StatusOK, mc)
}

// GetMasterClassTags возвращает все теги, используемые в мастер-классах.
// GET /master-classes/tags
func GetMasterClassTags(c *gin.Context) {
	query := `
		SELECT DISTINCT t.id, t.name_ru, t.name_en, t.slug, t.created_at
		FROM tags t
		JOIN product_tags pt ON t.id = pt.tag_id
		JOIN products p ON pt.product_id = p.id
		WHERE p.type = 'masterclass' AND p.status = 'published'
		ORDER BY t.name_ru
	`

	rows, err := db.Query(query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "database error: " + err.Error()})
		return
	}
	defer rows.Close()

	tags := []models.Tag{}
	for rows.Next() {
		var tag models.Tag
		err := rows.Scan(&tag.Id, &tag.NameRu, &tag.NameEn, &tag.Slug, &tag.CreatedAt)
		if err != nil {
			continue
		}
		tags = append(tags, tag)
	}

	if tags == nil {
		tags = []models.Tag{}
	}

	c.JSON(http.StatusOK, tags)
}

// getTagsForProduct загружает теги для конкретного продукта.
func getTagsForProduct(productID int) ([]models.TagDTO, error) {
	query := `
		SELECT t.name_ru, t.name_en, t.slug
		FROM tags t
		JOIN product_tags pt ON t.id = pt.tag_id
		WHERE pt.product_id = ?
		ORDER BY t.name_ru
	`

	rows, err := db.Query(query, productID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	tags := []models.TagDTO{}
	for rows.Next() {
		var tag models.TagDTO
		err := rows.Scan(&tag.NameRu, &tag.NameEn, &tag.Slug)
		if err != nil {
			continue
		}
		tags = append(tags, tag)
	}

	return tags, nil
}

// GetMasterClassesByTag возвращает мастер-классы по slug тега.
// GET /master-classes/tag/:tag_slug
func GetMasterClassesByTag(c *gin.Context) {
	tagSlug := c.Param("tag_slug")

	query := `
		SELECT
			p.id, p.title_ru, p.title_en,
			p.description_ru, p.description_en,
			p.short_description_ru, p.short_description_en,
			p.price, p.thumbnail_url, p.video_url,
			p.total_duration_minutes, p.difficulty,
			p.category_id, p.is_featured, p.view_count
		FROM products p
		JOIN product_tags pt ON p.id = pt.product_id
		JOIN tags t ON pt.tag_id = t.id
		WHERE p.type = 'masterclass' AND p.status = 'published' AND t.slug = ?
		ORDER BY p.is_featured DESC, p.created_at DESC
	`

	rows, err := db.Query(query, tagSlug)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "database error: " + err.Error()})
		return
	}
	defer rows.Close()

	masterClasses := []models.MasterClassResponse{}
	for rows.Next() {
		var mc models.MasterClassResponse
		var descriptionRu, descriptionEn sql.NullString
		var shortDescriptionRu, shortDescriptionEn sql.NullString
		var thumbnailUrl, videoUrl sql.NullString
		var difficulty sql.NullString
		var categoryID sql.NullInt64

		err := rows.Scan(
			&mc.Id, &mc.TitleRu, &mc.TitleEn,
			&descriptionRu, &descriptionEn,
			&shortDescriptionRu, &shortDescriptionEn,
			&mc.Price, &thumbnailUrl, &videoUrl,
			&mc.DurationMinutes, &difficulty,
			&categoryID, &mc.IsFeatured, &mc.ViewCount,
		)
		if err != nil {
			continue
		}

		if descriptionRu.Valid {
			mc.DescriptionRu = descriptionRu.String
		}
		if descriptionEn.Valid {
			mc.DescriptionEn = descriptionEn.String
		}
		if shortDescriptionRu.Valid {
			mc.ShortDescriptionRu = shortDescriptionRu.String
		}
		if shortDescriptionEn.Valid {
			mc.ShortDescriptionEn = shortDescriptionEn.String
		}
		if thumbnailUrl.Valid {
			mc.ThumbnailUrl = thumbnailUrl.String
		}
		if videoUrl.Valid {
			mc.VideoUrl = videoUrl.String
		}
		if difficulty.Valid {
			mc.Difficulty = difficulty.String
		}
		if categoryID.Valid {
			mc.CategoryId = int(categoryID.Int64)
		}

		mc.IsFree = mc.Price == 0

		tags, err := getTagsForProduct(mc.Id)
		if err != nil {
			mc.Tags = []models.TagDTO{}
		} else {
			mc.Tags = tags
		}

		masterClasses = append(masterClasses, mc)
	}

	if masterClasses == nil {
		masterClasses = []models.MasterClassResponse{}
	}

	c.JSON(http.StatusOK, masterClasses)
}
