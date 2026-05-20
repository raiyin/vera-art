package main

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/raiyin/artserver/models"
)

// GetProducts возвращает список продуктов с фильтрацией
func GetProducts(c *gin.Context) {
	// Параметры фильтрации
	typeParam := c.Query("type")
	categorySlug := c.Query("category")
	difficulty := c.Query("difficulty")
	status := c.Query("status")
	if status == "" {
		status = "published" // по умолчанию только опубликованные
	}
	limitStr := c.Query("limit")
	offsetStr := c.Query("offset")
	featured := c.Query("featured")

	// Базовый запрос с JOIN категорий
	query := `
		SELECT
			p.id, p.type, p.title_ru, p.title_en,
			p.description_ru, p.description_en,
			p.short_description_ru, p.short_description_en,
			p.price, p.duration_days, p.thumbnail_url, p.video_url,
			p.status, p.difficulty, p.total_lessons, p.total_duration_minutes,
			p.category_id, p.instructor_id, p.tags,
			p.prerequisites_ru, p.prerequisites_en,
			p.learning_outcomes_ru, p.learning_outcomes_en,
			p.certificate_available, p.max_students, p.start_date,
			p.language, p.is_featured, p.view_count,
			p.created_at, p.updated_at,
			pc.id, pc.name_ru, pc.name_en, pc.slug,
			pc.description_ru, pc.description_en,
			pc.sort_order, pc.is_active, pc.created_at, pc.updated_at
		FROM products p
		LEFT JOIN product_categories pc ON p.category_id = pc.id
		WHERE p.status = ?
	`

	args := []interface{}{status}
	argIndex := 2 // 1-based после WHERE

	// Фильтр по типу
	if typeParam != "" {
		query += " AND p.type = ?"
		args = append(args, typeParam)
		argIndex++
	}

	// Фильтр по категории (через slug)
	if categorySlug != "" {
		query += " AND pc.slug = ?"
		args = append(args, categorySlug)
		argIndex++
	}

	// Фильтр по сложности
	if difficulty != "" {
		query += " AND p.difficulty = ?"
		args = append(args, difficulty)
		argIndex++
	}

	// Фильтр по featured
	if featured == "true" {
		query += " AND p.is_featured = TRUE"
	}

	// Сортировка
	query += " ORDER BY p.is_featured DESC, p.created_at DESC"

	// Пагинация
	if limitStr != "" {
		limit, err := strconv.Atoi(limitStr)
		if err == nil && limit > 0 {
			query += " LIMIT ?"
			args = append(args, limit)
			argIndex++
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

	products := []models.Product{}
	for rows.Next() {
		var product models.Product
		var tagsJSON []byte
		var categoryID sql.NullInt64
		var instructorID sql.NullInt64
		var durationDays sql.NullInt64
		var maxStudents sql.NullInt64
		var startDate sql.NullTime
		var category models.ProductCategory
		var categoryCreatedAt, categoryUpdatedAt string

		// Nullable string fields
		var descriptionRu, descriptionEn sql.NullString
		var shortDescriptionRu, shortDescriptionEn sql.NullString
		var thumbnailUrl, videoUrl sql.NullString
		var difficulty sql.NullString
		var prerequisitesRu, prerequisitesEn sql.NullString
		var learningOutcomesRu, learningOutcomesEn sql.NullString
		var language sql.NullString

		err := rows.Scan(
			&product.Id, &product.Type, &product.TitleRu, &product.TitleEn,
			&descriptionRu, &descriptionEn,
			&shortDescriptionRu, &shortDescriptionEn,
			&product.Price, &durationDays, &thumbnailUrl, &videoUrl,
			&product.Status, &difficulty, &product.TotalLessons, &product.TotalDurationMinutes,
			&categoryID, &instructorID, &tagsJSON,
			&prerequisitesRu, &prerequisitesEn,
			&learningOutcomesRu, &learningOutcomesEn,
			&product.CertificateAvailable, &maxStudents, &startDate,
			&language, &product.IsFeatured, &product.ViewCount,
			&product.CreatedAt, &product.UpdatedAt,
			&category.Id, &category.NameRu, &category.NameEn, &category.Slug,
			&category.DescriptionRu, &category.DescriptionEn,
			&category.SortOrder, &category.IsActive,
			&categoryCreatedAt, &categoryUpdatedAt,
		)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "scan error: " + err.Error()})
			return
		}

		// Обработка nullable string полей
		if descriptionRu.Valid {
			product.DescriptionRu = descriptionRu.String
		}
		if descriptionEn.Valid {
			product.DescriptionEn = descriptionEn.String
		}
		if shortDescriptionRu.Valid {
			product.ShortDescriptionRu = shortDescriptionRu.String
		}
		if shortDescriptionEn.Valid {
			product.ShortDescriptionEn = shortDescriptionEn.String
		}
		if thumbnailUrl.Valid {
			product.ThumbnailUrl = thumbnailUrl.String
		}
		if videoUrl.Valid {
			product.VideoUrl = videoUrl.String
		}
		if difficulty.Valid {
			product.Difficulty = difficulty.String
		}
		if prerequisitesRu.Valid {
			product.PrerequisitesRu = prerequisitesRu.String
		}
		if prerequisitesEn.Valid {
			product.PrerequisitesEn = prerequisitesEn.String
		}
		if learningOutcomesRu.Valid {
			product.LearningOutcomesRu = learningOutcomesRu.String
		}
		if learningOutcomesEn.Valid {
			product.LearningOutcomesEn = learningOutcomesEn.String
		}
		if language.Valid {
			product.Language = language.String
		}

		// Обработка nullable полей
		if categoryID.Valid {
			catID := int(categoryID.Int64)
			product.CategoryId = &catID
			product.Category = &category
		}
		if instructorID.Valid {
			instID := int(instructorID.Int64)
			product.InstructorId = &instID
		}
		if durationDays.Valid {
			days := int(durationDays.Int64)
			product.DurationDays = &days
		}
		if maxStudents.Valid {
			students := int(maxStudents.Int64)
			product.MaxStudents = &students
		}
		if startDate.Valid {
			product.StartDate = &startDate.Time
		}

		// Парсинг JSON тегов
		if len(tagsJSON) > 0 {
			var tags []string
			if err := json.Unmarshal(tagsJSON, &tags); err == nil {
				product.Tags = tags
			}
		}

		products = append(products, product)
	}

	c.JSON(http.StatusOK, products)
}

