import js from '@eslint/js';
import globals from 'globals';
import tseslint from 'typescript-eslint';
import pluginVue from 'eslint-plugin-vue';
import { defineConfig, } from 'eslint/config';

export default defineConfig([
    {
        files: ['**/*.{js,mjs,cjs,ts,mts,cts,vue}',],
        plugins: { js, },
        extends: ['js/recommended',],
        languageOptions: { globals: { ...globals.browser, ...globals.node, }, },
        rules: {
            indent: ['error', 4,],
        },
    },
    tseslint.configs.recommended,
    {
        files: ['**/*.{ts,mts,cts}',],
        rules: {
            '@stylistic/indent': ['error', 4,],
            '@stylistic/indent-binary-ops': ['error', 4,],
        },
    },
    {
        ...pluginVue.configs['flat/essential'],
        rules: {
            ...pluginVue.configs['flat/essential'].rules,
            'vue/html-indent': ['error', 4,],
            'vue/script-indent': ['error', 4, {
                baseIndent: 1,
                switchCase: 1,
                ignores: [],
            },],
        },
    },
    {
        files: ['**/*.vue',],
        languageOptions: { parserOptions: { parser: tseslint.parser, }, },
    },
],);
