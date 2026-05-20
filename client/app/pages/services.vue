<script lang="ts">
export default {};
</script>

<template>
    <UContainer class="services">
        <div class="hero-section">
            <h1 class="text-4xl md:text-5xl font-bold text-center mb-6">
                {{ $t('header.services') }}
            </h1>
            <p
                class="text-lg text-center text-gray-600 dark:text-gray-300 max-w-3xl mx-auto mb-12"
            >
                {{ $t('services.p_one') }}
            </p>
        </div>

        <div class="grid grid-cols-1 md:grid-cols-3 gap-8 mb-12">
            <UCard class="service-card hover:shadow-xl transition-shadow duration-300">
                <template #header>
                    <div class="flex items-center justify-center mb-4">
                        <UIcon name="i-lucide-palette" class="w-12 h-12 text-green-600" />
                    </div>
                    <h3 class="text-xl font-semibold text-center">
                        {{ $t('header.artist') }}
                    </h3>
                </template>
                <p class="text-gray-700 dark:text-gray-300 text-justify">
                    {{ $t('services.p_two') }}
                </p>
            </UCard>

            <UCard class="service-card hover:shadow-xl transition-shadow duration-300">
                <template #header>
                    <div class="flex items-center justify-center mb-4">
                        <UIcon
                            name="i-lucide-book-open"
                            class="w-12 h-12 text-green-600"
                        />
                    </div>
                    <h3 class="text-xl font-semibold text-center">
                        {{ $t('header.teacher') }}
                    </h3>
                </template>
                <p class="text-gray-700 dark:text-gray-300 text-justify">
                    {{ $t('services.p_three') }}
                </p>
            </UCard>

            <UCard class="service-card hover:shadow-xl transition-shadow duration-300">
                <template #header>
                    <div class="flex items-center justify-center mb-4">
                        <UIcon name="i-lucide-brush" class="w-12 h-12 text-green-600" />
                    </div>
                    <h3 class="text-xl font-semibold text-center">
                        {{ $t('header.illustrator') }}
                    </h3>
                </template>
                <p class="text-gray-700 dark:text-gray-300 text-justify">
                    {{ $t('services.p_one') }}
                </p>
            </UCard>
        </div>

        <div
            class="materials-section bg-green-50 dark:bg-green-900/20 rounded-2xl p-8 mb-12"
        >
            <h2 class="text-2xl font-bold text-center mb-6">
                {{ $t('services.materials_heading') }}
            </h2>
            <div class="grid grid-cols-2 md:grid-cols-4 gap-4">
                <div
                    class="flex items-center justify-center p-4 bg-white dark:bg-gray-800 rounded-lg shadow"
                >
                    <UIcon name="i-lucide-palette" class="w-8 h-8 text-green-600 mr-3" />
                    <span class="font-medium">{{ $t('services.watercolor') }}</span>
                </div>
                <div
                    class="flex items-center justify-center p-4 bg-white dark:bg-gray-800 rounded-lg shadow"
                >
                    <UIcon name="i-lucide-droplets" class="w-8 h-8 text-green-600 mr-3" />
                    <span class="font-medium">{{ $t('services.oil') }}</span>
                </div>
                <div
                    class="flex items-center justify-center p-4 bg-white dark:bg-gray-800 rounded-lg shadow"
                >
                    <UIcon name="i-lucide-pencil" class="w-8 h-8 text-green-600 mr-3" />
                    <span class="font-medium">{{ $t('services.pencil_charcoal') }}</span>
                </div>
                <div
                    class="flex items-center justify-center p-4 bg-white dark:bg-gray-800 rounded-lg shadow"
                >
                    <UIcon name="i-lucide-layers" class="w-8 h-8 text-green-600 mr-3" />
                    <span class="font-medium">{{ $t('services.pastel_gouache') }}</span>
                </div>
            </div>
        </div>

        <div class="cta-section text-center">
            <h3 class="text-2xl font-bold mb-4">
                {{ $t('services.collaborate_heading') }}
            </h3>
            <p class="text-gray-600 dark:text-gray-300 mb-6 max-w-2xl mx-auto">
                {{ $t('services.collaborate_text') }}
            </p>
            <UButton
                size="lg"
                color="success"
                variant="solid"
                icon="i-lucide-mail"
                :to="{ path: '/', hash: '#contact' }"
            >
                {{ $t('services.contact_button') }}
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
    transition: background-color 0.3s ease;
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
