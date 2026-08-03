<template>
  <div class="card">
    <div class="card-body">
      <div class="table-responsive" v-loading="loading">
        <div class="search-container flex gap-2">
          <span class="search-item">
            <button class="btn btn-primary btn-sm" @click="add">
              <i class="iconfont icon-plus"></i> 新增
            </button>
          </span>
        </div>

        <el-table
          :data="items"
          border
          class="data-table"
          :row-key="(row) => row.id"
          @selection-change="handleSelectionChange"
          table-layout="auto"
        >
          <el-table-column type="selection" width="38"></el-table-column>
          <el-table-column label="账号">
            <template #default="scope">
              {{ scope.row.username }}
              <el-tooltip
                effect="dark"
                placement="right"
                content="VIP会员"
                v-if="scope.row.vip"
              >
                <i class="iconfont icon-vip-user text-warning"></i>
              </el-tooltip>
            </template>
          </el-table-column>
          <el-table-column prop="last_login_ip" label="最后登录IP" />
          <el-table-column prop="last_login_at" label="最后登录时间">
            <template #default="scope">
              <span>{{ dateFormat(scope.row.last_login_at) }}</span>
            </template>
          </el-table-column>
          <el-table-column label="启用状态">
            <template #default="scope">
              <el-switch
                v-model="scope.row.status"
                @change="enable(scope.row)"
              />
            </template>
          </el-table-column>

          <el-table-column label="创建时间">
            <template #default="scope">
              <span>{{ dateFormat(scope.row['created_at']) }}</span>
            </template>
          </el-table-column>

          <el-table-column label="操作" fixed="right">
            <template #default="scope">
              <el-dropdown placement="bottom" trigger="click">
                <button class="btn btn-primary btn-sm">
                  <i class="iconfont icon-more-vertical"></i>
                </button>
                <template #dropdown>
                  <el-dropdown-menu>
                    <el-dropdown-item @click="resetPass(scope.row)">
                      <i class="iconfont icon-password"> 重置密码</i>
                    </el-dropdown-item>
                    <el-dropdown-item @click="remove(scope.row)">
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
    </div>

    <model-dialog
      :modelValue="showDialog"
      title="新增管理员"
      @cancel="showDialog = false"
      @confirm="handleSubmit"
      confirm-text="保存"
      cancel-text="关闭"
      :width="400"
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
          />
          <div class="invalid-feedback">{{ errors.username }}</div>
        </div>

        <div class="mb-3">
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
          />
          <div class="invalid-feedback">{{ errors.password }}</div>
        </div>
      </form>
    </model-dialog>
  </div>
</template>

<script setup>
  import { onMounted } from 'vue'
  import { storeToRefs } from 'pinia'
  import ModelDialog from '@/components/ModelDialog.vue'
  import { useAdminManagerStore } from '@/js/store/admin/manager'
  import { dateFormat } from '@/js/utils/libs'

  const managerStore = useAdminManagerStore()
  const { loading, item, items, errors, showDialog, showResetPassDialog } =
    storeToRefs(managerStore)
  const {
    handleSubmit,
    remove,
    enable,
    add,
    handleResetPass,
    handleSelectionChange,
    initialize,
  } = managerStore

  const resetPass = (row) => {
    item.value = row
    showResetPassDialog.value = true
  }

  onMounted(() => {
    initialize()
  })
</script>

<style scoped lang="scss">
  @use '@/assets/css/admin/admin.scss';
</style>
