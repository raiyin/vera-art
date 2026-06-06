<template>
    <Teleport to="body">
        <Transition name="admin-confirm">
            <div v-if="visible" class="admin-confirm__overlay" @click.self="onCancel">
                <div
                    class="admin-confirm__dialog"
                    :class="{ 'admin-confirm__dialog--dark': isDark }"
                >
                    <div class="admin-confirm__header">
                        <div
                            class="admin-confirm__icon"
                            :class="`admin-confirm__icon--${type}`"
                        >
                            <svg
                                v-if="type === 'danger'"
                                viewBox="0 0 24 24"
                                fill="none"
                                stroke="currentColor"
                                stroke-width="2"
                                width="24"
                                height="24"
                            >
                                <circle cx="12" cy="12" r="10" />
                                <line x1="12" y1="8" x2="12" y2="12" />
                                <line x1="12" y1="16" x2="12.01" y2="16" />
                            </svg>
                            <svg
                                v-else-if="type === 'warning'"
                                viewBox="0 0 24 24"
                                fill="none"
                                stroke="currentColor"
                                stroke-width="2"
                                width="24"
                                height="24"
                            >
                                <path
                                    d="M10.29 3.86L1.82 18a2 2 0 001.71 3h16.94a2 2 0 001.71-3L13.71 3.86a2 2 0 00-3.42 0z"
                                />
                                <line x1="12" y1="9" x2="12" y2="13" />
                                <line x1="12" y1="17" x2="12.01" y2="17" />
                            </svg>
                            <svg
                                v-else
                                viewBox="0 0 24 24"
                                fill="none"
                                stroke="currentColor"
                                stroke-width="2"
                                width="24"
                                height="24"
                            >
                                <circle cx="12" cy="12" r="10" />
                                <line x1="12" y1="16" x2="12" y2="12" />
                                <line x1="12" y1="8" x2="12.01" y2="8" />
                            </svg>
                        </div>
                        <h3 class="admin-confirm__title">{{ title }}</h3>
                    </div>

                    <p class="admin-confirm__message">{{ message }}</p>

                    <div class="admin-confirm__actions">
                        <button
                            class="admin-confirm__btn admin-confirm__btn--cancel"
                            @click="onCancel"
                        >
                            {{ cancelText }}
                        </button>
                        <button
                            class="admin-confirm__btn"
                            :class="`admin-confirm__btn--${type}`"
                            :disabled="loading"
                            @click="onConfirm"
                        >
                            <span v-if="loading" class="admin-confirm__spinner" />
                            {{ loading ? loadingText : confirmText }}
                        </button>
                    </div>
                </div>
            </div>
        </Transition>
    </Teleport>
</template>

<script setup lang="ts">
import { ref, watch } from 'vue';

const props = withDefaults(
    defineProps<{
        visible: boolean;
        title?: string;
        message?: string;
        type?: 'danger' | 'warning' | 'info';
        confirmText?: string;
        cancelText?: string;
        loadingText?: string;
        loading?: boolean;
        isDark?: boolean;
    }>(),
    {
        title: 'Подтверждение',
        message: 'Вы уверены, что хотите выполнить это действие?',
        type: 'danger',
        confirmText: 'Подтвердить',
        cancelText: 'Отмена',
        loadingText: 'Выполнение...',
        loading: false,
        isDark: false,
    }
);

const emit = defineEmits<{
    confirm: [];
    cancel: [];
    'update:visible': [value: boolean];
}>();

function onConfirm() {
    emit('confirm');
}

function onCancel() {
    emit('cancel');
    emit('update:visible', false);
}
</script>

<style scoped>
.admin-confirm__overlay {
    position: fixed;
    inset: 0;
    background: rgba(0, 0, 0, 0.5);
    display: flex;
    align-items: center;
    justify-content: center;
    z-index: 1000;
    padding: 20px;
}

