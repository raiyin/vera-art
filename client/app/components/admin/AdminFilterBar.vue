<template>
    <div class="admin-filter-bar" :class="{ 'admin-filter-bar--dark': isDark }">
        <div class="admin-filter-bar__filters">
            <template v-for="filter in filters" :key="filter.key">
                <!-- Text filter -->
                <div v-if="filter.type === 'text'" class="admin-filter-bar__item">
                    <label class="admin-filter-bar__label">{{ filter.label }}</label>
                    <input
                        :value="filterValues[filter.key] || ''"
                        type="text"
                        class="admin-filter-bar__input"
                        :placeholder="filter.placeholder || filter.label"
                        @input="onFilterChange(filter.key, ($event.target as HTMLInputElement).value)"
                    />
                </div>

                <!-- Select filter -->
                <div v-else-if="filter.type === 'select'" class="admin-filter-bar__item">
                    <label class="admin-filter-bar__label">{{ filter.label }}</label>
                    <select
                        :value="filterValues[filter.key] || ''"
                        class="admin-filter-bar__select"
                        @change="onFilterChange(filter.key, ($event.target as HTMLSelectElement).value)"
                    >
                        <option value="">Все</option>
                        <option
                            v-for="opt in filter.options"
                            :key="String(opt.value)"
                            :value="String(opt.value)"
                        >
                            {{ opt.label }}
                        </option>
                    </select>
                </div>

                <!-- Number range filter -->
                <div
                    v-else-if="filter.type === 'number-range'"
                    class="admin-filter-bar__item admin-filter-bar__item--range"
                >
                    <label class="admin-filter-bar__label">{{ filter.label }}</label>
                    <div class="admin-filter-bar__range">
                        <input
                            :value="filterValues[`${filter.key}_min`] || ''"
                            type="number"
                            class="admin-filter-bar__input admin-filter-bar__input--sm"
                            placeholder="От"
                            @input="onFilterChange(`${filter.key}_min`, ($event.target as HTMLInputElement).value)"
                        />
                        <span class="admin-filter-bar__range-sep">—</span>
                        <input
                            :value="filterValues[`${filter.key}_max`] || ''"
                            type="number"
                            class="admin-filter-bar__input admin-filter-bar__input--sm"
                            placeholder="До"
                            @input="onFilterChange(`${filter.key}_max`, ($event.target as HTMLInputElement).value)"
                        />
                    </div>
                </div>
            </template>
        </div>

        <div class="admin-filter-bar__actions">
            <button
                v-if="hasActiveFilters"
                class="admin-filter-bar__clear-btn"
                @click="clearFilters"
            >
                <svg
                    viewBox="0 0 24 24"
                    fill="none"
                    stroke="currentColor"
                    stroke-width="2"
                    width="14"
                    height="14"
                >
                    <line x1="18" y1="6" x2="6" y2="18" />
                    <line x1="6" y1="6" x2="18" y2="18" />
                </svg>
                Сбросить
            </button>
        </div>
    </div>
</template>

<script setup lang="ts">
import { ref, computed, watch } from 'vue';

export interface FilterDef {
    key: string;
    label: string;
    type: 'text' | 'select' | 'date-range' | 'number-range';
    options?: { label: string; value: string | number | boolean }[];
    placeholder?: string;
}

const props = withDefaults(
    defineProps<{
        filters: FilterDef[];
        isDark?: boolean;
    }>(),
    {
        isDark: false,
    }
);

const emit = defineEmits<{
    'update:filter': [filters: Record<string, any>];
}>();

const filterValues = ref<Record<string, string>>({});

const hasActiveFilters = computed(() => {
    return Object.values(filterValues.value).some((v) => v !== '');
});

let debounceTimer: ReturnType<typeof setTimeout> | null = null;

function onFilterChange(key: string, value: string) {
    filterValues.value[key] = value;

    if (debounceTimer) clearTimeout(debounceTimer);
    debounceTimer = setTimeout(() => {
        emit('update:filter', { ...filterValues.value });
    }, 300);
}

function clearFilters() {
    filterValues.value = {};
    emit('update:filter', {});
}
</script>

<style scoped>
.admin-filter-bar {
    display: flex;
    align-items: flex-end;
    gap: 12px;
    padding: 12px 20px;
    border-bottom: 1px solid var(--admin-border, #e0e0e0);
    flex-wrap: wrap;
}

.admin-filter-bar--dark {
    border-color: var(--admin-border, #2d2d3d);
}

.admin-filter-bar__filters {
    display: flex;
    align-items: flex-end;
    gap: 12px;
    flex-wrap: wrap;
    flex: 1;
}

.admin-filter-bar__item {
    display: flex;
    flex-direction: column;
    gap: 4px;
    min-width: 160px;
}

.admin-filter-bar__item--range {
    min-width: 220px;
}

.admin-filter-bar__label {
    font-size: 12px;
    font-weight: 600;
    color: var(--admin-text-secondary, #636e72);
    text-transform: uppercase;
    letter-spacing: 0.3px;
}

.admin-filter-bar__input,
.admin-filter-bar__select {
    padding: 8px 12px;
    border: 1px solid var(--admin-border, #e0e0e0);
    border-radius: 8px;
    background: var(--admin-surface, #ffffff);
    color: var(--admin-text-primary, #2d3436);
    font-size: 13px;
    outline: none;
    transition: border-color 0.2s;
}

.admin-filter-bar--dark .admin-filter-bar__input,
.admin-filter-bar--dark .admin-filter-bar__select {
    background: rgba(255, 255, 255, 0.05);
    border-color: var(--admin-border, #2d2d3d);
    color: var(--admin-text-primary, #e0e0e0);
}

.admin-filter-bar__input:focus,
.admin-filter-bar__select:focus {
    border-color: var(--admin-primary, #6c5ce7);
}

.admin-filter-bar__input--sm {
    width: 80px;
    text-align: center;
}

.admin-filter-bar__range {
    display: flex;
    align-items: center;
    gap: 6px;
}

.admin-filter-bar__range-sep {
    color: var(--admin-text-secondary, #636e72);
    font-size: 14px;
}

.admin-filter-bar__actions {
    display: flex;
    align-items: center;
}

.admin-filter-bar__clear-btn {
    display: inline-flex;
    align-items: center;
    gap: 6px;
    padding: 8px 14px;
    border: 1px solid var(--admin-border, #e0e0e0);
    border-radius: 8px;
    background: var(--admin-surface, #ffffff);
    color: var(--admin-text-secondary, #636e72);
    font-size: 13px;
    cursor: pointer;
    transition: all 0.2s;
    white-space: nowrap;
}

.admin-filter-bar--dark .admin-filter-bar__clear-btn {
    background: rgba(255, 255, 255, 0.05);
    border-color: var(--admin-border, #2d2d3d);
    color: var(--admin-text-secondary, #a0a5b5);
}

.admin-filter-bar__clear-btn:hover {
    border-color: var(--admin-danger, #e17055);
    color: var(--admin-danger, #e17055);
}
</style>
