<template>
  <div class="settings-container">
    <custom-tabs v-model="activeName" @tab-click="handleTabClick">
      <custom-tab-pane name="profile">
        <template #label>
          <div class="flex items-center justify-center">
            <i class="iconfont icon-user-fill mr-1"></i>
            <span>账户设置</span>
          </div>
        </template>
        <div class="tab-content">
          <div class="user-info">
            <div class="mb-3 row text-center p-3 rounded-4">
              <el-upload
                :auto-upload="true"
                :show-file-list="false"
                :http-request="handleUpload"
                accept=".png,.jpg,.jpeg,.bmp"
              >
                <el-avatar :src="user.avatar" shape="circle" :size="100" />
              </el-upload>
            </div>
            <div class="mb-3 row">
              <label class="col-sm-3 col-form-label">用户名</label>
              <div class="col-sm-9 d-flex align-items-center">
                <input
                  type="text"
                  readonly
                  class="form-control-plaintext"
                  :value="username"
                />
                <el-button
                  class="ms-2"
                  size="small"
                  @click="openChangeUsernameDialog"
                >
                  更改
                </el-button>
              </div>
            </div>

            <div class="mb-3 row">
              <label class="col-sm-3 col-form-label">剩余积分</label>
              <div class="col-sm-9">
                <el-tag size="large">{{ user['scores'] }}</el-tag>
              </div>
            </div>

            <div class="mb-3 row" v-if="user['expired_time'] > 0">
              <label class="col-sm-3 col-form-label">会员有效期</label>
              <div class="col-sm-9">
                <el-tag size="large">{{
                  dateFormat(user['expired_time'], 'yyyy-MM-dd')
                }}</el-tag>
              </div>
            </div>

            <div class="mb-3 row">
              <label class="col-sm-3 col-form-label">昵称</label>
              <div class="col-sm-9">
                <input
                  type="text"
                  v-model="user.nickname"
                  class="form-control"
                />
              </div>
            </div>

            <div class="mt-4 text-center">
              <button class="btn btn-primary" @click="save">保存</button>
            </div>
          </div>
        </div>
      </custom-tab-pane>

      <custom-tab-pane name="password">
        <template #label>
          <div class="flex items-center justify-center">
            <i class="iconfont icon-password mr-1"></i>
            <span>修改密码</span>
          </div>
        </template>
        <div class="tab-content">
          <div class="settings-list">
            <div class="setting-item p-3">
              <form @submit.prevent="handleChangePassword">
                <div class="mb-3">
                  <label class="form-label">新密码</label>
                  <input
                    type="password"
                    class="form-control"
                    v-model="passwordForm.newPassword"
                    placeholder="请输入新密码"
                    :class="{ 'is-invalid': passwordErrors.newPassword }"
                  />
                  <div
                    class="invalid-feedback"
                    v-if="passwordErrors.newPassword"
                  >
                    {{ passwordErrors.newPassword }}
                  </div>
                </div>
                <div class="mb-3">
                  <label class="form-label">确认密码</label>
                  <input
                    type="password"
                    class="form-control"
                    v-model="passwordForm.confirmPassword"
                    placeholder="请确认新密码"
                    :class="{ 'is-invalid': passwordErrors.confirmPassword }"
                  />
                  <div
                    class="invalid-feedback"
                    v-if="passwordErrors.confirmPassword"
                  >
                    {{ passwordErrors.confirmPassword }}
                  </div>
                </div>
                <div class="d-grid">
                  <button type="submit" class="btn btn-primary">
                    确认修改
                  </button>
                </div>
              </form>
            </div>
          </div>
        </div>
      </custom-tab-pane>
    </custom-tabs>

    <!-- Bootstrap Modal -->
    <div
      class="modal fade"
      :class="{ show: showChangeUsernameDialog }"
      :style="{ display: showChangeUsernameDialog ? 'block' : 'none' }"
      tabindex="-1"
      @click.self="showChangeUsernameDialog = false"
    >
      <div class="modal-dialog modal-dialog-centered">
        <div class="modal-content">
          <div class="modal-header">
            <h5 class="modal-title">修改用户名</h5>
            <button
              type="button"
              class="btn-close"
              @click="showChangeUsernameDialog = false"
            ></button>
          </div>
          <div class="modal-body">
            <form @submit.prevent="submitUsernameChange">
              <div class="mb-3">
                <label class="form-label">手机号或邮箱</label>
                <input
                  type="text"
                  class="form-control"
                  v-model="userForm.username"
                  placeholder="请输入手机号或邮箱"
                  :class="{ 'is-invalid': usernameErrors.username }"
                />
                <div class="invalid-feedback" v-if="usernameErrors.username">
                  {{ usernameErrors.username }}
                </div>
              </div>
              <div class="mb-3">
                <label class="form-label">验证码</label>
                <div class="d-flex">
                  <input
                    type="text"
                    class="form-control me-2"
                    v-model="userForm.code"
                    maxlength="6"
                    placeholder="请输入验证码"
                    :class="{ 'is-invalid': usernameErrors.code }"
                  />
                  <send-msg :receiver="userForm.username" size="small" />
                </div>
                <div class="invalid-feedback" v-if="usernameErrors.code">
                  {{ usernameErrors.code }}
                </div>
              </div>
            </form>
          </div>
          <div class="modal-footer">
            <button
              type="button"
              class="btn btn-secondary"
              @click="showChangeUsernameDialog = false"
            >
              取消
            </button>
            <button
              type="button"
              class="btn btn-primary"
              @click="submitUsernameChange"
            >
              确认
            </button>
          </div>
        </div>
      </div>
    </div>
    <!-- Bootstrap Modal Backdrop -->
    <div
      v-if="showChangeUsernameDialog"
      class="modal-backdrop fade show"
      @click="showChangeUsernameDialog = false"
    ></div>
  </div>
