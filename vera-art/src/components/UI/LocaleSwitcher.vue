<script lang="ts">
import { useI18n } from 'vue-i18n';
import vSelect from 'vue-select';

export default {
    setup() {
        const { t, locale } = useI18n({ useScope: 'global' });
    },
    components: {
        vSelect,
    },
    data() {
        return {
            select: { lang: 'RUS', abbr: this.$t('locale.RUS') },
            languages: [
                { lang: 'RUS', abbr: this.$t('locale.RUS') },
                { lang: 'ENG', abbr: this.$t('locale.ENG') },
            ],
        };
    },
    methods: {
        onChange(event: { [x: string]: unknown }) {
            if (event['lang'] !== null) {
                i18n.global.locale.value = event['lang'];
            }
        },
    },
};
</script>

<template>
    <v-select
        :options="$i18n.availableLocales"
        :clearable="false"
        v-model="$i18n.locale"
        inputId="value"
    >
        {{ $i18n.locale }}
    </v-select>
</template>

<style scoped>
.v-select {
    width: 6rem;
    font-size: 1rem;
}
</style>
