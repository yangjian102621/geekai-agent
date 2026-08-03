<template>
  <div class="card">
    <div class="card-body">
      <div class="table-responsive" v-loading="loading">
        <!-- 已审核列表搜索 -->
        <div class="search-container d-flex flex-wrap align-items-end gap-2">
          <span class="search-item">
            <el-input
              v-model="query.name"
              placeholder="请输入创作者名称"
              style="width: 200px"
              clearable
              @clear="fetchData(1)"
              @keyup.enter="fetchData(1)"
            />
          </span>

          <span class="search-item">
            <el-select
              v-model="query.check"
              placeholder="审核状态"
              clearable
              @change="fetchData(1)"
              style="width: 120px"
            >
              <el-option label="待审核" value="0" />
              <el-option label="审核通过" value="1" />
              <el-option label="审核不通过" value="2" />
            </el-select>
          </span>

          <span class="search-item">
            <el-button type="primary" @click="fetchData(1)">
              <i class="iconfont icon-search"></i>
            </el-button>
          </span>
        </div>

        <!-- 已审核列表表格 -->
        <el-table
          border
          class="data-table"
          :data="dataSets.items"
          table-layout="auto"
        >
          <el-table-column prop="id" label="ID" width="50" />
          <el-table-column label="创作者名称">
            <template #default="{ row }">
              <div class="flex gap-2 align-items-center">
                <img
                  v-if="row.logo"
                  :src="row.logo"
                  class="rounded w-[50px] h-[50px] object-cover"
                />
                <div>
                  <div class="text-base font-bold px-1">{{ row.name }}</div>
                  <div class="text-sm ml-1 text-gray-500 line-clamp-2">
                    {{ row.description }}
                  </div>
                </div>
              </div>
            </template>
          </el-table-column>
          <el-table-column prop="user_id" label="用户ID" width="50" />
          <el-table-column label="启用状态" width="60">
            <template #default="{ row }">
              <el-switch v-model="row.enabled" @change="toggleEnable(row)" />
            </template>
          </el-table-column>
          <el-table-column prop="fee" label="提现费率(%)" width="120" />
          <el-table-column prop="scores" label="积分" width="100" />
          <el-table-column label="审核状态" width="100">
            <template #default="{ row }">
              <el-tag v-if="row.check === 1" type="success">审核通过</el-tag>
              <el-tag v-else-if="row.check === 2" type="danger"
                >审核不通过</el-tag
              >
              <el-tag v-else type="info">待审核</el-tag>
            </template>
          </el-table-column>
          <el-table-column label="申请时间" width="180">
            <template #default="{ row }">
              {{ dateFormat(row.created_at) }}
            </template>
          </el-table-column>
          <el-table-column label="操作" fixed="right" width="50">
            <template #default="{ row }">
              <el-dropdown>
                <button class="btn btn-primary btn-sm">
                  <i class="iconfont icon-more-vertical"></i>
                </button>
                <template #dropdown>
                  <el-dropdown-menu>
                    <el-dropdown-item @click="showCheckDialog(row)">
                      <i class="iconfont icon-check2 text-success"></i>
                      审核
                    </el-dropdown-item>

                    <el-dropdown-item @click="showEditDialog(row)">
                      <i class="iconfont icon-edit text-primary"></i>
                      编辑
                    </el-dropdown-item>
                    <el-dropdown-item @click="deleteCreator(row)">
                      <i class="iconfont icon-remove text-danger"></i>
                      <span class="text-danger">删除</span>
                    </el-dropdown-item>
                    <el-dropdown-item @click="showScoreDialog(row)">
                      <i class="iconfont icon-edit text-warning"></i>
                      修改积分
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

        <!-- 分页 -->
        <div class="pagination p-3" v-if="dataSets.total > 0">
          <Pagination
            :total="dataSets.total"
            :pageSize="dataSets.page_size"
            :currentPage="dataSets.page"
            :layout="['total', 'prev', 'pages', 'sizes', 'next']"
            @update:currentPage="fetchData"
            @update:pageSize="dataSets.page_size = $event"
          />
        </div>
      </div>
    </div>

    <!-- 审核对话框 -->
    <model-dialog
      :modelValue="checkDialog.visible"
      title="审核创作者"
      @cancel="checkDialog.visible = false"
      @confirm="submitCheck"
      confirm-text="确定"
      cancel-text="取消"
    >
      <form
        ref="checkFormRef"
        novalidate
        :class="{ 'was-validated': checkDialog.form.validated }"
        @submit.prevent
      >
        <div class="mb-3">
          <label class="form-label">创作者名称</label>
          <input
            type="text"
            class="form-control"
            :value="checkDialog.creator && checkDialog.creator.name"
            disabled
          />
        </div>
        <div class="mb-3">
          <label class="form-label">创作者简介</label>
          <textarea
            class="form-control"
            :value="checkDialog.creator && checkDialog.creator.description"
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

    <!-- 编辑对话框 -->
    <model-dialog
      :modelValue="editDialog.visible"
      title="编辑创作者"
      @cancel="editDialog.visible = false"
      @confirm="submitEdit"
      confirm-text="确定"
      cancel-text="取消"
    >
      <form
        ref="editFormRef"
        novalidate
        :class="{ 'was-validated': editDialog.form.validated }"
        @submit.prevent
      >
        <div class="mb-3">
          <label class="form-label"
            >创作者名称 <span class="text-danger">*</span></label
          >
          <input
            type="text"
            class="form-control"
            v-model="editDialog.form.name"
            :class="{ 'is-invalid': editFormErrors.name }"
            required
          />
          <div class="invalid-feedback" v-if="editFormErrors.name">
            {{ editFormErrors.name }}
          </div>
        </div>
        <div class="mb-3">
          <label class="form-label"
            >简介 <span class="text-danger">*</span></label
          >
          <textarea
            class="form-control"
            v-model="editDialog.form.description"
            rows="3"
            :class="{ 'is-invalid': editFormErrors.description }"
            required
          ></textarea>
          <div class="invalid-feedback" v-if="editFormErrors.description">
            {{ editFormErrors.description }}
          </div>
        </div>
        <div class="mb-3">
          <label class="form-label"
            >Logo <span class="text-danger">*</span></label
          >
          <div class="d-flex gap-2">
            <div class="flex flex-col w-full">
              <input
                type="text"
                class="form-control"
                v-model="editDialog.form.logo"
                :class="{ 'is-invalid': editFormErrors.logo }"
                required
              />
              <div class="invalid-feedback" v-if="editFormErrors.logo">
                {{ editFormErrors.logo }}
              </div>
            </div>

            <el-upload
              :auto-upload="true"
              :show-file-list="false"
              :http-request="uploadLogo"
              class="flex-end"
            >
              <button type="button" class="btn btn-primary">
                <i class="iconfont icon-upload"></i>
              </button>
            </el-upload>
          </div>
        </div>
        <div class="mb-3">
          <label class="form-label"
            >提现费率(%) <span class="text-danger">*</span></label
          >
          <input
            type="number"
            class="form-control"
            v-model="editDialog.form.fee"
            :class="{ 'is-invalid': editFormErrors.fee }"
            placeholder="请输入0-100之间的整数"
            required
          />
          <div class="invalid-feedback" v-if="editFormErrors.fee">
            {{ editFormErrors.fee }}
          </div>
        </div>
      </form>
    </model-dialog>

    <!-- 新增积分调整对话框 -->
    <model-dialog
      :modelValue="scoreDialog.visible"
      title="修改创作者积分"
      @cancel="scoreDialog.visible = false"
      @confirm="submitScore"
      :confirm-loading="scoreDialog.loading"
      confirm-text="确定"
      cancel-text="取消"
    >
      <form @submit.prevent>
        <div class="mb-3">
          <el-radio-group v-model="scoreDialog.form.action" size="large">
            <el-radio :value="'inc'" border>增加</el-radio>
            <el-radio :value="'dec'" border>减少</el-radio>
          </el-radio-group>
        </div>
        <div class="mb-3">
          <label class="form-label"
            >积分数值 <span class="text-danger">*</span></label
          >
          <input
            type="number"
            class="form-control"
            v-model="scoreDialog.form.score"
            min="1"
            :class="{ 'is-invalid': scoreFormErrors.score }"
            required
          />
          <div class="invalid-feedback" v-if="scoreFormErrors.score">
            {{ scoreFormErrors.score }}
          </div>
        </div>
        <div class="mb-3">
          <label class="form-label"
            >备注 <span class="text-danger">*</span></label
          >
          <textarea
            class="form-control"
            v-model="scoreDialog.form.remark"
            rows="2"
            placeholder="可填写积分调整原因"
            :class="{ 'is-invalid': scoreFormErrors.remark }"
          ></textarea>
          <div class="invalid-feedback" v-if="scoreFormErrors.remark">
            {{ scoreFormErrors.remark }}
          </div>
        </div>
      </form>
    </model-dialog>
  </div>
