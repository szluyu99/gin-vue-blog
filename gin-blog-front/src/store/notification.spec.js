import { createPinia, setActivePinia } from 'pinia'
import { beforeEach, describe, expect, it, vi } from 'vitest'
import { useNotificationStore } from '@/store/notification'

vi.mock('@/api', () => ({
  default: {
    getNotifications: vi.fn(),
    getUnreadNotificationCount: vi.fn(),
    readNotifications: vi.fn(),
  },
}))

const api = (await import('@/api')).default

const list = [
  { id: 1, is_read: false, content: 'a', article_id: 3 },
  { id: 2, is_read: false, content: 'b', article_id: 4 },
]

describe('useNotificationStore', () => {
  beforeEach(() => {
    setActivePinia(createPinia())
    vi.clearAllMocks()
    api.getUnreadNotificationCount.mockResolvedValue({ code: 0, data: 2 })
    api.getNotifications.mockResolvedValue({ code: 0, data: { page_data: list.map(e => ({ ...e })), total: 2 } })
    api.readNotifications.mockResolvedValue({ code: 0, data: 1 })
  })

  it('拉取未读数与列表', async () => {
    const store = useNotificationStore()

    await store.fetchUnreadCount()
    expect(store.unreadCount).toBe(2)

    await store.fetchList()
    expect(store.list).toHaveLength(2)
    expect(store.total).toBe(2)
    expect(store.loading).toBe(false)
  })

  // 后端直接返回数字, 接口异常或返回非数字时不能让红点变成 NaN / undefined
  it('未读数拿到非数字时退回 0', async () => {
    const store = useNotificationStore()

    api.getUnreadNotificationCount.mockResolvedValue({ code: 0, data: null })
    await store.fetchUnreadCount()
    expect(store.unreadCount).toBe(0)

    api.getUnreadNotificationCount.mockRejectedValue(new Error('boom'))
    store.unreadCount = 5
    await store.fetchUnreadCount()
    expect(store.unreadCount).toBe(5)
  })

  it('列表接口失败时不抛出, loading 复位', async () => {
    const store = useNotificationStore()
    api.getNotifications.mockRejectedValue(new Error('boom'))

    await store.fetchList()

    expect(store.list).toEqual([])
    expect(store.loading).toBe(false)
  })

  it('标记单条已读: 本地置位并重新取未读数', async () => {
    const store = useNotificationStore()
    await store.fetchList()

    api.getUnreadNotificationCount.mockResolvedValue({ code: 0, data: 1 })
    await store.read([1])

    expect(api.readNotifications).toHaveBeenCalledWith([1])
    expect(store.list[0].is_read).toBe(true)
    expect(store.list[1].is_read).toBe(false)
    // 未读数从后端重新取, 而不是本地减: 别的标签页可能也读过
    expect(store.unreadCount).toBe(1)
  })

  it('ids 为空表示全部已读', async () => {
    const store = useNotificationStore()
    await store.fetchList()

    api.getUnreadNotificationCount.mockResolvedValue({ code: 0, data: 0 })
    await store.read()

    expect(api.readNotifications).toHaveBeenCalledWith([])
    expect(store.list.every(e => e.is_read)).toBe(true)
    expect(store.unreadCount).toBe(0)
  })

  // 退出登录要连列表一起清掉, 否则下一个人登录会先看到上一个人的通知
  it('reset 清空未读数与列表', async () => {
    const store = useNotificationStore()
    await store.fetchUnreadCount()
    await store.fetchList()

    store.reset()

    expect(store.unreadCount).toBe(0)
    expect(store.list).toEqual([])
    expect(store.total).toBe(0)
  })
})
