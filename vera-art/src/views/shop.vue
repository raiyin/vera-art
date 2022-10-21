<script>
import Gallery from "@/components/Gallery.vue";
import axios from "axios";
import vSelect from "vue-select";
import "vue-select/dist/vue-select.css";

export default {
  inject: ["jsonserverhost"],
  components: {
    Gallery,
    vSelect,
  },
  data() {
    return {
      images: [],
      selectedSort: "",
      sortOptions: [
        { value: "name_ru", name: "По названию" },
        { value: "year", name: "По новизне" },
        { value: "height", name: "По высоте" },
        { value: "width", name: "По ширине" },
      ],
    };
  },
  methods: {
    async fetchData() {
      try {
        const response = await axios.get(this.jsonserverhost + "sale");
        this.images = response.data;
      } catch (e) {
        alert("Error");
        console.log(e);
      }
    },
  },
  mounted() {
    this.fetchData();
  },
  computed: {},
  watch: {
    selectedSort(newValue, oldValue) {
      this.images.sort((image_first, image_second) => {
        if (typeof image_first[this.selectedSort] === "string")
          return image_first[this.selectedSort]?.localeCompare(
            image_second[this.selectedSort]
          );
        if (typeof image_first[this.selectedSort] === "number")
          return (
            image_first[this.selectedSort] - image_second[this.selectedSort]
          );
      });
    },
  },
};
</script>

<template>
  <section class="container mt-3 px-0">
    <v-select
      :options="sortOptions"
      :reduce="(item) => item.value"
      label="name"
      v-model="selectedSort"
      inputId="value"
      placeholder="Выберите из списка"
    >
      Выберите из списка
    </v-select>
  </section>
  <Gallery :images="images" />
</template>

<style scoped>
.v-select {
  width: 20rem;
  box-sizing: border-box;
}
</style>
