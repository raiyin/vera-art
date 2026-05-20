<template>
    <div class="chat-widget" :class="{ 'chat-widget--collapsed': collapsed }">
        <div class="chat-widget__header" @click="toggleCollapse">
            <div class="chat-widget__header-left">
                <div class="chat-widget__icon">💬</div>
                <div class="chat-widget__title">
                    {{
                        thread
                            ? thread.purchase?.product?.title_ru || 'Чат с преподавателем'
                            : 'Чат'
                    }}
                </div>
                <div v-if="unreadCount > 0" class="chat-widget__badge">
                    {{ unreadCount }}
                </div>
            </div>
            <div class="chat-widget__header-right">
                <button class="chat-widget__header-button" @click.stop="toggleCollapse">
                    {{ collapsed ? '▶' : '▼' }}
                </button>
            </div>
        </div>

        <div v-if="!collapsed" class="chat-widget__body">
            <div v-if="loading" class="chat-widget__loading">Загрузка...</div>
            <div v-else-if="error" class="chat-widget__error">{{ error }}</div>
            <div v-else-if="!thread" class="chat-widget__no-thread">
                <p>У вас нет активного чата по этому курсу.</p>
                <button
                    v-if="purchaseId"
                    class="chat-widget__create-button"
                    @click="createThread"
                >
                    Создать чат
                </button>
            </div>
            <div v-else class="chat-widget__messages-container">
                <div ref="messagesContainer" class="chat-widget__messages">
                    <div
                        v-for="message in messages"
                        :key="message.id"
                        class="chat-widget__message"
                        :class="{
                            'chat-widget__message--outgoing':
                                message.sender_id === currentUserId,
                            'chat-widget__message--incoming':
                                message.sender_id !== currentUserId,
                        }"
                    >
                        <div class="chat-widget__message-sender">
                            {{
                                message.sender?.full_name ||
                                message.sender?.username ||
                                'Пользователь'
                            }}
                        </div>
                        <div class="chat-widget__message-content">
                            <div v-if="message.message_type === 'text'">
                                {{ message.content }}
                            </div>
                            <div
                                v-else-if="message.message_type === 'image'"
                                class="chat-widget__image"
                            >
                                <img :src="message.attachment_url" alt="Изображение" />
                            </div>
                            <div v-else class="chat-widget__file">
                                <a :href="message.attachment_url" target="_blank">Файл</a>
                            </div>
                        </div>
                        <div class="chat-widget__message-time">
                            {{ formatTime(message.created_at) }}
                        </div>
                    </div>
                </div>

                <div class="chat-widget__input-area">
                    <textarea
                        v-model="newMessage"
                        class="chat-widget__textarea"
                        placeholder="Введите сообщение..."
                        rows="2"
                        @keydown.enter.exact.prevent="sendMessage"
                    />
                    <button
                        class="chat-widget__send-button"
                        @click="sendMessage"
                        :disabled="!newMessage.trim()"
                    >
                        Отправить
                    </button>
                </div>
            </div>
        </div>
    </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, onUnmounted, watch, nextTick } from 'vue';
import { useAuthStore } from '../stores/AuthStore';
import type { ChatThread, ChatMessage } from '../types';

interface Props {
    purchaseId?: number;
    threadId?: number;
    autoLoad?: boolean;
}

const props = withDefaults(defineProps<Props>(), {
    autoLoad: true,
});

const authStore = useAuthStore();
const currentUserId = computed(() => authStore.currentUserId || 0);

const collapsed = ref(false);
const loading = ref(false);
const error = ref('');
const thread = ref<ChatThread | null>(null);
const messages = ref<ChatMessage[]>([]);
const newMessage = ref('');
const unreadCount = ref(0);
const messagesContainer = ref<HTMLElement | null>(null);
let pollInterval: unknown = null;

const fetchThread = async () => {
    if (props.threadId) {
        // Загружаем конкретный тред
        try {
            const response = await $fetch(`/api/chat/threads/${props.threadId}`, {
                headers: { Authorization: `Bearer ${authStore.token}` },
            });
            thread.value = response;
            await fetchMessages();
        } catch (err) {
            error.value = 'Не удалось загрузить чат';
        }
    } else if (props.purchaseId) {
        // Ищем тред по purchaseId
        try {
            const threads = await $fetch('/api/chat/threads', {
                headers: { Authorization: `Bearer ${authStore.token}` },
            });
            const found = threads.find(
                (t: ChatThread) => t.purchase_id === props.purchaseId
            );
            if (found) {
                thread.value = found;
                await fetchMessages();
            } else {
                thread.value = null;
            }
        } catch (err) {
            error.value = 'Не удалось загрузить список чатов';
        }
    }
};

const fetchMessages = async () => {
    if (!thread.value) return;
    try {
        const response = await $fetch(`/api/chat/threads/${thread.value.id}/messages`, {
            headers: { Authorization: `Bearer ${authStore.token}` },
        });
        messages.value = response;
        scrollToBottom();
        updateUnreadCount();
    } catch (err) {
        error.value = 'Не удалось загрузить сообщения';
    }
};

const createThread = async () => {
    if (!props.purchaseId) return;
    try {
        const response = await $fetch('/api/chat/threads', {
            method: 'POST',
            headers: { Authorization: `Bearer ${authStore.token}` },
            body: { purchase_id: props.purchaseId },
        });
        thread.value = response;
        await fetchMessages();
    } catch (err) {
        error.value = 'Не удалось создать чат';
    }
};

const sendMessage = async () => {
    const content = newMessage.value.trim();
    if (!content || !thread.value) return;

    try {
        await $fetch(`/api/chat/threads/${thread.value.id}/messages`, {
            method: 'POST',
            headers: { Authorization: `Bearer ${authStore.token}` },
            body: {
                content,
                message_type: 'text',
            },
        });
        newMessage.value = '';
        await fetchMessages();
    } catch (err) {
        error.value = 'Не удалось отправить сообщение';
    }
};

