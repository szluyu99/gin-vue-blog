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
    api.getArticles.mockReset().mockResolvedValue({ code: 0, data: { page_data: articles, total: 2 } })
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
    expect(wrapper.vm.total).toBe(0)
  })

  // 以前接口把 total 丢掉, 前台只取第一页, 分类下超过一页的文章没有入口能看到
  it('按分页参数请求, 并按总数算出页数', async () => {
    api.getArticles.mockResolvedValue({ code: 0, data: { page_data: articles, total: 25 } })
    const wrapper = mountPage()
    await vi.waitFor(() => expect(wrapper.vm.total).toBe(25))

    expect(api.getArticles).toHaveBeenCalledWith(expect.objectContaining({ page_num: 1, page_size: 9 }))
    expect(wrapper.vm.pageCount).toBe(3)
  })

  it('切页会带上新的页码重新请求', async () => {
    api.getArticles.mockResolvedValue({ code: 0, data: { page_data: articles, total: 25 } })
    const wrapper = mountPage()
    await vi.waitFor(() => expect(wrapper.vm.total).toBe(25))
    window.scrollTo = vi.fn()

    wrapper.vm.current = 2
    await vi.waitFor(() => expect(api.getArticles).toHaveBeenCalledWith(expect.objectContaining({ page_num: 2 })))
    expect(window.scrollTo).toHaveBeenCalled()
  })

  // 之前这里写的是 <NPagination>, 但前台没装 naive-ui, 组件解析不到,
  // 分页器一个像素都渲染不出来。上面那条用例直接改 vm.current, 所以没抓到
  it('多于一页时真的把分页器渲染出来', async () => {
    api.getArticles.mockResolvedValue({ code: 0, data: { page_data: articles, total: 25 } })
    const wrapper = mountPage()
    await vi.waitFor(() => expect(wrapper.vm.total).toBe(25))
    await wrapper.vm.$nextTick()

    const pager = wrapper.find('nav[aria-label="分页"]')
    expect(pager.exists()).toBe(true)
    expect(pager.find('[aria-current="page"]').text()).toBe('1')

    // 点分页器上的「2」要真的翻页, 而不是只有 vm.current 能改
    await pager.findAll('button').find(b => b.text() === '2').trigger('click')
    await vi.waitFor(() => expect(api.getArticles).toHaveBeenCalledWith(expect.objectContaining({ page_num: 2 })))
  })

  it('卡片错峰入场, 翻页后不再错峰', async () => {
    api.getArticles.mockResolvedValue({ code: 0, data: { page_data: articles, total: 25 } })
    const wrapper = mountPage()
    await vi.waitFor(() => expect(wrapper.vm.total).toBe(25))
    await wrapper.vm.$nextTick()

    const delays = () => wrapper.findAll('.card-enter').map(c => c.attributes('style'))
    expect(delays()).toEqual(['--i: 0;', '--i: 1;'])

    wrapper.vm.current = 2
    await vi.waitFor(() => expect(wrapper.vm.staggered).toBe(false))
    await wrapper.vm.$nextTick()

    expect(delays()).toEqual(['--i: 0;', '--i: 0;'])
  })

  it('只有一页时不渲染分页器', async () => {
    const wrapper = mountPage()
    await vi.waitFor(() => expect(wrapper.vm.total).toBe(2))

    expect(wrapper.vm.pageCount).toBe(1)
    expect(wrapper.find('nav[aria-label="分页"]').exists()).toBe(false)
    expect(wrapper.text()).not.toContain('这里还没有文章')
  })

  it('没有文章时给出空列表提示', async () => {
    api.getArticles.mockResolvedValue({ code: 0, data: { page_data: [], total: 0 } })
    const wrapper = mountPage()
    await vi.waitFor(() => expect(wrapper.vm.loading).toBe(false))
    await wrapper.vm.$nextTick()

    expect(wrapper.text()).toContain('这里还没有文章')
  })
})
