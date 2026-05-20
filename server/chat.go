package main

import (
	"database/sql"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/raiyin/artserver/models"
)

// GetChatThreads возвращает чат-треды текущего пользователя
func GetChatThreads(c *gin.Context) {
	claimsInterface, _ := c.Get("claims")
	claims, ok := claimsInterface.(*models.Claims)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "authentication required"})
		return
	}
	userID := claims.UserID
	role := claims.Role

	var rows *sql.Rows
	var err error

	if role == "admin" {
		// Админ видит все треды
		rows, err = db.Query(`
			SELECT
				ct.id, ct.purchase_id, ct.user_id, ct.admin_id,
				ct.last_message_at, ct.is_resolved, ct.created_at,
				u.username, u.full_name, u.email,
				p.product_id, pr.title_ru, pr.title_en
			FROM chat_threads ct
			JOIN users u ON ct.user_id = u.id
			JOIN purchases p ON ct.purchase_id = p.id
			JOIN products pr ON p.product_id = pr.id
			ORDER BY ct.last_message_at DESC
		`)
	} else {
		// Пользователь видит только свои треды
		rows, err = db.Query(`
			SELECT
				ct.id, ct.purchase_id, ct.user_id, ct.admin_id,
				ct.last_message_at, ct.is_resolved, ct.created_at,
				u.username, u.full_name, u.email,
				p.product_id, pr.title_ru, pr.title_en
			FROM chat_threads ct
			JOIN users u ON ct.user_id = u.id
			JOIN purchases p ON ct.purchase_id = p.id
			JOIN products pr ON p.product_id = pr.id
			WHERE ct.user_id = ?
			ORDER BY ct.last_message_at DESC
		`, userID)
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "database error"})
		return
	}
	defer rows.Close()

	threads := []models.ChatThread{}
	for rows.Next() {
		var ct models.ChatThread
		var adminID sql.NullInt64
		var username, fullName, email, titleRu, titleEn string
		var productID int
		err := rows.Scan(
			&ct.Id, &ct.PurchaseId, &ct.UserId, &adminID,
			&ct.LastMessageAt, &ct.IsResolved, &ct.CreatedAt,
			&username, &fullName, &email,
			&productID, &titleRu, &titleEn,
		)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "scan error"})
			return
		}
		if adminID.Valid {
			ct.AdminId = adminID
		}
		ct.User = &models.User{Username: username, FullName: fullName, Email: email}
		ct.Purchase = &models.Purchase{Id: ct.PurchaseId, ProductId: productID}
		if ct.Purchase.Product == nil {
			ct.Purchase.Product = &models.Product{TitleRu: titleRu, TitleEn: titleEn}
		}
		threads = append(threads, ct)
	}

	c.JSON(http.StatusOK, threads)
}

