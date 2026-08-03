<template>
  <div class="flex flex-column justify-content-center align-items-center w-100">
    <el-form
      ref="formRef"
      :model="data"
      :rules="rules"
      label-width="120px"
      class="app-form w-100"
      label-position="top"
    >
      <el-form-item label="应用类型" prop="type">
        <el-select
          v-model="data.type"
          placeholder="请选择应用类型"
          style="width: 100%"
        >
          <el-option
            v-for="item in appTypes"
            :key="item.value"
            :value="item.value"
            :label="item.label"
          >
            <div class="d-flex justify-content-between align-items-center">
              <span style="float: left">{{ item.label }}</span>
              <span class="text-xs text-gray-400">
                {{ item.info }}
              </span>
            </div>
          </el-option>
        </el-select>
      </el-form-item>

      <el-form-item label="分类" prop="cid">
        <el-select
          v-model="data.cid"
          placeholder="请选择分类"
          style="width: 100%"
        >
          <el-option
            v-for="cat in categories"
            :key="cat.id"
            :label="cat.name"
            :value="cat.id"
          />
        </el-select>
      </el-form-item>

      <el-form-item label="应用名称" prop="name">
        <el-input
          v-model="data.name"
          placeholder="请输入应用名称"
          maxlength="50"
          show-word-limit
        />
      </el-form-item>

      <el-form-item label="API 地址" prop="configs.api_url">
        <el-input
          v-model="data.configs.api_url"
          placeholder="请输入API BASE URL，如：https://api.geekai.pro"
        />
      </el-form-item>

      <template v-if="data.type === 'coze'">
        <el-form-item label="Bot ID" prop="configs.bot_id">
          <el-input
            v-model="data.configs.bot_id"
            placeholder="请输入Coze智能体ID"
          />
        </el-form-item>
        <el-form-item label="授权应用ID" prop="configs.app_id">
          <el-input
            v-model="data.configs.app_id"
            placeholder="请输入授权应用ID"
          />
        </el-form-item>
        <el-form-item label="授权公钥ID" prop="configs.public_key_id">
          <el-input
            v-model="data.configs.public_key_id"
            placeholder="请输入授权公钥ID"
          />
        </el-form-item>
        <el-form-item label="授权私钥" prop="configs.private_key">
          <el-input
            v-model="data.configs.private_key"
            type="textarea"
            placeholder="请输入授权私钥"
            :rows="3"
          />
        </el-form-item>
      </template>

      <template v-else>
        <el-form-item label="API KEY" prop="configs.token">
          <el-input v-model="data.configs.token" placeholder="请输入API KEY" />
        </el-form-item>
        <template v-if="data.type === 'openai'">
          <el-form-item label="模型名称" prop="configs.model_name">
            <el-input
              v-model="data.configs.model_name"
              placeholder="请输入模型名称，如: gpt-4o-mini"
            />
          </el-form-item>
          <el-form-item label="最大输出长度" prop="configs.max_length">
            <el-input-number
              v-model="data.configs.max_length"
              :min="1"
              :max="32768"
              placeholder="AI单次回复的最大长度"
              style="width: 100%"
            />
          </el-form-item>
          <el-form-item label="对话历史" prop="configs.enable_context">
            <el-switch v-model="data.configs.enable_context" />
          </el-form-item>
          <el-form-item
            v-if="data.configs.enable_context"
            prop="configs.history_deep"
          >
            <template #label>
              <span>会话轮数</span>
              <el-tooltip content="历史记录包含的最大对话轮数" placement="top">
                <i class="iconfont icon-info ml-1"></i>
              </el-tooltip>
            </template>
            <el-input-number
              v-model="data.configs.history_deep"
              :min="1"
              :max="20"
              placeholder="历史记录包含的最大对话轮数"
              style="width: 100%"
            />
          </el-form-item>
          <el-form-item
            v-if="data.configs.enable_context"
            prop="configs.max_context_length"
          >
            <template #label>
              <span>上下文长度</span>
              <el-tooltip
                content="模型允许的最大上下文长度，<br>如32K则填写 32768，如16K则填写 16384"
                raw-content
                placement="top"
              >
                <i class="iconfont icon-info ml-1"></i>
              </el-tooltip>
            </template>
            <el-input-number
              v-model="data.configs.max_context_length"
              :min="1"
              :max="32768"
              placeholder="如32K则填写 32768"
              style="width: 100%"
            />
          </el-form-item>
        </template>
      </template>

      <el-form-item prop="score">
        <template #label>
          <span>对话积分</span>
          <el-tooltip content="每次对话消耗积分，如：10" placement="top">
            <i class="iconfont icon-info ml-1"></i>
          </el-tooltip>
        </template>
        <el-input-number
          v-model="data.score"
          :min="0"
          :max="1000"
          placeholder="每次对话消耗积分，如：10"
          style="width: 100%"
        />
      </el-form-item>

      <el-card class="mb-3">
        <template #header>
          <div class="d-flex justify-content-between">
            <span>应用参数配置</span>
            <el-tooltip content="添加参数" placement="top">
              <i
                class="iconfont icon-plus-circle cursor-pointer"
                @click="addParam"
              ></i>
            </el-tooltip>
          </div>
        </template>
        <el-table :data="data.params">
          <el-table-column label="参数说明" width="120">
            <template #default="scope">
              <el-input v-model="scope.row.label" />
            </template>
          </el-table-column>
          <el-table-column label="参数名称" width="120">
            <template #default="scope">
              <el-input v-model="scope.row.name" />
            </template>
          </el-table-column>
          <el-table-column label="参数类型" width="130">
            <template #default="scope">
              <el-select
                v-model="scope.row.type"
                placeholder="请选择参数类型"
                @change="changeParamType(scope.row)"
              >
                <el-option
                  v-for="type in paramTypes"
                  :key="type.value"
                  :label="type.label"
                  :value="type.value"
                />
              </el-select>
            </template>
          </el-table-column>
          <el-table-column label="默认值">
            <template #default="scope">
              <el-input
                v-model="scope.row.default"
                v-if="scope.row.type === 'String'"
              />
              <el-input-number
                v-else-if="scope.row.type === 'Number'"
                v-model="scope.row.default"
              />
              <el-switch
                v-else-if="scope.row.type === 'Boolean'"
                v-model="scope.row.default"
              />
              <el-date-picker
                v-else-if="scope.row.type === 'Date'"
                format="YYYY-MM-DD"
                v-model="scope.row.default"
              />
              <el-date-picker
                v-else-if="scope.row.type === 'DateTime'"
                format="YYYY-MM-DD HH:mm:ss"
                type="datetime"
                v-model="scope.row.default"
              />
              <div
                v-else-if="uploadTypes.includes(scope.row.type)"
                class="flex justify-content-start"
              >
                <el-upload
                  :auto-upload="true"
                  :show-file-list="false"
                  :http-request="
                    (file) => handleParamFileUpload(file, scope.row)
                  "
                  :accept="getAcceptByType(scope.row.type)"
                >
                  <el-button
                    type="primary"
                    class="flex flex-row justify-content-center align-items-center"
                    plain
                  >
                    <span class="mr-2"
                      ><i class="iconfont icon-upload"></i
                    ></span>
                    <span class="text-xs"> 文件大小不超过 5MB </span>
                  </el-button>
                </el-upload>
                <div
                  class="flex flex-row justify-content-center align-items-center"
                  v-if="scope.row.default"
                >
                  <el-image
                    v-if="scope.row.type === 'Image'"
                    :src="scope.row.default"
                    class="cursor-pointer ml-2 rounded-1"
                    fit="cover"
                    style="height: 30px"
                    @click="openPreview(scope.row.default)"
                  />
                  <a
                    v-else
                    :href="scope.row.default"
                    target="_blank"
                    class="cursor-pointer ml-2"
                  >
                    <el-tooltip content="点击下载文件" placement="top">
                      <el-image
                        :src="GetFileIcon(getFileExt(scope.row.default))"
                        fit="cover"
                        class="rounded-1"
                        style="height: 30px"
                      />
                    </el-tooltip>
                  </a>
                </div>
              </div>
            </template>
          </el-table-column>
          <el-table-column label="必填" width="60">
            <template #default="scope">
              <div class="d-flex justify-content-center">
                <el-checkbox v-model="scope.row.required" />
              </div>
            </template>
          </el-table-column>
          <el-table-column label="操作" width="60">
            <template #default="scope">
              <span
                class="bg-gray-200 cursor-pointer p-1 rounded-1 ml-2"
                @click="removeParam(scope.row)"
              >
                <i class="iconfont icon-remove"></i>
              </span>
            </template>
          </el-table-column>
        </el-table>
      </el-card>

      <el-form-item
        label="系统预设提示词"
        v-if="data.type === 'openai'"
        prop="configs.system_prompt"
      >
        <el-input
          v-model="data.configs.system_prompt"
          type="textarea"
          placeholder="请输入系统预设提示词"
          maxlength="1024"
          :rows="3"
          show-word-limit
        />
      </el-form-item>

      <el-form-item label="应用图标" prop="icon">
        <template #label>
          <span>应用图标</span>
          <el-tooltip content="建议尺寸：200x200px" placement="top">
            <i class="iconfont icon-info ml-1"></i>
          </el-tooltip>
        </template>
        <div style="display: flex; gap: 8px; align-items: center; width: 100%">
          <el-image :src="data.icon" class="h-[45px] w-[45px] rounded-3" />
          <el-upload
            :show-file-list="false"
            :http-request="uploadIcon"
            accept="image/*"
          >
            <el-button type="primary" size="large" circle>
              <i class="iconfont icon-upload"></i>
            </el-button>
          </el-upload>
        </div>
      </el-form-item>

      <el-form-item label="应用简介" prop="summary">
        <el-input
          v-model="data.summary"
          type="textarea"
          placeholder="请输入应用简介"
          maxlength="255"
          show-word-limit
          :rows="3"
        />
      </el-form-item>

      <el-form-item label="是否启用" prop="enabled">
        <el-switch v-model="data.enabled" />
      </el-form-item>
    </el-form>
    <div class="w-100 flex justify-content-end pr-2 pt-3 border-top">
      <button class="btn btn-secondary mr-2" @click="emits('cancel')">
        取消
      </button>
      <button class="btn btn-primary" @click="submit">提交</button>
    </div>
  </div>

  <el-image-viewer
    v-if="previewURL"
    :url-list="[previewURL]"
    @close="closePreview"
  />
