import { mount } from '@vue/test-utils'
import { createPinia, setActivePinia } from 'pinia'
import { beforeEach, describe, expect, it, vi } from 'vitest'

import api from '@/api'
import HomePage from './index.vue'

vi.mock('@/api', () => ({
  default: {
    getArticles: vi.fn(),
  },
}))

// 无限滚动组件依赖 IntersectionObserver, 测试里用桩
vi.mock('v3-infinite-loading', () => ({
  default: { name: 'InfiniteLoading', template: '<div />' },
}))

const articles = [
  { id: 1, title: '第一篇', content: '# 标题\n正文**加粗**', img: '', created_at: '2026-09-01T00:00:00Z' },
  { id: 2, title: '第二篇', content: '正文', img: '', created_at: '2026-09-02T00:00:00Z' },
]

function mountPage() {
  return mount(HomePage, {
    global: {
      stubs: {
        HomeBanner: true,
        AuthorInfo: true,
        WebsiteInfo: true,
        Announcement: true,
        TalkingCarousel: true,
        ArticleCard: true,
        AppFooter: true,
        RouterLink: { template: '<a><slot /></a>' },
      },
    },
  })
}

describe('前台首页', () => {
  beforeEach(() => {
    setActivePinia(createPinia())
    api.getArticles.mockReset().mockResolvedValue({ code: 0, data: articles })
  })

  it('首屏加载文章并去掉正文里的 Markdown 记号', async () => {
    const wrapper = mountPage()
    await vi.waitFor(() => expect(wrapper.vm.articleList).toHaveLength(2))

    // params 是 reactive 对象且会被 page_num++ 改掉, 断言调用次数而不是当时的入参快照
    expect(api.getArticles).toHaveBeenCalledTimes(1)
    expect(wrapper.vm.articleList[0].content).not.toContain('#')
    expect(wrapper.vm.articleList[0].content).not.toContain('**')
    expect(wrapper.vm.params.page_num).toBe(2)
  })

  // 回归: 无限加载靠 !loading 才会继续请求, 首屏失败后 loading 停在 true
  // 就等于首页彻底不再加载文章
  it('首屏失败时 loading 复位', async () => {
    api.getArticles.mockRejectedValue(new Error('boom'))
    const wrapper = mountPage()
    await vi.waitFor(() => expect(wrapper.vm.loading).toBe(false))

    expect(wrapper.vm.articleList).toEqual([])
  })

  it('后端返回空数据时不抛异常', async () => {
    api.getArticles.mockResolvedValue({ code: 0, data: null })
    const wrapper = mountPage()
    await vi.waitFor(() => expect(wrapper.vm.loading).toBe(false))

    expect(wrapper.vm.articleList).toEqual([])
  })

  it('无限加载: 有数据时追加并翻页', async () => {
    const wrapper = mountPage()
    await vi.waitFor(() => expect(wrapper.vm.articleList).toHaveLength(2))

    const state = { loaded: vi.fn(), complete: vi.fn(), error: vi.fn() }
    api.getArticles.mockResolvedValue({ code: 0, data: [{ id: 3, title: '第三篇', content: '正文' }] })
    await wrapper.vm.getArticlesInfinite(state)

    expect(wrapper.vm.articleList).toHaveLength(3)
    expect(wrapper.vm.params.page_num).toBe(3)
    expect(state.loaded).toHaveBeenCalled()
  })

  it('无限加载: 空数据标记完成, 失败标记错误', async () => {
    const wrapper = mountPage()
    await vi.waitFor(() => expect(wrapper.vm.articleList).toHaveLength(2))

    const state = { loaded: vi.fn(), complete: vi.fn(), error: vi.fn() }
    api.getArticles.mockResolvedValue({ code: 0, data: [] })
    await wrapper.vm.getArticlesInfinite(state)
    expect(state.complete).toHaveBeenCalled()

    api.getArticles.mockRejectedValue(new Error('boom'))
    await wrapper.vm.getArticlesInfinite(state)
    expect(state.error).toHaveBeenCalled()
  })

  it('无限加载: 返回 null 也算加载完成, 不抛异常', async () => {
    const wrapper = mountPage()
    await vi.waitFor(() => expect(wrapper.vm.articleList).toHaveLength(2))

    const state = { loaded: vi.fn(), complete: vi.fn(), error: vi.fn() }
    api.getArticles.mockResolvedValue({ code: 0, data: null })
    await wrapper.vm.getArticlesInfinite(state)

    expect(state.complete).toHaveBeenCalled()
    expect(state.error).not.toHaveBeenCalled()
  })
})
