<template>
  <div class="card">
    <div class="card-header bg-white text-center p-3">
      <h5 class="card-title mb-0">邮件服务设置</h5>
    </div>
    <div class="card-body">
      <div class="container">
        <div class="alert alert-primary mt-2 mb-4" role="alert">
          如果你不知道怎么获取这些配置信息，请参考文档：
          <a
            href="https://docs.geekai.me/plus/config/sms.html#%E9%82%AE%E4%BB%B6%E6%9C%8D%E5%8A%A1%E9%85%8D%E7%BD%AE"
            target="_blank"
            >邮件服务配置</a
          >。
        </div>
        <form class="form">
          <div class="mb-3">
            <label class="form-label"
              >邮件服务器地址
              <el-tooltip placement="top">
                <template #content>
                  推荐使用网易邮箱，<br />
                  国外邮箱推荐使用 Google 邮箱
                </template>
                <i class="iconfont icon-info"></i>
              </el-tooltip>
            </label>
            <input
              type="text"
              class="form-control"
              v-model="smtpConfig.host"
              placeholder="请输入邮件服务器地址，推荐使用网易邮箱"
            />
          </div>
          <div class="mb-3">
            <label class="form-label"
              >邮件服务器端口
              <el-tooltip placement="top">
                <template #content>
                  线上推荐使用465端口，<br />本地测试推荐使用25端口
                </template>
                <i class="iconfont icon-info"></i>
              </el-tooltip>
            </label>
            <input
              type="number"
              class="form-control"
              v-model="smtpConfig.port"
            />
          </div>
          <div class="mb-3">
            <label class="form-label">是否使用 TLS</label>
            <div>
              <el-switch v-model="smtpConfig.use_tls" />
            </div>
          </div>
          <div class="mb-3">
            <label class="form-label"
              >应用名称
              <el-tooltip placement="top">
                <template #content> 应用名称会显示在邮件的抬头 </template>
                <i class="iconfont icon-info"></i>
              </el-tooltip>
            </label>
            <input
              type="text"
              class="form-control"
              v-model="smtpConfig.app_name"
            />
          </div>
          <div class="mb-3">
            <label class="form-label">发件人邮箱地址</label>
            <input
              type="email"
              class="form-control"
              v-model="smtpConfig.from"
            />
          </div>
          <div class="mb-3">
            <label class="form-label"
              >发件人邮箱密码
              <el-tooltip placement="top">
                <template #content> 如果使用授权码，请输入授权码 </template>
                <i class="iconfont icon-info"></i>
              </el-tooltip>
            </label>
            <PasswordInput v-model="smtpConfig.password" />
          </div>
          <div class="d-flex justify-content-center mt-4">
            <button type="button" class="btn btn-primary" @click="saveConfig">
              提交保存
            </button>
          </div>
        </form>
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
  import PasswordInput from '@/components/PasswordInput.vue'

  const smtpConfig = ref({})

  onMounted(async () => {
    try {
      const res = await httpGet('/api/admin/config/get?name=smtp')
      smtpConfig.value = res.data || {}
    } catch (e) {
      smtpConfig.value = {
        use_tls: false,
        host: 'smtp.163.com',
        port: 465,
        app_name: 'GeekAI',
        from: '',
        password: '',
      }
    }
  })

  // 保存配置
  const saveConfig = () => {
    showLoading()
    httpPost('/api/admin/config/update/smtp', smtpConfig.value)
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
  }
</style>
