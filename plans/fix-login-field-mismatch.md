# Fix Login "Invalid Request Body" Error

## Root Cause

The login request fails because of a **field name mismatch** between the client and server.

### Client sends (`client/app/pages/auth/login.vue:265-268`)

```json
{
  "username": "someuser",
  "password": "somepass"
}
```

### Server expects (`server/internal/dto/auth.go:11-14`)

```go
type LoginUserRequest struct {
    Email    string `json:"email" binding:"required,email"`
    Password string `json:"password" binding:"required"`
}
```

The server's `LoginUserRequest` DTO expects an `email` field (with `binding:"required,email"` validation), but the client sends `username`. Since `email` is missing and fails validation, `c.ShouldBindJSON(&req)` returns an error → HTTP 400 with `INVALID_REQUEST`.

### Additional issue: Register page has the same pattern

The register page (`client/app/pages/auth/register.vue:500-504`) sends:

```json
{
  "username": "someuser",
  "password": "somepass",
  "email": "some@email.com"
}
```

But the server expects (`server/internal/dto/auth.go:4-8`):

```go
type RegisterUserRequest struct {
    Email    string `json:"email" binding:"required,email"`
    Password string `json:"password" binding:"required,min=8"`
    Name     string `json:"name" binding:"required"`
}
```

The server expects `name` but receives `username`. However, the register page also sends `email` which IS expected, so registration might partially work if the server ignores unknown fields (Gin's default behavior). But the `name` field would be empty, causing potential issues.

## Decision: Align client with server

The server is the source of truth — it uses `email` for login (not username) and `name` for registration (not username). The client should be fixed to match.

## Changes Required

### 1. `client/app/pages/auth/login.vue`

| Change | Location | Description |
|--------|----------|-------------|
| Rename field `username` → `email` in form | Lines 43-56 | Change `formState.username` to `formState.email`, update label/placeholder i18n keys, change input type to `email` |
| Update `LoginForm` interface | Lines 228-231 | Change `username: string` to `email: string` |
| Update `formState` initialization | Lines 238-241 | Change `username: ''` to `email: ''` |
| Update `handleLogin` validation | Line 253 | Check `formState.email` instead of `formState.username` |
| Update API call | Lines 265-268 | Send `email` instead of `username` |
| Update "remember me" logic | Lines 271-279 | Store/retrieve `email` instead of `username` |
| Update `onMounted` | Lines 357-363 | Restore remembered email |

### 2. `client/app/pages/auth/register.vue`

| Change | Location | Description |
|--------|----------|-------------|
| Update API call | Lines 500-504 | Send `name: formState.username` instead of `username: formState.username` |

### 3. `client/app/api/auth.ts`

| Change | Location | Description |
|--------|----------|-------------|
| Update `login` parameter type | Line 36 | Change `{ username: string, password: string }` to `{ email: string, password: string }` |

## Files to Modify

1. `client/app/pages/auth/login.vue` — main fix
2. `client/app/pages/auth/register.vue` — minor fix for `name` field
3. `client/app/api/auth.ts` — type alignment

## Verification

After changes, the login request should send:

```json
{
  "email": "user@example.com",
  "password": "somepass"
}
```

Which matches the server's `LoginUserRequest` DTO exactly.
