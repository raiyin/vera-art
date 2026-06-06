package main

import (
	"database/sql"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/raiyin/artserver/models"
)

// buildMasterClassAPIURLs преобразует имена файлов из БД в API-URL для клиента.
func buildMasterClassAPIURLs(mc *models.MasterClassResponse) {
	if mc.ThumbnailUrl != "" {
		mc.ThumbnailUrl = fmt.Sprintf("/master-classes/%d/thumbnail", mc.Id)
	}
	if mc.VideoUrl != "" {
		mc.VideoUrl = fmt.Sprintf("/master-classes/%d/video", mc.Id)
	}
}

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

		buildMasterClassAPIURLs(&mc)
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

	buildMasterClassAPIURLs(&mc)

	// Увеличиваем счётчик просмотров
	_, _ = db.Exec("UPDATE products SET view_count = view_count + 1 WHERE id = ?", id)

	c.JSON(http.StatusOK, mc)
}

// storageBasePath — базовый путь к директории с файлами мастер-классов.
// Путь относительно корня сервера (server/).
const storageBasePath = "../storage/master-classes"

// ServeMasterClassVideo отдаёт видеофайл мастер-класса с проверкой доступа.
// GET /master-classes/:id/video
// Поддерживает авторизацию через Bearer-токен (заголовок Authorization) или query-параметр ?token=
// (query-параметр нужен для <video> элементов, которые не могут устанавливать кастомные заголовки).
func ServeMasterClassVideo(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	// Получаем информацию о мастер-классе из БД
	var price int
	var videoUrl sql.NullString
	err = db.QueryRow(
		"SELECT price, video_url FROM products WHERE id = ? AND type = 'masterclass' AND status = 'published'",
		id,
	).Scan(&price, &videoUrl)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "master class not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "database error"})
		return
	}

	if !videoUrl.Valid || videoUrl.String == "" {
		c.JSON(http.StatusNotFound, gin.H{"error": "video not found"})
		return
	}

	// Проверка доступа
	hasAccess := price == 0 // бесплатные доступны всем

	if !hasAccess {
		// Проверяем авторизацию: сначала из заголовка, потом из query-параметра
		userID, exists := c.Get("user_id")
		var userRole string
		if !exists {
			// Пробуем получить токен из query-параметра (для <video> элементов)
			tokenParam := c.Query("token")
			if tokenParam != "" {
				claims := &models.Claims{}
				token, parseErr := jwt.ParseWithClaims(tokenParam, claims, func(token *jwt.Token) (interface{}, error) {
					return jwtKey, nil
				})
				if parseErr == nil && token.Valid && claims.TokenType == "access" {
					userID = claims.UserID
					userRole = claims.Role
					exists = true
				}
			}
		} else {
			// Получаем роль из контекста (установлена AuthMiddlewareOptional)
			if claimsObj, ok := c.Get("claims"); ok {
				if claims, ok := claimsObj.(*models.Claims); ok {
					userRole = claims.Role
				}
			}
		}

		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "authentication required"})
			return
		}

		// Администратор имеет доступ ко всем видео
		if userRole == "admin" {
			hasAccess = true
		} else {
			// Проверяем активную покупку
			var purchaseID int
			err = db.QueryRow(`
				SELECT id FROM purchases
				WHERE user_id = ? AND product_id = ? AND status = 'active'
				AND (access_end IS NULL OR access_end > ?)
			`, userID, id, time.Now()).Scan(&purchaseID)
			if err != nil {
				c.JSON(http.StatusForbidden, gin.H{"error": "access denied. purchase required"})
				return
			}
			hasAccess = true
		}
	}

	if !hasAccess {
		c.JSON(http.StatusForbidden, gin.H{"error": "access denied"})
		return
	}

	// Формируем путь к видеофайлу
	videoPath := filepath.Join(storageBasePath, idStr, videoUrl.String)

	// Проверяем существование файла
	if _, err := os.Stat(videoPath); os.IsNotExist(err) {
		c.JSON(http.StatusNotFound, gin.H{"error": "video file not found"})
		return
	}

	// Открываем и отдаём файл
	file, err := os.Open(videoPath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to open video file"})
		return
	}
	defer file.Close()

	// Получаем информацию о файле для Content-Length и поддержки range-запросов
	fileInfo, err := file.Stat()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to stat video file"})
		return
	}

	// Определяем MIME-тип по расширению
	ext := filepath.Ext(videoPath)
	contentType := "video/mp4"
	switch ext {
	case ".webm":
		contentType = "video/webm"
	case ".ogg":
		contentType = "video/ogg"
	case ".avi":
		contentType = "video/x-msvideo"
	case ".mov":
		contentType = "video/quicktime"
	case ".mkv":
		contentType = "video/x-matroska"
	}

	c.Header("Content-Type", contentType)
	c.Header("Content-Length", fmt.Sprintf("%d", fileInfo.Size()))
	c.Header("Accept-Ranges", "bytes")
	c.Header("Cache-Control", "private, max-age=3600")

	// Поддержка range-запросов для прогрессивной загрузки видео
	rangeHeader := c.GetHeader("Range")
	if rangeHeader != "" {
		var start, end int64
		_, err := fmt.Sscanf(rangeHeader, "bytes=%d-%d", &start, &end)
		if err != nil || end == 0 {
			end = fileInfo.Size() - 1
		}
		if start > end || start < 0 || end >= fileInfo.Size() {
			c.Status(http.StatusRequestedRangeNotSatisfiable)
			return
		}

		file.Seek(start, io.SeekStart)
		c.Status(http.StatusPartialContent)
		c.Header("Content-Range", fmt.Sprintf("bytes %d-%d/%d", start, end, fileInfo.Size()))
		c.Header("Content-Length", fmt.Sprintf("%d", end-start+1))
		io.CopyN(c.Writer, file, end-start+1)
		return
	}

	// Отдаём весь файл
	io.Copy(c.Writer, file)
}

