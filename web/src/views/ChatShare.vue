<template>
  <div class="chat-share-container">
    <!-- 标题栏 -->
    <div class="header">
      <h3 class="title">{{ chatTitle }}</h3>
    </div>

    <!-- 消息列表 -->
    <div class="message-list" v-if="messages.length > 0">
      <div v-for="message in messages" :key="message.id" class="message-item">
        <message-user :data="message" v-if="message.role === 'user'" />
        <message-assistant :data="message" v-else />
      </div>
    </div>

    <!-- 加载状态 -->
    <div v-else-if="loading" class="loading-container">
      <div class="text-center">
        <div class="spinner-border" role="status">
          <span class="visually-hidden">Loading...</span>
        </div>
        <p class="mt-2">正在加载对话...</p>
      </div>
    </div>

    <!-- 错误状态 -->
    <div v-else class="error-container">
      <div class="text-center">
        <h4>对话不存在或已被删除</h4>
        <p class="text-muted">请检查分享链接是否正确</p>
      </div>
    </div>

    <!-- 底部按钮 -->
    <div class="footer">
      <button class="btn btn-primary btn-lg try-button" @click="goToHome">
        我也试试
      </button>
    </div>
  </div>
</template>

<script setup>
  import MessageAssistant from '@/components/MessageAssistant.vue'
  import MessageUser from '@/components/MessageUser.vue'
  import { httpGet } from '@/js/utils/http'
  import { onMounted, ref } from 'vue'
  import { useRoute, useRouter } from 'vue-router'

  const route = useRoute()
  const router = useRouter()

  // 响应式数据
  const chatTitle = ref('')
  const messages = ref([])
  const loading = ref(true)

  // 获取对话详情
  const getChatDetail = async (chatId) => {
    try {
      const response = await httpGet(`/api/chat/detail?chat_id=${chatId}`)
      if (response.data) {
        chatTitle.value = response.data.title || '未命名对话'
        document.title = `对话分享 - ${chatTitle.value}` // 设置页面标题
      }
    } catch (error) {
      console.error('获取对话详情失败:', error)
    }
  }

  // 获取对话消息
  const getChatMessages = async (chatId) => {
    try {
      const response = await httpGet(`/api/chat/messages?chat_id=${chatId}`)
      if (response.data) {
        messages.value = response.data
        messages.value.map((item) => {
          item.icon = item.icon || selectedApp.value.icon
          item.completed = true
          item.toolsCompleted = true
          item.share = true
        })
      }
    } catch (error) {
      console.error('获取对话消息失败:', error)
    }
  }

  // 跳转到首页
  const goToHome = () => {
    router.push('/')
  }

  // 初始化页面
  onMounted(async () => {
    const chatId = route.params.chat_id
    if (chatId) {
      try {
        await Promise.all([getChatDetail(chatId), getChatMessages(chatId)])
      } finally {
        loading.value = false
      }
    } else {
      loading.value = false
    }
  })
</script>

<style scoped>
  .chat-share-container {
    min-height: 100vh;
    display: flex;
    flex-direction: column;
    background-color: #f8f9fa;
  }

  .header {
    background: white;
    padding: 1rem 2rem;
    border-bottom: 1px solid #e9ecef;
    box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
  }

  .title {
    margin: 0;
    text-align: center;
    color: #333;
    font-weight: 600;
  }

  .message-list {
    flex: 1;
    padding: 2rem;
    max-width: 1000px;
    margin: 0 auto;
    width: 100%;
  }

  .message-item {
    margin-bottom: 1rem;
  }

  .loading-container,
  .error-container {
    flex: 1;
    display: flex;
    align-items: center;
    justify-content: center;
    padding: 2rem;
  }

  .footer {
    background: white;
    padding: 1.5rem 2rem;
    border-top: 1px solid #e9ecef;
    text-align: center;
    box-shadow: 0 -2px 4px rgba(0, 0, 0, 0.1);
  }

  .try-button {
    padding: 0.75rem 2rem;
    font-size: 1.1rem;
    font-weight: 600;
    border-radius: 25px;
    min-width: 150px;
  }

  /* 响应式设计 */
  @media (max-width: 768px) {
    .header {
      padding: 1rem;
    }

    .message-list {
      padding: 1rem;
    }

    .footer {
      padding: 1rem;
    }

    .try-button {
      width: 100%;
      max-width: 300px;
    }
  }
</style>
