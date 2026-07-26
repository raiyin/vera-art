package handler

import (
	"log/slog"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/raiyin/artserver/internal/dto"
	"github.com/raiyin/artserver/internal/port"
	"github.com/raiyin/artserver/pkg/apperror"
)

// ChatHandler handles chat HTTP requests.
type ChatHandler struct {
	chatService port.ChatService
}

// NewChatHandler creates a new ChatHandler.
func NewChatHandler(chatService port.ChatService) *ChatHandler {
	return &ChatHandler{chatService: chatService}
}

// GetThreads returns chat threads for the authenticated user (or all threads for admin).
func (h *ChatHandler) GetThreads(c *gin.Context) {
	userID := c.GetInt64("user_id")
	role := GetRoleFromContext(c)

	var threads []interface{}
	var err error

	if role == "admin" {
		adminThreads, adminErr := h.chatService.GetAllThreads(c.Request.Context())
		err = adminErr
		threads = make([]interface{}, len(adminThreads))
		for i, t := range adminThreads {
			threads[i] = t
		}
	} else {
		userThreads, userErr := h.chatService.GetThreads(c.Request.Context(), userID)
		err = userErr
		threads = make([]interface{}, len(userThreads))
		for i, t := range userThreads {
			threads[i] = t
		}
	}

	if err != nil {
		slog.Error("GetThreads: failed to list threads",
			"user_id", userID,
			"role", role,
			"error", err,
		)
		apiErr := apperror.FromError(err)
		c.JSON(apiErr.Status, apiErr)
		return
	}

	c.JSON(http.StatusOK, gin.H{"threads": threads})
}

// CreateThread creates a new chat thread.
func (h *ChatHandler) CreateThread(c *gin.Context) {
	userID := c.GetInt64("user_id")

	var req dto.CreateChatThreadRequest
	if !BindJSON(c, &req) {
		return
	}

	thread, err := h.chatService.CreateThread(c.Request.Context(), userID, strconv.FormatInt(req.PurchaseID, 10))
	if err != nil {
		slog.Error("CreateThread: failed to create thread",
			"user_id", userID,
			"purchase_id", req.PurchaseID,
			"error", err,
		)
		apiErr := apperror.FromError(err)
		c.JSON(apiErr.Status, apiErr)
		return
	}

	slog.Info("Chat thread created",
		"thread_id", thread.ID,
		"user_id", userID,
	)
	c.JSON(http.StatusCreated, dto.CreateChatThreadResponse{
		ID:         thread.ID,
		PurchaseID: req.PurchaseID,
		Message:    "Chat thread created successfully",
	})
}

// GetMessages returns messages for a chat thread (user).
func (h *ChatHandler) GetMessages(c *gin.Context) {
	threadID, ok := ParseInt64Param(c, "thread_id")
	if !ok {
		return
	}

	userID := c.GetInt64("user_id")

	messages, err := h.chatService.GetMessages(c.Request.Context(), threadID, userID)
	if err != nil {
		slog.Error("GetMessages: failed to get messages",
			"thread_id", threadID,
			"user_id", userID,
			"error", err,
		)
		apiErr := apperror.FromError(err)
		c.JSON(apiErr.Status, apiErr)
		return
	}

	responses := make([]dto.ChatMessageResponse, len(messages))
	for i, m := range messages {
		responses[i] = dto.ChatMessageResponse{
			ID:          m.ID,
			ThreadID:    m.ThreadID,
			SenderID:    m.UserID,
			MessageType: "text",
			Content:     m.Content,
			IsRead:      m.IsRead,
			CreatedAt:   m.CreatedAt,
		}
	}

	c.JSON(http.StatusOK, gin.H{"messages": responses})
}

// AdminGetMessages returns messages for a chat thread (admin).
func (h *ChatHandler) AdminGetMessages(c *gin.Context) {
	threadID, ok := ParseInt64Param(c, "thread_id")
	if !ok {
		return
	}

	messages, err := h.chatService.AdminGetMessages(c.Request.Context(), threadID)
	if err != nil {
		slog.Error("AdminGetMessages: failed to get messages",
			"thread_id", threadID,
			"error", err,
		)
		apiErr := apperror.FromError(err)
		c.JSON(apiErr.Status, apiErr)
		return
	}

	responses := make([]dto.ChatMessageResponse, len(messages))
	for i, m := range messages {
		responses[i] = dto.ChatMessageResponse{
			ID:          m.ID,
			ThreadID:    m.ThreadID,
			SenderID:    m.UserID,
			MessageType: "text",
			Content:     m.Content,
			IsRead:      m.IsRead,
			CreatedAt:   m.CreatedAt,
		}
	}

	c.JSON(http.StatusOK, gin.H{"messages": responses})
}

