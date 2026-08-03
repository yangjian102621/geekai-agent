<template>
  <div
    class="alert mb-4 d-flex align-items-center"
    :class="typeClass"
    role="alert"
  >
    <div>
      <slot name="icon">
        <i class="iconfont !text-xl me-2" :class="typeIcon"></i>
      </slot>
    </div>
    <div class="text-sm">
      <slot>{{ message }}</slot>
    </div>
  </div>
</template>

<script setup>
  import { computed } from 'vue'

  const props = defineProps({
    type: {
      type: String,
      default: 'info',
    },
  })

  const typeClass = computed(() => {
    return {
      info: 'alert-primary',
      success: 'alert-success',
      warning: 'alert-warning',
    }[props.type]
  })
  const typeIcon = computed(() => {
    return {
      info: 'icon-info',
      success: 'icon-success',
      warning: 'icon-warning',
    }[props.type]
  })
</script>

<style scoped lang="scss"></style>
