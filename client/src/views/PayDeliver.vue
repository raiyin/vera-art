<script lang="ts">
export default {};
</script>

<template>
    <section class="container paydelivery">
        <p>
            {{ $t('payment.p_one') }}
        </p>
        <p>
            {{ $t('payment.p_two') }}
        </p>
    </section>
</template>

<style scoped>
.paydelivery {
    margin-top: 5rem;
    flex: 1 0 auto;
}

p {
    text-align: justify;
}
</style>
