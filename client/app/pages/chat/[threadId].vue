<template>
    <div class="chat-thread-page">
        <div class="chat-thread-page__header">
            <button class="chat-thread-page__back" @click="goBack">← Назад</button>
            <h1 class="chat-thread-page__title">Чат</h1>
            <div v-if="thread" class="chat-thread-page__thread-info">
                <div class="chat-thread-page__product">
                    {{ thread.purchase?.product?.title_ru || 'Без названия' }}
                </div>
                <div class="chat-thread-page__status">
                    <span
                        v-if="thread.is_resolved"
                        class="chat-thread-page__status-resolved"
                        >Закрыт</span
                    >
                    <span v-else class="chat-thread-page__status-active">Активен</span>
                </div>
            </div>
        </div>

        <div v-if="loading" class="chat-thread-page__loading">Загрузка чата...</div>
        <div v-else-if="error" class="chat-thread-page__error">{{ error }}</div>
        <div v-else-if="!thread" class="chat-thread-page__not-found">Чат не найден</div>
        <div v-else class="chat-thread-page__content">
            <div ref="messagesContainer" class="chat-thread-page__messages">
                <div
                    v-for="message in messages"
                    :key="message.id"
                    class="chat-thread-page__message"
                    :class="{
                        'chat-thread-page__message--outgoing':
                            message.sender_id === currentUserId,
                        'chat-thread-page__message--incoming':
                            message.sender_id !== currentUserId,
                    }"
                >
                    <div class="chat-thread-page__message-avatar">
                        {{
                            message.sender?.full_name?.[0] ||
                            message.sender?.username?.[0] ||
                            '?'
                        }}
                    </div>
                    <div class="chat-thread-page__message-body">
                        <div class="chat-thread-page__message-header">
                            <span class="chat-thread-page__message-sender">
                                {{
                                    message.sender?.full_name ||
                                    message.sender?.username ||
                                    'Пользователь'
                                }}
                            </span>
                            <span class="chat-thread-page__message-time">{{
                                formatTime(message.created_at)
                            }}</span>
                        </div>
                        <div class="chat-thread-page__message-content">
                            <div v-if="message.message_type === 'text'">
                                {{ message.content }}
                            </div>
                            <div
                                v-else-if="
                                    message.message_type === 'image' &&
                                    message.attachment_url
                                "
                                class="chat-thread-page__image"
                            >
                                <img :src="message.attachment_url" alt="Изображение" />
                            </div>
                            <div
                                v-else-if="
                                    message.message_type === 'file' &&
                                    message.attachment_url
                                "
                                class="chat-thread-page__file"
                            >
                                <a :href="message.attachment_url" target="_blank">Файл</a>
                            </div>
                            <div v-else class="chat-thread-page__unknown">
                                [Неподдерживаемый тип сообщения]
                            </div>
                        </div>
                    </div>
                </div>
            </div>

            <div v-if="!thread.is_resolved" class="chat-thread-page__input">
                <textarea
                    v-model="newMessage"
                    class="chat-thread-page__textarea"
                    placeholder="Введите сообщение..."
                    rows="3"
                    @keydown.enter.exact.prevent="sendMessage"
                />
                <div class="chat-thread-page__input-actions">
                    <button
                        class="chat-thread-page__send-button"
                        @click="sendMessage"
                        :disabled="!newMessage.trim()"
                    >
                        Отправить
                    </button>
                </div>
            </div>
            <div v-else class="chat-thread-page__resolved">
                Этот чат закрыт. Новые сообщения отправлять нельзя.
            </div>
        </div>
    </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, onUnmounted, nextTick } from 'vue';
import { useRoute, useRouter } from 'vue-router';
import { useAuthStore } from '~/stores/AuthStore';
import type { ChatThread, ChatMessage } from '~/types';

const route = useRoute();
const router = useRouter();
const authStore = useAuthStore();

const threadId = computed(() => Number(route.params.threadId));
const currentUserId = computed(() => authStore.userId || 0);

const thread = ref<ChatThread | null>(null);
const messages = ref<ChatMessage[]>([]);
const loading = ref(true);
const error = ref('');
const newMessage = ref('');
const messagesContainer = ref<HTMLElement | null>(null);
let pollInterval: number | null = null;

const authToken = computed(() => authStore.token);

onMounted(async () => {
    await loadThread();
    if (thread.value) {
        await loadMessages();
        startPolling();
    }
});

onUnmounted(() => {
    stopPolling();
});

async function loadThread() {
    loading.value = true;
    error.value = '';
    try {
        const response = await $fetch<ChatThread>(`/api/chat/threads/${threadId.value}`, {
            headers: { Authorization: `Bearer ${authToken.value}` },
        });
        thread.value = response;
    } catch (err) {
        error.value = 'Не удалось загрузить чат';
        console.error(err);
    } finally {
        loading.value = false;
    }
}

async function loadMessages() {
    if (!thread.value) return;
    try {
        const response = await $fetch<ChatMessage[]>(
            `/api/chat/threads/${thread.value.id}/messages`,
            {
                headers: { Authorization: `Bearer ${authToken.value}` },
            }
        );
        messages.value = response;
        scrollToBottom();
    } catch (err) {
        error.value = 'Не удалось загрузить сообщения';
    }
}

async function sendMessage() {
    const content = newMessage.value.trim();
    if (!content || !thread.value || thread.value.is_resolved) return;

    try {
        await $fetch(`/api/chat/threads/${thread.value.id}/messages`, {
            method: 'POST',
            headers: { Authorization: `Bearer ${authToken.value}` },
            body: {
                content,
                message_type: 'text',
            },
        });
        newMessage.value = '';
        await loadMessages();
    } catch (err) {
        error.value = 'Не удалось отправить сообщение';
    }
}