// GetProductBySlug возвращает продукт по slug (через категорию + slug продукта)
func GetProductBySlug(c *gin.Context) {
	categorySlug := c.Param("category_slug")
	productSlug := c.Param("product_slug")

	query := `
		SELECT
			p.id, p.type, p.title_ru, p.title_en,
			p.description_ru, p.description_en,
			p.short_description_ru, p.short_description_en,
			p.price, p.duration_days, p.thumbnail_url, p.video_url,
			p.status, p.difficulty, p.total_lessons, p.total_duration_minutes,
			p.category_id, p.instructor_id, p.tags,
			p.prerequisites_ru, p.prerequisites_en,
			p.learning_outcomes_ru, p.learning_outcomes_en,
			p.certificate_available, p.max_students, p.start_date,
			p.language, p.is_featured, p.view_count,
			p.created_at, p.updated_at,
			pc.id, pc.name_ru, pc.name_en, pc.slug,
			pc.description_ru, pc.description_en,
			pc.sort_order, pc.is_active, pc.created_at, pc.updated_at
		FROM products p
		LEFT JOIN product_categories pc ON p.category_id = pc.id
		WHERE pc.slug = ? AND p.status = 'published'
		AND (p.title_ru LIKE ? OR p.title_en LIKE ?)
		LIMIT 1
	`

	searchPattern := "%" + strings.ReplaceAll(productSlug, "-", " ") + "%"

	var product models.Product
	var tagsJSON []byte
	var categoryID sql.NullInt64
	var instructorID sql.NullInt64
	var durationDays sql.NullInt64
	var maxStudents sql.NullInt64
	var startDate sql.NullTime
	var category models.ProductCategory
	var categoryCreatedAt, categoryUpdatedAt string

	err := db.QueryRow(query, categorySlug, searchPattern, searchPattern).Scan(
		&product.Id, &product.Type, &product.TitleRu, &product.TitleEn,
		&product.DescriptionRu, &product.DescriptionEn,
		&product.ShortDescriptionRu, &product.ShortDescriptionEn,
		&product.Price, &durationDays, &product.ThumbnailUrl, &product.VideoUrl,
		&product.Status, &product.Difficulty, &product.TotalLessons, &product.TotalDurationMinutes,
		&categoryID, &instructorID, &tagsJSON,
		&product.PrerequisitesRu, &product.PrerequisitesEn,
		&product.LearningOutcomesRu, &product.LearningOutcomesEn,
		&product.CertificateAvailable, &maxStudents, &startDate,
		&product.Language, &product.IsFeatured, &product.ViewCount,
		&product.CreatedAt, &product.UpdatedAt,
		&category.Id, &category.NameRu, &category.NameEn, &category.Slug,
		&category.DescriptionRu, &category.DescriptionEn,
		&category.SortOrder, &category.IsActive,
		&categoryCreatedAt, &categoryUpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "product not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "database error: " + err.Error()})
		return
	}

	// Обработка nullable полей
	if categoryID.Valid {
		catID := int(categoryID.Int64)
		product.CategoryId = &catID
		product.Category = &category
	}
	if instructorID.Valid {
		instID := int(instructorID.Int64)
		product.InstructorId = &instID
	}
	if durationDays.Valid {
		days := int(durationDays.Int64)
		product.DurationDays = &days
	}
	if maxStudents.Valid {
		students := int(maxStudents.Int64)
		product.MaxStudents = &students
	}
	if startDate.Valid {
		product.StartDate = &startDate.Time
	}

	// Парсинг JSON тегов
	if len(tagsJSON) > 0 {
		var tags []string
		if err := json.Unmarshal(tagsJSON, &tags); err == nil {
			product.Tags = tags
		}
	}

	c.JSON(http.StatusOK, product)
}

