<template>
  <div class="card">
    <div class="card-header bg-white text-center p-3">
      <h5 class="card-title mb-0">支付设置</h5>
    </div>
    <div class="card-body">
      <div class="container">
        <el-tabs v-model="activeTab" type="border-card">
          <el-tab-pane name="alipay">
            <template #label>
              <div class="d-flex align-items-center text-blue-600">
                <i class="iconfont icon-alipay"></i>
                <span class="ms-2">支付宝</span>
              </div>
            </template>
            <div class="alert alert-primary mt-2 mb-4" role="alert">
              如果你不知道怎么获取这些配置信息，请参考文档：
              <a
                href="https://docs.geekai.me/plus/config/payment.html#%E6%94%AF%E4%BB%98%E5%AE%9D%E9%85%8D%E7%BD%AE"
                target="_blank"
                >支付宝配置</a
              >。
            </div>

            <div class="mb-3">
              <label class="form-label">应用 ID</label>
              <input
                type="text"
                class="form-control"
                v-model="alipayConfig.app_id"
                placeholder="请输入应用 AppId"
              />
            </div>
            <div class="mb-3">
              <label class="form-label">应用私钥</label>
              <textarea
                rows="5"
                class="form-control"
                v-model="alipayConfig.private_key"
              />
            </div>

            <div class="mb-3">
              <label class="form-label">支付宝公钥</label>
              <textarea
                rows="5"
                class="form-control"
                v-model="alipayConfig.alipay_public_key"
              />
            </div>

            <div class="mb-3">
              <label class="form-label"
                >支付回调域名
                <el-tooltip
                  placement="top"
                  content="请确保回调域名已备案且在支付宝应用添加了白名单"
                >
                  <i class="iconfont icon-info"></i>
                </el-tooltip>
              </label>
              <input
                type="text"
                class="form-control"
                v-model="alipayConfig.domain"
                placeholder="请输入支付回调域名，如：https://www.geekai.me"
              />
            </div>
            <div class="mb-3">
              <label class="form-label">是否启用该支付渠道 </label>
              <div>
                <el-switch v-model="alipayConfig.enabled" />
              </div>
            </div>
          </el-tab-pane>
          <el-tab-pane name="wxpay">
            <template #label>
              <div class="d-flex align-items-center text-green-600">
                <i class="iconfont icon-wechat-pay"></i>
                <span class="ms-2">微信支付</span>
              </div>
            </template>

            <div class="alert alert-primary mt-2 mb-4" role="alert">
              如果你不知道怎么获取这些配置信息，请参考文档：
              <a
                href="https://docs.geekai.me/plus/config/payment.html#%E5%BE%AE%E4%BF%A1%E6%94%AF%E4%BB%98%E9%85%8D%E7%BD%AE"
                target="_blank"
                >微信支付配置</a
              >。
            </div>

            <div class="mb-3">
              <label class="form-label">应用ID</label>
              <input
                type="text"
                class="form-control"
                v-model="wxpayConfig.app_id"
                placeholder="请输入微信支付应用 AppId"
              />
            </div>
            <div class="mb-3">
              <label class="form-label">商户号(MchID)</label>
              <input
                type="text"
                class="form-control"
                v-model="wxpayConfig.mch_id"
              />
            </div>
            <div class="mb-3">
              <label class="form-label">商户证书序列号</label>
              <input
                type="text"
                class="form-control"
                v-model="wxpayConfig.serial_no"
              />
            </div>
            <div class="mb-3">
              <label class="form-label">商户证书私钥</label>
              <textarea
                rows="5"
                class="form-control"
                v-model="wxpayConfig.private_key"
              />
            </div>
            <div class="mb-3">
              <label class="form-label">API V3 秘钥</label>
              <PasswordInput v-model="wxpayConfig.api_v3_key" />
            </div>
            <div class="mb-3">
              <label class="form-label"
                >支付回调域名
                <el-tooltip
                  placement="top"
                  content="请确保回调域名已备案且在微信支付应用添加了白名单"
                >
                  <i class="iconfont icon-info"></i>
                </el-tooltip>
              </label>
              <input
                type="text"
                class="form-control"
                v-model="wxpayConfig.domain"
                placeholder="请输入支付回调域名，如：https://www.geekai.me"
              />
            </div>
            <div class="mb-3">
              <label class="form-label">是否启用该支付渠道 </label>
              <div>
                <el-switch v-model="wxpayConfig.enabled" />
              </div>
            </div>
          </el-tab-pane>
          <el-tab-pane name="geekpay">
            <template #label>
              <div class="d-flex align-items-center text-purple-600">
                <i class="iconfont icon-reward"></i>
                <span class="ms-2">易支付</span>
              </div>
            </template>
            <div class="alert alert-primary mt-2 mb-4" role="alert">
              如果你不知道怎么获取这些配置信息，请参考文档：
              <a
                href="https://docs.geekai.me/plus/config/payment.html#%E6%98%93%E6%94%AF%E4%BB%98%E5%BC%80%E9%80%9A"
                target="_blank"
                >易支付开通</a
              >。
            </div>
            <div class="mb-3">
              <label class="form-label">API URL</label>
              <input
                type="text"
                class="form-control"
                v-model="epayConfig.api_url"
                placeholder="请输入支付 API 地址，如：https://api.geekpay.com"
              />
            </div>
            <div class="mb-3">
              <label class="form-label">商户PID</label>
              <input
                type="text"
                class="form-control"
                v-model="epayConfig.app_id"
              />
            </div>
            <div class="mb-3">
              <label class="form-label">商户密钥</label>
              <PasswordInput v-model="epayConfig.private_key" />
            </div>
            <div class="mb-3">
              <label class="form-label">支付回调域名</label>
              <input
                type="text"
                class="form-control"
                v-model="epayConfig.domain"
              />
            </div>
            <div class="mb-3">
              <label class="form-label">是否启用该支付渠道 </label>
              <div>
                <el-switch v-model="epayConfig.enabled" />
              </div>
            </div>
          </el-tab-pane>
        </el-tabs>
        <div class="d-flex justify-content-center pt-4">
          <button type="button" class="btn btn-primary" @click="saveConfig">
            提交保存
          </button>
        </div>
      </div>
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
  import { httpGet, httpPost } from '@/js/utils/http.js'
  import { ElTabPane, ElTabs } from 'element-plus'
  import { onMounted, ref } from 'vue'
  import PasswordInput from '@/components/PasswordInput.vue'

  const activeTab = ref('alipay')

  const alipayConfig = ref({})
  const wxpayConfig = ref({})
  const epayConfig = ref({})

  onMounted(async () => {
    try {
      const res = await httpGet('/api/admin/config/get?name=payment')
      alipayConfig.value = res.data.alipay
      wxpayConfig.value = res.data.wxpay
      epayConfig.value = res.data.epay
    } catch (e) {
      alipayConfig.value = {
        api_url: '',
        private_key: '',
        app_id: '',
        public_key: '',
        alipay_public_key: '',
        root_cert: '',
        domain: '',
        enabled: true,
      }
      wxpayConfig.value = {
        api_url: '',
        mch_id: '',
        api_key: '',
        app_id: '',
        app_secret: '',
      }
      epayConfig.value = {
        api_url: '',
        pid: '',
        key: '',
      }
    }
  })

  const saveConfig = async () => {
    showLoading()
    try {
      await httpPost('/api/admin/config/update/payment', {
        alipay: alipayConfig.value,
        wxpay: wxpayConfig.value,
        geekpay: epayConfig.value,
      })
      showMessageOK('操作成功！')
    } catch (e) {
      showMessageError('操作失败：' + e.message)
    } finally {
      closeLoading()
    }
  }
</script>

<style scoped lang="scss">
  .card {
    padding: 0.5rem 2rem 1rem 2rem;
  }
</style>
