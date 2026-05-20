package models

import (
	"time"
)

// UserConsent представляет согласие пользователя на обработку данных
type UserConsent struct {
	ID             int       `json:"id"`
	UserID         int       `json:"user_id"`
	ConsentType    string    `json:"consent_type"` // cookie, privacy, terms, marketing
	ConsentGiven   bool      `json:"consent_given"`
	ConsentVersion string    `json:"consent_version"`
	IPAddress      string    `json:"ip_address,omitempty"`
	UserAgent      string    `json:"user_agent,omitempty"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

// CookieConsentDetail представляет детальные настройки согласия на файлы cookie
type CookieConsentDetail struct {
	ID             int       `json:"id"`
	UserConsentID  int       `json:"user_consent_id"`
	CookieCategory string    `json:"cookie_category"` // necessary, analytics, marketing, preferences
	IsAccepted     bool      `json:"is_accepted"`
	CreatedAt      time.Time `json:"created_at"`
}

// ConsentAuditLog представляет запись аудита изменений согласий
type ConsentAuditLog struct {
	ID             int       `json:"id"`
	UserID         int       `json:"user_id"`
	Action         string    `json:"action"` // given, withdrawn, updated
	ConsentType    string    `json:"consent_type"`
	ConsentVersion string    `json:"consent_version"`
	IPAddress      string    `json:"ip_address,omitempty"`
	UserAgent      string    `json:"user_agent,omitempty"`
	Metadata       string    `json:"metadata,omitempty"` // JSON-строка с дополнительными метаданными
	CreatedAt      time.Time `json:"created_at"`
}

// SystemSetting представляет системные настройки (версии политик)
type SystemSetting struct {
	ID          int       `json:"id"`
	Key         string    `json:"key"`
	Value       string    `json:"value"`
	Description string    `json:"description,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// ConsentRequest представляет запрос на сохранение согласия от фронтенда
type ConsentRequest struct {
	UserID           int    `json:"user_id"`
	AnalyticsConsent bool   `json:"analytics_consent"`
	MarketingConsent bool   `json:"marketing_consent"`
	IPAddress        string `json:"ip_address,omitempty"`
	UserAgent        string `json:"user_agent,omitempty"`
}

// ConsentResponse представляет ответ с информацией о согласии
type ConsentResponse struct {
	Success          bool      `json:"success"`
	Message          string    `json:"message,omitempty"`
	ConsentID        int       `json:"consent_id,omitempty"`
	ConsentGiven     bool      `json:"consent_given"`
	ConsentVersion   string    `json:"consent_version"`
	ConsentDate      time.Time `json:"consent_date,omitempty"`
	AnalyticsAllowed bool      `json:"analytics_allowed"`
	MarketingAllowed bool      `json:"marketing_allowed"`
}

// CookieConsentStatus представляет текущий статус согласия на файлы cookie
type CookieConsentStatus struct {
	Necessary   bool `json:"necessary"`
	Analytics   bool `json:"analytics"`
	Marketing   bool `json:"marketing"`
	Preferences bool `json:"preferences"`
}

// ConsentTypes предопределенные типы согласий
const (
	ConsentTypeCookie    = "cookie"
	ConsentTypePrivacy   = "privacy"
	ConsentTypeTerms     = "terms"
	ConsentTypeMarketing = "marketing"
)

// CookieCategories предопределенные категории файлов cookie
const (
	CookieCategoryNecessary   = "necessary"
	CookieCategoryAnalytics   = "analytics"
	CookieCategoryMarketing   = "marketing"
	CookieCategoryPreferences = "preferences"
)

// ConsentActions действия с согласием
const (
	ConsentActionGiven     = "given"
	ConsentActionWithdrawn = "withdrawn"
	ConsentActionUpdated   = "updated"
)

// SystemSettingKeys ключи системных настроек
const (
	SettingPrivacyPolicyVersion = "privacy_policy_version"
	SettingTermsVersion         = "terms_of_service_version"
	SettingCookiePolicyVersion  = "cookie_policy_version"
)
