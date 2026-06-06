<script lang="ts">
export default {};
</script>

<template>
    <UContainer class="paydelivery">
        <div class="hero-section">
            <h1 class="text-4xl md:text-5xl font-bold text-center mb-6">
                {{ $t('header.payment') }}
            </h1>
            <p
                class="text-lg text-center text-gray-600 dark:text-gray-300 max-w-3xl mx-auto mb-12"
            >
                {{ $t('payment.p_one') }}
            </p>
        </div>

        <div class="process-section mb-12">
            <h2 class="text-2xl font-bold text-center mb-8">
                {{ $t('payment.how_it_works') }}
            </h2>
            <div class="grid grid-cols-1 md:grid-cols-3 gap-8">
                <div
                    class="step-card text-center p-6 rounded-2xl bg-white dark:bg-gray-800 shadow-lg hover:shadow-xl transition-shadow"
                >
                    <div class="step-icon-wrapper mb-4">
                        <div
                            class="inline-flex items-center justify-center w-16 h-16 rounded-full bg-green-100 dark:bg-green-900/30 text-green-600"
                        >
                            <UIcon name="i-lucide-message-square" class="w-8 h-8" />
                        </div>
                        <div
                            class="step-number mt-2 text-sm font-semibold text-green-600"
                        >
                            {{ $t('payment.step1') }}
                        </div>
                    </div>
                    <h3 class="text-xl font-semibold mb-3">
                        {{ $t('payment.contact_discuss') }}
                    </h3>
                    <p class="text-gray-700 dark:text-gray-300">
                        {{ $t('payment.contact_text') }}
                    </p>
                </div>

                <div
                    class="step-card text-center p-6 rounded-2xl bg-white dark:bg-gray-800 shadow-lg hover:shadow-xl transition-shadow"
                >
                    <div class="step-icon-wrapper mb-4">
                        <div
                            class="inline-flex items-center justify-center w-16 h-16 rounded-full bg-green-100 dark:bg-green-900/30 text-green-600"
                        >
                            <UIcon name="i-lucide-credit-card" class="w-8 h-8" />
                        </div>
                        <div
                            class="step-number mt-2 text-sm font-semibold text-green-600"
                        >
                            {{ $t('payment.step2') }}
                        </div>
                    </div>
                    <h3 class="text-xl font-semibold mb-3">
                        {{ $t('payment.agree_pay') }}
                    </h3>
                    <p class="text-gray-700 dark:text-gray-300">
                        {{ $t('payment.agree_text') }}
                    </p>
                </div>

                <div
                    class="step-card text-center p-6 rounded-2xl bg-white dark:bg-gray-800 shadow-lg hover:shadow-xl transition-shadow"
                >
                    <div class="step-icon-wrapper mb-4">
                        <div
                            class="inline-flex items-center justify-center w-16 h-16 rounded-full bg-green-100 dark:bg-green-900/30 text-green-600"
                        >
                            <UIcon name="i-lucide-package" class="w-8 h-8" />
                        </div>
                        <div
                            class="step-number mt-2 text-sm font-semibold text-green-600"
                        >
                            {{ $t('payment.step3') }}
                        </div>
                    </div>
                    <h3 class="text-xl font-semibold mb-3">
                        {{ $t('payment.creation_delivery') }}
                    </h3>
                    <p class="text-gray-700 dark:text-gray-300">
                        {{ $t('payment.creation_text') }}
                    </p>
                </div>
            </div>
        </div>

        <div class="note-section bg-green-50 dark:bg-green-900/20 rounded-2xl p-8 mb-12">
            <h3 class="text-2xl font-bold mb-4 text-center">
                {{ $t('payment.important_notes') }}
            </h3>
            <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
                <div class="flex items-start">
                    <UIcon
                        name="i-lucide-info"
                        class="w-6 h-6 text-green-600 mr-3 mt-1 shrink-0"
                    />
                    <div>
                        <h4 class="font-semibold mb-1">
                            {{ $t('payment.custom_commissions') }}
                        </h4>
                        <p class="text-gray-700 dark:text-gray-300 text-sm">
                            {{ $t('payment.custom_text') }}
                        </p>
                    </div>
                </div>
                <div class="flex items-start">
                    <UIcon
                        name="i-lucide-clock"
                        class="w-6 h-6 text-green-600 mr-3 mt-1 shrink-0"
                    />
                    <div>
                        <h4 class="font-semibold mb-1">
                            {{ $t('payment.production_time') }}
                        </h4>
                        <p class="text-gray-700 dark:text-gray-300 text-sm">
                            {{ $t('payment.production_text') }}
                        </p>
                    </div>
                </div>
            </div>
        </div>

        <div class="cta-section text-center">
            <h3 class="text-2xl font-bold mb-4">{{ $t('payment.ready_order') }}</h3>
            <p class="text-gray-600 dark:text-gray-300 mb-6 max-w-2xl mx-auto">
                {{ $t('payment.ready_text') }}
            </p>
            <div class="flex flex-wrap justify-center gap-4">
                <UButton
                    size="lg"
                    color="success"
                    variant="solid"
                    icon="i-lucide-shopping-cart"
                    to="/art-store"
                >
                    {{ $t('payment.browse_shop_button') }}
                </UButton>
                <UButton
                    size="lg"
                    color="neutral"
                    variant="outline"
                    icon="i-lucide-mail"
                    :to="{ path: '/', hash: '#contact' }"
                >
                    {{ $t('payment.contact_button') }}
                </UButton>
            </div>
        </div>
    </UContainer>
</template>

<style scoped>
.paydelivery {
    margin-top: 3rem;
    flex: 1 0 auto;
    padding-bottom: 4rem;
}

.hero-section {
    margin-bottom: 4rem;
}

.step-card {
    transition: transform 0.3s ease;
}

.step-card:hover {
    transform: translateY(-5px);
}

.note-section {
    transition: background-color 0.3s ease;
}

@media (max-width: 768px) {
    .paydelivery {
        margin-top: 2rem;
    }
    .hero-section h1 {
        font-size: 2.5rem;
    }
    .details-section {
        grid-template-columns: 1fr;
    }
}
</style>
