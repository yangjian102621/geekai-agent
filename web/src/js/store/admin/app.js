import ClipboardJS from 'clipboard'
import { ElMessageBox } from 'element-plus'
import { defineStore } from 'pinia'
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import {
  closeLoading,
  showConfirm,
  showLoading,
  showMessageError,
  showMessageOK,
} from '@/js/utils/dialog.js'
import { httpGet, httpPost } from '@/js/utils/http.js'
import { copyObj } from '@/js/utils/libs.js'
import { adminUploadFile, validateForm } from '@/js/store/common.js'

const createDefaultItem = () => ({
  configs: {},
  params: [],
  enabled: true,
  icon: '/images/app-placeholder.png',
  summary: '我是  Geek-Agent 智能体助手，有什么问题可以问我！',
  score: 1,
  billing_mode: 'immediate',
  billing_config: {
    suffixes: [],
    marker: '',
  },
})

export const useAppStore = defineStore('admin-app', () => {
  // state
  const item = ref({ icon: '/images/app-placeholder.png' })
  const showDialog = ref(false)
  const loading = ref(false)
  const clipboard = ref(null)
  const title = ref('添加应用')
  const dataSets = ref({ total: 0, page: 1, pageSize: 10, items: [] })
  const query = ref({})
  const rules = {
    name: { required: true, message: '请输入应用名称' },
    type: { required: true, message: '请选择应用类型' },
    icon: { required: true, message: '请上传应用图标' },
    score: { required: true, message: '请输入每次对话积分' },
  }
  const errors = ref({})
  const showCozeAgentDialog = ref(false)
  const cozeAgentList = ref([])
  const router = useRouter()
  const appCategories = ref([])
  const importCategoryId = ref(0)
  const appIds = ref([])
  const cozeAgentIds = ref([])

  const ensureClipboard = () => {
    if (clipboard.value) {
      return
    }
    clipboard.value = new ClipboardJS('.icon-copy')
    clipboard.value.on('success', () => {
      showMessageOK('复制成功！')
    })
    clipboard.value.on('error', () => {
      showMessageError('复制失败！')
    })
  }

  const releaseClipboard = () => {
    if (clipboard.value) {
      clipboard.value.destroy()
      clipboard.value = null
    }
  }

  // 获取数据
  const fetchData = async (page) => {
    loading.value = true
    dataSets.value.page = page || 1
    query.value.page = dataSets.value.page
    query.value.page_size = dataSets.value.pageSize
    try {
      const res = await httpGet('/api/admin/app/list', query.value)
      dataSets.value.items = []
      if (res.data.items) {
        res.data.items.forEach((item) => {
          if (!item.configs.api_url) item.configs.api_url = ''
          if (!item.configs.token) item.configs.token = ''
          if (item.type === 'openai') {
            item.configs.model_name = item.configs.model_name || 'gpt-3.5-turbo'
            item.configs.system_prompt = item.configs.system_prompt || ''
            item.configs.max_length = item.configs.max_length || 1024
            item.configs.max_context_length =
              item.configs.max_context_length || 32000
            item.configs.history_deep = item.configs.history_deep || 3
          }
          if (!item.params) {
            item.params = []
          }
          // 初始化扣费配置
          if (!item.billing_mode) {
            item.billing_mode = 'immediate'
          }
          if (!item.billing_config) {
            item.billing_config = {
              suffixes: [],
              marker: '',
              case_sensitive: false,
            }
          }
          dataSets.value.items.push(item)
        })
        dataSets.value.total = res.data.total
      }
    } catch (e) {
      showMessageError('获取数据失败：' + e.message)
    } finally {
      loading.value = false
    }
  }

  const fetchAppCategories = async () => {
    try {
      const res = await httpGet('/api/admin/app/category/list?system=1')
      appCategories.value = res.data || []
      if (appCategories.value.length > 0) {
        importCategoryId.value = appCategories.value[0].id
      }
    } catch (e) {
      // 静默失败即可
    }
  }

  const initialize = async () => {
    ensureClipboard()
    await fetchData()
    await fetchAppCategories()
  }

  // 新增
  const add = () => {
    showDialog.value = true
    item.value = createDefaultItem()
    title.value = '添加应用'
  }

  // 编辑
  const edit = (row) => {
    showDialog.value = true
    item.value = copyObj(row)
    if (!item.value.params) {
      item.value.params = []
    }
    // 初始化扣费配置
    if (!item.value.billing_mode) {
      item.value.billing_mode = 'immediate'
    }
    if (!item.value.billing_config) {
      item.value.billing_config = {
        suffixes: [],
        marker: '',
        case_sensitive: false,
      }
    }
    title.value = '编辑应用'
  }

  // 复制应用
  const copy = async (row) => {
    showLoading()
    try {
      await httpPost('/api/admin/app/copy', { id: row.id })
      showMessageOK('操作成功！')
      showDialog.value = false
      await fetchData()
    } catch (e) {
      showMessageError('操作失败，' + e.message)
      showDialog.value = false
    }
  }

  // 提交
  const handleSubmit = async () => {
    if (validateForm(item.value, rules, errors.value)) {
      showLoading()
      try {
        await httpPost('/api/admin/app/save', item.value)
        showMessageOK('操作成功！')
        showDialog.value = false
        await fetchData()
      } catch (e) {
        showMessageError('操作失败，' + e.message)
        showDialog.value = false
      }
    }
  }

  // 设置字段
  const setValue = async (row, key, value) => {
    try {
      await httpPost('/api/admin/app/set', {
        id: row.id,
        name: key,
        value: value,
      })
      showMessageOK('操作成功！')
    } catch (e) {
      showMessageError('操作失败：' + e.message)
    }
  }

  // 删除
  const remove = (row) => {
    showConfirm('删除提示', '确定要删除当前记录吗?？', async () => {
      showLoading()
      try {
        await httpGet('/api/admin/app/remove?id=' + row.id)
        showMessageOK('删除成功！')
        await fetchData()
      } catch (e) {
        showMessageError('删除失败：' + e.message)
      }
    })
  }

  // 上传图标
  const uploadIcon = (file) => {
    adminUploadFile(file, (data) => {
      item.value.icon = data.url
    })
  }

  // 多选
  const handleSelectionChange = (selection) => {
    appIds.value = []
    selection.forEach((row) => {
      appIds.value.push(row.id)
    })
  }

  // Coze 多选
  const handleSelectionCoze = (selection) => {
    cozeAgentIds.value = []
    selection.forEach((row) => {
      cozeAgentIds.value.push(row.bot_id)
    })
  }

  // 导入 Coze 智能体
  const importCozeAgents = async () => {
    showLoading('正在获取智能体列表...')
    try {
      const res = await httpGet('/api/admin/app/coze/agents')
      closeLoading()
      if (res.code === 402) {
        ElMessageBox.alert('请先配置 Coze API 信息', '获取智能体失败', {
          confirmButtonText: '确定',
          callback: () => {
            router.push('/admin/settings/coze')
          },
        })
      } else {
        cozeAgentList.value = res.data
        showCozeAgentDialog.value = true
      }
    } catch (e) {
      closeLoading()
      showMessageError('获取智能体列表失败：' + e.message)
    }
  }

  // 执行导入
  const doImportCozeAgents = async () => {
    showLoading()
    const selectedAgents = cozeAgentList.value.filter((agent) =>
      cozeAgentIds.value.includes(agent.bot_id)
    )
    if (selectedAgents.length === 0) {
      closeLoading()
      showMessageError('请选择要导入的智能体')
      return
    }
    const agents = selectedAgents.map((agent) => ({
      bot_id: agent.bot_id,
      bot_name: agent.bot_name,
      description: agent.description,
      icon: agent.icon_url,
      cid: importCategoryId.value,
    }))
    try {
      await httpPost('/api/admin/app/coze/import', { agents })
      showMessageOK('导入成功！')
      showCozeAgentDialog.value = false
      await fetchData()
    } catch (e) {
      showMessageError('导入失败：' + e.message)
    } finally {
      closeLoading()
    }
  }

  // 批量删除
  const batchDelete = () => {
    if (appIds.value.length === 0) {
      showMessageError('请选择要删除的应用')
      return
    }
    showConfirm(
      '删除提示',
      '确定要删除选中的 ' + appIds.value.length + ' 条记录吗？',
      async () => {
        showLoading()
        try {
          await httpPost('/api/admin/app/batch-remove', { ids: appIds.value })
          showMessageOK('删除成功！')
          await fetchData()
        } catch (e) {
          showMessageError('删除失败：' + e.message)
        }
      }
    )
  }

  // 切换类型
  const changeAppType = (value) => {
    if (value === 'coze') {
      item.value.configs.api_url = 'https://api.coze.cn'
    } else if (value === 'bailian') {
      item.value.configs.api_url = ''
      if (!item.value.configs.bailian_api_key) {
        item.value.configs.bailian_api_key = ''
      }
      if (!item.value.configs.bailian_app_id) {
        item.value.configs.bailian_app_id = ''
      }
    } else {
      item.value.configs.api_url = ''
    }
  }

  return {
    loading,
    item,
    errors,
    showDialog,
    title,
    showCozeAgentDialog,
    cozeAgentList,
    dataSets,
    query,
    appCategories,
    importCategoryId,
    handleSubmit,
    remove,
    setValue,
    add,
    edit,
    copy,
    handleSelectionChange,
    handleSelectionCoze,
    uploadIcon,
    importCozeAgents,
    doImportCozeAgents,
    batchDelete,
    fetchData,
    fetchAppCategories,
    initialize,
    releaseClipboard,
    changeAppType,
  }
})
