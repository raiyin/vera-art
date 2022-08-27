<script setup lang="ts">
import Card from "./Card.vue"
import { ref, onMounted, computed, toRefs } from 'vue'

const props = defineProps<{
    gallaryDir: string;
}>();


const { gallaryDir } = toRefs(props);
let imgDirs = ref()

function makeFullFileName(imgDir: string, fileName: string) {
    console.log(gallaryDir.value + "-" + imgDir);
    return gallaryDir.value + imgDir + "/" + fileName
}

onMounted(async () => {
    const res = await fetch(gallaryDir.value + "/total.json")
    const total = await res.json()
    imgDirs.value = total.dirs
})

</script>

<template>
    <section class="container text-center works">
        <div class="row">
            <template v-for="value in imgDirs">
                <div class="col d-flex justify-content-center">
                    <Card :jsonFile="makeFullFileName(value, 'desc.json')"
                        :filename="makeFullFileName(value, '1.jpg')" />
                </div>
            </template>
        </div>
    </section>
</template>

<style scoped>
.works {
    margin-top: 10vh;
}
</style>
