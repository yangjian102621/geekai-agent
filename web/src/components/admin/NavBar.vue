<template>
  <!-- 顶部导航栏 -->
  <nav class="navbar">
    <div
      :class="
        store.collapsed
          ? 'container-fluid justify-content-end'
          : 'container-fluid'
      "
    >
      <button class="btn btn-toggle" type="button" @click="toggleSidebar()">
        <i class="iconfont icon-sub-menu"></i>
      </button>

      <div class="d-flex align-items-start ms-auto" v-if="!store.collapsed">
        <el-dropdown class="user-name">
          <a class="dropdown-toggle d-flex align-items-center" role="button">
            <img
              :src="avatar"
              class="rounded-circle"
              alt="用户头像"
              width="30"
              height="30"
            />
            <span class="nav-profile-name">Geek Master</span>
          </a>

          <template #dropdown>
            <el-dropdown-menu>
              <el-dropdown-item>
                <i class="iconfont icon-version"></i> 当前版本-{{ version }}
              </el-dropdown-item>
              <el-dropdown-item divided @click="logout">
                <i class="iconfont icon-logout"></i>
                <span>退出登录</span>
              </el-dropdown-item>
            </el-dropdown-menu>
          </template>
        </el-dropdown>
      </div>
    </div>
  </nav>
</template>

<script setup>
  import { removeAdminToken } from '@/js/cache/session.js'
  import { useSharedStore } from '@/js/cache/sharedata.js'
  import { showMessageError } from '@/js/utils/dialog.js'
  import { httpGet } from '@/js/utils/http.js'

  const avatar = ref('/images/admin/avatar.jpg')
  const store = useSharedStore()

  const toggleSidebar = () => {
    store.collapsed = !store.collapsed
  }

  const router = useRouter()
  const version = import.meta.env.VITE_APP_VERSION
  const logout = () => {
    httpGet('/api/admin/logout')
      .then(() => {
        removeAdminToken()
        router.push('/admin/login')
      })
      .catch((e) => {
        showMessageError('注销失败：' + e.message)
      })
  }
</script>

<style scoped lang="scss">
  .navbar {
    background: #fff;
    box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);

    .user-name {
      .nav-profile-name {
        margin-left: 0.5rem;
        margin-right: 0.5rem;
        color: #4a4a4a;
        font-weight: 500;
      }
    }

    .btn-toggle {
      border-color: #5a23c8;
      display: none;
    }
  }

  @media (max-width: 768px) {
    .navbar {
      .btn-toggle {
        display: block;
      }
    }
  }
</style>
<style lang="scss">
  // dropdown styles
  .el-dropdown__popper {
    --el-dropdown-menuItem-hover-color: #7c39ed;
    --el-dropdown-menuItem-hover-fill: #f8f2ff;
  }
</style>
