import axios from 'axios'
import { useAuthStore } from '@/store'

export const request = axios.create(
  {
    baseURL: import.meta.env.VITE_BASE_API,
    timeout: 12000,
  },
)

request.interceptors.request.use(
  // 请求成功拦截
  (config) => {
    if (config.noNeedToken) {
      return config
    }

    const { token } = useAuthStore()
    if (token) {
      config.headers.Authorization = `Bearer ${token}`
    }
    return config
  },
  // 请求失败拦截
  (error) => {
    return Promise.reject(error)
  },
)

request.interceptors.response.use(
  // 响应成功拦截
  (response) => {
    // 业务信息
    const responseData = response.data
    const { code, message, data } = responseData
    if (code !== 0) { // ! 与后端约定业务状态码
      if (data && message !== data) {
        window.$message.error(`${message} ${data}`)
      }
      else {
        window.$message.error(message)
      }
      console.error(responseData) // 控制台输出错误信息

      const authStore = useAuthStore()
      if (code === 1201) { // Token 存在问题
        authStore.toLogin()
        return
      }
      // 1202-Token 过期
      if (code === 1202 || code === 1203 || code === 1207) {
        authStore.forceOffline()
        return
      }
      return Promise.reject(responseData)
    }
    return Promise.resolve(responseData)
  },
  // 响应失败拦截
  (error) => {
    // 主要使用业务状态码决定状态, 一般不根据 HTTP 状态码进行操作
    const responseData = error.response?.data
    const { message, data } = responseData
    if (error.response.status === 500) {
      if (message && data) {
        window.$message.error(`${message} ${data}`)
      }
      else {
        window.$message.error('服务端异常')
      }
    }
    return Promise.reject(error)
  },
)
