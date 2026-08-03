import { defineStore } from 'pinia'
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { getSystemInfo, setUserToken } from '../cache/session.js'
import { showMessageError, showMessageOK } from '../utils/dialog.js'
import { httpPost } from '../utils/http.js'
import { replaceURL } from '../utils/libs.js'
import { validateForm } from './common.js'

const createDefaultForm = () => ({
  username: '',
  password: '',
  repass: '',
  code: '',
})

const RULES = {
  username: { required: true, message: '请输入用户名' },
  password: { required: true, message: '请输入密码' },
  repass: { required: true, message: '请输入重复密码' },
  code: { required: true, message: '请输入验证码' },
}

export const useRegisterStore = defineStore('user-register', () => {
  const router = useRouter()
  const data = ref(createDefaultForm())
  const errors = ref({
    username: '',
    password: '',
    repass: '',
    code: '',
  })
  const captchaRef = ref(null)
  const indexURL = ref('/')
  const loading = ref(false)
  const systemConfig = ref({})

  const loadConfig = async () => {
    try {
      systemConfig.value = await getSystemInfo()
      if (!systemConfig.value.logo) {
        systemConfig.value.logo = replaceURL('/images/logo.png')
      }
    } catch (error) {
      console.warn(error)
    }
  }

  const doRegister = async () => {
    loading.value = true
    try {
      const res = await httpPost('/api/user/register', {
        username: data.value.username,
        password: data.value.password,
        code: data.value.code,
      })
      setUserToken(res.data.token)
      showMessageOK('注册成功!')
      router.push(indexURL.value)
    } catch (error) {
      showMessageError('注册失败：' + error.message)
      throw error
    } finally {
      loading.value = false
    }
  }

  const handleSubmit = async () => {
    if (!validateForm(data.value, RULES, errors.value)) {
      showMessageError('请输入正确的信息')
      return
    }
    if (data.value.password !== data.value.repass) {
      errors.value.repass = '两次输入的密码不一致'
      return
    }
    await doRegister()
  }

  const initialize = async () => {
    await loadConfig()
  }

  return {
    loading,
    data,
    errors,
    indexURL,
    captchaRef,
    systemConfig,
    doRegister,
    handleSubmit,
    loadConfig,
    initialize,
  }
})
