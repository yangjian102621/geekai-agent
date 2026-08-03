<template>
  <div class="password-input-wrapper">
    <div class="input-group" :class="{ 'is-focus': isFocus }">
      <input
        :type="showPassword ? 'text' : 'password'"
        :value="modelValue"
        :placeholder="placeholder"
        :disabled="disabled"
        :readonly="readonly"
        class="form-control"
        :class="{ 'is-invalid': error }"
        @input="handleInput"
        @focus="isFocus = true"
        @blur="handleBlur"
        @keyup.enter="$emit('enter', $event)"
      />
      <button
        type="button"
        class="btn password-toggle-btn"
        :class="{ active: showPassword }"
        @click="togglePassword"
        tabindex="-1"
      >
        <i
          :class="
            showPassword ? 'iconfont icon-eye-open' : 'iconfont icon-eye-close'
          "
          class="password-icon"
        ></i>
      </button>
    </div>
    <div v-if="error" class="invalid-feedback d-block">
      {{ error }}
    </div>
  </div>
</template>

<script setup>
  import { ref } from 'vue'

  const props = defineProps({
    modelValue: {
      type: String,
      default: '',
    },
    placeholder: {
      type: String,
      default: '',
    },
    disabled: {
      type: Boolean,
      default: false,
    },
    readonly: {
      type: Boolean,
      default: false,
    },
    error: {
      type: String,
      default: '',
    },
  })

  const emit = defineEmits(['update:modelValue', 'enter', 'focus', 'blur'])

  const showPassword = ref(false)
  const isFocus = ref(false)

  const togglePassword = () => {
    showPassword.value = !showPassword.value
  }

  const handleInput = (event) => {
    emit('update:modelValue', event.target.value)
  }

  const handleBlur = (event) => {
    isFocus.value = false
    emit('blur', event)
  }
</script>

<style scoped lang="scss">
  .password-input-wrapper {
    width: 100%;

    .input-group {
      position: relative;
      display: flex;
      align-items: center;
      border: 1px solid #ced4da;
      border-radius: 0.375rem;
      background-color: #fff;
      transition: all 0.15s ease-in-out;

      &:hover {
        border-color: rgba(var(--bs-primary-rgb), 0.5);
      }

      &.is-focus {
        border-color: var(--bs-primary);
        box-shadow: 0 0 0 0.2rem rgba(var(--bs-primary-rgb), 0.25);
      }

      .form-control {
        border: none;
        border-radius: 0.375rem;
        padding: 0.375rem 3rem 0.375rem 0.75rem;
        flex: 1;
        background-color: transparent;
        box-shadow: none;

        &:focus {
          border: none;
          box-shadow: none;
          outline: none;
        }

        &:disabled,
        &[readonly] {
          background-color: #e9ecef;
          opacity: 1;
        }

        &.is-invalid {
          color: var(--bs-danger);
        }
      }

      .password-toggle-btn {
        position: absolute;
        right: 0;
        top: 0;
        bottom: 0;
        display: flex;
        align-items: center;
        justify-content: center;
        width: 2.5rem;
        padding: 0;
        border: none;
        background: transparent;
        color: #6c757d;
        cursor: pointer;
        transition: color 0.15s ease-in-out;
        z-index: 5;

        &:hover {
          color: var(--bs-primary);
          background-color: transparent;
        }

        &:focus {
          box-shadow: none;
          outline: none;
        }

        &.active {
          color: var(--bs-primary);
        }

        .password-icon {
          font-size: 1rem;
          line-height: 1;
        }
      }
    }

    .invalid-feedback {
      margin-top: 0.25rem;
      font-size: 0.875rem;
      color: var(--bs-danger);
    }
  }

  // 自定义主题色（紫色）支持
  :root {
    --bs-primary: #6f42c1;
    --bs-primary-rgb: 111, 66, 193;
    --bs-danger: #dc3545;
  }

  // 如果项目使用了紫色主题，可以这样覆盖
  .password-input-wrapper .input-group.is-focus {
    border-color: #6f42c1;
    box-shadow: 0 0 0 0.2rem rgba(111, 66, 193, 0.25);
  }

  .password-input-wrapper .input-group:hover {
    border-color: rgba(111, 66, 193, 0.5);
  }

  .password-input-wrapper .password-toggle-btn:hover,
  .password-input-wrapper .password-toggle-btn.active {
    color: #6f42c1;
  }
</style>
