<template>
  <div class="card">
    <div class="card-body">
      <div class="table-responsive">
        <div class="header">
          <el-button type="primary" @click="handleAdd">新增产品</el-button>
        </div>

        <el-table :data="products" border style="width: 100%; margin-top: 20px">
          <el-table-column prop="name" label="产品名称">
            <template #default="scope">
              <span class="sort" :data-id="scope.row.id">
                <i class="iconfont icon-drag cursor-move"></i>
                {{ scope.row.name }}
              </span>
            </template>
          </el-table-column>
          <el-table-column prop="price" label="价格" />
          <el-table-column prop="credit" label="额度" />
          <el-table-column prop="sales" label="销量" />
          <el-table-column label="启用状态">
            <template #default="scope">
              <el-switch
                v-model="scope.row.enabled"
                @change="handleStatusChange(scope.row)"
              />
            </template>
          </el-table-column>
          <el-table-column label="更新时间">
            <template #default="scope">
              <span>{{ dateFormat(scope.row.updated_at) }}</span>
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
                    <el-dropdown-item @click="handleEdit(scope.row)">
                      <i class="iconfont icon-edit"> 编辑</i>
                    </el-dropdown-item>
                    <el-dropdown-item @click="handleDelete(scope.row)">
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
      :modelValue="dialogVisible"
      :title="dialogTitle"
      @cancel="dialogVisible = false"
      @confirm="handleSubmit"
      confirm-text="保存"
      cancel-text="关闭"
      :width="500"
    >
      <form>
        <div class="mb-3">
          <label class="form-label"
            >产品名称 <i class="iconfont icon-xinghao text-danger"></i
          ></label>
          <input
            type="text"
            class="form-control"
            v-model="form.name"
            :class="{ 'is-invalid': errors.name }"
          />
          <div class="invalid-feedback">{{ errors.name }}</div>
        </div>

        <div class="mb-3">
          <label class="form-label"
            >原价 <i class="iconfont icon-xinghao text-danger"></i
          ></label>
          <input
            type="number"
            class="form-control"
            v-model.number="form.price"
            :class="{ 'is-invalid': errors.price }"
            step="0.01"
            min="0"
          />
          <div class="invalid-feedback">{{ errors.price }}</div>
        </div>

        <div class="mb-3">
          <label class="form-label"
            >额度 <i class="iconfont icon-xinghao text-danger"></i
          ></label>
          <input
            type="number"
            class="form-control"
            v-model.number="form.credit"
            :class="{ 'is-invalid': errors.credit }"
            min="0"
          />
          <div class="invalid-feedback">{{ errors.credit }}</div>
        </div>

        <div class="mb-3">
          <label class="form-label">启用状态</label>
          <div>
            <el-switch v-model="form.enabled" />
          </div>
        </div>
      </form>
    </model-dialog>
  </div>
</template>

<script setup>
  import { onMounted } from 'vue'
  import { storeToRefs } from 'pinia'
  import ModelDialog from '@/components/ModelDialog.vue'
  import { useAdminProductStore } from '@/js/store/admin/product'
  import { dateFormat } from '@/js/utils/libs'
  import Sortable from 'sortablejs'

  const productStore = useAdminProductStore()
  const { products, dialogVisible, dialogTitle, form, errors } =
    storeToRefs(productStore)
  const {
    handleAdd,
    handleEdit,
    handleDelete,
    handleStatusChange,
    handleSubmit,
    handleSort,
    initialize,
  } = productStore

  const setupSortable = () => {
    const tableBody = document.querySelector('.el-table__body tbody')
    if (!tableBody) return

    Sortable.create(tableBody, {
      sort: true,
      animation: 500,
      onEnd({ newIndex, oldIndex, from }) {
        if (oldIndex === newIndex) return

        const sortedData = Array.from(from.children).map((row) =>
          row.querySelector('.sort').getAttribute('data-id')
        )
        const ids = []
        const sorts = []
        sortedData.forEach((id, index) => {
          ids.push(parseInt(id))
          sorts.push(index + 1)
          products.value[index].sort_num = index + 1
        })

        handleSort(ids, sorts)
      },
    })
  }

  onMounted(async () => {
    await initialize()
    setupSortable()
  })
</script>
