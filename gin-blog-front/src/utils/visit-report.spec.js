import { beforeEach, describe, expect, it, vi } from 'vitest'

vi.mock('@/api', () => ({
  default: { report: vi.fn(() => Promise.resolve({ code: 0 })) },
}))

// 去重标记在 sessionStorage 里, 每个用例都要一份干净的
async function load() {
  vi.resetModules()
  sessionStorage.clear()
  const api = (await import('@/api')).default
  api.report.mockClear()
  return { reportVisit: (await import('./visit-report')).reportVisit, api }
}

describe('访客上报', () => {
  beforeEach(() => {
    sessionStorage.clear()
  })

  it('首次进入会上报一次', async () => {
    const { reportVisit, api } = await load()

    await reportVisit()

    expect(api.report).toHaveBeenCalledTimes(1)
  })

  // 同一次会话里刷新页面不该重复计数
  it('同一会话内只上报一次', async () => {
    const { reportVisit, api } = await load()

    await reportVisit()
    await reportVisit()
    await reportVisit()

    expect(api.report).toHaveBeenCalledTimes(1)
  })

  // 统计失败不能给访客任何可见的动静
  it('上报失败时静默', async () => {
    const { reportVisit, api } = await load()
    api.report.mockRejectedValueOnce(new Error('network down'))

    await expect(reportVisit()).resolves.toBeUndefined()
  })

  // 隐私模式下 sessionStorage 会抛错, 不能因此卡住上报
  it('sessionStorage 不可用时仍然上报', async () => {
    const { reportVisit, api } = await load()
    const spy = vi.spyOn(Storage.prototype, 'getItem').mockImplementation(() => {
      throw new Error('denied')
    })

    await reportVisit()

    expect(api.report).toHaveBeenCalledTimes(1)
    spy.mockRestore()
  })
})
