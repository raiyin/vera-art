<script setup lang="ts">
import type { DashboardStats } from '~/types/dashboard';

interface StatCardConfig {
    key: keyof DashboardStats;
    label: string;
    icon: string;
    color: string;
    prefix?: string;
    suffix?: string;
    format?: 'number' | 'price';
}

const props = defineProps<{
    stat: StatCardConfig;
    stats: DashboardStats;
}>();

function formatPrice(value: number): string {
    return new Intl.NumberFormat('ru-RU', {
        style: 'currency',
        currency: 'RUB',
        minimumFractionDigits: 0,
        maximumFractionDigits: 0,
    }).format(value);
}

function formatStatValue(): string {
    const value = props.stats[props.stat.key];
    if (value === undefined || value === null) return '0';
    if (props.stat.format === 'price') {
        return formatPrice(value as number);
    }
    return (value as number).toLocaleString('ru-RU');
}
</script>

<template>
    <UCard class="admin-dashboard__stat-card" :ui="{ body: 'p-0' }">
        <div class="admin-dashboard__stat-inner">
            <div
                class="admin-dashboard__stat-icon"
                :class="`admin-dashboard__stat-icon--${stat.color}`"
            >
                <UIcon :name="stat.icon" class="size-5" />
            </div>
            <div class="admin-dashboard__stat-info">
                <span class="admin-dashboard__stat-value">{{ formatStatValue() }}</span>
                <span class="admin-dashboard__stat-label">{{ stat.label }}</span>
            </div>
        </div>
    </UCard>
</template>

<style scoped>
.admin-dashboard__stat-card {
    transition: transform 0.2s ease, box-shadow 0.2s ease;
}

.admin-dashboard__stat-card:hover {
    transform: translateY(-2px);
    box-shadow: 0 4px 12px rgba(0, 0, 0, 0.08);
}

.admin-dashboard__stat-inner {
    display: flex;
    align-items: center;
    gap: 16px;
    padding: 20px;
}

.admin-dashboard__stat-icon {
    width: 48px;
    height: 48px;
    border-radius: 12px;
    display: flex;
    align-items: center;
    justify-content: center;
    flex-shrink: 0;
}

.admin-dashboard__stat-icon--purple {
    background: rgba(108, 92, 231, 0.1);
    color: #6c5ce7;
}

.admin-dashboard__stat-icon--green {
    background: rgba(0, 184, 148, 0.1);
    color: #00b894;
}

.admin-dashboard__stat-icon--blue {
    background: rgba(116, 185, 255, 0.1);
    color: #74b9ff;
}

.admin-dashboard__stat-icon--orange {
    background: rgba(253, 203, 110, 0.15);
    color: #e17055;
}

.admin-dashboard__stat-icon--pink {
    background: rgba(232, 67, 147, 0.1);
    color: #e84393;
}

.admin-dashboard__stat-icon--red {
    background: rgba(225, 112, 85, 0.1);
    color: #e17055;
}

.admin-dashboard__stat-icon--teal {
    background: rgba(0, 206, 201, 0.1);
    color: #00cec9;
}

.admin-dashboard__stat-icon--yellow {
    background: rgba(253, 203, 110, 0.15);
    color: #fdcb6e;
}

.admin-dashboard__stat-info {
    display: flex;
    flex-direction: column;
    gap: 2px;
}

.admin-dashboard__stat-value {
    font-size: 24px;
    font-weight: 700;
    color: var(--admin-text-primary, #2d3436);
    line-height: 1.2;
}

.admin-dashboard__stat-label {
    font-size: 13px;
    color: var(--admin-text-secondary, #636e72);
    line-height: 1.3;
}
</style>
