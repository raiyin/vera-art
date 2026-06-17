<script setup lang="ts">
import { ref, onMounted, onUnmounted } from 'vue';
import { useThemeStore } from '~/stores/ThemeStore';

const { t } = useI18n();

const themeStore = useThemeStore();

const scrollY = ref(0);
const mouseX = ref(0);
const mouseY = ref(0);

const toggleTheme = (event?: MouseEvent) => {
    if (event) {
        event.preventDefault();
        event.stopPropagation();
    }
    if (typeof window === 'undefined') return;
    themeStore.theme = themeStore.theme === 'light' ? 'dark' : 'light';
};

const handleScroll = () => {
    scrollY.value = window.scrollY;
};

const handleMouseMove = (event: MouseEvent) => {
    mouseX.value = event.clientX;
    mouseY.value = event.clientY;

    const visualElements = document.querySelectorAll('.visual-element');
    visualElements.forEach((el, index) => {
        const element = el as HTMLElement;
        const speed = 0.01 + index * 0.005;
        const x = (mouseX.value * speed) % 100;
        const y = (mouseY.value * speed) % 100;
        element.style.transform = `translate(${x * 0.3}px, ${y * 0.3}px)`;
    });
};

onMounted(() => {
    window.addEventListener('scroll', handleScroll);
    window.addEventListener('mousemove', handleMouseMove);
    handleScroll();
});

onUnmounted(() => {
    window.removeEventListener('scroll', handleScroll);
    window.removeEventListener('mousemove', handleMouseMove);
});
</script>

<template>
    <div class="hero-section">
        <button
            class="theme-toggle-btn"
            :aria-label="t('home.toggleTheme')"
            @click="toggleTheme"
        >
            {{ themeStore.theme === 'light' ? '🌙' : '☀️' }}
        </button>

        <div class="hero-overlay">
            <div class="hero-content">
                <h1 class="hero-title">{{ t('home.hi') }}</h1>
                <p class="hero-subtitle">{{ t('home.familiarity') }}</p>
                <div class="scroll-indicator">
                    <div class="mouse">
                        <div class="wheel" />
                    </div>
                    <div class="arrow-down" />
                </div>
            </div>
        </div>
        <img
            src="../assets/images/img_parallax.webp"
            alt="Pertsukova"
            class="hero-image"
        />
    </div>
</template>

<style scoped>
.hero-section {
    position: relative;
    height: 100vh;
    overflow: hidden;
    display: flex;
    align-items: center;
    justify-content: center;
}

.hero-image {
    position: absolute;
    top: 0;
    left: 0;
    width: 100%;
    height: 100%;
    object-fit: cover;
    z-index: 1;
    filter: brightness(0.7);
    animation: zoomIn 20s ease-in-out infinite alternate;
}

@media (max-width: 767px) {
    .hero-image {
        animation: zoomInMobile 25s ease-in-out infinite alternate;
        object-position: center center;
    }
}

@keyframes zoomInMobile {
    0% {
        transform: scale(1);
    }
    100% {
        transform: scale(1.05);
    }
}

.hero-overlay {
    position: absolute;
    top: 0;
    left: 0;
    width: 100%;
    height: 100%;
    background: linear-gradient(135deg, rgba(0, 0, 0, 0.7) 0%, rgba(0, 0, 0, 0.4) 100%);
    z-index: 2;
    display: flex;
    align-items: center;
    justify-content: center;
}

.theme-toggle-btn {
    position: absolute;
    top: 2rem;
    right: 2rem;
    z-index: 1000;
    width: 50px;
    height: 50px;
    border-radius: 50%;
    background: rgba(255, 255, 255, 0.2);
    border: 2px solid rgba(255, 255, 255, 0.3);
    color: white;
    font-size: 1.8rem;
    display: flex;
    align-items: center;
    justify-content: center;
    cursor: pointer;
    transition: all 0.3s ease;
    backdrop-filter: blur(10px);
    pointer-events: auto;
    outline: 2px solid rgba(255, 255, 255, 0.5);
}

.theme-toggle-btn:hover {
    background: rgba(255, 255, 255, 0.3);
    border-color: rgba(255, 255, 255, 0.5);
    transform: scale(1.1);
    outline-color: rgba(255, 255, 255, 0.8);
}

.theme-toggle-btn:active {
    transform: scale(0.95);
}

.hero-content {
    text-align: center;
    color: white;
    z-index: 3;
    max-width: min(95%, 1000px);
    padding: 2rem;
    animation: fadeInUp 1s ease-out;
    width: 100%;
    box-sizing: border-box;
}

