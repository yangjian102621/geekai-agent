<template>
  <div class="card">
    <div class="card-body">
      <div class="table-responsive" v-loading="appStore.loading">
        <div class="search-container d-flex flex-wrap align-items-end gap-2">
          <span class="search-item">
            <el-input
              type="text"
              v-model="appStore.query.name"
              placeholder="请输入应用名称"
              clearable
              @keyup.enter="appStore.fetchData(1)"
              @clear="appStore.fetchData(1)"
            />
          </span>
          <span class="search-item">
            <el-select
              v-model="appStore.query.cid"
              placeholder="请选择分类"
              style="width: 120px"
              clearable
              @change="appStore.fetchData(1)"
            >
              <el-option
                v-for="item in appStore.appCategories"
                :key="item.id"
                :label="item.name"
                :value="item.id"
              />
            </el-select>
          </span>
          <span class="search-item">
            <el-button type="primary" @click="appStore.fetchData(1)">
              <i class="iconfont icon-search"></i>
            </el-button>
          </span>
          <span class="search-item">
            <el-button type="primary" @click="appStore.add">
              <i class="iconfont icon-plus me-1"></i> 新增
            </el-button>
          </span>
          <span class="search-item">
            <el-button type="danger" @click="appStore.batchDelete">
              <i class="iconfont icon-remove me-1"></i> 删除
            </el-button>
          </span>
          <span class="search-item">
            <button
              class="btn btn-success btn-sm"
              @click="appStore.importCozeAgents"
            >
              <i class="iconfont icon-download"></i> 导入COZE智能体
            </button>
          </span>
        </div>

        <el-table
          :data="appStore.dataSets.items"
          border
          class="data-table"
          :row-key="(row) => row.id"
          @selection-change="appStore.handleSelectionChange"
          table-layout="auto"
        >
          <el-table-column type="selection" width="38"></el-table-column>
          <el-table-column label="名称">
            <template #default="scope">
              <div class="flex items-center justify-start gap-1">
                <el-image
                  :src="scope.row.icon"
                  class="rounded-circle max-w-10 max-h-10 min-w-10 min-h-10"
                >
                  <template #error>
                    <i class="iconfont icon-image text-4xl"></i>
                  </template>
                </el-image>
                <div>{{ scope.row.name }}</div>
              </div>
            </template>
          </el-table-column>
          <el-table-column label="分类" prop="cname" />
          <el-table-column label="类型" prop="type" />
          <el-table-column label="API URL">
            <template #default="scope">
              <span>{{ substr(scope.row.configs.api_url, 20) }}</span>
              <i
                class="iconfont icon-copy ms-1 mt-1"
                :data-clipboard-text="scope.row.configs.api_url"
              ></i>
            </template>
          </el-table-column>
          <el-table-column label="API Token">
            <template #default="scope">
              <span>{{ substr(scope.row.configs.token, 10) }}</span>
              <i
                class="iconfont icon-copy ms-1 mt-1"
                :data-clipboard-text="scope.row.configs.token"
              ></i>
            </template>
          </el-table-column>
          <el-table-column label="上架状态">
            <template #default="scope">
              <el-switch
                v-model="scope.row.enabled"
                @change="appStore.setValue(scope.row, 'enabled', $event)"
              />
            </template>
          </el-table-column>
          <el-table-column label="推荐状态">
            <template #default="scope">
              <el-switch
                v-model="scope.row.is_hot"
                @change="appStore.setValue(scope.row, 'is_hot', $event)"
              />
            </template>
          </el-table-column>
          <el-table-column label="每次对话积分" prop="score" />

          <el-table-column label="创建时间">
            <template #default="scope">
              <span>{{ dateFormat(scope.row['created_at']) }}</span>
            </template>
          </el-table-column>

          <el-table-column label="操作" fixed="right">
            <template #default="scope">
              <el-dropdown placement="bottom" trigger="hover">
                <button class="btn btn-primary btn-sm">
                  <i class="iconfont icon-more-vertical"></i>
                </button>
                <template #dropdown>
                  <el-dropdown-menu>
                    <el-dropdown-item @click="appStore.edit(scope.row)">
                      <i class="iconfont icon-edit"></i> 编辑
                    </el-dropdown-item>
                    <el-dropdown-item @click="appStore.copy(scope.row)">
                      <i class="iconfont icon-copy"></i> 复制
                    </el-dropdown-item>
                    <el-dropdown-item @click="appStore.remove(scope.row)">
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

      <div class="pagination p-3" v-if="appStore.dataSets.total > 0">
        <Pagination
          :total="appStore.dataSets.total"
          :pageSize="appStore.dataSets.pageSize"
          :currentPage="appStore.dataSets.page"
          :layout="['total', 'prev', 'pages', 'sizes', 'next']"
          @update:currentPage="appStore.fetchData"
          @update:pageSize="appStore.dataSets.pageSize = $event"
        />
      </div>
    </div>

    <model-dialog
      :modelValue="appStore.showDialog"
      :title="appStore.title"
      @cancel="appStore.showDialog = false"
      @confirm="appStore.handleSubmit"
      confirm-text="保存"
      cancel-text="关闭"
      :width="800"
    >
      <form>
        <div class="flex flex-row gap-2 mb-3">
          <div class="w-1/2">
            <label class="form-label"
              >应用类型 <i class="iconfont icon-xinghao text-danger"></i
            ></label>
            <el-select
              v-model="appStore.item.type"
              size="large"
              :class="{ 'is-invalid': appStore.errors.type }"
              placeholder="请选择应用类型"
              @change="appStore.changeAppType"
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
            <div class="invalid-feedback">{{ appStore.errors.type }}</div>
          </div>

          <div class="w-1/2">
            <label class="form-label"
              >分类 <i class="iconfont icon-xinghao text-danger"></i
            ></label>
            <el-select
              v-model="appStore.item.cid"
              size="large"
              :class="{ 'is-invalid': appStore.errors.cid }"
              placeholder="请选择分类"
            >
              <el-option
                v-for="category in appStore.appCategories"
                :key="category.id"
                :label="category.name"
                :value="category.id"
              />
            </el-select>
            <div class="invalid-feedback">{{ appStore.errors.category }}</div>
          </div>
        </div>

        <div class="flex flex-row gap-2 mb-3">
          <div class="w-1/2">
            <div class="mb-3">
              <label class="form-label"
                >应用名称 <i class="iconfont icon-xinghao text-danger"></i
              ></label>
              <input
                type="text"
                class="form-control"
                v-model="appStore.item.name"
                :class="{ 'is-invalid': appStore.errors.name }"
              />
              <div class="invalid-feedback">{{ appStore.errors.name }}</div>
            </div>

            <div class="flex flex-col">
              <label class="form-label">应用简介</label>
              <textarea
                class="form-control"
                v-model="appStore.item.summary"
                maxlength="255"
                rows="3"
              ></textarea>
            </div>
          </div>

          <div class="w-1/2">
            <div class="mb-3">
              <label class="form-label"
                >应用图标 <i class="iconfont icon-xinghao text-danger"></i
              ></label>
              <div class="flex flex-row gap-2 justify-center items-center">
                <el-image
                  :src="appStore.item.icon"
                  :preview-src-list="[appStore.item.icon]"
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
                  :http-request="appStore.uploadIcon"
                >
                  <button type="button" class="btn btn-primary">
                    <i class="iconfont icon-upload mr-2"></i>
                    <span class="text-base">上传图片</span>
                  </button>
                </el-upload>
              </div>
            </div>
          </div>
        </div>

        <div class="mb-3">
          <label class="form-label"
            >API Base URL
            <el-tooltip
              effect="dark"
              content="填写大模型API地址，如 https://api.deepseek.com <br /> 或者 Dify API服务地址，如：https://dify.geekai.me"
              placement="right"
              raw-content
            >
              <i class="iconfont icon-info ml-1 mr-1"></i>
            </el-tooltip>
            <i class="iconfont icon-xinghao text-danger"></i>
          </label>
          <input
            type="text"
            class="form-control"
            v-model="appStore.item.configs.api_url"
            :class="{ 'is-invalid': appStore.errors.api_url }"
          />
          <div class="invalid-feedback">{{ appStore.errors.api_url }}</div>
        </div>

        <div v-if="appStore.item.type === 'coze'">
          <div class="mb-3">
            <label class="form-label"
              >智能体 Bot ID
              <el-tooltip
                effect="dark"
                content="在Coze智能体搭建页面的地址栏获取"
                placement="right"
                raw-content
              >
                <i class="iconfont icon-info ml-1 mr-1"></i> </el-tooltip
              ><i class="iconfont icon-xinghao text-danger"></i
            ></label>
            <input
              type="text"
              class="form-control"
              v-model="appStore.item.configs.bot_id"
              placeholder="请输入智能体ID，如：7404176830909120538"
            />
          </div>
          <div class="mb-3">
            <label class="form-label"
              >授权应用 ID <i class="iconfont icon-xinghao text-danger"></i>

              <span class="text-xs ms-2 text-gray-400">
                获取方法参考：
                <a
                  href="https://docs.geekai.me/agent/config/coze.html#_3-%E5%88%9B%E5%BB%BA%E6%8E%88%E6%9D%83%E5%BA%94%E7%94%A8"
                  class="hover:underline"
                  target="_blank"
                  >创建Coze授权应用</a
                >
              </span>
            </label>
            <input
              type="text"
              class="form-control"
              v-model="appStore.item.configs.app_id"
            />
          </div>

          <div class="mb-3">
            <label class="form-label"
              >授权公钥 ID <i class="iconfont icon-xinghao text-danger"></i
            ></label>
            <input
              type="text"
              class="form-control"
              v-model="appStore.item.configs.public_key_id"
            />
          </div>
          <div class="mb-3">
            <label class="form-label"
              >授权私钥 <i class="iconfont icon-xinghao text-danger"></i
            ></label>
            <textarea
              class="form-control"
              v-model="appStore.item.configs.private_key"
              rows="3"
            ></textarea>
          </div>
        </div>

        <div v-else-if="appStore.item.type === 'bailian'">
          <div class="mb-3">
            <label class="form-label"
              >API Key <i class="iconfont icon-xinghao text-danger"></i
            ></label>
            <PasswordInput v-model="appStore.item.configs.bailian_api_key" />
          </div>
          <div class="mb-3">
            <label class="form-label"
              >应用 ID <i class="iconfont icon-xinghao text-danger"></i
            ></label>
            <input
              type="text"
              class="form-control"
              v-model="appStore.item.configs.bailian_app_id"
              placeholder="请输入百炼应用 ID"
            />
          </div>
        </div>

        <div v-else>
          <div class="mb-3">
            <label class="form-label"
              >API KEY <i class="iconfont icon-xinghao text-danger"></i
            ></label>
            <PasswordInput v-model="appStore.item.configs.token" />
          </div>

          <div v-if="appStore.item.type === 'openai'">
            <div class="mb-3">
              <label class="form-label">模型名称</label>
              <input
                type="text"
                class="form-control"
                placeholder="请输入模型名称，如: gpt-4o-mini"
                v-model="appStore.item.configs.model_name"
              />
            </div>

            <div class="mb-3">
              <label class="form-label">最大输出长度</label>
              <input
                type="number"
                class="form-control"
                placeholder="AI单次回复的最大长度"
                v-model="appStore.item.configs.max_length"
              />
            </div>

            <div class="mb-3">
              <label class="form-label">对话历史</label>
              <div>
                <el-switch v-model="appStore.item.configs.enable_context" />
              </div>
            </div>

            <div class="mb-3" v-if="appStore.item.configs.enable_context">
              <label class="form-label"
                >会话轮数
                <el-tooltip
                  content="轮数越大每次对话消耗的token越多"
                  placement="right"
                >
                  <i class="iconfont icon-info"></i></el-tooltip
              ></label>
              <input
                type="number"
                class="form-control"
                v-model="appStore.item.configs.history_deep"
                placeholder="历史记录包含的最大对话轮数"
              />
            </div>

            <div class="mb-3" v-if="appStore.item.configs.enable_context">
              <label class="form-label"
                >最大上下文长度
                <el-tooltip
                  effect="dark"
                  content="避免超出模型允许的最大上下文长度"
                  raw-content
                  placement="right"
                >
                  <i class="iconfont icon-info"></i>
                </el-tooltip>
              </label>
              <input
                type="number"
                class="form-control"
                v-model="appStore.item.configs.max_context_length"
                placeholder="如32K则填写 32768"
              />
            </div>
          </div>
        </div>

        <div class="mb-3">
          <label class="form-label"
            >每次对话消耗积分 <i class="iconfont icon-xinghao text-danger"></i
          ></label>
          <input
            type="number"
            class="form-control"
            v-model.number="appStore.item.score"
            :class="{ 'is-invalid': appStore.errors.score }"
          />
          <div class="invalid-feedback">{{ appStore.errors.score }}</div>
        </div>

        <div class="mb-3" v-if="appStore.item.score > 0">
          <label class="form-label">扣费触发方式</label>
          <div class="flex justify-start">
            <el-radio-group
              v-model="appStore.item.billing_mode"
              class="flex flex-column align-items-start gap-2"
            >
              <el-radio value="immediate"
                >立即扣费（AI输出内容就扣费）</el-radio
              >
              <el-radio value="file_suffix"
                >文件后缀触发（生成指定类型文件才扣费）</el-radio
              >
              <el-radio value="string_marker"
                >字符串标记触发（输出特定标记才扣费）</el-radio
              >
            </el-radio-group>
          </div>
        </div>

        <!-- 文件后缀配置 -->
        <div class="mb-3" v-if="appStore.item.billing_mode === 'file_suffix'">
          <label class="form-label">文件后缀列表</label>
          <el-select
            v-model="appStore.item.billing_config.suffixes"
            multiple
            filterable
            allow-create
            placeholder="选择或输入文件后缀，如 .pptx"
            clearable
          >
            <el-option
              v-for="suffix in fileSuffixes"
              :key="suffix"
              :label="suffix"
              :value="suffix"
            />
          </el-select>
          <div class="text-muted text-xs mt-1">
            AI返回的文本中包含这些文件后缀的URL时才扣费
          </div>
        </div>

        <!-- 字符串标记配置 -->
        <div class="mb-3" v-if="appStore.item.billing_mode === 'string_marker'">
          <label class="form-label">触发标记字符串</label>
          <input
            type="text"
            class="form-control"
            v-model="appStore.item.billing_config.marker"
            placeholder="如：####$$$$"
          />
          <div class="text-muted text-xs mt-1">AI输出包含此标记时才扣费</div>
        </div>

        <!-- <ParamConfig
          class="mb-3"
          v-model="appStore.item.params"
          title="应用表单配置"
        /> -->

        <div class="mb-3" v-if="appStore.item.type === 'openai'">
          <label class="form-label">系统预设提示词</label>
          <textarea
            class="form-control"
            v-model="appStore.item.configs.system_prompt"
            rows="3"
          ></textarea>
        </div>

        <div class="mb-3">
          <label class="form-label">是否上架</label>
          <div>
            <el-switch v-model="appStore.item.enabled" />
          </div>
        </div>

        <div class="mb-3">
          <label class="form-label">是否推荐到热门应用</label>
          <div>
            <el-switch v-model="appStore.item.is_hot" />
          </div>
        </div>
      </form>
    </model-dialog>

    <model-dialog
      :modelValue="appStore.showCozeAgentDialog"
      title="导入 COZE 智能体"
      @cancel="appStore.showCozeAgentDialog = false"
      @confirm="appStore.doImportCozeAgents"
      confirm-text="开始导入"
      :width="960"
    >
      <div class="table-responsive">
        <div
          class="d-flex justify-content-center align-items-center flex-wrap gap-2 mb-2 bg-purple-50 py-2 px-3 rounded"
        >
          <div class="text-muted">分类：</div>
          <div class="w-25">
            <el-select
              v-model="appStore.importCategoryId"
              placeholder="请选择分类"
            >
              <el-option
                v-for="category in appStore.appCategories"
                :key="category.id"
                :label="category.name"
                :value="category.id"
              />
            </el-select>
          </div>
        </div>
        <el-table
          :data="appStore.cozeAgentList"
          class="data-table"
          :row-key="(row) => row.bot_id"
          table-layout="auto"
          @selection-change="appStore.handleSelectionCoze"
        >
          <el-table-column type="selection" width="38"></el-table-column>
          <el-table-column label="智能体名称">
            <template #default="scope">
              <div class="flex items-center justify-start">
                <el-image
                  :src="scope.row.icon_url"
                  style="height: 32px; width: 32px"
                  class="rounded-circle"
                >
                  <template #error>
                    <i class="iconfont icon-image fs-2"></i>
                  </template>
                </el-image>
                <span class="ml-2">{{ scope.row.bot_name }}</span>
              </div>
            </template>
          </el-table-column>
          <el-table-column label="Bot ID" prop="bot_id" />

          <el-table-column label="发布时间">
            <template #default="scope">
              <span>{{ dateFormat(scope.row.publish_time) }}</span>
            </template>
          </el-table-column>
          <el-table-column label="简介">
            <template #default="scope">
              <el-popover
                placement="top-start"
                title="智能体介绍"
                :width="200"
                trigger="hover"
                :content="scope.row.description"
              >
                <template #reference>
                  <span class="cursor-pointer">{{
                    substr(scope.row.description, 20)
                  }}</span>
                </template>
              </el-popover>
            </template>
          </el-table-column>
          <template #empty>
            <el-empty description="暂无数据" />
          </template>
        </el-table>
      </div>
    </model-dialog>
  </div>
