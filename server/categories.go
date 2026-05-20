package main

import (
	"database/sql"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/raiyin/artserver/models"
)

// GetCategories возвращает список всех категорий
func GetCategories(c *gin.Context) {
	query := `SELECT id, name_ru, name_en, slug, description_ru, description_en,
		sort_order, is_active, created_at, updated_at
		FROM product_categories
		WHERE is_active = TRUE
		ORDER BY sort_order, name_ru`

	rows, err := db.Query(query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "database error"})
		return
	}
	defer rows.Close()

	categories := []models.ProductCategory{}
	for rows.Next() {
		var category models.ProductCategory
		err := rows.Scan(
			&category.Id, &category.NameRu, &category.NameEn,
			&category.Slug, &category.DescriptionRu, &category.DescriptionEn,
			&category.SortOrder, &category.IsActive,
			&category.CreatedAt, &category.UpdatedAt,
		)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "scan error"})
			return
		}
		categories = append(categories, category)
	}

	c.JSON(http.StatusOK, categories)
}

// GetCategoryBySlug возвращает категорию по slug
func GetCategoryBySlug(c *gin.Context) {
	slug := c.Param("slug")

	var category models.ProductCategory
	query := `SELECT id, name_ru, name_en, slug, description_ru, description_en,
		sort_order, is_active, created_at, updated_at
		FROM product_categories
		WHERE slug = ? AND is_active = TRUE`

	err := db.QueryRow(query, slug).Scan(
		&category.Id, &category.NameRu, &category.NameEn,
		&category.Slug, &category.DescriptionRu, &category.DescriptionEn,
		&category.SortOrder, &category.IsActive,
		&category.CreatedAt, &category.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "category not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "database error"})
		return
	}

	c.JSON(http.StatusOK, category)
}

// AdminGetCategories возвращает все категории (включая неактивные) для админки
func AdminGetCategories(c *gin.Context) {
	query := `SELECT id, name_ru, name_en, slug, description_ru, description_en,
		sort_order, is_active, created_at, updated_at
		FROM product_categories
		ORDER BY sort_order, name_ru`

	rows, err := db.Query(query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "database error"})
		return
	}
	defer rows.Close()

	categories := []models.ProductCategory{}
	for rows.Next() {
		var category models.ProductCategory
		err := rows.Scan(
			&category.Id, &category.NameRu, &category.NameEn,
			&category.Slug, &category.DescriptionRu, &category.DescriptionEn,
			&category.SortOrder, &category.IsActive,
			&category.CreatedAt, &category.UpdatedAt,
		)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "scan error"})
			return
		}
		categories = append(categories, category)
	}

	c.JSON(http.StatusOK, categories)
}

// AdminCreateCategory создает новую категорию
func AdminCreateCategory(c *gin.Context) {
	var req models.ProductCategoryCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Проверяем, существует ли категория с таким slug
	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM product_categories WHERE slug = ?", req.Slug).Scan(&count)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "database error"})
		return
	}
	if count > 0 {
		c.JSON(http.StatusConflict, gin.H{"error": "category with this slug already exists"})
		return
	}

	query := `INSERT INTO product_categories
		(name_ru, name_en, slug, description_ru, description_en, sort_order, is_active)
		VALUES (?, ?, ?, ?, ?, ?, ?)`

	result, err := db.Exec(query,
		req.NameRu, req.NameEn, req.Slug,
		req.DescriptionRu, req.DescriptionEn,
		req.SortOrder, req.IsActive,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not create category"})
		return
	}

	id, _ := result.LastInsertId()
	c.JSON(http.StatusCreated, gin.H{
		"id":      id,
		"message": "category created successfully",
	})
}

// AdminUpdateCategory обновляет категорию
func AdminUpdateCategory(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid category id"})
		return
	}

	var req models.ProductCategoryUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Проверяем существование категории
	var exists bool
	err = db.QueryRow("SELECT EXISTS(SELECT 1 FROM product_categories WHERE id = ?)", id).Scan(&exists)
	if err != nil || !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "category not found"})
		return
	}

	// Если обновляется slug, проверяем уникальность
	if req.Slug != "" {
		var count int
		err = db.QueryRow("SELECT COUNT(*) FROM product_categories WHERE slug = ? AND id != ?", req.Slug, id).Scan(&count)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "database error"})
			return
		}
		if count > 0 {
			c.JSON(http.StatusConflict, gin.H{"error": "category with this slug already exists"})
			return
		}
	}

	// Строим динамический запрос обновления
	query := "UPDATE product_categories SET updated_at = CURRENT_TIMESTAMP"
	args := []interface{}{}

	if req.NameRu != "" {
		query += ", name_ru = ?"
		args = append(args, req.NameRu)
	}
	if req.NameEn != "" {
		query += ", name_en = ?"
		args = append(args, req.NameEn)
	}
	if req.Slug != "" {
		query += ", slug = ?"
		args = append(args, req.Slug)
	}
	if req.DescriptionRu != "" {
		query += ", description_ru = ?"
		args = append(args, req.DescriptionRu)
	}
	if req.DescriptionEn != "" {
		query += ", description_en = ?"
		args = append(args, req.DescriptionEn)
	}
	if req.SortOrder != nil {
		query += ", sort_order = ?"
		args = append(args, *req.SortOrder)
	}
	if req.IsActive != nil {
		query += ", is_active = ?"
		args = append(args, *req.IsActive)
	}

	query += " WHERE id = ?"
	args = append(args, id)

	_, err = db.Exec(query, args...)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not update category"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "category updated successfully"})
}

// AdminDeleteCategory удаляет категорию
func AdminDeleteCategory(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid category id"})
		return
	}

	// Проверяем, есть ли продукты в этой категории
	var productCount int
	err = db.QueryRow("SELECT COUNT(*) FROM products WHERE category_id = ?", id).Scan(&productCount)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "database error"})
		return
	}
	if productCount > 0 {
		c.JSON(http.StatusConflict, gin.H{
			"error":         "cannot delete category with products",
			"product_count": productCount,
		})
		return
	}

	// Удаляем категорию
	result, err := db.Exec("DELETE FROM product_categories WHERE id = ?", id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not delete category"})
		return
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "category not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "category deleted successfully"})
}
