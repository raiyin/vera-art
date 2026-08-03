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
                { label: t('header.services',), to: '/services', },
                { label: t('header.payment',), to: '/pay-delivery', },
            ]);

            const serviceLinks = computed(() => [
                { label: t('footer.paintings',), to: '/art-store', },
                { label: t('footer.master_classes',), to: '/master-classes', },
                { label: t('footer.online_courses',), to: '/courses/online', },
                { label: t('footer.individual_lessons',), to: '/courses/individual', },
            ]);

            return {
                currentYear,
                t,
                locale,
                navLinks,
                serviceLinks,
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
                <div class="footer-section">
                    <h4 class="footer-section-title">
                        {{ t('footer.contacts',) }}
                    </h4>
                    <div class="footer-social">
                        <a
                            href="https://t.me/MilayaV"
                            target="_blank"
                            rel="noopener noreferrer"
                            class="footer-social-link"
                            aria-label="Telegram"
                        >
                            <svg
                                class="footer-social-icon"
                                viewBox="0 0 24 24"
                                fill="currentColor"
                            >
                                <path d="M11.944 0A12 12 0 000 12a12 12 0 0012 12 12 12 0 0012-12A12 12 0 0012 0a12 12 0 00-.056 0zm4.962 7.224c.1-.002.321.023.465.14a.506.506 0 01.171.325c.016.093.036.306.02.472-.18 1.898-.962 6.502-1.36 8.627-.168.9-.499 1.201-.82 1.23-.696.065-1.225-.46-1.9-.902-1.056-.693-1.653-1.124-2.678-1.8-1.185-.78-.417-1.21.258-1.91.177-.184 3.247-2.977 3.307-3.23.007-.032.014-.15-.056-.212s-.174-.041-.249-.024c-.106.024-1.793 1.14-5.061 3.345-.48.33-.913.49-1.302.48-.428-.008-1.252-.241-1.865-.44-.752-.245-1.349-.374-1.297-.789.027-.216.325-.437.893-.663 3.498-1.524 5.83-2.529 6.998-3.014 3.332-1.386 4.025-1.627 4.476-1.635z" />
                            </svg>
                        </a>
                        <a
                            href="https://vk.com/perczukowa"
                            target="_blank"
                            rel="noopener noreferrer"
                            class="footer-social-link"
                            aria-label="VK"
                        >
                            <svg
                                class="footer-social-icon"
                                viewBox="0 0 24 24"
                                fill="currentColor"
                            >
                                <path d="M15.684 0H8.316C2.755 0 0 2.755 0 8.316v7.368C0 21.245 2.755 24 8.316 24h7.368C21.245 24 24 21.245 24 15.684V8.316C24 2.755 21.245 0 15.684 0zm3.534 16.635h-1.447c-.788 0-1.039-.503-2.466-1.933-1.246-1.247-1.797-1.403-2.099-1.403-.428 0-.55.126-.55.607v1.593c0 .286-.13.454-1.216.454-1.791 0-3.774-1.095-5.167-3.14C4.045 10.127 3.418 8.15 3.418 7.51c0-.25.113-.478.606-.478h1.454c.454 0 .624.204.795.694.638 2.1 1.714 3.964 2.155 3.964.168 0 .242-.084.242-.547v-2.12c-.056-1.012-.588-1.096-.588-1.46 0-.168.126-.336.336-.336h2.32c.284 0 .378.168.378.504v2.727c0 .285.126.378.21.378.168 0 .294-.093.588-.378.924-1.033 1.59-2.634 1.59-2.634.084-.21.252-.378.504-.378h1.454c.168 0 .294.084.336.252.084.378-.378 2.155-1.77 3.745-1.413 1.582-1.413 1.624-.168 2.666 1.202 1.034 1.476 1.54 1.476 1.994.008.378-.21.546-.588.546z" />
                            </svg>
                        </a>
                        <a
                            href="mailto:perczukowa@yandex.ru"
                            class="footer-social-link"
                            aria-label="Email"
                        >
                            <svg
                                class="footer-social-icon"
                                viewBox="0 0 24 24"
                                fill="none"
                                stroke="currentColor"
                                stroke-width="2"
                            >
                                <path d="M4 4h16c1.1 0 2 .9 2 2v12c0 1.1-.9 2-2 2H4c-1.1 0-2-.9-2-2V6c0-1.1.9-2 2-2z" />
                                <polyline points="22,6 12,13 2,6" />
                            </svg>
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
