<script lang="ts">
export default {};
</script>

<template>
    <UContainer class="paydelivery">
        <p>
            {{ $t('payment.p_one') }}
        </p>
        <p>
            {{ $t('payment.p_two') }}
        </p>
    </UContainer>
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
