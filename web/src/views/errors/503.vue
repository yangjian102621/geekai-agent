<template>
  <div class="container-fluid error-page">
    <div class="error-content">
      <div class="error-number">503</div>
      <img :src="errImage" alt="503 illustration" class="img-fluid">
      <h2 class="mb-4">服务暂时不可用</h2>
      <p class="error-description">
        抱歉，服务器当前正在维护或负载过高，请稍后再试。
      </p>
      <div class="d-flex justify-content-center gap-3">
        <button class="btn btn-primary" onclick="window.history.back()">
          返回上一页
        </button>
        <a href="/" class="btn btn-outline-primary">
          返回首页
        </a>
      </div>
    </div>
  </div>
</template>

<script setup>
import {ref} from "vue";

const errImage = ref('images/logo.png')
</script>

<style lang="scss" scoped>
.error-page {
  height: 100vh;
  display: flex;
  align-items: center;
  justify-content: center;
  background: linear-gradient(135deg, #f5f7fa 0%, #D0BCF6 100%);
}

.error-content {
  text-align: center;
  padding: 2rem;
  background: rgba(255, 255, 255, 0.9);
  border-radius: 1rem;
  box-shadow: 0 0 20px rgba(0, 0, 0, 0.1);
}

.error-number {
  font-size: 8rem;
  font-weight: 700;
  color: #343a40;
  text-shadow: 2px 2px 4px rgba(0, 0, 0, 0.1);
}

.error-description {
  margin: 1.5rem 0;
  color: #6c757d;
}

@media (max-width: 768px) {
  .img-fluid {
    max-width: 200px;
    margin: 1rem 0;
  }
}
</style>