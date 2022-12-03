import type { AxiosRequestConfig } from 'axios'

interface RequestConfig extends AxiosRequestConfig {
  /** 接口是否需要token */
  needToken?: boolean
  /** 接口是否需要错误提醒 */
  noNeedTip?: boolean
}

interface ErrorResolveResponse {
  code?: number | string
  message: string
  data?: any
}
