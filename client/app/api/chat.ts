import { getHttpClient, } from '~/api/http-client';

// ─── Types matching server DTOs ───────────────────────────────────────

export interface ChatThreadItem {
    id: number
    purchase_id: number
    user_id: number
    admin_id: number | null
    last_message_at: string
    is_resolved: boolean
    created_at: string
    user?: {
        username: string
        full_name: string
        email: string
    } | null
    purchase?: {
        product_id: number
        product_name: string
    } | null
}

export interface ChatMessageItem {
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
    sender?: {
        username: string
        full_name: string
        email: string
    } | null
}

export interface CreateThreadPayload {
    purchase_id: number
}

export interface SendMessagePayload {
    content: string
    message_type: 'text' | 'image' | 'file'
    attachment_url?: string
}

export interface PollMessagesResult {
    messages: ChatMessageItem[]
    has_more: boolean
}

// ─── API functions ────────────────────────────────────────────────────

/**
 * Fetches all chat threads for the current user.
 */
export async function fetchChatThreads(): Promise<ChatThreadItem[]> {
    try {
        const { data, } = await getHttpClient().get<ChatThreadItem[]>('chat/threads',);
        return data;
    } catch (e) {
        console.error('fetchChatThreads error:', e,);
        return [];
    }
}

/**
 * Creates a new chat thread.
 */
export async function createChatThread(payload: CreateThreadPayload,): Promise<{ id: number } | null> {
    try {
        const { data, } = await getHttpClient().post<{ id: number }>('chat/threads', payload,);
        return data;
    } catch (e) {
        console.error('createChatThread error:', e,);
        return null;
    }
}

/**
 * Fetches messages for a specific thread.
 */
export async function fetchChatMessages(threadId: number,): Promise<ChatMessageItem[]> {
    try {
        const { data, } = await getHttpClient().get<ChatMessageItem[]>(
            `chat/threads/${threadId}/messages`,
        );
        return data;
    } catch (e) {
        console.error('fetchChatMessages error:', e,);
        return [];
    }
}

/**
 * Sends a message in a thread.
 */
export async function sendChatMessage(threadId: number, payload: SendMessagePayload,): Promise<{ id: number } | null> {
    try {
        const { data, } = await getHttpClient().post<{ id: number }>(
            `chat/threads/${threadId}/messages`,
            payload,
        );
        return data;
    } catch (e) {
        console.error('sendChatMessage error:', e,);
        return null;
    }
}

/**
 * Marks a message as read.
 */
export async function markMessageAsRead(messageId: number,): Promise<boolean> {
    try {
        await getHttpClient().post(`chat/messages/${messageId}/read`,);
        return true;
    } catch (e) {
        console.error('markMessageAsRead error:', e,);
        return false;
    }
}

/**
 * Polls for new messages (long-polling).
 */
export async function pollMessages(): Promise<PollMessagesResult> {
    try {
        const { data, } = await getHttpClient().get<PollMessagesResult>('chat/poll',);
        return data;
    } catch (e) {
        console.error('pollMessages error:', e,);
        return { messages: [], has_more: false, };
    }
}
