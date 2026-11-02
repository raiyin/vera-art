# APISyncExpert

Автоматизация синхронизации API‑контрактов между фронтендом (Nuxt.js) и бэкендом (Go), обеспечение согласованности взаимодействия компонентов.

## Назначение

Скилл обеспечивает:
* автоматическую генерацию TypeScript‑клиента для Nuxt.js на основе OpenAPI/Swagger‑спецификаций Go‑бэкенда;
* валидацию соответствия фронтенда и бэкенда API‑контракту;
* обнаружение breaking changes при изменениях API;
* автоматическое обновление документации и SDK;
* уведомления о несовместимости версий API.

## Функциональные возможности

### 1. Парсинг API‑спецификаций
* извлечение OpenAPI/Swagger‑схем из Go‑кода (с использованием `swag` или аналогичных инструментов);
* поддержка версий OpenAPI 3.0+;
* обработка аннотаций в Go‑коде (`// @openapi` и т. д.).

### 2. Генерация клиентского кода
* автоматическая генерация TypeScript‑клиента для Nuxt.js;
* поддержка Axios и Fetch API;
* генерация интерфейсов TypeScript для строгой типизации;
* обновление клиентского кода при изменениях API.

### 3. Валидация контрактов
* проверка соответствия запросов фронтенда API‑спецификации;
* валидация ответов бэкенда по схеме;
* обнаружение несоответствий типов данных;
* проверка HTTP‑методов и путей.

### 4. Управление версиями
* версионирование API (поддержка нескольких версий);
* отслеживание изменений между версиями;
* предупреждение о breaking changes (удаление полей, изменение типов);
* миграционные подсказки при критических изменениях.

### 5. Документация и отчётность
* автоматическая генерация интерактивной документации (Swagger UI, ReDoc);
* создание changelog API;
* отчёты о совместимости фронтенда и бэкенда;
* уведомления об устаревших эндпоинтах.

### 6. Интеграция и уведомления
* интеграция с системами контроля версий (GitHub, GitLab);
* уведомления в Slack/Telegram о критических изменениях API;
* автоматическая проверка совместимости при пулл‑реквестах.

## Команды скилла

| Команда | Описание | Пример использования |
|-------|-------|-------|
| `sync:api` | Синхронизация API‑контракта | `skill sync:api --source=go-backend` |
| `validate:api` | Валидация соответствия фронтенда и бэкенда | `skill validate:api --strict=true` |
| `generate:client` | Генерация TypeScript‑клиента | `skill generate:client --output=lib/api-client` |
| `docs:generate` | Создание документации API | `skill docs:generate --format=swagger-ui` |
| `changelog:create` | Создание changelog API | `skill changelog:create --version=v2.0` |
| `check:breaking` | Проверка на breaking changes | `skill check:breaking --base=v1.0 --target=v2.0` |
| `notify:changes` | Отправка уведомлений об изменениях | `skill notify:changes --channel=slack-dev` |

## Конфигурация

Пример файла `api-sync-expert.config.json`:

```json
{
  "apiSources": {
    "backend": {
      "type": "go",
      "path": "./backend/api",
      "specGenerator": "swag",
      "outputSpec": "openapi.yaml"
    },
    "frontend": {
      "type": "nuxt",
      "clientPath": "./frontend/lib/api-client",
      "template": "typescript-axios"
    }
  },
  "validation": {
    "strict": true,
    "checkRequests": true,
    "checkResponses": true
  },
  "breakingChanges": {
    "alertOnRemoval": true,
    "alertOnTypeChange": true,
    "alertOnRequiredField": true
  },
  "documentation": {
    "formats": ["swagger-ui", "redoc"],
    "outputPath": "./docs/api"
  },
  "notifications": {
    "slackWebhook": "https://...",
    "telegramBot": "token",
    "channels": ["dev-alerts"]
  },
  "versioning": {
    "strategy": "semantic",
    "currentVersion": "v1.0.0"
  }
}
