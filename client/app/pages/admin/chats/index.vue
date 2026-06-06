<template>
    <div class="admin-page admin-chats">
        <!-- Page Header -->
        <div class="admin-page__header">
            <div>
                <h1 class="admin-page__title">Чаты</h1>
                <p class="admin-page__subtitle">Управление обращениями пользователей</p>
            </div>
            <div class="admin-chats__header-actions">
                <UInput
                    v-model="searchQuery"
                    placeholder="Поиск по пользователю..."
                    leading-icon="i-heroicons-magnifying-glass-20-solid"
                    size="sm"
                    class="admin-chats__search"
                />
                <USelect
                    v-model="statusFilter"
                    :items="statusOptions"
                    size="sm"
                    class="admin-chats__filter"
                />
            </div>
        </div>

        <!-- Loading State -->
        <template v-if="loading">
            <div class="admin-chats__layout">
                <div class="admin-chats__sidebar">
                    <div v-for="i in 5" :key="i" class="admin-chats__skeleton">
                        <div class="admin-chats__skeleton-avatar" />
                        <div class="admin-chats__skeleton-content">
                            <div class="admin-chats__skeleton-line" style="width: 60%" />
                            <div class="admin-chats__skeleton-line" style="width: 40%" />
                        </div>
                    </div>
                </div>
                <div class="admin-chats__main">
                    <div class="admin-chats__skeleton-messages">
                        <div
                            v-for="i in 4"
                            :key="i"
                            class="admin-chats__skeleton-msg"
                            :class="{ 'admin-chats__skeleton-msg--right': i % 2 === 0 }"
                        >
                            <div class="admin-chats__skeleton-line" style="width: 70%" />
                            <div class="admin-chats__skeleton-line" style="width: 30%" />
                        </div>
                    </div>
                </div>
            </div>
        </template>

        <!-- Error State -->
        <template v-else-if="error">
            <UCard class="admin-page__table-card">
                <div class="admin-page__placeholder">
                    <UIcon
                        name="i-heroicons-exclamation-triangle"
                        class="admin-page__placeholder-icon admin-page__placeholder-icon--error"
                    />
                    <p class="admin-page__hint">{{ error }}</p>
                    <UButton color="primary" variant="solid" @click="loadThreads">
                        Повторить
                    </UButton>
                </div>
            </UCard>
        </template>

        <!-- Empty State -->
        <template v-else-if="filteredThreads.length === 0">
            <UCard class="admin-page__table-card">
                <div class="admin-page__placeholder">
                    <UIcon
                        name="i-heroicons-chat-bubble-left-right"
                        class="admin-page__placeholder-icon"
                    />
                    <p class="admin-page__hint">Нет активных чатов</p>
                </div>
            </UCard>
        </template>

        <!-- Chat Layout -->
        <template v-else>
            <div class="admin-chats__layout">
                <!-- Thread List (Sidebar) -->
                <div class="admin-chats__sidebar">
                    <div
                        v-for="thread in filteredThreads"
                        :key="thread.id"
                        class="admin-chats__thread"
                        :class="{
                            'admin-chats__thread--active':
                                selectedThread?.id === thread.id,
                            'admin-chats__thread--unresolved': !thread.is_resolved,
                        }"
                        @click="selectThread(thread)"
                    >
                        <UAvatar
                            :text="
                                (
                                    (thread.user?.username ||
                                        thread.user?.full_name ||
                                        '?')[0] || '?'
                                ).toUpperCase()
                            "
                            size="sm"
                            class="admin-chats__thread-avatar"
                        />
                        <div class="admin-chats__thread-info">
                            <div class="admin-chats__thread-name">
                                {{
                                    thread.user?.full_name ||
                                    thread.user?.username ||
                                    'Пользователь'
                                }}
                            </div>
                            <div class="admin-chats__thread-product">
                                {{
                                    thread.purchase?.product?.title_ru ||
                                    thread.purchase?.product?.title_en ||
                                    '—'
                                }}
                            </div>
                            <div class="admin-chats__thread-time">
                                {{ formatDate(thread.last_message_at) }}
                            </div>
                        </div>
                        <div class="admin-chats__thread-status">
                            <UBadge
                                v-if="thread.is_resolved"
                                color="neutral"
                                variant="solid"
                                size="xs"
                            >
                                Закрыт
                            </UBadge>
                            <UBadge v-else color="success" variant="solid" size="xs">
                                Открыт
                            </UBadge>
                        </div>
                    </div>
                </div>

                <!-- Messages Area -->
                <div class="admin-chats__main">
                    <template v-if="selectedThread">
                        <!-- Thread Header -->
                        <div class="admin-chats__main-header">
                            <div class="admin-chats__main-header-info">
                                <UAvatar
                                    :text="
                                        (
                                            (selectedThread.user?.username ||
                                                selectedThread.user?.full_name ||
                                                '?')[0] || '?'
                                        ).toUpperCase()
                                    "
                                    size="sm"
                                />
                                <div>
                                    <div class="admin-chats__main-header-name">
                                        {{
                                            selectedThread.user?.full_name ||
                                            selectedThread.user?.username ||
                                            'Пользователь'
                                        }}
                                    </div>
                                    <div class="admin-chats__main-header-product">
                                        {{
                                            selectedThread.purchase?.product?.title_ru ||
                                            selectedThread.purchase?.product?.title_en ||
                                            '—'
                                        }}
                                    </div>
                                </div>
                            </div>
                            <div class="admin-chats__main-header-actions">
                                <UButton
                                    v-if="!selectedThread.is_resolved"
                                    color="warning"
                                    variant="soft"
                                    size="sm"
                                    :loading="resolving"
                                    @click="resolveThread(selectedThread.id)"
                                >
                                    <template #leading>
                                        <UIcon name="i-heroicons-check-circle" />
                                    </template>
                                    Закрыть чат
                                </UButton>
                                <UButton
                                    v-else
                                    color="primary"
                                    variant="soft"
                                    size="sm"
                                    :loading="reopening"
                                    @click="reopenThread(selectedThread.id)"
                                >
                                    <template #leading>
                                        <UIcon name="i-heroicons-arrow-uturn-left" />
                                    </template>
                                    Открыть снова
                                </UButton>
                            </div>
                        </div>

                        <!-- Messages List -->
                        <div ref="messagesContainer" class="admin-chats__messages">
                            <div
                                v-for="msg in messages"
                                :key="msg.id"
                                class="admin-chats__message"
                                :class="{
                                    'admin-chats__message--admin':
                                        msg.sender_id !== selectedThread.user_id,
                                    'admin-chats__message--user':
                                        msg.sender_id === selectedThread.user_id,
                                }"
                            >
                                <div class="admin-chats__message-bubble">
                                    <div class="admin-chats__message-text">
                                        {{ msg.content }}
                                    </div>
                                    <div class="admin-chats__message-meta">
                                        <span class="admin-chats__message-sender">
                                            {{
                                                msg.sender?.full_name ||
                                                msg.sender?.username ||
                                                '—'
                                            }}
                                        </span>
                                        <span class="admin-chats__message-time">{{
                                            formatMessageTime(msg.created_at)
                                        }}</span>
                                        <span
                                            v-if="
                                                msg.is_read &&
                                                msg.sender_id !== selectedThread.user_id
                                            "
                                            class="admin-chats__message-read"
                                        >
                                            <UIcon
                                                name="i-heroicons-check-badge"
                                                class="text-green-500"
                                            />
                                        </span>
                                    </div>
                                </div>
                            </div>
                            <div
                                v-if="messagesLoading"
                                class="admin-chats__messages-loading"
                            >
                                <UIcon
                                    name="i-heroicons-arrow-path"
                                    class="animate-spin"
                                />
                                Загрузка сообщений...
                            </div>
                        </div>

                        <!-- Message Input -->
                        <div class="admin-chats__input-area">
                            <UInput
                                v-model="newMessage"
                                placeholder="Введите сообщение..."
                                size="lg"
                                class="admin-chats__input"
                                @keydown.enter.prevent="sendMessage"
                            />
                            <UButton
                                color="primary"
                                variant="solid"
                                :disabled="!newMessage.trim()"
                                :loading="sending"
                                @click="sendMessage"
                            >
                                <template #leading>
                                    <UIcon name="i-heroicons-paper-airplane" />
                                </template>
                                Отправить
                            </UButton>
                        </div>
                    </template>

                    <!-- No Thread Selected -->
                    <template v-else>
                        <div class="admin-page__placeholder">
                            <UIcon
                                name="i-heroicons-chat-bubble-left-right"
                                class="admin-page__placeholder-icon"
                            />
                            <p class="admin-page__hint">
                                Выберите чат для просмотра сообщений
                            </p>
                        </div>
                    </template>
                </div>
            </div>
        </template>
    </div>
