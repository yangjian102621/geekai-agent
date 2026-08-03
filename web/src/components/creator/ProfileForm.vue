<template>
  <div v-loading="loading">
    <CustomTabs v-model="activeTab">
      <CustomTabPane name="base">
        <template #label>
          <span><i class="iconfont icon-user-fill"></i> 基础信息</span>
        </template>
        <el-form
          ref="formRef"
          :model="form"
          :rules="rules"
          label-width="100px"
          class="p-2 pb-0"
          label-position="top"
        >
          <el-form-item label="创作者Logo" prop="logo">
            <div
              class="w-[60px] h-[60px] flex items-center justify-center border rounded-lg cursor-pointer mb-1"
            >
              <el-upload
                :show-file-list="false"
                :http-request="uploadLogo"
                accept="image/*"
              >
                <img
                  v-if="form.logo"
                  :src="form.logo"
                  class="w-[60px] h-[60px] rounded-lg"
                />
                <el-icon v-else>
                  <Plus />
                </el-icon>
              </el-upload>
            </div>
          </el-form-item>
          <el-form-item label="创作者名称" prop="name">
            <el-input
              v-model="form.name"
              placeholder="请输入创作者名称"
              maxlength="20"
              show-word-limit
            />
          </el-form-item>
          <el-form-item prop="username">
            <template #label>
              <span>用户名</span>
              <span class="text-xs text-gray-500"> 【只支持英文和数字】 </span>
            </template>
            <el-input
              v-model="form.username"
              placeholder="请输入创作者用户名"
              maxlength="30"
              show-word-limit
            />
          </el-form-item>
          <el-form-item label="创作者简介" prop="description">
            <el-input
              v-model="form.description"
              type="textarea"
              placeholder="请输入创作者简介"
              :rows="4"
              maxlength="200"
              show-word-limit
            />
          </el-form-item>
        </el-form>
      </CustomTabPane>
      <CustomTabPane name="withdraw">
        <template #label>
          <div class="flex items-center justify-center">
            <i class="iconfont icon-recharge text-xl"></i>
            <span class="ml-1">收款信息</span>
          </div>
        </template>
        <el-form
          ref="withdrawFormRef"
          :model="form.withdraw_configs"
          :rules="withdrawRules"
          label-width="100px"
          class="p-2 pb-1"
          label-position="top"
        >
          <el-form-item label="真实姓名" prop="name">
            <el-input
              v-model="form.withdraw_configs.name"
              placeholder="请输入真实姓名"
              maxlength="20"
            />
          </el-form-item>
          <el-form-item label="联系手机" prop="mobile">
            <el-input
              v-model="form.withdraw_configs.mobile"
              placeholder="请输入联系手机号"
              maxlength="11"
            />
          </el-form-item>
          <el-form-item label="收款方式" prop="method">
            <el-select
              v-model="form.withdraw_configs.method"
              placeholder="请选择收款方式"
            >
              <el-option label="支付宝" value="alipay" />
              <el-option label="微信" value="wxpay" />
            </el-select>
          </el-form-item>
          <el-form-item label="收款账号" prop="account">
            <el-input
              v-model="form.withdraw_configs.account"
              placeholder="请输入收款账号"
              maxlength="50"
            />
          </el-form-item>
          <el-form-item label="收款二维码" prop="qrcode">
            <div
              class="w-[60px] h-[60px] flex items-center justify-center border rounded-lg cursor-pointer mb-1"
            >
              <el-upload :show-file-list="false" :http-request="uploadQrcode">
                <img
                  v-if="form.withdraw_configs.qrcode"
                  :src="form.withdraw_configs.qrcode"
                  class="w-[60px] h-[60px] rounded-lg"
                />
                <el-icon v-else>
                  <Plus />
                </el-icon>
              </el-upload>
            </div>
          </el-form-item>
        </el-form>
      </CustomTabPane>
    </CustomTabs>
    <div class="w-100 flex justify-content-end pr-2 pt-3 border-top">
      <button class="btn btn-secondary mr-2" @click="emits('cancel')">
        取消
      </button>
      <button class="btn btn-primary" @click="submit">提交</button>
    </div>
  </div>
</template>

