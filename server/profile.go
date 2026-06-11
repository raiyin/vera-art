package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/raiyin/artserver/internal/config"
	"github.com/raiyin/artserver/models"
)

// GetProfile returns the authenticated user's profile data
func GetProfile(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	var user models.User
	query := "SELECT id, username, email, full_name, role, avatar, created_at, updated_at FROM users WHERE id = ?"
	err := db.QueryRow(query, userID).Scan(
		&user.Id,
		&user.Username,
		&user.Email,
		&user.FullName,
		&user.Role,
		&user.Avatar,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}

	avatarURL := ""
	if user.Avatar != "" {
		avatarURL = fmt.Sprintf("/profile/avatar/%d", user.Id)
	}

	c.JSON(http.StatusOK, gin.H{
		"id":         user.Id,
		"username":   user.Username,
		"email":      user.Email,
		"full_name":  user.FullName,
		"role":       user.Role,
		"avatar_url": avatarURL,
		"created_at": user.CreatedAt,
		"updated_at": user.UpdatedAt,
	})
}

// UpdateProfileRequest represents the request body for updating profile
type UpdateProfileRequest struct {
	Email    *string `json:"email"`
	FullName *string `json:"full_name"`
}

// UpdateProfile updates the authenticated user's profile data
func UpdateProfile(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	var req UpdateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	// Build dynamic update query
	updates := []string{}
	args := []interface{}{}

	if req.Email != nil {
		updates = append(updates, "email = ?")
		args = append(args, *req.Email)
	}

	if req.FullName != nil {
		updates = append(updates, "full_name = ?")
		args = append(args, *req.FullName)
	}

	if len(updates) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "no fields to update"})
		return
	}

	// Add updated_at
	updates = append(updates, "updated_at = ?")
	args = append(args, time.Now())

	// Add user ID as last argument for WHERE clause
	args = append(args, userID)

	query := "UPDATE users SET "
	for i, update := range updates {
		if i > 0 {
			query += ", "
		}
		query += update
	}
	query += " WHERE id = ?"

	_, err := db.Exec(query, args...)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not update profile"})
		return
	}

	// Return updated user data
	var user models.User
	selectQuery := "SELECT id, username, email, full_name, role, avatar, created_at, updated_at FROM users WHERE id = ?"
	err = db.QueryRow(selectQuery, userID).Scan(
		&user.Id,
		&user.Username,
		&user.Email,
		&user.FullName,
		&user.Role,
		&user.Avatar,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not fetch updated profile"})
		return
	}

	avatarURL := ""
	if user.Avatar != "" {
		avatarURL = fmt.Sprintf("/profile/avatar/%d", user.Id)
	}

	c.JSON(http.StatusOK, gin.H{
		"id":         user.Id,
		"username":   user.Username,
		"email":      user.Email,
		"full_name":  user.FullName,
		"role":       user.Role,
		"avatar_url": avatarURL,
		"created_at": user.CreatedAt,
		"updated_at": user.UpdatedAt,
	})
}

// UploadAvatar handles avatar file upload
func UploadAvatar(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	// Get the file from the request
	file, header, err := c.Request.FormFile("avatar")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "no avatar file provided"})
		return
	}
	defer file.Close()

	// Validate file size (max 5MB)
	if header.Size > 5*1024*1024 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "file too large: maximum size is 5MB"})
		return
	}

	// Validate file extension
	ext := strings.ToLower(filepath.Ext(header.Filename))
	allowedExts := map[string]bool{
		".jpg":  true,
		".jpeg": true,
		".png":  true,
		".gif":  true,
		".webp": true,
	}
	if !allowedExts[ext] {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid file type: allowed types are jpg, jpeg, png, gif, webp"})
		return
	}

	// Get avatars directory from config
	avatarsDir := config.AppConfigInstance.Directories.AbsAvatarsDir

	// Ensure directory exists
	if err := os.MkdirAll(avatarsDir, 0755); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not create avatars directory"})
		return
	}

	// Delete old avatar if exists
	var oldAvatar string
	err = db.QueryRow("SELECT avatar FROM users WHERE id = ?", userID).Scan(&oldAvatar)
	if err == nil && oldAvatar != "" {
		oldPath := filepath.Join(avatarsDir, oldAvatar)
		if err := os.Remove(oldPath); err != nil && !os.IsNotExist(err) {
			// Log but continue
			log.Printf("Warning: could not delete old avatar: %v", err)
		}
	}

	// Generate unique filename
	filename := fmt.Sprintf("user_%d_%d%s", userID, time.Now().UnixNano(), ext)
	filePath := filepath.Join(avatarsDir, filename)

	// Save the file
	out, err := os.Create(filePath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not save avatar file"})
		return
	}
	defer out.Close()

	if _, err := io.Copy(out, file); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not save avatar file"})
		return
	}

	// Update user record in database
	_, err = db.Exec("UPDATE users SET avatar = ?, updated_at = ? WHERE id = ?", filename, time.Now(), userID)
	if err != nil {
		// Clean up saved file on DB error
		os.Remove(filePath)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not update user avatar"})
		return
	}

	avatarURL := fmt.Sprintf("/profile/avatar/%d", userID)

	c.JSON(http.StatusOK, gin.H{
		"message":    "avatar uploaded successfully",
		"avatar_url": avatarURL,
	})
}

// DeleteAvatar removes the user's avatar
func DeleteAvatar(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	// Get current avatar filename
	var avatar string
	err := db.QueryRow("SELECT avatar FROM users WHERE id = ?", userID).Scan(&avatar)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}

	if avatar == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "no avatar to delete"})
		return
	}

	// Delete the file
	avatarsDir := config.AppConfigInstance.Directories.AbsAvatarsDir
	filePath := filepath.Join(avatarsDir, avatar)
	if err := os.Remove(filePath); err != nil && !os.IsNotExist(err) {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not delete avatar file"})
		return
	}

	// Clear avatar field in database
	_, err = db.Exec("UPDATE users SET avatar = '', updated_at = ? WHERE id = ?", time.Now(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not update user record"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "avatar deleted successfully",
	})
}

// ServeAvatar serves the avatar file for a given user ID
func ServeAvatar(c *gin.Context) {
	userID := c.Param("id")

	var avatar string
	err := db.QueryRow("SELECT avatar FROM users WHERE id = ?", userID).Scan(&avatar)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}

	if avatar == "" {
		c.JSON(http.StatusNotFound, gin.H{"error": "no avatar"})
		return
	}

	avatarsDir := config.AppConfigInstance.Directories.AbsAvatarsDir
	filePath := filepath.Join(avatarsDir, avatar)

	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		c.JSON(http.StatusNotFound, gin.H{"error": "avatar file not found"})
		return
	}

	c.File(filePath)
}