</template>

<script setup lang="ts">
import { ref, computed, watch, nextTick, onMounted } from 'vue';
import {
    fetchAdminChatThreads,
    fetchAdminChatMessages,
    sendAdminChatMessage,
    resolveAdminChatThread,
    reopenAdminChatThread,
} from '~/api/admin';
import type { AdminChatThread, AdminChatMessage } from '~/api/admin';

definePageMeta({
    layout: 'admin',
    middleware: 'admin-auth',
});

const loading = ref(true);
const error = ref<string | null>(null);
const threads = ref<AdminChatThread[]>([]);
const selectedThread = ref<AdminChatThread | null>(null);
const messages = ref<AdminChatMessage[]>([]);
const messagesLoading = ref(false);
const newMessage = ref('');
const sending = ref(false);
const resolving = ref(false);
const reopening = ref(false);
const searchQuery = ref('');
const statusFilter = ref('all');
const messagesContainer = ref<HTMLElement | null>(null);

const statusOptions = [
    { label: 'Все статусы', value: 'all' },
    { label: 'Открытые', value: 'open' },
    { label: 'Закрытые', value: 'resolved' },
];

const filteredThreads = computed(() => {
    let result = threads.value;

    // Filter by status
    if (statusFilter.value === 'open') {
        result = result.filter((t) => !t.is_resolved);
    } else if (statusFilter.value === 'resolved') {
        result = result.filter((t) => t.is_resolved);
    }

    // Filter by search
    if (searchQuery.value.trim()) {
        const q = searchQuery.value.toLowerCase();
        result = result.filter((t) => {
            const name = (t.user?.full_name || t.user?.username || '').toLowerCase();
            const email = (t.user?.email || '').toLowerCase();
            const product = (
                t.purchase?.product?.title_ru ||
                t.purchase?.product?.title_en ||
                ''
            ).toLowerCase();
            return name.includes(q) || email.includes(q) || product.includes(q);
        });
    }

    return result;
});

