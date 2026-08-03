<template>
  <div class="min-h-screen bg-gray-50 pb-12 md:pb-16">
    <!-- 页面加载状态 -->
    <div
      v-if="pageLoading"
      class="max-w-6xl mx-auto px-4 sm:px-6 lg:px-8 py-10"
    >
      <el-skeleton :rows="8" animated />
    </div>

    <div v-else>
      <!-- 创作者信息头部 -->
      <div
        class="bg-gradient-to-b from-violet-50 to-gray-50"
        v-if="creatorInfo"
      >
        <div class="max-w-6xl mx-auto px-4 sm:px-6 lg:px-8 py-5">
          <div class="bg-white rounded-xl p-6 sm:p-8 shadow-sm">
            <div
              class="flex flex-col md:flex-row md:items-start items-center gap-6 md:gap-8 text-center md:text-left"
            >
              <!-- Avatar -->
              <div
                class="w-full md:w-auto flex justify-center md:justify-start"
              >
                <el-image
                  :src="creatorInfo.logo || '/images/avatar/user.png'"
                  class="w-24 h-24 sm:w-28 sm:h-28 md:w-32 md:h-32 rounded-full border-4 border-white shadow-md"
                  fit="cover"
                >
                  <template #error>
                    <div
                      class="w-full h-full rounded-full bg-gray-100 flex items-center justify-center text-5xl text-gray-400"
                    >
                      <i class="iconfont icon-user"></i>
                    </div>
                  </template>
                </el-image>
              </div>

              <!-- Details -->
              <div class="flex-grow w-full">
                <div
                  class="flex flex-col md:flex-row md:items-baseline gap-2.5 mb-3"
                >
                  <h1
                    class="text-2xl sm:text-3xl font-semibold text-gray-800 m-0"
                  >
                    {{ creatorInfo.name }}
                  </h1>
                  <span class="text-sm sm:text-base text-gray-500"
                    >@{{ route.params.username }}</span
                  >
                </div>
                <p class="text-sm text-gray-600 mb-5 leading-relaxed">
                  {{
                    creatorInfo.description || '这个用户很懒，什么都没有留下'
                  }}
                </p>
                <div
                  class="flex flex-wrap justify-center md:justify-start gap-6"
                >
                  <div class="flex items-center gap-1.5 text-gray-600">
                    <span
                      class="text-lg font-semibold text-gray-800 leading-none"
                      >0</span
                    >
                    <span class="text-sm text-gray-500">关注</span>
                  </div>
                  <div class="flex items-center gap-1.5 text-gray-600">
                    <span
                      class="text-lg font-semibold text-gray-800 leading-none"
                      >0</span
                    >
                    <span class="text-sm text-gray-500">粉丝</span>
                  </div>
                  <div class="flex items-center gap-1.5 text-gray-600">
                    <span
                      class="text-lg font-semibold text-gray-800 leading-none"
                      >0</span
                    >
                    <span class="text-sm text-gray-500">获赞</span>
                  </div>
                </div>
              </div>

              <!-- Actions -->
              <div
                class="flex gap-3 md:gap-4 w-full md:w-auto justify-center md:justify-end mt-4 md:mt-0"
              >
                <el-button
                  circle
                  class="copy-share-url"
                  :data-clipboard-text="shareUrl"
                  ><i class="iconfont icon-share1"></i
                ></el-button>
              </div>
            </div>
          </div>
        </div>
      </div>

      <!-- 应用展示区域 -->
      <div class="max-w-6xl mx-auto px-4 sm:px-6 lg:px-8 pb-0">
        <div
          class="flex flex-col md:flex-row md:justify-between md:items-end items-center mb-8 gap-4"
        >
          <div class="flex-1 w-full md:min-w-[200px]">
            <h2
              class="text-2xl font-bold text-gray-800 m-0 flex items-center gap-2.5 justify-center md:justify-start"
            >
              <i class="iconfont icon-apps"></i>
              智能体 ({{ appsData.items?.length || 0 }})
            </h2>
          </div>
          <div class="flex items-center gap-4 flex-wrap w-full md:w-auto">
            <div class="filter-tabs w-full md:w-auto">
              <el-radio-group
                v-model="categoryId"
                @change="handleCategoryChange"
                class="filter-tabs__scroller flex gap-2 md:gap-3 w-full md:w-auto overflow-x-auto md:overflow-visible whitespace-nowrap md:whitespace-normal justify-start"
              >
                <el-radio-button value="0">全部</el-radio-button>
                <el-radio-button
                  v-for="category in appCategories"
                  :key="category.id"
                  :value="category.id"
                >
                  {{ category.name }}
                </el-radio-button>
              </el-radio-group>
            </div>
          </div>
        </div>

        <!-- 应用加载状态 -->
        <div v-if="appsLoading" class="mb-10">
          <div
            class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-4 sm:gap-6"
          >
            <el-skeleton
              v-for="i in 6"
              :key="i"
              animated
              class="bg-white rounded-2xl overflow-hidden shadow-sm"
            >
              <template #template>
                <el-skeleton-item variant="image" class="w-full h-[200px]" />
                <div class="p-4">
                  <el-skeleton-item variant="h3" class="w-1/2" />
                  <el-skeleton-item variant="text" class="mt-2.5" />
                  <el-skeleton-item variant="text" class="mt-2.5 w-4/5" />
                </div>
              </template>
            </el-skeleton>
          </div>
        </div>

        <!-- 应用列表 -->
        <div
          v-else-if="appsData.items && appsData.items.length > 0"
          class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-4 sm:gap-6 mb-10"
        >
          <div
            v-for="app in appsData.items"
            :key="app.id"
            class="bg-white rounded-2xl overflow-hidden shadow-sm transition-all duration-300 ease-in-out cursor-pointer border-2 border-transparent flex flex-col h-full md:hover:-translate-y-1 md:hover:shadow-lg md:hover:border-indigo-500"
          >
            <div class="relative pt-5 sm:pt-6 px-5 sm:px-6 pb-4">
              <el-image
                :src="app.icon || '/images/app-placeholder.png'"
                class="w-16 h-16 sm:w-20 sm:h-20 rounded-2xl shadow-md"
                fit="cover"
              >
                <template #error>
                  <div
                    class="w-16 h-16 sm:w-20 sm:h-20 rounded-2xl bg-gray-100 flex items-center justify-center text-4xl text-gray-300"
                  >
                    <i class="iconfont icon-app"></i>
                  </div>
                </template>
              </el-image>
              <div
                class="absolute top-4 right-4 sm:top-6 sm:right-6 bg-indigo-100 text-indigo-600 px-3 py-1 rounded-full text-xs font-medium"
              >
                {{ getAppTypeLabel(app.type) }}
              </div>
            </div>

            <div class="px-5 sm:px-6 pb-5 sm:pb-6 flex flex-col flex-grow">
              <h3
                class="text-lg sm:text-xl font-semibold mb-2 text-gray-800 truncate"
              >
                {{ app.name }}
              </h3>
              <p
                class="text-sm text-gray-500 leading-snug mb-4 h-[40px] overflow-hidden line-clamp-2"
              >
                {{ app.summary }}
              </p>

              <div class="flex gap-4 mb-3 flex-wrap">
                <div class="flex items-center gap-1 text-sm text-gray-500">
                  <span>累计使用次数：</span>
                  <span>{{ formatNumber(app.usage_count || 0) }}</span>
                </div>
              </div>

              <div
                class="flex flex-col sm:flex-row sm:justify-between sm:items-center gap-3 mt-auto"
              >
                <div
                  class="flex items-baseline gap-1 justify-center sm:justify-start"
                >
                  <span
                    class="text-lg font-semibold text-red-500 leading-none"
                    >{{ app.score }}</span
                  >
                  <span class="text-xs text-gray-400">积分/次</span>
                </div>
                <a
                  :href="`/chat/${UUID()}?app_id=${app.id}`"
                  class="w-full sm:w-auto"
                  target="_blank"
                >
                  <el-button
                    type="primary"
                    class="!rounded-full w-full sm:w-auto justify-center"
                  >
                    <i class="iconfont icon-play"></i>
                    <span class="ml-1">立即使用</span>
                  </el-button>
                </a>
              </div>
            </div>
          </div>
        </div>

        <!-- 空状态 -->
        <div v-else class="text-center py-10 px-5 bg-white rounded-2xl mb-10">
          <div class="text-gray-300 mb-4">
            <i class="iconfont icon-empty-box text-6xl"></i>
          </div>
          <h3 class="text-xl text-gray-800 mb-2">暂无应用</h3>
          <p class="text-sm text-gray-600 m-0">该创作者还没有发布任何应用</p>
        </div>

        <!-- 分页 -->
        <div class="flex justify-center py-3">
          <Pagination
            :total="appsData.total"
            :pageSize="appsData.page_size"
            :currentPage="appsData.page"
            @update:currentPage="fetchApps"
          />
        </div>

        <div class="flex justify-content-center py-3 text-gray-400">
          <Footer />
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
  import Footer from '@/components/Footer.vue'
  import Pagination from '@/components/Pagination.vue'
  import { useSharedStore } from '@/js/cache/sharedata'
  import { httpGet } from '@/js/utils/http'
  import { UUID } from '@/js/utils/libs'
  import ClipboardJS from 'clipboard'
  import { ElMessage } from 'element-plus'
  import { computed, onMounted, onUnmounted, ref } from 'vue'
  import { useRoute, useRouter } from 'vue-router'

  const route = useRoute()
  const router = useRouter()
  const store = useSharedStore()
  // 数据状态
  const pageLoading = ref(true)
  const appsLoading = ref(false)

  const creatorInfo = ref(null)
  const appsData = ref({
    items: [],
    total: 0,
    page: 1,
    page_size: 12,
  })

  const categoryId = ref(0)
  const clipboard = ref(null)
  const appCategories = ref([])

  const shareUrl = computed(() => {
    if (!creatorInfo.value) return ''
    return `${window.location.protocol}//${window.location.host}/creator/${creatorInfo.value.username}`
  })

  // 生命周期
  onMounted(async () => {
    await fetchCreatorInfo()
    await fetchApps()
    await fetchAppCategories()
    pageLoading.value = false
    clipboard.value = new ClipboardJS('.copy-share-url')
    clipboard.value.on('success', (e) => {
      ElMessage.success('复制链接成功，快分享给好友吧！')
      e.clearSelection()
    })
    clipboard.value.on('error', (e) => {
      ElMessage.error('复制分享链接失败!')
      e.clearSelection()
    })
  })

  onUnmounted(() => {
    clipboard.value.destroy()
  })

  // 获取创作者信息
  const fetchCreatorInfo = async () => {
    try {
      const res = await httpGet(`/api/creator/${route.params.username}`)
      creatorInfo.value = res.data
    } catch (error) {
      ElMessage.error('获取创作者信息失败：' + error.message)
    }
  }

  // 获取应用列表
  const fetchApps = async (page = 1) => {
    appsLoading.value = true
    try {
      const params = {
        page,
        page_size: appsData.value.page_size,
        category_id: categoryId.value,
      }
      const res = await httpGet(
        `/api/creator/${route.params.username}/apps`,
        params
      )
      appsData.value = res.data
    } catch (error) {
      ElMessage.error('获取应用列表失败：' + error.message)
    } finally {
      appsLoading.value = false
    }
  }

  // 获取应用分类
  const fetchAppCategories = async () => {
    try {
      const res = await httpGet(
        `/api/creator/app-categories/list?creator_id=${creatorInfo.value.id}`
      )
      appCategories.value = res.data
    } catch (error) {
      ElMessage.error('获取应用分类失败：' + error.message)
    }
  }

  const handleCategoryChange = () => {
    appsData.value.page = 1
    fetchApps(1)
  }

  // 工具函数
  const formatNumber = (num) => {
    if (num >= 10000) {
      return (num / 10000).toFixed(1) + 'w'
    } else if (num >= 1000) {
      return (num / 1000).toFixed(1) + 'k'
    }
    return num.toString()
  }

  const getAppTypeLabel = (type) => {
    const typeMap = {
      openai: '大模型',
      coze: 'Coze',
      dify: 'Dify',
    }
    return typeMap[type] || type
  }
</script>

<style scoped>
  .filter-tabs__scroller {
    padding: 0.25rem 0;
    -webkit-overflow-scrolling: touch;
    scroll-snap-type: x proximity;
  }

  .filter-tabs__scroller::-webkit-scrollbar {
    display: none;
  }

  .filter-tabs__scroller :deep(.el-radio-button__inner) {
    border-radius: 9999px;
    padding: 0.5rem 1.25rem;
    min-width: 96px;
    justify-content: center;
    scroll-snap-align: center;
  }

  .filter-tabs__scroller :deep(.is-active .el-radio-button__inner) {
    box-shadow: 0 4px 12px rgba(99, 102, 241, 0.25);
  }

  @media (min-width: 768px) {
    .filter-tabs__scroller {
      overflow: visible;
      padding: 0;
    }
  }
</style>
