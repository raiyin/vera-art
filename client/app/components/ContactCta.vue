<script setup lang="ts">
    import { useI18n, } from '#imports';
    import type { ButtonProps, DropdownMenuItem, } from '@nuxt/ui';

    const props = withDefaults(defineProps<{
        label?: string
        icon?: string
        size?: ButtonProps['size']
        color?: ButtonProps['color']
        variant?: ButtonProps['variant']
        align?: 'start' | 'center' | 'end'
    }>(), {
        label: '',
        icon: 'i-lucide-mail',
        size: 'lg',
        color: 'primary',
        variant: 'solid',
        align: 'center',
    },);

    const { t, } = useI18n({ useScope: 'global', },);
    const { channels, } = useContacts();

    const items = computed<DropdownMenuItem[]>(() => channels.value.map(channel => ({
        label: channel.label,
        href: channel.href,
        icon: channel.icon,
        external: channel.external,
        target: channel.external ? '_blank' : undefined,
    })));

    const buttonLabel = computed(() => props.label || t('contacts.button',));

    const content = computed(() => {
        return { align: props.align, };
    });
</script>

<template>
    <UDropdownMenu
        v-if="items.length"
        :items="items"
        :content="content"
    >
        <UButton
            :size="size"
            :color="color"
            :variant="variant"
            :icon="icon"
            trailing-icon="i-lucide-chevron-down"
            class="font-semibold"
        >
            {{ buttonLabel }}
        </UButton>
    </UDropdownMenu>
</template>