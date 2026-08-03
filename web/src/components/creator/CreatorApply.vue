<template>
  <div class="creator-apply container">
    <!-- 已申请状态显示 -->
    <div
      v-if="creatorStatus && creatorStatus.status !== 'not_applied'"
      class="alert my-4"
      :class="getStatusClass()"
      role="alert"
    >
      <div class="d-flex align-items-center">
        <i class="iconfont me-2 text-2xl" :class="getStatusIcon()"></i>
        <div>
          {{ getStatusMessage() }}
        </div>
      </div>
    </div>

    <!-- 申请表单 -->
    <div v-else>
      <div class="alert alert-primary my-4" role="alert">
        成为创作者后，您可以上传自己的智能体，并通过平台实现变现。
      </div>

      <form @submit.prevent="submit" novalidate>
        <div class="mb-4">
          <label for="creatorName" class="form-label fw-bold">创作者名称</label>
          <div class="input-group">
            <input
              id="creatorName"
              type="text"
              class="form-control"
              v-model="form.name"
              placeholder="请输入创作者名称"
              :class="{ 'is-invalid': errors.name }"
              maxlength="20"
              autocomplete="off"
              @blur="validate"
            />
            <el-tooltip content="随机生成一个名称" placement="top">
              <button
                class="btn btn-outline-primary"
                type="button"
                @click="randomName"
                :disabled="loadingName"
              >
                <span
                  v-if="loadingName"
                  class="spinner-border spinner-border-sm"
                ></span>
                <i v-else class="iconfont icon-chuangzuo"></i>
              </button>
            </el-tooltip>
          </div>
          <div class="invalid-feedback d-block" v-if="errors.name">
            {{ errors.name }}
          </div>
        </div>

        <div class="mb-4">
          <label for="creatorUsername" class="form-label fw-bold">
            用户名
            <el-tooltip
              content="用户名是创作者的唯一标识，<br />跟创作者首页地址关联，<br />建议使用好记的英文名"
              raw-content
              placement="right"
            >
              <i class="iconfont icon-info"></i>
            </el-tooltip>
          </label>
          <div class="input-group">
            <input
              id="creatorUsername"
              type="text"
              class="form-control"
              v-model="form.username"
              placeholder="请输入用户名, 只支持英文和数字"
              :class="{ 'is-invalid': errors.username }"
              maxlength="30"
              autocomplete="off"
              @blur="validate"
            />
            <el-tooltip content="随机生成一个用户名" placement="top">
              <button
                class="btn btn-outline-primary"
                type="button"
                @click="randomUsername"
              >
                <i class="iconfont icon-chuangzuo"></i>
              </button>
            </el-tooltip>
          </div>
          <div class="invalid-feedback d-block" v-if="errors.username">
            {{ errors.username }}
          </div>
        </div>

        <div class="mb-4">
          <label for="creatorDescription" class="form-label fw-bold"
            >创作者简介</label
          >
          <textarea
            id="creatorDescription"
            class="form-control"
            v-model="form.description"
            placeholder="请简单介绍一下您的创作背景和专长领域..."
            :class="{ 'is-invalid': errors.description }"
            rows="4"
            maxlength="200"
            @blur="validate"
          ></textarea>
          <div class="form-text">
            <span class="text-muted">{{ form.description.length }}/200</span>
          </div>
          <div class="invalid-feedback d-block" v-if="errors.description">
            {{ errors.description }}
          </div>
        </div>

        <div class="mb-4">
          <label class="form-label fw-bold">Logo</label>
          <div class="d-flex align-items-center gap-3">
            <div class="avatar-uploader position-relative">
              <el-upload
                :auto-upload="true"
                :show-file-list="false"
                :http-request="handleUpload"
                accept=".png,.jpg,.jpeg,.bmp"
                :disabled="uploadingLogo"
              >
                <el-avatar
                  :src="form.logo"
                  shape="circle"
                  :size="60"
                  v-if="form.logo"
                />
                <button
                  type="button"
                  class="w-[60px] h-[60px] flex items-center justify-center rounded-lg border border-dashed border-gray-300"
                  v-else
                  :disabled="uploadingLogo"
                >
                  <span
                    v-if="uploadingLogo"
                    class="spinner-border spinner-border-sm"
                  ></span>
                  <i v-else class="iconfont icon-plus text-xl"></i>
                </button>
              </el-upload>
            </div>
            <el-tooltip content="随机生成一个Logo" placement="top">
              <button
                type="button"
                class="w-8 h-8 flex items-center justify-center rounded-full text-white shadow bg-gradient-to-br from-blue-500 via-purple-500 to-pink-500 hover:from-blue-600 hover:via-purple-600 hover:to-pink-600 transition"
                @click="randomLogo"
                :disabled="loadingLogo"
              >
                <span
                  v-if="loadingLogo"
                  class="spinner-border spinner-border-sm"
                  style="font-size: 12px"
                ></span>
                <i v-else class="iconfont icon-chuangzuo text-xl"></i>
              </button>
            </el-tooltip>
          </div>
          <div class="invalid-feedback d-block" v-if="errors.logo">
            {{ errors.logo }}
          </div>
        </div>

        <div class="mt-2 text-center">
          <button type="submit" class="btn btn-primary" :disabled="loading">
            <span
              v-if="loading"
              class="spinner-border spinner-border-sm me-2"
            ></span>
            提交申请
          </button>
        </div>
      </form>
    </div>
  </div>
</template>

