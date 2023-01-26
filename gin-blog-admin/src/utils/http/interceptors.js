import { resolveResError } from './helper'
import { getToken } from '@/utils'

// 发送请求前的操作
export function reqResolve(config) {
  // 发送请求前的配置: 不需要 token, 直接发送即可
  if (config.noNeedToken)
    return config

  // 获取本地存储的 token, 不存在则登录过期
  const token = getToken()
  if (!token)
    return Promise.reject({ code: 401, message: '登录已过期，请重新登陆！' })

  // * JWT Bearer 认证: 在请求头 Authorization 添加 token
  config.headers.Authorization = config.headers.Authorization || `Bearer ${token}`
  return config
}

// 请求错误的操作
export function reqReject(err) {
  // console.log('请求失败')
  return Promise.reject(err)
}

// 响应成功的操作
export function resResolve(resp) {
  // console.log('响应成功', resp)
  // TODO: 处理不同的 resp.headers
  // 响应结构 http://axios-js.com/zh-cn/docs/#响应结构
  const { data, status, config, statusText } = resp
  // 处理请求异常: 自定义状态码异常
  if (data?.code !== 0) {
    const code = data?.code ?? status
    // 根据 code 处理对应操作, 返回处理后的 msg
    const msg = resolveResError(code, data?.message ?? statusText)
    // 需要错误提醒
    !config?.noNeedTip && window.$message.error(msg)
    return Promise.reject({ code, msg, error: data || resp })
  }
  // 请求正常
  return Promise.resolve(data)
}

// 响应错误的操作
export function resReject(err) {
  // console.log('响应错误', err)
  if (!err || !err.response) {
    const code = err?.code
    // 根据 code 处理对应操作, 返回处理后的 msg
    const msg = resolveResError(code, err.message)
    // 全局弹窗提示错误
    window.$message?.error(msg)
    return Promise.reject({ code, msg, err })
  }
  const { data, status, config } = err.response
  const code = status ?? data?.code
  // TODO: 优化, 前后端联动
  const msg = resolveResError(code, err.message)
  // 需要错误提醒
  !config?.noNeedTip && $message.error(msg)
  return Promise.reject({ code, msg, error: err.response?.data || err.response })
}
