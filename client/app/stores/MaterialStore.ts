import { defineStore, } from 'pinia';
import { ref, } from 'vue';
import type { Material, } from '~/types';

export const useMaterialStore = defineStore('materialStore', () => {
    const materials = ref<Material[]>([],);
    const isLoading = ref(false,);
    const error = ref<string | null>(null,);

    const serverUrl = import.meta.env.VITE_SERVER_URL || 'http://localhost:8000/';

    async function fetchMaterials() {
        isLoading.value = true;
        error.value = null;
        try {
            const response = await fetch(`${serverUrl}materials`,);
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

    return {
        materials,
        isLoading,
        error,
        fetchMaterials,
        getMaterialById,
        getMaterialName,
        getMaterialNames,
        getMaterialNamesString,
    };
},);
