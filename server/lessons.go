package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/raiyin/artserver/models"
)

// GetLessonsByProduct возвращает уроки для конкретного продукта
func GetLessonsByProduct(c *gin.Context) {
	productIdStr := c.Param("product_id")
	productId, err := strconv.Atoi(productIdStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid product id"})
		return
	}

	// Проверяем существование продукта и его статус
	var productStatus string
	err = db.QueryRow("SELECT status FROM products WHERE id = ?", productId).Scan(&productStatus)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "product not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "database error"})
		return
	}

	// Для неопубликованных продуктов проверяем права доступа
	if productStatus != "published" {
		// Проверяем, является ли пользователь админом или владельцем продукта
		claims, exists := c.Get("claims")
		if !exists {
			c.JSON(http.StatusForbidden, gin.H{"error": "access denied"})
			return
		}
		userClaims := claims.(*models.Claims)
		if userClaims.Role != "admin" {
			// Проверяем, является ли пользователь инструктором продукта
			var instructorId sql.NullInt64
			err = db.QueryRow("SELECT instructor_id FROM products WHERE id = ?", productId).Scan(&instructorId)
			if err != nil || !instructorId.Valid || int(instructorId.Int64) != userClaims.UserID {
				c.JSON(http.StatusForbidden, gin.H{"error": "access denied"})
				return
			}
		}
	}

	query := `
		SELECT
			id, product_id, title_ru, title_en, description_ru, description_en,
			content_type, content_url, duration_minutes, sort_order, is_preview,
			resources, homework_ru, homework_en, estimated_study_time, is_required,
			created_at, updated_at
		FROM lessons
		WHERE product_id = ?
		ORDER BY sort_order, id
	`

	rows, err := db.Query(query, productId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "database error: " + err.Error()})
		return
	}
	defer rows.Close()

	lessons := []models.Lesson{}
	for rows.Next() {
		var lesson models.Lesson
		var resourcesJSON []byte
		var estimatedStudyTime sql.NullInt64
		var homeworkRu, homeworkEn sql.NullString
		var descriptionRu, descriptionEn sql.NullString

		err := rows.Scan(
			&lesson.Id, &lesson.ProductId, &lesson.TitleRu, &lesson.TitleEn,
			&descriptionRu, &descriptionEn,
			&lesson.ContentType, &lesson.ContentUrl, &lesson.DurationMinutes,
			&lesson.SortOrder, &lesson.IsPreview,
			&resourcesJSON,
			&homeworkRu, &homeworkEn,
			&estimatedStudyTime,
			&lesson.IsRequired,
			&lesson.CreatedAt, &lesson.UpdatedAt,
		)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "scan error: " + err.Error()})
			return
		}

		// Обработка nullable полей
		if descriptionRu.Valid {
			lesson.DescriptionRu = descriptionRu.String
		}
		if descriptionEn.Valid {
			lesson.DescriptionEn = descriptionEn.String
		}
		if homeworkRu.Valid {
			lesson.HomeworkRu = homeworkRu.String
		}
		if homeworkEn.Valid {
			lesson.HomeworkEn = homeworkEn.String
		}
		if estimatedStudyTime.Valid {
			estTime := int(estimatedStudyTime.Int64)
			lesson.EstimatedStudyTime = &estTime
		}

		// Парсинг JSON ресурсов
		if len(resourcesJSON) > 0 {
			var resources []models.Resource
			if err := json.Unmarshal(resourcesJSON, &resources); err == nil {
				lesson.Resources = resources
			}
		}

		lessons = append(lessons, lesson)
	}

	c.JSON(http.StatusOK, lessons)
}

