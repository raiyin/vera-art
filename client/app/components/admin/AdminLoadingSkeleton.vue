<template>
    <div
        class="admin-skeleton"
        :class="{ 'admin-skeleton--dark': isDark, }"
    >
        <table class="admin-skeleton__table">
            <thead>
                <tr>
                    <th
                        v-for="col in columns"
                        :key="col"
                    >
                        <div class="admin-skeleton__cell admin-skeleton__cell--head" />
                    </th>
                </tr>
            </thead>
            <tbody>
                <tr
                    v-for="row in rows"
                    :key="row"
                >
                    <td
                        v-for="col in columns"
                        :key="col"
                    >
                        <div
                            class="admin-skeleton__cell"
                            :class="{
                                'admin-skeleton__cell--image': col === 0 && row % 2 === 0,
                                'admin-skeleton__cell--short': col === columns - 1,
                                'admin-skeleton__cell--medium': col === 1,
                            }"
                        />
                    </td>
                </tr>
            </tbody>
        </table>
    </div>
</template>

<script setup lang="ts">
    withDefaults(
        defineProps<{
            rows?: number
            columns?: number
            isDark?: boolean
        }>(),
        {
            rows: 5,
            columns: 6,
            isDark: false,
        }
    );
</script>

<style scoped>
.admin-skeleton {
    width: 100%;
}

.admin-skeleton__table {
    width: 100%;
    border-collapse: collapse;
}

.admin-skeleton__cell {
    height: 20px;
    background: linear-gradient(
        90deg,
        var(--admin-border, #e0e0e0) 25%,
        rgba(108, 92, 231, 0.08) 50%,
        var(--admin-border, #e0e0e0) 75%
    );
    background-size: 200% 100%;
    animation: admin-skeleton-shimmer 1.5s ease-in-out infinite;
    border-radius: 6px;
    margin: 12px 16px;
}

.admin-skeleton--dark .admin-skeleton__cell {
    background: linear-gradient(
        90deg,
        rgba(255, 255, 255, 0.06) 25%,
        rgba(108, 92, 231, 0.12) 50%,
        rgba(255, 255, 255, 0.06) 75%
    );
    background-size: 200% 100%;
}

.admin-skeleton__cell--head {
    height: 14px;
    margin: 14px 16px;
}

.admin-skeleton__cell--image {
    width: 48px;
    height: 48px;
    border-radius: 8px;
    margin: 8px 16px;
}

.admin-skeleton__cell--short {
    width: 60%;
}

.admin-skeleton__cell--medium {
    width: 75%;
}

@keyframes admin-skeleton-shimmer {
    0% {
        background-position: 200% 0;
    }
    100% {
        background-position: -200% 0;
    }
}
</style>