// GetProductById возвращает продукт по ID (для админки)
func GetProductById(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid product id"})
		return
	}

	query := `
		SELECT
			p.id, p.type, p.title_ru, p.title_en,
			p.description_ru, p.description_en,
			p.short_description_ru, p.short_description_en,
			p.price, p.duration_days, p.thumbnail_url, p.video_url,
			p.status, p.difficulty, p.total_lessons, p.total_duration_minutes,
			p.category_id, p.instructor_id, p.tags,
			p.prerequisites_ru, p.prerequisites_en,
			p.learning_outcomes_ru, p.learning_outcomes_en,
			p.certificate_available, p.max_students, p.start_date,
			p.language, p.is_featured, p.view_count,
			p.created_at, p.updated_at
		FROM products p
		WHERE p.id = ?
	`

	var product models.Product
	var tagsJSON []byte
	var categoryID sql.NullInt64
	var instructorID sql.NullInt64
	var durationDays sql.NullInt64
	var maxStudents sql.NullInt64
	var startDate sql.NullTime

	err = db.QueryRow(query, id).Scan(
		&product.Id, &product.Type, &product.TitleRu, &product.TitleEn,
		&product.DescriptionRu, &product.DescriptionEn,
		&product.ShortDescriptionRu, &product.ShortDescriptionEn,
		&product.Price, &durationDays, &product.ThumbnailUrl, &product.VideoUrl,
		&product.Status, &product.Difficulty, &product.TotalLessons, &product.TotalDurationMinutes,
		&categoryID, &instructorID, &tagsJSON,
		&product.PrerequisitesRu, &product.PrerequisitesEn,
		&product.LearningOutcomesRu, &product.LearningOutcomesEn,
		&product.CertificateAvailable, &maxStudents, &startDate,
		&product.Language, &product.IsFeatured, &product.ViewCount,
		&product.CreatedAt, &product.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "product not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "database error: " + err.Error()})
		return
	}

	// Обработка nullable полей
	if categoryID.Valid {
		catID := int(categoryID.Int64)
		product.CategoryId = &catID
	}
	if instructorID.Valid {
		instID := int(instructorID.Int64)
		product.InstructorId = &instID
	}
	if durationDays.Valid {
		days := int(durationDays.Int64)
		product.DurationDays = &days
	}
	if maxStudents.Valid {
		students := int(maxStudents.Int64)
		product.MaxStudents = &students
	}
	if startDate.Valid {
		product.StartDate = &startDate.Time
	}

	// Парсинг JSON тегов
	if len(tagsJSON) > 0 {
		var tags []string
		if err := json.Unmarshal(tagsJSON, &tags); err == nil {
			product.Tags = tags
		}
	}

	c.JSON(http.StatusOK, product)
}

