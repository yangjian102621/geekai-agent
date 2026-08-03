<template>
  <!-- 侧边栏 -->
  <div :class="store.collapsed ? 'sidebar show' : 'sidebar'">
    <div class="p-3">
      <h4 class="mb-4">{{ title }}</h4>
      <el-scrollbar style="height: calc(100vh - 100px)">
        <el-menu
          class="sidebar-el-menu"
          :default-active="onRoutes"
          unique-opened
          router
        >
          <template v-for="item in items">
            <template v-if="item.subs">
              <el-sub-menu :index="item.index" :key="item.index">
                <template #title>
                  <i
                    :class="'iconfont icon-' + item.icon"
                    :style="item.style"
                  ></i>
                  <span>{{ item.title }}</span>
                </template>
                <template v-for="subItem in item.subs">
                  <el-sub-menu
                    v-if="subItem.subs"
                    :index="subItem.index"
                    :key="subItem.index"
                  >
                    <template #title>{{ subItem.title }}</template>
                    <el-menu-item
                      v-for="(threeItem, i) in subItem.subs"
                      :key="i"
                      :index="threeItem.index"
                    >
                      {{ threeItem.title }}
                    </el-menu-item>
                  </el-sub-menu>
                  <el-menu-item
                    v-else
                    :index="subItem.index"
                    :key="subItem.url"
                  >
                    <i
                      v-if="subItem.icon"
                      :class="'iconfont icon-' + subItem.icon"
                      :style="subItem.style"
                    ></i>
                    {{ subItem.title }}
                  </el-menu-item>
                </template>
              </el-sub-menu>
            </template>
            <template v-else>
              <el-menu-item :index="item.index" :key="item.index">
                <i
                  :class="'iconfont icon-' + item.icon"
                  :style="item.style"
                ></i>
                <template #title>{{ item.title }}</template>
              </el-menu-item>
            </template>
          </template>
        </el-menu>
      </el-scrollbar>
    </div>
  </div>
</template>

<script setup>
  import { getSystemInfoAdmin } from '@/js/cache/session.js'
  import { useSharedStore } from '@/js/cache/sharedata.js'
  import { computed, watch } from 'vue'
  import { useRoute } from 'vue-router'

  const items = [
    {
      icon: 'home',
      index: '/admin/dashboard',
      title: '仪表盘',
    },

    {
      icon: 'user-fill',
      index: '/admin/users',
      title: '用户管理',
    },
    {
      icon: 'ai-agent',
      index: 'admin-apps',
      title: '智能体管理',
      subs: [
        {
          index: '/admin/apps/list',
          title: '智能体列表',
          icon: 'list',
        },
        {
          index: '/admin/apps/category',
          title: '智能体分类',
          icon: 'categroy',
        },
      ],
    },
    {
      icon: 'workflow',
      index: '/admin/workflows',
      title: '工作流管理',
    },
    {
      icon: 'chuangzuo',
      index: 'admin-creators',
      title: '创作者管理',
      style: '--icon-font-size:16px',
      subs: [
        {
          index: '/admin/creators',
          title: '创作者管理',
          icon: 'role',
        },
        {
          index: '/admin/creator/apps',
          title: '应用管理',
          icon: 'app',
        },
        {
          index: '/admin/creator/withdraws',
          title: '提现管理',
          icon: 'withdraw',
        },
      ],
    },
    {
      icon: 'redeem',
      index: '/admin/redeem',
      title: '兑换码',
      style: '--icon-font-size:14px',
    },
    {
      icon: 'recharge',
      index: '/admin/products',
      title: '充值产品',
      style: '--icon-font-size:20px',
    },
    {
      icon: 'order',
      index: '/admin/orders',
      title: '订单管理',
      style: '--icon-font-size:18px',
    },
    {
      icon: 'log',
      index: '/admin/score/log',
      title: '积分日志',
      style: '--icon-font-size:20px',
    },
    {
      icon: 'log',
      index: '/admin/users/loginLog',
      title: '登录日志',
      style: '--icon-font-size:20px',
    },
    {
      icon: 'role',
      index: '/admin/manager',
      title: '管理员',
    },
    {
      icon: 'control',
      index: 'control',
      title: '系统设置',
      subs: [
        {
          index: '/admin/settings/basic',
          title: '基础设置',
          icon: 'config',
        },
        {
          index: '/admin/settings/notice',
          title: '公告设置',
          icon: 'bell',
        },
        {
          index: '/admin/settings/geek',
          title: '增值服务配置',
          icon: 'model',
        },
        {
          index: '/admin/settings/coze',
          title: 'Coze 设置',
          icon: 'coze',
          style: '--icon-font-size:22px',
        },
        {
          index: '/admin/settings/payment',
          title: '支付设置',
          icon: 'alipay',
        },
        {
          index: '/admin/settings/sms',
          title: '短信设置',
          icon: 'duanxin',
        },
        {
          index: '/admin/settings/oss',
          title: '存储设置',
          icon: 'storage',
        },
        {
          index: '/admin/settings/smtp',
          title: '邮件设置',
          icon: 'email',
        },
      ],
    },
  ]

  const route = useRoute()
  const title = ref('')
  const onRoutes = computed(() => {
    return route.path
  })

  // 监听路由变化，移动端点击菜单跳转时候自动收起侧边栏
  const store = useSharedStore()
  watch(
    () => route.path,
    () => {
      store.collapsed = false
    }
  )
  onMounted(async () => {
    const systemConfig = await getSystemInfoAdmin()
    title.value = systemConfig.admin_title
  })
</script>

<style scoped lang="scss">
  .sidebar {
    --icon-font-size: 18px;
    background: #fff;
    box-shadow: 0 0 10px rgba(0, 0, 0, 0.1);
    height: 100vh;
    position: fixed;
    width: 250px;

    .sidebar-el-menu {
      --el-menu-item-color: #333;
      --el-menu-bg-color: #fff;
      --el-menu-hover-bg-color: #f8f2ff;
      --el-menu-hover-text-color: #7c39ed;
      --el-menu-border-color: 0;
    }

    .el-menu-item,
    .el-sub-menu {
      color: #333;
      margin: 0.2rem 0;

      i {
        font-size: var(--icon-font-size);
        margin-right: 10px;
      }
    }

    .el-menu-item.is-active {
      background: #f8f2ff;
      color: var(--el-menu-hover-text-color);
    }
  }

  @media (max-width: 768px) {
    .sidebar {
      transform: translateX(-100%);
      transition: transform 0.3s ease;
      z-index: 1000;
    }

    .sidebar.show {
      transform: translateX(0);
    }

    .main-content {
      margin-left: 0;
    }
  }
</style>
