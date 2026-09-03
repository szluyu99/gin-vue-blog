import { createPinia, setActivePinia } from 'pinia'
import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest'
import { useAuthStore } from '@/store'
import { request } from '@/utils/http'

/*
拦截器是模块内的私有函数, 这里不直接调用它们, 而是替换 axios 的 adapter:
请求不会真正发出, 但完整走一遍 请求拦截 -> adapter -> 响应拦截 的链路。
*/

// 真实的 router 会在导入时创建 history, 且这里要断言跳转行为
vi.mock('@/router', () => ({
  router: {
    replace: vi.fn(),
    push: vi.fn(),
    currentRoute: { query: {} },
    getRoutes: () => [],
    hasRoute: () => false,
    removeRoute: vi.fn(),
  },
  resetRouter: vi.fn(),
}))
vi.mock('@/api', () => ({ default: {} }))

// adapter 收到的 config, 用来断言请求头
let lastConfig
// adapter 返回的响应体, 每个用例按需设置
let responseBody

const originalAdapter = request.defaults.adapter

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

// 模拟后端返回非 200 的 HTTP 状态码
function errorAdapter(status, data) {
  return config => Promise.reject(Object.assign(new Error('Request failed'), {
    config,
    response: { status, data, config, headers: {} },
  }))
}

describe('http 拦截器', () => {
  beforeEach(() => {
    setActivePinia(createPinia())
    lastConfig = undefined
    responseBody = { code: 0, message: 'OK', data: null }
    request.defaults.adapter = fakeAdapter
    // 源码里直接用了全局的 $message (由 App.vue 注入)
    window.$message = { error: vi.fn(), success: vi.fn() }
  })

  afterEach(() => {
    request.defaults.adapter = originalAdapter
  })

  it('业务码为 0 时把后端响应体直接抛给调用方', async () => {
    responseBody = { code: 0, message: 'OK', data: { id: 1 } }

    const resp = await request.get('/user/info')

    // 注意: 拿到的是后端响应体本身, 不是 axios 的 response
    expect(resp).toEqual(responseBody)
    expect(window.$message.error).not.toHaveBeenCalled()
  })

  it('业务码非 0 时 reject 并弹错误提示', async () => {
    responseBody = { code: 9004, message: '数据库操作异常', data: null }

    await expect(request.get('/user/info')).rejects.toEqual(responseBody)
    expect(window.$message.error).toHaveBeenCalledWith('数据库操作异常')
  })

  it('message 与 data 不同时一起展示', async () => {
    responseBody = { code: 400, message: '请求参数错误', data: 'id 不能为空' }

    await expect(request.get('/user/info')).rejects.toEqual(responseBody)
    expect(window.$message.error).toHaveBeenCalledWith('请求参数错误 id 不能为空')
  })

  it('业务码 1201 (token 有问题) 会跳转登录页', async () => {
    const { router } = await import('@/router')
    responseBody = { code: 1201, message: 'token 不存在', data: null }

    // 跳登录页的同时必须 reject: resolve(undefined) 会让调用方的
    // const { data } = await ... / resp.data.token 二次抛错
    await expect(request.get('/user/info')).rejects.toEqual(responseBody)
    expect(router.replace).toHaveBeenCalledWith({ path: '/login', query: {} })
  })

  it('业务码 1202 (token 过期) 会清掉登录态', async () => {
    const authStore = useAuthStore()
    authStore.setToken('expired-token')
    responseBody = { code: 1202, message: 'token 已过期', data: null }

    await expect(request.get('/user/info')).rejects.toEqual(responseBody)
    expect(authStore.token).toBeNull()
  })

  it('默认请求会带上 Authorization', async () => {
    useAuthStore().setToken('my-token')

    await request.get('/user/info')

    expect(lastConfig.headers.Authorization).toBe('Bearer my-token')
  })

  it('noNeedToken 的请求不带 Authorization', async () => {
    useAuthStore().setToken('my-token')

    await request.get('/config', { noNeedToken: true })

    expect(lastConfig.headers.Authorization).toBeUndefined()
  })

  it('http 500 且带业务信息时弹出提示', async () => {
    request.defaults.adapter = errorAdapter(500, { message: '服务端异常', data: 'sql error' })

    await expect(request.get('/user/info')).rejects.toThrow('Request failed')
    expect(window.$message.error).toHaveBeenCalledWith('服务端异常 sql error')
  })

  it('http 500 无业务信息时弹出兜底提示', async () => {
    request.defaults.adapter = errorAdapter(500, undefined)

    await expect(request.get('/user/info')).rejects.toThrow('Request failed')
    expect(window.$message.error).toHaveBeenCalledWith('服务端异常')
  })

  it('网络异常 (没有 response) 时原样抛出, 不在拦截器里再抛错', async () => {
    request.defaults.adapter = () => Promise.reject(new Error('Network Error'))

    await expect(request.get('/user/info')).rejects.toThrow('Network Error')
    expect(window.$message.error).not.toHaveBeenCalled()
  })
})