</template>

<script setup>
  import { onMounted, onUnmounted, ref } from 'vue'
  import ModelDialog from '@/components/ModelDialog.vue'
  import Pagination from '@/components/Pagination.vue'
  import ParamConfig from '@/components/admin/ParamConfig.vue'
  import PasswordInput from '@/components/PasswordInput.vue'
  import { useAppStore } from '@/js/store/admin/app'
  import { dateFormat, substr } from '@/js/utils/libs.js'

  const appStore = useAppStore()

  onMounted(() => {
    appStore.initialize()
  })

  onUnmounted(() => {
    appStore.releaseClipboard()
  })

  const appTypes = ref([
    { label: '大模型', value: 'openai', info: '支持国内外各通用大模型' },
    { label: 'Dify', value: 'dify', info: '支持Agent和ChatFlow两种类型' },
    {
      label: 'COZE',
      value: 'coze',
      info: '支持单Agent和多Agent对话智能体',
    },
    {
      label: '百炼',
      value: 'bailian',
      info: '支持阿里云百炼智能体和工作流',
    },
  ])

  const fileSuffixes = ref([
    '.pptx',
    '.pdf',
    '.doc',
    '.docx',
    '.xlsx',
    '.xls',
    '.png',
    '.jpg',
    '.jpeg',
  ])
</script>

<style scoped lang="scss">
  @use '@/assets/css/admin/admin.scss';
</style>