// CreateChatThread создает новый чат-тред (пользователь)
func CreateChatThread(c *gin.Context) {
	claimsInterface, _ := c.Get("claims")
	claims, ok := claimsInterface.(*models.Claims)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "authentication required"})
		return
	}
	userID := claims.UserID

	var req struct {
		PurchaseID int `json:"purchase_id" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	// Проверяем, принадлежит ли покупка пользователю и активна ли она
	var exists bool
	err := db.QueryRow(`
		SELECT EXISTS(
			SELECT 1 FROM purchases
			WHERE id = ? AND user_id = ? AND status = 'active'
			AND (access_end IS NULL OR access_end > CURRENT_TIMESTAMP)
		)
	`, req.PurchaseID, userID).Scan(&exists)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "database error"})
		return
	}
	if !exists {
		c.JSON(http.StatusForbidden, gin.H{"error": "purchase not found or inactive"})
		return
	}

	// Проверяем, не существует ли уже тред для этой покупки
	var threadID int
	err = db.QueryRow("SELECT id FROM chat_threads WHERE purchase_id = ?", req.PurchaseID).Scan(&threadID)
	if err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "chat thread already exists for this purchase"})
		return
	} else if err != sql.ErrNoRows {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "database error"})
		return
	}

	// Создаем тред
	result, err := db.Exec(`
		INSERT INTO chat_threads (purchase_id, user_id, last_message_at, is_resolved, created_at)
		VALUES (?, ?, CURRENT_TIMESTAMP, FALSE, CURRENT_TIMESTAMP)
	`, req.PurchaseID, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create thread"})
		return
	}

	newID, _ := result.LastInsertId()
	c.JSON(http.StatusCreated, gin.H{
		"id":          newID,
		"purchase_id": req.PurchaseID,
		"message":     "chat thread created",
	})
}

// GetChatMessages возвращает сообщения треда
func GetChatMessages(c *gin.Context) {
	threadID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid thread id"})
		return
	}

	claimsInterface, _ := c.Get("claims")
	claims, ok := claimsInterface.(*models.Claims)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "authentication required"})
		return
	}
	userID := claims.UserID
	role := claims.Role

	// Проверяем доступ к треду
	var threadUserID, threadAdminID sql.NullInt64
	err = db.QueryRow("SELECT user_id, admin_id FROM chat_threads WHERE id = ?", threadID).Scan(&threadUserID, &threadAdminID)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "thread not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "database error"})
		return
	}
	// Пользователь может видеть тред, если он автор или админ, или если он назначенный админ
	if role != "admin" && threadUserID.Int64 != int64(userID) && (!threadAdminID.Valid || threadAdminID.Int64 != int64(userID)) {
		c.JSON(http.StatusForbidden, gin.H{"error": "access denied"})
		return
	}

	// Получаем сообщения
	rows, err := db.Query(`
		SELECT
			cm.id, cm.thread_id, cm.sender_id, cm.message_type,
			cm.content, cm.attachment_url, cm.attachment_size,
			cm.is_read, cm.read_at, cm.created_at,
			u.username, u.full_name
		FROM chat_messages cm
		JOIN users u ON cm.sender_id = u.id
		WHERE cm.thread_id = ?
		ORDER BY cm.created_at ASC
	`, threadID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "database error"})
		return
	}
	defer rows.Close()

	messages := []models.ChatMessage{}
	for rows.Next() {
		var cm models.ChatMessage
		var readAt sql.NullTime
		var attachmentURL, attachmentSize sql.NullString
		var username, fullName string
		err := rows.Scan(
			&cm.Id, &cm.ThreadId, &cm.SenderId, &cm.MessageType,
			&cm.Content, &attachmentURL, &attachmentSize,
			&cm.IsRead, &readAt, &cm.CreatedAt,
			&username, &fullName,
		)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "scan error"})
			return
		}
		cm.AttachmentUrl = attachmentURL
		if attachmentSize.Valid {
			size, _ := strconv.ParseInt(attachmentSize.String, 10, 64)
			cm.AttachmentSize = sql.NullInt64{Int64: size, Valid: true}
		} else {
			cm.AttachmentSize = sql.NullInt64{Valid: false}
		}
		cm.ReadAt = readAt
		cm.Sender = &models.User{Username: username, FullName: fullName}
		messages = append(messages, cm)
	}

	c.JSON(http.StatusOK, messages)
}

// SendChatMessage отправляет сообщение в тред
func SendChatMessage(c *gin.Context) {
	threadID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid thread id"})
		return
	}

	claimsInterface, _ := c.Get("claims")
	claims, ok := claimsInterface.(*models.Claims)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "authentication required"})
		return
	}
	userID := claims.UserID
	role := claims.Role

	// Проверяем доступ к треду
	var threadUserID, threadAdminID sql.NullInt64
	var isResolved bool
	err = db.QueryRow("SELECT user_id, admin_id, is_resolved FROM chat_threads WHERE id = ?", threadID).
		Scan(&threadUserID, &threadAdminID, &isResolved)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "thread not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "database error"})
		return
	}
	if isResolved {
		c.JSON(http.StatusConflict, gin.H{"error": "thread is resolved, cannot send messages"})
		return
	}
	// Пользователь может отправлять сообщения, если он автор, админ или назначенный админ
	if role != "admin" && threadUserID.Int64 != int64(userID) && (!threadAdminID.Valid || threadAdminID.Int64 != int64(userID)) {
		c.JSON(http.StatusForbidden, gin.H{"error": "access denied"})
		return
	}

	var req struct {
		Content       string `json:"content" binding:"required"`
		MessageType   string `json:"message_type" binding:"required,oneof=text image file"`
		AttachmentUrl string `json:"attachment_url,omitempty"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	// Вставляем сообщение
	result, err := db.Exec(`
		INSERT INTO chat_messages (thread_id, sender_id, message_type, content, attachment_url, created_at)
		VALUES (?, ?, ?, ?, ?, CURRENT_TIMESTAMP)
	`, threadID, userID, req.MessageType, req.Content, req.AttachmentUrl)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to send message"})
		return
	}

	// Обновляем last_message_at в треде
	_, err = db.Exec("UPDATE chat_threads SET last_message_at = CURRENT_TIMESTAMP WHERE id = ?", threadID)
	if err != nil {
		// Логируем, но не прерываем
		log.Println("Failed to update last_message_at:", err)
	}

	messageID, _ := result.LastInsertId()
	c.JSON(http.StatusCreated, gin.H{
		"id":         messageID,
		"thread_id":  threadID,
		"created_at": time.Now().Format(time.RFC3339),
	})
}