</template>

<script setup>
  import { frontUploadFile } from '@/js/store/common.js'
  import { httpGet, httpPost } from '@/js/utils/http'
  import { getFileExt, GetFileIcon } from '@/js/utils/libs.js'
  import { ElMessage } from 'element-plus'
  import { onMounted, reactive, ref } from 'vue'

  const props = defineProps({
    form: {
      type: Object,
      default: () => ({}),
    },
    creatorId: {
      type: Number,
      default: 0,
    },
  })
  const emits = defineEmits(['cancel', 'success'])

  const formRef = ref(null)
  const appTypes = [
    { label: '大模型', value: 'openai', info: '支持国内外各通用大模型' },
    { label: 'Dify', value: 'dify', info: '支持Agent和ChatFlow两种类型' },
    { label: 'Coze', value: 'coze', info: '支持单Agent和多Agent对话智能体' },
  ]
  const categories = ref([])
  const paramTypes = [
    { label: 'String', value: 'String' },
    { label: 'Number', value: 'Number' },
    { label: 'Boolean', value: 'Boolean' },
    { label: 'Date', value: 'Date' },
    { label: 'DateTime', value: 'DateTime' },
    { label: 'File', value: 'File' },
    { label: 'Image', value: 'Image' },
    { label: 'Audio', value: 'Audio' },
    { label: 'Video', value: 'Video' },
    { label: 'Doc', value: 'Doc' },
    { label: 'Zip', value: 'Zip' },
  ]
  const uploadTypes = ['Image', 'Audio', 'Video', 'Doc', 'Zip', 'File']
  const previewURL = ref('')

  onMounted(async () => {
    try {
      const res = await httpGet('/api/creator/app-categories/list', {
        creator_id: props.creatorId,
        enabled: 1,
      })
      categories.value = res.data || []
    } catch (e) {
      ElMessage.error('获取分类失败：' + e.message)
    }
  })

  const data = reactive(props.form)
  if (!data.params) {
    data.params = []
  }
  const rules = reactive({
    type: [{ required: true, message: '请选择应用类型', trigger: 'change' }],
    cid: [{ required: true, message: '请选择分类', trigger: 'change' }],
    name: [
      { required: true, message: '请输入应用名称', trigger: 'blur' },
      { min: 2, max: 50, message: '应用名称长度在2-50个字符', trigger: 'blur' },
    ],
    'configs.api_url': [
      { required: true, message: '请输入API地址', trigger: 'blur' },
    ],
    'configs.token': [
      { required: true, message: '请输入API KEY', trigger: 'blur' },
    ],
    'configs.bot_id': [
      { required: true, message: '请输入Bot ID', trigger: 'blur' },
    ],
    'configs.app_id': [
      { required: true, message: '请输入授权应用ID', trigger: 'blur' },
    ],
    'configs.public_key_id': [
      { required: true, message: '请输入授权公钥ID', trigger: 'blur' },
    ],
    'configs.private_key': [
      { required: true, message: '请输入授权私钥', trigger: 'blur' },
    ],
    'configs.model_name': [
      { required: false, message: '请输入模型名称', trigger: 'blur' },
    ],
    'configs.max_length': [
      {
        required: false,
        type: 'number',
        min: 1,
        max: 32768,
        message: '最大输出长度范围1-32768',
        trigger: 'blur',
      },
    ],
    'configs.history_deep': [
      {
        required: false,
        type: 'number',
        min: 1,
        max: 20,
        message: '会话轮数范围1-20',
        trigger: 'blur',
      },
    ],
    'configs.max_context_length': [
      {
        required: false,
        type: 'number',
        min: 1,
        max: 32768,
        message: '最大上下文长度范围1-32768',
        trigger: 'blur',
      },
    ],
    score: [
      { required: true, message: '请输入积分消耗', trigger: 'blur' },
      {
        type: 'number',
        min: 0,
        max: 1000,
        message: '积分消耗范围0-1000',
        trigger: 'blur',
      },
    ],
    icon: [{ required: true, message: '请上传应用图标', trigger: 'blur' }],
    summary: [
      { required: true, message: '请输入应用简介', trigger: 'blur' },
      { min: 10, max: 255, message: '简介长度在10-255个字符', trigger: 'blur' },
    ],
  })

  onMounted(async () => {})
  const uploadIcon = (file) => {
    frontUploadFile(file, (res) => {
      data.icon = res.url
    })
  }

  const addParam = () => {
    data.params.push({
      label: '',
      name: '',
      type: 'String',
      default: '',
      required: true,
    })
  }

  const changeParamType = (param) => {
    if (param.type === 'Number') {
      param.default = 0
      return
    }
    if (param.type === 'Boolean') {
      param.default = false
      return
    }
    param.default = ''
  }

  const removeParam = (param) => {
    data.params = data.params.filter((p) => p !== param)
  }

  const getAcceptByType = (type) => {
    if (type === 'Image') return '.png,.jpg,.jpeg,.bmp,.gif'
    if (type === 'Audio') return '.mp3,.wav,.ogg'
    if (type === 'Video') return '.mp4,.avi,.mov'
    if (type === 'Doc') return '.doc,.docx,.pdf,.txt,.xls,.xlsx,.csv,.ppt,.pptx'
    if (type === 'Zip') return '.zip,.rar,.7z'
    return ''
  }

  const handleParamFileUpload = (file, param) => {
    frontUploadFile(file, (res) => {
      param.default = res.url
    })
  }

  const openPreview = (url) => {
    previewURL.value = url || ''
  }

  const closePreview = () => {
    previewURL.value = ''
  }

  // 提交表单
  const submit = () => {
    formRef.value?.validate((valid) => {
      if (valid) {
        httpPost(`/api/creator/apps/save`, data)
          .then(() => {
            ElMessage.success('保存成功')
            emits('success')
          })
          .catch((err) => {
            ElMessage.error(err.message)
          })
      }
    })
  }
</script>

<style scoped>
  .app-form {
    max-width: 600px;
  }
</style>
