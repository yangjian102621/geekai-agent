import ClipboardJS from 'clipboard'
import { defineStore } from 'pinia'
import { ref } from 'vue'
import { httpGet, httpPost, httpPostDownload } from '../../utils/http'
import {
  closeLoading,
  showConfirm,
  showLoading,
  showMessageError,
  showMessageOK,
  showMessageWarning,
} from '../../utils/dialog'
import { UUID } from '../../utils/libs'
import { validateForm } from '../common'

const DEFAULT_DATASETS = {
  page: 1,
  page_size: 10,
  total: 0,
  items: [],
}

const DEFAULT_ITEM = {
  name: '100积分点卡',
  amount: 100,
  num: 1,
}

const RULES = {
  name: { required: true, message: '请输入兑换码名称' },
  amount: { required: true, message: '请输入兑换额度' },
  num: { required: true, message: '请输入生成数量' },
}

const REDEEM_STATUS = [
  { value: -1, label: '全部' },
  { value: 0, label: '未核销' },
  { value: 1, label: '已核销' },
]

export const useAdminRedeemStore = defineStore('admin-redeem', () => {
  const dataSets = ref({ ...DEFAULT_DATASETS })
  const loading = ref(false)
  const query = ref({ code: '', status: -1, page: 1, page_size: 10 })
  const showDialog = ref(false)
  const item = ref({ ...DEFAULT_ITEM })
  const itemIds = ref([])
  const exporting = ref(false)
  const errors = ref({})
  const title = ref('生成兑换码')
  const clipboard = ref(null)

  const ensureClipboard = () => {
    if (clipboard.value) return
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

  const fetchData = async (page = 1) => {
    loading.value = true
    query.value.page = page
    query.value.page_size =
      dataSets.value.page_size || DEFAULT_DATASETS.page_size
    try {
      const res = await httpGet('/api/admin/redeem/list', query.value)
      if (res?.data) {
        dataSets.value = { ...DEFAULT_DATASETS, ...res.data }
      }
    } catch (error) {
      showMessageError('获取数据失败：' + error.message)
    } finally {
      loading.value = false
    }
  }

  const add = () => {
    item.value = { ...DEFAULT_ITEM }
    showDialog.value = true
    errors.value = {}
  }

  const handleSubmit = async () => {
    if (!validateForm(item.value, RULES, errors.value)) {
      return
    }
    showLoading()
    try {
      const res = await httpPost('/api/admin/redeem/create', item.value)
      showMessageOK(`成功生成了${res.data.counter}个兑换码`)
      showDialog.value = false
      await fetchData()
    } catch (error) {
      showMessageError('生成失败：' + error.message)
    } finally {
      closeLoading()
    }
  }

  const set = async (field, row) => {
    try {
      await httpPost('/api/admin/redeem/set', {
        id: row.id,
        filed: field,
        value: row[field],
      })
      showMessageOK('操作成功！')
    } catch (error) {
      showMessageError('操作失败：' + error.message)
    }
  }

  const remove = (row) => {
    showConfirm('删除提示', '确定要删除当前记录吗？', async () => {
      try {
        await httpGet('/api/admin/redeem/remove?id=' + row.id)
        showMessageOK('删除成功！')
        await fetchData()
      } catch (error) {
        showMessageError('删除失败：' + error.message)
      }
    })
  }

  const batchRemove = () => {
    if (itemIds.value.length === 0) {
      showMessageWarning('请先选择要删除的记录')
      return
    }
    showConfirm('批量删除', '确定要删除所有选中的记录吗？', async () => {
      try {
        await httpPost('/api/admin/redeem/batchRemove', { ids: itemIds.value })
        showMessageOK('删除成功！')
        await fetchData()
      } catch (error) {
        showMessageError('删除失败：' + error.message)
      }
    })
  }

  const handleSelectionChange = (items) => {
    itemIds.value = items.map((item) => item.id)
  }

  const exportItems = async () => {
    exporting.value = true
    try {
      const response = await httpPostDownload('/api/admin/redeem/export', {
        ids: itemIds.value,
        status: query.value.status,
      })
      const url = window.URL.createObjectURL(new Blob([response.data]))
      const link = document.createElement('a')
      link.href = url
      link.setAttribute('download', UUID() + '.csv')
      document.body.appendChild(link)
      link.click()
      document.body.removeChild(link)
      window.URL.revokeObjectURL(url)
    } catch (error) {
      showMessageError('下载失败')
    } finally {
      exporting.value = false
    }
  }

  const initialize = async () => {
    ensureClipboard()
    await fetchData()
  }

  return {
    dataSets,
    loading,
    query,
    showDialog,
    item,
    itemIds,
    exporting,
    errors,
    redeemStatus: REDEEM_STATUS,
    title,
    add,
    handleSubmit,
    set,
    fetchData,
    remove,
    handleSelectionChange,
    exportItems,
    batchRemove,
    initialize,
    ensureClipboard,
    releaseClipboard,
  }
})
