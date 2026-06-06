package main

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/raiyin/artserver/models"
)

// AdminGetTags возвращает список всех тегов для админки
func AdminGetTags(c *gin.Context) {
	query := `SELECT id, name_ru, name_en, slug, created_at
		FROM tags
		ORDER BY name_ru`

	rows, err := db.Query(query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "database error"})
		return
	}
	defer rows.Close()

	tags := []models.Tag{}
	for rows.Next() {
		var tag models.Tag
		err := rows.Scan(
			&tag.Id, &tag.NameRu, &tag.NameEn,
			&tag.Slug, &tag.CreatedAt,
		)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "scan error"})
			return
		}
		tags = append(tags, tag)
	}

	c.JSON(http.StatusOK, tags)
}

// AdminCreateTag создает новый тег
func AdminCreateTag(c *gin.Context) {
	var req struct {
		NameRu string `json:"name_ru" binding:"required"`
		NameEn string `json:"name_en" binding:"required"`
		Slug   string `json:"slug" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Проверяем, существует ли тег с таким slug
	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM tags WHERE slug = ?", req.Slug).Scan(&count)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "database error"})
		return
	}
	if count > 0 {
		c.JSON(http.StatusConflict, gin.H{"error": "tag with this slug already exists"})
		return
	}

	result, err := db.Exec(
		"INSERT INTO tags (name_ru, name_en, slug) VALUES (?, ?, ?)",
		req.NameRu, req.NameEn, req.Slug,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not create tag"})
		return
	}

	id, _ := result.LastInsertId()
	c.JSON(http.StatusCreated, gin.H{
		"id":      id,
		"message": "tag created successfully",
	})
}

// AdminUpdateTag обновляет тег
func AdminUpdateTag(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid tag id"})
		return
	}

	var req struct {
		NameRu string `json:"name_ru,omitempty"`
		NameEn string `json:"name_en,omitempty"`
		Slug   string `json:"slug,omitempty"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Проверяем существование тега
	var exists bool
	err = db.QueryRow("SELECT EXISTS(SELECT 1 FROM tags WHERE id = ?)", id).Scan(&exists)
	if err != nil || !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "tag not found"})
		return
	}

	// Если обновляется slug, проверяем уникальность
	if req.Slug != "" {
		var count int
		err = db.QueryRow("SELECT COUNT(*) FROM tags WHERE slug = ? AND id != ?", req.Slug, id).Scan(&count)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "database error"})
			return
		}
		if count > 0 {
			c.JSON(http.StatusConflict, gin.H{"error": "tag with this slug already exists"})
			return
		}
	}

	// Строим динамический запрос обновления
	query := "UPDATE tags SET"
	args := []interface{}{}
	parts := []string{}

	if req.NameRu != "" {
		parts = append(parts, " name_ru = ?")
		args = append(args, req.NameRu)
	}
	if req.NameEn != "" {
		parts = append(parts, " name_en = ?")
		args = append(args, req.NameEn)
	}
	if req.Slug != "" {
		parts = append(parts, " slug = ?")
		args = append(args, req.Slug)
	}

	if len(parts) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "no fields to update"})
		return
	}

	for i, part := range parts {
		if i == 0 {
			query += part
		} else {
			query += "," + part
		}
	}

	query += " WHERE id = ?"
	args = append(args, id)

	_, err = db.Exec(query, args...)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not update tag"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "tag updated successfully"})
}

// AdminDeleteTag удаляет тег
func AdminDeleteTag(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid tag id"})
		return
	}

	// Проверяем существование тега
	var exists bool
	err = db.QueryRow("SELECT EXISTS(SELECT 1 FROM tags WHERE id = ?)", id).Scan(&exists)
	if err != nil || !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "tag not found"})
		return
	}

	// Удаляем связи product_tags (каскадно), затем сам тег
	_, err = db.Exec("DELETE FROM product_tags WHERE tag_id = ?", id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not delete tag relations"})
		return
	}

	result, err := db.Exec("DELETE FROM tags WHERE id = ?", id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not delete tag"})
		return
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "tag not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "tag deleted successfully"})
}
