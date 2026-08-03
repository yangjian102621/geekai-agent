import { defineStore } from 'pinia'
import { ref } from 'vue'
import {
  showConfirm,
  showLoading,
  showMessageError,
  showMessageOK,
} from '../../utils/dialog.js'
import { httpGet, httpPost } from '../../utils/http.js'
import { dateFormat } from '../../utils/libs.js'
import { validateEmail, validateMobile } from '../../utils/validate.js'
import { validateForm } from '../common.js'

const createDefaultItem = () => ({
  enabled: true,
})

export const useAdminUserStore = defineStore('admin-user', () => {
  const dataSets = ref({ total: 0, page: 1, pageSize: 10, items: [] })
  const item = ref(createDefaultItem())
  const showDialog = ref(false)
  const showResetPassDialog = ref(false)
  const loading = ref(false)
  const query = ref({})
  const errors = ref({})
  const userIds = ref([])

  const rules = {
    username: {
      required: true,
      message: '请输入用户名',
      validator: (value) => {
        if (
          item.value.id > 0 ||
          validateMobile(value) ||
          validateEmail(value)
        ) {
          return true
        }
        errors.value.username = '用户名必须是手机号或者邮箱地址'
        return false
      },
    },
    password: {
      required: true,
      message: '请输入密码',
      min: 8,
      max: 16,
      validator: (value) => {
        if (item.value.id > 0) {
          return true
        }
        if (value.length < 8 || value.length > 16) {
          errors.value.password = '密码长度必须为8-16位'
          return false
        }
        return true
      },
    },
    nickname: { required: true, message: '请输入昵称' },
  }

  const fetchData = async (page = 1) => {
    loading.value = true
    dataSets.value.page = page
    query.value.page = dataSets.value.page
    query.value.page_size = dataSets.value.pageSize
    try {
      const res = await httpGet('/api/admin/user/list', query.value)
      if (res.data) {
        dataSets.value.items = res.data.items.map((record) => ({
          ...record,
          expired_time: dateFormat(record.expired_time, 'yyyy-MM-dd'),
        }))
        dataSets.value.total = res.data.total
      }
    } catch (error) {
      showMessageError('加载用户列表失败：' + error.message)
    } finally {
      loading.value = false
    }
  }

  const add = () => {
    showDialog.value = true
    rules.password.required = true
    item.value = createDefaultItem()
    errors.value = {}
  }

  const edit = (row) => {
    showDialog.value = true
    item.value = { ...row }
    rules.password.required = false
    errors.value = {}
  }

  const validate = (form) => {
    return validateForm(form, rules, errors.value)
  }

  const handleResetPass = async () => {
    if (!validate(item.value)) {
      return
    }
    showLoading()
    try {
      await httpPost('/api/admin/user/resetPass', {
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
      await httpPost('/api/admin/user/save', item.value)
      showMessageOK('操作成功！')
      showDialog.value = false
      await fetchData(dataSets.value.page)
    } catch (error) {
      showMessageError('操作失败，' + error.message)
      showDialog.value = false
    }
  }

  const enable = async (row) => {
    try {
      await httpPost('/api/admin/user/enable', {
        id: row.id,
        enabled: row.enabled,
      })
      showMessageOK('操作成功！')
    } catch (error) {
      showMessageError('操作失败：' + error.message)
    }
  }

  const remove = (row) => {
    showConfirm('删除提示', '确定要删除当前记录吗?', async () => {
      showLoading()
      try {
        await httpGet('/api/admin/user/remove?id=' + row.id)
        showMessageOK('删除成功！')
        await fetchData(1)
      } catch (error) {
        showMessageError('删除失败：' + error.message)
      }
    })
  }

  const handleSelectionChange = (rows) => {
    userIds.value = rows.map((row) => row.id)
  }

  const multiRemove = () => {
    if (userIds.value.length === 0) {
      showMessageError('请选择你要删除的数据行')
      return
    }
    showConfirm(
      '删除提示',
      '此操作将会永久删除用户信息和聊天记录，确认操作吗?',
      async () => {
        showLoading()
        try {
          await httpGet('/api/admin/user/remove', { ids: userIds.value })
          showMessageOK('删除成功！')
          await fetchData(1)
        } catch (error) {
          showMessageError('删除失败：' + error.message)
        }
      }
    )
  }

  const initialize = async () => {
    await fetchData(1)
  }

  return {
    loading,
    item,
    dataSets,
    errors,
    showDialog,
    showResetPassDialog,
    query,
    handleSubmit,
    remove,
    enable,
    add,
    edit,
    fetchData,
    handleResetPass,
    multiRemove,
    handleSelectionChange,
    initialize,
  }
})
