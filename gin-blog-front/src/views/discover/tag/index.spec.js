import { mount } from '@vue/test-utils'
import { beforeEach, describe, expect, it, vi } from 'vitest'

import api from '@/api'
import TagPage from './index.vue'

vi.mock('@/api', () => ({
  default: { getTags: vi.fn() },
}))

const tags = [
  { id: 1, name: 'Go', article_count: 10 },
  { id: 2, name: 'Vue', article_count: 2 },
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

  // 曾经是 15~30px 纯随机, 标签云"字越大越热门"的语义完全没了
  it('字号按文章数映射, 最多的最大, 最少的最小', async () => {
    api.getTags.mockResolvedValue({ code: 0, data: [...tags, { id: 3, name: 'Gin', article_count: 6 }] })
    const wrapper = mountPage()
    await vi.waitFor(() => expect(wrapper.vm.tagList).toHaveLength(3))

    const [go, vue, gin] = wrapper.vm.tagList
    expect(wrapper.vm.fontSize(go)).toBe(30) // 10 篇, 最多
    expect(wrapper.vm.fontSize(vue)).toBe(15) // 2 篇, 最少
    // 6 篇落在中间: 15 + 15 * (6-2)/(10-2) = 22.5 → 23
    expect(wrapper.vm.fontSize(gin)).toBe(23)
  })

  it('文章数都一样时取中间字号, 而不是全部顶到最大', async () => {
    api.getTags.mockResolvedValue({ code: 0, data: [
      { id: 1, name: 'Go', article_count: 3 },
      { id: 2, name: 'Vue', article_count: 3 },
    ] })
    const wrapper = mountPage()
    await vi.waitFor(() => expect(wrapper.vm.tagList).toHaveLength(2))

    expect(wrapper.vm.fontSize(wrapper.vm.tagList[0])).toBe(23)
  })

  it('缺 article_count 时按 0 处理, 不产生 NaN', async () => {
    api.getTags.mockResolvedValue({ code: 0, data: [{ id: 1, name: 'Go' }, { id: 2, name: 'Vue', article_count: 4 }] })
    const wrapper = mountPage()
    await vi.waitFor(() => expect(wrapper.vm.tagList).toHaveLength(2))

    expect(wrapper.vm.fontSize(wrapper.vm.tagList[0])).toBe(15)
    expect(wrapper.vm.fontSize(wrapper.vm.tagList[1])).toBe(30)
  })

  // 原来是 `#${Math.floor(Math.random() * 16777215).toString(16)}`,
  // 约 6% 概率生成不足 6 位的无效色值, 而且每次进页面颜色都变
  it('颜色取自固定色板, 同一个标签颜色稳定', () => {
    const wrapper = mountPage()

    for (const tag of [{ id: 1 }, { id: 7 }, { id: 8 }, { id: 99 }]) {
      const c = wrapper.vm.color(tag)
      expect(c).toMatch(/^#[0-9a-f]{6}$/)
      expect(wrapper.vm.color(tag)).toBe(c)
    }
    // 按 id 取模, 色板长度 8
    expect(wrapper.vm.color({ id: 8 })).toBe(wrapper.vm.color({ id: 0 }))
  })
})
