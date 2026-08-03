<template>
  <div class="app-container">
    <SideBar />
    <!-- 主内容区 -->
    <div class="main-content">
      <NavBar />

      <!-- 表单区域 -->
      <div class="content" :style="{ height: winHeight + 'px' }">
        <router-view v-slot="{ Component }">
          <transition name="move" mode="out-in">
            <component :is="Component"></component>
          </transition>
        </router-view>
      </div>
    </div>
  </div>
</template>

<script setup>
  import NavBar from '@/components/admin/NavBar.vue'
  import SideBar from '@/components/admin/SideBar.vue'
  import { checkAdminSession } from '@/js/cache/session.js'

  const winHeight = ref(window.innerHeight - 50)

  const router = useRouter()
  checkAdminSession().catch(() => {
    router.push('/admin/login')
  })
</script>

<style scoped lang="scss">
  .app-container {
    .main-content {
      margin-left: 250px;
      padding: 0;
      background: #f8f9fa;
      min-height: 100vh;
      overflow: hidden;

      .content {
        overflow-y: auto;
        overflow-x: hidden;
        padding: 20px;
      }
    }
  }

  @media (max-width: 768px) {
    .app-container {
      .main-content {
        margin-left: 0;
      }
    }
  }
</style>
