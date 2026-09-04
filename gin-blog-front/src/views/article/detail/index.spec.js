import { mount } from '@vue/test-utils'
import { createPinia, setActivePinia } from 'pinia'
import { beforeEach, describe, expect, it, vi } from 'vitest'

import api from '@/api'
import ArticleDetail from './index.vue'

vi.mock('@/api', () => ({
  default: {
    getArticleDetail: vi.fn(),
    getComments: vi.fn().mockResolvedValue({ code: 0, data: { page_data: [], total: 0 } }),
  },
}))

vi.mock('vue-router', async importOriginal => ({
  ...await importOriginal(),
  useRoute: () => ({ params: { id: '1' } }),
}))

// MathJax 只在正文含公式时加载, 测试里不需要真的去加载
vi.mock('@/utils/mathjax', () => ({
  typesetMath: vi.fn().mockResolvedValue(undefined),
  hasMath: () => false,
}))

const article = {
  id: 1,
  title: '第一篇',
  content: '# 标题\n\n```go\nfmt.Println("hi")\n```',
  img: 'public/uploaded/a.png',
  like_count: 3,
  view_count: 10,
  comment_count: 0,
  created_at: '2026-09-01T00:00:00Z',
  updated_at: '2026-09-01T00:00:00Z',
  tags: [{ id: 1, name: 'Go' }],
  category: { id: 1, name: '项目' },
  newest_articles: [],
  recommend_articles: [],
  last_article: {},
  next_article: {},
}

function mountPage() {
  return mount(ArticleDetail, {
    global: {
      stubs: {
        BannerInfo: true,
        Catalogue: true,
        Copyright: true,
        Forward: true,
        LastNext: true,
        LatestList: true,
        Recommend: true,
        Reward: true,
        Comment: true,
        AppFooter: true,
        RouterLink: { template: '<a><slot /></a>' },
      },
    },
  })
}

