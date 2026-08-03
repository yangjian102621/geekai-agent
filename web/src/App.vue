<template>
  <router-view :key="routerKey" />
</template>

<script setup>
  import { useSharedStore } from '@/js/cache/sharedata'

  // 添加响应式 key 用于控制 router-view 刷新
  const routerKey = ref(0)
  const store = useSharedStore()

  watch(
    () => store.isLogin,
    () => {
      routerKey.value++
    }
  )

  console.log(
    '%c GeekAI-Agent 智能体平台',
    "color: green;font-size: 36px;font-family: '微软雅黑';"
  )

  console.log(
    '%c从现在开始，使用AI智能体开启工作效率倍增时代！',
    "color: #7c39ed;font-size: 18px;font-family: '微软雅黑';"
  )

  console.log(
    '%c 官方文档：https://docs.geekai.me',
    "color: green;font-size: 18px;font-family: '微软雅黑';"
  )
</script>

<style lang="scss">
  @use '@/assets/css/root.scss';
</style>
