<script setup lang="ts">
import { ref, onMounted, onUnmounted } from 'vue';
import { useThemeStore } from '../stores/ThemeStore';

const { locale, setLocale } = useI18n();
const switchLocale = (newLocale: 'ru' | 'en') => {
    setLocale(newLocale);
};

const themeStore = useThemeStore();
const toggleTheme = (event?: MouseEvent) => {
    if (event) {
        event.preventDefault();
        event.stopPropagation();
    }
    if (typeof window === 'undefined') return;
    themeStore.theme = themeStore.theme === 'light' ? 'dark' : 'light';
};

// Interactive elements
const scrollY = ref(0);
const isScrolled = ref(false);
const animatedStats = ref({
    years: 0,
    students: 0,
    works: 0,
});
const statsTarget = { years: 22, students: 2000, works: 150 };
const animationStarted = ref(false);

// Handle scroll for parallax and navbar effects
const handleScroll = () => {
    scrollY.value = window.scrollY;
    isScrolled.value = scrollY.value > 100;

    // Trigger stats animation when teaching section is in view
    if (!animationStarted.value) {
        const teachingSection = document.querySelector('.teaching-section');
        if (teachingSection) {
            const rect = teachingSection.getBoundingClientRect();
            if (rect.top < window.innerHeight * 0.8) {
                animationStarted.value = true;
                animateStats();
            }
        }
    }
};

// Animate statistics counters
const animateStats = () => {
    const duration = 2000; // 2 seconds
    const steps = 60;
    const stepDuration = duration / steps;

    const incrementYears = statsTarget.years / steps;
    const incrementStudents = statsTarget.students / steps;
    const incrementWorks = statsTarget.works / steps;

    let currentStep = 0;
    const timer = setInterval(() => {
        currentStep++;
        animatedStats.value.years = Math.min(
            Math.round(incrementYears * currentStep),
            statsTarget.years
        );
        animatedStats.value.students = Math.min(
            Math.round(incrementStudents * currentStep),
            statsTarget.students
        );
        animatedStats.value.works = Math.min(
            Math.round(incrementWorks * currentStep),
            statsTarget.works
        );

        if (currentStep >= steps) {
            clearInterval(timer);
            // Ensure final values are exact
            animatedStats.value.years = statsTarget.years;
            animatedStats.value.students = statsTarget.students;
            animatedStats.value.works = statsTarget.works;
        }
    }, stepDuration);
};

// Mouse move parallax effect for visual elements
const mouseX = ref(0);
const mouseY = ref(0);

const handleMouseMove = (event: MouseEvent) => {
    mouseX.value = event.clientX;
    mouseY.value = event.clientY;

    // Apply subtle parallax to visual elements
    const visualElements = document.querySelectorAll('.visual-element');
    visualElements.forEach((el, index) => {
        const element = el as HTMLElement;
        const speed = 0.01 + index * 0.005;
        const x = (mouseX.value * speed) % 100;
        const y = (mouseY.value * speed) % 100;
        element.style.transform = `translate(${x * 0.3}px, ${y * 0.3}px)`;
    });
};

// Initialize
onMounted(() => {
    window.addEventListener('scroll', handleScroll);
    window.addEventListener('mousemove', handleMouseMove);
    // Initial check
    handleScroll();
});

onUnmounted(() => {
    window.removeEventListener('scroll', handleScroll);
    window.removeEventListener('mousemove', handleMouseMove);
});
</script>

