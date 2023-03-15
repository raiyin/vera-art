<script>
import CalendarIcon from '@/components/icons/IconCalendar.vue';

export default {
    components: {
        CalendarIcon,
    },
    props: {
        newsObject: {
            type: Object,
            required: true,
        },
    },
    data() {
        return {
            news: [],
            currentNews: {},
        };
    },
};
</script>

<template>
    <div class="desc">
        <div class="title">
            <h2>{{ newsObject.title }}&nbsp;{{ newsObject.subTitle }}</h2>
            <div class="date">
                <CalendarIcon />
                <span>&nbsp;{{ newsObject.datetimehuman }}</span>
            </div>
        </div>
    </div>
</template>

<style scoped>
.desc {
    margin-bottom: 2rem;
}

.title {
    margin-top: 2rem;
    text-align: left;
    color: black;
    vertical-align: bottom;
}

.title>span {
    padding-bottom: 0;
    vertical-align: bottom;
}
</style>
