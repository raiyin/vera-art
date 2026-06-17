package service

import (
	"context"

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
	return s.threadRepo.ListByUser(ctx, userID)
}

// GetAllThreads retrieves all chat threads (admin).
func (s *ChatService) GetAllThreads(ctx context.Context) ([]domain.ChatThread, error) {
	return s.threadRepo.List(ctx)
}

// CreateThread creates a new chat thread.
func (s *ChatService) CreateThread(ctx context.Context, userID int64, subject string) (*domain.ChatThread, error) {
	thread := &domain.ChatThread{
		UserID:  userID,
		Subject: subject,
		Status:  "open",
	}
	if err := s.threadRepo.Create(ctx, thread); err != nil {
		return nil, err
	}
	return thread, nil
}

// GetMessages retrieves messages for a thread (user).
func (s *ChatService) GetMessages(ctx context.Context, threadID int64, userID int64) ([]domain.ChatMessage, error) {
	thread, err := s.threadRepo.GetByID(ctx, threadID)
	if err != nil {
		return nil, err
	}
	if thread.UserID != userID {
		return nil, domain.ErrForbidden
	}
	return s.messageRepo.ListByThread(ctx, threadID)
}

// AdminGetMessages retrieves messages for a thread (admin).
func (s *ChatService) AdminGetMessages(ctx context.Context, threadID int64) ([]domain.ChatMessage, error) {
	return s.messageRepo.ListByThread(ctx, threadID)
}

// SendMessage sends a message in a thread (user).
func (s *ChatService) SendMessage(ctx context.Context, threadID, userID int64, content string) (*domain.ChatMessage, error) {
	thread, err := s.threadRepo.GetByID(ctx, threadID)
	if err != nil {
		return nil, err
	}
	if thread.UserID != userID {
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
		return nil, err
	}
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
		return nil, err
	}
	return msg, nil
}

// MarkMessageAsRead marks a message as read.
func (s *ChatService) MarkMessageAsRead(ctx context.Context, messageID int64) error {
	return s.messageRepo.MarkAsRead(ctx, messageID)
}

// PollMessages polls for new messages (long-polling).
func (s *ChatService) PollMessages(ctx context.Context, threadID int64, lastMessageID int64, userID int64) ([]domain.ChatMessage, bool, error) {
	thread, err := s.threadRepo.GetByID(ctx, threadID)
	if err != nil {
		return nil, false, err
	}
	if thread.UserID != userID {
		return nil, false, domain.ErrForbidden
	}

	return s.messageRepo.PollNewMessages(ctx, threadID, lastMessageID)
}

// Ensure interface compliance.
var _ port.ChatService = (*ChatService)(nil)
