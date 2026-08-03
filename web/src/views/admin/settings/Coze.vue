<template>
  <div class="card">
    <div class="card-header bg-white text-center p-3">
      <h5 class="card-title mb-0">Coze API 设置</h5>
    </div>
    <div class="card-body">
      <div class="container">
        <div class="alert alert-primary mt-2 mb-4" role="alert">
          如果你不知道怎么获取这些配置信息，请参考文档：
          <a
            href="https://docs.geekai.me/agent/config/coze.html"
            target="_blank"
            class="hover:underline"
            >创建Coze授权应用</a
          >。
        </div>
        <form class="form">
          <div class="mb-3">
            <label class="form-label">API URL</label>
            <input
              type="text"
              class="form-control"
              v-model="cozeConfig.api_url"
            />
          </div>
          <div class="mb-3">
            <label class="form-label">Coze 空间ID</label>
            <input
              type="text"
              class="form-control"
              v-model="cozeConfig.space_id"
            />
          </div>
          <div class="mb-3">
            <label class="form-label">授权应用ID</label>
            <input
              type="text"
              class="form-control"
              v-model="cozeConfig.app_id"
            />
          </div>

          <div class="mb-3">
            <label class="form-label">授权公钥ID</label>
            <input
              type="text"
              class="form-control"
              v-model="cozeConfig.public_key_id"
            />
          </div>

          <div class="mb-3">
            <label class="form-label">授权私钥</label>
            <textarea
              type="textarea"
              rows="5"
              class="form-control"
              v-model="cozeConfig.private_key"
            />
          </div>

          <div class="d-flex justify-content-center pt-3">
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

  const configs = ref({})
  const cozeConfig = ref({
    api_url: '',
    private_key: '',
    app_id: '',
    public_key_id: '',
    space_id: '',
  })

  onMounted(async () => {
    const res = await httpGet('/api/admin/config/get?name=coze')
    cozeConfig.value = res.data || {}
  })

  // 保存配置
  const saveConfig = () => {
    showLoading()
    configs.value.coze_config = cozeConfig.value
    httpPost('/api/admin/config/update/coze', cozeConfig.value)
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
