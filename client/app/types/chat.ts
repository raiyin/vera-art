export interface ChatThread {
    id: number
    purchase_id: number
    user_id: number
    admin_id: number | null
    last_message_at: string
    is_resolved: boolean
    created_at: string

    // Joined fields
    purchase?: {
        id: number
        product_id: number
        product?: {
            title_ru: string
            title_en: string
        }
    }
    user?: {
        id: number
        username: string
        full_name: string | null
    }
    admin?: {
        id: number
        username: string
        full_name: string | null
    }
    messages?: ChatMessage[]
}

export interface ChatMessage {
    id: number
    thread_id: number
    sender_id: number
    message_type: 'text' | 'image' | 'file'
    content: string
    attachment_url: string | null
    attachment_size: number | null
    is_read: boolean
    read_at: string | null
    created_at: string

    // Joined fields
    sender?: {
        id: number
        username: string
        full_name: string | null
    }
    thread?: ChatThread
}

export interface CreateChatThreadDto {
    purchase_id: number
}

export interface SendMessageDto {
    content: string
    message_type: 'text' | 'image' | 'file'
    attachment_url?: string
}
