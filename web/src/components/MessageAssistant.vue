<template>
  <div class="vue-message chat-message-ai d-flex flex-column mt-1 w-100">
    <div class="d-flex min-w-100">
      <div class="icon me-3">
        <el-image :src="data.icon" />
      </div>
      <div
        v-if="
          data.content?.texts?.length === 0 &&
          data.content?.tools?.length === 0 &&
          !data.completed
        "
        class="alert alert-secondary"
        style="--bs-alert-padding-x: 0.5rem; --bs-alert-padding-y: 0.5rem"
      >
        <Loading :size="5" :spacing="2" />
      </div>
      <div v-else class="content-container">
        <div v-if="data.content?.tools?.length > 0">
          <div class="mb-3">
            <button
              class="btn btn-light btn-sm dropdown-toggle"
              v-if="data.completed"
              style="
                --bs-btn-hover-color: var(--bs-success);
                --bs-btn-padding-x: 0.9rem;
                --bs-btn-padding-y: 0.5rem;
                --bs-btn-color: var(--bs-success);
                --bs-btn-active-color: var(--bs-success);
                --bs-btn-bg: #f5f5f5;
                --bs-btn-hover-bg: #dddddd;
                --bs-btn-hover-border-color: #dddddd;
              "
              type="button"
              data-bs-toggle="collapse"
              :data-bs-target="`#tool-list-${data.id}`"
            >
              <i class="iconfont icon-success me-2"></i> 运行完毕
            </button>
            <button
              class="btn btn-light btn-sm dropdown-toggle"
              v-else
              style="
                --bs-btn-hover-color: var(--bs-primary);
                --bs-btn-padding-x: 0.9rem;
                --bs-btn-padding-y: 0.5rem;
                --bs-btn-color: var(--bs-primary);
                --bs-btn-active-color: var(--bs-primary);
                --bs-btn-bg: #f5f5f5;
                --bs-btn-hover-bg: #dddddd;
                --bs-btn-hover-border-color: #dddddd;
              "
              type="button"
              data-bs-toggle="collapse"
              :data-bs-target="`#tool-list-${data.id}`"
            >
              <span class="spinner-border spinner-border-sm me-2"></span>
              <span class="me-2"
                >正在调用
                <strong>{{
                  data.content.tools[data.content.tools.length - 1].name
                }}</strong
                >...</span
              >
            </button>
          </div>
          <div class="collapse mb-3 tool-list" :id="`tool-list-${data.id}`">
            <div class="card card-body">
              <div
                class="d-flex justify-content-between align-items-center tool-item"
                v-for="tool in data.content.tools"
                :key="tool.id"
              >
                <div class="left">
                  <i class="iconfont icon-app me-2"></i>
                  已调用 <span class="fs-14 fw-bold">{{ tool.name }}</span>
                </div>
                <div class="right">
                  <div v-if="tool.status === 'SUCCESS'">
                    <span class="me-3 fs-14 text-success"
                      >{{ (tool.spend / 1000).toFixed(2) }} s</span
                    >
                    <i class="iconfont icon-success text-success"></i>
                  </div>
                  <span
                    v-else
                    class="spinner-border spinner-border-sm me-2 text-primary"
                  ></span>
                </div>
              </div>
            </div>
          </div>
        </div>
        <div
          v-for="text in data.content.texts"
          :key="text"
          :id="`content-container-${data.id}`"
        >
          <div
            v-if="text.length > 2"
            class="content mb-2"
            :class="{ 'content-bg': data.content.type !== 'alert' }"
            v-html="md.render(processContent(text))"
          ></div>
        </div>
        <!-- 参数输入表单组件 -->
        <div v-if="inputParams.length" class="mb-3">
          <InputForm
            :params="inputParams"
            :message-id="data.id"
            @submit="handleFormSubmit"
          />
        </div>
        <!-- 问答组件 -->
        <div v-if="questionData" class="mb-3">
          <QuestionAnswer
            :question-data="questionData"
            :message-id="data.id"
            @submit="handleAnswerSubmit"
          />
        </div>
      </div>
    </div>

    <div
      class="tool-box text-black-50 mt-1 d-flex align-items-center"
      v-if="data.completed && !data.share"
    >
      <el-tooltip effect="dark" placement="top" content="复制文本">
        <i
          class="iconfont icon-copy copy-answer"
          :data-clipboard-text="data.content.texts.join('\n')"
        ></i>
      </el-tooltip>
      <el-tooltip effect="dark" placement="top" content="重新生成">
        <i class="iconfont icon-refresh" @click="regenerate"></i>
      </el-tooltip>
      <span class="ms-2 fs-14">{{ dateFormat(data.created_at) }}</span>
    </div>

    <!-- 图片预览组件 -->
    <el-image-viewer
      v-if="previewSrcList.length > 0"
      :url-list="previewSrcList"
      :initial-index="0"
      @close="previewSrcList = []"
    />
  </div>
