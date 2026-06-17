# Client-Side Refactoring Plan

## Overview

This plan addresses all identified issues in the Nuxt.js 4 client codebase. The refactoring is organized into **phases** with clear dependencies, ensuring each phase can be implemented and tested independently.

## Current Architecture Issues

| # | Issue | Severity | Files Affected |
|---|-------|----------|----------------|
| 1 | Mixed API clients (axios vs fetch) | High | `api/*`, `stores/*`, `components/PicGallery.vue` |
| 2 | No centralized HTTP client | High | All API consumers |
| 3 | Module-level `useRuntimeConfig()` calls | High | `api/requests.ts`, `stores/MaterialStore.ts` |
| 4 | Duplicated URLSearchParams / error handling | Medium | `api/admin.ts`, `api/master-classes.ts` |
| 5 | Inline admin types in `api/admin.ts` | Medium | `api/admin.ts` |
| 6 | Options API in `news/[id].vue` | Medium | `pages/news/[id].vue` |
| 7 | No composables directory | Medium | All pages with duplicated logic |
| 8 | Duplicated infinite scroll logic | Medium | `pages/gallery/index.vue`, `pages/art-store/index.vue` |
| 9 | Massive pages with inline CSS | Medium | `pages/index.vue` (1811 lines), `pages/admin/index.vue` (1328 lines) |
| 10 | Inconsistent import aliases | Low | Multiple files |
| 11 | Overlapping type files | Low | `types/base_work.ts`, `types/common_work.ts`, `types/work.ts`, `types/sale.ts` |
| 12 | PicGallery uses raw fetch for DELETE | Medium | `components/PicGallery.vue` |

## Target Directory Structure

```
client/app/
├── api/
│   ├── http-client.ts          # NEW: Centralized Axios instance
│   ├── auth.ts                 # REFACTOR: Use http-client
│   ├── admin.ts                # REFACTOR: Use http-client, extract types
│   ├── gallery.ts              # NEW: Gallery/works endpoints
│   ├── shop.ts                 # NEW: Shop/sales/products endpoints
│   ├── news.ts                 # NEW: News endpoints (replaces requests.ts)
│   ├── master-classes.ts       # REFACTOR: Use http-client
│   ├── chat.ts                 # NEW: Chat endpoints
│   ├── reviews.ts              # NEW: Reviews endpoints
│   ├── payments.ts             # NEW: Payments endpoints
│   └── learning.ts             # NEW: Learning/lessons endpoints
├── composables/                # NEW directory
│   ├── useApi.ts               # Composable wrapper for API calls
│   ├── useInfiniteScroll.ts    # Extract from gallery/art-store
│   ├── usePagination.ts        # Pagination logic
│   ├── useFormatting.ts        # Date, price, text formatting
│   └── useErrorHandler.ts      # Centralized error handling
├── plugins/                    # NEW directory
│   └── http-client.ts          # Nuxt plugin to provide http-client
├── types/
│   ├── index.ts                # Re-exports
│   ├── base.ts                 # Base interfaces
│   ├── work.ts                 # Work/gallery types
│   ├── sale.ts                 # Sale types
│   ├── product.ts              # Product types (from admin.ts inline)
│   ├── user.ts                 # User/auth types (from admin.ts inline)
│   ├── news.ts                 # News types (from admin.ts inline)
│   ├── lesson.ts               # Lesson types (from admin.ts inline)
│   ├── payment.ts              # Payment/purchase types (from admin.ts inline)
│   ├── promo-code.ts           # Promo code types (from admin.ts inline)
│   ├── chat.ts                 # Chat types
│   ├── review.ts               # Review types
│   ├── master-class.ts         # Master class types
│   ├── material.ts             # Material types
│   ├── category.ts             # Category types (from admin.ts inline)
│   ├── tag.ts                  # Tag types (from admin.ts inline)
│   └── dashboard.ts            # Dashboard stats types (from admin.ts inline)
├── stores/
│   ├── AuthStore.ts            # KEEP: Well-structured
│   ├── ProductStore.ts         # REFACTOR: Use http-client
│   ├── MaterialStore.ts        # REFACTOR: Use http-client
│   ├── NotificationStore.ts    # REFACTOR: Use http-client
│   ├── ThemeStore.ts           # KEEP: Clean
│   └── CookieConsentStore.ts   # KEEP: Clean
├── components/
│   ├── PicGallery.vue          # REFACTOR: Use http-client for DELETE
│   ├── ... (other components)  # Minor import fixes
│   └── admin/                  # KEEP: Well-structured
├── pages/
│   ├── index.vue               # REFACTOR: Extract sections, reduce size
│   ├── gallery/index.vue       # REFACTOR: Use useInfiniteScroll composable
│   ├── art-store/index.vue     # REFACTOR: Use useInfiniteScroll composable
│   ├── news/[id].vue           # REFACTOR: Composition API
│   └── admin/index.vue         # REFACTOR: Extract dashboard components
└── utils/
    └── jwt.ts                  # KEEP: Clean utility
```

