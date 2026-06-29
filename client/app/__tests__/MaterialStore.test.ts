import { describe, it, expect, beforeEach } from 'vitest'
import { setActivePinia, createPinia } from 'pinia'
import { useMaterialStore } from '~/stores/MaterialStore'

describe('MaterialStore', () => {
  beforeEach(() => {
    setActivePinia(createPinia())
  })

  it('starts with empty arrays', () => {
    const store = useMaterialStore()
    expect(store.materials).toEqual([])
    expect(store.bases).toEqual([])
    expect(store.isLoading).toBe(false)
    expect(store.error).toBeNull()
  })

  describe('getMaterialById', () => {
    it('returns undefined when materials is empty', () => {
      const store = useMaterialStore()
      expect(store.getMaterialById(1)).toBeUndefined()
    })

    it('finds a material by id', () => {
      const store = useMaterialStore()
      store.materials = [
        { id: 1, material_ru: 'Холст', material_en: 'Canvas' },
        { id: 2, material_ru: 'Акрил', material_en: 'Acrylic' },
      ]
      expect(store.getMaterialById(2)).toEqual({
        id: 2,
        material_ru: 'Акрил',
        material_en: 'Acrylic',
      })
    })
  })

  describe('getMaterialName', () => {
    it('returns russian name', () => {
      const store = useMaterialStore()
      store.materials = [{ id: 1, material_ru: 'Масло', material_en: 'Oil' }]
      expect(store.getMaterialName(1, 'ru')).toBe('Масло')
    })

    it('returns english name', () => {
      const store = useMaterialStore()
      store.materials = [{ id: 1, material_ru: 'Масло', material_en: 'Oil' }]
      expect(store.getMaterialName(1, 'en')).toBe('Oil')
    })

    it('returns empty string for unknown id', () => {
      const store = useMaterialStore()
      store.materials = [{ id: 1, material_ru: 'Масло', material_en: 'Oil' }]
      expect(store.getMaterialName(99, 'ru')).toBe('')
    })
  })

  describe('getMaterialNames', () => {
    it('returns names for multiple ids', () => {
      const store = useMaterialStore()
      store.materials = [
        { id: 1, material_ru: 'Холст', material_en: 'Canvas' },
        { id: 2, material_ru: 'Акрил', material_en: 'Acrylic' },
      ]
      expect(store.getMaterialNames([1, 2], 'ru')).toEqual(['Холст', 'Акрил'])
    })

    it('filters out unknown ids', () => {
      const store = useMaterialStore()
      store.materials = [{ id: 1, material_ru: 'Масло', material_en: 'Oil' }]
      expect(store.getMaterialNames([1, 99], 'ru')).toEqual(['Масло'])
    })
  })

  describe('getMaterialNamesString', () => {
    it('joins names with comma', () => {
      const store = useMaterialStore()
      store.materials = [
        { id: 1, material_ru: 'Холст', material_en: 'Canvas' },
        { id: 2, material_ru: 'Акрил', material_en: 'Acrylic' },
      ]
      expect(store.getMaterialNamesString([1, 2], 'ru')).toBe('Холст, Акрил')
    })
  })

  describe('getBaseById', () => {
    it('returns undefined when bases is empty', () => {
      const store = useMaterialStore()
      expect(store.getBaseById(1)).toBeUndefined()
    })

    it('finds a base by id', () => {
      const store = useMaterialStore()
      store.bases = [
        { id: 1, base_ru: 'Картон', base_en: 'Cardboard' },
        { id: 2, base_ru: 'Дерево', base_en: 'Wood' },
      ]
      expect(store.getBaseById(2)).toEqual({
        id: 2,
        base_ru: 'Дерево',
        base_en: 'Wood',
      })
    })
  })

  describe('getBaseName', () => {
    it('returns russian name', () => {
      const store = useMaterialStore()
      store.bases = [{ id: 1, base_ru: 'Холст', base_en: 'Canvas' }]
      expect(store.getBaseName(1, 'ru')).toBe('Холст')
    })

    it('returns english name', () => {
      const store = useMaterialStore()
      store.bases = [{ id: 1, base_ru: 'Холст', base_en: 'Canvas' }]
      expect(store.getBaseName(1, 'en')).toBe('Canvas')
    })

    it('returns empty string for unknown id', () => {
      const store = useMaterialStore()
      store.bases = [{ id: 1, base_ru: 'Холст', base_en: 'Canvas' }]
      expect(store.getBaseName(99, 'en')).toBe('')
    })
  })

  describe('fetchMaterials / fetchBases', () => {
    it('sets isLoading during fetch', async () => {
      const store = useMaterialStore()
      const fetchPromise = store.fetchMaterials()
      expect(store.isLoading).toBe(true)
      // The actual HTTP call will fail in tests — that's fine,
      // we just verify the loading state is managed
      await expect(fetchPromise).rejects.toThrow()
      expect(store.isLoading).toBe(false)
      expect(store.error).toBeTruthy()
    })
  })
})
