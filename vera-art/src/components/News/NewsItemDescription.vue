<script lang="ts">
import CalendarIcon from '@/components/Icons/IconCalendar.vue';
import NewsDescriptionSkeleton from '../UI/Skeletons/NewsDescriptionSkeleton.vue';

export default {
    components: {
        CalendarIcon,
        NewsDescriptionSkeleton,
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
            isLoaded: false,
        };
    },
    methods: {
        getHumanDate(inDate: string, locale: string) {
            const options = {
                year: 'numeric',
                month: 'long',
                day: 'numeric',
            } as const;
            const date = new Date(inDate);
            const stdLocale = locale === 'RUS' ? 'ru-RU' : 'en-EN';
            return date.toLocaleDateString(stdLocale, options);
        },
        onLoad() {
            this.isLoaded = true;
        },
    },
    mounted() {
        this.$nextTick(() => {
            this.onLoad();
        });
    },
};
</script>

<template>
    <div class="desc">
        <div class="title" @load="onLoad" v-show="isLoaded">
            <h2>
                {{
                    $i18n.locale === 'RUS'
                        ? newsObject.title_ru
                        : newsObject.title_en
                }}&nbsp;{{
                    $i18n.locale === 'RUS'
                        ? newsObject.subTitle_ru
                        : newsObject.subTitle_en
                }}
            </h2>
            <div class="date" v-show="!!newsObject.datetime">
                <CalendarIcon />
                <span>
                    &nbsp;{{
                        getHumanDate(newsObject.datetime, $i18n.locale)
                    }}</span
                >
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
