<template>
  <div class="charge-container p-3">
    <div class="mb-3 row text-center bg-light p-3 rounded-4">
      <div class="d-flex flex-column align-items-center">
        <i class="iconfont icon-reward fs-1"></i>
        <span class="fw-bold fs-4 d-flex align-items-center"
          >余额 :<span class="fs-3 ms-1">{{ user.scores }}</span></span
        >
      </div>
    </div>
    <el-divider class="mt-5 mb-5"><strong>兑换额度</strong></el-divider>

    <div class="mb-3 d-flex flex-column">
      <label class="d-flex fw-bold mb-2" style="width: 80px">兑换码：</label>
      <div class="d-flex mb-3">
        <input type="text" v-model="redeemCode" class="form-control" />
      </div>
      <button class="btn btn-primary" @click="doRedeem">兑 换</button>
    </div>

    <el-divider class="mt-5 mb-5"><strong>在线充值</strong></el-divider>

    <div class="mb-3 d-flex flex-row align-items-center">
      <label class="d-flex fw-bold" style="width: 110px">充值金额：</label>
      <div class="d-flex w-100">
        <el-select
          v-model="selectedPid"
          size="large"
          value-key="id"
          @change="handleProductChange()"
          clearable
        >
          <el-option
            v-for="item in products"
            :key="item.id"
            :label="item.name"
            :value="item.id"
          >
            <template #default>
              <div class="d-flex align-items-center justify-content-between">
                <span>{{ item.name }}</span>
                <el-tag type="danger" class="ms-2" size="small"
                  >￥{{ item.price }}
                </el-tag>
              </div>
            </template>
          </el-option>
        </el-select>
      </div>
    </div>

    <div class="mb-3 mt-5 row gx-3">
      <div class="col">
        <button
          class="w-100 bg-[#07C160] text-white hover:bg-[#06AD56] rounded-2 p-2"
          @click="wxPay"
        >
          <i class="iconfont icon-wechat-pay me-2"></i> 微信支付
        </button>
      </div>
      <div class="col">
        <button
          class="w-100 bg-[#1677ff] text-white hover:bg-[#0E5FD8] rounded-2 p-2"
          @click="alipay"
        >
          <i class="iconfont icon-alipay me-2"></i> 支付宝
        </button>
      </div>
    </div>

    <!--支付二维码-->
    <el-dialog
      v-model="showQrCode"
      :show-close="true"
      @close="showQrCode = false"
      style="width: 334px; height: 400px"
    >
      <template #title>
        <div class="d-flex align-items-center justify-content-center">
          <span>{{ title }}</span>
        </div>
      </template>
      <div class="qr-container">
        <div class="text-center">
          <el-tag type="danger" size="large"
            ><span class="text-2xl text-red-500">￥{{ totalFee }}</span></el-tag
          >
        </div>
        <el-image :src="qrImg" style="height: 300px; width: 300px" />
        <div class="qr-overlay" v-if="scanned">
          <div
            class="success-text d-flex align-items-center pt-2 pb-2 pl-4 pr-4"
          >
            <i class="iconfont icon-success me-2"></i> 扫码成功，请完成支付
          </div>
        </div>
      </div>
    </el-dialog>
  </div>
</template>

<script setup>
  import { checkSession } from '@/js/cache/session'
  import {
    closeLoading,
    showLoading,
    showMessageError,
  } from '@/js/utils/dialog'
  import { httpGet, httpPost } from '@/js/utils/http'
  import { ElMessage } from 'element-plus'
  import QRCode from 'qrcode'
  import { onMounted, ref } from 'vue'
  import { isMobile } from '@/js/utils/libs.js'

  const redeemCode = ref('')
  const products = ref([])
  const selectedPid = ref(0) // 选中的产品ID
  const user = ref({
    scores: 0,
  })
  const showQrCode = ref(false)
  const qrImg = ref('')
  const title = ref('请打开微信扫码支付')
  const totalFee = ref(0)
  const qrMargin = ref(2)
  const scanned = ref(false)

  onMounted(async () => {
    try {
      user.value = await checkSession()
      const res = await httpGet('/api/product/list')
      products.value = res.data
      selectedPid.value = products.value[0].id
      handleProductChange()
    } catch (e) {
      showMessageError(e.message)
    }
  })

  const wxPay = () => {
    title.value = '请打开微信扫码支付'
    qrMargin.value = 2
    generateOrder('wxpay')
  }

  const alipay = () => {
    title.value = '请打开支付宝扫码支付'
    qrMargin.value = 5
    generateOrder('alipay')
  }

  const generateOrder = (payWay) => {
    showLoading('正在生成支付订单...')
    // 生成支付订单
    httpPost('/api/payment/create', {
      pid: selectedPid.value,
      pay_way: payWay,
      domain: `${window.location.protocol}//${window.location.host}`,
      device: isMobile() ? 'mobile' : 'pc',
    })
      .then((res) => {
        closeLoading()

        if (isMobile()) {
          window.location.href = res.data.pay_url
        } else {
          QRCode.toDataURL(
            res.data.pay_url,
            { width: 300, height: 300, margin: qrMargin.value },
            (error, url) => {
              if (error) {
                console.error(error)
              } else {
                qrImg.value = url
              }
            }
          )
          // 查询订单状态
          if (handler.value) {
            clearTimeout(handler.value)
          }
          handler.value = setTimeout(() => queryOrder(res.data.order_no), 3000)
          showQrCode.value = true
        }
      })
      .catch((e) => {
        closeLoading()
        ElMessage.error('生成支付订单失败：' + e.message)
      })
  }

  const handler = ref(null)
  const queryOrder = async (orderNo) => {
    const res = await httpGet('/api/order/query?order_no=' + orderNo)
    if (res?.data.status === 1) {
      // 订单支付成功
      clearTimeout(handler.value)
      ElMessage.success('支付成功')
      showQrCode.value = false
      // 更新用户积分
      user.value.scores += res.data.credit
    } else {
      handler.value = setTimeout(() => queryOrder(orderNo), 3000)
    }
  }

  const doRedeem = () => {
    if (!redeemCode.value) {
      ElMessage.error('请输入兑换码')
      return
    }

    showLoading('正在兑换中...')
    httpPost('/api/redeem/verify', {
      code: redeemCode.value,
    })
      .then((res) => {
        ElMessage.success('兑换成功')
        user.value.scores += res.data
        redeemCode.value = ''
        closeLoading()
      })
      .catch((e) => {
        ElMessage.error('兑换失败：' + e.message)
        closeLoading()
      })
  }

  const handleProductChange = () => {
    totalFee.value = products.value.find(
      (item) => item.id === selectedPid.value
    )?.price
  }
</script>

<style lang="scss">
  .charge-container {
    .el-row {
      justify-content: center;
      margin-bottom: 10px;
    }

    .vip-icon {
      position: relative;
      top: 5px;
    }
  }

  .qr-container {
    position: relative;
    display: inline-block;
  }

  .qr-overlay {
    position: absolute;
    top: 0;
    left: 0;
    width: 100%;
    height: 100%;
    background-color: rgba(0, 0, 0, 0.7);
    display: flex;
    align-items: center;
    justify-content: center;
  }

  .el-tabs__nav {
    width: 100%;

    .el-tabs__item {
      width: 50%;
      margin: 0;
      padding: 0;
    }

    .el-tabs__item.is-active {
      background-color: rgba(var(--bs-primary-rgb), 0.1);
    }
  }
</style>
