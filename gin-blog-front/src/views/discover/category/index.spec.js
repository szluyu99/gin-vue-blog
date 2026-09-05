import { mount } from '@vue/test-utils'
import { beforeEach, describe, expect, it, vi } from 'vitest'

import api from '@/api'
import CategoryPage from './index.vue'

vi.mock('@/api', () => ({
  default: { getCategorys: vi.fn() },
}))

const categories = [
  { id: 1, name: '项目', article_count: 3 },
  { id: 2, name: '随笔', article_count: 1 },
]

function mountPage() {
  return mount(CategoryPage, {
    global: { stubs: { BannerPage: { template: '<div><slot /></div>' }, RouterLink: { template: '<a><slot /></a>' } } },
  })
}

describe('分类列表', () => {
  beforeEach(() => {
    api.getCategorys.mockReset().mockResolvedValue({ code: 0, data: categories })
  })

  it('渲染分类与数量', async () => {
    const wrapper = mountPage()
    await vi.waitFor(() => expect(wrapper.vm.categoryList).toHaveLength(2))

    expect(wrapper.text()).toContain('共 2 个分类')
    expect(wrapper.text()).toContain('项目')
    // 每个分类一张卡片, 文章数以 badge 形式渲染出来
    expect(wrapper.findAll('li')).toHaveLength(2)
    expect(wrapper.text()).toContain('3')
  })

  it('接口失败时 loading 复位, 列表保持空数组', async () => {
    api.getCategorys.mockRejectedValue(new Error('boom'))
    const wrapper = mountPage()
    await vi.waitFor(() => expect(wrapper.vm.loading).toBe(false))

    expect(wrapper.vm.categoryList).toEqual([])
  })

  // 模板里直接取 categoryList.length, 传 null 会抛
  it('后端返回空数据时列表不会变成 null', async () => {
    api.getCategorys.mockResolvedValue({ code: 0, data: null })
    const wrapper = mountPage()
    await vi.waitFor(() => expect(wrapper.vm.loading).toBe(false))

    expect(wrapper.vm.categoryList).toEqual([])
    expect(wrapper.text()).toContain('共 0 个分类')
    // 一个分类都没有时给提示, 而不是一张空白卡片
    expect(wrapper.text()).toContain('还没有分类')
  })
})
