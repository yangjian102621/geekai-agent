import Storage from 'good-storage'
import { httpGet } from '../utils/http'
import { randString } from '../utils/libs'
import { useSharedStore } from './sharedata'

/**
 * storage handler
 */

const UserTokenKey = import.meta.env.VITE_KEY_PREFIX + 'X-USER-TOKEN'
const UserDataKey = 'USER_INFO_CACHE_KEY'
const AdminTokenKey = import.meta.env.VITE_KEY_PREFIX + 'X-ADMIN-TOKEN'
const adminDataKey = 'ADMIN_INFO_CACHE_KEY'
const systemInfoKey = 'SYSTEM_INFO_CACHE_KEY'
const systemInfoAdminKey = 'SYSTEM_INFO_ADMIN_CACHE_KEY'
const licenseKey = 'LICENSE_INFO_CACHE_KEY'
let sessionPromise = null
let adminSessionPromise = null
let systemInfoPromise = null
const sharedStore = useSharedStore()

export function getSessionId() {
  return randString(42)
}

export function getUserToken() {
  return Storage.get(UserTokenKey) ?? ''
}

export function setUserToken(token) {
  // 刷新 session 缓存
  Storage.set(UserTokenKey, token)
}

export function removeUserToken() {
  Storage.remove(UserTokenKey)
  Storage.remove(UserDataKey)
}

export async function checkSession() {
  const item = Storage.get(UserDataKey) ?? { expire: 0, data: null }
  if (item.expire > Date.now()) {
    return item.data
  }

  // 如果已经有正在进行的请求，直接返回该 Promise
  if (sessionPromise) {
    return sessionPromise
  }

  // 创建新的请求 Promise
  sessionPromise = new Promise(async (resolve) => {
    try {
      const res = await httpGet('/api/user/session')
      item.data = res.data
      item.expire = Date.now() + 1000 * 3 // 3秒后过期
      Storage.set(UserDataKey, item)
      resolve(item.data)
    } catch (e) {
      Storage.remove(UserDataKey)
      sharedStore.setIsLogin(false)
      resolve(null)
    } finally {
      // 请求完成后清除 Promise 缓存
      sessionPromise = null
    }
  })

  return sessionPromise
}

export function getAdminToken() {
  return Storage.get(AdminTokenKey) ?? ''
}

export function setAdminToken(token) {
  // 刷新 session 缓存
  Storage.set(AdminTokenKey, token)
}

export function removeAdminToken() {
  Storage.remove(AdminTokenKey)
  Storage.remove(adminDataKey)
}

export async function checkAdminSession() {
  const item = Storage.get(adminDataKey) ?? { expire: 0, data: null }
  if (item.expire > Date.now()) {
    return Promise.resolve(item.data)
  }

  // 如果已经有正在进行的请求，直接返回该 Promise
  if (adminSessionPromise) {
    return adminSessionPromise
  }

  // 创建新的请求 Promise
  adminSessionPromise = new Promise((resolve, reject) => {
    httpGet('/api/admin/session')
      .then((res) => {
        item.data = res.data
        item.expire = Date.now() + 1000 * 3 // 3秒后过期
        Storage.set(adminDataKey, item)
        resolve(item.data)
      })
      .catch((e) => {
        Storage.remove(adminDataKey)
        reject(e)
      })
      .finally(() => {
        // 请求完成后清除 Promise 缓存
        adminSessionPromise = null
      })
  })

  return adminSessionPromise
}

export async function getSystemInfo() {
  const item = Storage.get(systemInfoKey) ?? { expire: 0, data: null }
  if (item.expire > Date.now()) {
    return Promise.resolve(item.data)
  }

  // 如果已经有正在进行的请求，直接返回该 Promise
  if (systemInfoPromise) {
    return systemInfoPromise
  }

  // 创建新的请求 Promise
  systemInfoPromise = new Promise((resolve, reject) => {
    httpGet('/api/config/get?name=system')
      .then((res) => {
        item.data = res.data
        item.expire = Date.now() + 1000 * 10 // 10秒后过期
        Storage.set(systemInfoKey, item)
        resolve(item.data)
      })
      .catch((err) => {
        reject(err)
      })
      .finally(() => {
        // 请求完成后清除 Promise 缓存
        systemInfoPromise = null
      })
  })

  return systemInfoPromise
}

export async function getSystemInfoAdmin() {
  const item = Storage.get(systemInfoAdminKey) ?? { expire: 0, data: null }
  if (item.expire > Date.now()) {
    return Promise.resolve(item.data)
  }
  return new Promise((resolve, reject) => {
    httpGet('/api/admin/config/get?name=system')
      .then((res) => {
        item.data = res.data
        item.expire = Date.now() + 1000 * 10 // 10秒后过期
        Storage.set(systemInfoAdminKey, item)
        resolve(item.data)
      })
      .catch((err) => {
        reject(err)
      })
  })
}

export async function getLicense() {
  const item = Storage.get(licenseKey) ?? { expire: 0, data: null }
  if (item.expire > Date.now()) {
    return Promise.resolve(item.data)
  }
  return new Promise((resolve, reject) => {
    httpGet('/api/license/get')
      .then((res) => {
        item.data = res.data
        item.expire = Date.now() + 1000 * 10 // 10秒后过期
        Storage.set(licenseKey, item)
        resolve(item.data)
      })
      .catch((err) => {
        reject(err)
      })
  })
}
