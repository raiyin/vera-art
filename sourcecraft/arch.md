# Архитектура проекта (Nuxt.js + Go)

## Общая архитектура

```
┌─────────────────────────────────────────────────────────────┐
│                    Клиент (Nuxt.js 4)                        │
│  ┌──────────┐  ┌──────────┐  ┌──────────┐  ┌────────────┐  │
│  │  Pages   │  │  Stores  │  │  API     │  │ Components │  │
│  │ (file-   │  │ (Pinia)  │  │ (Axios/  │  │ (Nuxt UI)  │  │
│  │  based)  │  │          │  │  fetch)  │  │            │  │
│  └────┬─────┘  └────┬─────┘  └────┬─────┘  └────┬───────┘  │
│       └──────────────┴─────────────┴──────────────┘          │
│                          │ HTTP/JSON                         │
└──────────────────────────┼──────────────────────────────────┘
                           │
                    ┌──────┴──────┐
                    │  Nginx/CDN  │ (статический прокси)
                    └──────┬──────┘
                           │
┌──────────────────────────┼──────────────────────────────────┐
│                    Сервер (Go + Gin)                         │
│  ┌──────────────────────────────────────────────────────┐   │
│  │                    Router                             │   │
│  │  /api/v1/auth/*  /api/v1/shop/*  /api/v1/admin/* ... │   │
│  └──────────┬───────────────────────────────────────────┘   │
│             │                                                │
│  ┌──────────┴──────────┐                                     │
│  │     Middleware       │                                     │
│  │  Auth | Admin | CORS│                                     │
│  │  RateLimit | Logger  │                                     │
│  └──────────┬──────────┘                                     │
│             │                                                │
│  ┌──────────┴──────────┐                                     │
│  │      Handlers        │  (HTTP-обработчики)                │
│  │  auth_handler.go     │                                     │
│  │  gallery_handler.go  │                                     │
│  │  shop_handler.go     │                                     │
│  │  learning_handler.go │                                     │
│  │  payment_handler.go  │                                     │
│  │  chat_handler.go     │                                     │
│  │  news_handler.go     │                                     │
│  │  misc_handler.go     │                                     │
│  │  profile_handler.go  │                                     │
│  └──────────┬──────────┘                                     │
│             │                                                │
│  ┌──────────┴──────────┐                                     │
│  │  Service Interfaces  │  (port/service.go)                 │
│  │  UserService         │                                     │
│  │  GalleryService      │                                     │
│  │  ShopService         │                                     │
│  │  LearningService     │                                     │
│  │  PaymentService      │                                     │
│  │  ChatService         │                                     │
│  │  NewsService         │                                     │
│  │  MiscService         │                                     │
│  └──────────┬──────────┘                                     │
│             │                                                │
│  ┌──────────┴──────────┐                                     │
│  │  Service Implement.  │  (service/*.go)                    │
│  │  user_service.go     │                                     │
│  │  shop_service.go     │                                     │
│  │  payment_service.go  │                                     │
│  │  news_service.go     │                                     │
│  └──────────┬──────────┘                                     │
│             │                                                │
│  ┌──────────┴──────────┐                                     │
│  │  Repository Interfaces│ (port/repository.go)              │
│  │  UserRepository      │                                     │
│  │  WorkRepository      │                                     │
│  │  SaleRepository      │                                     │
│  │  ProductRepository   │                                     │
│  │  ... (18 интерфейсов)│                                     │
│  └──────────┬──────────┘                                     │
│             │                                                │
│  ┌──────────┴──────────┐  ┌──────────────────────┐          │
│  │  SQLite Repository  │  │  File Repository     │          │
│  │  (database/sql)     │  │  (images, avatars,   │          │
│  │  WAL + foreign_keys │  │   news media)        │          │
│  └──────────┬──────────┘  └──────────────────────┘          │
│             │                                                │
│  ┌──────────┴──────────┐                                     │
│  │    SQLite3 DB        │                                     │
│  │  (1 connection max)  │                                     │
│  └─────────────────────┘                                     │
└──────────────────────────────────────────────────────────────┘
```

