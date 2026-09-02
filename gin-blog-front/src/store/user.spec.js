import { createPinia, setActivePinia } from 'pinia'
import { beforeEach, describe, expect, it, vi } from 'vitest'
import { useUserStore } from '@/store/user'

// store 里的 action 会调后端接口, 这里替换成假实现
vi.mock('@/api', () => ({
  default: {
    getUser: vi.fn(),
    logout: vi.fn(),
  },
}))

const api = (await import('@/api')).default

describe('useUserStore', () => {
  beforeEach(() => {
    setActivePinia(createPinia())
    vi.clearAllMocks()
  })

  it('setToken 后 resetLoginState 能恢复初始状态', () => {
    const store = useUserStore()
    expect(store.token).toBeNull()

    store.setToken('fake-token')
    expect(store.token).toBe('fake-token')

    store.resetLoginState()
    expect(store.token).toBeNull()
  })

  it('没有 token 时 getUserInfo 直接返回, 不发请求', async () => {
    const store = useUserStore()
    await store.getUserInfo()
    expect(api.getUser).not.toHaveBeenCalled()
  })

  it('getUserInfo 把后端字段映射成前端 userInfo', async () => {
    api.getUser.mockResolvedValue({
      code: 0,
      data: {
        id: 1,
        nickname: '阵雨',
        avatar: 'public/uploaded/a.png',
        website: 'https://example.com',
        intro: '简介',
        email: 'a@qq.com',
        article_like_set: ['1', '2'],
        comment_like_set: ['3'],
      },
    })

    const store = useUserStore()
    store.setToken('fake-token')
    await store.getUserInfo()

    expect(store.nickname).toBe('阵雨')
    // 相对路径的头像会被拼上后端地址
    expect(store.avatar).toBe('http://test-server/public/uploaded/a.png')
    // 点赞集合从字符串转成数字, 便于前端 includes 判断
    expect(store.articleLikeSet).toEqual([1, 2])
    expect(store.commentLikeSet).toEqual([3])
  })

  it('getUserInfo 遇到业务错误码时 reject', async () => {
    api.getUser.mockResolvedValue({ code: 1203, message: 'token 已过期' })

    const store = useUserStore()
    store.setToken('fake-token')
    await expect(store.getUserInfo()).rejects.toMatchObject({ code: 1203 })
  })

  it('头像为空时回退到默认头像', async () => {
    api.getUser.mockResolvedValue({
      code: 0,
      data: { id: 1, nickname: 'x', avatar: '', article_like_set: [], comment_like_set: [] },
    })

    const store = useUserStore()
    store.setToken('fake-token')
    await store.getUserInfo()
    expect(store.avatar).toContain('bing.com')
  })

  it('articleLike / commentLike 是切换语义, 重复调用会取消', () => {
    const store = useUserStore()

    store.articleLike(10)
    expect(store.articleLikeSet).toEqual([10])
    store.articleLike(10)
    expect(store.articleLikeSet).toEqual([])

    store.commentLike(5)
    expect(store.commentLikeSet).toEqual([5])
    store.commentLike(5)
    expect(store.commentLikeSet).toEqual([])
  })

  it('logout 会调接口并清空状态', async () => {
    api.logout.mockResolvedValue({ code: 0 })

    const store = useUserStore()
    store.setToken('fake-token')
    await store.logout()

    expect(api.logout).toHaveBeenCalledTimes(1)
    expect(store.token).toBeNull()
  })
})
