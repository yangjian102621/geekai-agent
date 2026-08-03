<template>
  <div class="chat-list d-flex flex-column flex-shrink-0">
    <div
      class="logo d-flex justify-content-md-center justify-content-between align-items-center text-dark"
    >
      <el-avatar :src="systemConfig.logo" class="bg-white me-2" />
      <span class="fs-4">{{ systemConfig.title }}</span>
      <span v-if="isMobile()" @click="emits('close')"
        ><i class="iconfont icon-colspan"></i
      ></span>
    </div>
    <button class="btn btn-primary mt-3 mb-3" @click="newChat">
      <i class="iconfont icon-new-chat me-1"></i>开启新对话
    </button>
    <el-scrollbar :max-height="maxHeight" :height="maxHeight">
      <div v-if="items.length > 0">
        <el-menu class="chat-item-menu" id="chat-list-menu">
          <el-menu-item
            :index="'index' + chat.id"
            v-for="chat in items"
            :key="chat.id"
            @click="changeChat(chat)"
            :class="{ 'is-active': chat.chat_id === chatId }"
          >
            <div class="d-flex w-100" v-if="chat.edit">
              <input
                class="form-control form-control-sm"
                v-model="chat.title"
                @blur="updateTitle(chat)"
                @keydown="handleKeyDown($event, chat)"
                ref="chatRef"
              />
            </div>
            <div class="chat-item" v-else>
              <el-image
                :src="chat.icon"
                class="w-[20px] h-[20px] rounded-full me-1"
                fit="cover"
                v-if="chat.icon"
              />
              <i v-else class="iconfont icon-chat me-1"></i>
              <div class="chat-title text-nowrap overflow-hidden">
                {{ chat.title }}
              </div>
              <span class="chat-opt">
                <el-dropdown trigger="click" popper-class="dropdown-round">
                  <span
                    class="el-dropdown-link me-1"
                    @click="stopPropagation($event)"
                  >
                    <i class="iconfont icon-more-horizontal"></i>
                  </span>
                  <template #dropdown>
                    <el-dropdown-menu>
                      <el-dropdown-item @click="editChat(chat)">
                        <i class="iconfont icon-edit me-1"></i
                        >重命名</el-dropdown-item
                      >
                      <el-dropdown-item @click="removeChat(chat)">
                        <span class="text-danger"
                          ><i class="iconfont icon-remove me-1"></i>删除</span
                        >
                      </el-dropdown-item>
                      <el-dropdown-item
                        class="copy-share-link"
                        @click="copyShareUrl(chat)"
                      >
                        <i class="iconfont icon-share1 me-1"></i
                        >分享</el-dropdown-item
                      >
                    </el-dropdown-menu>
                  </template>
                </el-dropdown>
              </span>
            </div>
          </el-menu-item>
        </el-menu>
      </div>
      <el-empty v-else description="暂无对话" :image-size="100" />
    </el-scrollbar>

    <el-button class="mt-3 mb-3" type="danger" plain @click="clearChatList">
      <i class="iconfont icon-remove me-1"></i>清空对话列表
    </el-button>
  </div>
</template>

