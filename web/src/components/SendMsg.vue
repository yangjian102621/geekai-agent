<template>
  <div class="d-flex">
    <button
      class="btn btn-primary"
      :class="{ 'btn-sm': size === 'small' }"
      :disabled="!canSend"
      @click="sendMsg"
      type="button"
    >
      {{ btnText }}
    </button>

    <captcha
      :type="captchaType"
      :show="enableVerify"
      @success="doSendMsg"
      ref="captchaRef"
    />
  </div>
</template>

<script setup>
  // 发送短信验证码组件
  import Captcha from '@/components/Captcha.vue'
  import {
    showLoading,
    showMessageError,
    showMessageOK,
  } from '@/js/utils/dialog'
  import { httpGet, httpPost } from '@/js/utils/http'
  import { validateEmail, validateMobile } from '@/js/utils/validate'
  import { onMounted } from 'vue'

  // eslint-disable-next-line no-undef
  const props = defineProps({
    receiver: String,
    size: String,
  })
  const btnText = ref('发送验证码')
  const canSend = ref(true)
  const captchaRef = ref(null)
  const enableVerify = ref(false)
  const captchaType = ref('slide')

  onMounted(async () => {
    const res = await httpGet('/api/config/captcha')
    enableVerify.value = res.data.enabled
    captchaType.value = res.data.type
  })

  const sendMsg = () => {
    if (!validateMobile(props.receiver) && !validateEmail(props.receiver)) {
      return showMessageError('请输入合法的手机号或者邮箱地址')
    }

    if (enableVerify.value) {
      captchaRef.value.loadCaptcha()
    } else {
      doSendMsg({})
    }
  }

  const doSendMsg = (data) => {
    if (!canSend.value) {
      return
    }
    showLoading('正在发送验证码')
    canSend.value = false
    httpPost('/api/sms/code', {
      receiver: props.receiver,
      key: data.key,
      dots: data.dots,
      x: data.x,
    })
      .then(() => {
        // 如果是邮箱
        if (validateEmail(props.receiver)) {
          showMessageOK('验证码发送成功，请到邮箱查看')
        } else {
          showMessageOK('验证码发送成功，请到手机查看')
        }
        let time = 60
        btnText.value = time
        const handler = setInterval(() => {
          time = time - 1
          if (time <= 0) {
            clearInterval(handler)
            btnText.value = '重新发送'
            canSend.value = true
          } else {
            btnText.value = time
          }
        }, 1000)
      })
      .catch((e) => {
        canSend.value = true
        showMessageError('验证码发送失败：' + e.message)
      })
  }
</script>

<style lang="scss" scoped>
  .btn {
    min-width: 120px;
  }
</style>
