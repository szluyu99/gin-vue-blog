import { mount } from '@vue/test-utils'
import { createPinia, setActivePinia } from 'pinia'
import { beforeEach, describe, expect, it, vi } from 'vitest'

import api from '@/api'
import LeaveMsgPage from './index.vue'

vi.mock('@/api', () => ({
  default: {
    getMessages: vi.fn().mockResolvedValue({ code: 0, data: { page_data: [], total: 0 } }),
    deleteMessages: vi.fn().mockResolvedValue({ code: 0 }),
    updateMessageReview: vi.fn().mockResolvedValue({ code: 0 }),
  },
}))

function mountPage() {
  return mount(LeaveMsgPage, {
    global: { stubs: { CommonPage: { template: '<div><slot name="action" /><slot /></div>' } } },
  })
}

const msg = { id: 5, nickname: '路人甲', content: '你好', avatar: 'public/uploaded/a.png', ip_address: '10.0.0.1', ip_source: '内网IP', is_review: false }

describe('留言管理', () => {
  beforeEach(() => {
    setActivePinia(createPinia())
    api.getMessages.mockReset().mockResolvedValue({ code: 0, data: { page_data: [], total: 0 } })
    api.updateMessageReview.mockReset().mockResolvedValue({ code: 0 })
    window.$message = { success: vi.fn(), error: vi.fn(), info: vi.fn() }
  })

  it('头像列把相对路径转成根相对路径', () => {
    const wrapper = mountPage()
    const column = wrapper.vm.columns.find(e => e.key === 'avatar')

    expect(column.render(msg).props.src).toBe('/public/uploaded/a.png')
  })

  it('没有选中数据时不发请求', async () => {
    const wrapper = mountPage()

    await wrapper.vm.handleUpdateReview([], true)

    expect(api.updateMessageReview).not.toHaveBeenCalled()
    expect(window.$message.info).toHaveBeenCalled()
  })

  it('审核通过后刷新列表', async () => {
    const wrapper = mountPage()
    api.getMessages.mockClear()

    await wrapper.vm.handleUpdateReview([5], true)

    expect(api.updateMessageReview).toHaveBeenCalledWith([5], true)
    expect(window.$message.success).toHaveBeenCalledWith('审核成功')
    await vi.waitFor(() => expect(api.getMessages).toHaveBeenCalled())
  })

  // 与评论页 A16 同类: 失败不能提示成功, 也不该刷新列表
  it('审核失败不提示成功也不刷新列表', async () => {
    api.updateMessageReview.mockRejectedValue(new Error('boom'))
    const wrapper = mountPage()
    api.getMessages.mockClear()

    await wrapper.vm.handleUpdateReview([5], true)

    expect(window.$message.success).not.toHaveBeenCalled()
    expect(api.getMessages).not.toHaveBeenCalled()
  })

  it('切换标签页改变查询参数', async () => {
    const wrapper = mountPage()

    wrapper.vm.handleChangeTab('has_review')
    expect(wrapper.vm.extraParams.is_review).toBe(1)

    wrapper.vm.handleChangeTab('not_review')
    expect(wrapper.vm.extraParams.is_review).toBe(0)

    wrapper.vm.handleChangeTab('all')
    expect(wrapper.vm.extraParams.is_review).toBe(null)
  })
})
