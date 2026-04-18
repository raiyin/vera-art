<template>
    <div class="register-container">
        <div class="register">
            <h2>{{ $t('auth.register') }}</h2>
            <form @submit.prevent="handleRegister">
                <div class="input-pair">
                    <label>{{ $t('auth.username') }}</label>
                    <UInput v-model="user.username" type="text" required />
                </div>
                <div class="input-pair">
                    <label>{{ $t('auth.password') }}</label>
                    <UInput v-model="user.password" type="password" required />
                </div>
                <BButton type="submit">{{ $t('auth.register') }}</BButton>
                <p v-if="error" class="error">{{ error }}</p>
            </form>
            <p class="have-account">
                {{ $t('auth.haveAccount') }}
                <router-link to="/login">{{ $t('auth.login') }}</router-link>
            </p>
        </div>
    </div>
</template>

<script lang="ts">
import auth from '@/api/auth';
import type { UserPassPair } from '../types';
import { useAuthStore } from '../stores/AuthStore';

export default {
    data() {
        return {
            user: {
                username: '',
                password: '',
            } as UserPassPair,
            error: '',
        };
    },
    methods: {
        async handleRegister() {
            try {
                await auth.register(this.user);
                this.$router.push('/login');
            } catch (err: any) {
                this.error =
                    err.response?.data?.message ||
                    'Registration failed. Username may be taken.';
            }
        },
    },
};
</script>

<style scoped>
.register-container {
    display: flex;
    justify-content: center;
    align-items: center;
    height: 100%;
}

.register {
    max-width: 400px;
    margin: 0 auto;
}
h2 {
    text-align: center;
}
form {
    display: flex;
    flex-direction: column;
}

.input-pair {
    margin-bottom: 10px;
}

.input-pair label {
    display: block;
}

.input-pair input {
    width: 100%;
}

label {
    margin-bottom: 5px;
}
input {
    margin-bottom: 10px;
}

.have-account {
    margin-top: 10px;
    text-align: center;
}
</style>