</template>

<script setup>
  import Loading from '@/components/Loading.vue'
  import InputForm from '@/components/InputForm.vue'
  import QuestionAnswer from '@/components/QuestionAnswer.vue'
  import { httpGet } from '@/js/utils/http.js'
  import { dateFormat, processContent } from '@/js/utils/libs.js'
  import hl from 'highlight.js'
  import 'highlight.js/styles/a11y-dark.css'
  import MarkdownIt from 'markdown-it'
  import mathjaxPlugin from 'markdown-it-mathjax3'
  import { computed, nextTick, onMounted, ref, watchEffect } from 'vue'

  const previewSrcList = ref([])
  const props = defineProps({
    data: {
      type: Object,
      default: {
        content: {},
        created_at: '',
      },
    },
  })
  const messageId = ref(0)
  const emit = defineEmits(['regenerate', 'sendMessage'])
  const avatar = ref('')
  const md = MarkdownIt({
    breaks: true,
    html: true,
    linkify: true,
    typographer: true,
    highlight: function (str, lang) {
      const codeIndex =
        parseInt(Date.now()) + Math.floor(Math.random() * 10000000)
      // 显示复制代码按钮和展开/收起按钮
      const copyBtn = `<div class="flex">
      <span class="text-[12px] mr-2 text-[#00e0e0] cursor-pointer expand-btn" data-code-id="${codeIndex}">展开</span>
      <span class="copy-code-btn" data-clipboard-action="copy" data-clipboard-target="#copy-target-${codeIndex}">复制</span>
      </div><textarea style="position: absolute;top: -9999px;left: -9999px;z-index: -9999;" id="copy-target-${codeIndex}">${str.replace(
        /<\/textarea>/g,
        '&lt;/textarea>'
      )}</textarea>`
      let langHtml = ''
      let preCode = ''
      // 处理代码高亮
      if (lang && hl.getLanguage(lang)) {
        langHtml = `<span class="lang-name">${lang}</span>`
        preCode = hl.highlight(str, { language: lang }).value
      } else {
        preCode = md.utils.escapeHtml(str)
      }

      // 将代码包裹在 pre 中，添加收起状态的类
      return `<pre class="code-container flex flex-col code-collapsed rounded-lg" data-code-id="${codeIndex}">
      <div class="flex justify-between bg-[#50505a] w-full rounded-tl-[10px] rounded-tr-[10px] px-3 py-1">${langHtml}${copyBtn}</div>
      <code class="language-${lang} hljs">${preCode}</code> 
      <span class="copy-code-btn absolute right-2 bottom-1" data-clipboard-action="copy" data-clipboard-target="#copy-target-${codeIndex}">复制</span></pre>`
    },
  })
  md.use(mathjaxPlugin)

  const inputParams = computed(() => {
    const content = props.data?.content || {}
    if (Array.isArray(content.inputs) && content.inputs.length) {
      return content.inputs
    }
    if (Array.isArray(content.inputForm) && content.inputForm.length) {
      return content.inputForm
    }
    return []
  })

  const questionData = computed(() => {
    const content = props.data?.content || {}
    const answersBlock = content.answers || content.questionAnswer || null

    if (
      answersBlock &&
      Array.isArray(answersBlock.answers) &&
      answersBlock.answers.length
    ) {
      return {
        title: answersBlock.title || answersBlock.question || '',
        answers: answersBlock.answers,
      }
    }
    return null
  })

  onMounted(async () => {
    nextTick(() => {
      const contentContainer = document.getElementById(
        `content-container-${props.data.id}`
      )
      if (contentContainer) {
        // 查找所有图片元素并添加点击事件
        const imgElements = contentContainer.querySelectorAll('img')
        imgElements.forEach((img) => {
          img.addEventListener('click', handleImageClick)
        })
      }

      setupCodeBlockEvents()
      messageId.value = props.data.id
    })
  })

  const handleImageClick = (e) => {
    if (e.target.tagName === 'IMG') {
      previewSrcList.value = [e.target.src]
    }
  }

  const regenerate = async () => {
    // 补全消息ID
    if (!messageId.value) {
      const res = await httpGet('/api/chat/last-message', {
        chat_id: props.data.chat_id,
      })
      messageId.value = res.data.id
    }
    emit('regenerate', messageId.value)
  }

  // 添加代码块展开/收起功能
  const toggleCodeBlock = (codeId) => {
    const codeContainer = document.querySelector(
      `pre[data-code-id="${codeId}"]`
    )
    const expandBtn = document.querySelector(
      `.expand-btn[data-code-id="${codeId}"]`
    )

    if (codeContainer && expandBtn) {
      if (codeContainer.classList.contains('code-collapsed')) {
        codeContainer.classList.remove('code-collapsed')
        codeContainer.classList.add('code-expanded')
        expandBtn.textContent = '收起'
      } else {
        codeContainer.classList.remove('code-expanded')
        codeContainer.classList.add('code-collapsed')
        expandBtn.textContent = '展开'
      }
    }
  }

  // 监听内容变化，重新绑定事件
  watchEffect(() => {
    if (props.data.content.text) {
      nextTick(() => {
        setupCodeBlockEvents()
      })
    }
  })

  const setupCodeBlockEvents = () => {
    // 移除旧的事件监听器
    const oldBtns = document.querySelectorAll('.expand-btn')
    oldBtns.forEach((btn) => {
      btn.removeEventListener('click', handleExpandClick)
    })

    // 为展开按钮添加点击事件
    const expandBtns = document.querySelectorAll('.expand-btn')
    expandBtns.forEach((btn) => {
      btn.addEventListener('click', handleExpandClick)

      // 检查对应的代码块是否需要展开功能
      const codeId = btn.getAttribute('data-code-id')
      const codeContainer = document.querySelector(
        `pre[data-code-id="${codeId}"]`
      )
      const codeElement = codeContainer?.querySelector('.hljs')

      if (codeElement) {
        // 临时移除高度限制来获取真实高度
        const originalMaxHeight = codeElement.style.maxHeight
        codeElement.style.maxHeight = 'none'
        const realHeight = codeElement.scrollHeight
        codeElement.style.maxHeight = originalMaxHeight

        // 如果代码块高度小于等于200px，隐藏展开按钮
        if (realHeight <= 200) {
          btn.style.display = 'none'
          // 移除收起状态的类，让短代码块完全展示
          codeContainer.classList.remove('code-collapsed')
        } else {
          btn.style.display = 'inline'
        }
      }
    })
  }

  const handleExpandClick = (e) => {
    const codeId = e.target.getAttribute('data-code-id')
    toggleCodeBlock(codeId)
  }

  // 处理表单提交
  const handleFormSubmit = (params) => {
    // 构建消息
    const messages = []
    for (const key in params) {
      messages.push(`${key}: ${params[key]}`)
    }
    emit('sendMessage', messages.join('\n'))
  }

  // 处理问答提交
  const handleAnswerSubmit = (answer) => {
    emit('sendMessage', answer)
  }