## Phase 1: Foundation — Centralized HTTP Client

**Goal**: Create a single Axios-based HTTP client with interceptors that all API modules and stores will use.

### Step 1.1: Create `api/http-client.ts`

Create a centralized Axios instance with:
- Base URL from `runtimeConfig.public.serverUrl`
- Request interceptor: inject `Authorization: Bearer` header from AuthStore
- Response interceptor: on 401, attempt token refresh, retry original request
- Unified error handling that extracts meaningful error messages
- Type-safe response wrapper

```typescript
// api/http-client.ts
import axios, { type AxiosInstance, type AxiosError, type InternalAxiosRequestConfig } from 'axios';

interface TokenRefreshResponse {
    access_token: string;
    refresh_token: string;
}

let isRefreshing = false;
let failedQueue: Array<{
    resolve: (token: string) => void;
    reject: (error: unknown) => void;
}> = [];

function processQueue(error: unknown, token: string | null = null): void {
    failedQueue.forEach((promise) => {
        if (error) {
            promise.reject(error);
        } else {
            promise.resolve(token!);
        }
    });
    failedQueue = [];
}

export function createHttpClient(): AxiosInstance {
    const config = useRuntimeConfig();
    const baseURL = config.public.serverUrl as string;

    const client = axios.create({
        baseURL,
        timeout: 15000,
        headers: { 'Content-Type': 'application/json' },
    });

    // Request interceptor: attach auth token
    client.interceptors.request.use(
        (config: InternalAxiosRequestConfig) => {
            const authStore = useAuthStore();
            if (authStore.accessToken) {
                config.headers.Authorization = `Bearer ${authStore.accessToken}`;
            }
            return config;
        },
        (error) => Promise.reject(error),
    );

    // Response interceptor: handle 401 with token refresh
    client.interceptors.response.use(
        (response) => response,
        async (error: AxiosError) => {
            const originalRequest = error.config as InternalAxiosRequestConfig & { _retry?: boolean };

            if (error.response?.status === 401 && !originalRequest._retry) {
                if (isRefreshing) {
                    return new Promise((resolve, reject) => {
                        failedQueue.push({ resolve, reject });
                    }).then((token) => {
                        originalRequest.headers.Authorization = `Bearer ${token}`;
                        return client(originalRequest);
                    });
                }

                originalRequest._retry = true;
                isRefreshing = true;

                try {
                    const authStore = useAuthStore();
                    const response = await axios.post<TokenRefreshResponse>(
                        `${baseURL}auth/refresh`,
                        { refresh_token: authStore.refreshToken },
                    );
                    const { access_token, refresh_token } = response.data;
                    authStore.saveTokens({ access_token, refresh_token, expiry: calculateExpiry(access_token) });
                    processQueue(null, access_token);
                    originalRequest.headers.Authorization = `Bearer ${access_token}`;
                    return client(originalRequest);
                } catch (refreshError) {
                    processQueue(refreshError, null);
                    const authStore = useAuthStore();
                    authStore.clearTokens();
                    if (import.meta.client) {
                        navigateTo('/auth/login');
                    }
                    return Promise.reject(refreshError);
                } finally {
                    isRefreshing = false;
                }
            }

            return Promise.reject(error);
        },
    );

    return client;
}

// Singleton pattern
let httpClientInstance: AxiosInstance | null = null;

export function getHttpClient(): AxiosInstance {
    if (!httpClientInstance) {
        httpClientInstance = createHttpClient();
    }
    return httpClientInstance;
}
```

