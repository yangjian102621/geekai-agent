<template>
  <div class="login-dialog w-full pt-2" id="login-dialog">
    <custom-tabs v-model="activeName" @tab-click="handleTabClick">
      <!-- 验证码登录 -->
      <custom-tab-pane name="code">
        <template #label>
          <div class="flex items-center justify-center">
            <i class="iconfont icon-duanxin mr-2"></i>
            <span>验证码登录</span>
          </div>
        </template>
        <div class="tab-content">
          <el-form :model="codeForm" class="login-form pt-1">
            <div class="text-center text-gray-500 text-xs mb-3">
              未注册用户登录后会自动注册账号
            </div>
            <div class="form-item">
              <el-input
                placeholder="请输入手机号或邮箱"
                size="large"
                v-model="codeForm.username"
                autocomplete="off"
              >
                <template #prefix>
                  <i class="iconfont icon-mobile"></i>
                </template>
              </el-input>
            </div>

            <div class="form-item">
              <div class="code-input-row">
                <el-input
                  placeholder="请输入验证码"
                  size="large"
                  maxlength="6"
                  v-model="codeForm.code"
                  autocomplete="off"
                  @keyup.enter="submitLogin"
                >
                  <template #prefix>
                    <i class="iconfont icon-yanzm"></i>
                  </template>
                </el-input>
                <send-msg
                  size="large"
                  :receiver="codeForm.username"
                  type="auto"
                  class="send-code-btn"
                />
              </div>
            </div>

            <div class="text-left text-gray-400 text-xs">
              注册登录即代表已阅读并同意我们的
              <a href="/user-agreement" target="_blank">《用户协议》</a>
              与
              <a href="/privacy-policy" target="_blank">《隐私政策》</a
              >，未注册的手机号/邮箱将自动注册
            </div>

            <div class="form-item mt-[30px]">
              <el-button
                class="login-btn"
                type="primary"
                size="large"
                @click="submitLogin"
              >
                登录
              </el-button>
            </div>
          </el-form>
        </div>
      </custom-tab-pane>

      <!-- 微信登录 -->
      <custom-tab-pane name="wechat">
        <template #label>
          <div class="flex items-center justify-center">
            <i class="iconfont icon-wechat mr-2"></i>
            <span>微信登录</span>
          </div>
        </template>
        <div class="tab-content">
          <div class="wechat-login pt-3">
            <div class="qr-code-container">
              <div
                class="qr-code-wrapper w-[300px] h-[300px]"
                v-loading="qrcodeLoading"
              >
                <img :src="wechatLoginQRCode" class="qr-frame" />
                <!-- 二维码过期蒙版 -->
                <div v-if="qrcodeExpired" class="qr-expired-mask">
                  <div class="expired-content">
                    <i class="iconfont icon-refresh-ccw expired-icon"></i>
                    <p class="expired-text">二维码已过期</p>
                    <el-button
                      @click="getWechatLoginURL"
                      class="refresh-btn"
                      circle
                    >
                      <i class="iconfont icon-refresh text-lg"></i>
                    </el-button>
                  </div>
                </div>
              </div>
              <p class="text-center mt-4 text-gray-600">
                请使用微信扫描二维码登录
              </p>
            </div>
          </div>
        </div>
      </custom-tab-pane>

      <!-- 密码登录 -->
      <custom-tab-pane name="password">
        <template #label>
          <div class="flex items-center justify-center">
            <i class="iconfont icon-password mr-2"></i>
            <span>密码登录</span>
          </div>
        </template>
        <div class="tab-content">
          <el-form :model="passwordForm" class="login-form">
            <div class="text-left text-gray-400 text-xs my-3">
              仅适用已注册用户，未注册的用户请通过手机号/邮箱/微信登录
            </div>

            <div class="form-item">
              <el-input
                placeholder="请输入用户名"
                size="large"
                v-model="passwordForm.username"
                autocomplete="off"
              >
                <template #prefix>
                  <i class="iconfont icon-user-fill"></i>
                </template>
              </el-input>
            </div>

            <div class="form-item">
              <el-input
                placeholder="请输入密码"
                size="large"
                v-model="passwordForm.password"
                show-password
                autocomplete="off"
                @keyup.enter="submitLogin"
              >
                <template #prefix>
                  <i class="iconfont icon-password"></i>
                </template>
              </el-input>
            </div>

            <div class="text-left text-gray-400 text-xs">
              登录即代表已阅读并同意我们的
              <a href="/user-agreement" target="_blank">《用户协议》</a>
              与
              <a href="/privacy-policy" target="_blank">《隐私政策》</a>
            </div>

            <div class="form-item mt-[30px]">
              <el-button
                class="login-btn"
                type="primary"
                size="large"
                @click="submitLogin"
              >
                登录
              </el-button>
            </div>
          </el-form>
        </div>
      </custom-tab-pane>
    </custom-tabs>
  </div>
