import { ref } from 'vue'
import { defineStore } from 'pinia'
import {
  showMessageError,
  showMessageOK,
  showMessageWarning,
} from '@/js/utils/dialog.js'
import { httpGet, httpPost } from '@/js/utils/http.js'

export const useWorkflowStore = defineStore('workflow', () => {
  const workflows = ref([])
  const tasks = ref([])
  const loadingTasks = ref(false)
  const pollingTimer = ref(null)
  const initialized = ref(false)

  const fetchWorkflows = async () => {
    try {
      const res = await httpGet('/api/workflow/list')
      workflows.value = res.data || []
    } catch (err) {
      showMessageWarning('获取工作流失败：' + err.message)
    }
  }

  const fetchTasks = async (silent = false) => {
    if (!silent) {
      loadingTasks.value = true
    }
    try {
      const res = await httpGet('/api/workflow/tasks')
      // 后端返回的是 Page 对象，包含 items, page, page_size, total, total_page
      tasks.value = res.data?.items || []
    } catch (err) {
      if (!silent) {
        showMessageError('获取任务失败：' + err.message)
      }
    } finally {
      if (!silent) {
        loadingTasks.value = false
      }
    }
  }

  const ensureInit = async () => {
    if (!initialized.value) {
      await Promise.all([fetchWorkflows(), fetchTasks()])
      initialized.value = true
    }
  }

  const createTask = async (workflowId, params) => {
    const res = await httpPost('/api/workflow/tasks', {
      workflow_id: workflowId,
      params,
    })
    const task = res.data
    tasks.value = [
      task,
      ...tasks.value.filter((item) => item.task_id !== task.task_id),
    ]
    showMessageOK('任务已创建，正在排队执行')
    return task
  }

  const refreshTask = async (taskId) => {
    try {
      const res = await httpGet(`/api/workflow/tasks/${taskId}`)
      const index = tasks.value.findIndex((item) => item.task_id === taskId)
      if (index >= 0) {
        tasks.value[index] = res.data
      } else {
        tasks.value.unshift(res.data)
      }
    } catch (err) {
      console.warn(err)
    }
  }

  const cancelTask = async (taskId) => {
    await httpPost(`/api/workflow/tasks/${taskId}/cancel`, {})
    showMessageOK('任务已取消')
    refreshTask(taskId)
  }

  const retryTask = async (taskId) => {
    await httpPost(`/api/workflow/tasks/${taskId}/retry`, {})
    showMessageOK('任务已重新排队')
    refreshTask(taskId)
  }

  const startPolling = () => {
    if (pollingTimer.value) return
    pollingTimer.value = setInterval(() => {
      fetchTasks(true)
    }, 5000)
  }

  const stopPolling = () => {
    if (pollingTimer.value) {
      clearInterval(pollingTimer.value)
      pollingTimer.value = null
    }
  }

  return {
    workflows,
    tasks,
    loadingTasks,
    ensureInit,
    fetchTasks,
    fetchWorkflows,
    createTask,
    refreshTask,
    cancelTask,
    retryTask,
    startPolling,
    stopPolling,
  }
})
