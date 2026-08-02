# CLAUDE.md — Art Server Project

## Project Overview

Nuxt.js 4 (frontend) + Go/Gin (backend) with SQLite. Art gallery, shop, courses, and admin platform.

## Context Files (read these first)

Before answering any query, read these files for full project context:

1. **sourcecraft/rules.md** — Coding standards, naming conventions, commit rules, review processes
2. **sourcecraft/arch.md** — System architecture, layer diagram, route tables, data flows, DB schema
3. **sourcecraft/requirements.md** — All functional/non-functional requirements, constraints, improvement recommendations

## Quick Reference

### Backend (Go)

- **Module**: `github.com/raiyin/artserver`
- **Entry**: `server/cmd/server/main.go`
- **Router**: `server/internal/router/router.go`
- **Handlers**: `server/internal/handler/` (auth, gallery, shop, learning, payment, chat, news, misc, profile)
- **Services**: `server/internal/service/` (user, shop, payment, news)
- **Ports**: `server/internal/port/` (14 service interfaces, 18 repository interfaces)
- **Domain**: `server/internal/domain/` (User, Work, Sale, Product, Lesson, Payment, Chat, etc.)
- **DTOs**: `server/internal/dto/` (9 files)
- **Config**: `server/internal/config/config.go` + `server/config.yaml`
- **DB**: SQLite3, WAL mode, single connection, inline migrations in main.go

### Frontend (Nuxt.js)

- **Entry**: `client/nuxt.config.ts`
- **Stores**: `client/app/stores/` (AuthStore, ProductStore, MaterialStore, NotificationStore, ThemeStore, CookieConsentStore)
- **API**: `client/app/api/` (auth.ts — Axios, admin.ts — fetch, master-classes.ts, requests.ts)
- **Types**: `client/app/types/`
- **Pages**: file-based routing in `client/app/pages/`
- **i18n**: `client/i18n/locales/` (ru.json, en.json)
- **Components**: `client/app/components/` + `client/app/components/admin/`

### Key Architecture

- Handler → Service (port interface) → Repository (port interface) → SQLite/File
- JWT: access token (15min) + refresh token (7d, stored in DB)
- Payments: YooKassa with webhook handling
- Admin: role-based middleware, panel at /admin
- Theme: light/dark with localStorage + cross-tab sync
- Chat: polling every 30s via NotificationStore

### Commands

- Server: `cd server && go run cmd/server/main.go`
- Client: `cd client && pnpm dev`
- Lint: `cd client && pnpm lint`
- Typecheck: `cd client && pnpm typecheck`
- Go vet: `cd server && go vet ./...`
- Go test: `cd server && go test ./...`
- Client test: `cd client && pnpm test`
- Client test watch: `cd client && pnpm test:watch`

### Documentation Update

When code changes, update sourcecraft/ files:

```bash
make update-docs   # Shows what changed and prompts AI to update docs
make check-docs    # Checks if docs need updating
```

GitHub Actions automatically:
- Comments on PRs when code changes without docs updates
- Creates an issue on merge to main if sourcecraft/ wasn't updated
