-- Миграция 001: Создание образовательной платформы
-- Добавление ролей пользователям и создание таблиц для курсов/мастер-классов

-- 1. Расширение таблицы users для RBAC
-- Проверяем существование колонок перед добавлением
-- role, email, full_name уже добавлены в предыдущих миграциях

-- Добавляем created_at и updated_at без DEFAULT (SQLite ограничение)
ALTER TABLE users ADD COLUMN created_at TIMESTAMP;
ALTER TABLE users ADD COLUMN updated_at TIMESTAMP;

-- Устанавливаем значения по умолчанию для существующих записей
UPDATE users SET created_at = datetime('now'), updated_at = datetime('now') WHERE created_at IS NULL OR created_at = '';

-- Индексы для пользователей (создаем если не существуют)
CREATE INDEX IF NOT EXISTS idx_users_role ON users(role);
CREATE INDEX IF NOT EXISTS idx_users_email ON users(email);

-- 2. Создание таблицы категорий продуктов
CREATE TABLE IF NOT EXISTS product_categories (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name_ru TEXT NOT NULL,
    name_en TEXT NOT NULL,
    slug TEXT UNIQUE NOT NULL,
    description_ru TEXT,
    description_en TEXT,
    sort_order INTEGER DEFAULT 0,
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Индексы для категорий
CREATE INDEX IF NOT EXISTS idx_product_categories_slug ON product_categories(slug);
CREATE INDEX IF NOT EXISTS idx_product_categories_is_active ON product_categories(is_active);

-- 3. Создание таблицы продуктов (курсы и мастер-классы)
CREATE TABLE IF NOT EXISTS products (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    type TEXT NOT NULL CHECK (type IN ('course', 'masterclass')),
    title_ru TEXT NOT NULL,
    title_en TEXT NOT NULL,
    description_ru TEXT,
    description_en TEXT,
    short_description_ru TEXT,
    short_description_en TEXT,
    price INTEGER NOT NULL, -- в копейках
    duration_days INTEGER, -- срок доступа в днях (NULL = бессрочно)
    thumbnail_url TEXT,
    video_url TEXT,
    status TEXT NOT NULL DEFAULT 'draft' CHECK (status IN ('draft', 'published', 'archived')),
    difficulty TEXT CHECK (difficulty IN ('beginner', 'intermediate', 'advanced')),
    total_lessons INTEGER DEFAULT 0,
    total_duration_minutes INTEGER DEFAULT 0,
    category_id INTEGER,
    instructor_id INTEGER,
    tags JSON,
    prerequisites_ru TEXT,
    prerequisites_en TEXT,
    learning_outcomes_ru TEXT,
    learning_outcomes_en TEXT,
    certificate_available BOOLEAN DEFAULT FALSE,
    max_students INTEGER,
    start_date TIMESTAMP,
    language TEXT DEFAULT 'ru' CHECK (language IN ('ru', 'en', 'both')),
    is_featured BOOLEAN DEFAULT FALSE,
    view_count INTEGER DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (category_id) REFERENCES product_categories(id) ON DELETE SET NULL,
    FOREIGN KEY (instructor_id) REFERENCES users(id) ON DELETE SET NULL
);

-- Индексы для продуктов
CREATE INDEX IF NOT EXISTS idx_products_type ON products(type);
CREATE INDEX IF NOT EXISTS idx_products_status ON products(status);
CREATE INDEX IF NOT EXISTS idx_products_difficulty ON products(difficulty);
CREATE INDEX IF NOT EXISTS idx_products_category_id ON products(category_id);
CREATE INDEX IF NOT EXISTS idx_products_instructor_id ON products(instructor_id);
CREATE INDEX IF NOT EXISTS idx_products_is_featured ON products(is_featured);
CREATE INDEX IF NOT EXISTS idx_products_created_at ON products(created_at);

-- 4. Создание таблицы уроков/модулей
CREATE TABLE IF NOT EXISTS lessons (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    product_id INTEGER NOT NULL,
    title_ru TEXT NOT NULL,
    title_en TEXT NOT NULL,
    description_ru TEXT,
    description_en TEXT,
    content_type TEXT NOT NULL CHECK (content_type IN ('video', 'text', 'pdf', 'quiz', 'assignment')),
    content_url TEXT,
    duration_minutes INTEGER DEFAULT 0,
    sort_order INTEGER NOT NULL DEFAULT 0,
    is_preview BOOLEAN DEFAULT FALSE,
    resources JSON,
    homework_ru TEXT,
    homework_en TEXT,
    estimated_study_time INTEGER,
    is_required BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (product_id) REFERENCES products(id) ON DELETE CASCADE
);

-- Индексы для уроков
CREATE INDEX IF NOT EXISTS idx_lessons_product_id ON lessons(product_id);
CREATE INDEX IF NOT EXISTS idx_lessons_sort_order ON lessons(product_id, sort_order);
CREATE INDEX IF NOT EXISTS idx_lessons_is_preview ON lessons(is_preview);

-- 5. Создание таблицы покупок (доступ пользователей)
CREATE TABLE IF NOT EXISTS purchases (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id INTEGER NOT NULL,
    product_id INTEGER NOT NULL,
    purchase_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    access_start TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    access_end TIMESTAMP,
    status TEXT NOT NULL DEFAULT 'active' CHECK (status IN ('active', 'expired', 'cancelled')),
    payment_id INTEGER,
    promo_code_id INTEGER,
    price_paid INTEGER NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (product_id) REFERENCES products(id) ON DELETE CASCADE
);

-- Индексы для покупок
CREATE INDEX IF NOT EXISTS idx_purchases_user_id ON purchases(user_id);
CREATE INDEX IF NOT EXISTS idx_purchases_product_id ON purchases(product_id);
CREATE INDEX IF NOT EXISTS idx_purchases_access_end ON purchases(access_end);
CREATE INDEX IF NOT EXISTS idx_purchases_status ON purchases(status);

-- 6. Создание таблицы прогресса обучения
CREATE TABLE IF NOT EXISTS learning_progress (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id INTEGER NOT NULL,
    purchase_id INTEGER NOT NULL,
    lesson_id INTEGER NOT NULL,
    completed BOOLEAN DEFAULT FALSE,
    completed_at TIMESTAMP,
    watch_duration_seconds INTEGER DEFAULT 0,
    last_position_seconds INTEGER DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(user_id, lesson_id),
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (purchase_id) REFERENCES purchases(id) ON DELETE CASCADE,
    FOREIGN KEY (lesson_id) REFERENCES lessons(id) ON DELETE CASCADE
);

-- Индексы для прогресса
CREATE INDEX IF NOT EXISTS idx_learning_progress_user_purchase ON learning_progress(user_id, purchase_id);
CREATE INDEX IF NOT EXISTS idx_learning_progress_completed ON learning_progress(completed);

-- Триггер для автоматического обновления статуса покупки при истечении срока
CREATE TRIGGER IF NOT EXISTS update_purchase_status_after_expiry
AFTER UPDATE OF access_end ON purchases
BEGIN
    UPDATE purchases
    SET status = 'expired'
    WHERE id = NEW.id
    AND NEW.access_end IS NOT NULL
    AND datetime(NEW.access_end) < datetime('now');
END;

-- Заполнение тестовыми данными (опционально)
INSERT OR IGNORE INTO product_categories (name_ru, name_en, slug, description_ru, description_en, sort_order) VALUES
('Акварель', 'Watercolor', 'watercolor', 'Курсы по акварельной живописи', 'Watercolor painting courses', 1),
('Масло', 'Oil Painting', 'oil-painting', 'Курсы по масляной живописи', 'Oil painting courses', 2),
('Графика', 'Graphics', 'graphics', 'Курсы по графике и рисунку', 'Graphics and drawing courses', 3),
('Цифровое искусство', 'Digital Art', 'digital-art', 'Курсы по цифровому искусству', 'Digital art courses', 4),
('Мастер-классы', 'Master Classes', 'master-classes', 'Разовые мастер-классы', 'One-time master classes', 5);

-- Добавление администратора (если нужно)
-- UPDATE users SET role = 'admin' WHERE username = 'admin_username';
