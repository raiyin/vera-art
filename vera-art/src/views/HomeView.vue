<script>
export default {
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
            <span class="border">ЗДРАВСТВУЙТЕ!</span>
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
        <h3 style="text-align: center">Давайте знакомиться</h3>
        <p>
            Меня зовут Перцукова Вера Олеговна. Я художник, учитель, иллюстратор,
            3d-дизайнер и просто хороший человек. Я рада вас видеть на своём сайте. Здесь
            вы можете познакомиться с моим творчеством, работами, достижениями, а также
            узнать о новостях из моей профессиональной жизни.
        </p>
        <p>
            Буду рада, если вам понравится какая-либо из моих работ, и тем более, если
            решите её купить или сделать заказ картины. Надеюсь вы с удовольствием
            проведёте время здесь и захотите посещать мой сайт снова и снова. Следите за
            обновлениями.
        </p>
    </div>

    <div class="bgimg-2">
        <div class="caption">
            <span class="border" style="font-size: 25px; color: #f0f0f0"
                >ЧЕРПАТЬ У МИРА</span
            >
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
                Мир - это потрясающее место, и у меня захватывает дух только от одной
                мысли, сколько ещё предстоит узнать и сделать.
            </p>
        </div>
    </div>

    <div class="bgimg-3">
        <div class="caption">
            <span
                class="border"
                style="background-color: transparent; font-size: 20px; color: #f7f7f7"
                >И ОТДАВАТЬ БОЛЬШЕ</span
            >
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
                Я убеждена, что человек должен не зарывать свои таланты, а приумножать их.
                Поэтому стараюсь постоянно развиваться сама и передавать накоплненные
                знания своим детям и ученикам.
            </p>
        </div>
    </div>

    <div class="bgimg-1">
        <div class="caption">
            <span class="border">НАЧНЁМ!</span>
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
