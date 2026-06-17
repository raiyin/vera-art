import { ref, onMounted, onUnmounted, type Ref, } from 'vue';

/**
 * Composable for infinite scroll using IntersectionObserver.
 *
 * Usage:
 * ```ts
 * const { sentinelRef, loadMore, reset } = useInfiniteScroll(async (page) => {
 *     const items = await fetchItems(page, limit);
 *     return items;
 * });
 * ```
 *
 * The `sentinelRef` should be bound to a sentinel element at the bottom of the list.
 * `loadMore` is called automatically when the sentinel enters the viewport.
 */
export function useInfiniteScroll<T,>(
    fetcher: (page: number,) => Promise<T[]> | T[],
    options?: {
        /** Root margin for IntersectionObserver (default: '0px') */
        rootMargin?: string
        /** Intersection threshold (default: 0) */
        threshold?: number
        /** Whether to load the first page immediately on mount (default: true) */
        immediate?: boolean
        /** Maximum number of pages to load (default: Infinity) */
        maxPages?: number
    },
) {
    const {
        rootMargin = '0px',
        threshold = 0,
        immediate = true,
        maxPages = Infinity,
    } = options ?? {};

    const items = ref<T[]>([],) as Ref<T[]>;
    const page = ref(-1,);
    const loading = ref(false,);
    const hasMore = ref(true,);
    const error = ref<string | null>(null,);
    const sentinelRef = ref<Element | null>(null,);

    let observer: IntersectionObserver | null = null;

    const loadMore = async (): Promise<void> => {
        if (loading.value || !hasMore.value) return;
        if (page.value + 1 >= maxPages) {
            hasMore.value = false;
            return;
        }

        loading.value = true;
        error.value = null;

        try {
            const nextPage = page.value + 1;
            const newItems = await fetcher(nextPage,);

            if (!newItems || newItems.length === 0) {
                hasMore.value = false;
            } else {
                page.value = nextPage;
                items.value = [...items.value, ...newItems,] as T[];
            }
        } catch (e) {
            error.value = e instanceof Error ? e.message : 'Failed to load items';
            console.error('useInfiniteScroll error:', e,);
        } finally {
            loading.value = false;
        }
    };

    const reset = (): void => {
        items.value = [] as T[];
        page.value = -1;
        hasMore.value = true;
        error.value = null;
        loading.value = false;
    };

    onMounted(() => {
        if (immediate) {
            loadMore();
        }

        observer = new IntersectionObserver(
            (entries,) => {
                if (entries[0]?.isIntersecting) {
                    loadMore();
                }
            },
            { rootMargin, threshold, },
        );

        if (sentinelRef.value) {
            observer.observe(sentinelRef.value,);
        }
    },);

    onUnmounted(() => {
        if (observer) {
            observer.disconnect();
            observer = null;
        }
    },);

    return {
        /** Reactive array of accumulated items */
        items,
        /** Current page index (starts at -1, first load sets to 0) */
        page,
        /** Whether a fetch is in progress */
        loading,
        /** Whether there are more pages to load */
        hasMore,
        /** Last error message, if any */
        error,
        /** Ref to bind to the sentinel element at the bottom of the list */
        sentinelRef,
        /** Manually trigger loading the next page */
        loadMore,
        /** Reset all state back to initial values */
        reset,
    };
}
