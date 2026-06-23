package service

import (
	"context"
	"log/slog"

	"github.com/raiyin/artserver/internal/domain"
	"github.com/raiyin/artserver/internal/port"
)

// ChatService implements port.ChatService.
type ChatService struct {
	threadRepo  port.ChatThreadRepository
	messageRepo port.ChatMessageRepository
}

// NewChatService creates a new ChatService.
func NewChatService(threadRepo port.ChatThreadRepository, messageRepo port.ChatMessageRepository) *ChatService {
	return &ChatService{
		threadRepo:  threadRepo,
		messageRepo: messageRepo,
	}
}

// GetThreads retrieves chat threads for a user.
func (s *ChatService) GetThreads(ctx context.Context, userID int64) ([]domain.ChatThread, error) {
	threads, err := s.threadRepo.ListByUser(ctx, userID)
	if err != nil {
		slog.Error("ChatService.GetThreads: failed to list threads",
			"user_id", userID,
			"error", err,
		)
		return nil, err
	}
	slog.Debug("ChatService.GetThreads: threads listed",
		"user_id", userID,
		"count", len(threads),
	)
	return threads, nil
}

// GetAllThreads retrieves all chat threads (admin).
func (s *ChatService) GetAllThreads(ctx context.Context) ([]domain.ChatThread, error) {
	threads, err := s.threadRepo.List(ctx)
	if err != nil {
		slog.Error("ChatService.GetAllThreads: failed to list all threads",
			"error", err,
		)
		return nil, err
	}
	slog.Debug("ChatService.GetAllThreads: threads listed",
		"count", len(threads),
	)
	return threads, nil
}

// CreateThread creates a new chat thread.
func (s *ChatService) CreateThread(ctx context.Context, userID int64, subject string) (*domain.ChatThread, error) {
	thread := &domain.ChatThread{
		UserID:  userID,
		Subject: subject,
		Status:  "open",
	}
	if err := s.threadRepo.Create(ctx, thread); err != nil {
		slog.Error("ChatService.CreateThread: failed to create thread",
			"user_id", userID,
			"subject", subject,
			"error", err,
		)
		return nil, err
	}
	slog.Info("ChatService.CreateThread: thread created",
		"user_id", userID,
		"thread_id", thread.ID,
		"subject", subject,
	)
	return thread, nil
}

// GetMessages retrieves messages for a thread (user).
func (s *ChatService) GetMessages(ctx context.Context, threadID int64, userID int64) ([]domain.ChatMessage, error) {
	thread, err := s.threadRepo.GetByID(ctx, threadID)
	if err != nil {
		slog.Error("ChatService.GetMessages: failed to get thread",
			"thread_id", threadID,
			"user_id", userID,
			"error", err,
		)
		return nil, err
	}
	if thread.UserID != userID {
		slog.Warn("ChatService.GetMessages: forbidden access",
			"thread_id", threadID,
			"user_id", userID,
			"thread_user_id", thread.UserID,
		)
		return nil, domain.ErrForbidden
	}
	messages, err := s.messageRepo.ListByThread(ctx, threadID)
	if err != nil {
		slog.Error("ChatService.GetMessages: failed to list messages",
			"thread_id", threadID,
			"user_id", userID,
			"error", err,
		)
		return nil, err
	}
	slog.Debug("ChatService.GetMessages: messages listed",
		"thread_id", threadID,
		"user_id", userID,
		"count", len(messages),
	)
	return messages, nil
}

// AdminGetMessages retrieves messages for a thread (admin).
func (s *ChatService) AdminGetMessages(ctx context.Context, threadID int64) ([]domain.ChatMessage, error) {
	messages, err := s.messageRepo.ListByThread(ctx, threadID)
	if err != nil {
		slog.Error("ChatService.AdminGetMessages: failed to list messages",
			"thread_id", threadID,
			"error", err,
		)
		return nil, err
	}
	slog.Debug("ChatService.AdminGetMessages: messages listed",
		"thread_id", threadID,
		"count", len(messages),
	)
	return messages, nil
}

