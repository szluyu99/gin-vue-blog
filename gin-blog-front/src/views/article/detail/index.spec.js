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
