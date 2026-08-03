// * +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
// * Copyright 2023 The Geek-AI Authors. All rights reserved.
// * Use of this source code is governed by a Apache-2.0 license
// * that can be found in the LICENSE file.
// * @Author yangjian102621@163.com
// * +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++

import { closeLoading, showConfirm, showLoading } from '@/js/utils/dialog'
import { httpGet, httpPost } from '@/js/utils/http'
import { dateFormat } from '@/js/utils/libs'
import { ElMessage } from 'element-plus'
import { defineStore } from 'pinia'
import { ref } from 'vue'

const createDefaultAppForm = () => ({
  cid: '',
  name: '',
  icon: '/images/app-placeholder.png',
  summary: '',
  enabled: true,
  score: 10,
  params: [],
  configs: {
    api_url: '',
    token: '',
    bot_id: '',
    app_id: '',
    public_key_id: '',
    private_key: '',
    model_name: '',
    max_length: 2048,
    enable_context: false,
    history_deep: 5,
    max_context_length: 16384,
    system_prompt: '',
  },
})

const createDefaultWithdrawForm = () => ({})

const createDefaultCategoryForm = () => ({
  enabled: true,
})

// 创作者控制台业务逻辑代码

export const useConsoleStore = defineStore('creator-console', () => {
  // state
  const creator = ref({})
  const appsData = ref({
    items: [],
    total: 0,
    page: 1,
    page_size: 10,
  })
  const appsLoading = ref(false)
  const appQuery = ref({ name: '', cid: '' })

  const appDialog = ref({
    visible: false,
    form: createDefaultAppForm(),
  })
  const withdrawDialog = ref({
    visible: false,
    form: createDefaultWithdrawForm(),
  })
  const scoreLogsDialog = ref({ visible: false })
  const withdrawLogsDialog = ref({ visible: false })
  const profileDialog = ref({
    visible: false,
    form: {},
  })

  const categories = ref([])
  const categoryDialog = ref({
    visible: false,
    form: createDefaultCategoryForm(),
  })
  const categoryLoading = ref(false)

  // actions
  async function fetchDashboard() {
    try {
      const res = await httpGet('/api/creator/info')
      creator.value = res.data
    } catch (error) {
      ElMessage.error('获取数据失败：' + error.message)
    }
  }

  async function fetchApps(page = 1) {
    appsLoading.value = true
    try {
      const params = {
        page,
        page_size: appsData.value.page_size,
        name: appQuery.value.name,
        cid: appQuery.value.cid,
      }
      const res = await httpGet('/api/creator/apps', params)
      appsData.value = res.data
    } catch (error) {
      ElMessage.error('获取应用列表失败：' + error.message)
    } finally {
      appsLoading.value = false
    }
  }

  function openAppDialog(app = null) {
    appDialog.value.visible = true
    const form = app ? { ...app } : createDefaultAppForm()
    form.params = form.params || []
    appDialog.value.form = form
  }

  // 启用/禁用应用
  async function enableApp(app) {
    try {
      await httpPost('/api/creator/apps/enable', {
        id: app.id,
        enabled: app.enabled,
      })
      ElMessage.success('操作成功')
    } catch (error) {
      ElMessage.error('操作失败：' + error.message)
      app.enabled = !app.enabled
    }
  }

  // 删除应用
  async function deleteApp(app) {
    showConfirm('删除提示', '确定要删除当前记录吗?？', async () => {
      showLoading()
      try {
        await httpGet('/api/creator/apps/remove', {
          id: app.id,
        })
        ElMessage.success('删除成功')
        fetchApps()
        // 更新创作者信息
        fetchDashboard()
      } catch (error) {
        ElMessage.error('删除失败：' + error.message)
      } finally {
        closeLoading()
      }
    })
  }

  function showWithdrawDialog() {
    if (!creator.value.scores || creator.value.scores <= 0) {
      ElMessage.warning('暂无可提现余额')
      return
    }
    withdrawDialog.value.visible = true
    withdrawDialog.value.form = createDefaultWithdrawForm()
  }

  function showScoreLogsDialog() {
    scoreLogsDialog.value.visible = true
  }

  function showWithdrawLogsDialog() {
    withdrawLogsDialog.value.visible = true
  }

  function showProfileDialog() {
    profileDialog.value.visible = true
    profileDialog.value.form = { ...(creator.value?.creator || {}) }
  }

  function submitProfile() {
    profileDialog.value.visible = false
    fetchDashboard()
  }

  async function fetchCategories() {
    categoryLoading.value = true
    try {
      const res = await httpGet('/api/creator/app-categories/list', {
        creator_id: creator.value.id,
      })
      categories.value = res.data || []
    } catch (e) {
      ElMessage.error('获取分类失败：' + e.message)
    } finally {
      categoryLoading.value = false
    }
  }

  function showCategoryDialog(category = null) {
    if (category) {
      categoryDialog.value.visible = true
      categoryDialog.value.form = { ...category }
    } else {
      categoryDialog.value.visible = true
      categoryDialog.value.form = createDefaultCategoryForm()
    }
  }

  async function submitCategory() {
    if (!categoryDialog.value.form.name) {
      ElMessage.warning('请输入分类名称')
      return
    }
    showLoading()
    try {
      if (categoryDialog.value.form.id) {
        await httpPost(
          `/api/creator/app-categories/update`,
          categoryDialog.value.form
        )
        ElMessage.success('更新成功')
      } else {
        await httpPost(
          `/api/creator/app-categories/create`,
          categoryDialog.value.form
        )
        ElMessage.success('创建成功')
      }
      categoryDialog.value.visible = false
      fetchCategories()
    } catch (e) {
      ElMessage.error('操作失败：' + e.message)
    } finally {
      closeLoading()
    }
  }

  async function deleteCategory(category) {
    showLoading()
    try {
      await httpPost('/api/creator/app-categories/delete', { id: category.id })
      ElMessage.success('删除成功')
      fetchCategories()
    } catch (e) {
      ElMessage.error('删除失败：' + e.message)
    } finally {
      closeLoading()
    }
  }

  function createAppSuccess() {
    fetchApps()
    // 更新创作者信息
    fetchDashboard()
    appDialog.value.visible = false
  }

  function withdrawSuccess() {
    fetchDashboard()
    withdrawDialog.value.visible = false
  }

  const getAppType = (type) => {
    switch (type) {
      case 'openai':
        return '通用大模型'
      case 'dify':
        return 'Dify智能体'
      case 'coze':
        return 'Coze智能体'
      case 'aliyun':
        return '百炼智能体'
      default:
        return '未知'
    }
  }

  // 额外导出 dateFormat 工具
  return {
    // state
    creator,
    appsData,
    appsLoading,
    appQuery,
    appDialog,
    withdrawDialog,
    scoreLogsDialog,
    withdrawLogsDialog,
    profileDialog,
    categories,
    categoryDialog,
    categoryLoading,
    // actions
    fetchDashboard,
    fetchApps,
    showAppDialog: () => openAppDialog(),
    editApp: (app) => openAppDialog(app),
    showWithdrawDialog,
    showScoreLogsDialog,
    showWithdrawLogsDialog,
    withdrawSuccess,
    showProfileDialog,
    fetchCategories,
    showCategoryDialog,
    submitCategory,
    deleteCategory,
    createAppSuccess,
    submitProfile,
    dateFormat,
    getAppType,
    enableApp,
    deleteApp,
  }
})
