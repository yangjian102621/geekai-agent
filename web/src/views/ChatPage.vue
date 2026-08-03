<template>
  <el-container class="chat-page-container d-flex">
    <!-- 左侧菜单栏 -->
    <el-aside v-show="!chatStore.isMobile" class="chat-sidebar">
      <chat-list
        v-model="chatStore.chatList"
        @new-chat="showWelcomePage"
        @change-chat="chatStore.changeChat"
        @delete="chatStore.deleteChat"
        @clear="chatStore.clearChats"
        :chat-id="chatStore.chatId"
        ref="chatListRef"
      />
    </el-aside>
    <!-- 右侧主内容区 -->
    <el-container direction="vertical" class="h-100">
      <!-- 标题栏 -->
      <el-header>
        <div class="header d-flex align-items-center pt-2 pb-2 gap-3">
          <span v-if="chatStore.isMobile" @click="chatStore.drawer = true"
            ><i class="iconfont icon-expand"></i
          ></span>
          <div class="header-center flex-grow-1">
            <div class="text-truncate fw-bold text-center">
              {{ chatStore.title }}
            </div>
          </div>

          <div
            class="new-chat-btn"
            v-if="chatStore.isMobile"
            @click="chatStore.newChat"
          >
            <i class="iconfont icon-new-chat fs-3 text-primary"></i>
          </div>
        </div>
      </el-header>
      <!-- 内容区域 -->
      <el-main class="main-container" id="main-container">
        <div class="chat-container">
          <div
            class="message-box pb-5"
            v-if="chatStore.messages.length > 0 && !chatStore.showWelcome"
          >
            <div
              v-for="item in chatStore.messages"
              :key="item.id"
              class="flex justify-content-center"
            >
              <message-user :data="item" v-if="item.role === 'user'" />
              <message-assistant
                :data="item"
                v-else
                @regenerate="chatStore.regenerateMessage"
                @send-message="chatStore.sendCustomMessage"
              />
            </div>

            <div>
              <back-top target=".el-main" :bottom="220" />
              <back-bottom target=".el-main" :bottom="160" />
            </div>
          </div>

          <div
            class="message-box pb-5 d-flex justify-content-center align-items-center flex-row"
            v-else
          >
            <div
              class="modal-app-list row p-3 w-100 d-flex justify-content-center"
              v-if="hotApps?.length > 0"
            >
              <div
                class="col-12 col-md-6 col-lg-4 mb-4"
                style="--bs-gutter-x: 1rem"
                v-for="item in hotApps"
                :key="item.id"
                @click="chatStore.selectApp(item)"
              >
                <app-list-item
                  :item="item"
                  :selected-app="chatStore.selectedApp"
                />
              </div>
            </div>

            <div
              class="flex fs-3 mb-3 flex-column gap-4 align-items-center"
              v-else
            >
              <div class="mb-2">
                <el-image
                  class="mb-4 mb-md-0"
                  :src="chatStore.systemConfig.logo"
                  style="
                    width: 100px;
                    height: 100px;
                    background-color: var(--bs-primary-light);
                    border-radius: 50%;
                  "
                />
              </div>
              <h2>Hi，很高兴见到你！</h2>

              <h4 class="text-center text-secondary">
                {{ chatStore.greet }}
              </h4>
            </div>
          </div>

          <div class="chat-footer-box d-flex justify-content-center">
            <div class="chat-footer w-100">
              <div class="input-box" :class="chatStore.isFocus ? 'focus' : ''">
                <div
                  class="d-flex w-100 pb-2"
                  v-if="chatStore.selectedFiles.length > 0"
                >
                  <file-list
                    :files="chatStore.selectedFiles"
                    @remove-file="chatStore.removeFile"
                  />
                </div>
                <el-input
                  class="input-custom"
                  v-model="chatStore.prompt"
                  placeholder="给 AI 发送消息，按 Enter 发送，Shift + Enter 换行"
                  type="textarea"
                  :autosize="{ minRows: 1, maxRows: 10 }"
                  @focus="chatStore.isFocus = true"
                  @blur="chatStore.isFocus = false"
                  @keydown="handleKeyDown"
                  @keyup.enter.exact="chatStore.sendMessage(0)"
                />
                <div
                  class="d-flex justify-content-between w-100 align-items-center"
                >
                  <div class="d-flex tools me-2 mt-2">
                    <div class="tool-item me-2">
                      <el-tooltip content="选择应用" placement="top">
                        <el-button
                          round
                          style="padding: 8px"
                          @click="chatStore.showAppDialog = true"
                        >
                          <el-image
                            class="me-1 rounded-circle"
                            :src="chatStore.selectedApp.icon"
                            style="width: 20px; height: 20px"
                          />

                          <span
                            class="text-truncate"
                            style="max-width: 100px"
                            >{{ chatStore.selectedApp.name }}</span
                          >
                        </el-button>
                      </el-tooltip>
                    </div>
                    <div class="tool-item">
                      <el-tooltip content="上传文件" placement="top">
                        <i
                          class="iconfont icon-attachment-cl"
                          @click="chatStore.showFileDialog = true"
                        ></i>
                      </el-tooltip>
                    </div>
                  </div>
                  <div
                    class="send ms-2"
                    :class="chatStore.prompt.length === 0 ? 'btn-disabled' : ''"
                  >
                    <button
                      class="btn btn-primary btn-sm"
                      @click="chatStore.sendMessage(0)"
                      v-if="!chatStore.isSending"
                      :disabled="chatStore.prompt.length === 0"
                    >
                      <i class="iconfont icon-send"></i>
                    </button>
                    <button
                      class="btn btn-primary btn-sm"
                      @click="chatStore.abortRequest"
                      v-else
                    >
                      <i class="iconfont icon-pause"></i>
                    </button>
                  </div>
                </div>
              </div>
              <div class="text-center text-body-tertiary tip p-1">
                内容由 AI 生成，请仔细甄别
                <Copyright />
                <span v-if="version"> | 当前版本 {{ version }}</span>
              </div>
            </div>
          </div>
        </div>
      </el-main>
    </el-container>
    <!-- 移动端菜单抽屉 -->
    <el-drawer
      :with-header="false"
      v-model="chatStore.drawer"
      size="300px"
      direction="ltr"
      style="--el-drawer-padding-primary: 0; border-radius: 0 20px 20px 0"
    >
      <chat-list
        v-model="chatStore.chatList"
        @new-chat="chatStore.newChat"
        @change-chat="chatStore.changeChat"
        @delete="chatStore.deleteChat"
        @clear="chatStore.clearChats"
        :chat-id="chatStore.chatId"
        @close="chatStore.drawer = false"
      />
    </el-drawer>

    <model-dialog
      v-model="chatStore.showAppDialog"
      title="请选择应用"
      :hide-footer="true"
      @cancel="chatStore.showAppDialog = false"
      :width="960"
    >
      <!-- 应用分类 -->
      <div class="d-flex px-3" v-if="appCategories.length > 0">
        <el-radio-group v-model="selectedCategory" @change="onCategoryChange">
          <el-radio-button label="全部" :value="0" />
          <el-radio-button
            v-for="category in appCategories"
            :key="category.id"
            :label="category.name"
            :value="category.id"
          />
        </el-radio-group>

        <!--搜索-->
        <div class="ms-auto">
          <el-input
            v-model="searchKeyword"
            placeholder="搜索应用"
            class="w-auto"
            @input="onSearch"
          >
            <template #prefix>
              <i class="iconfont icon-search"></i>
            </template>
          </el-input>
        </div>
      </div>
      <!-- 应用列表 -->
      <div class="modal-app-list row p-3" v-if="chatStore.appList.length > 0">
        <div
          class="col col-12 col-md-4 mb-4"
          style="--bs-gutter-x: 1rem"
          v-for="item in chatStore.appList"
          :key="item.id"
          @click="chatStore.selectApp(item)"
        >
          <app-list-item :item="item" :selected-app="chatStore.selectedApp" />
        </div>
      </div>
      <div class="text-center text-body-tertiary tip p-1" v-else>
        <el-empty description="暂无应用" />
      </div>
    </model-dialog>

    <model-dialog
      v-model="chatStore.showFileDialog"
      title="请选择文件"
      :hide-footer="true"
      @cancel="chatStore.showFileDialog = false"
      :width="800"
    >
      <file-select @selected="chatStore.selectFile" />
    </model-dialog>

    <model-dialog
      v-model="chatStore.showNoticeDialog"
      title="网站公告"
      @cancel="chatStore.showNoticeDialog = false"
      confirm-text="我知道了，不再提示"
      @confirm="chatStore.notShowNotice"
      :width="800"
    >
      <div class="p-4 pt-0 vue-message">
        <div v-html="chatStore.notice"></div>
      </div>
    </model-dialog>

    <el-dialog
      v-model="chatStore.showLoginDialog"
      width="500px"
      @close="store.setShowLoginDialog(false)"
    >
      <template #header>
        <div
          class="text-center text-xl"
          style="color: var(--theme-text-color-primary)"
        >
          登录后解锁更多功能
        </div>
      </template>
      <div class="p-4 pt-2 pb-2">
        <LoginDialog />
      </div>
    </el-dialog>
  </el-container>
