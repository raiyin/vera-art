<script>

import Header from "../components/Header.vue"
import Footer from "../components/Footer.vue"

export default {
    components: {
        Header,
        Footer,
    },
    props: [],
    data() {
        return {
        }
    },
    setup() {
    },
    methods: {
    },
    mounted() {
    },
    computed: {
    },
    watch: {
    },
}
</script>

<template>
    <div class="bgimg-1">
        <div class="caption">
            <span class="border">ПРОКРУТИТЬ ВНИЗ</span>
        </div>
    </div>

    <div style="color: #777;background-color:white;text-align:center;padding:50px 80px;text-align: justify;">
        <h3 style="text-align:center;">Паралакс демонстрация</h3>
        <p>Параллакс прокрутки веб-сайта тренд, где фоновое содержание перемещается с разной скоростью, чем содержание
            изображения во время прокрутки. Nascetur per nec posuere turpis, lectus nec libero turpis nunc at, sed
            posuere
            mollis ullamcorper libero ante lectus, blandit pellentesque a, magna turpis est sapien duis blandit
            dignissim.
            Viverra interdum mi magna mi, morbi sociis. Condimentum dui ipsum consequat morbi, curabitur aliquam pede,
            nullam vitae eu placerat eget et vehicula. Varius quisque non molestie dolor, nunc nisl dapibus vestibulum
            at,
            sodales tincidunt mauris ullamcorper, dapibus pulvinar, in in neque risus odio. Accumsan fringilla vulputate
            at
            quibusdam sociis eleifend, aenean maecenas vulputate, non id vehicula lorem mattis, ratione interdum sociis
            ornare. Suscipit proin magna cras vel, non sit platea sit, maecenas ante augue etiam maecenas, porta
            porttitor
            placerat leo.</p>
    </div>

    <div class="bgimg-2">
        <div class="caption">
            <span class="border" style="background-color:transparent;font-size:25px;color: #f7f7f7;">МЕНЬШЕ
                ВЫСОТЫ</span>
        </div>
    </div>

    <div style="position:relative;">
        <div style="color:#ddd;background-color:#282E34;text-align:center;padding:50px 80px;text-align: justify;">
            <p>Прокрутите вверх и вниз, чтобы действительно почувствовать, как работает прокрутка параллакса.</p>
        </div>
    </div>

    <div class="bgimg-3">
        <div class="caption">
            <span class="border" style="background-color:transparent;font-size:25px;color: #f7f7f7;">ПРОКРУЧИВАТЬ
                ВВЕРХ</span>
        </div>
    </div>

    <div style="position:relative;">
        <div style="color:#ddd;background-color:#282E34;text-align:center;padding:50px 80px;text-align: justify;">
            <p>Прокрутите вверх и вниз, чтобы действительно почувствовать, как работает прокрутка параллакса.</p>
        </div>
    </div>

    <div class="bgimg-1">
        <div class="caption">
            <span class="border">КРУТО!</span>
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
    background-image: url("../assets/img/home/img_parallax.jpg");
    min-height: 100%;
}

.bgimg-2 {
    background-image: url("../assets/img/home/img_parallax2.jpg");
    min-height: 400px;
}

.bgimg-3 {
    background-image: url("../assets/img/home/img_parallax3.jpg");
    min-height: 400px;
}

.caption {
    position: absolute;
    left: 0;
    top: 50%;
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

h3 {
    letter-spacing: 5px;
    text-transform: uppercase;
    font: 20px "Lato", sans-serif;
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