import { mount } from '@vue/test-utils'
import { createPinia, setActivePinia } from 'pinia'
import { beforeEach, describe, expect, it, vi } from 'vitest'

import api from '@/api'
import { useUserStore } from '@/store'
import HomePage from './index.vue'

vi.mock('@/api', () => ({
  default: {
    getHomeInfo: vi.fn(),
    getUserInfo: vi.fn().mockResolvedValue({ code: 0, data: {} }),
  },
}))

const homeInfo = { view_count: 10, user_count: 2, article_count: 4, message_count: 6 }

function mountPage() {
  return mount(HomePage, {
    global: { stubs: { AppPage: { template: '<div><slot /></div>' } } },
  })
}

describe('后台首页', () => {
  beforeEach(() => {
    setActivePinia(createPinia())
    api.getHomeInfo.mockReset().mockResolvedValue({ code: 0, data: { ...homeInfo } })
    window.$message = { success: vi.fn(), error: vi.fn(), info: vi.fn() }
  })

  it('挂载后填充统计数据', async () => {
    const wrapper = mountPage()
    await vi.waitFor(() => expect(wrapper.vm.homeInfo.view_count).toBe(10))

    expect(wrapper.vm.homeInfo.article_count).toBe(4)
  })

  it('接口失败或返回空数据时统计不会变成 null', async () => {
    api.getHomeInfo.mockRejectedValue(new Error('boom'))
    const wrapper = mountPage()
    await vi.waitFor(() => expect(api.getHomeInfo).toHaveBeenCalled())
    await wrapper.vm.$nextTick()
    expect(wrapper.vm.homeInfo).toMatchObject({ view_count: 0 })

    api.getHomeInfo.mockResolvedValue({ code: 0, data: null })
    const wrapper2 = mountPage()
    await wrapper2.vm.$nextTick()
    expect(wrapper2.vm.homeInfo).toMatchObject({ view_count: 0 })
  })

  // 以前是 const { nickname, avatar } = useUserStore(), 解构后失去响应性,
  // getUserInfo() 回来得比首屏渲染晚, 问候语里的昵称一直是空的
  it('用户信息后到也能显示在问候语里', async () => {
    const wrapper = mountPage()
    await vi.waitFor(() => expect(wrapper.vm.homeInfo.view_count).toBe(10))
    expect(wrapper.text()).not.toContain('管理员')

    const store = useUserStore()
    store.userInfo = { ...store.userInfo, nickname: '管理员' }
    await wrapper.vm.$nextTick()

    expect(wrapper.text()).toContain('管理员')
  })

  // 原来这句话来自 https://v1.hitokoto.cn, 现在从内置文案里随机取, 不再打外部请求
  it('问候语从内置文案里取, 不发外部请求', async () => {
    const fetchSpy = vi.fn()
    globalThis.fetch = fetchSpy

    const wrapper = mountPage()
    await vi.waitFor(() => expect(wrapper.vm.sentence).not.toBe(''))

    expect(wrapper.vm.SENTENCES).toContain(wrapper.vm.sentence)
    expect(fetchSpy).not.toHaveBeenCalled()
  })
})
