<template>
    <div class="my-courses-page">
        <div class="container">
            <h1 class="page-title">Мои курсы</h1>

            <div v-if="loading" class="loading">
                <p>Загрузка...</p>
            </div>

            <div v-else-if="courses.length === 0" class="empty-courses">
                <p>У вас пока нет активных курсов.</p>
                <NuxtLink to="/courses" class="btn btn-primary"
                    >Перейти к каталогу</NuxtLink
                >
            </div>

            <div v-else class="courses-grid">
                <div v-for="course in courses" :key="course.id" class="course-card">
                    <div class="course-image">
                        <img
                            :src="course.thumbnail_url || '/placeholder.jpg'"
                            :alt="course.title_ru"
                        />
                    </div>
                    <div class="course-info">
                        <h3>{{ course.title_ru }}</h3>
                        <p class="course-type">
                            {{ course.type === 'course' ? 'Курс' : 'Мастер-класс' }}
                        </p>
                        <p class="course-progress">
                            Прогресс: {{ course.completed_lessons }}/{{
                                course.total_lessons
                            }}
                            уроков
                        </p>
                        <div class="progress-bar">
                            <div
                                class="progress-fill"
                                :style="{ width: course.progress_percentage + '%' }"
                            ></div>
                        </div>
                        <div class="course-actions">
                            <NuxtLink
                                :to="`/learning/course/${course.id}`"
                                class="btn btn-primary"
                            >
                                Продолжить обучение
                            </NuxtLink>
                            <NuxtLink
                                :to="`/learning/progress/${course.purchase_id}`"
                                class="btn btn-secondary"
                            >
                                Детали прогресса
                            </NuxtLink>
                        </div>
                    </div>
                </div>
            </div>
        </div>
    </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue';

interface Course {
    id: number;
    title_ru: string;
    type: 'course' | 'masterclass';
    thumbnail_url?: string;
    total_lessons: number;
    completed_lessons: number;
    progress_percentage: number;
    purchase_id: number;
}

const courses = ref<Course[]>([]);
const loading = ref(true);

onMounted(async () => {
    await loadCourses();
});

async function loadCourses() {
    try {
        const response = (await $fetch('/api/learning/my-courses')) as any;
        courses.value = response.map((item: any) => ({
            id: item.product.id,
            title_ru: item.product.title_ru,
            type: item.product.type,
            thumbnail_url: item.product.thumbnail_url,
            total_lessons: item.product.total_lessons,
            completed_lessons: item.completed_lessons || 0,
            progress_percentage:
                item.total_lessons > 0
                    ? Math.round(
                          ((item.completed_lessons || 0) / item.total_lessons) * 100
                      )
                    : 0,
            purchase_id: item.purchase_id,
        }));
    } catch (error) {
        console.error('Ошибка загрузки курсов', error);
    } finally {
        loading.value = false;
    }
}
</script>

<style scoped>
.my-courses-page {
    padding: 2rem 0;
}

.page-title {
    margin-bottom: 2rem;
}

.loading,
.empty-courses {
    text-align: center;
    padding: 4rem;
}

.courses-grid {
    display: grid;
    grid-template-columns: repeat(auto-fill, minmax(300px, 1fr));
    gap: 2rem;
}

.course-card {
    border: 1px solid #ddd;
    border-radius: 8px;
    overflow: hidden;
    transition: box-shadow 0.3s;
}

.course-card:hover {
    box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
}

.course-image img {
    width: 100%;
    height: 180px;
    object-fit: cover;
}

.course-info {
    padding: 1.5rem;
}

.course-type {
    color: #666;
    font-size: 0.9rem;
    margin: 0.5rem 0;
}

.course-progress {
    margin: 1rem 0 0.5rem;
}

.progress-bar {
    height: 8px;
    background: #eee;
    border-radius: 4px;
    overflow: hidden;
}

.progress-fill {
    height: 100%;
    background: #4caf50;
    transition: width 0.3s;
}

.course-actions {
    display: flex;
    gap: 0.5rem;
    margin-top: 1rem;
}

.course-actions .btn {
    flex: 1;
    text-align: center;
}
</style>
