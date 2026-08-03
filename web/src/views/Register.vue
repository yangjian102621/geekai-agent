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

              <h2 class="mb-4">注册</h2>

              <div class="mb-3">
                <small class="text-secondary">
                  已有账号？
                  <a
                    @click="router.push('/login')"
                    class="text-primary text-decoration-none cursor-pointer"
                    >直接登录</a
                  >
                </small>
              </div>

              <form
                @submit.prevent="handleSubmit"
                v-if="systemConfig.enabled_register"
              >
                <div class="mb-3">
                  <label for="account" class="form-label">账号</label>
                  <input
                    type="text"
                    class="form-control"
                    v-model="data.username"
                    :class="{ 'is-invalid': errors.username }"
                    id="account"
                    placeholder="请输入手机号或者邮箱"
                  />
                  <div class="invalid-feedback">{{ errors.username }}</div>
                </div>

                <div class="mb-3">
                  <label for="verificationCode" class="form-label"
                    >验证码</label
                  >
                  <div class="d-flex gap-2">
                    <input
                      type="number"
                      class="form-control"
                      v-model="data.code"
                      :class="{ 'is-invalid': errors.code }"
                      id="verificationCode"
                      placeholder="请输入验证码"
                    />
                    <send-msg :receiver="data.username" />
                  </div>
                  <div class="invalid-feedback">{{ errors.code }}</div>
                </div>

                <div class="mb-3">
                  <label for="password" class="form-label">密码</label>
                  <input
                    type="password"
                    class="form-control"
                    v-model="data.password"
                    :class="{ 'is-invalid': errors.password }"
                    id="password"
                    placeholder="请输入密码(8-16位)"
                  />
                  <div class="invalid-feedback">{{ errors.password }}</div>
                </div>

                <div class="mb-3">
                  <label for="confirmPassword" class="form-label"
                    >重复密码</label
                  >
                  <input
                    type="password"
                    class="form-control"
                    v-model="data.repass"
                    :class="{ 'is-invalid': errors.repass }"
                    id="confirmPassword"
                    placeholder="请输入重复密码(8-16位)"
                  />
                  <div class="invalid-feedback">{{ errors.repass }}</div>
                </div>

                <div class="mb-3">
                  <label for="inviteCode" class="form-label">邀请码</label>
                  <input
                    type="text"
                    class="form-control"
                    v-model="data.inviteCode"
                    id="inviteCode"
                    placeholder="请输入邀请码(可选)"
                  />
                </div>

                <button
                  type="submit"
                  class="btn btn-primary w-100 mt-3 login-btn"
                >
                  注 册
                </button>
              </form>

              <div class="tip-result left" v-else>
                <el-result
                  icon="error"
                  title="暂停注册"
                  style="--el-result-padding: 10px 30px 20px 30px"
                  sub-title="抱歉，当前系统未开放注册功能，请联系管理员添加账号！"
                >
                  <template #icon>
                    <div class="wechat-card">
                      <el-image :src="systemConfig.wechat_card_url" />
                    </div>
                  </template>
                </el-result>
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
  </div>
</template>

<script setup>
  import { onMounted } from 'vue'
  import { storeToRefs } from 'pinia'
  import Footer from '@/components/Footer.vue'
  import SendMsg from '@/components/SendMsg.vue'
  import { useRegisterStore } from '@/js/store/register'
  import { useRouter } from 'vue-router'

  const registerStore = useRegisterStore()
  const { loading, data, errors, systemConfig } = storeToRefs(registerStore)
  const { doRegister, handleSubmit, captchaRef, initialize } = registerStore
  const router = useRouter()

  onMounted(() => {
    initialize()
  })
</script>

<style scoped lang="scss">
  @use '@/assets/css/login.scss';
</style>
