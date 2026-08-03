<template>
  <el-form
    ref="formRef"
    :model="form"
    :rules="rules"
    label-width="100px"
    v-loading="loading"
  >
    <div class="alert alert-secondary text-center" role="alert">
      <span class="text-gray-600">当前可提现积分：</span>
      <span class="text-red-500 text-xl font-bold mr-2">{{ balance }}</span>
      <span class="text-gray-600">积分</span>
    </div>

    <el-form-item label="提现金额" prop="scores">
      <div class="flex">
        <el-input-number
          v-model="form.scores"
          :min="withdrawConfigs.score_to_rmb_ratio"
          :step="10"
          placeholder="请输入提现金额"
        />
        <span class="text-red-500 ml-2">￥{{ form.totalMoney }} 元</span>
      </div>
    </el-form-item>

    <el-form-item label="到账金额">
      <div class="flex">
        <span class="text-red-500 ml-2"> ￥{{ form.realMoney }} 元</span>
        <span class="text-gray-500 ml-2"
          >（手续费： ￥{{ form.fee }} 元，费率：{{
            withdrawConfigs.fee
          }}%）</span
        >
      </div>
    </el-form-item>

    <el-form-item label="提现方式" prop="method">
      <el-select
        v-model="form.method"
        placeholder="请选择提现方式"
        style="width: 100%"
      >
        <el-option label="支付宝" value="alipay" />
        <el-option label="微信" value="wxpay" />
      </el-select>
    </el-form-item>

    <el-form-item label="收款账号" prop="account">
      <el-input
        v-model="form.account"
        placeholder="请输入收款账号"
        maxlength="50"
      />
    </el-form-item>

    <el-form-item label="收款人姓名" prop="account_name">
      <el-input
        v-model="form.account_name"
        placeholder="请输入收款人真实姓名"
        maxlength="20"
      />
    </el-form-item>

    <el-form-item label="收款二维码" prop="qrcode">
      <div
        class="w-[60px] h-[60px] flex items-center justify-center border rounded-lg cursor-pointer mb-1"
      >
        <el-upload :show-file-list="false" :http-request="uploadQrcode">
          <img
            v-if="form.qrcode"
            :src="form.qrcode"
            class="w-[60px] h-[60px] rounded-lg"
          />
          <el-icon v-else>
            <Plus />
          </el-icon>
        </el-upload>
      </div>
    </el-form-item>

    <el-form-item label="备注" prop="remark">
      <el-input
        v-model="form.remark"
        type="textarea"
        placeholder="请输入备注信息（可选）"
        :rows="3"
        maxlength="200"
        show-word-limit
      />
    </el-form-item>

    <div class="withdraw-tips">
      <div class="alert alert-primary" role="alert">
        <h3 class="text-primary text-bold text-lg">提现说明:</h3>
        <div class="text-sm line-height-1">
          <p class="mb-1">
            1. 当前兑换比例为 {{ withdrawConfigs.score_to_rmb_ratio }}:1，即
            {{ withdrawConfigs.score_to_rmb_ratio }} 积分兑换 1 元
          </p>
          <p class="mb-1">
            2. 实际到账金额 = 提现金额 - 手续费，不同创作者手续费不一样
          </p>
          <p class="mb-1">3. 工作日申请，一般1-3个工作日到账</p>
          <p class="mb-1">4. 请确保收款信息准确，避免提现失败</p>
          <p class="mb-1">
            5. 如有疑问，请:
            <span
              class="text-primary cursor-pointer"
              @click="serviceQrcodeDialog = true"
            >
              联系客服
            </span>
          </p>
          <el-image-viewer
            v-if="serviceQrcodeDialog"
            :url-list="[serviceQrcode]"
            @close="serviceQrcodeDialog = false"
          />
        </div>
      </div>
    </div>

    <div class="w-100 flex justify-content-end pr-2 pt-3 border-top">
      <button
        class="btn btn-secondary mr-2"
        type="button"
        @click="emits('cancel')"
      >
        取消
      </button>
      <button class="btn btn-primary" type="button" @click="submit">
        提交
      </button>
    </div>
  </el-form>