<template>
    <div class="hero-section">
        <div class="hero-overlay">
            <div class="hero-content">
                <h1 class="hero-title">{{ $t('home.hi') }}</h1>
                <p class="hero-subtitle">{{ $t('home.familiarity') }}</p>
                <div class="scroll-indicator">
                    <div class="mouse">
                        <div class="wheel"></div>
                    </div>
                    <div class="arrow-down"></div>
                </div>
            </div>
        </div>
        <img
            src="../assets/images/img_parallax.webp"
            alt="Pertsukova"
            class="hero-image"
        />
    </div>

    <div class="content-section intro-section">
        <div class="container">
            <div class="section-header">
                <h2 class="section-title">{{ $t('home.familiarity') }}</h2>
                <div class="section-divider"></div>
            </div>
            <div class="intro-grid">
                <div class="intro-card">
                    <div class="intro-icon">🎨</div>
                    <h3 class="intro-card-title">{{ $t('home.artist_title') }}</h3>
                    <p class="intro-card-text">{{ $t('home.whoami') }}</p>
                </div>
                <div class="intro-card">
                    <div class="intro-icon">💡</div>
                    <h3 class="intro-card-title">{{ $t('home.creativity_title') }}</h3>
                    <p class="intro-card-text">{{ $t('home.buy') }}</p>
                </div>
                <div class="intro-card">
                    <div class="intro-icon">🌟</div>
                    <h3 class="intro-card-title">{{ $t('home.development_title') }}</h3>
                    <p class="intro-card-text">
                        {{ $t('home.development_text') }}
                    </p>
                </div>
            </div>
        </div>
    </div>

    <div class="parallax-section parallax-1">
        <div class="parallax-overlay">
            <div class="parallax-content">
                <h2 class="parallax-title">{{ $t('home.draw') }}</h2>
                <p class="parallax-text">
                    {{ $t('home.inspiration_text') }}
                </p>
            </div>
        </div>
    </div>

    <div class="content-section philosophy-section">
        <div class="container">
            <div class="philosophy-content">
                <div class="philosophy-text">
                    <h2 class="section-title">{{ $t('home.philosophy_title') }}</h2>
                    <p class="philosophy-quote">{{ $t('home.world') }}</p>
                    <p class="philosophy-description">
                        {{ $t('home.philosophy_description') }}
                    </p>
                </div>
                <div class="philosophy-visual">
                    <div class="visual-element visual-1">
                        <div class="visual-content">
                            <div class="visual-icon">🎨</div>
                            <h3 class="visual-title">
                                {{ $t('home.expression_title') }}
                            </h3>
                            <p class="visual-text">{{ $t('home.expression_text') }}</p>
                        </div>
                    </div>
                    <div class="visual-element visual-2">
                        <div class="visual-content">
                            <div class="visual-icon">✨</div>
                            <h3 class="visual-title">{{ $t('home.harmony_title') }}</h3>
                            <p class="visual-text">
                                {{ $t('home.harmony_text') }}
                            </p>
                        </div>
                    </div>
                    <div class="visual-element visual-3">
                        <div class="visual-content">
                            <div class="visual-icon">🌱</div>
                            <h3 class="visual-title">{{ $t('home.growth_title') }}</h3>
                            <p class="visual-text">
                                {{ $t('home.growth_text') }}
                            </p>
                        </div>
                    </div>
                </div>
            </div>
        </div>
    </div>

    <div class="parallax-section parallax-2">
        <div class="parallax-overlay">
            <div class="parallax-content">
                <h2 class="parallax-title">{{ $t('home.givemore') }}</h2>
                <p class="parallax-text">
                    {{ $t('home.sharing_text') }}
                </p>
            </div>
        </div>
    </div>

    <div class="content-section teaching-section">
        <div class="container">
            <div class="section-header">
                <h2 class="section-title">{{ $t('home.teaching_title') }}</h2>
                <div class="section-divider"></div>
            </div>
            <div class="teaching-content">
                <p class="teaching-text">{{ $t('home.convinced') }}</p>
                <div class="teaching-stats">
                    <div class="stat-item">
                        <div class="stat-number">{{ animatedStats.years }}+</div>
                        <div class="stat-label">{{ $t('home.years_experience') }}</div>
                    </div>
                    <div class="stat-item">
                        <div class="stat-number">{{ animatedStats.students }}+</div>
                        <div class="stat-label">{{ $t('home.students') }}</div>
                    </div>
                    <div class="stat-item">
                        <div class="stat-number">{{ animatedStats.works }}+</div>
                        <div class="stat-label">{{ $t('home.works') }}</div>
                    </div>
                </div>
            </div>
        </div>
    </div>

    <div class="cta-section">
        <div class="container">
            <div class="cta-content">
                <h2 class="cta-title">{{ $t('home.start') }}</h2>
                <p class="cta-text">
                    {{ $t('home.cta_text') }}
                </p>
                <div class="cta-buttons">
                    <a href="/all-works" class="btn btn-primary">{{
                        $t('home.view_works')
                    }}</a>
                    <a href="/news" class="btn btn-secondary">{{
                        $t('home.read_news')
                    }}</a>
                    <a href="/services" class="btn btn-outline">{{
                        $t('home.services')
                    }}</a>
                </div>
            </div>
        </div>
    </div>
