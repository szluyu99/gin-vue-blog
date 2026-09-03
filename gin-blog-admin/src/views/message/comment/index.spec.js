import { mount } from '@vue/test-utils'
import { createPinia, setActivePinia } from 'pinia'
import { beforeEach, describe, expect, it, vi } from 'vitest'

import api from '@/api'
import CommentPage from './index.vue'

vi.mock('@/api', () => ({
  default: {
    getComments: vi.fn().mockResolvedValue({ code: 0, data: { page_data: [], total: 0 } }),
    deleteComments: vi.fn().mockResolvedValue({ code: 0 }),
    updateCommentReview: vi.fn().mockResolvedValue({ code: 0 }),
  },
}))

function mountPage() {
  return mount(CommentPage, {
    global: { stubs: { CommonPage: { template: '<div><slot name="action" /><slot /></div>' } } },
  })
}

describe('评论管理', () => {
  beforeEach(() => {
    setActivePinia(createPinia())
    api.getComments.mockReset().mockResolvedValue({ code: 0, data: { page_data: [], total: 0 } })
    api.updateCommentReview.mockReset().mockResolvedValue({ code: 0 })
    window.$message = { success: vi.fn(), error: vi.fn(), info: vi.fn(), warning: vi.fn() }
  })

  // 回归 A16: 「评论类型」和「来源」渲染的是同一个 type 字段, 前者 key 还是空串
  it('评论类型只有一列', () => {
    const wrapper = mountPage()
    const titles = wrapper.vm.columns.map(e => e.title)

    expect(titles.filter(t => t === '评论类型' || t === '来源')).toEqual(['来源'])
    expect(wrapper.vm.columns.every(e => e.key !== '')).toBe(true)
  })

  // 回归 A16: commentTypeMap[row.type].tag 没守卫, 后端多一种类型前端整页就崩
  it('未知的评论类型不会抛异常', () => {
    const wrapper = mountPage()
    const column = wrapper.vm.columns.find(e => e.key === 'type')

    expect(() => column.render({ type: 99 })).not.toThrow()
    expect(column.render({ type: 1 }).children.default()).toBe('文章')
  })

  // 回归 A16: 导出按 item[key] 取值, 老的 reply_nick_name 字段后端已经没有了
  it('回复对象列的 key 指向真实字段', () => {
    const wrapper = mountPage()
    const column = wrapper.vm.columns.find(e => e.title === '回复对象')

    expect(column.key).toBe('reply_user')
  })

  it('审核失败时不提示成功也不刷新列表', async () => {
    api.updateCommentReview.mockRejectedValue(new Error('boom'))
    const wrapper = mountPage()
    api.getComments.mockClear()

    await wrapper.vm.handleUpdateReview([1], true)

    expect(window.$message.success).not.toHaveBeenCalled()
    expect(api.getComments).not.toHaveBeenCalled()
  })
})
