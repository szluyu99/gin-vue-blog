import { mount } from '@vue/test-utils'
import { beforeEach, describe, expect, it, vi } from 'vitest'

import api from '@/api'
import TagPage from './index.vue'

vi.mock('@/api', () => ({
  default: { getTags: vi.fn() },
}))

const tags = [
  { id: 1, name: 'Go' },
  { id: 2, name: 'Vue' },
]

function mountPage() {
  return mount(TagPage, {
    global: { stubs: { BannerPage: { template: '<div><slot /></div>' }, RouterLink: { template: '<a><slot /></a>' } } },
  })
}

describe('标签列表', () => {
  beforeEach(() => {
    api.getTags.mockReset().mockResolvedValue({ code: 0, data: tags })
  })

  it('渲染标签与数量', async () => {
    const wrapper = mountPage()
    await vi.waitFor(() => expect(wrapper.vm.tagList).toHaveLength(2))

    expect(wrapper.text()).toContain('标签 - 2')
    expect(wrapper.text()).toContain('Go')
  })

  // 回归: 以前 loading 只在 then 里复位, 请求失败页面永远停在加载态
  it('接口失败时 loading 复位', async () => {
    api.getTags.mockRejectedValue(new Error('boom'))
    const wrapper = mountPage()
    await vi.waitFor(() => expect(wrapper.vm.loading).toBe(false))

    expect(wrapper.vm.tagList).toEqual([])
  })

  it('随机字号与颜色都在合理范围内', () => {
    const wrapper = mountPage()

    for (let i = 0; i < 20; i++) {
      const size = wrapper.vm.randomFontSize()
      expect(size).toBeGreaterThanOrEqual(15)
      expect(size).toBeLessThanOrEqual(30)
      expect(wrapper.vm.randomColorHex()).toMatch(/^#[0-9a-f]{1,6}$/)
    }
  })
})
