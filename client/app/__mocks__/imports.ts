// Mock for Nuxt auto-imports used in tests
// Resolved by vitest.config.ts alias `#imports`

export function useI18n() {
    return {
        locale: { value: 'ru', },
        t: (key: string,) => key,
    };
}

export function useRuntimeConfig() {
    return {
        public: {
            serverUrl: 'http://localhost:3000/api/',
            relWorksDir: '/content/works/',
            relSalesDir: '/content/sales/',
            limit: '20',
            contacts: {
                telegram: 'https://t.me/MilayaV',
                vk: 'https://vk.com/perczukowa',
                email: 'perczukowa@yandex.ru',
            },
        },
    };
}

export function useRouter() {
    return { push: () => {}, };
}

export function useRoute() {
    return { params: {}, };
}

export function useToast() {
    return { add: () => {}, };
}