## Слои архитектуры (Go)

### 1. Domain Layer (`server/internal/domain/`)

Чистые доменные модели без зависимостей от фреймворков и БД.

| Файл | Сущности |
|------|----------|
| [`user.go`](server/internal/domain/user.go) | `User` (ID, Username, Password, Email, FullName, Role, EmailVerified, AvatarPath, CreatedAt) |
| [`work.go`](server/internal/domain/work.go) | `Work` (ID, Title, Description, ImagePath, Materials, Bases, CreatedAt) |
| [`sale.go`](server/internal/domain/sale.go) | `Sale` (ID, WorkID, Price, IsSold, SoldAt) |
| [`product.go`](server/internal/domain/product.go) | `Product` (ID, Title, Slug, Description, FullDescription, Price, OldPrice, ImagePath, CategoryID, Status, IsDigital, IsMasterClass, SortOrder, Tags) |
| [`lesson.go`](server/internal/domain/lesson.go) | `Lesson` (ID, ProductID, Title, Description, VideoURL, Resources, SortOrder, Duration) |
| [`payment.go`](server/internal/domain/payment.go) | `Payment` (ID, UserID, Amount, Currency, Status, YooKassaID, YooKassaURL), `Purchase` (ID, UserID, ProductID, PaymentID, Status) |
| [`news.go`](server/internal/domain/news.go) | `News` (ID, Title, Description, ImagePath, VideoURL, CreatedAt) |
| [`chat.go`](server/internal/domain/chat.go) | `ChatThread` (ID, UserID, Status, CreatedAt), `ChatMessage` (ID, ThreadID, UserID, Content, CreatedAt) |
| [`misc.go`](server/internal/domain/misc.go) | `Tag`, `Material`, `Base`, `Category`, `PromoCode`, `Review`, `MasterClass`, `UserConsent` |
| [`errors.go`](server/internal/domain/errors.go) | `ErrNotFound`, `ErrUnauthorized`, `ErrForbidden`, `ErrConflict`, `ErrValidation`, `ErrInternal` |
| [`upload.go`](server/internal/domain/upload.go) | `UploadedFile` (Name, Path, Size, MimeType) |

### 2. Port Layer (`server/internal/port/`)

Интерфейсы для сервисов и репозиториев (принцип инверсии зависимостей).

* [`service.go`](server/internal/port/service.go) — 14 интерфейсов сервисов (UserService, GalleryService, ShopService, LearningService, PaymentService, ChatService, NewsService, MiscService, ProfileService, AuthService, ReviewService, PromoCodeService, MasterClassService, ConsentService)
* [`repository.go`](server/internal/port/repository.go) — 18 интерфейсов репозиториев (UserRepository, WorkRepository, SaleRepository, ProductRepository, LessonRepository, PaymentRepository, PurchaseRepository, NewsRepository, ChatThreadRepository, ChatMessageRepository, TagRepository, MaterialRepository, BaseRepository, CategoryRepository, PromoCodeRepository, ReviewRepository, MasterClassRepository, ConsentRepository)

### 3. Handler Layer (`server/internal/handler/`)

HTTP-обработчики на Gin. Принимают HTTP-запросы, валидируют DTO, вызывают сервисы, возвращают JSON-ответы.

