<template>
    <div class="protected">
        <h2>Protected Content</h2>
        <p v-if="content">{{ content }}</p>
        <BButton @click="logout">{{ $t('auth.logout') }}</BButton>
    </div>
</template>

<script>
import auth from '@/services/auth';

export default {
    data() {
        return {
            content: '',
        };
    },
    async created() {
        try {
            const response = await auth.getProtectedContent();
            this.content = `Hello, ${response.data.user}! ${response.data.message}`;
        } catch (err) {
            this.$router.push('/login');
        }
    },
    methods: {
        logout() {
            auth.logout();
            this.$router.push('/login');
        },
    },
};
</script>

<style>
.protected {
    display: flex;
    justify-content: center;
    align-items: center;
    flex-direction: column;
    height: 100%;
}
</style>
