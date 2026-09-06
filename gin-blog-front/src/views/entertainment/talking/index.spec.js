import { mount } from '@vue/test-utils'
import { createPinia, setActivePinia } from 'pinia'
import { beforeEach, describe, expect, it, vi } from 'vitest'

import api from '@/api'
import TalkPage from './index.vue'

vi.mock('@/api', () => ({
  default: {
    getTalks: vi.fn(),
  },
}))

const scrollTo = vi.fn()
vi.stubGlobal('scrollTo', scrollTo)

const talks = [
  { id: 1, content: '普通说说', nickname: '博主', avatar: 'a.png', is_top: false, comment_count: 0, created_at: '2026-09-05T00:00:00Z' },
  { id: 2, content: '置顶说说', nickname: '博主', avatar: 'a.png', is_top: true, comment_count: 3, created_at: '2026-09-06T00:00:00Z' },
]

function mountPage() {
  return mount(TalkPage, {
    global: {
      stubs: {
        BannerPage: { template: '<div><slot /></div>' },
        RouterLink: { props: ['to'], template: '<a :href="to"><slot /></a>' },
      },
    },
  })
}

describe('说说列表页', () => {
  beforeEach(() => {
    setActivePinia(createPinia())
    api.getTalks.mockReset().mockResolvedValue({
      code: 0,
      data: { page_data: talks.map(e => ({ ...e })), total: 2 },
    })
  })

  it('加载并渲染说说, 带评论数和置顶标记', async () => {
    const wrapper = mountPage()
    await vi.waitFor(() => expect(wrapper.vm.loading).toBe(false))

    expect(api.getTalks).toHaveBeenCalledWith({ page_num: 1, page_size: 10 })
    expect(wrapper.text()).toContain('共 2 条说说')
    expect(wrapper.text()).toContain('普通说说')
    expect(wrapper.text()).toContain('置顶说说')
    expect(wrapper.text()).toContain('置顶')
    expect(wrapper.text()).toContain('3 条评论')

    const links = wrapper.findAll('a')
    expect(links[0].attributes('href')).toBe('/talk/1')
    expect(links[1].attributes('href')).toBe('/talk/2')
  })

  it('接口失败时不清空 loading, 也不报未捕获异常', async () => {
    api.getTalks.mockRejectedValue(new Error('boom'))
    const wrapper = mountPage()
    await vi.waitFor(() => expect(wrapper.vm.loading).toBe(false))

    expect(wrapper.text()).toContain('还没有说说')
  })
})
