<template>
    <div class="admin-table" :class="{ 'admin-table--dark': isDark }">
        <!-- Toolbar -->
        <div
            v-if="$slots.toolbar || showSearch || showFilters || showExport"
            class="admin-table__toolbar"
        >
            <div class="admin-table__toolbar-left">
                <slot name="toolbar-left" />
            </div>
            <div class="admin-table__toolbar-right">
                <slot name="toolbar" />
                <UButton
                    v-if="showExport && data.length > 0"
                    icon="i-lucide-download"
                    color="neutral"
                    variant="ghost"
                    size="sm"
                    :title="exportLabel"
                    @click="exportToCSV"
                >
                    {{ exportLabel }}
                </UButton>
                <AdminSearchInput
                    v-if="showSearch"
                    v-model="internalSearch"
                    :placeholder="searchPlaceholder"
                    :is-dark="isDark"
                    @update:model-value="$emit('search', $event)"
                />
            </div>
        </div>

        <!-- Filter bar -->
        <AdminFilterBar
            v-if="showFilters && filters.length > 0"
            :filters="filters"
            :is-dark="isDark"
            @update:filter="$emit('filter', $event)"
        />

        <!-- Loading state -->
        <div v-if="loading" class="admin-table__loading">
            <AdminLoadingSkeleton
                :rows="loadingRows"
                :columns="columns.length"
                :is-dark="isDark"
            />
        </div>

        <!-- Empty state -->
        <AdminEmptyState
            v-else-if="!loading && (!data || data.length === 0)"
            :icon="emptyIcon"
            :title="emptyTitle"
            :description="emptyDescription"
            :is-dark="isDark"
        >
            <slot name="empty-action" />
        </AdminEmptyState>

        <!-- Table -->
        <div v-else class="admin-table__wrapper">
            <table class="admin-table__table">
                <thead class="admin-table__head">
                    <tr>
                        <th
                            v-if="selectable"
                            class="admin-table__cell admin-table__cell--checkbox"
                        >
                            <input
                                type="checkbox"
                                :checked="allSelected"
                                :indeterminate="someSelected"
                                @change="toggleSelectAll"
                            />
                        </th>
                        <th
                            v-for="col in columns"
                            :key="col.key"
                            class="admin-table__cell admin-table__cell--head"
                            :class="{
                                'admin-table__cell--sortable': col.sortable,
                                'admin-table__cell--sorted': sortKey === col.key,
                            }"
                            :style="col.width ? { width: col.width } : {}"
                            @click="col.sortable && toggleSort(col.key)"
                        >
                            <div class="admin-table__head-content">
                                <span>{{ col.label }}</span>
                                <span v-if="col.sortable" class="admin-table__sort-icon">
                                    <svg
                                        v-if="sortKey === col.key && sortDir === 'asc'"
                                        viewBox="0 0 24 24"
                                        fill="none"
                                        stroke="currentColor"
                                        stroke-width="2"
                                        width="14"
                                        height="14"
                                    >
                                        <polyline points="18 15 12 9 6 15" />
                                    </svg>
                                    <svg
                                        v-else-if="
                                            sortKey === col.key && sortDir === 'desc'
                                        "
                                        viewBox="0 0 24 24"
                                        fill="none"
                                        stroke="currentColor"
                                        stroke-width="2"
                                        width="14"
                                        height="14"
                                    >
                                        <polyline points="6 9 12 15 18 9" />
                                    </svg>
                                    <svg
                                        v-else
                                        viewBox="0 0 24 24"
                                        fill="none"
                                        stroke="currentColor"
                                        stroke-width="2"
                                        width="14"
                                        height="14"
                                        class="admin-table__sort-icon--inactive"
                                    >
                                        <polyline points="18 15 12 9 6 15" />
                                    </svg>
                                </span>
                            </div>
                        </th>
                        <th
                            v-if="$slots.actions"
                            class="admin-table__cell admin-table__cell--head admin-table__cell--actions"
                        >
                            {{ actionsLabel }}
                        </th>
                    </tr>
                </thead>
                <tbody class="admin-table__body">
                    <tr
                        v-for="(row, rowIdx) in data"
                        :key="row.id || rowIdx"
                        class="admin-table__row"
                        :class="{
                            'admin-table__row--selected': isSelected(row),
                            'admin-table__row--clickable': rowClickable,
                        }"
                        @click="rowClickable && $emit('rowClick', row)"
                    >
                        <td
                            v-if="selectable"
                            class="admin-table__cell admin-table__cell--checkbox"
                            @click.stop
                        >
                            <input
                                type="checkbox"
                                :checked="isSelected(row)"
                                @change="toggleSelect(row)"
                            />
                        </td>
                        <td
                            v-for="col in columns"
                            :key="`${row.id || rowIdx}-${col.key}`"
                            class="admin-table__cell"
                            :class="{
                                'admin-table__cell--mono': col.mono,
                                'admin-table__cell--nowrap': col.nowrap,
                            }"
                        >
                            <!-- Custom cell render -->
                            <slot
                                :name="`cell-${col.key}`"
                                :row="row"
                                :value="getNestedValue(row, col.key)"
                            >
                                <!-- Status badge -->
                                <AdminStatusBadge
                                    v-if="col.type === 'status'"
                                    :status="getNestedValue(row, col.key)"
                                    :is-dark="isDark"
                                />
                                <!-- Image preview -->
                                <div
                                    v-else-if="col.type === 'image'"
                                    class="admin-table__image"
                                >
                                    <img
                                        v-if="getNestedValue(row, col.key)"
                                        :src="getNestedValue(row, col.key)"
                                        alt="preview"
                                        class="admin-table__thumb"
                                        @error="onImgError"
                                    />
                                    <div v-else class="admin-table__image-placeholder">
                                        <svg
                                            viewBox="0 0 24 24"
                                            fill="none"
                                            stroke="currentColor"
                                            stroke-width="2"
                                            width="16"
                                            height="16"
                                        >
                                            <rect
                                                x="3"
                                                y="3"
                                                width="18"
                                                height="18"
                                                rx="2"
                                            />
                                            <circle cx="8.5" cy="8.5" r="1.5" />
                                            <path d="M21 15l-5-5L5 21" />
                                        </svg>
                                    </div>
                                </div>
                                <!-- Date -->
                                <span
                                    v-else-if="col.type === 'date'"
                                    class="admin-table__date"
                                >
                                    {{ formatDate(getNestedValue(row, col.key)) }}
                                </span>
                                <!-- Datetime -->
                                <span
                                    v-else-if="col.type === 'datetime'"
                                    class="admin-table__date"
                                >
                                    {{ formatDateTime(getNestedValue(row, col.key)) }}
                                </span>
                                <!-- Number with formatting -->
                                <span
                                    v-else-if="col.type === 'number'"
                                    class="admin-table__number"
                                >
                                    {{ formatNumber(getNestedValue(row, col.key)) }}
                                </span>
                                <!-- Price -->
                                <span
                                    v-else-if="col.type === 'price'"
                                    class="admin-table__price"
                                >
                                    {{ formatPrice(getNestedValue(row, col.key)) }}
                                </span>
                                <!-- Boolean -->
                                <span
                                    v-else-if="col.type === 'boolean'"
                                    class="admin-table__boolean"
                                    :class="{
                                        'admin-table__boolean--true': getNestedValue(
                                            row,
                                            col.key
                                        ),
                                    }"
                                >
                                    {{ getNestedValue(row, col.key) ? '✓' : '—' }}
                                </span>
                                <!-- Default text -->
                                <span v-else>
                                    {{ getNestedValue(row, col.key) ?? '—' }}
                                </span>
                            </slot>
                        </td>
                        <td
                            v-if="$slots.actions"
                            class="admin-table__cell admin-table__cell--actions"
                            @click.stop
                        >
                            <slot name="actions" :row="row" />
                        </td>
                    </tr>
                </tbody>
            </table>
        </div>

        <!-- Footer with pagination -->
        <div v-if="showPagination && total > perPage" class="admin-table__footer">
            <div class="admin-table__footer-info">
                {{ paginationInfo }}
            </div>
            <AdminPagination
                v-model:page="internalPage"
                :total="total"
                :per-page="perPage"
                :is-dark="isDark"
                @update:page="$emit('pageChange', $event)"
            />
        </div>
    </div>