// AdminGetProducts возвращает все продукты для админки (включая черновики)
func AdminGetProducts(c *gin.Context) {
	query := `
		SELECT
			p.id, p.type, p.title_ru, p.title_en,
			p.description_ru, p.description_en,
			p.short_description_ru, p.short_description_en,
			p.price, p.duration_days, p.thumbnail_url, p.video_url,
			p.status, p.difficulty, p.total_lessons, p.total_duration_minutes,
			p.category_id, p.instructor_id, p.tags,
			p.prerequisites_ru, p.prerequisites_en,
			p.learning_outcomes_ru, p.learning_outcomes_en,
			p.certificate_available, p.max_students, p.start_date,
			p.language, p.is_featured, p.view_count,
			p.created_at, p.updated_at,
			pc.name_ru, pc.name_en
		FROM products p
		LEFT JOIN product_categories pc ON p.category_id = pc.id
		ORDER BY p.created_at DESC
	`

	rows, err := db.Query(query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "database error: " + err.Error()})
		return
	}
	defer rows.Close()

	products := []models.Product{}
	for rows.Next() {
		var product models.Product
		var tagsJSON []byte
		var categoryID sql.NullInt64
		var instructorID sql.NullInt64
		var durationDays sql.NullInt64
		var maxStudents sql.NullInt64
		var startDate sql.NullTime
		var categoryNameRu, categoryNameEn sql.NullString

		err := rows.Scan(
			&product.Id, &product.Type, &product.TitleRu, &product.TitleEn,
			&product.DescriptionRu, &product.DescriptionEn,
			&product.ShortDescriptionRu, &product.ShortDescriptionEn,
			&product.Price, &durationDays, &product.ThumbnailUrl, &product.VideoUrl,
			&product.Status, &product.Difficulty, &product.TotalLessons, &product.TotalDurationMinutes,
			&categoryID, &instructorID, &tagsJSON,
			&product.PrerequisitesRu, &product.PrerequisitesEn,
			&product.LearningOutcomesRu, &product.LearningOutcomesEn,
			&product.CertificateAvailable, &maxStudents, &startDate,
			&product.Language, &product.IsFeatured, &product.ViewCount,
			&product.CreatedAt, &product.UpdatedAt,
			&categoryNameRu, &categoryNameEn,
		)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "scan error: " + err.Error()})
			return
		}

		// Обработка nullable полей
		if categoryID.Valid {
			catID := int(categoryID.Int64)
			product.CategoryId = &catID
			// Создаем минимальный объект категории
			if categoryNameRu.Valid || categoryNameEn.Valid {
				product.Category = &models.ProductCategory{
					Id:     int(categoryID.Int64),
					NameRu: categoryNameRu.String,
					NameEn: categoryNameEn.String,
				}
			}
		}
		if instructorID.Valid {
			instID := int(instructorID.Int64)
			product.InstructorId = &instID
		}
		if durationDays.Valid {
			days := int(durationDays.Int64)
			product.DurationDays = &days
		}
		if maxStudents.Valid {
			students := int(maxStudents.Int64)
			product.MaxStudents = &students
		}
		if startDate.Valid {
			product.StartDate = &startDate.Time
		}

		// Парсинг JSON тегов
		if len(tagsJSON) > 0 {
			var tags []string
			if err := json.Unmarshal(tagsJSON, &tags); err == nil {
				product.Tags = tags
			}
		}

		products = append(products, product)
	}

	c.JSON(http.StatusOK, products)
}

// AdminCreateProduct создает новый продукт (только для админов)
func AdminCreateProduct(c *gin.Context) {
	var req models.ProductCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request: " + err.Error()})
		return
	}

	// Валидация
	if req.Type != "course" && req.Type != "masterclass" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "type must be 'course' or 'masterclass'"})
		return
	}

	// Подготовка тегов в JSON
	tagsJSON, err := json.Marshal(req.Tags)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to marshal tags"})
		return
	}

	query := `
		INSERT INTO products (
			type, title_ru, title_en, description_ru, description_en,
			short_description_ru, short_description_en, price, duration_days,
			thumbnail_url, video_url, status, difficulty, category_id,
			instructor_id, tags, prerequisites_ru, prerequisites_en,
			learning_outcomes_ru, learning_outcomes_en, certificate_available,
			max_students, start_date, language, is_featured
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`

	result, err := db.Exec(query,
		req.Type, req.TitleRu, req.TitleEn,
		req.DescriptionRu, req.DescriptionEn,
		req.ShortDescriptionRu, req.ShortDescriptionEn,
		req.Price, req.DurationDays,
		req.ThumbnailUrl, req.VideoUrl,
		req.Status, req.Difficulty, req.CategoryId,
		req.InstructorId, tagsJSON,
		req.PrerequisitesRu, req.PrerequisitesEn,
		req.LearningOutcomesRu, req.LearningOutcomesEn,
		req.CertificateAvailable, req.MaxStudents, req.StartDate,
		req.Language, req.IsFeatured,
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "database error: " + err.Error()})
		return
	}

	id, err := result.LastInsertId()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get product id"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"id":      id,
		"message": "product created successfully",
	})
}

