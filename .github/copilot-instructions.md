# GitHub Copilot Instructions — Art Server Project

## Project Context

This is a **Nuxt.js 4 (frontend) + Go/Gin (backend)** project with SQLite database.

### Always Read These Files for Context

Before generating code or answering questions, read:
- [`sourcecraft/rules.md`](sourcecraft/rules.md) — coding standards, naming conventions, development processes
- [`sourcecraft/arch.md`](sourcecraft/arch.md) — system architecture, layers, data flows, component structure
- [`sourcecraft/requirements.md`](sourcecraft/requirements.md) — functional and non-functional requirements

### Tech Stack

**Backend (Go):**
- Framework: Gin (gin-gonic/gin)
- Database: SQLite3 (mattn/go-sqlite3), WAL mode, single connection
- Auth: JWT (golang-jwt/jwt/v5), access 15min + refresh 7d
- Payments: YooKassa
- Config: Viper + YAML + .env
- Architecture: Handler → Service (port interface) → Repository (port interface)

**Frontend (Nuxt.js):**
- Framework: Nuxt 4
- State: Pinia stores
- API: Axios (auth) + fetch (admin)
- Styling: Tailwind CSS v4 via Nuxt UI
- i18n: Russian (default), English
- Package Manager: pnpm

### Key Conventions

- Go: camelCase vars, PascalCase exports, snake_case JSON tags
- Nuxt: Composition API, `<script setup lang="ts">`, strict TypeScript
- Errors: domain.ErrNotFound, ErrUnauthorized, ErrForbidden, ErrConflict, ErrValidation
- DTOs: server/internal/dto/ for all request/response structs
- Types: client/app/types/ for all TypeScript interfaces
