<script lang="ts">
export default {};
</script>

<template>
    <div class="ava-wrapper">
        <img
            src="../assets/img/home/img_parallax.jpg"
            alt="Pertsukova"
            class="bgimg-1"
        />
    </div>

    <div class="bigblock">
        <h3 style="text-align: center">
            {{ $t('home.familiarity') }}
        </h3>
        <p>
            {{ $t('home.whoami') }}
        </p>
        <p>
            {{ $t('home.buy') }}
        </p>
    </div>

    <div class="bgimg-2">
        <div class="caption">
            <span class="border">
                {{ $t('home.draw') }}
            </span>
        </div>
    </div>

    <div style="position: relative">
        <div class="bigblock">
            <p style="text-align: center">
                {{ $t('home.world') }}
            </p>
        </div>
    </div>

    <div class="bgimg-3">
        <div class="caption">
            <span class="border">
                {{ $t('home.givemore') }}
            </span>
        </div>
    </div>

    <div class="bigblock">
        <p>
            {{ $t('home.convinced') }}
        </p>
    </div>

    <div class="bgimg-4">
        <div class="caption">
            <span class="border">
                {{ $t('home.start') }}
            </span>
        </div>
    </div>
</template>

<style scoped>
.bigblock {
    position: relative;
    margin: auto;
    color: var(--color-on-surface);
    text-align: center;
    padding: 5rem 2rem;
    text-align: justify;
    width: 80vw;
    font-size: 1rem;
}

.bgimg-1,
.bgimg-2,
.bgimg-3,
.bgimg-4 {
    position: relative;
    opacity: 0.65;
    background-attachment: fixed;
    background-position: center;
    background-repeat: no-repeat;
    background-size: cover;
}

.bgimg-1 {
    min-height: 100%;
    display: block;
    margin-left: auto;
    margin-right: auto;
    height: auto;
    max-width: 100%;
}

.bgimg-2 {
    background-image: url('../assets/img/home/img_parallax2.jpg');
    min-height: 60rem;
}

.bgimg-3 {
    background-image: url('../assets/img/home/img_parallax3.jpg');
    min-height: 60rem;
}

.bgimg-4 {
    background-image: url('../assets/img/home/img_parallax.jpg');
    display: block;
    min-height: 80vh;
    height: auto;
    max-width: 100%;
}

.caption {
    position: absolute;
    left: 0;
    top: 70%;
    width: 100%;
    text-align: center;
    color: var(--color-dark);
}

.caption span.border {
    background-color: var(--color-dark);
    color: var(--color-caption-border);
    padding: 1.8rem;
    font-size: 2.5rem;
    letter-spacing: 1rem;
}

@media (max-width: 641px) {
    .caption span.border {
        background-color: var(--color-dark);
        color: var(--color-light);
        padding: 1.8rem;
        font-size: 1.5rem;
        letter-spacing: 0.5rem;
    }
}

@media (max-width: 1200px) {
    .bgimg-2 {
        background-position: 80%;
    }

    .bgimg-3 {
        background-position: 30%;
    }

    .bgimg-4 {
        display: block;
        min-height: 80vw;
        width: auto;
        background-attachment: scroll;
    }
}

h3 {
    letter-spacing: 0.5rem;
    text-transform: uppercase;
    font:
        2rem 'Lato',
        sans-serif;
}

/* Отключение параллаксной прокрутки для планшетов и телефонов */
@media only screen and (max-device-width: 1024px) {
    .bgimg-1,
    .bgimg-2,
    .bgimg-3,
    .bgimg-4 {
        background-attachment: scroll;
    }
}
</style>