| Файл | Обработчики |
|------|-------------|
| [`auth_handler.go`](server/internal/handler/auth_handler.go) | Register, Login, RefreshToken, VerifyEmail, ResendVerification |
| [`profile_handler.go`](server/internal/handler/profile_handler.go) | GetProfile, UpdateProfile, UploadAvatar, DeleteAvatar |
| [`gallery_handler.go`](server/internal/handler/gallery_handler.go) | ListWorks, GetWork, CreateWork, UpdateWork, DeleteWork, ListSales, GetSale, CreateSale, UpdateSale, DeleteSale |
| [`shop_handler.go`](server/internal/handler/shop_handler.go) | ListProducts, GetProduct, CreateProduct, UpdateProduct, DeleteProduct, ListCategories, CreateCategory, UpdateCategory, DeleteCategory, ListPromoCodes, CreatePromoCode, UpdatePromoCode, DeletePromoCode, ListReviews, CreateReview, UpdateReview, DeleteReview |
| [`learning_handler.go`](server/internal/handler/learning_handler.go) | ListLessons, GetLesson, CreateLesson, UpdateLesson, DeleteLesson, GetMyCourses, GetLessonProgress, UpdateLessonProgress |
| [`payment_handler.go`](server/internal/handler/payment_handler.go) | CreatePayment, YooKassaWebhook, GetPaymentStatus, ListMyPurchases, CreatePurchase |
| [`chat_handler.go`](server/internal/handler/chat_handler.go) | GetThreads, CreateThread, GetMessages, SendMessage, AdminListThreads, AdminGetMessages, AdminSendMessage, ResolveThread, ReopenThread |
| [`news_handler.go`](server/internal/handler/news_handler.go) | ListNews, GetNews, CreateNews, UpdateNews, DeleteNews |
| [`misc_handler.go`](server/internal/handler/misc_handler.go) | ListTags, CreateTag, UpdateTag, DeleteTag, ListMaterials, CreateMaterial, ListBases, CreateBase, ListMasterClasses, CreateMasterClass, UpdateMasterClass, DeleteMasterClass, SaveConsent, GetConsent, GetDashboardStats, GetRecentActivity |
| [`middleware.go`](server/internal/handler/middleware.go) | AuthMiddleware, AdminMiddleware, CORSMiddleware, RateLimitMiddleware, Error handling |

### 4. DTO Layer (`server/internal/dto/`)

Структуры запросов и ответов для каждого эндпоинта.

| Файл | Назначение |
|------|------------|
| [`auth.go`](server/internal/dto/auth.go) | RegisterRequest, LoginRequest, RefreshRequest, AuthResponse |
| [`profile.go`](server/internal/dto/profile.go) | ProfileResponse, UpdateProfileRequest |
| [`gallery.go`](server/internal/dto/gallery.go) | WorkRequest, WorkResponse, SaleRequest, SaleResponse, ListWorksResponse, ListSalesResponse |
| [`shop.go`](server/internal/dto/shop.go) | ProductRequest, ProductResponse, CategoryRequest, PromoCodeRequest, ReviewRequest |
| [`learning.go`](server/internal/dto/learning.go) | LessonRequest, LessonResponse, LessonProgressResponse |
| [`payment.go`](server/internal/dto/payment.go) | CreatePaymentRequest, PaymentResponse, PurchaseResponse |
| [`chat.go`](server/internal/dto/chat.go) | CreateThreadRequest, SendMessageRequest, ThreadResponse, MessageResponse |
| [`news.go`](server/internal/dto/news.go) | NewsRequest, NewsResponse |
| [`misc.go`](server/internal/dto/misc.go) | TagRequest, MaterialRequest, BaseRequest, MasterClassRequest, ConsentRequest, DashboardResponse |

### 5. Service Layer (`server/internal/service/`)

Бизнес-логика. Реализует интерфейсы из `port/service.go`.

| Файл | Логика |
|------|--------|
| [`user_service.go`](server/internal/service/user_service.go) | Регистрация, логин, верификация email, JWT-токены, refresh-токены |
| [`shop_service.go`](server/internal/service/shop_service.go) | CRUD продуктов, категорий, промокодов, отзывов |
| [`payment_service.go`](server/internal/service/payment_service.go) | Интеграция с YooKassa, создание платежей, обработка webhook'ов, покупки |
| [`news_service.go`](server/internal/service/news_service.go) | CRUD новостей |

### 6. Repository Layer (SQLite)

Реализации интерфейсов из `port/repository.go`. Прямые SQL-запросы через `database/sql`.

### 7. Config Layer (`server/internal/config/`)

* [`config.go`](server/internal/config/config.go) — структура `Config` с подструктурами `App`, `CORS`, `Directories`, `SMTP`
* Загрузка из [`config.yaml`](server/config.yaml) + `.env` через Viper + godotenv

