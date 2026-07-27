<template>
    <div class="payment-form">
        <div v-if="status === 'pending'">
            <h3>Оплата</h3>
            <p>
                Сумма к оплате: <strong>{{ formatPrice(amount,) }} ₽</strong>
            </p>
            <p>Описание: {{ description }}</p>
            <button
                class="btn btn-primary"
                @click="initiatePayment"
            >
                Перейти к оплате
            </button>
        </div>
        <div v-else-if="status === 'redirecting'">
            <p>Перенаправление на страницу оплаты...</p>
            <div class="spinner" />
        </div>
        <div v-else-if="status === 'success'">
            <div class="success-message">
                <h3>Оплата успешно завершена!</h3>
                <p>Ваш доступ к курсу активирован.</p>
                <NuxtLink
                    to="/learning/my-courses"
                    class="btn btn-primary"
                >Перейти к моим курсам</NuxtLink>
            </div>
        </div>
        <div v-else-if="status === 'error'">
            <div class="error-message">
                <h3>Ошибка оплаты</h3>
                <p>{{ errorMessage }}</p>
                <button
                    class="btn btn-secondary"
                    @click="retry"
                >
                    Повторить
                </button>
            </div>
        </div>
    </div>
</template>

<script setup lang="ts">
    import { ref, } from 'vue';

    const props = defineProps<{
        paymentId: number
        amount: number // в копейках
        description: string
        confirmationUrl?: string
    }>();

    const status = ref<'pending' | 'redirecting' | 'success' | 'error'>('pending',);
    const errorMessage = ref('',);

    function formatPrice(price: number,) {
        return (price / 100).toFixed(2,);
    }

    async function initiatePayment() {
        if (props.confirmationUrl) {
            status.value = 'redirecting';
            // В реальности здесь может быть открытие iframe или redirect
            window.location.href = props.confirmationUrl;
            return;
        }

        // Имитация успешной оплаты для демо
        try {
            const response = await $fetch(`/api/payments/${props.paymentId}/capture`, {
                method: 'POST',
            });
            if ((response as Record<string, boolean>).success) {
                status.value = 'success';
            } else {
                throw new Error('Payment capture failed',);
            }
        } catch (error) {
            status.value = 'error';
            errorMessage.value = 'Не удалось завершить оплату. Попробуйте позже.';
        }
    }

    function retry() {
        status.value = 'pending';
    }
</script>

<style scoped>
.payment-form {
    border: 1px solid #ddd;
    border-radius: 8px;
    padding: 2rem;
    max-width: 500px;
    margin: 0 auto;
}

.spinner {
    border: 4px solid #f3f3f3;
    border-top: 4px solid #3498db;
    border-radius: 50%;
    width: 40px;
    height: 40px;
    animation: spin 1s linear infinite;
    margin: 1rem auto;
}

@keyframes spin {
    0% {
        transform: rotate(0deg);
    }
    100% {
        transform: rotate(360deg);
    }
}

.success-message {
    text-align: center;
    color: #2e7d32;
}

.error-message {
    text-align: center;
    color: #d32f2f;
}
</style>
