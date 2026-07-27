import { defineConfig, } from 'vitest/config';
import { resolve, } from 'path';

export default defineConfig({
    test: {
        environment: 'happy-dom',
        include: ['app/**/*.{test,spec}.{ts,tsx}',],
        exclude: ['node_modules', 'dist', '.nuxt',],
        globals: true,
        root: resolve(__dirname,),
        setupFiles: [],
    },
    resolve: {
        alias: {
            '~': resolve(__dirname, 'app',),
            '#imports': resolve(__dirname, 'app/__mocks__/imports.ts',),
        },
    },
},);
