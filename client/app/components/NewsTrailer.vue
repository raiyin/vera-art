<script setup lang="ts">
import CalendarIcon from './IconCalendar.vue';
import { useAuthStore } from '../../stores/AuthStore';
import { storeToRefs } from 'pinia';
import { ref, computed, onMounted } from 'vue';
import { useRouter, useI18n } from '#imports';

const props = defineProps<{
    newsObject: {
        id: string;
        dir: string;
        img_back: string;
        title_ru: string;
        title_en: string;
        subTitle_ru: string;
        subTitle_en: string;
        datetime: string;
    };
}>();

const emit = defineEmits<{
    'news-deleted': [id: string];
}>();

const authStore = useAuthStore();
const { isAuthenticated } = storeToRefs(authStore);
const router = useRouter();
const { locale } = useI18n();

// Reactive state
const isLoaded = ref(false);
const isDeleting = ref(false);
const isMounted = ref(false);

// Computed properties
const newsId = computed(() => '/news/' + props.newsObject.id);
const bgImage = computed(() => props.newsObject.dir + props.newsObject.img_back);

// Methods
const onImgLoad = () => {
    isLoaded.value = true;
};

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

const editNews = () => {
    router.push('/news/edit/' + props.newsObject.id + '/');
};

const deleteNews = async () => {
    // Better confirmation dialog
    if (!window.confirm('Вы уверены, что хотите удалить эту новость?')) {
        return;
    }

    // Set loading state
    isDeleting.value = true;

    try {
        const response = await fetch(
            `${import.meta.env.VITE_SERVER_URL}news/${props.newsObject.id}`,
            {
                method: 'DELETE',
                headers: {
                    Authorization: `Bearer ${localStorage.getItem('token')}`,
                },
            }
        );

        if (response.ok) {
            // Remove the news item from the UI
            emit('news-deleted', props.newsObject.id);
        } else {
            const errorData = await response.json();
            alert(
                `Ошибка при удалении новости: ${errorData.error || 'Неизвестная ошибка'}`
            );
        }
    } catch (error) {
        console.error('Error deleting news:', error);
        alert('Ошибка при удалении новости:_network_error');
    } finally {
        // Reset loading state
        isDeleting.value = false;
    }
};

// Lifecycle hooks
onMounted(() => {
    isMounted.value = true;
});
</script>

<template>
    <div class="news-item" :class="{ loading: !isLoaded }">
        <div class="img-holder">
            <router-link :to="newsId">
                <img :src="bgImage" @load="onImgLoad" v-show="isLoaded" />
                <div v-show="!isLoaded" class="image-stub" />
            </router-link>
        </div>

        <div class="news-content">
            <div v-show="isLoaded">
                {{ locale === 'ru' ? newsObject.title_ru : newsObject.title_en }}
            </div>
            <div v-show="isLoaded">
                {{ locale === 'ru' ? newsObject.subTitle_ru : newsObject.subTitle_en }}
            </div>
            <div v-show="isLoaded">
                <CalendarIcon />
                <span> &nbsp;{{ getHumanDate(newsObject.datetime, locale) }} </span>
            </div>

            <!-- Skeleton placeholders for SSR consistency -->
            <div v-show="!isLoaded" class="skeleton-text"></div>
            <div
                v-show="!isLoaded"
                class="skeleton-text"
                style="width: 70%; margin-top: 0.5rem"
            ></div>
            <div
                v-show="!isLoaded"
                class="skeleton-text"
                style="width: 50%; margin-top: 0.5rem"
            >
                <CalendarIcon style="visibility: hidden" />
                <span style="visibility: hidden">&nbsp;Placeholder</span>
            </div>
        </div>

        <div class="image-control" v-if="isAuthenticated && isMounted">
            <UButton class="btn btn-secondary w-100" type="button" v-on:click="editNews">
                Редактировать
            </UButton>
            <UButton
                class="btn btn-secondary w-100"
                type="button"
                v-on:click="deleteNews"
                :disabled="isDeleting"
            >
                <span v-if="isDeleting">Удаление...</span>
                <span v-else>Удалить</span>
            </UButton>
        </div>
    </div>
</template>

<style scoped>
.news-item {
    width: 25rem;
    border-radius: 0.3rem;
    border: 0.1rem solid var(--color-surface-secondary);
    overflow: hidden;
    display: flex;
    flex-direction: column;
    margin-bottom: 2rem;
}

.img-holder {
    width: 100%;
    height: 20rem;
    margin-bottom: 1rem;
    overflow: hidden;

    img {
        width: 100%;
        height: 100%;
        object-fit: cover; /* сохраняет пропорции */
        object-position: center; /* центрирует */
        transition: transform 0.2s;
    }

    img:hover {
        transform: scale(1.1);
    }
}

.image-stub {
    height: 100%;
    width: 100%;
}

.loading .img-holder,
.loading .news-content {
    background-color: var(--skeleton-gray);
    background: linear-gradient(
            100deg,
            rgba(255, 255, 255, 0) 40%,
            rgba(255, 255, 255, 0.5) 50%,
            rgba(255, 255, 255, 0) 60%
        )
        var(--skeleton-gray);
    background-size: 200% 100%;
    background-position-x: 180%;
    animation: 1s loading ease-in-out infinite;
    border-radius: 0.4rem;
}

@keyframes loading {
    to {
        background-position-x: -20%;
    }
}

.loading .news-content {
    height: 3rem;
    width: 80%;
    margin: auto;
}

.news-item > a {
    text-decoration: none;
}

.img-holder:hover {
    cursor: pointer;
}

.news-content {
    position: relative;
    font-weight: bold;
    font-size: larger;
    text-align: left;
    color: var(--color-on-surface);
    letter-spacing: 0.05rem;
    line-height: 1.4;
    padding-left: 1rem;
}

.news-content > p > span {
    font-size: small;
    font-weight: lighter;
}

.image-control {
    display: flex;
    flex-direction: column;
    justify-content: space-between;
    gap: 0.5rem;
    margin-top: 1rem;
}

.skeleton-text {
    height: 1.2rem;
    background-color: var(--skeleton-gray);
    background: linear-gradient(
            100deg,
            rgba(255, 255, 255, 0) 40%,
            rgba(255, 255, 255, 0.5) 50%,
            rgba(255, 255, 255, 0) 60%
        )
        var(--skeleton-gray);
    background-size: 200% 100%;
    background-position-x: 180%;
    animation: 1s loading ease-in-out infinite;
    border-radius: 0.4rem;
}
</style>
