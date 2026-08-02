# Feature Flag: REGISTRATION_ENABLED

## Назначение

Feature-флаг `registration_enabled` управляет доступностью регистрации и связанных элементов интерфейса. Когда флаг **выключен**:

- **Сервер**: эндпоинт `POST /register` не регистрируется в роутере — регистрация недоступна
- **Клиент**: в хедере скрываются ссылка «Услуги» и кнопка входа (авторизации)

## Принцип работы

- **Сервер**: флаг читается из `config.yaml` при каждом запуске — меняйте конфиг и перезапускайте бинарник
- **Клиент**: флаг читается из `runtimeConfig.public.registrationEnabled`. Значение по умолчанию задаётся в `nuxt.config.ts`. В **SSR-режиме** (по умолчанию) переопределяется переменной окружения `NUXT_PUBLIC_REGISTRATION_ENABLED` на сервере без пересборки

## Разработка

### Регистрация включена (по умолчанию)

```bash
# Сервер
cd server && go run cmd/server/main.go

# Клиент
cd client && pnpm dev
```

### Регистрация выключена

```bash
# Сервер
# В server/config.yaml установить features.registration_enabled: false
cd server && go run cmd/server/main.go

# Клиент
# В client/.env установить NUXT_PUBLIC_REGISTRATION_ENABLED=false
cd client && pnpm dev
```

## Продакшен

### Сборка сервера

```bash
cd server
go build -o bin/artserver cmd/server/main.go
```

На целевом сервере разместить `config.yaml` рядом с бинарником, при необходимости выставить флаг:

```yaml
features:
  registration_enabled: false
```

Запуск:

```bash
./bin/artserver
```

### Сборка клиента (Nuxt SSR)

```bash
cd client
pnpm build
pnpm preview
```

Чтобы выключить регистрацию в рантайме без пересборки, передать переменную окружения при запуске `node .output/server/index.mjs`:

```bash
NUXT_PUBLIC_REGISTRATION_ENABLED=false node .output/server/index.mjs
```

### Сборка клиента (static/SPA)

Если проект собран в статику (`npx nuxi generate`), значение флага фиксируется на момент сборки. Переопределить его без пересборки нельзя.

```bash
cd client
NUXT_PUBLIC_REGISTRATION_ENABLED=false pnpm generate
```

## Где находится

| Компонент | Файл | Что делает |
|---|---|---|
| Сервер: конфиг | `server/config.yaml` | Значение флага |
| Сервер: структура | `server/internal/config/config.go` | Поле `FeaturesConfig.RegistrationEnabled` |
| Сервер: роутер | `server/internal/router/router.go:42` | Условная регистрация `POST /register` |
| Сервер: точка входа | `server/cmd/server/main.go` | Передача `Features` в `NewRouter`, лог при старте |
| Клиент: .env | `client/.env` | Значение флага для разработки |
| Клиент: конфиг Nuxt | `client/nuxt.config.ts:27` | Проброс в `runtimeConfig.public` |
| Клиент: хедер | `client/app/components/Header.vue` | Условное отображение «Услуги» и кнопки входа |
