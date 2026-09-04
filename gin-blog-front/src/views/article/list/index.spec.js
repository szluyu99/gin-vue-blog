import { mount } from '@vue/test-utils'
import { createPinia, setActivePinia } from 'pinia'
import { beforeEach, describe, expect, it, vi } from 'vitest'

import api from '@/api'
import ArticleListPage from './index.vue'

vi.mock('@/api', () => ({
  default: {
    getArticles: vi.fn(),
  },
}))

vi.mock('vue-router', async importOriginal => ({
  ...await importOriginal(),
  useRoute: () => ({ params: { categoryId: '1' }, query: { name: '项目' }, meta: { title: '文章列表' } }),
}))

const articles = [
  { id: 1, title: '第一篇', img: 'public/uploaded/a.png', category_id: 1, category: { name: '项目' }, tags: [], created_at: '2026-09-01T00:00:00Z' },
  { id: 2, title: '没有分类的文章', img: '', category_id: 0, category: null, tags: null, created_at: '2026-09-02T00:00:00Z' },
]

function mountPage() {
  return mount(ArticleListPage, {
    global: { stubs: { BannerPage: { template: '<div><slot /></div>' }, RouterLink: { template: '<a><slot /></a>' } } },
  })
}

describe('前台文章列表', () => {
  beforeEach(() => {
    setActivePinia(createPinia())
    api.getArticles.mockReset().mockResolvedValue({ code: 0, data: articles })
  })

  it('渲染文章卡片, 图片走 convertImgUrl', async () => {
    const wrapper = mountPage()
    await vi.waitFor(() => expect(wrapper.vm.articleList).toHaveLength(2))

    expect(wrapper.text()).toContain('第一篇')
    const imgs = wrapper.findAll('img')
    expect(imgs[0].attributes('src')).toBe('/public/uploaded/a.png')
  })

  // 后端 category 是指针, 未分类的文章返回 null
  it('分类为 null 的文章不会让页面崩掉', async () => {
    const wrapper = mountPage()
    await vi.waitFor(() => expect(wrapper.vm.articleList).toHaveLength(2))

    expect(wrapper.text()).toContain('没有分类的文章')
  })

  // 回归: 以前 loading 只在成功路径上复位, 接口失败页面永远停在加载态
  it('接口失败时 loading 复位, 列表保持空数组', async () => {
    api.getArticles.mockRejectedValue(new Error('boom'))
    const wrapper = mountPage()
    await vi.waitFor(() => expect(wrapper.vm.loading).toBe(false))

    expect(wrapper.vm.articleList).toEqual([])
  })

  it('后端返回空数据时列表不会变成 null', async () => {
    api.getArticles.mockResolvedValue({ code: 0, data: null })
    const wrapper = mountPage()
    await vi.waitFor(() => expect(wrapper.vm.loading).toBe(false))

    expect(wrapper.vm.articleList).toEqual([])
  })
})
