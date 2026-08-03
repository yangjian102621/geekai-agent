<template>
  <div class="card">
    <div class="card-body">
      <div class="table-responsive" v-loading="loading">
        <div class="search-container flex flex-wrap gap-2">
          <span class="search-item">
            <el-input
              type="text"
              v-model="query.username"
              placeholder="请输入用户名"
              clearable
            />
          </span>
          <span class="search-item">
            <el-input
              type="text"
              v-model="query.user_id"
              placeholder="请输入用户ID"
              clearable
            />
          </span>
          <span class="search-item">
            <el-button type="primary" @click="fetchData(1)">
              <i class="iconfont icon-search"></i>
            </el-button>
          </span>
          <span class="search-item">
            <el-button type="primary" @click="add">
              <i class="iconfont icon-plus me-1"></i> 新增
            </el-button>
          </span>
          <span class="search-item">
            <el-button type="danger" @click="multiRemove">
              <i class="iconfont icon-remove me-1"></i> 删除
            </el-button>
          </span>
        </div>

        <el-table
          :data="dataSets.items"
          border
          class="data-table"
          :row-key="(row) => row.id"
          @selection-change="handleSelectionChange"
          table-layout="auto"
        >
          <el-table-column type="selection" width="38"></el-table-column>
          <el-table-column prop="id" label="ID" />
          <el-table-column label="账号">
            <template #default="{ row }">
              {{ row.username }}
              <el-tooltip effect="dark" content="复制用户名" placement="top">
                <i
                  class="iconfont icon-copy text-primary ms-1 copy-btn"
                  style="cursor: pointer"
                  :data-clipboard-text="row.username"
                ></i>
              </el-tooltip>
              <el-tooltip
                effect="dark"
                placement="right"
                content="VIP会员"
                v-if="row.vip"
              >
                <i class="iconfont icon-vip-user text-warning"></i>
              </el-tooltip>
            </template>
          </el-table-column>
          <el-table-column prop="nickname" label="昵称" />
          <el-table-column prop="scores" label="剩余积分" />
          <el-table-column label="启用状态">
            <template #default="{ row }">
              <el-switch v-model="row.enabled" @change="enable(row)" />
            </template>
          </el-table-column>
          <el-table-column label="有效期">
            <template #default="{ row }">
              <span v-if="row.expired_time">{{ row.expired_time }}</span>
              <span v-else class="badge rounded-pill text-bg-primary"
                >长期有效</span
              >
            </template>
          </el-table-column>

          <el-table-column label="创建时间">
            <template #default="{ row }">
              <span>{{ dateFormat(row['created_at']) }}</span>
            </template>
          </el-table-column>

          <el-table-column label="操作" fixed="right">
            <template #default="{ row }">
              <el-dropdown placement="bottom" trigger="click">
                <button class="btn btn-primary btn-sm">
                  <i class="iconfont icon-more-vertical"></i>
                </button>
                <template #dropdown>
                  <el-dropdown-menu>
                    <el-dropdown-item @click="edit(row)">
                      <i class="iconfont icon-edit"></i> 编辑
                    </el-dropdown-item>
                    <el-dropdown-item @click="resetPass(row)">
                      <i class="iconfont icon-password"> 重置密码</i>
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

      <div class="pagination-primary d-flex justify-content-end pt-3">
        <Pagination
          :total="dataSets.total"
          :pageSize="dataSets.page_size"
          :currentPage="dataSets.page"
          @update:currentPage="fetchData"
          @update:pageSize="dataSets.pageSize = $event"
        />
      </div>
    </div>

    <model-dialog
      :modelValue="showDialog"
      :title="item.id > 0 ? '编辑用户' : '新增用户'"
      @cancel="showDialog = false"
      @confirm="handleSubmit"
      confirm-text="保存"
      cancel-text="关闭"
    >
      <form>
        <div class="mb-3">
          <label class="form-label"
            >用户名 <i class="iconfont icon-xinghao text-danger"></i
          ></label>
          <input
            type="text"
            class="form-control"
            v-model="item.username"
            :class="{ 'is-invalid': errors.username }"
            placeholder="请输入手机号或者邮箱地址"
          />
          <div class="invalid-feedback">{{ errors.username }}</div>
        </div>

        <div class="mb-3">
          <label class="form-label"
            >昵称 <i class="iconfont icon-xinghao text-danger"></i
          ></label>
          <div class="flex align-items-top">
            <div class="flex-1">
              <input
                type="text"
                class="form-control"
                v-model="item.nickname"
                :class="{ 'is-invalid': errors.nickname }"
              />
              <div class="invalid-feedback">{{ errors.nickname }}</div>
            </div>
            <el-tooltip effect="dark" placement="top" content="随机生成">
              <el-button
                type="primary"
                @click="generateNickname"
                :loading="btnLoading"
                style="height: 38px"
                class="ms-2"
              >
                <i class="iconfont icon-refresh"></i>
              </el-button>
            </el-tooltip>
          </div>
        </div>

        <div class="mb-3" v-if="!item.id">
          <label class="form-label"
            >密码 <i class="iconfont icon-xinghao text-danger"></i
          ></label>
          <input
            type="password"
            class="form-control"
            v-model="item.password"
            :class="{ 'is-invalid': errors.password }"
          />
          <div class="invalid-feedback">{{ errors.password }}</div>
        </div>

        <div class="mb-3">
          <label class="form-label">积分</label>
          <input
            type="text"
            class="form-control"
            v-model.number="item.scores"
          />
        </div>

        <div class="mb-3">
          <label class="form-label">有效期</label>
          <div>
            <el-date-picker
              v-model="item.expired_time"
              type="datetime"
              placeholder="选择日期"
              format="YYYY-MM-DD"
              value-format="YYYY-MM-DD"
              size="large"
            />
          </div>
        </div>

        <div class="mb-3">
          <label class="form-label">是否启用</label>
          <div class="form-label">
            <el-switch v-model="item.enabled" size="large" />
          </div>
        </div>

        <div class="mb-3">
          <label class="form-label">开通VIP</label>
          <div class="form-label">
            <el-switch v-model="item.vip" size="large" />
          </div>
        </div>
      </form>
    </model-dialog>

    <model-dialog
      :modelValue="showResetPassDialog"
      title="重置密码"
      @cancel="showResetPassDialog = false"
      @confirm="handleResetPass"
      confirm-text="提交"
      cancel-text="关闭"
    >
      <form>
        <div class="mb-3">
          <input
            type="password"
            class="form-control"
            v-model="item.password"
            :class="{ 'is-invalid': errors.password }"
            placeholder="请输入新密码"
          />
          <div class="invalid-feedback">{{ errors.password }}</div>
        </div>
      </form>
    </model-dialog>
  </div>
