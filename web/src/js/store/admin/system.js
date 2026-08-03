import { defineStore } from 'pinia'
import { ref } from 'vue'
import {
  showLoading,
  showMessageError,
  showMessageOK,
} from '../../utils/dialog.js'
import { httpGet, httpPost } from '../../utils/http.js'
import { copyObj } from '../../utils/libs.js'
import { adminUploadFile } from '../common.js'

export const useAdminSystemStore = defineStore('admin-system', () => {
  const config = ref({})
  const configBackup = ref({})
  const appList = ref([])

  const loadConfig = async () => {
    try {
      const res = await httpGet('/api/admin/config/get?name=system')
      config.value = res.data || {}
      config.value.email_white_list = config.value.email_white_list || []
      configBackup.value = copyObj(config.value)
    } catch (error) {
      showMessageError('加载系统配置失败：' + error.message)
    }
  }

  const loadAppList = async () => {
    try {
      const res = await httpGet('/api/admin/app/list')
      appList.value = res.data?.items || []
    } catch (error) {
      showMessageError('获取应用列表失败：' + error.message)
    }
  }

  const saveConfig = async () => {
    showLoading()
    try {
      await httpPost('/api/admin/config/update/base', config.value)
      showMessageOK('操作成功！')
      configBackup.value = copyObj(config.value)
    } catch (error) {
      showMessageError('操作失败：' + error.message)
    }
  }

  const uploadLogo = (file) => {
    adminUploadFile(file, (data) => {
      config.value.logo = data.url
    })
  }

  const uploadBotAvatar = (file) => {
    adminUploadFile(file, (data) => {
      config.value.bot_avatar = data.url
    })
  }

  const uploadUserAvatar = (file) => {
    adminUploadFile(file, (data) => {
      config.value.user_avatar = data.url
    })
  }

  const uploadWechatCard = (file) => {
    adminUploadFile(file, (data) => {
      config.value.wechat_card_url = data.url
    })
  }

  const initialize = async () => {
    await Promise.all([loadConfig(), loadAppList()])
  }

  return {
    config,
    saveConfig,
    uploadLogo,
    uploadBotAvatar,
    uploadWechatCard,
    uploadUserAvatar,
    configBackup,
    appList,
    initialize,
    loadConfig,
    loadAppList,
  }
})
