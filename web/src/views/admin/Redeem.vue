<template>
  <div class="card">
    <div class="card-body">
      <div class="table-responsive" v-loading="loading">
        <div class="search-container d-flex flex-wrap align-items-end gap-2">
          <span class="search-item">
            <el-input
              type="text"
              v-model="query.code"
              placeholder="请输入兑换码"
              clearable
            />
          </span>
          <span class="search-item">
            <el-select
              v-model="query.status"
              placeholder="请选择状态"
              style="width: 100px"
              clearable
            >
              <el-option
                v-for="item in redeemStatus"
                :key="item.value"
                :label="item.label"
                :value="item.value"
              />
            </el-select>
          </span>
          <span class="search-item">
            <el-button type="primary" @click="fetchData">
              <i class="iconfont icon-search"></i>
            </el-button>
          </span>
          <span class="search-item">
            <el-button type="primary" @click="add">
              <i class="iconfont icon-plus me-1"></i> 新增
            </el-button>
          </span>
          <span class="search-item">
            <el-button type="danger" @click="batchRemove">
              <i class="iconfont icon-remove me-1"></i> 删除
            </el-button>
          </span>
          <span class="search-item">
            <el-button
              type="success"
              @click="exportItems"
              :disabled="exporting"
            >
              <i class="iconfont icon-export me-1"></i> 导出
            </el-button>
          </span>
        </div>

        <el-table
          border
          class="data-table"
          :data="dataSets.items"
          @selection-change="handleSelectionChange"
          table-layout="auto"
        >
          <el-table-column type="selection" width="38" />
          <el-table-column prop="name" label="名称" />
          <el-table-column label="兑换码">
            <template #default="{ row }">
              <span>{{ substr(row.code, 24) }}</span>
              <i
                class="iconfont icon-copy ms-1"
                :data-clipboard-text="row.code"
              ></i>
            </template>
          </el-table-column>
          <el-table-column label="兑换人">
            <template #default="{ row }">
              <span v-if="row.username">{{ row.username }}</span>
              <span v-else class="badge bg-secondary">未兑换</span>
            </template>
          </el-table-column>
          <el-table-column prop="amount" label="额度" />
          <el-table-column label="生成时间">
            <template #default="{ row }">
              {{ dateFormat(row.created_at) }}
            </template>
          </el-table-column>
          <el-table-column label="兑换时间">
            <template #default="{ row }">
              <span v-if="row.redeemed_at > 0">{{
                dateFormat(row.redeemed_at)
              }}</span>
              <span v-else class="badge bg-secondary">未兑换</span>
            </template>
          </el-table-column>
          <el-table-column label="启用状态">
            <template #default="{ row }">
              <el-switch
                v-model="row.enabled"
                @change="set('enabled', row)"
                :disabled="row.redeemed_at > 0"
              />
            </template>
          </el-table-column>
          <el-table-column label="操作" fixed="right">
            <template #default="{ row }">
              <el-dropdown>
                <button class="btn btn-primary btn-sm">
                  <i class="iconfont icon-more-vertical"></i>
                </button>
                <template #dropdown>
                  <el-dropdown-menu>
                    <el-dropdown-item @click="remove(row)" class="text-danger">
                      <i class="iconfont icon-remove"></i> 删除
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

    <model-dialog
      :modelValue="showDialog"
      :title="title"
      @cancel="showDialog = false"
      @confirm="handleSubmit"
      confirm-text="立即生成"
      cancel-text="关闭"
    >
      <form>
        <div class="mb-3">
          <label class="form-label"
            >名称 <i class="iconfont icon-xinghao text-danger"></i
          ></label>
          <input
            type="text"
            class="form-control"
            v-model="item.name"
            :class="{ 'is-invalid': errors.name }"
          />
          <div class="invalid-feedback">{{ errors.name }}</div>
        </div>

        <div class="mb-3">
          <label class="form-label"
            >兑换额度 <i class="iconfont icon-xinghao text-danger"></i
          ></label>
          <input
            type="number"
            class="form-control"
            v-model.number="item.amount"
            :class="{ 'is-invalid': errors.amount }"
          />
          <div class="invalid-feedback">{{ errors.amount }}</div>
        </div>

        <div class="mb-3">
          <label class="form-label"
            >生成数量 <i class="iconfont icon-xinghao text-danger"></i
          ></label>
          <input
            type="number"
            class="form-control"
            v-model.number="item.num"
            :class="{ 'is-invalid': errors.num }"
          />
          <div class="invalid-feedback">{{ errors.num }}</div>
        </div>
      </form>
    </model-dialog>
  </div>
</template>

<script setup>
  import { onMounted, onUnmounted } from 'vue'
  import { storeToRefs } from 'pinia'
  import ModelDialog from '@/components/ModelDialog.vue'
  import Pagination from '@/components/Pagination.vue'
  import { useAdminRedeemStore } from '@/js/store/admin/redeem'
  import { dateFormat, substr } from '@/js/utils/libs'

  const redeemStore = useAdminRedeemStore()
  const {
    dataSets,
    loading,
    query,
    showDialog,
    item,
    itemIds,
    exporting,
    errors,
    title,
  } = storeToRefs(redeemStore)
  const {
    redeemStatus,
    add,
    handleSubmit,
    set,
    fetchData,
    remove,
    handleSelectionChange,
    exportItems,
    batchRemove,
    initialize,
    releaseClipboard,
  } = redeemStore

  onMounted(() => {
    initialize()
  })

  onUnmounted(() => {
    releaseClipboard()
  })
</script>

<style scoped lang="scss">
  @use '@/assets/css/admin/admin.scss';
</style>