<script setup>
  import {
    closeLoading,
    showLoading,
    showMessageError,
    showMessageOK,
  } from '@/js/utils/dialog.js'
  import { httpGet, httpPost } from '@/js/utils/http'
  import { randString } from '@/js/utils/libs'
  import Compressor from 'compressorjs'
  import { onMounted, reactive, ref } from 'vue'

  const emits = defineEmits(['success'])

  const form = reactive({
    name: '',
    description: '',
    logo: '',
    username: '',
  })
  const errors = reactive({
    name: '',
    description: '',
    logo: '',
    username: '',
  })
  const loading = ref(false)
  const loadingName = ref(false)
  const loadingLogo = ref(false)
  const uploadingLogo = ref(false)
  const creatorStatus = ref(null)

  onMounted(() => {
    checkCreatorStatus()
  })

  // 检查创作者申请状态
  const checkCreatorStatus = async () => {
    try {
      const res = await httpGet('/api/creator/status')
      creatorStatus.value = res.data
    } catch (e) {
      console.error('检查创作者状态失败:', e)
    }
  }

  // 获取状态样式类
  const getStatusClass = () => {
    switch (creatorStatus.value?.status) {
      case 'pending':
        return 'alert-primary'
      case 'approved':
        return 'alert-success'
      case 'rejected':
        return 'alert-danger'
      default:
        return 'alert-info'
    }
  }

  // 获取状态图标
  const getStatusIcon = () => {
    switch (creatorStatus.value?.status) {
      case 'pending':
        return 'icon-clock'
      case 'approved':
        return 'icon-check'
      case 'rejected':
        return 'icon-refuse'
      default:
        return 'icon-info'
    }
  }

  // 获取状态标题
  const getStatusMessage = () => {
    switch (creatorStatus.value?.status) {
      case 'pending':
        return '申请审核中，请耐心等待...'
      case 'approved':
        return '您已成为创作者,请进入创作者中心上传智能体。'
      case 'rejected':
        return (
          '申请被拒绝：' + creatorStatus.value?.message + '，请联系管理员。'
        )
      default:
        return '申请状态'
    }
  }

  // 随机生成名称
  const randomName = async () => {
    loadingName.value = true
    try {
      const res = await httpGet('/api/creator/rand/name')
      form.name = res.data.name
      errors.name = ''
    } catch (e) {
      showMessageError('生成名称失败: ' + e.message)
    } finally {
      loadingName.value = false
    }
  }

  // 随机生成用户名

  const randomUsername = async () => {
    form.username = randString(12)
  }

  // 随机生成Logo
  const randomLogo = async () => {
    loadingLogo.value = true
    try {
      const res = await httpGet('/api/creator/rand/logo')
      form.logo = res.data.logo
      errors.logo = ''
    } catch (e) {
      showMessageError('生成Logo失败: ' + e.message)
    } finally {
      loadingLogo.value = false
    }
  }

  // 上传Logo
  const handleUpload = (file) => {
    uploadingLogo.value = true
    // 压缩图片并上传
    new Compressor(file.file, {
      quality: 0.6,
      success(result) {
        const formData = new FormData()
        formData.append('file', result, result.name)
        showLoading('上传中...')
        // 执行上传操作
        httpPost('/api/file/upload', formData)
          .then((res) => {
            form.logo = res.data.url
            errors.logo = ''
          })
          .catch((e) => {
            showMessageError('图片上传失败:' + e.message)
          })
          .finally(() => {
            closeLoading()
            uploadingLogo.value = false
          })
      },
      error(err) {
        console.log(err.message)
        uploadingLogo.value = false
      },
    })
  }

  const validate = () => {
    let valid = true
    errors.name = ''
    errors.description = ''
    errors.logo = ''
    errors.username = ''

    if (!form.name.trim()) {
      errors.name = '请输入创作者名称'
      valid = false
    } else if (form.name.trim().length < 2) {
      errors.name = '创作者名称至少2个字符'
      valid = false
    }

    if (!form.username.trim()) {
      errors.username = '请输入用户名'
      valid = false
    } else if (form.username.trim().length < 5) {
      errors.username = '用户名至少5个字符'
      valid = false
    } else if (!/^[a-zA-Z0-9]+$/.test(form.username.trim())) {
      errors.username = '用户名只支持英文和数字'
      valid = false
    }

    if (!form.description.trim()) {
      errors.description = '请输入创作者简介'
      valid = false
    } else if (form.description.trim().length < 10) {
      errors.description = '创作者简介至少10个字符'
      valid = false
    }

    if (!form.logo) {
      errors.logo = '请上传Logo'
      valid = false
    }

    return valid
  }

  const submit = async () => {
    if (!validate()) return

    loading.value = true
    try {
      await httpPost('/api/creator/apply', {
        name: form.name.trim(),
        description: form.description.trim(),
        logo: form.logo,
        username: form.username.trim(),
      })
      showMessageOK('申请已提交，等待审核')
      // 重新检查状态
      await checkCreatorStatus()
      // emits('success')
    } catch (e) {
      showMessageError(e.message || '申请失败')
    } finally {
      loading.value = false
    }
  }
</script>

<style scoped>
  .creator-apply {
    background: #fff;
    border-radius: 18px;
  }

  .avatar-uploader-icon {
    font-size: 28px;
    color: #999;
    border: 1px dashed #d9d9d9;
    background: #f8f9fa;
  }

  .alert {
    border-radius: 12px;
  }

  .spinner-border-sm {
    width: 1rem;
    height: 1rem;
  }

  textarea.form-control {
    resize: vertical;
    min-height: 100px;
  }
</style>
