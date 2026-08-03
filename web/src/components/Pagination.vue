<template>
  <nav class="bs-pagination d-flex" v-if="show">
    <div class="page-item" v-if="arrayContains(layout, 'total')">
      <button type="button" class="btn btn-light" disabled>
        Total: {{ total }}
      </button>
    </div>

    <!-- PageSize选择 -->
    <div class="me-3 ms-2" v-if="arrayContains(layout, 'sizes')">
      <select
        class="form-select"
        :class="{ 'form-select-sm': small }"
        v-model="pageSizeInternal"
        @change="changePageSize"
        style="width: 80px"
      >
        <option
          v-for="size in pageSizeOptionsInternal"
          :key="size"
          :value="size"
        >
          {{ size }}
        </option>
      </select>
    </div>

    <ul class="pagination" :class="{ 'pagination-sm': small }">
      <!-- 上一页 -->
      <li
        class="page-item"
        :class="{ disabled: currentPage === 1 }"
        v-if="arrayContains(layout, 'prev')"
      >
        <a
          class="page-link"
          href="#"
          @click.prevent="goToPage(currentPage - 1)"
        >
          <i class="iconfont icon-prev-page"></i>
        </a>
      </li>
      <!-- 页码 -->
      <template v-if="arrayContains(layout, 'pages')">
        <li
          v-for="page in pages"
          :key="page"
          class="page-item"
          :class="{ active: currentPage === page }"
        >
          <a
            v-if="page !== '...'"
            class="page-link"
            href="#"
            @click.prevent="goToPage(page)"
            >{{ page }}</a
          >
          <a
            v-else
            class="page-link page-more"
            href="#"
            @click.prevent="() => {}"
          >
            <i class="iconfont icon-more-horizontal"></i>
          </a>
        </li>
      </template>
      <!-- 下一页 -->
      <li
        class="page-item"
        :class="{ disabled: currentPage === totalPage }"
        v-if="arrayContains(layout, 'next')"
      >
        <a
          class="page-link"
          href="#"
          @click.prevent="goToPage(currentPage + 1)"
        >
          <i class="iconfont icon-next-page"></i>
        </a>
      </li>
      <!-- 跳转输入 -->
      <li
        class="page-item page-item-goto"
        v-if="arrayContains(layout, 'jumper')"
      >
        <div class="input-group" :class="{ 'input-group-sm': small }">
          <input type="text" class="form-control" v-model="inputPage" />
          <button class="btn btn-primary" @click="jumpToPage">GO</button>
        </div>
      </li>
    </ul>
  </nav>
</template>

<script setup>
  import { arrayContains } from '@/js/utils/libs.js'

  const props = defineProps({
    small: {
      type: Boolean,
      default: false,
    },
    hideOnSinglePage: {
      type: Boolean,
      default: true,
    },
    total: {
      type: Number,
      required: true,
    },
    pageSize: {
      type: Number,
      default: 10,
    },
    currentPage: {
      type: Number,
      default: 1,
    },
    pageSizeOptions: {
      type: Array,
      default: () => [10, 20, 30, 50],
    },
    layout: {
      type: Array,
      default: () => ['total', 'prev', 'pages', 'next'], // full values: ["total", "admin", "jumper", "sizes","prev","next"]
    },
  })

  const emits = defineEmits(['update:currentPage', 'update:pageSize'])

  const currentPageInternal = ref(props.currentPage)
  const pageSizeInternal = ref(props.pageSize)
  const totalInternal = ref(props.total)
  const pageSizeOptionsInternal = ref(props.pageSizeOptions)

  const totalPage = computed(() => {
    return Math.ceil(totalInternal.value / pageSizeInternal.value)
  })

  const show = computed(() => {
    return totalPage.value > 1 || !props.hideOnSinglePage
  })

  const pages = computed(() => {
    const pagesArray = []
    // 始终显示第一页
    pagesArray.push(1)
    if (totalPage.value <= 1) return pagesArray

    // 如果第一页和当前页之间有间隔，添加省略号
    if (currentPageInternal.value - 3 > 1) {
      pagesArray.push('...')
    }

    // 显示当前页前后各两个页码
    for (
      let i = Math.max(currentPageInternal.value - 2, 1);
      i <= Math.min(currentPageInternal.value + 2, totalPage.value);
      i++
    ) {
      if (!pagesArray.includes(i)) {
        pagesArray.push(i)
      }
    }

    // 如果当前页和最后一页之间有间隔，添加省略号
    if (currentPageInternal.value + 3 < totalPage.value) {
      pagesArray.push('...')
    }

    // 始终显示最后一页
    if (!pagesArray.includes(totalPage.value)) {
      pagesArray.push(totalPage.value)
    }

    return pagesArray
  })

  watch(
    () => props.currentPage,
    (newVal) => {
      currentPageInternal.value = newVal
    }
  )

  watch(
    () => props.pageSize,
    (newVal) => {
      pageSizeInternal.value = newVal
    }
  )

  watch(
    () => props.total,
    (newVal) => {
      totalInternal.value = newVal
    }
  )

  const goToPage = (page) => {
    if (page === currentPageInternal.value) {
      return
    }
    currentPageInternal.value = page
    if (page <= 0) {
      page = 1
    }
    if (page > totalPage.value) {
      page = totalPage
    }
    emits('update:currentPage', page)
  }

  const jumpToPage = () => {
    const page = parseInt(inputPage.value)
    if (!isNaN(page) && page >= 1 && page <= totalPage.value) {
      emits('update:currentPage', page)
      inputPage.value = ''
    } else {
      inputPage.value = ''
    }
  }

  const changePageSize = () => {
    emits('update:pageSize', pageSizeInternal.value)
    currentPageInternal.value = 1
    emits('update:currentPage', currentPageInternal.value)
  }

  const inputPage = ref('')
</script>

<style scoped lang="scss">
  .bs-pagination {
    .pagination {
      --bs-link-color: var(--bs-primary);
      --bs-pagination-active-bg: var(--bs-primary);
      --bs-pagination-hover-color: var(--bs-primary);
      --bs-pagination-focus-color: var(--bs-primary);
      --bs-pagination-active-border-color: var(--bs-primary);
      --bs-pagination-focus-box-shadow: 0 0 0 0.25rem rgba(124, 57, 237, 0.25);

      .page-more {
        --bs-pagination-padding-x: 0.5rem;
      }

      .disabled {
        cursor: not-allowed;
      }

      .page-item-goto {
        margin: 0 10px;

        .form-control,
        .el-select {
          width: 50px;
          text-align: center;
        }
      }
    }
  }
</style>
