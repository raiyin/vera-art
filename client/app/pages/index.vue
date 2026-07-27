<script setup lang="ts">
    const { locale, setLocale, } = useI18n();

    const switchLocale = (newLocale: 'ru' | 'en',) => {
        setLocale(newLocale,);
    };
</script>

<template>
    <div class="home-page">
        <HeroSection />

        <!-- Intro Section -->
        <div class="content-section intro-section">
            <div class="container">
                <div class="section-header">
                    <h2 class="section-title">
                        {{ $t('home.familiarity',) }}
                    </h2>
                    <div class="section-divider" />
                </div>
                <div class="intro-grid">
                    <div class="intro-card">
                        <div class="intro-icon">
                            <svg
                                class="icon-svg"
                                viewBox="0 0 24 24"
                                fill="none"
                                stroke="currentColor"
                                stroke-width="1.5"
                            >
                                <path d="M12 2L2 7l10 5 10-5-10-5zM2 17l10 5 10-5M2 12l10 5 10-5" />
                            </svg>
                        </div>
                        <h3 class="intro-card-title">
                            {{ $t('home.artist_title',) }}
                        </h3>
                        <p class="intro-card-text">
                            {{ $t('home.whoami',) }}
                        </p>
                    </div>
                    <div class="intro-card">
                        <div class="intro-icon">
                            <svg
                                class="icon-svg"
                                viewBox="0 0 24 24"
                                fill="none"
                                stroke="currentColor"
                                stroke-width="1.5"
                            >
                                <path d="M12 3v1m0 16v1m9-9h-1M4 12H3m15.364 6.364l-.707-.707M6.343 6.343l-.707-.707m12.728 0l-.707.707M6.343 17.657l-.707.707M16 12a4 4 0 11-8 0 4 4 0 018 0z" />
                            </svg>
                        </div>
                        <h3 class="intro-card-title">
                            {{ $t('home.creativity_title',) }}
                        </h3>
                        <p class="intro-card-text">
                            {{ $t('home.buy',) }}
                        </p>
                    </div>
                    <div class="intro-card">
                        <div class="intro-icon">
                            <svg
                                class="icon-svg"
                                viewBox="0 0 24 24"
                                fill="none"
                                stroke="currentColor"
                                stroke-width="1.5"
                            >
                                <path d="M13 10V3L4 14h7v7l9-11h-7z" />
                            </svg>
                        </div>
                        <h3 class="intro-card-title">
                            {{ $t('home.development_title',) }}
                        </h3>
                        <p class="intro-card-text">
                            {{ $t('home.development_text',) }}
                        </p>
                    </div>
                </div>
            </div>
        </div>

        <!-- Parallax Section 1 -->
        <div class="parallax-section parallax-1">
            <div class="parallax-overlay">
                <div class="parallax-content">
                    <h2 class="parallax-title">
                        {{ $t('home.draw',) }}
                    </h2>
                    <p class="parallax-text">
                        {{ $t('home.inspiration_text',) }}
                    </p>
                </div>
            </div>
        </div>

        <!-- Philosophy Section -->
        <div class="content-section philosophy-section">
            <div class="container">
                <div class="philosophy-content">
                    <div class="philosophy-text">
                        <h2 class="section-title">
                            {{ $t('home.philosophy_title',) }}
                        </h2>
                        <p class="philosophy-quote">
                            {{ $t('home.world',) }}
                        </p>
                        <p class="philosophy-description">
                            {{ $t('home.philosophy_description',) }}
                        </p>
                    </div>
                    <div class="philosophy-visual">
                        <div class="visual-element visual-1">
                            <div class="visual-content">
                                <div class="visual-icon">
                                    <svg
                                        class="icon-svg"
                                        viewBox="0 0 24 24"
                                        fill="none"
                                        stroke="currentColor"
                                        stroke-width="1.5"
                                    >
                                        <circle
                                            cx="12"
                                            cy="12"
                                            r="10"
                                        />
                                        <path d="M12 6v6l4 2" />
                                    </svg>
                                </div>
                                <h3 class="visual-title">
                                    {{ $t('home.expression_title',) }}
                                </h3>
                                <p class="visual-text">
                                    {{ $t('home.expression_text',) }}
                                </p>
                            </div>
                        </div>
                        <div class="visual-element visual-2">
                            <div class="visual-content">
                                <div class="visual-icon">
                                    <svg
                                        class="icon-svg"
                                        viewBox="0 0 24 24"
                                        fill="none"
                                        stroke="currentColor"
                                        stroke-width="1.5"
                                    >
                                        <path d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z" />
                                    </svg>
                                </div>
                                <h3 class="visual-title">
                                    {{ $t('home.harmony_title',) }}
                                </h3>
                                <p class="visual-text">
                                    {{ $t('home.harmony_text',) }}
                                </p>
                            </div>
                        </div>
                        <div class="visual-element visual-3">
                            <div class="visual-content">
                                <div class="visual-icon">
                                    <svg
                                        class="icon-svg"
                                        viewBox="0 0 24 24"
                                        fill="none"
                                        stroke="currentColor"
                                        stroke-width="1.5"
                                    >
                                        <path d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z" />
                                    </svg>
                                </div>
                                <h3 class="visual-title">
                                    {{ $t('home.growth_title',) }}
                                </h3>
                                <p class="visual-text">
                                    {{ $t('home.growth_text',) }}
                                </p>
                            </div>
                        </div>
                    </div>
                </div>
            </div>
        </div>

        <!-- Parallax Section 2 -->
        <div class="parallax-section parallax-2">
            <div class="parallax-overlay">
                <div class="parallax-content">
                    <h2 class="parallax-title">
                        {{ $t('home.givemore',) }}
                    </h2>
                    <p class="parallax-text">
                        {{ $t('home.sharing_text',) }}
                    </p>
                </div>
            </div>
        </div>

        <!-- Stats Section -->
        <StatsSection />

        <!-- CTA Section -->
        <div class="cta-section">
            <div class="container">
                <div class="cta-content">
                    <h2 class="cta-title">
                        {{ $t('home.start',) }}
                    </h2>
                    <p class="cta-text">
                        {{ $t('home.cta_text',) }}
                    </p>
                    <div class="cta-buttons">
                        <a
                            href="/all-works"
                            class="btn btn-primary"
                        >{{
                            $t('home.view_works',)
                        }}</a>
                        <a
                            href="/news"
                            class="btn btn-secondary"
                        >{{
                            $t('home.read_news',)
                        }}</a>
                        <a
                            href="/services"
                            class="btn btn-outline"
                        >{{
                            $t('home.services',)
                        }}</a>
                    </div>
                </div>
            </div>
        </div>
    </div>
