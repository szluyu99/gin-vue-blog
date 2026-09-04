import { mount } from '@vue/test-utils'
import { createPinia, setActivePinia } from 'pinia'
import { beforeEach, describe, expect, it, vi } from 'vitest'

import api from '@/api'
import LinkPage from './index.vue'

const { copy } = vi.hoisted(() => ({ copy: vi.fn() }))
vi.mock('@vueuse/core', async importOriginal => ({
  ...await importOriginal(),
  useClipboard: () => ({ copy }),
}))

vi.mock('@/api', () => ({
  default: {
    getLinks: vi.fn().mockResolvedValue({ code: 0, data: { page_data: [], total: 0 } }),
    saveOrUpdateLink: vi.fn().mockResolvedValue({ code: 0 }),
    deleteLinks: vi.fn().mockResolvedValue({ code: 0 }),
  },
}))

function mountPage() {
  return mount(LinkPage, {
    global: { stubs: { CommonPage: { template: '<div><slot name="action" /><slot /></div>' } } },
  })
}

const link = { id: 1, name: 'Gin', address: 'https://gin-gonic.com', intro: 'Go 框架', avatar: 'public/uploaded/a.png' }

describe('友链管理', () => {
  beforeEach(() => {
    setActivePinia(createPinia())
    copy.mockClear()
    api.getLinks.mockReset().mockResolvedValue({ code: 0, data: { page_data: [], total: 0 } })
    api.saveOrUpdateLink.mockReset().mockResolvedValue({ code: 0 })
    window.$message = { success: vi.fn(), error: vi.fn(), info: vi.fn() }
  })

  // 本地上传的友链头像存的是相对路径, 不走 convertImgUrl 会裂图
  it('头像列把相对路径转成根相对路径', () => {
    const wrapper = mountPage()
    const column = wrapper.vm.columns.find(e => e.key === 'avatar')

    expect(column.render(link).props.src).toBe('/public/uploaded/a.png')
    // 外链原样返回
    expect(column.render({ ...link, avatar: 'https://cdn.test/a.png' }).props.src).toBe('https://cdn.test/a.png')
  })

  it('点击链接地址复制到剪贴板', () => {
    const wrapper = mountPage()
    const column = wrapper.vm.columns.find(e => e.key === 'address')

    column.render(link).props.onClick()

    expect(copy).toHaveBeenCalledWith('https://gin-gonic.com')
    expect(window.$message.info).toHaveBeenCalled()
  })

  it('删除确认文案说的是友链', () => {
    const wrapper = mountPage()
    const actions = wrapper.vm.columns.find(e => e.key === 'actions')
    const popconfirm = actions.render(link)[1]

    expect(popconfirm.children.default().children).toBe('确定删除该友链吗?')
  })

  it('编辑带上原有行数据', async () => {
    const wrapper = mountPage()
    const actions = wrapper.vm.columns.find(e => e.key === 'actions')

    actions.render(link)[0].props.onClick()
    await wrapper.vm.$nextTick()

    expect(wrapper.vm.modalForm).toMatchObject({ id: 1, name: 'Gin', address: 'https://gin-gonic.com' })
  })
})
