import { mount } from '@vue/test-utils'
import { createPinia, setActivePinia } from 'pinia'
import { beforeEach, describe, expect, it, vi } from 'vitest'

import api from '@/api'
import PageSetting from './index.vue'

vi.mock('@/api', () => ({
  default: {
    getPages: vi.fn(),
    saveOrUpdatePage: vi.fn().mockResolvedValue({ code: 0 }),
    deletePage: vi.fn().mockResolvedValue({ code: 0 }),
  },
}))

const pages = [
  { id: 1, name: '首页', label: 'home', cover: 'public/uploaded/home.jpg' },
  { id: 2, name: '归档', label: 'archive', cover: 'https://cdn.test/archive.png' },
]

function mountPage() {
  return mount(PageSetting, {
    global: { stubs: { CommonPage: { template: '<div><slot name="action" /><slot /></div>' } } },
  })
}

describe('页面管理', () => {
  beforeEach(() => {
    setActivePinia(createPinia())
    api.getPages.mockReset().mockResolvedValue({ code: 0, data: pages })
    api.deletePage.mockReset().mockResolvedValue({ code: 0 })
    window.$message = { success: vi.fn(), error: vi.fn(), info: vi.fn() }
    window.$dialog = { confirm: ({ confirm }) => confirm() }
  })

  it('挂载后渲染页面列表', async () => {
    const wrapper = mountPage()
    await vi.waitFor(() => expect(wrapper.vm.pageList).toHaveLength(2))

    expect(wrapper.text()).toContain('首页')
    expect(wrapper.text()).toContain('归档')
  })

  it('拉取失败不抛出, 列表保持数组', async () => {
    api.getPages.mockRejectedValue(new Error('boom'))
    const wrapper = mountPage()
    await vi.waitFor(() => expect(api.getPages).toHaveBeenCalled())
    await wrapper.vm.$nextTick()

    expect(wrapper.vm.pageList).toEqual([])
  })

  it('后端返回空数据时列表不会变成 undefined', async () => {
    api.getPages.mockResolvedValue({ code: 0, data: null })
    const wrapper = mountPage()
    await vi.waitFor(() => expect(api.getPages).toHaveBeenCalled())
    await wrapper.vm.$nextTick()

    expect(wrapper.vm.pageList).toEqual([])
  })

  it('下拉菜单选编辑带上行数据, 选删除调删除接口', async () => {
    const wrapper = mountPage()
    await vi.waitFor(() => expect(wrapper.vm.pageList).toHaveLength(2))

    wrapper.vm.handleSelect('edit', pages[0])
    await wrapper.vm.$nextTick()
    expect(wrapper.vm.modalForm).toMatchObject({ id: 1, label: 'home' })

    wrapper.vm.handleSelect('delete', pages[1])
    await vi.waitFor(() => expect(api.deletePage).toHaveBeenCalled())
    expect(api.deletePage).toHaveBeenCalledWith(JSON.stringify([2]))
  })

  // 弹窗没打开时 UploadOne 还没挂载
  it('弹窗未打开时刷新预览图不会抛异常', async () => {
    const wrapper = mountPage()
    await vi.waitFor(() => expect(wrapper.vm.pageList).toHaveLength(2))

    expect(() => wrapper.vm.refreshImg('public/uploaded/x.png')).not.toThrow()
  })
})
