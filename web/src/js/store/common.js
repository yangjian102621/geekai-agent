import Compressor from 'compressorjs'
import { showMessageError, showMessageOK } from '../utils/dialog.js'
import { httpPost } from '../utils/http.js'

// 校验函数
export const validateForm = (form, rules, errorObj) => {
  let valid = true
  Object.keys(rules).forEach((key) => {
    const rule = rules[key]
    let required = rule.required
    if (typeof required === 'function') {
      required = required(form)
    }
    // 1. 必填校验
    if (required && !form[key]) {
      errorObj[key] = rule.message
      valid = false
    }
    // 2. 其他自定义校验
    else if (rule.validator && !rule.validator(form[key])) {
      valid = false
    } else {
      errorObj[key] = ''
    }
  })
  return valid
}

export const adminUploadFile = (file, callback) => {
  doUploadFile('/api/admin/upload', file, callback)
}

export const frontUploadFile = (file, callback) => {
  if (file.size / 1024 / 1024 > 5) {
    showMessageError('文件大小不能超过5MB')
    return
  }
  doUploadFile('/api/file/upload', file, callback)
}

export const doUploadFile = (url, file, callback) => {
  // 如果文件是图片，则压缩图片并上传
  if (file.file.type.startsWith('image/')) {
    new Compressor(file.file, {
      quality: 0.6,
      success(result) {
        const formData = new FormData()
        formData.append('file', result, result.name)
        // 执行上传操作
        httpPost(url, formData)
          .then((res) => {
            callback(res.data)
            showMessageOK('上传成功')
          })
          .catch((e) => {
            showMessageError('上传失败:' + e.message)
          })
      },
      error(e) {
        showMessageError('上传失败:' + e.message)
      },
    })
  } else {
    const formData = new FormData()
    formData.append('file', file.file, file.file.name)
    // 执行上传操作
    httpPost(url, formData)
      .then((res) => {
        callback(res.data)
        showMessageOK('上传成功')
      })
      .catch((e) => {
        showMessageError('上传失败:' + e.message)
      })
  }
}