// SendMessage sends a message in a thread (user).
func (h *ChatHandler) SendMessage(c *gin.Context) {
	threadID, ok := ParseInt64Param(c, "thread_id")
	if !ok {
		return
	}

	userID := c.GetInt64("user_id")

	var req dto.SendChatMessageRequest
	if !BindJSON(c, &req) {
		return
	}

	msg, err := h.chatService.SendMessage(c.Request.Context(), threadID, userID, req.Content)
	if err != nil {
		slog.Error("SendMessage: failed to send message",
			"thread_id", threadID,
			"user_id", userID,
			"error", err,
		)
		apiErr := apperror.FromError(err)
		c.JSON(apiErr.Status, apiErr)
		return
	}

	slog.Info("Message sent",
		"message_id", msg.ID,
		"thread_id", threadID,
		"user_id", userID,
	)
	c.JSON(http.StatusCreated, dto.SendChatMessageResponse{
		ID:        msg.ID,
		ThreadID:  msg.ThreadID,
		CreatedAt: msg.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
	})
}

// AdminSendMessage sends a message as admin.
func (h *ChatHandler) AdminSendMessage(c *gin.Context) {
	threadID, ok := ParseInt64Param(c, "thread_id")
	if !ok {
		return
	}

	var req dto.AdminSendChatMessageRequest
	if !BindJSON(c, &req) {
		return
	}

	msg, err := h.chatService.AdminSendMessage(c.Request.Context(), threadID, req.Content)
	if err != nil {
		slog.Error("AdminSendMessage: failed to send message",
			"thread_id", threadID,
			"error", err,
		)
		apiErr := apperror.FromError(err)
		c.JSON(apiErr.Status, apiErr)
		return
	}

	slog.Info("Admin message sent",
		"message_id", msg.ID,
		"thread_id", threadID,
	)
	c.JSON(http.StatusCreated, dto.SendChatMessageResponse{
		ID:        msg.ID,
		ThreadID:  msg.ThreadID,
		CreatedAt: msg.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
	})
}

// MarkMessageAsRead marks a message as read.
func (h *ChatHandler) MarkMessageAsRead(c *gin.Context) {
	messageID, ok := ParseInt64Param(c, "message_id")
	if !ok {
		return
	}

	if err := h.chatService.MarkMessageAsRead(c.Request.Context(), messageID); err != nil {
		slog.Error("MarkMessageAsRead: failed to mark message as read",
			"message_id", messageID,
			"error", err,
		)
		apiErr := apperror.FromError(err)
		c.JSON(apiErr.Status, apiErr)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Message marked as read"})
}

// PollMessages polls for new messages (long-polling).
func (h *ChatHandler) PollMessages(c *gin.Context) {
	threadID, ok := ParseInt64Param(c, "thread_id")
	if !ok {
		return
	}

	userID := c.GetInt64("user_id")

	lastMessageID, _ := ParseInt64Query(c, "last_message_id")

	messages, hasMore, err := h.chatService.PollMessages(c.Request.Context(), threadID, lastMessageID, userID)
	if err != nil {
		slog.Error("PollMessages: failed to poll messages",
			"thread_id", threadID,
			"user_id", userID,
			"error", err,
		)
		apiErr := apperror.FromError(err)
		c.JSON(apiErr.Status, apiErr)
		return
	}

	responses := make([]dto.ChatMessageResponse, len(messages))
	for i, m := range messages {
		responses[i] = dto.ChatMessageResponse{
			ID:          m.ID,
			ThreadID:    m.ThreadID,
			SenderID:    m.UserID,
			MessageType: "text",
			Content:     m.Content,
			IsRead:      m.IsRead,
			CreatedAt:   m.CreatedAt,
		}
	}

	c.JSON(http.StatusOK, dto.PollMessagesResponse{
		Messages: responses,
		HasMore:  hasMore,
	})
}

// ResolveThread resolves a chat thread (admin).
func (h *ChatHandler) ResolveThread(c *gin.Context) {
	id, ok := ParseInt64Param(c, "id")
	if !ok {
		return
	}

	if err := h.chatService.ResolveThread(c.Request.Context(), id); err != nil {
		slog.Error("ResolveThread: failed to resolve thread",
			"thread_id", id,
			"error", err,
		)
		apiErr := apperror.FromError(err)
		c.JSON(apiErr.Status, apiErr)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Thread resolved successfully"})
}

// ReopenThread reopens a chat thread (admin).
func (h *ChatHandler) ReopenThread(c *gin.Context) {
	id, ok := ParseInt64Param(c, "id")
	if !ok {
		return
	}

	if err := h.chatService.ReopenThread(c.Request.Context(), id); err != nil {
		slog.Error("ReopenThread: failed to reopen thread",
			"thread_id", id,
			"error", err,
		)
		apiErr := apperror.FromError(err)
		c.JSON(apiErr.Status, apiErr)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Thread reopened successfully"})
}
