export interface Review {
    id: number
    user_id: number
    product_id: number
    purchase_id: number
    rating: number
    title_ru: string | null
    title_en: string | null
    comment_ru: string | null
    comment_en: string | null
    is_approved: boolean
    is_visible: boolean
    created_at: string
    updated_at: string
    status?: string // computed status: 'pending', 'approved', 'rejected'

    // Joined fields
    user?: {
        id: number
        username: string
        full_name: string | null
        email: string | null
    }
    product?: {
        id: number
        title_ru: string
        title_en: string
    }
    purchase?: {
        id: number
        product_id: number
    }
}

export interface CreateReviewDto {
    rating: number
    title_ru?: string
    title_en?: string
    comment_ru?: string
    comment_en?: string
}

export interface UpdateReviewDto {
    rating?: number
    title_ru?: string
    title_en?: string
    comment_ru?: string
    comment_en?: string
}
