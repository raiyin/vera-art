# ErrorTracker

Централизованный сбор, анализ и устранение ошибок во фронтенд‑ (Nuxt.js) и бэкенд‑части (Go) приложения.

## Назначение

Скилл обеспечивает:
* централизованный сбор ошибок с фронтенда и бэкенда;
* автоматическую классификацию и приоритизацию ошибок;
* интеграцию с системами мониторинга и трекинга задач;
* генерацию отчётов по динамике ошибок;
* проактивные уведомления о критических проблемах.

## Функциональные возможности

### 1. Сбор ошибок
* сбор ошибок с Nuxt.js (клиентские ошибки, SSR‑ошибки);
* сбор ошибок из Go‑бэкенда (паники, HTTP‑ошибки, ошибки БД);
* захват стектрейсов и контекста выполнения;
* сбор метаданных (окружение, версия кода, пользователь).

### 2. Классификация и приоритизация
* группировка похожих ошибок (кластеризация по стектрейсу);
* определение критичности (критические, важные, информационные);
* автоматическое определение частоты возникновения;
* тегирование по модулям/сервисам.

### 3. Анализ и диагностика
* поиск корневых причин (root cause analysis);
* корреляция ошибок с деплоями и изменениями кода;
* обнаружение новых ошибок (new issues detection);
* выявление повторяющихся проблем.

### 4. Интеграция и автоматизация
* создание тикетов в Jira/GitHub Issues для критических ошибок;
* уведомления в Slack/Telegram/Email;
* интеграция с системами мониторинга (Sentry, Prometheus, Logstash);
* автоматические действия (перезапуск сервисов, откат деплоя).

### 5. Отчётность
* ежедневные/еженедельные отчёты по ошибкам;
* дашборды с метриками качества (MTBF, MTTR);
* тренды по типам ошибок;
* отчёты по эффективности исправлений.

## Команды скилла

| Команда | Описание | Пример использования |
|-------|-------|-------|
| `collect:errors` | Сбор ошибок за период | `skill collect:errors --from=24h` |
| `analyze:errors` | Анализ и кластеризация | `skill analyze:errors --cluster=stacktrace` |
| `report:errors` | Генерация отчёта | `skill report:errors --format=pdf --period=week` |
| `notify:alerts` | Отправка уведомлений | `skill notify:alerts --severity=critical` |
| `ticket:create` | Создание тикетов для ошибок | `skill ticket:create --priority=high` |
| `trend:analyze` | Анализ трендов ошибок | `skill trend:analyze --metric=frequency` |
| `clear:resolved` | Очистка обработанных ошибок | `skill clear:resolved --older=7d` |

## Конфигурация

Пример файла `error-tracker.config.json`:

```json
{
  "errorSources": {
    "frontend": "nuxt-client-errors",
    "backend": "go-service-logs",
    "database": "pg-errors"
  },
  "severityRules": {
    "critical": ["panic", "5xx", "database-deadlock"],
    "high": ["4xx", "timeout", "validation-error"],
    "medium": ["warning", "deprecated-api"]
  },
  "notificationChannels": {
    "slackWebhook": "https://...",
    "telegramBot": "token",
    "email": ["dev@company.com"]
  },
  "autoActions": {
    "criticalErrors": "create-ticket",
    "recurringErrors": "notify-team",
    "newErrors": "alert-lead"
  },
  "retentionPeriod": "30d",
  "samplingRate": 1.0
}
