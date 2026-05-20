package main

import (
	"database/sql"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/raiyin/artserver/models"
)

// GetMyCourses возвращает активные курсы пользователя
func GetMyCourses(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Не авторизован"})
		return
	}

	rows, err := db.Query(`
		SELECT p.*, pr.access_start, pr.access_end, pr.status as purchase_status
		FROM purchases pr
		JOIN products p ON pr.product_id = p.id
		WHERE pr.user_id = ? AND pr.status = 'active'
		AND (pr.access_end IS NULL OR pr.access_end > ?)
		ORDER BY pr.created_at DESC
	`, userID, time.Now())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка базы данных"})
		return
	}
	defer rows.Close()

	var courses []map[string]any
	for rows.Next() {
		var product models.Product
		var accessStart, accessEnd sql.NullTime
		var purchaseStatus string
		err := rows.Scan(
			&product.Id, &product.Type, &product.TitleRu, &product.TitleEn,
			&product.DescriptionRu, &product.DescriptionEn, &product.ShortDescriptionRu, &product.ShortDescriptionEn,
			&product.Price, &product.DurationDays, &product.ThumbnailUrl, &product.VideoUrl,
			&product.Status, &product.Difficulty, &product.TotalLessons, &product.TotalDurationMinutes,
			&product.CategoryId, &product.InstructorId, &product.Tags, &product.PrerequisitesRu,
			&product.PrerequisitesEn, &product.LearningOutcomesRu, &product.LearningOutcomesEn,
			&product.CertificateAvailable, &product.MaxStudents, &product.StartDate,
			&product.Language, &product.IsFeatured, &product.ViewCount, &product.CreatedAt, &product.UpdatedAt,
			&accessStart, &accessEnd, &purchaseStatus,
		)
		if err != nil {
			continue
		}
		courses = append(courses, map[string]any{
			"product":         product,
			"access_start":    accessStart,
			"access_end":      accessEnd,
			"purchase_status": purchaseStatus,
		})
	}

	c.JSON(http.StatusOK, courses)
}

// GetLearningProgress возвращает прогресс по конкретной покупке
func GetLearningProgress(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Не авторизован"})
		return
	}

	purchaseID := c.Param("purchase_id")
	// Проверяем, что покупка принадлежит пользователю
	var purchase models.Purchase
	err := db.QueryRow(`
		SELECT * FROM purchases WHERE id = ? AND user_id = ?
	`, purchaseID, userID).Scan(
		&purchase.Id, &purchase.UserId, &purchase.ProductId, &purchase.PurchaseDate,
		&purchase.AccessStart, &purchase.AccessEnd, &purchase.Status,
		&purchase.PaymentId, &purchase.PromoCodeId, &purchase.PricePaid, &purchase.CreatedAt,
	)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Покупка не найдена"})
		return
	}

	// Получаем уроки продукта
	rows, err := db.Query(`
		SELECT l.*, lp.completed, lp.completed_at, lp.watch_duration_seconds, lp.last_position_seconds
		FROM lessons l
		LEFT JOIN learning_progress lp ON l.id = lp.lesson_id AND lp.user_id = ?
		WHERE l.product_id = ?
		ORDER BY l.sort_order
	`, userID, purchase.ProductId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка базы данных"})
		return
	}
	defer rows.Close()

	var lessons []map[string]any
	for rows.Next() {
		var lesson models.Lesson
		var completed sql.NullBool
		var completedAt sql.NullTime
		var watchDuration, lastPosition sql.NullInt64
		err := rows.Scan(
			&lesson.Id, &lesson.ProductId, &lesson.TitleRu, &lesson.TitleEn,
			&lesson.DescriptionRu, &lesson.DescriptionEn, &lesson.ContentType,
			&lesson.ContentUrl, &lesson.DurationMinutes, &lesson.SortOrder,
			&lesson.IsPreview, &lesson.Resources, &lesson.HomeworkRu, &lesson.HomeworkEn,
			&lesson.EstimatedStudyTime, &lesson.IsRequired, &lesson.CreatedAt, &lesson.UpdatedAt,
			&completed, &completedAt, &watchDuration, &lastPosition,
		)
		if err != nil {
			continue
		}
		lessons = append(lessons, map[string]any{
			"lesson":         lesson,
			"completed":      completed.Bool,
			"completed_at":   completedAt,
			"watch_duration": watchDuration.Int64,
			"last_position":  lastPosition.Int64,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"purchase": purchase,
		"lessons":  lessons,
	})
}

// UpdateLessonProgress обновляет прогресс по уроку (отмечает как пройденный, сохраняет позицию)
func UpdateLessonProgress(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Не авторизован"})
		return
	}

	lessonID := c.Param("lesson_id")
	var req struct {
		Completed            *bool `json:"completed,omitempty"`
		LastPositionSeconds  *int  `json:"last_position_seconds,omitempty"`
		WatchDurationSeconds *int  `json:"watch_duration_seconds,omitempty"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный запрос"})
		return
	}

	// Проверяем, что урок существует и пользователь имеет к нему доступ (через покупку)
	var productID int
	err := db.QueryRow("SELECT product_id FROM lessons WHERE id = ?", lessonID).Scan(&productID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Урок не найден"})
		return
	}

	// Проверяем активную покупку
	var purchaseID int
	err = db.QueryRow(`
		SELECT id FROM purchases
		WHERE user_id = ? AND product_id = ? AND status = 'active'
		AND (access_end IS NULL OR access_end > ?)
	`, userID, productID, time.Now()).Scan(&purchaseID)
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": "Доступ к уроку запрещен"})
		return
	}

	// Обновляем или вставляем запись прогресса
	if req.Completed != nil && *req.Completed {
		_, err = db.Exec(`
			INSERT OR REPLACE INTO learning_progress
			(user_id, purchase_id, lesson_id, completed, completed_at, watch_duration_seconds, last_position_seconds, created_at, updated_at)
			VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
		`, userID, purchaseID, lessonID, true, time.Now(),
			req.WatchDurationSeconds, req.LastPositionSeconds, time.Now(), time.Now())
	} else {
		// Обновляем только позицию и длительность просмотра
		_, err = db.Exec(`
			INSERT OR REPLACE INTO learning_progress
			(user_id, purchase_id, lesson_id, completed, completed_at, watch_duration_seconds, last_position_seconds, created_at, updated_at)
			VALUES (?, ?, ?, COALESCE((SELECT completed FROM learning_progress WHERE user_id = ? AND lesson_id = ?), FALSE),
				(SELECT completed_at FROM learning_progress WHERE user_id = ? AND lesson_id = ?),
				?, ?, ?, ?)
		`, userID, purchaseID, lessonID, userID, lessonID, userID, lessonID,
			req.WatchDurationSeconds, req.LastPositionSeconds, time.Now(), time.Now())
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка обновления прогресса"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success"})
}
