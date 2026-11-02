import { useI18n, } from '#imports';

export type ContactChannelId = 'telegram' | 'vk' | 'email';

export interface ContactChannel {
    id: ContactChannelId;
    label: string;
    href: string;
    icon: string;
    external: boolean;
}

/**
 * Composable providing the contact channels of the artist.
 *
 * Single source of truth for the contact links used in the header,
 * the footer and the "contact me" CTA blocks: URLs come from
 * `runtimeConfig.public.contacts`, labels are localized.
 */
export function useContacts() {
    const { t, } = useI18n({ useScope: 'global', },);
    const config = useRuntimeConfig();

    const channels = computed<ContactChannel[]>(() => {
        const { telegram, vk, email, } = config.public.contacts;

        return [
            {
                id: 'telegram',
                label: t('contacts.telegram',),
                href: telegram,
                icon: 'i-simple-icons-telegram',
                external: true,
            },
            {
                id: 'vk',
                label: t('contacts.vk',),
                href: vk,
                icon: 'i-simple-icons-vk',
                external: true,
            },
            {
                id: 'email',
                label: t('contacts.email',),
                href: email ? `mailto:${email}` : '',
                icon: 'i-lucide-mail',
                external: false,
            },
        ].filter(channel => !!channel.href);
    });

    return {
        channels,
    };
}