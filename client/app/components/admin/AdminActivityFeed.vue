<script setup lang="ts">
    import type { RecentActivityItem, } from '~/types/dashboard';

    defineProps<{
        activities: RecentActivityItem[]
        loading?: boolean
    }>();

    const emit = defineEmits<{
        refresh: []
    }>();
</script>

<template>
    <UCard class="admin-dashboard__activity-card">
        <template #header>
            <div class="admin-dashboard__section-header">
                <h3 class="admin-dashboard__section-title">
                    Последняя активность
                </h3>
                <UButton
                    v-if="activities.length > 0"
                    color="neutral"
                    variant="ghost"
                    icon="i-lucide-refresh-cw"
                    size="xs"
                    :loading="loading"
                    @click="emit('refresh',)"
                />
            </div>
        </template>
        <div
            v-if="activities.length === 0"
            class="admin-dashboard__activity-empty"
        >
            <UIcon
                name="i-lucide-clock"
                class="size-6 text-gray-400"
            />
            <p>Активность пока отсутствует</p>
        </div>
        <div
            v-else
            class="admin-dashboard__activity-list"
        >
            <div
                v-for="activity in activities"
                :key="activity.id"
                class="admin-dashboard__activity-item"
            >
                <div
                    class="admin-dashboard__activity-dot"
                    :class="`admin-dashboard__activity-dot--${activity.type}`"
                />
                <div class="admin-dashboard__activity-content">
                    <span class="admin-dashboard__activity-text">{{
                        activity.text
                    }}</span>
                    <span class="admin-dashboard__activity-time">{{
                        activity.time
                    }}</span>
                </div>
            </div>
        </div>
    </UCard>
</template>

<style scoped>
.admin-dashboard__section-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
}

.admin-dashboard__section-title {
    font-size: 16px;
    font-weight: 600;
    color: var(--admin-text-primary, #2d3436);
    margin: 0;
}

.admin-dashboard__activity-empty {
    display: flex;
    flex-direction: column;
    align-items: center;
    gap: 8px;
    padding: 32px 20px;
    color: var(--admin-text-secondary, #636e72);
}

.admin-dashboard__activity-empty p {
    margin: 0;
    font-size: 14px;
}

.admin-dashboard__activity-list {
    display: flex;
    flex-direction: column;
}

.admin-dashboard__activity-item {
    display: flex;
    align-items: flex-start;
    gap: 12px;
    padding: 12px 16px;
    border-bottom: 1px solid var(--admin-border, #e0e0e0);
    transition: background 0.15s;
}

.admin-dashboard__activity-item:last-child {
    border-bottom: none;
}

.admin-dashboard__activity-item:hover {
    background: rgba(108, 92, 231, 0.03);
}

.admin-dashboard__activity-dot {
    width: 8px;
    height: 8px;
    border-radius: 50%;
    margin-top: 6px;
    flex-shrink: 0;
}

.admin-dashboard__activity-dot--gallery {
    background: #6c5ce7;
}
.admin-dashboard__activity-dot--shop {
    background: #00b894;
}
.admin-dashboard__activity-dot--news {
    background: #74b9ff;
}
.admin-dashboard__activity-dot--review {
    background: #fdcb6e;
}
.admin-dashboard__activity-dot--user {
    background: #e17055;
}
.admin-dashboard__activity-dot--purchase {
    background: #00cec9;
}

.admin-dashboard__activity-content {
    display: flex;
    flex-direction: column;
    gap: 2px;
    flex: 1;
    min-width: 0;
}

.admin-dashboard__activity-text {
    font-size: 13px;
    color: var(--admin-text-primary, #2d3436);
    line-height: 1.4;
    overflow: hidden;
    text-overflow: ellipsis;
    display: -webkit-box;
    -webkit-line-clamp: 2;
    -webkit-box-orient: vertical;
}

.admin-dashboard__activity-time {
    font-size: 11px;
    color: var(--admin-text-secondary, #636e72);
}
</style>
