<script setup lang="ts">
import CalendarIcon from './IconCalendar.vue';
import SideNewsTrailerSkeleton from './SideNewsTrailerSkeleton.vue';
import type { NewsDesc } from '../types';
import type { PropType } from 'vue';
import { ref, computed, onMounted, nextTick } from 'vue';
import { useI18n } from '#imports';

const props = defineProps({
    sideNewsObject: {
        type: Object as PropType<NewsDesc>,
        required: true,
    },
});

const { locale } = useI18n();

// Reactive state
const isLoaded = ref(false);

// Computed property
const background = computed(
    () => props.sideNewsObject.dir + props.sideNewsObject.img_backfull
);

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
    setTimeout(() => {
        isLoaded.value = true;
    }, 1000);
};

// Lifecycle hooks
onMounted(() => {
    nextTick(() => {
        onLoad();
    });
});
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
                    width="6.5rem"
                    height="5rem"
                />
            </div>
            <div class="other-news-desc">
                <h6>
                    {{
                        locale === 'ru'
                            ? sideNewsObject.title_ru
                            : sideNewsObject.title_en
                    }}
                </h6>
                <div class="date">
                    <CalendarIcon />
                    <span>
                        &nbsp;{{ getHumanDate(sideNewsObject.datetime, locale) }}
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
    padding: 0.5rem;
    border-radius: 0.5rem;
    transition: all 0.3s ease;
    background-color: transparent;
}

.other-news-item:hover {
    color: var(--color-on-surface-hover);
    background-color: var(--color-surface-secondary);
    transform: translateY(-1px);
    box-shadow: 0 2px 8px rgba(0, 0, 0, 0.16);
}

.other-news-img {
    margin-right: 0.75rem;
    flex-shrink: 0;

    img {
        max-width: 6.5rem;
        width: 6.5rem;
        height: 5rem;
        object-fit: cover;
        object-position: center;
        border-radius: 0.375rem;
        transition: all 0.3s ease;
        box-shadow: 0 1px 4px rgba(0, 0, 0, 0.08);
    }

    img:hover {
        transform: scale(1.03);
        box-shadow: 0 2px 8px rgba(0, 0, 0, 0.12);
    }
}

.other-news-desc {
    display: flex;
    flex-direction: column;
    align-items: flex-start;
    line-height: 1.4;
    flex: 1;
    min-width: 0; /* Allows text truncation to work properly */
    max-width: 100%;
}

.other-news-desc > h6 {
    text-align: left;
    font-size: 0.9rem;
    font-weight: 600;
    margin-bottom: 0.2rem;
    line-height: 1.25;
    display: -webkit-box;
    -webkit-line-clamp: 2;
    line-clamp: 2;
    -webkit-box-orient: vertical;
    overflow: hidden;
    text-overflow: ellipsis;
}

.other-news-desc .date {
    display: flex;
    align-items: center;
    gap: 0.375rem;
    font-size: 0.8125rem;
    color: var(--color-on-surface-hover);
    margin-top: 0.375rem;
    padding: 0.2rem 0.6rem;
    background-color: var(--color-surface-secondary);
    border-radius: 0.875rem;
    max-width: 100%;
    transition: all 0.2s ease;
    line-height: 1.3;
    align-self: flex-start;
    box-shadow: 0 1px 4px rgba(0, 0, 0, 0.08);
}

.other-news-desc .date:hover {
    background-color: var(--color-surface-secondary-solid);
    transform: translateY(-1px);
}

.other-news-desc .date svg {
    width: 12px;
    height: 12px;
    color: var(--color-on-surface);
    flex-shrink: 0;
}

.other-news-desc .date span {
    font-weight: 500;
    white-space: nowrap;
    overflow: hidden;
    text-overflow: ellipsis;
    max-width: calc(100% - 1rem);
}

@media screen and (max-width: 400px) {
    .other-news-item {
        display: flex;
        flex-direction: column;
        align-items: start;
    }

    .other-news-img > img {
        max-width: 100%;
        width: 100%;
        height: 6rem;
        object-fit: cover;
    }

    .other-news-desc {
        margin-top: 0.75rem;
        width: 100%;
    }

    .other-news-desc > h6 {
        font-size: 1rem;
        -webkit-line-clamp: 3;
        line-clamp: 3;
    }

    .other-news-desc .date {
        font-size: 0.8rem;
        padding: 0.2rem 0.6rem;
        align-self: flex-start;
        max-width: 100%;
    }

    .other-news-desc .date span {
        max-width: calc(100% - 1.2rem);
    }
}
</style>
