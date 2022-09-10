<script setup lang="ts">
import CalendarIcon from "@/components/icons/IconCalendar.vue"
import { computed } from "vue";

const props = defineProps<{
    sideNewsObject: any
}>();

var sideNews = props.sideNewsObject;
const newsId = computed(() => {
    return "/newsitem/" + sideNews.id
})
</script>
        
<template>

    <div class="other-news-item-wrapper">
        <router-link :to=newsId class="other-news-item">
            <div class="other-news-img">
                <img :src="sideNews.img_back" />
            </div>
            <div class="other-news-desc">
                <h6>{{sideNews.title}}&nbsp;{{sideNews.subTitle}}</h6>
                <div class="date">
                    <CalendarIcon />
                    <span>&nbsp;{{sideNews.datetimehuman}}</span>
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
        