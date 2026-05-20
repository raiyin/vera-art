-- Миграция 003: Таблица для отслеживания согласия на использование файлов cookie и обработку персональных данных
-- В соответствии с требованиями GDPR (ЕС) и Федерального закона № 152-ФЗ (Россия)

-- 1. Создание таблицы для хранения согласий пользователей
CREATE TABLE IF NOT EXISTS user_consents (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id INTEGER NOT NULL,
    consent_type TEXT NOT NULL CHECK (consent_type IN ('cookie', 'privacy', 'terms', 'marketing')),
    consent_given BOOLEAN NOT NULL DEFAULT FALSE,
    consent_version TEXT NOT NULL, -- Версия политики/условий на момент согласия
    ip_address TEXT, -- IP-адрес пользователя для аудита
    user_agent TEXT, -- User agent браузера
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    UNIQUE(user_id, consent_type, consent_version)
);

-- Индексы для быстрого поиска
CREATE INDEX IF NOT EXISTS idx_user_consents_user_id ON user_consents(user_id);
CREATE INDEX IF NOT EXISTS idx_user_consents_consent_type ON user_consents(consent_type);
CREATE INDEX IF NOT EXISTS idx_user_consents_created_at ON user_consents(created_at);

-- 2. Создание таблицы для детальных настроек cookie (если нужна гранулярность)
CREATE TABLE IF NOT EXISTS cookie_consent_details (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_consent_id INTEGER NOT NULL,
    cookie_category TEXT NOT NULL CHECK (cookie_category IN ('necessary', 'analytics', 'marketing', 'preferences')),
    is_accepted BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_consent_id) REFERENCES user_consents(id) ON DELETE CASCADE,
    UNIQUE(user_consent_id, cookie_category)
);

-- Индексы для детальных настроек
CREATE INDEX IF NOT EXISTS idx_cookie_consent_details_user_consent_id ON cookie_consent_details(user_consent_id);
CREATE INDEX IF NOT EXISTS idx_cookie_consent_details_cookie_category ON cookie_consent_details(cookie_category);

-- 3. Создание таблицы для аудита изменений согласий (для соответствия требованиям GDPR/152-ФЗ)
CREATE TABLE IF NOT EXISTS consent_audit_log (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id INTEGER NOT NULL,
    action TEXT NOT NULL CHECK (action IN ('given', 'withdrawn', 'updated')),
    consent_type TEXT NOT NULL,
    consent_version TEXT NOT NULL,
    ip_address TEXT,
    user_agent TEXT,
    metadata JSON, -- Дополнительные метаданные (например, какие именно cookie были изменены)
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

-- Индексы для аудита
CREATE INDEX IF NOT EXISTS idx_consent_audit_log_user_id ON consent_audit_log(user_id);
CREATE INDEX IF NOT EXISTS idx_consent_audit_log_action ON consent_audit_log(action);
CREATE INDEX IF NOT EXISTS idx_consent_audit_log_created_at ON consent_audit_log(created_at);

-- 4. Добавление поля в таблицу users для хранения даты последнего согласия (опционально, для быстрого доступа)
ALTER TABLE users ADD COLUMN last_consent_date TIMESTAMP;
ALTER TABLE users ADD COLUMN cookie_consent_version TEXT;
ALTER TABLE users ADD COLUMN privacy_consent_version TEXT;

-- 5. Создание триггеров для автоматического обновления updated_at
CREATE TRIGGER IF NOT EXISTS update_user_consents_updated_at
AFTER UPDATE ON user_consents
BEGIN
    UPDATE user_consents SET updated_at = CURRENT_TIMESTAMP WHERE id = NEW.id;
END;

-- 6. Создание триггера для автоматического логирования в аудит при изменении согласий
CREATE TRIGGER IF NOT EXISTS log_consent_changes
AFTER INSERT OR UPDATE ON user_consents
BEGIN
    INSERT INTO consent_audit_log (user_id, action, consent_type, consent_version, ip_address, user_agent, created_at)
    VALUES (
        NEW.user_id,
        CASE
            WHEN NEW.consent_given = 1 THEN 'given'
            ELSE 'withdrawn'
        END,
        NEW.consent_type,
        NEW.consent_version,
        NEW.ip_address,
        NEW.user_agent,
        CURRENT_TIMESTAMP
    );
END;

-- 7. Вставка начальных данных (версии политик)
-- Текущая версия политики конфиденциальности и условий использования
INSERT OR IGNORE INTO system_settings (key, value, description) VALUES
    ('privacy_policy_version', '1.0', 'Текущая версия политики конфиденциальности'),
    ('terms_of_service_version', '1.0', 'Текущая версия условий использования'),
    ('cookie_policy_version', '1.0', 'Текущая версия политики использования файлов cookie');

-- Примечание: Таблица system_settings должна существовать. Если её нет, создадим её.
CREATE TABLE IF NOT EXISTS system_settings (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    key TEXT UNIQUE NOT NULL,
    value TEXT NOT NULL,
    description TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Комментарии для документации
COMMENT ON TABLE user_consents IS 'Хранит согласия пользователей на обработку данных в соответствии с GDPR и 152-ФЗ';
COMMENT ON TABLE cookie_consent_details IS 'Детальные настройки согласия на разные категории файлов cookie';
COMMENT ON TABLE consent_audit_log IS 'Аудит всех изменений согласий для соответствия требованиям подотчётности';
