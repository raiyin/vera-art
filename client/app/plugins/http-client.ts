import { getHttpClient, } from '~/api/http-client';

export default defineNuxtPlugin(() => {
    const httpClient = getHttpClient();

    return {
        provide: {
            httpClient,
        },
    };
},);
