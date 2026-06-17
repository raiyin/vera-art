// ─── Public Types ───────────────────────────────────────────────────

export type { NewsDesc, } from './news-desc';
export type { NewsDescDto, } from './news-desc';
export type { Material, } from './material';
export type { SortOption, } from './sort-option';
export type { CreateWorkDto, TypedGetWorkDto as GetWorkDto, } from './work';
export type { CreateSaleDto, GetSaleDto, } from './sale';
export type { Base, } from './base';
export type { RequestResult, } from './request-result';
export type { UserPassPair, } from './userpasspair';
export type { CommonTypedGetWorkDto as CommonGetWorkDto, } from './common_work';
export type { GetBaseWork, } from './base_work';
export type { UpdateWorkResponse, UpdateWorkRequest, } from './work';
export type { UpdateSaleResponse, UpdateSaleRequest, } from './sale';
export type { Review, CreateReviewDto, UpdateReviewDto, } from './review';
export type { ChatThread, ChatMessage, CreateChatThreadDto, SendMessageDto, } from './chat';
export type { MasterClass, MasterClassTag, } from './master-class';

// ─── Dashboard Types ────────────────────────────────────────────────

export type { DashboardStats, SalesByMonthEntry, PopularCategoryEntry, RecentActivityItem, } from './dashboard';

// ─── Admin Types ────────────────────────────────────────────────────

export type {
    AdminWorkItem, AdminListWorksResponse,
    AdminSaleItem, AdminListSalesResponse,
    AdminNewsItem, AdminListNewsResponse,
} from './admin/gallery';

export type { AdminProductItem, AdminListProductsResponse, } from './admin/product';

export type { AdminLessonItem, AdminListLessonsResponse, } from './admin/lesson';

export type {
    AdminUserItem, AdminListUsersResponse,
    AdminUserPurchase, AdminUserReview, AdminUserDetail,
} from './admin/user';

export type { AdminReviewItem, AdminReviewsStats, AdminListReviewsResponse, } from './admin/review';

export type {
    AdminPurchaseItem, AdminListPurchasesResponse,
    AdminPaymentItem, AdminListPaymentsResponse,
} from './admin/payment';

export type { AdminPromoCodeItem, } from './admin/promo-code';

export type { AdminChatThread, AdminChatMessage, } from './admin/chat';

export type { AdminCategoryItem, } from './admin/category';

export type { AdminTagItem, } from './admin/tag';
