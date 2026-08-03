<template>
  <div class="card">
    <div class="card-body">
      <div class="table-responsive" v-loading="loading">
        <div class="flex flex-wrap gap-2 mb-3">
          <div class="search-item">
            <el-button type="primary" @click="add">
              <i class="iconfont icon-plus mr-1"></i> 新增
            </el-button>
          </div>
          <div class="search-item">
            <el-button type="success" @click="openImportDialog">
              <i class="iconfont icon-download mr-1"></i> 批量导入Coze工作流
            </el-button>
          </div>
          <div class="search-item">
            <el-button type="danger" @click="batchDelete">
              <i class="iconfont icon-remove mr-1"></i> 删除
            </el-button>
          </div>
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
          <el-table-column label="名称" width="200">
            <template #default="{ row }">
              <div class="flex items-center justify-start gap-1">
                <el-image
                  :src="row.icon"
                  class="w-10 h-10 min-w-10 min-h-10 rounded-circle"
                  fit="cover"
                >
                  <template #error>
                    <i class="iconfont icon-image text-4xl"></i>
                  </template>
                </el-image>
                <div>{{ row.name }}</div>
              </div>
            </template>
          </el-table-column>
          <el-table-column label="工作流ID" prop="workflow_id" />
          <el-table-column label="消耗积分" prop="score" />
          <el-table-column label="启用状态">
            <template #default="{ row }">
              <el-switch v-model="row.enabled" @change="enable(row)" />
            </template>
          </el-table-column>
          <el-table-column label="创建时间" width="200">
            <template #default="{ row }">
              <span>{{ dateFormat(row['created_at']) }}</span>
            </template>
          </el-table-column>

          <el-table-column label="最后运行时间" width="200">
            <template #default="{ row }">
              <span>{{ dateFormat(row['last_run_at']) }}</span>
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

      <div class="pagination p-3" v-if="dataSets.total > 0">
        <Pagination
          :total="dataSets.total"
          :pageSize="dataSets.pageSize"
          :currentPage="dataSets.page"
          :layout="['total', 'prev', 'pages', 'sizes', 'next']"
          @update:currentPage="fetchData"
          @update:pageSize="dataSets.pageSize = $event"
        />
      </div>
    </div>

    <model-dialog
      :modelValue="showDialog"
      :title="title"
      @cancel="showDialog = false"
      @confirm="handleSubmit"
      :width="800"
      confirm-text="保存"
      cancel-text="关闭"
    >
      <form>
        <div class="flex flex-row gap-2 mb-3">
          <div class="w-1/2">
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

          <div class="w-1/2">
            <label class="form-label"
              >工作流 ID <i class="iconfont icon-xinghao text-danger"></i
            ></label>
            <input
              type="text"
              class="form-control"
              v-model="item.workflow_id"
              :class="{ 'is-invalid': errors.workflow_id }"
            />
            <div class="invalid-feedback">{{ errors.workflow_id }}</div>
          </div>
        </div>

        <div class="flex flex-row gap-2 mb-3">
          <div class="w-1/2">
            <div class="mb-3">
              <label class="form-label"
                >消耗积分 <i class="iconfont icon-xinghao text-danger"></i
              ></label>
              <input
                type="number"
                class="form-control"
                v-model.number="item.score"
                :class="{ 'is-invalid': errors.score }"
              />
              <div class="invalid-feedback">{{ errors.score }}</div>
            </div>

            <div class="mb-0">
              <label class="form-label">简介</label>
              <textarea
                class="form-control"
                v-model="item.summary"
                maxlength="255"
                rows="3"
              ></textarea>
            </div>
          </div>

          <div class="mb-3">
            <label class="form-label"
              >工作流图标 <i class="iconfont icon-xinghao text-danger"></i
            ></label>
            <div class="flex flex-row gap-2 justify-center items-center">
              <el-image
                :src="item.icon"
                :preview-src-list="[item.icon]"
                fit="cover"
                class="rounded-lg h-[170px] w-[170px] border border-gray-100"
              >
                <template #error>
                  <div class="w-full h-full flex justify-center items-center">
                    <i class="iconfont icon-image text-4xl"></i>
                  </div>
                </template>
              </el-image>
              <el-upload
                :auto-upload="true"
                :show-file-list="false"
                :http-request="uploadIcon"
              >
                <button type="button" class="btn btn-primary">
                  <i class="iconfont icon-upload mr-2"></i>
                  <span class="text-base">上传图片</span>
                </button>
              </el-upload>
            </div>
            <div class="invalid-feedback">{{ errors.icon }}</div>
          </div>
        </div>

        <ParamConfig
          class="mb-3"
          v-model="item.params"
          title="工作流参数配置"
        />

        <div class="mb-3">
          <label class="form-label"
            >工作流类型 <i class="iconfont icon-xinghao text-danger"></i
          ></label>
          <el-select
            v-model="item.type"
            placeholder="请选择工作流类型"
            @change="changeWorkflowType"
          >
            <el-option label="Coze" value="coze" />
            <el-option label="百炼" value="bailian" />
          </el-select>
        </div>

        <div v-if="item.type === 'coze'" class="mb-3">
          <label class="form-label">授权配置</label>
          <div class="alert alert-info mb-2 p-2 small">
            配置 Coze API 授权信息，用于调用工作流
          </div>
          <div class="mb-2">
            <label class="form-label small">API URL</label>
            <input
              type="text"
              class="form-control form-control-sm"
              v-model="item.auth_config.api_url"
              placeholder="https://api.coze.cn"
            />
          </div>
          <div class="mb-2">
            <label class="form-label small">授权应用 ID</label>
            <input
              type="text"
              class="form-control form-control-sm"
              v-model="item.auth_config.app_id"
            />
          </div>
          <div class="mb-2">
            <label class="form-label small">授权公钥 ID</label>
            <input
              type="text"
              class="form-control form-control-sm"
              v-model="item.auth_config.public_key_id"
            />
          </div>
          <div class="mb-2">
            <label class="form-label small">授权私钥</label>
            <textarea
              class="form-control form-control-sm"
              v-model="item.auth_config.private_key"
              rows="3"
            />
          </div>
        </div>

        <div v-if="item.type === 'bailian'" class="mb-3">
          <label class="form-label">授权配置</label>
          <div class="alert alert-info mb-2 p-2 small">
            配置阿里云百炼 API 授权信息，用于调用工作流
          </div>
          <div class="mb-2">
            <label class="form-label small">API Key</label>
            <input
              type="text"
              class="form-control form-control-sm"
              v-model="item.bailian_auth_config.api_key"
              placeholder="请输入百炼 API Key"
            />
          </div>
          <div class="mb-2">
            <label class="form-label small">应用 ID</label>
            <input
              type="text"
              class="form-control form-control-sm"
              v-model="item.bailian_auth_config.app_id"
              placeholder="请输入百炼应用 ID"
            />
          </div>
        </div>

        <div class="mb-3">
          <label class="form-label">是否上架</label>
          <div>
            <el-switch v-model="item.enabled" />
          </div>
        </div>
      </form>
    </model-dialog>

    <!-- 导入工作流对话框 -->
    <model-dialog
      :modelValue="showImportDialog"
      title="从 Coze 导入工作流"
      :width="800"
      @cancel="showImportDialog = false"
      @confirm="importSelectedWorkflows"
      confirm-text="批量导入选中"
      cancel-text="取消"
      :hideConfirm="selectedImportWorkflows.length === 0 || importing"
    >
      <div class="mb-3">
        <label class="form-label">授权配置</label>
        <div class="mb-2">
          <label class="form-label small">API URL</label>
          <input
            type="text"
            class="form-control form-control-sm"
            v-model="importAuthConfig.api_url"
            placeholder="https://api.coze.cn"
          />
        </div>
        <div class="mb-2">
          <label class="form-label small">Coze 空间ID</label>
          <input
            type="text"
            class="form-control form-control-sm"
            v-model="importAuthConfig.space_id"
          />
        </div>
        <div class="mb-2">
          <label class="form-label small">授权应用 ID</label>
          <input
            type="text"
            class="form-control form-control-sm"
            v-model="importAuthConfig.app_id"
          />
        </div>
        <div class="mb-2">
          <label class="form-label small">授权公钥 ID</label>
          <input
            type="text"
            class="form-control form-control-sm"
            v-model="importAuthConfig.public_key_id"
          />
        </div>
        <div class="mb-2">
          <label class="form-label small">授权私钥</label>
          <textarea
            class="form-control form-control-sm"
            v-model="importAuthConfig.private_key"
            rows="3"
          />
        </div>
        <div class="pt-2">
          <el-button
            type="primary"
            @click="fetchWorkflows"
            :loading="importing"
          >
            <i class="iconfont icon-search mr-1"></i> 获取工作流列表
          </el-button>
        </div>
      </div>

      <div v-if="importedWorkflows.length > 0" class="mt-3">
        <label class="form-label">选择要导入的工作流</label>
        <el-table
          :data="importedWorkflows"
          border
          max-height="400"
          @selection-change="
            (selection) => (selectedImportWorkflows = selection)
          "
        >
          <el-table-column type="selection" width="55" />
          <el-table-column label="工作流ID" prop="workflow_id" />
          <el-table-column label="名称" prop="workflow_name" />
          <el-table-column label="描述" prop="description" />
        </el-table>
      </div>
    </model-dialog>

    <!-- 导入结果对话框 -->
    <model-dialog
      :modelValue="showImportResultDialog"
      title="导入结果"
      :width="800"
      @cancel="showImportResultDialog = false"
      :hideFooter="true"
    >
      <div class="p-3">
        <!-- 汇总信息 -->
        <div
          v-if="importResultSummary"
          class="mb-4 p-3 rounded"
          style="background-color: #f5f7fa"
        >
          <div class="flex items-center gap-4 text-sm">
            <span>
              共 <strong>{{ importResultSummary.total }}</strong> 个工作流
            </span>
            <span class="text-success">
              导入成功: <strong>{{ importResultSummary.imported }}</strong>
            </span>
            <span class="text-primary">
              更新成功: <strong>{{ importResultSummary.updated }}</strong>
            </span>
            <span v-if="importResultSummary.failed > 0" class="text-danger">
              失败: <strong>{{ importResultSummary.failed }}</strong>
            </span>
          </div>
        </div>

        <!-- 详细结果列表 -->
        <div class="space-y-2" style="max-height: 500px; overflow-y: auto">
          <div
            v-for="(result, index) in importResults"
            :key="index"
            class="flex items-start gap-3 p-3 rounded border"
            :style="{
              backgroundColor:
                result.status === 'imported'
                  ? 'var(--bs-success-bg-subtle)'
                  : result.status === 'updated'
                  ? 'var(--bs-primary-bg-subtle)'
                  : 'var(--bs-danger-bg-subtle)',
              borderColor:
                result.status === 'imported'
                  ? 'var(--bs-success-border-subtle)'
                  : result.status === 'updated'
                  ? 'var(--bs-primary-border-subtle)'
                  : 'var(--bs-danger-border-subtle)',
            }"
          >
            <!-- 状态图标 -->
            <div class="flex-shrink-0 mt-0.5">
              <i
                v-if="result.status === 'imported'"
                class="iconfont icon-success-line text-success text-xl"
              ></i>
              <i
                v-else-if="result.status === 'updated'"
                class="iconfont icon-refresh text-primary text-xl"
              ></i>
              <i
                v-else
                class="iconfont icon-error-line text-danger text-xl"
              ></i>
            </div>

            <!-- 结果内容 -->
            <div class="flex-1">
              <div
                class="font-medium mb-1"
                :class="{
                  'text-success': result.status === 'imported',
                  'text-primary': result.status === 'updated',
                  'text-danger': result.status === 'failed',
                }"
              >
                <span v-if="result.status === 'imported'">
                  {{ result.workflow_name }} 工作流导入成功
                </span>
                <span v-else-if="result.status === 'updated'">
                  {{ result.workflow_name }} 工作流更新成功
                </span>
                <span v-else> {{ result.workflow_name }} 导入失败 </span>
              </div>
              <div
                v-if="result.status === 'failed' && result.error"
                class="text-sm text-danger mt-1"
              >
                原因：{{ result.error }}
              </div>
            </div>
          </div>
        </div>
      </div>
    </model-dialog>
  </div>
