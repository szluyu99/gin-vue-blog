import axios from 'axios'
import { useAuthStore } from '@/store'

// 是否使用 mock 数据: 开启后不请求后端, 由 src/mock 返回假数据
const USE_MOCK = import.meta.env.VITE_USE_MOCK === 'true'

export const request = axios.create(
  {
    baseURL: import.meta.env.VITE_BASE_API,
    timeout: 12000,
  },
)

// mock 数据通过动态引入, 使其只出现在 mock 构建的产物中
// 需要在应用挂载前调用, 见 main.js
export async function setupMock() {
  if (!USE_MOCK) {
    return
  }
  const { mockAdapter } = await import('@/mock')
  request.defaults.adapter = mockAdapter
}

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
      }
      // 1202-Token 过期, 1209-账号被禁用: 手上的 token 已经没用了, 直接踢下线
      else if (code === 1202 || code === 1203 || code === 1207 || code === 1209) {
        authStore.forceOffline()
      }
      // 这里必须 reject: 以前两个鉴权分支直接 return, 等于 resolve(undefined),
      // 调用方的 const { data } = await ... / resp.data.token 会二次抛错
      return Promise.reject(responseData)
    }
    return Promise.resolve(responseData)
  },
  // 响应失败拦截
  (error) => {
    // 主要使用业务状态码决定状态, 一般不根据 HTTP 状态码进行操作
    // 网络异常 / 超时的时候没有 response, 不能直接解构
    const { message, data } = error.response?.data ?? {}
    if (error.response?.status === 500) {
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
