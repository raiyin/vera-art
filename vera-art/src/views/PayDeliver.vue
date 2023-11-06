<script lang="ts">
import { useI18n } from 'vue-i18n';
export default {
    setup() {
        const { t } = useI18n({ useScope: 'global' });
    },
};
</script>

<template>
    <section class="container paydelivery">
        <p>
            {{ this.$t('payment.p_one') }}
        </p>
        <p>
            {{ this.$t('payment.p_two') }}
        </p>
    </section>
</template>

<style scoped>
.paydelivery {
    margin-top: 5rem;
    height: 960px;
}
</style>