</template>

<script setup>
  import { onMounted, onUnmounted, ref } from 'vue'
  import { storeToRefs } from 'pinia'
  import ModelDialog from '@/components/ModelDialog.vue'
  import Pagination from '@/components/Pagination.vue'
  import ParamConfig from '@/components/admin/ParamConfig.vue'
  import { useWorkflowStore } from '@/js/store/admin/workflow'
  import { dateFormat } from '@/js/utils/libs.js'
  import {
    showLoading,
    closeLoading,
    showMessageError,
    showMessageOK,
  } from '@/js/utils/dialog.js'
  import { httpGet, httpPost } from '@/js/utils/http.js'

  const workflowStore = useWorkflowStore()
  const { loading, item, errors, showDialog, title, dataSets } =
    storeToRefs(workflowStore)
  const {
    handleSubmit,
    remove,
    enable,
    add,
    edit,
    handleSelectionChange,
    batchDelete,
    fetchData,
    initialize,
    releaseClipboard,
    uploadIcon,
    changeWorkflowType,
  } = workflowStore

  const showImportDialog = ref(false)
  const importing = ref(false)
  const importedWorkflows = ref([])
  const selectedImportWorkflows = ref([])
  const importAuthConfig = ref({
    api_url: '',
    space_id: '',
    app_id: '',
    public_key_id: '',
    private_key: '',
  })
  const showImportResultDialog = ref(false)
  const importResults = ref([])
  const importResultSummary = ref(null)

  // 打开导入对话框时，自动加载 coze 配置
  const openImportDialog = async () => {
    try {
      const res = await httpGet('/api/admin/config/get?name=coze')
      if (res.data) {
        // 将配置填充到导入授权配置中
        importAuthConfig.value = {
          api_url: res.data.api_url || '',
          space_id: res.data.space_id || '',
          app_id: res.data.app_id || '',
          public_key_id: res.data.public_key_id || '',
          private_key: res.data.private_key || '',
        }
      }
    } catch (e) {
      // 如果加载配置失败，不影响打开对话框，只是不填充默认值
      console.warn('加载 Coze 配置失败:', e.message)
    }
    showImportDialog.value = true
  }

  const fetchWorkflows = async () => {
    importing.value = true
    showLoading('正在获取工作流列表...')
    try {
      const res = await httpPost('/api/admin/workflow/import', {
        auth_config: importAuthConfig.value,
      })
      importedWorkflows.value = res.data || []
      selectedImportWorkflows.value = []
      if (importedWorkflows.value.length === 0) {
        showMessageError('未找到工作流')
      } else {
        showMessageOK(`找到 ${importedWorkflows.value.length} 个工作流`)
      }
    } catch (e) {
      showMessageError('获取工作流列表失败：' + e.message)
    } finally {
      importing.value = false
      closeLoading()
    }
  }

  const importSelectedWorkflows = async () => {
    if (selectedImportWorkflows.value.length === 0) {
      showMessageError('请选择要导入的工作流')
      return
    }

    importing.value = true
    showLoading('正在批量导入工作流...')
    try {
      const res = await httpPost('/api/admin/workflow/batch-import', {
        workflows: selectedImportWorkflows.value,
        auth_config: importAuthConfig.value,
      })

      // 保存导入结果
      importResults.value = res.data.results || []
      importResultSummary.value = res.data.summary || null

      // 关闭导入对话框，显示结果对话框
      showImportDialog.value = false
      showImportResultDialog.value = true

      // 刷新列表
      fetchData(dataSets.value.page)
    } catch (e) {
      showMessageError('批量导入失败：' + e.message)
    } finally {
      importing.value = false
      closeLoading()
    }
  }

  const createDefaultItem = () => ({
    type: 'coze',
    params: [],
    options: {},
    enabled: true,
    icon: '/images/app-placeholder.png',
    summary: '',
    auth_config: {
      api_url: '',
      app_id: '',
      public_key_id: '',
      private_key: '',
    },
    bailian_auth_config: {
      api_key: '',
      app_id: '',
    },
  })

  onMounted(() => {
    initialize()
  })

  onUnmounted(() => {
    releaseClipboard()
  })
</script>

<style scoped lang="scss"></style>
