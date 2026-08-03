<template>
  <el-card class="param-config-card">
    <template #header>
      <div class="param-config-header">
        <div class="header-left">
          <span>{{ title }}</span>
          <span class="subtitle" v-if="params.length">
            共 {{ params.length }} 个字段
          </span>
        </div>
        <el-button type="primary" size="small" @click="addParam">
          <i class="iconfont icon-plus mr-1"></i>
          新增字段
        </el-button>
      </div>
    </template>

    <div v-if="!params.length" class="flex px-3 justify-center items-center">
      <div class="text-gray-400 text-sm flex flex-col items-center gap-2">
        <i class="iconfont icon-empty-box text-4xl"></i>
        暂无参数
      </div>
    </div>

    <div v-else class="param-list">
      <div
        class="param-card"
        v-for="(row, index) in params"
        :key="row.id || index"
      >
        <div class="param-card__header">
          <div class="param-title">
            字段 {{ index + 1 }}
            <span class="param-key" v-if="row.name">({{ row.name }})</span>
          </div>
          <div class="param-actions">
            <el-tooltip content="删除字段" placement="top">
              <el-button
                type="danger"
                link
                @click="removeParam(row)"
                class="text-danger"
              >
                <i class="iconfont icon-remove"></i>
              </el-button>
            </el-tooltip>
          </div>
        </div>

        <div class="param-grid">
          <div class="form-item">
            <label>字段名称 *</label>
            <el-input v-model="row.label" placeholder="如：客户姓名" />
          </div>
          <div class="form-item">
            <label>字段标识 (key) *</label>
            <el-input
              v-model="row.name"
              placeholder="如：customer_name"
              maxlength="40"
              show-word-limit
            />
          </div>
          <div class="form-item">
            <label>字段类型 *</label>
            <el-select
              v-model="row.type"
              placeholder="请选择字段类型"
              @change="changeParamType(row)"
            >
              <el-option
                v-for="type in paramTypes"
                :key="type.value"
                :label="type.label"
                :value="type.value"
              />
            </el-select>
          </div>
          <div class="form-item">
            <label>是否必填</label>
            <el-radio-group v-model="row.required">
              <el-radio :value="true">是</el-radio>
              <el-radio :value="false">否</el-radio>
            </el-radio-group>
          </div>
        </div>

        <div class="param-card__section" v-if="uploadTypes.includes(row.type)">
          <label>上传配置</label>
          <div class="upload-row">
            <el-upload
              :auto-upload="true"
              :show-file-list="false"
              :http-request="(file) => handleUpload(file, row)"
              :accept="getAcceptByType(row.type)"
            >
              <el-button type="primary">
                <i class="iconfont icon-upload mr-1"></i> 上传文件
              </el-button>
            </el-upload>
            <span class="tips">单文件，最大</span>
            <el-input-number
              v-model="row.max_filesize"
              :min="1"
              :max="20"
              :step="1"
              controls-position="right"
              style="width: 100px"
            />
            <span class="tips">MB</span>
          </div>
          <div class="upload-preview" v-if="row.default">
            <el-image
              v-if="row.type === 'Image'"
              :src="row.default"
              fit="cover"
              class="preview-thumb"
              @click="handlePreview(row.default)"
            />
            <audio v-else-if="row.type === 'Audio'" :src="row.default" controls>
              您的浏览器不支持音频播放
            </audio>
            <video v-else-if="row.type === 'Video'" :src="row.default" controls>
              您的浏览器不支持视频播放
            </video>
            <a v-else :href="row.default" target="_blank" class="other-file">
              <el-tooltip content="点击下载文件" placement="top">
                <el-image
                  :src="GetFileIcon(getFileExt(row.default))"
                  fit="cover"
                  class="preview-thumb"
                />
              </el-tooltip>
            </a>
          </div>
        </div>

        <div class="param-card__section" v-else>
          <label>默认值</label>
          <div class="stack gap-2">
            <el-input
              v-if="row.type === 'String'"
              type="textarea"
              :autosize="{ minRows: 1, maxRows: 4 }"
              v-model="row.default"
              placeholder="默认内容"
            />
            <el-input-number
              v-else-if="row.type === 'Number'"
              v-model="row.default"
              :min="0"
              :step="1"
            />

            <el-switch
              v-else-if="row.type === 'Boolean'"
              v-model="row.default"
            />
            <el-date-picker
              v-else-if="row.type === 'DateTime'"
              v-model="row.default"
              type="datetime"
              value-format="YYYY-MM-DD HH:mm:ss"
              format="YYYY-MM-DD HH:mm:ss"
            />
            <el-select
              v-else-if="row.type === 'Select'"
              v-model="row.default"
              placeholder="请选择默认选项"
            >
              <el-option
                v-for="option in row.options"
                :key="option"
                :label="option"
                :value="option"
              />
            </el-select>
            <el-radio-group
              v-else-if="row.type === 'Radio'"
              v-model="row.default"
            >
              <el-radio
                v-for="option in row.options"
                :key="option"
                :value="option"
              >
                {{ option }}
              </el-radio>
            </el-radio-group>
            <el-checkbox-group
              v-else-if="row.type === 'CheckBox'"
              v-model="row.default"
            >
              <el-checkbox
                v-for="option in row.options"
                :key="option"
                :value="option"
              >
                {{ option }}
              </el-checkbox>
            </el-checkbox-group>
          </div>

          <div v-if="optionsTypes.includes(row.type)">
            <label class="mb-2">选项列表</label>
            <ItemsInput
              :value="row.options || []"
              @update:value="(tags) => (row.options = tags)"
            />
          </div>
        </div>
      </div>
    </div>

    <el-image-viewer
      @close="
        () => {
          previewURL = ''
        }
      "
      v-if="previewURL !== ''"
      :url-list="[previewURL]"
    />
  </el-card>