describe('文章详情', () => {
  beforeEach(() => {
    setActivePinia(createPinia())
    api.getArticleDetail.mockReset().mockResolvedValue({ code: 0, data: { ...article } })
  })

  it('markdown 正文被解析成 HTML', async () => {
    const wrapper = mountPage()
    await vi.waitFor(() => expect(wrapper.vm.loading).toBe(false))

    expect(wrapper.vm.data.content).toContain('<h1')
    expect(wrapper.vm.data.content).toContain('<code')
  })

  it('封面走 convertImgUrl, 没有封面时用灰底', async () => {
    const wrapper = mountPage()
    await vi.waitFor(() => expect(wrapper.vm.loading).toBe(false))

    expect(wrapper.vm.styleVal).toContain('/public/uploaded/a.png')

    api.getArticleDetail.mockResolvedValue({ code: 0, data: { ...article, img: '' } })
    const wrapper2 = mountPage()
    await vi.waitFor(() => expect(wrapper2.vm.loading).toBe(false))
    expect(wrapper2.vm.styleVal).toContain('rgba(0,0,0,0.1)')
  })

  // 后端这些字段可能是 null: 未分类、没有上一篇、没有推荐文章
  it('可空字段为 null 时退化成默认值, 页面不崩', async () => {
    api.getArticleDetail.mockResolvedValue({
      code: 0,
      data: {
        ...article,
        content: null,
        tags: null,
        category: null,
        last_article: null,
        next_article: null,
        newest_articles: null,
        recommend_articles: null,
      },
    })

    const wrapper = mountPage()
    await vi.waitFor(() => expect(wrapper.vm.loading).toBe(false))

    expect(wrapper.vm.data.tags).toEqual([])
    expect(wrapper.vm.data.recommend_articles).toEqual([])
    expect(wrapper.vm.data.category).toEqual({})
    expect(wrapper.vm.data.content).toBe('')
  })

  it('给代码块加上复制按钮', async () => {
    Object.defineProperty(navigator, 'clipboard', {
      value: { writeText: vi.fn().mockResolvedValue(undefined) },
      configurable: true,
      writable: true,
    })

    const wrapper = mountPage()
    await vi.waitFor(() => expect(wrapper.vm.loading).toBe(false))
    await wrapper.vm.$nextTick()

    expect(wrapper.find('[data-copy-btn]').exists()).toBe(true)
  })

  // 曾经的 bug: readProgress 在读 y 之前就因为"页面还撑不出滚动条"提前 return 0,
  // y 没被记成依赖, 之后滚动再也不会重算, 进度条永远停在 0 宽
  it('滚动后阅读进度条跟着变宽', async () => {
    const wrapper = mountPage()
    await vi.waitFor(() => expect(wrapper.vm.loading).toBe(false))

    const bar = wrapper.find('.z-999')
    expect(bar.exists()).toBe(true)
    expect(bar.attributes('style')).toContain('width: 0%')

    // jsdom 里 scrollHeight 恒为 0, 手动撑出一个可滚动的页面。
    // 注意不要动 innerHeight: 那是 readProgress 的另一个依赖, 改它会顺带
    // 触发重算, 就算漏了 y 也会算对, 这个用例就抓不到 bug 了
    const half = Math.round(window.innerHeight / 2)
    Object.defineProperty(document.documentElement, 'scrollHeight', {
      value: window.innerHeight * 3,
      configurable: true,
    })
    // vueuse 的 useWindowScroll 底层读的是滚动元素的 scrollTop
    document.documentElement.scrollTop = half * 2
    window.dispatchEvent(new Event('scroll'))
    await wrapper.vm.$nextTick()

    expect(wrapper.vm.readProgress).toBeCloseTo(50)
    expect(wrapper.find('.z-999').attributes('style')).toContain('width: 50%')

    // scrollHeight / scrollTop 是全局的, 留着会影响同文件后面的用例
    Reflect.deleteProperty(document.documentElement, 'scrollHeight')
    document.documentElement.scrollTop = 0
  })

  // 技术文章过期得快, 太久没更新要提醒读者
  it('文章太久没更新时给出提示', async () => {
    const old = new Date(Date.now() - 200 * 86400000).toISOString()
    api.getArticleDetail.mockResolvedValue({ code: 0, data: { ...article, updated_at: old } })

    const wrapper = mountPage()
    await vi.waitFor(() => expect(wrapper.vm.loading).toBe(false))
    await wrapper.vm.$nextTick()

    expect(wrapper.vm.staleDays).toBe(200)
    expect(wrapper.text()).toContain('本文最后更新于 200 天前')
  })

  it('最近更新过的文章不提示', async () => {
    const recent = new Date(Date.now() - 3 * 86400000).toISOString()
    api.getArticleDetail.mockResolvedValue({ code: 0, data: { ...article, updated_at: recent } })

    const wrapper = mountPage()
    await vi.waitFor(() => expect(wrapper.vm.loading).toBe(false))

    expect(wrapper.vm.staleDays).toBe(0)
    expect(wrapper.text()).not.toContain('部分内容可能已经过时')
  })

  it('阅读进度按滚动比例计算, 不会越界', async () => {
    const wrapper = mountPage()
    await vi.waitFor(() => expect(wrapper.vm.loading).toBe(false))

    // jsdom 里 scrollHeight 是 0, 算不出比例时退回 0 而不是 NaN
    expect(wrapper.vm.readProgress).toBe(0)
  })

  it('接口失败时 loading 复位, 数据保持默认结构', async () => {
    api.getArticleDetail.mockRejectedValue(new Error('boom'))
    const wrapper = mountPage()
    await vi.waitFor(() => expect(wrapper.vm.loading).toBe(false))

    expect(wrapper.vm.data.id).toBe(0)
    expect(wrapper.vm.data.tags).toEqual([])
  })

  it('后端返回空数据时不覆盖默认结构', async () => {
    api.getArticleDetail.mockResolvedValue({ code: 0, data: null })
    const wrapper = mountPage()
    await vi.waitFor(() => expect(wrapper.vm.loading).toBe(false))

    expect(wrapper.vm.data.title).toBe('')
    expect(wrapper.vm.data.newest_articles).toEqual([])
  })
})