</template>

<script setup>
  import CustomTabPane from '@/components/CustomTabPane.vue'
  import CustomTabs from '@/components/CustomTabs.vue'
  import SendMsg from '@/components/SendMsg.vue'
  import { setUserToken } from '@/js/cache/session'
  import { closeLoading, showLoading } from '@/js/utils/dialog'
  import { httpGet, httpPost } from '@/js/utils/http'
  import { validateEmail, validateMobile } from '@/js/utils/validate'
  import { ElMessage } from 'element-plus'
  import QRCode from 'qrcode'
  import { onUnmounted, ref } from 'vue'

  // eslint-disable-next-line no-undef
  const props = defineProps({
    show: Boolean,
  })

  const activeName = ref('code')
  const passwordForm = ref({
    username: import.meta.env.VITE_USER || '',
    password: import.meta.env.VITE_PASS || '',
  })
  const codeForm = ref({
    username: '',
    code: '',
  })

  // eslint-disable-next-line no-undef
  const wechatLoginQRCode = ref('')
  const wechatLoginState = ref('')
  const qrcodeLoading = ref(false)
  const pollingTimer = ref(null)
  const qrcodeExpired = ref(false)
  const qrcodeTimer = ref(null)

  const handleTabClick = (tabName) => {
    if (tabName === 'wechat') {
      getWechatLoginURL()
    } else {
      // 其他登录方式，清除定时器
      if (pollingTimer.value) {
        clearInterval(pollingTimer.value)
      }
      if (qrcodeTimer.value) {
        clearTimeout(qrcodeTimer.value)
      }
    }
  }

  // 登录主要逻辑
  const submitLogin = () => {
    let username = passwordForm.value.username
    // 密码登录
    if (activeName.value === 'password') {
      if (!passwordForm.value.username) {
        return ElMessage.error('请输入用户名')
      }
      if (!passwordForm.value.password) {
        return ElMessage.error('请输入密码')
      }
      codeForm.value.code = ''
    }

    // 验证码登录
    if (activeName.value === 'code') {
      username = codeForm.value.username
      if (!validateMobile(username) && !validateEmail(username)) {
        return ElMessage.error('请输入有效的手机号或邮箱地址')
      }

      if (!codeForm.value.code) {
        return ElMessage.error('请输入验证码')
      }
      passwordForm.value.password = ''
    }

    showLoading('登录中...', '#login-dialog')
    httpPost('/api/user/login', {
      username: username,
      password: passwordForm.value.password,
      code: codeForm.value.code,
      method: activeName.value,
    })
      .then((res) => {
        setUserToken(res.data.token)
        location.reload()
      })
      .catch((e) => {
        ElMessage.error('登录失败，' + e.message)
        closeLoading()
      })
  }

  // 获取微信登录 URL
  const getWechatLoginURL = () => {
    qrcodeLoading.value = true
    qrcodeExpired.value = false

    // 清除可能存在的旧定时器
    if (qrcodeTimer.value) {
      clearTimeout(qrcodeTimer.value)
    }

    httpGet('/api/user/login/qrcode')
      .then((res) => {
        // 生成二维码
        QRCode.toDataURL(
          res.data.url,
          { width: 300, height: 300, margin: 2 },
          (error, url) => {
            if (error) {
              console.error(error)
            } else {
              wechatLoginQRCode.value = url
            }
          }
        )
        wechatLoginState.value = res.data.state
        // 开始轮询状态
        startPolling()

        // 设置10分钟后二维码过期
        qrcodeTimer.value = setTimeout(() => {
          qrcodeExpired.value = true
          // 停止轮询
          if (pollingTimer.value) {
            clearInterval(pollingTimer.value)
          }
        }, 60 * 1000) // 1分钟过期
      })
      .catch((e) => {
        ElMessage.error('获取微信登录 URL 失败，' + e.message)
      })
      .finally(() => {
        qrcodeLoading.value = false
      })
  }

  // 开始轮询
  const startPolling = () => {
    // 清除可能存在的旧定时器
    if (pollingTimer.value) {
      clearInterval(pollingTimer.value)
    }

    pollingTimer.value = setInterval(() => {
      checkLoginStatus()
    }, 1000) // 每1秒检查一次
  }

  // 检查登录状态
  const checkLoginStatus = () => {
    if (!wechatLoginState.value) return

    httpGet(`/api/user/login/status?state=${wechatLoginState.value}`)
      .then((res) => {
        const status = res.data.status

        switch (status) {
          case 'success':
            // 登录成功
            clearInterval(pollingTimer.value)
            clearTimeout(qrcodeTimer.value)
            setUserToken(res.data.token)
            location.reload()
            break

          case 'expired':
            // 二维码过期
            clearInterval(pollingTimer.value)
            clearTimeout(qrcodeTimer.value)
            qrcodeExpired.value = true
            break

          case 'pending':
            // 继续轮询
            break

          default:
            // 其他错误情况
            clearInterval(pollingTimer.value)
            clearTimeout(qrcodeTimer.value)
            ElMessage.error('登录失败，请重试')
            break
        }
      })
      .catch((e) => {
        // 发生错误时显示过期状态
        clearInterval(pollingTimer.value)
        clearTimeout(qrcodeTimer.value)
        qrcodeExpired.value = true
      })
  }

  // 组件卸载时清除定时器
  onUnmounted(() => {
    if (pollingTimer.value) {
      clearInterval(pollingTimer.value)
    }
    if (qrcodeTimer.value) {
      clearTimeout(qrcodeTimer.value)
    }
  })
