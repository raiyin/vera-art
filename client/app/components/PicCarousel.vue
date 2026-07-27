<script setup lang="ts">
import type { CommonWork, } from '~/types';
import type { PropType, } from 'vue';
import { ref, computed, onBeforeMount, } from 'vue';
import { useI18n, } from '#imports';

const props = defineProps({
    imageObject: {
        type: Object as PropType<CommonWork>,
        default: {} as CommonWork,
    },
    imageId: {
        type: String,
    },
});

const { locale, t, } = useI18n();

// Reactive state
const imgCountGTOne = ref(false,);

// Computed property
const imgIdtoLink = computed(() => '#' + props.imageObject.str_id,);

// Methods
const makeFileName = (index: number,) => {
    return props.imageObject.dir + props.imageObject.images[index - 1];
};

// Lifecycle hook
onBeforeMount(() => {
    imgCountGTOne.value = props.imageObject.images.length > 1;
});
</script>

<template>
    <div
        :id="imageId"
        class="carousel slide"
        data-bs-ride="false"
    >
        <div
            v-if="imgCountGTOne"
            class="carousel-indicators"
        >
            <UButton
                v-for="index in imageObject.images.length"
                :key="index"
                type="button"
                :data-bs-target="imgIdtoLink"
                :class="{ active: index === 1, }"
                :data-bs-slide-to="index - 1"
                :aria-current="index === 1 ? true : false"
                :aria-label="imageObject.name_en"
            />
        </div>

        <div class="carousel-inner">
            <template
                v-for="index in imageObject.images.length"
                :key="index"
            >
                <div :class="index === 1 ? 'carousel-item active' : 'carousel-item'">
                    <img
                        :src="makeFileName(index,)"
                        class="d-block modal-image"
                        alt="..."
                    >

                    <div class="carousel-caption d-none d-md-block">
                        <h5>
                            {{
                                locale === 'ru'
                                    ? `${imageObject.name_ru}`
                                    : `${imageObject.name_en}`
                            }}
                        </h5>
                    </div>
                </div>
            </template>
        </div>

        <UButton
            v-if="imgCountGTOne"
            class="carousel-control-prev"
            type="button"
            :data-bs-target="imgIdtoLink"
            data-bs-slide="prev"
        >
            <span
                class="carousel-control-prev-icon"
                aria-hidden="true"
            />
            <span class="visually-hidden">
                {{ t('carousel.back',) }}
            </span>
        </UButton>

        <UButton
            v-if="imgCountGTOne"
            class="carousel-control-next"
            type="button"
            :data-bs-target="imgIdtoLink"
            data-bs-slide="next"
        >
            <span
                class="carousel-control-next-icon"
                aria-hidden="true"
            />
            <span class="visually-hidden">
                {{ t('carousel.next',) }}
            </span>
        </UButton>
    </div>
</template>
