# Образовательная платформа - Backend

## Начало работы

### 1. Применение миграций базы данных

```bash
# Сделайте скрипт исполняемым
chmod +x scripts/apply_migrations.sh

# Примените миграции
./scripts/apply_migrations.sh
```

### 2. Запуск сервера

```bash
# Установите зависимости (если еще не установлены)
go mod tidy

# Запустите сервер
go run main.go auth.go categories.go
```

Сервер будет доступен по адресу: `http://localhost:8000`

## API Endpoints

### Аутентификация
- `POST /register` - Регистрация нового пользователя
- `POST /login` - Вход, получение JWT токенов
- `POST /refresh` - Обновление access токена

### Категории продуктов (публичные)
- `GET /categories` - Список всех активных категорий
- `GET /categories/:slug` - Категория по slug

### Категории продуктов (админка)
- `GET /admin/categories` - Все категории (требуется роль admin)
- `POST /admin/categories` - Создание категории (требуется роль admin)
- `PUT /admin/categories/:id` - Обновление категории (требуется роль admin)
- `DELETE /admin/categories/:id` - Удаление категории (требуется роль admin)

### Продукты (курсы/мастер-классы) - скоро
- `GET /products` - Список продуктов с фильтрами
- `GET /products/:id` - Детали продукта
- `POST /products` - Создание продукта (admin)
- `PUT /products/:id` - Обновление продукта (admin)
- `DELETE /products/:id` - Удаление продукта (admin)

## Тестовые данные

После применения миграций создаются:

### Пользователи:
1. **Администратор**
   - Username: `admin`
   - Password: `Admin123!`
   - Role: `admin`

2. **Обычный пользователь**
   - Username: `user`
   - Password: `User123!`
   - Role: `user`

### Категории:
1. Акварель (watercolor)
2. Масло (oil-painting)
3. Графика (graphics)
4. Цифровое искусство (digital-art)
5. Мастер-классы (master-classes)

### Продукты:
1. Курс "Основы акварели" (90 дней, с сертификатом)
2. Мастер-класс "Масляная живопись за один день" (30 дней)
3. Курс "Цифровая иллюстрация" (120 дней, с сертификатом)

## Модели данных

### Основные таблицы:
1. `users` - Пользователи с ролями
2. `product_categories` - Категории продуктов
3. `products` - Курсы и мастер-классы
4. `lessons` - Уроки/модули
5. `purchases` - Покупки и доступ
6. `learning_progress` - Прогресс обучения
7. `reviews` - Отзывы и рейтинги
8. `chat_threads` и `chat_messages` - Чат с преподавателем
9. `payments` - Платежи (интеграция с ЮKassa)
10. `promo_codes` - Промокоды

## Middleware

### Доступные middleware:
1. `AuthMiddleware()` - Проверка JWT токена
2. `middleware.AdminMiddleware()` - Проверка роли admin
3. `middleware.UserMiddleware()` - Проверка роли user или admin
4. `middleware.RoleMiddleware("role1", "role2")` - Проверка конкретных ролей

### Пример использования:
```go
// Только для администраторов
r.GET("/admin/endpoint", AuthMiddleware(), middleware.AdminMiddleware(), handler)

// Для пользователей и администраторов
r.GET("/user/endpoint", AuthMiddleware(), middleware.UserMiddleware(), handler)
```

## Следующие шаги разработки

### Фаза 1 (текущая):
- [x] Расширение модели пользователей с ролями
- [x] Создание таблиц для категорий, продуктов, уроков
- [x] API для управления категориями
- [ ] API для управления продуктами
- [ ] API для управления уроками
- [ ] Базовые тесты API

### Фаза 2:
- [ ] Интеграция с ЮKassa
- [ ] Система покупок и доступа
- [ ] Личный кабинет пользователя
- [ ] Middleware проверки доступа к контенту

### Фаза 3:
- [ ] Система отзывов и рейтингов
- [ ] Чат с преподавателем (long-polling)
- [ ] Прогресс обучения
- [ ] Админ-панель

### Фаза 4:
- [ ] Оптимизация производительности
- [ ] Кэширование
- [ ] Мониторинг и аналитика
- [ ] Документация API

## Тестирование API

### Пример запроса на получение категорий:
```bash
curl -X GET "http://localhost:8000/categories"
```

### Пример запроса на создание категории (админ):
```bash
curl -X POST "http://localhost:8000/admin/categories" \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "name_ru": "Новая категория",
    "name_en": "New Category",
    "slug": "new-category",
    "description_ru": "Описание категории",
    "description_en": "Category description",
    "sort_order": 10,
    "is_active": true
  }'
```

## Устранение неполадок

### Ошибка "database is locked":
SQLite не поддерживает конкурентные записи. Убедитесь, что:
1. Не запущено несколько экземпляров сервера
2. Нет других процессов, работающих с базой данных

### Ошибка "no such table":
Миграции не применены. Запустите скрипт миграций:
```bash
./scripts/apply_migrations.sh
```

### Ошибка "invalid token":
Токен истек или невалиден. Получите новый токен через `/login`.