</template>

<script setup lang="ts">
import { ref, computed, watch } from 'vue';

export interface TableColumn {
    key: string;
    label: string;
    sortable?: boolean;
    type?:
        | 'text'
        | 'status'
        | 'image'
        | 'date'
        | 'datetime'
        | 'number'
        | 'price'
        | 'boolean';
    width?: string;
    mono?: boolean;
    nowrap?: boolean;
}

export interface TableFilter {
    key: string;
    label: string;
    type: 'text' | 'select' | 'date-range' | 'number-range';
    options?: { label: string; value: string | number | boolean }[];
    placeholder?: string;
}

function downloadBlob(content: string, filename: string, mimeType: string) {
    const blob = new Blob([content], { type: mimeType });
    const url = URL.createObjectURL(blob);
    const a = document.createElement('a');
    a.href = url;
    a.download = filename;
    document.body.appendChild(a);
    a.click();
    document.body.removeChild(a);
    URL.revokeObjectURL(url);
}

const props = withDefaults(
    defineProps<{
        columns: TableColumn[];
        data: Record<string, any>[];
        total?: number;
        page?: number;
        perPage?: number;
        loading?: boolean;
        loadingRows?: number;
        selectable?: boolean;
        selected?: any[];
        rowClickable?: boolean;
        showSearch?: boolean;
        showFilters?: boolean;
        showPagination?: boolean;
        showExport?: boolean;
        exportLabel?: string;
        exportFilename?: string;
        searchPlaceholder?: string;
        actionsLabel?: string;
        emptyIcon?: 'data' | 'search' | 'box' | 'info';
        emptyTitle?: string;
        emptyDescription?: string;
        filters?: TableFilter[];
        isDark?: boolean;
        sortKey?: string;
        sortDir?: 'asc' | 'desc';
        valueKey?: string;
    }>(),
    {
        total: 0,
        page: 1,
        perPage: 20,
        loading: false,
        loadingRows: 5,
        selectable: false,
        selected: () => [],
        rowClickable: false,
        showSearch: false,
        showFilters: false,
        showPagination: true,
        showExport: false,
        exportLabel: 'CSV',
        exportFilename: 'export',
        searchPlaceholder: 'Поиск...',
        actionsLabel: 'Действия',
        emptyIcon: 'data',
        emptyTitle: 'Нет данных',
        emptyDescription: 'По заданным критериям ничего не найдено',
        filters: () => [],
        isDark: false,
        sortKey: '',
        sortDir: 'asc',
        valueKey: 'id',
    }
);

