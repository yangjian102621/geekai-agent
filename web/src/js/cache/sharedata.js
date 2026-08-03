import Storage from 'good-storage'
import { defineStore } from 'pinia'

export const useSharedStore = defineStore('shared', {
  state: () => ({
    collapsed: false,
    chatSetting: Storage.get('chatSetting', {
      darkTheme: false,
      chatContext: true,
    }),
    showLoginDialog: false,
    isLogin: true,
    appId: '',
    creatorId: Storage.get('creatorId', ''),
  }),
  getters: {},
  actions: {
    setCollapsed(value) {
      this.collapsed = value
      Storage.set('collapsed', value)
    },
    setChatSetting(value) {
      this.chatSetting = value
      Storage.set('chatSetting', value)
    },
    setShowLoginDialog(value) {
      this.showLoginDialog = value
    },
    setIsLogin(value) {
      this.isLogin = value
    },
    setAppId(value) {
      this.appId = value
    },
    setCreatorId(value) {
      this.creatorId = value
      Storage.set('creatorId', value)
    },
  },
})
