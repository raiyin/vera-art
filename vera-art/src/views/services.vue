<script>
import { useI18n } from 'vue-i18n';
export default {
    setup() {
        const { t } = useI18n({ useScope: 'global' });
    },
};
</script>

<template>
    <section class="container services">
        <p>
            {{ $t('services.p_one') }}
        </p>
        <p>
            {{ $t('services.p_two') }}
        </p>
        <p>
            {{ $t('services.p_three') }}
        </p>
    </section>
</template>

<style scoped>
.services {
    margin-top: 5rem;
    height: 960px;
}
</style>
