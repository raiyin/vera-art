<template>
    <div class="cart-page">
        <div class="container">
            <h1 class="page-title">Корзина</h1>

            <div v-if="cartItems.length === 0" class="empty-cart">
                <p>Ваша корзина пуста</p>
                <NuxtLink to="/courses" class="btn btn-primary"
                    >Перейти к курсам</NuxtLink
                >
            </div>

            <div v-else class="cart-content">
                <div class="cart-items">
                    <div v-for="item in cartItems" :key="item.id" class="cart-item">
                        <div class="item-image">
                            <img
                                :src="item.thumbnail_url || '/placeholder.jpg'"
                                :alt="item.title_ru"
                            />
                        </div>
                        <div class="item-details">
                            <h3>{{ item.title_ru }}</h3>
                            <p class="item-type">
                                {{ item.type === 'course' ? 'Курс' : 'Мастер-класс' }}
                            </p>
                            <p class="item-price">{{ formatPrice(item.price) }} ₽</p>
                            <button
                                class="btn btn-remove"
                                @click="removeFromCart(item.id)"
                            >
                                Удалить
                            </button>
                        </div>
                    </div>
                </div>

                <div class="cart-summary">
                    <h3>Итого</h3>
                    <div class="summary-row">
                        <span>Стоимость:</span>
                        <span>{{ formatPrice(totalPrice) }} ₽</span>
                    </div>
                    <div class="summary-row" v-if="discount > 0">
                        <span>Скидка:</span>
                        <span class="discount">-{{ formatPrice(discount) }} ₽</span>
                    </div>
                    <div class="promo-code">
                        <input v-model="promoCode" type="text" placeholder="Промокод" />
                        <button class="btn btn-secondary" @click="applyPromoCode">
                            Применить
                        </button>
                    </div>
                    <UAlert
                        v-if="promoError"
                        :title="'Ошибка'"
                        :description="promoError"
                        icon="i-heroicons-exclamation-triangle"
                        color="error"
                        variant="outline"
                        class="mt-2"
                        @close="promoError = ''"
                    />
                    <div class="summary-row total">
                        <span>К оплате:</span>
                        <span class="total-price">{{ formatPrice(finalPrice) }} ₽</span>
                    </div>
                    <button class="btn btn-primary btn-checkout" @click="checkout">
                        Оформить заказ
                    </button>
                </div>
            </div>
        </div>
    </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue';
import { useRouter } from '#app';
import { useAuthStore } from '~/stores/AuthStore';
import { useProductStore } from '~/stores/ProductStore';

const router = useRouter();
const authStore = useAuthStore();
const productStore = useProductStore();

interface CartItem {
    id: number;
    title_ru: string;
    title_en: string;
    type: 'course' | 'masterclass';
    price: number;
    thumbnail_url?: string;
}

const cartItems = ref<CartItem[]>([]);
const promoCode = ref('');
const discount = ref(0);
const promoCodeValid = ref(false);
const promoError = ref('');

onMounted(() => {
    loadCart();
});

function loadCart() {
    // В реальности корзина может храниться в localStorage или в сторе
    const saved = localStorage.getItem('cart');
    if (saved) {
        cartItems.value = JSON.parse(saved);
    }
}

function formatPrice(price: number) {
    // цена в копейках, переводим в рубли
    return (price / 100).toFixed(2);
}

function removeFromCart(productId: number) {
    cartItems.value = cartItems.value.filter((item) => item.id !== productId);
    saveCart();
}

function saveCart() {
    localStorage.setItem('cart', JSON.stringify(cartItems.value));
}

async function applyPromoCode() {
    if (!promoCode.value.trim()) return;
    promoError.value = '';
    // Вызов API для проверки промокода
    try {
        const response = await $fetch('/api/promo-codes/validate', {
            method: 'POST',
            body: { code: promoCode.value },
        });
        const promoResponse = response as { is_valid: boolean; discount_amount?: number };
        if (promoResponse.is_valid) {
            discount.value = promoResponse.discount_amount || 0;
            promoCodeValid.value = true;
        } else {
            promoError.value = 'Промокод недействителен';
        }
    } catch (error) {
        console.error('Ошибка проверки промокода', error);
        promoError.value = 'Ошибка проверки промокода';
    }
}

const totalPrice = computed(() => {
    return cartItems.value.reduce((sum, item) => sum + item.price, 0);
});

const finalPrice = computed(() => {
    return totalPrice.value - discount.value;
});

async function checkout() {
    if (!authStore.isAuthenticated) {
        router.push('/auth/login?redirect=/learning/cart');
        return;
    }

    // Создание платежа
    try {
        const response = await $fetch('/api/payments/create', {
            method: 'POST',
            body: {
                product_id: cartItems.value[0]?.id ?? 0, // пока только один товар
                promo_code: promoCode.value || undefined,
            },
        });
        // Перенаправление на страницу оплаты ЮKassa
        const payResponse = response as { confirmation_url?: string };
        if (payResponse.confirmation_url) {
            window.location.href = payResponse.confirmation_url;
        }
    } catch (error) {
        console.error('Ошибка создания платежа', error);
    }
}
</script>

<style scoped>
.cart-page {
    padding: 2rem 0;
}

.page-title {
    margin-bottom: 2rem;
}

.empty-cart {
    text-align: center;
    padding: 4rem;
}

.cart-content {
    display: grid;
    grid-template-columns: 2fr 1fr;
    gap: 2rem;
}

.cart-items {
    display: flex;
    flex-direction: column;
    gap: 1rem;
}

.cart-item {
    display: flex;
    border: 1px solid #ddd;
    border-radius: 8px;
    padding: 1rem;
    gap: 1rem;
}

.item-image img {
    width: 120px;
    height: 80px;
    object-fit: cover;
    border-radius: 4px;
}

.item-details {
    flex: 1;
}

.item-type {
    color: #666;
    font-size: 0.9rem;
}

.item-price {
    font-weight: bold;
    font-size: 1.2rem;
    margin: 0.5rem 0;
}

.btn-remove {
    background: #f44336;
    color: white;
    border: none;
    padding: 0.5rem 1rem;
    border-radius: 4px;
    cursor: pointer;
}

.cart-summary {
    border: 1px solid #ddd;
    border-radius: 8px;
    padding: 1.5rem;
    background: #f9f9f9;
}

.summary-row {
    display: flex;
    justify-content: space-between;
    margin-bottom: 0.5rem;
}

.total {
    border-top: 1px solid #ddd;
    padding-top: 0.5rem;
    margin-top: 0.5rem;
    font-size: 1.2rem;
    font-weight: bold;
}

.total-price {
    color: #2e7d32;
}

.promo-code {
    display: flex;
    gap: 0.5rem;
    margin: 1rem 0;
}

.promo-code input {
    flex: 1;
    padding: 0.5rem;
    border: 1px solid #ccc;
    border-radius: 4px;
}

.btn-checkout {
    width: 100%;
    margin-top: 1rem;
    padding: 1rem;
    font-size: 1.1rem;
}
</style>
