<template>
  <transition name="fade">
    <div
      class="modal"
      v-if="showDialog"
      tabindex="-1"
      style="display: block; background-color: rgba(0, 0, 0, 0.5)"
    >
      <div
        :class="'modal-dialog modal-dialog-scrollable modal-dialog-centered modal-lg'"
        :style="{ maxWidth: maxWidth }"
      >
        <transition name="zoom">
          <div class="modal-content">
            <div class="modal-header" v-if="title">
              <h1 class="modal-title fs-5">
                {{ title }}
              </h1>
              <button
                type="button"
                class="btn-close"
                data-bs-dismiss="modal"
                @click="close"
                aria-label="Close"
              ></button>
            </div>
            <div class="modal-body" :style="{ padding: padding }">
              <slot></slot>
            </div>
            <div class="modal-footer" v-if="!hideFooter">
              <button
                type="button"
                class="btn btn-secondary"
                data-bs-dismiss="modal"
                @click="close"
                v-if="cancelText"
              >
                {{ cancelText }}
              </button>
              <button
                type="button"
                class="btn btn-primary"
                v-if="!hideConfirm && showConfirm && confirmText"
                @click="confirm()"
              >
                {{ confirmText }}
              </button>
            </div>
          </div>
        </transition>
      </div>
    </div>
  </transition>
</template>

<script setup>
  import { ref, watch } from 'vue'

  const props = defineProps({
    modelValue: Boolean,
    title: {
      type: String,
      default: '',
    },
    loading: {
      type: Boolean,
      default: false,
    },
    padding: {
      type: String,
      default: '1rem 1rem 1rem 1rem',
    },
    hideFooter: {
      type: Boolean,
      default: false,
    },
    hideConfirm: {
      type: Boolean,
      default: false,
    },
    confirmText: {
      type: String,
      default: '确定',
    },
    cancelText: {
      type: String,
      default: '取消',
    },
    width: {
      type: Number,
      default: 500,
    },
    showConfirm: {
      type: Boolean,
      default: true,
    },
  })
  const emits = defineEmits(['confirm', 'cancel'])
  const showDialog = ref(props.modelValue)
  const maxWidth = computed(() => {
    return props.width < 100 ? props.width + '%' : props.width + 'px'
  })

  watch(
    () => props.modelValue,
    (newValue) => {
      showDialog.value = newValue
    }
  )

  const close = () => {
    showDialog.value = false
    emits('cancel')
  }
  const confirm = () => {
    emits('confirm')
  }
</script>

<style scoped lang="scss">
  /* 淡入淡出动画 */
  .fade-enter-active,
  .fade-leave-active {
    transition: opacity 0.3s ease;
  }

  .fade-enter-from,
  .fade-leave-to {
    opacity: 0;
  }

  /* 缩放动画 */
  .zoom-enter-active,
  .zoom-leave-active {
    transition: transform 3s cubic-bezier(0.4, 0, 0.2, 1),
      opacity 3s cubic-bezier(0.4, 0, 0.2, 1);
  }

  .zoom-enter-from,
  .zoom-leave-to {
    opacity: 0;
    transform: scale(0.95);
  }

  @media (min-width: 1200px) {
    .modal-lg {
      --bs-modal-width: 800px;
    }
    .modal-xl {
      --bs-modal-width: 1000px;
    }
  }
  @media (min-width: 1400px) {
    .modal-lg {
      --bs-modal-width: 1000px;
    }
    .modal-xl {
      --bs-modal-width: 1200px;
    }
  }
</style>
