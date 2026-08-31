import axios from 'axios'
import { useAppStore, useUserStore } from '@/store'

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
    window.$message.error(message)
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
  if (code === 401) {
    window.$message.error(message)
    // 移除 token
    const userStore = useUserStore()
    userStore.resetLoginState()
    // 登录弹框
    const appStore = useAppStore()
    appStore.setLoginFlag(true)
  }
  return Promise.reject(error)
}
