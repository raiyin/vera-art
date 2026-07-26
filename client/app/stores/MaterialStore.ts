import { defineStore, } from 'pinia';
import { ref, } from 'vue';
import { getHttpClient, } from '~/api/http-client';
import type { Material, Base, } from '~/types';

export const useMaterialStore = defineStore('materialStore', () => {
    const materials = ref<Material[]>([],);
    const bases = ref<Base[]>([],);
    const isLoading = ref(false,);
    const error = ref<string | null>(null,);

    async function fetchMaterials() {
        isLoading.value = true;
        error.value = null;
        try {
            const { data, } = await getHttpClient().get('materials',);
            materials.value = data.materials ?? [];
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
            const { data, } = await getHttpClient().get('bases',);
            bases.value = data.bases ?? [];
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
        return locale === 'ru' ? material.name_ru : material.name_en;
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
        return locale === 'ru' ? base.name_ru : base.name_en;
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
