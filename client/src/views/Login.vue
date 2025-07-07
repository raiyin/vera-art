<template>
    <main class="login-container">
        <div class="login">
            <h2>{{ $t('auth.login') }}</h2>
            <form @submit.prevent="handleLogin">
                <div class="input-pair">
                    <label> {{ $t('auth.username') }}</label>
                    <input v-model="user.username" type="text" required />
                </div>
                <div class="input-pair">
                    <label> {{ $t('auth.password') }}</label>
                    <input v-model="user.password" type="password" required />
                </div>
                <BButton type="submit">{{ $t('auth.login') }}</BButton>
                <p v-if="error" class="error">{{ error }}</p>
            </form>
            <p class="no-account">
                {{ $t('auth.noAccount') }}
                <router-link to="/register"> {{ $t('auth.register') }}</router-link>
            </p>
        </div>
    </main>
</template>

<script>
import auth from '@/services/auth';

export default {
    data() {
        return {
            user: {
                username: '',
                password: '',
            },
            error: '',
        };
    },
    methods: {
        async handleLogin() {
            try {
                await auth.login(this.user);
                this.$router.push('/protected');
            } catch (err) {
                this.error = 'Login failed. Please check your credentials.';
            }
        },
    },
};
</script>

<style scoped>
.login-container {
    display: flex;
    justify-content: center;
    align-items: center;
    height: 100%;
}

.login {
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

.no-account {
    margin-top: 10px;
    text-align: center;
}
</style>
