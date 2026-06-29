// https://nuxt.com/docs/api/configuration/nuxt-config
export default defineNuxtConfig({
    modules: [
        '@nuxt/eslint',
        '@nuxt/ui',
        '@nuxt/fonts',
        '@nuxt/icon',
        '@nuxt/image',
        '@nuxtjs/i18n',
        '@pinia/nuxt',
    ],

    devtools: {
        enabled: true,
    },

    css: ['~/assets/css/main.css',],

    runtimeConfig: {
    // Private: Only available on the server-side
        apiSecret: process.env.API_SECRET,
        public: {
            serverUrl: process.env.NUXT_PUBLIC_SERVER_URL || 'http://localhost:8000/',
            limit: process.env.NUXT_PUBLIC_LIMIT || '9',
            relWorksDir: process.env.NUXT_PUBLIC_REL_WORKS_DIR || '/content/works/',
        },
    },

    features: {
        devLogs: false,
    },

    // routeRules: {
    //     '/': { prerender: true, },
    // },

    compatibilityDate: '2025-01-15',

    vite: {
        warmupEntry: false,
        optimizeDeps: {
            include: [
                '@vue/devtools-core',
                '@vue/devtools-kit',
            ],
        },
        server: {
            hmr: {
                protocol: 'ws',
                host: 'localhost',
                port: 24678,
                timeout: 30000,
            },
            watch: {
                usePolling: true,
            },
        },
    },

    eslint: {
        config: {
            stylistic: {
                commaDangle: 'always',
                braceStyle: '1tbs',
            },
        },
    },

    fonts: {
        families: [
            { name: 'Montserrat', provider: 'google', },
        ],
        // Disable fontsource provider to avoid network errors
        providers: {
            fontsource: false,
            bunny: false,
        },
    },

    i18n: {
        locales: [
            { code: 'ru', language: 'ru-RU', name: 'Русский', file: 'ru.json', },
            { code: 'en', language: 'en-US', name: 'English', file: 'en.json', },
        ],
        defaultLocale: 'ru',
        restructureDir: 'i18n',
        strategy: 'no_prefix',

        // Vue I18n configuration
        vueI18n: './i18n.config.ts', // optional external config file
    },
},);
