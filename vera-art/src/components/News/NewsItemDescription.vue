<script setup lang="ts">
    import CalendarIcon from "@/components/icons/IconCalendar.vue"
</script>
        
<template>
    <div class="desc">
        <div class="title">
            <h2>Учебный год начинается! Запрыгиваем!!!</h2>
            <div class="date">
                <CalendarIcon />
                <span> 28 августа, 2022</span>
            </div>
        </div>
    </div>
</template>
        
<style scoped>
.desc {
    margin-bottom: 2rem;
}

.title {
    margin-top: 2rem;
    text-align: left;
    color: black;
    vertical-align: bottom;
}

.title>span {
    padding-bottom: 0;
    vertical-align: bottom;
}
</style>
        