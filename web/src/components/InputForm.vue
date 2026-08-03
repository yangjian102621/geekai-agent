<template>
  <div
    class="input-form-container"
    :class="
      isSubmitting ? 'opacity-70 cursor-not-allowed pointer-events-none' : ''
    "
  >
    <div class="input-form-card">
      <div v-for="(param, index) in params" :key="index" class="input-field">
        <label class="input-label">{{ param.name }}</label>

        <!-- String 类型 -->
        <el-input
          v-if="getInputType(param) === 'string'"
          v-model="formData[param.name]"
          :placeholder="'Please enter content.'"
          class="input-custom"
          :class="{ 'input-error': errors[param.name] }"
        />

        <!-- Integer 类型 -->
        <el-input-number
          v-else-if="getInputType(param) === 'integer'"
          v-model="formData[param.name]"
          :precision="0"
          :placeholder="'Please enter content.'"
          class="input-custom"
          :class="{ 'input-error': errors[param.name] }"
          style="width: 100%"
        />

        <!-- Float 类型 -->
        <el-input-number
          v-else-if="getInputType(param) === 'float'"
          v-model="formData[param.name]"
          :placeholder="'Please enter content.'"
          class="input-custom"
          :class="{ 'input-error': errors[param.name] }"
          style="width: 100%"
        />

        <!-- Boolean 类型 -->
        <el-switch
          v-else-if="getInputType(param) === 'boolean'"
          v-model="formData[param.name]"
          class="input-custom"
        />

        <!-- Time 类型 (assistType: 10000) -->
        <el-date-picker
          v-else-if="getInputType(param) === 'time'"
          v-model="formData[param.name]"
          type="datetime"
          value-format="YYYY-MM-DD HH:mm:ss"
          format="YYYY-MM-DD HH:mm:ss"
          :placeholder="'Please select time.'"
          class="input-custom"
          :class="{ 'input-error': errors[param.name] }"
          style="width: 100%"
        />

        <!-- Object 类型 (JSON) -->
        <el-input
          v-else-if="getInputType(param) === 'object'"
          v-model="formData[param.name]"
          type="textarea"
          :autosize="{ minRows: 3, maxRows: 6 }"
          :placeholder="'Please enter JSON content.'"
          class="input-custom"
          :class="{ 'input-error': errors[param.name] }"
        />

        <!-- File 类型 (assistType: 2) -->
        <div
          v-else-if="getInputType(param) === 'file'"
          class="file-upload-container"
        >
          <div class="file-upload-row">
            <el-upload
              :auto-upload="true"
              :show-file-list="false"
              :http-request="(file) => handleFileUpload(file, param)"
              class="file-upload"
              :class="{ 'input-error': errors[param.name] }"
            >
              <el-button type="primary" plain class="upload-btn">
                <i class="iconfont icon-upload me-1"></i>
                上传文件
              </el-button>
            </el-upload>
            <span class="file-hint">支持图片、音频、视频及常见文档</span>
          </div>
          <div v-if="formData[param.name]" class="file-preview">
            <template v-if="getPreviewType(formData[param.name]) === 'image'">
              <el-image
                :src="formData[param.name]"
                fit="cover"
                class="preview-image"
                :preview-src-list="[formData[param.name]]"
              />
            </template>
            <template
              v-else-if="getPreviewType(formData[param.name]) === 'audio'"
            >
              <audio
                class="preview-audio"
                controls
                :src="formData[param.name]"
              ></audio>
            </template>
            <template
              v-else-if="getPreviewType(formData[param.name]) === 'video'"
            >
              <video
                class="preview-video"
                controls
                :src="formData[param.name]"
              ></video>
            </template>
            <template v-else>
              <div class="file-preview-generic">
                <el-image
                  :src="GetFileIcon(getFileExt(formData[param.name]))"
                  fit="cover"
                  class="preview-thumb"
                />
                <div class="file-info">
                  <div class="file-name">
                    {{ getFileName(formData[param.name]) }}
                  </div>
                  <a :href="formData[param.name]" target="_blank">查看文件</a>
                </div>
              </div>
            </template>
            <div class="file-preview-actions">
              <el-button
                type="danger"
                link
                size="small"
                @click="clearFile(param)"
              >
                清除
              </el-button>
            </div>
          </div>
        </div>

        <div v-if="errors[param.name]" class="error-message">
          {{ errors[param.name] }}
        </div>
      </div>
      <div class="submit-container">
        <el-button
          type="primary"
          class="submit-btn"
          :disabled="isSubmitting"
          @click="handleSubmit"
        >
          提交
        </el-button>
      </div>
    </div>
  </div>