## Маршрутизация (API Routes)

Все маршруты определены в [`router.go`](server/internal/router/router.go).

### Публичные эндпоинты (без аутентификации)

```
POST   /api/v1/auth/register
POST   /api/v1/auth/login
POST   /api/v1/auth/refresh
POST   /api/v1/auth/verify-email
POST   /api/v1/auth/resend-verification
GET    /api/v1/sales
GET    /api/v1/works
GET    /api/v1/works/:id
GET    /api/v1/news
GET    /api/v1/news/:id
GET    /api/v1/materials
GET    /api/v1/bases
GET    /api/v1/categories
GET    /api/v1/products
GET    /api/v1/products/:slug
GET    /api/v1/master-classes
GET    /api/v1/tags
```

### Эндпоинты с аутентификацией (AuthMiddleware)

```
GET    /api/v1/profile
PUT    /api/v1/profile
POST   /api/v1/profile/avatar
DELETE /api/v1/profile/avatar
POST   /api/v1/payments
GET    /api/v1/payments/:id
POST   /api/v1/payments/webhook
GET    /api/v1/purchases
POST   /api/v1/purchases
GET    /api/v1/learning/courses
GET    /api/v1/learning/lessons/:id
GET    /api/v1/learning/progress
PUT    /api/v1/learning/progress/:id
GET    /api/v1/chat/threads
POST   /api/v1/chat/threads
GET    /api/v1/chat/threads/:id/messages
POST   /api/v1/chat/threads/:id/messages
POST   /api/v1/consent
GET    /api/v1/consent
```

### Админ-эндпоинты (AuthMiddleware + AdminMiddleware)

```
POST   /api/v1/admin/works
PUT    /api/v1/admin/works/:id
DELETE /api/v1/admin/works/:id
POST   /api/v1/admin/sales
PUT    /api/v1/admin/sales/:id
DELETE /api/v1/admin/sales/:id
POST   /api/v1/admin/news
PUT    /api/v1/admin/news/:id
DELETE /api/v1/admin/news/:id
POST   /api/v1/admin/categories
PUT    /api/v1/admin/categories/:id
DELETE /api/v1/admin/categories/:id
POST   /api/v1/admin/tags
PUT    /api/v1/admin/tags/:id
DELETE /api/v1/admin/tags/:id
POST   /api/v1/admin/products
PUT    /api/v1/admin/products/:id
DELETE /api/v1/admin/products/:id
POST   /api/v1/admin/lessons
PUT    /api/v1/admin/lessons/:id
DELETE /api/v1/admin/lessons/:id
POST   /api/v1/admin/promo-codes
PUT    /api/v1/admin/promo-codes/:id
DELETE /api/v1/admin/promo-codes/:id
PUT    /api/v1/admin/reviews/:id
DELETE /api/v1/admin/reviews/:id
POST   /api/v1/admin/materials
POST   /api/v1/admin/bases
POST   /api/v1/admin/master-classes
PUT    /api/v1/admin/master-classes/:id
DELETE /api/v1/admin/master-classes/:id
GET    /api/v1/admin/chat/threads
GET    /api/v1/admin/chat/threads/:id/messages
POST   /api/v1/admin/chat/threads/:id/messages
PUT    /api/v1/admin/chat/threads/:id/resolve
PUT    /api/v1/admin/chat/threads/:id/reopen
GET    /api/v1/admin/payments
GET    /api/v1/admin/purchases
PUT    /api/v1/admin/purchases/:id/cancel
GET    /api/v1/admin/stats
GET    /api/v1/admin/recent-activity
```

## Клиентская архитектура (Nuxt.js 4)

### Структура директорий

