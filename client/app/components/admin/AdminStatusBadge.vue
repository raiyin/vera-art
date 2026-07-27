<template>
    <span
        class="admin-badge"
        :class="[`admin-badge--${variant}`, { 'admin-badge--dark': isDark, },]"
    >
        <span
            v-if="showDot"
            class="admin-badge__dot"
        />
        {{ label }}
    </span>
</template>

<script setup lang="ts">
    import { computed, } from 'vue';

    const props = withDefaults(
        defineProps<{
            status: string
            mapping?: Record<
                string,
                {
                    label: string
                    variant: 'success' | 'warning' | 'danger' | 'info' | 'neutral'
                }
            >
            isDark?: boolean
            showDot?: boolean
        }>(),
        {
            status: '',
            isDark: false,
            showDot: true,
        }
    );

    const defaultMapping: Record<
        string,
        { label: string, variant: 'success' | 'warning' | 'danger' | 'info' | 'neutral' }
    > = {
        published: { label: 'Опубликовано', variant: 'success', },
        active: { label: 'Активно', variant: 'success', },
        draft: { label: 'Черновик', variant: 'warning', },
        pending: { label: 'Ожидает', variant: 'warning', },
        archived: { label: 'Архив', variant: 'neutral', },
        cancelled: { label: 'Отменён', variant: 'danger', },
        refunded: { label: 'Возврат', variant: 'danger', },
        expired: { label: 'Истёк', variant: 'danger', },
        blocked: { label: 'Заблокирован', variant: 'danger', },
        resolved: { label: 'Решён', variant: 'success', },
        open: { label: 'Открыт', variant: 'info', },
        approved: { label: 'Одобрен', variant: 'success', },
        rejected: { label: 'Отклонён', variant: 'danger', },
        succeeded: { label: 'Успешно', variant: 'success', },
        beginner: { label: 'Начальный', variant: 'info', },
        intermediate: { label: 'Средний', variant: 'warning', },
        advanced: { label: 'Продвинутый', variant: 'danger', },
        true: { label: 'Да', variant: 'success', },
        false: { label: 'Нет', variant: 'neutral', },
    };

    const resolvedMapping = computed(() => ({
        ...defaultMapping,
        ...props.mapping,
    }),);

    const statusInfo = computed(() => {
        const key = String(props.status,).toLowerCase();
        return (
            resolvedMapping.value[key] || {
                label: props.status || '—',
                variant: 'neutral' as const,
            }
        );
    });

    const label = computed(() => statusInfo.value.label,);
    const variant = computed(() => statusInfo.value.variant,);
</script>

<style scoped>
.admin-badge {
    display: inline-flex;
    align-items: center;
    gap: 5px;
    padding: 3px 10px;
    border-radius: 20px;
    font-size: 12px;
    font-weight: 600;
    white-space: nowrap;
    line-height: 1.4;
}

.admin-badge__dot {
    width: 6px;
    height: 6px;
    border-radius: 50%;
    flex-shrink: 0;
}

/* Success */
.admin-badge--success {
    background: rgba(0, 184, 148, 0.12);
    color: #00b894;
}

.admin-badge--success .admin-badge__dot {
    background: #00b894;
}

/* Warning */
.admin-badge--warning {
    background: rgba(253, 203, 110, 0.15);
    color: #d4a017;
}

.admin-badge--warning .admin-badge__dot {
    background: #fdcb6e;
}

/* Danger */
.admin-badge--danger {
    background: rgba(225, 112, 85, 0.12);
    color: #e17055;
}

.admin-badge--danger .admin-badge__dot {
    background: #e17055;
}

/* Info */
.admin-badge--info {
    background: rgba(108, 92, 231, 0.12);
    color: #6c5ce7;
}

.admin-badge--info .admin-badge__dot {
    background: #6c5ce7;
}

/* Neutral */
.admin-badge--neutral {
    background: rgba(99, 110, 114, 0.1);
    color: var(--admin-text-secondary, #636e72);
}

.admin-badge--neutral .admin-badge__dot {
    background: var(--admin-text-secondary, #636e72);
}

/* Dark theme adjustments */
.admin-badge--dark.admin-badge--success {
    background: rgba(0, 184, 148, 0.2);
}

.admin-badge--dark.admin-badge--warning {
    background: rgba(253, 203, 110, 0.2);
}

.admin-badge--dark.admin-badge--danger {
    background: rgba(225, 112, 85, 0.2);
}

.admin-badge--dark.admin-badge--info {
    background: rgba(108, 92, 231, 0.2);
}

.admin-badge--dark.admin-badge--neutral {
    background: rgba(160, 165, 181, 0.15);
    color: #a0a5b5;
}
</style>
