<script lang="ts">
import ModalDialog from './ModalDialog.vue';
import { useThemeStore } from '../../stores/ThemeStore';
import type { PropType } from 'vue';
import PictureCardSkeleton from '../app-skeletons/PictureCardSkeleton.vue';
import PicCarousel from './PicCarousel.vue';
import { useAuthStore } from '../../stores/AuthStore';
import { storeToRefs } from 'pinia';
import { CommonGetWorkDto, GetSaleDto } from '@/types';
import Alert from './Alert.vue';

export default {
    setup() {
        const themeStore = useThemeStore();
        const authStore = useAuthStore();
        const { isAuthenticated } = storeToRefs(authStore);

        return {
            themeStore,
            isAuthenticated,
        };
    },
    components: {
        Modal: ModalDialog,
        CardSkeleton: PictureCardSkeleton,
        Carousel: PicCarousel,
        Alert,
    },
    props: {
        imageObject: {
            type: Object as PropType<CommonGetWorkDto>,
            default: {} as CommonGetWorkDto,
        },
    },
    emits: ['work-deleted'],
    data() {
        return {
            isLoaded: false,
            imagebasedir: import.meta.env.VITE_IMAGE_DIR,
            showAlert: false,
            alertType: 'success',
            alertTitle: '',
            alertMessage: '',
        };
    },
    methods: {
        onImgLoad() {
            setTimeout(() => {
                this.isLoaded = true;
            }, 1000);
        },
        edit() {
            // Check if the item is a shop item (has price) or gallery item (has type)
            if (this.imageObject.__type === 'GetSaleDto') {
                this.$router.push('/sales/edit/' + this.imageObject.id);
            } else if (this.imageObject.__type === 'GetWorkDto') {
                this.$router.push('/works/edit/' + this.imageObject.id);
            } else {
                console.error('Unknown item type', this.imageObject.__type);
            }
        },
        async onImageDelete(id: string) {
            try {
                const typeOfWork =
                    this.imageObject.__type === 'GetWorkDto' ? 'works/' : 'sales/';
                const token = localStorage.getItem('token');
                const response = await fetch(
                    import.meta.env.VITE_SERVER_URL + typeOfWork + id,
                    {
                        method: 'DELETE',
                        headers: {
                            Authorization: `Bearer ${token}`,
                        },
                    }
                );

                if (response.ok) {
                    // Emit event to parent component to update the list
                    this.$emit('work-deleted', this.imageObject.id);
                    // Show success alert
                    this.showAlertMessage('success', 'Успешно', 'Работа успешно удалена');
                } else {
                    console.error('Failed to delete work');
                    // Show error alert
                    this.showAlertMessage(
                        'danger',
                        'Ошибка',
                        'Не удалось удалить работу'
                    );
                }
            } catch (error) {
                console.error('Error deleting work:', error);
                // Show error alert
                this.showAlertMessage(
                    'danger',
                    'Ошибка',
                    'Произошла ошибка при удалении'
                );
            } finally {
                // Close the modal in all cases
                this.closeDeleteModal();
            }
        },
        showAlertMessage(type: 'success' | 'danger', title: string, message: string) {
            this.alertType = type;
            this.alertTitle = title;
            this.alertMessage = message;
            this.showAlert = true;

            // Hide alert after 5 seconds
            setTimeout(() => {
                this.showAlert = false;
            }, 5000);
        },
        closeDeleteModal() {
            const modal = document.getElementById(this.mapImgIdToDeleteId);
            // Try different ways to get modal instance
            let bsModal = null;

            if ((window as any).bootstrap?.Modal) {
                bsModal = (window as any).bootstrap.Modal.getInstance(modal);
            } else if (
                (window as any).jQuery &&
                (window as any).jQuery(modal).data('bs.modal')
            ) {
                // jQuery fallback for Bootstrap 4
                bsModal = (window as any).jQuery(modal).data('bs.modal');
            }

            console.log('bsModal', bsModal);
            if (bsModal) {
                bsModal.hide();
            } else {
                // Manual hide as fallback
                modal.classList.remove('show');
                modal.style.display = 'none';
                document.body.classList.remove('modal-open');
                const backdrop = document.querySelector('.modal-backdrop');
                if (backdrop) backdrop.remove();

                const restModal = document.getElementById('modal-backdrop');
                if (restModal) restModal.remove();
            }
        },
    },
    computed: {
        imgIdToModalIdSelector() {
            return '#' + this.imageObject.str_id + 'Modal';
        },
        imgIdToDeleteIdSelector() {
            return '#' + this.imageObject.str_id + 'DeleteModal';
        },
        imgIdToModalId() {
            return this.imageObject.str_id + 'Modal';
        },
        mapImgIdToDeleteId() {
            return this.imageObject.str_id + 'DeleteModal';
        },
        mainCardImage() {
            return this.imagebasedir + this.imageObject.dir + this.imageObject.images[0];
        },
        showCardShadow() {
            if (this.themeStore.theme === 'light') {
                return 'shadow-lg';
            } else return '';
        },
    },
};
</script>