// SendMessage sends a message in a thread (user).
func (s *ChatService) SendMessage(ctx context.Context, threadID, userID int64, content string) (*domain.ChatMessage, error) {
	thread, err := s.threadRepo.GetByID(ctx, threadID)
	if err != nil {
		slog.Error("ChatService.SendMessage: failed to get thread",
			"thread_id", threadID,
			"user_id", userID,
			"error", err,
		)
		return nil, err
	}
	if thread.UserID != userID {
		slog.Warn("ChatService.SendMessage: forbidden access",
			"thread_id", threadID,
			"user_id", userID,
			"thread_user_id", thread.UserID,
		)
		return nil, domain.ErrForbidden
	}

	msg := &domain.ChatMessage{
		ThreadID: threadID,
		UserID:   userID,
		Content:  content,
		IsAdmin:  false,
		IsRead:   false,
	}
	if err := s.messageRepo.Create(ctx, msg); err != nil {
		slog.Error("ChatService.SendMessage: failed to create message",
			"thread_id", threadID,
			"user_id", userID,
			"error", err,
		)
		return nil, err
	}
	slog.Info("ChatService.SendMessage: message sent",
		"thread_id", threadID,
		"user_id", userID,
		"message_id", msg.ID,
	)
	return msg, nil
}

// AdminSendMessage sends a message as admin.
func (s *ChatService) AdminSendMessage(ctx context.Context, threadID int64, content string) (*domain.ChatMessage, error) {
	msg := &domain.ChatMessage{
		ThreadID: threadID,
		Content:  content,
		IsAdmin:  true,
		IsRead:   false,
	}
	if err := s.messageRepo.Create(ctx, msg); err != nil {
		slog.Error("ChatService.AdminSendMessage: failed to create message",
			"thread_id", threadID,
			"error", err,
		)
		return nil, err
	}
	slog.Info("ChatService.AdminSendMessage: admin message sent",
		"thread_id", threadID,
		"message_id", msg.ID,
	)
	return msg, nil
}

// MarkMessageAsRead marks a message as read.
func (s *ChatService) MarkMessageAsRead(ctx context.Context, messageID int64) error {
	if err := s.messageRepo.MarkAsRead(ctx, messageID); err != nil {
		slog.Error("ChatService.MarkMessageAsRead: failed to mark message as read",
			"message_id", messageID,
			"error", err,
		)
		return err
	}
	slog.Debug("ChatService.MarkMessageAsRead: message marked as read",
		"message_id", messageID,
	)
	return nil
}

// PollMessages polls for new messages (long-polling).
func (s *ChatService) PollMessages(ctx context.Context, threadID int64, lastMessageID int64, userID int64) ([]domain.ChatMessage, bool, error) {
	thread, err := s.threadRepo.GetByID(ctx, threadID)
	if err != nil {
		slog.Error("ChatService.PollMessages: failed to get thread",
			"thread_id", threadID,
			"user_id", userID,
			"error", err,
		)
		return nil, false, err
	}
	if thread.UserID != userID {
		slog.Warn("ChatService.PollMessages: forbidden access",
			"thread_id", threadID,
			"user_id", userID,
			"thread_user_id", thread.UserID,
		)
		return nil, false, domain.ErrForbidden
	}

	messages, hasMore, err := s.messageRepo.PollNewMessages(ctx, threadID, lastMessageID)
	if err != nil {
		slog.Error("ChatService.PollMessages: failed to poll messages",
			"thread_id", threadID,
			"user_id", userID,
			"last_message_id", lastMessageID,
			"error", err,
		)
		return nil, false, err
	}
	slog.Debug("ChatService.PollMessages: messages polled",
		"thread_id", threadID,
		"user_id", userID,
		"count", len(messages),
		"has_more", hasMore,
	)
	return messages, hasMore, nil
}

// Ensure interface compliance.
var _ port.ChatService = (*ChatService)(nil)