<script setup>
  import { frontUploadFile } from '@/js/store/common.js'
  import { closeLoading, showLoading } from '@/js/utils/dialog'
  import { httpGet, httpPost } from '@/js/utils/http'
  import { Plus } from '@element-plus/icons-vue'
  import { ElMessage } from 'element-plus'
  import { reactive, ref } from 'vue'
  import CustomTabPane from '../CustomTabPane.vue'
  import CustomTabs from '../CustomTabs.vue'
  const props = defineProps({
    tab: {
      type: String,
      default: 'base',
    },
  })

  const formRef = ref(null)
  const withdrawFormRef = ref(null)
  const activeTab = ref(props.tab)
  const loading = ref(false)
  const form = ref({
    username: '',
    old_username: '', // 旧用户名
    withdraw_configs: {
      name: '',
      mobile: '',
      method: '',
      account: '',
      qrcode: '',
    },
  })

  const emits = defineEmits(['cancel', 'success'])

  onMounted(() => {
    // 获取创作者信息
    httpGet(`/api/creator/info`)
      .then((res) => {
        form.value = res.data
        form.value.old_username = res.data.username
      })
      .finally(() => {
        loading.value = false
      })
  })

  const rules = reactive({
    name: [
      { required: true, message: '请输入创作者名称', trigger: 'blur' },
      {
        min: 2,
        max: 20,
        message: '创作者名称长度在2-20个字符',
        trigger: 'blur',
      },
    ],
    username: [
      { required: true, message: '请输入用户名', trigger: 'blur' },
      {
        validator: (rule, value, callback) => {
          if (!/^[a-zA-Z0-9]+$/.test(value)) {
            callback(new Error('用户名只支持英文和数字'))
            return
          }
          if (value === form.value.old_username) {
            callback()
            return
          }
          httpGet(`/api/creator/check/username?username=${value}`)
            .then((res) => {
              if (res.data.message === '用户名可用') {
                callback()
              } else {
                callback(new Error(res.data.message))
              }
            })
            .catch((err) => {
              callback(new Error(err.message))
            })
        },
        trigger: 'blur',
      },
      {
        min: 2,
        max: 30,
        message: '用户名长度在2-30个字符',
        trigger: 'blur',
      },
    ],
    description: [
      { required: true, message: '请输入个人简介', trigger: 'blur' },
      {
        min: 10,
        max: 200,
        message: '个人简介长度在10-200个字符',
        trigger: 'blur',
      },
    ],
  })

  const withdrawRules = reactive({
    name: [
      { required: true, message: '请输入真实姓名', trigger: 'blur' },
      { min: 2, max: 20, message: '真实姓名长度在2-20个字符', trigger: 'blur' },
    ],
    mobile: [
      { required: true, message: '请输入联系手机号', trigger: 'blur' },
      {
        pattern: /^1[3-9]\d{9}$/,
        message: '请输入正确的手机号',
        trigger: 'blur',
      },
    ],
    account: [{ required: true, message: '请输入收款账号', trigger: 'blur' }],
    qrcode: [
      { required: true, message: '请上传收款二维码', trigger: 'change' },
    ],
    method: [{ required: true, message: '请选择收款方式', trigger: 'blur' }],
  })

  // 上传创作者Logo
  const uploadLogo = (file) => {
    showLoading('上传中...')
    frontUploadFile(file, (res) => {
      form.value.logo = res.url
      closeLoading()
    })
  }

  // 上传收款二维码
  const uploadQrcode = (file) => {
    showLoading('上传中...')
    frontUploadFile(file, (res) => {
      form.value.withdraw_configs.qrcode = res.url
      closeLoading()
    })
  }

  const validate = () => {
    return activeTab.value === 'base'
      ? formRef.value?.validate()
      : withdrawFormRef.value?.validate()
  }

  // 提交表单
  const submit = () => {
    validate().then((valid) => {
      if (valid) {
        loading.value = true
        httpPost(`/api/creator/update/profile`, {
          name: form.value.name,
          description: form.value.description,
          logo: form.value.logo,
          withdraw_configs: form.value.withdraw_configs,
          username: form.value.username,
        })
          .then(() => {
            ElMessage.success('更新成功')
            form.value.old_username = form.value.username
            emits('success')
          })
          .catch((err) => {
            ElMessage.error(err.message)
          })
          .finally(() => {
            loading.value = false
          })
      }
    })
  }
</script>

<style scoped></style>
