import { defineStore } from 'pinia'
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { getSystemInfo, setAdminToken } from '@/js/cache/session.js'
import { showMessageError, showMessageOK } from '@/js/utils/dialog.js'
import { httpGet, httpPost } from '@/js/utils/http.js'
import { replaceURL } from '@/js/utils/libs.js'
import { validateForm } from '@/js/store/common.js'

const createDefaultForm = () => ({
  username: import.meta.env.VITE_ADMIN_USER,
  password: import.meta.env.VITE_ADMIN_PASS,
})

const CAPTCHA_RULES = {
  username: { required: true, message: '请输入用户名' },
  password: { required: true, message: '请输入密码' },
}

export const useAdminLoginStore = defineStore('admin-login', () => {
  const router = useRouter()
  const indexURL = '/admin/dashboard'

  const data = ref(createDefaultForm())
  const errors = ref({
    username: '',
    password: '',
  })
  const captchaRef = ref(null)
  const loading = ref(false)
  const systemConfig = ref({})
  const captchaConfig = ref({ enabled: false })

  const loadConfigs = async () => {
    try {
      systemConfig.value = await getSystemInfo()
      if (!systemConfig.value.logo) {
        systemConfig.value.logo = replaceURL('/images/logo.png')
      }
    } catch (error) {
      console.error(error)
    }
    try {
      const res = await httpGet('/api/config/captcha')
      captchaConfig.value = res.data || { enabled: false }
    } catch (error) {
      console.error(error)
    }
  }

  const doLogin = async (captcha = {}) => {
    loading.value = true
    try {
      const res = await httpPost('/api/admin/login', {
        username: data.value.username,
        password: data.value.password,
        key: captcha.key,
        dots: captcha.dots,
        x: captcha.x,
      })
      setAdminToken(res.data.token)
      showMessageOK('登录成功!')
      router.push(indexURL)
    } catch (error) {
      showMessageError('登录失败：' + error.message)
      throw error
    } finally {
      loading.value = false
    }
  }

  const handleSubmit = async () => {
    if (!validateForm(data.value, CAPTCHA_RULES, errors.value)) {
      return
    }
    if (captchaConfig.value.enabled && captchaRef.value) {
      captchaRef.value.loadCaptcha()
      return
    }
    await doLogin({})
  }

  return {
    loading,
    data,
    errors,
    captchaRef,
    systemConfig,
    captchaConfig,
    loadConfigs,
    doLogin,
    handleSubmit,
  }
})
