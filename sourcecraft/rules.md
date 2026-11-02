# Правила разработки (Nuxt.js + Go)

## Стандарты кода

### Для Go (бэкенд)

* **Форматирование**: `gofmt` или `goimports`. Обязательно перед каждым коммитом.
* **Именование**:
  * `camelCase` для неэкспортируемых переменных и полей структур.
  * `PascalCase` для экспортируемых типов, функций, методов, полей структур.
  * `snake_case` для URL-параметров, query-параметров, JSON-полей (теги `json:"..."`).
  * `UPPER_SNAKE_CASE` для констант.
* **Пакеты**: имена пакетов — короткие, односложные, без подчёркиваний (`handler`, `service`, `domain`, `dto`, `port`, `config`, `router`).
* **Комментарии**: документирование всех публичных функций, типов, констант. Комментарии начинаются с имени сущности: `// Register создаёт нового пользователя`.
* **Ошибки**: все возвращаемые ошибки должны быть обработаны. Использование `_` для игнорирования ошибок запрещено.
* **Кастомные ошибки**: использовать пакет [`domain`](server/internal/domain/errors.go) с предопределёнными ошибками (`ErrNotFound`, `ErrUnauthorized`, `ErrForbidden`, `ErrConflict`, `ErrValidation`, `ErrInternal`). Маппинг в HTTP-статусы через `apperror` в [`middleware.go`](server/internal/handler/middleware.go).
* **Тесты**: покрытие критических путей ≥ 80%. Тесты для handler'ов через httptest, для service/repository — через mock-интерфейсы.
* **Архитектура**: строгое соблюдение слоёв: Handler → Service (порт-интерфейс) → Repository (порт-интерфейс). Handler не вызывает Repository напрямую.
* **Интерфейсы**: определять в пакете [`port`](server/internal/port/). Реализации — в пакете [`service`](server/internal/service/) и в репозиториях (SQLite/File).
* **DTO**: все структуры запросов/ответов — в пакете [`dto`](server/internal/dto/). Валидация входных данных — в handler'ах перед вызовом сервиса.
* **Конфигурация**: через Viper + YAML ([`config.yaml`](server/config.yaml)) + `.env` (godotenv). Структура конфига — в [`config.go`](server/internal/config/config.go).
* **Миграции БД**: инлайн SQL в [`main.go`](server/cmd/server/main.go) (ALTER TABLE ADD COLUMN). Новые миграции добавлять туда же, с проверкой `IF NOT EXISTS`.
* **SQLite**: WAL-режим, `busy_timeout=5000ms`, `foreign_keys=ON`, максимум 1 соединение. Все запросы через `database/sql`.
* **JWT**: использовать [`pkg/jwt`](server/pkg/jwt/) (golang-jwt/jwt/v5). Access token — 7 дней, refresh token — 30 дней. Refresh token хранится в БД (таблица `refresh_tokens`).
* **Логирование**: через `log/slog` (structured logging).
* **Rate Limiting**: встроенный in-memory limiter (5 запросов/мин для auth endpoints). При превышении — HTTP 429.
* **CORS**: через gin-contrib/cors-middleware. Настройки из конфига.
* **Загрузка файлов**: локальная файловая система. Пути: `{upload_dir}/images/`, `{upload_dir}/avatars/`, `{upload_dir}/news/`. Доступ через статический роут `/uploads/`.

### Для Nuxt.js (фронтенд)

