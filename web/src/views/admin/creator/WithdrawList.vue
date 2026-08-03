<template>
  <div class="card">
    <div class="card-body">
      <div class="table-responsive" v-loading="loading">
        <!-- 搜索筛选 -->
        <div class="search-container flex flex-wrap gap-2">
          <span class="search-item">
            <el-input
              v-model="query.creator_id"
              placeholder="创作者ID"
              style="width: 120px"
              clearable
              @clear="fetchData(1)"
            />
          </span>

          <span class="search-item">
            <el-select
              v-model="query.status"
              placeholder="结算状态"
              style="width: 100px"
              clearable
              @change="fetchData(1)"
            >
              <el-option label="待处理" value="pending" />
              <el-option label="已结算" value="success" />
              <el-option label="已拒绝" value="reject" />
            </el-select>
          </span>

          <span class="search-item">
            <el-select
              v-model="query.method"
              placeholder="提现方式"
              style="width: 100px"
              clearable
              @change="fetchData(1)"
            >
              <el-option label="支付宝" value="alipay" />
              <el-option label="微信" value="wxpay" />
            </el-select>
          </span>

          <span class="search-item">
            <el-date-picker
              v-model="query.dateRange"
              type="daterange"
              range-separator="至"
              start-placeholder="开始日期"
              end-placeholder="结束日期"
              style="width: 260px"
              value-format="YYYY-MM-DD"
              clearable
              @clear="fetchData(1)"
            />
          </span>

          <span class="search-item">
            <el-button type="primary" @click="fetchData()">
              <i class="iconfont icon-search"></i>
              搜索
            </el-button>
          </span>
        </div>

        <!-- 提现申请列表 -->
        <el-table
          border
          class="data-table"
          :data="dataSets.items"
          table-layout="auto"
        >
          <el-table-column label="创作者ID" width="50">
            <template #default="{ row }">
              <div class="text-muted small">{{ row.creator_id }}</div>
            </template>
          </el-table-column>

          <el-table-column label="提现积分" prop="scores" />
          <el-table-column label="提现金额" prop="total_money" />
          <el-table-column label="手续费" prop="fee" />
          <el-table-column label="实际到账" prop="real_money" />
          <el-table-column label="提现方式" width="120">
            <template #default="{ row }">
              <el-tag
                type="primary"
                size="small"
                v-if="row.method === 'alipay'"
              >
                支付宝
              </el-tag>
              <el-tag type="success" size="small" v-else> 微信 </el-tag>
            </template>
          </el-table-column>
          <el-table-column label="收款人" prop="account_name" />
          <el-table-column label="收款账号" prop="account" />
          <el-table-column label="收款二维码" prop="qr_code">
            <template #default="{ row }">
              <el-image
                :src="row.qr_code"
                :preview-src-list="[row.qr_code]"
                style="width: 50px; height: 50px"
                :preview-teleported="true"
              />
            </template>
          </el-table-column>

          <el-table-column label="状态" width="120">
            <template #default="{ row }">
              <el-tag v-if="row.status === 'success'" type="success"
                >已结算</el-tag
              >
              <el-tag v-else-if="row.status === 'reject'" type="danger"
                >已拒绝</el-tag
              >
              <el-tag v-else type="warning">待处理</el-tag>
            </template>
          </el-table-column>

          <el-table-column label="申请时间" width="180">
            <template #default="{ row }">
              {{ dateFormat(row.created_at) }}
            </template>
          </el-table-column>

          <el-table-column label="备注" width="200">
            <template #default="{ row }">
              <span class="text-muted">{{ row.note || '-' }}</span>
            </template>
          </el-table-column>

          <el-table-column label="操作" width="160" fixed="right">
            <template #default="{ row }">
              <el-button
                v-if="row.status === 'pending'"
                type="primary"
                @click="preccessWithdraw(row)"
              >
                结算
              </el-button>
              <el-tag v-else type="success" size="large"> 已处理 </el-tag>
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

    <!-- 结算对话框 -->
    <model-dialog
      :modelValue="auditDialog.visible"
      title="结算处理"
      @cancel="auditDialog.visible = false"
      @confirm="submitWithdraw"
      confirm-text="确定"
      cancel-text="取消"
    >
      <form>
        <div class="mb-3">
          <label class="form-label">积分数量</label>
          <input
            type="text"
            class="form-control"
            v-model="auditDialog.form.scores"
            disabled
          />
        </div>

        <div class="mb-3">
          <label class="form-label">提现金额</label>
          <input
            type="text"
            class="form-control"
            v-model="auditDialog.form.real_money"
            disabled
          />
        </div>

        <div class="mb-3">
          <label class="form-label">结算状态</label>
          <div>
            <el-radio-group
              v-model="auditDialog.form.status"
              @change="changeStatus"
            >
              <el-radio label="success" value="success" size="large" border
                >已结算</el-radio
              >
              <el-radio label="reject" value="reject" size="large" border
                >拒绝</el-radio
              >
            </el-radio-group>
          </div>
        </div>

        <div class="mb-3">
          <label class="form-label">备注</label>
          <textarea
            class="form-control"
            v-model="auditDialog.form.note"
            rows="3"
          />
        </div>
      </form>
    </model-dialog>
  </div>
</template>

<script setup>
  import ModelDialog from '@/components/ModelDialog.vue'
  import Pagination from '@/components/Pagination.vue'
  import { closeLoading, showLoading } from '@/js/utils/dialog'
  import { httpPost } from '@/js/utils/http'
  import { dateFormat } from '@/js/utils/libs'
  import { ElMessage } from 'element-plus'
  import { onMounted, ref } from 'vue'

  // 变量定义
  const query = ref({
    creator_id: '',
    status: '',
    type: '',
    dateRange: [],
  })
  const dataSets = ref({
    total: 0,
    page: 1,
    page_size: 10,
    items: [],
  })
  const loading = ref(true)

  // 结算对话框
  const auditDialog = ref({
    visible: false,
    form: {
      note: '',
      status: 'success',
    },
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
      creator_id: query.value.creator_id,
      status: query.value.status,
      method: query.value.method,
      start_time:
        query.value.dateRange && query.value.dateRange[0]
          ? query.value.dateRange[0]
          : '',
      end_time:
        query.value.dateRange && query.value.dateRange[1]
          ? query.value.dateRange[1]
          : '',
    }

    httpPost('/api/admin/creator/withdraws', params)
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

  // 通过提现申请
  const preccessWithdraw = (row) => {
    auditDialog.value = {
      visible: true,
      form: {
        ...row,
        status: 'success',
        note: '审核通过，予以结算',
      },
    }
  }

  // 提交审核
  const submitWithdraw = async () => {
    if (
      auditDialog.value.form.status === 'reject' &&
      !auditDialog.value.form.note
    ) {
      ElMessage.error('请填写拒绝原因')
      return
    }

    showLoading()
    httpPost('/api/admin/creator/withdraws/proccess', {
      id: auditDialog.value.form.id,
      status: auditDialog.value.form.status,
      note: auditDialog.value.form.note,
    })
      .then(() => {
        ElMessage.success('操作成功')
        fetchData(1)
        auditDialog.value.visible = false
      })
      .catch((e) => {
        ElMessage.error('操作失败：' + e.message)
      })
      .finally(() => {
        closeLoading()
      })
  }

  const changeStatus = (status) => {
    if (status === 'reject') {
      auditDialog.value.form.note = '审核失败，终止结算'
    } else {
      auditDialog.value.form.note = '审核通过，予以结算'
    }
  }
</script>

<style scoped lang="scss">
  @use '@/assets/css/admin/admin.scss';
</style>
