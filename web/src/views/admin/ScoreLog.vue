<template>
  <div class="card">
    <div class="card-body">
      <div class="table-responsive" v-loading="loading">
        <div class="search-container d-flex flex-wrap align-items-end gap-2">
          <span class="search-item">
            <label class="form-label">用户名</label>
            <el-input
              v-model="query.username"
              placeholder="请输入用户名"
              style="width: 150px"
              clearable
              @change="fetchData(1)"
            />
          </span>
          <span class="search-item">
            <label class="form-label">类型</label>
            <el-select
              v-model="query.type"
              style="width: 120px"
              value-key="value"
              clearable
              @change="fetchData(1)"
            >
              <el-option
                v-for="item in scoreTypes"
                :key="item.value"
                :value="item.value"
                :label="item.label"
              />
            </el-select>
          </span>
          <span class="search-item">
            <label class="form-label">时间范围</label>
            <div class="d-flex gap-2">
              <el-date-picker
                v-model="query.start_time"
                type="date"
                placeholder="开始日期"
                format="YYYY-MM-DD"
                value-format="YYYY-MM-DD"
                style="width: 120px"
                clearable
                @change="fetchData(1)"
              />
              <el-date-picker
                v-model="query.end_time"
                type="date"
                placeholder="结束日期"
                format="YYYY-MM-DD"
                value-format="YYYY-MM-DD"
                style="width: 120px"
                clearable
                @change="fetchData(1)"
              />
            </div>
          </span>
          <span class="search-item">
            <el-button type="primary" @click="fetchData()">
              <i class="iconfont icon-search"></i> 搜索
            </el-button>
          </span>
          <span class="search-item">
            <el-button type="danger" @click="batchRemove">
              <i class="iconfont icon-remove me-1"></i> 删除
            </el-button>
          </span>
        </div>

        <el-table
          border
          class="data-table"
          :data="dataSets.items"
          table-layout="auto"
          @selection-change="handleSelectionChange"
        >
          <el-table-column type="selection" width="38" />
          <el-table-column prop="username" label="用户" />
          <el-table-column prop="subject" label="主题" />
          <el-table-column label="类型">
            <template #default="{ row }">
              <span
                class="badge rounded-pill"
                :class="getScoreTypeClass(row.type)"
              >
                {{ getScoreTypeLabel(row.type) }}
              </span>
            </template>
          </el-table-column>
          <el-table-column label="数额">
            <template #default="{ row }">
              <span
                :class="{
                  'text-success': row.mark === 1,
                  'text-danger': row.mark === 0,
                }"
              >
                {{ row.mark === 1 ? '+' : '-' }}{{ row.amount }}
              </span>
            </template>
          </el-table-column>
          <el-table-column prop="balance" label="余额" />
          <el-table-column label="发生时间">
            <template #default="{ row }">
              {{ dateFormat(row.created_at) }}
            </template>
          </el-table-column>
          <el-table-column prop="remark" label="备注" />
          <template #empty>
            <el-empty description="暂无数据" />
          </template>
        </el-table>

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
  </div>
</template>

<script setup>
  import { onMounted } from 'vue'
  import { storeToRefs } from 'pinia'
  import Pagination from '@/components/Pagination.vue'
  import { useAdminScoreStore } from '@/js/store/admin/score'
  import { dateFormat } from '@/js/utils/libs'

  const scoreStore = useAdminScoreStore()
  const { dataSets, loading, query } = storeToRefs(scoreStore)
  const {
    scoreTypes,
    fetchData,
    batchRemove,
    handleSelectionChange,
    getScoreTypeLabel,
    getScoreTypeClass,
  } = scoreStore

  onMounted(() => {
    if (!dataSets.value.items.length) {
      fetchData()
    }
  })
</script>

<style scoped lang="scss">
  @use '@/assets/css/admin/admin.scss';
</style>
