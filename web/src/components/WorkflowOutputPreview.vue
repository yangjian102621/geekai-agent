<template>
  <div class="workflow-output-preview">
    <div v-if="parsedOutput && parsedOutput.type" class="preview-content">
      <!-- 文本类型 -->
      <div v-if="parsedOutput.type === 'text'" class="output-text">
        <div class="bg-light p-3 rounded">
          <p class="m-0" style="white-space: pre-wrap">
            {{ parsedOutput.text || '' }}
          </p>
        </div>
      </div>

      <!-- 图片类型 -->
      <div v-else-if="parsedOutput.type === 'image'" class="output-image">
        <div class="mb-3">
          <el-image
            :src="parsedOutput.url"
            :preview-src-list="[parsedOutput.url]"
            fit="contain"
            class="rounded border output-image-el"
            :lazy="false"
          >
            <template #error>
              <div
                class="d-flex justify-content-center align-items-center"
                style="height: 200px; width: 100%"
              >
                <i class="iconfont icon-image text-4xl text-muted"></i>
              </div>
            </template>
          </el-image>
        </div>
        <div v-if="parsedOutput.text" class="mb-3">
          <p class="m-0 text-muted">{{ parsedOutput.text }}</p>
        </div>
        <el-button
          type="primary"
          @click="downloadFile(parsedOutput.url, 'image')"
        >
          <i class="iconfont icon-download mr-1 !text-sm"></i> 下载图片
        </el-button>
      </div>

      <!-- 音频类型 -->
      <div v-else-if="parsedOutput.type === 'audio'" class="output-audio">
        <div class="mb-3">
          <audio :src="parsedOutput.url" controls class="w-100">
            您的浏览器不支持音频播放
          </audio>
        </div>
        <div v-if="parsedOutput.text" class="mb-3">
          <p class="m-0 text-muted">{{ parsedOutput.text }}</p>
        </div>
        <el-button
          type="primary"
          @click="downloadFile(parsedOutput.url, 'audio')"
        >
          <i class="iconfont icon-download mr-1 !text-sm"></i> 下载音频
        </el-button>
      </div>

      <!-- 视频类型 -->
      <div v-else-if="parsedOutput.type === 'video'" class="output-video">
        <div class="mb-3">
          <video
            :src="parsedOutput.url"
            controls
            class="w-100 rounded border"
            style="max-height: 500px"
          >
            您的浏览器不支持视频播放
          </video>
        </div>
        <div v-if="parsedOutput.text" class="mb-3">
          <p class="m-0 text-muted">{{ parsedOutput.text }}</p>
        </div>
        <el-button
          type="primary"
          @click="downloadFile(parsedOutput.url, 'video')"
        >
          <i class="iconfont icon-download mr-1 !text-sm"></i> 下载视频
        </el-button>
      </div>

      <!-- 文件类型 -->
      <div v-else-if="parsedOutput.type === 'file'" class="output-file">
        <div class="d-flex align-items-center gap-2 mb-3">
          <i class="iconfont icon-file text-3xl text-primary"></i>
          <div class="flex-grow-1">
            <div class="fw-bold">{{ getFileName(parsedOutput.url) }}</div>
            <div v-if="parsedOutput.text" class="text-muted small">
              {{ parsedOutput.text }}
            </div>
          </div>
        </div>
        <el-button
          type="primary"
          @click="downloadFile(parsedOutput.url, 'file')"
        >
          <i class="iconfont icon-download mr-1 !text-sm"></i> 下载文件
        </el-button>
      </div>

      <!-- 未知类型或未解析成功 -->
      <div v-else class="output-raw">
        <Alert type="warning">
          未解析到可预览实体，请检查输出是否符合要求。
        </Alert>

        <h6>工作流结束节点正确输出示例：</h6>
        <div
          class="bg-light p-3 rounded font-monospace small overflow-auto mb-3"
        >
          <p class="text-muted">图片类型：</p>
          <pre class="mt-2">{{
            JSON.stringify(
              { type: 'image', url: 'https://example.com/image.jpg' },
              null,
              2
            )
          }}</pre>

          <p class="text-muted">音频类型：</p>
          <pre class="mt-2">{{
            JSON.stringify(
              { type: 'audio', url: 'https://example.com/audio.mp3' },
              null,
              2
            )
          }}</pre>

          <p class="text-muted">视频类型：</p>
          <pre class="mt-2">{{
            JSON.stringify(
              { type: 'video', url: 'https://example.com/video.mp4' },
              null,
              2
            )
          }}</pre>

          <p class="text-muted">文件类型：</p>
          <pre class="mt-2">{{
            JSON.stringify(
              { type: 'file', url: 'https://example.com/file.bin' },
              null,
              2
            )
          }}</pre>
        </div>
      </div>
    </div>
    <div v-else class="text-muted">
      <p class="m-0">暂无执行结果</p>
    </div>
  </div>
