<script setup lang="ts">
import CalendarIcon from './IconCalendar.vue';
import NewsDescriptionSkeleton from './NewsDescriptionSkeleton.vue';
import { ref, onMounted, nextTick } from 'vue';
import { useI18n } from '#imports';

const props = defineProps<{
    newsObject: {
        title_ru: string;
        title_en: string;
        subTitle_ru: string;
        subTitle_en: string;
        datetime: string;
    };
}>();

const { locale } = useI18n();

// Reactive state
const isLoaded = ref(false);

// Methods
const getHumanDate = (inDate: string, locale: string) => {
    const options = {
        year: 'numeric',
        month: 'long',
        day: 'numeric',
    } as const;
    const date = new Date(inDate);
    const stdLocale = locale === 'ru' ? 'ru-RU' : 'en-EN';
    return date.toLocaleDateString(stdLocale, options);
};

const onLoad = () => {
    isLoaded.value = true;
};

// Lifecycle hooks
onMounted(() => {
    nextTick(() => {
        onLoad();
    });
});
</script>

<template>
    <div class="desc">
        <div class="title" @load="onLoad" v-show="isLoaded">
            <h2>
                {{ locale === 'ru' ? newsObject.title_ru : newsObject.title_en }}&nbsp;{{
                    locale === 'ru' ? newsObject.subTitle_ru : newsObject.subTitle_en
                }}
            </h2>
            <div class="date" v-show="!!newsObject.datetime">
                <CalendarIcon />
                <span> &nbsp;{{ getHumanDate(newsObject.datetime, locale) }}</span>
            </div>
        </div>
        <NewsDescriptionSkeleton v-show="!isLoaded" />
    </div>
</template>

<style scoped>
.desc {
    margin-bottom: 2rem;
}

.title {
    margin-top: 2rem;
    text-align: left;
    color: var(--color-on-surface);
    vertical-align: bottom;
}

.title > span {
    padding-bottom: 0;
    vertical-align: bottom;
}
</style>
