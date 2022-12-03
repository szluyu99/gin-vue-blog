import axios from 'axios'
import { reqReject, reqResolve, resReject, resResolve } from './interceptors'

export function createAxios(options = {}) {
  const service = axios.create({
    // 默认配置
    timeout: 12000,
    ...options,
  })
  service.interceptors.request.use(reqResolve, reqReject) // 请求拦截
  service.interceptors.response.use(resResolve, resReject) // 响应拦截
  return service
}

export const request = createAxios({
  baseURL: import.meta.env.VITE_BASE_API,
})

// 如果还有其他请求, 可以再创建一个 axios 实例
// export const request_xxx = createAxios({
//   baseURL: 'xxx',
// })
