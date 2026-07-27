<script lang="ts">
    export default {};
</script>

<template>
    <UContainer class="services">
        <!-- Hero Section -->
        <div class="hero-section text-center mb-16">
            <div class="inline-flex items-center justify-center w-16 h-16 bg-gradient-to-br from-teal-400 to-green-500 rounded-2xl shadow-lg mb-6">
                <UIcon
                    name="i-lucide-palette"
                    class="w-8 h-8 text-white"
                />
            </div>
            <h1 class="text-4xl md:text-5xl font-bold text-gray-900 dark:text-white mb-4">
                {{ $t('header.services',) }}
            </h1>
            <p class="text-lg text-gray-600 dark:text-gray-400 max-w-3xl mx-auto">
                {{ $t('services.p_one',) }}
            </p>
        </div>

        <!-- Service Cards -->
        <div class="grid grid-cols-1 md:grid-cols-3 gap-8 mb-16">
            <UCard
                class="service-card hover:shadow-xl hover:-translate-y-1 transition-all duration-300 border-t-4 border-t-teal-500"
            >
                <template #header>
                    <div class="flex items-center justify-center mb-4">
                        <div class="w-16 h-16 rounded-2xl bg-teal-50 dark:bg-teal-900/20 flex items-center justify-center">
                            <UIcon
                                name="i-lucide-palette"
                                class="w-8 h-8 text-teal-500"
                            />
                        </div>
                    </div>
                    <h3 class="text-xl font-semibold text-center text-gray-900 dark:text-white">
                        {{ $t('header.artist',) }}
                    </h3>
                </template>
                <p class="text-gray-600 dark:text-gray-400 text-center leading-relaxed">
                    {{ $t('services.p_two',) }}
                </p>
            </UCard>

            <UCard
                class="service-card hover:shadow-xl hover:-translate-y-1 transition-all duration-300 border-t-4 border-t-teal-500"
            >
                <template #header>
                    <div class="flex items-center justify-center mb-4">
                        <div class="w-16 h-16 rounded-2xl bg-teal-50 dark:bg-teal-900/20 flex items-center justify-center">
                            <UIcon
                                name="i-lucide-book-open"
                                class="w-8 h-8 text-teal-500"
                            />
                        </div>
                    </div>
                    <h3 class="text-xl font-semibold text-center text-gray-900 dark:text-white">
                        {{ $t('header.teacher',) }}
                    </h3>
                </template>
                <p class="text-gray-600 dark:text-gray-400 text-center leading-relaxed">
                    {{ $t('services.p_three',) }}
                </p>
            </UCard>

            <UCard
                class="service-card hover:shadow-xl hover:-translate-y-1 transition-all duration-300 border-t-4 border-t-teal-500"
            >
                <template #header>
                    <div class="flex items-center justify-center mb-4">
                        <div class="w-16 h-16 rounded-2xl bg-teal-50 dark:bg-teal-900/20 flex items-center justify-center">
                            <UIcon
                                name="i-lucide-brush"
                                class="w-8 h-8 text-teal-500"
                            />
                        </div>
                    </div>
                    <h3 class="text-xl font-semibold text-center text-gray-900 dark:text-white">
                        {{ $t('header.illustrator',) }}
                    </h3>
                </template>
                <p class="text-gray-600 dark:text-gray-400 text-center leading-relaxed">
                    {{ $t('services.p_one',) }}
                </p>
            </UCard>
        </div>

        <!-- Materials Section -->
        <div class="materials-section rounded-2xl p-8 md:p-12 mb-16">
            <h2 class="text-2xl md:text-3xl font-bold text-center mb-8 text-gray-900 dark:text-white">
                {{ $t('services.materials_heading',) }}
            </h2>
            <div class="grid grid-cols-2 md:grid-cols-4 gap-4">
                <div class="flex items-center justify-center p-5 bg-white dark:bg-gray-800 rounded-xl shadow-sm hover:shadow-md transition-shadow border border-gray-100 dark:border-gray-700">
                    <UIcon
                        name="i-lucide-palette"
                        class="w-6 h-6 text-teal-500 mr-3 shrink-0"
                    />
                    <span class="font-medium text-gray-800 dark:text-gray-200">{{ $t('services.watercolor',) }}</span>
                </div>
                <div class="flex items-center justify-center p-5 bg-white dark:bg-gray-800 rounded-xl shadow-sm hover:shadow-md transition-shadow border border-gray-100 dark:border-gray-700">
                    <UIcon
                        name="i-lucide-droplets"
                        class="w-6 h-6 text-teal-500 mr-3 shrink-0"
                    />
                    <span class="font-medium text-gray-800 dark:text-gray-200">{{ $t('services.oil',) }}</span>
                </div>
                <div class="flex items-center justify-center p-5 bg-white dark:bg-gray-800 rounded-xl shadow-sm hover:shadow-md transition-shadow border border-gray-100 dark:border-gray-700">
                    <UIcon
                        name="i-lucide-pencil"
                        class="w-6 h-6 text-teal-500 mr-3 shrink-0"
                    />
                    <span class="font-medium text-gray-800 dark:text-gray-200">{{ $t('services.pencil_charcoal',) }}</span>
                </div>
                <div class="flex items-center justify-center p-5 bg-white dark:bg-gray-800 rounded-xl shadow-sm hover:shadow-md transition-shadow border border-gray-100 dark:border-gray-700">
                    <UIcon
                        name="i-lucide-layers"
                        class="w-6 h-6 text-teal-500 mr-3 shrink-0"
                    />
                    <span class="font-medium text-gray-800 dark:text-gray-200">{{ $t('services.pastel_gouache',) }}</span>
                </div>
            </div>
        </div>

        <!-- CTA Section -->
        <div class="cta-section text-center mb-16 p-12 rounded-3xl bg-gradient-to-br from-teal-50 to-green-50 dark:from-gray-800 dark:to-gray-900 border border-teal-100 dark:border-teal-900/30">
            <div class="inline-flex items-center justify-center w-16 h-16 rounded-2xl bg-teal-500 shadow-lg mb-6">
                <UIcon
                    name="i-lucide-mail"
                    class="w-8 h-8 text-white"
                />
            </div>
            <h3 class="text-2xl md:text-3xl font-bold text-gray-900 dark:text-white mb-4">
                {{ $t('services.collaborate_heading',) }}
            </h3>
            <p class="text-gray-600 dark:text-gray-400 mb-8 max-w-2xl mx-auto text-lg">
                {{ $t('services.collaborate_text',) }}
            </p>
            <UButton
                size="lg"
                color="primary"
                variant="solid"
                icon="i-lucide-mail"
                :to="{ path: '/', hash: '#contact', }"
                class="font-semibold"
            >
                {{ $t('services.contact_button',) }}
            </UButton>
        </div>
    </UContainer>
</template>

<style scoped>
.services {
    margin-top: 3rem;
    flex: 1 0 auto;
    padding-bottom: 4rem;
}

.hero-section {
    margin-bottom: 4rem;
}

.service-card {
    height: 100%;
    display: flex;
    flex-direction: column;
}

.service-card :deep(.u-card-body) {
    flex-grow: 1;
}

.materials-section {
    background: linear-gradient(135deg, #f0fdfa 0%, #ecfdf5 100%);
    transition: background-color 0.3s ease;
}

.body_theme_dark .materials-section {
    background: linear-gradient(135deg, #111 0%, #1a1a1a 100%);
}

.body_theme_dark .materials-section > div > div {
    background: var(--color-surface-secondary-solid);
    border-color: var(--color-border);
}

@media (max-width: 768px) {
    .services {
        margin-top: 2rem;
    }
    .hero-section h1 {
        font-size: 2.5rem;
    }
}
</style>
