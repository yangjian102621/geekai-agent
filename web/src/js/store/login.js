import { defineStore } from 'pinia'
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { getSystemInfo, setUserToken } from '../cache/session.js'
import { showMessageError, showMessageOK } from '../utils/dialog.js'
import { httpGet, httpPost } from '../utils/http.js'
import { replaceURL } from '../utils/libs.js'
import { validateForm } from './common.js'

const createDefaultForm = () => ({
  username: import.meta.env.VITE_USER,
  password: import.meta.env.VITE_PASS,
})

const RULES = {
  username: { required: true, message: '请输入用户名' },
  password: { required: true, message: '请输入密码' },
}

export const useLoginStore = defineStore('user-login', () => {
  const router = useRouter()
  const data = ref(createDefaultForm())
  const errors = ref({
    username: '',
    password: '',
  })
  const captchaRef = ref(null)
  const indexURL = ref('/')
  const loading = ref(false)
  const systemConfig = ref({})
  const captchaConfig = ref({
    enabled: false,
    type: 'slide',
  })

  const loadConfigs = async () => {
    try {
      systemConfig.value = await getSystemInfo()
      if (!systemConfig.value.logo) {
        systemConfig.value.logo = replaceURL('/images/logo.png')
      }
    } catch (error) {
      console.warn(error)
    }
    try {
      const res = await httpGet('/api/config/captcha')
      captchaConfig.value = res.data || { enabled: false, type: 'slide' }
    } catch (error) {
      console.warn(error)
    }
  }

  const doLogin = async (captcha = {}) => {
    loading.value = true
    try {
      const res = await httpPost('/api/user/login', {
        username: data.value.username,
        password: data.value.password,
        key: captcha.key,
        dots: captcha.dots,
        x: captcha.x,
      })
      setUserToken(res.data.token)
      showMessageOK('登录成功!')
      router.push(indexURL.value)
    } catch (error) {
      showMessageError('登录失败：' + error.message)
      throw error
    } finally {
      loading.value = false
    }
  }

  const handleSubmit = async () => {
    if (!validateForm(data.value, RULES, errors.value)) {
      return
    }
    if (systemConfig.value.enabled_verify && captchaRef.value) {
      captchaRef.value.loadCaptcha()
      return
    }
    await doLogin({})
  }

  const initialize = async () => {
    await loadConfigs()
  }

  return {
    loading,
    data,
    errors,
    indexURL,
    captchaRef,
    systemConfig,
    captchaConfig,
    doLogin,
    handleSubmit,
    loadConfigs,
    initialize,
  }
})