.admin-confirm__dialog {
    background: var(--admin-surface, #ffffff);
    border: 1px solid var(--admin-border, #e0e0e0);
    border-radius: 16px;
    padding: 28px;
    max-width: 420px;
    width: 100%;
    box-shadow: 0 20px 60px rgba(0, 0, 0, 0.15);
}

.admin-confirm__dialog--dark {
    background: var(--admin-surface, #1a1d29);
    border-color: var(--admin-border, #2d2d3d);
}

.admin-confirm__header {
    display: flex;
    align-items: center;
    gap: 12px;
    margin-bottom: 12px;
}

.admin-confirm__icon {
    display: flex;
    align-items: center;
    justify-content: center;
    width: 44px;
    height: 44px;
    border-radius: 12px;
    flex-shrink: 0;
}

.admin-confirm__icon--danger {
    background: rgba(225, 112, 85, 0.15);
    color: #e17055;
}

.admin-confirm__icon--warning {
    background: rgba(253, 203, 110, 0.15);
    color: #fdcb6e;
}

.admin-confirm__icon--info {
    background: rgba(108, 92, 231, 0.15);
    color: var(--admin-primary, #6c5ce7);
}

.admin-confirm__title {
    font-size: 18px;
    font-weight: 700;
    color: var(--admin-text-primary, #2d3436);
    margin: 0;
}

.admin-confirm__dialog--dark .admin-confirm__title {
    color: var(--admin-text-primary, #e0e0e0);
}

.admin-confirm__message {
    font-size: 14px;
    color: var(--admin-text-secondary, #636e72);
    line-height: 1.5;
    margin: 0 0 24px;
    padding-left: 56px;
}

.admin-confirm__actions {
    display: flex;
    justify-content: flex-end;
    gap: 10px;
}

.admin-confirm__btn {
    display: inline-flex;
    align-items: center;
    gap: 8px;
    padding: 10px 20px;
    border-radius: 8px;
    font-size: 14px;
    font-weight: 600;
    cursor: pointer;
    transition: all 0.2s;
    border: none;
}

.admin-confirm__btn--cancel {
    background: var(--admin-bg, #f0f2f5);
    color: var(--admin-text-primary, #2d3436);
}

.admin-confirm__dialog--dark .admin-confirm__btn--cancel {
    background: rgba(255, 255, 255, 0.08);
    color: var(--admin-text-primary, #e0e0e0);
}

.admin-confirm__btn--cancel:hover {
    background: var(--admin-border, #e0e0e0);
}

.admin-confirm__dialog--dark .admin-confirm__btn--cancel:hover {
    background: rgba(255, 255, 255, 0.12);
}

.admin-confirm__btn--danger {
    background: #e17055;
    color: white;
}

.admin-confirm__btn--danger:hover:not(:disabled) {
    background: #d63031;
}

.admin-confirm__btn--warning {
    background: #fdcb6e;
    color: #2d3436;
}

.admin-confirm__btn--warning:hover:not(:disabled) {
    background: #f9ca24;
}

.admin-confirm__btn--info {
    background: var(--admin-primary, #6c5ce7);
    color: white;
}

.admin-confirm__btn--info:hover:not(:disabled) {
    background: var(--admin-primary-hover, #5a4bd1);
}

.admin-confirm__btn:disabled {
    opacity: 0.6;
    cursor: not-allowed;
}

.admin-confirm__spinner {
    width: 16px;
    height: 16px;
    border: 2px solid rgba(255, 255, 255, 0.3);
    border-top-color: white;
    border-radius: 50%;
    animation: admin-confirm-spin 0.6s linear infinite;
}

@keyframes admin-confirm-spin {
    to {
        transform: rotate(360deg);
    }
}

/* Transition */
.admin-confirm-enter-active,
.admin-confirm-leave-active {
    transition: opacity 0.2s ease;
}

.admin-confirm-enter-active .admin-confirm__dialog,
.admin-confirm-leave-active .admin-confirm__dialog {
    transition: transform 0.2s ease;
}

.admin-confirm-enter-from,
.admin-confirm-leave-to {
    opacity: 0;
}

.admin-confirm-enter-from .admin-confirm__dialog {
    transform: scale(0.95);
}

.admin-confirm-leave-to .admin-confirm__dialog {
    transform: scale(0.95);
}
</style>
