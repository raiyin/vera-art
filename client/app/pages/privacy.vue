<script lang="ts">
    export default {
        name: 'PrivacyPage',
        head() {
            return {
                title: this.$t('privacy.title',),
                meta: [
                    {
                        hid: 'description',
                        name: 'description',
                        content: this.$t('privacy.introduction',),
                    },
                ],
            };
        },
    };
</script>

<template>
    <UContainer class="privacy-page py-8 md:py-12">
        <!-- Hero Section -->
        <div class="hero-section mb-10">
            <h1 class="text-4xl md:text-5xl font-bold text-center mb-4">
                {{ $t('privacy.title',) }}
            </h1>
            <p
                class="text-lg text-center text-gray-600 dark:text-gray-300 max-w-3xl mx-auto"
            >
                {{ $t('privacy.lastUpdated',) }}
            </p>
        </div>

        <!-- Content -->
        <div class="max-w-4xl mx-auto">
            <!-- Introduction -->
            <section class="mb-10">
                <p class="text-gray-700 dark:text-gray-300 text-lg leading-relaxed">
                    {{ $t('privacy.introduction',) }}
                </p>
            </section>

            <!-- Information We Collect -->
            <section class="mb-10">
                <h2 class="text-2xl font-bold mb-4 text-gray-900 dark:text-white">
                    {{ $t('privacy.informationWeCollect',) }}
                </h2>
                <p class="text-gray-700 dark:text-gray-300">
                    {{ $t('privacy.informationWeCollectText',) }}
                </p>
            </section>

            <!-- How We Use -->
            <section class="mb-10">
                <h2 class="text-2xl font-bold mb-4 text-gray-900 dark:text-white">
                    {{ $t('privacy.howWeUse',) }}
                </h2>
                <p class="text-gray-700 dark:text-gray-300">
                    {{ $t('privacy.howWeUseText',) }}
                </p>
            </section>

            <!-- Data Protection -->
            <section class="mb-10">
                <h2 class="text-2xl font-bold mb-4 text-gray-900 dark:text-white">
                    {{ $t('privacy.dataProtection',) }}
                </h2>
                <p class="text-gray-700 dark:text-gray-300">
                    {{ $t('privacy.dataProtectionText',) }}
                </p>
            </section>

            <!-- Cookies -->
            <section class="mb-10">
                <h2 class="text-2xl font-bold mb-4 text-gray-900 dark:text-white">
                    {{ $t('privacy.cookies',) }}
                </h2>
                <p class="text-gray-700 dark:text-gray-300">
                    {{ $t('privacy.cookiesText',) }}
                </p>
            </section>

            <!-- Third-Party Services -->
            <section class="mb-10">
                <h2 class="text-2xl font-bold mb-4 text-gray-900 dark:text-white">
                    {{ $t('privacy.thirdParty',) }}
                </h2>
                <p class="text-gray-700 dark:text-gray-300">
                    {{ $t('privacy.thirdPartyText',) }}
                </p>
            </section>

            <!-- Your Rights -->
            <section class="mb-10">
                <h2 class="text-2xl font-bold mb-4 text-gray-900 dark:text-white">
                    {{ $t('privacy.yourRights',) }}
                </h2>
                <p class="text-gray-700 dark:text-gray-300">
                    {{ $t('privacy.yourRightsText',) }}
                </p>
            </section>

            <!-- Data Retention -->
            <section class="mb-10">
                <h2 class="text-2xl font-bold mb-4 text-gray-900 dark:text-white">
                    {{ $t('privacy.dataRetention',) }}
                </h2>
                <p class="text-gray-700 dark:text-gray-300">
                    {{ $t('privacy.dataRetentionText',) }}
                </p>
            </section>

            <!-- International Transfers -->
            <section class="mb-10">
                <h2 class="text-2xl font-bold mb-4 text-gray-900 dark:text-white">
                    {{ $t('privacy.internationalTransfers',) }}
                </h2>
                <p class="text-gray-700 dark:text-gray-300">
                    {{ $t('privacy.internationalTransfersText',) }}
                </p>
            </section>

            <!-- Changes -->
            <section class="mb-10">
                <h2 class="text-2xl font-bold mb-4 text-gray-900 dark:text-white">
                    {{ $t('privacy.changes',) }}
                </h2>
                <p class="text-gray-700 dark:text-gray-300">
                    {{ $t('privacy.changesText',) }}
                </p>
            </section>

            <!-- Contact -->
            <section class="mb-10">
                <h2 class="text-2xl font-bold mb-4 text-gray-900 dark:text-white">
                    {{ $t('privacy.contact',) }}
                </h2>
                <p class="text-gray-700 dark:text-gray-300">
                    {{ $t('privacy.contactText',) }}
                </p>
            </section>

            <!-- Back Link -->
            <div
                class="mt-12 pt-8 border-t border-gray-200 dark:border-gray-700 text-center"
            >
                <NuxtLink
                    to="/auth/register"
                    class="inline-flex items-center text-blue-600 hover:text-blue-700 dark:text-blue-400 dark:hover:text-blue-300 font-medium"
                >
                    <UIcon
                        name="i-heroicons-arrow-left"
                        class="w-5 h-5 mr-2"
                    />
                    {{ $t('auth.backToHome',) }}
                </NuxtLink>
            </div>
        </div>
    </UContainer>
</template>

<style scoped>
.privacy-page {
    min-height: calc(100vh - 200px);
}

.hero-section {
    background: linear-gradient(135deg, #f0f9ff 0%, #e0f2fe 100%);
    padding: 3rem 1.5rem;
    border-radius: 1.5rem;
    margin-bottom: 2rem;
}

.dark .hero-section {
    background: linear-gradient(135deg, #0c4a6e 0%, #082f49 100%);
}

section {
    background: white;
    padding: 2rem;
    border-radius: 1rem;
    box-shadow: 0 4px 6px -1px rgba(0, 0, 0, 0.1), 0 2px 4px -1px rgba(0, 0, 0, 0.06);
}

.dark section {
    background: #1f2937;
    box-shadow: 0 4px 6px -1px rgba(0, 0, 0, 0.3), 0 2px 4px -1px rgba(0, 0, 0, 0.2);
}

h2 {
    color: #1e40af;
}

.dark h2 {
    color: #93c5fd;
}

p {
    line-height: 1.7;
}
</style>
