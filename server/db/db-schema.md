# Database Schema Documentation — Educational Platform

## Overview

SQLite database (`server/db/db.sqlite`) for the art educational platform.
Supports: master-classes (single video lessons), courses (multi-lesson programs), user management, payments, reviews, chat, and e-commerce.

---

## Table: `products`

Центральная таблица для всех образовательных продуктов (мастер-классы и курсы).

| Column | Type | Description |
|--------|------|-------------|
| `id` | INTEGER PK | Auto-increment ID |
| `type` | TEXT | `'masterclass'` — одно видео, `'course'` — курс из нескольких уроков |
| `title_ru` | TEXT | Название на русском |
| `title_en` | TEXT | Название на английском |
| `description_ru` | TEXT | Полное описание на русском |
| `description_en` | TEXT | Полное описание на английском |
| `short_description_ru` | TEXT | Краткое описание на русском |
| `short_description_en` | TEXT | Краткое описание на английском |
| `price` | INTEGER | Цена в **копейках** (500000 = 5000₽). 0 = бесплатно |
| `duration_days` | INTEGER? | Срок доступа в днях. NULL = бессрочно |
| `thumbnail_url` | TEXT | URL обложки |
| `video_url` | TEXT | URL видео (для мастер-класса — основное видео) |
| `status` | TEXT | `'draft'`, `'published'`, `'archived'` |
| `difficulty` | TEXT? | `'beginner'`, `'intermediate'`, `'advanced'` |
| `total_lessons` | INTEGER | Кол-во уроков (для курсов) |
| `total_duration_minutes` | INTEGER | Общая длительность в минутах |
| `category_id` | INTEGER FK → `product_categories.id` | Категория |
| `instructor_id` | INTEGER FK → `users.id` | Преподаватель |
| `tags` | JSON | **Устаревшее поле.** Хранит теги как JSON-массив строк. Мигрируем на `product_tags` |
| `prerequisites_ru` | TEXT | Требования на русском |
| `prerequisites_en` | TEXT | Требования на английском |
| `learning_outcomes_ru` | TEXT | Результаты обучения на русском |
| `learning_outcomes_en` | TEXT | Результаты обучения на английском |
| `certificate_available` | BOOLEAN | Доступен ли сертификат |
| `max_students` | INTEGER? | Макс. кол-во студентов (для курсов с ограничением) |
| `start_date` | TIMESTAMP? | Дата старта (для курсов с расписанием) |
| `language` | TEXT | `'ru'`, `'en'`, `'both'` |
| `is_featured` | BOOLEAN | Рекомендованный продукт |
| `view_count` | INTEGER | Счётчик просмотров |
| `created_at` | TIMESTAMP | Дата создания |
| `updated_at` | TIMESTAMP | Дата обновления |

---

## Table: `product_categories`

Категории продуктов.

| Column | Type | Description |
|--------|------|-------------|
| `id` | INTEGER PK | |
| `name_ru` | TEXT | Название на русском |
| `name_en` | TEXT | Название на английском |
| `slug` | TEXT UNIQUE | URL-идентификатор |
| `description_ru` | TEXT? | Описание на русском |
| `description_en` | TEXT? | Описание на английском |
| `sort_order` | INTEGER | Порядок сортировки |
| `is_active` | BOOLEAN | Активна ли категория |
| `created_at` | TIMESTAMP | |
| `updated_at` | TIMESTAMP | |

**Seed data:**

| id | name_ru | slug |
|----|---------|------|
| 1 | Акварель | watercolor |
| 2 | Масло | oil-painting |
| 3 | Графика | graphics |
| 4 | Цифровое искусство | digital-art |
| 5 | Мастер-классы | master-classes |

---

## Table: `tags`

Теги для продуктов (связь многие-ко-многим через `product_tags`).

| Column | Type | Description |
|--------|------|-------------|
| `id` | INTEGER PK | |
| `name_ru` | TEXT NOT NULL UNIQUE | Название тега на русском |
| `name_en` | TEXT NOT NULL UNIQUE | Название тега на английском |
| `slug` | TEXT NOT NULL UNIQUE | URL-идентификатор |
| `created_at` | TIMESTAMP | |

