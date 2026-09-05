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
    global: { stubs: { BannerPage: { template: '<div><slot /></div>' }, RouterLink: { template: '<a><slot /></a>' } } },
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

    wrapper.vm.current = 2
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
  })
})
