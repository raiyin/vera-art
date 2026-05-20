package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/raiyin/artserver/models"
)

// SaveConsent сохраняет согласие пользователя на обработку данных
func SaveConsent(c *gin.Context) {
	var req models.ConsentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data"})
		return
	}

	// Получаем ID пользователя из токена (если не передан в запросе)
	userID := req.UserID
	if userID == 0 {
		claims, exists := c.Get("claims")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
			return
		}
		userClaims := claims.(*models.Claims)
		userID = userClaims.UserID
	}

	// Получаем IP-адрес и User-Agent, если не переданы
	ipAddress := req.IPAddress
	if ipAddress == "" {
		ipAddress = c.ClientIP()
	}

	userAgent := req.UserAgent
	if userAgent == "" {
		userAgent = c.Request.UserAgent()
	}

	// Получаем текущую версию политики cookie
	cookieVersion := getSystemSetting(models.SettingCookiePolicyVersion, "1.0")

	// Начинаем транзакцию
	tx, err := db.Begin()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to start transaction"})
		return
	}
	defer tx.Rollback()

	// Сохраняем основное согласие на cookie
	consentID, err := saveUserConsent(tx, userID, models.ConsentTypeCookie, true, cookieVersion, ipAddress, userAgent)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save consent"})
		return
	}

	// Сохраняем детальные настройки cookie
	err = saveCookieConsentDetails(tx, consentID, req.AnalyticsConsent, req.MarketingConsent)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save cookie details"})
		return
	}

	// Обновляем дату последнего согласия в таблице users
	err = updateUserLastConsent(tx, userID, cookieVersion)
	if err != nil {
		log.Printf("Warning: Failed to update user last consent: %v", err)
		// Продолжаем, так как это не критическая ошибка
	}

	// Коммитим транзакцию
	if err := tx.Commit(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to commit transaction"})
		return
	}

	// Возвращаем ответ
	response := models.ConsentResponse{
		Success:          true,
		Message:          "Consent saved successfully",
		ConsentID:        consentID,
		ConsentGiven:     true,
		ConsentVersion:   cookieVersion,
		ConsentDate:      time.Now(),
		AnalyticsAllowed: req.AnalyticsConsent,
		MarketingAllowed: req.MarketingConsent,
	}

	c.JSON(http.StatusOK, response)
}

// GetConsentStatus возвращает текущий статус согласия пользователя
func GetConsentStatus(c *gin.Context) {
	// Получаем ID пользователя из токена
	claims, exists := c.Get("claims")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}
	userClaims := claims.(*models.Claims)
	userID := userClaims.UserID

	// Получаем последнее согласие пользователя
	var consent models.UserConsent
	query := `
		SELECT id, user_id, consent_type, consent_given, consent_version, created_at
		FROM user_consents
		WHERE user_id = ? AND consent_type = ?
		ORDER BY created_at DESC
		LIMIT 1
	`

	err := db.QueryRow(query, userID, models.ConsentTypeCookie).Scan(
		&consent.ID, &consent.UserID, &consent.ConsentType, &consent.ConsentGiven,
		&consent.ConsentVersion, &consent.CreatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			// Согласие не найдено
			c.JSON(http.StatusOK, models.ConsentResponse{
				Success:        true,
				ConsentGiven:   false,
				ConsentVersion: getSystemSetting(models.SettingCookiePolicyVersion, "1.0"),
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get consent status"})
		return
	}

	// Получаем детальные настройки cookie
	var analyticsAllowed, marketingAllowed bool
	detailsQuery := `
		SELECT cookie_category, is_accepted
		FROM cookie_consent_details
		WHERE user_consent_id = ?
	`

	rows, err := db.Query(detailsQuery, consent.ID)
	if err != nil && err != sql.ErrNoRows {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get cookie details"})
		return
	}
	defer rows.Close()

	for rows.Next() {
		var category string
		var accepted bool
		if err := rows.Scan(&category, &accepted); err != nil {
			continue
		}
		switch category {
		case models.CookieCategoryAnalytics:
			analyticsAllowed = accepted
		case models.CookieCategoryMarketing:
			marketingAllowed = accepted
		}
	}

	response := models.ConsentResponse{
		Success:          true,
		ConsentID:        consent.ID,
		ConsentGiven:     consent.ConsentGiven,
		ConsentVersion:   consent.ConsentVersion,
		ConsentDate:      consent.CreatedAt,
		AnalyticsAllowed: analyticsAllowed,
		MarketingAllowed: marketingAllowed,
	}

	c.JSON(http.StatusOK, response)
}

// WithdrawConsent отзывает согласие пользователя
func WithdrawConsent(c *gin.Context) {
	// Получаем ID пользователя из токена
	claims, exists := c.Get("claims")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}
	userClaims := claims.(*models.Claims)
	userID := userClaims.UserID

	// Получаем IP-адрес и User-Agent
	ipAddress := c.ClientIP()
	userAgent := c.Request.UserAgent()
	cookieVersion := getSystemSetting(models.SettingCookiePolicyVersion, "1.0")

	// Начинаем транзакцию
	tx, err := db.Begin()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to start transaction"})
		return
	}
	defer tx.Rollback()

	// Сохраняем отзыв согласия
	consentID, err := saveUserConsent(tx, userID, models.ConsentTypeCookie, false, cookieVersion, ipAddress, userAgent)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save consent withdrawal"})
		return
	}

	// Сохраняем детальные настройки (все отключены)
	err = saveCookieConsentDetails(tx, consentID, false, false)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save cookie details"})
		return
	}

	// Коммитим транзакцию
	if err := tx.Commit(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to commit transaction"})
		return
	}

	c.JSON(http.StatusOK, models.ConsentResponse{
		Success:        true,
		Message:        "Consent withdrawn successfully",
		ConsentID:      consentID,
		ConsentGiven:   false,
		ConsentVersion: cookieVersion,
		ConsentDate:    time.Now(),
	})
}

