<template>
  <div class="workflow-container">
    <!-- 左侧任务列表侧边栏 -->
    <div class="task-sidebar d-flex flex-column">
      <div class="sidebar-header p-3 border-bottom">
        <div class="d-flex justify-content-between align-items-center mb-2">
          <h5 class="m-0">我的任务</h5>
          <el-button
            link
            type="primary"
            @click="workflowStore.fetchTasks"
            :loading="workflowStore.loadingTasks"
          >
            <i class="iconfont icon-refresh"></i>
          </el-button>
        </div>
        <div class="d-flex gap-2">
          <el-select
            v-model="taskFilter"
            size="small"
            placeholder="状态筛选"
            clearable
            class="flex-grow-1"
          >
            <el-option label="执行中" value="running" />
            <el-option label="已完成" value="completed" />
            <el-option label="失败" value="failed" />
          </el-select>
        </div>
      </div>

      <div class="task-list flex-grow-1 overflow-auto custom-scrollbar">
        <div
          v-if="filteredTasks.length === 0"
          class="text-center p-4 text-muted"
        >
          <el-empty :image-size="80" description="暂无任务" />
        </div>
        <div
          v-for="task in filteredTasks"
          :key="task.task_id"
          class="task-item p-3 border-bottom cursor-pointer"
          :class="{ active: currentTaskId === task.task_id }"
          @click="showTaskDetail(task)"
        >
          <div class="d-flex justify-content-between align-items-start mb-1">
            <span
              class="fw-bold text-truncate"
              style="max-width: 120px"
              :title="task.workflow_name"
            >
              {{ task.workflow_name || '未知工作流' }}
            </span>
            <el-tag size="small" :type="statusInfo(task.status).type">
              {{ statusInfo(task.status).label }}
            </el-tag>
          </div>
          <div
            class="d-flex justify-content-between align-items-center text-secondary small mb-2"
          >
            <span>{{ formatTime(task.created_at) }}</span>
          </div>
          <el-progress
            v-if="['running', 'pending'].includes(task.status)"
            :percentage="task.progress"
            :status="task.status === 'failed' ? 'exception' : ''"
            :stroke-width="4"
            :show-text="false"
          />
        </div>
      </div>

      <div class="sidebar-footer p-3 border-top">
        <!-- 返回对话按钮 (已移除，通过 Home 布局导航) -->
      </div>
    </div>

    <!-- 右侧主内容区 -->
    <div class="main-content flex-grow-1 d-flex flex-column bg-light">
      <div
        class="content-header d-flex justify-content-between align-items-center"
      >
        <div>
          <h4>工作流广场</h4>
          <small>选择工作流模板，开始创建自动化任务</small>
        </div>
        <div class="user-info d-flex align-items-center gap-3">
          <div
            class="rounded-pill d-flex justify-content-center align-items-center gap-1 bg-purple-100 text-primary px-3 py-2 text-sm"
          >
            <i class="iconfont icon-doller text-warning text-sm"></i>
            <span>{{ userStore.userInfo.scores }} 积分</span>
          </div>
        </div>
      </div>

      <div class="content-body overflow-auto custom-scrollbar">
        <div class="workflow-grid">
          <div
            v-for="workflow in workflows"
            :key="workflow.workflow_id"
            class="workflow-card h-100 d-flex flex-column"
            @click="openCreateDialog(workflow)"
          >
            <!-- 大图区域 -->
            <div class="workflow-image-wrapper">
              <el-image :src="workflow.icon" fit="cover" class="workflow-image">
                <template #error>
                  <div class="workflow-image-placeholder">
                    <i class="iconfont icon-image"></i>
                    <span>{{ workflow.name?.charAt(0) }}</span>
                  </div>
                </template>
              </el-image>
              <!-- 积分标签覆盖在图片上 -->
              <div v-if="workflow.score > 0" class="workflow-score-badge">
                <i class="iconfont icon-doller"></i>
                <span>{{ workflow.score }} 积分</span>
              </div>
            </div>

            <!-- 内容区域 -->
            <div class="workflow-content d-flex flex-column flex-grow-1 p-3">
              <h5 class="workflow-title mb-2" :title="workflow.name">
                {{ workflow.name }}
              </h5>
              <p
                class="workflow-description line-clamp-2 mb-3 flex-grow-1"
                :title="workflow.summary"
              >
                {{ workflow.summary || '暂无描述' }}
              </p>
              <el-button
                type="primary"
                class="workflow-action-btn"
                @click.stop="openCreateDialog(workflow)"
              >
                <i class="iconfont icon-play me-1"></i>
                立即使用
              </el-button>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- 创建任务抽屉 -->
    <el-drawer
      v-model="showCreateDrawer"
      :title="selectedWorkflow?.name"
      size="400px"
      destroy-on-close
    >
      <div class="create-form-container" v-if="selectedWorkflow">
        <div
          class="alert alert-info mb-3 p-2 small"
          v-if="selectedWorkflow.summary"
        >
          {{ selectedWorkflow.summary }}
        </div>

        <el-form :model="formData" ref="formRef" label-position="top">
          <el-form-item
            v-for="param in selectedWorkflow.params"
            :key="param.name"
            :label="param.label || param.name"
            :prop="param.name"
            :required="param.required"
            :rules="[
              {
                required: param.required,
                message: '此项不能为空',
                trigger: 'blur',
              },
            ]"
          >
            <template #label>
              <span>{{ param.label || param.name }}</span>
              <el-tooltip
                v-if="param.description"
                :content="param.description"
                placement="top"
              >
                <i class="iconfont icon-info text-secondary ms-1"></i>
              </el-tooltip>
            </template>

            <!-- String / textarea -->
            <el-input
              v-if="
                param.type === 'String' ||
                !param.type ||
                param.type === 'textarea'
              "
              v-model="formData[param.name]"
              :type="param.type === 'textarea' ? 'textarea' : 'text'"
              :rows="param.type === 'textarea' ? 3 : 1"
              :placeholder="
                param.placeholder || '请输入 ' + (param.label || param.name)
              "
            />
            <!-- Number -->
            <el-input-number
              v-else-if="param.type === 'Number' || param.type === 'number'"
              v-model="formData[param.name]"
              class="w-100"
              :placeholder="param.placeholder"
            />
            <!-- Boolean -->
            <el-switch
              v-else-if="param.type === 'Boolean'"
              v-model="formData[param.name]"
            />
            <!-- Select -->
            <el-select
              v-else-if="param.type === 'Select' || param.type === 'select'"
              v-model="formData[param.name]"
              class="w-100"
              :placeholder="param.placeholder || '请选择'"
            >
              <el-option
                v-for="opt in getSelectOptions(param)"
                :key="opt"
                :label="opt"
                :value="opt"
              />
            </el-select>
            <!-- Radio -->
            <el-radio-group
              v-else-if="param.type === 'Radio'"
              v-model="formData[param.name]"
            >
              <el-radio v-for="opt in param.options" :key="opt" :value="opt">
                {{ opt }}
              </el-radio>
            </el-radio-group>
            <!-- CheckBox -->
            <el-checkbox-group
              v-else-if="param.type === 'CheckBox'"
              v-model="formData[param.name]"
            >
              <el-checkbox v-for="opt in param.options" :key="opt" :value="opt">
                {{ opt }}
              </el-checkbox>
            </el-checkbox-group>
            <!-- Image / Audio / Video / File 上传 -->
            <div v-else-if="isUploadType(param.type)" class="upload-field">
              <el-upload
                :auto-upload="true"
                :show-file-list="false"
                :http-request="(file) => handleUpload(file, param)"
                :accept="getAcceptByType(param.type)"
                :before-upload="(file) => beforeUpload(file, param)"
              >
                <el-button type="primary">
                  <i class="iconfont icon-upload me-1"></i> 上传文件
                </el-button>
              </el-upload>
              <div v-if="formData[param.name]" class="upload-preview mt-2">
                <el-image
                  v-if="param.type === 'Image'"
                  :src="formData[param.name]"
                  fit="cover"
                  class="preview-thumb"
                  :preview-src-list="[formData[param.name]]"
                />
                <audio
                  v-else-if="param.type === 'Audio'"
                  :src="formData[param.name]"
                  controls
                  class="preview-media"
                >
                  您的浏览器不支持音频播放
                </audio>
                <video
                  v-else-if="param.type === 'Video'"
                  :src="formData[param.name]"
                  controls
                  class="preview-media"
                >
                  您的浏览器不支持视频播放
                </video>
                <div v-else class="file-link">
                  <a
                    :href="formData[param.name]"
                    target="_blank"
                    class="text-primary"
                  >
                    <i class="iconfont icon-file me-1"></i>查看文件
                  </a>
                </div>
              </div>
            </div>
            <!-- DateTime -->
            <el-date-picker
              v-else-if="param.type === 'DateTime'"
              v-model="formData[param.name]"
              type="datetime"
              value-format="YYYY-MM-DD HH:mm:ss"
              format="YYYY-MM-DD HH:mm:ss"
              class="w-100"
              :placeholder="param.placeholder || '请选择日期时间'"
            />
            <!-- 默认文本输入 -->
            <el-input
              v-else
              v-model="formData[param.name]"
              :placeholder="param.placeholder"
            />
          </el-form-item>
        </el-form>
      </div>
      <template #footer>
        <div class="d-flex gap-2">
          <el-button @click="showCreateDrawer = false">取消</el-button>
          <el-button type="primary" @click="submitTask" :loading="submitting">
            提交任务
          </el-button>
        </div>
      </template>
    </el-drawer>

    <!-- 任务详情抽屉 -->
    <el-drawer
      v-model="showDetailDrawer"
      title="任务详情"
      size="500px"
      destroy-on-close
    >
      <div v-if="currentTask" class="task-detail-container">
        <el-descriptions :column="1" border>
          <el-descriptions-item label="任务ID">{{
            currentTask.task_id
          }}</el-descriptions-item>
          <el-descriptions-item label="工作流">{{
            currentTask.workflow_name
          }}</el-descriptions-item>
          <el-descriptions-item label="状态">
            <el-tag :type="statusInfo(currentTask.status).type">
              {{ statusInfo(currentTask.status).label }}
            </el-tag>
          </el-descriptions-item>
          <el-descriptions-item label="创建时间">{{
            formatTime(currentTask.created_at)
          }}</el-descriptions-item>
          <el-descriptions-item label="进度">
            <el-progress
              :percentage="currentTask.progress"
              :status="currentTask.status === 'failed' ? 'exception' : ''"
            />
          </el-descriptions-item>
        </el-descriptions>

        <div class="mt-4" v-if="currentTask.error">
          <h6 class="text-danger">错误信息</h6>
          <div class="bg-light p-3 rounded text-danger text-break">
            {{ currentTask.error }}
          </div>
        </div>

        <div
          class="mt-4"
          v-if="currentTask.params && Object.keys(currentTask.params).length"
        >
          <h6>输入参数</h6>
          <div class="bg-light p-3 rounded font-monospace small overflow-auto">
            <pre class="m-0">{{
              JSON.stringify(currentTask.params, null, 2)
            }}</pre>
          </div>
        </div>

        <div
          class="mt-4"
          v-if="currentTask.output && Object.keys(currentTask.output).length"
        >
          <h6>执行结果</h6>
          <!-- 原始数据 -->
          <div
            class="bg-light p-3 rounded font-monospace small overflow-auto mb-3"
          >
            <pre class="m-0">{{
              JSON.stringify(currentTask.output, null, 2)
            }}</pre>
          </div>

          <!-- 预览组件 -->
          <div class="mb-3">
            <WorkflowOutputPreview :output="currentTask.output" />
          </div>
        </div>
      </div>
      <template #footer>
        <div class="d-flex gap-2 justify-content-end" v-if="currentTask">
          <el-button
            v-if="['running', 'pending'].includes(currentTask.status)"
            type="danger"
            plain
            @click="handleCancelTask(currentTask)"
          >
            取消任务
          </el-button>
          <el-button
            v-if="['failed', 'canceled'].includes(currentTask.status)"
            type="primary"
            plain
            @click="handleRetryTask(currentTask)"
          >
            重试任务
          </el-button>
        </div>
      </template>
    </el-drawer>
  </div>