const emit = defineEmits<{
    search: [query: string];
    filter: [filters: Record<string, any>];
    sort: [key: string, dir: 'asc' | 'desc'];
    pageChange: [page: number];
    select: [selected: any[]];
    rowClick: [row: any];
}>();

const internalSearch = ref('');
const internalPage = ref(props.page);

watch(
    () => props.page,
    (val) => {
        internalPage.value = val;
    }
);

const allSelected = computed(() => {
    if (!props.data.length) return false;
    return props.data.every((row) => isSelected(row));
});

const someSelected = computed(() => {
    if (!props.data.length) return false;
    return props.data.some((row) => isSelected(row)) && !allSelected.value;
});

function isSelected(row: any): boolean {
    const id = row[props.valueKey] ?? row;
    return props.selected.includes(id);
}

function toggleSelect(row: any) {
    const id = row[props.valueKey] ?? row;
    const newSelected = isSelected(row)
        ? props.selected.filter((s) => s !== id)
        : [...props.selected, id];
    emit('select', newSelected);
}

function toggleSelectAll() {
    if (allSelected.value) {
        emit('select', []);
    } else {
        const ids = props.data.map((row) => row[props.valueKey] ?? row);
        emit('select', ids);
    }
}

function toggleSort(key: string) {
    const newDir = props.sortKey === key && props.sortDir === 'asc' ? 'desc' : 'asc';
    emit('sort', key, newDir);
}

function getNestedValue(obj: any, path: string): any {
    return path
        .split('.')
        .reduce((acc, part) => (acc != null ? acc[part] : undefined), obj);
}

const paginationInfo = computed(() => {
    const start = (internalPage.value - 1) * props.perPage + 1;
    const end = Math.min(internalPage.value * props.perPage, props.total);
    return `${start}–${end} из ${props.total}`;
});

function formatDate(val: string | null | undefined): string {
    if (!val) return '—';
    try {
        return new Date(val).toLocaleDateString('ru-RU');
    } catch {
        return String(val);
    }
}

function formatDateTime(val: string | null | undefined): string {
    if (!val) return '—';
    try {
        return new Date(val).toLocaleString('ru-RU');
    } catch {
        return String(val);
    }
}

function formatNumber(val: number | null | undefined): string {
    if (val == null) return '—';
    return new Intl.NumberFormat('ru-RU').format(val);
}

function formatPrice(val: number | null | undefined): string {
    if (val == null) return '—';
    return new Intl.NumberFormat('ru-RU', {
        style: 'currency',
        currency: 'RUB',
        minimumFractionDigits: 0,
    }).format(val / 100);
}

function onImgError(e: Event) {
    const target = e.target as HTMLImageElement;
    target.style.display = 'none';
}

