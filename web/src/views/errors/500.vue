<template>
  <div class="container-fluid error-page">
    <div class="error-content error-500">
      <div class="error-number">500</div>
      <img :src="errImage" alt="500 illustration" class="img-fluid" />
      <h2 class="mb-4">服务器内部错误</h2>
      <p class="error-description">
        抱歉，服务器在处理您的请求时遇到了意外错误。
      </p>
      <div class="error-details">
        <p class="mb-2">您可以尝试以下操作：</p>
        <ul class="text-start">
          <li>刷新页面重试</li>
          <li>清除浏览器缓存</li>
          <li>稍后再试</li>
        </ul>
      </div>
      <div class="d-flex justify-content-center gap-3 mt-4">
        <button class="btn btn-danger" onclick="window.location.reload()">
          重新加载
        </button>
        <a href="/" class="btn btn-outline-danger"> 返回首页 </a>
      </div>
    </div>
  </div>
</template>

<script setup>
  import { replaceURL } from '@/js/utils/libs.js'

  const errImage = replaceURL('/images/error-page.png')
</script>

<style lang="scss" scoped>
  .error-page {
    height: 100vh;
    display: flex;
    align-items: center;
    justify-content: center;
    background: linear-gradient(135deg, #f5f7fa 0%, #d0bcf6 100%);
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

  .img-fluid {
    max-width: 250px;
    margin: 1rem 0;
  }

  @media (max-width: 768px) {
    .img-fluid {
      max-width: 200px;
      margin: 1rem 0;
    }
  }

  .error-500 .error-number {
    color: #dc3545;
  }

  .error-details {
    font-size: 0.9rem;
    color: #6c757d;
    margin-top: 1rem;
    padding: 1rem;
    background: rgba(0, 0, 0, 0.05);
    border-radius: 0.5rem;
  }
</style>
