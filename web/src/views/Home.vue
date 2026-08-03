<template>
  <div class="app-home-layout">
    <!-- 最左侧一级导航栏 -->
    <aside class="primary-sidebar">
      <div class="logo-container" @click="router.push('/')">
        <el-avatar :size="40" :src="logo" shape="square" class="logo-avatar" />
        <span class="logo-text" v-if="!collapsed">GeekAI</span>
      </div>

      <nav class="nav-menu">
        <div
          v-for="item in menuItems"
          :key="item.path"
          class="nav-item"
          :class="{ active: isRouteActive(item.path) }"
          @click="navigateTo(item.path)"
        >
          <div class="nav-icon">
            <i class="iconfont" :class="item.icon"></i>
          </div>
          <span class="nav-label">{{ item.label }}</span>
        </div>
      </nav>

      <div class="user-profile">
        <el-dropdown
          trigger="click"
          placement="top-start"
          v-if="userStore.isLogin"
        >
          <div class="avatar-wrapper">
            <el-avatar
              :size="36"
              :src="userStore.userInfo.avatar || '/images/avatar/user.png'"
            />
          </div>
          <template #dropdown>
            <el-dropdown-menu>
              <el-dropdown-item @click="showSetting = true">
                <i class="iconfont icon-config text-primary me-2"></i> 系统设置
              </el-dropdown-item>
              <el-dropdown-item @click="showScoreLog = true">
                <i class="iconfont icon-log text-info me-2"></i> 消费日志
              </el-dropdown-item>
              <el-dropdown-item @click="showCharge = true">
                <i class="iconfont icon-recharge text-success me-2"></i>
                积分充值
              </el-dropdown-item>
              <el-dropdown-item v-if="isCreator">
                <a
                  href="/creator/console"
                  target="_blank"
                  class="text-decoration-none text-dark"
                >
                  <i class="iconfont icon-role text-success me-2"></i>
                  创作中心
                </a>
              </el-dropdown-item>
              <el-dropdown-item v-else @click="showCreatorApply = true">
                <span class="cursor-pointer">
                  <i class="iconfont icon-role text-success me-2"></i>
                  申请创作者
                </span>
              </el-dropdown-item>
              <el-dropdown-item divided @click="handleLogout">
                <i class="iconfont icon-logout text-secondary me-2"></i>
                退出登录
              </el-dropdown-item>
            </el-dropdown-menu>
          </template>
        </el-dropdown>

        <div
          v-else
          class="flex flex-column align-items-center justify-content-center"
          @click="sharedStore.setShowLoginDialog(true)"
        >
          <button
            data-v-4002dbf7=""
            class="btn btn-primary !px-2 !py-0.5 !text-base"
          >
            登录
          </button>
        </div>
      </div>
    </aside>

    <!-- 二级路由内容区 -->
    <main class="main-content-area">
      <router-view v-slot="{ Component }">
        <keep-alive include="ChatPage">
          <component :is="Component" />
        </keep-alive>
      </router-view>
    </main>

    <el-dialog
      v-model="sharedStore.showLoginDialog"
      width="500px"
      @close="sharedStore.setShowLoginDialog(false)"
    >
      <template #header>
        <div
          class="text-center text-xl"
          style="color: var(--theme-text-color-primary)"
        >
          登录后解锁更多功能
        </div>
      </template>
      <div class="p-4 pt-2 pb-2">
        <LoginDialog />
      </div>
    </el-dialog>

    <!-- 弹窗组件 -->
    <model-dialog
      v-model="showSetting"
      title="系统设置"
      :hide-footer="true"
      @cancel="showSetting = false"
      padding="0 1rem 1rem 1rem"
    >
      <settings @update-user="updateUser" :active-name="activeName" />
    </model-dialog>

    <model-dialog
      v-model="showCharge"
      title="积分充值"
      :hide-footer="true"
      @cancel="showCharge = false"
      padding="0 1rem 1rem 1rem"
    >
      <charge />
    </model-dialog>

    <model-dialog
      v-model="showScoreLog"
      title="积分日志"
      :hide-footer="true"
      @cancel="showScoreLog = false"
      padding="0 1rem 1rem 1rem"
      :width="1500"
    >
      <score-log />
    </model-dialog>

    <model-dialog
      v-model="showCreatorApply"
      title="申请成为创作者"
      :hide-footer="true"
      @cancel="showCreatorApply = false"
      padding="0 1rem 1rem 1rem"
    >
      <creator-apply @success="showCreatorApply = false" />
    </model-dialog>
  </div>
</template>

