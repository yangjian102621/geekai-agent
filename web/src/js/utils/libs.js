// * +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
// * Copyright 2023 The Geek-AI Authors. All rights reserved.
// * Use of this source code is governed by a Apache-2.0 license
// * that can be found in the LICENSE file.
// * @Author yangjian102621@163.com
// * +++++++++++++++++++++++++++++++++++++++++++++++++++++++++++

/**
 * Util lib functions
 */
// generate a random string
export function randString(length) {
  const str = '0123456789abcdefghijklmnopqrstuvwxyz'
  const size = str.length
  let buf = []
  for (let i = 0; i < length; i++) {
    const rand = Math.random() * size
    buf.push(str.charAt(rand))
  }
  return buf.join('')
}

export function UUID() {
  let d = new Date().getTime()
  return 'xxxxxxxx-xxxx-4xxx-yxxx-xxxxxxxxxxxx'.replace(/[xy]/g, function (c) {
    const r = (d + Math.random() * 16) % 16 | 0
    d = Math.floor(d / 16)
    return (c === 'x' ? r : (r & 0x3) | 0x8).toString(16)
  })
}

// 判断是否是移动设备
export function isMobile() {
  const userAgent = navigator.userAgent
  const mobileRegex =
    /Android|webOS|iPhone|iPad|iPod|BlackBerry|IEMobile|Opera Mini|Mobile|mobile|CriOS/i
  return mobileRegex.test(userAgent)
}

// 格式化日期
export function dateFormat(timestamp, format) {
  if (!timestamp) {
    return ''
  } else if (timestamp < 9680917502) {
    timestamp = timestamp * 1000
  }
  let year, month, day, HH, mm, ss
  let time = new Date(timestamp)
  let timeDate
  year = time.getFullYear() // 年
  month = time.getMonth() + 1 // 月
  day = time.getDate() // 日
  HH = time.getHours() // 时
  mm = time.getMinutes() // 分
  ss = time.getSeconds() // 秒

  month = month < 10 ? '0' + month : month
  day = day < 10 ? '0' + day : day
  HH = HH < 10 ? '0' + HH : HH // 时
  mm = mm < 10 ? '0' + mm : mm // 分
  ss = ss < 10 ? '0' + ss : ss // 秒

  switch (format) {
    case 'yyyy':
      timeDate = String(year)
      break
    case 'yyyy-MM':
      timeDate = year + '-' + month
      break
    case 'yyyy-MM-dd':
      timeDate = year + '-' + month + '-' + day
      break
    case 'yyyy/MM/dd':
      timeDate = year + '/' + month + '/' + day
      break
    case 'yyyy-MM-dd HH:mm:ss':
      timeDate = year + '-' + month + '-' + day + ' ' + HH + ':' + mm + ':' + ss
      break
    case 'HH:mm:ss':
      timeDate = HH + ':' + mm + ':' + ss
      break
    case 'MM':
      timeDate = String(month)
      break
    default:
      timeDate = year + '-' + month + '-' + day + ' ' + HH + ':' + mm + ':' + ss
      break
  }
  return timeDate
}

export function formatTime(time) {
  const minutes = Math.floor(time / 60)
  const seconds = Math.floor(time % 60)
  return `${minutes}:${seconds < 10 ? '0' : ''}${seconds}`
}

// 删除数组中指定的元素
export function removeArrayItem(array, value, compare) {
  if (typeof compare !== 'function') {
    compare = function (v1, v2) {
      return v1 === v2
    }
  }
  return array.filter((item) => !compare(item, value))
}

// 拷贝对象
export function copyObj(origin) {
  return JSON.parse(JSON.stringify(origin))
}

// 字符串截取
export function substr(str, length) {
  if (!str) {
    return ''
  }

  let result = ''
  let count = 0

  for (let i = 0; i < str.length; i++) {
    const char = str.charAt(i)
    const charCode = str.charCodeAt(i)

    // 判断字符是否为中文字符
    if (charCode >= 0x4e00 && charCode <= 0x9fff) {
      // 中文字符算两个字符
      count += 2
    } else {
      count++
    }

    if (count <= length) {
      result += char
    } else {
      result += '...'
      break
    }
  }

  return result
}

export function isImage(url) {
  const expr = /\.(jpg|jpeg|png|gif|bmp|svg)$/i
  return expr.test(url)
}

// 将文件大小转成字符
export function FormatFileSize(bytes) {
  if (bytes === 0) return '0 Bytes'
  const k = 1024
  const sizes = ['Bytes', 'KiB', 'MiB', 'GiB', 'TiB', 'PiB']
  const i = Math.floor(Math.log(bytes) / Math.log(k))
  return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + ' ' + sizes[i]
}

