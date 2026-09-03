import { mount } from '@vue/test-utils'
import { createPinia, setActivePinia } from 'pinia'
import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest'

import api from '@/api'
import Comment from './Comment.vue'
import CommentField from './CommentField.vue'
import Paging from './Paging.vue'

vi.mock('@/api', () => ({
  default: {
    getComments: vi.fn(),
    getCommentReplies: vi.fn(),
    saveComment: vi.fn().mockResolvedValue({ code: 0 }),
    saveLikeComment: vi.fn().mockResolvedValue({ code: 0 }),
  },
}))

vi.mock('vue-router', async importOriginal => ({
  ...await importOriginal(),
  useRoute: () => ({ params: { id: '7' } }),
}))

function makeComment(id, replyCount = 0) {
  return {
    id,
    user_id: id * 10,
    content: `评论 ${id}`,
    created_at: '2026-09-01T00:00:00Z',
    like_count: 0,
    reply_count: replyCount,
    reply_list: [],
    user: { info: { nickname: `用户${id}`, avatar: '' } },
  }
}

// 列表加载有 0.8s 的延时, 用假定时器推过去
async function mountComment(comments) {
  api.getComments.mockResolvedValue({
    code: 0,
    data: { page_data: comments, total: comments.length },
  })
  const wrapper = mount(Comment, { props: { type: 1 } })
  await vi.advanceTimersByTimeAsync(900)
  await wrapper.vm.$nextTick()
  return wrapper
}

describe('前台评论列表', () => {
  beforeEach(() => {
    vi.useFakeTimers()
    setActivePinia(createPinia())
    api.getCommentReplies.mockReset().mockResolvedValue({ code: 0, data: [] })
    window.$message = { success: vi.fn(), error: vi.fn(), info: vi.fn() }
  })

  afterEach(() => {
    vi.useRealTimers()
  })

  // 回归: 原来按 v-for 下标取模板 ref 数组, Vue 不保证顺序与源数组一致,
  // 点第二条评论的「回复」可能打开别人的回复框
  it('点击回复只打开这条评论的回复框, 并带上正确的父评论', async () => {
    const wrapper = await mountComment([makeComment(1), makeComment(2)])

    // 顶部的评论框始终存在, 回复框还没有
    expect(wrapper.findAllComponents(CommentField)).toHaveLength(1)

    // 第二条评论的「回复」按钮
    await wrapper.findAll('button.color-\\#ef2f11')[1].trigger('click')

    const fields = wrapper.findAllComponents(CommentField)
    expect(fields).toHaveLength(2)
    const reply = fields[1]
    expect(reply.props('parentId')).toBe(2)
    expect(reply.props('replyUserId')).toBe(20)
    expect(reply.props('nickname')).toBe('用户2')
  })

  it('切换到另一条评论的回复框时, 上一个会关掉', async () => {
    const wrapper = await mountComment([makeComment(1), makeComment(2)])
    const buttons = wrapper.findAll('button.color-\\#ef2f11')

    await buttons[0].trigger('click')
    expect(wrapper.findAllComponents(CommentField)[1].props('parentId')).toBe(1)

    await buttons[1].trigger('click')
    const fields = wrapper.findAllComponents(CommentField)
    expect(fields).toHaveLength(2) // 仍然只有一个回复框
    expect(fields[1].props('parentId')).toBe(2)
  })

  // 回归: 原来用 checkRefs[idx].style.display 直接改 DOM
  it('点击查看后隐藏「点击查看」并按回复数决定是否显示分页', async () => {
    const c1 = makeComment(1, 4) // 4 条回复: 展开后不需要分页
    const c2 = makeComment(2, 8) // 8 条回复: 展开后有分页
    const wrapper = await mountComment([c1, c2])

    expect(wrapper.findAllComponents(Paging)).toHaveLength(0)

    // 展开第二条评论的回复
    await wrapper.findAll('button.color-\\#00a1d6')[1].trigger('click')
    await wrapper.vm.$nextTick()

    expect(api.getCommentReplies).toHaveBeenCalledWith(2, { page_num: 1, page_size: 5 })
    const pagings = wrapper.findAllComponents(Paging)
    expect(pagings).toHaveLength(1)
    expect(pagings[0].props('pageTotal')).toBe(2)
    // 第二条的「点击查看」被隐藏, 第一条不受影响
    const checks = wrapper.findAll('button.color-\\#00a1d6')
    expect(checks[0].isVisible()).toBe(true)
    expect(checks[1].isVisible()).toBe(false)
  })

  it('回复分页翻页带上对应评论的 id 与页码', async () => {
    const wrapper = await mountComment([makeComment(1, 8)])
    await wrapper.find('button.color-\\#00a1d6').trigger('click')
    await wrapper.vm.$nextTick()

    api.getCommentReplies.mockClear()
    wrapper.findComponent(Paging).vm.$emit('changeCurrent', 2)
    await vi.advanceTimersByTimeAsync(0)
    await wrapper.vm.$nextTick()

    expect(api.getCommentReplies).toHaveBeenCalledWith(1, { page_num: 2, page_size: 5 })
    expect(wrapper.findComponent(Paging).props('current')).toBe(2)
  })

  // 回归: 原来是 pageRefs.value[idx].current, 按下标取错组件就会刷到别人的回复
  it('提交回复后按当前页重新加载这条评论的回复', async () => {
    const wrapper = await mountComment([makeComment(1, 8)])
    await wrapper.find('button.color-\\#00a1d6').trigger('click')
    await wrapper.vm.$nextTick()
    wrapper.findComponent(Paging).vm.$emit('changeCurrent', 2)
    await vi.advanceTimersByTimeAsync(0)
    await wrapper.vm.$nextTick()

    await wrapper.findAll('button.color-\\#ef2f11')[0].trigger('click')
    api.getCommentReplies.mockClear()

    const reply = wrapper.findAllComponents(CommentField)[1]
    reply.vm.$emit('afterSubmit')
    await wrapper.vm.$nextTick()

    expect(api.getCommentReplies).toHaveBeenCalledWith(1, { page_size: 5, page_num: 2 })
  })
})
