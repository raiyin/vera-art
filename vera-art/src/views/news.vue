<script>
import NewsTrailer from "@/components/News/NewsTrailer.vue";
import axios from "axios";

export default {
  inject: ["jsonserverhost"],
  components: {
    NewsTrailer,
  },
  props: [],
  data() {
    return {
      news: [],
    };
  },
  methods: {
    async fetchPosts() {
      try {
        const response = await axios.get(this.jsonserverhost + "news");
        this.news = response.data;
      } catch (e) {
        console.log("Error");
      }
    },
  },
  mounted() {
    this.fetchPosts();
  },
  computed: {},
  watch: {},
};
</script>

<template>
  <section class="container text-center main-content">
    <div class="row gx-0">
      <template v-for="newsObject in news" v-bind:key="newsObject.id">
        <div class="col d-flex justify-content-center mb-5">
          <NewsTrailer :newsObject="newsObject" />
        </div>
      </template>
    </div>
  </section>
</template>

<style scoped></style>
