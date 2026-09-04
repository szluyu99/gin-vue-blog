import { mount } from '@vue/test-utils'
import { createPinia, setActivePinia } from 'pinia'
import { beforeEach, describe, expect, it, vi } from 'vitest'

import api from '@/api'
import { useAppStore, useUserStore } from '@/store'
import MessagePage from './index.vue'

vi.mock('@/api', () => ({
  default: {
    getMessages: vi.fn(),
    saveMessage: vi.fn(),
  },
}))

// 弹幕组件依赖真实 DOM 尺寸, 用桩替掉并暴露 push
const push = vi.fn()
vi.mock('vue3-danmaku', () => ({
  default: {
    name: 'VueDanmaku',
    props: ['danmus'],
    setup(_props, { expose }) {
      expose({ push, hide: vi.fn(), show: vi.fn() })
      return () => null
    },
  },
}))

const messages = [
  { id: 1, nickname: '路人甲', content: '你好', avatar: 'public/uploaded/a.png' },
]

function mountPage() {
  return mount(MessagePage, {
    global: { stubs: { AppFooter: true } },
  })
}

describe('留言板', () => {
  beforeEach(() => {
    setActivePinia(createPinia())
    push.mockClear()
    api.getMessages.mockReset().mockResolvedValue({ code: 0, data: messages })
    api.saveMessage.mockReset().mockResolvedValue({ code: 0 })
    window.$message = { success: vi.fn(), error: vi.fn(), info: vi.fn() }
  })

  it('作者的默认弹幕加上接口返回的留言', async () => {
    const wrapper = mountPage()
    await vi.waitFor(() => expect(wrapper.vm.danmus).toHaveLength(2))

    expect(wrapper.vm.danmus[0].nickname).toBe('阵、雨')
    expect(wrapper.vm.danmus[1].nickname).toBe('路人甲')
  })

  it('接口返回空数据时只保留默认弹幕', async () => {
    api.getMessages.mockResolvedValue({ code: 0, data: null })
    const wrapper = mountPage()
    await vi.waitFor(() => expect(api.getMessages).toHaveBeenCalled())
    await wrapper.vm.$nextTick()

    expect(wrapper.vm.danmus).toHaveLength(1)
  })

  it('内容为空时不发请求', async () => {
    const wrapper = mountPage()
    await vi.waitFor(() => expect(wrapper.vm.danmus).toHaveLength(2))

    wrapper.vm.content = '   '
    await wrapper.vm.send()

    expect(api.saveMessage).not.toHaveBeenCalled()
    expect(window.$message.info).toHaveBeenCalled()
  })

  it('发送成功后推入弹幕并清空输入框', async () => {
    const wrapper = mountPage()
    await vi.waitFor(() => expect(wrapper.vm.danmus).toHaveLength(2))

    wrapper.vm.content = '新留言'
    await wrapper.vm.send()

    expect(api.saveMessage).toHaveBeenCalledWith(expect.objectContaining({ content: '新留言' }))
    expect(push).toHaveBeenCalled()
    expect(wrapper.vm.content).toBe('')
  })

  // 回归: 以前是裸 await, 失败会留下未捕获的 rejection 且输入框被清空
  it('发送失败时不清空输入框也不推入弹幕', async () => {
    api.saveMessage.mockRejectedValue(new Error('boom'))
    const wrapper = mountPage()
    await vi.waitFor(() => expect(wrapper.vm.danmus).toHaveLength(2))

    wrapper.vm.content = '发不出去的留言'
    await wrapper.vm.send()

    expect(push).not.toHaveBeenCalled()
    expect(wrapper.vm.content).toBe('发不出去的留言')
  })

  it('没有配置留言页封面时用兜底图', async () => {
    const wrapper = mountPage()
    await vi.waitFor(() => expect(wrapper.vm.danmus).toHaveLength(2))

    expect(wrapper.vm.coverStyle).toContain('/images/page/message.jpeg')

    const appStore = useAppStore()
    appStore.page_list = [{ id: 1, label: 'message', cover: '/images/page/custom.jpg' }]
    await wrapper.vm.$nextTick()
    expect(wrapper.vm.coverStyle).toContain('/images/page/custom.jpg')
  })

  it('未登录时用默认头像发送', async () => {
    const wrapper = mountPage()
    await vi.waitFor(() => expect(wrapper.vm.danmus).toHaveLength(2))
    const userStore = useUserStore()

    wrapper.vm.content = '游客留言'
    await wrapper.vm.send()

    expect(api.saveMessage).toHaveBeenCalledWith(expect.objectContaining({ avatar: userStore.avatar }))
  })
})
