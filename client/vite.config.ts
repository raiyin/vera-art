import { fileURLToPath, URL } from 'node:url';
import { resolve, dirname } from 'node:path';
import { defineConfig } from 'vite';
import vue from '@vitejs/plugin-vue';
import vueJsx from '@vitejs/plugin-vue-jsx';
// import VueI18nPlugin from './node_modules/@intlify/unplugin-vue-i18n/vite';

// https://vitejs.dev/config/
export default defineConfig({
    plugins: [
        vue(),
        vueJsx(),
        // VueI18nPlugin({
        //     runtimeOnly: false,
        //     include: resolve(
        //         dirname(fileURLToPath(import.meta.url)),
        //         './src/i18n/locales/**'
        //     ),
        // }),
    ],
    resolve: {
        alias: {
            '@': fileURLToPath(new URL('./src', import.meta.url)),
        },
    },
});
