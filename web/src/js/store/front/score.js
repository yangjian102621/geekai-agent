import { defineStore } from 'pinia'
import { ref } from 'vue'
import { httpPost } from '../../utils/http'
import { showMessageError } from '../../utils/dialog'

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

export const useFrontScoreStore = defineStore('front-score', () => {
  const dataSets = ref({ ...DEFAULT_DATASETS })
  const loading = ref(false)
  const query = ref({
    type: 0,
    start_time: '',
    end_time: '',
    page: DEFAULT_DATASETS.page,
    page_size: DEFAULT_DATASETS.page_size,
  })

  const fetchData = async (page = 1) => {
    loading.value = true
    query.value.page = page
    query.value.page_size =
      dataSets.value.page_size || DEFAULT_DATASETS.page_size
    try {
      const res = await httpPost('/api/score/list', query.value)
      if (res?.data) {
        dataSets.value = {
          ...DEFAULT_DATASETS,
          ...res.data,
        }
      }
    } catch (e) {
      showMessageError('获取数据失败：' + e.message)
    } finally {
      loading.value = false
    }
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
    query,
    scoreTypes: SCORE_TYPES,
    fetchData,
    getScoreTypeLabel,
    getScoreTypeClass,
  }
})
