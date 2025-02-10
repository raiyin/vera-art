/// <reference types="vite/client" />

interface ImportMetaEnv {
    readonly VITE_PAGE_SIZE: number;
    // more env variables...
}

interface ImportMeta {
    readonly env: ImportMetaEnv;
}
