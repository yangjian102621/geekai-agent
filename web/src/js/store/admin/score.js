import { defineStore } from 'pinia'
import { ref } from 'vue'
import { httpPost } from '../../utils/http'
import {
  showConfirm,
  showMessageError,
  showMessageOK,
} from '../../utils/dialog'

const DEFAULT_DATASETS = {
  page: 1,
  page_size: 20,
  total: 0,
  items: [],
}

const SCORE_TYPES = [
  { value: 0, label: '全部' },
  { value: 1, label: '充值' },
  { value: 2, label: '消费' },
  { value: 3, label: '退款' },
  { value: 4, label: '奖励' },
  { value: 5, label: '兑换' },
  { value: 6, label: '系统' },
]

export const useAdminScoreStore = defineStore('admin-score', () => {
  const dataSets = ref({ ...DEFAULT_DATASETS })
  const loading = ref(false)
  const query = ref({
    username: '',
    type: 0,
    start_time: '',
    end_time: '',
    page: DEFAULT_DATASETS.page,
    page_size: DEFAULT_DATASETS.page_size,
  })
  const totalScores = ref(0)
  const selectedIds = ref([])

  const fetchData = async (page = 1) => {
    loading.value = true
    query.value.page = page
    query.value.page_size =
      dataSets.value.page_size || DEFAULT_DATASETS.page_size
    try {
      const res = await httpPost('/api/admin/score/list', query.value)
      if (res?.data) {
        dataSets.value = {
          ...DEFAULT_DATASETS,
          ...(res.data.data || {}),
        }
        totalScores.value = res.data.stat || 0
      }
    } catch (e) {
      showMessageError('获取数据失败：' + e.message)
    } finally {
      loading.value = false
    }
  }

  const batchRemove = async () => {
    if (selectedIds.value.length === 0) {
      showMessageError('请先选择要删除的记录')
      return
    }

    showConfirm('批量删除', '确定要删除所有选中的记录吗？', async () => {
      try {
        await httpPost('/api/admin/score/batchRemove', {
          ids: selectedIds.value,
        })
        showMessageOK('删除成功')
        selectedIds.value = []
        await fetchData(1)
      } catch (e) {
        showMessageError('删除失败：' + e.message)
      }
    })
  }

  const handleSelectionChange = (selection) => {
    selectedIds.value = selection.map((item) => item.id)
  }

  const getScoreTypeLabel = (type) => {
    const found = SCORE_TYPES.find((item) => item.value === type)
    return found ? found.label : '未知'
  }

  const getScoreTypeClass = (type) => {
    switch (type) {
      case 1:
        return 'text-bg-primary'
      case 2:
        return 'text-bg-danger'
      case 3:
        return 'text-bg-warning'
      case 4:
        return 'text-bg-success'
      case 5:
        return 'text-bg-info'
      case 6:
        return 'text-bg-secondary'
      default:
        return 'text-bg-secondary'
    }
  }

  return {
    dataSets,
    loading,
    totalScores,
    query,
    scoreTypes: SCORE_TYPES,
    fetchData,
    batchRemove,
    handleSelectionChange,
    getScoreTypeLabel,
    getScoreTypeClass,
  }
})
