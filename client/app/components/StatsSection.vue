<script setup lang="ts">
import { ref, onMounted, onUnmounted } from 'vue';

const { t } = useI18n();

const animatedStats = ref({
    years: 0,
    students: 0,
    works: 0,
});
const statsTarget = { years: 22, students: 2000, works: 150 };
const animationStarted = ref(false);

const animateStats = () => {
    const duration = 2000;
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
            animatedStats.value.years = statsTarget.years;
            animatedStats.value.students = statsTarget.students;
            animatedStats.value.works = statsTarget.works;
        }
    }, stepDuration);
};

const handleScroll = () => {
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

onMounted(() => {
    window.addEventListener('scroll', handleScroll);
    handleScroll();
});

onUnmounted(() => {
    window.removeEventListener('scroll', handleScroll);
});
</script>

<template>
    <div class="content-section teaching-section">
        <div class="container">
            <div class="section-header">
                <h2 class="section-title">{{ t('home.teaching_title') }}</h2>
                <div class="section-divider" />
            </div>
            <div class="teaching-content">
                <p class="teaching-text">{{ t('home.convinced') }}</p>
                <div class="teaching-stats">
                    <div class="stat-item">
                        <div class="stat-number">{{ animatedStats.years }}+</div>
                        <div class="stat-label">{{ t('home.years_experience') }}</div>
                    </div>
                    <div class="stat-item">
                        <div class="stat-number">{{ animatedStats.students }}+</div>
                        <div class="stat-label">{{ t('home.students') }}</div>
                    </div>
                    <div class="stat-item">
                        <div class="stat-number">{{ animatedStats.works }}+</div>
                        <div class="stat-label">{{ t('home.works') }}</div>
                    </div>
                </div>
            </div>
        </div>
    </div>
</template>

<style scoped>
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

@media (max-width: 767px) {
    .stat-item {
        flex: 0 0 calc(50% - 1.25rem);
        margin-bottom: 1.5rem;
    }

    .stat-number {
        font-size: 2.8rem;
    }
}

@media (max-width: 479px) {
    .content-section {
        padding: 2.5rem 0.8rem;
    }

    .section-title {
        font-size: 1.7rem;
        padding: 0 0.5rem;
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
}

@media (max-width: 320px) {
    .content-section {
        padding: 2rem 0.5rem;
    }

    .container {
        padding-left: 0.5rem;
        padding-right: 0.5rem;
    }

    .section-title {
        font-size: 1.5rem;
        padding: 0 0.3rem;
    }

    .stat-number {
        font-size: 2rem;
    }

    .stat-label {
        font-size: 0.85rem;
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
}

.content-section {
    padding-top: clamp(3rem, 6vw, 8rem);
    padding-bottom: clamp(3rem, 6vw, 8rem);
}

.container {
    padding-left: clamp(1rem, 3vw, 2rem);
    padding-right: clamp(1rem, 3vw, 2rem);
}
</style>