---

## Table: `product_tags`

Связь многие-ко-многим между продуктами и тегами.

| Column | Type | Description |
|--------|------|-------------|
| `id` | INTEGER PK | |
| `product_id` | INTEGER FK → `products.id` ON DELETE CASCADE | |
| `tag_id` | INTEGER FK → `tags.id` ON DELETE CASCADE | |
| `created_at` | TIMESTAMP | |

**Unique constraint:** `(product_id, tag_id)`

---

## Table: `lessons`

Уроки внутри курса (`type = 'course'`). Для мастер-классов (`type = 'masterclass'`) уроки не используются — всё видео хранится в `products.video_url`.

| Column | Type | Description |
|--------|------|-------------|
| `id` | INTEGER PK | |
| `product_id` | INTEGER FK → `products.id` | |
| `title_ru` | TEXT | Название на русском |
| `title_en` | TEXT | Название на английском |
| `description_ru` | TEXT? | Описание на русском |
| `description_en` | TEXT? | Описание на английском |
| `content_type` | TEXT | `'video'`, `'text'`, `'pdf'`, `'quiz'`, `'assignment'` |
| `content_url` | TEXT? | URL контента |
| `duration_minutes` | INTEGER | Длительность |
| `sort_order` | INTEGER | Порядок сортировки |
| `is_preview` | BOOLEAN | Доступен ли без покупки |
| `resources` | JSON | Доп. материалы |
| `homework_ru` | TEXT? | ДЗ на русском |
| `homework_en` | TEXT? | ДЗ на английском |
| `estimated_study_time` | INTEGER? | Примерное время изучения |
| `is_required` | BOOLEAN | Обязательный урок |
| `created_at` | TIMESTAMP | |
| `updated_at` | TIMESTAMP | |

---

## Table: `purchases`

Покупки продуктов пользователями.

| Column | Type | Description |
|--------|------|-------------|
| `id` | INTEGER PK | |
| `user_id` | INTEGER FK → `users.id` | |
| `product_id` | INTEGER FK → `products.id` | |
| `purchase_date` | TIMESTAMP | Дата покупки |
| `access_start` | TIMESTAMP | Начало доступа |
| `access_end` | TIMESTAMP? | Окончание доступа (NULL = бессрочно) |
| `status` | TEXT | `'active'`, `'expired'`, `'cancelled'` |
| `payment_id` | INTEGER? FK → `payments.id` | |
| `promo_code_id` | INTEGER? FK → `promo_codes.id` | |
| `price_paid` | INTEGER | Фактически уплаченная сумма в копейках |
| `created_at` | TIMESTAMP | |

---

## Table: `learning_progress`

Прогресс пользователя по урокам.

| Column | Type | Description |
|--------|------|-------------|
| `id` | INTEGER PK | |
| `user_id` | INTEGER FK → `users.id` | |
| `purchase_id` | INTEGER FK → `purchases.id` | |
| `lesson_id` | INTEGER FK → `lessons.id` | |
| `completed` | BOOLEAN | Завершён ли урок |
| `completed_at` | TIMESTAMP? | |
| `watch_duration_seconds` | INTEGER? | |
| `last_position_seconds` | INTEGER? | |
| `created_at` | TIMESTAMP | |
| `updated_at` | TIMESTAMP | |

---

## Table: `users`

| Column | Type | Description |
|--------|------|-------------|
| `id` | INTEGER PK | Auto-increment ID |
| `username` | TEXT NOT NULL UNIQUE | Login username |
| `password_hash` | TEXT NOT NULL | bcrypt hash of password |
| `role` | TEXT NOT NULL | `'user'`, `'admin'`. Default `'user'` |
| `email` | TEXT? | Email address |
| `name` | TEXT? | Display name (defaults to username on registration) |
| `avatar_path` | TEXT | Avatar file path. Default `''` |
| `email_verified` | INTEGER NOT NULL | `0` = unverified, `1` = verified. Default `0` |
| `verification_token` | TEXT? | Cryptographically secure random token for email verification |
| `verification_sent_at` | TIMESTAMP? | When the verification email was last sent |
| `created_at` | TIMESTAMP | Record creation timestamp |
| `updated_at` | TIMESTAMP | Record last update timestamp |

