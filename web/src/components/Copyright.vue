<template>
  <span>
    <span>{{ copyright }}</span>
    © {{ year }} All rights reserved
  </span>
</template>

<script setup>
  import { getLicense, getSystemInfo } from '@/js/cache/session.js'

  const systemConfig = ref({})
  const copyright = ref('')
  const license = ref({})
  const year = ref(new Date().getFullYear())

  onMounted(async () => {
    // 获取 license 信息
    license.value = await getLicense()
    // 获取 system 信息
    systemConfig.value = await getSystemInfo()
    if (license.value.is_active) {
      copyright.value = systemConfig.value.copyright
    } else {
      copyright.value = '极客学长'
    }
  })
</script>
