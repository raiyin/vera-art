<script>
import { useI18n } from 'vue-i18n';
export default {
    setup() {
        const { t } = useI18n({ useScope: 'global' });
    },
    components: {},
    props: [],
    data() {
        return {};
    },
    methods: {},
    computed: {},
    watch: {},
};
</script>

<template>
    <div class="bgimg-1">
        <div class="caption">
            <span class="border">
                {{ $t('home.hi') }}
            </span>
        </div>
    </div>

    <div
        style="
            margin: auto;
            color: #777;
            background-color: white;
            text-align: center;
            padding: 50px 20px;
            text-align: justify;
            width: 80vw;
        "
    >
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
            <span class="border" style="font-size: 25px; color: #f0f0f0">
                {{ $t('home.draw') }}
            </span>
        </div>
    </div>

    <div style="position: relative">
        <div
            style="
                text-align: center;
                color: #ddd;
                background-color: #282e34;
                text-align: center;
                padding: 50px 20px;
                text-align: justify;
            "
        >
            <p style="text-align: center">
                {{ $t('home.world') }}
            </p>
        </div>
    </div>

    <div class="bgimg-3">
        <div class="caption">
            <span
                class="border"
                style="background-color: transparent; font-size: 20px; color: #f7f7f7"
            >
                {{ $t('home.givemore') }}
            </span>
        </div>
    </div>

    <div style="position: relative; background-color: #282e34">
        <div
            style="
                margin: auto;
                color: #ddd;
                background-color: #282e34;
                text-align: center;
                padding: 50px 20px;
                text-align: justify;
                width: 80vw;
            "
        >
            <p>
                {{ $t('home.convinced') }}
            </p>
        </div>
    </div>

    <div class="bgimg-1">
        <div class="caption">
            <span class="border">
                {{ $t('home.start') }}
            </span>
        </div>
    </div>
</template>

<style scoped>
.bgimg-1,
.bgimg-2,
.bgimg-3 {
    position: relative;
    opacity: 0.65;
    background-attachment: fixed;
    background-position: center;
    background-repeat: no-repeat;
    background-size: cover;
}

.bgimg-1 {
    background-image: url('../assets/img/home/img_parallax.jpg');
    min-height: 100%;
}

.bgimg-2 {
    background-image: url('../assets/img/home/img_parallax2.jpg');
    min-height: 400px;
}

.bgimg-3 {
    background-image: url('../assets/img/home/img_parallax3.jpg');
    min-height: 400px;
}

.caption {
    position: absolute;
    left: 0;
    top: 70%;
    width: 100%;
    text-align: center;
    color: #000;
}

.caption span.border {
    background-color: #111;
    color: #fff;
    padding: 18px;
    font-size: 25px;
    letter-spacing: 10px;
}

@media (max-width: 641px) {
    .caption span.border {
        background-color: #111;
        color: #fff;
        padding: 18px;
        font-size: 15px;
        letter-spacing: 5px;
    }
}

h3 {
    letter-spacing: 5px;
    text-transform: uppercase;
    font: 20px 'Lato', sans-serif;
    color: #111;
}

/* Отключение параллаксной прокрутки для планшетов и телефонов */
@media only screen and (max-device-width: 1024px) {
    .bgimg-1,
    .bgimg-2,
    .bgimg-3 {
        background-attachment: scroll;
    }
}
</style>
