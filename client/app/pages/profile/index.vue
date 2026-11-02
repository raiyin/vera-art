<template>
    <div
        class="min-h-screen flex items-center justify-center bg-linear-to-br from-green-50 to-teal-50 dark:from-gray-900 dark:to-gray-800 p-4"
    >
        <div class="w-full max-w-2xl">
            <!-- Page Header -->
            <div class="text-center mb-6">
                <div
                    class="inline-flex items-center justify-center w-14 h-14 bg-teal-500 rounded-2xl shadow-lg mb-3"
                >
                    <Icon
                        name="i-heroicons-user-circle"
                        class="w-8 h-8 text-white"
                    />
                </div>
                <h1 class="text-3xl font-bold text-gray-900 dark:text-white">
                    {{ $t('profile.title',) }}
                </h1>
                <p class="text-gray-600 dark:text-gray-300 mt-2">
                    {{ $t('profile.subtitle',) }}
                </p>
            </div>

            <!-- Profile Card -->
            <UCard class="shadow-2xl border-0">
                <!-- Avatar Section -->
                <div class="flex flex-col items-center mb-8">
                    <div class="relative group">
                        <!-- Avatar Preview -->
                        <div
                            class="w-28 h-28 rounded-full overflow-hidden border-4 border-gray-200 dark:border-gray-600 bg-gray-100 dark:bg-gray-700 flex items-center justify-center"
                        >
                            <img
                                v-if="avatarUrl"
                                :src="avatarUrl"
                                alt="Avatar"
                                class="w-full h-full object-cover"
                            >
                            <Icon
                                v-else
                                name="i-heroicons-user-circle"
                                class="w-16 h-16 text-gray-400"
                            />
                        </div>

                        <!-- Upload overlay on hover -->
                        <label
                            for="avatar-upload"
                            class="absolute inset-0 flex items-center justify-center bg-black/40 rounded-full opacity-0 group-hover:opacity-100 cursor-pointer transition-opacity"
                        >
                            <Icon
                                name="i-heroicons-camera"
                                class="w-8 h-8 text-white"
                            />
                        </label>
                        <input
                            id="avatar-upload"
                            ref="avatarInputRef"
                            type="file"
                            accept="image/jpeg,image/png,image/gif,image/webp"
                            class="hidden"
                            @change="handleAvatarSelect"
                        >
                    </div>

                    <!-- Avatar info text -->
                    <p class="text-xs text-gray-500 dark:text-gray-400 mt-3 text-center">
                        {{ $t('profile.avatarMaxSize',) }}<br>
                        {{ $t('profile.avatarTypes',) }}
                    </p>

                    <!-- Avatar action buttons -->
                    <div class="flex gap-2 mt-3">
                        <UButton
                            size="sm"
                            color="primary"
                            variant="solid"
                            :loading="avatarUploading"
                            :disabled="avatarUploading"
                            @click="handleAvatarClick"
                        >
                            <template #leading>
                                <Icon
                                    name="i-heroicons-arrow-up-tray"
                                    class="w-4 h-4"
                                />
                            </template>
                            {{ $t('profile.changeAvatar',) }}
                        </UButton>

                        <UButton
                            v-if="avatarUrl"
                            size="sm"
                            color="error"
                            variant="outline"
                            :loading="avatarDeleting"
                            :disabled="avatarDeleting"
                            @click="handleAvatarDelete"
                        >
                            <template #leading>
                                <Icon
                                    name="i-heroicons-trash"
                                    class="w-4 h-4"
                                />
                            </template>
                            {{ $t('profile.deleteAvatar',) }}
                        </UButton>
                    </div>

                    <!-- Avatar error message -->
                    <UAlert
                        v-if="avatarError"
                        :title="$t('profile.error',)"
                        :description="avatarError"
                        icon="i-heroicons-exclamation-triangle"
                        color="error"
                        variant="outline"
                        class="mt-3 w-full"
                    />
                </div>

                <USeparator class="mb-6" />

                <UForm
                    :state="formState"
                    class="space-y-6"
                    @submit="handleSave"
                >
                    <!-- Username (read-only) -->
                    <UFormField
                        name="username"
                        :label="$t('auth.username',)"
                    >
                        <UInput
                            v-model="formState.username"
                            type="text"
                            icon="i-heroicons-user"
                            disabled
                            class="w-full"
                        />
                    </UFormField>

                    <!-- Email -->
                    <UFormField
                        name="email"
                        :label="$t('auth.email',)"
                    >
                        <UInput
                            v-model="formState.email"
                            type="email"
                            icon="i-heroicons-envelope"
                            :placeholder="$t('profile.emailPlaceholder',)"
                            class="w-full"
                        />
                    </UFormField>

                    <!-- Full Name -->
                    <UFormField
                        name="full_name"
                        :label="$t('profile.fullName',)"
                    >
                        <UInput
                            v-model="formState.full_name"
                            type="text"
                            icon="i-heroicons-identification"
                            :placeholder="$t('profile.fullNamePlaceholder',)"
                            class="w-full"
                        />
                    </UFormField>

                    <!-- Role (read-only) -->
                    <UFormField
                        name="role"
                        :label="$t('profile.role',)"
                    >
                        <UInput
                            v-model="formState.role"
                            type="text"
                            icon="i-heroicons-shield-check"
                            disabled
                            class="w-full"
                        />
                    </UFormField>

                    <!-- Member since (read-only) -->
                    <UFormField
                        name="created_at"
                        :label="$t('profile.memberSince',)"
                    >
                        <UInput
                            v-model="formState.created_at"
                            type="text"
                            icon="i-heroicons-calendar"
                            disabled
                            class="w-full"
                        />
                    </UFormField>

                    <!-- Success Message -->
                    <UAlert
                        v-if="success"
                        :title="$t('profile.success',)"
                        :description="success"
                        icon="i-heroicons-check-circle"
                        color="success"
                        variant="outline"
                    />

                    <!-- Error Message -->
                    <UAlert
                        v-if="error"
                        :title="$t('profile.error',)"
                        :description="error"
                        icon="i-heroicons-exclamation-triangle"
                        color="error"
                        variant="outline"
                    />

                    <!-- Action Buttons -->
                    <div class="flex gap-3 pt-2">
                        <UButton
                            type="submit"
                            :loading="saving"
                            :disabled="saving"
                            color="primary"
                            class="flex-1 justify-center"
                        >
                            <template #leading>
                                <Icon
                                    name="i-heroicons-check"
                                    class="w-5 h-5"
                                />
                            </template>
                            {{ $t('profile.save',) }}
                        </UButton>

                        <UButton
                            variant="outline"
                            color="neutral"
                            :disabled="saving"
                            class="justify-center"
                            @click="resetForm"
                        >
                            {{ $t('profile.reset',) }}
                        </UButton>
                    </div>
                </UForm>
            </UCard>

            <!-- Back to Home -->
            <div class="mt-6 text-center">
                <NuxtLink
                    to="/"
                    class="inline-flex items-center text-sm text-gray-500 hover:text-gray-700 dark:text-gray-400 dark:hover:text-gray-200 transition-colors"
                >
                    <Icon
                        name="i-heroicons-arrow-left"
                        class="w-4 h-4 mr-2"
                    />
                    {{ $t('auth.backToHome',) }}
                </NuxtLink>
            </div>
        </div>
    </div>