</template>

<script setup>
  import { frontUploadFile } from '@/js/store/common.js'
  import { getSystemInfo } from '@/js/cache/session'
  import { closeLoading, showLoading } from '@/js/utils/dialog'
  import { httpGet, httpPost } from '@/js/utils/http'
  import { Plus } from '@element-plus/icons-vue'
  import { ElMessage } from 'element-plus'
  import { onMounted, reactive, ref } from 'vue'

  const props = defineProps({})
  const emits = defineEmits(['cancel', 'success'])

  const formRef = ref(null)
  const form = ref({
    scores: 0,
    totalMoney: 0,
    realMoney: 0,
    fee: 0,
    method: '',
    account: '',
    qrcode: '',
    remark: '',
    account_name: '',
  })
  const balance = ref(0)
  const withdrawConfigs = ref({})
  const serviceQrcode = ref('')
  const serviceQrcodeDialog = ref(false)
  const loading = ref(false)

  onMounted(() => {
    httpGet('/api/creator/info').then((res) => {
      balance.value = res.data.scores
      withdrawConfigs.value = res.data.withdraw_configs
      withdrawConfigs.value.fee = res.data.fee
      form.value.account_name = withdrawConfigs.value.name
      form.value.qrcode = withdrawConfigs.value.qrcode
      form.value.method = withdrawConfigs.value.method
      form.value.account = withdrawConfigs.value.account
    })

    getSystemInfo().then((res) => {
      serviceQrcode.value = res.wechat_card_url
    })
  })

  // 监听 form.scores 变化，自动计算 amount 和 realAmount
  watch(
    () => form.value.scores,
    (scores) => {
      if (scores) {
        // 先计算金额,转为数字类型
        form.value.totalMoney = Number(
          (scores / withdrawConfigs.value.score_to_rmb_ratio).toFixed(2)
        )
        // 计算手续费
        form.value.fee = (
          (form.value.totalMoney * withdrawConfigs.value.fee) /
          100
        ).toFixed(2)
        // 计算实际到账金额
        form.value.realMoney = (form.value.totalMoney - form.value.fee).toFixed(
          2
        )
      }
    },
    { immediate: true }
  )

  const rules = reactive({
    scores: [
      { required: true, message: '请输入提现金额', trigger: 'blur' },
      {
        type: 'number',
        min: 1000,
        message: '提现金额不能小于1元',
        trigger: 'blur',
      },
      {
        validator: (rule, value, callback) => {
          if (value > balance.value) {
            callback(new Error('提现金额不能超过可提现余额'))
          } else {
            callback()
          }
        },
        trigger: 'blur',
      },
    ],
    method: [{ required: true, message: '请选择提现方式', trigger: 'change' }],
    account: [
      { required: true, message: '请输入收款账号', trigger: 'blur' },
      { min: 5, max: 50, message: '收款账号长度在5-50个字符', trigger: 'blur' },
    ],
    account_name: [
      { required: true, message: '请输入收款人姓名', trigger: 'blur' },
      {
        min: 2,
        max: 20,
        message: '收款人姓名长度在2-20个字符',
        trigger: 'blur',
      },
    ],
    qrcode: [
      { required: true, message: '请上传收款二维码', trigger: 'change' },
    ],
  })

  // 上传收款二维码
  const uploadQrcode = (file) => {
    showLoading('上传中...')
    frontUploadFile(file, (res) => {
      form.value.qrcode = res.url
      closeLoading()
    })
  }

  const validate = () => {
    return formRef.value?.validate()
  }

  const submit = () => {
    validate().then((valid) => {
      if (valid) {
        loading.value = true
        httpPost('/api/creator/withdraw', {
          scores: form.value.scores,
          method: form.value.method,
          account: form.value.account,
          account_name: form.value.account_name,
          qrcode: form.value.qrcode,
          total_money: Number(form.value.totalMoney),
          real_money: Number(form.value.realMoney),
          fee: Number(form.value.fee),
          note: form.value.remark,
        })
          .then(() => {
            ElMessage.success('提现申请成功，请等待审核')
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

<style lang="scss" scoped></style>
