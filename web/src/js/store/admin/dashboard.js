// * +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
// * Copyright 2023 The Geek-AI Authors. All rights reserved.
// * Use of this source code is governed by a Apache-2.0 license
// * that can be found in the LICENSE file.
// * @Author yangjian102621@163.com
// * +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++

import * as echarts from 'echarts'
import { ElMessage } from 'element-plus'
import { defineStore } from 'pinia'
import { computed, ref } from 'vue'
import { httpGet } from '@/js/utils/http'

export const useDashboardStore = defineStore('dashboard', () => {
  // 响应式数据
  const stats = ref({
    user_total: 0,
    user_today_new: 0,
    user_active: 0,
    app_total: 0,
    app_enabled: 0,
    chat_total: 0,
    chat_today: 0,
    order_total: 0,
    order_today: 0,
    revenue_total: 0,
    revenue_today: 0,
    score_consumed: 0,
    score_recharged: 0,
    creator_total: 0,
    creator_active: 0,
  })

  const trends = ref({
    user_trend: [],
    revenue_trend: [],
    chat_trend: [],
  })

  const recent = ref({
    recent_orders: [],
    recent_users: [],
    hot_apps: [],
  })

  const system = ref({
    database_status: false,
    redis_status: false,
    system_load: '',
    memory_usage: '',
    disk_usage: '',
  })

  const trendDays = ref(7)
  const refreshTimer = ref(null)
  const userChart = ref(null)
  const revenueChart = ref(null)

  // 计算属性
  const formattedStats = computed(() => ({
    user_total: formatNumber(stats.value.user_total),
    user_today_new: formatNumber(stats.value.user_today_new),
    user_active: formatNumber(stats.value.user_active),
    app_total: formatNumber(stats.value.app_total),
    app_enabled: formatNumber(stats.value.app_enabled),
    chat_total: formatNumber(stats.value.chat_total),
    chat_today: formatNumber(stats.value.chat_today),
    order_total: formatNumber(stats.value.order_total),
    order_today: formatNumber(stats.value.order_today),
    revenue_total: formatMoney(stats.value.revenue_total),
    revenue_today: formatMoney(stats.value.revenue_today),
    score_consumed: formatNumber(stats.value.score_consumed),
    score_recharged: formatNumber(stats.value.score_recharged),
    creator_total: formatNumber(stats.value.creator_total),
    creator_active: formatNumber(stats.value.creator_active),
  }))

  // 加载数据
  const loadData = async () => {
    try {
      // 并行加载所有数据
      const [statsRes, trendsRes, recentRes, systemRes] = await Promise.all([
        getDashboardStats(),
        getDashboardTrends(trendDays.value),
        getDashboardRecent(),
        getDashboardSystem(),
      ])

      if (statsRes.data) {
        stats.value = statsRes.data
      }
      if (trendsRes.data) {
        trends.value = trendsRes.data
        updateCharts()
      }
      if (recentRes.data) {
        recent.value = recentRes.data
      }
      if (systemRes.data) {
        system.value = systemRes.data
      }
    } catch (e) {
      console.error('加载Dashboard数据失败:', e)
      ElMessage.error('加载数据失败')
    }
  }

  // 获取仪表盘统计数据
  const getDashboardStats = () => {
    return httpGet('/api/admin/dashboard/stats')
  }

  // 获取趋势数据
  const getDashboardTrends = (days = 7) => {
    return httpGet(`/api/admin/dashboard/trends?days=${days}`)
  }

  // 获取最近数据
  const getDashboardRecent = () => {
    return httpGet('/api/admin/dashboard/recent')
  }

  // 获取系统状态
  const getDashboardSystem = () => {
    return httpGet('/api/admin/dashboard/system')
  }

  // 切换趋势天数
  const changeTrendDays = async (days) => {
    trendDays.value = days
    try {
      const res = await getDashboardTrends(days)
      if (res.data) {
        trends.value = res.data
        updateCharts()
      }
    } catch (e) {
      console.error('加载趋势数据失败:', e)
      ElMessage.error('加载趋势数据失败')
    }
  }

  // 初始化图表
  const initCharts = () => {
    // 用户趋势图表
    userChart.value = echarts.init(document.getElementById('userTrendChart'))
    revenueChart.value = echarts.init(
      document.getElementById('revenueTrendChart')
    )

    // 监听窗口大小变化
    const handleResize = () => {
      userChart.value?.resize()
      revenueChart.value?.resize()
    }
    window.addEventListener('resize', handleResize)

    updateCharts()

    // 返回清理函数
    return () => {
      window.removeEventListener('resize', handleResize)
    }
  }

  // 更新图表
  const updateCharts = () => {
    if (!userChart.value || !revenueChart.value) return

    // 用户趋势图表配置
    const userOption = {
      title: {
        show: false,
      },
      tooltip: {
        trigger: 'axis',
        formatter: '{b}<br/>新增用户: {c} 人',
      },
      grid: {
        left: '3%',
        right: '4%',
        bottom: '3%',
        containLabel: true,
      },
      xAxis: {
        type: 'category',
        boundaryGap: false,
        data: trends.value.user_trend.map((item) => item.date),
      },
      yAxis: {
        type: 'value',
        minInterval: 1,
      },
      series: [
        {
          name: '新增用户',
          type: 'line',
          smooth: true,
          areaStyle: {},
          itemStyle: {
            color: '#007bff',
          },
          areaStyle: {
            color: {
              type: 'linear',
              x: 0,
              y: 0,
              x2: 0,
              y2: 1,
              colorStops: [
                {
                  offset: 0,
                  color: 'rgba(0, 123, 255, 0.3)',
                },
                {
                  offset: 1,
                  color: 'rgba(0, 123, 255, 0.1)',
                },
              ],
            },
          },
          data: trends.value.user_trend.map((item) => item.value),
        },
      ],
    }

    // 收入趋势图表配置
    const revenueOption = {
      title: {
        show: false,
      },
      tooltip: {
        trigger: 'axis',
        formatter: '{b}<br/>收入: ￥{c}',
      },
      grid: {
        left: '3%',
        right: '4%',
        bottom: '3%',
        containLabel: true,
      },
      xAxis: {
        type: 'category',
        boundaryGap: false,
        data: trends.value.revenue_trend.map((item) => item.date),
      },
      yAxis: {
        type: 'value',
      },
      series: [
        {
          name: '收入',
          type: 'line',
          smooth: true,
          areaStyle: {},
          itemStyle: {
            color: '#28a745',
          },
          areaStyle: {
            color: {
              type: 'linear',
              x: 0,
              y: 0,
              x2: 0,
              y2: 1,
              colorStops: [
                {
                  offset: 0,
                  color: 'rgba(40, 167, 69, 0.3)',
                },
                {
                  offset: 1,
                  color: 'rgba(40, 167, 69, 0.1)',
                },
              ],
            },
          },
          data: trends.value.revenue_trend.map((item) => item.value),
        },
      ],
    }

    userChart.value.setOption(userOption)
    revenueChart.value.setOption(revenueOption)
  }

  // 格式化数字
  const formatNumber = (num) => {
    if (num >= 10000) {
      return (num / 10000).toFixed(1) + 'w'
    } else if (num >= 1000) {
      return (num / 1000).toFixed(1) + 'k'
    }
    return num?.toString() || '0'
  }

  // 格式化金额
  const formatMoney = (money) => {
    return parseFloat(money || 0).toFixed(2)
  }

  // 格式化时间
  const formatTime = (timestamp) => {
    if (!timestamp) return ''
    const date = new Date(timestamp * 1000)
    const now = new Date()
    const diff = now - date

    if (diff < 60000) return '刚刚'
    if (diff < 3600000) return Math.floor(diff / 60000) + '分钟前'
    if (diff < 86400000) return Math.floor(diff / 3600000) + '小时前'
    if (diff < 604800000) return Math.floor(diff / 86400000) + '天前'

    return date.toLocaleDateString()
  }

  // 获取排名样式
  const getRankClass = (index) => {
    if (index === 0) return 'bg-warning'
    if (index === 1) return 'bg-secondary'
    if (index === 2) return 'bg-dark'
    return 'bg-light text-dark'
  }

  // 开始自动刷新
  const startAutoRefresh = () => {
    refreshTimer.value = setInterval(() => {
      loadData()
    }, 30000) // 30秒刷新一次
  }

  // 停止自动刷新
  const stopAutoRefresh = () => {
    if (refreshTimer.value) {
      clearInterval(refreshTimer.value)
      refreshTimer.value = null
    }
  }

  // 清理图表
  const disposeCharts = () => {
    if (userChart.value) {
      userChart.value.dispose()
      userChart.value = null
    }
    if (revenueChart.value) {
      revenueChart.value.dispose()
      revenueChart.value = null
    }
  }

  // 初始化
  const init = () => {
    loadData()
    startAutoRefresh()
  }

  // 清理
  const cleanup = () => {
    stopAutoRefresh()
    disposeCharts()
  }

  return {
    // 状态
    stats,
    trends,
    recent,
    system,
    trendDays,

    // 计算属性
    formattedStats,

    // 方法
    loadData,
    changeTrendDays,
    initCharts,
    updateCharts,
    formatNumber,
    formatMoney,
    formatTime,
    getRankClass,
    startAutoRefresh,
    stopAutoRefresh,
    disposeCharts,
    init,
    cleanup,
  }
})
