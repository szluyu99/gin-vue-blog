import { mount } from '@vue/test-utils'
import { createPinia, setActivePinia } from 'pinia'
import { beforeEach, describe, expect, it, vi } from 'vitest'

import api from '@/api'
import LinkPage from './index.vue'

vi.mock('@/api', () => ({
  default: {
    getLinks: vi.fn(),
  },
}))

const links = [
  { id: 1, name: 'Gin', address: 'https://gin-gonic.com', intro: 'Go 框架', avatar: 'public/uploaded/a.png' },
]

function mountPage() {
  return mount(LinkPage, {
    global: {
      stubs: {
        BannerPage: { template: '<div><slot /></div>' },
        AddLink: true,
        Comment: true,
      },
    },
  })
}

describe('友情链接', () => {
  beforeEach(() => {
    setActivePinia(createPinia())
    api.getLinks.mockReset().mockResolvedValue({ code: 0, data: links })
  })

  it('渲染友链列表', async () => {
    const wrapper = mountPage()
    await vi.waitFor(() => expect(wrapper.vm.linkList).toHaveLength(1))

    expect(wrapper.text()).toContain('Gin')
  })

  it('接口失败时 loading 复位, 列表保持空数组', async () => {
    api.getLinks.mockRejectedValue(new Error('boom'))
    const wrapper = mountPage()
    await vi.waitFor(() => expect(wrapper.vm.loading).toBe(false))

    expect(wrapper.vm.linkList).toEqual([])
  })

  // LinkList 里直接取 linkList.length, 传 null 会抛
  it('后端返回空数据时列表不会变成 null', async () => {
    api.getLinks.mockResolvedValue({ code: 0, data: null })
    const wrapper = mountPage()
    await vi.waitFor(() => expect(wrapper.vm.loading).toBe(false))

    expect(wrapper.vm.linkList).toEqual([])
    expect(wrapper.text()).toContain('暂无友情链接')
  })
})