function exportToCSV() {
    const visibleColumns = props.columns.filter(
        (col) => col.type !== 'image' && col.key !== 'actions'
    );
    const headers = visibleColumns.map((col) => col.label);
    const rows = props.data.map((row) =>
        visibleColumns
            .map((col) => {
                const val = getNestedValue(row, col.key);
                if (val == null || val === undefined) return '';
                const str = String(val).replace(/"/g, '""');
                return `"${str}"`;
            })
            .join(',')
    );
    const csv = [headers.join(','), ...rows].join('\n');
    const bom = '\uFEFF';
    const filename = `${props.exportFilename}-${new Date()
        .toISOString()
        .slice(0, 10)}.csv`;
    downloadBlob(bom + csv, filename, 'text/csv;charset=utf-8;');
}
</script>

<style scoped>
.admin-table {
    background: var(--admin-surface, #ffffff);
    border: 1px solid var(--admin-border, #e0e0e0);
    border-radius: 12px;
    overflow: hidden;
    transition: background 0.3s ease, border-color 0.3s ease;
}

.admin-table--dark {
    background: var(--admin-surface, #1a1d29);
    border-color: var(--admin-border, #2d2d3d);
}

/* Toolbar */
.admin-table__toolbar {
    display: flex;
    align-items: center;
    justify-content: space-between;
    padding: 16px 20px;
    gap: 12px;
    border-bottom: 1px solid var(--admin-border, #e0e0e0);
    flex-wrap: wrap;
}

.admin-table--dark .admin-table__toolbar {
    border-color: var(--admin-border, #2d2d3d);
}

.admin-table__toolbar-left {
    display: flex;
    align-items: center;
    gap: 8px;
    flex: 1;
}

.admin-table__toolbar-right {
    display: flex;
    align-items: center;
    gap: 8px;
}

/* Loading */
.admin-table__loading {
    padding: 20px;
}

/* Table wrapper */
.admin-table__wrapper {
    overflow-x: auto;
    -webkit-overflow-scrolling: touch;
}

.admin-table__wrapper::-webkit-scrollbar {
    height: 6px;
}

.admin-table__wrapper::-webkit-scrollbar-track {
    background: transparent;
}

.admin-table__wrapper::-webkit-scrollbar-thumb {
    background: var(--admin-border, #e0e0e0);
    border-radius: 3px;
}

/* Table */
.admin-table__table {
    width: 100%;
    border-collapse: collapse;
    font-size: 14px;
}

/* Head */
.admin-table__head {
    background: var(--admin-bg, #f0f2f5);
    border-bottom: 1px solid var(--admin-border, #e0e0e0);
}

.admin-table--dark .admin-table__head {
    background: rgba(255, 255, 255, 0.03);
    border-color: var(--admin-border, #2d2d3d);
}

.admin-table__cell--head {
    padding: 12px 16px;
    font-weight: 600;
    font-size: 12px;
    text-transform: uppercase;
    letter-spacing: 0.5px;
    color: var(--admin-text-secondary, #636e72);
    white-space: nowrap;
    user-select: none;
}

.admin-table__cell--sortable {
    cursor: pointer;
    transition: color 0.2s;
}

.admin-table__cell--sortable:hover {
    color: var(--admin-primary, #6c5ce7);
}

.admin-table__cell--sorted {
    color: var(--admin-primary, #6c5ce7);
}

.admin-table__head-content {
    display: flex;
    align-items: center;
    gap: 4px;
}

.admin-table__sort-icon {
    display: inline-flex;
    align-items: center;
    flex-shrink: 0;
}

.admin-table__sort-icon--inactive {
    opacity: 0.3;
}

/* Body */
.admin-table__body .admin-table__row {
    border-bottom: 1px solid var(--admin-border, #e0e0e0);
    transition: background 0.15s ease, transform 0.15s ease;
    animation: adminTableRowIn 0.3s ease-out both;
}

.admin-table__body .admin-table__row:nth-child(1) {
    animation-delay: 0.01s;
}
.admin-table__body .admin-table__row:nth-child(2) {
    animation-delay: 0.02s;
}
.admin-table__body .admin-table__row:nth-child(3) {
    animation-delay: 0.03s;
}
.admin-table__body .admin-table__row:nth-child(4) {
    animation-delay: 0.04s;
}
.admin-table__body .admin-table__row:nth-child(5) {
    animation-delay: 0.05s;
}
.admin-table__body .admin-table__row:nth-child(6) {
    animation-delay: 0.06s;
}
.admin-table__body .admin-table__row:nth-child(7) {
    animation-delay: 0.07s;
}
.admin-table__body .admin-table__row:nth-child(8) {
    animation-delay: 0.08s;
}
.admin-table__body .admin-table__row:nth-child(9) {
    animation-delay: 0.09s;
}
.admin-table__body .admin-table__row:nth-child(10) {
    animation-delay: 0.1s;
}

@keyframes adminTableRowIn {
    from {
        opacity: 0;
        transform: translateX(-4px);
    }
    to {
        opacity: 1;
        transform: translateX(0);
    }
}

.admin-table--dark .admin-table__body .admin-table__row {
    border-color: var(--admin-border, #2d2d3d);
}

.admin-table__body .admin-table__row:last-child {
    border-bottom: none;
}

.admin-table__body .admin-table__row:hover {
    background: rgba(108, 92, 231, 0.03);
    transform: translateX(2px);
}

.admin-table--dark .admin-table__body .admin-table__row:hover {
    background: rgba(108, 92, 231, 0.08);
}

.admin-table__row--selected {
    background: rgba(108, 92, 231, 0.06) !important;
}

.admin-table--dark .admin-table__row--selected {
    background: rgba(108, 92, 231, 0.12) !important;
}

.admin-table__row--clickable {
    cursor: pointer;
}

/* Cells */
.admin-table__cell {
    padding: 12px 16px;
    color: var(--admin-text-primary, #2d3436);
    vertical-align: middle;
}

.admin-table--dark .admin-table__cell {
    color: var(--admin-text-primary, #e0e0e0);
}

.admin-table__cell--checkbox {
    width: 48px;
    text-align: center;
}

.admin-table__cell--checkbox input[type='checkbox'] {
    width: 16px;
    height: 16px;
    cursor: pointer;
    accent-color: var(--admin-primary, #6c5ce7);
}

.admin-table__cell--actions {
    text-align: right;
    white-space: nowrap;
}

.admin-table__cell--mono {
    font-family: 'JetBrains Mono', 'SF Mono', 'Fira Code', monospace;
    font-size: 13px;
}

.admin-table__cell--nowrap {
    white-space: nowrap;
}

/* Image cell */
.admin-table__image {
    width: 48px;
    height: 48px;
    border-radius: 8px;
    overflow: hidden;
    background: var(--admin-bg, #f0f2f5);
    display: flex;
    align-items: center;
    justify-content: center;
}

.admin-table--dark .admin-table__image {
    background: rgba(255, 255, 255, 0.05);
}

.admin-table__thumb {
    width: 100%;
    height: 100%;
    object-fit: cover;
}

.admin-table__image-placeholder {
    color: var(--admin-text-secondary, #636e72);
    opacity: 0.5;
}

/* Date */
.admin-table__date {
    white-space: nowrap;
    font-size: 13px;
    color: var(--admin-text-secondary, #636e72);
}

/* Number */
.admin-table__number {
    font-family: 'JetBrains Mono', 'SF Mono', 'Fira Code', monospace;
    font-size: 13px;
}

/* Price */
.admin-table__price {
    font-weight: 600;
    color: var(--admin-success, #00b894);
}

/* Boolean */
.admin-table__boolean {
    display: inline-flex;
    align-items: center;
    justify-content: center;
    width: 24px;
    height: 24px;
    border-radius: 50%;
    background: var(--admin-bg, #f0f2f5);
    font-size: 12px;
    color: var(--admin-text-secondary, #636e72);
}

.admin-table--dark .admin-table__boolean {
    background: rgba(255, 255, 255, 0.05);
}

.admin-table__boolean--true {
    background: rgba(0, 184, 148, 0.15);
    color: var(--admin-success, #00b894);
}

/* Footer */
.admin-table__footer {
    display: flex;
    align-items: center;
    justify-content: space-between;
    padding: 12px 20px;
    border-top: 1px solid var(--admin-border, #e0e0e0);
    gap: 16px;
    flex-wrap: wrap;
}

@media (max-width: 480px) {
    .admin-table__footer {
        flex-direction: column;
        align-items: stretch;
        text-align: center;
    }
}

.admin-table--dark .admin-table__footer {
    border-color: var(--admin-border, #2d2d3d);
}

.admin-table__footer-info {
    font-size: 13px;
    color: var(--admin-text-secondary, #636e72);
}

@media (max-width: 480px) {
    .admin-table__toolbar {
        flex-direction: column;
        align-items: stretch;
    }
    .admin-table__toolbar-right {
        flex-wrap: wrap;
    }
    .admin-table__cell--head,
    .admin-table__cell {
        padding: 8px 10px;
        font-size: 13px;
    }
}
</style>
