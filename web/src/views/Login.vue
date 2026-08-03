<template>
  <div class="container-fluid vh-100 login-page">
    <!-- 动态气泡背景 -->
    <div class="bubbles">
      <div v-for="n in 15" :key="n" class="bubble"></div>
    </div>

    <div class="container">
      <div
        class="row h-100 align-items-center justify-content-center position-relative"
      >
        <div class="col-12 col-lg-10">
          <div class="row form-container">
            <!-- Left Side - Form -->
            <div class="col-12 col-md-7 p-4 p-sm-5">
              <div class="mb-4">
                <a
                  @click="router.push('/')"
                  class="text-decoration-none text-secondary cursor-pointer"
                >
                  <i class="iconfont icon-home"></i> 首页
                </a>
              </div>

              <h2 class="mb-4">登录</h2>

              <div class="mb-3">
                <small class="text-secondary">
                  还没账号？
                  <a
                    @click="router.push('/register')"
                    class="text-primary text-decoration-none cursor-pointer"
                    >立即注册</a
                  >
                </small>
              </div>

              <form @submit.prevent="handleSubmit">
                <div class="mb-4">
                  <input
                    type="text"
                    class="form-control form-control-lg"
                    v-model="data.username"
                    :class="{ 'is-invalid': errors.username }"
                    placeholder="请输入手机号或者邮箱"
                  />
                  <div class="invalid-feedback">{{ errors.username }}</div>
                </div>

                <div class="mb-3">
                  <input
                    type="password"
                    class="form-control form-control-lg"
                    v-model="data.password"
                    :class="{ 'is-invalid': errors.password }"
                    placeholder="请输入密码(8-16位)"
                  />
                  <div class="invalid-feedback">{{ errors.password }}</div>
                </div>

                <button
                  type="submit"
                  class="btn btn-primary w-100 mt-3 login-btn btn-lg"
                  :disabled="loading"
                >
                  <span
                    class="spinner-border spinner-border-sm me-2"
                    v-if="loading"
                  ></span>
                  {{ loading ? 'Loading...' : '登 录' }}
                </button>
              </form>

              <!-- Social Login -->
              <div class="text-center mt-4">
                <p class="text-dark mb-3 fs-5">其他登录方式</p>
                <div class="d-flex justify-content-center gap-4">
                  <button
                    class="btn btn-light btn-lg social-btn"
                    @click="showComingSoon"
                  >
                    <i class="iconfont icon-github fs-4"></i>
                  </button>
                  <button
                    class="btn btn-light btn-lg social-btn"
                    @click="showComingSoon"
                  >
                    <i class="iconfont icon-wechat fs-4"></i>
                  </button>
                  <button
                    class="btn btn-light btn-lg social-btn"
                    @click="showComingSoon"
                  >
                    <i class="iconfont icon-google fs-4"></i>
                  </button>
                </div>
              </div>
            </div>

            <!-- Right Side - Logo -->
            <div
              class="col-md-5 d-none d-md-flex logo-section align-items-center justify-content-center text-center text-white p-5"
            >
              <div
                class="h-100 d-flex flex-column align-items-center justify-content-center"
              >
                <el-image
                  :src="systemConfig.logo"
                  alt="GeekAI Logo"
                  class="img-fluid mb-4"
                />
                <h5 class="mb-3">{{ systemConfig.slogan }}</h5>
                <div class="text-sm">
                  <Footer />
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>

    <captcha @success="doLogin" ref="captchaRef" :type="captchaConfig.type" />
  </div>
</template>

<script setup>
  import { onMounted } from 'vue'
  import { storeToRefs } from 'pinia'
  import Captcha from '@/components/Captcha.vue'
  import Footer from '@/components/Footer.vue'
  import { useLoginStore } from '@/js/store/login'
  import { checkSession } from '@/js/cache/session.js'
  import { showComingSoon } from '@/js/utils/dialog.js'
  import { useRouter } from 'vue-router'
  const router = useRouter()

  const loginStore = useLoginStore()
  const { loading, indexURL, data, errors, systemConfig, captchaConfig } =
    storeToRefs(loginStore)
  const { doLogin, handleSubmit, captchaRef, initialize } = loginStore

  onMounted(async () => {
    await initialize()
    const user = await checkSession()
    if (user) {
      router.push(indexURL.value)
    }
  })
</script>

<style scoped lang="scss">
  @use '@/assets/css/login.scss';
</style>
