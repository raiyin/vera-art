<template>
    <div
        class="admin-search"
        :class="{ 'admin-search--dark': isDark, }"
    >
        <svg
            class="admin-search__icon"
            viewBox="0 0 24 24"
            fill="none"
            stroke="currentColor"
            stroke-width="2"
            width="16"
            height="16"
        >
            <circle
                cx="11"
                cy="11"
                r="8"
            />
            <line
                x1="21"
                y1="21"
                x2="16.65"
                y2="16.65"
            />
        </svg>
        <input
            :value="modelValue"
            type="text"
            class="admin-search__input"
            :placeholder="placeholder"
            @input="onInput"
            @keydown.esc="onClear"
        >
        <button
            v-if="modelValue"
            class="admin-search__clear"
            @click="onClear"
        >
            <svg
                viewBox="0 0 24 24"
                fill="none"
                stroke="currentColor"
                stroke-width="2"
                width="14"
                height="14"
            >
                <line
                    x1="18"
                    y1="6"
                    x2="6"
                    y2="18"
                />
                <line
                    x1="6"
                    y1="6"
                    x2="18"
                    y2="18"
                />
            </svg>
        </button>
    </div>
</template>

<script setup lang="ts">
import { ref, } from 'vue';

const props = withDefaults(
    defineProps<{
        modelValue: string
        placeholder?: string
        isDark?: boolean
    }>(),
    {
        placeholder: 'Поиск...',
        isDark: false,
    }
);

const emit = defineEmits<{
    'update:modelValue': [value: string,]
}>();

let debounceTimer: ReturnType<typeof setTimeout> | null = null;

function onInput(e: Event,) {
    const val = (e.target as HTMLInputElement).value;
    if (debounceTimer) clearTimeout(debounceTimer,);
    debounceTimer = setTimeout(() => {
        emit('update:modelValue', val,);
    }, 300,);
}

function onClear() {
    emit('update:modelValue', '',);
}
</script>

<style scoped>
.admin-search {
    display: flex;
    align-items: center;
    gap: 8px;
    background: var(--admin-bg, #f0f2f5);
    border: 1px solid var(--admin-border, #e0e0e0);
    border-radius: 8px;
    padding: 8px 12px;
    min-width: 200px;
    max-width: 320px;
    transition: all 0.2s;
}

.admin-search--dark {
    background: rgba(255, 255, 255, 0.05);
    border-color: var(--admin-border, #2d2d3d);
}

.admin-search:focus-within {
    border-color: var(--admin-primary, #6c5ce7);
    box-shadow: 0 0 0 3px rgba(108, 92, 231, 0.1);
}

.admin-search__icon {
    color: var(--admin-text-secondary, #636e72);
    flex-shrink: 0;
}

.admin-search__input {
    flex: 1;
    background: none;
    border: none;
    outline: none;
    color: var(--admin-text-primary, #2d3436);
    font-size: 14px;
    min-width: 0;
}

.admin-search--dark .admin-search__input {
    color: var(--admin-text-primary, #e0e0e0);
}

.admin-search__input::placeholder {
    color: var(--admin-text-secondary, #636e72);
}

.admin-search__clear {
    display: flex;
    align-items: center;
    justify-content: center;
    background: none;
    border: none;
    color: var(--admin-text-secondary, #636e72);
    cursor: pointer;
    padding: 2px;
    border-radius: 4px;
    transition: all 0.2s;
}

.admin-search__clear:hover {
    color: var(--admin-text-primary, #2d3436);
    background: rgba(0, 0, 0, 0.05);
}

.admin-search--dark .admin-search__clear:hover {
    background: rgba(255, 255, 255, 0.1);
}
</style>
