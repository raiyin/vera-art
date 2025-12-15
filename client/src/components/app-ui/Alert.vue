<template>
    <div
        v-if="show"
        :class="['alert', `alert-${type}`, 'alert-dismissible', 'fade', 'show']"
        role="alert"
    >
        <div class="alert-content">
            <svg
                class="alert-icon"
                xmlns="http://www.w3.org/2000/svg"
                width="24"
                height="24"
                viewBox="0 0 24 24"
                fill="none"
                stroke="currentColor"
                stroke-width="2"
                stroke-linecap="round"
                stroke-linejoin="round"
            >
                <path
                    v-if="type === 'success'"
                    d="M22 11.08V12a10 10 0 1 1-5.93-9.14"
                ></path>
                <path v-if="type === 'success'" d="M22 4 12 14.01 9 11.01"></path>
                <circle v-if="type === 'danger'" cx="12" cy="12" r="10"></circle>
                <line v-if="type === 'danger'" x1="12" y1="8" x2="12" y2="12"></line>
                <line v-if="type === 'danger'" x1="12" y1="16" x2="12.01" y2="16"></line>
            </svg>
            <div>
                <strong>{{ title }}</strong>
                <p>{{ message }}</p>
            </div>
        </div>
        <button
            type="button"
            class="btn-close"
            @click="close"
            :aria-label="closeButtonText"
        ></button>
    </div>
</template>

<script lang="ts">
import { defineComponent } from 'vue';

export default defineComponent({
    name: 'Alert',
    props: {
        type: {
            type: String,
            required: true,
            validator: (value: string) => ['success', 'danger'].includes(value),
        },
        title: {
            type: String,
            required: true,
        },
        message: {
            type: String,
            required: true,
        },
        closeButtonText: {
            type: String,
            default: 'Закрыть',
        },
        modelValue: {
            type: Boolean,
            default: false,
        },
    },
    emits: ['update:modelValue'],
    computed: {
        show: {
            get(): boolean {
                return this.modelValue;
            },
            set(value: boolean) {
                this.$emit('update:modelValue', value);
            },
        },
    },
    methods: {
        close() {
            this.show = false;
        },
    },
});
</script>

<style scoped>
.alert {
    position: fixed;
    bottom: 1rem;
    right: 1rem;
    padding: 1rem;
    border-radius: 6px;
    box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
    display: flex;
    align-items: flex-start;
    gap: 0.75rem;
    max-width: 350px;
    z-index: 1000;
}

.alert-content {
    display: flex;
    align-items: flex-start;
    gap: 0.75rem;
    flex: 1;
}

.alert-icon {
    width: 1.5rem;
    height: 1.5rem;
    flex-shrink: 0;
}

.alert-success {
    background-color: #d4edda;
    border: 1px solid #c3e6cb;
    color: #155724;
}

.alert-success .alert-icon {
    color: #28a745;
}

.alert-danger {
    background-color: #f8d7da;
    border: 1px solid #f5c6cb;
    color: #721c24;
}

.alert-danger .alert-icon {
    color: #dc3545;
}

.btn-close {
    border: none;
    font-size: 1.25rem;
    cursor: pointer;
    padding: 0;
    width: 1.5rem;
    height: 1.5rem;
    display: flex;
    align-items: center;
    justify-content: center;
    color: inherit;
    opacity: 0.7;
    background: none;
}

.btn-close:hover {
    opacity: 1;
}

@media (max-width: 768px) {
    .alert {
        right: 0.5rem;
        left: 0.5rem;
        max-width: none;
    }
}
</style>