</template>

<script setup>
  import { onBeforeUnmount, onMounted, ref } from 'vue'
  import { storeToRefs } from 'pinia'
  import ModelDialog from '@/components/ModelDialog.vue'
  import Pagination from '@/components/Pagination.vue'
  import { useAdminUserStore } from '@/js/store/admin/user'
  import { dateFormat } from '@/js/utils/libs'
  import ClipboardJS from 'clipboard'
  import { ElMessage } from 'element-plus'
  import { httpGet } from '../../js/utils/http'

  const userStore = useAdminUserStore()
  const {
    loading,
    item,
    dataSets,
    errors,
    showDialog,
    showResetPassDialog,
    query,
  } = storeToRefs(userStore)
  const {
    handleSubmit,
    remove,
    enable,
    add,
    edit,
    handleResetPass,
    multiRemove,
    fetchData,
    handleSelectionChange,
    initialize,
  } = userStore

  const btnLoading = ref(false)

  let clipboard = null

  onMounted(async () => {
    await initialize()
    clipboard = new ClipboardJS('.copy-btn')
    clipboard.on('success', function (e) {
      ElMessage.success('用户名已复制')
      e.clearSelection()
    })
    clipboard.on('error', function () {
      ElMessage.error('复制失败，请手动复制')
    })
  })

  onBeforeUnmount(() => {
    if (clipboard) {
      clipboard.destroy()
      clipboard = null
    }
  })

  const resetPass = (row) => {
    item.value = { ...row }
    showResetPassDialog.value = true
  }
  const onSelectedDate = (dateStr) => {
    item.value.expired_time = dateStr
  }

  const generateNickname = () => {
    btnLoading.value = true
    httpGet('/api/admin/user/nickname')
      .then((res) => {
        item.value.nickname = res.data.nickname
      })
      .finally(() => {
        btnLoading.value = false
      })
  }
</script>

<style scoped lang="scss">
  @use '@/assets/css/admin/admin.scss';
</style>