</script>

<style lang="scss">
  @use '@/assets/markdown/vue.css';

  //font-family: Menlo, "Roboto Mono", "Courier New", Courier, monospace, "Inter", sans-serif !important;
  //src: url(https://fonts.gstatic.com/s/robotomono/v6/L0x5DF4xlVMF-BfR8bXMIjhLq-Yg.woff2);
  .chat-message-ai {
    --bs-md-zoom: 1.143;
    --bs-font: '微软雅黑', 'Microsoft YaHei';
    // 引用快样式
    --quote-bg-color: #e0dfff;
    --quote-text-color: #333;

    align-items: flex-start;
    margin-bottom: 32px;
    position: relative;
    .content-container {
      padding-right: 45px;
      // max-width: 800px;
      .collapse {
        visibility: inherit !important;
      }
    }

    .icon {
      .el-image {
        width: 30px;
        height: 30px;
        border-radius: 50%;
        box-shadow: 0 0 0 2px rgba(var(--bs-primary-rgb), 0.5);
      }
    }

    .content {
      color: #404040;
      font-size: calc(var(--bs-md-zoom) * 14px);
      line-height: calc(var(--bs-md-zoom) * 25px);
      font-family: var(--bs-font) sans-serif;
      max-width: max-content;
      word-break: break-all;
      overflow-wrap: break-word;

      p,
      li,
      h1,
      h2,
      h3,
      h4,
      h5,
      h6,
      span {
        code {
          font-size: 0.875em;
          font-weight: 600;
          font-family: var(--bs-font), sans-serif;
          background-color: #ececec;
          color: #000000;
          border-radius: 4px;
          padding: 0.15rem 0.3rem;
        }
      }

      // 代码快
      blockquote {
        margin: 0 0 0.5rem 0;
        background-color: var(--quote-bg-color);
        padding: 0.8rem 1.5rem;
        color: var(--quote-text-color);
        border-left: 0.4rem solid #6b50e1; /* 紫色边框 */
        font-size: 16px;
        line-height: 1.6;
      }

      img {
        max-width: 600px;
        height: auto;
      }

      p:last-child {
        margin-bottom: 0;
      }

      p:first-child {
        margin-top: 0;
      }

      .code-container {
        position: relative;
        margin: 0 !important;
        display: flex;
        .hljs {
          padding-top: 40px;
          width: 100%;
          white-space: pre-wrap;
          word-break: break-all;
        }

        .copy-code-btn {
          cursor: pointer;
          font-size: 12px;
          color: #c1c1c1;

          &:hover {
            color: #ffffff;
          }
        }
      }

      // 添加代码块展开/收起样式
      .code-collapsed {
        .hljs {
          max-height: 200px;
          overflow: hidden;
          position: relative;
          transition: max-height 0.3s ease;

          &::after {
            content: '';
            position: absolute;
            bottom: 0;
            left: 0;
            right: 0;
            height: 30px;
            background: linear-gradient(transparent, #2b2b2b);
            pointer-events: none;
          }
        }
      }

      .code-expanded {
        .hljs {
          max-height: none;
          overflow: auto;
          transition: max-height 0.3s ease;

          &::after {
            display: none;
          }
        }
      }

      .expand-btn {
        transition: color 0.2s ease;

        &:hover {
          color: #20a0ff !important;
        }
      }

      .lang-name {
        color: #00e0e0;
      }
    }

    .content-bg {
      background-color: #f5f5f5;
      padding: 12px 16px;
      border-radius: 10px;
      word-break: break-all;
      overflow-wrap: break-word;
    }

    .tool-box {
      padding-left: 45px;

      .iconfont {
        cursor: pointer;
        font-size: 18px;
        padding: 2px 5px;
        border-radius: 6px;

        &:hover {
          background-color: #f1f1f1;
        }
      }
    }

    .tool-list {
      .card-body {
        --bs-card-spacer-y: 0.5rem;
        --bs-card-spacer-x: 0.5rem;

        .tool-item {
          padding: 5px 10px;
          border-radius: 6px;
          &:hover {
            background-color: #f1f1f1;
            cursor: pointer;
          }
        }
      }
    }
  }
</style>