<script setup>
  import { getLicense, getSystemInfo } from '@/js/cache/session.js'
  import {
    showConfirm,
    showMessageError,
    showMessageOK,
  } from '@/js/utils/dialog.js'
  import { httpGet, httpPost } from '@/js/utils/http.js'
  import { isMobile } from '@/js/utils/libs'
  import { removeArrayItem } from '@/js/utils/libs.js'
  import { nextTick, onMounted, ref, watch } from 'vue'
  import { useRouter } from 'vue-router'

  const props = defineProps({
    modelValue: {
      type: Array,
      default: () => [],
    },
    chatId: {
      type: String,
      default: '',
    },
  })

  const emits = defineEmits([
    'delete',
    'close',
    'clear',
    'new-chat',
    'change-chat',
  ])

  const chatRef = ref(null)
  const router = useRouter()
  const items = ref(props.modelValue)
  const systemConfig = ref({})
  const license = ref({})

  onMounted(async () => {
    systemConfig.value = await getSystemInfo()
    license.value = await getLicense()
    if (!license.value.is_active) {
      systemConfig.value.logo = '/images/logo.png'
      systemConfig.value.title = 'GeekAI-智能体'
    }
  })

  // 新增：原生复制方法
  const copyShareUrl = async (chat) => {
    try {
      await navigator.clipboard.writeText(
        `${window.location.protocol}//${window.location.host}/share/${chat.chat_id}`
      )
      showMessageOK('分享链接已复制成功，赶快分享给你的朋友吧！')
    } catch (e) {
      showMessageError('复制分享链接失败')
    }
  }

  watch(
    () => props.modelValue,
    (newVal) => {
      items.value = newVal
    }
  )

  const stopPropagation = (e) => {
    e.stopPropagation()
  }
  const maxHeight = ref(window.innerHeight - 190)
  const editChat = (chat) => {
    chat.edit = true
    nextTick(() => {
      chatRef.value[0].focus()
      chat.titleBak = chat.title
    })
  }

  const handleKeyDown = (e, chat) => {
    if (e.key === 'Enter' && e.keyCode === 13) {
      e.preventDefault()
      updateTitle(chat)
    }
  }

  // 更新标题
  const updateTitle = (chat) => {
    chat.edit = false
    if (chat.title === chat.titleBak) {
      return
    }
    httpPost('/api/chat/update', {
      chat_id: chat.chat_id,
      title: chat.title,
    }).catch((e) => {
      showMessageError('更新对话标题失败：' + e.message)
    })
  }

  // 删除对话
  const removeChat = (chat) => {
    emits('close')
    showConfirm('删除对话', '删除后，该对话将不可恢复。确认删除吗？', () => {
      httpGet('/api/chat/remove', { chat_id: chat.chat_id })
        .then(() => {
          showMessageOK('删除对话成功')
          items.value = removeArrayItem(
            items.value,
            chat,
            (a, b) => a.chat_id === b.chat_id
          )
          emits('delete', chat)
        })
        .catch((e) => {
          showMessageError('删除对话失败：' + e.message)
        })
    })
  }

  const changeChat = (chat) => {
    emits('change-chat', chat)
  }

  const newChat = () => {
    emits('new-chat')
    // 移除所有对话的选中状态
    nextTick(() => {
      const activeChats = document.querySelectorAll(
        '#chat-list-menu .is-active'
      )
      activeChats.forEach((chat) => chat.classList.remove('is-active'))
    })
    // 隐藏弹框
    emits('close')
  }

  const clearChatList = () => {
    showConfirm(
      '清空对话列表',
      '清空后，所有对话将不可恢复。确认清空吗？',
      () => {
        httpGet('/api/chat/clear').then(() => {
          showMessageOK('清空对话列表成功')
          items.value = []
          emits('clear')
        })
      }
    )
  }
</script>

<style scoped lang="scss">
  .chat-list {
    padding: 1rem 1rem 0 1rem;

    .logo {
      .iconfont {
        font-size: 20px;
      }
    }

    .chat-item-menu {
      --el-menu-item-color: #333;
      --el-menu-bg-color: transparent;
      --el-menu-hover-bg-color: rgba(194, 164, 245, 0.3);
      --el-menu-hover-text-color: #7c39ed;
      --el-menu-border-color: 0;
      --el-menu-base-level-padding: 10px;

      .chat-item {
        display: flex;
        align-items: center;

        .chat-title {
          max-width: 202px;
          min-width: 202px;
        }

        .chat-opt {
          .iconfont {
            margin-left: 0.5rem;
            padding: 1px;
            border-radius: 5px;
            background-color: rgba(194, 164, 245, 0.3);
          }
        }
      }
    }

    .el-menu-item {
      color: #333;
      margin: 0.2rem 0;
      --el-menu-item-height: 40px;
      border-radius: 10px;
    }

    .el-menu-item.is-active {
      background-color: var(--el-menu-hover-bg-color);
      color: var(--el-menu-hover-text-color);
    }
  }
</style>