</template>

<script setup>
  import ModelDialog from '@/components/ModelDialog.vue'
  import Pagination from '@/components/Pagination.vue'
  import { adminUploadFile, validateForm } from '@/js/store/common.js'
  import { httpGet, httpPost } from '@/js/utils/http'
  import { dateFormat } from '@/js/utils/libs'
  import { ElMessage, ElMessageBox } from 'element-plus'
  import { nextTick, onMounted, ref } from 'vue'

  // 变量定义
  const query = ref({ name: '', check: '' })
  const dataSets = ref({
    total: 0,
    page: 1,
    page_size: 10,
    items: [],
  })
  const loading = ref(true)

  // 审核对话框
  const checkDialog = ref({
    visible: false,
    creator: null,
    form: {
      check: 1,
      check_note: '审核通过',
    },
  })

  // 编辑对话框
  const editDialog = ref({
    visible: false,
    creator: null,
    form: {
      id: 0,
      name: '',
      description: '',
      logo: '',
      fee: 0,
    },
  })

  const checkFormRef = ref(null)
  const editFormRef = ref(null)

  // 校验规则
  const checkFormRules = {
    check_note: {
      required: (form) => form.check === 2,
      message: '请填写不通过原因',
    },
  }

  const editFormRules = {
    name: {
      required: true,
      message: '创作者名称不能为空',
    },
    fee: {
      required: true,
      validator: (value) => {
        if (value < 0 || value > 100) {
          editFormErrors.value.fee = '提现费率必须在0-100之间'
          return false
        }
        // 判断必须为整数
        if (!Number.isInteger(value)) {
          editFormErrors.value.fee = '提现费率必须是正整数'
          return false
        }
        return true
      },
    },
    logo: {
      required: true,
      message: '请上传创作者Logo',
    },
    description: {
      required: true,
      message: '创作者简介不能为空',
    },
  }

  const scoreFormRules = {
    score: {
      required: true,
      message: '请输入积分数值',
    },
    remark: {
      required: true,
      message: '请输入备注',
    },
  }
  const checkFormErrors = ref({})
  const editFormErrors = ref({})
  const scoreFormErrors = ref({})

  const scoreDialog = ref({
    visible: false,
    creator: null,
    form: {
      action: 'inc',
      score: 0,
      remark: '',
    },
    loading: false,
  })

  onMounted(() => {
    fetchData(1)
  })

  // 获取数据
  const fetchData = (page) => {
    loading.value = true
    const params = {
      page: page || dataSets.value.page,
      page_size: dataSets.value.page_size,
      name: query.value.name,
      check: query.value.check,
    }

    httpGet('/api/admin/creator/list', params)
      .then((res) => {
        if (res.data) {
          dataSets.value = res.data
        }
        loading.value = false
      })
      .catch((e) => {
        ElMessage.error('获取数据失败：' + e.message)
        loading.value = false
      })
  }

  // 显示审核对话框
  const showCheckDialog = (row) => {
    checkDialog.value.creator = row
    checkDialog.value.form.check = row.check || 1
    checkDialog.value.form.check_note = row.check_note
    checkDialog.value.visible = true
    // 清空校验错误
    checkFormErrors.value = {}
  }

  // 提交审核
  const submitCheck = async () => {
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
    checkDialog.value.form.validated = true
    httpPost(
      `/api/admin/creator/check?id=${checkDialog.value.creator.id}`,
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

  // 显示编辑对话框
  const showEditDialog = (row) => {
    editDialog.value.form = Object.assign({}, row)
    editDialog.value.visible = true
    // 清空校验错误
    editFormErrors.value = {}
  }

  // 提交编辑
  const submitEdit = async () => {
    await nextTick()
    // 校验
    if (
      !validateForm(editDialog.value.form, editFormRules, editFormErrors.value)
    ) {
      return
    }
    editDialog.value.form.validated = true
    httpPost('/api/admin/creator/update', editDialog.value.form)
      .then(() => {
        ElMessage.success('更新成功')
        editDialog.value.visible = false
        fetchData()
      })
      .catch((e) => {
        ElMessage.error('更新失败：' + e.message)
      })
  }

  // 启用/禁用切换
  const toggleEnable = (row) => {
    httpPost('/api/admin/creator/enable', {
      id: row.id,
      enabled: row.enabled,
    })
      .then(() => {
        ElMessage.success('操作成功')
      })
      .catch((e) => {
        ElMessage.error('操作失败：' + e.message)
        row.enabled = !row.enabled
      })
  }

  // 删除创作者
  const deleteCreator = (row) => {
    if (row.check === 1) {
      ElMessage.error('审核通过的创作者不能删除')
      return
    }
    // 弹窗确认
    ElMessageBox.confirm('删除创作者后，将无法恢复，确定删除吗？', '警告', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning',
    }).then(() => {
      httpGet('/api/admin/creator/delete', { id: row.id })
        .then(() => {
          ElMessage.success('删除成功')
          fetchData()
        })
        .catch((e) => {
          ElMessage.error('删除失败：' + e.message)
        })
    })
  }

  const uploadLogo = (file) => {
    adminUploadFile(file, (data) => {
      editDialog.value.form.logo = data.url
    })
  }

  function showScoreDialog(row) {
    scoreDialog.value.creator = row
    scoreDialog.value.form = { action: 'inc', score: 0, remark: '' }
    scoreDialog.value.visible = true
    scoreFormErrors.value = {}
  }

  async function submitScore() {
    const { creator, form } = scoreDialog.value

    await nextTick()
    // 校验
    if (
      !validateForm(
        scoreDialog.value.form,
        scoreFormRules,
        scoreFormErrors.value
      )
    ) {
      return
    }
    scoreDialog.value.form.validated = true
    scoreDialog.value.loading = true
    try {
      await httpPost('/api/admin/creator/score', {
        creator_id: creator.id,
        action: form.action,
        score: form.score,
        remark: form.remark,
      })
      ElMessage.success('操作成功')
      scoreDialog.value.visible = false
      fetchData()
    } catch (e) {
      ElMessage.error('操作失败：' + (e.message || e))
    } finally {
      scoreDialog.value.loading = false
    }
  }
</script>

<style scoped lang="scss">
  @use '@/assets/css/admin/admin.scss';
</style>
