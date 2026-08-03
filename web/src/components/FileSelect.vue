<template>
  <el-container class="file-select-box d-flex">
    <div class="container-fluid">
      <div class="row g-2">
        <div class="item col-6 col-md-2 mb-2">
          <el-upload
            class="avatar-uploader"
            :auto-upload="true"
            :show-file-list="false"
            :http-request="uploadFile"
            accept=".doc,.docx,.jpg,.png,.jpeg,.xls,.xlsx,.ppt,.pptx,.pdf,.mp4,.mp3"
          >
            <div class="d-flex">
              <i class="iconfont icon-plus"></i>
            </div>
          </el-upload>
        </div>
        <div
          class="item col-6 col-md-2 mb-2"
          v-for="file in fileData.items"
          :key="file.url"
        >
          <el-tooltip
            class="box-item"
            effect="dark"
            :content="file.name"
            placement="top"
          >
            <el-image
              :src="file.url"
              fit="cover"
              v-if="isImage(file.ext)"
              @click="insertURL(file)"
            >
              <template #error>
                <div class="img-error">
                  <i class="iconfont icon-image"></i>
                </div>
              </template>
            </el-image>
            <el-image
              :src="GetFileIcon(file.ext)"
              fit="cover"
              v-else
              @click="insertURL(file)"
            />
          </el-tooltip>

          <div class="opt">
            <el-button
              type="danger"
              size="small"
              :icon="Delete"
              @click="removeFile(file)"
              circle
            />
          </div>
        </div>
      </div>

      <div class="d-flex justify-content-center pt-3 pb-2">
        <el-button
          class="text-primary"
          v-if="!fileData.isLastPage"
          @click="fetchFiles(fileData.page)"
          >加载更多</el-button
        >
        <span class="text-secondary" v-else>没有更多了</span>
      </div>
    </div>
  </el-container>
</template>

<script setup>
  import { frontUploadFile } from '@/js/store/common.js'
  import { httpGet, httpPost } from '@/js/utils/http'
  import { GetFileIcon, isImage, removeArrayItem } from '@/js/utils/libs'
  import { Delete } from '@element-plus/icons-vue'
  import { reactive, ref } from 'vue'
  import {
    showConfirm,
    showMessageError,
    showMessageOK,
  } from '../js/utils/dialog'

  const emits = defineEmits(['selected'])
  const show = ref(false)
  const fileData = reactive({
    items: [],
    page: 1,
    isLastPage: true,
  })

  onMounted(() => {
    fetchFiles(1)
  })

  const fetchFiles = (pageNo) => {
    if (pageNo === 1) show.value = true
    httpPost('/api/file/list', { page: pageNo || 1, page_size: 30 })
      .then((res) => {
        const { items, page, total_page } = res.data

        if (page === 1) {
          fileData.items = items
        } else {
          fileData.items = [...fileData.items, ...items]
        }

        fileData.isLastPage = page === total_page || total_page === 0

        if (!fileData.isLastPage) {
          fileData.page = page + 1
        }
      })
      .catch(() => {})
  }

  const uploadFile = (file) => {
    frontUploadFile(file, () => {
      fetchFiles(1)
    })
  }

  const removeFile = (file) => {
    showConfirm('删除文件', '确定要删除文件吗？此操作不可恢复！', () => {
      httpGet('/api/file/remove?id=' + file.id)
        .then(() => {
          fileData.items = removeArrayItem(fileData.items, file, (v1, v2) => {
            return v1.id === v2.id
          })
          showMessageOK('文件删除成功！')
          fetchFiles(1)
        })
        .catch((e) => {
          showMessageError('文件删除失败:' + e.message)
        })
    })
  }

  const insertURL = (file) => {
    show.value = false
    // 如果是相对路径，处理成绝对路径
    if (file.url.indexOf('http') === -1) {
      file.url = location.protocol + '//' + location.host + file.url
    }
    emits('selected', file)
  }
</script>

<style lang="scss">
  .file-select-box {
    .item {
      position: relative;
      display: flex;
      justify-content: center;
      align-items: center;

      .avatar-uploader {
        width: 100%;
        display: flex;
        justify-content: center;
        align-items: center;

        .el-upload {
          border: 1px dashed #e1e1e1;
          border-radius: 10px;
          width: 100%;
          max-width: 100px;
          height: 100px;
        }
      }

      .el-image {
        width: 100px;
        height: 100px;
        border: 1px solid #ffffff;
        border-radius: 10px;
        cursor: pointer;

        &:hover {
          border-color: var(--el-color-primary);
        }
      }

      .img-error {
        width: 100px;
        height: 100px;
        background: #f0f0f0;
        display: flex;
        justify-content: center;
        align-items: center;
      }

      .iconfont {
        color: var(--el-color-primary);
        font-size: 30px;
      }

      .opt {
        display: none;
        position: absolute;
        top: -8px;
        right: 8px;
      }

      @media (max-width: 768px) {
        .opt {
          display: block;
        }
      }

      &:hover {
        .opt {
          display: block;
        }
      }
    }
  }
</style>