</template>

<script setup>
  import { computed, ref, watch } from 'vue'
  import { getFileExt, GetFileIcon } from '@/js/utils/libs.js'
  import ItemsInput from '@/components/ItemsInput.vue'
  import {
    showLoading,
    closeLoading,
    showMessageError,
  } from '@/js/utils/dialog.js'
  import { adminUploadFile } from '@/js/store/common.js'

  const uploadTypes = ['Image', 'Audio', 'Video', 'Doc', 'Zip', 'File']
  const optionsTypes = ['Radio', 'CheckBox', 'Select']
  const previewURL = ref('')
  const props = defineProps({
    modelValue: {
      type: Array,
      default: () => [],
    },
    title: {
      type: String,
      default: '参数配置',
    },
  })

  const emit = defineEmits(['update:modelValue'])

  const paramTypes = [
    { label: 'String', value: 'String' },
    { label: 'Number', value: 'Number' },
    { label: 'Boolean', value: 'Boolean' },
    { label: 'DateTime', value: 'DateTime' },
    { label: 'Image', value: 'Image' },
    { label: 'Audio', value: 'Audio' },
    { label: 'Video', value: 'Video' },
    { label: 'Doc', value: 'Doc' },
    { label: 'Zip', value: 'Zip' },
    { label: 'File', value: 'File' },
    { label: 'Select', value: 'Select' },
    { label: 'Radio', value: 'Radio' },
    { label: 'CheckBox', value: 'CheckBox' },
  ]

  const params = computed({
    get: () => props.modelValue || [],
    set: (value) => emit('update:modelValue', value),
  })

  const ensureUploadConfig = (param) => {
    if (!uploadTypes.includes(param.type || '')) return
    const size =
      param.max_filesize ??
      param.maxFilesize ??
      (param.max_filesize === 0 ? 0 : 5)
    param.max_filesize = size
    param.maxFilesize = size
  }

  const addParam = () => {
    params.value.push({
      label: '',
      name: '',
      type: 'String',
      default: '',
      required: true,
    })
  }

  const changeParamType = (param) => {
    if (param.type === 'Number') {
      param.default = 0
    } else if (param.type === 'Boolean') {
      param.default = false
    } else if (param.type === 'CheckBox' && !Array.isArray(param.default)) {
      param.default = []
    } else {
      param.default = ''
    }

    if (optionsTypes.includes(param.type)) {
      if (!param.options) {
        param.options = []
      }
      if (param.type === 'CheckBox' && !Array.isArray(param.default)) {
        param.default = []
      } else if (param.type !== 'CheckBox' && Array.isArray(param.default)) {
        param.default = ''
      }
    } else {
      delete param.options
    }

    if (uploadTypes.includes(param.type)) {
      ensureUploadConfig(param)
    } else {
      delete param.max_filesize
      delete param.maxFilesize
    }
  }

  const removeParam = (target) => {
    params.value = params.value.filter((param) => param !== target)
  }

  const getAcceptByType = (type) => {
    if (type === 'Image') return '.png,.jpg,.jpeg,.bmp,.gif'
    if (type === 'Audio') return '.mp3,.wav,.ogg'
    if (type === 'Video') return '.mp4,.avi,.mov'
    if (type === 'Doc') return '.doc,.docx,.pdf,.txt,.xls,.xlsx,.csv,.ppt,.pptx'
    if (type === 'Zip') return '.zip,.rar,.7z'
    if (type === 'File')
      return '.doc,.docx,.pdf,.txt,.xls,.xlsx,.csv,.ppt,.pptx,.zip,.rar,.7z'
    return ''
  }

  const handleUpload = (file, param) => {
    // 这里检测文件大小是否超过最大文件大小
    if (file.file.size > param.max_filesize * 1024 * 1024) {
      showMessageError(`文件不能超过 ${param.max_filesize}MB`)
      return
    }

    const formData = new FormData()
    formData.append('file', file.file, file.file.name)
    showLoading('正在上传...')
    adminUploadFile(file, (data) => {
      param.default = data.url
      closeLoading()
    })
  }

  const handlePreview = (url) => {
    previewURL.value = url
  }

  watch(
    params,
    (val) => {
      if (!Array.isArray(val)) return
      val.forEach((item) => {
        ensureUploadConfig(item)
        if (
          optionsTypes.includes(item.type || '') &&
          item.type === 'CheckBox' &&
          !Array.isArray(item.default)
        ) {
          item.default = []
        }
      })
    },
    { immediate: true, deep: true }
  )
