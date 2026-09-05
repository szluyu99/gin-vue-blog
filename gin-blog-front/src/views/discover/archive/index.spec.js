import { mount } from '@vue/test-utils'
import { createPinia, setActivePinia } from 'pinia'
import { beforeEach, describe, expect, it, vi } from 'vitest'

import api from '@/api'
import ArchivePage from './index.vue'

vi.mock('@/api', () => ({
  default: {
    getArchives: vi.fn(),
  },
}))

vi.mock('vue-router', async importOriginal => ({
  ...await importOriginal(),
  useRouter: () => ({ push: vi.fn() }),
}))

const page1 = {
  page_data: [
    { id: 1, title: '第一篇', created_at: '2026-09-01T00:00:00Z' },
    { id: 2, title: '第二篇', created_at: '2026-08-01T00:00:00Z' },
  ],
  total: 60,
}

function mountPage() {
  return mount(ArchivePage, {
    global: {
      stubs: {
        BannerPage: { template: '<div><slot /></div>' },
        // 桩把 to 渲染成 href, 才能断言标题用的是 RouterLink 而不是 <a @click>
        RouterLink: { props: ['to'], template: '<a :href="to"><slot /></a>' },
      },
    },
  })
}

describe('归档', () => {
  beforeEach(() => {
    setActivePinia(createPinia())
    api.getArchives.mockReset().mockResolvedValue({ code: 0, data: { ...page1 } })
  })

  it('挂载后按第一页拉取并渲染', async () => {
    const wrapper = mountPage()
    await vi.waitFor(() => expect(wrapper.vm.archiveList).toHaveLength(2))

    expect(api.getArchives).toHaveBeenCalledWith({ page_num: 1, page_size: 50 })
    expect(wrapper.vm.total).toBe(60)
    expect(wrapper.text()).toContain('第一篇')
    // 计数文案要用 total: 以前用的是当页条数(page_size=50), 分页后会少报
    expect(wrapper.text()).toContain('目前共计 60 篇文章')
  })

  it('切换页数会重新拉取', async () => {
    const wrapper = mountPage()
    await vi.waitFor(() => expect(wrapper.vm.archiveList).toHaveLength(2))
    window.scrollTo = vi.fn()

    wrapper.vm.current = 2
    await vi.waitFor(() => expect(api.getArchives).toHaveBeenCalledWith({ page_num: 2, page_size: 50 }))
    expect(window.scrollTo).toHaveBeenCalled()
  })

  // 归档的意义就是时间轴, 之前是平铺一长串, 看不出年月节奏
  it('按年月分组, 同月的归到一组', async () => {
    api.getArchives.mockResolvedValue({ code: 0, data: {
      page_data: [
        { id: 1, title: '九月甲', created_at: '2026-09-05T00:00:00Z' },
        { id: 2, title: '九月乙', created_at: '2026-09-01T00:00:00Z' },
        { id: 3, title: '八月甲', created_at: '2026-08-20T00:00:00Z' },
        { id: 4, title: '去年', created_at: '2025-08-20T00:00:00Z' },
      ],
      total: 4,
    } })
    const wrapper = mountPage()
    await vi.waitFor(() => expect(wrapper.vm.archiveList).toHaveLength(4))

    // 接口已按时间倒序返回, 分组保持这个顺序
    expect(wrapper.vm.monthGroups.map(g => g.key)).toEqual(['2026-09', '2026-08', '2025-08'])
    expect(wrapper.vm.monthGroups[0].items).toHaveLength(2)
    expect(wrapper.text()).toContain('2026 年 9 月')
    // 同名月份跨年不能合并
    expect(wrapper.text()).toContain('2025 年 8 月')
  })

  // 标题以前是 <a @click="router.push()">, 没有 href, 中键新开/复制链接都用不了
  it('标题走 RouterLink, 带上 to', async () => {
    const wrapper = mountPage()
    await vi.waitFor(() => expect(wrapper.vm.archiveList).toHaveLength(2))

    const links = wrapper.findAll('a')
    expect(links.length).toBeGreaterThan(0)
    expect(links[0].attributes('href')).toBe('/article/1')
  })

  // 以前分页整段被注释掉, page_size 写死 50, 第 51 篇之后没有入口
  it('总数超过一页时渲染分页器', async () => {
    const wrapper = mountPage()
    await vi.waitFor(() => expect(wrapper.vm.total).toBe(60))
    await wrapper.vm.$nextTick()

    expect(wrapper.vm.pageCount).toBe(2)
    const pager = wrapper.find('nav[aria-label="分页"]')
    expect(pager.exists()).toBe(true)

    window.scrollTo = vi.fn()
    await pager.findAll('button').find(b => b.text() === '2').trigger('click')
    await vi.waitFor(() => expect(api.getArchives).toHaveBeenCalledWith({ page_num: 2, page_size: 50 }))
  })

  // 回归: 以前 loading 只在成功路径复位, 失败就一直停在加载态
  it('接口失败时 loading 复位', async () => {
    api.getArchives.mockRejectedValue(new Error('boom'))
    const wrapper = mountPage()
    await vi.waitFor(() => expect(wrapper.vm.loading).toBe(false))

    expect(wrapper.vm.archiveList).toEqual([])
  })

  // 回归: 以前直接取 resp.data.page_data, data 为空会抛在 onMounted 里
  it('后端返回空数据时不抛异常', async () => {
    api.getArchives.mockResolvedValue({ code: 0, data: null })
    const wrapper = mountPage()
    await vi.waitFor(() => expect(wrapper.vm.loading).toBe(false))

    expect(wrapper.vm.archiveList).toEqual([])
    expect(wrapper.vm.total).toBe(0)
    expect(wrapper.vm.monthGroups).toEqual([])
    expect(wrapper.text()).toContain('还没有文章')
    expect(wrapper.find('nav[aria-label="分页"]').exists()).toBe(false)
  })
})