// ServeMasterClassThumbnail отдаёт изображение-обложку мастер-класса.
// GET /master-classes/:id/thumbnail
func ServeMasterClassThumbnail(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	// Получаем имя файла обложки из БД
	var thumbnailUrl sql.NullString
	err = db.QueryRow(
		"SELECT thumbnail_url FROM products WHERE id = ? AND type = 'masterclass' AND status = 'published'",
		id,
	).Scan(&thumbnailUrl)
	if err != nil {
		if err == sql.ErrNoRows {
			// write error to console
			log.Printf("master class not found for master-class: %d", id)
			c.JSON(http.StatusNotFound, gin.H{"error": "master class not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "database error"})
		return
	}

	if !thumbnailUrl.Valid || thumbnailUrl.String == "" {
		log.Printf("thumbnail not found for master-class: %d", id)
		c.JSON(http.StatusNotFound, gin.H{"error": "thumbnail not found"})
		return
	}

	// Формируем путь к файлу
	thumbnailPath := filepath.Join(storageBasePath, idStr, thumbnailUrl.String)

	// Проверяем существование файла
	if _, err := os.Stat(thumbnailPath); os.IsNotExist(err) {
		log.Printf("thumbnail file with name %s not found for master-class: %d", thumbnailUrl.String, id)
		c.JSON(http.StatusNotFound, gin.H{"error": "thumbnail file not found"})
		return
	}

	// Определяем MIME-тип
	ext := filepath.Ext(thumbnailPath)
	contentType := "image/jpeg"
	switch ext {
	case ".png":
		contentType = "image/png"
	case ".webp":
		contentType = "image/webp"
	case ".gif":
		contentType = "image/gif"
	case ".avif":
		contentType = "image/avif"
	}

	c.Header("Content-Type", contentType)
	c.Header("Cache-Control", "public, max-age=86400") // кэш на сутки

	c.File(thumbnailPath)
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

		buildMasterClassAPIURLs(&mc)
		masterClasses = append(masterClasses, mc)
	}

	if masterClasses == nil {
		masterClasses = []models.MasterClassResponse{}
	}

	c.JSON(http.StatusOK, masterClasses)
}
