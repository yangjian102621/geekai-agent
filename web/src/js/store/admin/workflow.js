import ClipboardJS from 'clipboard'
import { defineStore } from 'pinia'
import { ref } from 'vue'
import {
  showConfirm,
  showLoading,
  showMessageError,
  showMessageOK,
  closeLoading,
} from '../../utils/dialog.js'
import { httpGet, httpPost } from '../../utils/http.js'
import { copyObj } from '../../utils/libs.js'
import { adminUploadFile, validateForm } from '../common.js'

const createDefaultItem = () => ({
  type: 'coze',
  params: [],
  options: {},
  enabled: true,
  icon: '/images/app-placeholder.png',
  summary: '',
  auth_config: {
    api_url: '',
    app_id: '',
    public_key_id: '',
    private_key: '',
  },
  bailian_auth_config: {
    api_key: '',
    app_id: '',
  },
})

const FORM_RULES = {
  name: { required: true, message: '请输入工作流名称' },
  workflow_id: { required: true, message: '请输入工作流ID' },
  score: { required: true, message: '请输入消耗积分' },
  summary: { required: true, message: '请输入工作流简介' },
  icon: { required: true, message: '请上传工作流图标' },
}

export const useWorkflowStore = defineStore('admin-workflow', () => {
  const item = ref(createDefaultItem())
  const showDialog = ref(false)
  const loading = ref(false)
  const clipboard = ref(null)
  const title = ref('添加工作流')
  const dataSets = ref({ total: 0, page: 1, pageSize: 10, items: [] })
  const query = ref({})
  const errors = ref({})
  const selectedIds = ref([])

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
    if (!clipboard.value) {
      return
    }
    clipboard.value.destroy()
    clipboard.value = null
  }

  const fetchData = async (page) => {
    loading.value = true
    dataSets.value.page = page || 1
    query.value.page = dataSets.value.page
    query.value.page_size = dataSets.value.pageSize
    try {
      const res = await httpGet('/api/admin/workflow/list', query.value)
      dataSets.value.items = []
      if (res.data.items) {
        res.data.items.forEach((workflow) => {
          if (!workflow.params) {
            workflow.params = []
          }
          if (!workflow.icon) {
            workflow.icon = '/images/app-placeholder.png'
          }
          if (!workflow.options) {
            workflow.options = {}
          }
          dataSets.value.items.push(workflow)
        })
        dataSets.value.total = res.data.total
      }
    } catch (e) {
      showMessageError('获取数据失败：' + e.message)
    } finally {
      loading.value = false
    }
  }

  const add = () => {
    showDialog.value = true
    item.value = createDefaultItem()
    title.value = '添加工作流'
  }

  const edit = (row) => {
    showDialog.value = true
    item.value = copyObj(row)
    if (!item.value.type) {
      item.value.type = 'coze'
    }
    if (!item.value.params) {
      item.value.params = []
    }
    if (!item.value.options) {
      item.value.options = {}
    }
    if (!item.value.icon) {
      item.value.icon = '/images/app-placeholder.png'
    }
    if (!item.value.auth_config) {
      item.value.auth_config = {
        api_url: '',
        app_id: '',
        public_key_id: '',
        private_key: '',
      }
    }
    if (!item.value.bailian_auth_config) {
      item.value.bailian_auth_config = {
        api_key: '',
        app_id: '',
      }
    }
    title.value = '编辑工作流'
  }

  const handleSubmit = async () => {
    if (!validateForm(item.value, FORM_RULES, errors.value)) {
      console.log(errors.value)
      return
    }
    console.log(item.value)
    showLoading()
    try {
      await httpPost('/api/admin/workflow/save', item.value)
      showMessageOK('操作成功！')
      showDialog.value = false
      await fetchData()
    } catch (e) {
      showMessageError('操作失败，' + e.message)
      showDialog.value = false
    }
  }

  const enable = async (row) => {
    try {
      await httpPost('/api/admin/workflow/enable', {
        id: row.id,
        enabled: row.enabled,
      })
      showMessageOK('操作成功！')
    } catch (e) {
      showMessageError('操作失败：' + e.message)
    }
  }

  const remove = (row) => {
    showConfirm('删除提示', '确定要删除当前记录吗?？', async () => {
      showLoading()
      try {
        await httpGet('/api/admin/workflow/remove?id=' + row.id)
        showMessageOK('删除成功！')
        await fetchData()
      } catch (e) {
        showMessageError('删除失败：' + e.message)
      }
    })
  }

  const handleSelectionChange = (selection) => {
    selectedIds.value = selection.map((row) => row.id)
  }

  const batchDelete = () => {
    if (selectedIds.value.length === 0) {
      showMessageError('请选择要删除的工作流')
      return
    }

    showConfirm(
      '删除提示',
      '确定要删除选中的 ' + selectedIds.value.length + ' 条记录吗？',
      async () => {
        showLoading()
        try {
          await httpPost('/api/admin/workflow/batch-remove', {
            ids: selectedIds.value,
          })
          showMessageOK('删除成功！')
          await fetchData()
        } catch (e) {
          showMessageError('删除失败：' + e.message)
        }
      }
    )
  }

  const changeWorkflowType = (value) => {
    if (value === 'bailian') {
      if (!item.value.bailian_auth_config) {
        item.value.bailian_auth_config = {
          api_key: '',
          app_id: '',
        }
      }
    } else {
      if (!item.value.auth_config) {
        item.value.auth_config = {
          api_url: '',
          app_id: '',
          public_key_id: '',
          private_key: '',
        }
      }
    }
  }

  const initialize = async () => {
    ensureClipboard()
    await fetchData()
  }

  const uploadIcon = (file) => {
    adminUploadFile(file, (data) => {
      item.value.icon = data.url
    })
  }

  return {
    loading,
    item,
    errors,
    showDialog,
    title,
    dataSets,
    handleSubmit,
    remove,
    enable,
    add,
    edit,
    handleSelectionChange,
    batchDelete,
    fetchData,
    initialize,
    releaseClipboard,
    uploadIcon,
    changeWorkflowType,
  }
})