</template>

<script setup>
  import { ref, computed, onMounted, onUnmounted, reactive } from 'vue'
  import { useRouter } from 'vue-router'
  import { useWorkflowStore } from '@/js/store/workflow'
  import { useUserStore } from '@/js/store/user'
  import { ElMessage, ElMessageBox } from 'element-plus'
  import { frontUploadFile } from '@/js/store/common.js'
  import {
    showLoading,
    closeLoading,
    showMessageError,
  } from '@/js/utils/dialog.js'
  import WorkflowOutputPreview from '@/components/WorkflowOutputPreview.vue'
  import { replaceURL } from '@/js/utils/libs.js'

  const router = useRouter()
  const workflowStore = useWorkflowStore()
  const userStore = useUserStore()

  const taskFilter = ref('')
  const showCreateDrawer = ref(false)
  const showDetailDrawer = ref(false)
  const selectedWorkflow = ref(null)
  const currentTaskId = ref(null)
  const submitting = ref(false)
  const formData = reactive({})
  const formRef = ref(null)

  const workflows = computed(() => workflowStore.workflows)
  const tasks = computed(() => workflowStore.tasks)

  const filteredTasks = computed(() => {
    if (!taskFilter.value) return tasks.value
    return tasks.value.filter((t) => t.status === taskFilter.value)
  })

  const currentTask = computed(() => {
    return tasks.value.find((t) => t.task_id === currentTaskId.value)
  })

  // 状态映射
  const statusInfo = (status) => {
    const map = {
      pending: { label: '排队中', type: 'warning' },
      running: { label: '执行中', type: 'primary' },
      completed: { label: '已完成', type: 'success' },
      failed: { label: '失败', type: 'danger' },
      canceled: { label: '已取消', type: 'info' },
    }
    return map[status] || { label: status, type: 'info' }
  }

  const formatTime = (timestamp) => {
    if (!timestamp) return '-'
    const date = new Date(timestamp * 1000)
    return date.toLocaleString()
  }

  // 判断是否是上传类型
  const isUploadType = (type) => {
    return ['Image', 'Audio', 'Video', 'Doc', 'Zip', 'File'].includes(type)
  }

  // 获取文件类型限制
  const getAcceptByType = (type) => {
    if (type === 'Image') return '.png,.jpg,.jpeg,.bmp,.gif,.webp'
    if (type === 'Audio') return '.mp3,.wav,.ogg,.m4a'
    if (type === 'Video') return '.mp4,.avi,.mov,.wmv,.flv'
    if (type === 'Doc') return '.doc,.docx,.pdf,.txt,.xls,.xlsx,.csv,.ppt,.pptx'
    if (type === 'Zip') return '.zip,.rar,.7z,.tar,.gz'
    if (type === 'File')
      return '.doc,.docx,.pdf,.txt,.xls,.xlsx,.csv,.ppt,.pptx,.zip,.rar,.7z'
    return ''
  }

  // 上传前验证
  const beforeUpload = (file, param) => {
    const maxSize = (param.max_filesize || 5) * 1024 * 1024 // 默认5MB
    if (file.size > maxSize) {
      showMessageError(`文件大小不能超过 ${param.max_filesize || 5}MB`)
      return false
    }
    return true
  }

  // 处理文件上传
  const handleUpload = (file, param) => {
    showLoading('正在上传...')
    frontUploadFile(file, (data) => {
      formData[param.name] = replaceURL(data.url)
      closeLoading()
    })
  }

  // 获取选择框选项（兼容字符串数组和对象数组）
  const getSelectOptions = (param) => {
    if (!param.options) return []
    // 如果是字符串数组，直接返回
    if (Array.isArray(param.options) && param.options.length > 0) {
      if (typeof param.options[0] === 'string') {
        return param.options
      }
      // 如果是对象数组，提取 value
      return param.options.map((opt) => opt.value || opt)
    }
    return []
  }

  // 根据类型设置默认值
  const getDefaultValue = (param) => {
    if (param.default !== undefined && param.default !== null) {
      return param.default
    }
    // 根据类型设置默认值
    switch (param.type) {
      case 'Number':
      case 'number':
        return 0
      case 'Boolean':
        return false
      case 'CheckBox':
        return []
      case 'String':
      case 'Select':
      case 'Radio':
      case 'Image':
      case 'Audio':
      case 'Video':
      case 'File':
      case 'Doc':
      case 'Zip':
      default:
        return ''
    }
  }

  // 打开创建任务表单
  const openCreateDialog = (workflow) => {
    selectedWorkflow.value = workflow
    // 重置表单数据
    Object.keys(formData).forEach((key) => delete formData[key])
    if (workflow.params) {
      workflow.params.forEach((p) => {
        formData[p.name] = getDefaultValue(p)
      })
    }
    showCreateDrawer.value = true
  }

  // 提交任务
  const submitTask = async () => {
    if (!formRef.value) return

    await formRef.value.validate(async (valid) => {
      if (valid) {
        submitting.value = true
        try {
          const task = await workflowStore.createTask(
            selectedWorkflow.value.workflow_id,
            { ...formData }
          )
          showCreateDrawer.value = false
          currentTaskId.value = task.task_id
          showDetailDrawer.value = true // 创建后自动打开详情
        } catch (e) {
          const errorMessage = e.message || e.toString()
          showMessageError('创建任务失败：' + errorMessage)
        } finally {
          submitting.value = false
        }
      }
    })
  }

  // 显示任务详情
  const showTaskDetail = (task) => {
    currentTaskId.value = task.task_id
    showDetailDrawer.value = true
  }

  // 取消任务
  const handleCancelTask = (task) => {
    ElMessageBox.confirm('确定要取消该任务吗？', '提示', {
      type: 'warning',
    }).then(async () => {
      await workflowStore.cancelTask(task.task_id)
    })
  }

  // 重试任务
  const handleRetryTask = async (task) => {
    await workflowStore.retryTask(task.task_id)
  }

  // 生命周期
  let pollingTimer = null

  onMounted(async () => {
    await userStore.fetchUserInfo()
    await workflowStore.ensureInit()

    // 启动轮询
    pollingTimer = setInterval(() => {
      const hasRunning = tasks.value.some((t) =>
        ['pending', 'running'].includes(t.status)
      )
      if (hasRunning || showDetailDrawer.value) {
        workflowStore.fetchTasks(true) // 传入 true 表示静默刷新
      }
    }, 3000)
  })

  onUnmounted(() => {
    if (pollingTimer) clearInterval(pollingTimer)
  })
</script>

<style scoped lang="scss">
  @use '@/assets/css/workflow.scss';
</style>
