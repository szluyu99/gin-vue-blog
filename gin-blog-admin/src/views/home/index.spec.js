import { mount } from '@vue/test-utils'
import { createPinia, setActivePinia } from 'pinia'
import { beforeEach, describe, expect, it, vi } from 'vitest'

import api from '@/api'
import { useUserStore } from '@/store'
import HomePage from './index.vue'

vi.mock('@/api', () => ({
  default: {
    getHomeInfo: vi.fn(),
    getUserInfo: vi.fn().mockResolvedValue({ code: 0, data: {} }),
    // 仪表盘各块用的都是现成的列表接口, 没有新增后端接口
    getComments: vi.fn(),
    getMessages: vi.fn(),
    getArticles: vi.fn(),
    getCategorys: vi.fn(),
    getLoginLogs: vi.fn(),
    getOnlineUsers: vi.fn(),
  },
}))

vi.mock('vue-router', async importOriginal => ({
  ...await importOriginal(),
  useRouter: () => ({ push: vi.fn() }),
}))

const page = (total, list = []) => ({ code: 0, data: { total, page_data: list } })

const homeInfo = { view_count: 10, user_count: 2, article_count: 4, message_count: 6 }

function mountPage() {
  return mount(HomePage, {
    global: { stubs: { AppPage: { template: '<div><slot /></div>' } } },
  })
}

describe('后台首页', () => {
  beforeEach(() => {
    setActivePinia(createPinia())
    api.getHomeInfo.mockReset().mockResolvedValue({ code: 0, data: { ...homeInfo } })
    api.getComments.mockReset().mockResolvedValue(page(2))
    api.getMessages.mockReset().mockResolvedValue(page(1))
    api.getArticles.mockReset().mockResolvedValue(page(0))
    api.getCategorys.mockReset().mockResolvedValue(page(0))
    api.getLoginLogs.mockReset().mockResolvedValue(page(0))
    api.getOnlineUsers.mockReset().mockResolvedValue({ code: 0, data: [] })
    window.$message = { success: vi.fn(), error: vi.fn(), info: vi.fn() }
  })

  it('挂载后填充统计数据', async () => {
    const wrapper = mountPage()
    await vi.waitFor(() => expect(wrapper.vm.homeInfo.view_count).toBe(10))

    expect(wrapper.vm.homeInfo.article_count).toBe(4)
  })

  it('接口失败或返回空数据时统计不会变成 null', async () => {
    api.getHomeInfo.mockRejectedValue(new Error('boom'))
    const wrapper = mountPage()
    await vi.waitFor(() => expect(api.getHomeInfo).toHaveBeenCalled())
    await wrapper.vm.$nextTick()
    expect(wrapper.vm.homeInfo).toMatchObject({ view_count: 0 })

    api.getHomeInfo.mockResolvedValue({ code: 0, data: null })
    const wrapper2 = mountPage()
    await wrapper2.vm.$nextTick()
    expect(wrapper2.vm.homeInfo).toMatchObject({ view_count: 0 })
  })

  // 以前是 const { nickname, avatar } = useUserStore(), 解构后失去响应性,
  // getUserInfo() 回来得比首屏渲染晚, 问候语里的昵称一直是空的
  it('用户信息后到也能显示在问候语里', async () => {
    const wrapper = mountPage()
    await vi.waitFor(() => expect(wrapper.vm.homeInfo.view_count).toBe(10))
    expect(wrapper.text()).not.toContain('管理员')

    const store = useUserStore()
    store.userInfo = { ...store.userInfo, nickname: '管理员' }
    await wrapper.vm.$nextTick()

    expect(wrapper.text()).toContain('管理员')
  })

  // 仪表盘各块: 待处理数字取各列表接口的 total, 在线用户接口返回的是数组
  it('待处理数字来自各列表接口的 total', async () => {
    api.getComments.mockResolvedValue(page(3))
    api.getMessages.mockResolvedValue(page(2))
    // 回收站那次调用是 is_delete: true, 最新文章那次是 false
    api.getArticles.mockImplementation(params => Promise.resolve(
      params.is_delete ? page(4) : page(1, [{ id: 9, title: '最新一篇', status: 1, created_at: '2026-09-05T00:00:00Z' }]),
    ))
    api.getOnlineUsers.mockResolvedValue({ code: 0, data: [{ id: 1 }, { id: 2 }] })

    const wrapper = mountPage()
    await vi.waitFor(() => expect(wrapper.vm.loading).toBe(false))

    expect(wrapper.vm.todo).toEqual({ comment: 3, message: 2, recycle: 4, online: 2 })
    expect(wrapper.vm.latestArticles).toHaveLength(1)
    expect(wrapper.text()).toContain('最新一篇')
  })

  // 某个接口挂掉不该让整个面板空着, 也不该抛出未捕获的 rejection
  it('单个接口失败时其他块仍然填充', async () => {
    api.getComments.mockRejectedValue(new Error('boom'))
    api.getMessages.mockResolvedValue(page(7))

    const wrapper = mountPage()
    await vi.waitFor(() => expect(wrapper.vm.loading).toBe(false))

    expect(wrapper.vm.todo.comment).toBe(0)
    expect(wrapper.vm.todo.message).toBe(7)
  })

  // 分类分布按文章数倒序取前 5, 条形按最大值归一化
  it('分类分布按文章数倒序取前五', async () => {
    api.getCategorys.mockResolvedValue(page(6, [
      { id: 1, name: 'A', article_count: 2 },
      { id: 2, name: 'B', article_count: 9 },
      { id: 3, name: 'C', article_count: 5 },
      { id: 4, name: 'D', article_count: 1 },
      { id: 5, name: 'E', article_count: 7 },
      { id: 6, name: 'F', article_count: 3 },
    ]))

    const wrapper = mountPage()
    await vi.waitFor(() => expect(wrapper.vm.loading).toBe(false))

    expect(wrapper.vm.categoryStat.map(c => c.name)).toEqual(['B', 'E', 'C', 'F', 'A'])
    expect(wrapper.vm.maxCategoryCount).toBe(9)
  })

  it('没有数据时各块给出空状态', async () => {
    const wrapper = mountPage()
    await vi.waitFor(() => expect(wrapper.vm.loading).toBe(false))

    expect(wrapper.text()).toContain('还没有文章')
    expect(wrapper.text()).toContain('还没有分类')
    expect(wrapper.text()).toContain('还没有登录记录')
  })

  it('一言接口正常时用接口返回的句子', async () => {
    globalThis.fetch = vi.fn().mockResolvedValue({ json: () => Promise.resolve({ hitokoto: '接口给的一句话' }) })

    const wrapper = mountPage()
    await vi.waitFor(() => expect(wrapper.vm.sentence).toBe('接口给的一句话'))
  })

  // 这个接口在部分网络下不通, 兜底要从内置文案里随机取, 而不是固定一句
  it('一言接口挂了从内置文案里随机取', async () => {
    globalThis.fetch = vi.fn().mockRejectedValue(new Error('offline'))

    const wrapper = mountPage()
    await vi.waitFor(() => expect(wrapper.vm.sentence).not.toBe(''))

    expect(wrapper.vm.FALLBACK_SENTENCES).toContain(wrapper.vm.sentence)
  })

  it('一言接口返回空内容时也走兜底', async () => {
    globalThis.fetch = vi.fn().mockResolvedValue({ json: () => Promise.resolve({}) })

    const wrapper = mountPage()
    await vi.waitFor(() => expect(wrapper.vm.sentence).not.toBe(''))

    expect(wrapper.vm.FALLBACK_SENTENCES).toContain(wrapper.vm.sentence)
  })
})
