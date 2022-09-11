<script>
import CalendarIcon from "@/components/icons/IconCalendar.vue"
import { computed } from "vue";

export default {
    components: {
        CalendarIcon
    },
    props: {
        sideNewsObject: {
            type: Object,
            required: true
        }
    },
    data() {
    },
    computed: {
    },
    //created() {
    //    this.$watch(
    //        () => this.$route.params,
    //        (toParams, previousParams) => {
    //            console.log(toParams);
    //        }
    //    )
    //},
}
</script>
        
<template>

    <div class="other-news-item-wrapper">
        <router-link :to=sideNewsObject.id class="other-news-item">
            <div class="other-news-img">
                <img :src="sideNewsObject.img_back" />
            </div>
            <div class="other-news-desc">
                <h6>{{sideNewsObject.title}}&nbsp;{{sideNewsObject.subTitle}}</h6>
                <div class="date">
                    <CalendarIcon />
                    <span>&nbsp;{{sideNewsObject.datetimehuman}}</span>
                </div>
            </div>
        </router-link>
    </div>

</template>
        
<style scoped>
.other-news-item-wrapper {
    margin-bottom: 1rem;
}

.other-news-item {
    display: flex;
    text-decoration: none;
    color: black;
}

.other-news-img {
    width: 70px;
    height: 70px;
    margin-right: 1rem;
}

.other-news-img>img {
    max-width: 100%;
    max-height: 100%;
}

.other-news-desc {
    display: flex;
    flex-direction: column;
    align-items: flex-start;
    line-height: 0.1em;
}
</style>
        