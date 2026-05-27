import { defineStore, } from 'pinia';
import { ref, } from 'vue';
import type { Material, Base, } from '~/types';

export const useMaterialStore = defineStore('materialStore', () => {
    const materials = ref<Material[]>([],);
    const bases = ref<Base[]>([],);
    const isLoading = ref(false,);
    const error = ref<string | null>(null,);

    const config = useRuntimeConfig();
    const SERVER_URL = config.public.serverUrl;

    async function fetchMaterials() {
        isLoading.value = true;
        error.value = null;
        try {
            const response = await fetch(`${SERVER_URL}materials`,);
            if (!response.ok) {
                throw new Error(`HTTP ${response.status}: ${response.statusText}`,);
            }
            const data = await response.json();
            materials.value = data;
        } catch (err) {
            error.value = err instanceof Error ? err.message : 'Failed to fetch materials';
            console.error('Error fetching materials:', err,);
        } finally {
            isLoading.value = false;
        }
    }

    async function fetchBases() {
        isLoading.value = true;
        error.value = null;
        try {
            const response = await fetch(`${SERVER_URL}bases`,);
            if (!response.ok) {
                throw new Error(`HTTP ${response.status}: ${response.statusText}`,);
            }
            const data = await response.json();
            bases.value = data;
        } catch (err) {
            error.value = err instanceof Error ? err.message : 'Failed to fetch bases';
            console.error('Error fetching bases:', err,);
        } finally {
            isLoading.value = false;
        }
    }

    async function fetchAll() {
        await Promise.all([fetchMaterials(), fetchBases(),],);
    }

    function getMaterialById(id: number,): Material | undefined {
        return materials.value.find(m => m.id === id,);
    }

    function getMaterialName(id: number, locale: string,): string {
        const material = getMaterialById(id,);
        if (!material) return '';
        return locale === 'ru' ? material.material_ru : material.material_en;
    }

    function getMaterialNames(ids: number[], locale: string,): string[] {
        return ids.map(id => getMaterialName(id, locale,),).filter(Boolean,);
    }

    function getMaterialNamesString(ids: number[], locale: string,): string {
        return getMaterialNames(ids, locale,).join(', ',);
    }

    function getBaseById(id: number,): Base | undefined {
        return bases.value.find(b => b.id === id,);
    }

    function getBaseName(id: number, locale: string,): string {
        const base = getBaseById(id,);
        if (!base) return '';
        return locale === 'ru' ? base.base_ru : base.base_en;
    }

    return {
        materials,
        bases,
        isLoading,
        error,
        fetchMaterials,
        fetchBases,
        fetchAll,
        getMaterialById,
        getMaterialName,
        getMaterialNames,
        getMaterialNamesString,
        getBaseById,
        getBaseName,
    };
},);
