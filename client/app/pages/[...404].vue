<script lang="ts">
export default {
    setup() {
        const { t, locale } = useI18n();
        const searchMessage = ref('');

        const handleSearch = (event: KeyboardEvent) => {
            const input = event.target as HTMLInputElement;
            if (input.value.trim()) {
                // In a real app, you would implement search functionality
                console.log('Searching for:', input.value);
                // Show search info via UAlert
                searchMessage.value =
                    locale.value === 'ru'
                        ? `Поиск: ${input.value} (функция поиска в разработке)`
                        : `Search: ${input.value} (search functionality in development)`;
            }
        };

        return { t, locale, handleSearch, searchMessage };
    },
};
</script>

<template>
    <UContainer class="not-found min-h-screen flex flex-col max-w-none w-full">
        <div class="flex-1 flex flex-col items-center justify-center py-12 px-0">
            <!-- Decorative elements -->
            <div class="relative w-full max-w-4xl mx-auto">
                <!-- Paint splatter background elements -->
                <div class="absolute -top-20 -left-20 w-64 h-64">
                    <div class="w-full h-full rounded-full paint-glow-1 blur-3xl"></div>
                </div>
                <div class="absolute -bottom-20 -right-20 w-80 h-80">
                    <div class="w-full h-full rounded-full paint-glow-2 blur-3xl"></div>
                </div>

                <!-- Main content -->
                <div class="relative z-10 text-center">
                    <!-- Animated 404 number -->
                    <div class="relative mb-8">
                        <h1
                            class="text-9xl md:text-[12rem] font-bold tracking-tighter text-gray-900 dark:text-white"
                        >
                            <span class="relative inline-block">
                                <span
                                    class="absolute -inset-4 digit-glow-1 rounded-3xl blur-xl"
                                ></span>
                                <span class="relative">4</span>
                            </span>
                            <span class="relative inline-block mx-2 md:mx-4">
                                <span
                                    class="absolute -inset-4 digit-glow-2 rounded-3xl blur-xl"
                                ></span>
                                <span class="relative">0</span>
                            </span>
                            <span class="relative inline-block">
                                <span
                                    class="absolute -inset-4 digit-glow-3 rounded-3xl blur-xl"
                                ></span>
                                <span class="relative">4</span>
                            </span>
                        </h1>

                        <!-- Brush stroke underline -->
                        <div class="relative mt-8 mb-12">
                            <div
                                class="h-2 w-64 mx-auto bg-linear-to-r from-transparent via-green-400 to-transparent rounded-full"
                            ></div>
                            <div
                                class="h-1 w-48 mx-auto bg-linear-to-r from-transparent via-emerald-500 to-transparent rounded-full mt-1"
                            ></div>
                        </div>
                    </div>

                    <!-- Title -->
                    <h2
                        class="text-4xl md:text-5xl font-bold text-gray-900 dark:text-white mb-6"
                    >
                        {{ $t('notfound.title') }}
                    </h2>

                    <!-- Description -->
                    <p
                        class="text-xl text-gray-600 dark:text-gray-300 max-w-2xl mx-auto mb-10"
                    >
                        {{ $t('notfound.description') }}
                    </p>

                    <!-- Search info alert -->
                    <UAlert
                        v-if="searchMessage"
                        :title="locale === 'ru' ? 'Информация' : 'Information'"
                        :description="searchMessage"
                        icon="i-heroicons-information-circle"
                        color="info"
                        variant="outline"
                        class="mb-6"
                        @close="searchMessage = ''"
                    />

                    <!-- Action buttons -->
                    <div class="flex flex-col sm:flex-row gap-4 justify-center mb-12">
                        <UButton
                            to="/"
                            icon="i-heroicons-home"
                            size="xl"
                            color="primary"
                            variant="solid"
                            class="px-8 py-3 text-lg"
                        >
                            {{ $t('notfound.back_home') }}
                        </UButton>

                        <UButton
                            to="/all-works"
                            icon="i-heroicons-paint-brush"
                            size="xl"
                            color="success"
                            variant="outline"
                            class="px-8 py-3 text-lg"
                        >
                            {{ $t('notfound.browse_works') }}
                        </UButton>

                        <UButton
                            to="/shop"
                            icon="i-heroicons-shopping-bag"
                            size="xl"
                            color="info"
                            variant="ghost"
                            class="px-8 py-3 text-lg"
                        >
                            {{ $t('notfound.visit_shop') }}
                        </UButton>
                    </div>

                    <!-- Artistic quote -->
                    <div class="mt-16 pt-8 border-t border-gray-200 dark:border-gray-800">
                        <p
                            class="text-lg italic text-gray-500 dark:text-gray-400 max-w-2xl mx-auto"
                        >
                            "{{ $t('notfound.quote') }}"
                        </p>
                        <p class="text-gray-400 dark:text-gray-500 mt-2">
                            — {{ $t('notfound.author') }}
                        </p>
                    </div>
                </div>
            </div>
        </div>
    </UContainer>
