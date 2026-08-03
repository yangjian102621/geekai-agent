<template>
  <div class="card">
    <div class="card-header bg-white text-center p-3">
      <h5 class="card-title mb-0">GeekAI 增值服务</h5>
    </div>
    <div class="card-body">
      <div class="container">
        <el-tabs v-model="activeTab" type="border-card">
          <el-tab-pane name="license">
            <template #label>
              <div class="d-flex align-items-center">
                <i class="iconfont icon-license"></i>
                <span class="ms-2">许可证授权</span>
              </div>
            </template>

            <h4 class="mt-2">License 授权说明：</h4>
            <ol class="active-info">
              <li>1、授权后可以解锁 Geek-Agent 所有高级功能。</li>
              <li>2、授权后可以去品牌化，也就是可以更换 Logo 和 版权信息。</li>
              <li>3、系统部署后会自动申请一个免费授权，时间为一个月。</li>
            </ol>

            <el-descriptions
              v-if="license.is_active"
              class="mt-2"
              title="已授权信息"
              :column="1"
              border
            >
              <el-descriptions-item>
                <template #label>
                  <div class="cell-item">授权对象</div>
                </template>
                {{ license.name }}
              </el-descriptions-item>
              <el-descriptions-item>
                <template #label>
                  <div class="cell-item">License Key</div>
                </template>
                <div class="d-flex align-items-center">
                  {{
                    showFullLicense
                      ? license.license
                      : maskString(license.license)
                  }}
                  <el-tooltip content="显示/隐藏完整License" placement="top">
                    <i
                      class="ms-2 iconfont"
                      :class="
                        showFullLicense ? 'icon-eye-close' : 'icon-eye-open'
                      "
                      style="cursor: pointer"
                      @click="showFullLicense = !showFullLicense"
                    ></i>
                  </el-tooltip>
                </div>
              </el-descriptions-item>
              <el-descriptions-item>
                <template #label>
                  <div class="cell-item">机器码</div>
                </template>
                {{ license.mid }}
              </el-descriptions-item>
              <el-descriptions-item>
                <template #label>
                  <div class="cell-item">授权类型</div>
                </template>
                <el-tag type="primary" v-if="license.type === 'free'">
                  {{ licenseType[license.type] }}
                </el-tag>
                <el-tag type="success" v-else>
                  {{ licenseType[license.type] }}
                </el-tag>
              </el-descriptions-item>
              <el-descriptions-item>
                <template #label>
                  <div class="cell-item">到期时间</div>
                </template>
                <span class="me-2">{{ dateFormat(license.expired_at) }}</span>
              </el-descriptions-item>
            </el-descriptions>

            <form class="form" v-else>
              <div class="mb-3">
                <label class="form-label"
                  >当前系统未授权或者授权已过期，请输入许可授权码激活：</label
                >
                <PasswordInput v-model="licenseKey" />
              </div>

              <div class="d-flex justify-content-center mt-4">
                <button type="button" class="btn btn-primary" @click="active">
                  立即激活
                </button>
              </div>
            </form>
          </el-tab-pane>
          <el-tab-pane name="geekservice">
            <template #label>
              <div class="d-flex align-items-center">
                <i class="iconfont icon-yanzm"></i>
                <span class="ms-2">行为验证码</span>
              </div>
            </template>
            <div class="alert alert-primary mt-2 mb-4" role="alert">
              行为验证码服务，开启后用户登录的时候需要进行行为验证，可以有效防止恶意登录。<br />
              请联系作者免费领取令牌，填入下面输入框开通验证服务。
            </div>

            <form class="form">
              <div class="mb-3">
                <label class="form-label">服务令牌</label>
                <PasswordInput v-model="captchaConfig.api_key" />
              </div>

              <div class="mb-3">
                <label class="form-label">验证码类型</label>
                <el-radio-group v-model="captchaConfig.type" class="w-100">
                  <el-radio value="dot" border>点选验证码</el-radio>
                  <el-radio value="slide" border>滑动验证码</el-radio>
                </el-radio-group>
              </div>

              <div class="mb-3">
                <label class="form-label">启用验证码</label>
                <div>
                  <el-switch size="large" v-model="captchaConfig.enabled" />
                </div>
              </div>

              <div class="d-flex justify-content-center mt-4">
                <button
                  type="button"
                  class="btn btn-primary"
                  @click="saveCaptchaConfig"
                >
                  提交保存
                </button>
              </div>
            </form>
          </el-tab-pane>
          <el-tab-pane name="wechat">
            <template #label>
              <div class="d-flex align-items-center">
                <i class="iconfont icon-wechat"></i>
                <span class="ms-2">微信登录</span>
              </div>
            </template>
            <div class="alert alert-primary mt-2 mb-4" role="alert">
              微信登录服务，开启后用户可以使用微信扫码登录。<br />
              请联系作者免费领取令牌，填入下面输入框开通验证服务。
            </div>

            <form class="form">
              <div class="mb-3">
                <label class="form-label">服务令牌</label>
                <PasswordInput v-model="wechatConfig.api_key" />
              </div>

              <div class="mb-3">
                <label class="form-label">登录成功回调URL</label>
                <input
                  type="text"
                  class="form-control"
                  v-model="wechatConfig.notify_url"
                />
              </div>

              <div class="mb-3">
                <label class="form-label">启用微信登录</label>
                <div>
                  <el-switch size="large" v-model="wechatConfig.enabled" />
                </div>
              </div>

              <div class="d-flex justify-content-center mt-4">
                <button
                  type="button"
                  class="btn btn-primary"
                  @click="saveWechatConfig"
                >
                  提交保存
                </button>
              </div>
            </form>
          </el-tab-pane>
        </el-tabs>
      </div>
    </div>
  </div>
