import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest'

// 每个用例都要一份新模块: 去重用的 Set 是模块级状态
async function load() {
  vi.resetModules()
  return import('./error-report')
}

// setupErrorReport 返回的清理函数: 监听器留在 window 上会让后续用例重复上报
let teardown = null

describe('前端错误上报', () => {
  beforeEach(() => {
    vi.stubGlobal('fetch', vi.fn().mockResolvedValue({ ok: true }))
  })

  afterEach(() => {
    teardown?.()
    teardown = null
  })

  it('上报时带上来源/信息/栈/页面地址', async () => {
    const { reportError } = await load()

    await reportError('front', 'boom', 'at foo (a.js:1:1)')

    expect(fetch).toHaveBeenCalledTimes(1)
    const [url, init] = fetch.mock.calls[0]
    expect(url).toContain('/front/error/report')
    expect(init.method).toBe('POST')
    expect(JSON.parse(init.body)).toMatchObject({
      source: 'front',
      message: 'boom',
      stack: 'at foo (a.js:1:1)',
    })
    expect(JSON.parse(init.body).url).toBeTruthy()
  })

  // 上报接口有按 IP 配额, 一个死循环里的报错会把配额吃光
  it('同一条错误在本次会话里只报一次', async () => {
    const { reportError } = await load()

    await reportError('front', 'boom', 'at foo (a.js:1:1)')
    await reportError('front', 'boom', 'at foo (a.js:1:1)')
    await reportError('front', 'boom', 'at foo (a.js:1:1)')

    expect(fetch).toHaveBeenCalledTimes(1)
  })

  it('栈顶不同视为不同错误, 分别上报', async () => {
    const { reportError } = await load()

    await reportError('front', 'boom', 'at foo (a.js:1:1)')
    await reportError('front', 'boom', 'at bar (b.js:2:2)')

    expect(fetch).toHaveBeenCalledTimes(2)
  })

  it('没有 message 时不上报', async () => {
    const { reportError } = await load()

    await reportError('front', '')

    expect(fetch).not.toHaveBeenCalled()
  })

  // 上报失败再抛错就成环了
  it('上报失败时静默, 不抛出异常', async () => {
    vi.stubGlobal('fetch', vi.fn().mockRejectedValue(new Error('network down')))
    const { reportError } = await load()

    await expect(reportError('front', 'boom')).resolves.toBeUndefined()
  })

  it('捕获未处理的 Promise 拒绝', async () => {
    const { setupErrorReport } = await load()
    teardown = setupErrorReport('front')

    // jsdom 不会自动派发 unhandledrejection, 手动触发
    const event = new Event('unhandledrejection')
    event.reason = new Error('promise 炸了')
    window.dispatchEvent(event)

    expect(fetch).toHaveBeenCalledTimes(1)
    expect(JSON.parse(fetch.mock.calls[0][1].body).message).toBe('promise 炸了')
  })

  it('捕获同步异常', async () => {
    const { setupErrorReport } = await load()
    teardown = setupErrorReport('front')

    const event = new Event('error')
    event.message = 'sync 炸了'
    event.error = new Error('sync 炸了')
    window.dispatchEvent(event)

    expect(fetch).toHaveBeenCalledTimes(1)
    expect(JSON.parse(fetch.mock.calls[0][1].body).message).toBe('sync 炸了')
  })

  // 图片/脚本加载失败也走 error 事件, 但没有 error 对象, 不该占用配额
  it('资源加载失败不上报', async () => {
    const { setupErrorReport } = await load()
    teardown = setupErrorReport('front')

    window.dispatchEvent(new Event('error'))

    expect(fetch).not.toHaveBeenCalled()
  })
})