// GetCookiePolicyVersion возвращает текущую версию политики cookie
func GetCookiePolicyVersion(c *gin.Context) {
	version := getSystemSetting(models.SettingCookiePolicyVersion, "1.0")
	c.JSON(http.StatusOK, gin.H{
		"version":    version,
		"policy_url": "/privacy",
	})
}

// Вспомогательные функции

func saveUserConsent(tx *sql.Tx, userID int, consentType string, consentGiven bool, version, ipAddress, userAgent string) (int, error) {
	query := `
		INSERT INTO user_consents (user_id, consent_type, consent_given, consent_version, ip_address, user_agent, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)
	`

	result, err := tx.Exec(query, userID, consentType, consentGiven, version, ipAddress, userAgent, time.Now(), time.Now())
	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(id), nil
}

func saveCookieConsentDetails(tx *sql.Tx, consentID int, analytics, marketing bool) error {
	// Необходимые cookie всегда активны
	details := []struct {
		category string
		accepted bool
	}{
		{models.CookieCategoryNecessary, true},
		{models.CookieCategoryAnalytics, analytics},
		{models.CookieCategoryMarketing, marketing},
		{models.CookieCategoryPreferences, true}, // Настройки предпочтений
	}

	for _, detail := range details {
		query := `
			INSERT OR REPLACE INTO cookie_consent_details (user_consent_id, cookie_category, is_accepted, created_at)
			VALUES (?, ?, ?, ?)
		`
		_, err := tx.Exec(query, consentID, detail.category, detail.accepted, time.Now())
		if err != nil {
			return err
		}
	}

	return nil
}

func updateUserLastConsent(tx *sql.Tx, userID int, version string) error {
	query := `
		UPDATE users
		SET last_consent_date = ?, cookie_consent_version = ?
		WHERE id = ?
	`
	_, err := tx.Exec(query, time.Now(), version, userID)
	return err
}

func getSystemSetting(key, defaultValue string) string {
	var value string
	query := `SELECT value FROM system_settings WHERE key = ?`
	err := db.QueryRow(query, key).Scan(&value)
	if err != nil {
		if err == sql.ErrNoRows {
			// Создаем настройку со значением по умолчанию
			_, err = db.Exec(`INSERT OR IGNORE INTO system_settings (key, value, created_at, updated_at) VALUES (?, ?, ?, ?)`,
				key, defaultValue, time.Now(), time.Now())
			if err != nil {
				log.Printf("Failed to create system setting %s: %v", key, err)
			}
			return defaultValue
		}
		log.Printf("Failed to get system setting %s: %v", key, err)
		return defaultValue
	}
	return value
}

// GetConsentAuditLog возвращает журнал аудита согласий (только для админов)
func GetConsentAuditLog(c *gin.Context) {
	// Проверяем, что пользователь админ
	claims, exists := c.Get("claims")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}
	userClaims := claims.(*models.Claims)
	if userClaims.Role != "admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Admin access required"})
		return
	}

	// Получаем параметры запроса
	limit := c.DefaultQuery("limit", "100")
	offset := c.DefaultQuery("offset", "0")
	userID := c.Query("user_id")

	// Строим запрос
	query := `
		SELECT cal.id, cal.user_id, cal.action, cal.consent_type, cal.consent_version,
		       cal.ip_address, cal.user_agent, cal.metadata, cal.created_at,
		       u.username, u.email
		FROM consent_audit_log cal
		LEFT JOIN users u ON cal.user_id = u.id
	`

	var args []interface{}
	if userID != "" {
		query += " WHERE cal.user_id = ?"
		args = append(args, userID)
	}

	query += " ORDER BY cal.created_at DESC LIMIT ? OFFSET ?"
	args = append(args, limit, offset)

	rows, err := db.Query(query, args...)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get audit log"})
		return
	}
	defer rows.Close()

	var logs []gin.H
	for rows.Next() {
		var logEntry models.ConsentAuditLog
		var username, email sql.NullString
		var metadata sql.NullString

		err := rows.Scan(
			&logEntry.ID, &logEntry.UserID, &logEntry.Action, &logEntry.ConsentType,
			&logEntry.ConsentVersion, &logEntry.IPAddress, &logEntry.UserAgent,
			&metadata, &logEntry.CreatedAt, &username, &email,
		)
		if err != nil {
			continue
		}

		var metadataMap map[string]interface{}
		if metadata.Valid && metadata.String != "" {
			json.Unmarshal([]byte(metadata.String), &metadataMap)
		}

		logs = append(logs, gin.H{
			"id":              logEntry.ID,
			"user_id":         logEntry.UserID,
			"username":        username.String,
			"email":           email.String,
			"action":          logEntry.Action,
			"consent_type":    logEntry.ConsentType,
			"consent_version": logEntry.ConsentVersion,
			"ip_address":      logEntry.IPAddress,
			"user_agent":      logEntry.UserAgent,
			"metadata":        metadataMap,
			"created_at":      logEntry.CreatedAt,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"logs":   logs,
		"count":  len(logs),
		"limit":  limit,
		"offset": offset,
	})
}
