<template>
  <div class="card">
    <div class="card-body">
      <div class="table-responsive" v-loading="loading">
        <div class="search-container d-flex flex-wrap align-items-end gap-2">
          <span class="search-item">
            <el-input
              type="text"
              v-model="query.name"
              placeholder="请输入应用名称"
              clearable
              @keyup.enter="fetchData(1)"
              @clear="fetchData(1)"
            />
          </span>

          <span class="search-item">
            <el-select
              v-model="query.check"
              placeholder="审核状态"
              style="width: 120px"
              clearable
              @change="fetchData(1)"
            >
              <el-option label="审核通过" value="1" />
              <el-option label="审核不通过" value="2" />
              <el-option label="待审核" value="0" />
            </el-select>
          </span>

          <span class="search-item">
            <el-button type="primary" @click="fetchData(1)">
              <i class="iconfont icon-search"></i>
            </el-button>
          </span>
        </div>

        <el-table
          :data="dataSets.items"
          border
          class="data-table"
          :row-key="(row) => row.id"
          table-layout="auto"
        >
          <el-table-column type="selection" width="38"></el-table-column>
          <el-table-column label="应用名称" prop="name" />
          <el-table-column label="分类" prop="cname" />
          <el-table-column label="图标">
            <template #default="{ row }">
              <el-image
                :src="row.icon"
                style="height: 32px; width: 32px"
                class="rounded-circle"
              >
                <template #error>
                  <i class="iconfont icon-image fs-2"></i>
                </template>
              </el-image>
            </template>
          </el-table-column>
          <el-table-column label="类型" prop="type" />
          <el-table-column label="API URL">
            <template #default="{ row }">
              <span>{{ substr(row.configs.api_url, 20) }}</span>
              <i
                class="iconfont icon-copy ms-1 mt-1"
                :data-clipboard-text="row.configs.api_url"
              ></i>
            </template>
          </el-table-column>
          <el-table-column label="API Token">
            <template #default="{ row }">
              <span>{{ substr(row.configs.token, 10) }}</span>
              <i
                class="iconfont icon-copy ms-1 mt-1"
                :data-clipboard-text="row.configs.token"
              ></i>
            </template>
          </el-table-column>
          <el-table-column label="上架状态">
            <template #default="{ row }">
              <el-switch v-model="row.enabled" @change="enable(row)" />
            </template>
          </el-table-column>
          <el-table-column label="审核状态">
            <template #default="{ row }">
              <el-tag
                :type="
                  row.check === 0
                    ? 'info'
                    : row.check === 1
                    ? 'success'
                    : 'danger'
                "
              >
                {{
                  row.check === 0
                    ? '待审核'
                    : row.check === 1
                    ? '审核通过'
                    : '审核不通过'
                }}
              </el-tag>
            </template>
          </el-table-column>
          <el-table-column label="每次对话积分" prop="score" />

          <el-table-column label="创建时间">
            <template #default="{ row }">
              <span>{{ dateFormat(row['created_at']) }}</span>
            </template>
          </el-table-column>

          <el-table-column label="操作" fixed="right">
            <template #default="{ row }">
              <el-dropdown placement="bottom" trigger="hover">
                <button class="btn btn-primary btn-sm">
                  <i class="iconfont icon-more-vertical"></i>
                </button>
                <template #dropdown>
                  <el-dropdown-menu>
                    <el-dropdown-item @click="showCheckDialog(row)">
                      <i class="iconfont icon-check2"></i> 审核
                    </el-dropdown-item>
                    <el-dropdown-item @click="remove(row)">
                      <span class="text-danger">
                        <i class="iconfont icon-remove"></i> 删除
                      </span>
                    </el-dropdown-item>
                  </el-dropdown-menu>
                </template>
              </el-dropdown>
            </template>
          </el-table-column>
          <template #empty>
            <el-empty description="暂无数据" />
          </template>
        </el-table>
      </div>

      <div class="pagination p-3" v-if="dataSets.total > 0">
        <Pagination
          :total="dataSets.total"
          :pageSize="dataSets.pageSize"
          :currentPage="dataSets.page"
          :layout="['total', 'prev', 'pages', 'sizes', 'next']"
          @update:currentPage="fetchData"
          @update:pageSize="dataSets.pageSize = $event"
        />
      </div>
    </div>

    <!-- 审核对话框 -->
    <model-dialog
      :modelValue="checkDialog.visible"
      title="创作者应用审核"
      @cancel="checkDialog.visible = false"
      @confirm="submitCheck"
      confirm-text="确定"
      cancel-text="取消"
    >
      <form
        :key="checkDialog.key"
        novalidate
        :class="{ 'was-validated': checkDialog.form.validated }"
        @submit.prevent
      >
        <div class="mb-3">
          <label class="form-label">应用名称</label>
          <input
            type="text"
            class="form-control"
            :value="checkDialog.app && checkDialog.app.name"
            disabled
          />
        </div>
        <div class="mb-3">
          <label class="form-label">应用简介</label>
          <textarea
            class="form-control"
            :value="checkDialog.app && checkDialog.app.summary"
            disabled
            rows="3"
          ></textarea>
        </div>
        <div class="mb-3">
          <div>
            <label class="form-label"
              >审核结果 <span class="text-danger">*</span></label
            >
          </div>
          <el-radio-group v-model="checkDialog.form.check" size="large">
            <el-radio :value="1" border>审核通过</el-radio>
            <el-radio :value="2" border>审核不通过</el-radio>
          </el-radio-group>
        </div>
        <div class="mb-3">
          <label class="form-label"
            >审核备注
            <span class="text-danger" v-if="checkDialog.form.check === 2"
              >*</span
            ></label
          >
          <textarea
            class="form-control"
            v-model="checkDialog.form.check_note"
            :required="checkDialog.form.check === 2"
            rows="3"
            :placeholder="
              checkDialog.form.check === 2
                ? '请填写不通过的原因'
                : '可选填写备注信息'
            "
          ></textarea>
          <div class="invalid-feedback" v-if="checkDialog.form.check === 2">
            请填写不通过原因
          </div>
        </div>
      </form>
    </model-dialog>
  </div>