</script>

<style lang="scss" scoped>
  .login-dialog {
    max-width: 480px;
    margin: 0 auto;

    .login-container {
      background: white;
      border-radius: 12px;
      padding: 32px;
      box-shadow: 0 4px 20px rgba(0, 0, 0, 0.1);
    }

    .login-header {
      margin-bottom: 24px;

      h2 {
        font-size: 24px;
        font-weight: 600;
        color: #2c3e50;
        margin: 0;
      }
    }

    .login-form {
      .form-item {
        margin-bottom: 20px;

        .login-btn {
          width: 100%;
          height: 48px;
          font-size: 16px;
          font-weight: 500;
          border-radius: 8px;
        }
      }

      .code-input-row {
        display: flex;
        gap: 12px;

        .el-input {
          flex: 1;
        }

        .send-code-btn {
          flex-shrink: 0;
        }
      }
    }

    .wechat-login {
      display: flex;
      justify-content: center;
      align-items: center;
      min-height: 240px;

      .qr-code-container {
        text-align: center;

        .qr-code-wrapper {
          display: inline-block;
          border: 1px solid #e1e5e9;
          border-radius: 8px;
          overflow: hidden;
          position: relative;

          .qr-frame {
            display: block;
          }

          .qr-expired-mask {
            position: absolute;
            top: 0;
            left: 0;
            right: 0;
            bottom: 0;
            background: rgba(0, 0, 0, 0.7);
            display: flex;
            align-items: center;
            justify-content: center;
            border-radius: 8px;

            .expired-content {
              text-align: center;
              color: white;

              .expired-icon {
                font-size: 48px;
                color: #f56565;
                margin-bottom: 12px;
                display: block;
              }

              .expired-text {
                font-size: 16px;
                margin: 0 0 16px 0;
                font-weight: 500;
              }

              .refresh-btn {
                border-radius: 6px;
                padding: 8px 16px;

                .iconfont {
                  font-size: 14px;
                }
              }
            }
          }
        }
      }
    }

    // Element Plus 组件样式优化
    :deep(.el-input) {
      .el-input__wrapper {
        border-radius: 8px;
        box-shadow: 0 0 0 1px #dcdfe6;
        transition: all 0.3s;

        &:hover {
          box-shadow: 0 0 0 1px #c0c4cc;
        }

        &.is-focus {
          box-shadow: 0 0 0 1px var(--el-color-primary);
        }
      }
    }

    :deep(.el-tabs) {
      .el-tabs__nav {
        width: 100%;

        .el-tabs__item {
          flex: 1;
          text-align: center;
          padding: 16px 20px;
          font-size: 14px;
          font-weight: 500;

          .iconfont {
            font-size: 16px;
          }
        }
      }
    }
  }

  // 响应式设计
  @media (max-width: 576px) {
    .login-dialog {
      .login-container {
        padding: 24px 20px;
      }

      .tab-content {
        padding: 16px 0;
        min-height: 240px;
      }

      .wechat-login {
        .qr-code-wrapper {
          // .qr-frame {
          //   width: 240px !important;
          //   height: 240px !important;
          // }

          .qr-expired-mask {
            .expired-content {
              .expired-icon {
                font-size: 36px;
                margin-bottom: 8px;
              }

              .expired-text {
                font-size: 14px;
                margin: 0 0 12px 0;
              }

              .refresh-btn {
                padding: 6px 12px;
                font-size: 12px;
              }
            }
          }
        }
      }
    }
  }
</style>
