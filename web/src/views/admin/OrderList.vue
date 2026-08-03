<template>
  <div class="card">
    <div class="card-body">
      <div class="table-responsive" v-loading="loading">
        <div class="search-container d-flex flex-wrap align-items-end gap-2">
          <span class="search-item">
            <label class="form-label">订单号</label>
            <el-input
              v-model="query.order_no"
              placeholder="请输入订单号"
              style="width: 150px"
              clearable
              @keyup.enter="fetchData(1)"
            />
          </span>
          <span class="search-item">
            <label class="form-label">订单状态</label>
            <el-select
              v-model="query.status"
              style="width: 120px"
              value-key="value"
              clearable
              @keyup.enter="fetchData(1)"
            >
              <el-option
                v-for="item in orderStatus"
                :key="item.value"
                :value="item.value"
                :label="item.label"
              />
            </el-select>
          </span>
          <span class="search-item">
            <label class="form-label">时间范围</label>
            <el-date-picker
              v-model="query.pay_time"
              type="daterange"
              start-placeholder="开始日期"
              end-placeholder="结束日期"
              format="YYYY-MM-DD"
              value-format="YYYY-MM-DD"
              style="width: 250px"
              @keyup.enter="fetchData(1)"
            />
          </span>
          <span class="search-item">
            <el-button type="primary" @click="fetchData(1)">
              <i class="iconfont icon-search"></i> 搜索
            </el-button>
          </span>
          <span class="search-item">
            <el-button type="danger" @click="clearOrders">
              <i class="iconfont icon-remove me-1"></i> 清空未支付订单
            </el-button>
          </span>
        </div>

        <el-table
          border
          class="data-table"
          :data="dataSets.items"
          table-layout="auto"
          style="width: 100%"
        >
          <el-table-column prop="order_no" label="订单号" min-width="120" />
          <el-table-column prop="trade_no" label="交易号" min-width="120" />
          <el-table-column prop="username" label="下单用户" min-width="100" />
          <el-table-column prop="subject" label="产品名称" min-width="150" />
          <el-table-column prop="amount" label="订单金额" min-width="100" />
          <el-table-column label="充值积分" min-width="100">
            <template #default="{ row }">
              <span>{{ row.remark && row.remark.credit }}</span>
            </template>
          </el-table-column>
          <el-table-column label="下单时间" min-width="120">
            <template #default="{ row }">
              {{ dateFormat(row.created_at) }}
            </template>
          </el-table-column>
          <el-table-column label="支付时间" min-width="120">
            <template #default="{ row }">
              <span v-if="row.pay_time">{{ dateFormat(row.pay_time) }}</span>
              <el-tag v-else>未支付</el-tag>
            </template>
          </el-table-column>
          <el-table-column label="订单状态" min-width="100">
            <template #default="{ row }">
              <el-tag v-if="row.status === 1" type="success">已支付</el-tag>
              <el-tag v-else-if="row.status === 2" type="danger">已关闭</el-tag>
              <el-tag v-else type="warning">未支付</el-tag>
            </template>
          </el-table-column>
          <el-table-column
            prop="channel_name"
            label="支付渠道"
            min-width="100"
          />
          <el-table-column prop="pay_name" label="支付名称" min-width="100" />
          <el-table-column label="操作" width="150" fixed="right">
            <template #default="{ row }">
              <el-dropdown placement="bottom" trigger="click">
                <button class="btn btn-primary btn-sm">
                  <i class="iconfont icon-more-vertical"></i>
                </button>
                <template #dropdown>
                  <el-dropdown-menu>
                    <el-dropdown-item @click="remove(row)">
                      <span class="text-danger">
                        <i class="iconfont icon-remove"></i> 删除
                      </span>
                    </el-dropdown-item>
                    <el-dropdown-item
                      v-if="row.status === 0"
                      @click="recharge(row)"
                    >
                      <i class="iconfont icon-edit"></i> 补单
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
  import { httpGet, httpPost } from '@/js/utils/http'
  import { dateFormat } from '@/js/utils/libs'
  import { ElMessage, ElMessageBox } from 'element-plus'
  import { onMounted, ref } from 'vue'

  // 变量定义
  const query = ref({ order_no: '', pay_time: [], status: -1 })
  const dataSets = ref({
    total: 0,
    page: 1,
    page_size: 20,
  })
  const loading = ref(true)
  const orderStatus = ref([
    { value: -1, label: '全部' },
    { value: 0, label: '未支付' },
    { value: 1, label: '已支付' },
    { value: 2, label: '已关闭' },
  ])

  onMounted(() => {
    fetchData(1)
  })

  // 获取数据
  const fetchData = (page) => {
    loading.value = true
    query.value.page = page
    query.value.page_size = dataSets.value.page_size
    httpPost('/api/admin/order/list', query.value)
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

  const remove = function (row) {
    ElMessageBox.confirm('确定要删除该订单吗?', '删除提示', {
      confirmButtonText: '确认',
      cancelButtonText: '取消',
      type: 'warning',
    }).then(() => {
      httpGet('/api/admin/order/remove?id=' + row.id)
        .then(() => {
          ElMessage.success('删除成功！')
          fetchData(1)
        })
        .catch((e) => {
          ElMessage.error('删除失败：' + e.message)
        })
    })
  }

  const clearOrders = () => {
    ElMessageBox.confirm(
      '此操作将会删除所有未支付订单，继续操作吗?',
      '删除提示',
      {
        confirmButtonText: '确认',
        cancelButtonText: '取消',
        type: 'warning',
      }
    ).then(() => {
      httpGet('/api/admin/order/clear').then(() => {
        ElMessage.success('订单删除成功')
        dataSets.value.page = 1
        fetchData(1)
      })
    })
  }

  const recharge = (row) => {
    httpGet('/api/admin/order/recharge?id=' + row.id)
      .then(() => {
        ElMessage.success('补单成功,请稍后刷新页面查看结果！')
        fetchData(1)
      })
      .catch((e) => {
        ElMessage.error('补单失败：' + e.message)
      })
  }
</script>

<style scoped lang="scss">
  @use '@/assets/css/admin/admin.scss';
</style>
