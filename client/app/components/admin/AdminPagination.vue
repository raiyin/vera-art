<template>
    <div
        class="admin-pagination"
        :class="{ 'admin-pagination--dark': isDark, }"
    >
        <button
            class="admin-pagination__btn"
            :disabled="page <= 1"
            title="Предыдущая"
            @click="goTo(page - 1,)"
        >
            <svg
                viewBox="0 0 24 24"
                fill="none"
                stroke="currentColor"
                stroke-width="2"
                width="16"
                height="16"
            >
                <polyline points="15 18 9 12 15 6" />
            </svg>
        </button>

        <template
            v-for="(p, idx) in visiblePages"
            :key="idx"
        >
            <span
                v-if="p === '...'"
                class="admin-pagination__ellipsis"
            >...</span>
            <button
                v-else
                class="admin-pagination__btn"
                :class="{ 'admin-pagination__btn--active': p === page, }"
                @click="goTo(p as number,)"
            >
                {{ p }}
            </button>
        </template>

        <button
            class="admin-pagination__btn"
            :disabled="page >= totalPages"
            title="Следующая"
            @click="goTo(page + 1,)"
        >
            <svg
                viewBox="0 0 24 24"
                fill="none"
                stroke="currentColor"
                stroke-width="2"
                width="16"
                height="16"
            >
                <polyline points="9 18 15 12 9 6" />
            </svg>
        </button>

        <div class="admin-pagination__per-page">
            <select
                :value="perPage"
                class="admin-pagination__select"
                @change="onPerPageChange"
            >
                <option :value="10">
                    10
                </option>
                <option :value="20">
                    20
                </option>
                <option :value="50">
                    50
                </option>
                <option :value="100">
                    100
                </option>
            </select>
        </div>
    </div>
</template>

<script setup lang="ts">
    import { computed, } from 'vue';

    const props = withDefaults(
        defineProps<{
            page: number
            total: number
            perPage?: number
            isDark?: boolean
        }>(),
        {
            perPage: 20,
            isDark: false,
        }
    );

    const emit = defineEmits<{
        'update:page': [page: number,]
        'update:perPage': [perPage: number,]
    }>();

    const totalPages = computed(() => Math.max(1, Math.ceil(props.total / props.perPage,),),);

    const visiblePages = computed(() => {
        const pages: (number | string)[] = [];
        const total = totalPages.value;
        const current = props.page;

        if (total <= 7) {
            for (let i = 1; i <= total; i++) pages.push(i,);
            return pages;
        }

        pages.push(1,);

        if (current > 3) pages.push('...',);

        const start = Math.max(2, current - 1,);
        const end = Math.min(total - 1, current + 1,);

        for (let i = start; i <= end; i++) {
            pages.push(i,);
        }

        if (current < total - 2) pages.push('...',);

        pages.push(total,);

        return pages;
    });

    function goTo(p: number,) {
        if (p >= 1 && p <= totalPages.value) {
            emit('update:page', p,);
        }
    }

    function onPerPageChange(e: Event,) {
        const val = parseInt((e.target as HTMLSelectElement).value,);
        emit('update:perPage', val,);
        emit('update:page', 1,);
    }
</script>

<style scoped>
.admin-pagination {
    display: flex;
    align-items: center;
    gap: 4px;
}

.admin-pagination__btn {
    display: inline-flex;
    align-items: center;
    justify-content: center;
    min-width: 36px;
    height: 36px;
    padding: 0 8px;
    border: 1px solid var(--admin-border, #e0e0e0);
    border-radius: 8px;
    background: var(--admin-surface, #ffffff);
    color: var(--admin-text-primary, #2d3436);
    font-size: 14px;
    font-weight: 500;
    cursor: pointer;
    transition: all 0.2s;
}

.admin-pagination--dark .admin-pagination__btn {
    background: var(--admin-surface, #1a1d29);
    border-color: var(--admin-border, #2d2d3d);
    color: var(--admin-text-primary, #e0e0e0);
}

.admin-pagination__btn:hover:not(:disabled):not(.admin-pagination__btn--active) {
    border-color: var(--admin-primary, #6c5ce7);
    color: var(--admin-primary, #6c5ce7);
}

.admin-pagination__btn--active {
    background: var(--admin-primary, #6c5ce7);
    border-color: var(--admin-primary, #6c5ce7);
    color: white;
}

.admin-pagination__btn:disabled {
    opacity: 0.4;
    cursor: not-allowed;
}

.admin-pagination__ellipsis {
    padding: 0 4px;
    color: var(--admin-text-secondary, #636e72);
    font-size: 14px;
}

.admin-pagination__per-page {
    margin-left: 12px;
}

.admin-pagination__select {
    padding: 6px 10px;
    border: 1px solid var(--admin-border, #e0e0e0);
    border-radius: 8px;
    background: var(--admin-surface, #ffffff);
    color: var(--admin-text-primary, #2d3436);
    font-size: 13px;
    cursor: pointer;
    outline: none;
}

.admin-pagination--dark .admin-pagination__select {
    background: var(--admin-surface, #1a1d29);
    border-color: var(--admin-border, #2d2d3d);
    color: var(--admin-text-primary, #e0e0e0);
}

.admin-pagination__select:focus {
    border-color: var(--admin-primary, #6c5ce7);
}
</style>