```
client/
├── app/
│   ├── api/              # API-клиенты (auth.ts, admin.ts, master-classes.ts, requests.ts)
│   ├── assets/           # CSS, изображения, иконки
│   ├── components/       # Переиспользуемые компоненты
│   │   └── admin/        # Админ-компоненты (таблицы, фильтры, пагинация)
│   ├── layouts/          # default.vue, admin.vue
│   ├── middleware/        # admin-auth.ts
│   ├── pages/            # Файловая маршрутизация
│   │   ├── admin/        # Админ-панель (dashboard, reviews, categories, chats, courses, etc.)
│   │   ├── art-store/    # Магазин (index, add, edit/[id])
│   │   ├── auth/         # Аутентификация (login, register, verify-email)
│   │   ├── chat/         # Чат (index, [threadId])
│   │   ├── courses/      # Курсы (index, individual, online)
│   │   ├── gallery/      # Галерея (index, add, edit/[id])
│   │   ├── learning/     # Обучение (index, cart)
│   │   ├── master-classes/ # Мастер-классы
│   │   ├── news/         # Новости (index, [id], add, edit/[id])
│   │   ├── profile/      # Профиль
│   │   └── reviews/      # Отзывы (write)
│   ├── stores/           # Pinia stores (AuthStore, ProductStore, MaterialStore, etc.)
│   ├── types/            # TypeScript интерфейсы
│   └── utils/            # Утилиты (jwt.ts)
├── i18n/                 # Локализация (locales/ru.json, locales/en.json)
├── public/               # Статические файлы
├── nuxt.config.ts        # Конфигурация Nuxt
└── package.json          # Зависимости
```

### Pinia Stores

| Store | Файл | Назначение |
|-------|------|------------|
| [`AuthStore`](client/app/stores/AuthStore.ts) | Управление JWT-токенами, localStorage, роль/userId |
| [`ProductStore`](client/app/stores/ProductStore.ts) | CRUD продуктов, фильтры, категории, избранное |
| [`MaterialStore`](client/app/stores/MaterialStore.ts) | Справочники (materials, bases) с двуязычными названиями |
| [`NotificationStore`](client/app/stores/NotificationStore.ts) | Polling новых сообщений чата (30 сек) |
| [`ThemeStore`](client/app/stores/ThemeStore.ts) | Светлая/тёмная тема, localStorage, cross-tab sync |
| [`CookieConsentStore`](client/app/stores/CookieConsentStore.ts) | GDPR-согласие (analytics, marketing) |

### Компонентная архитектура

```
App.vue
├── default.vue (layout)
│   ├── Header.vue
│   ├── <NuxtPage />
│   ├── ChatWidget.vue
│   ├── CookieConsent.vue
│   └── AppFooter.vue
│
├── admin.vue (layout)
│   ├── AdminSidebar.vue
│   ├── AdminHeader.vue
│   └── <NuxtPage />
│
├── Компоненты галереи
│   ├── PicGallery.vue
│   ├── PicCarousel.vue
│   ├── PictureCardSkeleton.vue
│   └── PhotoSection.vue
│
├── Компоненты новостей
│   ├── NewsCarousel.vue
│   ├── NewsItemDescription.vue
│   ├── NewsPhotoItem.vue
│   ├── NewsTrailer.vue
│   ├── SideNewsTrailer.vue
│   └── NewsDescriptionSkeleton.vue / SideNewsTrailerSkeleton.vue
│
├── Компоненты магазина
│   ├── ProductCard.vue
│   ├── ProductList.vue
│   ├── PriceBadge.vue
│   └── PaymentForm.vue
│
├── Компоненты отзывов
│   ├── ReviewCard.vue
│   └── ReviewForm.vue
│
├── Компоненты обучения
│   └── ProgressTracker.vue
│
└── Admin-компоненты
    ├── AdminBreadcrumbs.vue
    ├── AdminConfirmDialog.vue
    ├── AdminDataTable.vue
    ├── AdminEmptyState.vue
    ├── AdminFilterBar.vue
    ├── AdminImageUploader.vue
    ├── AdminLoadingSkeleton.vue
    ├── AdminPagination.vue
    ├── AdminRichEditor.vue
    ├── AdminSearchInput.vue
    ├── AdminSidebar.vue
    ├── AdminStatusBadge.vue
    └── AdminHeader.vue
```