</template>

<script setup>
  import ModelDialog from '@/components/ModelDialog.vue'
  import Pagination from '@/components/Pagination.vue'
  import { validateForm } from '@/js/store/common.js'
  import { closeLoading, showConfirm, showLoading } from '@/js/utils/dialog.js'
  import { httpGet, httpPost } from '@/js/utils/http.js'
  import { dateFormat, substr } from '@/js/utils/libs.js'
  import ClipboardJS from 'clipboard'
  import { ElMessage } from 'element-plus'
  import { nextTick, onMounted, onUnmounted, ref } from 'vue'

  const dataSets = ref({ total: 0, page: 1, pageSize: 10, items: [] })
  const query = ref({})
  const loading = ref(true)
  const clipboard = ref(null)

  // 获取数据
  const fetchData = (page) => {
    dataSets.value.page = page || 1
    query.value.page = dataSets.value.page
    query.value.page_size = dataSets.value.pageSize
    query.value.check = query.value.check
    query.value.name = query.value.name
    httpGet('/api/admin/creator/app/list', query.value)
      .then((res) => {
        dataSets.value = res.data
        loading.value = false
      })
      .catch((e) => {
        ElMessage.error('获取数据失败：' + e.message)
        loading.value = false
      })
  }

  onMounted(() => {
    clipboard.value = new ClipboardJS('.icon-copy')
    clipboard.value.on('success', () => {
      ElMessage.success('复制成功！')
    })

    clipboard.value.on('error', () => {
      ElMessage.error('复制失败！')
    })
    fetchData()
  })

  onUnmounted(() => {
    clipboard.value.destroy()
  })

  const enable = (row) => {
    httpPost('/api/admin/app/set', {
      id: row.id,
      name: 'enabled',
      value: row.enabled,
    })
      .then(() => {
        ElMessage.success('操作成功！')
      })
      .catch((e) => {
        ElMessage.error('操作失败：' + e.message)
      })
  }

  const remove = function (row) {
    showConfirm('删除提示', '确定要删除当前记录吗?？', async () => {
      showLoading()
      try {
        await httpGet('/api/admin/creator/app/remove', {
          id: row.id,
          creator_id: row.creator_id,
        })
        ElMessage.success('删除成功！')
        fetchData()
      } catch (e) {
        ElMessage.error('删除失败：' + e.message)
      } finally {
        closeLoading()
      }
    })
  }

  // 审核对话框
  const checkDialog = ref({
    visible: false,
    key: 0,
    app: {},
    form: {
      check_note: '审核通过',
    },
  })
  // 校验规则
  const checkFormRules = {
    check_note: {
      required: (form) => form.check === 2,
      message: '请填写不通过原因',
    },
  }
  const checkFormErrors = ref({})
  // 显示审核对话框
  const showCheckDialog = (row) => {
    checkDialog.value.app = row
    checkDialog.value.form.check = row.check
    checkDialog.value.form.check_note = row.check_note
    checkDialog.value.visible = true
    // 清空校验错误
    checkFormErrors.value = {}
  }

  // 提交审核
  const submitCheck = async () => {
    checkDialog.value.form.validated = true
    await nextTick()
    // 校验
    if (
      !validateForm(
        checkDialog.value.form,
        checkFormRules,
        checkFormErrors.value
      )
    ) {
      return
    }

    httpPost(
      `/api/admin/creator/app/check?id=${checkDialog.value.app.id}`,
      checkDialog.value.form
    )
      .then(() => {
        ElMessage.success('审核成功')
        checkDialog.value.visible = false
        fetchData()
      })
      .catch((e) => {
        ElMessage.error('审核失败：' + e.message)
      })
  }
</script>

<style scoped lang="scss">
  @use '@/assets/css/admin/admin.scss';
</style>
