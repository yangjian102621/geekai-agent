<template>
  <div
    class="w-full max-w-[600px] min-w-[300px]"
    :class="
      isSubmitting ? 'opacity-70 cursor-not-allowed pointer-events-none' : ''
    "
  >
    <div class="bg-[#f5f5f5] rounded-[10px] p-3 mb-2.5">
      <div class="text-lg text-gray-800 mb-4 leading-normal text-left">
        {{ questionData.title }}
      </div>
      <div class="flex flex-col gap-3">
        <div
          v-for="(answer, index) in questionData.answers"
          :key="index"
          class="w-full rounded-lg text-[15px] font-normal text-gray-800 bg-white text-left py-2 px-3 transition-all shadow-sm flex items-center justify-start cursor-pointer hover:bg-white hover:-translate-y-0.5 hover:shadow-md hover:text-blue-500 active:translate-y-0 !font-semibold"
          @click="!isSubmitting && handleSelectAnswer(answer)"
        >
          {{ answer }}
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
  import { ref } from 'vue'
  import { useChatStore } from '@/js/store/chat.js'

  const props = defineProps({
    questionData: {
      type: Object,
      required: true,
      default: () => ({
        title: '',
        answers: [],
      }),
    },
    messageId: {
      type: Number,
      default: 0,
    },
  })

  const emit = defineEmits(['submit'])

  const chatStore = useChatStore()
  const selectedAnswer = ref('')
  const isSubmitting = ref(false)

  // 选择答案
  const handleSelectAnswer = (answer) => {
    if (isSubmitting.value) return

    selectedAnswer.value = answer
    isSubmitting.value = true

    // 将选中的内容发送给后端
    const originalPrompt = chatStore.prompt
    chatStore.prompt = answer
    chatStore.sendMessage(props.messageId || 0)

    // 恢复原始 prompt（如果有的话）
    setTimeout(() => {
      chatStore.prompt = originalPrompt
    }, 100)

    // 触发提交事件
    emit('submit', answer)
  }
</script>