async function pollNewMessages() {
    if (!thread.value) return;
    try {
        const lastId =
            messages.value.length > 0 ? Math.max(...messages.value.map((m) => m.id)) : 0;
        const response = await $fetch<{ messages: ChatMessage[] }>(
            `/api/chat/poll?thread_id=${thread.value.id}&last_message_id=${lastId}`,
            {
                headers: { Authorization: `Bearer ${authToken.value}` },
            }
        );
        if (response.messages && response.messages.length > 0) {
            messages.value.push(...response.messages);
            scrollToBottom();
        }
    } catch (err) {
        // Игнорируем ошибки long-polling
    }
}

function startPolling() {
    pollInterval = window.setInterval(pollNewMessages, 5000);
}

function stopPolling() {
    if (pollInterval) {
        clearInterval(pollInterval);
        pollInterval = null;
    }
}

function scrollToBottom() {
    nextTick(() => {
        if (messagesContainer.value) {
            messagesContainer.value.scrollTop = messagesContainer.value.scrollHeight;
        }
    });
}

function formatTime(dateString: string) {
    const date = new Date(dateString);
    return date.toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' });
}

function goBack() {
    router.push('/chat');
}
</script>

<style scoped>
.chat-thread-page {
    max-width: 800px;
    margin: 0 auto;
    padding: 20px;
    height: 100vh;
    display: flex;
    flex-direction: column;
}

.chat-thread-page__header {
    display: flex;
    align-items: center;
    gap: 16px;
    margin-bottom: 24px;
    padding-bottom: 16px;
    border-bottom: 1px solid #eee;
}

.chat-thread-page__back {
    background: none;
    border: none;
    color: #6a11cb;
    font-size: 16px;
    cursor: pointer;
    padding: 8px;
}

.chat-thread-page__title {
    font-size: 24px;
    font-weight: 700;
    color: #222;
    flex: 1;
}

.chat-thread-page__thread-info {
    display: flex;
    flex-direction: column;
    align-items: flex-end;
}

.chat-thread-page__product {
    font-weight: 600;
    color: #333;
}

.chat-thread-page__status {
    font-size: 12px;
    padding: 2px 8px;
    border-radius: 12px;
}

.chat-thread-page__status-active {
    background: #e8f5e9;
    color: #388e3c;
}

.chat-thread-page__status-resolved {
    background: #f5f5f5;
    color: #757575;
}

.chat-thread-page__loading,
.chat-thread-page__error,
.chat-thread-page__not-found {
    text-align: center;
    padding: 40px;
    color: #666;
}

.chat-thread-page__error {
    color: #d32f2f;
}

.chat-thread-page__content {
    flex: 1;
    display: flex;
    flex-direction: column;
    overflow: hidden;
}

.chat-thread-page__messages {
    flex: 1;
    overflow-y: auto;
    padding: 16px;
    background: #fafafa;
    border-radius: 12px;
    margin-bottom: 16px;
    display: flex;
    flex-direction: column;
    gap: 16px;
}

.chat-thread-page__message {
    display: flex;
    gap: 12px;
    max-width: 80%;
}

.chat-thread-page__message--outgoing {
    align-self: flex-end;
    flex-direction: row-reverse;
}

.chat-thread-page__message--incoming {
    align-self: flex-start;
}

.chat-thread-page__message-avatar {
    width: 40px;
    height: 40px;
    border-radius: 50%;
    background: linear-gradient(135deg, #6a11cb 0%, #2575fc 100%);
    color: white;
    display: flex;
    align-items: center;
    justify-content: center;
    font-weight: bold;
    flex-shrink: 0;
}

.chat-thread-page__message-body {
    background: white;
    padding: 12px 16px;
    border-radius: 12px;
    box-shadow: 0 2px 4px rgba(0, 0, 0, 0.05);
    max-width: 100%;
}

.chat-thread-page__message--outgoing .chat-thread-page__message-body {
    background: #e3f2fd;
    border-bottom-right-radius: 4px;
}

.chat-thread-page__message--incoming .chat-thread-page__message-body {
    border-bottom-left-radius: 4px;
}

.chat-thread-page__message-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 6px;
}

.chat-thread-page__message-sender {
    font-weight: 600;
    font-size: 14px;
    color: #555;
}

.chat-thread-page__message-time {
    font-size: 12px;
    color: #999;
}

.chat-thread-page__message-content {
    font-size: 15px;
    line-height: 1.4;
    word-wrap: break-word;
}

.chat-thread-page__image img {
    max-width: 100%;
    border-radius: 8px;
}

.chat-thread-page__file a {
    color: #6a11cb;
    text-decoration: underline;
}

.chat-thread-page__input {
    border-top: 1px solid #eee;
    padding-top: 16px;
}

.chat-thread-page__textarea {
    width: 100%;
    padding: 12px 16px;
    border: 1px solid #ddd;
    border-radius: 8px;
    font-family: inherit;
    font-size: 15px;
    resize: none;
    margin-bottom: 12px;
}

.chat-thread-page__textarea:focus {
    outline: none;
    border-color: #6a11cb;
    box-shadow: 0 0 0 3px rgba(106, 17, 203, 0.1);
}

.chat-thread-page__input-actions {
    display: flex;
    justify-content: flex-end;
}

.chat-thread-page__send-button {
    padding: 10px 24px;
    background: #6a11cb;
    color: white;
    border: none;
    border-radius: 8px;
    cursor: pointer;
    font-weight: 500;
}

.chat-thread-page__send-button:disabled {
    opacity: 0.5;
    cursor: not-allowed;
}

.chat-thread-page__resolved {
    text-align: center;
    padding: 20px;
    background: #f5f5f5;
    border-radius: 8px;
    color: #666;
    margin-top: 16px;
}
</style>
