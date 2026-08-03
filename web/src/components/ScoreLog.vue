<template>
  <div class="flex pt-3 pb-3">
    <div class="card-body">
      <div class="table-responsive" v-loading="loading">
        <div class="search-container flex gap-2 flex-wrap align-items-end">
          <span class="search-item">
            <el-select
              v-model="query.type"
              style="width: 120px"
              value-key="value"
              clearable
              placeholder="类型"
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
            <div class="d-flex gap-2">
              <el-date-picker
                v-model="query.start_time"
                type="date"
                placeholder="开始日期"
                format="YYYY-MM-DD"
                value-format="YYYY-MM-DD"
                style="width: 120px"
                clearable
              />
              <el-date-picker
                v-model="query.end_time"
                type="date"
                placeholder="结束日期"
                format="YYYY-MM-DD"
                value-format="YYYY-MM-DD"
                style="width: 120px"
                clearable
              />
            </div>
          </span>
          <span class="search-item">
            <el-button type="primary" @click="fetchData()">
              <i class="iconfont icon-search"></i>
            </el-button>
          </span>

          <span class="search-item">
            <el-tooltip content="刷新" placement="top">
              <el-button type="success" @click="fetchData()">
                <i class="iconfont icon-refresh"></i>
              </el-button>
            </el-tooltip>
          </span>
        </div>

        <el-table
          border
          class="data-table"
          :data="dataSets.items"
          table-layout="auto"
        >
          <el-table-column prop="subject" label="主题" width="100" />
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
          <el-table-column label="发生时间" width="175">
            <template #default="{ row }">
              {{ dateFormat(row.created_at) }}
            </template>
          </el-table-column>
          <el-table-column prop="remark" label="备注" />
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
  import { useFrontScoreStore } from '@/js/store/front/score'
  import { dateFormat } from '@/js/utils/libs'

  const scoreStore = useFrontScoreStore()
  const { dataSets, loading, query } = storeToRefs(scoreStore)
  const { scoreTypes, fetchData, getScoreTypeLabel, getScoreTypeClass } =
    scoreStore

  onMounted(() => {
    if (!dataSets.value.items.length) {
      fetchData()
    }
  })
</script>

<style scoped lang="scss">
  @use '../assets/css/admin/admin.scss';
</style>