## Потоки данных (Data Flows)

### Аутентификация

```
Register:
  Client → POST /api/v1/auth/register {username, password, email}
  → AuthHandler.Register → UserService.Register
  → UserRepository.Create (bcrypt password hash)
  → JWT: generate access + refresh tokens
  → Email: send verification email (SMTP)
  ← Response: {access_token, refresh_token, user}

Login:
  Client → POST /api/v1/auth/login {username, password}
  → AuthHandler.Login → UserService.Login
  → UserRepository.FindByUsername → bcrypt.CompareHashAndPassword
  → JWT: generate access + refresh tokens
  ← Response: {access_token, refresh_token, user}

Token Refresh:
  Client → POST /api/v1/auth/refresh {refresh_token}
  → AuthHandler.RefreshToken → UserService.RefreshToken
  → Validate refresh token → Check in DB (refresh_tokens table)
  → Generate new access + refresh tokens
  ← Response: {access_token, refresh_token}
```

### Платежи (YooKassa)

```
Create Payment:
  Client → POST /api/v1/payments {product_id, promo_code?}
  → PaymentHandler.CreatePayment → PaymentService.CreatePayment
  → Validate product exists → Calculate price (with promo code)
  → YooKassa API: create payment → Get confirmation URL
  → Save Payment record in DB (status: pending)
  ← Response: {payment_id, confirmation_url, amount}

YooKassa Webhook:
  YooKassa → POST /api/v1/payments/webhook {event, object}
  → PaymentHandler.YooKassaWebhook → PaymentService.HandleWebhook
  → Verify signature → Update Payment status
  → If succeeded: Create Purchase record → Update product status
  ← Response: 200 OK

Purchase Flow:
  Client → POST /api/v1/purchases {payment_id}
  → PaymentHandler.CreatePurchase → PaymentService.CreatePurchase
  → Verify payment succeeded → Create Purchase record
  ← Response: {purchase_id, product_id, access_granted}
```

### Чат (поддержка)

```
Create Thread:
  Client → POST /api/v1/chat/threads {subject?}
  → ChatHandler.CreateThread → ChatService.CreateThread
  → ChatThreadRepository.Create
  ← Response: {thread_id, status: "open"}

Send Message:
  Client → POST /api/v1/chat/threads/:id/messages {content}
  → ChatHandler.SendMessage → ChatService.SendMessage
  → ChatMessageRepository.Create
  ← Response: {message_id, content, created_at}

Polling (Client):
  NotificationStore → GET /api/v1/chat/threads (every 30s)
  → Check for unread messages → Update badge count
```

## База данных (SQLite)

* **Движок**: SQLite3 через `mattn/go-sqlite3`
* **Режим**: WAL (Write-Ahead Logging)
* **Параметры**: `busy_timeout=5000ms`, `foreign_keys=ON`
* **Соединения**: максимум 1 (однопоточный доступ)
* **Миграции**: инлайн SQL в [`main.go`](server/cmd/server/main.go)

### Основные таблицы

| Таблица | Назначение |
|---------|------------|
| `users` | Пользователи (username, pass_hash, email, role, email_verified) |
| `refresh_tokens` | Refresh-токены (user_id, token_hash, expires_at) |
| `works` | Работы галереи (title, description, image_path) |
| `sales` | Продажи работ (work_id, price, is_sold) |
| `news` | Новости (title, description, image_path, video_url) |
| `products` | Товары/курсы (title, slug, price, category_id, status, is_digital) |
| `categories` | Категории товаров (name, slug) |
| `tags` | Теги (name_ru, name_en, slug) |
| `product_tags` | Связь товаров и тегов |
| `lessons` | Уроки курсов (product_id, title, video_url, sort_order) |
| `lesson_progress` | Прогресс уроков (user_id, lesson_id, completed, score) |
| `payments` | Платежи (user_id, amount, status, yookassa_id) |
| `purchases` | Покупки (user_id, product_id, payment_id, status) |
| `chat_threads` | Темы чата (user_id, status: open/resolved) |
| `chat_messages` | Сообщения чата (thread_id, user_id, content) |
| `promo_codes` | Промокоды (code, discount_type, discount_value, expires_at) |
| `reviews` | Отзывы (user_id, product_id, rating, text, moderated) |
| `master_classes` | Мастер-классы (title, description, price, date) |
| `materials` | Материалы (name_ru, name_en) |
| `bases` | Основы (name_ru, name_en) |
| `user_consent` | Согласия пользователей (analytics, marketing) |

