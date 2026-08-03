import { defineStore } from 'pinia'
import { ref } from 'vue'
import { checkSession } from '@/js/cache/session'

export const useUserStore = defineStore('user', () => {
  const userInfo = ref({
    username: '',
    power: 0,
    email: '',
    mobile: '',
    avatar: '',
    id: 0,
  })
  const isLogin = ref(false)

  const fetchUserInfo = async () => {
    try {
      const data = await checkSession()
      if (data) {
        userInfo.value = data
        isLogin.value = true
      } else {
        isLogin.value = false
      }
    } catch (e) {
      console.error(e)
      isLogin.value = false
    }
  }

  return {
    userInfo,
    isLogin,
    fetchUserInfo,
  }
})