</template>

<script setup lang="ts">
    import { ref, reactive, onMounted, computed, } from 'vue';
    import { useI18n, } from 'vue-i18n';
    import authApi from '~/api/auth';
    import { useAuthStore, } from '~/stores/AuthStore';

    const { t, } = useI18n();
    const authStore = useAuthStore();

    interface ProfileForm {
        username: string
        email: string
        full_name: string
        role: string
        created_at: string
    }

    const formState = reactive<ProfileForm>({
        username: '',
        email: '',
        full_name: '',
        role: '',
        created_at: '',
    });

    const originalData = ref<ProfileForm | null>(null,);
    const loading = ref(true,);
    const saving = ref(false,);
    const error = ref('',);
    const success = ref('',);

    // Avatar state
    const avatarUrl = ref('',);
    const avatarUploading = ref(false,);
    const avatarDeleting = ref(false,);
    const avatarError = ref('',);
    const avatarInputRef = ref<HTMLInputElement | null>(null,);

    const config = useRuntimeConfig();
    const serverUrl = (config.public.serverUrl as string).replace(/\/+$/, '',);

    const formatDate = (dateStr: string,): string => {
        const date = new Date(dateStr,);
        return date.toLocaleDateString(undefined, {
            year: 'numeric',
            month: 'long',
            day: 'numeric',
        });
    };

    const formatRole = (role: string,): string => {
        if (role === 'admin') return t('admin_pages.role_admin',);
        return t('admin_pages.role_user',);
    };

    const loadProfile = async () => {
        loading.value = true;
        error.value = '';
        try {
            const data = await authApi.getProfile();
            formState.username = data.username || '';
            formState.email = data.email || '';
            formState.full_name = data.name || '';
            formState.role = formatRole(data.role,);
            formState.created_at = formatDate(data.created_at,);
            avatarUrl.value = data.avatar_url ? `${serverUrl}${data.avatar_url}` : '';
            // Store original for reset
            originalData.value = { ...formState, };
        } catch (err: any) {
            error.value = err.error || t('profile.loadError',);
        } finally {
            loading.value = false;
        }
    };

    const handleSave = async () => {
        saving.value = true;
        error.value = '';
        success.value = '';

        try {
            const payload: { email?: string, full_name?: string } = {};
            if (formState.email !== originalData.value?.email) {
                payload.email = formState.email;
            }
            if (formState.full_name !== originalData.value?.full_name) {
                payload.full_name = formState.full_name;
            }

            if (Object.keys(payload,).length === 0) {
                success.value = t('profile.noChanges',);
                saving.value = false;
                return;
            }

            const data = await authApi.updateProfile(payload,);
            formState.email = data.email || '';
            formState.full_name = data.name || '';
            avatarUrl.value = data.avatar_url ? `${serverUrl}${data.avatar_url}?t=${Date.now()}` : avatarUrl.value;
            originalData.value = { ...formState, };
            success.value = t('profile.updateSuccess',);
        } catch (err: any) {
            error.value = err.error || t('profile.updateError',);
        } finally {
            saving.value = false;
        }
    };

    const resetForm = () => {
        if (originalData.value) {
            formState.username = originalData.value.username;
            formState.email = originalData.value.email;
            formState.full_name = originalData.value.full_name;
            formState.role = originalData.value.role;
            formState.created_at = originalData.value.created_at;
        }
        error.value = '';
        success.value = '';
    };

    // Avatar handlers
    const handleAvatarSelect = async (event: Event,) => {
        const input = event.target as HTMLInputElement;
        const files = input.files;
        if (!files || files.length === 0) return;

        const file = files[0] as File;

        // Validate file size (5MB)
        if (file.size > 5 * 1024 * 1024) {
            avatarError.value = t('profile.avatarTooLarge',);
            input.value = '';
            return;
        }

        // Validate file type
        const allowedTypes = ['image/jpeg', 'image/png', 'image/gif', 'image/webp',];
        if (!allowedTypes.includes(file.type,)) {
            avatarError.value = t('profile.avatarInvalidType',);
            input.value = '';
            return;
        }

        // Show preview
        const reader = new FileReader();
        reader.onload = (e,) => {
            if (e.target?.result) {
                avatarUrl.value = e.target.result as string;
            }
        };
        reader.readAsDataURL(file,);

        // Upload immediately
        avatarUploading.value = true;
        avatarError.value = '';
        success.value = '';

        try {
            const formData = new FormData();
            formData.append('avatar', file, file.name,);
            const data = await authApi.uploadAvatar(formData,);
            avatarUrl.value = `${serverUrl}${data.avatar_url}?t=${Date.now()}`;
            success.value = t('profile.avatarUploadSuccess',);
        } catch (err: any) {
            avatarError.value = err.error || t('profile.avatarUploadError',);
            loadProfile();
        } finally {
            avatarUploading.value = false;
            input.value = '';
        }
    };

    const handleAvatarClick = () => {
        avatarInputRef.value?.click();
    };

    const handleAvatarDelete = async () => {
        avatarDeleting.value = true;
        avatarError.value = '';
        success.value = '';

        try {
            await authApi.deleteAvatar();
            avatarUrl.value = '';
            success.value = t('profile.avatarDeleteSuccess',);
        } catch (err: any) {
            avatarError.value = err.error || t('profile.avatarDeleteError',);
        } finally {
            avatarDeleting.value = false;
        }
    };

    onMounted(() => {
        loadProfile();
    });
</script>

<style scoped>
/* Custom scrollbar for dark mode */
.dark ::-webkit-scrollbar {
    width: 8px;
}

.dark ::-webkit-scrollbar-track {
    background: #1f2937;
}

.dark ::-webkit-scrollbar-thumb {
    background: #4b5563;
    border-radius: 4px;
}

.dark ::-webkit-scrollbar-thumb:hover {
    background: #6b7280;
}

/* Smooth transitions */
* {
    transition: background-color 0.2s ease, border-color 0.2s ease;
}

/* Focus styles */
:focus-visible {
    outline: 2px solid var(--color-primary, #4B9E90);
    outline-offset: 2px;
}
</style>