</template>

<script setup>
  import { ref, reactive } from 'vue'
  import { frontUploadFile } from '@/js/store/common.js'
  import {
    GetFileIcon,
    getFileExt,
    isImage,
    replaceURL,
  } from '@/js/utils/libs.js'

  const props = defineProps({
    params: {
      type: Array,
      required: true,
      default: () => [],
    },
    messageId: {
      type: Number,
      default: 0,
    },
  })

  const emit = defineEmits(['submit'])

  const formData = reactive({})
  const errors = reactive({})
  const isSubmitting = ref(false)

  const audioExts = ['.mp3', '.wav', '.ogg', '.m4a']
  const videoExts = ['.mp4', '.avi', '.mov', '.mkv', '.webm']

  // 判断输入类型
  const getInputType = (param) => {
    // Time 类型：type 为 string 且 assistType 为 10000
    if (param.type === 'string' && param.assistType === 10000) {
      return 'time'
    }
    // File 类型：type 为 string 且 assistType 为 2
    if (param.type === 'string' && param.assistType === 2) {
      return 'file'
    }
    // 其他类型直接返回 type
    return param.type
  }

  // 初始化表单数据
  props.params.forEach((param) => {
    const inputType = getInputType(param)
    switch (inputType) {
      case 'integer':
      case 'float':
        formData[param.name] = 0
        break
      case 'boolean':
        formData[param.name] = false
        break
      case 'time':
        formData[param.name] = null
        break
      case 'object':
      case 'string':
      case 'file':
      default:
        formData[param.name] = ''
        break
    }
  })

  // 验证表单
  const validateForm = () => {
    let isValid = true
    // 清除之前的错误
    Object.keys(errors).forEach((key) => {
      errors[key] = ''
    })

    props.params.forEach((param) => {
      const inputType = getInputType(param)
      const value = formData[param.name]

      // 必填验证
      if (param.required) {
        if (value === null || value === undefined || value === '') {
          // boolean 类型允许 false 值
          if (inputType === 'boolean' && value === false) {
            // boolean 类型的 false 是有效值，不报错
          } else {
            errors[param.name] = '此字段为必填项'
            isValid = false
            return
          }
        }
      }

      // 类型验证
      if (value !== null && value !== undefined && value !== '') {
        // Integer 和 Float 类型验证
        if (inputType === 'integer' || inputType === 'float') {
          if (isNaN(value)) {
            errors[param.name] = '请输入有效的数字'
            isValid = false
          }
        }

        // Object 类型验证 JSON 格式
        if (inputType === 'object') {
          try {
            JSON.parse(value)
          } catch (e) {
            errors[param.name] = '请输入有效的 JSON 格式'
            isValid = false
          }
        }
      }
    })

    return isValid
  }

  // 处理文件上传
  const handleFileUpload = (file, param) => {
    frontUploadFile(file, (res) => {
      // 如果是相对路径，则转换为绝对路径
      formData[param.name] = replaceURL(res.url)
      // 清除可能的错误
      if (errors[param.name]) {
        errors[param.name] = ''
      }
    })
  }

  // 清除文件
  const clearFile = (param) => {
    formData[param.name] = ''
  }

  // 获取文件名
  const getFileName = (url) => {
    if (!url) return ''
    const parts = url.split('/')
    return parts[parts.length - 1]
  }

  const getPreviewType = (url) => {
    if (!url) return 'other'
    if (isImage(url)) {
      return 'image'
    }
    const ext = getFileExt(url)
    if (audioExts.includes(ext)) {
      return 'audio'
    }
    if (videoExts.includes(ext)) {
      return 'video'
    }
    return 'other'
  }

  // 提交表单
  const handleSubmit = () => {
    if (isSubmitting.value) return

    // 清除之前的错误
    Object.keys(errors).forEach((key) => {
      errors[key] = ''
    })

    // 验证表单
    if (!validateForm()) {
      return
    }

    // 将参数名称和值组成键值对
    const params = {}
    props.params.forEach((param) => {
      const inputType = getInputType(param)
      const value = formData[param.name]

      // boolean 类型的 false 是有效值，需要包含
      if (inputType === 'boolean') {
        params[param.name] = value
        return
      }

      // 跳过空值
      if (value === null || value === undefined || value === '') {
        return
      }

      // 根据类型格式化值
      let formattedValue = value

      // Time 类型已经是字符串格式（YYYY-MM-DD HH:mm:ss），直接使用
      // Object 类型保持 JSON 字符串格式
      // File 类型保持 URL 字符串格式
      // 其他类型直接使用原值

      params[param.name] = formattedValue
    })

    // 将键值对格式化为字符串，格式：key1:value1\nkey2:value2
    const paramString = Object.entries(params)
      .map(([key, value]) => {
        // 对于对象类型，保持 JSON 格式
        if (typeof value === 'object' && value !== null) {
          return `${key}:${JSON.stringify(value)}`
        }
        return `${key}:${value}`
      })
      .join('\n')

    // 触发提交事件
    emit('submit', params)
    isSubmitting.value = true
  }