// AdminUpdateProduct обновляет продукт (только для админов)
func AdminUpdateProduct(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid product id"})
		return
	}

	var req models.ProductUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request: " + err.Error()})
		return
	}

	// Проверяем существование продукта
	var exists bool
	err = db.QueryRow("SELECT 1 FROM products WHERE id = ?", id).Scan(&exists)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "product not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "database error"})
		return
	}

	// Подготовка тегов в JSON
	var tagsJSON []byte
	if req.Tags != nil {
		tagsJSON, err = json.Marshal(req.Tags)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to marshal tags"})
			return
		}
	}

	// Динамическое построение запроса обновления
	query := "UPDATE products SET updated_at = CURRENT_TIMESTAMP"
	args := []interface{}{}

	if req.TitleRu != nil && *req.TitleRu != "" {
		query += ", title_ru = ?"
		args = append(args, *req.TitleRu)
	}
	if req.TitleEn != nil && *req.TitleEn != "" {
		query += ", title_en = ?"
		args = append(args, *req.TitleEn)
	}
	if req.DescriptionRu != nil {
		query += ", description_ru = ?"
		args = append(args, *req.DescriptionRu)
	}
	if req.DescriptionEn != nil {
		query += ", description_en = ?"
		args = append(args, *req.DescriptionEn)
	}
	if req.ShortDescriptionRu != nil {
		query += ", short_description_ru = ?"
		args = append(args, *req.ShortDescriptionRu)
	}
	if req.ShortDescriptionEn != nil {
		query += ", short_description_en = ?"
		args = append(args, *req.ShortDescriptionEn)
	}
	if req.Price != nil && *req.Price > 0 {
		query += ", price = ?"
		args = append(args, *req.Price)
	}
	if req.DurationDays != nil {
		query += ", duration_days = ?"
		args = append(args, *req.DurationDays)
	}
	if req.ThumbnailUrl != nil {
		query += ", thumbnail_url = ?"
		args = append(args, *req.ThumbnailUrl)
	}
	if req.VideoUrl != nil {
		query += ", video_url = ?"
		args = append(args, *req.VideoUrl)
	}
	if req.Status != nil && *req.Status != "" {
		query += ", status = ?"
		args = append(args, *req.Status)
	}
	if req.Difficulty != nil && *req.Difficulty != "" {
		query += ", difficulty = ?"
		args = append(args, *req.Difficulty)
	}
	if req.CategoryId != nil {
		query += ", category_id = ?"
		args = append(args, *req.CategoryId)
	}
	if req.InstructorId != nil {
		query += ", instructor_id = ?"
		args = append(args, *req.InstructorId)
	}
	if req.Tags != nil {
		query += ", tags = ?"
		args = append(args, tagsJSON)
	}
	if req.PrerequisitesRu != nil {
		query += ", prerequisites_ru = ?"
		args = append(args, *req.PrerequisitesRu)
	}
	if req.PrerequisitesEn != nil {
		query += ", prerequisites_en = ?"
		args = append(args, *req.PrerequisitesEn)
	}
	if req.LearningOutcomesRu != nil {
		query += ", learning_outcomes_ru = ?"
		args = append(args, *req.LearningOutcomesRu)
	}
	if req.LearningOutcomesEn != nil {
		query += ", learning_outcomes_en = ?"
		args = append(args, *req.LearningOutcomesEn)
	}
	if req.CertificateAvailable != nil {
		query += ", certificate_available = ?"
		args = append(args, *req.CertificateAvailable)
	}
	if req.MaxStudents != nil {
		query += ", max_students = ?"
		args = append(args, *req.MaxStudents)
	}
	if req.StartDate != nil {
		query += ", start_date = ?"
		args = append(args, *req.StartDate)
	}
	if req.Language != nil && *req.Language != "" {
		query += ", language = ?"
		args = append(args, *req.Language)
	}
	if req.IsFeatured != nil {
		query += ", is_featured = ?"
		args = append(args, *req.IsFeatured)
	}

	query += " WHERE id = ?"
	args = append(args, id)

	_, err = db.Exec(query, args...)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "database error: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "product updated successfully"})
}

// AdminDeleteProduct удаляет продукт (только для админов)
func AdminDeleteProduct(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid product id"})
		return
	}

	// Проверяем существование продукта
	var exists bool
	err = db.QueryRow("SELECT 1 FROM products WHERE id = ?", id).Scan(&exists)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "product not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "database error"})
		return
	}

	// Удаляем продукт (каскадное удаление уроков через FOREIGN KEY)
	_, err = db.Exec("DELETE FROM products WHERE id = ?", id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "database error: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "product deleted successfully"})
}
