package handler

import (
	"io"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/raiyin/artserver/internal/dto"
	"github.com/raiyin/artserver/internal/port"
	"github.com/raiyin/artserver/pkg/apperror"
)

// ProfileHandler handles user profile HTTP requests.
type ProfileHandler struct {
	userService port.UserService
}

// NewProfileHandler creates a new ProfileHandler.
func NewProfileHandler(userService port.UserService) *ProfileHandler {
	return &ProfileHandler{userService: userService}
}

// GetProfile returns the authenticated user's profile.
func (h *ProfileHandler) GetProfile(c *gin.Context) {
	userID := c.GetInt64("user_id")

	user, err := h.userService.GetProfile(c.Request.Context(), userID)
	if err != nil {
		apiErr := apperror.FromError(err)
		c.JSON(apiErr.Status, apiErr)
		return
	}

	avatarURL := ""
	if user.AvatarPath != "" {
		avatarURL = "/profile/avatar/" + strconv.FormatInt(user.ID, 10)
	}

	c.JSON(http.StatusOK, dto.ProfileResponse{
		ID:        user.ID,
		Email:     user.Email,
		Name:      user.Name,
		Role:      user.Role,
		AvatarURL: avatarURL,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	})
}

// UpdateProfile updates the authenticated user's profile.
func (h *ProfileHandler) UpdateProfile(c *gin.Context) {
	userID := c.GetInt64("user_id")

	var req dto.UpdateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, apperror.APIError{
			Status:  http.StatusBadRequest,
			Code:    "INVALID_REQUEST",
			Message: "Invalid request body",
		})
		return
	}

	if req.Name == nil {
		c.JSON(http.StatusBadRequest, apperror.APIError{
			Status:  http.StatusBadRequest,
			Code:    "INVALID_REQUEST",
			Message: "No fields to update",
		})
		return
	}

	if err := h.userService.UpdateProfile(c.Request.Context(), userID, *req.Name); err != nil {
		apiErr := apperror.FromError(err)
		c.JSON(apiErr.Status, apiErr)
		return
	}

	// Return updated profile
	user, err := h.userService.GetProfile(c.Request.Context(), userID)
	if err != nil {
		apiErr := apperror.FromError(err)
		c.JSON(apiErr.Status, apiErr)
		return
	}

	avatarURL := ""
	if user.AvatarPath != "" {
		avatarURL = "/profile/avatar/" + strconv.FormatInt(user.ID, 10)
	}

	c.JSON(http.StatusOK, dto.ProfileResponse{
		ID:        user.ID,
		Email:     user.Email,
		Name:      user.Name,
		Role:      user.Role,
		AvatarURL: avatarURL,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	})
}

// UploadAvatar handles avatar file upload.
func (h *ProfileHandler) UploadAvatar(c *gin.Context) {
	userID := c.GetInt64("user_id")

	file, header, err := c.Request.FormFile("avatar")
	if err != nil {
		c.JSON(http.StatusBadRequest, apperror.APIError{
			Status:  http.StatusBadRequest,
			Code:    "INVALID_REQUEST",
			Message: "No avatar file provided",
		})
		return
	}
	defer file.Close()

	avatarPath, err := h.userService.UploadAvatar(c.Request.Context(), userID, header.Filename, file)
	if err != nil {
		apiErr := apperror.FromError(err)
		c.JSON(apiErr.Status, apiErr)
		return
	}

	avatarURL := "/profile/avatar/" + strconv.FormatInt(userID, 10)

	c.JSON(http.StatusOK, dto.UploadAvatarResponse{
		Message:   "avatar uploaded successfully",
		AvatarURL: avatarURL,
	})
	_ = avatarPath
}

// DeleteAvatar removes the user's avatar.
func (h *ProfileHandler) DeleteAvatar(c *gin.Context) {
	userID := c.GetInt64("user_id")

	if err := h.userService.DeleteAvatar(c.Request.Context(), userID); err != nil {
		apiErr := apperror.FromError(err)
		c.JSON(apiErr.Status, apiErr)
		return
	}

	c.JSON(http.StatusOK, dto.DeleteAvatarResponse{
		Message: "avatar deleted successfully",
	})
}

// ServeAvatar serves the avatar file for a given user ID.
func (h *ProfileHandler) ServeAvatar(c *gin.Context) {
	userIDStr := c.Param("id")
	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, apperror.APIError{
			Status:  http.StatusBadRequest,
			Code:    "INVALID_REQUEST",
			Message: "Invalid user ID",
		})
		return
	}

	avatarPath, err := h.userService.ServeAvatar(c.Request.Context(), userID)
	if err != nil {
		apiErr := apperror.FromError(err)
		c.JSON(apiErr.Status, apiErr)
		return
	}

	c.File(avatarPath)
}

// GetUserIDFromContext is a helper to extract user ID from context.
func GetUserIDFromContext(c *gin.Context) int64 {
	userID, exists := c.Get("user_id")
	if !exists {
		return 0
	}
	return userID.(int64)
}

// GetRoleFromContext is a helper to extract role from context.
func GetRoleFromContext(c *gin.Context) string {
	role, exists := c.Get("role")
	if !exists {
		return ""
	}
	return role.(string)
}

// BindJSON is a helper to bind JSON and return error response.
func BindJSON(c *gin.Context, obj interface{}) bool {
	if err := c.ShouldBindJSON(obj); err != nil {
		c.JSON(http.StatusBadRequest, apperror.APIError{
			Status:  http.StatusBadRequest,
			Code:    "INVALID_REQUEST",
			Message: "Invalid request body",
		})
		return false
	}
	return true
}

// BindQuery is a helper to bind query params and return error response.
func BindQuery(c *gin.Context, obj interface{}) bool {
	if err := c.ShouldBindQuery(obj); err != nil {
		c.JSON(http.StatusBadRequest, apperror.APIError{
			Status:  http.StatusBadRequest,
			Code:    "INVALID_REQUEST",
			Message: "Invalid query parameters",
		})
		return false
	}
	return true
}

// ParseInt64Param parses an int64 path parameter.
func ParseInt64Param(c *gin.Context, name string) (int64, bool) {
	valStr := c.Param(name)
	val, err := strconv.ParseInt(valStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, apperror.APIError{
			Status:  http.StatusBadRequest,
			Code:    "INVALID_PARAM",
			Message: "Invalid " + name,
		})
		return 0, false
	}
	return val, true
}

// ParseInt64Query parses an int64 query parameter.
func ParseInt64Query(c *gin.Context, name string) (int64, bool) {
	valStr := c.Query(name)
	if valStr == "" {
		return 0, true
	}
	val, err := strconv.ParseInt(valStr, 10, 64)
	if err != nil {
		return 0, false
	}
	return val, true
}

// ParseIntQuery parses an int query parameter with default.
func ParseIntQuery(c *gin.Context, name string, defaultVal int) int {
	valStr := c.Query(name)
	if valStr == "" {
		return defaultVal
	}
	val, err := strconv.Atoi(valStr)
	if err != nil {
		return defaultVal
	}
	return val
}

// GetUploadedFile extracts a single uploaded file from multipart form.
func GetUploadedFile(c *gin.Context, fieldName string) (string, io.ReadCloser, error) {
	file, header, err := c.Request.FormFile(fieldName)
	if err != nil {
		return "", nil, err
	}
	return header.Filename, file, nil
}
