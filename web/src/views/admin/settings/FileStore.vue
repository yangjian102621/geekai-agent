<template>
  <div class="card">
    <div class="card-header bg-white text-center p-3">
      <h5 class="card-title mb-0">文件存储设置</h5>
    </div>
    <div class="card-body">
      <div class="container">
        <el-tabs v-model="activeTab" type="border-card">
          <el-tab-pane name="local">
            <template #label>
              <div class="d-flex align-items-center">
                <i class="iconfont icon-localstorage"></i>
                <span class="ms-2">本地存储</span>
              </div>
            </template>
            <div class="mb-3">
              <label class="form-label"
                >文件存储根目录
                <el-tooltip placement="top">
                  <template #content>
                    可以是绝对路径，如：/data/static/upload<br />也可以是相对路径，如：./static/upload
                  </template>
                  <i class="iconfont icon-info"></i>
                </el-tooltip>
              </label>
              <input
                type="text"
                class="form-control"
                v-model="localConfig.base_path"
                placeholder="请输入文件存储根目录，如：./static/upload"
              />
            </div>
            <div class="mb-3">
              <label class="form-label"
                >文件访问根 URL
                <el-tooltip placement="top">
                  <template #content>
                    可以是绝对路径，如：https://oss.geekai.me/static/upload
                    <br />也可以是相对路径，如：/static/upload
                  </template>
                  <i class="iconfont icon-info"></i>
                </el-tooltip>
              </label>
              <input
                type="text"
                class="form-control"
                v-model="localConfig.base_url"
              />
            </div>
          </el-tab-pane>
          <el-tab-pane name="qiniu">
            <template #label>
              <div class="d-flex align-items-center">
                <i class="iconfont icon-qiniu"></i>
                <span class="ms-2">七牛云</span>
              </div>
            </template>
            <div class="alert alert-primary mt-2 mb-4" role="alert">
              如果你不知道怎么获取这些配置信息，请参考文档：
              <a
                href="https://docs.geekai.me/plus/config/oss.html#%E4%B8%83%E7%89%9B%E4%BA%91-oss-%E9%85%8D%E7%BD%AE"
                target="_blank"
                >七牛云配置</a
              >。
            </div>

            <div class="mb-3">
              <label class="form-label">Access Key</label>
              <input
                type="text"
                class="form-control"
                v-model="qiniuConfig.access_key"
                placeholder="请输入七牛云 AccessKey"
              />
            </div>
            <div class="mb-3">
              <label class="form-label">Access Secret</label>
              <PasswordInput v-model="qiniuConfig.access_secret" />
            </div>
            <div class="mb-3">
              <label class="form-label">Bucket</label>
              <input
                type="text"
                class="form-control"
                v-model="qiniuConfig.bucket"
              />
            </div>
            <div class="mb-3">
              <label class="form-label"
                >区域（Zone）
                <el-tooltip
                  placement="top"
                  content="华南：z2，华东：z0，华北：z1，北美：na0，新加坡：as0"
                >
                  <i class="iconfont icon-info"></i>
                </el-tooltip>
              </label>
              <input
                type="text"
                class="form-control"
                placeholder="请输入区域，如：z2"
                v-model="qiniuConfig.zone"
              />
            </div>
            <div class="mb-3">
              <label class="form-label">域名</label>
              <input
                type="text"
                class="form-control"
                v-model="qiniuConfig.domain"
                placeholder="请输入七牛云Bucket绑定的域名"
              />
            </div>
          </el-tab-pane>
          <el-tab-pane name="aliyun">
            <template #label>
              <div class="d-flex align-items-center">
                <i class="iconfont icon-aliyun"></i>
                <span class="ms-2">阿里云</span>
              </div>
            </template>
            <div class="alert alert-primary mt-2 mb-4" role="alert">
              如果你不知道怎么获取这些配置信息，请参考文档：
              <a
                href="https://docs.geekai.me/plus/config/oss.html#%E9%98%BF%E9%87%8C%E4%BA%91-oss-%E9%85%8D%E7%BD%AE"
                target="_blank"
                >阿里云配置</a
              >。
            </div>
            <div class="mb-3">
              <label class="form-label">Endpoint</label>
              <input
                type="text"
                class="form-control"
                v-model="aliyunConfig.endpoint"
                placeholder="请输入阿里云 Endpoint，如：oss-cn-shenzhen.aliyuncs.com"
              />
            </div>
            <div class="mb-3">
              <label class="form-label">Access Key</label>
              <input
                type="text"
                class="form-control"
                v-model="aliyunConfig.access_key"
                placeholder="请输入阿里云 AccessKey"
              />
            </div>
            <div class="mb-3">
              <label class="form-label">Access Secret</label>
              <PasswordInput
                v-model="aliyunConfig.access_secret"
                placeholder="请输入阿里云 AccessSecret"
              />
            </div>
            <div class="mb-3">
              <label class="form-label">Bucket</label>
              <input
                type="text"
                class="form-control"
                v-model="aliyunConfig.bucket"
                placeholder="请输入阿里云 Bucket 名称"
              />
            </div>
            <div class="mb-3">
              <label class="form-label">域名</label>
              <input
                type="text"
                class="form-control"
                v-model="aliyunConfig.domain"
                placeholder="请输入阿里云Bucket绑定的域名"
              />
            </div>
          </el-tab-pane>
          <el-tab-pane name="minio">
            <template #label>
              <div class="d-flex align-items-center">
                <i class="iconfont icon-minio"></i>
                <span class="ms-2">Minio</span>
              </div>
            </template>
            <div class="alert alert-primary mt-2 mb-4" role="alert">
              如果你不知道怎么获取这些配置信息，请参考文档：
              <a
                href="https://docs.geekai.me/plus/config/oss.html#%E6%90%AD%E5%BB%BA-minio-%E5%AD%98%E5%82%A8%E6%9C%8D%E5%8A%A1"
                target="_blank"
                >Minio 配置</a
              >。
            </div>
            <div class="mb-3">
              <label class="form-label">Endpoint</label>
              <input
                type="text"
                class="form-control"
                v-model="minioConfig.endpoint"
                placeholder="请输入 Minio Endpoint，如：https://oss.geekai.me"
              />
            </div>
            <div class="mb-3">
              <label class="form-label">Access Key</label>
              <input
                type="text"
                class="form-control"
                v-model="minioConfig.access_key"
                placeholder="请输入 Minio AccessKey"
              />
            </div>
            <div class="mb-3">
              <label class="form-label">Access Secret</label>
              <PasswordInput
                v-model="minioConfig.access_secret"
                placeholder="请输入 Minio AccessSecret"
              />
            </div>
            <div class="mb-3">
              <label class="form-label">Bucket</label>
              <input
                type="text"
                class="form-control"
                v-model="minioConfig.bucket"
                placeholder="请输入 Minio Bucket 名称"
              />
            </div>
            <div class="mb-3">
              <label class="form-label"
                >是否启用SSL
                <el-tooltip placement="top">
                  <template #content>
                    如果启用SSL，则需要配置域名，如：https://oss.geekai.me
                  </template>
                  <i class="iconfont icon-info"></i>
                </el-tooltip>
              </label>
              <div class="d-flex align-items-center">
                <el-switch v-model="minioConfig.use_ssl" />
              </div>
            </div>
            <div class="mb-3">
              <label class="form-label">域名</label>
              <input
                type="text"
                class="form-control"
                v-model="minioConfig.domain"
                placeholder="请输入 Minio Bucket 绑定的域名"
              />
            </div>
          </el-tab-pane>
        </el-tabs>

        <div class="mt-3">
          <label class="form-label mr-2">默认使用</label>
          <el-radio-group v-model="active" size="large">
            <el-radio value="local" border>本地存储</el-radio>
            <el-radio value="aliyun" border>阿里云</el-radio>
            <el-radio value="qiniu" border>七牛云</el-radio>
            <el-radio value="minio" border>Minio</el-radio>
          </el-radio-group>
        </div>

        <div class="d-flex justify-content-center mt-4 mb-2">
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

  const activeTab = ref('local')

  const localConfig = ref({})
  const qiniuConfig = ref({})
  const aliyunConfig = ref({})
  const minioConfig = ref({})
  const active = ref('')

  onMounted(async () => {
    try {
      const res = await httpGet('/api/admin/config/get?name=oss')
      localConfig.value = res.data.local
      qiniuConfig.value = res.data.qiniu
      aliyunConfig.value = res.data.aliyun
      minioConfig.value = res.data.minio
      active.value = res.data.active
    } catch (e) {
      active.value = 'local'
      localConfig.value = {
        base_path: './static/upload',
        base_url: '/static/upload',
      }
      qiniuConfig.value = {
        access_key: '',
        access_secret: '',
        bucket: '',
        zone: 'z2',
        domain: '',
      }
      aliyunConfig.value = {
        endpoint: 'oss-cn-hangzhou.aliyuncs.com',
        access_key: '',
        access_secret: '',
        bucket: '',
        domain: '',
      }
      minioConfig.value = {
        endpoint: '',
        access_key: '',
        access_secret: '',
        bucket: '',
        domain: '',
      }
    }
  })

  const saveConfig = async () => {
    if (active.value === '') {
      showMessageError('请选择默认存储介质')
      return
    }
    showLoading()
    try {
      await httpPost('/api/admin/config/update/oss', {
        local: localConfig.value,
        qiniu: qiniuConfig.value,
        aliyun: aliyunConfig.value,
        minio: minioConfig.value,
        active: active.value,
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
