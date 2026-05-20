import { defineStore, } from 'pinia';
import { ref, computed, } from 'vue';
import { useAuthStore, } from './AuthStore';
import api from '../api/auth';

export interface ProductCategory {
    id: number
    name_ru: string
    name_en: string
    slug: string
    description_ru?: string
    description_en?: string
    parent_id?: number
    sort_order: number
    is_active: boolean
    created_at: string
    updated_at: string
}

export interface Product {
    id: number
    category_id: number
    name_ru: string
    name_en: string
    slug: string
    description_ru?: string
    description_en?: string
    short_description_ru?: string
    short_description_en?: string
    price: number
    discount_price?: number
    currency: string
    type: 'course' | 'masterclass'
    status: 'draft' | 'published' | 'archived'
    difficulty?: 'beginner' | 'intermediate' | 'advanced'
    duration_days?: number
    duration_hours?: number
    total_lessons: number
    thumbnail_url?: string
    video_url?: string
    certificate_available: boolean
    tags: string[]
    prerequisites_ru?: string
    prerequisites_en?: string
    learning_outcomes_ru?: string
    learning_outcomes_en?: string
    language?: string
    is_featured: boolean
    sort_order: number
    created_at: string
    updated_at: string
    category?: ProductCategory
}

interface ProductFilters {
    category_slug?: string
    type?: 'course' | 'masterclass'
    status?: 'published'
    difficulty?: 'beginner' | 'intermediate' | 'advanced'
    language?: string
    min_price?: number
    max_price?: number
    search?: string
    is_featured?: boolean
    sort_by?: 'price' | 'created_at' | 'name' | 'popularity'
    sort_order?: 'asc' | 'desc'
    page?: number
    limit?: number
}

