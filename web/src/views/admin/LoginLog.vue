<template>
  <div class="card">
    <div class="card-body">
      <div class="table-responsive" v-loading="loading">
        <el-table
          :data="dataSets.items"
          border
          class="data-table"
          table-layout="auto"
        >
          <el-table-column prop="username" label="用户名" />
          <el-table-column prop="nickname" label="昵称" />
          <el-table-column prop="login_ip" label="登录IP" />
          <el-table-column prop="login_address" label="登录地址" />
          <el-table-column label="登录时间">
            <template #default="{ row }">
              {{ dateFormat(row.created_at) }}
            </template>
          </el-table-column>
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
  import Pagination from '@/components/Pagination.vue'
  import { showMessageError } from '@/js/utils/dialog'
  import { httpGet } from '@/js/utils/http'
  import { dateFormat } from '@/js/utils/libs'
  import { onMounted, ref } from 'vue'

  const dataSets = ref({ page: 1, page_size: 20, total: 0, items: [] })
  const loading = ref(true)

  onMounted(() => {
    fetchData(1)
  })

  const fetchData = (page) => {
    loading.value = true
    httpGet('/api/admin/user/loginLog', {
      page: page,
      page_size: dataSets.value.page_size,
    })
      .then((res) => {
        if (res.data) {
          dataSets.value = res.data
        }
      })
      .catch((e) => {
        showMessageError('获取数据失败：' + e.message)
      })
      .finally(() => {
        loading.value = false
      })
  }
</script>

<style scoped lang="scss">
  @use '@/assets/css/admin/admin.scss';
</style>
