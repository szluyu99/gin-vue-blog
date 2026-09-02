import { createPinia, setActivePinia } from 'pinia'
import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest'
import { useAppStore, useUserStore } from '@/store'
import { baseRequest } from '@/utils/http'

/*
拦截器是模块内的私有函数, 这里不直接调用它们, 而是替换 axios 的 adapter:
请求不会真正发出, 但完整走一遍 请求拦截 -> adapter -> 响应拦截 的链路。
*/

// adapter 收到的 config, 用来断言请求头
let lastConfig
// adapter 返回的响应体, 每个用例按需设置
let responseBody

const originalAdapter = baseRequest.defaults.adapter

function fakeAdapter(config) {
  lastConfig = config
  return Promise.resolve({
    data: responseBody,
    status: 200,
    statusText: 'OK',
    headers: {},
    config,
  })
}

describe('http 拦截器', () => {
  beforeEach(() => {
    setActivePinia(createPinia())
    lastConfig = undefined
    responseBody = { code: 0, message: 'OK', data: null }
    baseRequest.defaults.adapter = fakeAdapter
    // 源码里直接用了全局的 $message (由 App.vue 注入)
    window.$message = { error: vi.fn(), success: vi.fn() }
  })

  afterEach(() => {
    baseRequest.defaults.adapter = originalAdapter
  })

  it('业务码为 0 时把 data 层直接抛给调用方', async () => {
    responseBody = { code: 0, message: 'OK', data: { id: 1 } }

    const resp = await baseRequest.get('/config')

    // 注意: 拿到的是后端响应体本身, 不是 axios 的 response
    expect(resp).toEqual(responseBody)
    expect(window.$message.error).not.toHaveBeenCalled()
  })

  it('业务码非 0 时 reject 并弹错误提示', async () => {
    responseBody = { code: 9004, message: '数据库操作异常', data: null }

    await expect(baseRequest.get('/config')).rejects.toEqual(responseBody)
    expect(window.$message.error).toHaveBeenCalledWith('数据库操作异常')
  })

  it('业务码 1203 (token 失效) 会清掉登录态', async () => {
    const userStore = useUserStore()
    userStore.setToken('expired-token')
    responseBody = { code: 1203, message: 'token 已失效', data: null }

    await expect(baseRequest.get('/user/info')).rejects.toMatchObject({ code: 1203 })
    expect(userStore.token).toBeNull()
  })

  it('needToken 的请求会带上 Authorization', async () => {
    const userStore = useUserStore()
    userStore.setToken('my-token')

    await baseRequest.get('/user/info', { needToken: true })

    expect(lastConfig.headers.Authorization).toBe('Bearer my-token')
  })

  it('不需要 token 的请求不会带 Authorization', async () => {
    await baseRequest.get('/config')
    expect(lastConfig.headers.Authorization).toBeUndefined()
  })

  it('needToken 但没登录时请求不发出, 并弹出登录框', async () => {
    const appStore = useAppStore()
    expect(appStore.loginFlag).toBe(false)

    await expect(baseRequest.get('/user/info', { needToken: true })).rejects.toThrow('当前没有登录')

    // 请求被拦在 adapter 之前
    expect(lastConfig).toBeUndefined()
    expect(appStore.loginFlag).toBe(true)
    expect(window.$message.error).toHaveBeenCalledWith('当前没有登录，请先登录！')
  })

  it('超时和断网会给出对应提示', async () => {
    baseRequest.defaults.adapter = () =>
      Promise.reject(Object.assign(new Error('timeout'), { code: 'ECONNABORTED' }))
    await expect(baseRequest.get('/config')).rejects.toThrow('timeout')
    expect(window.$message.error).toHaveBeenCalledWith('请求超时，请稍后重试')

    baseRequest.defaults.adapter = () =>
      Promise.reject(Object.assign(new Error('offline'), { code: 'ERR_NETWORK' }))
    await expect(baseRequest.get('/config')).rejects.toThrow('offline')
    expect(window.$message.error).toHaveBeenCalledWith('网络异常，请检查网络连接')
  })

  it('带状态码的请求失败会把状态码带进提示', async () => {
    baseRequest.defaults.adapter = () =>
      Promise.reject(Object.assign(new Error('boom'), { code: 'ERR_BAD_RESPONSE', response: { status: 502 } }))

    await expect(baseRequest.get('/config')).rejects.toThrow('boom')
    expect(window.$message.error).toHaveBeenCalledWith('请求失败 (502)')
  })

  it('$message 还没挂载时不会抛 TypeError', async () => {
    window.$message = undefined
    responseBody = { code: 9004, message: '数据库操作异常', data: null }

    await expect(baseRequest.get('/config')).rejects.toMatchObject({ code: 9004 })
  })
})
