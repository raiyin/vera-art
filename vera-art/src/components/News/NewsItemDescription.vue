<script>
import CalendarIcon from '@/components/Icons/IconCalendar.vue';
import { useI18n } from 'vue-i18n';

export default {
    setup() {
        const { t, locale } = useI18n({ useScope: 'global' });
    },
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
        getHumanDate(inDate, locale) {
            const options = {
                year: 'numeric',
                month: 'long',
                day: 'numeric',
            };
            const date = new Date(inDate);
            return date.toLocaleDateString(locale + '-' + locale.toUpperCase(), options);
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
                    $i18n.locale === 'ru' ? newsObject.title_ru : newsObject.title_en
                }}&nbsp;{{
                    $i18n.locale === 'ru'
                        ? newsObject.subTitle_ru
                        : newsObject.subTitle_en
                }}
            </h2>
            <div class="date">
                <CalendarIcon />
                <span>&nbsp;{{ getHumanDate(newsObject.datetime, $i18n.locale) }}</span>
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

.title > span {
    padding-bottom: 0;
    vertical-align: bottom;
}
</style>