<script setup>
  import { ref, computed, onMounted } from 'vue'
  import { useRouter, useRoute } from 'vue-router'
  import { useUserStore } from '@/js/store/user'
  import { useSharedStore } from '@/js/cache/sharedata'
  import { removeUserToken } from '@/js/cache/session'
  import { httpGet } from '@/js/utils/http'

  import LoginDialog from '@/components/LoginDialog.vue'
  import ModelDialog from '@/components/ModelDialog.vue'
  import Settings from '@/components/Settings.vue'
  import Charge from '@/components/Charge.vue'
  import ScoreLog from '@/components/ScoreLog.vue'
  import CreatorApply from '@/components/creator/CreatorApply.vue'

  const router = useRouter()
  const route = useRoute()
  const userStore = useUserStore()
  const sharedStore = useSharedStore()
  const collapsed = ref(false)
  const logo = ref('/images/logo.png')

  // 状态变量
  const showSetting = ref(false)
  const showCharge = ref(false)
  const showScoreLog = ref(false)
  const showCreatorApply = ref(false)
  const isCreator = ref(false)
  const activeName = ref('profile')

  const menuItems = [
    { label: '智能体', path: '/chat', icon: 'icon-ai-agent' },
    { label: '工作流', path: '/workflow', icon: 'icon-workflow' },
  ]

  const isRouteActive = (path) => {
    if (path === '/' && route.path === '/') return true
    if (path !== '/' && route.path.startsWith(path)) return true
    return false
  }

  const navigateTo = (path) => {
    router.push(path)
  }

  const handleLogout = () => {
    httpGet('/api/user/logout').then(() => {
      removeUserToken()
      location.reload()
    })
  }

  const updateUser = (user) => {
    userStore.userInfo = { ...userStore.userInfo, ...user }
    showSetting.value = false
  }

  onMounted(async () => {
    await userStore.fetchUserInfo()
    // 检查创作者状态
    if (userStore.isLogin) {
      try {
        const res = await httpGet('/api/creator/status')
        isCreator.value = res.data.status === 'approved'
      } catch (e) {
        isCreator.value = false
      }
    }
  })
</script>

<style scoped lang="scss">
  .app-home-layout {
    display: flex;
    width: 100vw;
    height: 100vh;
    background-color: #f7f9fb;

    .primary-sidebar {
      width: 72px;
      background-color: #fff;
      border-right: 1px solid #e5e7eb;
      display: flex;
      flex-direction: column;
      align-items: center;
      padding: 16px 0;
      flex-shrink: 0;
      z-index: 100;

      .logo-container {
        margin-bottom: 24px;
        cursor: pointer;
        display: flex;
        flex-direction: column;
        align-items: center;

        .logo-avatar {
          background: transparent;
        }
        .logo-text {
          display: none;
        }
      }

      .nav-menu {
        flex: 1;
        display: flex;
        flex-direction: column;
        gap: 16px;
        width: 100%;
        align-items: center;

        .nav-item {
          display: flex;
          flex-direction: column;
          align-items: center;
          justify-content: center;
          width: 56px;
          height: 56px;
          border-radius: 12px;
          cursor: pointer;
          color: #64748b;
          transition: all 0.2s;

          .nav-icon {
            font-size: 24px;
            margin-bottom: 2px;
            i {
              font-size: 22px;
            }
          }

          .nav-label {
            font-size: 10px;
            line-height: 1.2;
          }

          &:hover {
            background-color: #f1f5f9;
            color: var(--el-color-primary);
          }

          &.active {
            background-color: #eff6ff;
            color: var(--el-color-primary);
          }
        }
      }

      .user-profile {
        margin-top: auto;

        .avatar-wrapper {
          cursor: pointer;
          border: 2px solid transparent;
          border-radius: 50%;
          transition: border-color 0.2s;

          &:hover {
            border-color: var(--el-color-primary);
          }
        }
      }
    }

    .main-content-area {
      flex: 1;
      height: 100%;
      overflow: hidden;
      background-color: #fff;
      position: relative;
    }
  }

  /* 移动端适配 */
  @media (max-width: 768px) {
    .app-home-layout {
      flex-direction: column-reverse;

      .primary-sidebar {
        width: 100%;
        height: 60px;
        flex-direction: row;
        border-right: none;
        border-top: 1px solid #e5e7eb;
        padding: 0 16px;
        justify-content: space-around;

        .logo-container,
        .user-profile {
          display: none;
        }

        .nav-menu {
          flex-direction: row;
          justify-content: space-around;
          gap: 0;

          .nav-item {
            width: auto;
            height: 100%;
            padding: 0 12px;
            border-radius: 0;

            &.active {
              background-color: transparent;
              color: var(--el-color-primary);
              position: relative;

              &::after {
                content: '';
                position: absolute;
                top: 0;
                left: 50%;
                transform: translateX(-50%);
                width: 20px;
                height: 3px;
                background-color: var(--el-color-primary);
                border-radius: 0 0 3px 3px;
              }
            }
          }
        }
      }

      .main-content-area {
        height: calc(100vh - 60px);
      }
    }
  }
</style>
