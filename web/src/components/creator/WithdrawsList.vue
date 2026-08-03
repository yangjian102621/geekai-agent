<template>
  <div class="withdraws-list" v-loading="loading">
    <div class="mb-3">
      <div class="flex items-center gap-2">
        <el-select
          v-model="query.status"
          placeholder="结算状态"
          clearable
          style="width: 100px"
          @change="fetchData(1)"
        >
          <el-option label="待处理" value="pending" />
          <el-option label="已结算" value="success" />
          <el-option label="已拒绝" value="rejected" />
        </el-select>
        <div class="flex">
          <el-date-picker
            v-model="query.date_range"
            type="daterange"
            range-separator="至"
            start-placeholder="开始日期"
            end-placeholder="结束日期"
            value-format="YYYY-MM-DD"
            @change="fetchData(1)"
          />
        </div>
        <el-button type="primary" @click="fetchData(1)">
          <i class="iconfont icon-search"></i>
        </el-button>
      </div>
    </div>

    <el-table
      :data="data.items"
      v-loading="loading"
      border
      class="withdraws-table"
    >
      <el-table-column label="提现积分" prop="scores" width="100" />
      <el-table-column label="提现金额" prop="total_money" width="100" />
      <el-table-column label="手续费" prop="fee" width="100" />
      <el-table-column label="实际到账" prop="real_money" width="100" />
      <el-table-column label="提现方式" width="100">
        <template #default="{ row }">
          <el-tag type="primary" size="small" v-if="row.method === 'alipay'">
            支付宝
          </el-tag>
          <el-tag type="success" size="small" v-else> 微信 </el-tag>
        </template>
      </el-table-column>
      <el-table-column label="收款人" prop="account_name" width="100" />
      <el-table-column label="收款账号" prop="account" width="100" />
      <el-table-column label="收款二维码" prop="qr_code" width="100">
        <template #default="{ row }">
          <el-image
            :src="row.qr_code"
            :preview-src-list="[row.qr_code]"
            style="width: 50px; height: 50px"
            :preview-teleported="true"
          />
        </template>
      </el-table-column>

      <el-table-column label="状态" width="100">
        <template #default="{ row }">
          <el-tag v-if="row.status === 'success'" type="success">已结算</el-tag>
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

      <el-table-column label="备注">
        <template #default="{ row }">
          <span class="text-muted">{{ row.note || '-' }}</span>
        </template>
      </el-table-column>
    </el-table>

    <!-- 分页 -->
    <div class="pagination" v-if="data.total > 0">
      <Pagination
        :total="data.total"
        :pageSize="data.page_size"
        :currentPage="data.page"
        @update:currentPage="fetchData"
        @update:pageSize="data.page_size = $event"
      />
    </div>
  </div>
</template>

<script setup>
  import Pagination from '@/components/Pagination.vue'
  import { httpPost } from '@/js/utils/http'
  import { dateFormat } from '@/js/utils/libs'
  import { ElMessage } from 'element-plus'
  import { onMounted, reactive, ref } from 'vue'

  const data = ref({
    items: [],
    total: 0,
    page: 1,
    page_size: 10,
    stats: null,
  })

  const loading = ref(false)

  const query = reactive({
    status: '',
    date_range: [],
  })

  onMounted(() => {
    fetchData(1)
  })

  const fetchData = async (page = 1) => {
    loading.value = true
    try {
      const res = await httpPost('/api/creator/withdraws/list', {
        page,
        page_size: data.value.page_size,
        status: query.status,
        start_date: query.date_range ? query.date_range[0] : '',
        end_date: query.date_range ? query.date_range[1] : '',
      })
      data.value = res.data
    } catch (error) {
      ElMessage.error('获取提现记录失败：' + error.message)
    } finally {
      loading.value = false
    }
  }

  defineExpose({
    fetchData,
  })
</script>

<style lang="scss">
  .withdraws-list {
    font-size: 14px;

    .el-table {
      .cell {
        font-size: 14px;
      }
    }
  }
</style>