// MarkMessageAsRead отмечает сообщение как прочитанное
func MarkMessageAsRead(c *gin.Context) {
	messageID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid message id"})
		return
	}

	claimsInterface, _ := c.Get("claims")
	claims, ok := claimsInterface.(*models.Claims)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "authentication required"})
		return
	}
	userID := claims.UserID

	// Проверяем, принадлежит ли сообщение треду, к которому у пользователя есть доступ
	var threadID int
	err = db.QueryRow("SELECT thread_id FROM chat_messages WHERE id = ?", messageID).Scan(&threadID)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "message not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "database error"})
		return
	}

	// Проверяем доступ к треду (аналогично SendChatMessage)
	var threadUserID, threadAdminID sql.NullInt64
	err = db.QueryRow("SELECT user_id, admin_id FROM chat_threads WHERE id = ?", threadID).
		Scan(&threadUserID, &threadAdminID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "database error"})
		return
	}
	if claims.Role != "admin" && threadUserID.Int64 != int64(userID) && (!threadAdminID.Valid || threadAdminID.Int64 != int64(userID)) {
		c.JSON(http.StatusForbidden, gin.H{"error": "access denied"})
		return
	}

	// Отмечаем как прочитанное
	_, err = db.Exec("UPDATE chat_messages SET is_read = TRUE, read_at = CURRENT_TIMESTAMP WHERE id = ?", messageID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to mark as read"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "message marked as read"})
}

// GetAdminChatThreads возвращает все активные треды (админ)
func GetAdminChatThreads(c *gin.Context) {
	// Уже реализовано в GetChatThreads с проверкой роли, но для админского маршрута можно добавить фильтры
	// Просто вызываем GetChatThreads, middleware уже проверил роль admin
	GetChatThreads(c)
}

// ResolveChatThread закрывает тред (админ)
func ResolveChatThread(c *gin.Context) {
	threadID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid thread id"})
		return
	}

	// Обновляем тред
	result, err := db.Exec("UPDATE chat_threads SET is_resolved = TRUE WHERE id = ?", threadID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "database error"})
		return
	}
	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "thread not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "thread resolved"})
}

