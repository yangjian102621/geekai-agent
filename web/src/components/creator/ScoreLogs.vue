<template>
  <div class="score-logs">
    <div class="list-header">
      <div class="mb-3">
        <el-select
          v-model="query.type"
          placeholder="类型"
          clearable
          style="width: 120px; margin-right: 8px"
          @change="fetchData(1)"
          @clear="fetchData(1)"
        >
          <el-option label="收入" value="income" />
          <el-option label="提现" value="withdraw" />
          <el-option label="退款" value="refund" />
        </el-select>
        <el-date-picker
          v-model="query.date_range"
          type="daterange"
          range-separator="至"
          start-placeholder="开始日期"
          end-placeholder="结束日期"
          value-format="YYYY-MM-DD"
          style="width: 240px; margin-right: 8px"
          @change="fetchData(1)"
        />
        <el-button type="primary" @click="fetchData(1)">
          <i class="iconfont icon-search"></i>
        </el-button>
      </div>
    </div>

    <el-table
      :data="data.items"
      v-loading="loading"
      border
      class="earnings-table"
    >
      <el-table-column label="名称" width="120" prop="subject" />
      <el-table-column label="积分" width="100" prop="score">
        <template #default="{ row }">
          <span class="text-success" v-if="row.mark === 1">
            + {{ row.score }}
          </span>
          <span class="text-danger" v-else> - {{ row.score }} </span>
        </template>
      </el-table-column>
      <el-table-column label="余额" width="100" prop="balance" />
      <el-table-column label="收支类型" width="100" prop="type">
        <template #default="{ row }">
          <el-tag :type="getTypeColor(row.type)" size="small">
            {{ getScoreType(row.type) }}
          </el-tag>
        </template>
      </el-table-column>
      <el-table-column prop="created_at" label="时间" width="160">
        <template #default="{ row }">
          {{ dateFormat(row.created_at * 1000) }}
        </template>
      </el-table-column>
      <el-table-column label="备注" min-width="120">
        <template #default="{ row }">
          <span class="remark">{{ row.remark || '-' }}</span>
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
    page_size: 20,
    stats: null,
  })

  const loading = ref(false)

  const query = reactive({
    type: '',
    date_range: [],
  })

  onMounted(() => {
    fetchData(1)
  })

  const fetchData = async (page = 1) => {
    loading.value = true
    try {
      const params = {
        page,
        page_size: data.value.page_size,
        type: query.type,
        start_date:
          query.date_range && query.date_range.length > 0
            ? query.date_range[0]
            : '',
        end_date:
          query.date_range && query.date_range.length > 1
            ? query.date_range[1]
            : '',
      }
      const res = await httpPost('/api/creator/scores/list', params)
      data.value = res.data
    } catch (error) {
      ElMessage.error('获取收益明细失败：' + error.message)
    } finally {
      loading.value = false
    }
  }

  const getTypeColor = (type) => {
    if (type === 'income' || type === 'refund') {
      return 'success'
    } else if (type === 'withdraw') {
      return 'warning'
    }
    return 'info'
  }

  const getScoreType = (type) => {
    switch (type) {
      case 'income':
        return '收入'
      case 'withdraw':
        return '提现'
      case 'refund':
        return '退款'
      case 'fine_tune':
        return '系统调整'
    }
  }

  defineExpose({
    fetchData,
  })
</script>

<style lang="scss">
  .score-logs {
    .el-table {
      .cell {
        font-size: 14px;
      }
    }
    .text-success {
      color: #21ba45;
    }
    .text-danger {
      color: #f56c6c;
    }
  }
</style>
