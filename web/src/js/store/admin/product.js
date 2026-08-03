import { defineStore } from 'pinia'
import { ref } from 'vue'
import {
  closeLoading,
  showConfirm,
  showLoading,
  showMessageError,
  showMessageOK,
} from '../../utils/dialog.js'
import { httpGet, httpPost } from '../../utils/http.js'
import { validateForm } from '../common.js'

const createDefaultForm = () => ({
  name: '',
  price: 0,
  credit: 0,
  enabled: true,
})

const rules = {
  name: { required: true, message: '请输入产品名称' },
  price: { required: true, message: '请输入原价' },
  credit: { required: true, message: '请输入额度' },
}

export const useAdminProductStore = defineStore('admin-product', () => {
  const products = ref([])
  const dialogVisible = ref(false)
  const dialogTitle = ref('')
  const errors = ref({})
  const form = ref(createDefaultForm())

  const fetchProducts = async () => {
    showLoading()
    try {
      const res = await httpGet('/api/admin/product/list')
      products.value = res.data || []
    } catch (error) {
      showMessageError('获取产品列表失败：' + error.message)
    } finally {
      closeLoading()
    }
  }

  const handleAdd = () => {
    form.value = createDefaultForm()
    dialogTitle.value = '新增产品'
    dialogVisible.value = true
  }

  const handleEdit = (row) => {
    form.value = { ...row }
    dialogTitle.value = '编辑产品'
    dialogVisible.value = true
  }

  const handleDelete = (row) => {
    showConfirm('删除提示', '确认删除该产品？', async () => {
      showLoading()
      try {
        await httpGet(`/api/admin/product/remove`, {
          id: row.id,
        })
        showMessageOK('删除成功')
        await fetchProducts()
      } catch (error) {
        showMessageError('删除失败：' + error.message)
      } finally {
        closeLoading()
      }
    })
  }

  const handleStatusChange = async (row) => {
    showLoading()
    try {
      await httpPost(`/api/admin/product/enable`, {
        id: row.id,
        enabled: row.enabled,
      })
      showMessageOK('状态更新成功')
    } catch (error) {
      row.enabled = !row.enabled
      showMessageError('状态更新失败：' + error.message)
    } finally {
      closeLoading()
    }
  }

  const handleSubmit = async () => {
    if (!validateForm(form.value, rules, errors.value)) {
      return
    }
    showLoading()
    try {
      await httpPost('/api/admin/product/save', form.value)
      showMessageOK('保存成功')
      dialogVisible.value = false
      await fetchProducts()
    } catch (error) {
      showMessageError('保存失败：' + error.message)
    } finally {
      closeLoading()
    }
  }

  const handleSort = async (ids, sorts) => {
    showLoading()
    try {
      await httpPost('/api/admin/product/sort', { ids, sorts })
      showMessageOK('排序成功')
    } catch (error) {
      showMessageError('排序失败：' + error.message)
    } finally {
      closeLoading()
    }
  }

  const initialize = async () => {
    await fetchProducts()
  }

  return {
    products,
    dialogVisible,
    dialogTitle,
    form,
    errors,
    handleAdd,
    handleEdit,
    handleDelete,
    handleStatusChange,
    handleSubmit,
    fetchProducts,
    handleSort,
    initialize,
  }
})
