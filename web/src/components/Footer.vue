<template>
  <div class="footer text-center text-sm">
    <p class="mb-0">
      <Copyright />
    </p>
    <div class="flex justify-center items-center">
      <span v-if="license.is_active">
        {{ systemConfig.title }} - {{ version }}
      </span>
      <a
        href="https://docs.geekai.me"
        target="_blank"
        class="hover:underline"
        v-else
      >
        GeekAI-Agent - {{ version }}
      </a>
    </div>
  </div>
</template>

<script setup>
  import Copyright from '@/components/Copyright.vue'
  import { getLicense, getSystemInfo } from '@/js/cache/session.js'

  const systemConfig = ref({})
  const version = import.meta.env.VITE_APP_VERSION
  const license = ref({})

  onMounted(async () => {
    systemConfig.value = await getSystemInfo()
    license.value = await getLicense()
  })
</script>

<style lang="scss" scoped></style>
