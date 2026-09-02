import axios from 'axios'
import { useAppStore, useUserStore } from '@/store'

// 是否使用 mock 数据: 开启后不请求后端, 由 src/mock 返回假数据
const USE_MOCK = import.meta.env.VITE_USE_MOCK === 'true'

// 通用请求
export const baseRequest = axios.create(
  {
    baseURL: import.meta.env.VITE_API,
    timeout: 12000,
  },
)

baseRequest.interceptors.request.use(requestSuccess, requestFail)
baseRequest.interceptors.response.use(responseSuccess, responseFail)

// 前台请求
export const request = axios.create(
  {
    baseURL: `${import.meta.env.VITE_API}/front`,
    timeout: 12000,
  },
)

request.interceptors.request.use(requestSuccess, requestFail)
request.interceptors.response.use(responseSuccess, responseFail)

// mock 数据通过动态引入, 使其只出现在 mock 构建的产物中
// 需要在应用挂载前调用, 见 main.js
export async function setupMock() {
  if (!USE_MOCK) {
    return
  }
  const { mockAdapter } = await import('@/mock')
  baseRequest.defaults.adapter = mockAdapter
  request.defaults.adapter = mockAdapter
}

/**
 * 请求成功拦截
 * @param {import('axios').InternalAxiosRequestConfig} config
 */
function requestSuccess(config) {
  if (config.needToken) {
    const { token } = useUserStore()
    if (!token) {
      return Promise.reject(new axios.AxiosError('当前没有登录，请先登录！', 401))
    }
    config.headers.Authorization = config.headers.Authorization || `Bearer ${token}`
  }
  return config
}

/**
 * 请求失败拦截
 * @param {any} error
 */
function requestFail(error) {
  return Promise.reject(error)
}

/**
 * 响应成功拦截
 * @param {import('axios').AxiosResponse} response
 */
function responseSuccess(response) {
  const responseData = response.data
  const { code, message } = responseData
  if (code !== 0) { // 与后端约定业务状态码
    if (code === 1203) {
      // 移除 token
      const userStore = useUserStore()
      userStore.resetLoginState()
    }
    // $message 在 App.vue 的 onMounted 里才挂上, 早期失败的请求可能取不到
    window.$message?.error(message)
    return Promise.reject(responseData)
  }
  return Promise.resolve(responseData)
}

/**
 * 响应失败拦截
 * @param {any} error
 */
function responseFail(error) {
  const { code, message } = error
  // 401 只可能来自 requestSuccess 里自己构造的 AxiosError(没有 token)
  if (code === 401) {
    window.$message?.error(message)
    // 移除 token
    const userStore = useUserStore()
    userStore.resetLoginState()
    // 登录弹框
    const appStore = useAppStore()
    appStore.setLoginFlag(true)
  }
  else {
    // 超时(ECONNABORTED)、断网(ERR_NETWORK)、5xx 等以前只进 console, 用户看不到任何反馈
    window.$message?.error(networkErrorText(error))
  }
  return Promise.reject(error)
}

/**
 * 把 axios 的错误翻译成给用户看的文案
 * @param {any} error
 */
function networkErrorText(error) {
  if (error?.code === 'ECONNABORTED') {
    return '请求超时，请稍后重试'
  }
  if (error?.code === 'ERR_NETWORK') {
    return '网络异常，请检查网络连接'
  }
  const status = error?.response?.status
  return status ? `请求失败 (${status})` : '请求失败，请稍后重试'
}
