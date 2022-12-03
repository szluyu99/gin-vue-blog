import axios from 'axios'
import { reqReject, reqResolve, resReject, resResolve } from './interceptors'

export function createAxios(options = {}) {
  const defaultOptions = {
    // baseURL: import.meta.env.VITE_BASE_API,
    timeout: 12000,
  }
  const service = axios.create({
    ...defaultOptions,
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
// export const backRequest = createAxios(
//   {
//     baseURL: '/back/api',
//     timeout: 12000,
//   },
// )
