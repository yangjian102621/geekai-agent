<template>
  <div class="card">
    <div class="card-header bg-white text-center p-3">
      <h5 class="card-title mb-0">公告设置</h5>
    </div>
    <div class="card-body">
      <form class="form">
        <md-editor
          class="mgb20"
          v-model="notice"
          @on-upload-img="onUploadImg"
        />

        <div class="d-flex justify-content-center mt-5">
          <button type="button" class="btn btn-primary" @click="saveNotice">
            提交
          </button>
        </div>
      </form>
    </div>
  </div>
</template>

<script setup>
  import { showMessageError, showMessageOK } from '@/js/utils/dialog'
  import { httpGet, httpPost } from '@/js/utils/http.js'
  import MdEditor from 'md-editor-v3'
  import 'md-editor-v3/lib/style.css'

  const notice = ref('')
  httpGet('/api/admin/config/get?name=notice')
    .then((res) => {
      notice.value = res.data['content']
    })
    .catch((e) => {
      showMessageError('公告信息获取失败: ' + e.message)
    })

  const saveNotice = () => {
    httpPost('/api/admin/config/update/notice', {
      content: notice.value,
    })
      .then(() => {
        showMessageOK('操作成功！')
      })
      .catch((e) => {
        showMessageError('操作失败：' + e.message)
      })
  }

  // 编辑期文件上传处理
  const onUploadImg = (files, callback) => {
    Promise.all(
      files.map((file) => {
        return new Promise((rev, rej) => {
          const formData = new FormData()
          formData.append('file', file, file.name)
          // 执行上传操作
          httpPost('/api/admin/upload', formData)
            .then((res) => rev(res))
            .catch((error) => rej(error))
        })
      })
    )
      .then((res) => {
        showMessageOK('上传成功')
        callback(res.map((item) => item.data.url))
      })
      .catch((e) => {
        showMessageError('图片上传失败:' + e.message)
      })
  }
</script>

<style scoped lang="scss">
  .card {
    padding: 0.5rem 2rem 1rem 2rem;

    .iconfont {
      cursor: help;
    }
  }
</style>