async function loadThreads() {
    loading.value = true;
    error.value = null;
    try {
        threads.value = await fetchAdminChatThreads();
    } catch (e) {
        error.value = e instanceof Error ? e.message : 'Ошибка загрузки чатов';
    } finally {
        loading.value = false;
    }
}

async function selectThread(thread: AdminChatThread) {
    selectedThread.value = thread;
    messagesLoading.value = true;
    messages.value = [];
    try {
        messages.value = await fetchAdminChatMessages(thread.id);
        await nextTick();
        scrollToBottom();
    } catch (e) {
        console.error('Failed to load messages:', e);
    } finally {
        messagesLoading.value = false;
    }
}

async function sendMessage() {
    if (!newMessage.value.trim() || !selectedThread.value || sending.value) return;

    sending.value = true;
    try {
        const msg = await sendAdminChatMessage(
            selectedThread.value.id,
            newMessage.value.trim()
        );
        messages.value.push(msg);
        newMessage.value = '';
        await nextTick();
        scrollToBottom();
    } catch (e) {
        console.error('Failed to send message:', e);
    } finally {
        sending.value = false;
    }
}

async function resolveThread(id: number) {
    resolving.value = true;
    try {
        await resolveAdminChatThread(id);
        const thread = threads.value.find((t) => t.id === id);
        if (thread) {
            thread.is_resolved = true;
        }
        if (selectedThread.value?.id === id) {
            selectedThread.value.is_resolved = true;
        }
    } catch (e) {
        console.error('Failed to resolve thread:', e);
    } finally {
        resolving.value = false;
    }
}

async function reopenThread(id: number) {
    reopening.value = true;
    try {
        await reopenAdminChatThread(id);
        const thread = threads.value.find((t) => t.id === id);
        if (thread) {
            thread.is_resolved = false;
        }
        if (selectedThread.value?.id === id) {
            selectedThread.value.is_resolved = false;
        }
    } catch (e) {
        console.error('Failed to reopen thread:', e);
    } finally {
        reopening.value = false;
    }
}

