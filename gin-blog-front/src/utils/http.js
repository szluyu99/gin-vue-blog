import axios, { AxiosError } from 'axios'
import { getToken } from './token'

export const request = axios.create(
  {
    baseURL: import.meta.env.VITE_API,
    timeout: 12000,
  },
)

request.interceptors.request.use(
  // 请求成功拦截
  (config) => {
    if (config.needToken) {
      const token = getToken()
      if (!token) {
        return Promise.reject(new AxiosError('当前没有登录，请先登录！', 401))
      }
      config.headers.Authorization = config.headers.Authorization || `Bearer ${token}`
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
    const { data } = response
    if (data?.code !== 0) { // 与后端约定业务状态码
      window.$message.error(data?.message)
    }
    return Promise.resolve(data)
  },
  // 响应失败拦截
  (error) => {
    return Promise.reject(error)
  },
)