</script>

<style scoped lang="scss">
  .input-form-container {
    width: 100% !important;
    max-width: 400px;
  }

  .input-form-card {
    background-color: #ffffff;
    border: 1px solid #e1e1e1;
    border-radius: 12px;
    padding: 20px;
    margin-bottom: 10px;
  }

  .input-field {
    margin-bottom: 16px;

    &:last-of-type {
      margin-bottom: 20px;
    }
  }

  .input-label {
    display: block;
    font-size: 14px;
    font-weight: 500;
    color: #333;
    margin-bottom: 8px;
  }

  .error-message {
    color: #f56c6c;
    font-size: 12px;
    margin-top: 4px;
  }

  .file-upload-container {
    width: 100%;
  }

  .file-upload-row {
    display: flex;
    align-items: center;
    gap: 12px;
  }

  .file-upload {
    width: auto;

    &.input-error {
      :deep(.el-button) {
        border-color: #f56c6c;
        color: #f56c6c;
      }
    }
  }

  .upload-btn {
    display: inline-flex;
    align-items: center;
    gap: 6px;
  }

  .file-hint {
    font-size: 12px;
    color: #909399;
  }

  .file-preview {
    margin-top: 8px;
    padding: 12px;
    background-color: #f5f7fa;
    border-radius: 8px;
    display: flex;
    flex-direction: column;
    gap: 10px;
  }

  .preview-image {
    width: 100%;
    max-height: 200px;
    border-radius: 8px;
    object-fit: cover;
  }

  .preview-audio,
  .preview-video {
    width: 100%;
    border-radius: 8px;
    background-color: #fff;
  }

  .file-preview-generic {
    display: flex;
    align-items: center;
    gap: 12px;
  }

  .preview-thumb {
    width: 48px;
    height: 48px;
    border-radius: 6px;
    object-fit: cover;
  }

  .file-info {
    flex: 1;
    display: flex;
    flex-direction: column;
    gap: 4px;
  }

  .file-info .file-name {
    font-size: 13px;
    color: #303133;
    word-break: break-all;
  }

  .file-preview-actions {
    display: flex;
    justify-content: flex-end;
  }

  .submit-container {
    display: flex;
    justify-content: center;
    margin-top: 20px;
  }

  .submit-btn {
    width: 100%;
    height: 44px;
    border-radius: 8px;
    font-size: 16px;
    font-weight: 500;
    background-color: #6b50e1;
    border-color: #6b50e1;
    color: #ffffff;

    &:hover {
      background-color: #5a3fd0;
      border-color: #5a3fd0;
    }

    &:disabled {
      background-color: #c0c4cc;
      border-color: #c0c4cc;
      cursor: not-allowed;
      color: #ffffff;
    }
  }
</style>
