<script lang="ts">
import CalendarIcon from '@/components/Icons/IconCalendar.vue';
import SideNewsTrailerSkeleton from '../UI/Skeletons/SideNewsTrailerSkeleton.vue';
import { NewsItemType } from '@/types';
import { PropType } from 'vue';

export default {
    components: {
        CalendarIcon,
        SideNewsTrailerSkeleton,
    },
    props: {
        sideNewsObject: {
            type: Object as PropType<NewsItemType>,
            required: true,
        },
    },
    data() {
        return {
            isLoaded: false,
            imagebasedir: import.meta.env.VITE_IMAGE_DIR,
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
            setTimeout(() => {
                this.isLoaded = true;
            }, 1000);
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
    mounted() {
        this.$nextTick(() => {
            this.onLoad();
        });
    },
};
</script>

<template>
    <div>
        <router-link
            :to="sideNewsObject.id.toString()"
            class="other-news-item"
            v-show="isLoaded"
        >
            <div class="other-news-img">
                <img
                    :src="background"
                    :alt="sideNewsObject.title_en"
                    width="9rem"
                    height="6rem"
                />
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
                        &nbsp;{{
                            getHumanDate(sideNewsObject.datetime, $i18n.locale)
                        }}
                    </span>
                </div>
            </div>
        </router-link>
        <SideNewsTrailerSkeleton v-show="!isLoaded" />
    </div>
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
    max-width: 9rem;
    width: 9rem;
    height: 6rem;
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
    min-width: 15rem;
    line-height: 1.3;
    text-wrap: nowrap;
}

@media screen and (max-width: 400px) {
    .other-news-item {
        display: flex;
        flex-direction: column;
        align-items: start;
    }

    .other-news-img > img {
        max-width: 9rem;
        width: 9rem;
        height: 6rem;
    }

    .other-news-desc {
        margin-top: .5rem;
}
}
</style>
