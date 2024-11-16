<script lang="ts">
import CalendarIcon from '@/components/Icons/IconCalendar.vue';

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
    methods: {
        getHumanDate(inDate: string, locale: string) {
            if (!inDate || !locale) {
                return '';
            }

            const options = {
                year: 'numeric',
                month: 'long',
                day: 'numeric',
            } as const;
            const date = new Date(inDate);
            const stdLocale = locale === 'RUS' ? 'ru-RU' : 'en-EN';
            return date.toLocaleDateString(stdLocale, options);
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
            <div class="date">
                <CalendarIcon />
                <span
                    >&nbsp;{{
                        getHumanDate(newsObject.datetime, $i18n.locale)
                    }}</span
                >
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
    color: var(--color-on-surface);
    vertical-align: bottom;
}

.title > span {
    padding-bottom: 0;
    vertical-align: bottom;
}
</style>
