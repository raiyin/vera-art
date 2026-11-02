<script setup lang="ts">
    interface NavItem {
        to: string
        label: string
        desc: string
        icon: string
        color: string
    }

    interface NavSection {
        label: string
        items: NavItem[]
    }

    const navSections: NavSection[] = [
        {
            label: 'Контент',
            items: [
                {
                    to: '/admin/gallery',
                    label: 'Галерея',
                    desc: 'Управление работами в галерее',
                    icon: 'i-lucide-image',
                    color: 'purple',
                },
                {
                    to: '/admin/shop',
                    label: 'Магазин',
                    desc: 'Управление товарами в магазине',
                    icon: 'i-lucide-shopping-bag',
                    color: 'green',
                },
                {
                    to: '/admin/news',
                    label: 'Новости',
                    desc: 'Управление новостями',
                    icon: 'i-lucide-newspaper',
                    color: 'blue',
                },
                {
                    to: '/admin/courses',
                    label: 'Курсы',
                    desc: 'Управление курсами',
                    icon: 'i-lucide-graduation-cap',
                    color: 'pink',
                },
                {
                    to: '/admin/master-classes',
                    label: 'Мастер-классы',
                    desc: 'Управление мастер-классами',
                    icon: 'i-lucide-video',
                    color: 'orange',
                },
                {
                    to: '/admin/lessons',
                    label: 'Уроки',
                    desc: 'Управление уроками курсов и МК',
                    icon: 'i-lucide-book-open',
                    color: 'teal',
                },
            ],
        },
        {
            label: 'Пользователи',
            items: [
                {
                    to: '/admin/users',
                    label: 'Пользователи',
                    desc: 'Управление пользователями',
                    icon: 'i-lucide-users',
                    color: 'orange',
                },
                {
                    to: '/admin/reviews',
                    label: 'Отзывы',
                    desc: 'Модерация отзывов пользователей',
                    icon: 'i-lucide-star',
                    color: 'red',
                },
            ],
        },
        {
            label: 'Финансы',
            items: [
                {
                    to: '/admin/purchases',
                    label: 'Покупки',
                    desc: 'Управление покупками',
                    icon: 'i-lucide-shopping-cart',
                    color: 'teal',
                },
                {
                    to: '/admin/payments',
                    label: 'Платежи',
                    desc: 'Управление платежами и возвратами',
                    icon: 'i-lucide-credit-card',
                    color: 'yellow',
                },
                {
                    to: '/admin/promo-codes',
                    label: 'Промокоды',
                    desc: 'Управление промокодами и скидками',
                    icon: 'i-lucide-ticket-percent',
                    color: 'green',
                },
            ],
        },
        {
            label: 'Коммуникация',
            items: [
                {
                    to: '/admin/chats',
                    label: 'Чаты',
                    desc: 'Управление обращениями пользователей',
                    icon: 'i-lucide-message-square',
                    color: 'blue',
                },
            ],
        },
        {
            label: 'Справочники',
            items: [
                {
                    to: '/admin/categories',
                    label: 'Категории',
                    desc: 'Управление категориями продуктов',
                    icon: 'i-lucide-folder-tree',
                    color: 'purple',
                },
                {
                    to: '/admin/tags',
                    label: 'Теги',
                    desc: 'Управление тегами продуктов',
                    icon: 'i-lucide-tags',
                    color: 'pink',
                },
                {
                    to: '/admin/settings',
                    label: 'Настройки',
                    desc: 'Управление справочниками',
                    icon: 'i-lucide-settings',
                    color: 'neutral',
                },
            ],
        },
    ];
</script>

<template>
    <div class="admin-dashboard__nav-section">
        <h3 class="admin-dashboard__nav-section-title">
            Навигация по разделам
        </h3>
        <div
            v-for="(section, sIdx) in navSections"
            :key="sIdx"
            class="admin-dashboard__nav-group"
        >
            <span class="admin-dashboard__nav-group-label">{{ section.label }}</span>
            <div class="admin-dashboard__nav-grid">
                <NuxtLink
                    v-for="item in section.items"
                    :key="item.to"
                    :to="item.to"
                    class="admin-dashboard__nav-tile"
                >
                    <div
                        class="admin-dashboard__nav-tile-icon"
                        :class="`admin-dashboard__nav-tile-icon--${item.color}`"
                    >
                        <UIcon
                            :name="item.icon"
                            class="size-5"
                        />
                    </div>
                    <div class="admin-dashboard__nav-tile-info">
                        <span class="admin-dashboard__nav-tile-label">{{
                            item.label
                        }}</span>
                        <span class="admin-dashboard__nav-tile-desc">{{
                            item.desc
                        }}</span>
                    </div>
                    <UIcon
                        name="i-lucide-chevron-right"
                        class="admin-dashboard__nav-tile-arrow size-4"
                    />
                </NuxtLink>
            </div>
        </div>
    </div>
