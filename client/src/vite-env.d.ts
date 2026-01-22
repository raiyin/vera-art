/// <reference types="vite/client" />

interface ImportMetaEnv {
    readonly VITE_PAGE_SIZE: number;
    readonly VITE_SERVER_URL: string;
    // more env variables...
}

interface ImportMeta {
    readonly env: ImportMetaEnv;
}
