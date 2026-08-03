<template>
  <div
    class="d-flex align-items-center px-3 py-1 rounded-3 item h-[98px]"
    :class="{ active: item.id === selectedApp.id }"
  >
    <div class="d-flex align-items-center relative">
      <el-image
        :src="item.icon"
        class="rounded"
        style="width: 70px; height: 70px"
        fit="cover"
      />
      <div
        class="checked-bg absolute left-0 top-0 w-[70px] h-[70px] bg-black/50 rounded flex justify-center items-center"
        v-if="item.id === selectedApp.id"
      >
        <i class="iconfont icon-success text-success bg-white checked"></i>
      </div>
    </div>
    <div class="flex flex-col p-2 overflow-hidden">
      <div class="d-flex flex-row align-items-center">
        <div class="me-2 text-truncate text-[#2b3674] text-h3 fw-bold">
          {{ item.name }}
        </div>
      </div>
      <div class="py-1 h-[42px]">
        <el-text line-clamp="2">
          {{ item.summary }}
        </el-text>
      </div>
    </div>
  </div>
</template>

<script setup>
  const props = defineProps({
    item: {
      type: Object,
      required: true,
    },
    selectedApp: {
      type: Object,
      required: true,
    },
  })
</script>

<style scoped lang="scss">
  .item {
    cursor: pointer;
    position: relative;
    background-color: rgba(var(--bs-primary-rgb), 0.06);
    &:hover {
      background-color: rgba(var(--bs-primary-rgb), 0.1);
    }
    .checked {
      font-size: 24px;
      height: 24px;
      line-height: 24px;
      border-radius: 50%;
    }
  }
  .item.active {
    background-color: rgba(var(--bs-primary-rgb), 0.1);
  }
</style>