function scrollToBottom() {
    if (messagesContainer.value) {
        messagesContainer.value.scrollTop = messagesContainer.value.scrollHeight;
    }
}

function formatDate(dateStr: string): string {
    if (!dateStr) return '';
    const d = new Date(dateStr);
    const now = new Date();
    const diff = now.getTime() - d.getTime();
    const days = Math.floor(diff / (1000 * 60 * 60 * 24));

    if (days === 0) {
        return d.toLocaleTimeString('ru-RU', { hour: '2-digit', minute: '2-digit' });
    } else if (days === 1) {
        return 'Вчера';
    } else if (days < 7) {
        return `${days} дн. назад`;
    }
    return d.toLocaleDateString('ru-RU', { day: 'numeric', month: 'short' });
}

function formatMessageTime(dateStr: string): string {
    if (!dateStr) return '';
    const d = new Date(dateStr);
    return d.toLocaleTimeString('ru-RU', { hour: '2-digit', minute: '2-digit' });
}

onMounted(() => {
    loadThreads();
});
</script>

<style scoped>
@import '../_shared.css';

.admin-chats__header-actions {
    display: flex;
    align-items: center;
    gap: 12px;
}

.admin-chats__search {
    width: 240px;
}

.admin-chats__filter {
    width: 160px;
}

