import { defineStore, } from 'pinia';
import { ref, computed, } from 'vue';
import { useAuthStore, } from './AuthStore';
import { getHttpClient, } from '~/api/http-client';
import type { ChatThread, } from '~/types';

export const useNotificationStore = defineStore('notificationStore', () => {
    const authStore = useAuthStore();

    // Unread message count
    const unreadMessageCount = ref(0,);

    // Last checked timestamp
    const lastChecked = ref<Date | null>(null,);

    // Whether notifications are enabled
    const notificationsEnabled = ref(true,);

    // Polling interval ID
    let pollingInterval: number | null = null;

    // Check for new messages
    const checkNewMessages = async () => {
        if (!authStore.isAuthenticated) {
            unreadMessageCount.value = 0;
            return;
        }

        try {
            const { data, } = await getHttpClient().get('chat/threads',);

            const threads = data as ChatThread[];

            let totalUnread = 0;
            if (Array.isArray(threads,)) {
                // Simple logic: count threads with unread messages
                // In a real app, we'd have a proper unread flag per message
                threads.forEach((thread: ChatThread,) => {
                    if (thread.messages && thread.messages.length > 0) {
                        // Check if last message is from someone else and not read
                        const lastMessage = thread.messages[thread.messages.length - 1];
                        if (lastMessage?.sender_id !== authStore.userId && !lastMessage?.is_read) {
                            totalUnread++;
                        }
                    }
                },);
            }

            unreadMessageCount.value = totalUnread;
            lastChecked.value = new Date();
        } catch (error) {
            console.error('Failed to check new messages:', error,);
        }
    };

    // Start polling for new messages
    const startPolling = (intervalMs = 30000,) => {
        if (pollingInterval) {
            window.clearInterval(pollingInterval,);
        }

        // Initial check
        checkNewMessages();

        // Set up interval
        pollingInterval = window.setInterval(checkNewMessages, intervalMs,);
    };

    // Stop polling
    const stopPolling = () => {
        if (pollingInterval) {
            window.clearInterval(pollingInterval,);
            pollingInterval = null;
        }
    };

    // Mark all messages as read
    const markAllAsRead = async () => {
        unreadMessageCount.value = 0;
        lastChecked.value = new Date();

    // In a real implementation, we'd call an API endpoint
    // to mark all messages as read
    };

    // Reset notifications
    const reset = () => {
        unreadMessageCount.value = 0;
        lastChecked.value = null;
        stopPolling();
    };

    // Computed properties
    const hasUnread = computed(() => unreadMessageCount.value > 0,);
    const unreadCount = computed(() => unreadMessageCount.value,);
    const lastCheckedTime = computed(() => lastChecked.value,);

    return {
    // State
        unreadMessageCount,
        lastChecked,
        notificationsEnabled,

        // Computed
        hasUnread,
        unreadCount,
        lastCheckedTime,

        // Actions
        checkNewMessages,
        startPolling,
        stopPolling,
        markAllAsRead,
        reset,
    };
},);
