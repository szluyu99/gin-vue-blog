import { mount } from '@vue/test-utils'
import { createPinia, setActivePinia } from 'pinia'
import { beforeEach, describe, expect, it, vi } from 'vitest'

import api from '@/api'
import { useNotificationStore, useUserStore } from '@/store'
import NotificationPage from './index.vue'

vi.mock('@/api', () => ({
  default: {
    getNotifications: vi.fn(),
    getUnreadNotificationCount: vi.fn().mockResolvedValue({ code: 0, data: 0 }),
    readNotifications: vi.fn().mockResolvedValue({ code: 0, data: 1 }),
  },
}))

const push = vi.fn()
vi.mock('vue-router', async importOriginal => ({
  ...await importOriginal(),
  useRouter: () => ({ push }),
}))

const items = [
  { id: 1, type: 1, is_read: false, from_nickname: '小明', content: '回复内容', article_id: 3, article_title: '某篇文章', created_at: '2026-09-05T00:00:00Z' },
  { id: 2, type: 2, is_read: true, from_nickname: '小红', content: '评论内容', article_id: 4, article_title: '另一篇', created_at: '2026-09-04T00:00:00Z' },
]

function mountPage() {
  return mount(NotificationPage, {
    global: {
      stubs: {
        BannerPage: { template: '<div><slot /></div>' },
        RouterLink: { props: ['to'], template: '<a :href="to"><slot /></a>' },
      },
    },
  })
}

function login() {
  const userStore = useUserStore()
  userStore.userInfo = { ...userStore.userInfo, id: 1 }
}

describe('站内通知页', () => {
  beforeEach(() => {
    setActivePinia(createPinia())
    push.mockReset()
    api.getNotifications.mockReset().mockResolvedValue({
      code: 0,
      data: { page_data: items.map(e => ({ ...e })), total: 2 },
    })
    api.readNotifications.mockReset().mockResolvedValue({ code: 0, data: 1 })
  })

  // 未登录不该发请求, 也不该显示空列表, 而是引导登录
  it('未登录时提示去登录, 不请求接口', async () => {
    const wrapper = mountPage()
    await vi.waitFor(() => expect(wrapper.vm.loading).toBe(false))

    expect(api.getNotifications).not.toHaveBeenCalled()
    expect(wrapper.text()).toContain('登录后才能查看站内通知')
  })

  it('登录后拉取并渲染列表', async () => {
    login()
    const wrapper = mountPage()
    await vi.waitFor(() => expect(wrapper.vm.list).toHaveLength(2))

    expect(api.getNotifications).toHaveBeenCalledWith({ page_num: 1, page_size: 10 })
    expect(wrapper.text()).toContain('共 2 条通知')
    expect(wrapper.text()).toContain('小明')
    expect(wrapper.text()).toContain('回复了你')
    expect(wrapper.text()).toContain('评论了你的文章')
    expect(wrapper.text()).toContain('来自《某篇文章》')
  })

  it('点未读通知先标已读再跳文章', async () => {
    login()
    const wrapper = mountPage()
    await vi.waitFor(() => expect(wrapper.vm.list).toHaveLength(2))

    await wrapper.vm.open(wrapper.vm.list[0])

    expect(api.readNotifications).toHaveBeenCalledWith([1])
    expect(wrapper.vm.list[0].is_read).toBe(true)
    expect(push).toHaveBeenCalledWith('/article/3')
  })

  it('点已读通知不再重复请求标记', async () => {
    login()
    const wrapper = mountPage()
    await vi.waitFor(() => expect(wrapper.vm.list).toHaveLength(2))

    await wrapper.vm.open(wrapper.vm.list[1])

    expect(api.readNotifications).not.toHaveBeenCalled()
    expect(push).toHaveBeenCalledWith('/article/4')
  })

  it('全部已读把本页每条都置位', async () => {
    login()
    const notificationStore = useNotificationStore()
    notificationStore.unreadCount = 1
    const wrapper = mountPage()
    await vi.waitFor(() => expect(wrapper.vm.list).toHaveLength(2))

    await wrapper.vm.readAll()

    expect(api.readNotifications).toHaveBeenCalledWith([])
    expect(wrapper.vm.list.every(e => e.is_read)).toBe(true)
  })

  it('翻页带上新页码重新请求', async () => {
    login()
    api.getNotifications.mockResolvedValue({
      code: 0,
      data: { page_data: items.map(e => ({ ...e })), total: 25 },
    })
    const wrapper = mountPage()
    await vi.waitFor(() => expect(wrapper.vm.total).toBe(25))
    window.scrollTo = vi.fn()

    expect(wrapper.vm.pageCount).toBe(3)
    wrapper.vm.current = 2
    await vi.waitFor(() => expect(api.getNotifications).toHaveBeenCalledWith({ page_num: 2, page_size: 10 }))
    expect(window.scrollTo).toHaveBeenCalled()
  })

  it('没有通知时给出空状态', async () => {
    login()
    api.getNotifications.mockResolvedValue({ code: 0, data: { page_data: [], total: 0 } })
    const wrapper = mountPage()
    await vi.waitFor(() => expect(wrapper.vm.loading).toBe(false))
    await wrapper.vm.$nextTick()

    expect(wrapper.text()).toContain('还没有通知')
  })

  it('接口失败时 loading 复位, 列表保持空数组', async () => {
    login()
    api.getNotifications.mockRejectedValue(new Error('boom'))
    const wrapper = mountPage()
    await vi.waitFor(() => expect(wrapper.vm.loading).toBe(false))

    expect(wrapper.vm.list).toEqual([])
  })
})