const pollNewMessages = async () => {
    if (!thread.value) return;
    try {
        const response = await $fetch(
            `/api/chat/poll?thread_id=${thread.value.id}&last_message_id=${lastMessageId}`,
            {
                headers: { Authorization: `Bearer ${authStore.token}` },
            }
        );
        if (response.messages && response.messages.length > 0) {
            messages.value.push(...response.messages);
            scrollToBottom();
            updateUnreadCount();
        }
    } catch (err) {
        // Игнорируем ошибки long-polling
    }
};

const lastMessageId = computed(() => {
    if (messages.value.length === 0) return 0;
    return Math.max(...messages.value.map((m) => m.id));
});

const updateUnreadCount = () => {
    unreadCount.value = messages.value.filter(
        (m) => !m.is_read && m.sender_id !== currentUserId.value
    ).length;
};

const formatTime = (dateString: string) => {
    const date = new Date(dateString);
    return date.toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' });
};

const scrollToBottom = () => {
    nextTick(() => {
        if (messagesContainer.value) {
            messagesContainer.value.scrollTop = messagesContainer.value.scrollHeight;
        }
    });
};

const toggleCollapse = () => {
    collapsed.value = !collapsed.value;
};

onMounted(() => {
    if (props.autoLoad) {
        fetchThread();
    }
    // Запускаем опрос каждые 10 секунд
    pollInterval = setInterval(pollNewMessages, 10000);
});

onUnmounted(() => {
    if (pollInterval) {
        clearInterval(pollInterval);
    }
});

watch(() => props.threadId, fetchThread);
watch(() => props.purchaseId, fetchThread);
</script>

<style scoped>
.chat-widget {
    border: 1px solid #ddd;
    border-radius: 12px;
    overflow: hidden;
    background: white;
    box-shadow: 0 4px 12px rgba(0, 0, 0, 0.05);
    max-width: 400px;
    min-width: 300px;
}

.chat-widget--collapsed {
    max-height: 60px;
}

.chat-widget__header {
    background: linear-gradient(135deg, #6a11cb 0%, #2575fc 100%);
    color: white;
    padding: 16px;
    display: flex;
    justify-content: space-between;
    align-items: center;
    cursor: pointer;
    user-select: none;
}

.chat-widget__header-left {
    display: flex;
    align-items: center;
    gap: 12px;
}

.chat-widget__icon {
    font-size: 20px;
}

.chat-widget__title {
    font-weight: 600;
    font-size: 16px;
}

.chat-widget__badge {
    background: #ff4757;
    color: white;
    font-size: 12px;
    font-weight: bold;
    border-radius: 50%;
    width: 20px;
    height: 20px;
    display: flex;
    align-items: center;
    justify-content: center;
}

.chat-widget__header-right {
    display: flex;
    align-items: center;
}

.chat-widget__header-button {
    background: none;
    border: none;
    color: white;
    font-size: 18px;
    cursor: pointer;
    padding: 4px;
}

.chat-widget__body {
    padding: 16px;
    max-height: 400px;
    display: flex;
    flex-direction: column;
}

.chat-widget__loading,
.chat-widget__error,
.chat-widget__no-thread {
    text-align: center;
    padding: 20px;
    color: #666;
}

.chat-widget__error {
    color: #d32f2f;
}

.chat-widget__create-button {
    margin-top: 12px;
    padding: 8px 16px;
    background: #6a11cb;
    color: white;
    border: none;
    border-radius: 6px;
    cursor: pointer;
    font-weight: 500;
}

.chat-widget__messages-container {
    flex: 1;
    display: flex;
    flex-direction: column;
}

.chat-widget__messages {
    flex: 1;
    overflow-y: auto;
    padding: 8px;
    margin-bottom: 16px;
    border: 1px solid #eee;
    border-radius: 8px;
    background: #fafafa;
    max-height: 300px;
}

.chat-widget__message {
    margin-bottom: 12px;
    max-width: 80%;
    padding: 8px 12px;
    border-radius: 12px;
    position: relative;
}

.chat-widget__message--outgoing {
    align-self: flex-end;
    background: #e3f2fd;
    border-bottom-right-radius: 4px;
}

.chat-widget__message--incoming {
    align-self: flex-start;
    background: white;
    border: 1px solid #e0e0e0;
    border-bottom-left-radius: 4px;
}

.chat-widget__message-sender {
    font-size: 12px;
    font-weight: 600;
    color: #555;
    margin-bottom: 4px;
}

.chat-widget__message-content {
    font-size: 14px;
    line-height: 1.4;
    word-wrap: break-word;
}

.chat-widget__image img {
    max-width: 100%;
    border-radius: 8px;
}

.chat-widget__file a {
    color: #6a11cb;
    text-decoration: underline;
}

.chat-widget__message-time {
    font-size: 10px;
    color: #999;
    text-align: right;
    margin-top: 4px;
}

.chat-widget__input-area {
    display: flex;
    gap: 8px;
    align-items: flex-end;
}

.chat-widget__textarea {
    flex: 1;
    padding: 10px;
    border: 1px solid #ddd;
    border-radius: 8px;
    font-family: inherit;
    font-size: 14px;
    resize: none;
}

.chat-widget__send-button {
    padding: 10px 16px;
    background: #6a11cb;
    color: white;
    border: none;
    border-radius: 8px;
    cursor: pointer;
    font-weight: 500;
    flex-shrink: 0;
}

.chat-widget__send-button:disabled {
    opacity: 0.5;
    cursor: not-allowed;
}
</style>
