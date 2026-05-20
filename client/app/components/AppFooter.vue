<script lang="ts">
import { useI18n } from 'vue-i18n';

export default {
    setup() {
        const { t } = useI18n({ useScope: 'global' });
        const currentYear = new Date().getFullYear();
        return {
            currentYear,
            t,
        };
    },
};
</script>

<template>
    <footer class="footer">
        <div class="footer-container">
            <div class="footer-content">
                <div class="footer-text">
                    <span class="copyright"
                        >© {{ currentYear }} {{ t('footer.copyright') }}</span
                    >
                    <span class="separator">•</span>
                    <span class="studio-link">
                        {{ t('footer.developed') }}
                        <a
                            href="https://www.publicmaders.ru"
                            target="_blank"
                            rel="noopener noreferrer"
                            class="studio-link-anchor"
                        >
                            www.publicmaders.ru
                        </a>
                    </span>
                </div>
            </div>
        </div>
    </footer>
</template>

<style scoped>
.footer {
    background-color: var(--color-surface);
    color: var(--color-on-surface);
    border-top: 1px solid var(--color-border);
    padding: 1.5rem 0;
    margin-top: auto;
    flex-shrink: 0;
}

.footer-container {
    max-width: 1200px;
    margin: 0 auto;
    padding: 0 1rem;
}

.footer-content {
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
    text-align: center;
    gap: 0.75rem;
}

.footer-text {
    display: flex;
    flex-wrap: wrap;
    align-items: center;
    justify-content: center;
    gap: 0.75rem;
    font-size: 0.95rem;
    line-height: 1.5;
}

.copyright {
    font-weight: 500;
    color: var(--color-on-surface);
}

.separator {
    color: var(--color-border);
    font-weight: 300;
}

.studio-link {
    color: var(--color-on-surface);
}

.studio-link-anchor {
    color: var(--color-primary);
    text-decoration: none;
    font-weight: 600;
    transition: color 0.2s ease;
    margin-left: 0.25rem;
}

.studio-link-anchor:hover {
    color: var(--color-primary-dark);
    text-decoration: underline;
}

.footer-note {
    font-size: 0.85rem;
    color: var(--color-on-surface-muted);
    max-width: 600px;
    line-height: 1.4;
}

.protection-text {
    font-style: italic;
}

/* Responsive adjustments */
@media (max-width: 640px) {
    .footer-text {
        flex-direction: column;
        gap: 0.5rem;
    }

    .separator {
        display: none;
    }

    .footer {
        padding: 1.25rem 0;
    }
}
</style>
