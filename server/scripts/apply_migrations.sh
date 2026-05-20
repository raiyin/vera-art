#!/bin/bash

# Скрипт для применения миграций базы данных
# Использование: ./apply_migrations.sh [путь_к_бд]

set -e

DB_PATH="${1:-./db/db.sqlite}"
MIGRATIONS_DIR="./migrations"

echo "Применение миграций к базе данных: $DB_PATH"

# Проверяем существование базы данных
if [ ! -f "$DB_PATH" ]; then
    echo "База данных не найдена: $DB_PATH"
    echo "Создаем новую базу данных..."
    touch "$DB_PATH"
fi

# Применяем миграции по порядку
for migration in $(ls $MIGRATIONS_DIR/*.sql | sort); do
    echo "Применяем миграцию: $(basename $migration)"
    sqlite3 "$DB_PATH" < "$migration"
done

echo "Все миграции успешно применены!"

# Создаем тестовые данные (опционально)
echo "Создаем тестовые данные..."
sqlite3 "$DB_PATH" <<EOF
-- Добавляем тестового администратора (пароль: Admin123!)
INSERT OR IGNORE INTO users (username, pass_hash, role, email, full_name)
VALUES ('admin', '\$2a\$10\$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi', 'admin', 'admin@example.com', 'Администратор');

-- Добавляем тестового пользователя (пароль: User123!)
INSERT OR IGNORE INTO users (username, pass_hash, role, email, full_name)
VALUES ('user', '\$2a\$10\$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi', 'user', 'user@example.com', 'Тестовый Пользователь');

-- Добавляем тестовые продукты
INSERT OR IGNORE INTO products (type, title_ru, title_en, description_ru, description_en, price, duration_days, status, difficulty, category_id, certificate_available, language) VALUES
('course', 'Основы акварели', 'Watercolor Basics', 'Научитесь основам акварельной живописи', 'Learn the basics of watercolor painting', 500000, 90, 'published', 'beginner', 1, TRUE, 'ru'),
('masterclass', 'Масляная живопись за один день', 'Oil Painting in One Day', 'Мастер-класс по масляной живописи', 'Oil painting master class', 300000, 30, 'published', 'intermediate', 2, FALSE, 'ru'),
('course', 'Цифровая иллюстрация', 'Digital Illustration', 'Курс по цифровой иллюстрации в Procreate', 'Digital illustration course in Procreate', 700000, 120, 'published', 'intermediate', 4, TRUE, 'both');

-- Добавляем тестовые уроки
INSERT OR IGNORE INTO lessons (product_id, title_ru, title_en, content_type, sort_order, duration_minutes) VALUES
(1, 'Введение в акварель', 'Introduction to Watercolor', 'video', 1, 45),
(1, 'Основные техники', 'Basic Techniques', 'video', 2, 60),
(2, 'Подготовка материалов', 'Materials Preparation', 'video', 1, 30),
(3, 'Настройка Procreate', 'Procreate Setup', 'video', 1, 40);

EOF

echo "Тестовые данные созданы!"
echo ""
echo "Доступные учетные записи:"
echo "  Администратор: username=admin, password=Admin123!"
echo "  Пользователь:  username=user, password=User123!"
