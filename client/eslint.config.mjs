// @ts-check
import withNuxt from './.nuxt/eslint.config.mjs';

export default withNuxt(
    {
        files: ['**/*.{js,mjs,cjs,ts,mts,cts,vue}',],
        rules: {
            'indent': ['error', 4,],
            '@stylistic/semi': ['error', 'always',],
            '@stylistic/indent': ['error', 4,],
            '@stylistic/indent-binary-ops': ['error', 4,],
            'vue/html-indent': ['error', 4,],
            'vue/script-indent': ['error', 4, {
                baseIndent: 1,
                switchCase: 1,
                ignores: [],
            },],
        },
    },
);
