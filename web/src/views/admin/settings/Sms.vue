<template>
  <div class="card">
    <div class="card-header bg-white text-center p-3">
      <h5 class="card-title mb-0">短信服务配置</h5>
    </div>
    <div class="card-body">
      <div class="container">
        <el-tabs v-model="activeTab" type="border-card">
          <el-tab-pane name="aliyun">
            <template #label>
              <div class="d-flex align-items-center">
                <i class="iconfont icon-aliyun"></i>
                <span class="ms-2">阿里云短信</span>
              </div>
            </template>

            <div class="alert alert-primary mt-2 mb-4" role="alert">
              如果你不知道怎么获取这些配置信息，请参考文档：
              <a
                href="https://docs.geekai.me/plus/config/sms.html#%E9%98%BF%E9%87%8C%E4%BA%91"
                target="_blank"
                >阿里云短信配置</a
              >。
            </div>

            <form class="form">
              <div class="mb-3">
                <label class="form-label">AccessKey</label>
                <input
                  type="text"
                  class="form-control"
                  v-model="configs.aliyun.access_key"
                />
              </div>
              <div class="mb-3">
                <label class="form-label">AccessSecret</label>
                <PasswordInput v-model="configs.aliyun.access_secret" />
              </div>
              <div class="mb-3">
                <label class="form-label">短信签名</label>
                <input
                  type="text"
                  class="form-control"
                  v-model="configs.aliyun.sign"
                />
              </div>
              <div class="mb-3">
                <label class="form-label">验证码模板ID</label>
                <input
                  type="text"
                  class="form-control"
                  v-model="configs.aliyun.code_temp_id"
                />
              </div>
            </form>
          </el-tab-pane>
          <el-tab-pane name="bao">
            <template #label>
              <div class="d-flex align-items-center">
                <i class="iconfont icon-sms"></i>
                <span class="ms-2">短信宝</span>
              </div>
            </template>
            <div class="alert alert-primary mt-2 mb-4" role="alert">
              如果你不知道怎么获取这些配置信息，请参考文档：
              <a
                href="https://docs.geekai.me/plus/config/sms.html#%E7%9F%AD%E4%BF%A1%E5%AE%9D"
                target="_blank"
                >短信宝配置</a
              >。
            </div>

            <form class="form">
              <div class="mb-3">
                <label class="form-label">用户名</label>
                <input
                  type="text"
                  class="form-control"
                  v-model="configs.bao.username"
                />
              </div>
              <div class="mb-3">
                <label class="form-label">密码</label>
                <PasswordInput v-model="configs.bao.password" />
              </div>
              <div class="mb-3">
                <label class="form-label">短信签名</label>
                <input
                  type="text"
                  class="form-control"
                  v-model="configs.bao.sign"
                />
              </div>
              <div class="mb-3">
                <label class="form-label">验证码模板</label>
                <input
                  type="text"
                  class="form-control"
                  v-model="configs.bao.code_temp"
                />
              </div>
            </form>
          </el-tab-pane>
        </el-tabs>
        <div class="mt-3">
          <label class="form-label mr-2">默认使用</label>
          <el-radio-group v-model="configs.active" size="large">
            <el-radio value="aliyun" border>阿里云</el-radio>
            <el-radio value="bao" border>短信宝</el-radio>
          </el-radio-group>
        </div>
        <div class="d-flex justify-content-center mt-4">
          <button type="button" class="btn btn-primary" @click="saveSmsConfig">
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
    showConfirm,
    showLoading,
    showMessageError,
    showMessageOK,
  } from '@/js/utils/dialog.js'
  import { httpGet, httpPost } from '@/js/utils/http.js'
  import { dateFormat } from '@/js/utils/libs.js'
  import { ElTabPane, ElTabs } from 'element-plus'
  import { el } from 'element-plus/es/locale/index.mjs'
  import { onMounted, ref } from 'vue'
  import PasswordInput from '@/components/PasswordInput.vue'

  const activeTab = ref('aliyun')
  const configs = ref({
    active: 'aliyun',
    aliyun: {
      access_key: '',
      access_secret: '',
      sign: '',
      code_temp_id: '',
    },
    bao: {
      username: '',
      password: '',
      sign: '',
      code_temp_id: '',
    },
  })

  onMounted(async () => {
    try {
      const res = await httpGet('/api/admin/config/get?name=sms')
      configs.value = Object.assign(configs.value, res.data)
    } catch (e) {
      configs.value = {
        active: 'aliyun',
        aliyun: {
          access_key: '',
          access_secret: '',
          sign: '',
          code_temp_id: '',
        },
        bao: {
          username: '',
          password: '',
          sign: '',
          code_temp: '',
        },
      }
    }
  })

  const saveSmsConfig = async () => {
    showLoading()
    try {
      await httpPost('/api/admin/config/update/sms', configs.value)
      showMessageOK('操作成功！')
    } catch (e) {
      showMessageError('操作失败！')
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
