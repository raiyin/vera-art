<script setup lang="ts">

import Card from "../components/Card.vue"
import Header from "../components/Header.vue"
import Footer from "../components/Footer.vue"

</script>

<template>
    <Header />

    <section class="container slogan">
        <p>Мы - то,что нас окружает...</p>
        <hr>
    </section>

    <section class="container text-center works">
        <div class="row">
            <div class="col d-flex justify-content-center">
                <Card filename="/src/assets/img/pictures/triptych/1.jpg" name="Триптих" base="Холст" material="Масло"
                    size="..." year="2019" />
            </div>
        </div>
    </section>

    <Footer />
</template>