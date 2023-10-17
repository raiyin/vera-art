<script lang="ts">
import CalendarIcon from '@/components/Icons/IconCalendar.vue';
import { useI18n } from 'vue-i18n';

export default {
    setup() {
        const { t, locale } = useI18n({ useScope: 'global' });
    },
    inject: ['imagebasedir'],
    components: {
        CalendarIcon,
    },
    props: {
        sideNewsObject: {
            type: Object,
            required: true,
        },
    },
    methods: {
        getHumanDate(inDate: string, locale: string) {
            const options = {
                year: 'numeric',
                month: 'long',
                day: 'numeric',
            } as const;
            const date = new Date(inDate);
            let stdLocale = locale === 'RUS' ? 'ru-RU' : 'en-EN';
            return date.toLocaleDateString(stdLocale, options);
        },
    },
    computed: {
        background() {
            return (
                this.imagebasedir +
                this.sideNewsObject.dir +
                this.sideNewsObject.img_backfull
            );
        },
    },
};
</script>

<template>
    <router-link :to="sideNewsObject.id" class="other-news-item">
        <div class="other-news-img">
            <img :src="background" alt="one more news todo" />
        </div>
        <div class="other-news-desc">
            <h6>
                {{
                    $i18n.locale === 'RUS'
                        ? sideNewsObject.title_ru
                        : sideNewsObject.title_en
                }}&nbsp;{{
                    $i18n.locale === 'RUS'
                        ? sideNewsObject.subTitle_ru
                        : sideNewsObject.subTitle_en
                }}
            </h6>
            <div class="date">
                <CalendarIcon />
                <span>
                    &nbsp;{{ getHumanDate(sideNewsObject.datetime, $i18n.locale) }}
                </span>
            </div>
        </div>
    </router-link>
</template>

<style scoped>
.other-news-item {
    display: flex;
    flex: 1 0 auto;
    text-decoration: none;
    color: var(--color-on-surface);
}
.other-news-item:hover {
    color: var(--color-on-surface-hover);
}

.other-news-img {
    margin-right: 1rem;
}

.other-news-img > img {
    max-width: 90px;
}

.other-news-desc {
    display: flex;
    flex-direction: column;
    align-items: flex-start;
    line-height: 0.1em;
}

.other-news-desc > h6 {
    text-align: left;
}

.other-news-desc .date {
    text-align: left;
    min-width: 150px;
}
</style>