export const useProductStore = defineStore('productStore', () => {
    const authStore = useAuthStore();

    // Helper function to extract error message
    const extractErrorMessage = (err: unknown,): string => {
        if (typeof err === 'object' && err !== null) {
            const errorObj = err as { response?: { data?: { error?: string } } };
            return errorObj?.response?.data?.error || 'An unknown error occurred';
        }
        return 'An unknown error occurred';
    };

    // State
    const categories = ref<ProductCategory[]>([],);
    const products = ref<Product[]>([],);
    const featuredProducts = ref<Product[]>([],);
    const currentProduct = ref<Product | null>(null,);
    const currentCategory = ref<ProductCategory | null>(null,);
    const isLoading = ref(false,);
    const error = ref<string | null>(null,);

    // Computed
    const courses = computed(() =>
        products.value.filter(p => p.type === 'course',),
    );

    const masterclasses = computed(() =>
        products.value.filter(p => p.type === 'masterclass',),
    );

    const publishedProducts = computed(() =>
        products.value.filter(p => p.status === 'published',),
    );

    const activeCategories = computed(() =>
        categories.value.filter(c => c.is_active,),
    );

    // Actions
    const fetchCategories = async () => {
        try {
            isLoading.value = true;
            error.value = null;

            const apiInstance = api.getApiInstance();
            const response = await apiInstance.get('categories',);

            categories.value = response.data;
            return categories.value;
        } catch (err: unknown) {
            const errorMessage = (err as { response?: { data?: { error?: string } } })?.response?.data?.error || 'Failed to fetch categories';
            error.value = errorMessage;
            console.error('Error fetching categories:', err,);
            throw err;
        } finally {
            isLoading.value = false;
        }
    };

    const fetchCategoryBySlug = async (slug: string,) => {
        try {
            isLoading.value = true;
            error.value = null;

            const apiInstance = api.getApiInstance();
            const response = await apiInstance.get(`categories/${slug}`,);

            currentCategory.value = response.data;
            return currentCategory.value;
        } catch (err: unknown) {
            error.value = extractErrorMessage(err,) || 'Failed to fetch category';
            console.error('Error fetching category:', err,);
            throw err;
        } finally {
            isLoading.value = false;
        }
    };

    const fetchProducts = async (filters: ProductFilters = {},) => {
        try {
            isLoading.value = true;
            error.value = null;

            const apiInstance = api.getApiInstance();
            const params = new URLSearchParams();

            // Add filter params
            if (filters.category_slug) params.append('category_slug', filters.category_slug,);
            if (filters.type) params.append('type', filters.type,);
            if (filters.status) params.append('status', filters.status,);
            if (filters.difficulty) params.append('difficulty', filters.difficulty,);
            if (filters.language) params.append('language', filters.language,);
            if (filters.min_price !== undefined) params.append('min_price', filters.min_price.toString(),);
            if (filters.max_price !== undefined) params.append('max_price', filters.max_price.toString(),);
            if (filters.search) params.append('search', filters.search,);
            if (filters.is_featured !== undefined) params.append('is_featured', filters.is_featured.toString(),);
            if (filters.sort_by) params.append('sort_by', filters.sort_by,);
            if (filters.sort_order) params.append('sort_order', filters.sort_order,);
            if (filters.page) params.append('page', filters.page.toString(),);
            if (filters.limit) params.append('limit', filters.limit.toString(),);

            const response = await apiInstance.get(`products?${params.toString()}`,);

            products.value = response.data;
            return products.value;
        } catch (err: unknown) {
            error.value = extractErrorMessage(err,) || 'Failed to fetch products';
            console.error('Error fetching products:', err,);
            throw err;
        } finally {
            isLoading.value = false;
        }
    };

    const fetchProductBySlug = async (slug: string,) => {
        try {
            isLoading.value = true;
            error.value = null;

            const apiInstance = api.getApiInstance();
            const response = await apiInstance.get(`product/${slug}`,);

            currentProduct.value = response.data;
            return currentProduct.value;
        } catch (err: unknown) {
            error.value = extractErrorMessage(err,) || 'Failed to fetch product';
            console.error('Error fetching product:', err,);
            throw err;
        } finally {
            isLoading.value = false;
        }
    };

    const fetchProductById = async (id: number,) => {
        try {
            isLoading.value = true;
            error.value = null;

            const apiInstance = api.getApiInstance();
            const response = await apiInstance.get(`products/${id}`,);

            currentProduct.value = response.data;
            return currentProduct.value;
        } catch (err: unknown) {
            error.value = extractErrorMessage(err,) || 'Failed to fetch product';
            console.error('Error fetching product:', err,);
            throw err;
        } finally {
            isLoading.value = false;
        }
    };

    const fetchFeaturedProducts = async () => {
        try {
            isLoading.value = true;
            error.value = null;

            const apiInstance = api.getApiInstance();
            const response = await apiInstance.get('products', {
                params: {
                    is_featured: true,
                    status: 'published',
                    limit: 6,
                },
            },);

            featuredProducts.value = response.data;
            return featuredProducts.value;
        } catch (err: unknown) {
            error.value = extractErrorMessage(err,) || 'Failed to fetch featured products';
            console.error('Error fetching featured products:', err,);
            throw err;
        } finally {
            isLoading.value = false;
        }
    };

    const fetchProductsByCategory = async (categorySlug: string, filters: Omit<ProductFilters, 'category_slug'> = {},) => {
        return fetchProducts({
            ...filters,
            category_slug: categorySlug,
        },);
    };

    // Admin actions
    const createProduct = async (productData: Partial<Product>,) => {
        if (!authStore.isAdmin) {
            throw new Error('Insufficient permissions',);
        }

        try {
            isLoading.value = true;
            error.value = null;

            const apiInstance = api.getApiInstance();
            const response = await apiInstance.post('admin/products', productData,);

            // Add to local state
            products.value.push(response.data,);
            return response.data;
        } catch (err: unknown) {
            error.value = extractErrorMessage(err,) || 'Failed to create product';
            console.error('Error creating product:', err,);
            throw err;
        } finally {
            isLoading.value = false;
        }
    };

    const updateProduct = async (id: number, productData: Partial<Product>,) => {
        if (!authStore.isAdmin) {
            throw new Error('Insufficient permissions',);
        }

        try {
            isLoading.value = true;
            error.value = null;

            const apiInstance = api.getApiInstance();
            const response = await apiInstance.put(`admin/products/${id}`, productData,);

            // Update in local state
            const index = products.value.findIndex(p => p.id === id,);
            if (index !== -1) {
                products.value[index] = response.data;
            }

            if (currentProduct.value?.id === id) {
                currentProduct.value = response.data;
            }

            return response.data;
        } catch (err: unknown) {
            error.value = extractErrorMessage(err,) || 'Failed to update product';
            console.error('Error updating product:', err,);
            throw err;
        } finally {
            isLoading.value = false;
        }
    };

    const deleteProduct = async (id: number,) => {
        if (!authStore.isAdmin) {
            throw new Error('Insufficient permissions',);
        }

        try {
            isLoading.value = true;
            error.value = null;

            const apiInstance = api.getApiInstance();
            await apiInstance.delete(`admin/products/${id}`,);

            // Remove from local state
            products.value = products.value.filter(p => p.id !== id,);

            if (currentProduct.value?.id === id) {
                currentProduct.value = null;
            }

            return true;
        } catch (err: unknown) {
            error.value = extractErrorMessage(err,) || 'Failed to delete product';
            console.error('Error deleting product:', err,);
            throw err;
        } finally {
            isLoading.value = false;
        }
    };

    // Helper functions
    const getCategoryById = (id: number,) => {
        return categories.value.find(c => c.id === id,);
    };

    const getCategoryBySlug = (slug: string,) => {
        return categories.value.find(c => c.slug === slug,);
    };

    const getProductsByCategoryId = (categoryId: number,) => {
        return products.value.filter(p => p.category_id === categoryId,);
    };

    const clearError = () => {
        error.value = null;
    };

    const clearCurrentProduct = () => {
        currentProduct.value = null;
    };

    const clearCurrentCategory = () => {
        currentCategory.value = null;
    };

    return {
    // State
        categories,
        products,
        featuredProducts,
        currentProduct,
        currentCategory,
        isLoading,
        error,

        // Computed
        courses,
        masterclasses,
        publishedProducts,
        activeCategories,

        // Actions
        fetchCategories,
        fetchCategoryBySlug,
        fetchProducts,
        fetchProductBySlug,
        fetchProductById,
        fetchFeaturedProducts,
        fetchProductsByCategory,
        createProduct,
        updateProduct,
        deleteProduct,

        // Helper functions
        getCategoryById,
        getCategoryBySlug,
        getProductsByCategoryId,
        clearError,
        clearCurrentProduct,
        clearCurrentCategory,
    };
},);
