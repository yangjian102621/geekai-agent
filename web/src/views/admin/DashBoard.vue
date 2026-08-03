<template>
  <div class="dashboard">
    <div class="row g-4">
      <!-- 第一行统计卡片 -->
      <div class="col-12 col-sm-6 col-xl-3">
        <div class="card stat-card">
          <div class="card-body">
            <div class="d-flex justify-content-between align-items-center">
              <div>
                <h6 class="card-title text-muted mb-2">用户总数</h6>
                <h2 class="mb-0 text-primary">
                  {{ dashboardStore.formattedStats.user_total }}
                </h2>
                <small class="text-muted"
                  >今日新增:
                  {{ dashboardStore.formattedStats.user_today_new }}</small
                >
              </div>
              <div class="stat-icon">
                <i class="iconfont text-2xl icon-user-fill"></i>
              </div>
            </div>
          </div>
        </div>
      </div>

      <div class="col-12 col-sm-6 col-xl-3">
        <div class="card stat-card">
          <div class="card-body">
            <div class="d-flex justify-content-between align-items-center">
              <div>
                <h6 class="card-title text-muted mb-2">应用总数</h6>
                <h2 class="mb-0 text-success">
                  {{ dashboardStore.formattedStats.app_total }}
                </h2>
                <small class="text-muted"
                  >启用: {{ dashboardStore.formattedStats.app_enabled }}</small
                >
              </div>
              <div class="stat-icon">
                <i class="iconfont text-2xl icon-app"></i>
              </div>
            </div>
          </div>
        </div>
      </div>

      <div class="col-12 col-sm-6 col-xl-3">
        <div class="card stat-card">
          <div class="card-body">
            <div class="d-flex justify-content-between align-items-center">
              <div>
                <h6 class="card-title text-muted mb-2">对话总数</h6>
                <h2 class="mb-0 text-info">
                  {{ dashboardStore.formattedStats.chat_total }}
                </h2>
                <small class="text-muted"
                  >今日: {{ dashboardStore.formattedStats.chat_today }}</small
                >
              </div>
              <div class="stat-icon">
                <i class="iconfont text-2xl icon-chat"></i>
              </div>
            </div>
          </div>
        </div>
      </div>

      <div class="col-12 col-sm-6 col-xl-3">
        <div class="card stat-card">
          <div class="card-body">
            <div class="d-flex justify-content-between align-items-center">
              <div>
                <h6 class="card-title text-muted mb-2">总收入</h6>
                <h2 class="mb-0 text-warning">
                  ￥{{ dashboardStore.formattedStats.revenue_total }}
                </h2>
                <small class="text-muted"
                  >今日: ￥{{
                    dashboardStore.formattedStats.revenue_today
                  }}</small
                >
              </div>
              <div class="stat-icon">
                <i class="iconfont text-3xl icon-money"></i>
              </div>
            </div>
          </div>
        </div>
      </div>

      <!-- 第二行统计卡片 -->
      <div class="col-12 col-sm-6 col-xl-3">
        <div class="card stat-card">
          <div class="card-body">
            <div class="d-flex justify-content-between align-items-center">
              <div>
                <h6 class="card-title text-muted mb-2">订单总数</h6>
                <h2 class="mb-0 text-purple">
                  {{ dashboardStore.formattedStats.order_total }}
                </h2>
                <small class="text-muted"
                  >今日: {{ dashboardStore.formattedStats.order_today }}</small
                >
              </div>
              <div class="stat-icon">
                <i class="iconfont text-2xl icon-order"></i>
              </div>
            </div>
          </div>
        </div>
      </div>

      <div class="col-12 col-sm-6 col-xl-3">
        <div class="card stat-card">
          <div class="card-body">
            <div class="d-flex justify-content-between align-items-center">
              <div>
                <h6 class="card-title text-muted mb-2">活跃用户</h6>
                <h2 class="mb-0 text-danger">
                  {{ dashboardStore.formattedStats.user_active }}
                </h2>
                <small class="text-muted">近7天内</small>
              </div>
              <div class="stat-icon">
                <i class="iconfont text-2xl icon-user-circle"></i>
              </div>
            </div>
          </div>
        </div>
      </div>

      <div class="col-12 col-sm-6 col-xl-3">
        <div class="card stat-card">
          <div class="card-body">
            <div class="d-flex justify-content-between align-items-center">
              <div>
                <h6 class="card-title text-muted mb-2">积分消费</h6>
                <h2 class="mb-0 text-dark">
                  {{ dashboardStore.formattedStats.score_consumed }}
                </h2>
                <small class="text-muted"
                  >充值:
                  {{ dashboardStore.formattedStats.score_recharged }}</small
                >
              </div>
              <div class="stat-icon">
                <i class="iconfont text-3xl icon-doller"></i>
              </div>
            </div>
          </div>
        </div>
      </div>

      <div class="col-12 col-sm-6 col-xl-3">
        <div class="card stat-card">
          <div class="card-body">
            <div class="d-flex justify-content-between align-items-center">
              <div>
                <h6 class="card-title text-muted mb-2">创作者</h6>
                <h2 class="mb-0 text-indigo">
                  {{ dashboardStore.formattedStats.creator_total }}
                </h2>
                <small class="text-muted"
                  >活跃:
                  {{ dashboardStore.formattedStats.creator_active }}</small
                >
              </div>
              <div class="stat-icon">
                <i class="iconfont text-2xl icon-creator"></i>
              </div>
            </div>
          </div>
        </div>
      </div>

      <!-- 图表区域 -->
      <div class="col-12 col-lg-6">
        <div class="card">
          <div class="card-header">
            <div class="d-flex justify-content-between align-items-center">
              <h5 class="card-title mb-0">用户增长趋势</h5>
              <div class="btn-group btn-group-sm" role="group">
                <button
                  type="button"
                  class="btn"
                  :class="{
                    'btn-primary': dashboardStore.trendDays === 7,
                    'btn-outline-primary': dashboardStore.trendDays !== 7,
                  }"
                  @click="dashboardStore.changeTrendDays(7)"
                >
                  7天
                </button>
                <button
                  type="button"
                  class="btn"
                  :class="{
                    'btn-primary': dashboardStore.trendDays === 30,
                    'btn-outline-primary': dashboardStore.trendDays !== 30,
                  }"
                  @click="dashboardStore.changeTrendDays(30)"
                >
                  30天
                </button>
              </div>
            </div>
          </div>
          <div class="card-body">
            <div id="userTrendChart" style="height: 300px"></div>
          </div>
        </div>
      </div>

      <div class="col-12 col-lg-6">
        <div class="card">
          <div class="card-header">
            <div class="d-flex justify-content-between align-items-center">
              <h5 class="card-title mb-0">收入趋势</h5>
              <div class="btn-group btn-group-sm" role="group">
                <button
                  type="button"
                  class="btn"
                  :class="{
                    'btn-success': dashboardStore.trendDays === 7,
                    'btn-outline-success': dashboardStore.trendDays !== 7,
                  }"
                  @click="dashboardStore.changeTrendDays(7)"
                >
                  7天
                </button>
                <button
                  type="button"
                  class="btn"
                  :class="{
                    'btn-success': dashboardStore.trendDays === 30,
                    'btn-outline-success': dashboardStore.trendDays !== 30,
                  }"
                  @click="dashboardStore.changeTrendDays(30)"
                >
                  30天
                </button>
              </div>
            </div>
          </div>
          <div class="card-body">
            <div id="revenueTrendChart" style="height: 300px"></div>
          </div>
        </div>
      </div>

      <!-- 数据表格区域 -->
      <div class="col-12 col-lg-4">
        <div class="card">
          <div class="card-header">
            <h5 class="card-title mb-0">最近订单</h5>
          </div>
          <div class="card-body p-0">
            <div
              class="list-group list-group-flush"
              v-if="
                dashboardStore.recent.recent_orders &&
                dashboardStore.recent.recent_orders.length > 0
              "
            >
              <div
                class="list-group-item"
                v-for="order in dashboardStore.recent.recent_orders"
                :key="order.id"
              >
                <div class="d-flex justify-content-between align-items-start">
                  <div class="me-auto">
                    <div class="fw-bold">
                      #{{ order.order_no.substring(0, 8) }}...
                      <span class="text-muted">
                        <i
                          class="iconfont icon-copy cursor-pointer ml-1"
                          @click="copyText(order.order_no)"
                        ></i>
                      </span>
                    </div>
                    <small class="text-muted">{{ order.username }}</small>
                    <br />
                    <small class="text-muted">{{
                      dashboardStore.formatTime(order.pay_time)
                    }}</small>
                  </div>
                  <div class="text-end">
                    <span class="badge bg-success">{{ order.status }}</span>
                    <div class="text-primary fw-bold">￥{{ order.amount }}</div>
                  </div>
                </div>
              </div>
            </div>
            <div v-else class="text-center p-4">
              <el-empty description="暂无订单" :image-size="80" />
            </div>
          </div>
        </div>
      </div>

      <div class="col-12 col-lg-4">
        <div class="card">
          <div class="card-header">
            <h5 class="card-title mb-0">最近用户</h5>
          </div>
          <div class="card-body p-0">
            <div
              class="list-group list-group-flush"
              v-if="
                dashboardStore.recent.recent_users &&
                dashboardStore.recent.recent_users.length > 0
              "
            >
              <div
                class="list-group-item"
                v-for="user in dashboardStore.recent.recent_users"
                :key="user.id"
              >
                <div class="d-flex align-items-center">
                  <el-image
                    :src="user.avatar"
                    class="rounded-circle me-3 w-[40px] h-[40px]"
                    fit="cover"
                    alt="头像"
                  />
                  <div class="me-auto">
                    <div class="fw-bold">{{ user.nickname }}</div>
                    <small class="text-muted">{{ user.username }}</small>
                    <br />
                    <small class="text-muted">{{
                      dashboardStore.formatTime(user.created_at)
                    }}</small>
                  </div>
                  <span class="badge bg-info" v-if="user.platform">{{
                    user.platform
                  }}</span>
                </div>
              </div>
            </div>
            <div v-else class="text-center p-4">
              <el-empty description="暂无用户" :image-size="80" />
            </div>
          </div>
        </div>
      </div>

      <div class="col-12 col-lg-4">
        <div class="card">
          <div class="card-header">
            <h5 class="card-title mb-0">热门应用</h5>
          </div>
          <div class="card-body p-0">
            <div
              class="list-group list-group-flush"
              v-if="
                dashboardStore.recent.hot_apps &&
                dashboardStore.recent.hot_apps.length > 0
              "
            >
              <div
                class="list-group-item"
                v-for="(app, index) in dashboardStore.recent.hot_apps"
                :key="app.id"
              >
                <div class="d-flex align-items-center">
                  <div class="me-3">
                    <span
                      class="badge badge-rank"
                      :class="dashboardStore.getRankClass(index)"
                      >{{ index + 1 }}</span
                    >
                  </div>
                  <el-image
                    :src="app.icon || '/images/app-placeholder.png'"
                    class="rounded me-3 w-[40px] h-[40px]"
                    fit="cover"
                    alt="应用图标"
                  />
                  <div class="me-auto">
                    <div class="fw-bold">{{ app.name }}</div>
                    <small class="text-muted"
                      >使用 {{ app.use_count }} 次</small
                    >
                  </div>
                  <div class="text-end">
                    <span class="badge bg-warning">{{ app.score }} 积分</span>
                  </div>
                </div>
              </div>
            </div>
            <div v-else class="text-center p-4">
              <el-empty description="暂无应用" :image-size="80" />
            </div>
          </div>
        </div>
      </div>

      <!-- 系统状态 -->
      <div class="col-12">
        <div class="card">
          <div class="card-header">
            <h5 class="card-title mb-0">系统状态</h5>
          </div>
          <div class="card-body">
            <div class="row">
              <div class="col-12 col-md-3">
                <div class="d-flex align-items-center mb-2">
                  <i
                    class="fas fa-database me-2"
                    :class="
                      dashboardStore.system.database_status
                        ? 'text-success'
                        : 'text-danger'
                    "
                  ></i>
                  <span>数据库</span>
                  <span
                    class="ms-auto badge"
                    :class="
                      dashboardStore.system.database_status
                        ? 'bg-success'
                        : 'bg-danger'
                    "
                  >
                    {{
                      dashboardStore.system.database_status ? '正常' : '异常'
                    }}
                  </span>
                </div>
              </div>
              <div class="col-12 col-md-3">
                <div class="d-flex align-items-center mb-2">
                  <i
                    class="fas fa-server me-2"
                    :class="
                      dashboardStore.system.redis_status
                        ? 'text-success'
                        : 'text-danger'
                    "
                  ></i>
                  <span>Redis</span>
                  <span
                    class="ms-auto badge"
                    :class="
                      dashboardStore.system.redis_status
                        ? 'bg-success'
                        : 'bg-danger'
                    "
                  >
                    {{ dashboardStore.system.redis_status ? '正常' : '异常' }}
                  </span>
                </div>
              </div>
              <div class="col-12 col-md-3">
                <div class="d-flex align-items-center mb-3">
                  <i class="fas fa-microchip me-2 text-info"></i>
                  <span class="me-2">内存使用</span>
                  <div class="flex-grow-1 mx-2">
                    <el-progress
                      :percentage="
                        getUsagePercentage(dashboardStore.system.memory_usage)
                      "
                      :status="
                        getUsageStatus(dashboardStore.system.memory_usage)
                      "
                      :text-inside="true"
                      :stroke-width="20"
                      :border-radius="5"
                    />
                  </div>
                </div>
              </div>
              <div class="col-12 col-md-3">
                <div class="d-flex align-items-center mb-3">
                  <i class="fas fa-hdd me-2 text-warning"></i>
                  <span class="me-2">磁盘使用</span>
                  <div class="flex-grow-1 mx-2">
                    <el-progress
                      :percentage="
                        getUsagePercentage(dashboardStore.system.disk_usage)
                      "
                      :status="getUsageStatus(dashboardStore.system.disk_usage)"
                      :text-inside="true"
                      :stroke-width="20"
                      :border-radius="5"
                    />
                  </div>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
  import { useDashboardStore } from '@/js/store/admin/dashboard'
  import { ElMessage } from 'element-plus'
  import { onMounted, onUnmounted } from 'vue'

  // 使用 dashboard store
  const dashboardStore = useDashboardStore()

  // 生命周期
  onMounted(() => {
    // 初始化 dashboard
    dashboardStore.init()

    // 延迟初始化图表，确保DOM已渲染
    setTimeout(() => {
      dashboardStore.initCharts()
    }, 100)
  })

  onUnmounted(() => {
    // 清理资源
    dashboardStore.cleanup()
  })

  const copyText = (text) => {
    navigator.clipboard.writeText(text)
    ElMessage.success('复制成功')
  }

  const getUsagePercentage = (usage) => {
    if (usage === null || usage === undefined) {
      return 0
    }
    if (typeof usage === 'string') {
      usage = parseFloat(usage)
    }
    if (Number.isNaN(usage)) {
      return 0
    }
    return Math.min(100, Math.max(0, Math.round(usage)))
  }

  const getUsageStatus = (usage) => {
    if (usage === null || usage === undefined) {
      return '' // 默认颜色
    }
    if (typeof usage === 'string') {
      usage = parseFloat(usage)
    }
    if (Number.isNaN(usage)) {
      return ''
    }
    if (usage >= 80) {
      return 'exception' // 红色
    } else if (usage >= 50) {
      return 'warning' // 黄色
    } else {
      return 'success' // 绿色
    }
  }
</script>

<style scoped lang="scss">
  .dashboard {
    .stat-card {
      transition: all 0.3s ease;
      border: none;
      box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);

      &:hover {
        transform: translateY(-2px);
        box-shadow: 0 4px 8px rgba(0, 0, 0, 0.15);
      }

      .stat-icon {
        width: 60px;
        height: 60px;
        border-radius: 50%;
        display: flex;
        align-items: center;
        justify-content: center;
        background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
        color: white;
        font-size: 24px;
        opacity: 0.8;
      }
    }

    .text-purple {
      color: #6f42c1 !important;
    }

    .text-indigo {
      color: #6610f2 !important;
    }

    .badge-rank {
      font-size: 12px;
      padding: 4px 8px;
      border-radius: 50%;
      min-width: 24px;
      text-align: center;
    }

    .card {
      border: none;
      box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);

      .card-header {
        background-color: #f8f9fa;
        border-bottom: 1px solid #e9ecef;
      }
    }

    .list-group-item {
      border-left: none;
      border-right: none;

      &:first-child {
        border-top: none;
      }

      &:hover {
        background-color: #f8f9fa;
      }
    }
  }
</style>
