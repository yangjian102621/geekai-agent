<template>
  <div class="container-fluid vh-100 login-container">
    <!-- 动态气泡背景 -->
    <div class="bubbles">
      <div v-for="n in 15" :key="n" class="bubble"></div>
    </div>
    <!-- 水波纹效果 -->
    <div class="wave-group">
      <div class="wave wave1"></div>
      <div class="wave wave2"></div>
      <div class="wave wave3"></div>
      <div class="wave wave4"></div>
    </div>
    <div class="background-overlay"></div>
    <div
      class="row h-100 align-items-center justify-content-center position-relative"
    >
      <div class="inner-container col-11 col-sm-8 col-md-6 col-lg-4">
        <div class="card shadow-lg border-0">
          <div class="card-body p-5">
            <!-- Logo and Title -->
            <div class="d-flex flex-column align-items-center mb-4">
              <img
                :src="systemConfig.logo"
                class="logo rounded img-fluid w-[100px] h-[100px]"
                alt="GeekAI-Agent"
              />
              <h3 class="fw-bold">欢迎登录</h3>
              <div class="text-center text-sm text-muted">
                登录 GeekAI-Agent 管理控制台
              </div>
            </div>

            <!-- Login Form -->
            <form @submit.prevent="handleSubmit">
              <div class="mb-4">
                <input
                  type="text"
                  class="form-control"
                  id="username"
                  v-model="data.username"
                  :class="{ 'is-invalid': errors.username }"
                  placeholder="Enter your username"
                />
                <div class="invalid-feedback">{{ errors.username }}</div>
              </div>

              <div class="mb-4">
                <input
                  type="password"
                  class="form-control"
                  id="password"
                  v-model="data.password"
                  :class="{ 'is-invalid': errors.password }"
                  placeholder="Enter your password"
                />
                <div class="invalid-feedback">{{ errors.password }}</div>
              </div>

              <button
                type="submit"
                class="btn btn-primary w-100 mb-4 login-btn"
                :disabled="loading"
              >
                <span
                  class="spinner-border spinner-border-sm me-2"
                  v-if="loading"
                ></span>
                {{ loading ? 'Loading...' : '登 录' }}
              </button>
            </form>
          </div>
        </div>

        <div class="text-white footer-container">
          <Footer />
        </div>
      </div>
    </div>

    <captcha
      v-if="captchaConfig.enabled"
      @success="doLogin"
      :type="captchaConfig.type"
      ref="captchaRef"
    />
  </div>
</template>

<script setup>
  import { onMounted } from 'vue'
  import { storeToRefs } from 'pinia'
  import { useRouter } from 'vue-router'
  import Captcha from '@/components/Captcha.vue'
  import Footer from '@/components/Footer.vue'
  import { useAdminLoginStore } from '@/js/store/admin/login'
  import { checkAdminSession } from '@/js/cache/session.js'

  const router = useRouter()
  const indexURL = '/admin/dashboard'
  checkAdminSession().then(() => {
    router.push(indexURL)
  })

  const loginStore = useAdminLoginStore()
  const { loading, data, errors, systemConfig, captchaConfig, captchaRef } =
    storeToRefs(loginStore)
  const { doLogin, handleSubmit, loadConfigs } = loginStore

  onMounted(() => {
    loadConfigs()
  })
</script>

<style scoped lang="scss">
  @use '@/assets/css/admin/login.scss';

  .footer-container {
    position: fixed;
    bottom: 0;
    left: 0;
    width: 100%;
    text-align: center;
    padding: 15px 0;
  }
</style>