</template>

<style scoped>
/* Hero Section */
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

/* Optimize image loading and size for mobile */
@media (max-width: 767px) {
    .hero-image {
        animation: zoomInMobile 25s ease-in-out infinite alternate;
        object-position: center center;
    }

    .parallax-1,
    .parallax-2 {
        background-size: cover;
        background-position: center center;
        background-attachment: scroll;
    }

    .parallax-section {
        background-attachment: scroll;
    }
}

/* Extra small devices portrait phones */
@media (max-width: 480px) {
    .parallax-section {
        min-height: 40vh;
        background-size: cover !important;
        background-position: center center !important;
    }

    .parallax-content {
        padding: 1.5rem;
    }
}

/* Very small devices (320px and below) */
@media (max-width: 320px) {
    .parallax-section {
        min-height: 35vh;
    }

    .parallax-overlay {
        background: linear-gradient(
            135deg,
            rgba(0, 0, 0, 0.8) 0%,
            rgba(0, 0, 0, 0.5) 100%
        );
    }

    .parallax-content {
        padding: 1rem;
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

/* Theme toggle button */
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

/* Content Sections */
.content-section {
    padding: 6rem 2rem;
    position: relative;
}

.container {
    max-width: 1200px;
    margin: 0 auto;
    padding: 0 1.5rem;
}

.section-header {
    text-align: center;
    margin-bottom: 4rem;
}

.section-title {
    font-size: clamp(1.8rem, 3vw + 0.5rem, 3.5rem);
    font-weight: 700;
    color: var(--color-on-surface);
    margin-bottom: 1.5rem;
    position: relative;
    display: inline-block;
    line-height: 1.2;
}

.section-divider {
    width: 80px;
    height: 4px;
    background: linear-gradient(90deg, #4b9e90, #73d1be);
    margin: 0 auto;
    border-radius: 2px;
}

/* Intro Section */
.intro-grid {
    display: grid;
    grid-template-columns: repeat(auto-fit, minmax(min(100%, 300px), 1fr));
    gap: clamp(1.5rem, 3vw, 3rem);
    margin-top: 3rem;
}

/* Ensure proper grid behavior on very small screens */
@media (max-width: 350px) {
    .intro-grid {
        grid-template-columns: 1fr;
    }
}

.intro-card {
    background: var(--color-surface);
    border-radius: 20px;
    padding: 2.5rem;
    text-align: center;
    box-shadow: 0 10px 30px rgba(0, 0, 0, 0.08);
    transition: all 0.4s ease;
    border: 1px solid rgba(255, 255, 255, 0.1);
    position: relative;
    overflow: hidden;
}

.intro-card::before {
    content: '';
    position: absolute;
    top: 0;
    left: 0;
    width: 100%;
    height: 5px;
    background: linear-gradient(90deg, #4b9e90, #73d1be);
}

.intro-card:hover {
    transform: translateY(-10px);
    box-shadow: 0 20px 40px rgba(0, 0, 0, 0.15);
}

.intro-icon {
    font-size: 3.5rem;
    margin-bottom: 1.5rem;
    display: inline-block;
}

.intro-card-title {
    font-size: 1.8rem;
    font-weight: 600;
    color: var(--color-on-surface);
    margin-bottom: 1rem;
}

.intro-card-text {
    font-size: 1.1rem;
    line-height: 1.7;
    color: var(--color-on-surface);
    opacity: 0.9;
}

/* Parallax Sections */
.parallax-section {
    position: relative;
    min-height: 70vh;
    background-attachment: fixed;
    background-position: center;
    background-repeat: no-repeat;
    background-size: cover;
    display: flex;
    align-items: center;
    justify-content: center;
}

.parallax-1 {
    background-image: url('../assets/images/img_parallax2.webp');
}

.parallax-2 {
    background-image: url('../assets/images/img_parallax3.webp');
}

/* Fallback for browsers that don't support webp */
@supports not (background-image: url('../assets/images/img_parallax2.webp')) {
    .parallax-1 {
        background-image: url('../assets/images/img_parallax2.jpg');
    }
    .parallax-2 {
        background-image: url('../assets/images/img_parallax3.jpg');
    }
}

.parallax-overlay {
    position: absolute;
    top: 0;
    left: 0;
    width: 100%;
    height: 100%;
    background: linear-gradient(135deg, rgba(0, 0, 0, 0.7) 0%, rgba(0, 0, 0, 0.4) 100%);
    display: flex;
    align-items: center;
    justify-content: center;
}

.parallax-content {
    text-align: center;
    color: white;
    max-width: 800px;
    padding: 2rem;
    z-index: 2;
}

.parallax-title {
    font-size: clamp(2rem, 4vw + 0.5rem, 4.5rem);
    font-weight: 800;
    margin-bottom: 1.5rem;
    text-transform: uppercase;
    letter-spacing: clamp(0.1rem, 0.3vw, 0.3rem);
    animation: fadeIn 1.5s ease-out;
    line-height: 1.1;
}

.parallax-text {
    font-size: clamp(1.1rem, 1.5vw + 0.5rem, 1.8rem);
    font-weight: 300;
    line-height: 1.6;
    opacity: 0.9;
    max-width: min(90%, 800px);
    margin: 0 auto;
}

/* Philosophy Section */
.philosophy-section {
    background: linear-gradient(135deg, var(--color-surface) 0%, #f0f0f0 100%);
}

.philosophy-content {
    display: grid;
    grid-template-columns: repeat(auto-fit, minmax(min(100%, 400px), 1fr));
    gap: clamp(2rem, 4vw, 5rem);
    align-items: center;
}

.philosophy-text {
    padding-right: 2rem;
}

.philosophy-quote {
    font-size: 1.8rem;
    font-weight: 500;
    line-height: 1.8;
    color: var(--color-on-surface);
    margin: 2rem 0;
    padding-left: 2rem;
    border-left: 4px solid #4b9e90;
    font-style: italic;
}

.philosophy-description {
    font-size: 1.1rem;
    line-height: 1.7;
    color: var(--color-on-surface);
    opacity: 0.9;
}

.philosophy-visual {
    display: flex;
    flex-direction: column;
    gap: 1.5rem;
}

.visual-element {
    height: 130px;
    border-radius: 12px;
    background: linear-gradient(
        135deg,
        rgba(75, 158, 144, 0.9),
        rgba(115, 209, 190, 0.9)
    );
    opacity: 0.9;
    transition: all 0.4s ease;
    display: flex;
    align-items: center;
    justify-content: center;
    padding: 1.2rem 1.5rem;
    position: relative;
    overflow: hidden;
}

.visual-element::before {
    content: '';
    position: absolute;
    top: 0;
    left: 0;
    width: 100%;
    height: 100%;
    background: linear-gradient(
        135deg,
        rgba(255, 255, 255, 0.1),
        rgba(255, 255, 255, 0.05)
    );
    z-index: 1;
}

.visual-element:hover {
    opacity: 1;
    transform: translateY(-5px);
    box-shadow: 0 10px 25px rgba(75, 158, 144, 0.3);
}

.visual-content {
    position: relative;
    z-index: 2;
    text-align: center;
    color: white;
    width: 100%;
    display: flex;
    flex-direction: column;
    justify-content: center;
    align-items: center;
    padding-bottom: 0.5rem;
}

.visual-icon {
    font-size: 2.2rem;
    margin-bottom: 0.6rem;
    display: block;
    line-height: 1;
}

.visual-title {
    font-size: 1.3rem;
    font-weight: 700;
    margin-bottom: 0.4rem;
    letter-spacing: 0.03rem;
    line-height: 1.2;
}

.visual-text {
    font-size: 0.9rem;
    opacity: 0.9;
    line-height: 1.3;
    margin: 0;
    padding: 0 0.5rem;
}

.visual-1 {
    width: 85%;
    align-self: flex-start;
}

.visual-2 {
    width: 95%;
    align-self: center;
}

.visual-3 {
    width: 75%;
    align-self: flex-end;
}

/* Teaching Section */
.teaching-content {
    text-align: center;
    max-width: 800px;
    margin: 0 auto;
}

.teaching-text {
    font-size: 1.3rem;
    line-height: 1.8;
    color: var(--color-on-surface);
    margin-bottom: 3rem;
}

.teaching-stats {
    display: flex;
    justify-content: center;
    gap: clamp(2rem, 4vw, 6rem);
    margin-top: 3rem;
    flex-wrap: wrap;
}

/* Better wrapping for stats on medium screens */
@media (max-width: 767px) {
    .teaching-stats {
        gap: clamp(1.5rem, 3vw, 3rem);
    }
}

.stat-item {
    text-align: center;
}

.stat-number {
    font-size: 3.5rem;
    font-weight: 800;
    color: #4b9e90;
    margin-bottom: 0.5rem;
}

.stat-label {
    font-size: 1.1rem;
    color: var(--color-on-surface);
    opacity: 0.8;
    text-transform: uppercase;
    letter-spacing: 0.1rem;
}

/* CTA Section */
.cta-section {
    background: linear-gradient(135deg, #4b9e90 0%, #73d1be 100%);
    padding: 6rem 2rem;
    color: white;
    text-align: center;
}

.cta-title {
    font-size: 3.2rem;
    font-weight: 800;
    margin-bottom: 1.5rem;
}

.cta-text {
    font-size: 1.3rem;
    line-height: 1.7;
    max-width: 700px;
    margin: 0 auto 3rem;
    opacity: 0.9;
}

.cta-buttons {
    display: flex;
    gap: clamp(1rem, 2vw, 2rem);
    justify-content: center;
    flex-wrap: wrap;
}

/* Button sizing for better touch targets on mobile */
@media (max-width: 767px) {
    .cta-buttons {
        gap: 1rem;
    }

    .btn {
        min-height: 50px;
        display: flex;
        align-items: center;
        justify-content: center;
    }
}

.btn {
    padding: 1rem 2.5rem;
    border-radius: 50px;
    font-size: 1.1rem;
    font-weight: 600;
    text-decoration: none;
    transition: all 0.3s ease;
    display: inline-block;
    border: 2px solid transparent;
    cursor: pointer;
}

.btn-primary {
    background: white;
    color: #4b9e90;
}

.btn-primary:hover {
    background: #f0f0f0;
    transform: translateY(-3px);
    box-shadow: 0 10px 20px rgba(0, 0, 0, 0.2);
}

.btn-secondary {
    background: rgba(255, 255, 255, 0.2);
    color: white;
    border: 2px solid white;
}

.btn-secondary:hover {
    background: white;
    color: #4b9e90;
    transform: translateY(-3px);
}

.btn-outline {
    background: transparent;
    color: white;
    border: 2px solid white;
}

.btn-outline:hover {
    background: white;
    color: #4b9e90;
    transform: translateY(-3px);
}

/* Animations */
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

@keyframes fadeIn {
    0% {
        opacity: 0;
    }
    100% {
        opacity: 1;
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

/* ============================================
   COMPREHENSIVE RESPONSIVE DESIGN
   ============================================ */

/* 1. EXTRA LARGE SCREENS (1440px and above) */
@media (min-width: 1440px) {
    .container {
        max-width: 1400px;
    }

    .hero-content {
        max-width: min(95%, 1100px);
    }

    .hero-title {
        font-size: clamp(4rem, 4.5vw, 5rem);
        line-height: 1.2;
        letter-spacing: 0.15rem;
    }

    .hero-subtitle {
        font-size: clamp(1.8rem, 1.8vw, 2.2rem);
        max-width: 900px;
        margin-left: auto;
        margin-right: auto;
    }

    .intro-grid {
        grid-template-columns: repeat(3, 1fr);
        gap: 3rem;
    }

    .philosophy-content {
        gap: 5rem;
    }

    .parallax-title {
        font-size: 4.5rem;
    }

    .parallax-text {
        font-size: 1.8rem;
        max-width: 900px;
    }

    .teaching-stats {
        gap: 6rem;
    }

    .stat-number {
        font-size: 4.5rem;
    }
}

/* Handle long words in hero title (especially for Russian locale) */
.hero-title {
    /* Ensure long words can break if necessary */
    overflow-wrap: break-word;
    word-break: keep-all;
    hyphens: auto;
}

/* Specific adjustments for very large screens where Russian text might overflow */
@media (min-width: 1600px) {
    .hero-title {
        font-size: clamp(3.8rem, 4vw, 4.8rem);
        letter-spacing: 0.1rem;
    }

    .hero-content {
        max-width: min(98%, 1200px);
    }
}

/* For ultra-wide screens */
@media (min-width: 2000px) {
    .hero-title {
        font-size: clamp(3.5rem, 3.5vw, 4.5rem);
    }
}

/* 2. LARGE SCREENS / LAPTOPS (1024px - 1439px) */
@media (max-width: 1439px) and (min-width: 1024px) {
    .hero-title {
        font-size: 4.2rem;
    }

    .hero-subtitle {
        font-size: 1.9rem;
    }

    .section-title {
        font-size: 3rem;
    }

    .intro-grid {
        grid-template-columns: repeat(3, 1fr);
        gap: 2.5rem;
    }

    .parallax-title {
        font-size: 3.8rem;
    }

    .parallax-text {
        font-size: 1.5rem;
    }

    .teaching-stats {
        gap: 5rem;
    }

    .stat-number {
        font-size: 4rem;
    }

    .cta-title {
        font-size: 3.5rem;
    }
}

/* 3. TABLETS (768px - 1023px) */
@media (max-width: 1023px) and (min-width: 768px) {
    .hero-title {
        font-size: 3.5rem;
    }

    .hero-subtitle {
        font-size: 1.6rem;
    }

    .section-title {
        font-size: 2.5rem;
    }

    .intro-grid {
        grid-template-columns: repeat(2, 1fr);
        gap: 2rem;
    }

    .philosophy-content {
        grid-template-columns: 1fr;
        gap: 3rem;
    }

    .philosophy-text {
        padding-right: 0;
    }

    .parallax-section {
        background-attachment: scroll;
        min-height: 60vh;
    }

    .parallax-title {
        font-size: 3rem;
    }

    .parallax-text {
        font-size: 1.3rem;
    }

    .teaching-stats {
        gap: 3rem;
    }

    .stat-number {
        font-size: 3.2rem;
    }

    .cta-title {
        font-size: 2.8rem;
    }

    .cta-buttons {
        gap: 1.2rem;
    }

    .btn {
        padding: 0.9rem 2rem;
        font-size: 1rem;
    }

    .visual-element {
        height: 120px;
    }

    .visual-1,
    .visual-2,
    .visual-3 {
        width: 90%;
    }
}

/* 4. MOBILE LANDSCAPE / SMALL TABLETS (480px - 767px) */
@media (max-width: 767px) and (min-width: 480px) {
    .hero-title {
        font-size: 2.8rem;
    }

    .hero-subtitle {
        font-size: 1.4rem;
    }

    .content-section {
        padding: 4rem 1.5rem;
    }

    .section-title {
        font-size: 2.2rem;
    }

    .intro-grid {
        grid-template-columns: 1fr;
        gap: 1.8rem;
    }

    .intro-card {
        padding: 2rem;
    }

    .philosophy-content {
        gap: 2.5rem;
    }

    .philosophy-quote {
        font-size: 1.6rem;
        padding-left: 1.5rem;
    }

    .parallax-section {
        min-height: 50vh;
    }

    .parallax-title {
        font-size: 2.5rem;
        letter-spacing: 0.15rem;
    }

    .parallax-text {
        font-size: 1.2rem;
    }

    .teaching-stats {
        gap: 2.5rem;
        flex-wrap: wrap;
    }

    .stat-item {
        flex: 0 0 calc(50% - 1.25rem);
        margin-bottom: 1.5rem;
    }

    .stat-number {
        font-size: 2.8rem;
    }

    .cta-title {
        font-size: 2.5rem;
    }

    .cta-text {
        font-size: 1.2rem;
    }

    .cta-buttons {
        flex-direction: column;
        align-items: center;
        gap: 1rem;
    }

    .btn {
        width: 100%;
        max-width: 320px;
        padding: 1rem 2rem;
    }

    .visual-element {
        height: 115px;
        padding: 1.2rem;
    }

    .visual-icon {
        font-size: 2rem;
    }

    .visual-title {
        font-size: 1.2rem;
    }

    .visual-text {
        font-size: 0.85rem;
    }

    .visual-1,
    .visual-2,
    .visual-3 {
        width: 100%;
    }

    .theme-toggle-btn {
        top: 1.5rem;
        right: 1.5rem;
        width: 45px;
        height: 45px;
        font-size: 1.6rem;
    }
}

/* 5. SMALL MOBILE (up to 479px) */
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

    .content-section {
        padding: 2.5rem 0.8rem;
    }

    .section-title {
        font-size: 1.7rem;
        padding: 0 0.5rem;
    }

    .intro-grid {
        grid-template-columns: 1fr;
        gap: 1.2rem;
    }

    .intro-card {
        padding: 1.2rem;
    }

    .intro-icon {
        font-size: 2.5rem;
    }

    .intro-card-title {
        font-size: 1.4rem;
    }

    .intro-card-text {
        font-size: 0.95rem;
    }

    .philosophy-content {
        gap: 1.5rem;
    }

    .philosophy-quote {
        font-size: 1.3rem;
        padding-left: 0.8rem;
        margin: 1.2rem 0;
    }

    .philosophy-description {
        font-size: 0.95rem;
    }

    .parallax-section {
        min-height: 35vh;
        background-position: center center;
        background-size: cover;
        background-attachment: scroll;
    }

    .parallax-title {
        font-size: 1.7rem;
        letter-spacing: 0.05rem;
        padding: 0 0.8rem;
        line-height: 1.2;
    }

    .parallax-text {
        font-size: 1rem;
        padding: 0 0.8rem;
        line-height: 1.4;
    }

    .parallax-content {
        padding: 1rem;
        width: 100%;
    }

    .teaching-stats {
        flex-direction: column;
        gap: 1.2rem;
    }

    .stat-item {
        flex: 0 0 100%;
        margin-bottom: 1.2rem;
    }

    .stat-number {
        font-size: 2.2rem;
    }

    .stat-label {
        font-size: 0.9rem;
    }

    .cta-title {
        font-size: 1.8rem;
    }

    .cta-text {
        font-size: 1rem;
        padding: 0 0.5rem;
    }

    .cta-buttons {
        flex-direction: column;
        gap: 0.7rem;
    }

    .btn {
        width: 100%;
        max-width: 260px;
        padding: 0.8rem 1.2rem;
        font-size: 0.95rem;
    }

    .visual-element {
        height: 100px;
        padding: 0.8rem;
    }

    .visual-icon {
        font-size: 1.6rem;
        margin-bottom: 0.3rem;
    }

    .visual-title {
        font-size: 1rem;
        margin-bottom: 0.2rem;
    }

    .visual-text {
        font-size: 0.75rem;
        padding: 0 0.2rem;
    }

    .visual-1,
    .visual-2,
    .visual-3 {
        width: 100%;
        align-self: center;
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

/* 6. EXTRA SMALL MOBILE (320px and below) */
@media (max-width: 320px) {
    /* Hero section adjustments */
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

    /* Parallax sections */
    .parallax-title {
        font-size: 1.5rem;
        padding: 0 0.5rem;
        line-height: 1.2;
        margin-bottom: 0.8rem;
    }

    .parallax-text {
        font-size: 0.9rem;
        padding: 0 0.5rem;
        line-height: 1.3;
    }

    .parallax-section {
        min-height: 30vh;
    }

    /* Content sections */
    .section-title {
        font-size: 1.5rem;
        padding: 0 0.3rem;
    }

    .content-section {
        padding: 2rem 0.5rem;
    }

    .container {
        padding-left: 0.5rem;
        padding-right: 0.5rem;
    }

    /* Intro cards */
    .intro-card {
        padding: 1rem;
    }

    .intro-card-title {
        font-size: 1.2rem;
    }

    .intro-card-text {
        font-size: 0.9rem;
        line-height: 1.5;
    }

    .intro-icon {
        font-size: 2.2rem;
        margin-bottom: 0.8rem;
    }

    /* Buttons */
    .btn {
        max-width: 240px;
        padding: 0.7rem 1rem;
        font-size: 0.9rem;
        min-height: 44px;
    }

    /* Stats */
    .stat-number {
        font-size: 2rem;
    }

    .stat-label {
        font-size: 0.85rem;
    }

    /* CTA section */
    .cta-title {
        font-size: 1.6rem;
    }

    .cta-text {
        font-size: 0.9rem;
        padding: 0 0.3rem;
    }

    /* Visual elements */
    .visual-element {
        height: 90px;
        padding: 0.7rem;
    }

    .visual-icon {
        font-size: 1.4rem;
        margin-bottom: 0.2rem;
    }

    .visual-title {
        font-size: 0.95rem;
        margin-bottom: 0.1rem;
    }

    .visual-text {
        font-size: 0.7rem;
        padding: 0 0.1rem;
    }

    /* Theme toggle */
    .theme-toggle-btn {
        top: 0.6rem;
        right: 0.6rem;
        width: 34px;
        height: 34px;
        font-size: 1.1rem;
    }

    /* Scroll indicator */
    .scroll-indicator {
        bottom: 1rem;
    }

    /* Ensure background images are properly centered */
    .hero-image {
        object-position: center center;
    }

    .parallax-1,
    .parallax-2 {
        background-position: center center !important;
        background-size: cover !important;
    }

    /* Prevent text overflow in all text elements */
    .philosophy-quote,
    .philosophy-description,
    .teaching-text,
    .intro-card-text {
        word-break: break-word;
        overflow-wrap: break-word;
        hyphens: auto;
    }

    /* Adjust padding for philosophy section */
    .philosophy-quote {
        font-size: 1.2rem;
        padding-left: 0.6rem;
        margin: 1rem 0;
    }

    .philosophy-description {
        font-size: 0.9rem;
    }
}

/* ============================================
   TOUCH INTERACTION ENHANCEMENTS
   ============================================ */

/* Improve touch interactions for mobile/touch devices */
@media (hover: none) and (pointer: coarse) {
    .intro-card:hover {
        transform: none;
    }

    .intro-card:active {
        transform: scale(0.98);
        transition: transform 0.1s ease;
    }

    .btn:hover {
        transform: none;
    }

    .btn:active {
        transform: scale(0.95);
        opacity: 0.9;
    }

    .visual-element:hover {
        transform: none;
        opacity: 0.8;
    }

    .visual-element:active {
        opacity: 0.9;
        transform: scale(0.98);
    }

    .theme-toggle-btn:hover {
        transform: none;
    }

    .theme-toggle-btn:active {
        transform: scale(0.85);
        transition: transform 0.1s ease;
    }
}

/* Better touch targets for mobile navigation */
@media (max-width: 767px) {
    .intro-card,
    .visual-element,
    .btn {
        -webkit-tap-highlight-color: rgba(0, 0, 0, 0.1);
        tap-highlight-color: rgba(0, 0, 0, 0.1);
    }

    .btn {
        min-height: 52px;
        padding-top: 0.9rem;
        padding-bottom: 0.9rem;
    }

    .theme-toggle-btn {
        min-width: 50px;
        min-height: 50px;
    }

    /* Prevent accidental horizontal scroll on mobile */
    .hero-section,
    .content-section,
    .parallax-section {
        overflow-x: hidden;
    }
}

/* Smooth scrolling for better UX */
@media (prefers-reduced-motion: no-preference) {
    html {
        scroll-behavior: smooth;
    }
}

/* ============================================
   PERFORMANCE & ACCESSIBILITY OPTIMIZATIONS
   ============================================ */

/* Prevent animation jank on mobile and respect reduced motion */
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

/* Optimize container padding for different viewports */
.container {
    padding-left: clamp(1rem, 3vw, 2rem);
    padding-right: clamp(1rem, 3vw, 2rem);
}

/* Responsive spacing for sections */
.content-section {
    padding-top: clamp(3rem, 6vw, 8rem);
    padding-bottom: clamp(3rem, 6vw, 8rem);
}

.cta-section {
    padding-top: clamp(4rem, 8vw, 10rem);
    padding-bottom: clamp(4rem, 8vw, 10rem);
}

/* Ensure images don't cause layout shift */
.hero-image,
.parallax-1,
.parallax-2 {
    aspect-ratio: 16/9;
}

@media (max-width: 768px) {
    .hero-image,
    .parallax-1,
    .parallax-2 {
        aspect-ratio: 4/3;
    }
}

/* Improve readability on very large screens */
@media (min-width: 1600px) {
    .container {
        max-width: 1500px;
    }

    .intro-card-text,
    .philosophy-description,
    .teaching-text,
    .cta-text {
        font-size: 1.2rem;
        line-height: 1.8;
    }
}

/* High DPI screen optimizations */
@media (-webkit-min-device-pixel-ratio: 2), (min-resolution: 192dpi) {
    .hero-image {
        filter: brightness(0.7) contrast(1.05);
    }
}

/* Dark theme adjustments */
.body_theme_dark .intro-card {
    background: var(--color-surface);
    border-color: rgba(255, 255, 255, 0.05);
}

.body_theme_dark .philosophy-section {
    background: linear-gradient(135deg, var(--color-surface) 0%, #2a2a2a 100%);
}

.body_theme_dark .intro-card-text,
.body_theme_dark .philosophy-quote,
.body_theme_dark .philosophy-description,
.body_theme_dark .teaching-text {
    color: var(--color-on-surface);
}
</style>
