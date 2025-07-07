<template>
    <div class="admin-container container">
        <h1 class="title">{{ $t('admin.title') }}</h1>
        <div class="buttons-wrapper">
            <router-link to="/add-news" class="action-button news-button">
                <div class="button-content">
                    <svg
                        class="button-icon"
                        xmlns="http://www.w3.org/2000/svg"
                        viewBox="0 0 24 24"
                        fill="none"
                        stroke="currentColor"
                    >
                        <path
                            stroke-linecap="round"
                            stroke-linejoin="round"
                            stroke-width="2"
                            d="M19 20H5a2 2 0 01-2-2V6a2 2 0 012-2h10a2 2 0 012 2v1m2 13a2 2 0 01-2-2V7m2 13a2 2 0 002-2V9a2 2 0 00-2-2h-2m-4-3H9M7 16h6M7 8h6v4H7V8z"
                        />
                    </svg>
                    <span class="button-title">{{ $t('admin.add_news') }}</span>
                </div>
            </router-link>

            <router-link to="/add-shop-item" class="action-button shop-button">
                <div class="button-content">
                    <svg
                        class="button-icon"
                        xmlns="http://www.w3.org/2000/svg"
                        viewBox="0 0 24 24"
                        fill="none"
                        stroke="currentColor"
                    >
                        <path
                            stroke-linecap="round"
                            stroke-linejoin="round"
                            stroke-width="2"
                            d="M3 3h2l.4 2M7 13h10l4-8H5.4M7 13L5.4 5M7 13l-2.293 2.293c-.63.63-.184 1.707.707 1.707H17m0 0a2 2 0 100 4 2 2 0 000-4zm-8 2a2 2 0 11-4 0 2 2 0 014 0z"
                        />
                    </svg>
                    <span class="button-title">{{ $t('admin.add_work_to_shop') }}</span>
                </div>
            </router-link>

            <router-link to="/add-gallery-item" class="action-button gallery-button">
                <div class="button-content">
                    <svg
                        class="button-icon"
                        xmlns="http://www.w3.org/2000/svg"
                        viewBox="0 0 24 24"
                        fill="none"
                        stroke="currentColor"
                    >
                        <path
                            stroke-linecap="round"
                            stroke-linejoin="round"
                            stroke-width="2"
                            d="M4 16l4.586-4.586a2 2 0 012.828 0L16 16m-2-2l1.586-1.586a2 2 0 012.828 0L20 14m-6-6h.01M6 20h12a2 2 0 002-2V6a2 2 0 00-2-2H6a2 2 0 00-2 2v12a2 2 0 002 2z"
                        />
                    </svg>
                    <span class="button-title">{{
                        $t('admin.add_work_to_gallery')
                    }}</span>
                </div>
            </router-link>
        </div>
    </div>
</template>

<script>
export default {
    name: 'AdminDashboard',
};
</script>

<style scoped>
.admin-container {
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
    min-height: 100%;
    padding: 2rem;
}

.title {
    margin-bottom: 3rem;
    font-size: 2.5rem;
    font-weight: 600;
    text-shadow: 1px 1px 3px rgba(0, 0, 0, 0.1);
    text-transform: capitalize;
}

.buttons-wrapper {
    display: flex;
    gap: 2rem;
    flex-wrap: wrap;
    justify-content: center;
}

.button-title {
    text-transform: capitalize;
}

.action-button {
    display: inline-block;
    padding: 1.5rem 2.5rem;
    border-radius: 12px;
    text-decoration: none;
    font-weight: 600;
    font-size: 1.2rem;
    transition: all 0.3s ease;
    box-shadow: 0 4px 6px rgba(0, 0, 0, 0.1);
    width: 280px;
    text-align: center;
    position: relative;
    overflow: hidden;
}

.action-button::before {
    content: '';
    position: absolute;
    top: 0;
    left: 0;
    width: 100%;
    height: 100%;
    background: rgba(255, 255, 255, 0.1);
    transform: translateX(-100%);
    transition: transform 0.3s ease;
}

.action-button:hover::before {
    transform: translateX(0);
}

.news-button {
    background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
    color: white;
}

.shop-button {
    background: linear-gradient(135deg, #11998e 0%, #38ef7d 100%);
    color: white;
}

.gallery-button {
    background: linear-gradient(135deg, #a1c4fd 0%, #c2e9fb 100%);
    color: white;
}

.action-button:hover {
    transform: translateY(-5px);
    box-shadow: 0 10px 20px rgba(0, 0, 0, 0.2);
}

.button-content {
    display: flex;
    flex-direction: column;
    align-items: center;
    gap: 0.8rem;
}

.button-icon {
    width: 2.5rem;
    height: 2.5rem;
}

/* Responsive design */
@media (max-width: 768px) {
    .buttons-wrapper {
        flex-direction: column;
    }

    .action-button {
        width: 100%;
    }
}
</style>
