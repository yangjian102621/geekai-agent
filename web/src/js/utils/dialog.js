/**
 * Util lib functions
 */
import { ElMessage, ElMessageBox } from 'element-plus'
import {
  closeToast,
  setToastDefaultOptions,
  showConfirmDialog,
  showDialog,
  showLoadingToast,
  showToast,
} from 'vant'
import { isMobile } from './libs.js'

setToastDefaultOptions({ duration: 3000 })

export function showMessageOK(message) {
  if (isMobile()) {
    showToast({
      message: message,
      icon: 'passed',
      className: 'van-toast--success',
    })
  } else {
    ElMessage.success(message, closeLoading())
  }
}

export function showMessageInfo(message) {
  if (!isMobile()) {
    showToast({
      message: message,
      wordBreak: 'break-word',
    })
  } else {
    ElMessage.info(message, closeLoading())
  }
}

export function showMessageWarning(message) {
  if (isMobile()) {
    showToast({
      message: message,
      icon: 'warning-o',
      className: 'van-toast--warn',
    })
  } else {
    ElMessage.warning(message, closeLoading())
  }
}

export function showMessageError(message) {
  if (isMobile()) {
    showToast({
      message: message,
      icon: 'close-o',
      className: 'van-toast--fail',
    })
  } else {
    ElMessage.error(message, closeLoading())
  }
}

export function showLoading(message = '正在处理...', appendTo = 'body') {
  showLoadingToast({
    message: message,
    forbidClick: true,
    duration: 0,
    teleport: appendTo,
    zIndex: 9999,
    className: 'loading-toast',
  })
}

export function closeLoading() {
  closeToast()
}

export function showConfirm(title, message, onConfirm, onCancel) {
  onConfirm = onConfirm || function () {}
  onCancel = onCancel || function () {}
  if (isMobile()) {
    showConfirmDialog({
      title: title,
      message: message,
    })
      .then(() => onConfirm())
      .catch(() => onCancel())
  } else {
    ElMessageBox.confirm(message, title, {
      confirmButtonText: '确认',
      cancelButtonText: '取消',
      type: '',
      center: true,
    })
      .then(() => onConfirm())
      .catch(() => onCancel())
  }
}

export function showMessageBox(title, message) {
  if (isMobile()) {
    showDialog({ title: title, message: message })
  } else {
    ElMessageBox.alert(message, title, {
      confirmButtonText: '确认',
      center: true,
    })
  }
}
export function showComingSoon() {
  showMessageInfo('当前功能正在开发中，敬请期待！')
}