.hero-title {
    font-size: clamp(2.2rem, 4.2vw + 0.7rem, 4.8rem);
    font-weight: 800;
    margin-bottom: 1rem;
    letter-spacing: clamp(0.03rem, 0.3vw, 0.2rem);
    text-transform: uppercase;
    background: linear-gradient(90deg, #fff, #f8f8f8);
    -webkit-background-clip: text;
    -webkit-text-fill-color: transparent;
    background-clip: text;
    text-shadow: 0 5px 15px rgba(0, 0, 0, 0.2);
    line-height: 1.2;
    word-break: break-word;
    overflow-wrap: break-word;
    hyphens: auto;
    padding: 0 0.5rem;
    max-width: 100%;
    box-sizing: border-box;
}

.hero-subtitle {
    font-size: clamp(1.2rem, 2vw + 0.5rem, 2.2rem);
    font-weight: 300;
    margin-bottom: 3rem;
    opacity: 0.9;
    line-height: 1.6;
    max-width: 90%;
    margin-left: auto;
    margin-right: auto;
}

.scroll-indicator {
    position: absolute;
    bottom: 2rem;
    left: 50%;
    transform: translateX(-50%);
    display: flex;
    flex-direction: column;
    align-items: center;
    gap: 0.5rem;
    animation: bounce 2s infinite;
}

.mouse {
    width: 30px;
    height: 50px;
    border: 2px solid white;
    border-radius: 20px;
    display: flex;
    justify-content: center;
    padding-top: 10px;
}

.wheel {
    width: 4px;
    height: 10px;
    background-color: white;
    border-radius: 2px;
    animation: scroll 1.5s infinite;
}

.arrow-down {
    width: 20px;
    height: 20px;
    border-right: 2px solid white;
    border-bottom: 2px solid white;
    transform: rotate(45deg);
}

@keyframes zoomIn {
    0% {
        transform: scale(1);
    }
    100% {
        transform: scale(1.1);
    }
}

@keyframes fadeInUp {
    0% {
        opacity: 0;
        transform: translateY(30px);
    }
    100% {
        opacity: 1;
        transform: translateY(0);
    }
}

@keyframes bounce {
    0%,
    20%,
    50%,
    80%,
    100% {
        transform: translateY(0) translateX(-50%);
    }
    40% {
        transform: translateY(-10px) translateX(-50%);
    }
    60% {
        transform: translateY(-5px) translateX(-50%);
    }
}

@keyframes scroll {
    0% {
        opacity: 1;
        transform: translateY(0);
    }
    100% {
        opacity: 0;
        transform: translateY(20px);
    }
}

@media (max-width: 767px) {
    .theme-toggle-btn {
        top: 1.5rem;
        right: 1.5rem;
        width: 45px;
        height: 45px;
        font-size: 1.6rem;
    }
}

@media (max-width: 479px) {
    .hero-title {
        font-size: clamp(1.8rem, 6vw, 2.2rem);
        line-height: 1.1;
        padding: 0 0.5rem;
        word-break: break-word;
        overflow-wrap: break-word;
    }

    .hero-subtitle {
        font-size: clamp(1rem, 3.5vw, 1.2rem);
        line-height: 1.4;
        padding: 0 1rem;
        max-width: 100%;
    }

    .hero-content {
        padding: 1rem;
        width: 100%;
    }

    .theme-toggle-btn {
        top: 0.8rem;
        right: 0.8rem;
        width: 36px;
        height: 36px;
        font-size: 1.2rem;
    }

    .scroll-indicator {
        bottom: 1.2rem;
    }
}

@media (max-width: 320px) {
    .hero-title {
        font-size: 1.6rem;
        letter-spacing: 0.05rem;
        padding: 0 0.3rem;
        line-height: 1.1;
        margin-bottom: 0.8rem;
    }

    .hero-subtitle {
        font-size: 0.95rem;
        padding: 0 0.5rem;
        margin-bottom: 2rem;
        line-height: 1.3;
    }

    .hero-content {
        padding: 0.5rem;
        width: 100%;
    }

    .hero-section {
        height: 95vh;
        min-height: 500px;
    }
}

@media (hover: none) and (pointer: coarse) {
    .theme-toggle-btn:hover {
        transform: none;
    }

    .theme-toggle-btn:active {
        transform: scale(0.85);
        transition: transform 0.1s ease;
    }
}

@media (max-width: 767px) {
    .theme-toggle-btn {
        min-width: 50px;
        min-height: 50px;
    }

    .hero-section {
        overflow-x: hidden;
    }
}

@media (prefers-reduced-motion: reduce) {
    *,
    *::before,
    *::after {
        animation-duration: 0.01ms !important;
        animation-iteration-count: 1 !important;
        transition-duration: 0.01ms !important;
        scroll-behavior: auto !important;
    }

    .hero-image {
        animation: none !important;
    }

    .scroll-indicator {
        animation: none !important;
    }
}

.hero-image {
    aspect-ratio: 16/9;
}

@media (max-width: 768px) {
    .hero-image {
        aspect-ratio: 4/3;
    }
}

@media (-webkit-min-device-pixel-ratio: 2), (min-resolution: 192dpi) {
    .hero-image {
        filter: brightness(0.7) contrast(1.05);
    }
}
</style>
