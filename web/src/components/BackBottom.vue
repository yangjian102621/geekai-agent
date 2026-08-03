<template>
  <div ref="myElement">
    <div v-if="showButton" @click="scrollToBottom" class="scroll-to-bottom text-white"
         :style="styles">
      <i class="iconfont icon-arrow-down"></i>
    </div>
  </div>
</template>

<script setup>

const props = defineProps({
  bottom: {
    type: Number,
    default: 30
  },
  right: {
    type: Number,
    default: 30
  },
  color: {
    type: String,
    default: 'var(--bs-primary)'
  },
  follow: {
    type: Boolean,
    default: false
  },
  target: {
    type: String,
    default: ''
  }
});

const showButton = ref(false);
const container = ref(null);
const myElement = ref(null);
const styles = ref({
  '--bg-color': props.color,
  'bottom': props.bottom + 'px'
});

onMounted(() => {
  if (props.target) {
    container.value = document.querySelector(props.target);
  } else {
    container.value = myElement.value.parentElement;
  }
  if (props.follow) {
    styles.value['left'] = (myElement.value.parentElement.offsetLeft + myElement.value.parentElement.offsetWidth) - props.right + 'px'
  } else {
    styles.value['right'] = props.right + 'px'
  }

  checkScroll();
  window.addEventListener('resize', checkScroll);
  container.value.addEventListener('scroll', checkScroll);
})

const scrollToBottom = () => {
  container.value.scrollTo({
    top: container.value.scrollHeight,
    behavior: 'smooth'
  });
};
const checkScroll = () => {
  showButton.value = container.value.scrollHeight - container.value.scrollTop > container.value.clientHeight + window.innerHeight;
}

</script>

<style scoped lang="scss">
.scroll-to-bottom {
  position: fixed;
  border: none;
  border-radius: 50%;
  width: 40px;
  height: 40px;
  cursor: pointer;
  outline: none;
  transition: opacity 0.3s;
  display: flex;
  justify-content: center;
  align-items: center;
  background-color: var(--bg-color);

  &:hover {
    opacity: 0.6;
  }

  .iconfont {
    font-size: 24px;
    color: #ffffff;
  }
}
</style>