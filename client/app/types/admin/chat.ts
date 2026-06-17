// ─── Admin Chat Types ───────────────────────────────────────────────

export interface AdminChatThread {
    id: number
    purchase_id: number
    user_id: number
    admin_id: number | null
    last_message_at: string
    is_resolved: boolean
    created_at: string
    user: {
        username: string
        full_name: string
        email: string
    } | null
    purchase: {
        id: number
        product_id: number
        product: {
            title_ru: string
            title_en: string
        } | null
    } | null
}

export interface AdminChatMessage {
    id: number
    thread_id: number
    sender_id: number
    message_type: string
    content: string
    attachment_url: string | null
    attachment_size: number | null
    is_read: boolean
    read_at: string | null
    created_at: string
    sender: {
        username: string
        full_name: string
    } | null
}