</template>

<style scoped>
.not-found {
    background: linear-gradient(
        to bottom,
        transparent 0%,
        #f8fafc 20%,
        #f0fdf4 80%,
        transparent 100%
    );
}

.dark .not-found {
    background: linear-gradient(
        to bottom,
        transparent 0%,
        #0f172a 20%,
        #052e16 80%,
        transparent 100%
    );
}

/* Custom gradients for transparent edges */
.digit-glow-1 {
    background: radial-gradient(
        ellipse at center,
        rgba(74, 222, 128, 0.4) 0%,
        rgba(74, 222, 128, 0.3) 20%,
        rgba(74, 222, 128, 0.1) 40%,
        transparent 60%
    );
}

.digit-glow-2 {
    background: radial-gradient(
        ellipse at center,
        rgba(52, 211, 153, 0.4) 0%,
        rgba(52, 211, 153, 0.3) 20%,
        rgba(52, 211, 153, 0.1) 40%,
        transparent 60%
    );
}

.digit-glow-3 {
    background: radial-gradient(
        ellipse at center,
        rgba(45, 212, 191, 0.4) 0%,
        rgba(45, 212, 191, 0.3) 20%,
        rgba(45, 212, 191, 0.1) 40%,
        transparent 60%
    );
}

.paint-glow-1 {
    background: radial-gradient(
        ellipse at center,
        rgba(74, 222, 128, 0.15) 0%,
        rgba(74, 222, 128, 0.1) 20%,
        rgba(74, 222, 128, 0.05) 40%,
        transparent 60%
    );
}

.paint-glow-2 {
    background: radial-gradient(
        ellipse at center,
        rgba(52, 211, 153, 0.15) 0%,
        rgba(52, 211, 153, 0.1) 20%,
        rgba(52, 211, 153, 0.05) 40%,
        transparent 60%
    );
}

/* Animation for the 404 numbers */
@keyframes float {
    0%,
    100% {
        transform: translateY(0px);
    }
    50% {
        transform: translateY(-10px);
    }
}

h1 span:nth-child(1) {
    animation: float 3s ease-in-out infinite;
}

h1 span:nth-child(2) {
    animation: float 3s ease-in-out infinite 0.5s;
}

h1 span:nth-child(3) {
    animation: float 3s ease-in-out infinite 1s;
}

/* Brush stroke animation */
@keyframes brushStroke {
    0% {
        width: 0;
        opacity: 0;
    }
    100% {
        width: 64px;
        opacity: 1;
    }
}

.relative.mt-8.mb-12 div:first-child {
    animation: brushStroke 1.5s ease-out forwards;
    animation-delay: 0.5s;
}

.relative.mt-8.mb-12 div:last-child {
    animation: brushStroke 1s ease-out forwards;
    animation-delay: 1s;
}
</style>