* **TypeScript**: строгая типизация. Запрет `any`. Все API-ответы должны иметь интерфейс в [`types/`](client/app/types/).
* **ESLint**: соблюдение правил из [`eslint.config.mjs`](client/eslint.config.mjs).
* **Компоненты**: однофайловые компоненты (`.vue`), Composition API + `<script setup lang="ts">`. Чёткое разделение: шаблон, скрипт, стили.
* **Стилизация**: Tailwind CSS v4 через Nuxt UI. Кастомные стили — в [`assets/css/`](client/app/assets/css/). Scoped-стили для изолированных компонентов.
* **Состояние**: Pinia stores в [`stores/`](client/app/stores/). Каждый store — отдельный файл, Composition API (`defineStore('name', () => {...})`).
* **API-клиент**:
  * Для публичных и auth-запросов: Axios-инстанс из [`api/auth.ts`](client/app/api/auth.ts) с перехватчиками для инъекции токена и автоматического refresh'а при 401.
  * Для админ-запросов: Axios-инстанс [`api/admin.ts`](client/app/api/admin.ts) с ручным управлением refresh'ом токена.
  * Для справочных данных: fAxios-инстансы [`api/master-classes.ts`](client/app/api/master-classes.ts), [`api/requests.ts`](client/app/api/requests.ts).
* **Маршрутизация**: файловая маршрутизация Nuxt. Middleware для защиты админ-роутов — [`middleware/admin-auth.ts`](client/app/middleware/admin-auth.ts).
* **Интернационализация (i18n)**: через `@nuxtjs/i18n`. Локали: `ru` (по умолчанию), `en`. Файлы переводов — в [`i18n/locales/`](client/i18n/locales/). Конфигурация — [`i18n/i18n.config.ts`](client/i18n/i18n.config.ts).
* **Тема**: светлая/тёмная через [`ThemeStore`](client/app/stores/ThemeStore.ts) с localStorage + cross-tab sync.
* **JWT на клиенте**: декодирование через [`utils/jwt.ts`](client/app/utils/jwt.ts). Извлечение роли и userId из payload. Проверка expiry.
* **Уведомления**: polling чата каждые 30 секунд через [`NotificationStore`](client/app/stores/NotificationStore.ts).
* **Cookie Consent**: GDPR-согласие через [`CookieConsentStore`](client/app/stores/CookieConsentStore.ts). Аналитические и маркетинговые cookie разделены.
* **Сборка**: pnpm. Скрипты: `dev`, `build`, `generate`, `preview`, `lint`, `typecheck`.

## Правила ревью кода

* Обязательное ревью PR перед мерджем.
* Минимальное количество ревьюеров: 1.
* Проверка:
  * соответствия стандартам кода (Go: `gofmt`, `go vet`, `staticcheck`; Nuxt: ESLint, TypeScript strict);
  * покрытия тестами (Go: `go test ./...`);
  * отсутствия регрессий API (проверка обратной совместимости эндпоинтов);
  * обработки ошибок (все пути возврата ошибок);
  * безопасности (SQL-инъекции, XSS, CSRF, валидация входных данных);
  * i18n (все строки для пользователя через `$t()` или локализованные поля `name_ru`/`name_en`);
  * отсутствия захардкоженных URL и секретов (использовать конфиг / runtimeConfig / .env).

## Соглашения по коммитам

* Формат: [Conventional Commits](https://www.conventionalcommits.org/).
* Типы: `feat:`, `fix:`, `docs:`, `chore:`, `refactor:`, `test:`, `style:`, `perf:`.
* Пример: `feat(api): add user authentication endpoints`.
* Пример: `fix(shop): validate promo code expiry date`.
* Пример: `refactor(auth): extract JWT manager to pkg/jwt`.

## Процессы

* **API-изменения**: при добавлении/изменении эндпоинта — обновить DTO в [`dto/`](server/internal/dto/) и типы в [`types/`](client/app/types/).
* **Миграции БД**: добавлять инлайн SQL в [`main.go`](server/cmd/server/main.go) с проверкой `IF NOT EXISTS`. Обновлять [`db/db-schema.md`](server/db/db-schema.md).
* **Перед деплоем**: `go vet ./...`, `go test ./...`, `pnpm typecheck`, `pnpm lint`.
* **Критические ошибки**: приоритет «высокий», срок устранения ≤ 4 ч.
* **Документация**: все публичные эндпоинты должны быть описаны. Рекомендуется OpenAPI/Swagger (на данный момент отсутствует).
* **Безопасность**: не коммитить `.env` файлы, токены, пароли. Использовать переменные окружения через godotenv.
