<script>
import Carousel from "@/components/Carousel.vue";

export default {
  components: {
    Carousel,
  },
  props: {
    imageObject: {
      type: Object,
      required: true,
    },
  },
  data() {
    return {};
  },
  methods: {},
  computed: {
    imgIdModal() {
      return this.imageObject.id + "Modal";
    },
  },
  watch: {},
};
</script>

<template>
  <div
    class="modal fade"
    :id="imgIdModal"
    tabindex="-1"
    aria-labelledby="exampleModalLabel"
    aria-hidden="true"
  >
    <div class="modal-dialog">
      <div class="modal-content">
        <div class="modal-body">
          <Carousel :imageObject="imageObject" :imageId="imageObject.id" />
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
.modal-dialog {
  position: relative;
  display: table;
  /* This is important */
  overflow-y: auto;
  overflow-x: auto;
  width: auto;
}
</style>