---

## Table: `payments`

| Column | Type | Description |
|--------|------|-------------|
| `id` | INTEGER PK | |
| `user_id` | INTEGER FK → `users.id` | |
| `external_id` | TEXT? | ID платежа в ЮKassa |
| `status` | TEXT | `'pending'`, `'waiting_for_capture'`, `'succeeded'`, `'canceled'`, `'refunded'` |
| `amount` | INTEGER | Сумма в копейках |
| `currency` | TEXT | `'RUB'` |
| `description` | TEXT? | |
| `payment_method` | TEXT? | |
| `metadata` | JSON? | |
| `created_at` | TIMESTAMP | |
| `updated_at` | TIMESTAMP | |

---

## Table: `promo_codes`

| Column | Type | Description |
|--------|------|-------------|
| `id` | INTEGER PK | |
| `code` | TEXT UNIQUE | |
| `discount_type` | TEXT | `'percentage'`, `'fixed'` |
| `discount_value` | INTEGER | |
| `max_uses` | INTEGER? | |
| `used_count` | INTEGER | |
| `valid_from` | TIMESTAMP? | |
| `valid_until` | TIMESTAMP? | |
| `is_active` | BOOLEAN | |
| `created_at` | TIMESTAMP | |

---

## Table: `reviews`

| Column | Type | Description |
|--------|------|-------------|
| `id` | INTEGER PK | |
| `user_id` | INTEGER FK → `users.id` | |
| `product_id` | INTEGER FK → `products.id` | |
| `purchase_id` | INTEGER FK → `purchases.id` | |
| `rating` | INTEGER | 1–5 |
| `title_ru` | TEXT? | |
| `title_en` | TEXT? | |
| `comment_ru` | TEXT? | |
| `comment_en` | TEXT? | |
| `is_approved` | BOOLEAN | |
| `is_visible` | BOOLEAN | |
| `created_at` | TIMESTAMP | |
| `updated_at` | TIMESTAMP | |

---

## Table: `chat_threads`

| Column | Type | Description |
|--------|------|-------------|
| `id` | INTEGER PK | |
| `purchase_id` | INTEGER FK → `purchases.id` | |
| `user_id` | INTEGER FK → `users.id` | |
| `admin_id` | INTEGER? FK → `users.id` | |
| `last_message_at` | TIMESTAMP | |
| `is_resolved` | BOOLEAN | |
| `created_at` | TIMESTAMP | |

---

## Table: `chat_messages`

| Column | Type | Description |
|--------|------|-------------|
| `id` | INTEGER PK | |
| `thread_id` | INTEGER FK → `chat_threads.id` | |
| `sender_id` | INTEGER FK → `users.id` | |
| `message_type` | TEXT | `'text'`, `'image'`, `'file'` |
| `content` | TEXT | |
| `attachment_url` | TEXT? | |
| `attachment_size` | INTEGER? | |
| `is_read` | BOOLEAN | |
| `read_at` | TIMESTAMP? | |
| `created_at` | TIMESTAMP | |

---

## Legacy Tables (not used for educational platform)

- `materials` — материалы для картин (холст, бумага и т.д.)
- `bases` — основы (масло, акварель)
- `works` — картины (галерея)
- `sales` — картины в продаже
- `illustrations` — иллюстрации
- `threeds` — 3D-работы
- `work_types` — типы работ
- `news` — новости
- `consent_log` — логи согласий на cookie
- `sales_materials`, `works_materials` — связи материалов с работами

---

## Entity Relationship Summary

```
products ──1:N──> lessons
products ──M:N──> tags (via product_tags)
products ──N:1──> product_categories
products ──1:N──> purchases ──N:1──> users
purchases ──1:N──> learning_progress ──N:1──> lessons
purchases ──1:N──> chat_threads ──1:N──> chat_messages
products ──1:N──> reviews ──N:1──> users
purchases ──N:1──> payments
purchases ──N:1──> promo_codes
