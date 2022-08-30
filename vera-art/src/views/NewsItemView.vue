<script setup lang="ts">

import Header from "../components/Header.vue"
import Footer from "../components/Footer.vue"
import Slogan from "../components/Slogan.vue"
import SideNewsItem from "@/components/SideNewsItem.vue";
import NewsItemDescription from "@/components/NewsItemDescription.vue"

</script>
    
<template>

    <Header />
    <Slogan />

    <section class="container text-center px-0">
        <article>
            <div class="news-header">
                <div class="news_img">
                    <img src="../assets/img/news/2022/08/28/back_full.png" />
                </div>
                <div class="other-news">
                    <SideNewsItem/>
                    <SideNewsItem/>
                </div>
            </div>
            <NewsItemDescription/>
            <div class="news-text">
                <p>Lorem ipsum dolor sit amet, consectetur adipiscing elit. Aliquam eu bibendum urna, vel pharetra urna.
                    Etiam iaculis ligula vel lacus eleifend vestibulum. Donec mattis ligula sit amet lectus ultrices
                    facilisis. Proin iaculis vitae sapien et ultricies. Nunc mattis quam id mi suscipit, id euismod eros
                    gravida. Ut dictum mauris in massa molestie, vitae facilisis nisl cursus. Nulla odio lorem,
                    tincidunt quis odio sed, rutrum cursus nulla. Morbi facilisis, tortor vel lobortis efficitur, enim
                    sem tristique nibh, vel dapibus mi massa sed urna. Etiam volutpat laoreet elit vitae pretium.
                    Maecenas ac lobortis purus. Nulla vel arcu nisi. Pellentesque vehicula tortor non nunc blandit, id
                    imperdiet ex viverra. Sed mattis nisl lectus, quis consequat risus tincidunt non. Sed arcu odio,
                    malesuada vitae sapien vitae, bibendum tincidunt sapien. Proin mauris urna, consectetur id venenatis
                    sed, placerat eu est. Duis tempus neque id feugiat semper.</p>
                <p>Morbi ultrices vitae justo et posuere. Phasellus a sollicitudin enim. Donec porta finibus purus, sit
                    amet feugiat mauris malesuada at. Donec volutpat hendrerit dui vitae tincidunt. Cras in nisi in
                    justo scelerisque lobortis venenatis vel magna. Nulla et fermentum leo, nec feugiat purus. Duis sed
                    tempor ex.</p>
                <p>Nulla quis ex consectetur, auctor ipsum ac, ullamcorper felis. Donec eget erat sed lectus imperdiet
                    pharetra non eu lacus. Pellentesque eu diam mi. Aliquam vitae quam eleifend, cursus erat eget,
                    euismod lacus. Phasellus laoreet nunc sit amet varius faucibus. Integer molestie risus sagittis,
                    auctor lacus eget, blandit augue. Proin a est vitae purus fringilla viverra a eu metus.</p>
                <p>Nullam ullamcorper varius metus, eget lobortis arcu lacinia eget. Duis id nisl non mauris facilisis
                    efficitur in vel tellus. Sed at sem vitae enim sollicitudin tempor ut non sapien. Proin pretium nisi
                    non facilisis commodo. Duis et magna eu augue interdum pharetra at et ipsum. Fusce ipsum quam,
                    iaculis sit amet dignissim ac, hendrerit id turpis. Nam purus mauris, finibus in purus et, faucibus
                    sodales lorem. Suspendisse diam lorem, molestie vel dui ac, aliquet pulvinar enim. Maecenas varius
                    accumsan velit sit amet sagittis. Proin hendrerit nisi sit amet sollicitudin vehicula. Aliquam in
                    pretium nunc. Aliquam sed nibh ligula. Fusce efficitur condimentum nibh ut egestas. Vivamus nec
                    cursus turpis, id luctus velit.</p>
                <p>Etiam commodo magna id tincidunt aliquet. Aliquam sem nulla, bibendum at mi eu, iaculis ornare
                    turpis. Cras commodo condimentum quam, ut lacinia quam accumsan eu. Duis placerat scelerisque lectus
                    sed consequat. Class aptent taciti sociosqu ad litora torquent per conubia nostra, per inceptos
                    himenaeos. Nam vel quam ac lacus luctus ornare convallis sed eros. Sed vel facilisis felis. Donec
                    non velit arcu. Vestibulum sit amet ex nulla.</p>
                <p>Praesent convallis quis tellus vitae cursus. Pellentesque magna elit, luctus rhoncus tincidunt quis,
                    pretium ac felis. Nulla eleifend velit eget dolor finibus, eget volutpat nulla luctus. Cras laoreet
                    lacus massa, id tempus ante lacinia ut. Quisque ut tellus consectetur, commodo nulla nec, pretium
                    magna. Curabitur a euismod nunc, quis cursus turpis. Donec eu imperdiet libero. Sed libero massa,
                    suscipit quis imperdiet vitae, tempor vel massa.</p>
                <p>Nam fermentum consequat justo facilisis venenatis. Proin euismod enim a orci tincidunt, id suscipit
                    nisi ultrices. Pellentesque facilisis dui sed dui porta accumsan. Nulla facilisi. Nullam fermentum
                    interdum massa in sodales. Nulla finibus metus et dolor lobortis, id aliquam elit ultricies.
                    Phasellus iaculis ut felis eu vehicula.</p>
                <p>Maecenas auctor nunc nec posuere tempor. Proin sapien neque, commodo sit amet magna vitae, congue
                    suscipit sapien. Duis vitae ipsum quis enim placerat venenatis sed quis orci. Nunc vestibulum nibh
                    ac turpis aliquet, a auctor lacus rhoncus. Aenean sagittis id purus nec mollis. Quisque ac pretium
                    nisl, sed cursus enim. Nunc dictum odio in auctor commodo. Fusce ut laoreet magna. Curabitur
                    imperdiet auctor nibh sit amet egestas. Sed venenatis bibendum mauris et tempus. Vestibulum tempor
                    volutpat purus, non tempus ex imperdiet et. Aenean eu maximus ante. Cras quam magna, dignissim a
                    lacus id, vehicula dapibus turpis.</p>
                <p>Ut molestie lacus ut ante pulvinar molestie. Ut rutrum, odio et cursus gravida, libero ante tincidunt
                    nunc, at viverra sapien nulla a leo. Duis sapien purus, vestibulum a lacinia nec, dapibus vel leo.
                    Vivamus eget ultrices libero. Fusce tortor ex, ultricies sit amet blandit vel, pharetra a eros. Nam
                    mi lorem, feugiat non interdum sed, rutrum at ex. Sed tristique, mauris quis iaculis molestie, justo
                    eros rutrum sapien, nec hendrerit neque lectus quis ligula. Vestibulum feugiat, turpis at
                    consectetur congue, lacus eros gravida massa, id bibendum nunc nisi quis nulla.</p>
                <p>Mauris semper tortor sed volutpat ultrices. Cras arcu erat, ornare eget ullamcorper eget, hendrerit
                    eget felis. Proin sed eros at ligula pretium malesuada. Vivamus id nulla condimentum, sodales urna
                    quis, finibus lectus. Phasellus molestie, tellus ac mattis porta, nunc elit dictum orci, sit amet
                    posuere felis ipsum ac risus. Duis non libero tellus. Aliquam porta elementum accumsan. Fusce
                    vulputate nisl urna, vel tempor metus ullamcorper a. Vestibulum et odio id lectus mattis sodales eu
                    id metus. </p>
            </div>
        </article>
    </section>
    <Footer />
</template>
    
<style scoped>
.news-header {
    display: flex;
}

.other-news {
    display: flex;
    flex-direction: column;
}

.other-news-item-wrapper {
    margin-bottom: 1rem;
}

.other-news-item {
    display: flex;
    text-decoration: none;
    color: black;
}

.other-news-img {
    width: 70px;
    height: 70px;
    margin-right: 1rem;
}

.other-news-img>img {
    max-width: 100%;
    max-height: 100%;
}

.other-news-desc {
    display: flex;
    flex-direction: column;
    align-items: flex-start;
    line-height: 0.1em;
}

.news_img {
    display: flex;
    flex-direction: row;
    justify-content: left;
    margin-right: 2rem;
}

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

.news-text {
    text-align: left;
    margin-bottom: 2rem;
}
</style>
    