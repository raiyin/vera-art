<template>
    <div class="chat-page">
        <div class="chat-page__header">
            <h1 class="chat-page__title">Мои чаты</h1>
            <p class="chat-page__subtitle">Общайтесь с преподавателями по вашим курсам</p>
        </div>

        <div v-if="loading" class="chat-page__loading">Загрузка...</div>
        <div v-else-if="error" class="chat-page__error">{{ error }}</div>
        <div v-else-if="threads.length === 0" class="chat-page__empty">
            <p>У вас пока нет активных чатов.</p>
            <p>После покупки курса вы сможете создать чат с преподавателем.</p>
        </div>
        <div v-else class="chat-page__threads">
            <div
                v-for="thread in threads"
                :key="thread.id"
                class="chat-page__thread"
                :class="{ 'chat-page__thread--unread': hasUnread(thread) }"
                @click="openThread(thread.id)"
            >
                <div class="chat-page__thread-avatar">
                    {{ thread.purchase?.product?.title_ru?.[0] || 'Ч' }}
                </div>
                <div class="chat-page__thread-info">
                    <div class="chat-page__thread-title">
                        {{ thread.purchase?.product?.title_ru || 'Без названия' }}
                    </div>
                    <div class="chat-page__thread-last-message">
                        {{ lastMessagePreview(thread) }}
                    </div>
                    <div class="chat-page__thread-meta">
                        <span class="chat-page__thread-date">{{
                            formatDate(thread.last_message_at)
                        }}</span>
                        <span v-if="hasUnread(thread)" class="chat-page__thread-unread"
                            >Новое</span
                        >
                    </div>
                </div>
                <div class="chat-page__thread-arrow">→</div>
            </div>
        </div>
    </div>
</template>

<script setup lang="ts">
import { ref, onMounted, computed } from 'vue';
import { useRouter } from 'vue-router';
import { useAuthStore } from '../../stores/AuthStore';
import type { ChatThread } from '../../types';

const router = useRouter();
const authStore = useAuthStore();

const threads = ref<ChatThread[]>([]);
const loading = ref(true);
const error = ref('');

const authToken = computed(() => authStore.token);

onMounted(async () => {
    await loadThreads();
});

async function loadThreads() {
    loading.value = true;
    error.value = '';
    try {
        const response = await $fetch<ChatThread[]>('/api/chat/threads', {
            headers: { Authorization: `Bearer ${authToken.value}` },
        });
        threads.value = response;
    } catch (err) {
        error.value = 'Не удалось загрузить чаты';
        console.error(err);
    } finally {
        loading.value = false;
    }
}

function hasUnread(thread: ChatThread): boolean {
    // Простая проверка: если есть сообщения, не прочитанные текущим пользователем
    // В реальном приложении нужно учитывать поле is_read
    return false;
}

function lastMessagePreview(thread: ChatThread): string {
    if (!thread.messages || thread.messages.length === 0) {
        return 'Нет сообщений';
    }
    const last = thread.messages[thread.messages.length - 1];
    const content = last?.content || '';
    return content.length > 50 ? content.substring(0, 50) + '...' : content;
}

function formatDate(dateString: string): string {
    const date = new Date(dateString);
    const now = new Date();
    const diffMs = now.getTime() - date.getTime();
    const diffDays = Math.floor(diffMs / (1000 * 60 * 60 * 24));
    if (diffDays === 0) {
        return 'Сегодня';
    } else if (diffDays === 1) {
        return 'Вчера';
    } else if (diffDays < 7) {
        return `${diffDays} дн. назад`;
    } else {
        return date.toLocaleDateString();
    }
}

function openThread(threadId: number) {
    router.push(`/chat/${threadId}`);
}
</script>

<style scoped>
.chat-page {
    max-width: 800px;
    margin: 0 auto;
    padding: 40px 20px;
}

.chat-page__header {
    text-align: center;
    margin-bottom: 40px;
}

.chat-page__title {
    font-size: 36px;
    font-weight: 700;
    color: #222;
    margin-bottom: 8px;
}

.chat-page__subtitle {
    font-size: 16px;
    color: #666;
    line-height: 1.5;
}

.chat-page__loading,
.chat-page__error,
.chat-page__empty {
    text-align: center;
    padding: 40px;
    background: #f9f9f9;
    border-radius: 12px;
    color: #666;
}

.chat-page__error {
    color: #d32f2f;
    background: #ffebee;
}

.chat-page__threads {
    display: flex;
    flex-direction: column;
    gap: 16px;
}

.chat-page__thread {
    display: flex;
    align-items: center;
    gap: 16px;
    padding: 20px;
    background: white;
    border-radius: 12px;
    box-shadow: 0 2px 8px rgba(0, 0, 0, 0.05);
    border: 1px solid #eaeaea;
    cursor: pointer;
    transition: all 0.2s;
}

.chat-page__thread:hover {
    box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
    border-color: #6a11cb;
}

.chat-page__thread--unread {
    border-left: 4px solid #6a11cb;
    background-color: #f9f5ff;
}

.chat-page__thread-avatar {
    width: 56px;
    height: 56px;
    border-radius: 50%;
    background: linear-gradient(135deg, #6a11cb 0%, #2575fc 100%);
    color: white;
    display: flex;
    align-items: center;
    justify-content: center;
    font-size: 24px;
    font-weight: bold;
    flex-shrink: 0;
}

.chat-page__thread-info {
    flex: 1;
    min-width: 0;
}

.chat-page__thread-title {
    font-size: 18px;
    font-weight: 600;
    color: #222;
    margin-bottom: 4px;
    white-space: nowrap;
    overflow: hidden;
    text-overflow: ellipsis;
}

.chat-page__thread-last-message {
    font-size: 14px;
    color: #666;
    margin-bottom: 6px;
    white-space: nowrap;
    overflow: hidden;
    text-overflow: ellipsis;
}

.chat-page__thread-meta {
    display: flex;
    gap: 12px;
    font-size: 12px;
    color: #999;
}

.chat-page__thread-unread {
    color: #6a11cb;
    font-weight: 600;
}

.chat-page__thread-arrow {
    color: #999;
    font-size: 20px;
    flex-shrink: 0;
}
</style>