**Note**: `useRuntimeConfig()` inside `createHttpClient()` works because it's called lazily (not at module level).

### Step 1.2: Create `plugins/http-client.ts`

```typescript
// plugins/http-client.ts
import { getHttpClient } from '~/api/http-client';

export default defineNuxtPlugin(() => {
    const httpClient = getHttpClient();
    return {
        provide: {
            httpClient,
        },
    };
});
```

### Step 1.3: Create `composables/useApi.ts`

```typescript
// composables/useApi.ts
import { getHttpClient } from '~/api/http-client';

export function useApi() {
    const client = getHttpClient();

    async function get<T>(url: string, params?: Record<string, unknown>): Promise<T> {
        const response = await client.get<T>(url, { params });
        return response.data;
    }

    async function post<T>(url: string, data?: unknown): Promise<T> {
        const response = await client.post<T>(url, data);
        return response.data;
    }

    async function put<T>(url: string, data?: unknown): Promise<T> {
        const response = await client.put<T>(url, data);
        return response.data;
    }

    async function del<T>(url: string): Promise<T> {
        const response = await client.delete<T>(url);
        return response.data;
    }

    return { get, post, put, del };
}
```

## Phase 2: Extract Admin Types

**Goal**: Move all inline interfaces from `api/admin.ts` to dedicated type files.

### Step 2.1: Create type files for admin entities

Extract these interfaces from `api/admin.ts`:

| Current Location | Target File |
|-----------------|-------------|
| `DashboardStats`, `SalesByMonthEntry`, `PopularCategoryEntry`, `RecentActivityItem` | `types/dashboard.ts` |
| `AdminWorkItem`, `AdminListWorksResponse` | `types/work.ts` |
| `AdminSaleItem`, `AdminListSalesResponse` | `types/sale.ts` |
| `AdminNewsItem`, `AdminListNewsResponse` | `types/news.ts` |
| `AdminProductItem`, `AdminListProductsResponse` | `types/product.ts` |
| `AdminLessonItem`, `AdminListLessonsResponse` | `types/lesson.ts` |
| `AdminUserItem`, `AdminListUsersResponse`, `AdminUserPurchase`, `AdminUserReview`, `AdminUserDetail` | `types/user.ts` |
| `AdminReviewItem`, `AdminReviewsStats`, `AdminListReviewsResponse` | `types/review.ts` |
| `AdminPurchaseItem`, `AdminListPurchasesResponse` | `types/payment.ts` |
| `AdminPaymentItem`, `AdminListPaymentsResponse` | `types/payment.ts` |
| `AdminPromoCodeItem` | `types/promo-code.ts` |
| `AdminChatThread`, `AdminChatMessage` | `types/chat.ts` |
| `AdminCategoryItem` | `types/category.ts` |
| `AdminTagItem` | `types/tag.ts` |

### Step 2.2: Refactor `api/admin.ts`

After extracting types, refactor `admin.ts` to:
- Use `getHttpClient()` instead of raw `fetch`
- Remove duplicated `getAuthHeaders()`, `tryRefreshToken()`, `fetchApi()`, `postApi()` helpers
- Remove inline type definitions (import from `~/types/` instead)
- Use consistent endpoint paths (all with `admin/` prefix)

## Phase 3: Refactor API Modules

**Goal**: Convert all API modules to use the centralized HTTP client.

### Step 3.1: Refactor `api/auth.ts`

Current: Uses axios directly with its own interceptor logic.
Target: Use `getHttpClient()` and remove duplicate interceptor logic.

Changes:
- Remove standalone `api` axios instance creation
- Remove duplicate interceptor logic (now in http-client.ts)
- Use `getHttpClient()` for all requests
- Keep the same exported interface for backward compatibility

### Step 3.2: Refactor `api/master-classes.ts`

