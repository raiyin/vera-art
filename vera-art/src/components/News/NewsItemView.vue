<script>

import SideNewsTrailer from "@/components/News/SideNewsTrailer.vue";
import NewsItemDescription from "@/components/News/NewsItemDescription.vue"
import { useRouter, useRoute } from 'vue-router'
import axios from "axios";


export default {
    components: {
        SideNewsTrailer,
        NewsItemDescription
    },
    data() {
        return {
            news: [],
            currentNews: {}
        }
    },
    setup() {
        const route = useRoute();
        const router = useRouter()
    },
    methods: {
        async fetchData() {
            try {
                let newsid = this.$route.path.substring(this.$route.path.lastIndexOf('/') + 1)
                var response = await axios.get('http://localhost:3001/news', { params: { id: newsid } });
                this.currentNews = response.data[0];

                response = await axios.get('http://localhost:3001/news', { params: { id_ne: newsid, _limit: 5 } });
                this.news = response.data;
            }
            catch (e) {
                console.log(e);
                var propertyNames = Object.getOwnPropertyNames(e);
                propertyNames.forEach(function (property) {
                    var descriptor = Object.getOwnPropertyDescriptor(e, property);
                    console.log(property + ":" + e[property] + ":" + propsToStr(descriptor));
                });
            }
        },
        makeImageName(index) {
            return this.currentNews.base_dir + index + ".jpg"
        },
        makeVideoSlideLabel(index) {
            return "Видеослайд " + index
        },
        makeVideoName(index) {
            return this.currentNews.base_dir + index + ".mp4"
        }
    },
    mounted() {
        this.fetchData();
        let mdbScript = document.createElement('script')
        mdbScript.setAttribute('src', '/src/assets/js/mdb.min.js')
        document.head.appendChild(mdbScript)
    },
    async onBeforeMount() {
        await router.isReady()
    },
    computed: {
        background() {
            return this.currentNews.base_dir + this.currentNews.img_backfull;
        }
    }
}

</script>
    
<template>

    <section class="container text-center px-0 main-content">
        <article>
            <div class="news-header">
                <div class="news_img">
                    <img :src="background" />
                </div>
                <div class="other-news">
                    <template v-for="newsItem in news">
                        <SideNewsTrailer :sideNewsObject="newsItem" />
                    </template>
                </div>
            </div>
            <NewsItemDescription :newsObject="currentNews" />
            <div class="news-text">
                <p v-html="currentNews.text" />
            </div>
        </article>

        <div class="row">
            <div class="col-lg-4 col-md-12 mb-4 mb-lg-0">
                <template v-for="image_index in currentNews.imagescount">
                    <img v-if="image_index%3==1" :src=makeImageName(image_index)
                        class="w-100 shadow-1-strong rounded mb-4" />
                </template>
            </div>
            <div class="col-lg-4 col-md-12 mb-4 mb-lg-0">
                <template v-for="image_index in currentNews.imagescount">
                    <img v-if="image_index%3==2" :src=makeImageName(image_index)
                        class="w-100 shadow-1-strong rounded mb-4" />
                </template>
            </div>
            <div class="col-lg-4 col-md-12 mb-4 mb-lg-0">
                <template v-for="image_index in currentNews.imagescount">
                    <img v-if="image_index%3==0" :src=makeImageName(image_index)
                        class="w-100 shadow-1-strong rounded mb-4" />
                </template>
            </div>
        </div>

        <div id="carouselVideoExample" class="carousel slide carousel-fade" data-mdb-ride="carousel">
            <div class="carousel-indicators">
                <template v-for="video_index in currentNews.videoscount">
                    <button v-if="video_index==1" type="button" data-mdb-target="#carouselVideoExample"
                        :data-mdb-slide-to="video_index" class="active" aria-current="true"
                        :aria-label="makeVideoSlideLabel(video_index)">
                    </button>
                    <button v-if="video_index!=1" type="button" data-mdb-target="#carouselVideoExample"
                        :data-mdb-slide-to="video_index" :aria-label="makeVideoSlideLabel(video_index)">
                    </button>
                </template>
            </div>

            <div class="carousel-inner">
                <template v-for="video_index in currentNews.videoscount">
                    <div v-if="video_index==1" class="carousel-item active">
                        <video class="img-fluid" autoplay loop muted>
                            <source :src=makeVideoName(video_index) type="video/mp4" />
                        </video>
                    </div>

                    <div v-if="video_index!=1" class="carousel-item">
                        <video class="img-fluid" autoplay loop muted>
                            <source :src=makeVideoName(video_index) type="video/mp4" />
                        </video>
                    </div>
                </template>
            </div>

            <button class="carousel-control-prev" type="button" data-mdb-target="#carouselVideoExample"
                data-mdb-slide="prev">
                <span class="carousel-control-prev-icon" aria-hidden="true"></span>
                <span class="visually-hidden">Previous</span>
            </button>
            <button class="carousel-control-next" type="button" data-mdb-target="#carouselVideoExample"
                data-mdb-slide="next">
                <span class="carousel-control-next-icon" aria-hidden="true"></span>
                <span class="visually-hidden">Next</span>
            </button>
        </div>

    </section>
</template>
    
<style scoped>
@import '@/assets/css/mdb.min.css';

.news-header {
    display: flex;
}

.other-news {
    display: flex;
    flex-direction: column;
}

.news_img {
    display: flex;
    flex-direction: row;
    justify-content: left;
    margin-right: 2rem;
}

.news-text {
    text-align: left;
    margin-bottom: 2rem;
}

.news_img {
    width: 696px;
    height: 468px;
}
</style>
    