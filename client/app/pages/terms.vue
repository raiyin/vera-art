<script lang="ts">
export default {
    name: 'TermsPage',
    head() {
        return {
            title: this.$t('terms.title'),
            meta: [
                {
                    hid: 'description',
                    name: 'description',
                    content: this.$t('terms.introduction'),
                },
            ],
        };
    },
};
</script>

<template>
    <UContainer class="terms-page py-8 md:py-12">
        <!-- Hero Section -->
        <div class="hero-section mb-10">
            <h1 class="text-4xl md:text-5xl font-bold text-center mb-4">
                {{ $t('terms.title') }}
            </h1>
            <p
                class="text-lg text-center text-gray-600 dark:text-gray-300 max-w-3xl mx-auto"
            >
                {{ $t('terms.lastUpdated') }}
            </p>
        </div>

        <!-- Content -->
        <div class="max-w-4xl mx-auto">
            <!-- Introduction -->
            <section class="mb-10">
                <p class="text-gray-700 dark:text-gray-300 text-lg leading-relaxed">
                    {{ $t('terms.introduction') }}
                </p>
            </section>

            <!-- Acceptance of Terms -->
            <section class="mb-10">
                <h2 class="text-2xl font-bold mb-4 text-gray-900 dark:text-white">
                    {{ $t('terms.acceptance') }}
                </h2>
                <p class="text-gray-700 dark:text-gray-300">
                    {{ $t('terms.acceptanceText') }}
                </p>
            </section>

            <!-- Intellectual Property -->
            <section class="mb-10">
                <h2 class="text-2xl font-bold mb-4 text-gray-900 dark:text-white">
                    {{ $t('terms.intellectualProperty') }}
                </h2>
                <p class="text-gray-700 dark:text-gray-300">
                    {{ $t('terms.intellectualPropertyText') }}
                </p>
            </section>

            <!-- Use License -->
            <section class="mb-10">
                <h2 class="text-2xl font-bold mb-4 text-gray-900 dark:text-white">
                    {{ $t('terms.useLicense') }}
                </h2>
                <p class="text-gray-700 dark:text-gray-300">
                    {{ $t('terms.useLicenseText') }}
                </p>
            </section>

            <!-- Restrictions -->
            <section class="mb-10">
                <h2 class="text-2xl font-bold mb-4 text-gray-900 dark:text-white">
                    {{ $t('terms.restrictions') }}
                </h2>
                <p class="text-gray-700 dark:text-gray-300">
                    {{ $t('terms.restrictionsList') }}
                </p>
            </section>

            <!-- User Content -->
            <section class="mb-10">
                <h2 class="text-2xl font-bold mb-4 text-gray-900 dark:text-white">
                    {{ $t('terms.userContent') }}
                </h2>
                <p class="text-gray-700 dark:text-gray-300">
                    {{ $t('terms.userContentText') }}
                </p>
            </section>

            <!-- Limitations of Liability -->
            <section class="mb-10">
                <h2 class="text-2xl font-bold mb-4 text-gray-900 dark:text-white">
                    {{ $t('terms.limitations') }}
                </h2>
                <p class="text-gray-700 dark:text-gray-300">
                    {{ $t('terms.limitationsText') }}
                </p>
            </section>

            <!-- Indemnification -->
            <section class="mb-10">
                <h2 class="text-2xl font-bold mb-4 text-gray-900 dark:text-white">
                    {{ $t('terms.indemnification') }}
                </h2>
                <p class="text-gray-700 dark:text-gray-300">
                    {{ $t('terms.indemnificationText') }}
                </p>
            </section>

            <!-- Changes to Terms -->
            <section class="mb-10">
                <h2 class="text-2xl font-bold mb-4 text-gray-900 dark:text-white">
                    {{ $t('terms.changes') }}
                </h2>
                <p class="text-gray-700 dark:text-gray-300">
                    {{ $t('terms.changesText') }}
                </p>
            </section>

            <!-- Governing Law -->
            <section class="mb-10">
                <h2 class="text-2xl font-bold mb-4 text-gray-900 dark:text-white">
                    {{ $t('terms.governingLaw') }}
                </h2>
                <p class="text-gray-700 dark:text-gray-300">
                    {{ $t('terms.governingLawText') }}
                </p>
            </section>

            <!-- Contact -->
            <section class="mb-10">
                <h2 class="text-2xl font-bold mb-4 text-gray-900 dark:text-white">
                    {{ $t('terms.contact') }}
                </h2>
                <p class="text-gray-700 dark:text-gray-300">
                    {{ $t('terms.contactText') }}
                </p>
            </section>

            <!-- Back Link -->
            <div
                class="mt-12 pt-8 border-t border-gray-200 dark:border-gray-700 text-center"
            >
                <NuxtLink
                    to="/register"
                    class="inline-flex items-center text-blue-600 hover:text-blue-700 dark:text-blue-400 dark:hover:text-blue-300 font-medium"
                >
                    <UIcon name="i-heroicons-arrow-left" class="w-5 h-5 mr-2" />
                    {{ $t('auth.backToHome') }}
                </NuxtLink>
            </div>
        </div>
    </UContainer>
</template>

<style scoped>
.terms-page {
    min-height: calc(100vh - 200px);
}

.hero-section {
    background: linear-gradient(135deg, #f3f4f6 0%, #e5e7eb 100%);
    padding: 3rem 1.5rem;
    border-radius: 1.5rem;
    margin-bottom: 2rem;
}

.dark .hero-section {
    background: linear-gradient(135deg, #1f2937 0%, #111827 100%);
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
    color: #1f2937;
}

.dark h2 {
    color: #f9fafb;
}

p {
    line-height: 1.7;
}
</style>