<template>
    <div class="col-12 col-sm-6 col-md-4 d-flex justify-content-center">
        <div class="card p-3 mb-5 rounded showCardShadow" v-show="isLoaded">
            <img
                :src="mainCardImage"
                @load="onImgLoad"
                class="card-img-top"
                :alt="$i18n.locale === 'RUS' ? imageObject.name_ru : imageObject.name_en"
                data-bs-toggle="modal"
                :data-bs-target="imgIdToModalIdSelector"
            />

            <div class="card-body">
                <div class="desc">
                    <h5 class="card-title">
                        {{
                            $i18n.locale === 'RUS'
                                ? imageObject.name_ru
                                : imageObject.name_en
                        }}
                    </h5>

                    <div class="card-text">
                        <span v-if="imageObject.base_ru">
                            {{
                                $i18n.locale === 'RUS'
                                    ? imageObject.base_ru
                                    : imageObject.base_en
                            }}
                        </span>
                        <span v-if="imageObject.material_ru">
                            {{
                                $i18n.locale === 'RUS'
                                    ? `, ${imageObject.material_ru}`
                                    : `, ${imageObject.material_en}`
                            }}
                        </span>
                        <span v-if="imageObject.width && imageObject.width != '0'">
                            {{ `, ${imageObject.width + 'x' + imageObject.height}` }}
                        </span>
                        <span v-if="imageObject.year">
                            {{ `, ${imageObject.year}` }}
                        </span>
                    </div>

                    <p v-if="imageObject.__type === 'GetSaleDto' && imageObject.price">
                        {{ $t('card.price') + ` ${imageObject.price} ` + $t('card.rub') }}
                    </p>

                    <div class="image-control" v-if="isAuthenticated">
                        <UButton
                            class="btn btn-secondary w-100"
                            type="button"
                            v-on:click="edit"
                        >
                            Редактировать
                        </UButton>
                        <UButton
                            class="btn btn-secondary w-100"
                            type="button"
                            data-bs-toggle="modal"
                            :data-bs-target="imgIdToDeleteIdSelector"
                        >
                            Удалить
                        </UButton>
                    </div>
                </div>
            </div>

            <UModal :modalId="imgIdToModalId">
                <Carousel :imageObject="imageObject" :imageId="imageObject.str_id" />
            </Modal>
        </div>
        <CardSkeleton v-if="!isLoaded" />

        <div
            class="modal fade"
            :id="mapImgIdToDeleteId"
            tabindex="-1"
            aria-labelledby="imageDeleteModalLabel"
            aria-hidden="true"
        >
            <div class="modal-dialog modal-dialog-centered">
                <div class="modal-content">
                    <div class="modal-header">
                        <h5 class="modal-title" id="imageDeleteModalLabel">
                            Подтверждение удаления
                        </h5>
                        <UButton
                            type="button"
                            class="btn-close"
                            data-bs-dismiss="modal"
                            aria-label="Close"
                        ></UButton>
                    </div>
                    <div class="modal-body">
                        Вы действительно хотите удалить работу "{{
                            imageObject.name_ru
                        }}"?
                    </div>
                    <div class="modal-footer">
                        <UButton
                            type="button"
                            class="btn btn-secondary"
                            data-bs-dismiss="modal"
                        >
                            Отменить
                        </UButton>
                        <UButton
                            type="button"
                            class="btn btn-danger"
                            @click="onImageDelete(imageObject.id)"
                        >
                            Удалить
                        </UButton>
                    </div>
                </div>
            </div>
        </div>

        <UAlert
            v-model="showAlert"
            :type="alertType"
            :title="alertTitle"
            :message="alertMessage"
            close-button-text="Закрыть"
        />
    </div>
</template>

<style scoped>
.card-body {
    display: flex;
    align-items: center;
    padding: 0;
}

.card {
    --bs-card-bg: var(--color-surface);
    opacity: 1;
    width: 100%;
}

.desc {
    color: var(--color-on-surface);
    align-self: flex-end;
    width: 100%;
}

img {
    display: block;
    width: 100%;
    height: inherit;
    object-fit: cover;
}

.image-control {
    display: flex;
    flex-direction: column;
    justify-content: space-between;
    gap: 0.5rem;
    margin-top: 1rem;
}

.card > img:hover {
    cursor: pointer;
}
</style>
