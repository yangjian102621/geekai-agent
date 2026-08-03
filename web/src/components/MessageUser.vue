<template>
  <div class="chat-message-user d-flex flex-column mt-1 w-100">
    <div
      class="d-flex w-100 justify-content-end"
      :id="`content-container-${data.id}`"
    >
      <div class="content me-3 d-flex flex-column align-items-end">
        <!-- 文件列表 -->
        <div v-if="fileList.length > 0" class="file-list-box">
          <div
            v-for="file in fileList"
            :key="file.ext"
            class="d-flex justify-content-end mb-2"
          >
            <div class="image" v-if="isImage(file.ext)">
              <el-image :src="file.url" fit="cover" />
            </div>
            <div class="item" v-else>
              <div class="icon">
                <el-image :src="GetFileIcon(file.ext)" fit="cover" />
              </div>
              <div class="body">
                <div class="title">
                  <el-link
                    :href="file.url"
                    target="_blank"
                    style="--el-font-weight-primary: bold"
                    >{{ file.name }}
                  </el-link>
                </div>
                <div class="info">
                  <span>{{ getFileExtType(file.ext) }}</span>
                  <span class="ms-2">{{ FormatFileSize(file.size) }}</span>
                </div>
              </div>
            </div>
          </div>
        </div>
        <!-- 文本内容 -->
        <div
          class="text-wrap"
          v-if="data.content.texts[0] != ''"
          v-html="md.render(data.content.texts[0])"
        ></div>
      </div>
      <div class="icon">
        <el-image :src="avatar" />
      </div>
    </div>
    <div
      class="tool-box text-black-50 d-flex justify-content-end align-items-center w-100 mt-2"
    >
      <el-tooltip effect="dark" placement="top" content="复制文本">
        <i
          class="iconfont icon-copy copy-answer"
          :data-clipboard-text="data.content.texts[0] || ''"
        ></i>
      </el-tooltip>
      <el-tag type="info" round class="ms-2">{{
        dateFormat(data.created_at)
      }}</el-tag>
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
  import { ref, computed, onMounted, nextTick } from 'vue'
  import MarkdownIt from 'markdown-it'
  import { checkSession } from '@/js/cache/session.js'
  import {
    FormatFileSize,
    GetFileIcon,
    dateFormat,
    getFileExt,
    getFileExtType,
    isImage,
    processContent,
  } from '@/js/utils/libs.js'

  const previewSrcList = ref([])
  const props = defineProps({
    data: {
      type: Object,
      default: {
        content: '',
        created_at: '',
      },
    },
  })
  const avatar = ref('')

  const md = new MarkdownIt({
    breaks: true,
    linkify: true,
  })

  const fileList = computed(() => {
    if (!props.data.content.files) {
      return []
    }

    return props.data.content.files.map((file) => {
      return {
        ...file,
        ext: getFileExt(file.url),
      }
    })
  })
  onMounted(async () => {
    const user = await checkSession()
    avatar.value = user.avatar || '/images/avatar/user.png'

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
    })
  })

  const handleImageClick = (e) => {
    if (e.target.tagName === 'IMG') {
      previewSrcList.value = [e.target.src]
    }
  }
</script>

<style lang="scss">
  .chat-message-user {
    --bs-md-zoom: 1.143;
    align-items: flex-start;
    margin-bottom: 32px;
    position: relative;

    .icon {
      .el-image {
        width: 30px;
        height: 30px;
        border-radius: 50%;
        box-shadow: 0 0 0 2px rgba(var(--bs-primary-rgb), 0.3);
      }
    }

    .content {
      margin-left: 45px;
      color: #262626;
      font-size: calc(var(--bs-md-zoom) * 14px);
      line-height: calc(var(--bs-md-zoom) * 25px);

      .file-list-box {
        display: flex;
        flex-flow: column;

        .image {
          display: flex;
          flex-flow: row;
          margin-right: 10px;
          position: relative;
          max-width: 50%;

          .el-image {
            border: 1px solid #e3e3e3;
            border-radius: 10px;
          }
        }

        .item {
          display: flex;
          font-size: 14px;
          flex-flow: row;
          border-radius: 10px;
          background-color: #f5f7fc;
          border: 1px solid #e3e3e3;
          color: var(--theme-text-color-primary);
          padding: 10px;

          .icon {
            .el-image {
              width: 44px;
              height: 44px;
              border-radius: 6px;
            }
          }

          .body {
            margin-left: 8px;
            font-size: 14px;

            .title {
              font-weight: bold;
              line-height: 24px;
            }

            .info {
              color: #b4b4b4;
              font-size: 12px;
              line-height: 20px;
            }
          }
        }
      }

      .text-wrap {
        background-color: var(--bs-primary-light);
        border-radius: 10px;
        padding: 10px 20px;
        max-width: max-content;
        word-break: break-all;
        overflow-wrap: break-word;

        p {
          margin: 0;
          // white-space: pre-wrap;
        }

        img {
          max-width: 500px;
          height: auto;
          max-height: 500px;
          object-fit: cover;
        }
      }

      img {
        max-width: 100%;
        height: auto;
      }
    }

    .tool-box {
      padding-right: 45px;

      .iconfont {
        cursor: pointer;
        font-size: 18px;
        padding: 5px;
        border-radius: 6px;

        &:hover {
          background-color: #f1f1f1;
        }
      }
    }
  }
</style>