// PollChatMessages long-polling эндпоинт для получения новых сообщений
func PollChatMessages(c *gin.Context) {
	claimsInterface, _ := c.Get("claims")
	claims, ok := claimsInterface.(*models.Claims)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "authentication required"})
		return
	}
	userID := claims.UserID
	role := claims.Role

	// Параметры запроса
	threadIDStr := c.Query("thread_id")
	lastMessageIDStr := c.Query("last_message_id")
	timeoutStr := c.Query("timeout")
	if timeoutStr == "" {
		timeoutStr = "25" // секунд по умолчанию
	}
	timeoutSec, err := strconv.Atoi(timeoutStr)
	if err != nil || timeoutSec < 1 || timeoutSec > 60 {
		timeoutSec = 25
	}

	// Если передан thread_id, проверяем доступ к треду
	var threadID int
	if threadIDStr != "" {
		threadID, err = strconv.Atoi(threadIDStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid thread_id"})
			return
		}
		var threadUserID, threadAdminID sql.NullInt64
		err = db.QueryRow("SELECT user_id, admin_id FROM chat_threads WHERE id = ?", threadID).
			Scan(&threadUserID, &threadAdminID)
		if err != nil {
			if err == sql.ErrNoRows {
				c.JSON(http.StatusNotFound, gin.H{"error": "thread not found"})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{"error": "database error"})
			return
		}
		if role != "admin" && threadUserID.Int64 != int64(userID) && (!threadAdminID.Valid || threadAdminID.Int64 != int64(userID)) {
			c.JSON(http.StatusForbidden, gin.H{"error": "access denied"})
			return
		}
	}

	lastMessageID := 0
	if lastMessageIDStr != "" {
		lastMessageID, _ = strconv.Atoi(lastMessageIDStr)
	}

	// Функция проверки новых сообщений
	checkNewMessages := func() ([]models.ChatMessage, bool, error) {
		var rows *sql.Rows
		var err error
		if threadIDStr != "" {
			// Проверяем новые сообщения в конкретном треде
			rows, err = db.Query(`
				SELECT
					cm.id, cm.thread_id, cm.sender_id, cm.message_type,
					cm.content, cm.attachment_url, cm.attachment_size,
					cm.is_read, cm.read_at, cm.created_at,
					u.username, u.full_name
				FROM chat_messages cm
				JOIN users u ON cm.sender_id = u.id
				WHERE cm.thread_id = ? AND cm.id > ?
				ORDER BY cm.created_at ASC
			`, threadID, lastMessageID)
		} else {
			// Проверяем новые сообщения во всех тредах пользователя
			rows, err = db.Query(`
				SELECT
					cm.id, cm.thread_id, cm.sender_id, cm.message_type,
					cm.content, cm.attachment_url, cm.attachment_size,
					cm.is_read, cm.read_at, cm.created_at,
					u.username, u.full_name
				FROM chat_messages cm
				JOIN users u ON cm.sender_id = u.id
				JOIN chat_threads ct ON cm.thread_id = ct.id
				WHERE (ct.user_id = ? OR ct.admin_id = ? OR ? = 'admin')
					AND cm.id > ?
				ORDER BY cm.created_at ASC
			`, userID, userID, role, lastMessageID)
		}
		if err != nil {
			return nil, false, err
		}
		defer rows.Close()

		messages := []models.ChatMessage{}
		for rows.Next() {
			var cm models.ChatMessage
			var readAt sql.NullTime
			var attachmentURL, attachmentSize sql.NullString
			var username, fullName string
			err := rows.Scan(
				&cm.Id, &cm.ThreadId, &cm.SenderId, &cm.MessageType,
				&cm.Content, &attachmentURL, &attachmentSize,
				&cm.IsRead, &readAt, &cm.CreatedAt,
				&username, &fullName,
			)
			if err != nil {
				return nil, false, err
			}
			if attachmentURL.Valid {
				cm.AttachmentUrl = sql.NullString{String: attachmentURL.String, Valid: true}
			} else {
				cm.AttachmentUrl = sql.NullString{Valid: false}
			}
			if attachmentSize.Valid {
				size, _ := strconv.ParseInt(attachmentSize.String, 10, 64)
				cm.AttachmentSize = sql.NullInt64{Int64: size, Valid: true}
			} else {
				cm.AttachmentSize = sql.NullInt64{Valid: false}
			}
			if readAt.Valid {
				cm.ReadAt = sql.NullTime{Time: readAt.Time, Valid: true}
			} else {
				cm.ReadAt = sql.NullTime{Valid: false}
			}
			cm.Sender = &models.User{Username: username, FullName: fullName}
			messages = append(messages, cm)
		}
		return messages, len(messages) > 0, nil
	}

	// Первая проверка
	messages, hasNew, err := checkNewMessages()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "database error"})
		return
	}
	if hasNew {
		c.JSON(http.StatusOK, gin.H{"messages": messages, "has_more": false})
		return
	}

	// Long-polling: ждем новые сообщения с интервалами
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()
	timeout := time.After(time.Duration(timeoutSec) * time.Second)

	for {
		select {
		case <-ticker.C:
			messages, hasNew, err := checkNewMessages()
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "database error"})
				return
			}
			if hasNew {
				c.JSON(http.StatusOK, gin.H{"messages": messages, "has_more": false})
				return
			}
		case <-timeout:
			// Таймаут - возвращаем пустой ответ
			c.JSON(http.StatusOK, gin.H{"messages": []models.ChatMessage{}, "has_more": false})
			return
		case <-c.Request.Context().Done():
			// Клиент отключился
			return
		}
	}
}
