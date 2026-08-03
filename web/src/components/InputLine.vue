<template>
  <div class="input-line">
    <div class="input-box" :class="isFocus?'focus':''">
      <div class="tools me-2">
        <i class="iconfont icon-attachment-cl"></i>
      </div>
      <el-input
          class="input-custom"
          v-model="logo"
          placeholder="给 GeekAI 智能体发送消息"
          type="textarea"
          :autosize="{ minRows: 1,maxRows: 10 }"
          @focus="isFocus = true"
          @blur="isFocus = false"
          style=""
      />
      <div class="send ms-2">
        <button class="btn btn-primary btn-sm"><i class="iconfont icon-send"></i></button>
      </div>
    </div>
  </div>
</template>

<script setup>

const isFocus = ref(false)

</script>


<style scoped lang="scss">
.input-line {
  .input-box {
    flex-direction: row;
    justify-content: center;
    align-items: flex-end;
    padding: 10px 15px;
    display: flex;
    overflow: hidden;
    background: transparent;
    border-radius: 28px;
    border: 1px solid #e1e1e1;
    box-shadow: 0 0 .075rem .05rem rgba(33, 37, 41, 0.2);

    width: 100%;

    .tools {
      padding-bottom: 2px;

      .iconfont {
        font-size: 20px;
        padding: 2px;
        border-radius: 50%;
        cursor: pointer;

        &:hover {
          background-color: #f1f1f1;
        }
      }
    }

    .input-custom {
      padding-bottom: 4px;
    }

    &:hover {
      border: 1px solid rgba(var(--bs-primary-rgb), 0.5);
    }
  }

  .input-box.focus {
    box-shadow: 0 0 .25rem .05rem rgba(var(--bs-primary-rgb), 0.8);
  }
}
</style>