Current: Uses raw fetch with duplicated `getServerUrl()`.
Target: Use `getHttpClient()`.

Changes:
- Remove `getServerUrl()` function
- Replace all `fetch()` calls with `getHttpClient()` calls
- Keep function signatures the same

### Step 3.3: Refactor `api/requests.ts`

Current: Uses axios but calls `useRuntimeConfig()` at module level.
Target: Use `getHttpClient()` (lazy initialization).

Changes:
- Remove module-level `useRuntimeConfig()` call
- Use `getHttpClient()` inside functions
- Rename to `api/news.ts` for consistency

### Step 3.4: Create new API modules

Create dedicated API modules for domains that currently lack them:

- `api/gallery.ts` — Gallery/works endpoints (currently scattered)
- `api/shop.ts` — Shop/products/sales endpoints
- `api/chat.ts` — Chat endpoints
- `api/reviews.ts` — Reviews endpoints
- `api/payments.ts` — Payments endpoints
- `api/learning.ts` — Learning/lessons endpoints

## Phase 4: Refactor Stores

**Goal**: All stores use the centralized HTTP client instead of raw fetch or direct axios.

### Step 4.1: Refactor `stores/ProductStore.ts`

Current: Uses `api.getApiInstance()` directly.
Target: Use `getHttpClient()` or `useApi()` composable.

Changes:
- Replace `api.getApiInstance()` with `getHttpClient()`
- Remove `extractErrorMessage()` (use centralized error handling)
- Keep store structure and state management

### Step 4.2: Refactor `stores/MaterialStore.ts`

Current: Uses raw fetch with module-level `useRuntimeConfig()`.
Target: Use `getHttpClient()`.

Changes:
- Remove module-level `useRuntimeConfig()` call
- Replace `fetch()` with `getHttpClient().get()`
- Keep `getMaterialName()`, `getBaseName()` helper methods

### Step 4.3: Refactor `stores/NotificationStore.ts`

Current: Uses raw fetch with inline URL construction.
Target: Use `getHttpClient()`.

Changes:
- Replace `fetch()` with `getHttpClient().get()`
- Keep polling logic and state management

## Phase 5: Create Composables

**Goal**: Extract reusable logic from pages into composables.

### Step 5.1: `composables/useInfiniteScroll.ts`

Extract from `pages/gallery/index.vue` and `pages/art-store/index.vue`.

```typescript
export function useInfiniteScroll(fetchFn: (page: number) => Promise<unknown[]>) {
    const items = ref<unknown[]>([]);
    const page = ref(1);
    const loading = ref(false);
    const hasMore = ref(true);
    const observer = ref<IntersectionObserver | null>(null);

    const loadMore = async () => {
        if (loading.value || !hasMore.value) return;
        loading.value = true;
        try {
            const newItems = await fetchFn(page.value);
            if (newItems.length === 0) {
                hasMore.value = false;
            } else {
                items.value.push(...newItems);
                page.value++;
            }
        } finally {
            loading.value = false;
        }
    };

    const setupObserver = (element: Ref<HTMLElement | null>) => {
        // IntersectionObserver logic
    };

    return { items, loading, hasMore, loadMore, setupObserver };
}
```

### Step 5.2: `composables/useFormatting.ts`

Extract formatting helpers used across components.

### Step 5.3: `composables/useErrorHandler.ts`

Centralized error handling with toast notifications.

## Phase 6: Refactor Pages

**Goal**: Reduce page sizes, use Composition API consistently, extract duplicated logic.

### Step 6.1: Refactor `pages/news/[id].vue`

Current: Options API, uses `axios` directly, imports from `vue-router`.
Target: Composition API, uses `getHttpClient()`, uses Nuxt's `useRoute()`.

Changes:
- Convert `<script lang="ts">` to `<script setup lang="ts">`
- Replace `this.$route` with `useRoute()`
- Replace `this.$axios` with `getHttpClient()`
- Replace `this.$t()` with `useI18n().t`
- Extract large template sections into child components

### Step 6.2: Refactor `pages/gallery/index.vue` and `pages/art-store/index.vue`

Current: Nearly identical infinite scroll logic duplicated.
Target: Both use `useInfiniteScroll` composable.

