<template>
    <UContainer class="py-8 md:py-12">
        <!-- Page Header -->
        <div class="text-center mb-8">
            <div class="inline-flex items-center justify-center w-14 h-14 bg-teal-500 rounded-2xl shadow-lg mb-4">
                <UIcon
                    name="i-heroicons-chat-bubble-left-right"
                    class="w-8 h-8 text-white"
                />
            </div>
            <h1 class="text-3xl md:text-4xl font-bold text-gray-900 dark:text-white mb-2">
                Мои чаты
            </h1>
            <p class="text-gray-600 dark:text-gray-400">
                Общайтесь с преподавателями по вашим курсам
            </p>
        </div>

        <!-- Loading State -->
        <div
            v-if="loading"
            class="text-center py-16"
        >
            <UIcon
                name="i-heroicons-arrow-path"
                class="w-12 h-12 mx-auto text-gray-300 dark:text-gray-600 mb-4 animate-spin"
            />
            <p class="text-gray-500 dark:text-gray-400">
                Загрузка чатов...
            </p>
        </div>

        <!-- Error State -->
        <div
            v-else-if="error"
            class="text-center py-16"
        >
            <UCard class="max-w-md mx-auto">
                <div class="flex flex-col items-center gap-4">
                    <UIcon
                        name="i-heroicons-exclamation-triangle"
                        class="w-12 h-12 text-red-400"
                    />
                    <p class="text-gray-600 dark:text-gray-400">
                        {{ error }}
                    </p>
                    <UButton
                        color="primary"
                        variant="outline"
                        @click="loadThreads"
                    >
                        Попробовать снова
                    </UButton>
                </div>
            </UCard>
        </div>

        <!-- Empty State -->
        <div
            v-else-if="threads.length === 0"
            class="text-center py-16"
        >
            <UCard class="max-w-md mx-auto">
                <div class="flex flex-col items-center gap-4 py-8">
                    <div class="w-20 h-20 rounded-full bg-teal-50 dark:bg-teal-900/20 flex items-center justify-center">
                        <UIcon
                            name="i-heroicons-chat-bubble-left-right"
                            class="w-10 h-10 text-teal-400"
                        />
                    </div>
                    <h3 class="text-xl font-semibold text-gray-900 dark:text-white">
                        У вас пока нет активных чатов
                    </h3>
                    <p class="text-gray-500 dark:text-gray-400 max-w-sm">
                        После покупки курса вы сможете создать чат с преподавателем.
                    </p>
                </div>
            </UCard>
        </div>

        <!-- Threads List -->
        <div
            v-else
            class="max-w-3xl mx-auto space-y-4"
        >
            <div
                v-for="thread in threads"
                :key="thread.id"
                class="group bg-white dark:bg-gray-800 rounded-xl border border-gray-200 dark:border-gray-700 p-5 cursor-pointer transition-all duration-200 hover:shadow-lg hover:border-teal-300 dark:hover:border-teal-600"
                :class="{ 'border-l-4 border-l-teal-500 bg-teal-50/50 dark:bg-teal-900/10': hasUnread(thread,), }"
                @click="openThread(thread.id,)"
            >
                <div class="flex items-center gap-4">
                    <div class="w-12 h-12 rounded-full bg-gradient-to-br from-teal-400 to-green-500 flex items-center justify-center text-white font-bold text-lg shrink-0">
                        {{ thread.purchase?.product?.title_ru?.[0] || 'Ч' }}
                    </div>
                    <div class="flex-1 min-w-0">
                        <div class="flex items-center justify-between gap-3 mb-1">
                            <h3 class="font-semibold text-gray-900 dark:text-white truncate">
                                {{ thread.purchase?.product?.title_ru || 'Без названия' }}
                            </h3>
                            <span class="text-xs text-gray-400 dark:text-gray-500 shrink-0">
                                {{ formatDate(thread.last_message_at,) }}
                            </span>
                        </div>
                        <p class="text-sm text-gray-500 dark:text-gray-400 truncate">
                            {{ lastMessagePreview(thread,) }}
                        </p>
                        <div
                            v-if="hasUnread(thread,)"
                            class="mt-2"
                        >
                            <span class="inline-flex items-center gap-1 text-xs font-medium text-teal-600 dark:text-teal-400 bg-teal-50 dark:bg-teal-900/20 px-2 py-0.5 rounded-full">
                                <span class="w-1.5 h-1.5 rounded-full bg-teal-500" />
                                Новое сообщение
                            </span>
                        </div>
                    </div>
                    <UIcon
                        name="i-heroicons-chevron-right"
                        class="w-5 h-5 text-gray-300 dark:text-gray-600 group-hover:text-teal-500 transition-colors shrink-0"
                    />
                </div>
            </div>
        </div>
    </UContainer>
</template>

<script setup lang="ts">
    import { ref, onMounted, computed, } from 'vue';
    import { useRouter, } from 'vue-router';
    import { useAuthStore, } from '~/stores/AuthStore';
    import type { ChatThread, } from '~/types';

    const router = useRouter();
    const authStore = useAuthStore();

    const threads = ref<ChatThread[]>([],);
    const loading = ref(true,);
    const error = ref('',);

    const authToken = computed(() => authStore.token,);

    onMounted(async () => {
        await loadThreads();
    });

    async function loadThreads() {
        loading.value = true;
        error.value = '';
        try {
            const response = await $fetch<ChatThread[]>('/api/chat/threads', {
                headers: { Authorization: `Bearer ${authToken.value}`, },
            });
            threads.value = response;
        } catch (err) {
            error.value = 'Не удалось загрузить чаты';
            console.error(err,);
        } finally {
            loading.value = false;
        }
    }

    function hasUnread(thread: ChatThread,): boolean {
        return false;
    }

    function lastMessagePreview(thread: ChatThread,): string {
        if (!thread.messages || thread.messages.length === 0) {
            return 'Нет сообщений';
        }
        const last = thread.messages[thread.messages.length - 1];
        const content = last?.content || '';
        return content.length > 50 ? content.substring(0, 50,) + '...' : content;
    }

    function formatDate(dateString: string,): string {
        const date = new Date(dateString,);
        const now = new Date();
        const diffMs = now.getTime() - date.getTime();
        const diffDays = Math.floor(diffMs / (1000 * 60 * 60 * 24),);
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

    function openThread(threadId: number,) {
        router.push(`/chat/${threadId}`,);
    }
</script>