## Ключевые архитектурные решения

### 1. SQLite с одним соединением
* **Решение**: одно соединение с WAL-режимом вместо PostgreSQL/MySQL.
* **Обоснование**: простота деплоя, отсутствие необходимости в отдельном сервере БД, достаточная производительность для однопользовательского/admin-нагруженного приложения.
* **Компромисс**: ограниченная масштабируемость при росте числа одновременных записей.

### 2. Два API-клиента на фронтенде
* **Решение**: Axios для публичных/auth-запросов, fetch для админ-запросов.
* **Обоснование**: Axios предоставляет удобные перехватчики для автоматического refresh'а токена. Админ-модуль использует fetch для изоляции и ручного контроля.
* **Компромисс**: дублирование логики refresh'а токена (автоматический в Axios, ручной в admin.ts).

### 3. In-memory rate limiter
* **Решение**: простая карта `map[string]*rateLimiter` в памяти.
* **Обоснование**: не требует внешних зависимостей (Redis), достаточно для одного инстанса.
* **Компромисс**: сброс лимитов при перезапуске сервера, не работает при горизонтальном масштабировании.

### 4. Локальное файловое хранилище
* **Решение**: загрузка файлов на локальную файловую систему.
* **Обоснование**: простота, отсутствие зависимости от S3/Cloud Storage.
* **Компромисс**: необходимость синхронизации файлов между инстансами при масштабировании, отсутствие CDN.

### 5. JWT без refresh-ротации
* **Решение**: access token (15 мин) + refresh token (7 дней), refresh token хранится в БД.
* **Обоснование**: refresh token можно отозвать (удалить из БД). Access token живёт 15 минут для минимизации окна уязвимости.
* **Компромисс**: отсутствие ротации refresh-токенов (уязвимость к краже refresh-токена).

### 6. Отсутствие OpenAPI/Swagger
* **Решение**: документация API отсутствует.
* **Рекомендация**: добавить генерацию OpenAPI-спецификации для автоматической генерации клиента и документации.

### 7. Два паттерна API на клиенте
* **Решение**: Axios (auth.ts) и fetch (admin.ts, MaterialStore, master-classes.ts).
* **Проблема**: разные подходы к обработке ошибок, auth-заголовкам и refresh'у токена.
* **Рекомендация**: унифицировать API-клиент.

## Зависимости

### Сервер (Go)

| Пакет | Назначение |
|-------|------------|
| `github.com/gin-gonic/gin` | HTTP-фреймворк |
| `github.com/mattn/go-sqlite3` | SQLite3 драйвер |
| `github.com/golang-jwt/jwt/v5` | JWT-токены |
| `github.com/spf13/viper` | Конфигурация |
| `github.com/joho/godotenv` | .env файлы |
| `github.com/gin-contrib/cors` | CORS- middleware |
| `golang.org/x/crypto` | bcrypt хеширование |

### Клиент (Nuxt.js)

| Пакет | Назначение |
|-------|------------|
| `nuxt` (v4) | Фреймворк |
| `@nuxt/ui` (v4) | UI-компоненты |
| `pinia` | Управление состоянием |
| `@pinia/nuxt` | Интеграция Pinia с Nuxt |
| `axios` | HTTP-клиент |
| `@nuxtjs/i18n` | Интернационализация |
| `@nuxtjs/fonts` | Шрифты |
| `@nuxt/icon` | Иконки |
| `@nuxt/image` | Оптимизация изображений |
| `tailwindcss` (v4) | CSS-фреймворк |
