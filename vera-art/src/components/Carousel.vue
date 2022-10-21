<script>
export default {
  components: {},
  props: {
    imageObject: {
      type: Object,
    },
    imageId: {
      type: String,
    },
  },
  data() {
    return {
      parts: Object,
      imgCountGTOne: false,
      extension: ".jpg",
    };
  },
  methods: {
    makeFileName(dir, index) {
      return dir + index + this.extension;
    },
  },
  computed: {
    imgIdtoLink() {
      return "#" + this.imageObject.id;
    },
  },
  beforeMount() {
    this.parts = this.imageObject.parts;
    this.imgCountGTOne = this.parts.length > 1;
  },
  watch: {},
};
</script>

<template>
  <div :id="imageId" class="carousel slide" data-bs-ride="false">
    <div v-if="imgCountGTOne" class="carousel-indicators">
      <template v-for="(value, index) in parts" v-bind:key="index">
        <button
          type="button"
          :data-bs-target="imgIdtoLink"
          :class="index === 0 ? 'active' : null"
          :data-bs-slide-to="index"
          :aria-current="index === 0 ? true : null"
          :aria-label="value.name_ru"
        ></button>
      </template>
    </div>
    <div class="carousel-inner">
      <template v-for="(value, index) in parts" v-bind:key="index">
        <div :class="index === 0 ? 'carousel-item active' : 'carousel-item'">
          <img
            :src="makeFileName(imageObject.dir, index + 1)"
            class="d-block"
            alt="..."
          />
          <div class="carousel-caption d-none d-md-block">
            <h5>{{ value.name_ru }}</h5>
          </div>
        </div>
      </template>
    </div>
    <button
      v-if="imgCountGTOne"
      class="carousel-control-prev"
      type="button"
      :data-bs-target="imgIdtoLink"
      data-bs-slide="prev"
    >
      <span class="carousel-control-prev-icon" aria-hidden="true"></span>
      <span class="visually-hidden">Назад</span>
    </button>
    <button
      v-if="imgCountGTOne"
      class="carousel-control-next"
      type="button"
      :data-bs-target="imgIdtoLink"
      data-bs-slide="next"
    >
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
