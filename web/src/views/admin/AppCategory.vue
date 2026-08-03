<template>
  <div class="card">
    <div class="card-body">
      <div class="table-responsive" v-loading="loading">
        <div class="search-container flex flex-wrap gap-2">
          <span class="search-item">
            <el-button type="primary" @click="add">
              <i class="iconfont icon-plus me-1"></i> 新增
            </el-button>
          </span>
        </div>

        <el-table
          :data="items"
          border
          class="data-table"
          :row-key="(row) => row.id"
          table-layout="auto"
        >
          <el-table-column label="分类名称" prop="name" />
          <el-table-column label="启用状态">
            <template #default="scope">
              <el-switch
                v-model="scope.row.enabled"
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
                    <el-dropdown-item @click="edit(scope.row)">
                      <span class="text-primary">
                        <i class="iconfont icon-edit"></i> 编辑
                      </span>
                    </el-dropdown-item>

                    <el-dropdown-item>
                      <span class="text-danger" @click="remove(scope.row)">
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
      title="新增/编辑分类"
      @cancel="showDialog = false"
      @confirm="handleSubmit"
      confirm-text="保存"
      cancel-text="关闭"
      :width="400"
    >
      <form>
        <div class="mb-3">
          <label class="form-label"
            >分类名称 <i class="iconfont icon-xinghao text-danger"></i
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
          <label class="form-label">是否启用 </label>
          <div>
            <el-switch v-model="item.enabled" />
          </div>
        </div>
      </form>
    </model-dialog>
  </div>
</template>

<script setup>
  import ModelDialog from '@/components/ModelDialog.vue'
  import {
    showConfirm,
    showMessageError,
    showMessageOK,
  } from '@/js/utils/dialog'
  import { httpGet, httpPost } from '@/js/utils/http'
  import { copyObj, dateFormat } from '@/js/utils/libs'

  const loading = ref(false)
  const item = ref({})
  const items = ref([])
  const errors = ref({})
  const showDialog = ref(false)

  onMounted(() => {
    fetchData()
  })

  const fetchData = () => {
    httpGet('/api/admin/app/category/list', { system: true }).then((res) => {
      items.value = res.data
    })
  }

  const add = () => {
    showDialog.value = true
    item.value = {
      enabled: true,
    }
  }

  const edit = (row) => {
    item.value = copyObj(row)
    showDialog.value = true
  }

  const remove = (row) => {
    showConfirm('删除提示', '确定要删除当前记录吗?', () => {
      httpGet('/api/admin/app/category/remove', {
        id: row.id,
      })
        .then(() => {
          showMessageOK('删除成功')
          fetchData()
        })
        .catch((err) => {
          showMessageError(err.message)
        })
    })
  }

  const enable = (row) => {
    httpPost('/api/admin/app/category/enable', {
      id: row.id,
      enabled: row.enabled,
    })
      .then(() => {
        showMessageOK('操作成功')
        fetchData()
      })
      .catch((err) => {
        showMessageError(err.message)
      })
  }

  const handleSubmit = () => {
    httpPost('/api/admin/app/category/save', item.value)
      .then(() => {
        showMessageOK('操作成功')
        fetchData()
        showDialog.value = false
      })
      .catch((e) => {
        showMessageError(e.message)
      })
  }
</script>

<style scoped lang="scss">
  @use '@/assets/css/admin/admin.scss';
</style>