// GetLessonById возвращает урок по ID
func GetLessonById(c *gin.Context) {
	lessonIdStr := c.Param("id")
	lessonId, err := strconv.Atoi(lessonIdStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid lesson id"})
		return
	}

	query := `
		SELECT
			l.id, l.product_id, l.title_ru, l.title_en, l.description_ru, l.description_en,
			l.content_type, l.content_url, l.duration_minutes, l.sort_order, l.is_preview,
			l.resources, l.homework_ru, l.homework_en, l.estimated_study_time, l.is_required,
			l.created_at, l.updated_at,
			p.status as product_status
		FROM lessons l
		JOIN products p ON l.product_id = p.id
		WHERE l.id = ?
	`

	var lesson models.Lesson
	var resourcesJSON []byte
	var estimatedStudyTime sql.NullInt64
	var homeworkRu, homeworkEn sql.NullString
	var descriptionRu, descriptionEn sql.NullString
	var productStatus string

	err = db.QueryRow(query, lessonId).Scan(
		&lesson.Id, &lesson.ProductId, &lesson.TitleRu, &lesson.TitleEn,
		&descriptionRu, &descriptionEn,
		&lesson.ContentType, &lesson.ContentUrl, &lesson.DurationMinutes,
		&lesson.SortOrder, &lesson.IsPreview,
		&resourcesJSON,
		&homeworkRu, &homeworkEn,
		&estimatedStudyTime,
		&lesson.IsRequired,
		&lesson.CreatedAt, &lesson.UpdatedAt,
		&productStatus,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "lesson not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "database error: " + err.Error()})
		return
	}

	// Проверка доступа для неопубликованных продуктов
	if productStatus != "published" {
		claims, exists := c.Get("claims")
		if !exists {
			c.JSON(http.StatusForbidden, gin.H{"error": "access denied"})
			return
		}
		userClaims := claims.(*models.Claims)
		if userClaims.Role != "admin" {
			// Проверяем, является ли пользователь инструктором продукта
			var instructorId sql.NullInt64
			err = db.QueryRow("SELECT instructor_id FROM products WHERE id = ?", lesson.ProductId).Scan(&instructorId)
			if err != nil || !instructorId.Valid || int(instructorId.Int64) != userClaims.UserID {
				c.JSON(http.StatusForbidden, gin.H{"error": "access denied"})
				return
			}
		}
	}

	// Обработка nullable полей
	if descriptionRu.Valid {
		lesson.DescriptionRu = descriptionRu.String
	}
	if descriptionEn.Valid {
		lesson.DescriptionEn = descriptionEn.String
	}
	if homeworkRu.Valid {
		lesson.HomeworkRu = homeworkRu.String
	}
	if homeworkEn.Valid {
		lesson.HomeworkEn = homeworkEn.String
	}
	if estimatedStudyTime.Valid {
		estTime := int(estimatedStudyTime.Int64)
		lesson.EstimatedStudyTime = &estTime
	}

	// Парсинг JSON ресурсов
	if len(resourcesJSON) > 0 {
		var resources []models.Resource
		if err := json.Unmarshal(resourcesJSON, &resources); err == nil {
			lesson.Resources = resources
		}
	}

	c.JSON(http.StatusOK, lesson)
}

// AdminGetLessons возвращает все уроки для админки
func AdminGetLessons(c *gin.Context) {
	productIdStr := c.Query("product_id")

	var query string
	var args []interface{}

	if productIdStr != "" {
		productId, err := strconv.Atoi(productIdStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid product id"})
			return
		}
		query = `
			SELECT
				id, product_id, title_ru, title_en, description_ru, description_en,
				content_type, content_url, duration_minutes, sort_order, is_preview,
				resources, homework_ru, homework_en, estimated_study_time, is_required,
				created_at, updated_at
			FROM lessons
			WHERE product_id = ?
			ORDER BY sort_order, id
		`
		args = append(args, productId)
	} else {
		query = `
			SELECT
				id, product_id, title_ru, title_en, description_ru, description_en,
				content_type, content_url, duration_minutes, sort_order, is_preview,
				resources, homework_ru, homework_en, estimated_study_time, is_required,
				created_at, updated_at
			FROM lessons
			ORDER BY product_id, sort_order, id
		`
	}

	rows, err := db.Query(query, args...)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "database error: " + err.Error()})
		return
	}
	defer rows.Close()

	lessons := []models.Lesson{}
	for rows.Next() {
		var lesson models.Lesson
		var resourcesJSON []byte
		var estimatedStudyTime sql.NullInt64
		var homeworkRu, homeworkEn sql.NullString
		var descriptionRu, descriptionEn sql.NullString

		err := rows.Scan(
			&lesson.Id, &lesson.ProductId, &lesson.TitleRu, &lesson.TitleEn,
			&descriptionRu, &descriptionEn,
			&lesson.ContentType, &lesson.ContentUrl, &lesson.DurationMinutes,
			&lesson.SortOrder, &lesson.IsPreview,
			&resourcesJSON,
			&homeworkRu, &homeworkEn,
			&estimatedStudyTime,
			&lesson.IsRequired,
			&lesson.CreatedAt, &lesson.UpdatedAt,
		)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "scan error: " + err.Error()})
			return
		}

		// Обработка nullable полей
		if descriptionRu.Valid {
			lesson.DescriptionRu = descriptionRu.String
		}
		if descriptionEn.Valid {
			lesson.DescriptionEn = descriptionEn.String
		}
		if homeworkRu.Valid {
			lesson.HomeworkRu = homeworkRu.String
		}
		if homeworkEn.Valid {
			lesson.HomeworkEn = homeworkEn.String
		}
		if estimatedStudyTime.Valid {
			estTime := int(estimatedStudyTime.Int64)
			lesson.EstimatedStudyTime = &estTime
		}

		// Парсинг JSON ресурсов
		if len(resourcesJSON) > 0 {
			var resources []models.Resource
			if err := json.Unmarshal(resourcesJSON, &resources); err == nil {
				lesson.Resources = resources
			}
		}

		lessons = append(lessons, lesson)
	}

	c.JSON(http.StatusOK, lessons)
}