</script>

<style scoped lang="scss">
  .param-config-card {
    .param-config-header {
      display: flex;
      align-items: center;
      justify-content: space-between;

      .header-left {
        display: flex;
        align-items: baseline;
        gap: 12px;

        .subtitle {
          font-size: 13px;
          color: #909399;
        }
      }
    }
  }

  .param-list {
    display: flex;
    flex-direction: column;
    gap: 16px;
  }

  .param-card {
    border: 1px solid #ebeef5;
    border-radius: 12px;
    padding: 20px;
    background-color: #fff;

    &__header {
      display: flex;
      justify-content: space-between;
      align-items: center;
      margin-bottom: 16px;
      border-bottom: 1px dashed #e4e7ed;
      padding-bottom: 12px;
    }

    .param-title {
      font-weight: 600;
      font-size: 15px;

      .param-key {
        font-weight: 400;
        color: #909399;
        margin-left: 6px;
      }
    }

    .param-grid {
      display: grid;
      grid-template-columns: repeat(auto-fit, minmax(240px, 1fr));
      gap: 16px;
    }

    .form-item {
      display: flex;
      flex-direction: column;

      label {
        font-size: 13px;
        color: #606266;
        margin-bottom: 6px;
      }
    }

    &__section {
      margin-top: 18px;
      display: flex;
      flex-direction: column;
      gap: 8px;

      label {
        font-size: 13px;
        color: #606266;
      }

      .stack {
        display: flex;
        flex-direction: column;
      }

      .gap-2 {
        gap: 8px;
      }

      .gap-3 {
        gap: 16px;
      }
    }
  }

  .upload-row {
    display: flex;
    flex-wrap: wrap;
    align-items: center;
    gap: 8px;

    .tips {
      font-size: 13px;
      color: #909399;
    }
  }

  .upload-preview {
    display: flex;
    align-items: center;
    gap: 12px;

    .preview-thumb {
      width: 60px;
      height: 60px;
      border-radius: 6px;
      cursor: pointer;
      object-fit: cover;
    }

    audio,
    video {
      max-width: 260px;
      border-radius: 6px;
      background-color: #f5f7fa;
    }

    .other-file {
      display: inline-flex;
      align-items: center;
    }
  }
  @media (max-width: 768px) {
    .param-card {
      padding: 14px;
    }

    .param-card__header {
      flex-direction: column;
      align-items: flex-start;
      gap: 8px;
    }
  }
</style>