</template>

<script setup>
  import CustomTabPane from '@/components/CustomTabPane.vue'
  import CustomTabs from '@/components/CustomTabs.vue'
  import SendMsg from '@/components/SendMsg.vue'
  import { checkSession } from '@/js/cache/session'
  import { useSharedStore } from '@/js/cache/sharedata'
  import { httpGet, httpPost } from '@/js/utils/http'
  import { validateEmail, validateMobile } from '@/js/utils/validate'
  import Compressor from 'compressorjs'
  import { ElMessage } from 'element-plus'
  import { onMounted, ref } from 'vue'

  const user = ref({
    vip: false,
    nickname: '',
    avatar: '',
    scores: 0,
    expired_time: 0,
  })
  const username = ref('')

  const props = defineProps({
    activeName: {
      type: String,
      default: 'profile',
    },
  })

  const activeName = ref(props.activeName)
  const emit = defineEmits(['update-avatar', 'update-user'])
  const sharedata = useSharedStore()

  const settings = ref(sharedata.chatSetting)
  const showChangeUsernameDialog = ref(false)
  const userForm = ref({
    username: '',
    code: '',
  })
  const passwordForm = ref({
    newPassword: '',
    confirmPassword: '',
  })

  const passwordErrors = ref({
    newPassword: '',
    confirmPassword: '',
  })

  const usernameErrors = ref({
    username: '',
    code: '',
  })

  // 验证密码表单
  const validatePasswordForm = () => {
    passwordErrors.value = {
      newPassword: '',
      confirmPassword: '',
    }

    let isValid = true

    // 验证新密码
    if (!passwordForm.value.newPassword) {
      passwordErrors.value.newPassword = '请输入新密码'
      isValid = false
    } else if (passwordForm.value.newPassword.length < 6) {
      passwordErrors.value.newPassword = '密码长度不能少于6位'
      isValid = false
    }

    // 验证确认密码
    if (!passwordForm.value.confirmPassword) {
      passwordErrors.value.confirmPassword = '请确认新密码'
      isValid = false
    } else if (
      passwordForm.value.confirmPassword !== passwordForm.value.newPassword
    ) {
      passwordErrors.value.confirmPassword = '两次输入的密码不一致'
      isValid = false
    }

    return isValid
  }

  // 验证用户名表单
  const validateUsernameForm = () => {
    usernameErrors.value = {
      username: '',
      code: '',
    }

    let isValid = true

    // 验证用户名
    if (!userForm.value.username) {
      usernameErrors.value.username = '请输入用户名'
      isValid = false
    } else if (
      !validateEmail(userForm.value.username) &&
      !validateMobile(userForm.value.username)
    ) {
      usernameErrors.value.username = '请输入有效的手机号或邮箱地址'
      isValid = false
    }

    // 验证验证码
    if (!userForm.value.code) {
      usernameErrors.value.code = '请输入验证码'
      isValid = false
    } else if (userForm.value.code.length !== 6) {
      usernameErrors.value.code = '验证码长度应为6位'
      isValid = false
    }

    return isValid
  }

  const handleTabClick = (tabName) => {
    console.log('Tab clicked:', tabName)
  }

  onMounted(async () => {
    const _user = await checkSession()
    if (_user) {
      const res = await httpGet('/api/user/profile')
      user.value = res.data
      username.value = res.data.username
    }
  })

  const handleUpload = (file) => {
    // 压缩图片并上传
    new Compressor(file.file, {
      quality: 0.6,
      success(result) {
        const formData = new FormData()
        formData.append('file', result, result.name)
        // 执行上传操作
        httpPost('/api/file/upload', formData)
          .then((res) => {
            user.value.avatar = res.data.url
            emit('update-avatar', res.data.url)
            httpPost('/api/user/update/profile', user.value)
              .then(() => {
                ElMessage.success({ message: '更新头像成功', duration: 500 })
              })
              .catch((e) => {
                ElMessage.error('更新头像失败：' + e.message)
              })
          })
          .catch((e) => {
            ElMessage.error('图片上传失败:' + e.message)
          })
      },
      error(err) {
        console.log(err.message)
      },
    })
  }

  const save = () => {
    httpPost('/api/user/update/profile', user.value)
      .then(() => {
        ElMessage.success({ message: '更新成功', duration: 500 })
        emit('update-user', user.value)
      })
      .catch((e) => {
        ElMessage.error('更新失败：' + e.message)
      })
  }

  const changeTheme = () => {
    sharedata.setChatSetting({
      ...settings.value,
      darkTheme: settings.value.darkTheme,
    })
  }

  const changeChatContext = () => {
    sharedata.setChatSetting({
      ...settings.value,
      chatContext: settings.value.chatContext,
    })
  }

  const submitUsernameChange = async () => {
    // 验证表单
    if (!validateUsernameForm()) {
      return
    }

    try {
      await httpPost('/api/user/update/username', userForm.value)
      ElMessage.success('用户名修改成功')
      showChangeUsernameDialog.value = false
      // 重置表单和错误信息
      userForm.value = {
        username: '',
        code: '',
      }
      usernameErrors.value = {
        username: '',
        code: '',
      }
      // 刷新用户信息
      const res = await httpGet('/api/user/profile')
      user.value = res.data
      username.value = res.data.username
    } catch (error) {
      if (error.message) {
        ElMessage.error(error.message)
      }
    }
  }

  // 在用户名输入框旁边添加修改按钮
  const openChangeUsernameDialog = () => {
    showChangeUsernameDialog.value = true
    userForm.value = {
      username: '',
      code: '',
    }
    usernameErrors.value = {
      username: '',
      code: '',
    }
  }

  // 处理密码修改
  const handleChangePassword = async () => {
    // 验证表单
    if (!validatePasswordForm()) {
      return
    }

    try {
      await httpPost('/api/user/update/password', {
        password: passwordForm.value.newPassword,
      })
      ElMessage.success('密码修改成功')
      // 重置表单
      passwordForm.value = {
        newPassword: '',
        confirmPassword: '',
      }
      passwordErrors.value = {
        newPassword: '',
        confirmPassword: '',
      }
    } catch (error) {
      if (error.message) {
        ElMessage.error(error.message)
      }
    }
  }
</script>

<style lang="scss">
  .settings-container {
    padding: 20px;

    .tab-content {
      background-color: white;
      border-radius: 8px;

      .user-info {
        .el-row {
          justify-content: center;
          margin-bottom: 10px;
        }

        .vip-icon {
          position: relative;
          top: 5px;
        }
        // .col-form-label {
        //   font-size: 16px;
        //   color: #606266;
        // }
      }

      .settings-list {
        .setting-item {
          border-bottom: 1px solid #ebeef5;

          &:last-child {
            border-bottom: none;
          }

          // .setting-label {
          //   font-size: 16px;
          //   color: #606266;
          // }

          .setting-control {
            .el-switch {
              --el-switch-on-color: var(--el-color-primary);
            }
          }
        }
      }
    }
  }

  .dialog-footer {
    display: flex;
    justify-content: flex-end;
    gap: 10px;
  }

  // Bootstrap Modal 样式修正
  .modal {
    z-index: 1050;

    &.show {
      display: block !important;
    }
  }

  .modal-backdrop {
    z-index: 1040;
  }
</style>
