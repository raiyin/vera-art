<script lang="ts">
    import { useI18n, } from 'vue-i18n';
    import { computed, } from 'vue';

    export default {
        setup() {
            const { t, locale, } = useI18n({ useScope: 'global', },);
            const currentYear = new Date().getFullYear();

            const navLinks = computed(() => [
                { label: t('header.main',), to: '/', },
                { label: t('header.all_works',), to: '/gallery', },
                { label: t('header.news',), to: '/news', },
                { label: t('footer.paintings',), to: '/art-store', },
                { label: t('header.payment',), to: '/pay-delivery', },
            ]);

            const serviceLinks = computed(() => [
                { label: t('header.teaching',), to: '/teaching', },
                { label: t('footer.master_classes',), to: '/master-classes', },
                { label: t('footer.online_courses',), to: '/courses/online', },
                { label: t('footer.individual_lessons',), to: '/courses/individual', },
            ]);

            const { channels, } = useContacts();

            return {
                currentYear,
                t,
                locale,
                navLinks,
                serviceLinks,
                channels,
            };
        },
    };
</script>

<template>
    <footer class="footer">
        <div class="footer-container">
            <div class="footer-grid">
                <!-- Brand -->
                <div class="footer-brand">
                    <NuxtLink
                        to="/"
                        class="footer-logo"
                    >
                        <svg
                            class="footer-logo-icon"
                            viewBox="0 0 24 24"
                            fill="none"
                            stroke="currentColor"
                            stroke-width="2"
                        >
                            <circle
                                cx="12"
                                cy="12"
                                r="10"
                            />
                            <path d="M12 8v4l3 3" />
                        </svg>
                        <span class="footer-logo-text">Vera Pertsukova</span>
                    </NuxtLink>
                    <p class="footer-description">
                        {{ t('footer.copyright',) }}
                    </p>
                </div>

                <!-- Navigation -->
                <div class="footer-section">
                    <h4 class="footer-section-title">
                        {{ t('footer.navigation',) }}
                    </h4>
                    <ul class="footer-links">
                        <li
                            v-for="link in navLinks"
                            :key="link.to"
                        >
                            <NuxtLink
                                :to="link.to"
                                class="footer-link"
                            >
                                {{ link.label }}
                            </NuxtLink>
                        </li>
                    </ul>
                </div>

                <!-- Services -->
                <div class="footer-section">
                    <h4 class="footer-section-title">
                        {{ t('footer.services',) }}
                    </h4>
                    <ul class="footer-links">
                        <li
                            v-for="link in serviceLinks"
                            :key="link.to"
                        >
                            <NuxtLink
                                :to="link.to"
                                class="footer-link"
                            >
                                {{ link.label }}
                            </NuxtLink>
                        </li>
                    </ul>
                </div>

                <!-- Social -->
                <div
                    id="contacts"
                    class="footer-section scroll-mt-24"
                >
                    <h4 class="footer-section-title">
                        {{ t('footer.contacts',) }}
                    </h4>
                    <div class="footer-social">
                        <a
                            v-for="channel in channels"
                            :key="channel.id"
                            :href="channel.href"
                            :target="channel.external ? '_blank' : undefined"
                            :rel="channel.external ? 'noopener noreferrer' : undefined"
                            class="footer-social-link"
                            :aria-label="channel.label"
                        >
                            <UIcon
                                :name="channel.icon"
                                class="footer-social-icon"
                            />
                        </a>
                    </div>
                </div>
            </div>

            <div class="footer-bottom">
                <p class="footer-copyright">
                    © {{ currentYear }} {{ t('footer.copyright',) }}.
                    {{ t('footer.developed',) }}
                    <a
                        href="https://www.publicmaders.ru"
                        target="_blank"
                        rel="noopener noreferrer"
                        class="footer-studio-link"
                    >
                        www.publicmaders.ru
                    </a>
                </p>
            </div>
        </div>
    </footer>
</template>

<style scoped>
.footer {
    background-color: var(--color-surface);
    color: var(--color-on-surface);
    border-top: 1px solid var(--color-border);
    padding: 3rem 0 0;
    margin-top: auto;
    flex-shrink: 0;
}

.footer-container {
    max-width: 1200px;
    margin: 0 auto;
    padding: 0 1.5rem;
}

.footer-grid {
    display: grid;
    grid-template-columns: 2fr 1fr 1fr 1fr;
    gap: 2.5rem;
    padding-bottom: 2.5rem;
}

.footer-brand {
    max-width: 300px;
}

.footer-logo {
    display: inline-flex;
    align-items: center;
    gap: 0.5rem;
    text-decoration: none;
    margin-bottom: 1rem;
}

.footer-logo-icon {
    width: 32px;
    height: 32px;
    color: var(--color-primary, #4B9E90);
}

.footer-logo-text {
    font-size: 1.1rem;
    font-weight: 700;
    color: var(--color-on-surface);
}

.footer-description {
    font-size: 0.9rem;
    color: var(--color-on-surface);
    opacity: 0.7;
    line-height: 1.5;
}

.footer-section-title {
    font-size: 0.85rem;
    font-weight: 700;
    text-transform: uppercase;
    letter-spacing: 0.05rem;
    color: var(--color-on-surface);
    margin-bottom: 1rem;
    opacity: 0.8;
}

.footer-links {
    list-style: none;
    padding: 0;
    margin: 0;
    display: flex;
    flex-direction: column;
    gap: 0.6rem;
}

.footer-link {
    color: var(--color-on-surface);
    text-decoration: none;
    font-size: 0.9rem;
    opacity: 0.7;
    transition: all 0.2s ease;
}

.footer-link:hover {
    opacity: 1;
    color: var(--color-primary, #4B9E90);
}

.footer-social {
    display: flex;
    gap: 0.75rem;
}

.footer-social-link {
    display: flex;
    align-items: center;
    justify-content: center;
    width: 40px;
    height: 40px;
    border-radius: 50%;
    background: var(--color-surface-secondary);
    border: 1px solid var(--color-border);
    color: var(--color-on-surface);
    transition: all 0.2s ease;
}

.footer-social-link:hover {
    background: var(--color-primary, #4B9E90);
    border-color: var(--color-primary, #4B9E90);
    color: white;
    transform: translateY(-2px);
}

.footer-social-icon {
    width: 18px;
    height: 18px;
}

.footer-bottom {
    border-top: 1px solid var(--color-border);
    padding: 1.5rem 0;
    text-align: center;
}

.footer-copyright {
    font-size: 0.85rem;
    color: var(--color-on-surface);
    opacity: 0.6;
}

.footer-studio-link {
    color: var(--color-primary, #4B9E90);
    text-decoration: none;
    font-weight: 600;
    transition: color 0.2s ease;
}

.footer-studio-link:hover {
    text-decoration: underline;
}

@media (max-width: 900px) {
    .footer-grid {
        grid-template-columns: 1fr 1fr;
        gap: 2rem;
    }

    .footer-brand {
        grid-column: 1 / -1;
        max-width: none;
    }
}

@media (max-width: 600px) {
    .footer-grid {
        grid-template-columns: 1fr;
        gap: 1.5rem;
    }

    .footer {
        padding: 2rem 0 0;
    }
}
</style>
