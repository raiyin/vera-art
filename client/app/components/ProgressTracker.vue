<template>
    <div class="progress-tracker">
        <h3 class="tracker-title">
            Прогресс обучения
        </h3>
        <div class="progress-summary">
            <div class="summary-item">
                <span class="label">Всего уроков:</span>
                <span class="value">{{ totalLessons }}</span>
            </div>
            <div class="summary-item">
                <span class="label">Пройдено:</span>
                <span class="value">{{ completedLessons }}</span>
            </div>
            <div class="summary-item">
                <span class="label">Прогресс:</span>
                <span class="value">{{ progressPercentage }}%</span>
            </div>
        </div>
        <div class="progress-bar">
            <div
                class="progress-fill"
                :style="{ width: progressPercentage + '%', }"
            />
        </div>
        <div class="lessons-list">
            <div
                v-for="lesson in lessons"
                :key="lesson.id"
                class="lesson-item"
                :class="{ completed: lesson.completed, }"
            >
                <div
                    class="lesson-checkbox"
                    @click="toggleLesson(lesson,)"
                >
                    <span v-if="lesson.completed">✓</span>
                </div>
                <div class="lesson-details">
                    <h4>{{ lesson.title_ru }}</h4>
                    <p class="lesson-duration">
                        {{ lesson.duration_minutes }} мин.
                    </p>
                    <p
                        v-if="lesson.completed"
                        class="completed-at"
                    >
                        Завершено: {{ formatDate(lesson.completed_at,) }}
                    </p>
                </div>
                <div class="lesson-actions">
                    <NuxtLink
                        v-if="!lesson.completed"
                        :to="`/learning/lesson/${lesson.id}`"
                        class="btn btn-small"
                    >
                        Начать
                    </NuxtLink>
                    <button
                        v-else
                        class="btn btn-small btn-outline"
                        @click="toggleLesson(lesson,)"
                    >
                        Отметить как непройденный
                    </button>
                </div>
            </div>
        </div>
    </div>
</template>

<script setup lang="ts">
    import { computed, } from 'vue';

    interface Lesson {
        id: number
        title_ru: string
        duration_minutes: number
        completed: boolean
        completed_at?: string
    }

    const props = defineProps<{
        lessons: Lesson[]
    }>();

    const emit = defineEmits<{
        'update-lesson': [lessonId: number, completed: boolean,]
    }>();

    const totalLessons = computed(() => props.lessons.length,);
    const completedLessons = computed(() => props.lessons.filter(l => l.completed,).length,);
    const progressPercentage = computed(() => {
        return totalLessons.value > 0
            ? Math.round((completedLessons.value / totalLessons.value) * 100,)
            : 0;
    });

    function toggleLesson(lesson: Lesson,) {
        emit('update-lesson', lesson.id, !lesson.completed,);
    }

    function formatDate(dateString?: string,) {
        if (!dateString) return '';
        const date = new Date(dateString,);
        return date.toLocaleDateString('ru-RU',);
    }
</script>

<style scoped>
.progress-tracker {
    border: 1px solid #ddd;
    border-radius: 8px;
    padding: 1.5rem;
    background: #f9f9f9;
}

.tracker-title {
    margin-top: 0;
    margin-bottom: 1rem;
}

.progress-summary {
    display: flex;
    justify-content: space-between;
    margin-bottom: 1rem;
}

.summary-item {
    display: flex;
    flex-direction: column;
    align-items: center;
}

.label {
    font-size: 0.9rem;
    color: #666;
}

.value {
    font-size: 1.2rem;
    font-weight: bold;
}

.progress-bar {
    height: 10px;
    background: #eee;
    border-radius: 5px;
    overflow: hidden;
    margin-bottom: 1.5rem;
}

.progress-fill {
    height: 100%;
    background: #4caf50;
    transition: width 0.3s;
}

.lessons-list {
    display: flex;
    flex-direction: column;
    gap: 1rem;
}

.lesson-item {
    display: flex;
    align-items: center;
    padding: 1rem;
    border: 1px solid #ddd;
    border-radius: 6px;
    background: white;
}

.lesson-item.completed {
    border-color: #4caf50;
    background: #f1f8e9;
}

.lesson-checkbox {
    width: 24px;
    height: 24px;
    border: 2px solid #ccc;
    border-radius: 50%;
    display: flex;
    align-items: center;
    justify-content: center;
    margin-right: 1rem;
    cursor: pointer;
    user-select: none;
}

.lesson-item.completed .lesson-checkbox {
    border-color: #4caf50;
    background: #4caf50;
    color: white;
}

.lesson-details {
    flex: 1;
}

.lesson-details h4 {
    margin: 0 0 0.25rem;
}

.lesson-duration {
    font-size: 0.9rem;
    color: #666;
    margin: 0;
}

.completed-at {
    font-size: 0.85rem;
    color: #888;
    margin: 0.25rem 0 0;
}

.lesson-actions {
    margin-left: auto;
}

.btn-small {
    padding: 0.25rem 0.75rem;
    font-size: 0.9rem;
}

.btn-outline {
    background: transparent;
    border: 1px solid #ccc;
    color: #666;
}
</style>
