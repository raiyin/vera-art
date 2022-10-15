<script>

import { ref, onMounted, computed, toRefs, onBeforeMount } from 'vue'

// const props = defineProps<{
//     imageObject: any;
//     imageId: string
// }>();

// var imageObject = props.imageObject;
// const { imageId } = toRefs(props);
// let parts = ref()
// let imgCountGTOne = ref()
// const extension = ".jpg"

// function makeFileName(dir: string, index: number) {
//     return dir + index + extension
// }

// const imgIdtoLink = computed(() => {
//     return "#" + imageObject.id
// })

// onBeforeMount(async () => {
//     parts.value = imageObject.parts
//     imgCountGTOne.value = parts.value.length > 1
// })



export default {
    components: {
    },
    props: {
        imageObject: {
            type: Object
        },
        imageId: {
            type: String
        }
    },
    data() {
        return {
            parts: Object,
            imgCountGTOne: false,
            extension: ".jpg",
        }
    },
    setup() {
    },
    methods: {
        makeFileName(dir, index) {
            return dir + index + this.extension
        }
    },
    mounted() {
    },
    computed: {
        imgIdtoLink() {
            return "#" + this.imageObject.id
        }
    },
    beforeMount() {
        //console.log(this.imageObject.parts);
        this.parts = this.imageObject.parts;
        this.imgCountGTOne = this.parts.length > 1;
    },
    watch: {
    },
}
</script>

<template>

    <div :id=imageId class="carousel slide" data-bs-ride="false">
        <div v-if="imgCountGTOne" class="carousel-indicators">
            <template v-for="value, index in parts">
                <button type="button" :data-bs-target=imgIdtoLink :class="index === 0 ? 'active' : null"
                    :data-bs-slide-to=index :aria-current="index === 0 ? true : null"
                    :aria-label=value.name_ru></button>
            </template>

        </div>
        <div class="carousel-inner">
            <template v-for="value, index in parts">
                <div :class="index === 0 ? 'carousel-item active' : 'carousel-item'">
                    <img :src="makeFileName(imageObject.dir, index + 1)" class="d-block" alt="...">
                    <div class="carousel-caption d-none d-md-block">
                        <h5>{{ value.name_ru }}</h5>
                    </div>
                </div>
            </template>

        </div>
        <button v-if="imgCountGTOne" class="carousel-control-prev" type="button" :data-bs-target=imgIdtoLink
            data-bs-slide="prev">
            <span class="carousel-control-prev-icon" aria-hidden="true"></span>
            <span class="visually-hidden">Назад</span>
        </button>
        <button v-if="imgCountGTOne" class="carousel-control-next" type="button" :data-bs-target=imgIdtoLink
            data-bs-slide="next">
            <span class="carousel-control-next-icon" aria-hidden="true"></span>
            <span class="visually-hidden">Далее</span>
        </button>
    </div>

</template>

<style scoped>
@media (orientation: landscape) {
    img {
        max-height: 90vh;
    }
}

@media (orientation: portrait) {
    img {
        max-width: 90vw;
    }
}
</style>
