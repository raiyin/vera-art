<script>
import CalendarIcon from "@/components/icons/IconCalendar.vue";

export default {
  inject: ["imagebasedir"],
  components: {
    CalendarIcon,
  },
  props: {
    sideNewsObject: {
      type: Object,
      required: true,
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
};
</script>

<template>
  <router-link :to="sideNewsObject.id" class="other-news-item">
    <div class="other-news-img">
      <img :src="background" />
    </div>
    <div class="other-news-desc">
      <h6>{{ sideNewsObject.title }}&nbsp;{{ sideNewsObject.subTitle }}</h6>
      <div class="date">
        <CalendarIcon />
        <span>&nbsp;{{ sideNewsObject.datetimehuman }}</span>
      </div>
    </div>
  </router-link>
</template>

<style scoped>
.other-news-item {
  display: flex;
  text-decoration: none;
  color: black;
}

.other-news-img {
  width: 90px;
  margin-right: 1rem;
}

.other-news-img > img {
  max-width: 100%;
  max-height: 100%;
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
</style>
