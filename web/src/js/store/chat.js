import { fetchEventSource } from '@microsoft/fetch-event-source'
import Clipboard from 'clipboard'
import MarkdownIt from 'markdown-it'
import emoji from 'markdown-it-emoji'
import { defineStore } from 'pinia'
import { nextTick, onMounted, onUnmounted, ref, watch } from 'vue'
import { useRouter } from 'vue-router'
import { getSystemInfo, getUserToken } from '@/js/cache/session.js'
import { useSharedStore } from '@/js/cache/sharedata'
import {
  showMessageError,
  showMessageInfo,
  showMessageOK,
  showMessageWarning,
} from '@/js/utils/dialog.js'
import { httpGet, httpPost } from '@/js/utils/http.js'
import { getFileType, UUID } from '@/js/utils/libs.js'

export const useChatStore = defineStore('chat', () => {
  const logo = ref('/images/logo.png')
  const isMobile = ref(window.innerWidth < 768)
  const drawer = ref(false)
  const messages = ref([])
  const chatList = ref([])
  const prompt = ref('')
  const clipboard = ref(null)
  const isNewChat = ref(true)
  const chatId = ref()
  const appId = ref(0)
  const isFocus = ref(false)
  const isSending = ref(false)
  const abortController = ref(null)
  const defaultTitle = ref('')
  const title = ref('')
  const router = useRouter()
  const greet = ref('今天有什么可以帮忙的？')
  const systemConfig = ref({})
  // app list
  const showAppDialog = ref(false)
  const appList = ref([])
  const showLoginDialog = ref(false)
  const showNoticeDialog = ref(false)
  const noticeKey = ref('SYSTEM_NOTICE')
  const notice = ref('')
  const selectedApp = ref({
    icon: '/images/app-placeholder.png',
    name: '通用助手',
    summary: '今天有什么可以帮忙的？',
  })
  const showWelcome = ref(true)
  // 文件上传
  const showFileDialog = ref(false)
  const selectedFiles = ref([])
  const store = useSharedStore()

  const getSelectedApp = (appId) => {
    if (appId) {
      return appList.value.find((item) => item.id === Number(appId))
    }
    return appList.value[0]
  }

  onMounted(async () => {
    window.addEventListener('resize', () => {
      isMobile.value = window.innerWidth < 768
    })
    if (router.currentRoute.value.params.id) {
      chatId.value = router.currentRoute.value.params.id
      showWelcome.value = false
    } else {
      chatId.value = UUID()
    }

    // 获取系统配置
    systemConfig.value = await getSystemInfo()
    defaultTitle.value = `我是 ${systemConfig.value.title}，很高兴见到你！`
    title.value = defaultTitle.value

    try {
      await fetchAppList()
    } catch (e) {
      console.info('获取应用列表失败:' + e.message)
    }

    try {
      // 获取聊天记录
      const res = await httpGet('/api/chat/list')
      chatList.value = res.data
    } catch (e) {
      console.info('获取聊天记录失败:' + e.message)
    }

    // 获取网站公告
    const md = new MarkdownIt({
      breaks: true,
      html: true,
      linkify: true,
      typographer: true,
    }).use(emoji)
    let res = await httpGet('/api/config/get?name=notice')
    try {
      notice.value = md.render(res.data['content'])
      const oldNotice = localStorage.getItem(noticeKey.value)
      // 如果公告有更新，则显示公告
      if (oldNotice !== notice.value && notice.value.length > 10) {
        showNoticeDialog.value = true
      }
    } catch (e) {
      console.warn(e)
    }

    // 继续旧的对话
    const chat = chatList.value.find((item) => item.chat_id === chatId.value)
    if (chat) {
      title.value = chat.title
      appId.value = chat.app_id
      selectedApp.value = getSelectedApp(chat.app_id)
      if (!selectedApp.value) {
        selectedApp.value = appList.value[0]
      }
    } else {
      title.value = defaultTitle
      appId.value = Number(router.currentRoute.value.query.app_id)
      if (appId.value > 0) {
        selectedApp.value = getSelectedApp(appId.value)
        if (!selectedApp.value) {
          showMessageWarning('应用未审核通过或者已禁用')
        }
      }

      // 读取系统配置，获取默认 App
      if (!appId.value) {
        appId.value = systemConfig.value.app_id
      }

      // 如果系统配置没有设置默认 App，则使用第一个 App
      if (!appId.value) {
        appId.value = appList.value[0].id
      }
      selectedApp.value = getSelectedApp(appId.value)

      messages.value = setHelloMessage(selectedApp.value.summary)
    }

    // 加载聊天记录
    if (store.isLogin) {
      await loadMessages(chatId.value)
    }

    // 初始化复制组件 - 使用 try-catch 包装
    try {
      clipboard.value = new Clipboard('.copy-answer, .copy-code-btn')
      clipboard.value.on('success', () => {
        showMessageOK('复制成功！')
      })
      clipboard.value.on('error', () => {
        showMessageError('复制失败！')
      })
    } catch (e) {
      console.warn('Clipboard initialization failed:', e)
    }
  })

  onUnmounted(() => {
    if (clipboard.value) {
      clipboard.value.destroy()
    }
  })

  watch(
    () => store.showLoginDialog,
    (newValue) => {
      showLoginDialog.value = newValue
    }
  )

  const fetchAppList = async (cid = 0, keyword = '') => {
    const res = await httpGet('/api/app/list', {
      cid: cid,
      keyword: keyword,
    })
    appList.value = res.data || []
  }

  const notShowNotice = () => {
    localStorage.setItem(noticeKey.value, notice.value)
    showNoticeDialog.value = false
  }

  const setHelloMessage = (message) => {
    if (!message) {
      message = '你好，今天有什么可以帮忙的？'
    }
    return [
      {
        role: 'assistant',
        content: {
          texts: [
            `<div class="alert alert-primary d-flex align-items-center" role="alert">${message}</div>`,
          ],
          type: 'alert',
        },
        icon: selectedApp.value.icon,
      },
    ]
  }

  const loadMessages = async (chatId) => {
    httpGet('/api/chat/messages', { chat_id: chatId })
      .then((res) => {
        if (res.data) {
          messages.value = res.data
          messages.value.map((item) => {
            item.icon = item.icon || selectedApp.value.icon
            item.completed = true
            item.toolsCompleted = true
          })
          scrollToBottom()
        }
      })
      .catch((e) => {
        showMessageError('获取聊天信息失败：' + e.message)
      })
  }

  const messageStart = () => {
    if (messages.value.length > 1) {
      messages.value[messages.value.length - 1].start = true
    }
  }
  const messageEnd = () => {
    messages.value[messages.value.length - 1].completed = true
    isSending.value = false
  }

  const messageUpdate = (message) => {
    const lastMessageTexts =
      messages.value[messages.value.length - 1].content.texts
    if (lastMessageTexts.length === 0) {
      lastMessageTexts.push('')
    }
    lastMessageTexts[lastMessageTexts.length - 1] += message.content
    scrollToBottom()
  }

  const messageAlert = (message) => {
    messages.value[messages.value.length - 1] = {
      role: 'assistant',
      icon: selectedApp.value.icon,
      content: {
        texts: [
          `<div class="alert alert-danger d-flex align-items-center" role="alert">${message}</div>`,
        ],
        type: 'alert',
      },
    }
    scrollToBottom()
  }

  const toolUpdate = (tool) => {
    const tools = messages.value[messages.value.length - 1].content.tools
    if (tools.length === 0) {
      tools.push(tool)
    } else {
      const toolItem = tools.find((item) => item.name === tool.name)
      if (toolItem) {
        toolItem.status = tool.status
        toolItem.spend = tool.spend
        toolItem.response = tool.response
      } else {
        tools.push(tool)
      }
    }
    scrollToBottom()
  }

  const updateTitle = (value) => {
    if (title.value !== value) {
      httpGet('/api/chat/list')
        .then((res) => {
          chatList.value = res.data
        })
        .catch((e) => {
          showMessageError('获取聊天记录失败:' + e.message)
        })
    }
  }

  // 处理参数输入消息
  const handleInputMessage = (data) => {
    // 确保最后一条消息存在
    if (messages.value.length === 0) {
      messages.value.push({
        role: 'assistant',
        icon: selectedApp.value.icon,
        content: {
          texts: [],
          tools: [],
          files: [],
          type: 'text',
        },
        completed: false,
      })
    }
    const lastMessage = messages.value[messages.value.length - 1]
    // 解析数据（可能是字符串或已经是对象）
    let inputData = data
    if (typeof data === 'string') {
      try {
        inputData = JSON.parse(data)
      } catch (e) {
        console.error('Failed to parse input data:', e)
        return
      }
    }
    // 添加到消息内容中
    lastMessage.content.inputForm = inputData
    lastMessage.completed = true
    scrollToBottom()
  }

  // 处理问答消息
  const handleAnswerMessage = (data) => {
    // 确保最后一条消息存在
    if (messages.value.length === 0) {
      messages.value.push({
        role: 'assistant',
        icon: selectedApp.value.icon,
        content: {
          texts: [],
          tools: [],
          files: [],
          type: 'text',
        },
        completed: false,
      })
    }
    const lastMessage = messages.value[messages.value.length - 1]
    // 解析数据（可能是字符串或已经是对象）
    let answerData = data
    if (typeof data === 'string') {
      try {
        answerData = JSON.parse(data)
      } catch (e) {
        console.error('Failed to parse answer data:', e)
        return
      }
    }
    // 添加到消息内容中
    lastMessage.content.questionAnswer = answerData
    lastMessage.completed = true
    scrollToBottom()
  }

  const newChat = () => {
    isNewChat.value = true
    isSending.value = false
    messages.value = setHelloMessage(selectedApp.value.summary)
    chatId.value = UUID()
    title.value = defaultTitle
    router.push('/')
  }

  // 发送消息
  const sendMessage = (messageId = 0) => {
    if (
      isSending.value ||
      (prompt.value.trim() === '' && selectedFiles.value.length === 0)
    ) {
      showMessageInfo('请输入你的问题')
      return
    }

    if (!store.isLogin) {
      showLoginDialog.value = true
      return
    }

    // 离线状态直接阻止发送，避免触发连续错误
    if (typeof navigator !== 'undefined' && navigator.onLine === false) {
      showMessageError('网络已断开，请检查网络后重试')
      return
    }

    showWelcome.value = false
    router.push('/chat/' + chatId.value)
    isSending.value = true
    messages.value.push({
      role: 'user',
      content: {
        texts: [prompt.value],
        type: 'text',
        files: selectedFiles.value,
      },
      created_at: new Date().getTime() / 1000,
    })
    messages.value.push({
      role: 'assistant',
      chat_id: chatId.value,
      content: {
        texts: [],
        tools: [],
        type: 'text',
      },
      created_at: new Date().getTime() / 1000,
      icon: selectedApp.value.icon,
      completed: false,
    })
    scrollToBottom()

    abortController.value = new AbortController()
    let errorNotified = false

    fetchEventSource('/api/chat/message', {
      method: 'POST',
      headers: {
        Authorization: getUserToken(),
      },
      body: JSON.stringify({
        prompt: prompt.value,
        app_id: appId.value,
        chat_id: chatId.value,
        files: selectedFiles.value.map((file) => ({
          type: getFileType(file.name),
          url: file.url,
          name: file.name,
          size: file.size,
        })),
        last_msg_id: messageId,
      }),
      openWhenHidden: true,
      signal: abortController.value.signal,
      // 设置失败重试间隔，避免离线时频繁快速重试
      retry: 3000,
      async onopen(res) {
        const contentType = res.headers.get('content-type') || ''
        if (res.ok && contentType.includes('text/event-stream')) {
          return
        }
        if (res.status === 429) {
          let msg = '请求过于频繁，请稍后重试'
          try {
            const data = await res.clone().json()
            msg = data.message || msg
          } catch (_) {
            try {
              const text = await res.clone().text()
              if (text) msg = text
            } catch (_) {}
          }
          isSending.value = false
          messageAlert(msg)
        }
      },
      onmessage(msg) {
        switch (msg.event) {
          case 'error':
            messageAlert(msg.data)
            isSending.value = false
            break
          case 'start':
            messageStart()
            break
          case 'tool':
            toolUpdate(JSON.parse(msg.data))
            break
          case 'end':
            messageEnd()
            break
          case 'title':
            updateTitle(msg.data)
            break
          case 'message.delta':
            messageUpdate(JSON.parse(msg.data))
            break
          case 'message.completed':
            const lastMessage = messages.value[messages.value.length - 1]
            lastMessage.content.texts.push('')
            break
          case 'input':
            handleInputMessage(msg.data)
            break
          case 'answer':
            handleAnswerMessage(msg.data)
            break
          case 'follow_up':
            console.log('follow_up', msg.data)
            break
        }
      },
      onerror(err) {
        console.log('onerror', err)
        // 用户主动中断不提示错误
        if (
          err &&
          (err.name === 'AbortError' ||
            err.message === 'The user aborted a request.')
        ) {
          return
        }
        if (errorNotified) return
        errorNotified = true
        isSending.value = false
        try {
          abortController.value && abortController.value.abort()
        } catch (_) {}
        const msg =
          typeof navigator !== 'undefined' && navigator.onLine === false
            ? '网络已断开'
            : err.message
        showMessageError('发送消息失败:' + msg)
        // 抛出错误以停止后续重试
        throw err
      },
      onclose() {},
    })
    prompt.value = ''
    selectedFiles.value = []
  }

  const regenerateMessage = (messageId) => {
    messages.value = messages.value.filter(
      (item) => !item.id || item.id < messageId
    )
    // 删除用户消息
    const lastMessage = messages.value.pop()
    // 填入输入框
    prompt.value = lastMessage.content.texts[0]
    sendMessage(lastMessage.id)
  }
  const sendCustomMessage = (message) => {
    prompt.value = message
    sendMessage(0)
  }

  // 中断对话
  const abortRequest = () => {
    if (abortController.value) {
      abortController.value.abort()
      isSending.value = false
      httpPost('/api/chat/cancel', {
        chat_id: chatId.value,
        app_id: appId.value,
      })
        .then(() => {
          console.log('会话已中断')
        })
        .catch((e) => {
          showMessageError('中断对话失败:' + e.message)
        })
    }
  }

  // 切换对话
  const changeChat = (chat) => {
    router.push('/chat/' + chat.chat_id)
    chatId.value = chat.chat_id
    title.value = chat.title
    loadMessages(chat.chat_id)
    drawer.value = false
    appId.value = chat.app_id
    selectedApp.value = getSelectedApp(appId.value)
    if (!selectedApp.value) {
      selectedApp.value = appList.value[0]
    }
    showWelcome.value = false
  }

  // 删除聊天记录
  const deleteChat = (chat) => {
    if (chat.chat_id === chatId.value) {
      newChat()
    }
  }

  // 清空聊天记录
  const clearChats = () => {
    chatList.value = []
    newChat()
    showWelcome.value = true
  }

  // 滚动到页面底部
  const scrollToBottom = () => {
    nextTick(() => {
      const mainContainer = document.getElementById('main-container')
      mainContainer.scrollTop = mainContainer.scrollHeight
    })
  }

  // 选择应用
  const selectApp = (app) => {
    selectedApp.value = app
    appId.value = app.id
    showAppDialog.value = false
    showWelcome.value = false
    newChat()
  }

  // 选中文件回调
  const selectFile = (file) => {
    if (selectedFiles.value.length >= 3) {
      showMessageWarning('最多只能上传 3 个文件')
      return
    }
    if (!selectedFiles.value.find((item) => item.url === file.url)) {
      selectedFiles.value.push(file)
      showFileDialog.value = false
    } else {
      showMessageWarning('文件已选中')
    }
  }

  // 移除文件
  const removeFile = (file) => {
    console.log('removeFile', file)
    selectedFiles.value = selectedFiles.value.filter(
      (item) => item.url !== file.url
    )
  }

  return {
    prompt,
    chatId,
    isFocus,
    logo,
    isMobile,
    isNewChat,
    title,
    drawer,
    messages,
    chatList,
    isSending,
    greet,
    systemConfig,
    showAppDialog,
    appList,
    showLoginDialog,
    showNoticeDialog,
    showFileDialog,
    notice,
    selectedApp,
    selectedFiles,
    showWelcome,
    newChat,
    sendMessage,
    abortRequest,
    changeChat,
    deleteChat,
    clearChats,
    selectApp,
    selectFile,
    removeFile,
    notShowNotice,
    fetchAppList,
    regenerateMessage,
    sendCustomMessage,
  }
})