</template>

<style scoped>
.admin-dashboard__nav-section {
    margin-bottom: 24px;
}

.admin-dashboard__nav-section-title {
    font-size: 18px;
    font-weight: 700;
    color: var(--admin-text-primary, #2d3436);
    margin: 0 0 16px;
}

.admin-dashboard__nav-group {
    margin-bottom: 20px;
}

.admin-dashboard__nav-group:last-child {
    margin-bottom: 0;
}

.admin-dashboard__nav-group-label {
    display: block;
    font-size: 11px;
    font-weight: 700;
    text-transform: uppercase;
    letter-spacing: 1px;
    color: var(--admin-text-secondary, #636e72);
    margin-bottom: 8px;
    padding-left: 4px;
}

.admin-dashboard__nav-grid {
    display: grid;
    grid-template-columns: repeat(auto-fill, minmax(280px, 1fr));
    gap: 8px;
}

.admin-dashboard__nav-tile {
    display: flex;
    align-items: center;
    gap: 12px;
    padding: 12px 16px;
    background: var(--admin-surface, #ffffff);
    border: 1px solid var(--admin-border, #e0e0e0);
    border-radius: 10px;
    text-decoration: none;
    transition: all 0.2s ease;
    position: relative;
}

.admin-dashboard__nav-tile:hover {
    border-color: var(--admin-primary, #6c5ce7);
    box-shadow: 0 2px 8px rgba(108, 92, 231, 0.1);
    transform: translateY(-1px);
}

.admin-dashboard__nav-tile:active {
    transform: translateY(0);
}

.admin-dashboard__nav-tile-icon {
    width: 40px;
    height: 40px;
    border-radius: 10px;
    display: flex;
    align-items: center;
    justify-content: center;
    flex-shrink: 0;
}

.admin-dashboard__nav-tile-icon--purple {
    background: rgba(108, 92, 231, 0.1);
    color: #6c5ce7;
}

.admin-dashboard__nav-tile-icon--green {
    background: rgba(0, 184, 148, 0.1);
    color: #00b894;
}

.admin-dashboard__nav-tile-icon--blue {
    background: rgba(116, 185, 255, 0.1);
    color: #74b9ff;
}

.admin-dashboard__nav-tile-icon--orange {
    background: rgba(253, 203, 110, 0.15);
    color: #e17055;
}

.admin-dashboard__nav-tile-icon--pink {
    background: rgba(232, 67, 147, 0.1);
    color: #e84393;
}

.admin-dashboard__nav-tile-icon--red {
    background: rgba(225, 112, 85, 0.1);
    color: #e17055;
}

.admin-dashboard__nav-tile-icon--teal {
    background: rgba(0, 206, 201, 0.1);
    color: #00cec9;
}

.admin-dashboard__nav-tile-icon--yellow {
    background: rgba(253, 203, 110, 0.15);
    color: #fdcb6e;
}

.admin-dashboard__nav-tile-icon--neutral {
    background: rgba(99, 110, 114, 0.1);
    color: #636e72;
}

.admin-dashboard__nav-tile-info {
    flex: 1;
    min-width: 0;
    display: flex;
    flex-direction: column;
    gap: 2px;
}

.admin-dashboard__nav-tile-label {
    font-size: 14px;
    font-weight: 600;
    color: var(--admin-text-primary, #2d3436);
    line-height: 1.3;
}

.admin-dashboard__nav-tile-desc {
    font-size: 12px;
    color: var(--admin-text-secondary, #636e72);
    line-height: 1.3;
    overflow: hidden;
    text-overflow: ellipsis;
    white-space: nowrap;
}

.admin-dashboard__nav-tile-arrow {
    color: var(--admin-text-secondary, #636e72);
    flex-shrink: 0;
    opacity: 0;
    transition: opacity 0.2s ease, transform 0.2s ease;
}

.admin-dashboard__nav-tile:hover .admin-dashboard__nav-tile-arrow {
    opacity: 1;
    transform: translateX(2px);
}

@media (max-width: 640px) {
    .admin-dashboard__nav-grid {
        grid-template-columns: 1fr;
    }
}
</style>