</template>

<script setup>
  import { computed } from 'vue'
  import Alert from '@/components/Alert.vue'
  const props = defineProps({
    output: {
      type: [Object, String],
      default: null,
    },
  })

  // 解析 Output 字段
  const parsedOutput = computed(() => {
    if (!props.output) return null

    // 如果 output 是字符串，尝试解析
    if (typeof props.output === 'string') {
      try {
        const parsed = JSON.parse(props.output)
        // 检查是否有 Output 字段（嵌套的 JSON 字符串）
        if (parsed.Output && typeof parsed.Output === 'string') {
          try {
            return JSON.parse(parsed.Output)
          } catch (e) {
            console.warn('Failed to parse nested Output:', e)
          }
        }
        return parsed
      } catch (e) {
        console.warn('Failed to parse output string:', e)
        return null
      }
    }

    // 如果 output 是对象
    if (typeof props.output === 'object') {
      // 检查是否有 Output 字段（可能是字符串）
      if (props.output.Output && typeof props.output.Output === 'string') {
        try {
          return JSON.parse(props.output.Output)
        } catch (e) {
          console.warn('Failed to parse Output field:', e)
        }
      }
      // 如果直接有 type 字段，说明已经是解析后的格式
      if (props.output.type) {
        return props.output
      }
      // 否则返回原对象
      return props.output
    }

    return null
  })

  // 下载文件
  const downloadFile = (url, type) => {
    if (!url) return

    // 创建一个临时的 a 标签来触发下载
    const link = document.createElement('a')
    link.href = url
    link.download = getFileName(url) || `download.${getFileExtension(type)}`
    link.target = '_blank'
    document.body.appendChild(link)
    link.click()
    document.body.removeChild(link)
  }

  // 获取文件名
  const getFileName = (url) => {
    if (!url) return ''
    try {
      const urlObj = new URL(url)
      const pathname = urlObj.pathname
      const fileName = pathname.split('/').pop()
      return fileName || 'download'
    } catch (e) {
      // 如果不是有效的 URL，尝试从路径中提取
      const parts = url.split('/')
      return parts[parts.length - 1] || 'download'
    }
  }

  // 根据类型获取文件扩展名
  const getFileExtension = (type) => {
    const extensions = {
      image: 'jpg',
      audio: 'mp3',
      video: 'mp4',
      file: 'bin',
    }
    return extensions[type] || 'bin'
  }
</script>

<style scoped lang="scss">
  .workflow-output-preview {
    .preview-content {
      // 移除默认的 margin 和 padding，避免空白
      margin: 0;
      padding: 0;
    }

    .output-text {
      p {
        line-height: 1.6;
      }
    }

    .output-image {
      .output-image-el {
        background-color: #f5f7fa;
        // 让图片容器紧凑显示，不留下多余空白
        display: block;
        max-width: 100%;
        max-height: 500px;

        :deep(.el-image__inner) {
          max-width: 100%;
          max-height: 500px;
          width: auto !important;
          height: auto !important;
          object-fit: contain;
        }

        :deep(.el-image__wrapper) {
          display: inline-block;
          max-width: 100%;
        }
      }
    }

    .output-video {
      video {
        background-color: #f5f7fa;
        display: block;
        width: 100%;
      }
    }

    .output-audio {
      audio {
        height: 40px;
        display: block;
        width: 100%;
      }
    }

    .output-file {
      .icon-file {
        font-size: 2rem;
      }
    }

    .output-raw {
      pre {
        max-height: 400px;
        overflow: auto;
      }
    }
  }
</style>
