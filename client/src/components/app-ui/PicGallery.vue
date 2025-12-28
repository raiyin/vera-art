<script lang="ts">
import type { ImageProps } from '@/props/image-props';
import type { PropType } from 'vue';
import PictureCard from './PictureCard.vue';

export default {
    components: {
        Card: PictureCard,
    },
    props: {
        images: {
            type: [Array] as PropType<ImageProps[]>,
            required: true,
        },
    },
    emits: ['work-deleted'],
    data() {
        return {
            selectedWorkType: 'all',
        };
    },
    computed: {
        filteredImages(): ImageProps[] {
            if (this.selectedWorkType === 'all') {
                return this.images;
            }
            return this.images.filter((image) => {
                return image.type === parseInt(this.selectedWorkType);
            });
        },
    },
    methods: {
        handleWorkDeleted(str_id: string) {
            // Emit event to parent component to update the list
            this.$emit('work-deleted', str_id);
        },
    },
};
</script>

<template>
    <section class="container text-center main-content px-0">
        <div class="work-types-selector mb-4">
            <h3 class="selector-title">Выберите тип работы</h3>
            <div class="selector-container">
                <select
                    id="work-type-select"
                    name="work-type-choice"
                    class="form-control work-type-dropdown drop-down-arrow"
                    v-model="selectedWorkType"
                >
                    <option value="all">Все работы</option>
                    <option value="1">Картины</option>
                    <option value="2">Иллюстрации</option>
                    <option value="3">3D работы</option>
                </select>
            </div>
        </div>

        <div class="row">
            <Card
                v-for="imgObject in filteredImages"
                :imageObject="imgObject"
                :key="imgObject.id"
                @work-deleted="handleWorkDeleted"
            />
        </div>
    </section>
</template>

<style scoped>
.work-types-selector {
    display: flex;
    padding: 1rem 1rem;
    align-items: center;
    justify-content: start;
    background-color: var(--color-surface);
    border-radius: 0.5rem;
    box-shadow: 0 0.125rem 0.25rem rgba(0, 0, 0, 0.075);
    margin-bottom: 2rem;
    /* border: 1px solid var(--color-caption-border); */
    width: fit-content;
}

.selector-title {
    color: var(--color-on-surface);
    font-family: 'Montserrat', sans-serif;
    font-weight: 500;
    font-size: 1.3rem;
}

.selector-container {
    position: relative;
    margin-left: 1rem;
}

.work-type-dropdown {
    background-color: var(--color-surface-secondary-solid);
    color: var(--color-on-surface);
    border: 1px solid #ced4da;
    border-radius: 0.25rem;
    padding: 0.75rem;
    font-family: 'Montserrat', sans-serif;
    font-size: 1rem;
    width: 100%;
    transition: border-color 0.3s, box-shadow 0.3s;
}

.work-type-dropdown:focus {
    border-color: #4a90e2;
    outline: none;
    box-shadow: 0 0 0 3px rgba(74, 144, 226, 0.1);
}

.work-type-dropdown option {
    background-color: var(--color-surface);
    color: var(--color-on-surface);
}

.drop-down-arrow {
    background-image: url("data:image/svg+xml;charset=UTF-8,%3Csvg xmlns='http://www.w3.org/2000/svg' viewBox='0 0 24 24' fill='none' stroke='currentColor' stroke-width='2' stroke-linecap='round' stroke-linejoin='round'%3E%3Cpolyline points='6 9 12 15 18 9'%3E%3C/polyline%3E%3C/svg%3E");
    background-repeat: no-repeat;
    background-position: right 0.75rem center;
    background-size: 1rem;
    padding-right: 2.5rem; /* Make space for the arrow */
    -webkit-appearance: none;
    -moz-appearance: none;
    appearance: none;
}

@media (max-width: 576px) {
    .work-types-selector {
        padding: 1rem 0;
    }

    .selector-container {
        max-width: 100%;
        padding: 0 1rem;
    }
}
</style>