// 判断数组中是否包含某个元素
export function arrayContains(array, value, compare) {
  if (!array) {
    return false
  }
  if (typeof compare !== 'function') {
    compare = function (v1, v2) {
      return v1 === v2
    }
  }
  return array.some((item) => compare(item, value))
}

// 获取资源真实路径
export function replaceURL(url) {
  if (url.includes('http')) {
    return url
  }
  const prefix = location.protocol + '//' + location.host
  return prefix + url
}

// 获取文件图标
export function GetFileIcon(ext) {
  const files = {
    '.docx': 'doc.png',
    '.doc': 'doc.png',
    '.xls': 'xls.png',
    '.xlsx': 'xls.png',
    '.csv': 'xls.png',
    '.ppt': 'ppt.png',
    '.pptx': 'ppt.png',
    '.md': 'md.png',
    '.pdf': 'pdf.png',
    '.sql': 'sql.png',
    '.mp3': 'mp3.png',
    '.wav': 'mp3.png',
    '.mp4': 'mp4.png',
    '.avi': 'mp4.png',
  }
  if (files[ext]) {
    return '/images/ext/' + files[ext]
  }

  return '/images/ext/file.png'
}

// 获取文件扩展名类型
export function getFileExtType(ext) {
  return ext.replace('.', '').toUpperCase()
}

// 根据文件的后缀名判断文件的类型：text,image,audio,file 文档类型，图片类型，音频类型，其他类型
export function getFileType(fileName) {
  const ext = fileName.split('.').pop().toLowerCase()
  const textExtensions = ['txt', 'md', 'pdf', 'doc', 'docx']
  const imageExtensions = ['jpg', 'jpeg', 'png', 'gif', 'bmp']
  const audioExtensions = ['mp3', 'wav', 'ogg']

  if (textExtensions.includes(ext)) {
    return 'text'
  } else if (imageExtensions.includes(ext)) {
    return 'image'
  } else if (audioExtensions.includes(ext)) {
    return 'audio'
  }
  return 'file'
}

// 根据文件的URL获取文件的扩展名
export function getFileExt(url) {
  return '.' + url.split('.').pop().toLowerCase()
}

export function processContent(content) {
  if (!content) {
    return ''
  }

  // 如果是图片链接地址，则直接替换成图片标签
  const linkRegex = /(https?:\/\/\S+)/g
  const links = content.match(linkRegex)
  if (links) {
    for (let link of links) {
      if (isImage(link)) {
        const index = content.indexOf(link)
        if (content.substring(index - 1, 2) !== ']') {
          content = content.replace(link, '\n![](' + link + ')\n')
        }
      }
    }
  }
  // 处理推理标签
  if (content.includes('<think>')) {
    content = content.replace(/<think>(.*?)<\/think>/gs, (match, content) => {
      if (content.length > 10) {
        return `<blockquote>${content}</blockquote>`
      }
      return ''
    })
    content = content.replace(/<think>(.*?)$/gs, (match, content) => {
      if (content.length > 10) {
        return `<blockquote>${content}</blockquote>`
      }
      return ''
    })
  }

  // 支持 \[ 公式标签
  content = content.replace(/\\\[/g, '$$').replace(/\\\]/g, '$$')
  content = content.replace(
    /\\\(\\boxed\{(\d+)\}\\\)/g,
    '<span class="boxed">$1</span>'
  )
  return content
}

// 判断是否微信浏览器
export function isWechatBrowser() {
  const userAgent = navigator.userAgent
  const wechatRegex = /MicroMessenger/i
  return wechatRegex.test(userAgent)
}

// 给字符串打码，隐藏中间三分之一部分信息
export function maskString(str) {
  if (!str) return ''

  const len = str.length
  if (len <= 3) return str

  const oneThird = Math.floor(len / 3)
  const start = str.substring(0, oneThird)
  const end = str.substring(len - oneThird)
  const maskLen = len - 2 * oneThird
  const mask = '*'.repeat(maskLen)

  return start + mask + end
}

// 生成随机颜色，10个深色系
export function getRandomColor() {
  const colors = [
    '#16a34a',
    '#2563eb',
    '#ef4444',
    '#f59e0b',
    '#8b5cf6',
    '#f43f5e',
    '#0ea5e9',
    '#d946ef',
    '#f97316',
    '#10b981',
    '#8b5cf6',
  ]
  return colors[Math.floor(Math.random() * colors.length)]
}