</template>

<script setup>
  import AppListItem from '@/components/AppListItem.vue'
  import BackBottom from '@/components/BackBottom.vue'
  import BackTop from '@/components/BackTop.vue'
  import ChatList from '@/components/ChatList.vue'
  import Copyright from '@/components/Copyright.vue'
  import FileList from '@/components/FileList.vue'
  import FileSelect from '@/components/FileSelect.vue'
  import LoginDialog from '@/components/LoginDialog.vue'
  import MessageAssistant from '@/components/MessageAssistant.vue'
  import MessageUser from '@/components/MessageUser.vue'
  import ModelDialog from '@/components/ModelDialog.vue'
  import { useChatStore } from '@/js/store/chat.js'
  import { useSharedStore } from '@/js/cache/sharedata'
  import { httpGet } from '@/js/utils/http.js'
  import { nextTick, onMounted, ref } from 'vue'
  import { useRouter } from 'vue-router'

  const chatStore = useChatStore()
  const router = useRouter()
  const chatListRef = ref(null)
  const hotApps = ref([])
  const appCategories = ref([])
  const selectedCategory = ref(0)
  const store = useSharedStore()
  const version = import.meta.env.VITE_APP_VERSION

  const handleKeyDown = (e) => {
    if (e.key === 'Enter' && e.keyCode === 13) {
      e.preventDefault()
    }
    // 按下 shift + enter 换行
    if (e.key === 'Enter' && e.shiftKey) {
      const start = e.target.selectionStart
      const end = e.target.selectionEnd
      chatStore.prompt =
        chatStore.prompt.substring(0, start) +
        '\n' +
        chatStore.prompt.substring(end)
      nextTick(() => {
        e.target.selectionStart = e.target.selectionEnd = start + 1
      })
      e.preventDefault()
    }
  }

  const showWelcomePage = () => {
    chatStore.newChat()
    chatStore.showWelcome = true
    chatStore.title.value = `我是 ${chatStore.systemConfig.title}，很高兴见到你！`
  }

  const onCategoryChange = () => {
    chatStore.fetchAppList(selectedCategory.value)
  }

  onMounted(() => {
    httpGet('/api/app/hot').then((res) => {
      hotApps.value = res.data
    })
    httpGet('/api/app/category/list', {
      creator_id: 0,
    }).then((res) => {
      if (res.data) {
        appCategories.value = res.data
      }
    })
  })

  const searchKeyword = ref('')
  const onSearch = () => {
    chatStore.fetchAppList(0, searchKeyword.value)
  }
</script>

<style scoped lang="scss">
  @use '@/assets/css/chat.scss';
  @use '@/assets/markdown/vue.css';

  .header-center {
    text-align: center;
  }

  @media (max-width: 768px) {
    .header-center {
      width: 100%;
    }
  }
</style>

<style lang="scss">
  .chat-footer {
    .input-custom {
      .el-textarea__inner {
        --el-input-bg-color: transparent;
        --el-input-focus-border-color: transparent;
        --el-input-border-color: transparent;
        --el-border-color: transparent;
        color: #404040;
        caret-color: #404040;
        font-size: 16px;
        padding: 0;
        resize: none;
        box-shadow: none;
      }
    }
  }

  .default-theme h1 {
    margin: 1rem 0.8em;
  }
</style>
