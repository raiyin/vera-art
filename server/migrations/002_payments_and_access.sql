-- Миграция 002: Платежи, промокоды, отзывы, чат и проверка доступа
-- Фаза 2 образовательной платформы

-- 1. Таблица платежей (интеграция с ЮKassa)
CREATE TABLE IF NOT EXISTS payments (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id INTEGER NOT NULL,
    external_id TEXT UNIQUE, -- ID платежа в ЮKassa
    status TEXT NOT NULL CHECK (status IN ('pending', 'waiting_for_capture', 'succeeded', 'canceled', 'refunded')),
    amount INTEGER NOT NULL, -- в копейках
    currency TEXT DEFAULT 'RUB',
    description TEXT,
    payment_method TEXT,
    metadata JSON, -- дополнительные данные
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

-- Индексы для платежей
CREATE INDEX IF NOT EXISTS idx_payments_user_id ON payments(user_id);
CREATE INDEX IF NOT EXISTS idx_payments_external_id ON payments(external_id);
CREATE INDEX IF NOT EXISTS idx_payments_status ON payments(status);
CREATE INDEX IF NOT EXISTS idx_payments_created_at ON payments(created_at);

-- 2. Таблица промокодов
CREATE TABLE IF NOT EXISTS promo_codes (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    code TEXT UNIQUE NOT NULL,
    discount_type TEXT NOT NULL CHECK (discount_type IN ('percentage', 'fixed')),
    discount_value INTEGER NOT NULL,
    max_uses INTEGER, -- NULL = без ограничений
    used_count INTEGER DEFAULT 0,
    valid_from TIMESTAMP,
    valid_until TIMESTAMP,
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Индексы для промокодов
CREATE INDEX IF NOT EXISTS idx_promo_codes_code ON promo_codes(code);
CREATE INDEX IF NOT EXISTS idx_promo_codes_is_active ON promo_codes(is_active);

-- 3. Таблица отзывов и рейтингов
CREATE TABLE IF NOT EXISTS reviews (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id INTEGER NOT NULL,
    product_id INTEGER NOT NULL,
    purchase_id INTEGER NOT NULL,
    rating INTEGER NOT NULL CHECK (rating >= 1 AND rating <= 5),
    title_ru TEXT,
    title_en TEXT,
    comment_ru TEXT,
    comment_en TEXT,
    is_approved BOOLEAN DEFAULT FALSE, -- модерация админом
    is_visible BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (product_id) REFERENCES products(id) ON DELETE CASCADE,
    FOREIGN KEY (purchase_id) REFERENCES purchases(id) ON DELETE CASCADE
);

-- Индексы для отзывов
CREATE INDEX IF NOT EXISTS idx_reviews_product_id ON reviews(product_id);
CREATE INDEX IF NOT EXISTS idx_reviews_user_id ON reviews(user_id);
CREATE INDEX IF NOT EXISTS idx_reviews_rating ON reviews(rating);
CREATE INDEX IF NOT EXISTS idx_reviews_is_approved ON reviews(is_approved);

-- 4. Таблицы чата с преподавателем
CREATE TABLE IF NOT EXISTS chat_threads (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    purchase_id INTEGER NOT NULL,
    user_id INTEGER NOT NULL,
    admin_id INTEGER, -- ID администратора/преподавателя
    last_message_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    is_resolved BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (purchase_id) REFERENCES purchases(id) ON DELETE CASCADE,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (admin_id) REFERENCES users(id) ON DELETE SET NULL
);

CREATE TABLE IF NOT EXISTS chat_messages (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    thread_id INTEGER NOT NULL,
    sender_id INTEGER NOT NULL,
    message_type TEXT NOT NULL DEFAULT 'text' CHECK (message_type IN ('text', 'image', 'file')),
    content TEXT NOT NULL,
    attachment_url TEXT,
    attachment_size INTEGER,
    is_read BOOLEAN DEFAULT FALSE,
    read_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (thread_id) REFERENCES chat_threads(id) ON DELETE CASCADE,
    FOREIGN KEY (sender_id) REFERENCES users(id) ON DELETE CASCADE
);

-- Индексы для чата
CREATE INDEX IF NOT EXISTS idx_chat_threads_purchase_id ON chat_threads(purchase_id);
CREATE INDEX IF NOT EXISTS idx_chat_threads_user_id ON chat_threads(user_id);
CREATE INDEX IF NOT EXISTS idx_chat_messages_thread_id ON chat_messages(thread_id);
CREATE INDEX IF NOT EXISTS idx_chat_messages_created_at ON chat_messages(thread_id, created_at);
CREATE INDEX IF NOT EXISTS idx_chat_messages_is_read ON chat_messages(is_read) WHERE is_read = FALSE;

-- 5. Обновление таблицы purchases: добавление внешних ключей на payments и promo_codes
-- (Если внешние ключи уже есть, пропускаем)
-- Проверяем существование колонок payment_id и promo_code_id (они уже есть в миграции 001)
-- Добавляем внешние ключи, если они отсутствуют
-- SQLite не поддерживает ADD CONSTRAINT FOREIGN KEY, поэтому нужно пересоздать таблицу.
-- Вместо этого мы просто убедимся, что данные согласованы, и добавим ограничения на уровне приложения.

-- 6. Триггер для автоматического обновления статуса покупки при истечении срока (если не создан)
CREATE TRIGGER IF NOT EXISTS update_purchase_status_after_expiry
AFTER UPDATE OF access_end ON purchases
BEGIN
    UPDATE purchases
    SET status = 'expired'
    WHERE id = NEW.id
    AND NEW.access_end IS NOT NULL
    AND datetime(NEW.access_end) < datetime('now');
END;

-- 7. Добавление тестовых данных (опционально)
-- Промокоды
INSERT OR IGNORE INTO promo_codes (code, discount_type, discount_value, max_uses, valid_from, valid_until, is_active) VALUES
('WELCOME10', 'percentage', 10, 100, datetime('now'), datetime('now', '+1 year'), 1),
('FIRSTBUY', 'fixed', 500, 50, datetime('now'), datetime('now', '+6 months'), 1),
('STUDENT20', 'percentage', 20, NULL, datetime('now'), datetime('now', '+3 months'), 1);

-- Комментарий: Для активации миграции выполните:
-- sqlite3 db/db.sqlite < migrations/002_payments_and_access.sql