</template>

<style scoped>
.home-page {
    min-height: 100vh;
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
    display: inline-flex;
    align-items: center;
    justify-content: center;
    width: 80px;
    height: 80px;
    border-radius: 50%;
    background: linear-gradient(135deg, rgba(75, 158, 144, 0.1), rgba(115, 209, 190, 0.1));
    margin-left: auto;
    margin-right: auto;
}

.icon-svg {
    width: 40px;
    height: 40px;
    color: var(--color-primary, #4B9E90);
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
    margin-bottom: 0.6rem;
    display: flex;
    align-items: center;
    justify-content: center;
    width: 48px;
    height: 48px;
    border-radius: 50%;
    background: rgba(255, 255, 255, 0.2);
}

.visual-icon .icon-svg {
    width: 28px;
    height: 28px;
    color: white;
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
@keyframes fadeIn {
    0% {
        opacity: 0;
    }
    100% {
        opacity: 1;
    }
}

/* Responsive */
@media (max-width: 767px) {
    .parallax-section {
        background-attachment: scroll;
        min-height: 50vh;
    }

    .parallax-1,
    .parallax-2 {
        background-size: cover;
        background-position: center center;
        background-attachment: scroll;
    }
}

@media (max-width: 479px) {
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
}

@media (max-width: 320px) {
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

    .btn {
        max-width: 240px;
        padding: 0.7rem 1rem;
        font-size: 0.9rem;
        min-height: 44px;
    }

    .cta-title {
        font-size: 1.6rem;
    }

    .cta-text {
        font-size: 0.9rem;
        padding: 0 0.3rem;
    }
}

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
}

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

    .content-section,
    .parallax-section {
        overflow-x: hidden;
    }
}

@media (prefers-reduced-motion: no-preference) {
    html {
        scroll-behavior: smooth;
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

    .scroll-indicator {
        animation: none !important;
    }
}

.container {
    padding-left: clamp(1rem, 3vw, 2rem);
    padding-right: clamp(1rem, 3vw, 2rem);
}

.content-section {
    padding-top: clamp(3rem, 6vw, 8rem);
    padding-bottom: clamp(3rem, 6vw, 8rem);
}

.cta-section {
    padding-top: clamp(4rem, 8vw, 10rem);
    padding-bottom: clamp(4rem, 8vw, 10rem);
}

.parallax-1,
.parallax-2 {
    aspect-ratio: 16/9;
}

@media (max-width: 768px) {
    .parallax-1,
    .parallax-2 {
        aspect-ratio: 4/3;
    }
}

@media (min-width: 1600px) {
    .container {
        max-width: 1500px;
    }

    .intro-card-text,
    .philosophy-description,
    .cta-text {
        font-size: 1.2rem;
        line-height: 1.8;
    }
}

.body_theme_dark .intro-card {
    background: var(--color-surface);
    border-color: rgba(255, 255, 255, 0.05);
}

.body_theme_dark .philosophy-section {
    background: linear-gradient(135deg, var(--color-surface) 0%, #2a2a2a 100%);
}

.body_theme_dark .intro-card-text,
.body_theme_dark .philosophy-quote,
.body_theme_dark .philosophy-description {
    color: var(--color-on-surface);
}
</style>
