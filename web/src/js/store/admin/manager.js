import { defineStore } from 'pinia'
import { ref } from 'vue'
import {
  showConfirm,
  showLoading,
  showMessageError,
  showMessageOK,
} from '../../utils/dialog.js'
import { httpGet, httpPost } from '../../utils/http.js'
import { validateForm } from '../common.js'

const rules = {
  username: { required: true, message: '请输入用户名' },
  password: { required: true, message: '请输入密码' },
}

export const useAdminManagerStore = defineStore('admin-manager', () => {
  const items = ref([])
  const item = ref({})
  const showDialog = ref(false)
  const showResetPassDialog = ref(false)
  const loading = ref(false)
  const errors = ref({})

  const fetchData = async () => {
    loading.value = true
    try {
      const res = await httpGet('/api/admin/list')
      items.value = res.data || []
    } catch (error) {
      showMessageError('获取数据失败：' + error.message)
    } finally {
      loading.value = false
    }
  }

  const add = () => {
    showDialog.value = true
    item.value = {}
  }

  const handleResetPass = async () => {
    if (!validateForm(item.value, rules, errors.value)) {
      return
    }
    showLoading()
    try {
      await httpPost('/api/admin/resetPass', {
        id: item.value.id,
        password: item.value.password,
      })
      showResetPassDialog.value = false
      showMessageOK('操作成功')
    } catch (error) {
      showResetPassDialog.value = false
      showMessageError('操作失败：' + error.message)
    }
  }

  const handleSubmit = async () => {
    if (!validateForm(item.value, rules, errors.value)) {
      return
    }
    showLoading()
    try {
      await httpPost('/api/admin/save', item.value)
      showMessageOK('操作成功！')
      showDialog.value = false
      await fetchData()
    } catch (error) {
      showMessageError('操作失败：' + error.message)
      showDialog.value = false
    }
  }

  const enable = async (row) => {
    try {
      await httpPost('/api/admin/enable', { id: row.id, enabled: row.status })
      showMessageOK('操作成功！')
    } catch (error) {
      showMessageError('操作失败：' + error.message)
      row.status = !row.status
    }
  }

  const remove = (row) => {
    showConfirm('删除提示', '确定要删除当前记录吗?？', async () => {
      showLoading()
      try {
        await httpGet('/api/admin/remove?id=' + row.id)
        showMessageOK('删除成功！')
        await fetchData()
      } catch (error) {
        showMessageError('删除失败：' + error.message)
      }
    })
  }

  const initialize = async () => {
    await fetchData()
  }

  const handleSelectionChange = () => {}

  return {
    loading,
    item,
    items,
    errors,
    showDialog,
    showResetPassDialog,
    fetchData,
    initialize,
    handleSubmit,
    remove,
    enable,
    add,
    handleResetPass,
    handleSelectionChange,
  }
})