Changes:
- Extract infinite scroll logic to composable
- Keep page-specific template and styling
- Use `getHttpClient()` for API calls

### Step 6.3: Refactor `pages/index.vue`

Current: 1811 lines with extensive inline CSS, mixed theme toggle, parallax, stats animation.
Target: Extract sections into components.

Changes:
- Extract hero/parallax section → `components/HeroSection.vue`
- Extract stats/animation section → `components/StatsSection.vue`
- Extract gallery preview → reuse `PicGallery.vue` or create `GalleryPreview.vue`
- Keep page as orchestrator of sections

### Step 6.4: Refactor `pages/admin/index.vue`

Current: 1328 lines with inline CSS, stat cards, charts, navigation, activity feed.
Target: Extract dashboard sections into components.

Changes:
- Extract stat cards → `components/admin/AdminStatCard.vue`
- Extract charts section → `components/admin/AdminChartsSection.vue`
- Extract activity feed → `components/admin/AdminActivityFeed.vue`
- Extract quick navigation → `components/admin/AdminQuickNav.vue`

## Phase 7: Component Cleanup

**Goal**: Fix component-level issues.

### Step 7.1: Refactor `components/PicGallery.vue`

Current: Uses raw `fetch()` for DELETE operations.
Target: Use `getHttpClient()`.

Changes:
- Replace `fetch()` DELETE with `getHttpClient().delete()`
- Remove direct `localStorage.getItem('token')` usage (use AuthStore)

### Step 7.2: Fix import aliases

Standardize on `~/` prefix for all imports (Nuxt convention):
- `@/components/` → `~/components/`
- Relative imports like `../stores/` → `~/stores/`

## Phase 8: Clean Up Types

**Goal**: Consolidate overlapping type files.

### Step 8.1: Consolidate work/sale types

Current: `base_work.ts`, `common_work.ts`, `work.ts`, `sale.ts` have overlapping definitions.
Target: Clean, non-overlapping type hierarchy.

```typescript
// types/work.ts
export interface WorkBase {
    id: string;
    str_id: string;
    dir: string;
    name_ru: string;
    name_en: string;
    year: string;
    descr: string;
    base_ru: string;
    base_en: string;
    width: number;
    height: number;
    type?: number;
    images: string[];
    materials_ids: number[];
    materials_en?: string[];
    materials_ru?: string[];
}

export interface Work extends WorkBase {
    __type: 'GetWorkDto';
}

// types/sale.ts
export interface Sale extends WorkBase {
    price: number;
    __type: 'GetSaleDto';
}

export type CommonWork = Work | Sale;
```

## Migration Strategy

### Order of Execution

```
Phase 1 (Foundation)
  └─► Phase 2 (Types)
       └─► Phase 3 (API Modules)
            └─► Phase 4 (Stores)
                 └─► Phase 5 (Composables)
                      └─► Phase 6 (Pages)
                           └─► Phase 7 (Components)
                                └─► Phase 8 (Type Cleanup)
```

Each phase depends on the previous one. Within each phase, steps can be done in parallel where they don't share dependencies.

### Testing Strategy

After each phase:
1. Run `pnpm typecheck` to verify TypeScript
2. Run `pnpm lint` to verify ESLint
3. Manually test affected pages in browser
4. Fix any issues before moving to next phase

### Rollback Plan

If a phase introduces regressions:
1. Git revert the specific phase commits
2. Document what went wrong
3. Adjust the plan and retry

## Risk Assessment

| Risk | Likelihood | Impact | Mitigation |
|------|-----------|--------|------------|
| Token refresh breaks during migration | Medium | High | Keep old auth.ts working until Phase 3 is complete and tested |
| Module-level useRuntimeConfig() calls missed | Low | Medium | Search for all `useRuntimeConfig()` calls at module level |
| Breaking changes to store interfaces | Low | Medium | Keep public store interfaces unchanged during refactoring |
| Infinite scroll regression | Medium | Medium | Test gallery and art-store pages thoroughly after composable extraction |
| Import alias changes break components | Low | Low | Use search-and-replace carefully, verify with typecheck |