// AdminCreateLesson создает новый урок (только для админов)
func AdminCreateLesson(c *gin.Context) {
	var req models.LessonCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request: " + err.Error()})
		return
	}

	// Проверяем существование продукта
	var productExists bool
	err := db.QueryRow("SELECT 1 FROM products WHERE id = ?", req.ProductId).Scan(&productExists)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusBadRequest, gin.H{"error": "product not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "database error"})
		return
	}

	// Подготовка ресурсов в JSON
	resourcesJSON, err := json.Marshal(req.Resources)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to marshal resources"})
		return
	}

	query := `
		INSERT INTO lessons (
			product_id, title_ru, title_en, description_ru, description_en,
			content_type, content_url, duration_minutes, sort_order, is_preview,
			resources, homework_ru, homework_en, estimated_study_time, is_required
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`

	result, err := db.Exec(query,
		req.ProductId, req.TitleRu, req.TitleEn,
		req.DescriptionRu, req.DescriptionEn,
		req.ContentType, req.ContentUrl, req.DurationMinutes,
		req.SortOrder, req.IsPreview,
		resourcesJSON,
		req.HomeworkRu, req.HomeworkEn,
		req.EstimatedStudyTime, req.IsRequired,
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "database error: " + err.Error()})
		return
	}

	id, err := result.LastInsertId()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get lesson id"})
		return
	}

	// Обновляем счетчик уроков в продукте
	_, err = db.Exec("UPDATE products SET total_lessons = total_lessons + 1 WHERE id = ?", req.ProductId)
	if err != nil {
		// Логируем ошибку, но не прерываем выполнение
		log.Printf("Failed to update product lesson count: %v", err)
	}

	c.JSON(http.StatusCreated, gin.H{
		"id":      id,
		"message": "lesson created successfully",
	})
}

// AdminUpdateLesson обновляет урок (только для админов)
func AdminUpdateLesson(c *gin.Context) {
	lessonIdStr := c.Param("id")
	lessonId, err := strconv.Atoi(lessonIdStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid lesson id"})
		return
	}

	var req models.LessonUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request: " + err.Error()})
		return
	}

	// Проверяем существование урока
	var exists bool
	err = db.QueryRow("SELECT 1 FROM lessons WHERE id = ?", lessonId).Scan(&exists)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "lesson not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "database error"})
		return
	}

	// Подготовка ресурсов в JSON
	var resourcesJSON []byte
	if req.Resources != nil {
		resourcesJSON, err = json.Marshal(req.Resources)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to marshal resources"})
			return
		}
	}

	// Динамическое построение запроса обновления
	query := "UPDATE lessons SET updated_at = CURRENT_TIMESTAMP"
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
	if req.ContentType != nil && *req.ContentType != "" {
		query += ", content_type = ?"
		args = append(args, *req.ContentType)
	}
	if req.ContentUrl != nil {
		query += ", content_url = ?"
		args = append(args, *req.ContentUrl)
	}
	if req.DurationMinutes != nil {
		query += ", duration_minutes = ?"
		args = append(args, *req.DurationMinutes)
	}
	if req.SortOrder != nil {
		query += ", sort_order = ?"
		args = append(args, *req.SortOrder)
	}
	if req.IsPreview != nil {
		query += ", is_preview = ?"
		args = append(args, *req.IsPreview)
	}
	if req.Resources != nil {
		query += ", resources = ?"
		args = append(args, resourcesJSON)
	}
	if req.HomeworkRu != nil {
		query += ", homework_ru = ?"
		args = append(args, *req.HomeworkRu)
	}
	if req.HomeworkEn != nil {
		query += ", homework_en = ?"
		args = append(args, *req.HomeworkEn)
	}
	if req.EstimatedStudyTime != nil {
		query += ", estimated_study_time = ?"
		args = append(args, *req.EstimatedStudyTime)
	}
	if req.IsRequired != nil {
		query += ", is_required = ?"
		args = append(args, *req.IsRequired)
	}

	query += " WHERE id = ?"
	args = append(args, lessonId)

	_, err = db.Exec(query, args...)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "database error: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "lesson updated successfully"})
}

// AdminDeleteLesson удаляет урок (только для админов)
func AdminDeleteLesson(c *gin.Context) {
	lessonIdStr := c.Param("id")
	lessonId, err := strconv.Atoi(lessonIdStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid lesson id"})
		return
	}

	// Получаем product_id для обновления счетчика
	var productId int
	err = db.QueryRow("SELECT product_id FROM lessons WHERE id = ?", lessonId).Scan(&productId)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "lesson not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "database error"})
		return
	}

	// Удаляем урок
	_, err = db.Exec("DELETE FROM lessons WHERE id = ?", lessonId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "database error: " + err.Error()})
		return
	}

	// Обновляем счетчик уроков в продукте
	_, err = db.Exec("UPDATE products SET total_lessons = total_lessons - 1 WHERE id = ?", productId)
	if err != nil {
		// Логируем ошибку, но не прерываем выполнение
		log.Printf("Failed to update product lesson count: %v", err)
	}

	c.JSON(http.StatusOK, gin.H{"message": "lesson deleted successfully"})
}