.admin-chats__layout {
    display: grid;
    grid-template-columns: 340px 1fr;
    gap: 0;
    background: var(--admin-surface, #ffffff);
    border-radius: 12px;
    border: 1px solid var(--admin-border, #e0e0e0);
    overflow: hidden;
    min-height: 600px;
    height: calc(100vh - 200px);
}

@media (max-width: 768px) {
    .admin-chats__layout {
        grid-template-columns: 1fr;
        height: auto;
        min-height: auto;
    }
}

/* Sidebar */
.admin-chats__sidebar {
    border-right: 1px solid var(--admin-border, #e0e0e0);
    overflow-y: auto;
    background: var(--admin-bg, #f8f9fa);
}

.admin-chats__thread {
    display: flex;
    align-items: flex-start;
    gap: 12px;
    padding: 14px 16px;
    cursor: pointer;
    transition: background 0.15s;
    border-bottom: 1px solid var(--admin-border, #e0e0e0);
}

.admin-chats__thread:hover {
    background: var(--admin-hover, #f0f0f0);
}

.admin-chats__thread--active {
    background: var(--admin-primary-light, #e8e5f9);
    border-left: 3px solid var(--admin-primary, #6c5ce7);
}

.admin-chats__thread--unresolved {
    font-weight: 600;
}

.admin-chats__thread-avatar {
    flex-shrink: 0;
    margin-top: 2px;
}

.admin-chats__thread-info {
    flex: 1;
    min-width: 0;
}

.admin-chats__thread-name {
    font-size: 14px;
    font-weight: 600;
    color: var(--admin-text-primary, #2d3436);
    white-space: nowrap;
    overflow: hidden;
    text-overflow: ellipsis;
}

.admin-chats__thread-product {
    font-size: 12px;
    color: var(--admin-text-secondary, #636e72);
    white-space: nowrap;
    overflow: hidden;
    text-overflow: ellipsis;
    margin-top: 2px;
}

.admin-chats__thread-time {
    font-size: 11px;
    color: var(--admin-text-tertiary, #b2bec3);
    margin-top: 4px;
}

.admin-chats__thread-status {
    flex-shrink: 0;
    margin-top: 2px;
}

/* Main Area */
.admin-chats__main {
    display: flex;
    flex-direction: column;
    height: 100%;
    overflow: hidden;
}

.admin-chats__main-header {
    display: flex;
    align-items: center;
    justify-content: space-between;
    padding: 14px 20px;
    border-bottom: 1px solid var(--admin-border, #e0e0e0);
    background: var(--admin-surface, #ffffff);
    flex-shrink: 0;
}

.admin-chats__main-header-info {
    display: flex;
    align-items: center;
    gap: 10px;
}

.admin-chats__main-header-name {
    font-size: 14px;
    font-weight: 600;
    color: var(--admin-text-primary, #2d3436);
}

.admin-chats__main-header-product {
    font-size: 12px;
    color: var(--admin-text-secondary, #636e72);
}

.admin-chats__main-header-actions {
    display: flex;
    gap: 8px;
}

/* Messages */
.admin-chats__messages {
    flex: 1;
    overflow-y: auto;
    padding: 20px;
    display: flex;
    flex-direction: column;
    gap: 12px;
    background: var(--admin-bg, #f8f9fa);
}

.admin-chats__messages-loading {
    display: flex;
    align-items: center;
    justify-content: center;
    gap: 8px;
    padding: 20px;
    color: var(--admin-text-secondary, #636e72);
    font-size: 13px;
}

.admin-chats__message {
    display: flex;
    max-width: 75%;
}

.admin-chats__message--user {
    align-self: flex-start;
}

.admin-chats__message--admin {
    align-self: flex-end;
}

.admin-chats__message-bubble {
    padding: 10px 14px;
    border-radius: 12px;
    font-size: 14px;
    line-height: 1.5;
    position: relative;
}

.admin-chats__message--user .admin-chats__message-bubble {
    background: var(--admin-surface, #ffffff);
    border: 1px solid var(--admin-border, #e0e0e0);
    color: var(--admin-text-primary, #2d3436);
    border-bottom-left-radius: 4px;
}

.admin-chats__message--admin .admin-chats__message-bubble {
    background: var(--admin-primary, #6c5ce7);
    color: #ffffff;
    border-bottom-right-radius: 4px;
}

.admin-chats__message-text {
    word-break: break-word;
    white-space: pre-wrap;
}

.admin-chats__message-meta {
    display: flex;
    align-items: center;
    gap: 6px;
    margin-top: 4px;
    font-size: 11px;
}

.admin-chats__message--user .admin-chats__message-meta {
    color: var(--admin-text-tertiary, #b2bec3);
}

.admin-chats__message--admin .admin-chats__message-meta {
    color: rgba(255, 255, 255, 0.7);
}

.admin-chats__message-sender {
    font-weight: 500;
}

.admin-chats__message-read {
    display: inline-flex;
    align-items: center;
}

/* Input Area */
.admin-chats__input-area {
    display: flex;
    align-items: center;
    gap: 10px;
    padding: 14px 20px;
    border-top: 1px solid var(--admin-border, #e0e0e0);
    background: var(--admin-surface, #ffffff);
    flex-shrink: 0;
}

.admin-chats__input {
    flex: 1;
}

/* Skeleton */
.admin-chats__skeleton {
    display: flex;
    align-items: center;
    gap: 12px;
    padding: 14px 16px;
    border-bottom: 1px solid var(--admin-border, #e0e0e0);
}

.admin-chats__skeleton-avatar {
    width: 32px;
    height: 32px;
    border-radius: 50%;
    background: var(--admin-skeleton, #e0e0e0);
    flex-shrink: 0;
}

.admin-chats__skeleton-content {
    flex: 1;
    display: flex;
    flex-direction: column;
    gap: 6px;
}

.admin-chats__skeleton-line {
    height: 12px;
    border-radius: 4px;
    background: var(--admin-skeleton, #e0e0e0);
    animation: shimmer 1.5s infinite;
}

.admin-chats__skeleton-messages {
    padding: 20px;
    display: flex;
    flex-direction: column;
    gap: 16px;
}

.admin-chats__skeleton-msg {
    display: flex;
    flex-direction: column;
    gap: 6px;
    padding: 12px 16px;
    border-radius: 12px;
    background: var(--admin-skeleton, #e0e0e0);
    max-width: 60%;
}

.admin-chats__skeleton-msg--right {
    align-self: flex-end;
}

.admin-chats__skeleton-msg .admin-chats__skeleton-line {
    background: rgba(255, 255, 255, 0.3);
}

@keyframes shimmer {
    0% {
        opacity: 0.6;
    }
    50% {
        opacity: 1;
    }
    100% {
        opacity: 0.6;
    }
}

/* Dark mode overrides */
:root.dark .admin-chats__sidebar {
    background: var(--admin-bg, #1a1a2e);
}

:root.dark .admin-chats__thread--active {
    background: rgba(108, 92, 231, 0.15);
}

:root.dark .admin-chats__message--user .admin-chats__message-bubble {
    background: var(--admin-surface, #2d2d44);
    border-color: var(--admin-border, #3d3d5c);
}

:root.dark .admin-chats__messages {
    background: var(--admin-bg, #1a1a2e);
}
</style>
