import axios from 'axios'
import { reqReject, reqResolve, resReject, resResolve } from './interceptors'

export function createAxios(options = {}) {
  const service = axios.create({
    // 默认配置
    timeout: 12000,
    ...options,
  })
  service.interceptors.request.use(reqResolve, reqReject)
  service.interceptors.response.use(resResolve, resReject)
  return service
}

export const request = createAxios(
  {
    baseURL: import.meta.env.VITE_API,
    timeout: 12000,
  },
)