</template>

<script setup>
  import {
    closeLoading,
    showConfirm,
    showLoading,
    showMessageError,
    showMessageOK,
  } from '@/js/utils/dialog.js'
  import { httpGet, httpPost } from '@/js/utils/http.js'
  import { dateFormat, maskString } from '@/js/utils/libs.js'
  import { ElTabPane, ElTabs } from 'element-plus'
  import { onMounted, ref } from 'vue'
  import PasswordInput from '@/components/PasswordInput.vue'

  const activeTab = ref('license')
  const license = ref({ configs: {}, expired_at: 0, is_active: false })
  const licenseKey = ref('')
  const licenseType = {
    free: '免费版',
    pro: '旗舰版',
  }
  const showFullLicense = ref(false)

  // 行为验证码配置
  const captchaConfig = ref({
    api_key: '',
    type: 'slide',
    enabled: false,
  })

  // 微信登录配置
  const wechatConfig = ref({
    app_id: '',
    app_secret: '',
    enabled: false,
  })

  onMounted(() => {
    fetchLicense()
    fetchCaptchaConfig()
    fetchWechatConfig()
  })

  const fetchLicense = () => {
    httpGet('/api/admin/config/license/get')
      .then((res) => {
        license.value = res.data
      })
      .catch((e) => {
        showMessageError('许可证信息获取失败: ' + e.message)
      })
  }

  const fetchCaptchaConfig = async () => {
    try {
      const res = await httpGet('/api/admin/config/get?name=captcha')
      captchaConfig.value = res.data
    } catch (e) {
      console.warn(e)
    }
  }

  const fetchWechatConfig = async () => {
    try {
      const res = await httpGet('/api/admin/config/get?name=wx_login')
      wechatConfig.value = res.data
    } catch (e) {
      console.warn(e)
    } finally {
      wechatConfig.value.notify_url =
        wechatConfig.value.notify_url ||
        window.location.origin + '/api/user/login/callback'
    }
  }

  // 保存行为验证码配置
  const saveCaptchaConfig = () => {
    showLoading()
    httpPost('/api/admin/config/update/captcha', captchaConfig.value)
      .then(() => {
        showMessageOK('操作成功！')
        closeLoading()
      })
      .catch((e) => {
        showMessageError('操作失败：' + e.message)
        closeLoading()
      })
  }

  // 激活 License
  const active = () => {
    if (licenseKey.value === '') {
      return showMessageError('请输入授权码')
    }
    httpPost('/api/admin/config/license/active', { license: licenseKey.value })
      .then((res) => {
        showMessageOK('授权成功，机器编码为：' + res.data)
        licenseKey.value = ''
        fetchLicense()
      })
      .catch((e) => {
        showMessageError(e.message)
      })
  }

  // 申请免费License
  const applyFreeLicense = () => {
    showConfirm(
      '提示',
      '每账号只能申请一次免费License，请确认是否申请？',
      () => {
        httpPost('/api/admin/config/license/apply')
          .then(() => {
            showMessageOK('申请成功!')
            fetchLicense()
          })
          .catch((e) => {
            showMessageError(e.message)
          })
      }
    )
  }

  // 保存微信登录配置
  const saveWechatConfig = () => {
    showLoading()
    httpPost('/api/admin/config/update/wechat', wechatConfig.value)
      .then(() => {
        showMessageOK('操作成功！')
        closeLoading()
      })
      .catch((e) => {
        showMessageError('操作失败：' + e.message)
        closeLoading()
      })
  }
</script>

<style scoped lang="scss">
  .card {
    padding: 0.5rem 2rem 1rem 2rem;

    .active-info {
      line-height: 1.5;
      padding: 10px 0 30px 0;
    }

    .el-descriptions {
      margin-bottom: 20px;

      .iconfont {
        font-size: 20px;
      }

      .selected {
        color: #0bc15f;
      }

      .closed {
        color: #da0d54;
      }
    }
  }
</style>
