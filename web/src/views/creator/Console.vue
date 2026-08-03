<template>
  <div class="creator-dashboard">
    <!-- 顶部统计卡片 -->
    <div class="stats-grid">
      <div class="stat-card">
        <div class="stat-icon bg-purple-500">
          <i class="iconfont icon-app text-3xl"></i>
        </div>
        <div class="stat-info">
          <h3>{{ dashStore.creator.app_count || 0 }}</h3>
          <p>我的应用</p>
        </div>
      </div>
      <div class="stat-card">
        <div class="stat-icon bg-blue-500">
          <i class="iconfont icon-reward text-3xl"></i>
        </div>
        <div class="stat-info">
          <h3>{{ dashStore.creator.total_earnings || 0 }}</h3>
          <p>总积分收益</p>
        </div>
      </div>
      <div class="stat-card">
        <div class="stat-icon bg-green-500">
          <i class="iconfont icon-doller text-3xl"></i>
        </div>
        <div class="stat-info">
          <h3>{{ dashStore.creator.today_earnings || 0 }}</h3>
          <p>今日积分收益</p>
        </div>
      </div>
      <div class="stat-card">
        <div class="stat-icon bg-red-500">
          <i class="iconfont icon-wallet text-3xl"></i>
        </div>
        <div class="stat-info">
          <h3>{{ dashStore.creator.scores || 0 }}</h3>
          <p>可提现积分</p>
        </div>
      </div>
    </div>

    <!-- 操作按钮 -->
    <div class="action-buttons">
      <el-button
        type="success"
        size="large"
        style="--el-button-size: 42px"
        @click="dashStore.showWithdrawDialog"
      >
        <i class="iconfont icon-withdraw mr-1 text-xl"></i>
        申请提现
      </el-button>
      <el-button
        type="primary"
        size="large"
        style="--el-button-size: 42px"
        @click="dashStore.showScoreLogsDialog"
      >
        <i class="iconfont icon-list mr-1 text-xl"></i>
        收益明细
      </el-button>
      <el-button
        type="info"
        size="large"
        style="--el-button-size: 42px"
        @click="dashStore.showWithdrawLogsDialog"
      >
        <i class="iconfont icon-withdraw-log mr-1 text-xl"></i>
        提现记录
      </el-button>
      <el-button
        size="large"
        style="--el-button-size: 42px"
        color="#5983f4"
        class="text-white"
        @click="dashStore.showProfileDialog"
      >
        <i class="iconfont icon-config mr-1 text-xl"></i>
        个人设置
      </el-button>

      <el-tooltip content="个人主页" placement="top">
        <a :href="`/creator/${dashStore.creator.username}`" target="_blank">
          <el-button
            size="large"
            style="--el-button-size: 42px"
            class="ml-2"
            circle
          >
            <i class="iconfont icon-home text-lg"></i> </el-button
        ></a>
      </el-tooltip>
      <el-tooltip content="分享个人主页" placement="top">
        <el-button
          size="large"
          style="--el-button-size: 42px"
          class="copy-share-url ml-2"
          circle
          :data-clipboard-text="shareUrl"
        >
          <i class="iconfont icon-share1 text-base"></i>
        </el-button>
      </el-tooltip>
    </div>

    <!-- 应用分类模块 -->
    <div class="apps-section" style="margin-bottom: 24px">
      <div class="section-header">
        <h3>应用分类</h3>
        <el-button type="primary" @click="dashStore.showCategoryDialog()">
          <i class="iconfont icon-plus mr-1"></i>
          添加应用分类
        </el-button>
      </div>
      <el-table
        :data="dashStore.categories"
        v-loading="dashStore.categoryLoading"
        border
      >
        <el-table-column prop="id" label="ID" width="80" />
        <el-table-column prop="name" label="分类名称" />
        <el-table-column prop="enabled" label="启用状态" width="100">
          <template #default="{ row }">
            <el-tag v-if="row.enabled" type="success">启用</el-tag>
            <el-tag v-else type="danger">禁用</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="created_at" label="创建时间" width="180">
          <template #default="{ row }">
            {{ dashStore.dateFormat(row.created_at) }}
          </template>
        </el-table-column>

        <el-table-column label="操作" width="180">
          <template #default="{ row }">
            <el-button size="small" @click="dashStore.showCategoryDialog(row)"
              >编辑</el-button
            >
            <el-popconfirm
              title="确认要删除该分类吗?"
              @confirm="dashStore.deleteCategory(row)"
              confirm-button-text="确定"
              cancel-button-text="取消"
              width="200px"
            >
              <template #reference>
                <el-button size="small" type="danger">删除</el-button>
              </template>
            </el-popconfirm>
          </template>
        </el-table-column>
      </el-table>
    </div>

    <!-- 应用列表 -->
    <div class="apps-section">
      <div class="section-header">
        <h3>我的应用</h3>
        <div class="flex">
          <el-select
            v-model="dashStore.appQuery.cid"
            placeholder="选择分类"
            class="mr-2"
            clearable
            @change="dashStore.fetchApps(1)"
          >
            <el-option
              v-for="category in dashStore.categories"
              :key="category.id"
              :label="category.name"
              :value="category.id"
            />
          </el-select>
          <el-input
            v-model="dashStore.appQuery.name"
            placeholder="搜索应用"
            class="mr-2"
            clearable
            @keyup.enter="dashStore.fetchApps(1)"
            @clear="dashStore.fetchApps(1)"
          />
          <el-button type="primary" @click="dashStore.fetchApps(1)">
            <i class="iconfont icon-search"></i>
          </el-button>
          <el-button type="primary" @click="dashStore.showAppDialog">
            <i class="iconfont icon-plus mr-1"></i>
            创建应用
          </el-button>
        </div>
      </div>

      <el-table
        :data="dashStore.appsData.items"
        v-loading="dashStore.appsLoading"
        border
        class="apps-table"
      >
        <el-table-column prop="cname" label="应用分类" width="100" />
        <el-table-column label="应用信息">
          <template #default="{ row }">
            <div class="flex items-center gap-2">
              <img
                v-if="row.icon"
                :src="row.icon"
                class="h-[45px] w-[45px] rounded-md"
                fit="cover"
              />
              <div>
                <div class="text-base font-bold p-1">{{ row.name }}</div>
                <div class="text-xs text-gray-500 p-1">{{ row.summary }}</div>
              </div>
            </div>
          </template>
        </el-table-column>
        <el-table-column prop="type" label="应用类型" width="110">
          <template #default="{ row }">
            <el-tag type="primary">{{ dashStore.getAppType(row.type) }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="score" label="积分消耗" width="100" />
        <el-table-column label="状态" width="120">
          <template #default="{ row }">
            <el-tag v-if="row.check === 1" type="success">审核通过</el-tag>
            <el-tag v-else-if="row.check === 2" type="danger"
              >审核不通过</el-tag
            >
            <el-tag v-else type="info">待审核</el-tag>
          </template>
        </el-table-column>
        <el-table-column label="启用状态" width="100">
          <template #default="{ row }">
            <el-switch
              v-model="row.enabled"
              :disabled="row.check !== 1"
              @change="dashStore.enableApp(row)"
            />
          </template>
        </el-table-column>
        <el-table-column label="操作" width="60" fixed="right">
          <template #default="{ row }">
            <el-dropdown placement="bottom" trigger="hover">
              <button
                class="btn btn-primary btn-sm"
                style="--bs-btn-padding-y: 0.1rem; --bs-btn-padding-x: 0.3rem"
              >
                <i class="iconfont icon-more-vertical"></i>
              </button>
              <template #dropdown>
                <el-dropdown-menu>
                  <el-dropdown-item @click="dashStore.editApp(row)">
                    <i class="iconfont icon-edit"></i> 编辑
                  </el-dropdown-item>
                  <el-dropdown-item @click="dashStore.deleteApp(row)">
                    <span class="text-danger">
                      <i class="iconfont icon-remove"></i> 删除
                    </span>
                  </el-dropdown-item>
                </el-dropdown-menu>
              </template>
            </el-dropdown>
          </template>
        </el-table-column>
      </el-table>

      <!-- 分页 -->
      <div class="pagination" v-if="dashStore.appsData.total > 0">
        <Pagination
          :total="dashStore.appsData.total"
          :pageSize="dashStore.appsData.page_size"
          :currentPage="dashStore.appsData.page"
          @update:currentPage="dashStore.fetchApps"
          @update:pageSize="dashStore.appsData.page_size = $event"
        />
      </div>
    </div>

    <div class="flex justify-content-center mt-4 text-gray-400">
      <Footer />
    </div>

    <!-- 应用创建/编辑对话框 -->
    <ModelDialog
      :modelValue="dashStore.appDialog.visible"
      :title="dashStore.appDialog.form.id ? '编辑应用' : '创建应用'"
      @cancel="dashStore.appDialog.visible = false"
      :hide-footer="true"
    >
      <AppForm
        ref="appFormRef"
        :form="dashStore.appDialog.form"
        :creatorId="dashStore.creator.id"
        @cancel="dashStore.appDialog.visible = false"
        @success="dashStore.createAppSuccess"
      />
    </ModelDialog>

    <!-- 提现申请对话框 -->
    <ModelDialog
      :modelValue="dashStore.withdrawDialog.visible"
      title="申请提现"
      @cancel="dashStore.withdrawDialog.visible = false"
      :hide-footer="true"
    >
      <WithdrawForm
        ref="withdrawFormRef"
        @success="dashStore.withdrawSuccess"
        @cancel="dashStore.withdrawDialog.visible = false"
      />
    </ModelDialog>

    <!-- 收益明细对话框 -->
    <ModelDialog
      :modelValue="dashStore.scoreLogsDialog.visible"
      title="积分收益明细"
      @cancel="dashStore.scoreLogsDialog.visible = false"
      width="1200px"
      :hide-footer="true"
    >
      <ScoreLogs ref="scoreLogsRef" />
    </ModelDialog>

    <!-- 提现记录对话框 -->
    <ModelDialog
      :modelValue="dashStore.withdrawLogsDialog.visible"
      title="提现记录"
      @cancel="dashStore.withdrawLogsDialog.visible = false"
      width="1500px"
      :hide-footer="true"
    >
      <WithdrawsList ref="withdrawsListRef" />
    </ModelDialog>

    <!-- 个人设置对话框 -->
    <ModelDialog
      :modelValue="dashStore.profileDialog.visible"
      title="个人设置"
      @cancel="dashStore.profileDialog.visible = false"
      :hide-footer="true"
    >
      <ProfileForm
        ref="profileFormRef"
        :data="dashStore.profileDialog.form"
        @cancel="dashStore.profileDialog.visible = false"
        @success="dashStore.submitProfile"
      />
    </ModelDialog>

    <!-- 分类对话框 -->
    <ModelDialog
      :modelValue="dashStore.categoryDialog.visible"
      :title="
        dashStore.categoryDialog.form.id ? '编辑应用分类' : '添加应用分类'
      "
      @cancel="dashStore.categoryDialog.visible = false"
      @confirm="dashStore.submitCategory"
    >
      <el-form :model="dashStore.categoryDialog.form" label-position="top">
        <el-form-item label="分类名称">
          <el-input v-model="dashStore.categoryDialog.form.name" />
        </el-form-item>
        <el-form-item label="是否启用">
          <el-switch v-model="dashStore.categoryDialog.form.enabled" />
        </el-form-item>
      </el-form>
    </ModelDialog>
  </div>
</template>

<script setup>
  import Footer from '@/components/Footer.vue'
  import ModelDialog from '@/components/ModelDialog.vue'
  import Pagination from '@/components/Pagination.vue'
  import AppForm from '@/components/creator/AppForm.vue'
  import ProfileForm from '@/components/creator/ProfileForm.vue'
  import ScoreLogs from '@/components/creator/ScoreLogs.vue'
  import WithdrawForm from '@/components/creator/WithdrawForm.vue'
  import WithdrawsList from '@/components/creator/WithdrawsList.vue'
  import { useConsoleStore } from '@/js/store/creator/console'
  import ClipboardJS from 'clipboard'
  import { ElMessage } from 'element-plus'
  import { onMounted, onUnmounted, ref } from 'vue'

  const dashStore = useConsoleStore()
  const clipboard = ref(null)

  onMounted(async () => {
    try {
      await dashStore.fetchDashboard()
      await dashStore.fetchApps()
      await dashStore.fetchCategories()
      clipboard.value = new ClipboardJS('.copy-share-url')
      clipboard.value.on('success', (e) => {
        ElMessage.success('复制链接成功，快分享给好友吧！')
        e.clearSelection()
      })
      clipboard.value.on('error', (e) => {
        ElMessage.error('复制分享链接失败!')
        e.clearSelection()
      })
    } catch (error) {
      ElMessage.error(error.message)
    }
  })

  onUnmounted(() => {
    clipboard.value.destroy()
  })

  const shareUrl = computed(() => {
    return `${window.location.protocol}//${window.location.host}/creator/${dashStore.creator.username}`
  })
</script>

<style scoped>
  .creator-dashboard {
    padding: 20px;
  }

  .stats-grid {
    display: grid;
    grid-template-columns: repeat(auto-fit, minmax(250px, 1fr));
    gap: 20px;
    margin-bottom: 30px;
  }

  .stat-card {
    background: white;
    border-radius: 8px;
    padding: 24px;
    display: flex;
    align-items: center;
    box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
  }

  .stat-icon {
    width: 60px;
    height: 60px;
    border-radius: 12px;
    display: flex;
    align-items: center;
    justify-content: center;
    margin-right: 16px;
    font-size: 24px;
    color: white;
  }

  .stat-info h3 {
    font-size: 28px;
    font-weight: bold;
    margin: 0 0 8px 0;
    color: #333;
  }

  .stat-info p {
    margin: 0;
    color: #666;
    font-size: 14px;
  }

  .action-buttons {
    margin-bottom: 30px;
    display: flex;
    gap: 0px;
    flex-wrap: wrap;
  }

  .apps-section {
    background: white;
    border-radius: 8px;
    padding: 24px;
    box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
  }

  .section-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 20px;
  }

  .section-header h3 {
    margin: 0;
    font-size: 18px;
    font-weight: 600;
  }

  .pagination {
    margin-top: 20px;
    display: flex;
    justify-content: center;
  }
</style>
