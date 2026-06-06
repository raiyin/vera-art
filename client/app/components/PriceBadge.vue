<script setup lang="ts">
/**
 * PriceBadge — аккуратный бейдж цены, используемый на карточках
 * мастер-классов и в галерее картин (art-store).
 *
 * Props:
 *   price       — числовое значение цены
 *   isFree      — если true, показывается «Бесплатно» (зелёный)
 *   inKopecks   — если true, цена делится на 100 (режим мастер-классов)
 *   icon        — кастомная иконка (по умолчанию lock-open / lock-closed)
 *   color       — цвет UBadge (по умолчанию success/free, warning/paid)
 */
withDefaults(
    defineProps<{
        price: number;
        isFree?: boolean;
        inKopecks?: boolean;
        icon?: string;
        color?: 'success' | 'warning' | 'primary';
    }>(),
    {
        isFree: false,
        inKopecks: false,
        icon: undefined,
        color: undefined,
    }
);

function formatPrice(price: number, inKopecks: boolean): string {
    const value = inKopecks ? price / 100 : price;
    if (value === 0) return 'Бесплатно';
    return `${value.toLocaleString('ru-RU')} ₽`;
}

function resolveColor(
    isFree: boolean,
    color?: string
): 'success' | 'warning' | 'primary' {
    if (color) return color as 'success' | 'warning' | 'primary';
    return isFree ? 'success' : 'warning';
}

function resolveIcon(isFree: boolean, icon?: string): string {
    if (icon) return icon;
    return isFree ? 'i-heroicons-lock-open' : 'i-heroicons-lock-closed';
}
</script>

<template>
    <UBadge :color="resolveColor(isFree, color)" variant="solid" size="sm">
        <UIcon :name="resolveIcon(isFree, icon)" class="w-3.5 h-3.5 mr-1" />
        {{ formatPrice(price, inKopecks) }}
    </UBadge>
</template>
